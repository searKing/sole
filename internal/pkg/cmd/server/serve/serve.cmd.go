// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package serve

import (
	"context"
	"fmt"
	"strings"

	filepath_ "github.com/searKing/golang/go/path/filepath"
	"github.com/searKing/sole/pkg/appinfo"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"

	"github.com/searKing/sole/internal/pkg/cmd/server/serve/all"
	"github.com/searKing/sole/internal/pkg/cmd/server/serve/web"
	"github.com/searKing/sole/internal/pkg/provider"
)

// New represent the serve command
func New(ctx context.Context) *cobra.Command {
	serveCmd := &cobra.Command{
		Use:   "serve",
		Short: "Parent command for starting public and administrative HTTP/2 and GRPC APIs",
		Long: fmt.Sprintf(`%[1]s exposes one port, for HTTP and GRPC Server. 
It is recommended to run "%[1]s serve all". If you need granular control over CORS settings or similar, you may
want to run "%[1]s serve web" separately.


To learn more about each individual command, run:

- %[1]s help serve all
- %[1]s help serve web
`, appinfo.ServiceName),
		// stop printing usage when the command errors
		SilenceUsage: true,
		Run:          nil,
	}

	serveCmd.AddCommand(all.New(ctx))
	serveCmd.AddCommand(web.New(ctx))
	// Here you will define your flags and configuration settings.

	var cfgFile string
	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.
	serveCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", DefaultConfigPath(),
		fmt.Sprintf("Config file (default is %q)", DefaultConfigPath()))
	serveCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		// viper allows you to load config from default, config path、env and so on, but dies on failure.
		jwalterweatherman.SetLogOutput(logrus.StandardLogger().Writer())
		jwalterweatherman.SetLogThreshold(jwalterweatherman.LevelWarn)

		viper.AutomaticEnv()                    // read in environment variables that match
		viper.SetEnvPrefix(appinfo.ServiceName) // will be uppercase automatically
		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
		viper.SetConfigFile(cfgFile)

		logger := logrus.WithField("config_path", cfgFile).
			WithField("env_prefix", strings.ToUpper(appinfo.ServiceName))
		err := viper.ReadInConfig()
		if err != nil {
			logger.WithError(err).Errorf("load config file")
			return err
		}
		logger.Infof("load config file")
		return nil
	}

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")
	serveCmd.PersistentFlags().BoolVar(&provider.ForceDisableTls, "dangerous-force-disable-tls", false, "Disable HTTP/2 over TLS (HTTPS) and serve HTTP instead. Never use this in production.")

	return serveCmd
}

// DefaultConfigPath returns config file's default path
func DefaultConfigPath() string {
	// 	return filepath_.Pathify(fmt.Sprintf("$HOME/.%s.yaml", version.ServiceName))
	return filepath_.Pathify(fmt.Sprintf("./conf/%s.yaml", appinfo.ServiceName))
}
