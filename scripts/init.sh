#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail
# set -o xtrace

PROJECT_ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")"/.. && pwd -P)

source "${PROJECT_ROOT}/scripts/common.sh"

log_info "Installing go tools"
go install \
  github.com/golang/mock/mockgen
