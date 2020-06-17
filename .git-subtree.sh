#!/usr/bin/env bash
#
# .git-subtree.sh: add, init, update or list git subtrees
#
# Copyright (c) 2020 searKing Chan

# masking all out put info!
set -u
set -o pipefail

error_report() {
  echo "error($2) on line $1: see failed command above"
  exit "$2"
}

exit_report() {
  if [ "$2" -ne 0 ]; then
    echo "exit($2) on line $1: see failed command above"
    exit "$2"
  fi
}
trap 'error_report $LINENO $?' ERR
trap 'exit_report $LINENO $?' EXIT

# Having this variable in your environment would break scripts because
# you would cause "cd" to be taken to unexpected places.  If you
# like CDPATH, define it for your interactive shell sessions without
# exporting it.
# But we protect ourselves from such a user mistake nevertheless.
unset CDPATH

# Similarly for IFS, but some shells (e.g. FreeBSD 7.2) are buggy and
# do not equate an unset IFS with IFS with the default, so here is
# an explicit SP HT LF.
IFS='
'

g_shell_name="$(basename "$0")" #获取当前脚本名称
#切换并获取当前脚本所在路径
g_shell_abs_dir="$(
  cd "$(dirname "$0")" || exit
  pwd
)"

if test $# -eq 0; then
  set -- -h
fi

#dashless=$(basename "$0" | sed -e 's/-/ /')
dashless=$(basename "$0")

USAGE="[--quiet] [--debug]
   or: $dashless [--quiet] [--debug] add --update --prefix <prefix> [--] <url> [<branch>] [<disabled>]
   or: $dashless [--quiet] [--debug] tidy --update --prefix <prefix> [--]
   or: $dashless [--quiet] [--debug] update --install --prefix <prefix> [--]
   or: $dashless [--quiet] [--debug] download --prefix <prefix> [--]
   or: $dashless [--quiet] [--debug] summary --prefix <prefix> [--]"

OPTIONS_SPEC="
${dashless} add       --prefix=<prefix> <url> <branch>
${dashless} add       --prefix=<prefix> <url> <branch> <disabled>
${dashless} tidy      --prefix=<prefix>
${dashless} update    --prefix=<prefix>
${dashless} download  --prefix=<prefix>
${dashless} summary   --prefix=<prefix>
--
h,help        show the help
q,quiet       quiet
d,debug       show debug messages
P,prefix=     the name of the subdir to split out
m,message=    use the given message as the commit message for the merge commit
 options for 'add', 'tidy'
u,update     update tracked subtrees
 options for 'update'
i,intall     install missing subtrees
squash        merge subtree changes as a single commit
"
#eval "$(echo "$OPTIONS_SPEC" | git rev-parse --parseopt -- "$@" || echo exit $?)"

OPTIONS_KEEPDASHDASH=
OPTIONS_STUCKLONG=
GITTREES_FILE=".gittrees"

GIT_QUIET=

GIT_TEXTDOMAINDIR=
GIT_TEST_GETTEXT_POISON=
#LONG_USAGE=
#SUBDIRECTORY_OK=Yes
NONGIT_OK=Yes

PATH=$PATH:$(git --exec-path)
## shellcheck disable=SC1090
. git-sh-setup
## shellcheck disable=SC1090
. git-parse-remote
require_work_tree
wt_prefix=$(git rev-parse --show-prefix)
cd_to_toplevel

# Tell the rest of git that any URLs we get don't come
# directly from the user, so it can apply policy as appropriate.
GIT_PROTOCOL_FROM_USER=0
export GIT_PROTOCOL_FROM_USER

debug=
command=
message=
prefix=
update=
install=

debug() {
  if test -n "$debug"; then
    printf "%s\n" "$*" >&2
  fi
}

say() {
  if test -z "$GIT_QUIET"; then
    printf "%s\n" "$*" >&2
  fi
}

progress() {
  if test -z "$GIT_QUIET"; then
    printf "%s\r" "$*" >&2
  fi
}

