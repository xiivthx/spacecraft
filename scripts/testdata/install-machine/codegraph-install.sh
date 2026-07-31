#!/bin/sh
# Fixture substitute for codegraph official install.sh (network-free).
# Writes SPACECRAFT_COMPANION_MARKERS/codegraph.install when invoked.
# Real wire step is separate: `codegraph install --target=cursor --yes` (see bin/codegraph).
set -e
markers="${SPACECRAFT_COMPANION_MARKERS:?SPACECRAFT_COMPANION_MARKERS required}"
mkdir -p "$markers"
printf '%s\n' "codegraph-install-ok" > "$markers/codegraph.install"
printf '%s\n' "v0.0.3-fixture" > "$markers/codegraph.version"
