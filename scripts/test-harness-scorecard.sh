#!/bin/sh
# test-harness-scorecard.sh - RED/GREEN for scripts/harness-scorecard.sh (T2+T3).
#
# Frozen oracles (approved-scenarios S1, S2, S10):
#   S1  clean run → exit 0; SCORECARD <dim> pass for all four required dims
#   S2  HARNESS_SCORECARD_FORCE_FAIL=false-completion → exit ≠0;
#       SCORECARD false-completion fail
#   S10 install-smoke invokes smoke.sh with project dir = repo ROOT
#
# Usage: sh scripts/test-harness-scorecard.sh <repo-root> [spacecraft-binary]
set -e

ROOT="${1:?usage: test-harness-scorecard.sh <repo-root> [spacecraft-binary]}"
ROOT="$(cd "$ROOT" && pwd)"
BIN="${2:-}"
RUNNER="$ROOT/scripts/harness-scorecard.sh"

if [ -z "$BIN" ]; then
  if [ -x "$ROOT/spacecraft" ]; then
    BIN="$ROOT/spacecraft"
  elif command -v spacecraft >/dev/null 2>&1; then
    BIN="$(command -v spacecraft)"
  else
    echo "FAIL: spacecraft binary required (pass as arg2 or build ./spacecraft)"
    exit 1
  fi
fi
case "$BIN" in
  /*) ;;
  *) BIN="$(cd "$(dirname "$BIN")" && pwd)/$(basename "$BIN")" ;;
esac

echo "harness-scorecard test ROOT=$ROOT BIN=$BIN"

if [ ! -f "$RUNNER" ]; then
  echo "FAIL: $RUNNER does not exist yet"
  exit 1
fi

# --- S1 / T2-a: clean green SCORECARD lines + exit 0 ---
# Run from an unrelated cwd so S10 can prove PROJECT=ROOT (not caller cwd).
mkdir -p "$ROOT/.tmp"
tmp=$(mktemp -d "$ROOT/.tmp/harness-scorecard-test.XXXXXX")
cleanup() { rm -rf "$tmp"; }
trap cleanup EXIT INT TERM

out_clean="$tmp/clean.out"
rc_clean=0
(
  cd "$tmp"
  # shellcheck disable=SC2030
  sh "$RUNNER" "$ROOT" "$BIN"
) >"$out_clean" 2>&1 || rc_clean=$?

if [ "$rc_clean" -ne 0 ]; then
  echo "FAIL: clean harness-scorecard.sh exited $rc_clean (want 0)"
  echo "--- stdout/stderr ---"
  cat "$out_clean"
  exit 1
fi

missing=0
for line in \
  'SCORECARD install-smoke pass' \
  'SCORECARD false-completion pass' \
  'SCORECARD judge-skill pass' \
  'SCORECARD process-grammar pass'
do
  if ! grep -Fq "$line" "$out_clean"; then
    echo "FAIL: clean stdout missing: $line"
    missing=1
  fi
done
if [ "$missing" -ne 0 ]; then
  echo "--- stdout/stderr ---"
  cat "$out_clean"
  exit 1
fi
echo "  ok   S1 clean exit 0 + SCORECARD * pass lines"

# --- S10 / T2-b: smoke project dir = ROOT (harness source tree) ---
# smoke.sh prints "Smoke check: $TARGET" when invoked; TARGET must be ROOT.
if ! grep -Fq "Smoke check: $ROOT" "$out_clean"; then
  echo "FAIL: S10 install-smoke did not invoke smoke.sh with project dir=$ROOT"
  echo "      expected stdout to include: Smoke check: $ROOT"
  echo "--- stdout/stderr ---"
  cat "$out_clean"
  exit 1
fi
echo "  ok   S10 install-smoke smoke.sh project dir = ROOT"

# --- S2 / T3-a: force-fail false-completion → non-zero + fail line ---
out_fail="$tmp/force-fail.out"
rc_fail=0
(
  cd "$tmp"
  HARNESS_SCORECARD_FORCE_FAIL=false-completion \
    sh "$RUNNER" "$ROOT" "$BIN"
) >"$out_fail" 2>&1 || rc_fail=$?

if [ "$rc_fail" -eq 0 ]; then
  echo "FAIL: HARNESS_SCORECARD_FORCE_FAIL=false-completion exited 0 (want ≠0)"
  echo "--- stdout/stderr ---"
  cat "$out_fail"
  exit 1
fi

if ! grep -Fq 'SCORECARD false-completion fail' "$out_fail"; then
  echo "FAIL: force-fail stdout missing: SCORECARD false-completion fail"
  echo "--- stdout/stderr ---"
  cat "$out_fail"
  exit 1
fi
echo "  ok   S2 force-fail exit $rc_fail + SCORECARD false-completion fail"

echo "PASS: harness-scorecard.sh (T2+T3 acceptances S1/S2/S10)"
