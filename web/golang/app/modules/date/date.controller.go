// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package date

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/searKing/golang/third_party/github.com/grpc-ecosystem/grpc-gateway-v2/grpc"
	"github.com/searKing/sole/api/protobuf-spec/v1/date"
	"github.com/sirupsen/logrus"
	grpc_ "google.golang.org/grpc"
)

type Controller struct {

	// Embed the unimplemented server
	date.UnimplementedDateServiceServer
}

func NewController() *Controller {
	return &Controller{}
}

// SetRoutes registers this handler's routes.
func (c *Controller) SetRoutes(ginRouter gin.IRouter, grpcRouter *grpc.Gateway) {
	logrus.Info("installing router")

	grpcRouter.RegisterGRPCFunc(func(srv *grpc_.Server) {
		date.RegisterDateServiceServer(srv, c)
	})
	_ = grpcRouter.RegisterHTTPFunc(context.Background(), func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc_.DialOption) error {
		return date.RegisterDateServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
	})
}

// Now 日期查询
func (c *Controller) Now(_ context.Context, req *date.DateRequest) (resp *date.DateResponse, err error) {
	return &date.DateResponse{
		RequestId: req.GetRequestId(),
		Date:      time.Now().String(),
	}, nil
}
