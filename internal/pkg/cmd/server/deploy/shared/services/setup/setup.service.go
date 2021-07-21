// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package setup

import (
	"fmt"

	"github.com/common-nighthawk/go-figure"
	"github.com/kardianos/service"
	"github.com/searKing/sole/internal/pkg/banner"
	"github.com/searKing/sole/pkg/appinfo"
	"github.com/sirupsen/logrus"
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
	fmt.Println(banner.Banner(appinfo.ServiceName, appinfo.GetVersion().String()))
	switch action {
	case "install", "stop":
		figure.NewFigure(appinfo.ServiceDisplayName, "", false).Print()
	}

	svcConfig := &service.Config{
		Name:        appinfo.ServiceName,
		DisplayName: appinfo.ServiceDisplayName,
		Description: appinfo.ServiceDescription,
		Arguments:   []string{"serve", "all"},
	}
	s, err := service.New(&program{}, svcConfig)
	if err != nil {
		logrus.Fatalln(err)
		return
	}
	logrus.Infoln(svcConfig.Name, action, "...")
	if err := service.Control(s, action); err != nil {
		logrus.Fatalln(err)
		return
	}
	logrus.Infoln(svcConfig.Name, action, "ok")
	return
}
