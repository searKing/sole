// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package web

import (
	"github.com/ory/graceful"
	"github.com/searKing/golang/third_party/github.com/grpc-ecosystem/grpc-gateway/grpc"
	"github.com/searKing/sole/internal/pkg/net/serve"
	"github.com/searKing/sole/internal/pkg/provider"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/grpclog"
)

func ServeGRPC(srv *grpc.Gateway) (graceful.StartFunc, graceful.ShutdownFunc) {
	return serve.Serve(srv.Addr, srv)
}

func setupGRPC() *grpc.Gateway {
	c := provider.GlobalProvider()
	logger := provider.GlobalProvider().Logger().WithField("module", "grpc_gateway")
	grpclog.SetLoggerV2(grpclog.NewLoggerV2(
		logger.WriterLevel(logrus.InfoLevel),
		logger.WriterLevel(logrus.WarnLevel),
		logger.WriterLevel(logrus.ErrorLevel)))
	opts := grpc.WithDefaultMarsherOption()
	opts = append(opts, grpc.WithLogrusLogger(c.Logger()))
	return grpc.NewGatewayTLS(c.GetBackendBindHostPort(), c.TLSConfig(), opts...)
}
