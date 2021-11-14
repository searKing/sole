// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package web

import (
	"context"

	"github.com/go-playground/validator/v10"
	grpcopentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	grpcprometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	grpc_ "github.com/searKing/golang/third_party/github.com/grpc-ecosystem/grpc-gateway/v2/grpc"
	"github.com/searKing/sole/internal/pkg/provider"
	"github.com/searKing/sole/pkg/appinfo"
	"github.com/searKing/sole/pkg/webserver"
	"github.com/sirupsen/logrus"
)

// ServerRunOptions runs a kubernetes api server.
type ServerRunOptions struct {
	Provider *provider.Config

	// GRPC+HTTP
	WebServerOptions *webserver.Config

	AppInfo *appinfo.Config
}

type completedServerRunOptions struct {
	*ServerRunOptions
}

// CompletedServerRunOptions is a private wrapper that enforces a call of Complete() before Run can be invoked.
type CompletedServerRunOptions struct {
	// Embed a private pointer that cannot be instantiated outside of this package.
	*completedServerRunOptions
}

// Validate checks ServerRunOptions and return a slice of found errs.
func (s *ServerRunOptions) Validate(validate *validator.Validate) []error {
	var errs []error
	return errs
}

// Complete set default ServerRunOptions.
func (s *ServerRunOptions) Complete() (CompletedServerRunOptions, error) {
	s.WebServerOptions.Proto.ForceDisableTls = provider.ForceDisableTls
	s.WebServerOptions.GatewayOptions = append(s.WebServerOptions.GatewayOptions,
		grpc_.WithGrpcUnaryServerChain(grpcopentracing.UnaryServerInterceptor()),
		grpc_.WithGrpcStreamServerChain(grpcopentracing.StreamServerInterceptor()),
		grpc_.WithGrpcUnaryServerChain(grpcprometheus.UnaryServerInterceptor),
		grpc_.WithGrpcStreamServerChain(grpcprometheus.StreamServerInterceptor))
	return CompletedServerRunOptions{&completedServerRunOptions{s}}, nil
}

// Run runs the specified APIServer.  This should never exit.
func (s *CompletedServerRunOptions) Run(ctx context.Context) error {
	// To help debugging, immediately log version
	logrus.Infof("Version: %+v", appinfo.GetVersion())
	//isDSNAllowedOrDie(completeOptions.Provider.Proto.GetDatabase().GetDsn())

	server, cleanup, err := NewWebServer(ctx, s.ServerRunOptions)
	if err != nil {
		return err
	}
	defer cleanup()

	prepared, err := server.PrepareRun()
	if err != nil {
		return err
	}

	return prepared.Run(ctx)
}
