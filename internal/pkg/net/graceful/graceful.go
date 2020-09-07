// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package graceful

import (
	"context"
	"net/http"
	"sync"

	"github.com/ory/graceful"
	"github.com/searKing/sole/internal/pkg/provider"
)

type GracefulFunc struct {
	Name         string
	StartFunc    graceful.StartFunc
	ShutdownFunc graceful.ShutdownFunc
}

func Graceful(wg *sync.WaitGroup, gracefulFuncs []GracefulFunc) error {
	if len(gracefulFuncs) == 0 {
		return nil
	}

	return graceful.Graceful(func() error {
		logger := provider.GlobalProvider().Logger()

		for idx, g := range gracefulFuncs {
			start := g.StartFunc
			name := g.Name
			if start == nil {
				continue
			}
			wg.Add(1)
			go func() {
				defer wg.Done()
				if err := start(); err != nil && err != http.ErrServerClosed {
					logger.WithError(err).Fatalf("Could not gracefully run servers[%d] [%s]", idx, name)
				}
			}()
		}
		return nil
	}, func(ctx context.Context) error {
		logger := provider.GlobalProvider().Logger()

		for idx, g := range gracefulFuncs {
			shutdown := g.ShutdownFunc
			if shutdown == nil {
				continue
			}

			if err := shutdown(ctx); err != nil {
				logger.WithError(err).Fatalf("Could not gracefully shutdown servers[%d] [%s]", idx, g.Name)
			} else {
				logger.Infof("Gracefully shutdown servers[%d] [%s]", idx, g.Name)
			}
		}

		return nil
	})
}
