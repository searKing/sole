// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package config

import (
	"github.com/searKing/sole/api/protobuf-spec/sole/types/v1/configuration"
	v1 "github.com/searKing/sole/api/protobuf-spec/soletemplate/v1"
)

// ForceDisableTls ...
var ForceDisableTls bool

// SetDefaults assigns default values for the Config
func (c *Config) SetDefaults() {
	c.Proto.Category = &v1.Configuration_Category{
		System: &configuration.System{
			SecretKeeperUrl: "base64key://",
		},
	}
}
