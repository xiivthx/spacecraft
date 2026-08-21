#!/bin/sh
# install-cursor.sh - install the project-layer spacecraft .cursor surface.
#
# Copies domain/glob rules, hard-contract rule, domain-pack skills, and hooks
# (session-start + safety) into <target>/.cursor, scaffolds <target>/.space,
# and merges mcp.json without clobbering unrelated user MCP servers. Agents and
# lean-core skills stay User-layer only (~/.cursor via install-global).
#
# Usage: install-cursor.sh <target-project-dir> <source-repo-dir>
set -e

TARGET="${1:?usage: install-cursor.sh <target-project-dir> <source-repo-dir>}"
SRC="${2:?usage: install-cursor.sh <target-project-dir> <source-repo-dir>}"

# Lean-core User-layer skills (keep identical to scripts/global-sync.sh).
# Project layer gets domain packs only; lean skills live under ~/.cursor/skills.
LEAN_SKILLS="sc-discuss sc-run sc-ship sc-quick sc-mission sc-planning sc-tdd sc-verification sc-judge sc-clarify sc-git sc-search sc-storm sc-writer"

# Project hooks: session-start + hard safety (cloud agents only load project hooks).
PROJECT_HOOKS="session-start.sh check-main-write.sh check-ship-commands.sh block-secrets-read.sh block-destructive.sh"

is_lean_skill() {
  case " $LEAN_SKILLS " in
    *" $1 "*) return 0 ;;
    *) return 1 ;;
  esac
}

abspath() { (cd "$1" 2>/dev/null && pwd) || echo "$1"; }

mkdir -p "$TARGET"
TARGET_ABS=$(abspath "$TARGET")
SRC_ABS=$(abspath "$SRC")

if [ ! -d "$SRC_ABS/.cursor/rules" ]; then
  echo "error: source $SRC_ABS is not a spacecraft repo (.cursor/rules missing)" >&2
  exit 1
fi

echo "Installing spacecraft into $TARGET_ABS"

# .space mission state scaffold (idempotent). Detect first create before mkdir
# so ensure matches CLI: ready (template gitignore + git init) vs ignore (append).
SPACE_WAS_MISSING=0
[ -d "$TARGET_ABS/.space" ] || SPACE_WAS_MISSING=1
mkdir -p "$TARGET_ABS/.space/missions" "$TARGET_ABS/.space/archive" "$TARGET_ABS/.space/roadmaps"
if [ "$SPACE_WAS_MISSING" -eq 1 ]; then
  node "$SRC_ABS/cli/ensure-project-git.mjs" --ready "$TARGET_ABS"
else
  node "$SRC_ABS/cli/ensure-project-git.mjs" --ignore "$TARGET_ABS"
fi


if [ "$TARGET_ABS" = "$SRC_ABS" ]; then
  echo "  source == target; config already in place, scaffolding .space only"
else
  mkdir -p "$TARGET_ABS/.cursor/rules" "$TARGET_ABS/.cursor/skills"
  # Project layer: domain/glob rules 300-620 + alwaysApply hard-contract (010).
  # Soft User-layer depth rules (000/026/027/050/100/200) stay out of the project.
  USER_LAYER="000-spacecraft.mdc 026-intent-coach.mdc 027-th-en-hil.mdc 050-style.mdc 100-conventions.mdc 200-workflow.mdc"
  for rule in "$SRC_ABS"/.cursor/rules/*.mdc; do
    [ -f "$rule" ] || continue
    base=$(basename "$rule")
    case " $USER_LAYER " in
      *" $base "*) continue ;;
    esac
    cp "$rule" "$TARGET_ABS/.cursor/rules/"
  done

  # Agents stay User-layer only; prune any leftover spacecraft agents.
  if [ -d "$TARGET_ABS/.cursor/agents" ]; then
    for agent in "$TARGET_ABS"/.cursor/agents/sc-*.md; do
      [ -f "$agent" ] || continue
      rm -f "$agent"
    done
  fi

  # Domain packs only: skip lean-core and prune any leftover lean skill dirs.
  for skill_dir in "$SRC_ABS"/.cursor/skills/sc-*; do
    [ -d "$skill_dir" ] || continue
    name=$(basename "$skill_dir")
    if is_lean_skill "$name"; then
      rm -rf "$TARGET_ABS/.cursor/skills/$name"
      continue
    fi
    rm -rf "$TARGET_ABS/.cursor/skills/$name"
    cp -R "$skill_dir" "$TARGET_ABS/.cursor/skills/$name"
  done
  echo "  domain rules, hard-contract, domain skills -> $TARGET_ABS/.cursor"

  if [ -f "$SRC_ABS/.cursor/hooks.json" ]; then
    mkdir -p "$TARGET_ABS/.cursor/hooks" "$SRC_ABS/.tmp"
    # Copy session-start + safety hook scripts (dual-layer: project for cloud).
    for hook in $PROJECT_HOOKS; do
      if [ -f "$SRC_ABS/.cursor/hooks/$hook" ]; then
        cp "$SRC_ABS/.cursor/hooks/$hook" "$TARGET_ABS/.cursor/hooks/$hook"
      fi
    done

    tmp_project=$(mktemp "$SRC_ABS/.tmp/project-hooks-src.XXXXXX")
    # shellcheck disable=SC2064
    trap 'rm -f "$tmp_project"' EXIT INT TERM

    # Filter source hooks.json to PROJECT_HOOKS (relative .cursor/hooks/ paths).
    python3 "$SRC_ABS/scripts/rewrite-global-hooks.py" \
      "$SRC_ABS/.cursor/hooks.json" ".cursor/hooks" "$tmp_project" $PROJECT_HOOKS

    if [ -f "$TARGET_ABS/.cursor/hooks.json" ]; then
      python3 "$SRC_ABS/scripts/hooks-merge.py" merge \
        "$TARGET_ABS/.cursor/hooks.json" "$tmp_project" \
        | sed 's/^/  hooks: /'
    else
      cp "$tmp_project" "$TARGET_ABS/.cursor/hooks.json"
      echo "  hooks -> $TARGET_ABS/.cursor (session-start + safety)"
    fi

    rm -f "$tmp_project"
    trap - EXIT INT TERM
  fi
fi

# Merge MCP servers instead of overwriting the whole file.
if [ -f "$SRC_ABS/.cursor/mcp.json" ]; then
  python3 "$SRC_ABS/scripts/mcp-merge.py" merge \
    "$TARGET_ABS/.cursor/mcp.json" "$SRC_ABS/.cursor/mcp.json" \
    | sed 's/^/  mcp: /'
fi

echo "Done."
