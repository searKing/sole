// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package options

import (
	"context"
	"fmt"
	"os"
	"time"

	math_ "github.com/searKing/golang/go/exp/math"
	os_ "github.com/searKing/golang/go/os"
	time_ "github.com/searKing/golang/go/time"
	configpb "github.com/searKing/sole/api/protobuf-spec/v1/config"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
)

type _fileCleaner struct {
}

func NewFileCleaner(ctx context.Context, config *configpb.Configuration) (_ *_fileCleaner, err error) {
	spanName := "NewFileCleaner"
	ctx, span := otel.Tracer("").Start(ctx, spanName)
	defer span.End()
	logger := logrus.WithField("trace_id", span.SpanContext().TraceID()).
		WithField("span_id", span.SpanContext().SpanID())
	defer func() {
		if err != nil {
			logger.WithError(err).Error("load plugin failed")
			return
		}
		logger.Info("load plugin successfully")
	}()

	cleaners := config.GetFileCleaners()
	for _, cleaner := range cleaners {
		var interval = math_.Max(cleaner.GetCleanInterval().AsDuration(), time.Second)
		logger := logrus.WithField("clean_interval", interval).
			WithField("file_pattern", cleaner.GetFilePattern()).
			WithField("max_age", cleaner.GetMaxAge().AsDuration()).
			WithField("min_age", cleaner.GetMinAge().AsDuration()).
			WithField("max_count", cleaner.GetMaxCount()).
			WithField("max_used", cleaner.GetMaxUsedPercent()).
			WithField("max_iused", cleaner.GetMaxIusedPercent())
		quota := os_.DiskQuota{
			MaxAge:             cleaner.GetMaxAge().AsDuration(),
			MaxCount:           int(cleaner.GetMaxCount()),
			MaxUsedProportion:  cleaner.GetMaxUsedPercent(),
			MaxIUsedProportion: cleaner.GetMaxIusedPercent(),
		}
		minAge := cleaner.GetMinAge().AsDuration()
		go time_.Until(ctx, func(ctx context.Context) {
			err := os_.UnlinkOldestFilesFunc(cleaner.GetFilePattern(), quota, func(name string) (clean bool) {
				total, free, avail, inodes, inodesFree, err := os_.DiskUsage(name)
				if err != nil {
					logger.WithError(err).WithField("file", name).Errorf("clean file")
					return true
				}
				defer func() {
					if clean {
						logger.WithField("file", name).
							WithField("df", fmt.Sprintf("total :%d B, free: %d B, avail: %d B, inodes: %d, inodesFree: %d",
								total, free, avail, inodes, inodesFree)).Infof("clean file")
					}
				}()

				if minAge > 0 {
					fi, err := os.Stat(name)
					if err != nil {
						return true
					}
					age := time.Now().Sub(fi.ModTime())

					if age < minAge {
						logger.WithField("file", name).WithField("age", age).
							WithField("df", fmt.Sprintf("total :%d B, free: %d B, avail: %d B, inodes: %d, inodesFree: %d",
								total, free, avail, inodes, inodesFree)).
							Warnf("clean file ignored, in min age, please check your avaliable disk space")
						return false
					}
				}

				return true
			})
			if err != nil {
				logger.WithError(err).Error("clean files failed")
				return
			}
			logger.Info("clean files successfully")
		}, interval)
	}

	return &_fileCleaner{}, nil
}
