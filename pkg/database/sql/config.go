// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sql

import (
	"context"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/opentracing/opentracing-go"
	"github.com/searKing/golang/go/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/searKing/golang/third_party/github.com/go-sql-driver/mysql"
	"github.com/searKing/golang/third_party/github.com/golang/go/database/sql"
	viper_ "github.com/searKing/golang/third_party/github.com/spf13/viper"
	"github.com/searKing/sole/pkg/protobuf"
)

type Config struct {
	Proto     Database
	Validator *validator.Validate

	viper     *viper.Viper
	viperKeys []string
}

type completedConfig struct {
	*Config

	// for Complete Only
	completeError error
}

type CompletedConfig struct {
	// Embed a private pointer that cannot be instantiated outside of this package.
	*completedConfig
}

// NewConfig returns a Config struct with the default values
func NewConfig() *Config {
	return &Config{
		Proto: Database{
			Dsn:               "memory",
			MaxWaitDuration:   durationpb.New(5 * time.Second),
			FailAfterDuration: durationpb.New(5 * time.Minute),
		},
	}
}

// NewViperConfig returns a Config struct with the global viper instance
// key representing a sub tree of this instance.
// NewViperConfig is case-insensitive for a key.
func NewViperConfig(v *viper.Viper, keys ...string) *Config {
	c := NewConfig()
	c.viper = v
	c.viperKeys = keys
	return c
}

// Validate checks Config.
func (c *completedConfig) Validate() error {
	var errs []error
	errs = append(errs, c.Validator.Struct(c))

	dsnUrl := c.Proto.GetDsn()
	switch dsnUrl {
	case "memory":
		// ignore
		break
	case "":
		errs = append(errs, fmt.Errorf(`config.database.dsn is not set, use "memory" for an in memory storage or the documented database adapters`))
	}
	return errors.Multi(errs...)
}

// Complete fills in any fields not set that are required to have valid data and can be derived
// from other fields. If you're going to ApplyOptions, do that first. It's mutating the receiver.
// ApplyOptions is called inside.
func (c *Config) Complete() CompletedConfig {
	if c.viper != nil {
		err := viper_.UnmarshalKeys(c.viperKeys, &c.Proto)
		if err != nil {
			return CompletedConfig{&completedConfig{completeError: err}}
		}
	}
	if c.Validator == nil {
		c.Validator = validator.New()
	}
	return CompletedConfig{&completedConfig{Config: c}}
}

// New creates a new server which logically combines the handling chain with the passed server.
// name is used to differentiate for logging. The handler chain in particular can be difficult as it starts delgating.
// New usually called after Complete
func (c completedConfig) New(ctx context.Context) *sqlx.DB {
	if c.completeError != nil {
		logrus.WithError(c.completeError).Fatalf("complete config")
		return nil
	}
	err := c.Validate()
	if err != nil {
		logrus.WithError(err).Fatalf("validate config")
		return nil
	}
	return c.installSqlDBOrDie(ctx)
}

func (c *Config) installSqlDBOrDie(ctx context.Context) *sqlx.DB {
	var options []sql.DBOption
	if opentracing.IsGlobalTracerRegistered() {
		options = append(options, sql.WithDistributedTracing(), sql.WithOmitArgsFromTraceSpans())
	}

	dsnUrl := c.Proto.GetDsn()
	switch dsnUrl {
	case "memory":
		// ignore
		return nil
	case "":
		logrus.Fatalf(`dsn is not set, use "memory" for an in memory storage or the documented database adapters.`)
	}

	maxWait := protobuf.DurationOrDefault(c.Proto.GetMaxWaitDuration(), 5*time.Second, "max_wait")
	failAfter := protobuf.DurationOrDefault(c.Proto.GetFailAfterDuration(), 5*time.Minute, "fail_after")

	schema, sdnConfig, err := mysql.ParseDSN(dsnUrl)
	if err != nil {
		logrus.WithField("dsn", dsnUrl).
			Fatalf(`malformed DSN, must set as %s `, "schema://[user[:password]@][net[(addr)]]/dbname[?param1=value1&paramN=valueN]")
		return nil
	}
	if sdnConfig.DBName != "" {
		dbName := sdnConfig.DBName
		sdnConfig.DBName = ""
		trimDatabaseDSN := mysql.GetDSN(schema, sdnConfig)
		connection, _ := sql.Open(trimDatabaseDSN, options...)
		db, err := connection.GetDatabaseRetry(ctx, maxWait, failAfter)
		if err != nil {
			logrus.WithField("dsn", dsnUrl).
				WithField("max_wait", maxWait).
				WithField("fail_after", failAfter).
				Fatalf(`unable to initialize database`)
			return nil
		}

		createDatabaseSql := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %[1]s;", dbName)

		_, err = db.Exec(createDatabaseSql)
		if err != nil {
			logrus.WithField("dsn", dsnUrl).
				WithField("sql", createDatabaseSql).
				WithError(err).
				Fatalf(`unable to find or create database`)
		}
	}

	connection, _ := sql.Open(dsnUrl, options...)
	db, err := connection.GetDatabaseRetry(ctx, maxWait, failAfter)
	if err != nil {
		logrus.WithField("dsn", dsnUrl).
			WithField("max_wait", maxWait).
			WithField("fail_after", failAfter).
			Fatalf(`unable to initialize database`)
	}
	return db
}
