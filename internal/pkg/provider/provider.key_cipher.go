// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package provider

import (
	"sync/atomic"

	"github.com/searKing/sole/pkg/crypto/pasta"
)

//go:generate go-atomicvalue -type "keyCipher<*github.com/searKing/sole/app/core/crypto/pasta.Pasta>"
type keyCipher atomic.Value

// 加密解密类
func (p *Provider) KeyCipher() *pasta.Pasta {
	return p.keyCipher.Load()
}

func (p *Provider) updateKeyCipher() {
	p.keyCipher.Store(pasta.New(&innerKeyCipherProvider{}))
}

type innerKeyCipherProvider struct{}

func (c *innerKeyCipherProvider) GetRotatedSystemSecrets() [][]byte {
	return nil
}

func (c *innerKeyCipherProvider) GetSystemSecret() []byte {
	return GlobalProvider().SystemSecret()
}
