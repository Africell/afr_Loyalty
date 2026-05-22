# Endpoint Testing Guide — Post-Migration (daoc → mongox + redisx)

Use this guide after running the migration to verify the service starts correctly and all migrated endpoints behave as expected. Fill in your project-specific values wherever `<placeholders>` appear.

---

## Project Configuration Reference

Fill this in before starting:

| Variable | Value |
|---|---|
| Management port | `<ManagementPort>` e.g. `9802` |
| App port | `<AppPort>` e.g. `9801` |
| Service port | `<ServicePort>` e.g. `9803` |
| Module | `<Module>` e.g. `AppCMS` |
| Version | `<Version>` e.g. `V1` |
| Base URL | `http://localhost:<ManagementPort>/<Module>/<Version>/` |
| MongoDB port | `27017` (default) |
| Redis port | `6379` (default) |
| AUC / Auth service port | `<AUCPort>` e.g. `9001` |
| MongoDB database | `<DatabaseName>` |
| AutoIncrement collection | `<AutoIncrementCollection>` e.g. `Col_AutoIncrement` |

---

## 1. Pre-Flight: Verify Required Services Are Running

Before starting the application, confirm all dependencies are reachable.

### MongoDB

```powershell
# Windows PowerShell
$tcp = New-Object System.Net.Sockets.TcpClient
try { $tcp.Connect("localhost", 27017); "MongoDB: UP" } catch { "MongoDB: DOWN" } finally { $tcp.Close() }
```

```bash
# Linux / macOS
nc -z localhost 27017 && echo "MongoDB: UP" || echo "MongoDB: DOWN"
```

If down: start MongoDB before proceeding. CRUD operations, index creation, and auto-increment will all fail without it.

### Redis

```powershell
$tcp = New-Object System.Net.Sockets.TcpClient
try { $tcp.Connect("localhost", 6379); "Redis: UP" } catch { "Redis: DOWN" } finally { $tcp.Close() }
```

```bash
nc -z localhost 6379 && echo "Redis: UP" || echo "Redis: DOWN"
```

### Auth Center / AUC (if used)

```powershell
$tcp = New-Object System.Net.Sockets.TcpClient
try { $tcp.Connect("localhost", <AUCPort>); "AUC: UP" } catch { "AUC: DOWN" } finally { $tcp.Close() }
```

The management AUC must be running for route registration and authenticated management endpoint calls. If it is down, the service still starts but access entries are not created and authenticated routes fail.

### Redis ACL

If your Redis instance was started with ACL users configured, a successful TCP/PING connection does **not** mean key access works. The `default` user may have `nokeys`.

Verify:

```bash
redis-cli ACL WHOAMI     # which user the app connects as
redis-cli ACL LIST       # all users and their key permissions
```

If the `default` user has no `~*` key pattern, fix it before starting the service:

```bash
# DEV only — grant full key access to default user
redis-cli ACL SETUSER default ~* +@all nopass
```

For non-DEV environments, configure `Username` and `Password` in `redisx.Config` to match a user that has permissions for the key prefix your project uses.

---

## 2. Build and Start the Service

```powershell
# Windows
go build -o service.exe .
Start-Process -FilePath ".\service.exe" `
    -RedirectStandardOutput "stdout.log" `
    -RedirectStandardError  "stderr.log" `
    -PassThru | Select-Object -ExpandProperty Id
```

```bash
# Linux / macOS
go build -o service .
./service > stdout.log 2> stderr.log &
echo $!
```

Wait 2–3 seconds then check startup logs:

```powershell
Get-Content stderr.log | Select-Object -Last 20
```

```bash
tail -20 stderr.log
```

---

## 3. Verify Startup Log

A healthy startup looks like:

```
Connected successfuly to Redis           ← redisx.New + ping succeeded
Authenticated Successfuly                ← AUC S2S auth OK (if AUC is used)
DB index maintenance process started...
DB index maintenance process finished
Add <type> routers to the web service
  Read existing routes: # <N>
  Created Access Entry: <Route>|<Method>  ← new entries registered in AUC
