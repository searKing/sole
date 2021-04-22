// Copyright 2021 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package appinfo

import "github.com/searKing/golang/go/version"

var (
	// Version
	// NOTE: The $Format strings are replaced during 'git archive' thanks to the
	// companion .gitattributes file containing 'export-subst' in this same
	// directory.  See also https://git-scm.com/docs/gitattributes
	Version   = "v0.0.0-master+$Format:%h$" // git describe --long --tags --dirty --tags --always
	BuildTime = "1970-01-01T00:00:00Z"      // build date in ISO8601 format, output of $(date -u +'%Y-%m-%dT%H:%M:%SZ')
	GitHash   = "$Format:%H$"               // sha1 from git, output of $(git rev-parse HEAD)

	ServiceName        = "" // 服务名称
	ServiceDisplayName = "" // 服务全称
	ServiceDescription = "" // 服务描述
	ServiceId          = "" // 服务实例ID
)

// GetVersion ...
func GetVersion() version.Version {
	return version.Version{
		RawVersion: Version,
		BuildTime:  BuildTime,
		GitHash:    GitHash,
	}
}
