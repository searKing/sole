// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"time"

	context_ "github.com/searKing/golang/go/context"
	runtime_ "github.com/searKing/golang/go/runtime"
	"github.com/sirupsen/logrus"

	"github.com/searKing/sole/internal/pkg/cmd/server"
	"github.com/searKing/sole/pkg/logs"
	"github.com/searKing/sole/pkg/runtime/profile"
)

func main() {
	defer runtime_.LogPanic.Recover()
	rand.Seed(time.Now().UnixNano())

	logs.InitLog()
	logrus.WithTime(time.Now()).WithField("cmdline", os.Args).Infof("boosting")
	defer func() {
		logrus.WithTime(time.Now()).WithField("cmdline", os.Args).Infof("exited")
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ctx = context_.WithShutdownSignal(ctx)
	rootCmd := server.NewCommand(ctx)
	// profile
	{
		defer profile.Profile().Stop()
		rootCmd.SetHelpTemplate(fmt.Sprintf(`%s
%s`, rootCmd.HelpTemplate(), profile.HelpMessage()))
	}
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		logrus.WithError(err).Fatalf("exited.")
		return
	}
}
