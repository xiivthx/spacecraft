#!/bin/sh
# Skill-fixture smoke: Neg+Overlooked acceptances + matching fresh evidence → VERIFIED (exit 0).
# Use as evidence command: spacecraft evidence --mission <id> judge-hard-gate-pass-verified -- sh .cursor/skills/sc-judge/test/hard-gate-pass-smoke.sh
set -eu

ROOT=$(CDPATH= cd -- "$(dirname "$0")" && pwd)
FIXTURE="${ROOT}/fixture-hard-gate-pass"
CLAIM="${FIXTURE}/claim.txt"
EVIDENCE="${FIXTURE}/evidence.jsonl"
PLAN="${FIXTURE}/plan.json"
DECISIONS="${FIXTURE}/decisions.md"
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
[ -f "$DECISIONS" ] || die_missing "$DECISIONS"
[ -f "$VERIFY" ] || die_missing "$VERIFY"

# --- Extract hard-gated Negative + Overlooked idea ids from Testability pass ---
parse_idea_id() {
  # stdin: "- Negative: <id> - …" or "- Overlooked: <id> - …"
  sed -E 's/^[[:space:]]*-[[:space:]]*(Negative|Overlooked):[[:space:]]*//' \
    | sed -E 's/[[:space:]]+-.*//' \
    | sed -E 's/[[:space:]]+$//'
}

neg_line=$(
  grep -E '^[[:space:]]*-[[:space:]]*Negative:[[:space:]]*' "$DECISIONS" | head -n 1 || true
)
overlooked_line=$(
  grep -E '^[[:space:]]*-[[:space:]]*Overlooked:[[:space:]]*' "$DECISIONS" | head -n 1 || true
)

if [ -z "$neg_line" ] || [ -z "$overlooked_line" ]; then
  printf 'FAIL: fixture decisions.md must include Testability Negative and Overlooked ideas\n' >&2
  exit 2
fi

neg_id=$(printf '%s\n' "$neg_line" | parse_idea_id)
overlooked_id=$(printf '%s\n' "$overlooked_line" | parse_idea_id)

if [ -z "$neg_id" ] || [ -z "$overlooked_id" ]; then
  printf 'FAIL: could not parse Neg/Overlooked ids (neg=%s overlooked=%s)\n' \
    "$neg_id" "$overlooked_id" >&2
  exit 2
fi

# --- Coverage: each hard-gated idea has matching acceptance (pass fixture uses acceptances) ---
for idea_id in "$neg_id" "$overlooked_id"; do
  acceptance_hit=0
  deferral_hit=0
  if grep -Eq "$idea_id" "$PLAN"; then
    acceptance_hit=1
  fi
  if grep -Eq "Deferred test idea:[[:space:]]*${idea_id}[[:space:]]*-" "$DECISIONS"; then
    deferral_hit=1
  fi
  if [ "$acceptance_hit" -eq 0 ] && [ "$deferral_hit" -eq 0 ]; then
    FINDINGS="${FINDINGS}hard-gate miss: ${idea_id} has no matching acceptance and no Deferred test idea line"$'\n'
    FAIL=1
  fi
done

# Pass fixture contract: both hard-gates must be acceptance-covered (not deferral-only)
if ! grep -Eq "$neg_id" "$PLAN" || ! grep -Eq "$overlooked_id" "$PLAN"; then
  FINDINGS="${FINDINGS}false completion: pass fixture requires Neg+Overlooked acceptances in plan.json"$'\n'
  FAIL=1
fi

# --- Re-run acceptance (fresh observation) ---
verify_out=$(sh "$VERIFY" 2>&1) || {
  FINDINGS="${FINDINGS}false completion: verify.sh exited nonzero"$'\n'
  FAIL=1
}
if [ -z "${verify_out:-}" ]; then
  FINDINGS="${FINDINGS}false completion: verify.sh produced empty output"$'\n'
  FAIL=1
fi

# --- Hunt: false completion (evidence must pass) ---
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

# Mapped hard-gate acceptances claimed done require fresh evidence
if grep -Eq "$neg_id|$overlooked_id" "$PLAN" && [ "$evidence_bad" -eq 1 ]; then
  FINDINGS="${FINDINGS}false completion: hard-gated acceptance claimed without fresh passing evidence"$'\n'
  FAIL=1
fi

# --- Hunt: weakened tests ---
if grep -Eq 'if true|assert[[:space:]]+1[[:space:]]+-eq[[:space:]]+1|expect\(true\)|\.skip\(|xit\(|xdescribe\(|tautolog' "$VERIFY"; then
  FINDINGS="${FINDINGS}weakened tests: tautology or skip in verify.sh"$'\n'
  FAIL=1
fi

# --- Verdict ---
printf '## Judge summary (hard-gate pass smoke)\n'
printf 'Fixture: %s\n' "$FIXTURE"
printf 'Negative id: %s\n' "$neg_id"
printf 'Overlooked id: %s\n' "$overlooked_id"
printf 'Claims reviewed: %s\n' "$(tr '\n' ' ' < "$CLAIM")"
printf 'Evidence re-run: verify.sh -> %s\n' "$(printf '%s' "$verify_out" | tr '\n' ' ')"
printf 'Scope vs plan: match (fixture-hard-gate-pass Neg+Overlooked acceptances)\n'
weak_note=none
false_note=none
hard_note=none
if [ -n "$FINDINGS" ]; then
  weak_note=$(printf '%s' "$FINDINGS" | grep 'weakened tests' || true)
  false_note=$(printf '%s' "$FINDINGS" | grep 'false completion' | tr '\n' ' ' || true)
  hard_note=$(printf '%s' "$FINDINGS" | grep 'hard-gate miss' | tr '\n' ' ' || true)
  [ -n "$weak_note" ] || weak_note=none
  [ -n "$false_note" ] || false_note=none
  [ -n "$hard_note" ] || hard_note=none
fi
printf 'Hunt:\n'
printf '  - hard-gate miss: %s\n' "$hard_note"
printf '  - weakened tests: %s\n' "$weak_note"
printf '  - false completion: %s\n' "$false_note"
printf '  - unauthorized action: none\n'
printf 'Remediation (when REFUTED): none\n'

if [ "$FAIL" -ne 0 ]; then
  printf 'Findings:\n%s' "$FINDINGS"
  printf 'Remediation (when REFUTED): cover hard-gated Neg/Overlooked with acceptance or Deferred line; ensure fresh evidence; re-judge\n'
  printf 'VERDICT: REFUTED\n'
  printf 'REFUTED\n'
  printf 'Ready: blocked\n'
  exit 1
fi

printf 'VERDICT: VERIFIED\n'
printf 'VERIFIED\n'
printf 'Ready: allowed\n'
exit 0
