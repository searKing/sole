module github.com/searKing/sole

go 1.14

require (
	github.com/common-nighthawk/go-figure v0.0.0-20200609044655-c4b36f998cf2
	github.com/gin-gonic/gin v1.6.3
	github.com/go-playground/validator/v10 v10.2.0
	github.com/go-redis/redis/v8 v8.5.0
	github.com/go-sql-driver/mysql v1.5.0
	github.com/google/uuid v1.2.0
	github.com/grpc-ecosystem/go-grpc-middleware v1.2.0
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.1.0
	github.com/hashicorp/consul/api v1.8.1
	github.com/jmoiron/sqlx v1.3.1
	github.com/julienschmidt/httprouter v1.3.0
	github.com/kardianos/service v1.2.0
	github.com/opentracing/opentracing-go v1.2.0
	github.com/pborman/uuid v1.2.1
	github.com/pkg/errors v0.9.1
	github.com/pkg/profile v1.5.0
	github.com/prometheus/client_golang v1.9.0
	github.com/rs/cors v1.7.0
	github.com/searKing/golang v1.0.176
	github.com/segmentio/kafka-go v0.4.15
	github.com/sirupsen/logrus v1.7.0
	github.com/spf13/cobra v1.1.1
	github.com/spf13/jwalterweatherman v1.1.0
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.7.1
	github.com/syndtr/goleveldb v1.0.0
	github.com/uber/jaeger-client-go v2.25.0+incompatible
	go.uber.org/automaxprocs v1.4.0
	golang.org/x/net v0.0.0-20210119194325-5f4716e94777
	google.golang.org/genproto v0.0.0-20210204154452-deb828366460
	google.golang.org/grpc v1.35.0
	google.golang.org/grpc/cmd/protoc-gen-go-grpc v1.1.0
	google.golang.org/protobuf v1.26.0
	gopkg.in/square/go-jose.v2 v2.5.1

)

//replace go.opentelemetry.io/otel v0.16.0 => github.com/open-telemetry/opentelemetry-go v0.16.0

//replace github.com/searKing/golang v1.0.170 => ../../../github.com/searKing/golang

//
//replace github.com/searKing/golang/tools/cmd/protoc-gen-go-tag v0.0.0-20210305033433-3133a74ad7bd => ../../../github.com/searKing/golang/tools/cmd/protoc-gen-go-tag
