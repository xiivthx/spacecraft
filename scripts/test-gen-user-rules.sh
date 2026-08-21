#!/bin/sh
# test-gen-user-rules.sh - RED/GREEN test for scripts/gen-user-rules.sh (T1).
#
# Runs the generator over CORE source (010-hard-contract) and asserts:
#   1. output body contains CORE markers
#   2. output strips rule frontmatter
#   3. output stays short (always-on budget)
#
# Usage: sh scripts/test-gen-user-rules.sh <repo-root>
set -e

ROOT="${1:?usage: test-gen-user-rules.sh <repo-root>}"
GEN="$ROOT/scripts/gen-user-rules.sh"
RULES_DIR="$ROOT/.cursor/rules"

mkdir -p "$ROOT/.tmp"
tmp=$(mktemp -d "$ROOT/.tmp/gen-user-rules-test.XXXXXX")
cleanup() { rm -rf "$tmp"; }
trap cleanup EXIT INT TERM

out="$tmp/USER-RULES.txt"

echo "gen-user-rules test in $tmp"

if [ ! -f "$GEN" ]; then
  echo "FAIL: $GEN does not exist yet"
  exit 1
fi

if ! sh "$GEN" "$RULES_DIR" "$out"; then
  echo "FAIL: gen-user-rules.sh exited non-zero"
  exit 1
fi

if [ ! -f "$out" ]; then
  echo "FAIL: gen-user-rules.sh did not write $out"
  exit 1
fi

missing=0
for marker in 'hard contract' 'AUTH:' 'INTENT:' 'block-secrets-read' 'HIL language'; do
  if ! grep -qiF "$marker" "$out"; then
    echo "FAIL: missing marker from output body: $marker"
    missing=1
  fi
done
if [ "$missing" -ne 0 ]; then
  exit 1
fi
echo "  ok   output body contains CORE markers"

lines=$(wc -l < "$out" | tr -d ' ')
if [ "$lines" -gt 60 ]; then
  echo "FAIL: USER-RULES CORE too long ($lines lines; want <=60)"
  exit 1
fi
echo "  ok   CORE length <=60 lines ($lines)"

if grep -n '^---$' "$out" >/dev/null; then
  echo "FAIL: output still contains a '---' frontmatter fence line"
  grep -n '^---$' "$out"
  exit 1
fi
echo "  ok   no leading '---' fences in output"

if grep -n '^alwaysApply:' "$out" >/dev/null; then
  echo "FAIL: output still contains an 'alwaysApply:' frontmatter line"
  grep -n '^alwaysApply:' "$out"
  exit 1
fi
if grep -n '^description:' "$out" >/dev/null; then
  echo "FAIL: output still contains a 'description:' frontmatter line"
  grep -n '^description:' "$out"
  exit 1
fi
echo "  ok   no alwaysApply:/description: frontmatter lines in output"

echo "PASS: gen-user-rules.sh (T1 acceptances)"
