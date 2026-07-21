#!/bin/sh
# Assert sc-planning/SKILL.md explicitly rejects soft "prefer ≤7" and any
# 8-9 / 8–9 exception band. Hard Must alone is not enough.
set -e

ROOT="${1:-.}"
FILE="$ROOT/.cursor/skills/sc-planning/SKILL.md"

if [ ! -f "$FILE" ]; then
  echo "FAIL: missing $FILE"
  exit 1
fi

# Explicit rejection of soft "prefer ≤7" / "prefer <=7" (not merely hard Must).
if ! grep -Eiq 'reject.*prefer[[:space:]]*(≤|<=)[[:space:]]*7|prefer[[:space:]]*(≤|<=)[[:space:]]*7.*(reject|must[[:space:]]+not|forbidden|disallowed)|soft[[:space:]]+prefer[[:space:]]*(≤|<=)[[:space:]]*7|no[[:space:]]+soft.*"prefer[[:space:]]*(≤|<=)[[:space:]]*7"|exclude.*prefer[[:space:]]*(≤|<=)[[:space:]]*7|Must[[:space:]]+not:.*prefer[[:space:]]*(≤|<=)[[:space:]]*7' "$FILE"; then
  echo "FAIL: $FILE must explicitly reject soft prefer ≤7 (or prefer <=7)"
  exit 1
fi

# Explicit rejection of any 8-9 / 8–9 exception band.
if ! grep -Eiq 'reject.*(8[-–]9|8[[:space:]]*[-–][[:space:]]*9).*(exception|band)|(8[-–]9|8[[:space:]]*[-–][[:space:]]*9).*(exception[[:space:]]+band).*(reject|must[[:space:]]+not|forbidden|disallowed|no)|(no|reject|exclude|must[[:space:]]+not).*(8[-–]9|8[[:space:]]*[-–][[:space:]]*9).*(exception|band)|Must[[:space:]]+not:.*(8[-–]9|8[[:space:]]*[-–][[:space:]]*9)' "$FILE"; then
  echo "FAIL: $FILE must explicitly reject any 8-9 / 8–9 exception band"
  exit 1
fi

echo "ok: sc-planning explicitly rejects soft prefer ≤7 and any 8-9 exception band"
exit 0
