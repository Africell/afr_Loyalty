# Generic Migration Guide: `daoc` → `mongox` + `redisx`

Use this guide as the main instruction reference for Claude Code when migrating a Go project from the old `daoc` implementation to the newer `mongox` + `redisx` pattern.

The goal is to remove the old DAO/cache implementation and replace it with:

- `mongox` for MongoDB client initialization, repositories, CRUD, indexes, GridFS, and direct collection access.
- `redisx` for Redis client initialization, JSON cache operations, Mongo → Redis loading, and auto-increment IDs.
- `sync.Map` from the Go standard library only when a process-local in-memory cache is still required.

This guide is intentionally generic. Apply the patterns to the current project’s actual package names, struct names, collection names, startup files, and runtime behavior. Do not blindly copy example names such as `CMS_Application`, `MobileCache`, or `AccessEntry` unless the current project actually has those concepts.

**Important adaptation rule:** example receiver names like `UserControl`, field names like `MongoClient`, Redis globals like `RedisClient`, and direct client field access like `.Mongo` are placeholders. Substitute the current project’s actual container type, naming style, and confirmed `mongox` / `redisx` API. Do not rename an existing project container just to match this guide.

---

## 1. Migration Objectives

Claude Code should migrate the project by doing the following:

1. Remove all direct `daoc` imports and dependencies.
2. Remove old `daoc.DAO`, `daoc.MongoDB`, `daoc.GridFS`, `daoc.CacheRegistry`, and `daoc.Cache_Synch` usage.
3. Replace MongoDB access with `mongox.Client` and `*mongox.Repository`.
4. Replace Redis/cache access with `*redisx.Client` and `redisx` JSON helpers.
5. Replace process-local cache maps with `sync.Map` only when the old cache was truly local and in-memory.
6. Replace old GridFS/file helpers with `mongox.GridFS` methods.
7. Replace old startup/init sequence with explicit Mongo/Redis client initialization, repository initialization, optional Redis loading, and background index creation.
8. Run formatting, module cleanup, and tests after migration.

Keep existing project structure and function names where possible. Do not introduce unnecessary wrappers, helper context functions, or new abstractions unless the current project clearly requires them.

---

## 2. First Step: Scan the Project

Before editing code, scan the whole project and identify all files containing old `daoc` patterns.

Search terms:

```text
daoc
DAO_
Map_
Cache_Synch
CacheRegistry
MongoDB
GridFS
SysAdminInit
InitializeDAO
InitializeCache
CheckAndCreateIndex
HTTPWriteFile
ReadFileInBase64
HTTPUploadFileWithMimeType
DAOFindParams
DAOFindCriteria
DAOPaginate
AutoIncrement
sync.Map
context.Background()
```

Create a short migration checklist from the findings:

- Files importing `daoc`.
- Files importing `daoc/SysAdmin`.
- Files defining old DAO globals.
- Files defining old cache globals.
- Files initializing old DAO/cache.
- Files using old DAO CRUD methods.
- Files using old cache methods.
- Files using old file upload/download/base64 helpers.
- Files using old auto-increment logic.
- Structs that need `RedisKey() string` methods.
- Any old cache usage whose replacement is unclear.

Then migrate file by file.

If a pattern is unusual and this guide does not cover it, report it clearly and choose the closest safe migration pattern. Do not silently skip it.

If a project already uses `sync.Map`, inspect the existing usage before adding new globals so names and responsibilities do not collide.

### Pre-flight API compatibility check

Before changing many files, inspect the actual local `mongox` and `redisx` packages used by this repo. Confirm:

- `mongox.Client` field/method used to access the underlying Mongo client. Examples may show `.Mongo`, but the real API may differ.
- `mongox.NewDB`, `mongox.NewRepository`, `mongox.NewGridFS`, `mongox.FindPage`, and `mongox.InsertOne` signatures.
- `redisx.Config` fields and `DefaultTTL` semantics.
- `redisx.SetJSON`, `redisx.GetJSON`, `redisx.DelJSON`, `redisx.LoadMongoToRedis`, and `redisx.NextAutoIncrementID` signatures.
- Whether the project already has wrapper helpers or lifecycle/shutdown conventions that should be preserved.

If the actual package API differs from the examples, adapt the examples to the package API. Do not force the examples if they do not compile.

---

## 3. `go.mod` Changes

Remove old direct `daoc` dependency entries if present — **unless another local module that this project replaces (`replace X => ../X/`) still imports `daoc` internally.** In that case keep `daoc` as an indirect entry; see the transitive dependency note below.

```go
replace daoc => ../daoc/
daoc v0.0.0-00010101000000-000000000000
```

Remove direct MongoDB driver v1 usage from project code:

```go
go.mongodb.org/mongo-driver v1.x.x
```

Add or keep the new dependencies. Version numbers below are examples; prefer the versions required by this repo, the local modules, or the latest compatible versions:

```go
replace mongox => ../mongox/
replace redisx => ../redisx/

require (
    mongox v0.0.0-00010101000000-000000000000
    redisx v0.0.0-00010101000000-000000000000
    go.mongodb.org/mongo-driver/v2 vX.Y.Z // example only
    github.com/redis/go-redis/v9 vX.Y.Z // example only
)
```

If the project handles MIME detection for uploads, also include:

```go
github.com/gabriel-vasile/mimetype vX.Y.Z // example only
```

Then run `go mod tidy`. It is also acceptable to run `go mod tidy` between edit cycles if dependency errors cascade while replacing imports.

Important:

- Do not manually fight `go mod tidy` if it keeps indirect dependencies required by unrelated packages.
- The migration goal is no direct `daoc` usage and no direct v1 MongoDB driver imports in project code.
- All project MongoDB imports should use `go.mongodb.org/mongo-driver/v2/...` after migration.

### Transitive `daoc` dependency via local module replacements

If this project uses `replace` directives for local modules (e.g., `replace some_shared_lib => ../some_shared_lib/`), and those local modules internally import `daoc`, Go will require a `replace daoc => ../daoc/` directive in **this project's** `go.mod` as well — even though no project code directly imports `daoc`. Go does not propagate `replace` directives from dependencies.

In this case, keep `daoc` as an **indirect** dependency only:

```go
replace daoc => ../daoc/

require (
    daoc v0.0.0-00010101000000-000000000000 // indirect
)
```

`go mod tidy` will add this automatically if the transitive dependency is detected. Do not remove it manually — the build will fail with "missing dot in first path element" or similar if it is absent.

---

## 4. Import Migration

Remove imports like:

```go
"daoc"
sysadmin "daoc/SysAdmin"
"reflect" // if only used by old daoc patterns
```

Add imports as needed:

```go
import (
    "context"
    "fmt"
    "log"
    "sync"
    "time"

    "go.mongodb.org/mongo-driver/v2/bson"
    "go.mongodb.org/mongo-driver/v2/mongo/options"
    "mongox"
    "redisx"
)
```