assert() {
  if ! "$@"; then
    die "assertion failed: " "$@"
  fi
}

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
  -P)
    prefix+=$'\n'
    prefix+="${1%/}"
    shift
    ;;
  -m)
    message="$1"
    shift
    ;;
  -u | --update)
    update=1
    ;;
  -i | --install)
    install=1
    ;;
  --)
    break
    ;;
  *)
    die "Unexpected option: $opt"
    ;;
  esac
done

command="$1"
shift

case "$command" in
add | tidy | update | summary) ;;

*)
  die "Unknown command '$command'"
  ;;
esac

#
# Print a subtree configuration setting
#
# $1 = subtree prefix
# $2 = option
# $3 = default value
#
# Checks in the usual git-config places first (for overrides),
# otherwise it falls back on .gittrees.  This allows you to
# distribute project-wide defaults in .gittrees, while still
# customizing individual repositories if necessary.  If the option is
# not in .gittrees either, print a default value.
#
get_subtree_config() {
  prefix="$1"
  option="$2"
  default="$3"
  value=$(git config subtree."$prefix"."$option")
  if test -z "$value"; then
    value=$(git config -f ${GITTREES_FILE} subtree."$prefix"."$option")
  fi
  printf '%s' "${value:-$default}"
}

#
# Print subtree prefixes merged already
#
get_subtree_prefixes_merged() {
  prefixes="$(git log | grep git-subtree-dir | awk '{ print $2 }' | sort | uniq | sed "s/'//g" | sed "s/\/$//g")"
  printf '%s' "${prefixes}"
}

#
# Map subtree prefix to subtree prefix
#
# $1 = prefix, .* for all prefixes
#
get_subtree_prefixes_config() {
  # .* for all prefixes
  default_re=".*"
  # Do we have "subtree.<something>.prefix = $1" defined in .gittrees file?
  sm_path="${1}"
  re=$(printf '%s\n' "${sm_path}" | sed -e 's/[].[^$\\*]/\\&/g')
  re="${re:-${default_re}}"
  prefixes=$(git config -f "${GITTREES_FILE}" --get-regexp '^subtree\..*\.prefix$' |
    sed -n -e 's|^subtree\.\(.*\)\.prefix '"${re}"'$|\1|p')
  test -z "$prefixes" &&
    die "$(eval_gettext "No subtree mapping found in %s for prefix %s" "${GITTREES_FILE}" "${sm_path}")"
  echo "$prefixes"
}

# Sanitize the local git environment for use within a submodule. We
# can't simply use clear_local_git_env since we want to preserve some
# of the settings from GIT_CONFIG_PARAMETERS.
sanitize_subtree_env() {
  save_config=$GIT_CONFIG_PARAMETERS
  clear_local_git_env
  GIT_CONFIG_PARAMETERS=$save_config
  export GIT_CONFIG_PARAMETERS
}

# Print a subtree status for subtrees in index or working tree
#
# $1 = subtree prefix
# $2 = option
# $3 = default value
#
# print "<merged>\n<configed>"
get_subtree_status() {
  prefix="$1"
  merged_prefixes="$2"
  configed_prefixes="$3"

  merged="N"
  # check to see if the prefix is already in merged commits
  for p in ${merged_prefixes}; do
    if [ "${p}"x = "${prefix}"x ]; then
      merged="Y"
      break
    fi
  done

  configed="N"
  # check to see if the prefix is already in merged commits
  for p in ${configed_prefixes}; do
    if [ "${p}"x = "${prefix}"x ]; then
      configed="Y"
      break
    fi
  done
  printf "%s\n%s\n%s" "${merged}" "${configed}"
}

function cleanup() {
  tmpdir="$1"
  printf "Cleaning up %s..." "${tmpdir}"
  [ -d "${tmpdir}" ] && rm -Rf "${tmpdir}"
  printf "\r\033[KCleaning done."
  printf "\r\033[K"
}

