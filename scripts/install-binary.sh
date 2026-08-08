#!/bin/sh
# install-binary.sh - install the spacecraft Node CLI entry.
#
# Links cli/spacecraft.mjs into the target directory. Requires Node.js on PATH.
# Soft-skips (exit 0) when Node is missing so bootstrap can continue to smoke.
#
# Usage: install-binary.sh <target-dir> <source-repo-dir> [repo-url]
set -e

USAGE="usage: install-binary.sh <target-dir> <source-repo-dir> [repo-url]"
TARGET="${1:?$USAGE}"
SRC="${2:?$USAGE}"
# [repo-url] accepted for bootstrap compatibility; install uses the source tree.

# Resolve absolute paths so the symlink works from any target directory.
TARGET=$(CDPATH= cd -- "$TARGET" 2>/dev/null && pwd || { mkdir -p "$TARGET" && CDPATH= cd -- "$TARGET" && pwd; })
SRC=$(CDPATH= cd -- "$SRC" && pwd)
BIN="$TARGET/spacecraft"
ENTRY="$SRC/cli/spacecraft.mjs"

if [ ! -f "$ENTRY" ]; then
  echo "error: missing Node CLI entry at $ENTRY" >&2
  exit 1
fi

if ! command -v node >/dev/null 2>&1; then
  echo "  skipped - install Node.js and re-run, or run 'make build' from a checkout."
  exit 0
fi

echo "Installing spacecraft CLI (Node)..."
mkdir -p "$TARGET"
chmod +x "$ENTRY"
ln -sf "$ENTRY" "$BIN"
echo "  spacecraft -> $BIN"
