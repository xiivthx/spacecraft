#!/bin/sh
# Fixture substitute for caveman official install.sh (network-free).
# Writes SPACECRAFT_COMPANION_MARKERS/caveman.install when invoked.
set -e
markers="${SPACECRAFT_COMPANION_MARKERS:?SPACECRAFT_COMPANION_MARKERS required}"
mkdir -p "$markers"
printf '%s\n' "caveman-install-ok" > "$markers/caveman.install"
printf '%s\n' "v0.0.1-fixture" > "$markers/caveman.version"
