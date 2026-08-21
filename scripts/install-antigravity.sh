#!/bin/sh
# install-antigravity.sh - install Spacecraft into Antigravity (global plugin or project-local).
#
# Usage:
#   sh scripts/install-antigravity.sh global              # Install to ~/.gemini/config/plugins/spacecraft
#   sh scripts/install-antigravity.sh project [dir]       # Install into <dir>/.agents + <dir>/GEMINI.md
set -e

MODE="${1:-global}"
TARGET="${2:-.}"
SCRIPT_DIR=$(CDPATH= cd -- "$(dirname -- "$0")" 2>/dev/null && pwd)
SRC=$(CDPATH= cd -- "$SCRIPT_DIR/.." 2>/dev/null && pwd)
LOCAL_BIN="$HOME/.local/bin"
GEMINI_CONFIG="$HOME/.gemini/config"
GLOBAL_PLUGIN_DIR="$GEMINI_CONFIG/plugins/spacecraft"

abspath() { (cd "$1" 2>/dev/null && pwd) || echo "$1"; }

echo "Spacecraft Antigravity Installer"
echo "================================"

# 1. Sync and generate Antigravity assets first
node "$SRC/scripts/sync-antigravity.mjs"

# 2. Build and link CLI
echo "Linking Spacecraft CLI -> $LOCAL_BIN/spacecraft"
mkdir -p "$LOCAL_BIN"
chmod +x "$SRC/cli/spacecraft.mjs"
ln -sf "$SRC/cli/spacecraft.mjs" "$LOCAL_BIN/spacecraft"
ln -sf "$SRC/cli/spacecraft.mjs" "$SRC/spacecraft"

if [ "$MODE" = "global" ]; then
  echo "Installing Global Spacecraft Plugin to $GLOBAL_PLUGIN_DIR"
  mkdir -p "$GLOBAL_PLUGIN_DIR"
  
  # Copy plugin assets
  cp -R "$SRC/plugins/spacecraft/"* "$GLOBAL_PLUGIN_DIR/"
  
  # Ensure hooks script executable
  chmod +x "$GLOBAL_PLUGIN_DIR/hooks/safety-check.mjs" 2>/dev/null || true
  
  echo ""
  echo "Global Antigravity Plugin installed successfully!"
  echo "Location: $GLOBAL_PLUGIN_DIR"
  echo "Skills:   25 Spacecraft skills ready for on-demand activation."
  echo "Rules:    $GLOBAL_PLUGIN_DIR/rules/AGENTS.md"
  echo "Hooks:    $GLOBAL_PLUGIN_DIR/hooks.json"
elif [ "$MODE" = "project" ]; then
  TARGET_ABS=$(abspath "$TARGET")
  echo "Installing Spacecraft into Project: $TARGET_ABS"
  
  mkdir -p "$TARGET_ABS/.agents/rules" "$TARGET_ABS/.agents/skills" "$TARGET_ABS/.agents/hooks"
  
  # Copy rules, hooks, subagents, and skills
  cp "$SRC/plugins/spacecraft/rules/AGENTS.md" "$TARGET_ABS/.agents/rules/AGENTS.md"
  cp "$SRC/plugins/spacecraft/rules/AGENTS.md" "$TARGET_ABS/GEMINI.md"
  cp "$SRC/plugins/spacecraft/hooks.json" "$TARGET_ABS/.agents/hooks.json"
  cp "$SRC/plugins/spacecraft/hooks/safety-check.mjs" "$TARGET_ABS/.agents/hooks/safety-check.mjs"
  chmod +x "$TARGET_ABS/.agents/hooks/safety-check.mjs"
  
  # Copy skills
  cp -R "$SRC/plugins/spacecraft/skills/"* "$TARGET_ABS/.agents/skills/"
  
  # Scaffold .space
  mkdir -p "$TARGET_ABS/.space/missions" "$TARGET_ABS/.space/archive" "$TARGET_ABS/.space/roadmaps"
  node "$SRC/cli/ensure-project-git.mjs" --ready "$TARGET_ABS"
  
  echo ""
  echo "Project installation complete at $TARGET_ABS"
  echo "Rules:   $TARGET_ABS/GEMINI.md and $TARGET_ABS/.agents/rules/AGENTS.md"
  echo "Skills:  $TARGET_ABS/.agents/skills"
  echo "State:   $TARGET_ABS/.space"
else
  echo "Unknown mode: $MODE (must be 'global' or 'project')" >&2
  exit 1
fi

echo ""
echo "Done."
