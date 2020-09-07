// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package key

import (
	"github.com/ory/x/dbal"
	"github.com/ory/x/logrusx"
	"github.com/pkg/errors"
	"github.com/rubenv/sql-migrate"
)

var Migrations = map[string]*dbal.PackrMigrationSource{
	dbal.DriverMySQL: dbal.NewMustPackerMigrationSource(logrusx.New("", ""), AssetNames(), Asset, []string{
		"migrations/sql/shared",
		"migrations/sql/mysql",
	}, true),
	dbal.DriverPostgreSQL: dbal.NewMustPackerMigrationSource(logrusx.New("", ""), AssetNames(), Asset, []string{
		"migrations/sql/shared",
		"migrations/sql/postgres",
	}, true),
	dbal.DriverCockroachDB: dbal.NewMustPackerMigrationSource(logrusx.New("", ""), AssetNames(), Asset, []string{
		"migrations/sql/cockroach",
	}, true),
}

func (m *SQLManager) PlanMigration(dbName string, direct migrate.MigrationDirection) ([]*migrate.PlannedMigration, error) {
	migrate.SetTable(m.MigrateTableName())
	plans, _, err := migrate.PlanMigration(m.db.DB, dbal.Canonicalize(m.db.DriverName()), Migrations[dbName], direct, 0)
	return plans, errors.WithStack(err)
}

func (m *SQLManager) MigrateSchemas(dbName string, direct migrate.MigrationDirection) (int, error) {
	migrate.SetTable(m.MigrateTableName())
	n, err := migrate.Exec(m.db.DB, dbal.Canonicalize(m.db.DriverName()), Migrations[dbName], direct)
	if err != nil {
		return 0, errors.Wrapf(err, "Could not migrate sql schema, applied %d migrations", n)
	}
	return n, nil
}
