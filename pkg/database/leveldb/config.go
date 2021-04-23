// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package leveldb

import (
	_ "github.com/go-sql-driver/mysql"
	viper_ "github.com/searKing/golang/third_party/github.com/spf13/viper"
	"github.com/searKing/golang/third_party/github.com/syndtr/goleveldb/leveldb"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

type Config struct {
	GetViper func() *viper.Viper // If set, overrides params below
	Proto    LevelDB
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
	return &Config{
		Proto: LevelDB{
			PathPrefix: "./database/leveldb/db_",
		},
	}
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
func (c *Config) Validate() []error {
	return nil
}

// Complete fills in any fields not set that are required to have valid data and can be derived
// from other fields. If you're going to ApplyOptions, do that first. It's mutating the receiver.
// ApplyOptions is called inside.
func (c *Config) Complete() CompletedConfig {
	if err := c.loadViper(); err != nil {
		return CompletedConfig{&completedConfig{
			Config:        c,
			completeError: err,
		}}
	}
	return CompletedConfig{&completedConfig{Config: c}}
}

// New creates a new server which logically combines the handling chain with the passed server.
// name is used to differentiate for logging. The handler chain in particular can be difficult as it starts delgating.
// New usually called after Complete
func (c completedConfig) New() (*leveldb.ConsistentDB, error) {
	if c.completeError != nil {
		return nil, c.completeError
	}

	return leveldb.NewConsistentDB(c.Proto.GetPathPrefix(), int(c.Proto.GetPoolSize()), &opt.Options{
		BlockCacheCapacity: int(c.Proto.GetCacheInMb() / 2 * opt.MiB),
		WriteBuffer:        int(c.Proto.GetCacheInMb() / 4 * opt.MiB),
	})
}

func (c *Config) loadViper() error {
	var v *viper.Viper
	if c.GetViper != nil {
		v = c.GetViper()
	}

	if err := viper_.UnmarshalProtoMessageByJsonpb(v, &c.Proto); err != nil {
		logrus.WithError(err).Errorf("load redis config from viper")
		return err
	}
	return nil
}
