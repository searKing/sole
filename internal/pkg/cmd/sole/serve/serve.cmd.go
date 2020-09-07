// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package serve

import (
	"fmt"

	"github.com/searKing/sole/internal/pkg/cmd/sole/serve/all"
	"github.com/searKing/sole/internal/pkg/cmd/sole/serve/web"
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
		Run: nil,
	}

	serveCmd.AddCommand(all.New())
	serveCmd.AddCommand(web.New())
	// Here you will define your flags and configuration settings.

	var cfgFile string
	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.
	serveCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "",
		fmt.Sprintf("Config file (default is %q)", viper.DefaultConfigPath()))
	serveCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		provider.InitGlobalProvider(provider.NewProvider(cmd.Context(), cfgFile))
	}

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")
	serveCmd.PersistentFlags().BoolVar(&viper.ForceDisableTls, "dangerous-force-disable-tls", false, "Disable HTTP/2 over TLS (HTTPS) and serve HTTP instead. Never use this in production.")

	return serveCmd
}
