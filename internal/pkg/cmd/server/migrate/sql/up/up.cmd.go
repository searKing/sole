// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package up

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/searKing/sole/internal/pkg/provider"
)

// represent the up command
func New() *cobra.Command {
	return &cobra.Command{
		Use:   "up",
		Short: "Migrates upgrade sql",
		Long: fmt.Sprintf(`Run this command on a fresh SQL installation and when you upgrade %[1]s to a new minor version. For example,
upgrading %[1]s 0.7.0 to 0.8.0 requires running this command.

It is recommended to run this command close to the SQL instance (e.g. same subnet) instead of over the public internet.
This decreases risk of failure and decreases time required.

### WARNING ###

Before running this command on an existing database, create a back up!
`, provider.ServiceName),
		// stop printing usage when the command errors
		SilenceUsage: true,
		//Run:          controller(),
	}
}
