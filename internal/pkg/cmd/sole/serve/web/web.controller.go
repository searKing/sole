// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package web

import (
	"github.com/searKing/sole/internal/pkg/provider"
	"github.com/spf13/cobra"
)

func controller() func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		provider.GlobalProvider().EnableServeModule()
		isDSNAllowed()
		service()
	}
}

func isDSNAllowed() {
	if provider.GlobalProvider().Proto().GetDatabase().GetDsn() == "memory" {
		provider.GlobalProvider().Logger().Fatalf(`When using "sole serve web" the DSN can not be set to "memory".`)
	}
}
