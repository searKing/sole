#! /bin/bash
#
# Copyright 2022 The searKing Author. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.
#

set -euo pipefail

cur_path=$(cd "$(dirname "$0")"; pwd)
cd "${cur_path}"

echo "$0" "$*"

# ./build.sh -t soletemplate -l -d
target=soletemplate
enable_pprof=OFF
conf_name=""
dockerfile_name=""
go_build_tag=""

while getopts t:ldc:n:s:pb: option; do
  case $option in
  t) target=$OPTARG ;;
  p) enable_pprof=ON ;;
  c) conf_name=$OPTARG ;;
  s) dockerfile_name=$OPTARG ;;
  b) go_build_tag=$OPTARG ;;
  ?) exit 1 ;;
  esac
done

if [ "${enable_pprof}"x = "ON"x ]; then
  echo "go pprof is enabled, never use it in production environment"
fi

command -v go
go version

target_repo_cmd_root="${target}/cmd/${target}"
if [ -f "${target_repo_cmd_root}"/prebuild.sh ]; then
  ln -sf "${target_repo_cmd_root}"/prebuild.sh prebuild.sh
  chmod a+x prebuild.sh
fi
if [ -f "${target_repo_cmd_root}"/postbuild.sh ]; then
  ln -sf "${target_repo_cmd_root}"/postbuild.sh postbuild.sh
  chmod a+x postbuild.sh
fi

#make generate
make -C "${target_repo_cmd_root}" unpack
make -C "${target_repo_cmd_root}" pack ENABLE_PPROF="${enable_pprof}" CONF_NAME="${conf_name}" DOCKERFILE_NAME="${dockerfile_name}" GO_BUILD_TAG="${go_build_tag}"

rm -rf pack/
mv "${target_repo_cmd_root}/pack" pack

echo "pack done"
