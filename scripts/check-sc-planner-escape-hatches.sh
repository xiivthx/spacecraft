#!/bin/sh
# Assert sc-planner.md documents both escape hatches when scope exceeds ≤7:
# (1) same-mission phase split via plan-phaseN.json (or equivalent explicit naming)
# (2) roadmap/multi-mission via spacecraft map
# Vague "Phase 1 / Phase 2" alone is not enough.
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

# Require same-mission framing for the phase-split escape hatch.
if ! grep -Eiq 'same[- ]mission' "$FILE"; then
  echo "FAIL: $FILE must document same-mission phase split (not vague Phase 1 / Phase 2 alone)"
  exit 1
fi

# Require the CLI escape hatch name (not map.json alone).
if ! grep -Eq 'spacecraft[[:space:]]+map' "$FILE"; then
  echo "FAIL: $FILE must name spacecraft map (CLI) for roadmap/multi-mission split"
  exit 1
fi

# Require roadmap or multi-mission framing.
if ! grep -Eiq 'multi[- ]mission|roadmap' "$FILE"; then
  echo "FAIL: $FILE must document roadmap/multi-mission split"
  exit 1
fi

# Require escape-hatch framing for both paths.
if ! grep -Eiq 'escape[[:space:]]+hatch' "$FILE"; then
  echo "FAIL: $FILE must frame phase split and roadmap/multi-mission as escape hatches"
  exit 1
fi

echo "ok: sc-planner documents phase-split and roadmap/multi-mission escape hatches"
exit 0
