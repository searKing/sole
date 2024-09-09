#! /bin/bash
#
# Copyright 2023 The searKing Author. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.
#

set -euo pipefail

cur_path=$(cd "$(dirname "$0")"; pwd)
cd "${cur_path}"

echo "$0" "$*"

git_root="$(git rev-parse --show-toplevel)"

while getopts p: option; do
  case $option in
  ?)  echo '-h'
      echo 'build soletemplate:'
      echo "cd \"${git_root}\""
      echo './build.sh -t soletemplate -p'
      exit 1 ;;
  esac
done