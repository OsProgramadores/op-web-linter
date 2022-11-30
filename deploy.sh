#!/bin/bash
# Simple deployment script for op-web-linter.

readonly PROGNAME="${0##*/}"
export BUILD_DIR="/tmp/build.$$"
export REPO_DIR="/tmp/build.$$/op-web-linter"

function deploy_repo_copy() {
  local dest="${1?}"

  ssh -t -o ClearAllForwardings=yes -oBatchMode=yes "${dest}" '
    BUILD_DIR="/tmp/build.$$"
    rm -rf "${BUILD_DIR}"
    mkdir -p "${BUILD_DIR}"
    set -eu
    cd "${BUILD_DIR}"
    rm -rf op-web-linter
    git clone https://github.com/osprogramadores/op-web-linter
    cd op-web-linter
    make clean
    make docker
    echo "Docker finished with code ${?}"
    sudo systemctl restart op-web-linter
  '
}

function deploy_local_copy() {
  local dest="${1?}"

  ssh -t -o ClearAllForwardings=yes -oBatchMode=yes "${dest}" "
    rm -rf \"${BUILD_DIR}\"
    mkdir -p \"${REPO_DIR}\"
  "

  rsync -avp . "${dest}:$REPO_DIR"

  ssh -t -o ClearAllForwardings=yes -oBatchMode=yes "${dest}" "
    set -eu
    cd \"${REPO_DIR}\"
    make clean
    make docker
    echo \"Docker finished with code ${?}\"
    sudo systemctl restart op-web-linter
    echo \"Finished\"
  "
}

# Die with an error message.
function die() {
  local msg="$1"
  echo >&2 "$PROGNAME: $msg"
  exit 1
}

# Prints the program usage and quit.
function usage() {
  echo >&2 "Use: $PROGNAME [-h|--help] [-d|--destination=destination_host] [-l|--local]"
  exit 2
}

# Prints important information and waits for confirmation
function info() {
  echo '
OP-web-linter deployment script.

IMPORTANT: This script installs op-web-linter on a remote host and restarts
the docker container there to reflect the change. It assumes a few things:

1) You can ssh to the destination host with your default SSH account.
2) You have a sudo enabled account on the destination host.
3) op-web-linter is started by systemd, unit name op-web-linter.
4) You can build containers using your account at the destination host.

This script will FAIL IN UNPREDICTABLE WAYS if one of the criteria above
is not met.

Please press ENTER to acknowledge this.
'
  read -r _
}

function main() {
  TEMP=$(getopt -o d:hl --long destination,help,local -n "${0##*/}" -- "$@")
  # shellcheck disable=SC2181
  if [ $? != 0 ] ; then echo "Terminating..." >&2 ; exit 2 ; fi

  localcopy=0
  destination=""

  eval set -- "$TEMP"
  while :; do
    case "$1" in
      -d|--destination)
        destination="$2"; shift 2 ;;
      -h|--help)
        usage ;;
      -l|--local)
        localcopy=1; shift ;;
      --)
        shift; break ;;
      *)
        echo "Internal Error."; exit 2 ;;
    esac
  done

  [[ -z "$destination" ]] && die "Please specify destination host with --destination (-d)"

  info

  if (( localcopy == 1 )); then
    deploy_local_copy "${destination}"
  else
    deploy_repo_copy "${destination}"
  fi

}

main "${@}"
