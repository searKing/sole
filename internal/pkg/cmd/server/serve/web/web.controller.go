// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package web

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/searKing/sole/internal/pkg/version"
)

func isDSNAllowedOrDie(dsn string) {
	if dsn == "memory" {
		logrus.Fatalf(`When using "sole serve web" the DSN can not be set to "memory".`)
	}
}

// Run runs the specified APIServer.  This should never exit.
func Run(ctx context.Context, completeOptions CompletedServerRunOptions) error {
	// To help debugging, immediately log version
	logrus.Infof("Version: %+v", version.GetVersion())
	//isDSNAllowedOrDie(completeOptions.Provider.Proto().GetDatabase().GetDsn())

	server, err := completeOptions.WebServerOptions.Complete().New("sole")
	if err != nil {
		return err
	}

	prepared, err := server.PrepareRun()
	if err != nil {
		return err
	}

	return prepared.Run(ctx)
}
