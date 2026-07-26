#!/bin/sh
# Assert sc-planning/SKILL.md hands multi-mission splits to /sc-discuss +
# mission-sizing (discuss owns spacecraft map). Planning must not create maps.
# Vague "roadmap" alone, or map.json / Map integration, is not enough.
set -e

ROOT="${1:-.}"
FILE="$ROOT/.cursor/skills/sc-planning/SKILL.md"

if [ ! -f "$FILE" ]; then
  echo "FAIL: missing $FILE"
  exit 1
fi

# Require handoff to discuss for multi-mission / independent seams.
if ! grep -Eq '/sc-discuss' "$FILE"; then
  echo "FAIL: $FILE must hand multi-mission splits to /sc-discuss"
  exit 1
fi

if ! grep -Eq 'mission-sizing' "$FILE"; then
  echo "FAIL: $FILE must name mission-sizing"
  exit 1
fi

# Require multi-mission framing (not map.json / Map integration alone).
if ! grep -Eiq 'multi[- ]mission' "$FILE"; then
  echo "FAIL: $FILE must document multi-mission handoff"
  exit 1
fi

# Require explicit ban on planning-owned map create.
if ! grep -Eq 'map new' "$FILE"; then
  echo "FAIL: $FILE must mention map new (as a Must-not for planning)"
  exit 1
fi

if ! grep -Eiq 'must not.*map new|never.*map new|do not.*map new' "$FILE"; then
  echo "FAIL: $FILE must forbid planning-owned map new"
  exit 1
fi

echo "ok: sc-planning hands multi-mission to /sc-discuss + mission-sizing; forbids planning map new"
exit 0
