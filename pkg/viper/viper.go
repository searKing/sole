// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package viper

import (
	viperhelper "github.com/searKing/golang/third_party/github.com/spf13/viper"
	"github.com/spf13/viper"
)

func GetViper(subName string, envPrefix string) func() *viper.Viper {
	return func() *viper.Viper {
		v := viper.GetViper()
		if v == nil {
			return nil
		}
		if subName != "" {
			v = v.Sub(subName)
			envPrefix = envPrefix + "." + subName
		}
		if v == nil {
			return nil
		}
		v.AutomaticEnv()
		viperhelper.MergeConfigFromENV(v, envPrefix)
		return v
	}
}
