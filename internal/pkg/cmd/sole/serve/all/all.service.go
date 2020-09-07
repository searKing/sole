// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package all

import (
	"sync"

	"github.com/searKing/sole/internal/pkg/cmd/sole/serve/web"
	"github.com/searKing/sole/internal/pkg/net/graceful"
	"github.com/searKing/sole/internal/pkg/provider"
)

func service() {
	_, grpcBackend := web.Setup()

	startWeb, shutdownWeb := web.ServeGRPC(grpcBackend)
	//startServiceDiscovery, shutdownServiceDiscovery := service_discovery.Serve()

	var wg sync.WaitGroup
	err := graceful.Graceful(&wg, []graceful.GracefulFunc{{
		Name:         "frontend",
		StartFunc:    startWeb,
		ShutdownFunc: shutdownWeb,
		//}, {
		//	Name:         "service_discovery",
		//	StartFunc:    startServiceDiscovery,
		//	ShutdownFunc: shutdownServiceDiscovery,
	}})
	if err != nil {
		provider.GlobalProvider().Logger().WithError(err).Fatal("Could not gracefully run servers")
	}
	wg.Wait()
}
