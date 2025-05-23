#!/usr/bin/env bash
#
# Copyright 2020 The searKing Author. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.
#
function die() {
  echo 1>&2 "$*"
  exit 1
}

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

# 打印脚本入参
echo "$0" "$*"
GIT_QUIET=
debug=
g_proto_headers=
g_with_cpp=
g_with_python=
g_with_go=
g_with_go_tag=
g_with_govalidators=
g_with_go_grpc=
g_with_go_grpc_gateway=
g_with_openapiv2=

# //go:generate bash scripts/proto-gen.sh ./v1 --with_go -I "."

g_protos_dir="$1"
shift
while test $# -gt 0; do
  opt="$1"
  shift

  case "$opt" in
  -q)
    GIT_QUIET=1
    ;;
  -d)
    debug=1
    ;;
  -I)
    g_proto_headers+=" -I ${1}"
    shift
    ;;
  --with_cpp)
    g_with_cpp=1
    ;;
  --with_python)
    g_with_python=1
    ;;
  --with_go)
    g_with_go=1
    ;;
  --with_go_tag)
    g_with_go_tag=1
    ;;
  --with_govalidators)
    g_with_govalidators=1
    ;;
  --with_go_grpc)
    g_with_go_grpc=1
    ;;
  --with_go_grpc_gateway)
    g_with_go_grpc_gateway=1
    ;;
  --with_go_grpc_openapiv2)
    g_with_openapiv2=1
    ;;
  --)
    break
    ;;
  *)
    die "Unexpected option: $opt"
    ;;
  esac
done

g_proto_headers="${g_proto_headers} -I ."

# Directory and file names that begin with "." or "_" are ignored
# by the go tool, as are directories named "testdata".
# See: https://github.com/golang/go/issues/30058#issuecomment-459888562
# Also See: https://github.com/golang/go/wiki/Modules
if [ -d "${THIS_BASH_FILE_ABS_DIR}/../../../third_party/" ]; then
  g_proto_headers="${g_proto_headers} -I ${THIS_BASH_FILE_ABS_DIR}/../../../third_party/"
fi
if [ -d "${THIS_BASH_FILE_ABS_DIR}/../../../.third_party/" ]; then
  g_proto_headers="${g_proto_headers} -I ${THIS_BASH_FILE_ABS_DIR}/../../../.third_party/"
fi
if [ -d "${THIS_BASH_FILE_ABS_DIR}/../../../_third_party/" ]; then
  g_proto_headers="${g_proto_headers} -I ${THIS_BASH_FILE_ABS_DIR}/../../../_third_party/"
fi
if [ -n "${g_with_go_grpc_gateway}" ] || [ -n "${g_with_openapiv2}" ]; then
  if [ -d "${THIS_BASH_FILE_ABS_DIR}/../../../third_party/github.com/grpc-ecosystem/grpc-gateway" ]; then
    g_proto_headers="${g_proto_headers} -I ${THIS_BASH_FILE_ABS_DIR}/../../../third_party/github.com/grpc-ecosystem/grpc-gateway"
  fi
  if [ -d "${THIS_BASH_FILE_ABS_DIR}/../../../.third_party/github.com/grpc-ecosystem/grpc-gateway" ]; then
    g_proto_headers="${g_proto_headers} -I ${THIS_BASH_FILE_ABS_DIR}/../../../.third_party/github.com/grpc-ecosystem/grpc-gateway"
  fi
  if [ -d "${THIS_BASH_FILE_ABS_DIR}/../../../_third_party/" ]; then
    g_proto_headers="${g_proto_headers} -I ${THIS_BASH_FILE_ABS_DIR}/../../../_third_party/github.com/grpc-ecosystem/grpc-gateway"
  fi
fi

if [ -n "$g_with_go_tag" ]; then
  g_with_go=
fi

# Sanity check that the right tools are accessible.
for tool in protoc protoc-gen-go protoc-gen-go-tag protoc-gen-govalidators protoc-gen-go-grpc protoc-gen-grpc-gateway protoc-gen-openapiv2; do
  case $tool in
  "protoc-gen-go")
    if [ -z "${g_with_go}" ]; then
      continue
    fi
    ;;
  "protoc-gen-go-tag")
    if [ -z "${g_with_go_tag}" ]; then
      continue
    fi
    ;;
  "protoc-gen-govalidators")
    if [ -z "${g_with_govalidators}" ]; then
      continue
    fi
    ;;
  "protoc-gen-go-grpc")
    if [ -z "${g_with_go_grpc}" ]; then
      continue
    fi
    ;;
  "protoc-gen-grpc-gateway")
    if [ -z "${g_with_go_grpc_gateway}" ]; then
      continue
    fi
    ;;
  "protoc-gen-openapiv2")
    if [ -z "${g_with_openapiv2}" ]; then
      continue
    fi
    ;;
  esac

  # https://grpc.io/docs/protoc-installation/
  # http://google.github.io/proto-lens/installing-protoc.html
  # https://github.com/grpc/grpc-go
  # https://www.grpc.io/docs/languages/go/quickstart/
  # protoc-gen-go: go install github.com/golang/protobuf/protoc-gen-go@latest
  # protoc-gen-grpc-gateway: go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
  # protoc-gen-openapiv2: go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
  q=$(command -v $tool) || die "didn't find $tool
  protoc: dnf install protobuf protobuf-compiler.
  protoc-gen-go: go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
  protoc-gen-go-grpc: go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
  protoc-gen-grpc-gateway: go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
  protoc-gen-openapiv2: go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
  protoc-gen-go-tag: go install github.com/searKing/golang/tools/protoc-gen-go-tag@latest
  protoc-gen-govalidators: go install github.com/mwitkow/go-proto-validators/protoc-gen-govalidators@latest
  "
  echo 1>&2 "$tool: $q"
