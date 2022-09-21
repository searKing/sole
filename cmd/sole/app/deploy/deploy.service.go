// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package deploy

import (
	"fmt"

	"github.com/common-nighthawk/go-figure"
	"github.com/kardianos/service"
	"github.com/searKing/golang/go/version"
	"github.com/searKing/sole/internal/pkg/banner"
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

func RunService(action string) error {
	fmt.Println(banner.Banner(version.ServiceName, version.Get().String()))
	switch action {
	case "install", "stop":
		figure.NewFigure(version.ServiceDisplayName, "", false).Print()
	}
	logger := logrus.WithField("service_action", action).WithField("service_name", version.ServiceName)

	svcConfig := &service.Config{
		Name:        version.ServiceName,
		DisplayName: version.ServiceDisplayName,
		Description: version.ServiceDescription,
		Arguments:   []string{"serve", "all"},
	}
	s, err := service.New(&program{}, svcConfig)
	if err != nil {
		logger.WithError(err).Error("creates service failed")
		return err
	}
	logger.Infof("service is ready to control...")
	if err := service.Control(s, action); err != nil {
		logger.WithError(err).Infof("service controled failed...")
		return err
	}
	logger.Infof("service controled successfully...")
	return nil
}
