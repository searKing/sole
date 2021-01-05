// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package logs

import (
	"time"

	logrus_ "github.com/searKing/golang/third_party/github.com/sirupsen/logrus"
	"github.com/searKing/sole/internal/pkg/provider/viper"
	"github.com/sirupsen/logrus"
)

type Config struct {
	ReportCaller bool             // 是否打印调用者堆栈
	Level        logrus.Level     // 日志最低打印等级
	Formatter    logrus.Formatter // 日志格式

	Path           string        // //日志存储路径
	RotateDuration time.Duration // 日志循环覆盖分片时间
	RotateMaxAge   time.Duration // 文件最大保存时间
	RotateMaxCount int           //日志循环覆盖保留分片个数
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
		ReportCaller:   true,
		Level:          logrus.InfoLevel,
		Formatter:      &logrus.TextFormatter{CallerPrettyfier: logrus_.ShortCallerPrettyfier},
		Path:           "./log/" + viper.ServiceName,
		RotateDuration: 24 * time.Hour,
		RotateMaxAge:   7 * 24 * time.Hour,
		RotateMaxCount: 0,
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

// New creates a new server which logically combines the handling chain with the passed server.
// name is used to differentiate for logging. The handler chain in particular can be difficult as it starts delgating.
func (c completedConfig) Apply() error {
	return installLogrus(c.Config)
}

func installLogrus(c *Config) error {
	logrus.SetReportCaller(c.ReportCaller)
	if c.Formatter != nil {
		logrus.SetFormatter(c.Formatter)
	}
	logrus.SetLevel(c.Level)

	if err := logrus_.WithRotation(logrus.StandardLogger(),
		c.Path, c.RotateDuration, uint(c.RotateMaxCount), c.RotateMaxAge); err != nil {
		logrus.WithField("module", "log").WithField("path", c.Path).
			WithField("duration", c.RotateDuration).
			WithField("max_count", c.RotateMaxCount).
			WithField("max_age", c.RotateMaxAge).
			WithError(err).Error("add rotation wrapper for log")
		return err
	}
	return nil
}
