// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package app

import (
	"context"

	"github.com/google/uuid"
	"github.com/spf13/cobra"

	"github.com/searKing/golang/go/version"

	deploycmd "github.com/searKing/sole/soletemplate/cmd/soletemplate/app/deploy"
	servecmd "github.com/searKing/sole/soletemplate/cmd/soletemplate/app/serve"
	versioncmd "github.com/searKing/sole/soletemplate/cmd/soletemplate/app/version"
)

func init() {
	version.ServiceName = "soletemplate"
	version.ServiceDescription = "soletemplate is a cloud native high throughput service manager server, allowing you to manage all services."
	version.ServiceId = uuid.New().String()
}

// NewCommand This represents the base command when called without any sub commands
func NewCommand(ctx context.Context) *cobra.Command {
	// This represents the base command when called without any subcommands
	rootCmd := &cobra.Command{
		Use:     version.ServiceName,
		Short:   version.ServiceDescription,
		Version: version.Get().String(),
		// Uncomment the following line if your bare application
		// has an action associated with it:
		//Run: func(cmd *cobra.Command, args []string) {},

		// stop printing usage when the command errors
		SilenceUsage: true,
	}
	rootCmd.AddCommand(versioncmd.New())
	rootCmd.AddCommand(deploycmd.New())
	rootCmd.AddCommand(servecmd.New(ctx))

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	cobra.OnInitialize(func() {})

	return rootCmd
}
