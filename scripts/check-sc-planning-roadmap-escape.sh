#!/bin/sh
# Assert sc-planning/SKILL.md documents roadmap/multi-mission split via
# spacecraft map as the other escape hatch (alongside same-mission phases).
# Vague "roadmap" alone, or map.json / Map integration, is not enough.
set -e

ROOT="${1:-.}"
FILE="$ROOT/.cursor/skills/sc-planning/SKILL.md"

if [ ! -f "$FILE" ]; then
  echo "FAIL: missing $FILE"
  exit 1
fi

# Require the CLI escape hatch name (not outputs/map.json / "Map integration").
if ! grep -Eq 'spacecraft[[:space:]]+map' "$FILE"; then
  echo "FAIL: $FILE must name spacecraft map (CLI), not map.json alone"
  exit 1
fi

# Require roadmap or multi-mission framing for the split.
if ! grep -Eiq 'multi[- ]mission|roadmap' "$FILE"; then
  echo "FAIL: $FILE must document roadmap/multi-mission split"
  exit 1
fi

# Require escape-hatch framing (the other hatch alongside phases).
if ! grep -Eiq 'escape[[:space:]]+hatch' "$FILE"; then
  echo "FAIL: $FILE must frame roadmap/multi-mission split via spacecraft map as an escape hatch"
  exit 1
fi

echo "ok: sc-planning documents roadmap/multi-mission split via spacecraft map as escape hatch"
exit 0
