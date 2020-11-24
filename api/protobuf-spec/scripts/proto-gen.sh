#!/usr/bin/env bash
#
# Copyright 2020 The searKing Author. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.
#

set -o pipefail
set -o errexit
set -o nounset
# set -o xtrace

# 获取输入参数
THIS_BASE_PARAM="$*"
# 获取当前脚本的相对路径文件名称
THIS_BASH_FILE="${BASH_SOURCE-$0}"
# 获取当前脚本的相对路径目录
THIS_BASH_FILE_REF_DIR=$(dirname "${THIS_BASH_FILE}")
# 获取当前脚本的绝对路径目录
THIS_BASH_FILE_ABS_DIR=$(
  cd "${THIS_BASH_FILE_REF_DIR}" || exit
  pwd
)
# 获取当前脚本的名称
THIS_BASH_FILE_BASE_NAME=$(basename "${THIS_BASH_FILE}")
# 获取当前脚本绝对路径
THIS_BASH_FILE_ABS_PATH="${THIS_BASH_FILE_ABS_DIR}/${THIS_BASH_FILE_BASE_NAME}"
# 备份当前路径
STACK_ABS_DIR=$(pwd)
# 路径隔离
#cd "${THIS_BASH_FILE_ABS_DIR}" 1>/dev/null 2>&1 || exit

# This script rebuilds the generated code for the protocol buffers.
# To run this you will need protoc and goprotobuf installed;
# see https://github.com/golang/protobuf for instructions.
# You also need Go and Git installed.
# masking all out put info!

g_protos_dir="$1"
g_proto_headers="-I ${THIS_BASH_FILE_ABS_DIR}/../../../third_party/"
g_proto_headers="${g_proto_headers} -I ${THIS_BASH_FILE_ABS_DIR}/../../../third_party/github.com/grpc-ecosystem/grpc-gateway"
g_proto_headers="${g_proto_headers} -I ${THIS_BASH_FILE_ABS_DIR}/../../../"

function die() {
  echo 1>&2 "$*"
  exit 1
}

# Sanity check that the right tools are accessible.
for tool in protoc protoc-gen-go-grpc protoc-gen-grpc-gateway protoc-gen-openapiv2 protoc-gen-go-tag; do
  # http://google.github.io/proto-lens/installing-protoc.html
  # https://github.com/grpc/grpc-go
  # https://www.grpc.io/docs/languages/go/quickstart/
  # protoc-gen-go: go get -u github.com/golang/protobuf/protoc-gen-go
  # protoc-gen-grpc-gateway: go get -u github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway
  # protoc-gen-openapiv2: go get -u github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2
  q=$(command -v $tool) || die "didn't find $tool
  protoc: brew install protobuf.
  protoc-gen-go: go get -u google.golang.org/protobuf/cmd/protoc-gen-go
  protoc-gen-go-grpc: go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc
  protoc-gen-grpc-gateway: go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
  protoc-gen-openapiv2: go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-openapiv2
  protoc-gen-go-tag: go get -u github.com/searKing/golang/tools/cmd/protoc-gen-go-tag
  "
  echo 1>&2 "$tool: $q"
done

find "${g_protos_dir}" -name "*.proto" -print0 | while read -r -d $'\0' proto_file; do
  proto_base_name="$(basename "${proto_file}" .proto)"
  proto_dir="$(dirname "${proto_file}")"
  pushd "${proto_dir}" 1>/dev/null 2>&1 || exit

  #  go_option="--go_out=plugins=grpc,paths=source_relative:."
  #  go_option="--go_out=paths=source_relative:."
  go_grpc_option="--go-grpc_out=paths=source_relative:."
  grpc_gateway_option="--grpc-gateway_out=logtostderr=true"
  openapiv2_option="--openapiv2_out=logtostderr=true"
  go_tag_option="--go-tag_out=paths=source_relative:."

  api_conf_yaml="${proto_base_name}.yaml"
  if [[ -f "${api_conf_yaml}" ]]; then
    grpc_gateway_option="${grpc_gateway_option},grpc_api_configuration=${api_conf_yaml}"
    openapiv2_option="${openapiv2_option},grpc_api_configuration=${api_conf_yaml}"
  fi
  grpc_gateway_option="${grpc_gateway_option},paths=source_relative:."
  openapiv2_option="${openapiv2_option}:."

  service_openapiv2_json="${proto_base_name}.openapiv2.json"
  [[ -f "${service_openapiv2_json}" ]] && rm -f "${service_openapiv2_json}"

  printf "\r\033[K%s compiling " "${proto_file}"
  #  protoc -I . ${g_proto_headers} --go-grpc_out=paths=source_relative:. "${grpc_gateway_option}" "${openapiv2_option}" "${go_tag_option}" *.proto || exit
  protoc -I . ${g_proto_headers} "${go_grpc_option}" "${grpc_gateway_option}" "${openapiv2_option}" "${go_tag_option}" *.proto || exit
  printf "\r\033[K%s compilied " "${proto_file}"

  popd 1>/dev/null 2>&1 || exit
done
printf "\r\033[Kproto-gen done...\n"

#protoc -I . -I .. --go_out=. *.proto

#protoc -I . -I .. -I ../github.com/googleapis/googleapis/ --go_out=plugins=grpc:. --grpc-gateway_out=logtostderr=true:. *.proto

# 编译google api，新版编译器可以省略M参数
#protoc -I . --go_out=plugins=grpc,Mgoogle/protobuf/descriptor.proto=github.com/golang/protobuf/protoc-gen-go/descriptor:. google/api/*.proto

# 编译hello_http.proto
#protoc -I . --go_out=plugins=grpc,Mgoogle/api/annotations.proto=github.com/searKing/go-grpc-example/proto/google/api:. *.proto
