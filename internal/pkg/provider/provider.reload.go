// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package provider

import (
	"context"
	"time"

	time_ "github.com/searKing/golang/go/time"
	viperhelper "github.com/searKing/golang/third_party/github.com/spf13/viper"
	"github.com/sirupsen/logrus"
	"github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"

	"github.com/searKing/sole/internal/pkg/version"
)

const (
	DefaultTimeout = time.Minute
)

func (p *Provider) ReloadForever() {
	p.reloadOnce.Do(func() {
		func() {
			time_.NonSlidingUntil(p.Context(), func(ctx context.Context) {
				// viper allows you to load config from default, config path、env and so on, but dies on failure.
				jwalterweatherman.SetLogOutput(logrus.StandardLogger().Writer())
				jwalterweatherman.SetLogThreshold(jwalterweatherman.LevelWarn)
				if err := viperhelper.MergeAll(viper.GetViper(), p.ConfigFile, version.ServiceName); err != nil {
					logrus.WithError(err).WithField("config_path", p.ConfigFile).Fatalf("load config")
				}
				providerReloads.WithLabelValues(p.Proto.String()).Inc()
			}, DefaultTimeout)
		}()
	})
}
