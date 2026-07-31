#!/bin/sh
# Failing fixture for rtk install.sh (soft-fail smoke).
# Exits non-zero without writing success markers.
set -e
markers="${SPACECRAFT_COMPANION_MARKERS:?SPACECRAFT_COMPANION_MARKERS required}"
mkdir -p "$markers"
printf '%s\n' "rtk-install-fail" > "$markers/rtk.install.attempted"
exit 1
