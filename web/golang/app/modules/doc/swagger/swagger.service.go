// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package swagger

import (
	"github.com/searKing/sole/internal/pkg/provider"
	"github.com/searKing/sole/web/golang/app/configs/values"
)

type IndexTemplateInfo struct {
	Name           string
	BaseUrl        string
	SwaggerJsonUrl string
}

func GetIndexTemplateInfo(webPath string) IndexTemplateInfo {
	return IndexTemplateInfo{
		Name:           provider.GlobalProvider().Proto().GetService().GetDisplayName(),
		BaseUrl:        webPath,
		SwaggerJsonUrl: values.SwaggerJson,
	}
}