#
# Add a subtree path to correct revision, using clone and checkout as needed
#
# $1 = url
# $2 = branch
# $3 = disabled
#
cmd_add() {
  # parse $args after "... add".

  while test $# -ne 0; do
    case "$1" in
    -P | --prefix)
      case "$2" in '') usage ;; esac
      prefix+=$'\n'
      prefix+=$2
      shift
      ;;
    -q | --quiet)
      GIT_QUIET=1
      ;;
    --)
      shift
      break
      ;;
    -*)
      usage
      ;;
    *)
      break
      ;;
    esac
    shift
  done

  for p in $prefix; do
    prefix="$p"
  done

  if test $# -ne 2 && test $# -ne 3; then
    usage
  fi
  url="$1"
  branch="$2"
  disabled="${3:-false}"

  if [ -z "${prefix}" ] || [ -z "${url}" ] || [ -z "${branch}" ]; then
    echo "some arguments are missing!
    ignore git subtree add|pull --prefix=${prefix} ${url} ${branch} --squash"
    return
  fi

  # check to see if the subpath is already in "${GITTREES_TEMP_FILE}"
  if git config -f "${GITTREES_FILE}" "subtree.${prefix}.url" &>/dev/null; then
    echo "$prefix already configured, ignore"
    return
  fi

  tmpdir=$(mktemp -d -t "git-subtree-XXXXXX" 2>/dev/null) ||
    die "error: mktemp is needed"

  trap "cleanup \"${tmpdir}\"" EXIT
  [ -d "${tmpdir}" ] && rm -Rf "${tmpdir}"
  mkdir -p "${tmpdir}"
  GITTREES_TEMP_FILE="${tmpdir}/${GITTREES_FILE}"
  [ -f "${GITTREES_FILE}" ] && cp "${GITTREES_FILE}" "${GITTREES_TEMP_FILE}"
  echo "${GITTREES_TEMP_FILE}"
  echo "$prefix"

  # record the subtree info
  if git config -f "${GITTREES_TEMP_FILE}" --get subtree."$prefix" &>/dev/null; then
    git config -f "${GITTREES_TEMP_FILE}" --remove-section subtree."$prefix" ||
      die "$(eval_gettext "failed to remove section subtree.'\$prefix'")"
  fi
  git config -f "${GITTREES_TEMP_FILE}" --add subtree."$prefix".prefix "${prefix}" ||
    die "$(eval_gettext "failed to add subtree.'\$prefix'.prefix = '\$prefix'")"
  git config -f "${GITTREES_TEMP_FILE}" --add subtree."$prefix".url "${url}" ||
    die "$(eval_gettext "failed to add subtree.'\$prefix'.url = '\$url'")"
  git config -f "${GITTREES_TEMP_FILE}" --add subtree."$prefix".branch "${branch}" ||
    die "$(eval_gettext "failed to add subtree.'\$prefix'.branch = '\$branch'")"
  git config -f "${GITTREES_TEMP_FILE}" --add subtree."$prefix".disabled "${disabled}" ||
    die "$(eval_gettext "failed to add subtree.'\$prefix'.disabled= '\$disabled'")"

  [ -f "${GITTREES_TEMP_FILE}" ] && mv "${GITTREES_TEMP_FILE}" "${GITTREES_FILE}"

  if [ -z "$GIT_QUIET" ]; then
    say "$(printf "add %s %s %s\n" "${prefix}" "${url}" "${branch}")"
    say "__________"
  fi
}