For file upload/download/base64 logic, add only when needed:

```go
import (
    "bytes"
    "encoding/base64"
    "io"
    "net/http"

    "github.com/gabriel-vasile/mimetype"
)
```

Avoid unused imports. Run `go fmt` and fix compile errors.

---

## 5. Configuration Pattern

Add Mongo and Redis configuration to the project configuration struct. Use the pattern that best matches the current project.

### Redis `Mode` is required — no default, omitting it crashes the service

`redisx.New` validates the `Mode` field at startup. A zero-value `redisx.Config{}` — or any config where `Mode` is not explicitly set — causes `redisx.New` to return an error such as `"unknown redis mode:"`, which will crash the service before it binds to any port.

**Every environment config function that calls `redisx.New` must set `Mode` explicitly.** There is no safe default.

**Responsibility split:**

- The **dev** config function (e.g., `setDefaultConfiguration_Dev`) in `Configuration.go` **must** set all Redis fields in code, because developers run the service locally without external config injection.
- **UAT, staging, and production** Redis config is the responsibility of **DevOps** — values are injected via environment variables, secrets managers, or config files at deploy time. Do not hardcode non-dev credentials in source code.

### Dev config — set all Redis fields in code

In the dev config function, always set every Redis field explicitly:

```go
Configuration.Redis.Mode     = redisx.ModeSingle
Configuration.Redis.Addr     = "127.0.0.1:6379"
Configuration.Redis.Username = "dev_user"
Configuration.Redis.Password = "dev_password"
Configuration.Redis.DB       = 0
Configuration.Redis.KeyPrefix = "appname:dev:"
Configuration.Redis.DefaultTTL = -1 // See TTL note below.
```

Adapt `Addr`, `Username`, `Password`, and `KeyPrefix` to the actual dev Redis instance for this project. Use `redisx.ModeSentinel` or `redisx.ModeCluster` if the dev instance requires it.

### Option A — Direct package config structs

Use this if the project can store `mongox.Config` and `redisx.Config` directly:

```go
type ConfigType struct {
    Mongo mongox.Config
    Redis redisx.Config
}
```

Example:

```go
Configuration.Mongo = mongox.Config{
    URI:                    "mongodb://localhost:27017",
    AppName:                Configuration.Module,
    ConnectTimeout:         10 * time.Second,
    ServerSelectionTimeout: 10 * time.Second,
    SocketTimeout:          30 * time.Second,
    MinPoolSize:            1,
    MaxPoolSize:            100,
    RetryReads:             true,
    RetryWrites:            true,
}

Configuration.Redis = redisx.Config{
    Mode:       redisx.ModeSingle,
    Addr:       "127.0.0.1:6379",
    Username:   "dev_user",
    Password:   "dev_password",
    DB:         0,
    KeyPrefix:  "appname:dev:",
    DefaultTTL: -1, // See TTL note below.
}
```

### Option B — Rename old custom Mongo struct fields, keep individual fields

Use this when the old project had a custom MongoDB struct (host, port, username, password, replica set) and DevOps manages the per-environment credentials directly in the config functions.

**Do NOT convert the individual-field environment configs to URI format. Only modify the dev config function.**

In `ConfigType`, rename the old structs from `MongoDB`/`LoyaltyMongoDB` to `Mongo`/`LoyaltyMongo`, keeping the old struct definitions commented out for reference:

```go
import "redisx"

type ConfigType struct {
    Mongo struct {
        ReplicaSet string
        UserName   string
        Password   string
        HostIP_1   string
        HostPort_1 string
        HostIP_2   string
        HostPort_2 string
        HostIP_3   string
        HostPort_3 string
        HostIP_4   string
        HostPort_4 string
    }
    LoyaltyMongo struct {
        ReplicaSet string
        UserName   string
        Password   string
        HostIP_1   string
        HostPort_1 string
        HostIP_2   string
        HostPort_2 string
        HostIP_3   string
        HostPort_3 string
        HostIP_4   string
        HostPort_4 string
    }
    // MongoDB struct { ... old fields ... }
    // LoyaltyMongoDB struct { ... old fields ... }
    Redis redisx.Config
}
```

Update `buildMongoURI` and `buildLoyaltyMongoURI` helpers to reference `cfg.Mongo` and `cfg.LoyaltyMongo` instead of `cfg.MongoDB` and `cfg.LoyaltyMongoDB`. These helpers are still needed in `UserControl.go`:

```go
func buildMongoURI(cfg ConfigType) string {
    c := cfg.Mongo
    // ... same logic as before, using c.ReplicaSet, c.UserName, etc.
}

func buildLoyaltyMongoURI(cfg ConfigType) string {
    c := cfg.LoyaltyMongo
    // ... same logic
}
```

In `UserControl.go` (or equivalent container init), keep using the helpers and add a `Ping` after each connect to verify reachability:

```go
mongoClient, err := mongox.Connect(ctx, mongox.Config{
    URI:     buildMongoURI(Configuration),
    AppName: "appname",
})
if err != nil {
    log.Fatal("mongox.Connect (MongoDB):", err)
}
if err := mongoClient.Ping(ctx); err != nil {
    log.Fatal("Mongo not reachable:", err)
}

loyaltyMongoClient, err := mongox.Connect(ctx, mongox.Config{
    URI:     buildLoyaltyMongoURI(Configuration),
    AppName: "appname",
})
if err != nil {
    log.Fatal("mongox.Connect (LoyaltyMongoDB):", err)
}
if err := loyaltyMongoClient.Ping(ctx); err != nil {
    log.Fatal("LoyaltyMongo not reachable:", err)
}
```

> **Note on Ping:** `mongox.Connect` already validates connectivity internally via Ping before returning. The explicit `Ping` call above is a required pattern for clarity and explicit startup fail-fast behaviour. It is not redundant for operational reasons.

#### Only modify the Dev config function

**Do not touch any UAT, staging, or Live config functions.** Those are managed by DevOps and should keep their existing individual-field assignments. Only `setDefaultConfiguration_Dev` (or equivalent) needs to be updated to use the new field names (`Mongo.` instead of `MongoDB.`) and to set all Redis fields for local development:

```go
// setDefaultConfiguration_Dev — only this function is modified during migration
Configuration.Mongo.ReplicaSet = ""
Configuration.Mongo.UserName   = ""
Configuration.Mongo.Password   = ""
Configuration.Mongo.HostIP_1   = "localhost"
Configuration.Mongo.HostPort_1 = "27017"
Configuration.Mongo.HostIP_2   = ""
Configuration.Mongo.HostPort_2 = ""
Configuration.Mongo.HostIP_3   = ""
Configuration.Mongo.HostPort_3 = ""
Configuration.Mongo.HostIP_4   = ""
Configuration.Mongo.HostPort_4 = ""
Configuration.LoyaltyMongo = Configuration.Mongo

Configuration.Redis.Mode       = redisx.ModeSingle  // always use the constant, not the string "single"
Configuration.Redis.Addr       = "localhost:6379"
Configuration.Redis.Username   = ""
Configuration.Redis.Password   = ""
Configuration.Redis.DB         = 0
Configuration.Redis.KeyPrefix  = "appname:dev:"     // replace appname with the service name
Configuration.Redis.DefaultTTL = -1                 // -1 = no expiry; 0 would be overridden to 5 min by redisx defaults
```

