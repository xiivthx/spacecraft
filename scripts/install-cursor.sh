#!/bin/sh
# install-cursor.sh - install the full spacecraft .cursor surface into a project.
#
# Copies rules, agents, skills, and hooks (if present) into <target>/.cursor,
# scaffolds <target>/.space, and merges mcp.json without clobbering unrelated
# user MCP servers. Safe default: everything is project-local.
#
# Usage: install-cursor.sh <target-project-dir> <source-repo-dir>
set -e

TARGET="${1:?usage: install-cursor.sh <target-project-dir> <source-repo-dir>}"
SRC="${2:?usage: install-cursor.sh <target-project-dir> <source-repo-dir>}"

abspath() { (cd "$1" 2>/dev/null && pwd) || echo "$1"; }

mkdir -p "$TARGET"
TARGET_ABS=$(abspath "$TARGET")
SRC_ABS=$(abspath "$SRC")

if [ ! -d "$SRC_ABS/.cursor/rules" ]; then
  echo "error: source $SRC_ABS is not a spacecraft repo (.cursor/rules missing)" >&2
  exit 1
fi

echo "Installing spacecraft into $TARGET_ABS"

# .space mission state scaffold (idempotent).
mkdir -p "$TARGET_ABS/.space/missions" "$TARGET_ABS/.space/archive" "$TARGET_ABS/.space/roadmaps"

if [ "$TARGET_ABS" = "$SRC_ABS" ]; then
  echo "  source == target; config already in place, scaffolding .space only"
else
  mkdir -p "$TARGET_ABS/.cursor/rules" "$TARGET_ABS/.cursor/agents" "$TARGET_ABS/.cursor/skills"
  # Project layer gets domain/glob rules only (300-620); alwaysApply rules
  # (000/025/050/100/200) are User layer via install-global's USER-RULES.txt.
  for rule in "$SRC_ABS"/.cursor/rules/*.mdc; do
    [ -f "$rule" ] || continue
    grep -q '^alwaysApply: true$' "$rule" && continue
    cp "$rule" "$TARGET_ABS/.cursor/rules/"
  done
  cp -R "$SRC_ABS/.cursor/agents/." "$TARGET_ABS/.cursor/agents/"
  cp -R "$SRC_ABS/.cursor/skills/." "$TARGET_ABS/.cursor/skills/"
  echo "  domain rules, agents, skills -> $TARGET_ABS/.cursor"

  if [ -f "$SRC_ABS/.cursor/hooks.json" ]; then
    if [ -d "$SRC_ABS/.cursor/hooks" ]; then
      mkdir -p "$TARGET_ABS/.cursor/hooks"
      cp -R "$SRC_ABS/.cursor/hooks/." "$TARGET_ABS/.cursor/hooks/"
    fi
    if [ -f "$TARGET_ABS/.cursor/hooks.json" ]; then
      # Merge into existing config so unrelated user hooks survive.
      python3 "$SRC_ABS/scripts/hooks-merge.py" merge \
        "$TARGET_ABS/.cursor/hooks.json" "$SRC_ABS/.cursor/hooks.json" \
        | sed 's/^/  hooks: /'
    else
      cp "$SRC_ABS/.cursor/hooks.json" "$TARGET_ABS/.cursor/hooks.json"
      echo "  hooks -> $TARGET_ABS/.cursor"
    fi
  fi
fi

# Merge MCP servers instead of overwriting the whole file.
if [ -f "$SRC_ABS/.cursor/mcp.json" ]; then
  python3 "$SRC_ABS/scripts/mcp-merge.py" merge \
    "$TARGET_ABS/.cursor/mcp.json" "$SRC_ABS/.cursor/mcp.json" \
    | sed 's/^/  mcp: /'
fi

echo "Done."