#
# Tidy each subtree path to correct revision, using clone and checkout as needed
#
cmd_tidy() {
  while test $# -ne 0; do
    case "$1" in
    -P | --prefix)
      case "$2" in '') usage ;; esac
      prefix+=$'\n'
      prefix+=$2
      shift
      ;;
    -q | --quiet)
      GIT_QUIET=1
      ;;
    -u | --update)
      update=1
      ;;
    --)
      shift
      break
      ;;
    -*)
      usage
      ;;
    *)
      break
      ;;
    esac
    shift
  done

  # for each prefix that was subtree merged
  subtree_merged_prefixes=$(get_subtree_prefixes_merged) ||
    die "$(eval_gettext "failed to pasred git log")"
  subtree_configed_prefixes=$(get_subtree_prefixes_config "") ||
    die "$(eval_gettext "failed to pasred \$GITTREES_FILE")"
  subtree_prefixes=
  subtree_prefixes+="${subtree_merged_prefixes}"
  subtree_prefixes+=$'\n'
  subtree_prefixes+="${subtree_configed_prefixes}"
  subtree_prefixes=$(echo "${subtree_configed_prefixes}" | sort -u -k 1,1)

  for prefix in $subtree_prefixes; do
    if [ -z "${prefix}" ]; then
      continue
    fi

    say "__________"

    status=$(get_subtree_status "${prefix}" "${subtree_merged_prefixes}" "${subtree_configed_prefixes}") || exit $?
    commited="$(echo "$status" | awk 'BEGIN{RS="";FS="\n"} { print $1 }')" || exit $?
    configed="$(echo "$status" | awk 'BEGIN{RS="";FS="\n"} { print $2 }')" || exit $?

    if [ "${commited}"x = "Y"x ]; then
      # look for the most recent commit
      set +o pipefail
      commit=$(git log --grep "Squashed '$prefix/'" --oneline | head -n 1 | awk '{ print $NF }') || exit $?
      if [[ "$commit" =~ .. ]]; then
        commit=$(echo "$commit" | cut -d . -f 3) || exit $?
      fi
      set -o pipefail

      say "last commit for $prefix is $commit"
    fi

    url=
    branch=
    disabled=
    if [ "$configed"x = "Y"x ]; then
      set +o pipefail
      url="$(git config -f .gittrees --get subtree."$prefix".url)" || exit $?
      branch="$(git config -f .gittrees --get subtree."$prefix".branch)" || exit $?
      disabled="$(git config -f .gittrees --get subtree."$prefix".disabled)" || exit $?
      set -o pipefail
      if [ "${disabled}"x = "true"x ]; then
        continue
      fi
    fi

    if [ "$configed"x = "N"x ]; then
      # ask for the git url
      printf "adding %s\n" "${prefix}"

      read -r -p "Enter url: " URL
      read -r -p "Enter branch(master as default): " Branch
      Branch="${Branch:=master}"
      options=(
        "Enable"
        "Disable"
      )
      echo "Enable or Disable this subtree"
      select opt in "${options[@]}"; do
        case "$opt" in
        "Enable")
          disabled=false
          break
          ;;
        "Disable")
          disabled=true
          break
          ;;
        *)
          echo "Try again!"
          ;;
        esac
      done

      tmpdir=$(mktemp -d -t "git-subtree-XXXXXX" 2>/dev/null) ||
        die "error: mktemp is needed"

      trap "cleanup \"${tmpdir}\"" EXIT
      [ -d "${tmpdir}" ] && rm -Rf "${tmpdir}"
      mkdir -p "${tmpdir}"
      GITTREES_TEMP_FILE="${tmpdir}/${GITTREES_FILE}"
      [ -f "${GITTREES_FILE}" ] && cp "${GITTREES_FILE}" "${GITTREES_TEMP_FILE}"

      # record the subtree info
      if git config -f "${GITTREES_TEMP_FILE}" --get subtree."$prefix" &>/dev/null; then
        git config -f "${GITTREES_TEMP_FILE}" --remove-section subtree."$prefix" ||
          die "$(eval_gettext "failed to remove section subtree.'\$prefix'")"
      fi
      git config -f "${GITTREES_TEMP_FILE}" --add subtree."$prefix".prefix "${prefix}" ||
        die "$(eval_gettext "failed to add subtree.'\$prefix'.prefix = '\$prefix'")"
      git config -f "${GITTREES_TEMP_FILE}" --add subtree."$prefix".url "${url}" ||
        die "$(eval_gettext "failed to add subtree.'\$prefix'.url = '\$url'")"
      git config -f "${GITTREES_TEMP_FILE}" --add subtree."$prefix".branch "${branch}" ||
        die "$(eval_gettext "failed to add subtree.'\$prefix'.branch = '\$branch'")"
      git config -f "${GITTREES_TEMP_FILE}" --add subtree."$prefix".disabled "${disabled}" ||
        die "$(eval_gettext "failed to add subtree.'\$prefix'.disabled= '\$disabled'")"
    fi

    if [ -z "${prefix}" ] || [ -z "${url}" ] || [ -z "${branch}" ]; then
      echo "some arguments are missing!
    ignore git subtree add|pull --prefix=${prefix} ${url} ${branch} --squash"
      continue
    fi

    if [ "${disabled}"x = "true"x ]; then
      if [ -d "${prefix}" ]; then
        git rm "${prefix}" || exit $?
      fi
      continue
    fi

    if [ -d "${prefix}" ]; then
      if [ -z "$update" ]; then
        continue
      fi

      say "$(printf "pulling %s %s %s\n" "${prefix}" "${url}" "${branch}")"
      if test -z "${message}"; then
        git subtree pull --prefix="${prefix}" "${url}" "${branch}" --squash || exit $?
      else
        git subtree pull --prefix="${prefix}" "${url}" "${branch}" --squash -m "${message}" || exit $?
      fi
      say "$(printf "pull %s %s %s\n" "${prefix}" "${url}" "${branch}")"
    else
      say "$(printf "adding %s %s %s\n" "${prefix}" "${url}" "${branch}")"
      if test -z "${message}"; then
        git subtree add --prefix="${prefix}" "${url}" "${branch}" --squash -m "${message}" || exit $?
      else
        git subtree add --prefix="${prefix}" "${url}" "${branch}" --squash || exit $?
      fi
      say "$(printf "add %s %s %s\n" "${prefix}" "${url}" "${branch}")"
    fi
  done
  say "__________"
}

