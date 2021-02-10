// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package all

import (
	"context"
	"sync"

	"github.com/searKing/golang/go/errors"
	"github.com/spf13/cobra"

	"github.com/searKing/sole/internal/pkg/cmd/server/serve/web"
)

func CommandE(ctx context.Context) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		var errs []error
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			errs = append(errs, web.CommandE(ctx)(cmd, args))
		}()
		wg.Wait()
		return errors.Multi(errs...)
	}
}