<type> listen and serve on port: <Port>  ← all expected ports up
```

### Expected warnings in DEV (not errors)

```
Failed to authenticate ... <port> ... refused  ← optional AUC instance not running
Error Reading Existing Access Entries from AUC  ← same reason, non-fatal
```

### Startup failures to investigate

| Log message | Likely cause |
|---|---|
| `redis ping failed` | Redis not running, wrong address, or ACL blocks connection |
| `mongox.Connect` / `NewDB` error | MongoDB not running or URI wrong |
| `mongox.NewRepository` error | DB name or collection name mismatch in config |
| `RedisDataLoader` error | Redis ACL blocks key access; see section 1 |
| `log.Fatalln` / service exits immediately | AUC required at startup is not running |
| nil pointer panic on first request | A repository global was not initialized |

---

## 4. Bypass Auth for Testing (Development Only)

Management endpoints are protected by auth middleware (e.g., `ValidateAccess_AUC`, `ValidateJWEToken`). To test without a user token, temporarily remove those middlewares from the route definition file.

In the management routes function, change each protected handler from:

```go
Use(UC.HTTP_<Entity>, UC.ValidateAccess_AUC, UC.ValidateJWEToken)
```

to:

```go
Use(UC.HTTP_<Entity>)
```

Rebuild and restart after editing.

> **Restore before commit.** After testing, add the auth middleware back to every route you changed. Use `git diff <routes-file>` to confirm the file is clean before committing.

Public endpoints (mobile cache, file downloads, health/test) have no auth by design — do not modify those.

---

## 5. No-Auth Smoke Test

The project should have at least one unauthenticated endpoint for health/test purposes. Identify it from the service routes and call it:

```powershell
# Replace with the actual no-auth endpoint path
Invoke-WebRequest -Uri "http://localhost:<ServicePort>/<Module>/<Version>/<TestRoute>/" `
    -Method GET -UseBasicParsing
# Expect: 200 {"Status":"successful",...}
```

Also verify the Prometheus metrics endpoint (if present):

```powershell
Invoke-WebRequest -Uri "http://localhost:<ManagementPort>/metrics" `
    -Method GET -UseBasicParsing
```

In the metrics output, check connection status for MongoDB and Redis:

```
PortStatus{DestinationHost="MongoDB 01"} 1   ← 1 = up, 2 = down
```

---

## 6. CRUD Test Pattern

Repeat this pattern for each entity in the project. Replace `<Entity>`, `<EntityRoute>`, `<Port>`, and field names with your project's actual values.

### POST — Create

```powershell
$body = '{"Key":"<TestKey>", <other required fields>}'
$r = Invoke-WebRequest `
    -Uri "http://localhost:<Port>/<Module>/<Version>/<EntityRoute>/" `
    -Method POST -ContentType "application/json" -Body $body -UseBasicParsing
$r.StatusCode   # expect 200
$r.Content      # expect {"Data":{"Key":"<TestKey>","Id":<N>,...},"Status":"successful"}
```

**Verifies:** `redisx.NextAutoIncrementID` (Redis INCR + MongoDB upsert), `Mdb_<Entity>.InsertOne`, `redisx.SetJSON`.

If `Id` is `0` or you get a Redis permission error:
- Check Redis ACL (see section 1)
- Confirm `Mdb_AutoIncrement` repository is initialized in the startup sequence

### GET — Read

```powershell
$r = Invoke-WebRequest `
    -Uri "http://localhost:<Port>/<Module>/<Version>/<EntityRoute>/" `
    -Method GET -UseBasicParsing
$r.StatusCode   # expect 200
$r.Content      # expect {"Data":[...],"Status":"successful"}
```

**Verifies:** `mongox.FindPage` on `Mdb_<Entity>`.

### PUT — Update

```powershell
$body = '{"Key":"<TestKey>","Id":<N>, <updated fields>}'
$r = Invoke-WebRequest `
    -Uri "http://localhost:<Port>/<Module>/<Version>/<EntityRoute>/" `
    -Method PUT -ContentType "application/json" -Body $body -UseBasicParsing
$r.StatusCode   # expect 200
```

Follow with a GET to confirm the change persisted in MongoDB and Redis.

**Verifies:** `Mdb_<Entity>.UpdateOne` with `$set` + upsert, `redisx.SetJSON` cache update.

### DELETE

```powershell
$r = Invoke-WebRequest `
    -Uri "http://localhost:<Port>/<Module>/<Version>/<EntityRoute>/<TestKey>/" `
    -Method DELETE -UseBasicParsing
$r.StatusCode   # expect 200

# Confirm empty after delete
$r = Invoke-WebRequest `
    -Uri "http://localhost:<Port>/<Module>/<Version>/<EntityRoute>/" `
    -Method GET -UseBasicParsing
$r.Content      # expect {"Data":null,...}
```

**Verifies:** `Mdb_<Entity>.DeleteOne`, `redisx.DelJSON` cache eviction.

---

## 7. Verify MongoDB Directly

After CRUD tests, connect to MongoDB and confirm:

```javascript
use <DatabaseName>

