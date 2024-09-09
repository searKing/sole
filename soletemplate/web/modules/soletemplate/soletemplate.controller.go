// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package soletemplate

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	v1 "github.com/searKing/sole/api/protobuf-spec/soletemplate/v1"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"gocloud.dev/secrets"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	slog_ "github.com/searKing/golang/go/log/slog"
	grpc_gw "github.com/searKing/golang/third_party/github.com/grpc-ecosystem/grpc-gateway-v2/grpc"
	grpc_ "github.com/searKing/golang/third_party/google.golang.org/grpc"
	"github.com/searKing/sole/pkg/domain/logging"
	"github.com/searKing/sole/pkg/domain/truncate"

	"github.com/searKing/sole/soletemplate/pkg/application"
)

// Controller  gRPC协议适配接口层
type Controller struct {
	app application.Application

	keeper *secrets.Keeper

	PreferRegisterHTTPFromEndpoint bool

	// Embed the unimplemented server
	v1.UnimplementedSoleTemplateServiceServer
}

func NewController(app application.Application, keeper *secrets.Keeper) *Controller {
	return &Controller{
		app:    app,
		keeper: keeper,
	}
}

// SetRoutes registers this handler's routes.
func (c *Controller) SetRoutes(ginRouter gin.IRouter, grpcRouter *grpc_gw.Gateway) {
	slog.Info("installing router")
	grpcRouter.RegisterGRPCFunc(func(srv *grpc.Server) { v1.RegisterSoleTemplateServiceServer(srv, c) })
	if c.PreferRegisterHTTPFromEndpoint {
		_ = grpcRouter.RegisterHTTPFunc(context.Background(), v1.RegisterSoleTemplateServiceHandlerFromEndpoint)
	} else {
		_ = grpcRouter.RegisterHTTPNoForwardFunc(context.Background(), func(ctx context.Context, mux *runtime.ServeMux, decorators ...grpc_.UnaryHandlerDecorator) error {
			return v1.RegisterSoleTemplateServiceHandlerServer(ctx, mux, &HttpController{Controller: c, decorators: decorators})
		})
	}
}

// Health 健康监测
func (c *Controller) Health(_ context.Context, req *v1.HealthRequest) (
	resp *v1.HealthResponse, err error) {
	return &v1.HealthResponse{
		RequestId: req.GetRequestId(),
		Status:    "ok",
		Date:      time.Now().String(),
	}, nil
}

// Encrypt 文本加密
func (c *Controller) Encrypt(ctx context.Context, req *v1.EncryptRequest) (*v1.EncryptResponse, error) {
	spanName := "Encrypt"
	ctx, span := otel.Tracer("").Start(ctx, spanName)
	defer span.End()
	span.SetAttributes(attribute.String("request_id", req.GetRequestId()))
	span.SetAttributes(attribute.String("app_id", req.GetAppId()))
	logger := slog.With(logging.SpanAttrs(span)...).With("request_id", req.GetRequestId()).With("app_id", req.GetAppId())

	logger = logger.With("req", truncate.DefaultTruncate(req))

	ciphertext, err := c.keeper.Encrypt(ctx, req.GetPlainText())
	if err != nil {
		logger.With(slog_.Error(err)).Error("failed to encrypt plain text")
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("encrypt plain text: %s", err))
	}
	resp := &v1.EncryptResponse{
		RequestId:  req.GetRequestId(),
		CipherText: ciphertext,
	}

	logger.With("resp", truncate.DefaultTruncate(resp)).Info("success to Encrypt")
	return resp, nil
}
