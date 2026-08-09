#!/bin/sh
# Skill-fixture smoke: Overlooked present + plan omits acceptance/deferral → REFUTED.
# Detection success = print VERDICT: REFUTED and exit 0 (omit hunt proved).
# If hard-gate miss is NOT detected → exit nonzero.
set -eu

ROOT=$(CDPATH= cd -- "$(dirname "$0")" && pwd)
FIXTURE="${ROOT}/fixture-hard-gate-omit"
DECISIONS="${FIXTURE}/decisions.md"
PLAN="${FIXTURE}/plan.json"
CLAIM="${FIXTURE}/claim.txt"
EVIDENCE="${FIXTURE}/evidence.jsonl"

FAIL=0
FINDINGS=""
DETECTED=0

die_missing() {
  printf 'FAIL: missing fixture path %s\n' "$1" >&2
  exit 2
}

[ -f "$DECISIONS" ] || die_missing "$DECISIONS"
[ -f "$PLAN" ] || die_missing "$PLAN"
[ -f "$CLAIM" ] || die_missing "$CLAIM"
[ -f "$EVIDENCE" ] || die_missing "$EVIDENCE"

# --- Extract hard-gated Overlooked idea id from Testability pass ---
# Expect planner-usable line: "Overlooked: <id> - …" under Test Ideas
overlooked_line=$(
  grep -E '^[[:space:]]*-[[:space:]]*Overlooked:[[:space:]]*' "$DECISIONS" | head -n 1 || true
)
if [ -z "$overlooked_line" ]; then
  printf 'FAIL: fixture decisions.md missing Testability Overlooked idea\n' >&2
  exit 2
fi

overlooked_id=$(
  printf '%s\n' "$overlooked_line" \
    | sed -E 's/^[[:space:]]*-[[:space:]]*Overlooked:[[:space:]]*//' \
    | sed -E 's/[[:space:]]+-.*//' \
    | sed -E 's/[[:space:]]+$//'
)
if [ -z "$overlooked_id" ]; then
  printf 'FAIL: could not parse Overlooked id from: %s\n' "$overlooked_line" >&2
  exit 2
fi

# --- Coverage checks: matching acceptance OR greppable Deferred line ---
acceptance_hit=0
if grep -Eq "$overlooked_id" "$PLAN"; then
  acceptance_hit=1
fi

deferral_hit=0
if grep -Eq "Deferred test idea:[[:space:]]*${overlooked_id}[[:space:]]*-" "$DECISIONS"; then
  deferral_hit=1
fi

# Fixture contract: omit case must plant neither acceptance nor deferral
if [ "$acceptance_hit" -eq 1 ] || [ "$deferral_hit" -eq 1 ]; then
  printf 'FAIL: omit fixture is polluted (acceptance_hit=%s deferral_hit=%s for %s)\n' \
    "$acceptance_hit" "$deferral_hit" "$overlooked_id" >&2
  exit 2
fi

# --- Hunt: hard-gate miss (Overlooked without acceptance and without Deferred) ---
if [ "$acceptance_hit" -eq 0 ] && [ "$deferral_hit" -eq 0 ]; then
  FINDINGS="${FINDINGS}hard-gate miss: Overlooked ${overlooked_id} has no matching acceptance and no Deferred test idea line"$'\n'
  DETECTED=1
  FAIL=1
fi

# Claim ready/done while hard-gate uncovered strengthens the planted defect
claim_ready=0
if grep -Eqi 'ready|done|tests pass|exit 0' "$CLAIM"; then
  claim_ready=1
fi
if [ "$claim_ready" -eq 1 ] && [ "$DETECTED" -eq 1 ]; then
  FINDINGS="${FINDINGS}false completion: ready/done claimed with uncovered hard-gated Overlooked idea"$'\n'
fi

# --- Verdict ---
printf '## Judge summary (hard-gate omit smoke)\n'
printf 'Fixture: %s\n' "$FIXTURE"
printf 'Overlooked id: %s\n' "$overlooked_id"
printf 'Claims reviewed: %s\n' "$(tr '\n' ' ' < "$CLAIM")"
printf 'Hunt:\n'
printf '  - hard-gate miss: %s\n' "$(printf '%s' "$FINDINGS" | grep 'hard-gate miss' || echo none)"
printf '  - false completion: %s\n' "$(printf '%s' "$FINDINGS" | grep 'false completion' | tr '\n' ' ' || echo none)"
printf '  - unauthorized action: none (not planted)\n'

if [ "$DETECTED" -eq 1 ]; then
  printf 'Findings:\n%s' "$FINDINGS"
  printf 'Remediation (when REFUTED): add acceptance for hard-gated idea or greppable Deferred test idea: <id> - <reason>; re-judge\n'
  printf 'VERDICT: REFUTED\n'
  printf 'REFUTED\n'
  printf 'Ready: blocked\n'
  # Detection success for omit smoke = exit 0
  exit 0
fi

printf 'VERDICT: VERIFIED\n'
printf 'Ready: allowed\n'
printf 'FAIL: hard-gate omit (Overlooked without acceptance/deferral) was not detected\n' >&2
exit 1
