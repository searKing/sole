// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package web

import (
	"context"

	"github.com/searKing/golang/go/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func isDSNAllowedOrDie(dsn string) {
	if dsn == "memory" {
		logrus.Fatalf(`When using "sole serve web" the DSN can not be set to "memory".`)
	}
}

func CommandE(ctx context.Context) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		s := NewServerRunOptions()

		// set default options
		completedOptions, err := s.Complete()
		if err != nil {
			return err
		}

		// validate options
		if err := errors.Multi(completedOptions.Validate()...); err != nil {
			return err
		}
		return completedOptions.Run(ctx)
	}
}
