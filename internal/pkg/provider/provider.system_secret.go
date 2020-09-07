// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package provider

import (
	"sync/atomic"

	"github.com/searKing/sole/pkg/crypto/pasta"
)

//go:generate go-atomicvalue -type "systemSecret<bytes>"
type systemSecret atomic.Value
type bytes []byte

func (p *Provider) SystemSecret() []byte {
	return p.systemSecret.Load()
}

func (p *Provider) updateSystemSecret() {
	proto := p.Proto()
	logger := p.Logger().WithField("module", "provider.system_secret")
	secret := proto.GetSystemSecret()

	if len(secret) > 0 {
		if len(secret) >= 16 {
			// hashes the secret for consumption by the pasta encryption algorithm which expects exactly 32 bytes.
			p.systemSecret.Store(pasta.HashStringSecret(secret))
			return
		}

		logger.Fatalf("system secret must be undefined or have at least 16 characters but only has %d characters.", len(secret))
		return
	}

	if len(p.systemSecret.Load()) != 0 {
		return
	}

	logger.Warnf("Configuration secrets.system is not set, generating a temporary, random secret...")
	secretBytes := pasta.GenerateSecret(32)
	logger.Warnf("Generated secret: %s", string(secretBytes))
	secretBytes = pasta.HashByteSecret(secretBytes)

	logger.Warnln("Do not use generate secrets in production. The secret will be leaked to the logs.")
	p.systemSecret.Store(secretBytes)
}
