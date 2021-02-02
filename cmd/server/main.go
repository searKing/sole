// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/searKing/sole/internal/pkg/cmd/server"
	"github.com/searKing/sole/pkg/logs"
	"github.com/searKing/sole/pkg/runtime/profile"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	logs.InitLog()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	rootCmd := server.NewCommand()
	// profile
	{
		defer profile.Profile().Stop()
		rootCmd.SetHelpTemplate(fmt.Sprintf(`%s
%s`, rootCmd.HelpTemplate(), profile.HelpMessage()))
	}
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		log.Fatal(err)
	}
}
