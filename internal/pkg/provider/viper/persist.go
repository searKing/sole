// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package viper

import (
	"log"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	filepath_ "github.com/searKing/golang/go/path/filepath"
	"github.com/spf13/viper"
)

// persistConfig writes config using into .use.<name>.yml
func persistConfig() error {
	// persist using config
	f := viper.ConfigFileUsed() // /root/.sole.yml
	if f == "" {
		log.Printf("[WARN] persist skiped, for no config file used\n")
		return nil
	}
	dir := filepath.Dir(f)
	base := filepath.Base(f)
	ext := filepath.Ext(f)
	name := strings.TrimPrefix(strings.TrimSuffix(base, ext), ".")

	configFileUsing := filepath.Join(dir, ".use."+name+".yml") // /root/.use.sole.yml

	err := viper.WriteConfigAs(configFileUsing)
	if err != nil {
		log.Printf("[WARN] %s\n",
			errors.WithMessagef(err, "write using config file [%s] failed...", filepath_.Pathify(configFileUsing)))
	}
	return err
}
