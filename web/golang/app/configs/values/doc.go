// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package values

import "github.com/searKing/sole/api/protobuf-spec/v1/doc/swagger"

var (
	SwaggerJson = swagger.Pattern_SwaggerService_Json_0.String() //"/doc/swagger/swagger.json"
	SwaggerYaml = swagger.Pattern_SwaggerService_Yaml_0.String() //"/doc/swagger/swagger.yaml"
	SwaggerUis  = []string{
		swagger.Pattern_SwaggerService_UI_0.String(), //"/doc/swagger/swagger-ui"
		swagger.Pattern_SwaggerService_UI_1.String(), //"/doc/swagger/swagger-ui/index.html"
	}
)
