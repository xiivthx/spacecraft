#!/bin/sh
# bootstrap.sh - install spacecraft into a project (safe, project-local default).
#
# Usage:
#   ./bootstrap.sh [project-dir]                 # from a spacecraft clone
#   curl -fsSL <raw-url>/bootstrap.sh | sh       # remote (clones the repo)
#   curl -fsSL <raw-url>/bootstrap.sh | sh -s -- /path/to/project
#
# Installs the full .cursor surface (rules, agents, skills, hooks, merged MCP)
# and a .space scaffold into the target project, links the Node CLI when Node is
# available, and runs post-install smoke checks. Never writes ~/.cursorrules.
set -e

REPO_URL="${SPACECRAFT_REPO:-https://github.com/xiivthx/spacecraft.git}"
REPO_REF="${SPACECRAFT_REF:-main}"

echo "Spacecraft bootstrap"
echo "===================="

# Parse options: check if --antigravity is requested
ANTIGRAVITY_MODE=0
TARGET_DIR="."
for arg in "$@"; do
  case "$arg" in
    --antigravity|--agy) ANTIGRAVITY_MODE=1 ;;
    -*) ;;
    *) TARGET_DIR="$arg" ;;
  esac
done

mkdir -p "$TARGET_DIR"
TARGET=$(CDPATH= cd -- "$TARGET_DIR" && pwd)

# Resolve the source repo: use the local clone if this script runs from one,
# otherwise clone into a temp dir. $0 is a real file path when run directly
# (./bootstrap.sh or sh bootstrap.sh) and is not when piped via curl|sh, so
# checking for an existing file first tells the two cases apart.
SRC=""
CLEANUP=""
SCRIPT_DIR=""
# Prefer a real on-disk script path (./bootstrap.sh / sh path/to/bootstrap.sh).
# curl|sh often sets $0 to sh or /bin/sh; those resolve to a dir without
# spacecraft sources, so we fall through to a temp clone below.
if [ -f "$0" ]; then
  SCRIPT_DIR=$(CDPATH= cd -- "$(dirname -- "$0")" 2>/dev/null && pwd || true)
fi

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

# Nothing to clean up when SRC was already local; a bare `[ -n ... ]` failure
# would otherwise become the script's exit status under `set -e` + EXIT trap.
cleanup() { [ -z "$CLEANUP" ] || rm -rf "$CLEANUP"; }
trap cleanup EXIT INT TERM

# Sync Antigravity assets
node "$SRC/scripts/sync-antigravity.mjs" >/dev/null 2>&1 || true

# Install the config surface + scaffold + merged MCP for Cursor.
sh "$SRC/scripts/install-cursor.sh" "$TARGET" "$SRC"

# If Antigravity mode or detected, install project Antigravity surface (.agents + GEMINI.md)
if [ "$ANTIGRAVITY_MODE" = "1" ] || [ -f "$TARGET/GEMINI.md" ] || [ -d "$TARGET/.agents" ] || [ -d "$HOME/.gemini" ]; then
  sh "$SRC/scripts/install-antigravity.sh" project "$TARGET"
fi

# Link the Node CLI (cli/spacecraft.mjs) into the target when Node is available.
BIN="$TARGET/spacecraft"
sh "$SRC/scripts/install-binary.sh" "$TARGET" "$SRC" "$REPO_URL"

# Smoke check.
sh "$SRC/scripts/smoke.sh" "$TARGET" "$BIN"

echo ""
echo "Done. $TARGET is spacecraft-ready (Cursor + Antigravity)."
echo ""
echo "Optional global CLI:  ln -sf \"$BIN\" ~/.local/bin/spacecraft"
echo "Restart your IDE / Agent to pick up configuration."
