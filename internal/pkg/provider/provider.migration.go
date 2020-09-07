// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package provider

import (
	"fmt"
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
	migrate "github.com/rubenv/sql-migrate"
)

type schemaCreator interface {
	MigrateSchemas(dbName string, direct migrate.MigrationDirection) (int, error)
	PlanMigration(dbName string, direct migrate.MigrationDirection) ([]*migrate.PlannedMigration, error)
}

func (p *Provider) SchemaMigrationPlan(dbName string, direct migrate.MigrationDirection) (*tablewriter.Table, error) {

	table := tablewriter.NewWriter(os.Stdout)
	table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	table.SetCenterSeparator("|")
	table.SetAutoMergeCells(true)
	table.SetRowLine(true)
	table.SetColMinWidth(4, 20)
	table.SetHeader([]string{
		"Driver",
		"Module",
		"ID",
		"#",
		"Query",
	})

	for _, s := range p.schemaMigratorModules() {
		plans, err := s.migrator.PlanMigration(dbName, direct)
		if err != nil {
			return nil, err
		}

		for _, plan := range plans {
			var todos []string
			switch direct {
			case migrate.Up:
				todos = plan.Up
			case migrate.Down:
				todos = plan.Down
			}
			for k, todo := range todos {
				todo = strings.Replace(strings.TrimSpace(todo), "\n", "", -1)
				todo = strings.Join(strings.Fields(todo), " ")
				if len(todo) > 0 {
					table.Append([]string{p.SqlDB().DriverName(), s.moduleName, plan.Id + ".sql", fmt.Sprintf("%d", k), todo})
				}
			}
		}
	}

	return table, nil
}

func (p *Provider) MigrateSchemas(dbName string, direct migrate.MigrationDirection) (int, error) {
	logger := p.Logger().WithField("module", "provider.migration")

	var total int

	logger.Debugf("Applying %s SQL migrations...", dbName)
	for k, s := range p.schemaMigratorModules() {
		logger.Debugf("Applying %s SQL migrations for manager: %T (%d)", dbName, s, k)
		if c, err := s.migrator.MigrateSchemas(dbName, direct); err != nil {
			return c, err
		} else {
			logger.Debugf("Successfully applied %d %s SQL migrations from manager: %T (%d)", c, dbName, s, k)
			total += c
		}
	}
	logger.Debugf("Applied %d %s SQL migrations", total, dbName)

	return total, nil
}

func (p *Provider) schemaMigratorModules() []schemaCreatorModule {
	var modules []schemaCreatorModule

	key, ok := p.KeyManager().(schemaCreator)
	if ok {
		modules = append(modules, schemaCreatorModule{
			moduleName: "JSON Web Keys",
			migrator:   key,
		})
	}

	return modules
}

type schemaCreatorModule struct {
	moduleName string
	migrator   schemaCreator
}
