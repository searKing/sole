// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package serve

import (
	"fmt"

	"github.com/searKing/sole/internal/pkg/cmd/server/serve/all"
	"github.com/searKing/sole/internal/pkg/cmd/server/serve/web"
	"github.com/searKing/sole/internal/pkg/provider"
	"github.com/searKing/sole/internal/pkg/provider/viper"
	"github.com/spf13/cobra"
)

// represent the serve command
func New() *cobra.Command {
	serveCmd := &cobra.Command{
		Use:   "serve",
		Short: "Parent command for starting public and administrative HTTP/2 and GRPC APIs",
		Long: fmt.Sprintf(`%[1]s exposes one port, for HTTP and GRPC Server. 
It is recommended to run "%[1]s serve all". If you need granular control over CORS settings or similar, you may
want to run "%[1]s serve web" separately.


To learn more about each individual command, run:

- %[1]s help serve all
- %[1]s help serve web
`, viper.ServiceName),
		// stop printing usage when the command errors
		SilenceUsage: true,
		Run:          nil,
	}

	serveCmd.AddCommand(all.New())
	serveCmd.AddCommand(web.New())
	// Here you will define your flags and configuration settings.

	var cfgFile string
	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.
	serveCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", viper.DefaultConfigPath(),
		fmt.Sprintf("Config file (default is %q)", viper.DefaultConfigPath()))
	serveCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		cfg := provider.NewConfig()
		cfg.ConfigFile = cfgFile
		return cfg.Complete().Apply(cmd.Context())
	}

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	return serveCmd
}
