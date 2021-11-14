// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build wireinject
// +build wireinject

package web

import (
	"context"

	"github.com/google/wire"
	"github.com/searKing/sole/pkg/webserver"
	"github.com/searKing/sole/web/golang"
	"github.com/spf13/viper"
)

func NewWebServerConfig(v *viper.Viper) *webserver.Config {
	return webserver.NewViperConfig(v, "web")
}

//go:generate wire
// NewWebServer is a Wire injector function that sets up the server using WebServer(grpc+http).
func NewWebServer(ctx context.Context, opt *ServerRunOptions) (ws *webserver.WebServer, cleanup func(), err error) {
	// This will be filled in by Wire with providers from the provider sets in
	// wire.Build.
	wire.Build(
		wire.FieldsOf(new(*ServerRunOptions), "WebServerOptions"),
		golang.NewHandler,
		NewConfig,
		NewSecrets,
		setupWebServer)

	return nil, nil, nil
}

func setupWebServer(ctx context.Context, config *webserver.Config, handler *golang.Handler) (ws *webserver.WebServer, err error) {
	ws, err = webserver.NewWebServer(ctx, config)
	if err != nil {
		return nil, err
	}
	ws.InstallWebHandlers(handler)
	return ws, nil
}
