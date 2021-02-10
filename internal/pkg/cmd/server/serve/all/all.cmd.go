// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package all

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/searKing/sole/internal/pkg/version"
)

// represent the all command
func New(ctx context.Context) *cobra.Command {
	return &cobra.Command{
		Use:   "all",
		Short: "Serves both public, administrative HTTP/2 and GRPC APIs",
		Long: fmt.Sprintf(`Starts a process which listens on three ports for public, administrative HTTP/2 and GRPC API requests.

If you want more granular control (e.g. different TLS settings) over each API group (administrative, public) you
can run "serve admin", "serve public", "serve rtsp" separately.

This command exposes a variety of controls via environment variables. You can
set environments using "export KEY=VALUE" (Linux/macOS) or "set KEY=VALUE" (Windows). On Linux,
you can also set environments by prepending key value pairs: "KEY=VALUE KEY2=VALUE2 %[1]s"

service possible controls are listed below. This command exposes exposes command line flags, which are listed below
the controls section.

`, version.ServiceName),
		// stop printing usage when the command errors
		SilenceUsage: true,
		RunE:         CommandE(ctx),
	}
}
