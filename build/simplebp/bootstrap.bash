#!/bin/bash

set -e

export BOOTSTRAP=$(readlink -f "${BASH_SOURCE[0]}")
BOOTSTRAP_DIR=$(dirname "${BOOTSTRAP}")
export SRCDIR=$(readlink -f "${BOOTSTRAP_DIR}/../..")
export BOOTSTRAP_MANIFEST="${BOOTSTRAP_DIR}/build.ninja.in"

"${SRCDIR}/build/blueprint/bootstrap.bash" "$@"
