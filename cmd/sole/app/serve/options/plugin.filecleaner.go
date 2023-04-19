// Copyright 2022 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package options

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	math_ "github.com/searKing/golang/go/exp/math"
	os_ "github.com/searKing/golang/go/os"
	time_ "github.com/searKing/golang/go/time"
	configpb "github.com/searKing/sole/api/protobuf-spec/v1/config"
	"github.com/sirupsen/logrus"
)

type _fileCleaner struct {
}

func NewFileCleaner(ctx context.Context, config *configpb.Configuration) (_ *_fileCleaner, err error) {
	defer func() {
		if err != nil {
			logrus.WithError(err).Error("load plugin failed")
			return
		}
		logrus.Info("load plugin successfully")
	}()

	cleaners := config.GetFileCleaners()
	for i := range cleaners {
		cleaner := cleaners[i]
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
			err := UnlinkOldestFilesFunc(cleaner.GetFilePattern(), quota, func(name string) (clean bool) {
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

func GlobFunc(pattern string, handler func(name string) bool) (matches []string, err error) {
	matches, err = filepath.Glob(pattern)
	if err != nil {
		logrus.WithError(err).WithField("pattern", pattern).Errorf("filepath.Glob")
		return nil, err
	}

	logrus.WithField("pattern", pattern).WithField("matches", len(matches)).Info("GlobFunc")
	if handler == nil {
		return matches, err
	}

	var a []string
	for _, match := range matches {
		if handler(match) {
			a = append(a, match)
		}
	}
	return a, nil
}

// UnlinkOldestFilesFunc unlink old files satisfying f(c) if need
func UnlinkOldestFilesFunc(pattern string, quora os_.DiskQuota, f func(name string) bool) error {
	if quora.NoLimit() {
		return nil
	}

	now := time.Now()

	// find old files
	var filesNotExpired []string
	_, err := GlobFunc(pattern, func(name string) bool {
		fi, err := os.Stat(name)
		if err != nil {
			logrus.WithError(err).WithField("name", name).Errorf("os.Stat")
			return false
		}

		fl, err := os.Lstat(name)
		if err != nil {
			logrus.WithError(err).WithField("name", name).Errorf("os.Lstat")
			return false
		}
		logrus.WithError(err).WithField("name", name).Infof("os.File")
		if quora.MaxAge <= 0 {
			filesNotExpired = append(filesNotExpired, name)
			return false
		}

		if now.Sub(fi.ModTime()) < quora.MaxAge {
			filesNotExpired = append(filesNotExpired, name)
			return false
		}

		if fl.Mode()&os.ModeSymlink == os.ModeSymlink {
			return false
		}
		return true
	})
	if err != nil {
		return err
	}
	return nil
}
