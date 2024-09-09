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

if [[ -d "pack" ]]; then
  if [[ -f "scripts/http_probe.sh" ]]; then
    echo "post_build: copy scripts/http_probe.sh to pack"
    cp -r -d scripts/http_probe.sh pack/;
    chmod a+x pack/http_probe.sh
  fi
fi
