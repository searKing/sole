// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package provider

import (
	"context"
	"sync"

	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"github.com/searKing/golang/third_party/github.com/syndtr/goleveldb/leveldb"
	"github.com/searKing/sole/pkg/consul"
	"github.com/searKing/sole/pkg/crypto/pasta"

	viper_ "github.com/searKing/sole/api/protobuf-spec/v1/viper"
)

type Provider struct {
	ConfigFile string

	Proto *viper_.ViperProto

	KeyCipher       *pasta.Pasta
	SqlDB           *sqlx.DB
	Redis           redis.UniversalClient
	LevelDB         *leveldb.ConsistentDB
	ServiceRegister *consul.ServiceRegister
	ServiceResolver *consul.ServiceResolver

	ctx        context.Context
	reloadOnce sync.Once
}

func (p *Provider) Context() context.Context {
	if p.ctx == nil {
		return context.Background()
	}
	return p.ctx
}

func (p *Provider) SqlDBPing() error {
	if p.SqlDB == nil {
		return nil
	}
	return p.SqlDB.PingContext(p.Context())
}
