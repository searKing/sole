// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build wireinject
// +build wireinject

package web

import (
	"context"

	"github.com/google/wire"
	"github.com/searKing/sole/internal/pkg/provider"
	"github.com/spf13/viper"
)

// This file wires the generic interfaces up to Amazon Web Services (AWS). It
// won't be directly included in the final binary, since it includes a Wire
// injector template function (setupAWS), but the declarations will be copied
// into wire_gen.go when Wire is run.

func NewProviderConfig(v *viper.Viper) *provider.Config {
	return provider.NewViperConfig(v)
}

//go:generate wire
// NewConfig is a Wire injector function that sets up the server using config file.
func NewConfig(ctx context.Context, opt *ServerRunOptions) (config *provider.Provider, err error) {
	// This will be filled in by Wire with providers from the provider sets in
	// wire.Build.
	wire.Build(
		wire.FieldsOf(new(*ServerRunOptions), "Provider"),
		provider.NewProvider)
	return nil, nil
}
