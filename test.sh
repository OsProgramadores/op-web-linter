#!/bin/bash
#
# Simple test script. To use:
# 1. Build the tool with 'make'
# 2. Run the linter with ./op-web-linter
# 3. Run ./test.sh lang program (E.g: test.sh go program.go)

if [[ $# -ne 2 ]]; then
  echo >&2 "Use: ${0##*/} language source_file"
  exit 1
fi

LANG="$1"
FILE="$2"

curl -v --json "{ \"lang\":\"${LANG}\", \"text\":\"$(base64 -w 0 < "${FILE}")\" }" http://localhost:10000/lint
