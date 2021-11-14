// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package pasta

import (
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	viper_ "github.com/searKing/golang/third_party/github.com/spf13/viper"
)

type Config struct {
	Proto     Secret
	Validator *validator.Validate

	viper     *viper.Viper
	viperKeys []string
}

type completedConfig struct {
	*Config

	// for Complete Only
	completeError error
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
func NewViperConfig(v *viper.Viper, keys ...string) *Config {
	c := NewConfig()
	c.viper = v
	c.viperKeys = keys
	return c
}

// Validate checks Config.
func (c *completedConfig) Validate() error {
	return c.Validator.Struct(c)
}

// Complete fills in any fields not set that are required to have valid data and can be derived
// from other fields. If you're going to ApplyOptions, do that first. It's mutating the receiver.
// ApplyOptions is called inside.
func (c *Config) Complete() CompletedConfig {
	if c.viper != nil {
		err := viper_.UnmarshalKeys(c.viperKeys, &c.Proto)
		if err != nil {
			return CompletedConfig{&completedConfig{completeError: err}}
		}
	}
	if c.Validator == nil {
		c.Validator = validator.New()
	}
	return CompletedConfig{&completedConfig{Config: c}}
}

// New creates a new server which logically combines the handling chain with the passed server.
// name is used to differentiate for logging. The handler chain in particular can be difficult as it starts delgating.
// New usually called after Complete
func (c completedConfig) New() *Pasta {
	if c.completeError != nil {
		logrus.WithError(c.completeError).Errorf("complete secret config")
		return nil
	}
	err := c.Validate()
	if err != nil {
		logrus.WithError(err).Errorf("validate secret config")
		return nil
	}
	return c.installKeyCipherOrDie()
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
