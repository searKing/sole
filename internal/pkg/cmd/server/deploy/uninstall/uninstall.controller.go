// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uninstall

import (
	"github.com/searKing/sole/internal/pkg/cmd/server/deploy/shared/services/setup"
	"github.com/spf13/cobra"
)

func controller() func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		setup.Setup(setup.ServiceActionUninstall)
	}
}
