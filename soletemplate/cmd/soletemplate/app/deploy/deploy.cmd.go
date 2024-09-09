// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package deploy

import (
	"fmt"

	"github.com/searKing/golang/go/version"
	"github.com/spf13/cobra"
)

// New represents the deployment command
func New() *cobra.Command {
	deployCmd := &cobra.Command{
		Use:   "deploy",
		Short: "Parent command for controlling a service's lifecycle",
		Long: fmt.Sprintf(`%[1]s exposes four deployments for running %[1]s as a service, install, uninstall, start and stop.

To learn more about each individual command, run:

- %[1]s help deploy install
- %[1]s help deploy uninstall
- %[1]s help deploy start
- %[1]s help deploy stop
`, version.ServiceName),
		// stop printing usage when the command errors
		SilenceUsage: true,
	}
	deployCmd.AddCommand(NewInstallCommand())
	deployCmd.AddCommand(NewUninstallCommand())
	deployCmd.AddCommand(NewStartCommand())
	deployCmd.AddCommand(NewStopCommand())
	return deployCmd
}
