// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package migrate

import (
	"fmt"

	"github.com/searKing/sole/internal/pkg/cmd/sole/migrate/sql"
	"github.com/searKing/sole/internal/pkg/provider"
	"github.com/searKing/sole/internal/pkg/provider/viper"
	"github.com/spf13/cobra"
)

// represent the migrate command
func New() *cobra.Command {
	migrateCmd := &cobra.Command{
		Use:   "migrate",
		Short: "Apply migration plans",
		Long: fmt.Sprintf(`Run this command on a fresh %[1]s installation and when you upgrade %[1]s to a new minor version. For example,
upgrading %[1]s 0.7.0 to 0.8.0 requires running this command.

It is recommended to run this command close to the SQL instance (e.g. same subnet) instead of over the public internet.
This decreases risk of failure and decreases time required.

### WARNING ###

Before running this command on an existing database, create a back up!
`, viper.ServiceName),
		Run: nil,
	}

	// Here you will define your flags and configuration settings.

	var cfgFile string
	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.
	migrateCmd.AddCommand(sql.New())
	migrateCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "",
		fmt.Sprintf("Config file (default is %q)", viper.DefaultConfigPath()))
	migrateCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		provider.InitGlobalProvider(provider.NewProvider(cmd.Context(), cfgFile))
	}

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// migrateCmd.PersistentFlags().String("foo", "", "A help for foo")
	migrateCmd.PersistentFlags().BoolP("yes", "y", false, "If set all confirmation requests are accepted without user interaction.")
	migrateCmd.PersistentFlags().String("dsn", "", "SQL-compatible dsn for sql migrations.")
	return migrateCmd
}
