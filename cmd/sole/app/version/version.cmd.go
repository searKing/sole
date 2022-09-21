// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package version

import (
	"fmt"

	"github.com/searKing/golang/go/version"
	"github.com/spf13/cobra"
)

// New represent the version command
func New() *cobra.Command {
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Display this binary's version, build time and git hash of this build",
		// stop printing usage when the command errors
		SilenceUsage: true,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("%+v\n", version.Get())
		},
	}
	return versionCmd
}
