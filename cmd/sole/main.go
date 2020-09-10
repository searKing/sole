// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/searKing/sole/internal/pkg/cmd/sole"
	"github.com/searKing/sole/internal/pkg/provider/viper"
	"github.com/searKing/sole/pkg/runtime/profile"
)

func main() {
	log.SetPrefix(fmt.Sprintf("[%s] ", viper.ServiceName))
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	rand.Seed(time.Now().UnixNano())

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // cancel when we are finished consuming integers
	rootCmd := sole.NewCommand()
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
