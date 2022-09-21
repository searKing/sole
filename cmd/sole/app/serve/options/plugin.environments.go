// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package options

import (
	"context"
	"os"

	"github.com/pkg/errors"
	configpb "github.com/searKing/sole/api/protobuf-spec/v1/config"
	"github.com/sirupsen/logrus"
)

type _env struct{}

// NewEnv 根据服务设置环境变量
func NewEnv(ctx context.Context, config *configpb.Configuration) (_ *_env, err error) {
	defer func() { err = errors.WithStack(err) }()
	logrus.Infof("Environments updating")
	for k, v := range config.GetDynamicEnvironments() {
		err := os.Setenv(k, v)
		if err != nil {
			return nil, err
		}
	}

	return &_env{}, nil
}
