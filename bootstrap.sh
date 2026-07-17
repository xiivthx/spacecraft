#!/bin/sh
# bootstrap.sh — install spacecraft into a project (safe, project-local default).
#
# Usage:
#   ./bootstrap.sh [project-dir]                 # from a spacecraft clone
#   curl -fsSL <raw-url>/bootstrap.sh | sh       # remote (clones the repo)
#   curl -fsSL <raw-url>/bootstrap.sh | sh -s -- /path/to/project
#
# Installs the full .cursor surface (rules, agents, skills, hooks, merged MCP)
# and a .space scaffold into the target project, builds the CLI when Go is
# available, and runs post-install smoke checks. Never writes ~/.cursorrules.
set -e

REPO_URL="${SPACECRAFT_REPO:-https://github.com/xiivthx/spacecraft.git}"
REPO_REF="${SPACECRAFT_REF:-main}"
TARGET="${1:-.}"

echo "Spacecraft bootstrap"
echo "===================="

# Resolve the source repo: use the local clone if this script lives in one,
# otherwise clone into a temp dir.
SRC=""
CLEANUP=""
if [ -n "${BASH_SOURCE:-}" ]; then
  SELF="${BASH_SOURCE}"
else
  SELF="$0"
fi
SCRIPT_DIR=$(CDPATH= cd -- "$(dirname -- "$SELF")" 2>/dev/null && pwd || true)

if [ -n "$SCRIPT_DIR" ] && [ -d "$SCRIPT_DIR/.cursor/rules" ] && [ -f "$SCRIPT_DIR/scripts/install-cursor.sh" ]; then
  SRC="$SCRIPT_DIR"
  echo "Using local source: $SRC"
else
  if ! command -v git >/dev/null 2>&1; then
    echo "error: git is required to bootstrap remotely" >&2
    exit 1
  fi
  SRC=$(mktemp -d 2>/dev/null || mktemp -d -t spacecraft)
  CLEANUP="$SRC"
  echo "Cloning $REPO_URL ($REPO_REF)..."
  git clone --depth 1 --branch "$REPO_REF" "$REPO_URL" "$SRC" >/dev/null 2>&1 \
    || git clone --depth 1 "$REPO_URL" "$SRC" >/dev/null 2>&1
fi

cleanup() { [ -n "$CLEANUP" ] && rm -rf "$CLEANUP"; }
trap cleanup EXIT INT TERM

# Install the config surface + scaffold + merged MCP.
sh "$SRC/scripts/install-cursor.sh" "$TARGET" "$SRC"

# Build (preferred) or fetch the CLI binary.
BIN="$TARGET/spacecraft"
if command -v go >/dev/null 2>&1; then
  echo "Building spacecraft CLI from source..."
  ( cd "$SRC/cmd/spacecraft" && go build -o "$(cd "$TARGET" && pwd)/spacecraft" . )
  echo "  spacecraft -> $BIN"
else
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
    echo "  skipped — install Go and run 'make build', or build cmd/spacecraft manually."
  fi
fi

# Smoke check.
sh "$SRC/scripts/smoke.sh" "$TARGET" "$BIN"

echo ""
echo "Done. $TARGET is spacecraft-ready."
echo ""
echo "Optional global CLI:  ln -sf \"$(cd "$TARGET" && pwd)/spacecraft\" ~/.local/bin/spacecraft"
echo "Restart Cursor to pick up config."
