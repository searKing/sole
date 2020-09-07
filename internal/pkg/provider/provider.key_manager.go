// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package provider

import (
	"sync/atomic"

	"github.com/searKing/sole/pkg/database/models/key"
)

//go:generate go-atomicvalue -type "keyManager<github.com/searKing/sole/app/core/models/keys.Manager>"
type keyManager atomic.Value

func (p *Provider) KeyManager() key.Manager {
	return p.keyManager.Load()
}

func (p *Provider) updateKeyManager() {
	logger := p.Logger().WithField("module", "provider.key_manager")
	dsn := p.Proto().GetDatabase().GetDsn()
	switch dsn {
	case "memory":
		p.keyManager.Store(key.NewMemoryManager())
	case "":
		logger.Fatalf(`config.database.dsn is not set, use "export SOLE_DATABASE_DSN=memory" for an in memory storage or the documented database adapters.`)
	default:
		p.keyManager.Store(key.NewSQLManager(p.SqlDB(), p))
	}
}
