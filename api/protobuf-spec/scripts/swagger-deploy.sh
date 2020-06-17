#!/usr/bin/env bash
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

g_swagger_dir="$1"
g_deploy_target_dir="$2"

if [[ -z "${g_swagger_dir}" ]]; then
  echo "param 1, that is swagger dir, must not empty string"
  exit
fi

if [[ -z "${g_deploy_target_dir}" ]]; then
  echo "param 2, that is deploy target dir, must not empty string"
  exit
fi

if [ ! -d "${g_deploy_target_dir}" ]; then
  mkdir -p "${g_deploy_target_dir}" || exit
fi

find "${g_swagger_dir}" -maxdepth 1 -type f -name "swagger.json" -o -name "swagger.yaml" | while read -r swagger_file; do
  printf "\r\033[K%s deploying" "${swagger_file}"
  cp "${swagger_file}" "${g_deploy_target_dir}" || exit
  printf "\r\033[K%s deployed" "${swagger_file}"
done

printf "\r\033[Kswagger deploy done...\n"

set +o pipefail
