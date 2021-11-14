// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build wireinject
// +build wireinject

package web

import (
	"context"

	"github.com/google/wire"
	"github.com/searKing/sole/pkg/appinfo"
	"github.com/spf13/viper"
)

// This file wires the generic interfaces up to Amazon Web Services (AWS). It
// won't be directly included in the final binary, since it includes a Wire
// injector template function (setupAWS), but the declarations will be copied
// into wire_gen.go when Wire is run.

func NewAppInfoConfig(v *viper.Viper) *appinfo.Config {
	return appinfo.NewViperConfig(v, "app_info")
}

//go:generate wire
// NewAppInfo is a Wire injector function that sets up the server using config file.
func NewAppInfo(ctx context.Context, opt *ServerRunOptions) (err error) {
	// This will be filled in by Wire with providers from the provider sets in
	// wire.Build.
	wire.Build(
		wire.FieldsOf(new(*ServerRunOptions), "AppInfo"),
		appinfo.NewAppInfo)
	return nil
}
