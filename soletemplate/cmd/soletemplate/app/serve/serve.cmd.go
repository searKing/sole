// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package serve

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"

	slog_ "github.com/searKing/golang/go/log/slog"
	"github.com/searKing/sole/soletemplate/cmd/soletemplate/app/serve/config"
	"github.com/searKing/sole/soletemplate/cmd/soletemplate/app/serve/options"
	"github.com/spf13/cobra"
	"github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"

	filepath_ "github.com/searKing/golang/go/path/filepath"
	"github.com/searKing/golang/go/version"
	"github.com/searKing/golang/go/version/verflag"
	"github.com/searKing/golang/third_party/github.com/spf13/pflag"
)

// New represent the serve command
func New(ctx context.Context) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Parent command for starting public and administrative HTTP/2 and GRPC APIs",
		Long: fmt.Sprintf(`%[1]s exposes one port, for HTTP and GRPC Server. 
It is recommended to run "%[1]s serve all". If you need granular control over CORS settings or similar, you may
want to run "%[1]s serve web" separately.


To learn more about each individual command, run:

- %[1]s help serve all
- %[1]s help serve web
`, version.ServiceName),
		// stop printing usage when the command errors
		SilenceUsage: true,
		Version:      version.Get().String(),

		// Uncomment the following line if your bare application
		// has an action associated with it:
		RunE: func(cmd *cobra.Command, args []string) error {
			slog.With("cmdline", os.Args).Info("boosting")
			defer func() {
				slog.With("cmdline", os.Args).Info("exited")
			}()

			// To help debugging, immediately log version
			slog.Info(fmt.Sprintf("Version: %+v", version.Get()))

			ctx, cancel := context.WithCancel(ctx)
			_, f, err := options.RunServer(ctx, cancel)
			if err != nil {
				return err
			}
			if f != nil {
				defer f()
			}
			return nil
		},
		Args: func(cmd *cobra.Command, args []string) error {
			for _, arg := range args {
				if len(arg) > 0 {
					return fmt.Errorf("%q does not take any arguments, got %q", cmd.CommandPath(), args)
				}
			}
			return nil
		},
	}

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	cobra.OnInitialize(func() {})

	var cfgFiles []string
	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.
	cmd.PersistentFlags().StringArrayVarP(&cfgFiles, "config", "c", []string{DefaultConfigPath()},
		fmt.Sprintf("Config file (default is %q)", DefaultConfigPath()))
	cmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		verflag.PrintAndExitIfRequested()
		pflag.PrintFlags(cmd.Flags())

		// 通过viper从配置文件、环境变量等加载配置
		jwalterweatherman.SetLogOutput(slog.NewLogLogger(slog.Default().Handler(), slog.LevelWarn).Writer())
		jwalterweatherman.SetLogThreshold(jwalterweatherman.LevelWarn)

		viper.AutomaticEnv()                    // read in environment variables that match
		viper.SetEnvPrefix(version.ServiceName) // will be uppercase automatically
		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
		if len(cfgFiles) > 0 {
			viper.SetConfigFile(cfgFiles[0])
		}
		logger := slog.With("env_prefix", strings.ToUpper(version.ServiceName))
		for _, c := range cfgFiles {
			if c == "" {
				continue
			}
			v := viper.New()
			v.SetConfigFile(c)
			err := v.ReadInConfig()
			if err != nil {
				logger.With(slog_.Error(err)).With("config_file", viper.ConfigFileUsed()).Error("load config file")
				return err
			}
			err = viper.MergeConfigMap(v.AllSettings())
			if err != nil {
				logger.With(slog_.Error(err)).With("config_file", viper.ConfigFileUsed()).Error("load config file")
				return err
			}
			logger.With("config_file", viper.ConfigFileUsed()).Info("load config file")
		}
		logger.With("config_files", cfgFiles).Info("load all config files finished")
		return nil
	}

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// cmd.PersistentFlags().String("foo", "", "A help for foo")
	cmd.PersistentFlags().BoolVar(&config.ForceDisableTls,
		"dangerous-force-disable-tls", false,
		"Disable HTTP/2 over TLS (HTTPS) and serve HTTP instead. Never use this in production.")

	return cmd
}

// DefaultConfigPath returns config file's default path
func DefaultConfigPath() string {
	// 	return filepath_.Pathify(fmt.Sprintf("$HOME/.%s.yaml", version.ServiceName))
	return filepath_.Pathify(fmt.Sprintf("./conf/%s.yaml", version.ServiceName))
}