All other environment config functions need only a field-name rename (`MongoDB.` → `Mongo.`, `LoyaltyMongoDB.` → `LoyaltyMongo.`) — their credential values stay untouched.

### Redis TTL note

Preserve the old cache expiry behavior.

- If old `daoc` cache data did not expire, configure Redis keys with no expiry. In some `redisx` versions, `DefaultTTL = -1` means no default expiry.
- If old cache data had a TTL, preserve the same TTL.
- If the target `redisx` implementation treats `DefaultTTL = 0` as an internal default, do not use `0` when you need non-expiring keys.
- Confirm the actual `redisx.Config` semantics in the current repo before applying a TTL globally.
- Note: `redisx.SetJSON` uses the config `DefaultTTL`; `redisx.SetJSONWithTTL` takes an explicit duration. Use the one that matches the required expiry behavior.

Do not assume one TTL policy fits all projects.

---

## 6. Project Container Migration (`UserControl` / Equivalent)

`UserControl` is only an example name. Use the existing project container type, such as `UserControl`, `AppContext`, `ServiceLayer`, `Server`, package-level globals, or whatever the repo already uses. Do not rename the project’s architecture to match this guide.

Remove old fields:

```go
MongoDB  *daoc.MongoDB
CacheDir *daoc.CacheRegistry
```

Replace with:

```go
type ProjectContainer struct {
    MongoClient *mongox.Client
    Redis       *redisx.Client

    // Keep existing project-specific clients here.
}
```

Example using direct `mongox.Config` pass-through (single Mongo database):

```go
func NewProjectContainer() *ProjectContainer {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    mongoClient, err := mongox.Connect(ctx, Configuration.Mongo)
    if err != nil {
        log.Fatal("mongox.Connect:", err)
    }

    redisClient, err := redisx.New(Configuration.Redis)
    if err != nil {
        log.Fatal("redisx.New:", err)
    }

    return &ProjectContainer{
        MongoClient: mongoClient,
        Redis:       redisClient,
    }
}
```

Example with two Mongo clients (primary + secondary database on the same cluster):

```go
func NewProjectContainer() *ProjectContainer {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    mongoClient, err := mongox.Connect(ctx, Configuration.Mongo)
    if err != nil {
        log.Fatal("mongox.Connect (primary):", err)
    }

    secondaryMongoClient, err := mongox.Connect(ctx, Configuration.LoyaltyMongo)
    if err != nil {
        log.Fatal("mongox.Connect (secondary):", err)
    }

    redisClient, err := redisx.New(Configuration.Redis)
    if err != nil {
        log.Fatal("redisx.New:", err)
    }

    return &ProjectContainer{
        MongoClient:          mongoClient,
        SecondaryMongoClient: secondaryMongoClient,
        Redis:                redisClient,
    }
}
```

Pass `Configuration.Mongo` (or `Configuration.LoyaltyMongo`) directly — do not build URIs by hand or use a `buildMongoURI` helper. The URI is already set in each environment config function.

Keep existing project-specific fields in the project container if they are still needed. If the project uses package-level globals instead of a container struct, add the clients using that existing style.

---

## 7. Old Cache Replacement Decision Table

Do not replace every old `daoc.Cache_Synch` with the same thing. First understand how it was used.

| Old usage | Recommended replacement |
|---|---|
| Startup-only deduplication map | Local `map[string]T` inside the startup/setup function |
| Runtime shared cache across requests | Redis via `redisx` |
| Process-local derived cache rebuilt from Redis/Mongo | `sync.Map` |
| Persistent source of truth | MongoDB via `mongox.Repository` |
| Auto-increment counter | `redisx.NextAutoIncrementID` |
| File storage helper | `mongox.GridFS` |

Rules:

- Use Redis when the cache must be shared across app instances or survive process restarts.
- Use `sync.Map` only for local derived data that can be rebuilt.
- Use a local map for temporary startup deduplication.
- Do not create custom cache wrapper structs unless the project already has a justified abstraction.

---

## 8. Global Repository, Redis, GridFS, and Local Cache Naming

Replace old globals like:

```go
var DAO_Entity daoc.DAO
var Map_Entity daoc.Cache_Synch
var LocalDerivedCache daoc.Cache_Synch
```

With project-specific names following these patterns:

```go
var Mdb_Entity *mongox.Repository
var Gfs_Entity *mongox.GridFS // only for file buckets
var RedisClient *redisx.Client // or Rc / the name consistent with this project's style
var LocalDerivedCache sync.Map // only if local in-memory cache is needed
```

General naming rules:

| Old pattern | New pattern | Type |
|---|---|---|
| `DAO_<Entity>` | `Mdb_<Entity>` | `*mongox.Repository` |
| `Map_<Entity>` runtime cache | Redis via `RedisClient` | `*redisx.Client` |
| local in-memory cache | `<Name> sync.Map` | `sync.Map` |
| DAO file helpers | `Gfs_<Entity>` or `Gfs_<Bucket>` | `*mongox.GridFS` |

Example only:

```go
var Mdb_User *mongox.Repository
var Mdb_Order *mongox.Repository
var Mdb_AutoIncrement *mongox.Repository
var Mdb_EventLog *mongox.Repository

var Gfs_UserFiles *mongox.GridFS

var RedisClient *redisx.Client // or Rc / the name consistent with this project's style
var AppLocalCache sync.Map
```

Adapt names to the current project’s real collections and buckets.

---

## 9. Repository and GridFS Initialization

Delete old functions such as:

```go
InitializeDAO()
InitializeCache()
```

Replace DAO initialization with repository initialization:

```go
func (app *ProjectContainer) InitializeMongoxRepositories() error {
    // Replace app.MongoClient.Mongo with the confirmed way this mongox version exposes the underlying Mongo client.
    db, err := mongox.NewDB(app.MongoClient.Mongo, Configuration.DB_Name, 10*time.Second)
    if err != nil {
        return fmt.Errorf("mongox.NewDB: %w", err)
    }

    Mdb_Entity, err = mongox.NewRepository(db, "Col_Entity")
    if err != nil {
        return fmt.Errorf("Col_Entity: %w", err)
    }

    // Repeat for every Mongo collection.

    RedisClient = app.Redis
    return nil
}
```

For GridFS buckets:

