// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package setup

import (
	"fmt"

	"github.com/common-nighthawk/go-figure"
	"github.com/kardianos/service"
	"github.com/searKing/sole/internal/pkg/banner"
	"github.com/searKing/sole/internal/pkg/provider"
)

const (
	ServiceActionInstall   = "install"
	ServiceActionUninstall = "uninstall"
	ServiceActionStart     = "start"
	ServiceActionStop      = "stop"
)

type program struct{}

func (p *program) Start(s service.Service) (err error) {
	return
}
func (p *program) Stop(s service.Service) (err error) {
	return
}
func Setup(action string) {
	c := provider.GlobalProvider()
	proto := c.Proto()
	fmt.Println(banner.Banner(proto.GetService().GetName(), proto.GetAppInfo().GetBuildVersion()))
	switch action {
	case "install", "stop":
		figure.NewFigure(proto.GetService().GetDisplayName(), "", false).Print()
	}

	svcConfig := &service.Config{
		Name:        proto.GetService().GetName(),
		DisplayName: proto.GetService().GetDisplayName(),
		Description: proto.GetService().GetDescription(),
		Arguments:   []string{"serve", "all"},
	}
	s, err := service.New(&program{}, svcConfig)
	if err != nil {
		c.Logger().Fatalln(err)
		return
	}
	c.Logger().Infoln(svcConfig.Name, action, "...")
	if err := service.Control(s, action); err != nil {
		c.Logger().Fatalln(err)
		return
	}
	c.Logger().Infoln(svcConfig.Name, action, "ok")
	return
}
