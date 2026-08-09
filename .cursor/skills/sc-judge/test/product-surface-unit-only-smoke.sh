#!/bin/sh
# Skill-fixture smoke: UI/user-visible acceptance + unit-only verify → REFUTED.
# Detection success = print VERDICT: REFUTED and exit 0 (unit-only miss proved).
# If unit-only miss is NOT detected → exit nonzero.
set -eu

ROOT=$(CDPATH= cd -- "$(dirname "$0")" && pwd)
FIXTURE="${ROOT}/fixture-product-surface-unit-only"
DECISIONS="${FIXTURE}/decisions.md"
PLAN="${FIXTURE}/plan.json"
CLAIM="${FIXTURE}/claim.txt"
EVIDENCE="${FIXTURE}/evidence.jsonl"

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

# --- Product-surface markers (exact convention tokens) ---
PRODUCT_SURFACE_RE='verify\.product|[[:space:]]browser[[:space:]]|/browser|[[:space:]]browser$|curl|composition'

# --- Extract UI / user-visible task surface from plan ---
# Fixture contract: title or acceptance claims UI / user-visible / workflow behavior.
ui_surface=0
if grep -Eqi 'user-visible|user visible|[[:space:]]UI[[:space:]]|^[[:space:]]*"[^"]*UI[^"]*"' "$PLAN"; then
  ui_surface=1
fi
if grep -Eqi 'workflow|settings page|visible to users|shows the' "$PLAN"; then
  ui_surface=1
fi

if [ "$ui_surface" -eq 0 ]; then
  printf 'FAIL: fixture plan.json missing UI/user-visible acceptance claim\n' >&2
  exit 2
fi

# Gather verify + acceptance text for marker hunt (plan-wide; single-task fixture)
verify_line=$(
  grep -E '"verify"[[:space:]]*:' "$PLAN" | head -n 1 || true
)
acceptance_blob=$(
  # acceptance strings + verify string are the greppable surfaces per decisions.md
  grep -E '"acceptance"|"verify"|user-visible|UI|workflow' "$PLAN" || true
)

if [ -z "$verify_line" ]; then
  printf 'FAIL: fixture plan.json missing verify field\n' >&2
  exit 2
fi

# --- Marker presence in verify and/or acceptance text ---
marker_hit=0
if printf '%s\n%s\n' "$verify_line" "$acceptance_blob" | grep -Eqi "$PRODUCT_SURFACE_RE"; then
  marker_hit=1
fi
# Also allow exact token verify.product anywhere in plan verify/acceptance region
if printf '%s\n' "$verify_line" "$acceptance_blob" | grep -Eq 'verify\.product'; then
  marker_hit=1
fi

# Fixture contract: unit-only miss must plant no product-surface marker
if [ "$marker_hit" -eq 1 ]; then
  printf 'FAIL: unit-only fixture is polluted (product-surface marker present)\n' >&2
  exit 2
fi

# Unit-only smell: verify looks like unit/test runner without product-surface path
unit_only=0
if printf '%s\n' "$verify_line" | grep -Eqi 'npm test|vitest|jest|[[:space:]]test[[:space:]]|-- .*\.test\.|unit'; then
  unit_only=1
fi

# --- Hunt: user-visible claim + no product-surface marker (+ unit-only verify) ---
if [ "$ui_surface" -eq 1 ] && [ "$marker_hit" -eq 0 ]; then
  FINDINGS="${FINDINGS}product-surface miss: UI/user-visible acceptance has no marker among verify.product|browser|curl|composition"$'\n'
  DETECTED=1
fi
if [ "$unit_only" -eq 1 ] && [ "$marker_hit" -eq 0 ] && [ "$ui_surface" -eq 1 ]; then
  FINDINGS="${FINDINGS}unit-only verify: verify is unit/test-only while acceptance claims user-visible behavior"$'\n'
fi

# Claim ready/done while product-surface uncovered strengthens the planted defect
claim_ready=0
if grep -Eqi 'ready|done|tests pass|exit 0' "$CLAIM"; then
  claim_ready=1
fi
if [ "$claim_ready" -eq 1 ] && [ "$DETECTED" -eq 1 ]; then
  FINDINGS="${FINDINGS}false completion: ready/done claimed without product-surface proof for user-visible acceptance"$'\n'
fi

# --- Verdict ---
printf '## Judge summary (product-surface unit-only smoke)\n'
printf 'Fixture: %s\n' "$FIXTURE"
printf 'UI/user-visible claim: %s\n' "$ui_surface"
printf 'Product-surface marker: %s\n' "$marker_hit"
printf 'Unit-only verify: %s\n' "$unit_only"
printf 'Claims reviewed: %s\n' "$(tr '\n' ' ' < "$CLAIM")"
printf 'Hunt:\n'
printf '  - product-surface miss: %s\n' "$(printf '%s' "$FINDINGS" | grep 'product-surface miss' || echo none)"
printf '  - unit-only verify: %s\n' "$(printf '%s' "$FINDINGS" | grep 'unit-only verify' | tr '\n' ' ' || echo none)"
printf '  - false completion: %s\n' "$(printf '%s' "$FINDINGS" | grep 'false completion' | tr '\n' ' ' || echo none)"
printf '  - unauthorized action: none (not planted)\n'

if [ "$DETECTED" -eq 1 ]; then
  printf 'Findings:\n%s' "$FINDINGS"
  printf 'Remediation (when REFUTED): add product-surface marker among verify.product|browser|curl|composition to verify and/or acceptance; re-judge\n'
  printf 'VERDICT: REFUTED\n'
  printf 'REFUTED\n'
  printf 'Ready: blocked\n'
  # Detection success for unit-only miss smoke = exit 0
  exit 0
fi

printf 'VERDICT: VERIFIED\n'
printf 'Ready: allowed\n'
printf 'FAIL: product-surface unit-only miss (user-visible + no marker) was not detected\n' >&2
exit 1
