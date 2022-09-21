// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"fmt"
	"math/rand"
	"os/signal"
	"time"

	os_ "github.com/searKing/golang/go/os"
	runtime_ "github.com/searKing/golang/go/runtime"
	"github.com/searKing/sole/cmd/sole/app"
	"github.com/sirupsen/logrus"

	_ "go.uber.org/automaxprocs"

	"github.com/searKing/sole/pkg/runtime/profile"
)

func main() {
	defer runtime_.LogPanic.Recover()
	rand.Seed(time.Now().UnixNano())

	ctx, cancel := signal.NotifyContext(context.Background(), os_.ShutdownSignals...)
	defer cancel()

	rootCmd := app.NewCommand(ctx)
	// profile
	{
		defer profile.Profile().Stop()
		rootCmd.SetHelpTemplate(fmt.Sprintf("%s\n%s", rootCmd.HelpTemplate(), profile.HelpMessage()))
	}
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		logrus.WithError(err).Errorf("exited.")
		return
	}
}
