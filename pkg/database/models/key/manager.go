// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package key

import (
	"context"

	"gopkg.in/square/go-jose.v2"
)

// type check when compiling
var _, _ Manager = new(SQLManager), new(MemoryManager)

type Manager interface {
	AddKey(ctx context.Context, setId string, key *jose.JSONWebKey) error

	AddKeySet(ctx context.Context, setId string, keys *jose.JSONWebKeySet) error

	GetKey(ctx context.Context, setId, keyId string) (*jose.JSONWebKeySet, error)

	GetKeySet(ctx context.Context, setId string) (*jose.JSONWebKeySet, error)

	DeleteKey(ctx context.Context, setId, keyId string) error

	DeleteKeySet(ctx context.Context, setId string) error
}
