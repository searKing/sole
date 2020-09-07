// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package webapp

import (
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/searKing/sole/web/golang/app/configs/values"
)

const (
	RelativeFileStoragePathPrefix = "web/webapp"
)

var (
	// nodes below are not part of the publicdocument tree of the application.
	// No file contained in the WEB-INF|META-INF directory maybe served directly to a client by the container
	// Also, any requests from the client to access the resources in WEB-INF/ or META-INF/ directory must be
	// returned with a SC_NOT_FOUND(404) response.
	// see http://download.oracle.com/otn-pub/jcp/servlet-2.4-fr-spec-oth-JSpec/servlet-2_4-fr-spec.pdf
	ExcludedPathPrefixes = []string{"/WEB-INF/", "/META-INF/"}
)

func WebAppRouter(router gin.IRouter) gin.IRouter {
	router.StaticFS(values.WebApp, Dir(RelativeFileStoragePathPrefix))
	return router
}

func ResolveWeb(filePath string) (webPath string) {
	return path.Join(values.WebApp, strings.TrimPrefix(filePath, RelativeFileStoragePathPrefix))
}

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
