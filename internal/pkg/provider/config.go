// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package provider

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang/protobuf/ptypes"
	"github.com/jmoiron/sqlx"
	logrus_ "github.com/searKing/golang/third_party/github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus"
	jaegerConfig "github.com/uber/jaeger-client-go/config"

	viper_ "github.com/searKing/sole/api/protobuf-spec/v1/viper"
	"github.com/searKing/sole/internal/pkg/provider/viper"
	pasta2 "github.com/searKing/sole/pkg/crypto/pasta"
	"github.com/searKing/sole/pkg/database/sql"
	"github.com/searKing/sole/pkg/logs"
	"github.com/searKing/sole/pkg/net/cors"
	"github.com/searKing/sole/pkg/opentrace"
)

//go:generate go-option -type=Config
type Config struct {
	ConfigFile string

	proto *viper_.ViperProto

	Logs       *logs.Config
	OpenTracer *opentrace.Config

	CORS      *cors.Config
	KeyCipher *pasta2.Config
	Sql       *sql.Config
}

// NewConfig returns a Config struct with the default values
func NewConfig() *Config {
	return &Config{
		Logs:       logs.NewConfig(),
		OpenTracer: opentrace.NewConfig(),
		CORS:       cors.NewConfig(),
		KeyCipher:  pasta2.NewConfig(),
		Sql:        sql.NewConfig(),
	}
}

// Complete fills in any fields not set that are required to have valid data and can be derived
// from other fields. If you're going to ApplyOptions, do that first. It's mutating the receiver.
// ApplyOptions is called inside.
func (o *Config) Complete(options ...ConfigOption) CompletedConfig {
	o.ApplyOptions(options...)
	o.installViperProtoOrDie()
	o.completeLogs()
	o.completeOpenTraceOrDie()
	o.completeCors()
	o.completeKeyCipherOrDie()
	o.completeSqlDBOrDie()

	return CompletedConfig{&completedConfig{o}}
}

// New creates a new server which logically combines the handling chain with the passed server.
// name is used to differentiate for logging. The handler chain in particular can be difficult as it starts delgating.
// New usually called after Complete
func (c completedConfig) New(ctx context.Context) (*Provider, error) {
	var err error
	var corsHandler func(http.Handler) http.Handler
	var sqlDB *sqlx.DB

	if err := c.Logs.Complete().Apply(); err != nil {
		return nil, err
	}
	if closer, err := c.OpenTracer.Complete().Apply(); err != nil {
		go func() {
			select {
			case <-ctx.Done():
				err := closer.Close()
				if err != nil {
					logrus.WithError(err).Error("openTracing closed")
					return
				}
				logrus.Info("openTracing closed")
			}
		}()
		return nil, err
	}
	if c.CORS != nil {
		corsHandler, err = c.CORS.Complete().New()
		if err != nil {
			return nil, err
		}
	}

	if c.Sql != nil {
		sqlDB = c.Sql.Complete().New(ctx)
	}
	return &Provider{
		proto:       c.proto,
		sqlDB:       sqlDB,
		keyCipher:   c.KeyCipher.Complete().New(),
		corsHandler: corsHandler,
		ctx:         ctx,
	}, nil
}

// Apply set options and something else as global init, act likes New but without Config's instance
// Apply usually called after Complete
func (c completedConfig) Apply(ctx context.Context) error {
	provider, err := c.New(ctx)
	if err != nil {
		return err
	}
	InitGlobalProvider(provider)
	return nil
}

// installViperProtoOrDie allows you to load config from default, config path、env and so on, but dies on failure.
func (c *Config) installViperProtoOrDie() {
	tmp, err := viper.Load(c.ConfigFile, viper.NewDefaultViperProto())
	if err != nil {
		logrus.WithError(err).WithField("config_path", c.ConfigFile).Error("load")
	}
	c.proto = tmp
}

func (c *Config) completeLogs() {
	log := c.proto.GetLog()
	c.Logs.ReportCaller = log.GetReportCaller()

	if log.GetFormat() == viper_.Log_json {
		c.Logs.Formatter = &logrus.JSONFormatter{
			CallerPrettyfier: logrus_.ShortCallerPrettyfier,
		}
	} else if log.GetFormat() == viper_.Log_text {
		c.Logs.Formatter = &logrus.TextFormatter{
			CallerPrettyfier: logrus_.ShortCallerPrettyfier,
		}
	}

	level, err := logrus.ParseLevel(log.GetLevel().String())
	if err != nil {
		level = c.Logs.Level
		logrus.WithField("module", "log").WithField("log_level", log.GetLevel()).
			WithError(err).
			Warnf("malformed log level, use %s instead", level)
	}
	c.Logs.Level = level

	duration, err := ptypes.Duration(log.GetRotationDuration())
	if err != nil {
		duration = c.Logs.RotateDuration
		logrus.WithField("module", "log").WithField("rotation_duration", log.GetRotationDuration()).
			WithError(err).
			Warnf("malformed rotation duration, use %s instead", duration)
	}
	c.Logs.RotateDuration = duration

	maxAge, err := ptypes.Duration(log.GetRotationMaxAge())
	if err != nil {
		maxAge = c.Logs.RotateMaxAge
		logrus.WithField("module", "log").WithField("max_age", log.GetRotationMaxAge()).
			WithError(err).
			Warnf("malformed max age, use %s instead", maxAge)
	}
	c.Logs.RotateMaxAge = maxAge
	c.Logs.RotateMaxCount = int(log.GetRotationMaxCount())
}

