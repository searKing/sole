#! /bin/bash
#
# Copyright 2022 The searKing Author. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.
#

set -euo pipefail

cur_path=$(cd "$(dirname "$0")";pwd)
cd "${cur_path}"

echo "$0" "$*"

args_target=""
args_clean_only=OFF
while getopts t:ch option; do
  case $option in
  c) args_clean_only=ON ;;
  t) args_target=$OPTARG ;;
  h) echo './touch.sh -t YOUR-NEW-REPO-NAME' ;;
  ?) exit 1 ;;
  esac
done

if [[ -z "$args_target" ]]; then
    echo "target name should not be empty"
    exit 1
fi

#target=$(echo "$args_target"|tr A-Z a-z |tr -d '_-')
target=$(echo "$args_target" | awk 'BEGIN{RS="[_-]";ORS="";} {print tolower($1)}')
repo_name="sole$target"

if [ "${args_clean_only}"x = "ON"x ]; then
  rm -rvf "./api/protobuf-spec/$repo_name"
  rm -rvf "./$repo_name"
  echo "'$repo_name' cleaned. recommend to create with './touch.sh -t $target'"
  exit 0
fi

echo "copy template to ${target}"
find ./api/protobuf-spec/soletemplate -type f -exec bash -c "__new_file__=\${0//template/${target}}; mkdir -p \$(dirname \${__new_file__}); cp -v \$0 \${__new_file__}" {} \;
find ./soletemplate -type f -exec bash -c "__new_file__=\${0//template/${target}}; mkdir -p \$(dirname \${__new_file__}); cp -v \$0 \${__new_file__}" {} \;

args=(-i)
if [[ $OSTYPE == 'darwin'* ]]; then
  args=(-i '""')
fi

echo "replace template to ${target}"
cmd="sed ${args[*]} \"s#template#${target}#g\" {}"
find "./api/protobuf-spec/$repo_name" -type f -not -name "*.DS_Store" -print0 |xargs -0 -I {} bash -c "$cmd"
find "./$repo_name" -type f -not -name "*.DS_Store" -not -name "*.DS_Store" -print0 |xargs -0 -I {} bash -c "$cmd"

target_upper_first_letter=$(echo "$args_target" | awk 'BEGIN{RS="[_-]";ORS="";} {print toupper(substr($0,0,1)) substr($0,2,length($0))}')
echo "replace Template to $target_upper_first_letter"
cmd="sed ${args[*]} \"s#Template#${target_upper_first_letter}#g\" {}"
find "./api/protobuf-spec/$repo_name" -type f -not -name "*.DS_Store" -print0 |xargs -0 -I {} bash -c "$cmd"
find "./$repo_name" -type f -not -name "*.DS_Store" -print0 |xargs -0 -I {} bash -c "$cmd"

echo "replace Help${target_upper_first_letter} to HelpTemplate"
cmd="sed ${args[*]} \"s#Help${target_upper_first_letter}#HelpTemplate#g\" {}"
find "./$repo_name" -type f -not -name "*.DS_Store" -not -name "*.DS_Store" -print0 |xargs -0 -I {} bash -c "$cmd"

echo "'$repo_name' create finished. recommend to build with './build.sh -t $repo_name'"
