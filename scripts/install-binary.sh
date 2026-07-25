#!/bin/sh
# install-binary.sh - build or fetch the spacecraft CLI binary.
#
# Prefers `go build` from source. Falls back to an optional prebuilt release
# download (sha256-verified when a checksum is published). Otherwise prints a
# skip message and exits 0 - this is a soft skip; the caller (bootstrap.sh)
# continues to its smoke check either way.
#
# Usage: install-binary.sh <target-dir> <source-repo-dir> [repo-url]
set -e

USAGE="usage: install-binary.sh <target-dir> <source-repo-dir> [repo-url]"
TARGET="${1:?$USAGE}"
SRC="${2:?$USAGE}"
REPO_URL="${3:-https://github.com/xiivthx/spacecraft.git}"

BIN="$TARGET/spacecraft"

if command -v go >/dev/null 2>&1; then
  echo "Building spacecraft CLI from source..."
  ( cd "$SRC/cmd/spacecraft" && go build -o "$BIN" . )
  echo "  spacecraft -> $BIN"
  exit 0
fi

echo "Go not found; attempting optional prebuilt binary..."
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)
case "$ARCH" in
  x86_64|amd64) ARCH=amd64 ;;
  arm64|aarch64) ARCH=arm64 ;;
esac
ASSET="spacecraft-${OS}-${ARCH}"
BASE_URL="${REPO_URL%.git}/releases/latest/download"

if command -v curl >/dev/null 2>&1 && curl -fsSL "$BASE_URL/$ASSET" -o "$BIN" 2>/dev/null; then
  chmod +x "$BIN"
  if curl -fsSL "$BASE_URL/$ASSET.sha256" -o "$BIN.sha256" 2>/dev/null; then
    EXPECTED=$(awk '{print $1}' "$BIN.sha256")
    if command -v shasum >/dev/null 2>&1; then
      ACTUAL=$(shasum -a 256 "$BIN" | awk '{print $1}')
    else
      ACTUAL=$(sha256sum "$BIN" | awk '{print $1}')
    fi
    rm -f "$BIN.sha256"
    if [ "$EXPECTED" != "$ACTUAL" ]; then
      echo "error: checksum mismatch for $ASSET" >&2
      rm -f "$BIN"
      exit 1
    fi
    echo "  spacecraft -> $BIN (checksum verified)"
  else
    echo "  spacecraft -> $BIN (no checksum published; unverified)"
  fi
else
  echo "  skipped - install Go and run 'make build', or build cmd/spacecraft manually."
fi
