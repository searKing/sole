package date

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/searKing/golang/third_party/github.com/grpc-ecosystem/grpc-gateway/grpc"
	"github.com/searKing/sole/api/protobuf-spec/v1/date"
	grpc2 "google.golang.org/grpc"
)

func Router(router *grpc.Gateway) *grpc.Gateway {
	s := Controller{}
	router.RegisterGRPCFunc(func(srv *grpc2.Server) {
		date.RegisterDateServiceServer(srv, &s)
	})
	router.RegisterHTTPFunc(context.Background(), func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc2.DialOption) error {
		return date.RegisterDateServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	})
	return router
}
