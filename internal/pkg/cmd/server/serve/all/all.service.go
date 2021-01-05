// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package all

//import (
//	"context"
//	"sync"
//
//	"github.com/searKing/sole/internal/pkg/cmd/server/serve/dns"
//	"github.com/searKing/sole/internal/pkg/cmd/server/serve/web"
//	"github.com/searKing/sole/internal/pkg/net/graceful"
//	"github.com/searKing/sole/internal/pkg/provider"
//)
//
//func service() {
//	_, grpcBackend := web.Setup()
//	sd, err := dns.SetUp()
//	if err != nil {
//		logrus.WithError(err).Fatal("setup dns")
//	}
//	startWeb, shutdownWeb := web.ServeGRPC(grpcBackend)
//	startDns, shutdownDns := dns.Serve(context.Background(), sd)
//
//	var wg sync.WaitGroup
//	err = graceful.Graceful(&wg, []graceful.GracefulFunc{{
//		Name:         "frontend",
//		StartFunc:    startWeb,
//		ShutdownFunc: shutdownWeb,
//	}, {
//		Name:         "dns",
//		StartFunc:    startDns,
//		ShutdownFunc: shutdownDns,
//	}})
//	if err != nil {
//		logrus.WithError(err).Fatal("Could not gracefully run servers")
//	}
//	wg.Wait()
//}
