// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package provider

import (
	"context"

	viper_ "github.com/searKing/sole/api/protobuf-spec/v1/viper"
)

var ForceDisableTls bool

type Provider struct {
	Proto *viper_.ViperProto

	ctx context.Context
}

func NewProvider(ctx context.Context, config *Config) (*Provider, error) {
	return config.Complete().New(ctx)
}

func (p *Provider) Context() context.Context {
	if p.ctx == nil {
		return context.Background()
	}
	return p.ctx
}
