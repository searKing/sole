// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package start

import (
	"fmt"

	"github.com/searKing/sole/internal/pkg/provider/viper"
	"github.com/spf13/cobra"
)

// represent the start command
func New() *cobra.Command {
	return &cobra.Command{
		Use:   "start",
		Short: "start signals to the OS service manager the given service should start",
		Long: fmt.Sprintf(`start signals to the OS service manager the given service should start.

To learn more about each individual command, run:

- %[1]s help deploy start
`, viper.ServiceName),
		Run: controller(),
	}
}