func (c *Config) completeOpenTraceOrDie() {
	trace := c.proto.GetTracing()
	if !trace.GetEnable() {
		c.OpenTracer.Enabled = false
		return
	}

	c.OpenTracer.Enabled = true
	switch trace.GetType() {
	case viper_.Tracing_urber_jaeger:
		c.OpenTracer.Type = opentrace.TypeJeager
	case viper_.Tracing_zipkin:
		c.OpenTracer.Type = opentrace.TypeZipkin
	default:
		c.OpenTracer.Type = opentrace.TypeButt
		logrus.Fatalf("malformed trace type: %s", trace.GetType())
		return
	}
	c.OpenTracer.ServiceName = c.proto.GetService().GetName()

	reporter := trace.GetJaeger().GetReporter()
	if reporter != nil {
		c.OpenTracer.Configuration.Reporter = &jaegerConfig.ReporterConfig{}
		c.OpenTracer.Configuration.Reporter.LocalAgentHostPort = reporter.GetLocalAgentHostPort()
	}

	sampler := trace.GetJaeger().GetSampler()
	if sampler != nil {
		c.OpenTracer.Configuration.Sampler = &jaegerConfig.SamplerConfig{}
		c.OpenTracer.Configuration.Sampler.SamplingServerURL = sampler.GetServerUrl()
		c.OpenTracer.Configuration.Sampler.Type = sampler.GetType().String()
		c.OpenTracer.Configuration.Sampler.Param = float64(sampler.GetParam())
	}
	return
}

func (c *Config) completeCors() {
	corsConfig := c.CORS
	corsInfo := c.proto.GetWeb().GetCors()
	if corsInfo != nil {
		if corsInfo.Enable {
			maxAge, err := ptypes.Duration(corsInfo.GetMaxAge())
			if err != nil {
				maxAge = corsConfig.MaxAge
				logrus.WithField("max_age", corsInfo.GetMaxAge()).
					WithError(err).
					Warnf("malformed max_age, use %s instead", maxAge)
			}
			corsConfig.UseConditional = corsInfo.GetUseConditional()
			corsConfig.AllowedOrigins = corsInfo.GetAllowedOrigins()
			corsConfig.AllowedMethods = corsInfo.GetAllowedOrigins()
			corsConfig.AllowedHeaders = corsInfo.GetAllowedHeaders()
			corsConfig.ExposedHeaders = corsInfo.GetExposedHeaders()

			corsConfig.MaxAge = maxAge
			corsConfig.AllowCredentials = corsInfo.GetAllowCredentials()
		} else {
			corsConfig = nil
		}
	}
	c.CORS = corsConfig
}

// completeKeyCipherOrDie allows you to generate a key cipher.
func (c *Config) completeKeyCipherOrDie() {
	pastaConfig := c.KeyCipher

	secrets := c.proto.GetSecret().GetSystem()
	if len(secrets) > 0 {
		pastaConfig.SystemSecret = []byte(secrets[0])
	}
	if len(secrets) > 1 {
		for _, secret := range secrets[1:] {
			pastaConfig.RotatedSystemSecrets = append(pastaConfig.RotatedSystemSecrets, []byte(secret))
		}
	}
}

func (c *Config) completeSqlDBOrDie() {
	sqlConfig := c.Sql
	sqlConfig.Dsn = c.proto.GetDatabase().GetDsn()

	dsnUrl := c.proto.GetDatabase().GetDsn()
	switch dsnUrl {
	case "memory":
		// ignore
		return
	case "":
		logrus.Fatalf(`config.database.dsn is not set, use "export %s_DATABASE_DSN=memory" for an in memory storage or the documented database adapters.`,
			strings.ToUpper(viper.ServiceName))
	}

	maxWait, err := ptypes.Duration(c.proto.GetDatabase().GetMaxWaitDuration())
	if err != nil {
		maxWait = sqlConfig.MaxWait
		logrus.WithField("max_wait", c.proto.GetDatabase().GetMaxWaitDuration()).
			WithError(err).
			Warnf("malformed max_wait, use %s instead", maxWait)
	}

	failAfter, err := ptypes.Duration(c.proto.GetDatabase().GetFailAfterDuration())
	if err != nil {
		failAfter = sqlConfig.FailAfter
		logrus.WithField("fail_after", c.proto.GetDatabase().GetFailAfterDuration()).
			WithError(err).
			Warnf("malformed fail_after, use %s instead", failAfter)
	}
}
