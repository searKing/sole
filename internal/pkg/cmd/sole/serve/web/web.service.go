// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package web

import (
	"context"
	"fmt"

	"github.com/common-nighthawk/go-figure"
	"github.com/searKing/golang/go/x/graceful"
	"github.com/searKing/sole/internal/pkg/banner"
	"github.com/searKing/sole/internal/pkg/provider"
)

func service() {
	c := provider.GlobalProvider()
	proto := c.Proto()

	fmt.Println(banner.Banner(proto.GetService().GetName(), proto.GetAppInfo().GetBuildVersion()))
	figure.NewFigure(proto.GetService().GetDisplayName(), "", false).Print()

	_, grpcBackend := Setup()
	start, shutdown := ServeGRPC(grpcBackend)

	err := graceful.Graceful(context.Background(), graceful.Handler{
		Name:         "frontend",
		StartFunc:    start,
		ShutdownFunc: shutdown,
	})
	if err != nil {
		c.Logger().WithError(err).Fatal("Could not gracefully run servers")
	}
}
