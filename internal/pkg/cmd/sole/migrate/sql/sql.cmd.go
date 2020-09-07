// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sql

import (
	"fmt"

	"github.com/searKing/sole/internal/pkg/cmd/sole/migrate/sql/down"
	"github.com/searKing/sole/internal/pkg/cmd/sole/migrate/sql/up"
	"github.com/searKing/sole/internal/pkg/provider/viper"
	"github.com/spf13/cobra"
)

// represent the all command
func New() *cobra.Command {

	sqlCmd := &cobra.Command{
		Use:   "sql",
		Short: "Migrates sql",
		Long: fmt.Sprintf(`Run this command on a fresh SQL installation and when you upgrade|downgrade %[1]s to a new| an old minor version.

It is recommended to run this command close to the SQL instance (e.g. same subnet) instead of over the public internet.
This decreases risk of failure and decreases time required.

### WARNING ###

Before running this command on an existing database, create a back up!
`, viper.ServiceName),
	}
	sqlCmd.AddCommand(up.New())
	sqlCmd.AddCommand(down.New())
	return sqlCmd
}
