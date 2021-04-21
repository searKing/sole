// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package consul

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/searKing/golang/go/errors"
	net_ "github.com/searKing/golang/go/net"
	runtime_ "github.com/searKing/golang/go/runtime"
	"github.com/searKing/golang/go/sync/atomic"
	time_ "github.com/searKing/golang/go/time"
	"github.com/sirupsen/logrus"
)

//go:generate go-syncmap -type "serviceRegistrationMap<string, ServiceRegistration>"
type serviceRegistrationMap sync.Map

type ServiceRegistration struct {
	Name                string
	Id                  string // default is <Name>-<Ip>-<Port>
	Tags                []string
	Ip                  string
	Port                int
	HealthCheckUrl      string
	HealthCheckInterval time.Duration
	TTL                 time.Duration
}

func (r *ServiceRegistration) SetDefault() *ServiceRegistration {
	r.TTL = 300 * time.Second
	r.HealthCheckInterval = 10 * time.Second
	return r
}

func (r *ServiceRegistration) SetAddr(addr string) error {

	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return err
	}

	nport, err := strconv.Atoi(port)
	if err != nil {
		return fmt.Errorf("malformed service port: %w", err)
	}

	r.Ip = host
	r.Port = nport
	return nil
}

func (r *ServiceRegistration) Complete() {
	if r.Id == "" {
		r.Id = fmt.Sprintf("%s-%s-%d", r.Name, r.Ip, r.Port)
	}
}

func (r *ServiceRegistration) GetCheck() (*api.AgentServiceCheck, error) {
	if r.HealthCheckUrl == "" {
		return nil, nil
	}
	u, err := url.Parse(r.HealthCheckUrl)
	if err != nil {
		return nil, err
	}
	u.Host = net_.HostportOrDefault(u.Host, net.JoinHostPort(r.Ip, fmt.Sprintf("%d", r.Port)))
	if u.Scheme == "" {
		u.Scheme = "http"
	}
	return &api.AgentServiceCheck{
		Interval:                       r.HealthCheckInterval.String(),
		HTTP:                           u.String(),
		DeregisterCriticalServiceAfter: r.TTL.String(),
	}, nil
}

type ServiceRegister struct {
	ConsulAddress    string
	RegisterInterval time.Duration
	serviceById      serviceRegistrationMap

	inShutdown atomic.Bool // true when when server is in shutdown

	mu     sync.Mutex
	cancel func()
}

func NewServiceRegister(consulAddr string, services ...ServiceRegistration) (*ServiceRegister, error) {
	register := &ServiceRegister{
		ConsulAddress:    consulAddr,
		RegisterInterval: 30 * time.Second,
	}
	for _, s := range services {
		_ = register.AddService(s)
	}
	return register, nil
}

// Run will initialize the backend. It must not block, but may run go routines in the background.
func (srv *ServiceRegister) Run(ctx context.Context) error {
	logger := srv.logger()
	logger.Infof("ServiceRegister Run")

	if srv.inShutdown.Load() {
		err := fmt.Errorf("server closed")
		logger.WithError(err).Errorf("ServiceRegister Run canceled")
		return err
	}
	go func() {
		errors.HandleError(srv.Serve(ctx))
	}()
	return nil
}

func (srv *ServiceRegister) Serve(ctx context.Context) error {
	logger := srv.logger()
	logger.Infof("ServiceRegister Serve")

	if srv.inShutdown.Load() {
		err := fmt.Errorf("server closed")
		logger.WithError(err).Errorf("ServiceRegister Serve canceled")
		return err
	}

	defer srv.inShutdown.Store(true)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	srv.mu.Lock()
	srv.cancel = cancel
	srv.mu.Unlock()

	time_.Until(ctx, func(ctx context.Context) {
		logger.Infof("registering services to consul")
		err := srv.Register()
		if err != nil {
			logger.WithError(err).Errorf("register services failed")
			return
		}
		logger.Infof("register services by consul")
	}, srv.RegisterInterval)

	logger.Infoln("unregistering services to consul")
	err := srv.UnRegister()
	if err != nil {
		logger.WithError(err).Errorf("unregisters service failed")
		return err
	}
	logger.Info("unregister services by consul")
	return nil
}

func (srv *ServiceRegister) Shutdown() {
	srv.inShutdown.Store(true)
	srv.mu.Lock()
	defer srv.mu.Unlock()
	if srv.cancel != nil {
		srv.cancel()
	}
}

func (srv *ServiceRegister) Register() error {
	defer runtime_.NeverPanicButLog.Recover()
	var errs []error
	srv.serviceById.Range(func(id string, service ServiceRegistration) bool {
		errs = append(errs, srv.RegisterService(service))
		return true
	})
	return errors.Multi(errs...)
}

func (srv *ServiceRegister) UnRegister() error {
	defer runtime_.NeverPanicButLog.Recover()
	var errs []error
	srv.serviceById.Range(func(id string, service ServiceRegistration) bool {
		errs = append(errs, srv.UnregisterService(service))
		return true
	})
	return errors.Multi(errs...)
}

func (srv *ServiceRegister) logger() logrus.FieldLogger {
	return logrus.
		WithField("module", "service_registry").
		WithField("consul", srv.ConsulAddress)
}

func (srv *ServiceRegister) AddService(service ServiceRegistration) error {
	_, loaded := srv.serviceById.LoadOrStore(service.Id, service)
	if loaded {
		return fmt.Errorf("service entry already installed")
	}
	return nil
}

func (srv *ServiceRegister) DeleteService(serviceId string) error {
	service, loaded := srv.serviceById.LoadAndDelete(serviceId)
	if !loaded {
		return nil
	}
	return srv.UnregisterService(service)
}

func (srv *ServiceRegister) GetConsulAgent() (*api.Agent, error) {
	config := api.DefaultConfig()
	config.Address = srv.ConsulAddress
	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}
	return client.Agent(), nil
}

func (srv *ServiceRegister) RegisterService(service ServiceRegistration) error {
	logger := srv.logger().
		WithField("service_name", service.Name).
		WithField("service_id", service.Id).
		WithField("service_ip", service.Ip).
		WithField("service_port", service.Port).
		WithField("health_check_url", service.HealthCheckUrl)
	agent, err := srv.GetConsulAgent()
	if err != nil {
		logger.WithError(err).Errorf("get consul agent to register service to consul")
		return err
	}

	check, err := service.GetCheck()
	if err != nil {
		logger.WithError(err).Errorf("get check to register service to consul")
		return err
	}

	reg := &api.AgentServiceRegistration{
		ID:      service.Id,
		Name:    service.Name,
		Tags:    service.Tags,
		Port:    service.Port,
		Address: service.Ip,
		Check:   check,
	}

	err = agent.ServiceRegister(reg)
	if err != nil {
		logger.WithError(err).Errorf("register service to consul")
		return err
	}
	logger.Infof("register service to consul")
	return nil
}

func (srv *ServiceRegister) UnregisterService(service ServiceRegistration) error {
	logger := srv.logger().
		WithField("service_name", service.Name).
		WithField("service_id", service.Id).
		WithField("service_ip", service.Ip).
		WithField("service_port", service.Port).
		WithField("health_check_url", service.HealthCheckUrl)
	agent, err := srv.GetConsulAgent()
	if err != nil {
		logger.WithError(err).Errorf("get consul agent to register service to consul")
		return err
	}

	err = agent.ServiceDeregister(service.Id)
	if err != nil {
		logger.WithError(err).Errorf("unregister service to consul")
		return err
	}
	logger.Infof("unregister service to consul")
	return nil

}
