// Copyright 2024 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package http

import (
	"errors"
	"net/http"
	"os"
	"strings"
)

var (
	// ExcludedPathPrefixes
	// nodes below are not part of the public document tree of the application.
	// No file contained in the WEB-INF|META-INF directory maybe served directly to a client by the container
	// Also, any requests from the client to access the resources in WEB-INF/ or META-INF/ directory must be
	// returned with a SC_NOT_FOUND(404) response.
	// see http://download.oracle.com/otn-pub/jcp/servlet-2.4-fr-spec-oth-JSpec/servlet-2_4-fr-spec.pdf
	ExcludedPathPrefixes = []string{"/WEB-INF/", "/META-INF/"}
)

type Dir string

// Open implements FileSystem using os.Open, opening files for reading rooted
// and relative to the directory d.
func (d Dir) Open(name string) (http.File, error) {
	for _, prefix := range ExcludedPathPrefixes {
		if strings.HasPrefix(name, prefix) {
			return nil, os.ErrNotExist
		}
	}
	return http.Dir(d).Open(name)
}

type Dirs []string

// Open implements FileSystem using os.Open, opening files for reading rooted
// and relative to the directory d.
func (dirs Dirs) Open(name string) (http.File, error) {
	if len(dirs) == 0 {
		return nil, os.ErrNotExist
	}
	var errs []error
	for _, d := range dirs {
		file, err := Dir(d).Open(name)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		return file, nil
	}
	return nil, errors.Join(errs...)
}
