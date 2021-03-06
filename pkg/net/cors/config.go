// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cors

import (
	"net/http"
	"time"

	"github.com/gorilla/handlers"
)

type Config struct {
	// returns Access-Control-Allow-Origin: * if false
	UseConditional bool

	// AllowedOrigins is a list of origins a cross-domain request can be executed from.
	// If the special "*" value is present in the list, all origins will be allowed.
	// An origin may contain a wildcard (*) to replace 0 or more characters
	// (i.e.: http://*.domain.com). Usage of wildcards implies a small performance penalty.
	// Only one wildcard can be used per origin.
	// Default value is ["*"]
	AllowedOrigins []string
	// AllowedMethods is a list of methods the client is allowed to use with
	// cross-domain requests. Default value is simple methods (HEAD, GET and POST).
	AllowedMethods []string
	// AllowedHeaders is list of non simple headers the client is allowed to use with
	// cross-domain requests.
	// If the special "*" value is present in the list, all headers will be allowed.
	// Default value is [] but "Origin" is always appended to the list.
	AllowedHeaders []string
	// ExposedHeaders indicates which headers are safe to expose to the API of a CORS
	// API specification
	ExposedHeaders []string
	// MaxAge indicates how long (in seconds) the results of a preflight request
	// can be cached
	MaxAge time.Duration
	// AllowCredentials indicates whether the request can include user credentials like
	// cookies, HTTP authentication or client side SSL certificates.
	AllowCredentials bool
	// OptionsPassthrough instructs preflight to let other potential next handlers to
	// process the OPTIONS method. Turn this on if your application handles OPTIONS.
	OptionsPassthrough bool
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
	return &Config{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{
			http.MethodHead,
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete},
		AllowedHeaders: []string{"*"},
	}
}

// Validate checks Config and return a slice of found errs.
func (s *Config) Validate() []error {
	return nil
}

// Complete fills in any fields not set that are required to have valid data and can be derived
// from other fields. If you're going to `ApplyOptions`, do that first. It's mutating the receiver.
func (s *Config) Complete() CompletedConfig {
	var options completedConfig

	// set defaults
	options.Config = s
	return CompletedConfig{&completedConfig{s}}
}

func (c completedConfig) New() (func(http.Handler) http.Handler, error) {
	return installCors(c.Config)
}

func installCors(c *Config) (func(http.Handler) http.Handler, error) {
	if !c.UseConditional {
		return func(h http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Access-Control-Allow-Origin", "*")
				h.ServeHTTP(w, r)
			})
		}, nil
	}

	var opts []handlers.CORSOption
	if c.AllowCredentials {
		opts = append(opts, handlers.AllowCredentials())
	}
	if c.OptionsPassthrough {
		opts = append(opts, handlers.IgnoreOptions())
	}
	if c.AllowedMethods != nil {
		opts = append(opts, handlers.AllowedMethods(c.AllowedMethods))
	}
	if c.AllowedHeaders != nil {
		opts = append(opts, handlers.AllowedHeaders(c.AllowedHeaders))
	}
	if c.AllowedOrigins != nil {
		opts = append(opts, handlers.AllowedOrigins(c.AllowedOrigins))
	}
	if c.ExposedHeaders != nil {
		opts = append(opts, handlers.ExposedHeaders(c.ExposedHeaders))
	}
	if c.MaxAge >= 0 {
		opts = append(opts, handlers.MaxAge(int(c.MaxAge.Seconds())))
	}

	return handlers.CORS(opts...), nil
}
