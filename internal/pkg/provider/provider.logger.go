// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package provider

import (
	"sync/atomic"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/pkg/errors"
	logrus_ "github.com/searKing/golang/third_party/github.com/sirupsen/logrus"
	"github.com/searKing/sole/api/protobuf-spec/v1/viper"
	"github.com/sirupsen/logrus"
)

//go:generate go-atomicvalue -type "logger<*github.com/sirupsen/logrus.Logger>"
type logger atomic.Value

func (p *Provider) Logger() *logrus.Logger {
	logger := p.logger.Load()
	if logger != nil {
		return logger
	}
	return logrus.StandardLogger()
}

func (p *Provider) updateLogger() {
	proto := p.Proto()
	logger := logrus.New()
	if proto.GetLog().GetFormat() == viper.Log_json {
		logger.Formatter = &logrus.JSONFormatter{}
	} else if proto.GetLog().GetFormat() == viper.Log_text {
		logger.Formatter = &logrus.TextFormatter{}
	}

	level, err := logrus.ParseLevel(proto.GetLog().GetLevel().String())
	logger.Level = level
	if err != nil {
		logger.Level = logrus.InfoLevel
		logger.WithField("module", "provider.logger").WithField("log_level", proto.Log.Level).
			Warnf("malformed log level, use %s instead", logger.Level)
	}

	duration, err := ptypes.Duration(proto.GetLog().GetRotationDuration())
	if err != nil {
		duration = 24 * time.Hour
		logger.WithField("module", "provider.logger").WithField("rotation_duration", proto.GetLog().GetRotationDuration()).
			WithError(errors.WithStack(err)).
			Warnf("malformed rotation duration, use %s instead", duration)
	}

	maxAge, err := ptypes.Duration(proto.GetLog().GetRotationMaxAge())
	if err != nil {
		maxAge = 3 * 24 * time.Hour
		logger.WithField("module", "provider.logger").WithField("max_age", proto.Log.RotationMaxAge).
			WithError(errors.WithStack(err)).
			Warnf("malformed max age, use %s instead", maxAge)
	}
	if err := logrus_.WithRotation(logger,
		proto.GetLog().GetPath(), duration, uint(proto.GetLog().GetRotationMaxCount()), maxAge); err != nil {
		go logger.WithField("module", "provider.logger").WithField("path", proto.GetLog().GetPath()).
			WithField("duration", duration).
			WithField("max_count", proto.GetLog().GetRotationMaxCount()).
			WithField("max_age", maxAge).
			WithError(errors.WithStack(err)).Error("add rotation wrapper for log failed, ignore it")
	}
	p.logger.Store(logger)
}
