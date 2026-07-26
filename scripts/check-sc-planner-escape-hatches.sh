#!/bin/sh
# Assert sc-planner.md documents: same-mission phase split, and multi-mission
# handoff to /sc-discuss + mission-sizing (planner must not create maps).
set -e

ROOT="${1:-.}"
FILE="$ROOT/.cursor/agents/sc-planner.md"

if [ ! -f "$FILE" ]; then
  echo "FAIL: missing $FILE"
  exit 1
fi

# Require explicit plan-phaseN.json naming (not only "Phase 1 / Phase 2").
if ! grep -Eq 'plan-phaseN\.json|plan-phase<N>\.json|plan-phase[0-9]+\.json' "$FILE"; then
  echo "FAIL: $FILE must name plan-phaseN.json (or plan-phase<N>.json / plan-phase1.json)"
  exit 1
fi

# Require same-mission framing for the phase-split path.
if ! grep -Eiq 'same[- ]mission' "$FILE"; then
  echo "FAIL: $FILE must document same-mission phase split (not vague Phase 1 / Phase 2 alone)"
  exit 1
fi

# Require handoff to discuss + mission-sizing for multi-mission.
if ! grep -Eq '/sc-discuss' "$FILE"; then
  echo "FAIL: $FILE must hand multi-mission splits to /sc-discuss"
  exit 1
fi

if ! grep -Eq 'mission-sizing' "$FILE"; then
  echo "FAIL: $FILE must name mission-sizing"
  exit 1
fi

if ! grep -Eiq 'multi[- ]mission|roadmap' "$FILE"; then
  echo "FAIL: $FILE must document multi-mission or roadmap handoff"
  exit 1
fi

# Require ban on planner-owned map create (spacecraft map named in Must-not sense).
if ! grep -Eq 'spacecraft[[:space:]]+map' "$FILE"; then
  echo "FAIL: $FILE must name spacecraft map (as discuss-owned / planner Must-not)"
  exit 1
fi

if ! grep -Eiq 'never create|discuss owns|must not.*map' "$FILE"; then
  echo "FAIL: $FILE must forbid planner-owned map create"
  exit 1
fi

echo "ok: sc-planner documents phase-split and discuss handoff for multi-mission"
exit 0