done

find "${g_protos_dir}" -name "*.proto" -print0 | while read -r -d $'\0' proto_file; do
  proto_base_name="$(basename "${proto_file}" .proto)"
  proto_dir="$(dirname "${proto_file}")"

  cpp_option=""
  python_option=""
  go_option=""
  go_grpc_option=""
  grpc_gateway_option=""
  openapiv2_option=""
  go_tag_option=""
  go_validators_option=""

  if [ -n "${g_with_cpp}" ]; then
    cpp_option="--cpp_out=."
  fi
  if [ -n "${g_with_python}" ]; then
    python_option="--python_out=."
  fi
  #  go_option="--go_out=plugins=grpc,paths=source_relative:."
  if [ -n "${g_with_go}" ]; then
    go_option="--go_out=paths=source_relative:."
  fi
  if [ -n "${g_with_go_tag}" ]; then
    go_tag_option="--go-tag_out=paths=source_relative:."
  fi
  if [ -n "${g_with_govalidators}" ]; then
    go_validators_option="--govalidators_out=paths=source_relative:."
  fi

  if [ -n "${g_with_go_grpc}" ]; then
    go_grpc_option="--go-grpc_out=paths=source_relative:."
  fi
  if [ -n "${g_with_go_grpc_gateway}" ]; then
    grpc_gateway_option="--grpc-gateway_out=logtostderr=true"
  fi
  if [ -n "${g_with_openapiv2}" ]; then
    openapiv2_option="--openapiv2_out=logtostderr=true"
  fi

  api_conf_yaml="${proto_base_name}.yaml"
  if [ -n "${proto_dir}" ]; then
    api_conf_yaml="${proto_dir}/${api_conf_yaml}"
  fi
  if [[ -f "${api_conf_yaml}" ]]; then
    if [ -n "${g_with_go_grpc_gateway}" ]; then
      grpc_gateway_option="${grpc_gateway_option},grpc_api_configuration=${api_conf_yaml}"
    fi
    if [ -n "${g_with_openapiv2}" ]; then
      openapiv2_option="${openapiv2_option},grpc_api_configuration=${api_conf_yaml}"
    fi
  fi
  if [ -n "${g_with_go_grpc_gateway}" ]; then
    grpc_gateway_option="${grpc_gateway_option},paths=source_relative:. --grpc-gateway_opt=allow_delete_body=true"
  fi
  if [ -n "${g_with_openapiv2}" ]; then
    openapiv2_option="${openapiv2_option}:."
  fi

  service_openapiv2_json="${proto_base_name}.openapiv2.json"
  [[ -f "${service_openapiv2_json}" ]] && rm -f "${service_openapiv2_json}"

  printf "\r\033[K%s compiling " "${proto_file}"
  #  protoc -I . ${g_proto_headers} --go-grpc_out=paths=source_relative:. "${grpc_gateway_option}" "${openapiv2_option}" "${go_tag_option}" *.proto || exit
  protoc -I . ${g_proto_headers} ${cpp_option} ${python_option} ${go_option} ${go_tag_option} ${go_validators_option} ${go_grpc_option} ${grpc_gateway_option} ${openapiv2_option} "${proto_file}" || exit
  printf "\r\033[K%s compilied \n" "${proto_file}"
done
printf "\r\033[Kproto-gen done...\n"

#protoc -I . -I .. --go_out=. *.proto

#protoc -I . -I .. -I ../github.com/googleapis/googleapis/ --go_out=plugins=grpc:. --grpc-gateway_out=logtostderr=true:. *.proto

# 编译google api，新版编译器可以省略M参数
#protoc -I . --go_out=plugins=grpc,Mgoogle/protobuf/descriptor.proto=github.com/golang/protobuf/protoc-gen-go/descriptor:. google/api/*.proto

# 编译hello_http.proto
#protoc -I . --go_out=plugins=grpc,Mgoogle/api/annotations.proto=github.com/searKing/go-grpc-example/proto/google/api:. *.proto
