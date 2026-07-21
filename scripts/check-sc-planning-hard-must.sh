#!/bin/sh
# Assert sc-planning/SKILL.md states ≤7 tasks per phase as a hard Must
# (not preference-only wording). Bare "**Must**: ≤7" is not enough.
set -e

ROOT="${1:-.}"
FILE="$ROOT/.cursor/skills/sc-planning/SKILL.md"

if [ ! -f "$FILE" ]; then
  echo "FAIL: missing $FILE"
  exit 1
fi

# Cap must still be present (≤7 / <=7 / 7 tasks).
if ! grep -Eq '≤7|<=7|7 tasks' "$FILE"; then
  echo "FAIL: $FILE does not mention ≤7 / 7 tasks per phase"
  exit 1
fi

# Stronger than bare Must: require the phrase "hard Must".
if ! grep -Eiq 'hard[[:space:]]+Must' "$FILE"; then
  echo "FAIL: $FILE must state ≤7 tasks per phase as a hard Must (not preference-only)"
  echo "      found bare Must/≤7 wording is not sufficient"
  exit 1
fi

# Explicit non-preference / non-soft stance (acceptance: not preference-only).
if ! grep -Eiq 'not[[:space:]]+(a[[:space:]]+)?preference|not[[:space:]]+preference-only|not[[:space:]]+soft' "$FILE"; then
  echo "FAIL: $FILE must reject preference-only / soft wording for the ≤7 cap"
  exit 1
fi

echo "ok: sc-planning states ≤7 as hard Must (not preference-only)"
exit 0
