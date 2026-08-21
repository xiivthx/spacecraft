#!/bin/sh
# Skill-fixture smoke: Deferred covers overlooked-x + sibling hard-gate gap → REFUTED.
# Detection success = deferred idea covered, overall VERDICT: REFUTED, exit 0.
# If deferral is ignored, or sibling gap is missed → exit nonzero.
set -eu

ROOT=$(CDPATH= cd -- "$(dirname "$0")" && pwd)
FIXTURE="${ROOT}/fixture-hard-gate-deferral"
DECISIONS="${FIXTURE}/decisions.md"
PLAN="${FIXTURE}/plan.json"
CLAIM="${FIXTURE}/claim.txt"
EVIDENCE="${FIXTURE}/evidence.jsonl"

FAIL=0
FINDINGS=""
DEFERRED_COVERED=0
SIBLING_GAP=0

die_missing() {
  printf 'FAIL: missing fixture path %s\n' "$1" >&2
  exit 2
}

[ -f "$DECISIONS" ] || die_missing "$DECISIONS"
[ -f "$PLAN" ] || die_missing "$PLAN"
[ -f "$CLAIM" ] || die_missing "$CLAIM"
[ -f "$EVIDENCE" ] || die_missing "$EVIDENCE"

# --- Deferred idea: overlooked-x must be greppably deferred ---
DEFERRED_ID="overlooked-x"
if ! grep -Eq "Deferred test idea:[[:space:]]*${DEFERRED_ID}[[:space:]]*-" "$DECISIONS"; then
  printf 'FAIL: fixture decisions.md missing greppable Deferred test idea: %s - <reason>\n' \
    "$DEFERRED_ID" >&2
  exit 2
fi

# Fixture contract: deferred idea must NOT also have a plan acceptance (edge = deferral-only)
if grep -Eq "$DEFERRED_ID" "$PLAN"; then
  printf 'FAIL: deferral fixture polluted (plan acceptance mentions %s)\n' "$DEFERRED_ID" >&2
  exit 2
fi

DEFERRED_COVERED=1

# --- Extract sibling hard-gated idea (Negative preferred; else Overlooked != deferred id) ---
sibling_line=$(
  grep -E '^[[:space:]]*-[[:space:]]*Negative:[[:space:]]*' "$DECISIONS" | head -n 1 || true
)
sibling_bucket="Negative"
if [ -z "$sibling_line" ]; then
  sibling_line=$(
    grep -E '^[[:space:]]*-[[:space:]]*Overlooked:[[:space:]]*' "$DECISIONS" \
      | grep -v "Overlooked:[[:space:]]*${DEFERRED_ID}[[:space:]]*-" \
      | head -n 1 || true
  )
  sibling_bucket="Overlooked"
fi
if [ -z "$sibling_line" ]; then
  printf 'FAIL: fixture decisions.md missing sibling hard-gated Negative/Overlooked idea\n' >&2
  exit 2
fi

sibling_id=$(
  printf '%s\n' "$sibling_line" \
    | sed -E "s/^[[:space:]]*-[[:space:]]*${sibling_bucket}:[[:space:]]*//" \
    | sed -E 's/[[:space:]]+-.*//' \
    | sed -E 's/[[:space:]]+$//'
)
if [ -z "$sibling_id" ]; then
  printf 'FAIL: could not parse sibling %s id from: %s\n' "$sibling_bucket" "$sibling_line" >&2
  exit 2
fi
if [ "$sibling_id" = "$DEFERRED_ID" ]; then
  printf 'FAIL: sibling id must differ from deferred id %s\n' "$DEFERRED_ID" >&2
  exit 2
fi

sibling_acceptance=0
if grep -Eq "$sibling_id" "$PLAN"; then
  sibling_acceptance=1
fi
sibling_deferral=0
if grep -Eq "Deferred test idea:[[:space:]]*${sibling_id}[[:space:]]*-" "$DECISIONS"; then
  sibling_deferral=1
fi

# Fixture contract: sibling must remain uncovered
if [ "$sibling_acceptance" -eq 1 ] || [ "$sibling_deferral" -eq 1 ]; then
  printf 'FAIL: deferral fixture polluted (sibling %s acceptance=%s deferral=%s)\n' \
    "$sibling_id" "$sibling_acceptance" "$sibling_deferral" >&2
  exit 2
fi

# --- Hunt: deferred idea covered; sibling hard-gate miss remains ---
FINDINGS="${FINDINGS}deferred coverage: ${DEFERRED_ID} covered by Deferred test idea line (no evidence required)"$'\n'
FINDINGS="${FINDINGS}hard-gate miss: ${sibling_bucket} ${sibling_id} has no matching acceptance and no Deferred test idea line"$'\n'
SIBLING_GAP=1
FAIL=1

# Claim ready/done while sibling hard-gate uncovered strengthens the planted defect
claim_ready=0
if grep -Eqi 'ready|done|tests pass|exit 0' "$CLAIM"; then
  claim_ready=1
fi
if [ "$claim_ready" -eq 1 ] && [ "$SIBLING_GAP" -eq 1 ]; then
  FINDINGS="${FINDINGS}false completion: ready/done claimed with uncovered sibling hard-gated idea"$'\n'
fi

# --- Verdict ---
printf '## Judge summary (hard-gate deferral smoke)\n'
printf 'Fixture: %s\n' "$FIXTURE"
printf 'Deferred id: %s (covered=%s)\n' "$DEFERRED_ID" "$DEFERRED_COVERED"
printf 'Sibling %s id: %s (gap=%s)\n' "$sibling_bucket" "$sibling_id" "$SIBLING_GAP"
printf 'Claims reviewed: %s\n' "$(tr '\n' ' ' < "$CLAIM")"
printf 'Hunt:\n'
printf '  - deferred coverage: %s\n' "$(printf '%s' "$FINDINGS" | grep 'deferred coverage' || echo none)"
printf '  - hard-gate miss: %s\n' "$(printf '%s' "$FINDINGS" | grep 'hard-gate miss' || echo none)"
printf '  - false completion: %s\n' "$(printf '%s' "$FINDINGS" | grep 'false completion' | tr '\n' ' ' || echo none)"
printf '  - unauthorized action: none (not planted)\n'

# Success = deferred covered AND sibling gap ⇒ overall REFUTED, exit 0
if [ "$DEFERRED_COVERED" -eq 1 ] && [ "$SIBLING_GAP" -eq 1 ] && [ "$FAIL" -eq 1 ]; then
  printf 'Findings:\n%s' "$FINDINGS"
  printf 'Remediation (when REFUTED): cover sibling hard-gated idea via acceptance or Deferred test idea: <id> - <reason>; re-judge\n'
  printf 'VERDICT: REFUTED\n'
  printf 'REFUTED\n'
  printf 'Ready: blocked\n'
  # Detection success for deferral-edge smoke = exit 0
  exit 0
fi

printf 'VERDICT: VERIFIED\n'
printf 'Ready: allowed\n'
printf 'FAIL: deferral edge (deferred covered + sibling hard-gate gap → REFUTED) was not detected\n' >&2
exit 1
