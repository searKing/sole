// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime

import (
	"reflect"

	runtime_ "github.com/searKing/golang/go/runtime"
	"github.com/spf13/viper"
)

// Scheme defines methods for serializing and deserializing API objects, a type
// registry for converting group, version, and kind information to and from Go
// schemas, and mappings between Go schemas of different versions. A scheme is the
// foundation for a versioned API and versioned configuration over time.
//
// In a Scheme, a Type is a particular Go struct, a Version is a point-in-time
// identifier for a particular representation of that Type (typically backwards
// compatible), a Kind is the unique name for that Type within the Version, and a
// Group identifies a set of Versions, Kinds, and Types that evolve over time. An
// Unversioned Type is one that is not yet formally bound to a type and is promised
// to be backwards compatible (effectively a "v1" of a Type that does not expect
// to break in the future).
//
// Schemes are not expected to change at runtime and are only threadsafe after
// registration is complete.
type Scheme struct {
	// defaulterFuncs is an array of interfaces to be called with an object to provide defaulting
	// the provided object must be a pointer.
	defaulterFuncs map[reflect.Type]func(interface{})

	// viperFuncs is an array of interfaces to be called with an object to provide viper loading
	// the provided object must be a pointer.
	viperFuncs map[reflect.Type]func(obj interface{}, v *viper.Viper) error

	// schemeName is the name of this scheme.  If you don't specify a name, the stack of the NewScheme caller will be used.
	// This is useful for error reporting to indicate the origin of the scheme.
	schemeName string
}

// NewScheme creates a new Scheme. This scheme is pluggable by default.
func NewScheme() *Scheme {
	s := &Scheme{
		defaulterFuncs: map[reflect.Type]func(interface{}){},
		schemeName:     runtime_.GetShortCaller(2),
	}
	return s
}

// AddTypeDefaultingFunc registers a function that is passed a pointer to an
// object and can default fields on the object. These functions will be invoked
// when Default() is called. The function will never be called unless the
// defaulted object matches srcType. If this function is invoked twice with the
// same srcType, the fn passed to the later call will be used instead.
func (s *Scheme) AddTypeDefaultingFunc(srcType interface{}, fn func(interface{})) {
	s.defaulterFuncs[reflect.TypeOf(srcType)] = fn
}

// Default sets defaults on the provided Object.
func (s *Scheme) Default(src interface{}) {
	if fn, ok := s.defaulterFuncs[reflect.TypeOf(src)]; ok {
		fn(src)
	}
}

// AddTypeViperLoadingFunc registers a function that is passed a pointer to an
// object and can set fields from viper on the object. These functions will be invoked
// when Viper() is called. The function will never be called unless the
// object matches srcType. If this function is invoked twice with the
// same srcType, the fn passed to the later call will be used instead.
func (s *Scheme) AddTypeViperLoadingFunc(srcType interface{}, fn func(obj interface{}, v *viper.Viper) error) {
	s.viperFuncs[reflect.TypeOf(srcType)] = fn
}

// Viper sets defaults on the provided Object.
func (s *Scheme) Viper(src interface{}, v *viper.Viper) error {
	if fn, ok := s.viperFuncs[reflect.TypeOf(src)]; ok {
		return fn(src, v)
	}
	return nil
}
