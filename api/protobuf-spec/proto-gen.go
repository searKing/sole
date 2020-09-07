// Copyright 2020 The searKing Author. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package protobuf_spec

//go:generate bash scripts/proto-gen.sh ./v1
//go:generate bash scripts/swagger-mix.sh ./v1
//go:generate bash scripts/swagger-deploy.sh ./v1 "../../web/webapp/static/swagger"
