#!/bin/sh
# Fixture substitute for rtk official install.sh (network-free).
# Writes SPACECRAFT_COMPANION_MARKERS/rtk.install when invoked.
# Real wire step is separate: `rtk init -g --agent cursor` (see bin/rtk).
set -e
markers="${SPACECRAFT_COMPANION_MARKERS:?SPACECRAFT_COMPANION_MARKERS required}"
mkdir -p "$markers"
printf '%s\n' "rtk-install-ok" > "$markers/rtk.install"
printf '%s\n' "v0.0.2-fixture" > "$markers/rtk.version"
