#!/bin/sh
# global-sync.sh - copy or remove spacecraft's sc-* skills/agents in a global
# Cursor dir.
#
# Usage: global-sync.sh <repo-root> <global-cursor-dir> skills|agents install|uninstall
#
# Skills install profile (SPACECRAFT_SKILL_PROFILE, default lean):
#   lean - User-layer lifecycle + process skills only (see LEAN_SKILLS)
#   full - every repo sc-* skill (domain encyclopedias included)
set -e

USAGE="usage: global-sync.sh <repo-root> <global-cursor-dir> skills|agents install|uninstall"
ROOT="${1:?$USAGE}"
GLOBAL="${2:?$USAGE}"
KIND="${3:?$USAGE}"
MODE="${4:?$USAGE}"

# Lean-core User-layer skills: lifecycle + process. Domain packs stay
# project-layer unless SPACECRAFT_SKILL_PROFILE=full.
LEAN_SKILLS="sc-discuss sc-run sc-ship sc-quick sc-mission sc-planning sc-tdd sc-verification sc-judge sc-clarify sc-git sc-search sc-storm sc-writer"

case "$MODE" in
  install|uninstall) ;;
  *) echo "error: $USAGE" >&2; exit 1 ;;
esac

is_lean_skill() {
  case " $LEAN_SKILLS " in
    *" $1 "*) return 0 ;;
    *) return 1 ;;
  esac
}

case "$KIND" in
  skills)
    mkdir -p "$GLOBAL/skills"
    profile="${SPACECRAFT_SKILL_PROFILE:-lean}"
    case "$profile" in
      lean|full) ;;
      *)
        echo "error: SPACECRAFT_SKILL_PROFILE must be lean or full (got: $profile)" >&2
        exit 1
        ;;
    esac
    for d in "$ROOT"/.cursor/skills/sc-*; do
      [ -d "$d" ] || continue
      name=$(basename "$d")
      # Lean install: copy allowlist only; prune repo-managed packs outside it
      # (destructive for those domain encyclopedias under GLOBAL/skills only).
      # Uninstall still clears every repo sc-* skill under the global dir.
      if [ "$MODE" = install ] && [ "$profile" = lean ] && ! is_lean_skill "$name"; then
        rm -rf "$GLOBAL/skills/$name"
        continue
      fi
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
