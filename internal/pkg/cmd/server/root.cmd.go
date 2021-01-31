// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package server

import (
	"github.com/spf13/cobra"

	"github.com/searKing/sole/internal/pkg/cmd/server/serve"
	"github.com/searKing/sole/internal/pkg/cmd/server/version"
	"github.com/searKing/sole/internal/pkg/provider/viper"
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

		// stop printing usage when the command errors
		SilenceUsage: true,
	}
	rootCmd.AddCommand(version.New())
	rootCmd.AddCommand(serve.New())

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	cobra.OnInitialize(func() {})

	return rootCmd
}
