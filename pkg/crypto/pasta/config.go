// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pasta

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	viper_ "github.com/searKing/golang/third_party/github.com/spf13/viper"
)

type Config struct {
	KeyInViper string
	Viper      *viper.Viper // If set, overrides params below
	Proto      Secret
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

// NewViperConfig returns a Config struct with the global viper instance
// key representing a sub tree of this instance.
// NewViperConfig is case-insensitive for a key.
func NewViperConfig(key string) *Config {
	c := NewConfig()
	c.Viper = viper.GetViper()
	c.KeyInViper = key
	return c
}

// Validate checks Config and return a slice of found errs.
func (s *Config) Validate() []error {
	return nil
}

// Complete fills in any fields not set that are required to have valid data and can be derived
// from other fields. If you're going to ApplyOptions, do that first. It's mutating the receiver.
// ApplyOptions is called inside.
func (o *Config) Complete() CompletedConfig {
	return CompletedConfig{&completedConfig{o}}
}

// New creates a new server which logically combines the handling chain with the passed server.
// name is used to differentiate for logging. The handler chain in particular can be difficult as it starts delgating.
// New usually called after Complete
func (c completedConfig) New() *Pasta {
	c.loadViperOrDie()
	return c.installKeyCipherOrDie()
}

func (c *completedConfig) loadViperOrDie() {
	v := c.Viper
	if v != nil && c.KeyInViper != "" {
		v = v.Sub(c.KeyInViper)
	}

	if err := viper_.UnmarshalProtoMessageByJsonpb(c.Viper, &c.Proto); err != nil {
		logrus.WithError(err).Fatalf("load secret config from viper")
		return
	}
}

// installSystemSecretOrDie allows you to check or generate system secret, but dies on failure.
func (c *Config) installSystemSecretOrDie() {
	var SystemSecret []byte
	var RotatedSystemSecrets [][]byte

	secrets := c.Proto.GetSystemSecrets()
	if len(secrets) > 0 {
		SystemSecret = secrets[0]
	}
	if len(secrets) > 1 {
		for _, secret := range secrets[1:] {
			RotatedSystemSecrets = append(RotatedSystemSecrets, secret)
		}
	}

	RotatedSystemSecrets = append(RotatedSystemSecrets, c.Proto.GetRotatedSystemSecrets()...)
	logger := logrus.WithField("module", "provider.system_secret")
	if len(SystemSecret) == 0 {
		logger.Warnf("Configuration secrets.system is not set, generating a temporary, random secret...")
		secretBytes := GenerateSecret(32)
		logger.Warnf("Generated secret: %s", string(secretBytes))
		logger.Warnln("Do not use generate secrets in production. The secret will be leaked to the logs.")
		SystemSecret = secretBytes
	}

	if len(SystemSecret) >= 16 {
		// hashes the secret for consumption by the pasta encryption algorithm which expects exactly 32 bytes.
		SystemSecret = HashByteSecret(SystemSecret)
		c.Proto.SystemSecrets = [][]byte{SystemSecret}
		return
	}

	logger.Fatalf("system secret must be undefined or have at least 16 characters but only has %d characters.", len(SystemSecret))
	return
}

// installRotatedSystemSecret allows you to check rotated system secret.
func (c *Config) installRotatedSystemSecret() {
	secrets := c.Proto.GetRotatedSystemSecrets()
	for i, secret := range secrets {
		// hashes the secret for consumption by the pasta encryption algorithm which expects exactly 32 bytes.
		c.Proto.RotatedSystemSecrets[i] = HashByteSecret(secret)
	}
}

// installKeyCipherOrDie allows you to generate a key cipher.
func (c *Config) installKeyCipherOrDie() *Pasta {
	c.installSystemSecretOrDie()
	c.installRotatedSystemSecret()

	return NewFromKey(c.Proto.GetSystemSecrets()[0], c.Proto.GetRotatedSystemSecrets())
}
