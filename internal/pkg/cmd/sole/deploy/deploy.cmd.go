// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package deploy

import (
	"fmt"

	"github.com/searKing/sole/internal/pkg/cmd/sole/deploy/install"
	"github.com/searKing/sole/internal/pkg/cmd/sole/deploy/start"
	"github.com/searKing/sole/internal/pkg/cmd/sole/deploy/stop"
	"github.com/searKing/sole/internal/pkg/cmd/sole/deploy/uninstall"
	"github.com/searKing/sole/internal/pkg/provider/viper"
	"github.com/spf13/cobra"
)

// represent the deploy command
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
`, viper.ServiceName),
	}
	deployCmd.AddCommand(install.New())
	deployCmd.AddCommand(uninstall.New())
	deployCmd.AddCommand(start.New())
	deployCmd.AddCommand(stop.New())
	return deployCmd
}
