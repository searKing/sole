// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stop

import (
	"fmt"

	"github.com/searKing/sole/internal/pkg/provider/viper"
	"github.com/spf13/cobra"
)

// represent the stop command
func New() *cobra.Command {
	return &cobra.Command{
		Use:   "start",
		Short: "stop signals to the OS service manager the given service should stop",
		Long: fmt.Sprintf(`Stop signals to the OS service manager the given service should stop.

To learn more about each individual command, run:

- %[1]s help deploy stop
`, viper.ServiceName),
		Run: controller(),
	}
}
