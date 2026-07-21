#!/bin/sh
# Assert 200-workflow.mdc states ≤7 jigsaw tasks per phase as a hard Must
# (not preference-only wording). Bare ">7 tasks" roadmap wording is not enough.
set -e

ROOT="${1:-.}"
FILE="$ROOT/.cursor/rules/200-workflow.mdc"

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
  echo "FAIL: $FILE must state ≤7 jigsaw tasks per phase as a hard Must (not preference-only)"
  echo "      found bare Must/≤7/>7 wording is not sufficient"
  exit 1
fi

# Explicit non-preference / non-soft stance (acceptance: not preference-only).
if ! grep -Eiq 'not[[:space:]]+(a[[:space:]]+)?preference|not[[:space:]]+preference-only|not[[:space:]]+soft' "$FILE"; then
  echo "FAIL: $FILE must reject preference-only / soft wording for the ≤7 jigsaw cap"
  exit 1
fi

echo "ok: 200-workflow.mdc states ≤7 jigsaw tasks per phase as hard Must (not preference-only)"
exit 0
