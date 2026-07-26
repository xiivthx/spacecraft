#!/bin/sh
# Assert 200-workflow.mdc documents both paths when scope exceeds ≤7:
# (1) same-mission phase split via plan-phaseN.json
# (2) multi-mission via mission-sizing (discuss owns spacecraft map)
set -e

ROOT="${1:-.}"
FILE="$ROOT/.cursor/rules/200-workflow.mdc"

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

# Require mission-sizing for multi-mission (discuss-owned).
if ! grep -Eq 'mission-sizing' "$FILE"; then
  echo "FAIL: $FILE must name mission-sizing for multi-mission"
  exit 1
fi

if ! grep -Eiq 'multi[- ]mission' "$FILE"; then
  echo "FAIL: $FILE must document multi-mission split"
  exit 1
fi

# Require discuss owns map; planning must not map new.
if ! grep -Eq 'spacecraft[[:space:]]+map|map new' "$FILE"; then
  echo "FAIL: $FILE must name spacecraft map / map new (as discuss-owned plumbing)"
  exit 1
fi

if ! grep -Eiq 'discuss owns|/sc-discuss' "$FILE"; then
  echo "FAIL: $FILE must say discuss owns map or hand to /sc-discuss"
  exit 1
fi

if ! grep -Eiq 'must not.*map new|never.*map new|planning must not' "$FILE"; then
  echo "FAIL: $FILE must forbid planning-owned map new"
  exit 1
fi

echo "ok: 200-workflow.mdc documents phase-split and discuss-owned multi-mission"
exit 0
