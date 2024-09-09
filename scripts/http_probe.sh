#! /bin/bash
#
# Copyright 2023 The searKing Author. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.
#

set -euo pipefail

cur_path=$(
  cd "$(dirname "$0")"
  pwd
)
cd "${cur_path}"

# print args
echo "$0" "$*"

target=sole
pid=$(pgrep "${target:0:15}" || true)
while getopts t: option; do
  case $option in
  t) target=$OPTARG
    echo "exact match target name: ${target}"
    pid=$(ps -C "${target}" -o pid |awk '!/PID/ {print $1}' || true)
    ;;
  ?) echo "unknown option -$option=$OPTARG"; exit 1 ;;
  esac
done

echo "liveness probe: ${target}"

#[[ $(pgrep "${target}") -eq $(netstat -tlnp | awk '/tcp */ {split($NF,a,"/"); print a[1]}') ]] || exit 1
if [ -z "$pid" ]; then
  echo "${target} is not running, no pid"
  exit 1
fi

if [ "$(echo "$pid"|wc -l)" -gt 1 ]; then
  echo "${target} is running, more than one pid: $(echo "$pid"|tr '\n' ' ')"
  exit 1
fi

listen_pid_cnt=$(netstat -tlnp | awk '/tcp */ {split($NF,a,"/"); { print a[1] }}'|awk -v "pid=$pid" '$0 == pid {print $0}'|wc -l)

if [ "$listen_pid_cnt" -gt 0 ]; then
  echo "${target}'s pid(${pid}), listen on port: $(lsof -p "${pid}" |grep LISTEN| awk '{print $9}')..."
else
  echo "${target}'s pid(${pid}), no listen port"
  exit 1
fi