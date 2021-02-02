// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pasta

import (
	"github.com/sirupsen/logrus"

	"github.com/searKing/sole/pkg/crypto/pasta"
)

type Config struct {
	SystemSecret         []byte
	RotatedSystemSecrets [][]byte
}

type completedConfig struct {
	*Config
}

type CompletedConfig struct {
	// Embed a private pointer that cannot be instantiated outside of this package.
	*completedConfig
}

// NewConfig returns a Config struct with the default values
func NewConfig() *Config {
	return &Config{}
}

// Complete fills in any fields not set that are required to have valid data and can be derived
// from other fields. If you're going to ApplyOptions, do that first. It's mutating the receiver.
// ApplyOptions is called inside.
func (o *Config) Complete() CompletedConfig {
	o.installSystemSecretOrDie()
	o.installRotatedSystemSecret()

	return CompletedConfig{&completedConfig{o}}
}

// New creates a new server which logically combines the handling chain with the passed server.
// name is used to differentiate for logging. The handler chain in particular can be difficult as it starts delgating.
// New usually called after Complete
func (c completedConfig) New() *pasta.Pasta {
	return c.installKeyCipherOrDie()
}

// installSystemSecretOrDie allows you to check or generate system secret, but dies on failure.
func (c *Config) installSystemSecretOrDie() {
	logger := logrus.WithField("module", "provider.system_secret")
	secret := c.SystemSecret
	if len(secret) == 0 {
		logger.Warnf("Configuration secrets.system is not set, generating a temporary, random secret...")
		secretBytes := pasta.GenerateSecret(32)
		logger.Warnf("Generated secret: %s", string(secretBytes))
		secretBytes = pasta.HashByteSecret(secretBytes)

		logger.Warnln("Do not use generate secrets in production. The secret will be leaked to the logs.")
		c.SystemSecret = secretBytes
		return
	}

	if len(secret) >= 16 {
		// hashes the secret for consumption by the pasta encryption algorithm which expects exactly 32 bytes.
		c.SystemSecret = pasta.HashByteSecret(secret)
		return
	}

	logger.Fatalf("system secret must be undefined or have at least 16 characters but only has %d characters.", len(secret))
	return
}

// installRotatedSystemSecret allows you to check rotated system secret.
func (c *Config) installRotatedSystemSecret() {
	secrets := c.RotatedSystemSecrets
	if len(secrets) < 2 {
		return
	}
	for _, secret := range secrets[1:] {
		// hashes the secret for consumption by the pasta encryption algorithm which expects exactly 32 bytes.
		c.RotatedSystemSecrets = append(c.RotatedSystemSecrets, pasta.HashByteSecret(secret))
	}
}

// installKeyCipherOrDie allows you to generate a key cipher.
func (c *Config) installKeyCipherOrDie() *pasta.Pasta {
	return pasta.NewFromKey(c.SystemSecret, c.RotatedSystemSecrets)
}
