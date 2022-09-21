// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package deploy

import (
	"fmt"

	"github.com/searKing/golang/go/version"
	"github.com/spf13/cobra"
)

// NewStartCommand represents the start command
func NewStartCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "start",
		Short: "start signals to the OS service manager the given service should start",
		Long: fmt.Sprintf(`start signals to the OS service manager the given service should start.

To learn more about each individual command, run:

- %[1]s help deploy start
`, version.ServiceName),
		// stop printing usage when the command errors
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return RunService(ServiceActionStart)
		},
	}
}