#
# Update each subtree path to correct revision, using clone and checkout as needed
#
cmd_update() {
  while test $# -ne 0; do
    case "$1" in
    -P | --prefix)
      case "$2" in '') usage ;; esac
      prefix+=$'\n'
      prefix+=$2
      shift
      ;;
    -q | --quiet)
      GIT_QUIET=1
      ;;
    --)
      shift
      break
      ;;
    -*)
      usage
      ;;
    *)
      break
      ;;
    esac
    shift
  done

  subtree_configed_prefixes="${prefix}"
  # for each prefix that was subtree merged
  if [ -z "$subtree_configed_prefixes" ]; then
    subtree_configed_prefixes=$(get_subtree_prefixes_config "")
  fi

  for prefix in $subtree_configed_prefixes; do
    if [ -z "${prefix}" ]; then
      continue
    fi

    say "__________"

    url=
    branch=
    disabled=
    url="$(git config -f .gittrees --get subtree."$prefix".url)"
    branch="$(git config -f .gittrees --get subtree."$prefix".branch)"
    disabled="$(git config -f .gittrees --get subtree."$prefix".disabled)"

    if [ "${disabled}"x = "true"x ]; then
      continue
    fi

    if [ -z "${prefix}" ] || [ -z "${url}" ] || [ -z "${branch}" ]; then
      echo "some arguments are missing!
    ignore git subtree add|pull --prefix=${prefix} ${url} ${branch} --squash"
      continue
    fi

    if [ -d "${prefix}" ]; then
      say "$(printf "pulling %s %s %s\n" "${prefix}" "${url}" "${branch}")"
      if test -z "${message}"; then
        git subtree pull --prefix="${prefix}" "${url}" "${branch}" --squash || exit $?
      else
        git subtree pull --prefix="${prefix}" "${url}" "${branch}" --squash -m "${message}" || exit $?
      fi
      say "$(printf "pull %s %s %s\n" "${prefix}" "${url}" "${branch}")"
    else
      if [ -z "$install" ]; then
        continue
      fi
      say "$(printf "adding %s %s %s\n" "${prefix}" "${url}" "${branch}")"
      if test -z "${message}"; then
        git subtree add --prefix="${prefix}" "${url}" "${branch}" --squash || exit $?
      else
        git subtree add --prefix="${prefix}" "${url}" "${branch}" --squash -m "${message}" || exit $?
      fi
      say "$(printf "add %s %s %s\n" "${prefix}" "${url}" "${branch}")"
    fi
  done

  say "__________"
}

