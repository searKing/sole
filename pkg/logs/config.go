// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package logs

import (
	"time"

	logrus_ "github.com/searKing/golang/third_party/github.com/sirupsen/logrus"
	"github.com/searKing/sole/pkg/protobuf"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	GetViper func() *viper.Viper // If set, overrides params below
	Proto    Log
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
	var c Config
	c.SetDefaultsConfig()
	return &c
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
// from other fields. If you're going to `ApplyOptions`, do that first. It's mutating the receiver.
func (c *Config) Complete() CompletedConfig {
	if err := c.loadViper(); err != nil {
		return CompletedConfig{&completedConfig{
			Config:        c,
			completeError: err,
		}}
	}
	return CompletedConfig{&completedConfig{Config: c}}
}

// Apply creates a new server which logically combines the handling chain with the passed server.
// name is used to differentiate for logging. The handler chain in particular can be difficult as it starts delgating.
func (c completedConfig) Apply() error {
	if c.completeError != nil {
		return c.completeError
	}
	return c.install()
}

func (c *Config) loadViper() error {
	var v *viper.Viper
	if c.GetViper != nil {
		v = c.GetViper()
	}
	err := c.SetViperConfig(v)
	if err != nil {
		logrus.WithError(err).Errorf("load logs config from viper")
		return err
	}
	return nil
}

func (c *completedConfig) install() error {
	if c.Proto.GetFormat() == Log_json {
		logrus.SetFormatter(&logrus.JSONFormatter{
			CallerPrettyfier: logrus_.ShortCallerPrettyfier,
		})
	} else if c.Proto.GetFormat() == Log_text {
		logrus.SetFormatter(&logrus.TextFormatter{
			CallerPrettyfier: logrus_.ShortCallerPrettyfier,
		})
	}

	level, err := logrus.ParseLevel(c.Proto.GetLevel().String())
	if err != nil {
		level = logrus.InfoLevel
		logrus.WithField("module", "log").WithField("log_level", c.Proto.GetLevel()).
			WithError(err).
			Warnf("malformed log level, use %s instead", level)
	}
	logrus.SetLevel(level)

	var RotateDuration = protobuf.DurationOrDefault(c.Proto.GetRotationDuration(), 24*time.Hour, "rotation_duration")
	var RotateMaxAge = protobuf.DurationOrDefault(c.Proto.GetRotationMaxAge(), 7*24*time.Hour, "rotation_max_age")
	var RotateMaxCount = int(c.Proto.GetRotationMaxCount())

	logrus.SetReportCaller(c.Proto.GetReportCaller())

	if err := logrus_.WithRotate(logrus.StandardLogger(),
		c.Proto.GetPath(),
		logrus_.WithRotateInterval(RotateDuration),
		logrus_.WithMaxCount(RotateMaxCount),
		logrus_.WithMaxAge(RotateMaxAge),
		logrus_.WithMuteDirectlyOutput(c.Proto.GetMuteDirectlyOutput())); err != nil {
		logrus.WithField("module", "log").WithField("path", c.Proto.GetPath()).
			WithField("duration", RotateDuration).
			WithField("max_count", RotateMaxCount).
			WithField("max_age", RotateMaxAge).
			WithField("mute_directly_output", c.Proto.GetMuteDirectlyOutput()).
			WithError(err).Error("add rotation wrapper for log")
		return err
	}
	logrus.WithField("module", "log").
		WithField("path", c.Proto.GetPath()).
		WithField("duration", RotateDuration).
		WithField("max_count", RotateMaxCount).
		WithField("max_age", RotateMaxAge).
		WithField("mute_directly_output", c.Proto.GetMuteDirectlyOutput()).
		Infof("add rotation wrapper for log")
	return nil
}
