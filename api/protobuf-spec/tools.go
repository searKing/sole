// +build tools

package protobuf_spec

import (
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway"
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2"
	_ "google.golang.org/grpc/cmd/protoc-gen-go-grpc"
	_ "google.golang.org/protobuf/cmd/protoc-gen-go"
)

//go:generate go get -u github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway
//go:generate go get -u github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2
//go:generate go get -u google.golang.org/protobuf/cmd/protoc-gen-go
//go:generate go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc
//go:generate go get -u github.com/searKing/golang/tools/cmd/protoc-gen-go-tag
//go:generate go get -u github.com/go-swagger/go-swagger/cmd/swagger