#
# Show configed summary for subtrees in index or working tree
#
cmd_summary() {
  while test $# -ne 0; do
    case "$1" in
    -P | --prefix)
      case "$2" in '') usage ;; esac
      prefix+=$'\n'
      prefix+=$2
      shift
      ;;
    -q | --quiet)
      GIT_QUIET=1
      ;;
    --)
      shift
      break
      ;;
    -*)
      usage
      ;;
    *)
      break
      ;;
    esac
    shift
  done

  subtree_merged_prefixes=$(get_subtree_prefixes_merged)
  subtree_configed_prefixes=$(get_subtree_prefixes_config "")

  subtree_prefixes="${prefix}"
  if [ -z "$subtree_prefixes" ]; then
    # for each prefix that was subtree merged
    subtree_prefixes=
    subtree_prefixes+="${subtree_merged_prefixes}"
    subtree_prefixes+=$'\n'
    subtree_prefixes+="${subtree_configed_prefixes}"
  fi
  subtree_prefixes=$(echo "${subtree_prefixes}" | sort -u -k 1,1)
  printf "subtree_prefixes list below...\n__________\n"

  if [ -z "$GIT_QUIET" ]; then
    msg=$(printf "%s %s %s %s %s %s" "PREFIX" "COMMITED" "CONFIGED" "URL" "BRANCH" "DISABLED")
  else
    msg=$(printf "%s %s %s %s %s" "PREFIX" "COMMITED" "CONFIGED" "BRANCH" "DISABLED")
  fi
  msg+=$'\n'
  for prefix in $subtree_prefixes; do
    if [ -z "${prefix}" ]; then
      continue
    fi
    status=$(get_subtree_status "${prefix}" "${subtree_merged_prefixes}" "${subtree_configed_prefixes}")
    commited="$(echo "$status" | awk 'BEGIN{RS="";FS="\n"} { print $1 }')"
    configed="$(echo "$status" | awk 'BEGIN{RS="";FS="\n"} { print $2 }')"

    url=
    branch=
    disabled=
    if [ "$configed"x = "Y"x ]; then
      url="$(git config -f .gittrees --get subtree."$prefix".url)"
      branch="$(git config -f .gittrees --get subtree."$prefix".branch)"
      disabled="$(git config -f .gittrees --get subtree."$prefix".disabled)"
    fi
    if [ -z "$GIT_QUIET" ]; then
      msg+=$(printf "%s %s %s %s %s %s" "${prefix}" "${commited}" "${configed}" "${url:-<empty>}" "${branch:-<empty>}" "${disabled:-false}")
    else
      msg+=$(printf "%s %s %s %s %s" "${prefix}" "${commited}" "${configed}" "${branch:-<empty>}" "${disabled:-false}")
    fi
    msg+=$'\n'
  done
  echo "${msg}" | column -t
  say "__________"
}

# This loop parses the command line arguments to find the
# subcommand name to dispatch.  Parsing of the subcommand specific
# options are primarily done by the subcommand implementations.
# Subcommand specific options such as --branch and --cached are
# parsed here as well, for backward compatibility.

while test $# != 0 && test -z "$command"; do
  case "$1" in
  add | tidy | update | summary)
    command=$1
    ;;
  -q | --quiet)
    GIT_QUIET=1
    ;;
  --)
    break
    ;;
  -*)
    usage
    ;;
  *)
    break
    ;;
  esac
  shift
done

# No command word defaults to "summary"
if test -z "$command"; then
  if test $# = 0; then
    command=summary
  else
    usage
  fi
fi

#"cmd_$(echo $command | sed -e s/-/_/g)" "$@"
"cmd_$command" "$@"
