// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build wireinject
// +build wireinject

package web

import (
	"github.com/google/wire"
	"github.com/spf13/viper"
)

// This file wires the generic interfaces up to Amazon Web Services (AWS). It
// won't be directly included in the final binary, since it includes a Wire
// injector template function (setupAWS), but the declarations will be copied
// into wire_gen.go when Wire is run.

//go:generate wire
func NewServerRunOptions() (opt *ServerRunOptions, err error) {
	// This will be filled in by Wire with providers from the provider sets in
	// wire.Build.
	wire.Build(NewServerRunOptionsSet)
	return nil, nil
}

var NewServerRunOptionsSet = wire.NewSet(
	wire.Struct(new(ServerRunOptions), "Provider", "WebServerOptions", "AppInfo"),
	viper.GetViper,
	NewProviderConfig,
	NewWebServerConfig,
	NewAppInfoConfig)
