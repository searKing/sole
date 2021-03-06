// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uninstall

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/searKing/sole/internal/pkg/provider"
)

// represent the uninstall command
func New() *cobra.Command {
	return &cobra.Command{
		Use:   "uninstall",
		Short: "uninstall removes the given service from the OS service manager",
		Long: fmt.Sprintf(`install removes the given service from the OS service manager. This may require
greater rights. Will return an error if the service is not present.

To learn more about each individual command, run:

- %[1]s help uninstall
`, provider.ServiceName),
		// stop printing usage when the command errors
		SilenceUsage: true,
		Run:          controller(),
	}
}
