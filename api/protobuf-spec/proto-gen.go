// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// +build tools

package protobuf_spec

// example to get latest scripts
////go:generate bash -c "mkdir -p ./scripts"
////go:generate bash -c "curl -s -L -o ./scripts/proto-gen.sh https://raw.githubusercontent.com/searKing/sole/master/api/protobuf-spec/scripts/proto-gen.sh"
////go:generate bash -c "chmod a+x ./scripts/proto-gen.sh"
////go:generate bash -c "curl -s -L -o ./scripts/swagger-mix.sh https://raw.githubusercontent.com/searKing/sole/master/api/protobuf-spec/scripts/swagger-mix.sh"
////go:generate bash -c "chmod a+x ./scripts/swagger-mix.sh"
////go:generate bash -c "curl -s -L -o ./scripts/swagger-deploy.sh https://raw.githubusercontent.com/searKing/sole/master/api/protobuf-spec/scripts/swagger-deploy.sh"
////go:generate bash -c "chmod a+x ./scripts/swagger-deploy.sh"

//go:generate bash scripts/proto-gen.sh ./v1 --with_go_tag --with_go_grpc --with_go_grpc_gateway --with_go_grpc_openapiv2
//go:generate bash scripts/swagger-mix.sh ./v1
//go:generate bash scripts/swagger-deploy.sh ./v1 "../../web/webapp/static/swagger"
//go:generate bash scripts/swagger-deploy.sh ./v1 "../../api/openapi-spec"
