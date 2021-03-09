module github.com/searKing/sole

go 1.14

require (
	github.com/common-nighthawk/go-figure v0.0.0-20200609044655-c4b36f998cf2
	github.com/gin-gonic/gin v1.6.3
	github.com/go-redis/redis/v8 v8.5.0
	github.com/go-sql-driver/mysql v1.5.0
	github.com/golang/protobuf v1.4.3
	github.com/google/uuid v1.2.0
	github.com/gorilla/handlers v1.5.1
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.1.0
	github.com/hashicorp/consul/api v1.8.1
	github.com/jmoiron/sqlx v1.3.1
	github.com/julienschmidt/httprouter v1.3.0
	github.com/kardianos/service v1.2.0
	github.com/magiconair/properties v1.8.4 // indirect
	github.com/mitchellh/mapstructure v1.4.1 // indirect
	github.com/opentracing/opentracing-go v1.2.0
	github.com/pborman/uuid v1.2.1
	github.com/pelletier/go-toml v1.8.1 // indirect
	github.com/pkg/errors v0.9.1
	github.com/pkg/profile v1.5.0
	github.com/prometheus/client_golang v1.9.0
	github.com/rs/cors v1.7.0
	github.com/searKing/golang v0.0.134
	github.com/searKing/golang/tools/cmd/protoc-gen-go-tag v0.0.0-20210305033433-3133a74ad7bd
	github.com/sirupsen/logrus v1.7.0
	github.com/spf13/afero v1.5.1 // indirect
	github.com/spf13/cast v1.3.1 // indirect
	github.com/spf13/cobra v1.1.1
	github.com/spf13/jwalterweatherman v1.1.0
	github.com/spf13/viper v1.7.1 // indirect
	github.com/uber/jaeger-client-go v2.25.0+incompatible
	go.uber.org/automaxprocs v1.4.0
	golang.org/x/crypto v0.0.0-20201124201722-c8d3bf9c5392 // indirect
	golang.org/x/mod v0.4.1 // indirect
	golang.org/x/net v0.0.0-20210119194325-5f4716e94777
	golang.org/x/sys v0.0.0-20210124154548-22da62e12c0c // indirect
	golang.org/x/text v0.3.5 // indirect
	golang.org/x/tools v0.1.0 // indirect
	google.golang.org/genproto v0.0.0-20210204154452-deb828366460
	google.golang.org/grpc v1.35.0
	google.golang.org/grpc/cmd/protoc-gen-go-grpc v1.1.0
	google.golang.org/protobuf v1.25.0
	gopkg.in/ini.v1 v1.62.0 // indirect
	gopkg.in/square/go-jose.v2 v2.5.1
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

//replace go.opentelemetry.io/otel v0.16.0 => github.com/open-telemetry/opentelemetry-go v0.16.0

//replace github.com/searKing/golang v0.0.134 => ../../../github.com/searKing/golang
