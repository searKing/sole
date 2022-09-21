// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package deploy

import (
	"fmt"

	"github.com/searKing/golang/go/version"
	"github.com/spf13/cobra"
)

// NewStopCommand represents the stop command
func NewStopCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "start",
		Short: "stop signals to the OS service manager the given service should stop",
		Long: fmt.Sprintf(`Stop signals to the OS service manager the given service should stop.

To learn more about each individual command, run:

- %[1]s help deploy stop
`, version.ServiceName),
		// stop printing usage when the command errors
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunService(ServiceActionStop)
		},
	}
}
