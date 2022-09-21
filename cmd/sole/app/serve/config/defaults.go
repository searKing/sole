// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package config

// ForceDisableTls ...
var ForceDisableTls bool

// SetDefaults assigns default values for the Config
func (c *Config) SetDefaults() {
	c.Proto.SecretKeeperUrl = "base64key://"
}
