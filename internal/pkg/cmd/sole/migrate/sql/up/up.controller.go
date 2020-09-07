// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package up

import (
	"log"
	"os"

	"github.com/ory/x/flagx"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/spf13/cobra"

	migrate_ "github.com/searKing/sole/internal/pkg/cmd/sole/migrate/sql/shared/services/migrate"
)

func controller() func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		yes := flagx.MustGetBool(cmd, "yes")
		dsn := flagx.MustGetString(cmd, "dsn")
		if err := migrate_.Migrate(dsn, yes, migrate.Up); err != nil {
			log.Println(cmd.UsageString())
			log.Println("")
			log.Println(err)
			os.Exit(1)
			return
		}
	}
}