```go
Gfs_EntityFiles, err = mongox.NewGridFS(db, mongox.GridFSConfig{BucketName: "EntityFiles"})
if err != nil {
    return fmt.Errorf("GridFS EntityFiles: %w", err)
}
```

Only create GridFS globals for collections/buckets that actually store files.

---

## 10. Redis Data Loader

Replace old cache initialization with Redis loading only for collections that should be cached in Redis.

```go
func (app *ProjectContainer) RedisDataLoader() error {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
    defer cancel()

    if _, _, err := redisx.LoadMongoToRedis[Entity](
        ctx,
        RedisClient,
        Mdb_Entity.Coll,
        redisx.MongoLoadOptions{
            BatchSize:       2000,
            // TTL is example-only. Preserve old cache expiry behavior and confirm redisx semantics.
            // Use no-expiry, zero, or a duration according to this repo's redisx API.
            TTL:             0,
            FlushBeforeLoad: true,
            FlushPattern:    "Entity:*",
            UseUnlink:       true,
        },
    ); err != nil {
        return fmt.Errorf("load Entity: %w", err)
    }

    // Repeat only for cached collections.
    // Rebuild local derived caches after Redis is loaded, if needed.

    return nil
}
```

Rules:

- Load only collections that were previously cached or need fast Redis lookup.
- Do not load high-volume write-only logs into Redis.
- Match `FlushPattern` to the entity’s `RedisKey()` prefix.
- Preserve old TTL behavior. Use `TTL: 0`, `TTL: -1`, or a specific duration according to the actual `redisx` API semantics in this repo.

---

## 11. Redis Key Methods

Every entity stored in Redis must implement:

```go
RedisKey() string
```

General format:

```go
func (e EntityName) RedisKey() string {
    return "EntityName:" + e.Key
}
```

Examples:

```go
func (e User) RedisKey() string {
    return "User:" + e.Key
}

func (e Product) RedisKey() string {
    return "Product:" + e.Key
}
```

Rules:

- Use `entry.RedisKey()` when writing an entity to Redis.
- Use `EntityType{Key: keyVar}.RedisKey()` when reading or deleting by key.
- Do not hardcode keys like `"Entity:" + Key` at call sites when an entity method exists.
- Use collection-specific patterns only for scans, for example `"Entity:*"`.
- If an entity does not have a `Key` field, choose the existing unique identifier and make the key method reflect that clearly.

---

## 12. Context Timeout Rules

Do not create helper functions like:

```go
crudCtx()
gridCtx()
cacheCtx()
```

Use `context.WithTimeout(...)` inline at each call site. For startup jobs and background workers, `context.Background()` is fine. Inside HTTP handlers, prefer `r.Context()` as the parent context so cancellations and graceful shutdown propagate correctly.

Recommended timeouts:

| Operation | Timeout |
|---|---:|
| Single CRUD read/write/delete | `10 * time.Second` |
| GridFS upload/download | `30 * time.Second` |
| Redis scan / bulk load / startup cache | `5 * time.Minute` |
| Index creation | `context.Background()` |

Examples:

```go
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()
```

```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()
```

```go
// In an HTTP handler, prefer the request context as the parent:
ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
defer cancel()
```

---

## 13. Error Handling Convention

Prefer preserving the project’s existing error style, but use these rules when migrating code:

- When returning an underlying error to a caller, prefer wrapping with `fmt.Errorf("operation: %w", err)` so callers can inspect it.
- Use `errors.New(...)` only for new business errors that do not wrap another error.
- Avoid replacing meaningful existing error messages unless needed for compile correctness.
- Log at boundaries such as startup or background workers; avoid logging and returning the same error multiple times from deep helper functions unless the project already does this.

Example:

```go
if err != nil {
    return fmt.Errorf("load Entity into redis: %w", err)
}
```

---

## 14. CRUD Migration Patterns

### 14.1 Create

Old pattern may look like:

```go
DAO_Entity.PutOne(...)
Map_Entity.Put(...)
```

New pattern:

```go
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

if _, err := Mdb_Entity.InsertOne(ctx, newEntry); err != nil {
    if mongox.IsDuplicateKey(err) {
        return 0, errors.New("key already exists")
    }
    return 0, err
}

if err := redisx.SetJSON(ctx, RedisClient, newEntry.RedisKey(), newEntry); err != nil {
    return 0, err
}
```

Only write to Redis if this entity is part of the Redis cache.

### 14.2 Read single from Redis

Old pattern:

```go
Map_Entity.CheckThenGet(Key, &entry)
```

New pattern:

```go
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

entry, err := redisx.GetJSON[Entity](ctx, RedisClient, Entity{Key: Key}.RedisKey())
if redisx.IsNil(err) {
    return entry, errors.New("key does not exist")
}
if err != nil {
    return entry, err
}

return entry, nil
```

### 14.3 Read all from Redis

Old pattern:

```go
Map_Entity.ConvertToArray(...)
Map_Entity.ConvertToArraySortedByStringField(...)
```

New pattern:

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
defer cancel()

entries, err := redisx.GetAllJSONByPattern[Entity](
    ctx,
    RedisClient,
    redisx.ScanJSONOptions{
        Pattern:      "Entity:*",
        ScanCount:    500,
        PipelineSize: 250,
        Limit:        10000,
    },
)
if err != nil {
    return nil, err
}

return entries, nil
```

Sort in Go after retrieval if the old cache method returned sorted results.

### 14.4 Read paginated from MongoDB

Old pattern:

```go
DAO_Entity.FindPaginate(...)
```

New pattern:

```go
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

result, err := mongox.FindPage[Entity](
    ctx,
    Mdb_Entity.Coll,
    bson.M{},
    mongox.PageRequest{
        Page:        Page,
        PageSize:    Limit,
        MaxPageSize: 50000,
    },
)
if err != nil {
    return result, err
}

return result, nil
```

### 14.5 Update / Upsert

Old pattern:

```go
DAO_Entity.PutOne(...)
Map_Entity.Put(...)
```

New pattern:

```go
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

_, err := Mdb_Entity.UpdateOne(
    ctx,
    bson.M{"Key": entry.Key},
    bson.M{"$set": entry},
    options.UpdateOne().SetUpsert(true),
)
if err != nil {
    return entry.Id, errors.New("update failed: " + err.Error())
}

if err := redisx.SetJSON(ctx, RedisClient, entry.RedisKey(), entry); err != nil {
    return entry.Id, err
}

return entry.Id, nil
```

Only update Redis if this entity is part of the Redis cache.

### 14.6 Delete

Old pattern:

```go
DAO_Entity.Delete(...)
Map_Entity.Delete(...)
```

New pattern:

```go
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

if _, err := Mdb_Entity.DeleteOne(ctx, bson.M{"Key": Key}); err != nil {
    return err
}

if _, err := redisx.DelJSON(ctx, RedisClient, Entity{Key: Key}.RedisKey()); err != nil {
    return err
}

