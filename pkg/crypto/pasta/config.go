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
	GetViper func() *viper.Viper // If set, overrides params below
	Proto    Secret
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
func NewViperConfig(getViper func() *viper.Viper) *Config {
	c := NewConfig()
	c.GetViper = getViper
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
	var v *viper.Viper
	if c.GetViper != nil {
		v = c.GetViper()
	}

	if err := viper_.UnmarshalProtoMessageByJsonpb(v, &c.Proto); err != nil {
		logrus.WithError(err).Fatalf("load secret config from viper")
		return
	}
}

// installSystemSecretOrDie allows you to check or generate system secret, but dies on failure.
func (c *Config) installSystemSecretOrDie() {
	var SystemSecret = c.Proto.GetSystemSecret()

	logger := logrus.WithField("module", "provider.system_secret")
	if len(SystemSecret) == 0 {
		logger.Warnf("Configuration secrets.system is not set, generating a temporary, random secret...")
		secretBytes := GenerateSecret(32)
		logger.Warnf("Generated secret: %s", string(secretBytes))
		logger.Warnln("Do not use generate secrets in production. The secret will be leaked to the logs.")
		SystemSecret = string(secretBytes)
	}

	// hashes the secret for consumption by the pasta encryption algorithm which expects exactly 32 bytes.
	SystemSecret = string(HashByteSecret([]byte(SystemSecret)))
	c.Proto.SystemSecret = SystemSecret
	return
}

// installRotatedSystemSecret allows you to check rotated system secret.
func (c *Config) installRotatedSystemSecret() {
	secrets := c.Proto.GetRotatedSystemSecrets()
	for i, secret := range secrets {
		// hashes the secret for consumption by the pasta encryption algorithm which expects exactly 32 bytes.
		c.Proto.RotatedSystemSecrets[i] = string(HashByteSecret([]byte(secret)))
	}
}

// installKeyCipherOrDie allows you to generate a key cipher.
func (c *Config) installKeyCipherOrDie() *Pasta {
	c.installSystemSecretOrDie()
	c.installRotatedSystemSecret()

	var secrets [][]byte
	for _, s := range c.Proto.GetRotatedSystemSecrets() {
		secrets = append(secrets, []byte(s))
	}
	return NewFromKey([]byte(c.Proto.GetSystemSecret()), secrets)
}
