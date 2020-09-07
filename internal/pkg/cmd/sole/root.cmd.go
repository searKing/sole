// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package sole

import (
	"github.com/searKing/sole/internal/pkg/cmd/sole/deploy"
	"github.com/searKing/sole/internal/pkg/cmd/sole/migrate"
	"github.com/searKing/sole/internal/pkg/cmd/sole/serve"
	"github.com/searKing/sole/internal/pkg/cmd/sole/version"
	"github.com/searKing/sole/internal/pkg/provider/viper"
	"github.com/spf13/cobra"
)

// This represents the base command when called without any sub commands
func NewCommand() *cobra.Command {

	// This represents the base command when called without any subcommands
	rootCmd := &cobra.Command{
		Use:     viper.ServiceName,
		Short:   viper.ServiceDescription,
		Version: viper.Version,
		// Uncomment the following line if your bare application
		// has an action associated with it:
		//Run: func(cmd *cobra.Command, args []string) {},
	}
	rootCmd.AddCommand(version.New())
	rootCmd.AddCommand(deploy.New())
	rootCmd.AddCommand(migrate.New())
	rootCmd.AddCommand(serve.New())

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	cobra.OnInitialize(func() {})

	return rootCmd
}