// Check entity records were written
db.<EntityCollection>.find().pretty()

// Most important: auto-increment counters exist and are incrementing
db.<AutoIncrementCollection>.find().pretty()
// Each entity that uses NextAutoIncrementID should have a document here
// with the key name matching what the code passes to NextAutoIncrementID

// If access entries are stored in MongoDB, verify they exist
db.<AccessEntryCollection>.find().count()
```

---

## 8. Verify Redis Directly

```bash
# List all keys (include KeyPrefix if configured, e.g. "pdc:dev:")
redis-cli KEYS "*"

# Check a specific cached entity
redis-cli GET "<EntityName>:<TestKey>"
# With prefix: redis-cli GET "<KeyPrefix><EntityName>:<TestKey>"

# Check auto-increment counters
redis-cli KEYS "*seq*"
redis-cli GET "<KeyPrefix>seq:<EntityName>-Id"

# Verify RedisDataLoader loaded startup data
redis-cli KEYS "<EntityName>:*"
```

If MongoDB collections were empty at startup, `RedisDataLoader` will load nothing — that is expected. Run a POST first, then check that the key appears in Redis.

---

## 9. Test File Upload / GridFS (if applicable)

If the project stores files in GridFS, test upload and download.

**Upload (multipart form):**

```powershell
$boundary = "----TestBoundary"
$filePath  = "C:\path\to\testfile.png"
$fileBytes = [System.IO.File]::ReadAllBytes($filePath)
$fileName  = [System.IO.Path]::GetFileName($filePath)

$head    = [System.Text.Encoding]::UTF8.GetBytes(
    "--$boundary`r`nContent-Disposition: form-data; name=`"<fieldName>`"; filename=`"$fileName`"`r`nContent-Type: image/png`r`n`r`n"
)
$tail    = [System.Text.Encoding]::UTF8.GetBytes("`r`n--$boundary--`r`n")
$payload = $head + $fileBytes + $tail

$r = Invoke-WebRequest `
    -Uri "http://localhost:<Port>/<Module>/<Version>/<UploadRoute>/" `
    -Method POST `
    -ContentType "multipart/form-data; boundary=$boundary" `
    -Body $payload -UseBasicParsing
$r.Content   # expect {"Data":{"FileId":"<hex-oid>","ContentType":"image/png"},...}
```

**Download:**

```powershell
$fileId = "<hex-oid from upload response>"
Invoke-WebRequest `
    -Uri "http://localhost:<Port>/<Module>/<Version>/<DownloadRoute>/$fileId" `
    -Method GET -OutFile "downloaded_test.png" -UseBasicParsing
# Verify downloaded_test.png matches the uploaded file
```

---

## 10. Common Failures and Fixes

| Symptom | Likely cause | Fix |
|---|---|---|
| `404 page not found` | Wrong URL — module/version mismatch | Check `Configuration.Module` and `Configuration.Version`; use `/<Module>/<Version>/<Route>/` |
| `400 <field> is required` | Missing required field in request body | Include the field; check the request struct in `AppCMS_Structures.go` or equivalent |
| `400 id is not matching` | PUT body missing `Id` | Include `"Id":<n>` matching the existing record |
| `NOPERM No permissions to access a key` | Redis ACL blocks key pattern | Grant key access to the Redis user; see section 1 |
| `redis next id failed` | Redis ACL or nil Redis client | Check ACL and startup log for Redis init error |
| `mongo read AutoIncrement failed` | `Mdb_AutoIncrement` not initialized | Check `InitializeMongoxRepositories` for missing repository |
| `500` or nil pointer panic | A repository global is nil | One repository was not initialized — check startup log |
| GET returns empty after POST | Redis cache not written | Check `redisx.SetJSON` call in the POST handler; verify key in Redis |
| Counter resets to 1 on restart | Auto-increment key name mismatch | Align key names with what is in the `AutoIncrement` collection |

---

## 11. Cleanup After Testing

1. Delete test records from MongoDB:

```javascript
use <DatabaseName>
db.<EntityCollection>.deleteMany({"Key": /^<TestKey>/})
// repeat for each entity tested
```

2. Redis cache entries are cleared on next `RedisDataLoader` restart, or delete manually:

```bash
redis-cli DEL "<EntityName>:<TestKey>"
```

3. If auth was bypassed (section 4), restore it:

```bash
git diff <routes-definition-file>
# must show no diff
```

4. Rebuild to confirm the restored file compiles cleanly:

```powershell
go build -o service.exe .
```
