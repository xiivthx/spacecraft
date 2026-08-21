#!/bin/sh
# install-global-hooks.sh - install safety-only Cursor hooks into a global
# ~/.cursor dir, merge-safe with unrelated user hooks.
#
# Copies safety hooks (check-main-write, check-ship-commands, block-secrets-read,
# block-destructive - never session-start.sh, which is project-layer only) into
# <global-cursor-dir>/hooks/, rewrites their hooks.json commands to absolute
# paths under that dir (never the project's repo-relative .cursor/hooks/...),
# and merges the result into <global-cursor-dir>/hooks.json without
# clobbering unrelated hooks already there.
#
# Usage: install-global-hooks.sh <repo-root> <global-cursor-dir>
set -e

ROOT="${1:?usage: install-global-hooks.sh <repo-root> <global-cursor-dir>}"
GLOBAL="${2:?usage: install-global-hooks.sh <repo-root> <global-cursor-dir>}"

abspath() { (cd "$1" 2>/dev/null && pwd) || echo "$1"; }
ROOT_ABS=$(abspath "$ROOT")
mkdir -p "$GLOBAL"
GLOBAL_ABS=$(abspath "$GLOBAL")

SAFETY_HOOKS="check-main-write.sh check-ship-commands.sh block-secrets-read.sh block-destructive.sh"

mkdir -p "$GLOBAL_ABS/hooks"
for hook in $SAFETY_HOOKS; do
  cp "$ROOT_ABS/.cursor/hooks/$hook" "$GLOBAL_ABS/hooks/$hook"
done
echo "  hooks -> $GLOBAL_ABS/hooks ($SAFETY_HOOKS)"

mkdir -p "$ROOT_ABS/.tmp"
tmp_src=$(mktemp "$ROOT_ABS/.tmp/global-hooks-src.XXXXXX")
trap 'rm -f "$tmp_src"' EXIT INT TERM

python3 "$ROOT_ABS/scripts/rewrite-global-hooks.py" \
  "$ROOT_ABS/.cursor/hooks.json" "$GLOBAL_ABS/hooks" "$tmp_src" $SAFETY_HOOKS

python3 "$ROOT_ABS/scripts/hooks-merge.py" merge "$GLOBAL_ABS/hooks.json" "$tmp_src" \
  | sed 's/^/  hooks: /'
