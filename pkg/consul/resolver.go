// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package consul

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/searKing/golang/go/container/hashring"
	"github.com/searKing/golang/go/errors"
	"github.com/searKing/golang/go/sync/atomic"
	time_ "github.com/searKing/golang/go/time"
	"github.com/sirupsen/logrus"
)

//go:generate go-syncmap -type "serviceQueryMap<string, ServiceQuery>"
type serviceQueryMap sync.Map

type ResolverType int

const (
	ResolverTypeReserve ResolverType = iota
	ResolverTypeRandom  ResolverType = iota
	ResolverTypeConsist ResolverType = iota
)

type ServiceQuery struct {
	Name string `validate:"required"`

	// optional
	ResolverType       ResolverType
	Tags               []string
	PassingOnly        bool
	QueryOptions       *api.QueryOptions
	NodeLocatorOptions []hashring.NodeLocatorOption

	// update from consul server
	updateAt  time.Time
	services  []*api.ServiceEntry
	nodeAddrs *hashring.StringNodeLocator // for Consistent by hash ring
}

func (r *ServiceQuery) SetDefault() *ServiceQuery {
	r.ResolverType = ResolverTypeConsist
	r.PassingOnly = true
	return r
}
func (r *ServiceQuery) Complete() {}

type ServiceResolver struct {
	ConsulAddress    string
	ResolverInterval time.Duration
	inShutdown       atomic.Bool

	serviceByName serviceQueryMap

	mu     sync.Mutex
	cancel func()
}

func NewServiceResolver(address string, services ...ServiceQuery) *ServiceResolver {
	c := &ServiceResolver{
		ConsulAddress:    address,
		ResolverInterval: 10 * time.Second,
	}
	for _, s := range services {
		_ = c.AddService(s)
	}
	return c
}

func (srv *ServiceResolver) logger() logrus.FieldLogger {
	return logrus.
		WithField("module", "ServiceResolver").
		WithField("consul", srv.ConsulAddress)
}

// Run will initialize the backend. It must not block, but may run go routines in the background.
func (srv *ServiceResolver) Run(ctx context.Context) error {
	logger := srv.logger()
	logger.Infoln("ServiceResolver Run")
	if srv.inShutdown.Load() {
		logger.Infoln("ServiceResolver Shutdown")
		return fmt.Errorf("server closed")
	}
	go func() {
		errors.HandleError(srv.Serve(ctx))
	}()
	return nil
}

func (srv *ServiceResolver) Serve(ctx context.Context) error {
	logger := srv.logger()
	logger.Infoln("ServiceResolver Serve")

	if srv.inShutdown.Load() {
		err := fmt.Errorf("server closed")
		logger.WithError(err).Errorf("ServiceResolver Serve canceled")
		return err
	}

	defer srv.inShutdown.Store(true)
	ctx, cancel := context.WithCancel(ctx)
	srv.mu.Lock()
	srv.cancel = cancel
	srv.mu.Unlock()

	time_.Until(ctx, func(ctx context.Context) {
		logger.Infof("querying services to consul")
		err := srv.QueryServices()
		if err != nil {
			logger.WithError(err).Errorf("query services failed")
			return
		}
		logger.Infof("query services by consul")
	}, srv.ResolverInterval)
	logger.Info("stopped query services from consul")
	return nil
}

func (srv *ServiceResolver) Shutdown() {
	srv.inShutdown.Store(true)
	srv.mu.Lock()
	defer srv.mu.Unlock()
	if srv.cancel != nil {
		srv.cancel()
	}
}

func (srv *ServiceResolver) QueryServices() error {
	config := api.DefaultConfig()
	config.Address = srv.ConsulAddress
	client, err := api.NewClient(config)
	if err != nil {
		return err
	}

	var errs []error
	srv.serviceByName.Range(func(name string, service ServiceQuery) bool {
		nodes, _, err := client.Health().ServiceMultipleTags(
			service.Name,
			service.Tags,
			service.PassingOnly,
			service.QueryOptions)
		if err != nil {
			errs = append(errs, fmt.Errorf("query service %s: %w", name, err))
			return true
		}
		service.services = nodes
		service.updateAt = time.Now()
		if service.ResolverType == ResolverTypeConsist {
			var nodeAddrs []string
			for _, node := range nodes {
				nodeAddrs = append(nodeAddrs, node.Node.Address)
			}
			service.nodeAddrs = hashring.NewStringNodeLocator(service.NodeLocatorOptions...)
			service.nodeAddrs.AddNodes(nodeAddrs...)
		}
		srv.serviceByName.Store(name, service)
		return true
	})

	return errors.Multi(errs...)
}

func (srv *ServiceResolver) LookupService(name string) []*api.ServiceEntry {
	service, has := srv.serviceByName.Load(name)
	if !has {
		return nil
	}
	return service.services
}

func (srv *ServiceResolver) PickNode(name string, consistKey string) (addr string, has bool) {
	service, has := srv.serviceByName.Load(name)
	if !has {
		return "", false
	}
	if len(service.services) == 0 {
		return "", false
	}

	if service.ResolverType == ResolverTypeConsist {
		return service.nodeAddrs.Get(consistKey)
	}

	return service.services[rand.Intn(len(service.services))].Node.Address, true
}

func (srv *ServiceResolver) AddService(service ServiceQuery) error {
	_, loaded := srv.serviceByName.LoadOrStore(service.Name, service)
	if loaded {
		return fmt.Errorf("service entry already installed")
	}
	return nil
}

func (srv *ServiceResolver) DeleteService(name string) {
	srv.serviceByName.Delete(name)
}
