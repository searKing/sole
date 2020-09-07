// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package web

import (
	"github.com/spf13/cobra"
)

// represent the web command
func New() *cobra.Command {
	return &cobra.Command{
		Use:   "web",
		Short: "Serves Administrative and Public HTTP/2 and GRPC APIs",
		Long: `This command opens one port and listens to HTTP/2 and GRPC API requests. The exposed API handles administrative and publiuc
requests like managing as an administrator.

This command is configurable using the same options available to "serve all".

It is generally recommended to use this command only if you require granular control over the administrative and public APIs.
For example, you might want to run different TLS certificates or CORS settings on the public and administrative API.

This command does not work with the "memory" database. Both services (administrative, public) MUST use the same database
connection to be able to synchronize.

`,
		Run: controller(),
	}
}
