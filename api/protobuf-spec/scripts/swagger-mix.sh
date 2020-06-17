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

g_protos_dir="$1"

die() {
  die_with_status 1 "$@"
}

die_with_status() {
  status=$1
  shift
  printf >&2 '%s\n' "$*"
  exit "$status"
}

# 临时文件
# Install the working tree in a tempdir.
tmpdir=$(mktemp -d -t .build.XXXXXX)
function cleanup() {
  printf "Cleaning up %s..." "${tmpdir}"
  [ -d "${tmpdir}" ] && rm -Rf "${tmpdir}"
  printf "\r\033[KCleaning done."
  printf "\r\033[K"
}
trap cleanup EXIT
[ -d "${tmpdir}" ] && rm -Rf "${tmpdir}"
mkdir -p "${tmpdir}"

# Sanity check that the right tools are accessible.
for tool in swagger; do
  q=$(command -v $tool) || die "didn't find $tool"
  echo 1>&2 "$tool: $q"
done
pushd "${g_protos_dir}" 1>/dev/null 2>&1 || exit

mixed_swagger_json="swagger.json"
mixed_swagger_yaml="swagger.yaml"

printf "%s generating" "${mixed_swagger_json}"
swagger -q generate spec -o "${tmpdir}/${mixed_swagger_json}" "${tmpdir}" || exit
#swagger -q init spec --format="json" && cp swagger.json "${mixed_swagger_json}" || exit
printf "\r\033[K%s generated" "${mixed_swagger_json}"

find "." -name "*.swagger.yaml" -o -name "*.swagger.json" -o -name "*.swagger.init.yaml" -o -name "*.swagger.init.json" | while read -r swagger_file; do
  printf "\r\033[K%s mixing" "${swagger_file}"
  if [[ "${swagger_file}" == *.init.yaml ]]; then
    tmp="${swagger_file}.json"
    swagger -q generate spec -i "${swagger_file}" -o "${tmpdir}/${tmp}" || exit

    swagger -q mixin "${tmpdir}/${tmp}" "${tmpdir}/${mixed_swagger_json}" -o "${tmpdir}/${mixed_swagger_json}" || true
    [[ -f "${tmpdir}/${tmp}" ]] && rm -f "${tmpdir}/${tmp}"

    printf "\r\033[K%s mixed" "${swagger_file}"
    continue
  fi
  swagger -q mixin "${tmpdir}/${mixed_swagger_json}" "${swagger_file}" -o "${tmpdir}/${mixed_swagger_json}" || true
  printf "\r\033[K%s mixed" "${swagger_file}"
done
printf "\r\033[K%s generating" "${mixed_swagger_yaml}"
swagger -q generate spec -i "${tmpdir}/${mixed_swagger_json}" -o "${tmpdir}/${mixed_swagger_yaml}" || exit
printf "\r\033[K%s generated" "${mixed_swagger_yaml}"

mv "${tmpdir}/${mixed_swagger_yaml}" "${mixed_swagger_yaml}"
mv "${tmpdir}/${mixed_swagger_json}" "${mixed_swagger_json}"

popd 1>/dev/null 2>&1 || exit
printf "\r\033[Kswagger mixin done...\n"
