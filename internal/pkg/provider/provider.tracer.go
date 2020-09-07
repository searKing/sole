// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package provider

import (
	"sync/atomic"

	"github.com/searKing/sole/api/protobuf-spec/v1/viper"
	"github.com/searKing/sole/pkg/micro/tracing"
	jaegerConfig "github.com/uber/jaeger-client-go/config"
)

//go:generate go-atomicvalue -type "tracer<*github.com/searKing/sole/app/core/tracing.Tracer>"
type tracer atomic.Value

func (p *Provider) Tracer() *tracing.Tracer {
	return p.tracer.Load()
}

func (p *Provider) updateTracing() {
	proto := p.Proto()
	tracingInfo := proto.GetTracing()

	_ = p.tracer.Load().Close()

	if !tracingInfo.GetEnable() {
		p.tracer.Store(nil)
		return
	}

	tracer := tracing.New()
	switch tracingInfo.GetType() {
	case viper.Tracing_urber_jaeger:
		tracer.Type = tracing.TypeJeager
	case viper.Tracing_zipkin:
		tracer.Type = tracing.TypeZipkin
	default:
		p.tracer.Store(nil)
		tracer.Type = tracing.TypeButt
		return
	}
	tracer.ServiceName = proto.GetService().GetName()

	reporter := tracingInfo.GetJaeger().GetReporter()
	if reporter != nil {
		tracer.Reporter = &jaegerConfig.ReporterConfig{}
		tracer.Reporter.LocalAgentHostPort = reporter.GetLocalAgentHostPort()
	}

	sampler := tracingInfo.GetJaeger().GetSampler()
	if sampler != nil {
		tracer.Sampler = &jaegerConfig.SamplerConfig{}
		tracer.Sampler.SamplingServerURL = sampler.GetServerUrl()
		tracer.Sampler.Type = sampler.GetType().String()
		tracer.Sampler.Param = float64(sampler.GetParam())
	}
	p.tracer.Store(tracer)
}
