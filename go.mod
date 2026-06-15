module afr_Loyalty

go 1.25.3

replace afr_lendme => ../afr_Lendme/

replace afr_auth_center => ../afr_auth_center/

replace daoc => ../daoc/

replace afr_sb_in => ../afr_sb_in/

replace afr_unified_charging_gateway => ../afr_unified_charging_gateway/

replace afr_SpinAndWin_be => ../afr_SpinAndWin_be/

replace afr_propylaea => ../afr_propylaea/

replace afr_sb_mm => ../afr_sb_mm/

replace afr_ao_apgw_v2 => ../afr_ao_apgw_v2/

replace afr_sb_mm_jigsaw => ../afr_sb_mm_jigsaw/

replace afr_sb_MM_KD => ../afr_sb_MM_KD/

replace mongox => ../mongox/

replace redisx => ../redisx/

require (
	afr_SpinAndWin_be v0.0.0-00010101000000-000000000000
	afr_ao_apgw_v2 v0.0.0-00010101000000-000000000000
	afr_auth_center v0.0.0-00010101000000-000000000000
	afr_lendme v0.0.0-00010101000000-000000000000
	afr_propylaea v0.0.0-00010101000000-000000000000
	afr_sb_in v0.0.0-00010101000000-000000000000
	afr_sb_mm v0.0.0-00010101000000-000000000000
	afr_unified_charging_gateway v0.0.0-00010101000000-000000000000
	github.com/gocarina/gocsv v0.0.0-20240520201108-78e41c74b4b1
	github.com/google/uuid v1.6.0
	github.com/gorilla/mux v1.8.1
	github.com/jinzhu/copier v0.4.0
	github.com/kardianos/service v1.2.4
	github.com/prometheus/client_golang v1.23.2
	github.com/rs/cors v1.11.1
	go.mongodb.org/mongo-driver v1.17.9
	go.mongodb.org/mongo-driver/v2 v2.5.0
	mongox v0.0.0-00010101000000-000000000000
	redisx v0.0.0-00010101000000-000000000000
)

require (
	afr_sb_MM_KD v0.0.0-00010101000000-000000000000 // indirect
	afr_sb_mm_jigsaw v0.0.0-00010101000000-000000000000 // indirect
	daoc v0.0.0-00010101000000-000000000000 // indirect
	github.com/360EntSecGroup-Skylar/excelize v1.4.1 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/azcore v1.21.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/azidentity v1.13.1 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/internal v1.11.2 // indirect
	github.com/AzureAD/microsoft-authentication-library-for-go v1.6.0 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/gabriel-vasile/mimetype v1.4.13 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/golang-jwt/jwt/v5 v5.3.1 // indirect
	github.com/golang/snappy v1.0.0 // indirect
	github.com/gorilla/websocket v1.5.3 // indirect
	github.com/klauspost/compress v1.18.3 // indirect
	github.com/kr/fs v0.1.0 // indirect
	github.com/kylelemons/godebug v1.1.0 // indirect
	github.com/microsoft/kiota-abstractions-go v1.9.3 // indirect
	github.com/microsoft/kiota-authentication-azure-go v1.3.1 // indirect
	github.com/microsoft/kiota-http-go v1.5.4 // indirect
	github.com/microsoft/kiota-serialization-form-go v1.1.2 // indirect
	github.com/microsoft/kiota-serialization-json-go v1.1.2 // indirect
	github.com/microsoft/kiota-serialization-multipart-go v1.1.2 // indirect
	github.com/microsoft/kiota-serialization-text-go v1.1.3 // indirect
	github.com/microsoftgraph/msgraph-sdk-go v1.94.0 // indirect
	github.com/microsoftgraph/msgraph-sdk-go-core v1.4.0 // indirect
	github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826 // indirect
	github.com/montanaflynn/stats v0.7.1 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/pkg/browser v0.0.0-20240102092130-5ac0b6a4141c // indirect
	github.com/pkg/sftp v1.13.10 // indirect
	github.com/prometheus/client_model v0.6.2 // indirect
	github.com/prometheus/common v0.67.5 // indirect
	github.com/prometheus/procfs v0.19.2 // indirect
	github.com/redis/go-redis/v9 v9.17.2 // indirect
	github.com/std-uritemplate/std-uritemplate/go/v2 v2.0.8 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.2.0 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/youmark/pkcs8 v0.0.0-20240726163527-a2c0da244d78 // indirect
	go.opentelemetry.io/auto/sdk v1.2.1 // indirect
	go.opentelemetry.io/otel v1.40.0 // indirect
	go.opentelemetry.io/otel/metric v1.40.0 // indirect
	go.opentelemetry.io/otel/trace v1.40.0 // indirect
	go.yaml.in/yaml/v2 v2.4.3 // indirect
	golang.org/x/crypto v0.49.0 // indirect
	golang.org/x/image v0.35.0 // indirect
	golang.org/x/net v0.51.0 // indirect
	golang.org/x/sync v0.20.0 // indirect
	golang.org/x/sys v0.42.0 // indirect
	golang.org/x/text v0.35.0 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
	gopkg.in/alexcesaro/quotedprintable.v3 v3.0.0-20150716171945-2caba252f4dc // indirect
	gopkg.in/ezzarghili/recaptcha-go.v4 v4.3.0 // indirect
	gopkg.in/gomail.v2 v2.0.0-20160411212932-81ebce5c23df // indirect
	gopkg.in/square/go-jose.v2 v2.6.0 // indirect
)
