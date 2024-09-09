// Copyright 2023 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package options

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"go.opentelemetry.io/otel"

	slog_ "github.com/searKing/golang/go/log/slog"
	os_ "github.com/searKing/golang/go/os"
	time_ "github.com/searKing/golang/go/time"

	v1 "github.com/searKing/sole/api/protobuf-spec/soletemplate/v1"
	"github.com/searKing/sole/pkg/domain/logging"
)

type _fileCleaner struct {
}

func NewFileCleaner(ctx context.Context, config *v1.Configuration) (_ *_fileCleaner, err error) {
	spanName := "NewFileCleaner"
	ctx, span := otel.Tracer("").Start(ctx, spanName)
	defer span.End()
	logger := slog.With(logging.SpanAttrs(span)...)
	logger.Debug("loading plugin")
	defer func() {
		if err != nil {
			logger.With(slog_.Error(err)).Error("load plugin failed")
			return
		}
		logger.Info("load plugin successfully")
	}()

	cleaners := config.GetCategory().GetFileCleaners()
	for _, cleaner := range cleaners {
		var interval = max(cleaner.GetCleanInterval().AsDuration(), time.Second)
		logger := slog.With("clean_interval", interval).
			With("file_pattern", cleaner.GetFilePattern()).
			With("max_age", cleaner.GetMaxAge().AsDuration()).
			With("min_age", cleaner.GetMinAge().AsDuration()).
			With("max_count", cleaner.GetMaxCount()).
			With("max_used", cleaner.GetMaxUsedPercent()).
			With("max_iused", cleaner.GetMaxIusedPercent())
		quota := os_.DiskQuota{
			MaxAge:             cleaner.GetMaxAge().AsDuration(),
			MaxCount:           int(cleaner.GetMaxCount()),
			MaxUsedProportion:  cleaner.GetMaxUsedPercent(),
			MaxIUsedProportion: cleaner.GetMaxIusedPercent(),
		}
		pattern := cleaner.GetFilePattern()
		minAge := cleaner.GetMinAge().AsDuration()
		go time_.Until(ctx, func(ctx context.Context) {
			err := os_.UnlinkOldestFilesFunc(pattern, quota, func(name string) (clean bool) {
				total, free, avail, inodes, inodesFree, err := os_.DiskUsage(name)
				if err != nil {
					logger.With(slog_.Error(err)).With("file", name).Error("clean file")
					return true
				}
				defer func() {
					if clean {
						logger.With("file", name).With("clean", clean).
							With("df", fmt.Sprintf("total :%d B, free: %d B, avail: %d B, inodes: %d, inodesFree: %d",
								total, free, avail, inodes, inodesFree)).Info("clean file")
					} else {
						logger.With("file", name).With("clean", clean).
							With("df", fmt.Sprintf("total :%d B, free: %d B, avail: %d B, inodes: %d, inodesFree: %d",
								total, free, avail, inodes, inodesFree)).Debug("clean file")
					}
				}()

				if minAge > 0 {
					fi, err := os.Stat(name)
					if err != nil {
						return true
					}
					age := time.Now().Sub(fi.ModTime())

					if age < minAge {
						logger.With("file", name).With("age", age).
							With("df", fmt.Sprintf("total :%d B, free: %d B, avail: %d B, inodes: %d, inodesFree: %d",
								total, free, avail, inodes, inodesFree)).
							Warn("clean file ignored, in min age, please check your available disk space")
						return false
					}
				}

				return true
			})
			if err != nil {
				logger.With(slog_.Error(err)).Error("clean files failed")
				return
			}
			logger.Debug("clean files successfully")
		}, interval)
	}

	return &_fileCleaner{}, nil
}