return nil
```

Important:

```go
redisx.DelJSON(...)
```

returns two values. Always use:

```go
_, err = redisx.DelJSON(...)
```

not:

```go
err = redisx.DelJSON(...)
```

### 14.7 Key rename during edit

If the old key changes, delete the old Mongo/Redis record first:

```go
if request.Key != entry.Key {
    if _, err := Mdb_Entity.DeleteOne(ctx, bson.M{"Key": request.Key}); err != nil {
        return entry.Id, err
    }

    if _, err := redisx.DelJSON(ctx, RedisClient, Entity{Key: request.Key}.RedisKey()); err != nil {
        return entry.Id, err
    }
}
```

Then upsert the new entry and write it to Redis if cached.

---

## 15. Auto-Increment ID Migration

Old pattern:

```go
Map_AutoIncrement.GetNextAI(...)
daoc.AutoIncrement{}
```

New pattern:

```go
func (app *ProjectContainer) GetNewId(ctx context.Context, key string) (int64, error) {
    return redisx.NextAutoIncrementID(
        ctx,
        RedisClient,
        Mdb_AutoIncrement.Coll,
        key,
        redisx.NextIDOptions{
            RedisBase:          "AutoIncrement",
            EmitReconcileEvent: true,
            ReconcileStream:    "AutoIncrement:reconcile",
            MongoRetries:       3,
            RetryBackoff:       500 * time.Millisecond,
        },
    )
}
```

If the current project uses a simpler option set, keep only the options required by that project.

Usage:

```go
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

id, err := app.GetNewId(ctx, "Entity")
if err != nil {
    return 0, err
}
```

Redis is the primary counter. MongoDB is the fallback/reconciliation store if configured that way by `redisx`.

### Key naming for auto-increment

The `key` string passed to `NextAutoIncrementID` must match the `Key` field stored in the MongoDB `AutoIncrement` collection. Using a different name re-starts the counter from 1 and loses the existing sequence.

> **Human action required after migration:** Before running the service for the first time, check the existing `AutoIncrement` collection and confirm that the key names used in all `GetNewId` / `NextAutoIncrementID` calls match what is already there. Adjust the key strings in code if they differ.
>
> ```javascript
> db.<AutoIncrementCollection>.find().pretty()
> ```
>
> If the collection is empty (fresh environment), choose a consistent naming convention for this project and document it in the project-specific appendix (section 32).

---

## 16. Startup-Only Dedup Cache Migration

Some old projects used `daoc.Cache_Synch` as a temporary deduplication map during route setup, permission setup, or other startup registration.

Do not automatically replace this with a global `sync.Map`.

Generic replacement:

```go
existing := make(map[string]ExistingEntity)
for _, item := range sourceItems {
    existing[item.Key] = item
}

if _, ok := existing[newItem.Key]; !ok {
    // create/register/upsert new item
}
```

Rules:

- Use a local map if the data is only needed during one startup/setup function.
- Use Redis if the data must be shared at runtime across instances.
- Use MongoDB via `mongox.Repository` if it is persistent source-of-truth data.
- Pass the local map explicitly to helper functions if needed.

Example only:

```go
func (Uc *UserControl) AddEntityIfMissing(existing map[string]Entity, entry Entity) error {
    if _, ok := existing[entry.Key]; ok {
        return nil
    }

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    _, err := Mdb_Entity.UpdateOne(
        ctx,
        bson.M{"Key": entry.Key},
        bson.M{"$set": entry},
        options.UpdateOne().SetUpsert(true),
    )
    return err
}
```

### 16.1 External Service + MongoDB Dual-Write (Access Entries Pattern)

When the project registers entries in an **external service** (such as an Auth Center / AUC) at startup and also needs to keep a local MongoDB copy as source of truth, use a local map for deduplication and persist newly created entries to MongoDB via upsert.

Reading step (build local dedup map from external service):

```go
accessEntries, err := Uc.ExternalService.ReadEntries("")
if err != nil {
    log.Println("Error reading access entries from external service")
}

var existingEntries = make(map[string]AuthCenter.AccessEntry)
for _, acc := range accessEntries.Data {
    existingEntries[acc.Key] = acc
}
```

Then pass `existingEntries` to the router setup loop and to creation helpers.

Creation helper (create in external service, then upsert to local MongoDB and write to Redis):

```go
func (Uc *UserControl) AddToExternalServiceAndMongo(existing map[string]AuthCenter.AccessEntry, accessEntry AuthCenter.AccessEntry) {
    accessEntry.Key = accessEntry.AccessKey + "|" + accessEntry.AccessMethod
    if _, ok := existing[accessEntry.Key]; ok {
        return
    }

    accessEntry.Id = 0
    accessEntry.Status = ""
    accessEntry.AddUser = "Auto Add"
    accessEntry.AddDate = time.Now().UTC()
    accessEntry.LastModifyUser = "Auto Add"
    accessEntry.LastModifyDate = time.Now().UTC()

    _, err := Uc.ExternalService.CreateEntry(accessEntry)
    if err != nil {
        log.Println("Error creating entry in external service:", err)
        return
    }

    ctx := context.Background()

    Mdb_AccessEntry.UpdateOne(
        ctx,
        bson.M{"Key": accessEntry.Key},
        bson.M{"$set": accessEntry},
        options.UpdateOne().SetUpsert(true),
    )
    redisx.SetJSON(ctx, Uc.Redis, accessEntry.RedisKey(), accessEntry)
    log.Println("Created Access Entry:", accessEntry.Key)
}
```

Rules:

- Do **not** use a global `sync.Map` as a dedup cache for access entries. Use a local `map[string]T` populated once per router setup call.
- Always read existing entries from the external service first, build the local map, then call the creation helper.
- After creating in the external service, upsert to MongoDB with `options.UpdateOne().SetUpsert(true)` (idempotent — safe to re-run on restart), then write to Redis with `redisx.SetJSON`.
- Use `context.Background()` for the MongoDB upsert and Redis write during startup; no per-operation timeout is needed since these are fire-and-forget persistence writes after the external service call already succeeded.
- The old `MapAccessEntry.Clear()` + `MapAccessEntry.Put()` + `MapAccessEntry.Check()` pattern is fully replaced by a plain local `map[string]T`.

---

## 17. GridFS and File Handling Migration

Remove old file helper patterns:

```go
DAO_*.HTTPUploadFileWithMimeType(...)
DAO_*.HTTPWriteFile(...)
DAO_*.ReadFileInBase64(...)
```

Replace with `mongox.GridFS`.

### 17.1 Upload file from HTTP request

```go
func (Uc *UserControl) HTTP_EntityUpload(r *http.Request, fieldName string) (string, string) {
    file, handler, err := r.FormFile(fieldName)
    if err != nil {
        return "", ""
    }
    defer file.Close()

    data, err := io.ReadAll(file)
    if err != nil {
        return "", ""
    }

    mtype := mimetype.Detect(data)

    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    oid, err := Gfs_Entity.UploadFromReader(
        ctx,
        handler.Filename,
        bytes.NewReader(data),
        mongox.UploadOptions{ContentType: mtype.String()},
    )
    if err != nil {
        return "", ""
    }

    return oid.Hex(), mtype.String()
}
```

Use the correct bucket variable for the file type.

### 17.2 Serve/download file to HTTP response

```go
oid, err := bson.ObjectIDFromHex(fileID)
if err != nil {
    http.Error(w, "invalid file id", http.StatusBadRequest)
    return
}

ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

