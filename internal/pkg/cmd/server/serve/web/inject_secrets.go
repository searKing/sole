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
	"gocloud.dev/secrets"
	_ "gocloud.dev/secrets/localsecrets"
)

//go:generate wire
// NewSecrets is a Wire injector function that sets up the server using config file.
func NewSecrets(ctx context.Context, config *provider.Provider) (keeper *secrets.Keeper, cleanup func(), err error) {
	// This will be filled in by Wire with providers from the provider sets in
	// wire.Build.
	wire.Build(openSecrets)
	return nil, nil, nil
}

func openSecrets(ctx context.Context, opt *provider.Provider) (keeper *secrets.Keeper, cleanup func(), err error) {
	keeper, err = secrets.OpenKeeper(ctx, opt.Proto.GetSecret())
	if err != nil {
		return nil, nil, err
	}
	return keeper, func() {
		keeper.Close()
	}, nil
}
