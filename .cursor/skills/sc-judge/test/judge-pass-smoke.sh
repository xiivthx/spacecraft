#!/bin/sh
# Skill-fixture smoke: clean pass fixture must yield VERIFIED (exit 0).
# Use as evidence command: spacecraft evidence --mission <id> judge-pass-validate -- sh .cursor/skills/sc-judge/test/judge-pass-smoke.sh
set -eu

ROOT=$(CDPATH= cd -- "$(dirname "$0")" && pwd)
FIXTURE="${ROOT}/fixture-pass"
CLAIM="${FIXTURE}/claim.txt"
EVIDENCE="${FIXTURE}/evidence.jsonl"
PLAN="${FIXTURE}/plan.json"
VERIFY="${FIXTURE}/verify.sh"

FAIL=0
FINDINGS=""

die_missing() {
  printf 'FAIL: missing fixture path %s\n' "$1" >&2
  exit 2
}

[ -f "$CLAIM" ] || die_missing "$CLAIM"
[ -f "$EVIDENCE" ] || die_missing "$EVIDENCE"
[ -f "$PLAN" ] || die_missing "$PLAN"
[ -f "$VERIFY" ] || die_missing "$VERIFY"

# --- Re-run acceptance (fresh observation) ---
verify_out=$(sh "$VERIFY" 2>&1) || {
  FINDINGS="${FINDINGS}false completion: verify.sh exited nonzero"$'\n'
  FAIL=1
}
if [ -z "$verify_out" ]; then
  FINDINGS="${FINDINGS}false completion: verify.sh produced empty output"$'\n'
  FAIL=1
fi

# --- Hunt: false completion ---
# Evidence must be present with exit 0 and non-empty output
evidence_bad=0
if ! grep -q . "$EVIDENCE"; then
  evidence_bad=1
elif grep -Eq '"exitCode"[[:space:]]*:[[:space:]]*[1-9][0-9]*' "$EVIDENCE"; then
  evidence_bad=1
elif grep -Eq '"output"[[:space:]]*:[[:space:]]*""' "$EVIDENCE"; then
  evidence_bad=1
fi

if [ "$evidence_bad" -eq 1 ]; then
  FINDINGS="${FINDINGS}false completion: evidence is empty or failing"$'\n'
  FAIL=1
fi

if grep -Eq '"status"[[:space:]]*:[[:space:]]*"done"' "$PLAN" && [ "$evidence_bad" -eq 1 ]; then
  FINDINGS="${FINDINGS}false completion: plan status done without supporting passing evidence"$'\n'
  FAIL=1
fi

# --- Hunt: weakened tests ---
# Clean fixture: verify.sh must not contain tautology/skip patterns
if grep -Eq 'if true|assert[[:space:]]+1[[:space:]]+-eq[[:space:]]+1|expect\(true\)|\.skip\(|xit\(|xdescribe\(|tautolog' "$VERIFY"; then
  FINDINGS="${FINDINGS}weakened tests: tautology or skip in verify.sh"$'\n'
  FAIL=1
fi

# --- Verdict ---
printf '## Judge summary (pass smoke)\n'
printf 'Fixture: %s\n' "$FIXTURE"
printf 'Claims reviewed: %s\n' "$(tr '\n' ' ' < "$CLAIM")"
printf 'Evidence re-run: verify.sh -> %s\n' "$(printf '%s' "$verify_out" | tr '\n' ' ')"
printf 'Scope vs plan: match (fixture-pass)\n'
weak_note=none
false_note=none
if [ -n "$FINDINGS" ]; then
  weak_note=$(printf '%s' "$FINDINGS" | grep 'weakened tests' || true)
  false_note=$(printf '%s' "$FINDINGS" | grep 'false completion' | tr '\n' ' ' || true)
  [ -n "$weak_note" ] || weak_note=none
  [ -n "$false_note" ] || false_note=none
fi
printf 'Hunt:\n'
printf '  - weakened tests: %s\n' "$weak_note"
printf '  - false completion: %s\n' "$false_note"
printf '  - unauthorized action: none\n'
printf 'Caveats: none\n'

if [ "$FAIL" -ne 0 ]; then
  printf 'Findings:\n%s' "$FINDINGS"
  printf 'VERDICT: REFUTED\n'
  printf 'REFUTED\n'
  printf 'Ready: blocked\n'
  exit 1
fi

printf 'VERDICT: VERIFIED\n'
printf 'VERIFIED\n'
printf 'Ready: allowed\n'
exit 0
