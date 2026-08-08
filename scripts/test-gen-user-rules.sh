#!/bin/sh
# test-gen-user-rules.sh - RED/GREEN test for scripts/gen-user-rules.sh (T1).
#
# Runs the generator over the alwaysApply rule sources
# (000/025/026/027/050/100/200) and asserts:
#   1. output body contains marker text from all seven sources
#   2. output strips rule frontmatter (no leading '---' fences, no
#      'alwaysApply:' / 'description:' lines)
#
# Usage: sh scripts/test-gen-user-rules.sh <repo-root>
# Standalone: does not require `make install-global` or a built binary.
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

# --- Acceptance 1: marker text from all seven sources present in the body ---
missing=0
for marker in 'Spacecraft' 'English prompt coach' 'Intent coach' 'Thai + simple English HIL' 'Coding Standards' 'Project Structure' 'Lane Detection'; do
  if ! grep -qF "$marker" "$out"; then
    echo "FAIL: missing marker from output body: $marker"
    missing=1
  fi
done
if [ "$missing" -ne 0 ]; then
  exit 1
fi
echo "  ok   output body contains all seven source markers"

# --- Acceptance 2: rule frontmatter stripped ---
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
