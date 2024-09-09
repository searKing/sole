// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package date

import (
	"context"
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc"

	time_ "github.com/searKing/golang/go/time"
	grpcgateway "github.com/searKing/golang/third_party/github.com/grpc-ecosystem/grpc-gateway-v2/grpc"
	grpc_ "github.com/searKing/golang/third_party/google.golang.org/grpc"

	v1 "github.com/searKing/sole/api/protobuf-spec/sole/date/v1"
	"github.com/searKing/sole/pkg/domain/errorgoogle"
	"github.com/searKing/sole/pkg/domain/logging"
	"github.com/searKing/sole/pkg/domain/truncate"
)

type Controller struct {
	// 是否将HTTP注册为gRPC中转，支持gRPC stream
	PreferRegisterHTTPFromEndpoint bool

	// Embed the unimplemented server
	v1.UnimplementedDateServiceServer
}

func NewController() *Controller {
	return &Controller{}
}

// SetRoutes registers this handler's routes.
func (c *Controller) SetRoutes(ginRouter gin.IRouter, grpcRouter *grpcgateway.Gateway) {
	slog.Info("installing router")
	grpcRouter.RegisterGRPCFunc(func(srv *grpc.Server) { v1.RegisterDateServiceServer(srv, c) })
	if c.PreferRegisterHTTPFromEndpoint {
		_ = grpcRouter.RegisterHTTPFunc(context.Background(), v1.RegisterDateServiceHandlerFromEndpoint)
	} else {
		_ = grpcRouter.RegisterHTTPNoForwardFunc(context.Background(), func(ctx context.Context, mux *runtime.ServeMux, decorators ...grpc_.UnaryHandlerDecorator) error {
			return v1.RegisterDateServiceHandlerServer(ctx, mux, &HttpController{Controller: c, decorators: decorators})
		})
	}
}

// Now 日期查询
func (c *Controller) Now(ctx context.Context, req *v1.DateRequest) (resp *v1.DateResponse, err error) {
	ctx, span := otel.Tracer("").Start(ctx, "Now")
	defer span.End()
	logger := slog.With(logging.SpanAttrs(span)...).With("cmd", truncate.DefaultTruncate(req))
	logger.Info("Now executed")
	return &v1.DateResponse{
		Date: time.Now().String(),
	}, nil
}

// Error 日期查询，只返回错误，测试使用
func (c *Controller) Error(_ context.Context, req *v1.DateRequest) (resp *v1.DateResponse, err error) {
	httpReqs.With(nil).Inc()
	var cost time_.Cost
	cost.Start()
	defer cost.ElapseFunc(func(d time.Duration) {
		httpReqCostInMilliSeconds.Set(float64(d) / float64(time.Millisecond))
	})
	return nil, errorgoogle.Errorf(v1.DateErrorEnum_UNIMPLEMENTED, "method %q not implemented", "Error")
}
