// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log/slog"
	"math/rand"
	"os"
	"os/signal"
	"time"

	"github.com/searKing/sole/pkg/domain/runtime/profile"
	_ "go.uber.org/automaxprocs"
	"golang.org/x/net/context"

	slog_ "github.com/searKing/golang/go/log/slog"
	os_ "github.com/searKing/golang/go/os"
	runtime_ "github.com/searKing/golang/go/runtime"

	"github.com/searKing/sole/soletemplate/cmd/soletemplate/app"
)

func main() {
	defer runtime_.LogPanic.Recover()
	rand.New(rand.NewSource(time.Now().UnixNano()))

	ctx, cancel := signal.NotifyContext(context.Background(), os_.ShutdownSignals...)
	defer cancel()

	rootCmd := app.NewCommand(ctx)
	// profile
	{
		defer profile.Profile().Stop()
		rootCmd.SetHelpTemplate(fmt.Sprintf("%s", rootCmd.HelpTemplate()))
	}
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		slog.With(slog_.Error(err)).Error("exited.")
		os.Exit(1)
	}
}
