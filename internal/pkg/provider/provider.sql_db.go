// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package provider

import (
	"fmt"
	"strings"
	"sync/atomic"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/golang/protobuf/ptypes"
	"github.com/jmoiron/sqlx"
	"github.com/ory/x/logrusx"
	"github.com/ory/x/sqlcon"
	"github.com/pkg/errors"
)

//go:generate go-atomicvalue -type "sqlDB<*github.com/jmoiron/sqlx.DB>"
type sqlDB atomic.Value

func (p *Provider) SqlDB() *sqlx.DB {
	return p.sqlDB.Load()
}

func (p *Provider) SqlDBPing() error {
	dsn := p.Proto().GetDatabase().GetDsn()
	switch dsn {
	case "memory":
		// ignore
		return nil
	default:
		return p.SqlDB().PingContext(p.Context())
	}
}

func (p *Provider) updateSqlDB() {
	proto := p.Proto()
	logger := p.Logger().WithField("module", "provider.database")

	var options []sqlcon.OptionModifier
	if p.Tracer() != nil {
		options = append(options, sqlcon.WithDistributedTracing(), sqlcon.WithOmitArgsFromTraceSpans())
	}

	dsn := p.Proto().GetDatabase().GetDsn()
	switch dsn {
	case "memory":
		// ignore
		return
	case "":
		logger.Fatalf(`config.database.dsn is not set, use "export SOLE_DATABASE_DSN=memory" for an in memory storage or the documented database adapters.`)
	}

	maxWait, err := ptypes.Duration(proto.GetDatabase().GetMaxWaitDuration())
	if err != nil {
		maxWait = 5 * time.Second
		logger.WithField("max_wait", proto.GetDatabase().GetMaxWaitDuration()).
			WithError(errors.WithStack(err)).
			Warnf("malformed max_wait, use %s instead", maxWait)
	}
	failAfter, err := ptypes.Duration(proto.GetDatabase().GetFailAfterDuration())
	if err != nil {
		failAfter = 5 * time.Minute
		logger.WithField("fail_after", proto.GetDatabase().GetFailAfterDuration()).
			WithError(errors.WithStack(err)).
			Warnf("malformed fail_after, use %s instead", failAfter)
	}

	schema, sdnConfig, err := ParseDSN(dsn)
	if err != nil {
		logger.WithField("dsn", dsn).
			Fatalf(`malformed DSN, must set as %s `, "schema://[user[:password]@][net[(addr)]]/dbname[?param1=value1&paramN=valueN]")
	}
	if sdnConfig.DBName != "" {
		dbName := sdnConfig.DBName
		sdnConfig.DBName = ""
		trimDatabaseDSN := GetDSN(schema, sdnConfig)
		logger := logrusx.New("", "")
		logger.Logger = p.Logger()
		connection, _ := sqlcon.NewSQLConnection(trimDatabaseDSN, logger, options...)
		db, err := connection.GetDatabaseRetry(maxWait, failAfter)
		if err != nil {
			logger.WithField("dsn", dsn).
				WithField("max_wait", maxWait).
				WithField("fail_after", failAfter).
				Fatalf(`unable to initialize database`)
			return
		}

		createDatabaseSql := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %[1]s;", dbName)

		_, err = db.Exec(createDatabaseSql)
		if err != nil {
			logger.WithField("dsn", dsn).
				WithField("sql", createDatabaseSql).
				WithError(err).
				Fatalf(`unable to find or create database`)
		}
	}

	loggerx := logrusx.New("", "")
	loggerx.Logger = p.Logger()
	connection, _ := sqlcon.NewSQLConnection(dsn, loggerx, options...)
	db, err := connection.GetDatabaseRetry(maxWait, failAfter)
	if err != nil {
		logger.WithField("dsn", dsn).
			WithField("max_wait", maxWait).
			WithField("fail_after", failAfter).
			Fatalf(`unable to initialize database`)
	}
	p.sqlDB.Store(db)
}

// schema://[user[:password]@][net[(addr)]]/dbname[?param1=value1&paramN=valueN]
func ParseDSN(dsn string) (schema string, cfg *mysql.Config, err error) {
	schema = sqlcon.GetDriverName(dsn)
	prefix := fmt.Sprintf("%s://", schema)
	if strings.HasPrefix(dsn, prefix) {
		dsn = strings.TrimPrefix(dsn, prefix)
	}
	// [user[:password]@][net[(addr)]]/dbname[?param1=value1&paramN=valueN]
	cfg, err = mysql.ParseDSN(dsn)
	return schema, cfg, err
}

// schema://[user[:password]@][net[(addr)]]/dbname[?param1=value1&paramN=valueN]
func GetDSN(schema string, cfg *mysql.Config) string {
	return fmt.Sprintf("%s://%s", schema, cfg.FormatDSN())
}
