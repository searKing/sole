// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sql

import (
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/opentracing/opentracing-go"
	"github.com/ory/x/sqlcon"
	"github.com/searKing/sole/internal/pkg/provider/viper"
	"github.com/searKing/sole/pkg/database/dsn"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Dsn       string
	MaxWait   time.Duration
	FailAfter time.Duration
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
		Dsn:       "memory",
		MaxWait:   5 * time.Second,
		FailAfter: 5 * time.Minute,
	}
}

// Complete fills in any fields not set that are required to have valid data and can be derived
// from other fields. If you're going to ApplyOptions, do that first. It's mutating the receiver.
// ApplyOptions is called inside.
func (o *Config) Complete() CompletedConfig {
	o.installSqlDBOrDie()

	return CompletedConfig{&completedConfig{o}}
}

// New creates a new server which logically combines the handling chain with the passed server.
// name is used to differentiate for logging. The handler chain in particular can be difficult as it starts delgating.
// New usually called after Complete
func (c completedConfig) New() *sqlx.DB {
	return c.installSqlDBOrDie()
}

func (c *Config) installSqlDBOrDie() *sqlx.DB {
	var options []sqlcon.OptionModifier
	if opentracing.IsGlobalTracerRegistered() {
		options = append(options, sqlcon.WithDistributedTracing(), sqlcon.WithOmitArgsFromTraceSpans())
	}

	dsnUrl := c.Dsn
	switch dsnUrl {
	case "memory":
		// ignore
		return nil
	case "":
		logrus.Fatalf(`config.database.dsn is not set, use "export %s_DATABASE_DSN=memory" for an in memory storage or the documented database adapters.`,
			strings.ToUpper(viper.ServiceName))
	}

	maxWait := c.MaxWait
	failAfter := c.FailAfter

	schema, sdnConfig, err := dsn.ParseDSN(dsnUrl)
	if err != nil {
		logrus.WithField("dsn", dsnUrl).
			Fatalf(`malformed DSN, must set as %s `, "schema://[user[:password]@][net[(addr)]]/dbname[?param1=value1&paramN=valueN]")
		return nil
	}
	if sdnConfig.DBName != "" {
		dbName := sdnConfig.DBName
		sdnConfig.DBName = ""
		trimDatabaseDSN := dsn.GetDSN(schema, sdnConfig)
		connection, _ := sqlcon.NewSQLConnection(trimDatabaseDSN, nil, options...)
		db, err := connection.GetDatabaseRetry(maxWait, failAfter)
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

	connection, _ := sqlcon.NewSQLConnection(dsnUrl, nil, options...)
	db, err := connection.GetDatabaseRetry(maxWait, failAfter)
	if err != nil {
		logrus.WithField("dsn", dsnUrl).
			WithField("max_wait", maxWait).
			WithField("fail_after", failAfter).
			Fatalf(`unable to initialize database`)
	}
	return db
}
