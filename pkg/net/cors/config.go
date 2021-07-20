// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cors

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/rs/cors"
	gincors "github.com/rs/cors/wrapper/gin"
	viper_ "github.com/searKing/golang/third_party/github.com/spf13/viper"
	"github.com/searKing/sole/pkg/protobuf"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Config struct {
	GetViper func() *viper.Viper // If set, overrides params below
	Proto    CORS
	Validator            *validator.Validate
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
		Proto: CORS{
			AllowedOrigins: []string{"*"},
			AllowedMethods: []string{
				http.MethodHead,
				http.MethodGet,
				http.MethodPost,
				http.MethodPut,
				http.MethodPatch,
				http.MethodDelete},
			AllowedHeaders: []string{"*"},
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
func (c *completedConfig) Validate() error {
	return c.Validator.Struct(c)
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
	if c.Validator == nil {
		c.Validator = validator.New()
	}
	return CompletedConfig{&completedConfig{Config: c}}
}

func (c *Config) loadViper() error {
	var v *viper.Viper
	if c.GetViper != nil {
		v = c.GetViper()
	}

	if err := viper_.UnmarshalProtoMessageByJsonpb(v, &c.Proto); err != nil {
		logrus.WithError(err).Errorf("load cors config from viper")
		return err
	}
	return nil
}

func (c completedConfig) options() cors.Options {
	maxAge := protobuf.DurationOrDefault(c.Proto.GetMaxAge(), 0, "max_age")

	return cors.Options{
		AllowedOrigins:     c.Proto.AllowedOrigins,
		AllowedMethods:     c.Proto.AllowedMethods,
		AllowedHeaders:     c.Proto.AllowedHeaders,
		ExposedHeaders:     c.Proto.ExposedHeaders,
		AllowCredentials:   c.Proto.AllowCredentials,
		MaxAge:             int(maxAge.Truncate(time.Second).Seconds()),
		OptionsPassthrough: c.Proto.OptionsPassthrough,
	}
}

func (c completedConfig) New() *cors.Cors {
	return cors.New(c.options())
}

func (c completedConfig) NewWrapper() func(http.Handler) http.Handler {
	return c.New().Handler
}

func (c completedConfig) NewGinHandler() gin.HandlerFunc {
	return gincors.New(c.options())
}
