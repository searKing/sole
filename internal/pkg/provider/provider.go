// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package provider

import (
	"context"
	"net/http"
	"sync"

	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	viper_ "github.com/searKing/sole/api/protobuf-spec/v1/viper"
	"github.com/searKing/sole/pkg/crypto/pasta"
)

type Provider struct {
	proto *viper_.ViperProto

	sqlDB *sqlx.DB

	redis redis.UniversalClient

	keyCipher *pasta.Pasta

	corsHandler func(handler http.Handler) http.Handler

	ctx        context.Context
	reloadOnce sync.Once
}

func (p *Provider) Context() context.Context {
	if p.ctx == nil {
		return context.Background()
	}
	return p.ctx
}

func (p *Provider) Proto() *viper_.ViperProto {
	return p.proto
}

func (p *Provider) KeyCipher() *pasta.Pasta {
	return p.keyCipher
}

func (p *Provider) SqlDB() *sqlx.DB {
	return p.sqlDB
}

func (p *Provider) SqlDBPing() error {
	dsn := p.Proto().GetDatabase().GetDsn()
	switch dsn {
	case "memory":
		// ignore
		return nil
	default:
		return p.SqlDB().PingContext(p.Context())
	}
}
