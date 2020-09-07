// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sql

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/ory/x/sqlcon"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/searKing/sole/internal/pkg/provider"
)

func Migrate(dsn string, yes bool, direct migrate.MigrationDirection) error {
	if len(dsn) == 0 {
		dsn = provider.GlobalProvider().Proto().GetDatabase().GetDsn()
	}
	if len(dsn) == 0 {
		return fmt.Errorf("DSN must be set with --dsn or config file")
	}
	provider.GlobalProvider().Proto().Database.Dsn = dsn
	provider.GlobalProvider().EnableMigrateDependencyModules()

	logger := provider.GlobalProvider().Logger().WithField("module", "cmd.migrate.sql")

	scheme := sqlcon.GetDriverName(provider.GlobalProvider().Proto().GetDatabase().GetDsn())

	plan, err := provider.GlobalProvider().SchemaMigrationPlan(scheme, direct)
	if err != nil {
		return fmt.Errorf("an error occurred planning migrations: %w", err)
	}

	fmt.Println("The following migration is planned:")
	fmt.Println("")
	plan.Render()

	if !yes {
		fmt.Println("")
		fmt.Println("To skip the next question use flag --yes (at your own risk).")
		confirmed, err := askForConfirmation("Do you wish to execute this migration plan?")
		if err != nil {
			return fmt.Errorf("ask for confirmation failed: %w", err)
		}
		if !confirmed {
			fmt.Println("Migration aborted.")
			return nil
		}
	}

	n, err := provider.GlobalProvider().MigrateSchemas(scheme, direct)
	if err != nil {
		return fmt.Errorf("an error occurred while connecting to SQL: %w", err)
	}
	logger.Printf("Successfully applied %d SQL migrations!\n", n)
	return nil
}

func askForConfirmation(s string) (bool, error) {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s [y/n]: ", s)

		response, err := reader.ReadString('\n')
		if err != nil {
			return false, fmt.Errorf("read string failed %w", err)
		}

		response = strings.ToLower(strings.TrimSpace(response))
		if response == "y" || response == "yes" {
			return true, nil
		} else if response == "n" || response == "no" {
			return false, nil
		}
	}
}
