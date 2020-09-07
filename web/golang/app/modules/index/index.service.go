// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package index

import (
	"path"

	"github.com/searKing/sole/internal/pkg/provider"
)

type IndexTemplateInfo struct {
	Name        string
	Version     string
	Description string
	BaseUrl     string
}

func GetIndexTemplateInfo(prefix string, filename string) IndexTemplateInfo {
	return IndexTemplateInfo{
		Name:        provider.GlobalProvider().Proto().GetService().GetDisplayName(),
		Version:     provider.GlobalProvider().Proto().GetAppInfo().GetBuildVersion(),
		Description: provider.GlobalProvider().Proto().GetService().GetDescription(),
		BaseUrl:     path.Join(prefix, filename),
	}
}
