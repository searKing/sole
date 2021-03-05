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
	"github.com/searKing/golang/third_party/github.com/spf13/viper"
	"github.com/sirupsen/logrus"
	"github.com/spf13/jwalterweatherman"
	jaegerConfig "github.com/uber/jaeger-client-go/config"

	viper_ "github.com/searKing/sole/api/protobuf-spec/v1/viper"
	"github.com/searKing/sole/internal/pkg/version"
	"github.com/searKing/sole/pkg/crypto/pasta"
	redis_ "github.com/searKing/sole/pkg/database/redis"
	"github.com/searKing/sole/pkg/database/sql"
	"github.com/searKing/sole/pkg/logs"
	"github.com/searKing/sole/pkg/net/cors"
	"github.com/searKing/sole/pkg/opentrace"
	"github.com/searKing/sole/pkg/protobuf"
)

//go:generate go-option -type=Config
type Config struct {
	ConfigFile string

	proto *viper_.ViperProto

	Logs       *logs.Config
	OpenTracer *opentrace.Config

	CORS      *cors.Config
	KeyCipher *pasta.Config
	Sql       *sql.Config
	Redis     *redis_.Config
}

// NewConfig returns a Config struct with the default values
func NewConfig() *Config {
	return &Config{
		Logs:       logs.NewConfig(),
		OpenTracer: opentrace.NewConfig(),
		CORS:       cors.NewConfig(),
		KeyCipher:  pasta.NewConfig(),
		Sql:        sql.NewConfig(),
		Redis:      redis_.NewConfig(),
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
	o.completeRedis()

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
		redis:       c.Redis.Complete().New(),
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
	var v viper_.ViperProto
	jwalterweatherman.SetLogOutput(logrus.StandardLogger().Writer())
	jwalterweatherman.SetLogThreshold(jwalterweatherman.LevelWarn)

	if err := viper.LoadGlobalConfig(&v, c.ConfigFile, version.ServiceName, NewDefaultViperProto()); err != nil {
		logrus.WithError(err).WithField("config_path", c.ConfigFile).Fatalf("load config")
	}

	if err := viper.PersistGlobalConfig(); err != nil {
		logrus.WithError(err).WithField("config_path", c.ConfigFile).Warnf("persist config ignored")
	}

	c.proto = &v
}

func (c *Config) completeLogs() {
	log := c.proto.GetLog()
	c.Logs.Path = log.GetPath()
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
			strings.ToUpper(version.ServiceName))
	}

	sqlConfig.MaxWait = protobuf.DurationOrDefault(c.proto.GetDatabase().GetMaxWaitDuration(), sqlConfig.MaxWait, "max_wait")
	sqlConfig.FailAfter = protobuf.DurationOrDefault(c.proto.GetDatabase().GetFailAfterDuration(), sqlConfig.FailAfter, "fail_after")
}

func (c *Config) completeRedis() {
	redisConfig := c.Redis
	redis := c.proto.GetRedis()
	redisConfig.Addrs = redis.GetAddrs()
	redisConfig.DB = int(redis.GetDb())
	redisConfig.Username = redis.GetUsername()
	redisConfig.Password = redis.GetPassword()
	redisConfig.SentinelPassword = redis.GetSentinelPassword()
	redisConfig.MaxRetries = int(redis.GetMaxRetries())
	redisConfig.MinRetryBackoff = protobuf.DurationOrDefault(redis.GetMinRetryBackoff(), redisConfig.MinRetryBackoff, "min_retry_backoff")
	redisConfig.MaxRetryBackoff = protobuf.DurationOrDefault(redis.GetMaxRetryBackoff(), redisConfig.MaxRetryBackoff, "max_retry_backoff")
	redisConfig.DialTimeout = protobuf.DurationOrDefault(redis.GetDialTimeout(), redisConfig.DialTimeout, "dial_timeout")
	redisConfig.ReadTimeout = protobuf.DurationOrDefault(redis.GetReadTimeout(), redisConfig.ReadTimeout, "read_timeout")
	redisConfig.WriteTimeout = protobuf.DurationOrDefault(redis.GetWriteTimeout(), redisConfig.WriteTimeout, "write_timeout")
	redisConfig.PoolSize = int(redis.GetPoolSize())
	redisConfig.MinIdleConns = int(redis.GetMinIdleConns())
	redisConfig.MaxConnAge = protobuf.DurationOrDefault(redis.GetMaxConnAge(), redisConfig.MaxConnAge, "max_conn_age")
	redisConfig.PoolTimeout = protobuf.DurationOrDefault(redis.GetPoolTimeout(), redisConfig.PoolTimeout, "pool_timeout")
	redisConfig.IdleTimeout = protobuf.DurationOrDefault(redis.GetIdleTimeout(), redisConfig.IdleTimeout, "idle_timeout")
	redisConfig.IdleCheckFrequency = protobuf.DurationOrDefault(redis.GetIdleCheckFrequency(), redisConfig.IdleCheckFrequency, "idle_check_frequency")
	redisConfig.MaxRedirects = int(redis.GetMaxRedirects())
	redisConfig.ReadOnly = redis.GetReadOnly()
	redisConfig.RouteByLatency = redis.GetRouteByLatency()
	redisConfig.RouteRandomly = redis.GetRouteRandomly()
	redisConfig.MasterName = redis.GetMasterName()
}
