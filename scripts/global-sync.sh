#!/bin/sh
# global-sync.sh - copy or remove spacecraft's sc-* skills/agents in a global
# Cursor dir.
#
# Usage: global-sync.sh <repo-root> <global-cursor-dir> skills|agents install|uninstall
set -e

USAGE="usage: global-sync.sh <repo-root> <global-cursor-dir> skills|agents install|uninstall"
ROOT="${1:?$USAGE}"
GLOBAL="${2:?$USAGE}"
KIND="${3:?$USAGE}"
MODE="${4:?$USAGE}"

case "$MODE" in
  install|uninstall) ;;
  *) echo "error: $USAGE" >&2; exit 1 ;;
esac

case "$KIND" in
  skills)
    mkdir -p "$GLOBAL/skills"
    for d in "$ROOT"/.cursor/skills/sc-*; do
      [ -d "$d" ] || continue
      name=$(basename "$d")
      rm -rf "$GLOBAL/skills/$name"
      if [ "$MODE" = install ]; then
        cp -R "$d" "$GLOBAL/skills/$name"
      fi
    done
    ;;
  agents)
    mkdir -p "$GLOBAL/agents"
    for f in "$ROOT"/.cursor/agents/sc-*.md; do
      [ -f "$f" ] || continue
      name=$(basename "$f")
      rm -f "$GLOBAL/agents/$name"
      if [ "$MODE" = install ]; then
        cp "$f" "$GLOBAL/agents/$name"
      fi
    done
    ;;
  *) echo "error: $USAGE" >&2; exit 1 ;;
esac
