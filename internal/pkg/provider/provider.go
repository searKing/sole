// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package provider

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync/atomic"
	"syscall"

	"github.com/fsnotify/fsnotify"
	atomic_ "github.com/searKing/golang/go/sync/atomic"
	viper_ "github.com/searKing/sole/api/protobuf-spec/v1/viper"
	"github.com/searKing/sole/internal/pkg/provider/viper"
)

//go:generate go-atomicvalue -type "proto<*github.com/searKing/sole/api/v1/viper.ViperProto>"
type proto atomic.Value

type Provider struct {
	proto                 proto
	configPath            string
	reloadOnConfigChanged bool

	modulesEnabled           atomic_.Bool
	logger                   logger
	tracer                   tracer
	systemSecret             systemSecret
	sqlDB                    sqlDB
	keyManager               keyManager
	tlsConfig                tlsConfig
	keyCipher                keyCipher
	prometheusMetricsManager prometheusMetricsManager

	ctx context.Context
}

// NewProvider returns a provider with default proto
func NewProvider(ctx context.Context, cfgFile string) *Provider {
	provider := &Provider{
		ctx:        ctx,
		configPath: cfgFile,
	}
	provider.watchSignal()
	if err := provider.reload(); err != nil {
		return nil
	}
	return provider
}

func (p *Provider) Context() context.Context {
	if p.ctx == nil {
		return context.Background()
	}
	return p.ctx
}

func (p *Provider) Proto() *viper_.ViperProto {
	return p.proto.Load()
}

// onConfigChange refreshes proto, and is triggered when the config file has been loaded or changed.
// reload params and all modules
func (p *Provider) onConfigChange(in fsnotify.Event) {
	p.Logger().WithField("notify_event", in.String()).Infof("ignore config file changed")
	if !p.reloadOnConfigChanged {
		return
	}
	_ = syscall.Kill(syscall.Getpid(), syscall.SIGHUP)
}

func (p *Provider) watchSignal() {
	if !p.reloadOnConfigChanged {
		return
	}
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGHUP)
	go func() {
		defer signal.Stop(ch)
		for {
			select {
			case <-p.Context().Done():
				return
			case sig, ok := <-ch:
				if !ok {
					return
				}
				err := p.reload()
				if err != nil {
					p.Logger().WithField("signal", sig.String()).WithError(err).Error("reload")
					return
				}
			}
		}
	}()
}

func (p *Provider) TestConfig() error {
	_, err := p.test()
	return err
}

func (p *Provider) test() (*viper_.ViperProto, error) {
	tmp, err := viper.Reload(p.configPath, p.onConfigChange, viper.NewDefaultViperProto())
	if err != nil {
		p.Logger().WithField("config_path", p.configPath).WithError(err).Error("test")
		return nil, err
	}
	return tmp, nil
}

func (p *Provider) reload() error {
	tmp, err := p.test()
	if err != nil {
		p.Logger().WithError(err).Error("reload")
		return err
	}

	p.proto.Store(tmp)
	p.RegisterServeModules()
	return nil
}

func (p *Provider) EnableServeModule() {
	if p.modulesEnabled.CAS(false, true) {
		p.RegisterServeModules()
	}
}

func (p *Provider) EnableMigrateDependencyModules() {
	if p.modulesEnabled.CAS(false, true) {
		p.RegisterMigrateModules()
	}
}

func (p *Provider) RegisterServeModules() {
	if !p.modulesEnabled.Load() {
		return
	}
	p.updateLogger()
	p.updateTracing()
	p.updateSystemSecret()
	p.updateKeyCipher()
	p.updateSqlDB()
	p.updateKeyManager()
	p.updateTLSConfig()
	p.updatePrometheusMetricsManager()
}

func (p *Provider) RegisterMigrateModules() {
	if !p.modulesEnabled.Load() {
		return
	}
	p.updateLogger()
	p.updateSqlDB()
	p.updateKeyManager()
}

func (p *Provider) GetBackendBindHostPort() string {
	local := p.Proto().GetWeb().GetBindAddr()
	return getHostPort(local.GetHost(), local.GetPort())
}

func (p *Provider) GetBackendAdvertiseHostPort() string {
	www := p.Proto().GetWeb().GetAdvertiseAddr()
	if www.GetHost() == "" {
		return p.GetBackendBindHostPort()
	}
	return getHostPort(www.GetHost(), www.GetPort())
}

func getHostPort(hostname string, port int32) string {
	if strings.HasPrefix(hostname, "unix:") {
		return hostname
	}
	return fmt.Sprintf("%s:%d", hostname, port)
}
