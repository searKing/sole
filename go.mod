module github.com/searKing/sole

go 1.14

require (
	github.com/common-nighthawk/go-figure v0.0.0-20200609044655-c4b36f998cf2
	github.com/gin-gonic/gin v1.7.2
	github.com/go-playground/validator/v10 v10.7.0
	github.com/go-redis/redis/v8 v8.11.0
	github.com/go-sql-driver/mysql v1.6.0
	github.com/google/uuid v1.3.0
	github.com/google/wire v0.5.0
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.6.0
	github.com/hashicorp/consul/api v1.10.1
	github.com/jmoiron/sqlx v1.3.4
	github.com/julienschmidt/httprouter v1.3.0
	github.com/kardianos/service v1.2.0
	github.com/klauspost/compress v1.13.6 // indirect
	github.com/kr/pretty v0.3.0 // indirect
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/opentracing/opentracing-go v1.2.0
	github.com/pborman/uuid v1.2.1
	github.com/pkg/errors v0.9.1
	github.com/pkg/profile v1.5.0
	github.com/prometheus/client_golang v1.11.0
	github.com/rogpeppe/go-internal v1.8.0 // indirect
	github.com/rs/cors v1.7.0
	github.com/searKing/golang v1.1.47
	github.com/segmentio/kafka-go v0.4.15
	github.com/sirupsen/logrus v1.7.0
	github.com/spf13/cobra v1.2.1
	github.com/spf13/jwalterweatherman v1.1.0
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.9.0
	github.com/syndtr/goleveldb v1.0.0
	go.uber.org/automaxprocs v1.4.0
	golang.org/x/crypto v0.0.0-20210921155107-089bfa567519 // indirect
	golang.org/x/net v0.0.0-20211101193420-4a448f8816b3
	google.golang.org/genproto v0.0.0-20211102202547-e9cf271f7f2c
	google.golang.org/grpc v1.42.0
	google.golang.org/grpc/cmd/protoc-gen-go-grpc v1.1.0
	google.golang.org/protobuf v1.27.1
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect
	gopkg.in/square/go-jose.v2 v2.6.0
)

//replace github.com/searKing/golang v1.1.47 => ../../../github.com/searKing/golang
