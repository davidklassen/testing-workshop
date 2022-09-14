#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail
# set -o xtrace

COLOR_GREEN='\033[0;32m'
COLOR_RED='\033[0;31m'
COLOR_BOLD='\033[1m'
COLOR_NONE='\033[0m'

function log_info {
  echo >&2 -n -e "${COLOR_GREEN}"
  echo >&2 "$@"
  echo >&2 -n -e "${COLOR_NONE}"
}

function log_error {
  echo >&2 -n -e "${COLOR_BOLD}${COLOR_RED}"
  echo >&2 "$@"
  echo >&2 -n -e "${COLOR_NONE}"
}
