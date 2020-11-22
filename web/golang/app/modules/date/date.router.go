// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package date

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/searKing/golang/third_party/github.com/grpc-ecosystem/grpc-gateway/grpc"
	"github.com/searKing/sole/api/protobuf-spec/v1/date"
	grpc_ "google.golang.org/grpc"
)

func Router(router *grpc.Gateway) *grpc.Gateway {
	s := Controller{}
	router.RegisterGRPCFunc(func(srv *grpc_.Server) {
		date.RegisterDateServiceServer(srv, &s)
	})
	_ = router.RegisterHTTPFunc(context.Background(), func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc_.DialOption) error {
		return date.RegisterDateServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	})
	return router
}
