// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package index

import (
	_ "embed"
	"path"

	"github.com/searKing/golang/go/version"
)

type TemplateInfo struct {
	Name        string
	Version     string
	Description string
	BaseUrl     string
}

func GetTemplateInfo(prefix string, filename string) TemplateInfo {
	return TemplateInfo{
		Name:        version.ServiceDisplayName,
		Version:     version.Get().String(),
		Description: version.ServiceDescription,
		BaseUrl:     path.Join(prefix, filename),
	}
}
