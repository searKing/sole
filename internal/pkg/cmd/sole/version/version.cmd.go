// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package version

import (
	"github.com/spf13/cobra"
)

// represent the version command
func New() *cobra.Command {
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Display this binary's version, build time and git hash of this build",
		Run:   controller(),
	}
	return versionCmd
}