if _, err := Gfs_Entity.DownloadToWriter(ctx, oid, w); err != nil {
    http.Error(w, "file not found", http.StatusNotFound)
    return
}
```

### 17.3 Read file as base64

```go
func readFileBase64(gfs *mongox.GridFS, fileID string) (string, error) {
    if fileID == "" {
        return "", nil
    }

    oid, err := bson.ObjectIDFromHex(fileID)
    if err != nil {
        return "", err
    }

    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    data, err := gfs.DownloadBytes(ctx, oid)
    if err != nil {
        return "", err
    }

    return base64.StdEncoding.EncodeToString(data), nil
}
```

Usage example:

```go
imageBase64, err := readFileBase64(Gfs_EntityFiles, entity.FileID)
```

---

## 18. Local In-Memory Cache Migration with `sync.Map`

Use `sync.Map` only for process-local derived cache that can be rebuilt from Redis or MongoDB.

```go
var LocalCache sync.Map
```

Basic operations:

```go
LocalCache.Store(key, value)

val, ok := LocalCache.Load(key)
if ok {
    cached, _ := val.(CacheValue)
}

LocalCache.Delete(key)

LocalCache.Clear() // Go 1.23+

LocalCache.Range(func(k, v any) bool {
    // process k, v
    return true
})
```

If the project uses Go older than 1.23, replace `Clear()` with:

```go
LocalCache.Range(func(k, v any) bool {
    LocalCache.Delete(k)
    return true
})
```

Rules:

- Do not store source-of-truth data only in `sync.Map`.
- Rebuild the local cache after Redis/Mongo startup loading if needed.
- Handle type assertions safely.

---

## 19. Derived App/Mobile Cache Pattern

Some projects build a local derived cache from Redis data. Use this section only if the current project has that behavior.

Generic pattern:

```go
func Local_CacheInit() {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
    defer cancel()

    entries, err := redisx.GetAllJSONByPattern[Entity](
        ctx,
        RedisClient,
        redisx.ScanJSONOptions{Pattern: "Entity:*"},
    )
    if err != nil {
        log.Println("Local_CacheInit error:", err)
        return
    }

    // Go 1.23+: LocalCache.Clear()
    // For older Go versions, Range + Delete each key.
    LocalCache.Clear()

    for _, entry := range entries {
        LocalCache.Store(entry.Key, transformEntry(entry))
    }
}
```

When transforming cache data:

- Use `redisx.GetAllJSONByPattern` for full collection reads.
- Use `redisx.GetJSON` for per-key reads.
- Use a fresh 10-second context for each single Redis read.
- Use a 5-minute context for scans/bulk cache building.

Example per-key Redis read:

```go
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
item, err := redisx.GetJSON[Item](ctx, RedisClient, Item{Key: itemKey}.RedisKey())
cancel()
```

Reading from local cache:

```go
val, ok := LocalCache.Load(cacheKey)
if !ok {
    return CacheValue{}, errors.New("cache entry not found")
}

value, ok := val.(CacheValue)
if !ok {
    return CacheValue{}, errors.New("invalid cache entry type")
}
```

---

## 20. Index Maintenance

Remove old patterns:

```go
DAO_*.CheckAndCreateIndex(...)
CheckAndCreateIndex(...)
go sysadmin.SysAdminInit(...)
```

Create an index maintenance method and run it as a goroutine after repository initialization.

```go
func (app *ProjectContainer) IndexesMaintenanceProcess() {
    keyIdx, keyOpt := mongox.UniqueIndex("Key", true)

    mongox.CreateIndex(context.Background(), Mdb_Entity.Coll, keyIdx, keyOpt)
    mongox.CreateIndex(context.Background(), Mdb_OtherEntity.Coll, keyIdx, keyOpt)

    // Repeat for keyed collections.
}
```

Call it from startup:

```go
go app.IndexesMaintenanceProcess()
```

Rules:

- Use `context.Background()` for index creation unless the project has a stricter standard.
- Create indexes only after repositories are initialized.
- Do not assume every collection has a `Key` field. Use the correct unique field for each collection.

---

## 21. Startup Sequence

Replace old startup flow such as:

```go
UserControl := Module.NewUserControl()
UserControl.InitializeDAO()
UserControl.InitializeCache()
go sysadmin.SysAdminInit(...)
```

With:

```go
UserControl := Module.NewUserControl()

if err := UserControl.InitializeMongoxRepositories(); err != nil {
    log.Fatal("InitializeMongoxRepositories:", err)
}

if err := UserControl.RedisDataLoader(); err != nil {
    log.Fatal("RedisDataLoader:", err)
}

go UserControl.IndexesMaintenanceProcess()
```

If the project does not need Redis preload, skip `RedisDataLoader()` or make it load only required collections.

General order:

```text
GetDefaultConfiguration()
→ NewProjectContainer() / existing constructor
→ InitializeMongoxRepositories()
→ RedisDataLoader() if needed
→ Rebuild local derived caches if needed
→ IndexesMaintenanceProcess() as goroutine
→ Other project-specific startup tasks
→ Start HTTP server / workers
```

---

## 22. Direct Collection Access for Logs or Dynamic Collections

For high-volume write-only logs, do not cache in Redis and do not force repository usage if dynamic collection names are needed.

Pattern:

```go
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

// Replace app.MongoClient.Mongo with the confirmed underlying Mongo client accessor.
col := app.MongoClient.Mongo.Database(Configuration.DB_Name).Collection("Col_Event_Log")
_, err := col.InsertOne(ctx, record)
if err != nil {
    return err
}
```

For date-partitioned logs or dynamic databases:

```go
db := Mdb_BaseCollection.DB.Client.Database(dynamicDbName)
coll := db.Collection(dynamicCollectionName)

