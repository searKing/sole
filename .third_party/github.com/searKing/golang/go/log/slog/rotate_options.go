// Code generated by "go-option -type rotate"; DO NOT EDIT.
// Install go-option by "go get install github.com/searKing/golang/tools/go-option"

package slog

import "time"

// A RotateOption sets options.
type RotateOption interface {
	apply(*rotate)
}

// EmptyRotateOption does not alter the configuration. It can be embedded
// in another structure to build custom options.
//
// This API is EXPERIMENTAL.
type EmptyRotateOption struct{}

func (EmptyRotateOption) apply(*rotate) {}

// RotateOptionFunc wraps a function that modifies rotate into an
// implementation of the RotateOption interface.
type RotateOptionFunc func(*rotate)

func (f RotateOptionFunc) apply(do *rotate) {
	f(do)
}

// ApplyOptions call apply() for all options one by one
func (o *rotate) ApplyOptions(options ...RotateOption) *rotate {
	for _, opt := range options {
		if opt == nil {
			continue
		}
		opt.apply(o)
	}
	return o
}

// withRotate sets rotate.
func withRotate(v rotate) RotateOption {
	return RotateOptionFunc(func(o *rotate) {
		*o = v
	})
}

// WithRotateFilePathRotateLayout sets FilePathRotateLayout in rotate.
// Time layout to format rotate file
func WithRotateFilePathRotateLayout(v string) RotateOption {
	return RotateOptionFunc(func(o *rotate) {
		o.FilePathRotateLayout = v
	})
}

// WithRotateFileLinkPath sets FileLinkPath in rotate.
// sets the symbolic link name that gets linked to the current file name being used.
func WithRotateFileLinkPath(v string) RotateOption {
	return RotateOptionFunc(func(o *rotate) {
		o.FileLinkPath = v
	})
}

// WithRotateRotateInterval sets RotateInterval in rotate.
// Rotate files are rotated until RotateInterval expired before being removed
// take effects if only RotateInterval is bigger than 0.
func WithRotateRotateInterval(v time.Duration) RotateOption {
	return RotateOptionFunc(func(o *rotate) {
		o.RotateInterval = v
	})
}

// WithRotateRotateSize sets RotateSize in rotate.
// Rotate files are rotated if they grow bigger then size bytes.
// take effects if only RotateSize is bigger than 0.
func WithRotateRotateSize(v int64) RotateOption {
	return RotateOptionFunc(func(o *rotate) {
		o.RotateSize = v
	})
}

// WithRotateMaxAge sets MaxAge in rotate.
// max age of a log file before it gets purged from the file system.
// Remove rotated logs older than duration. The age is only checked if the file is
// to be rotated.
// take effects if only MaxAge is bigger than 0.
func WithRotateMaxAge(v time.Duration) RotateOption {
	return RotateOptionFunc(func(o *rotate) {
		o.MaxAge = v
	})
}

// WithRotateMaxCount sets MaxCount in rotate.
// Rotate files are rotated MaxCount times before being removed
// take effects if only MaxCount is bigger than 0.
func WithRotateMaxCount(v int) RotateOption {
	return RotateOptionFunc(func(o *rotate) {
		o.MaxCount = v
	})
}

// WithRotateForceNewFileOnStartup sets ForceNewFileOnStartup in rotate.
// Force File Rotate when start up
func WithRotateForceNewFileOnStartup(v bool) RotateOption {
	return RotateOptionFunc(func(o *rotate) {
		o.ForceNewFileOnStartup = v
	})
}