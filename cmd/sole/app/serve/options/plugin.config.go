// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package options

import (
	"github.com/pkg/errors"
	configpb "github.com/searKing/sole/api/protobuf-spec/v1/config"
	"github.com/searKing/sole/cmd/sole/app/serve/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// NewConfig 根据全局配置加载服务配置选项
func NewConfig(v *viper.Viper) (c *configpb.Configuration, err error) {
	defer func() { err = errors.WithStack(err) }()
	logrus.Infof("Installing Config")
	return config.NewViperConfig(v, "").Complete().New()
}