// Prefer the project's dominant style. Use raw coll.InsertOne for simple direct writes,
// or mongox.InsertOne if the project already uses that wrapper consistently.
_, err := coll.InsertOne(ctx, record)
if err != nil {
    return err
}
```

Rules:

- Do not load high-volume logs into Redis.
- Preserve existing partitioning rules if the old code wrote to dynamic databases or monthly collections.
- Prefer raw `coll.InsertOne` for simple direct writes unless the project already standardizes on `mongox.InsertOne`. Do not mix styles randomly.

---

## 23. Common Old → New Replacements

| Old `daoc` pattern | New `mongox` / `redisx` pattern |
|---|---|
| `daoc.MongoDB` | `*mongox.Client` |
| `daoc.DAO` | `*mongox.Repository` |
| `daoc.GridFS` / DAO file helpers | `*mongox.GridFS` |
| `daoc.CacheRegistry` | `*redisx.Client` |
| `daoc.Cache_Synch` startup-only dedup cache | local `map[string]T` |
| `daoc.Cache_Synch` runtime shared cache | Redis via `redisx` |
| `daoc.Cache_Synch` process-local derived cache | `sync.Map` |
| `MapAccessEntry.Clear()` + `MapAccessEntry.Put()` + `MapAccessEntry.Check()` | local `map[string]T` populated from external service; see section 16.1 |
| AUC `CreateAccessEntries` without local persistence | `CreateAccessEntries` + `Mdb_AccessEntry.UpdateOne($set, upsert:true)`; see section 16.1 |
| `DAO_*.Initialize(...)` | `mongox.NewRepository(...)` |
| `Map_*.Initialize(...)` | `redisx.LoadMongoToRedis(...)` if Redis preload is needed |
| `Map_*.CheckThenGet(...)` | `redisx.GetJSON[T](...)` |
| `Map_*.Check(...)` | `redisx.GetJSON[T](...)` then check `redisx.IsNil(err)` |
| `Map_*.Put(...)` | `redisx.SetJSON(...)` |
| `Map_*.Delete(...)` | `redisx.DelJSON(...)` |
| `Map_*.ConvertToArray(...)` | `redisx.GetAllJSONByPattern[T](...)` |
| `DAO_*.FindPaginate(...)` | `mongox.FindPage[T](...)` |
| `DAO_*.PutOneLogs(...)` | direct `InsertOne` or `mongox.InsertOne(...)` |
| `DAO_*.HTTPUploadFileWithMimeType(...)` | `Gfs_*.UploadFromReader(...)` |
| `DAO_*.HTTPWriteFile(...)` | `Gfs_*.DownloadToWriter(...)` |
| `DAO_*.ReadFileInBase64(...)` | `readFileBase64(Gfs_*, fileID)` |
| `Map_*.GetNextAI(...)` | `redisx.NextAutoIncrementID(...)` |
| `DAO_*.CheckAndCreateIndex(...)` | `mongox.CreateIndex(...)` |
| `go sysadmin.SysAdminInit(...)` | remove |

---

## 24. What Must Be Removed

Remove all of the following if present:

```go
replace daoc => ../daoc/
daoc v0.0.0-00010101000000-000000000000
```

```go
"daoc"
sysadmin "daoc/SysAdmin"
```

```go
MongoDB *daoc.MongoDB
CacheDir *daoc.CacheRegistry
```

```go
daoc.InitMongoHost(...)
daoc.NewMongoDBClient(...)
daoc.NewCacheRegistry(...)
```

```go
var Map_* daoc.Cache_Synch
var DAO_* daoc.DAO
```

```go
InitializeDAO()
InitializeCache()
CheckAndCreateIndex()
go sysadmin.SysAdminInit(...)
```

```go
daoc.DAOFindParams{}
daoc.DAOFindCriteria{}
daoc.DAOPaginate{}
daoc.AutoIncrement{}
```

Old cache methods:

```go
Initialize
ConvertToArray
ConvertToArraySortedByStringField
CheckThenGet
Check
Put
Delete
GetNextAI
```

Old DAO methods:

```go
Initialize
FindPaginate
CheckAndCreateIndex
PutOneLogs
HTTPWriteFile
ReadFileInBase64
HTTPUploadFileWithMimeType
```

Remove custom cache wrappers like:

```go
mapCache
custom cache structs
custom cache registry wrappers
```

Remove context helper functions like:

```go
crudCtx()
gridCtx()
cacheCtx()
```

Use inline `context.WithTimeout` instead.

---

## 25. Graceful Shutdown and Lifecycle

If the project has server shutdown handling, wire Mongo and Redis clients into that lifecycle. Do not invent a large lifecycle framework, but avoid leaving obvious client cleanup missing.

Generic rules:

- Check the actual `mongox.Client` and `redisx.Client` APIs for `Close`, `Disconnect`, or equivalent methods.
- On SIGTERM/SIGINT or server shutdown, close Redis and disconnect Mongo using the project’s existing pattern.
- If the old `daoc` client had explicit cleanup, migrate that cleanup to the new clients.
- If the project has no shutdown pattern, report this as a remaining operational improvement instead of adding a risky framework.

Example shape only:

```go
func (app *ProjectContainer) Close(ctx context.Context) error {
    // Adapt to actual APIs.
    // _ = app.Redis.Close()
    // _ = app.MongoClient.Disconnect(ctx)
    return nil
}
```

---

## 26. Test Migration Notes

Existing tests may depend on old `daoc` mocks, stubs, or globals. During migration:

- Search tests for `daoc`, `DAO_`, `Map_`, `Cache_Synch`, and old Mongo driver imports.
- Replace old DAO mocks with interfaces or test fixtures that match the project’s existing testing style.
- Prefer testing business logic through repository/client abstractions if the project already has them.
- For integration tests, ensure Redis/Mongo test setup matches the new clients.
- Do not delete failing tests just because they reference `daoc`; migrate or clearly report why they need manual follow-up.

---

## 27. Validation Checklist

After migration, Claude Code must run:

```bash
go fmt ./...
go mod tidy
go vet ./...
go test ./...
```

If `go vet ./...` or `go test ./...` is too broad or the repo has known unrelated test failures, run at least:

```bash
go vet ./path/to/affected/package
go test ./path/to/affected/package
```

Then verify:

- No `daoc` import remains in project code.
- No direct `daoc` dependency remains in `go.mod`.
- No old `DAO_* daoc.DAO` globals remain.
- No old `Map_* daoc.Cache_Synch` globals remain.
- MongoDB driver v1 is not mixed into project code with MongoDB driver v2.
- All Mongo imports in migrated files use `go.mongodb.org/mongo-driver/v2/...`.
- `options.UpdateOne()` and other Mongo options come from the v2 driver.
- `redisx.DelJSON` is handled as a two-return-value function.
- Redis-stored structs have `RedisKey() string` methods.
- Redis key call sites use `entry.RedisKey()` or `Entity{Key: key}.RedisKey()`.
- Redis scans use collection patterns like `"Entity:*"`.
- Startup initializes Mongo/Redis before repositories or Redis loaders are used.
- Index maintenance runs after repositories are initialized.
- File upload/download/base64 logic uses `mongox.GridFS`.
- Local process cache uses `sync.Map` only when appropriate.
- TTL choices preserve the old cache expiry behavior.
- HTTP handlers use `r.Context()` as the parent context where practical.
- Mongo and Redis clients have a shutdown/close path if the project has lifecycle management.

---

## 28. Expected Final Claude Code Report

After editing, provide a concise report with:

```text
Migration summary:
- Updated go.mod dependencies.
- Replaced UserControl daoc fields with mongox/redisx clients.
- Added/updated repository initialization.
- Added/updated Redis data loader if needed.
- Replaced old DAO CRUD calls.
- Replaced old cache calls.
- Replaced GridFS helpers.
- Updated startup sequence.
- Removed sysadmin/daoc initialization.

