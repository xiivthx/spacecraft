#!/bin/sh
# Assert sc-planning/SKILL.md documents same-mission phase split via
# plan-phaseN.json (or plan-phase<N>.json). Vague "Phase 1 / Phase 2"
# or bare plan-phase filenames without same-mission framing is not enough.
set -e

ROOT="${1:-.}"
FILE="$ROOT/.cursor/skills/sc-planning/SKILL.md"

if [ ! -f "$FILE" ]; then
  echo "FAIL: missing $FILE"
  exit 1
fi

# Require explicit plan-phaseN.json naming (not only "Phase 1, Phase 2").
if ! grep -Eq 'plan-phaseN\.json|plan-phase<N>\.json|plan-phase[0-9]+\.json' "$FILE"; then
  echo "FAIL: $FILE must name plan-phaseN.json (or plan-phase<N>.json / plan-phase1.json)"
  exit 1
fi

# Require explicit same-mission framing for the phase-split escape hatch.
if ! grep -Eiq 'same[- ]mission' "$FILE"; then
  echo "FAIL: $FILE must document same-mission phase split (not vague Phase 1 / Phase 2 alone)"
  exit 1
fi

echo "ok: sc-planning documents same-mission phase split via plan-phaseN.json"
exit 0
