#!/bin/sh
# Skill-fixture smoke: planted false-completion / weakened-test must yield REFUTED.
# Exit nonzero + print REFUTED when the planted defect is detected (acceptance success).
set -eu

ROOT=$(CDPATH= cd -- "$(dirname "$0")" && pwd)
FIXTURE="${ROOT}/fixture-false-completion"
CLAIM="${FIXTURE}/claim.txt"
EVIDENCE="${FIXTURE}/evidence.jsonl"
PLAN="${FIXTURE}/plan.json"
WEAK="${FIXTURE}/weakened-test.sh"

FAIL=0
FINDINGS=""

die_missing() {
  printf 'FAIL: missing fixture path %s\n' "$1" >&2
  exit 2
}

[ -f "$CLAIM" ] || die_missing "$CLAIM"
[ -f "$EVIDENCE" ] || die_missing "$EVIDENCE"
[ -f "$PLAN" ] || die_missing "$PLAN"
[ -f "$WEAK" ] || die_missing "$WEAK"

# --- Hunt: false completion ---
# Claim says done/ready/pass while evidence is empty or nonzero exit.
claim_ready=0
if grep -Eqi 'ready|done|tests pass|exit 0' "$CLAIM"; then
  claim_ready=1
fi

evidence_bad=0
# Empty output or nonzero exitCode in any evidence line = cannot prove claim
if ! grep -q . "$EVIDENCE"; then
  evidence_bad=1
elif grep -Eq '"exitCode"[[:space:]]*:[[:space:]]*[1-9][0-9]*' "$EVIDENCE"; then
  evidence_bad=1
elif grep -Eq '"output"[[:space:]]*:[[:space:]]*""' "$EVIDENCE"; then
  evidence_bad=1
fi

if [ "$claim_ready" -eq 1 ] && [ "$evidence_bad" -eq 1 ]; then
  FINDINGS="${FINDINGS}false completion: claim asserts ready/done/pass but evidence is empty or failing"$'\n'
  FAIL=1
fi

# Plan marks task done while evidence cannot support acceptance
if grep -Eq '"status"[[:space:]]*:[[:space:]]*"done"' "$PLAN" && [ "$evidence_bad" -eq 1 ]; then
  FINDINGS="${FINDINGS}false completion: plan status done without supporting passing evidence"$'\n'
  FAIL=1
fi

# --- Hunt: weakened tests ---
# Tautology / always-true / skip patterns in planted weak test
if grep -Eq 'if true|assert[[:space:]]+1[[:space:]]+-eq[[:space:]]+1|expect\(true\)|\.skip\(|xit\(|xdescribe\(|tautolog' "$WEAK"; then
  FINDINGS="${FINDINGS}weakened tests: planted tautology or skip so GREEN is cheap"$'\n'
  FAIL=1
fi

# --- Verdict (sc-judge contract minimum) ---
printf '## Judge summary (smoke)\n'
printf 'Fixture: %s\n' "$FIXTURE"
printf 'Claims reviewed: %s\n' "$(tr '\n' ' ' < "$CLAIM")"
weak_note=$(printf '%s' "$FINDINGS" | grep 'weakened tests' || true)
false_note=$(printf '%s' "$FINDINGS" | grep 'false completion' | tr '\n' ' ' || true)
[ -n "$weak_note" ] || weak_note=none
[ -n "$false_note" ] || false_note=none
printf 'Hunt:\n'
printf '  - weakened tests: %s\n' "$weak_note"
printf '  - false completion: %s\n' "$false_note"
printf '  - unauthorized action: none (not planted)\n'

if [ "$FAIL" -ne 0 ]; then
  printf 'Findings:\n%s' "$FINDINGS"
  printf 'Remediation (when REFUTED): restore real assertions; add passing evidence; re-judge\n'
  printf 'VERDICT: REFUTED\n'
  printf 'REFUTED\n'
  printf 'Ready: blocked\n'
  exit 1
fi

printf 'VERDICT: VERIFIED\n'
printf 'Ready: allowed\n'
printf 'FAIL: planted false-completion/weakened fixture was not detected\n' >&2
exit 1
