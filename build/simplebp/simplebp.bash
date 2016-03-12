#!/bin/bash

# simplebp.bash
#
# This script is a convenience wrapper that invokes the Blueprint blueprint.bash
# script to make sure our bootstrap inputs didn't change, and then (optionally)
# runs ninja. It must be copied/linked into the build dir to get the paths
# correct.

# We move to the build dir to pick up our extra bootstrap variables like
# SRCDIR_FROM_BUILDDIR and BUILDDIR.
cd $(dirname ${BASH_SOURCE[0]})
source .simplebp.bootstrap

# Now we move back to the root source dir and invoke Blueprint's wrapper script.
cd ${SRCDIR_FROM_BUILDDIR}

# But first make sure we've built ninja.
if [[ ! -x build/ninja/ninja ]]; then
    pushd build/ninja
    ./configure.py --bootstrap
    popd
fi
BUILDDIR="${BUILDDIR}" SKIP_NINJA=true build/blueprint/blueprint.bash

# Run ninja, using the build.ninja file in the build dir. We keep the root
# source dir as the root directory though, to keep paths predictable during the
# build process. Make multiple build lines for a single target an error, to
# catch potential bugs in the simplebp generator.
build/ninja/ninja -f "${BUILDDIR}/build.ninja" -w dupbuild=err "$@"
