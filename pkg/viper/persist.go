// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package viper

import (
	"path/filepath"
	"strings"

	filepath_ "github.com/searKing/golang/go/path/filepath"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// PersistConfig writes config using into .use.<name>.yaml
func PersistConfig() error {
	// persist using config
	f := viper.ConfigFileUsed() // ./conf/.sole.yaml
	if f == "" {
		logrus.Warnf("persist skiped, for no config file used")
		return nil
	}
	dir := filepath.Dir(f)
	base := filepath.Base(f)
	ext := filepath.Ext(f)
	name := strings.TrimPrefix(strings.TrimSuffix(base, ext), ".")

	configFileUsing := filepath.Join(dir, ".use."+name+".yaml") // /root/.use.sole.yaml

	err := viper.WriteConfigAs(configFileUsing)
	if err != nil {
		logrus.WithField("file", filepath_.Pathify(configFileUsing)).WithError(err).Errorf("write using config file")
		return err
	}
	return nil
}