Validation:
- go fmt ./...: passed/failed
- go mod tidy: passed/failed
- go vet ./...: passed/failed
- go test ./...: passed/failed

Remaining issues:
- List any compile errors, missing structs, unclear old behavior, or manual checks needed.
```

Do not claim the migration is complete if formatting, module cleanup, build, or tests still fail.

After printing this report, prompt the human owner with the following checklist **before they start testing**:

```text
Before testing — action required:

1. Redis config in dev
   Confirm that the dev config function (e.g., setDefaultConfiguration_Dev in Configuration.go)
   sets all redisx.Config fields: Mode, Addr, Username, Password, DB, KeyPrefix, DefaultTTL.
   An empty or partially-set redisx.Config causes redisx.New to fail with "unknown redis mode:"
   and the service will exit before binding to any port.
   UAT and production Redis config is the responsibility of DevOps — do not hardcode those values.

2. Redis ACL
   The app connects to Redis using the credentials in redisx.Config.
   If your Redis instance has ACL users configured, verify that the user has key access:
     redis-cli ACL WHOAMI
     redis-cli ACL LIST
   If the user has no key pattern (~*), grant it:
     redis-cli ACL SETUSER <username> ~* +@all
   A successful PING at startup does NOT mean key access works — key reads/writes may still fail.

3. Auto-increment key names
   The key strings passed to NextAutoIncrementID must match the existing Key values in the
   AutoIncrement MongoDB collection. Mismatched names re-start counters from 1.
   Check the collection and adjust code key strings if needed:
     db.<AutoIncrementCollection>.find().pretty()
```

---

## 29. Guide Improvement Policy

Claude Code should not automatically edit this guide during a project migration.

Instead, if Claude Code discovers a new pattern or correction, it should report proposed additions separately at the end of the migration.

Report format:

```text
Proposed guide updates:
1. Section: [section name]
   Finding: [new daoc pattern or mongox/redisx behavior]
   Proposed addition: [short text or code pattern]
2. Section: [section name]
   Finding: ...
   Proposed addition: ...
```

Examples of valid proposed updates:

- A new `daoc` usage that has no equivalent section.
- A `mongox` or `redisx` API difference in this repo.
- A Go version compatibility issue.
- A startup ordering constraint.
- A TTL behavior confirmed from the actual `redisx` implementation.

The human owner decides whether to update the guide.

---

## 30. New Project Scan Report

When starting a migration on a new project, Claude Code should scan first and print a short report before editing. Then continue with the migration unless there is a blocking ambiguity.

Report format:

```text
New project scan:
- daoc imports found in: [...]
- old DAO globals found in: [...]
- old cache globals found in: [...]
- old GridFS/file helpers found in: [...]
- old auto-increment usage found in: [...]
- startup/init files affected: [...]
- unclear patterns: [...]

Migration approach:
- DAO replacement: mongox repositories
- Runtime cache replacement: redisx
- Local derived cache replacement: sync.Map if needed
- Startup-only dedup replacement: local map if needed
- File storage replacement: mongox GridFS if needed
```

Do not stop after the scan unless the project has missing information that makes safe migration impossible.

---

## 31. Claude Code Prompt to Use with This File

```text
Read the attached migration guide markdown completely before editing.

Migrate this Go project from daoc to mongox + redisx using the guide as the source of truth.

Start by scanning the repo for all daoc imports, DAO globals, cache globals, old cache methods, old DAO methods, old GridFS helpers, auto-increment usage, and startup initialization.

Print a short scan report, then continue migration unless something is truly blocking.

Then migrate file by file:
1. Update go.mod.
2. Update configuration for Mongo and Redis.
3. Update UserControl to store mongox and redisx clients.
4. Replace DAO globals with *mongox.Repository globals.
5. Replace GridFS DAO helpers with *mongox.GridFS globals when file storage is used.
6. Replace cache globals based on behavior:
   - Redis for runtime/shared cache.
   - sync.Map for process-local derived cache.
   - local map for startup-only deduplication.
7. Add RedisKey() methods for all Redis-stored entities.
8. Replace CRUD logic using the patterns in the guide.
9. Replace auto-increment logic with redisx.NextAutoIncrementID.
10. Replace file upload/download/base64 with mongox.GridFS.
11. Replace startup sequence with NewUserControl, InitializeMongoxRepositories, RedisDataLoader if needed, and IndexesMaintenanceProcess.
12. Remove sysadmin and all daoc leftovers.

Do not create unnecessary wrappers or helper context functions.
Use inline context.WithTimeout calls.
Keep the existing project structure and function names where possible.
Preserve old cache TTL behavior instead of assuming one TTL policy for all projects.
Do not automatically edit the migration guide; report proposed guide updates separately.

After editing, run:
- go fmt ./...
- go mod tidy
- go test ./...

Fix migration-related compile errors.
Finally, report changed files, validation results, remaining issues, and any proposed guide updates.
```

---

## 32. Project-Specific Appendix Template

Use this section only when preparing the guide for a specific project. Keep project-specific names here instead of mixing them into the generic sections.

```md
# Project-Specific Appendix: [Project Name]

## Known startup files
- [file]

## Known config files
- [file]

## Old daoc globals
- DAO_X → Mdb_X
- Map_X → Redis or sync.Map depending on usage

## Mongo repositories to initialize
- Mdb_X → Col_X
- Mdb_Y → Col_Y

## GridFS buckets to initialize
- Gfs_X → bucket name X

## Redis cached entities
- EntityX → Redis prefix EntityX:
- EntityY → Redis prefix EntityY:

## Local derived caches
- CacheName → sync.Map

## Startup-only dedup maps
- Old MapX → local map in [function]

## Dev Redis config (set in code — setDefaultConfiguration_Dev or equivalent)
- Mode:       redisx.ModeSingle
- Addr:       [dev Redis host:port]
- Username:   [dev Redis username]
- Password:   [dev Redis password]
- DB:         [dev Redis DB index]
- KeyPrefix:  [appname:dev:]
- DefaultTTL: -1

## UAT / LIVE Redis config
- Managed by DevOps via environment variables or external config injection.
- Do not hardcode UAT or production credentials in source code.

## Special notes
- [Any deviations from the generic guide]
```

Do not put project-specific assumptions in the generic sections unless they apply to all projects.
