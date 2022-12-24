#!/bin/bash
# Simple deployment script for op-web-linter.

readonly PROGNAME="${0##*/}"

# Also check the systemd unit file if you change any of these.
export PROJECT="op-web-linter"
export BUILD_DIR="/tmp/build.$$"
export REPO_DIR="/tmp/build.$$/${PROJECT}"

# Remote data directory (will be created if it doesn't exist).
export DATA_DIR="/var/lib/${PROJECT}"

# Remote bot userid (ideally, an user with this ID should exist).
export PROJECT_UID=501

function deploy_local_copy() {
  local dest="${1?}"

  ssh -t -o ClearAllForwardings=yes -oBatchMode=yes "${dest}" "
    rm -rf \"${BUILD_DIR}\"
    mkdir -p \"${REPO_DIR}\"
  "

  rsync -avp --exclude .git/ . "${dest}:$REPO_DIR"

  ssh -t -o ClearAllForwardings=yes -oBatchMode=yes "${dest}" "
    set -eu
    cd \"${REPO_DIR}\"
    make clean
    make docker
    echo \"Docker finished with code ${?}\"
    cat >sudo.sh <<EOF
      mkdir -p \"${DATA_DIR}\";
      chown -R \"${PROJECT_UID}\" \"${DATA_DIR}\";
      chmod 755 \"${DATA_DIR}\";
      cp site-configs/*.service /etc/systemd/system;
      systemctl daemon-reload;
      systemctl restart \"${PROJECT}\"
EOF
    sudo bash ./sudo.sh
    cd /tmp
    rm -rf \"${BUILD_DIR}\"
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
  echo >&2 "Use: $PROGNAME [-h|--help] [-d|--destination=destination_host]"
  exit 2
}

# Prints important information and waits for confirmation
function info() {
  echo "
Deployment script for: ${PROJECT}.

IMPORTANT: This script installs ${PROJECT} on a remote host using the local
repository contents, and restarts the docker container there to reflect the
change. It assumes a few things:

1) You can ssh to the destination host with this SSH account.
2) Your local repo has the site-configs directory unlocked (git-crypt unlock).
3) You have a sudo enabled account on the destination host.
4) ${PROJECT} is started by systemd, unit name ${PROJECT}.
5) You can build containers using your account at the destination host.

This script will FAIL IN UNPREDICTABLE WAYS if one of the criteria above
is not met.

Please press ENTER to acknowledge this.
"
  read -r _
}

function main() {
  if [[ ! -d ./.git ]]; then
    die "Run this program from the git repository root directory."
  fi

  TEMP=$(getopt -o d:h --long destination,help -n "${0##*/}" -- "$@")
  # shellcheck disable=SC2181
  if [ $? != 0 ] ; then echo "Terminating..." >&2 ; exit 2 ; fi

  destination=""

  eval set -- "$TEMP"
  while :; do
    case "$1" in
      -d|--destination)
        destination="$2"; shift 2 ;;
      -h|--help)
        usage ;;
      --)
        shift; break ;;
      *)
        echo "Internal Error."; exit 2 ;;
    esac
  done

  [[ -z "$destination" ]] && die "Please specify destination host with --destination (-d)"

  info
  deploy_local_copy "${destination}"
}

main "${@}"
