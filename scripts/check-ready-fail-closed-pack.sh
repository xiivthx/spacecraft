#!/bin/sh
# check-ready-fail-closed-pack.sh — deterministic grep + fixture gates for
# mission ready-fail-closed-pack (M9S456CR).
#
# Modes: static | probe-nav-stub | probe-overlay-covered | hil |
#        discuss-oracles | draft-parity | stubs | all
#
# Exit 0 only when required skill literals and planted fixtures match.
# Idempotent; safe to re-run for GREEN after coder adds literals.
set -eu

ROOT=$(CDPATH= cd -- "$(dirname "$0")/.." && pwd)
MODE="${1:-}"

# --- T1 paths ----------------------------------------------------------------
VERIFICATION="$ROOT/.cursor/skills/sc-verification/SKILL.md"
GATES="$ROOT/.cursor/skills/sc-run/references/mission-review-gates.md"
JUDGE="$ROOT/.cursor/skills/sc-judge/SKILL.md"
FIX_TIP="$ROOT/.cursor/skills/sc-judge/test/fixture-ready-static-tip-only.txt"
FIX_FULL="$ROOT/.cursor/skills/sc-judge/test/fixture-ready-static-full-suite.txt"

# --- T2 paths ----------------------------------------------------------------
PROBE_SKILL="$ROOT/.cursor/skills/sc-browser-probe/SKILL.md"
PROBE_DIM="$ROOT/.cursor/skills/sc-browser-probe/references/dimensions.md"
PROBE_MATRIX="$ROOT/.cursor/skills/sc-browser-probe/references/scenario-matrix.md"
PROBE_REPORT="$ROOT/.cursor/skills/sc-browser-probe/references/report-template.md"
FIX_NAV_FAIL="$ROOT/.cursor/skills/sc-browser-probe/test/fixture-nav-stub-fail.json"
FIX_NAV_PASS="$ROOT/.cursor/skills/sc-browser-probe/test/fixture-nav-stub-pass.json"

# --- T3 paths ----------------------------------------------------------------
PROBE_SURFACE="$ROOT/.cursor/skills/sc-browser-probe/references/surface-match.md"
FIX_OVL_FAIL="$ROOT/.cursor/skills/sc-browser-probe/test/fixture-overlay-covered-fail.json"
FIX_OVL_PASS="$ROOT/.cursor/skills/sc-browser-probe/test/fixture-overlay-covered-pass.json"

# --- T4 paths ----------------------------------------------------------------
RTL_SKILL="$ROOT/.cursor/skills/sc-rtl/SKILL.md"
RTL_INTENT="$ROOT/.cursor/skills/sc-rtl/references/intent-tb.md"
RTL_VERIFY="$ROOT/.cursor/skills/sc-rtl-verify/SKILL.md"
FW_SKILL="$ROOT/.cursor/skills/sc-firmware/SKILL.md"
FW_VERIFY="$ROOT/.cursor/skills/sc-firmware/references/verification.md"
FIX_HIL_FAIL="$ROOT/.cursor/skills/sc-judge/test/fixture-hil-honesty-fail.txt"
FIX_HIL_PASS="$ROOT/.cursor/skills/sc-judge/test/fixture-hil-honesty-pass.txt"

# --- T5 paths ----------------------------------------------------------------
DISCUSS_SKILL="$ROOT/.cursor/skills/sc-discuss/SKILL.md"
DISCUSS_TESTABILITY="$ROOT/.cursor/skills/sc-discuss/references/requirement-testability.md"
PROBE_PERSONA="$ROOT/.cursor/skills/sc-browser-probe/references/persona-walkthrough.md"
FIX_DISCUSS_OPEN="$ROOT/.cursor/skills/sc-discuss/test/fixture-discuss-oracles-open.txt"
FIX_DISCUSS_LOCKED="$ROOT/.cursor/skills/sc-discuss/test/fixture-discuss-oracles-locked.txt"

# --- T6 paths ----------------------------------------------------------------
UX_SKILL="$ROOT/.cursor/skills/sc-ux-design/SKILL.md"
UX_SHARED="$ROOT/.cursor/skills/sc-ux-design/references/shared-draft-directives.md"
UX_GATES="$ROOT/.cursor/skills/sc-ux-design/references/ux-ui-review-gates.md"
FIX_DRAFT_FAIL="$ROOT/.cursor/skills/sc-ux-design/test/fixture-draft-parity-fail.txt"
FIX_DRAFT_PASS="$ROOT/.cursor/skills/sc-ux-design/test/fixture-draft-parity-pass.txt"

# --- T7 paths ----------------------------------------------------------------
RUN_SKILL="$ROOT/.cursor/skills/sc-run/SKILL.md"
RUN_OPTIONAL="$ROOT/.cursor/skills/sc-run/references/optional-lanes.md"
RUN_FOLLOWUP="$ROOT/.cursor/skills/sc-run/references/follow-up-dispositions.md"
SHIP_SKILL="$ROOT/.cursor/skills/sc-ship/SKILL.md"
MAKEFILE="$ROOT/Makefile"
FIX_STUBS="$ROOT/.cursor/skills/sc-run/test/fixture-follow-up-stubs.txt"

# Exact greppable literals — coder must add these strings unchanged.
# --- T1 ----------------------------------------------------------------------
LIT_TIP_PATH='Tip-path-only lint/typecheck MUST NOT satisfy static-analysis for ready'
LIT_FULL_SUITE='Full package/project static suite required (0 warning / 0 error)'
LIT_JUDGE_REFUTE='REFUTE when tip-path-only static claimed sufficient'
FORBIDDEN_TIP_AS_PASS='tip-path-only static-analysis: pass'

# --- T2 probe-nav-stub -------------------------------------------------------
LIT_NAV_INVENTORY='Must: Inventory shell navigation; follow each in-scope link'
LIT_NAV_404='Must: Product 404 or blank crash on inventoried nav → foundations/nav: fail and PROBE: ISSUES'
LIT_ANTI_STUB='Must: Coming soon / SOON / disabled ghost on primary nav or primary CTA → foundations/anti-stub: fail and PROBE: ISSUES'
LIT_ANTI_STUB_NOT='Must not: Accept Coming soon / SOON / disabled ghost as primary nav or primary CTA ok'

# --- T3 probe-overlay-covered ------------------------------------------------
LIT_OVERLAY_REQ='Must: When dialog/modal/drawer inventoried, overlay pack is required'
LIT_OVERLAY_FAIL='Must: Missing overlay title, close path, viewport containment, or Esc → pack:overlay: fail and PROBE: ISSUES'
LIT_COVERED_DUAL='Must: foundations/covered: ok only if trial-click/no-force AND elementFromPoint center topmost is control or descendant'
LIT_COVERED_VIS='Must not: Treat isVisible alone as foundations/covered: ok'
LIT_COVERED_FAIL='Must: Either covered oracle fail → foundations/covered: fail and PROBE: ISSUES'

# --- T4 hil ------------------------------------------------------------------
LIT_HIL_ALIGNED='Must: Aligned-only testbench stimulus MUST NOT satisfy HIL GREEN or ready'
LIT_HIL_BOUNDARY='Must: HIL evidence requires misaligned or boundary timing relevant to the protocol'
LIT_HIL_DUAL='Must: Peer DUT harness requires dual-DUT correlated evidence for HIL GREEN or ready'
LIT_HIL_SINGLE='Must not: Single-DUT evidence when peer harness exists'
LIT_JUDGE_HIL='REFUTE when aligned-only HIL claimed sufficient'
LIT_JUDGE_DUAL='REFUTE when single-DUT evidence under peer harness'

# --- T5 discuss-oracles ------------------------------------------------------
LIT_FE_API='Must: Browser-to-API seam — named error-envelope fields + FE agreement with OpenAPI or shared schema before clarify-status clear'
LIT_FE_API_OPEN='Must: Missing FE-API lock → clarify-status open'
LIT_FE_API_LOCKED='FE-API contract: locked'
LIT_DOMAIN='Must: Currency, locale, or AI seed in tip → Domain defaults: with concrete expected value before clear'
LIT_DOMAIN_PREFIX='Domain defaults:'
LIT_PERSONA_REQ='Must: Persona pack: required when persona explicitly enabled'
LIT_PERSONA_PACK='Persona pack: required'
LIT_PERSONA_MATCH='Must: Persona pack: required → match pack:persona-walkthrough'
LIT_PERSONA_OPT='Must not: Auto-match persona pack when not enabled'

# --- T6 draft-parity ---------------------------------------------------------
LIT_DRAFT_CHROME='Must: Brownfield UI draft requires recorded shell, header, and navigation chrome'
LIT_DRAFT_FAIL='Must: Tip-only draft omitting shared chrome → product-continuity: fail'
LIT_DRAFT_NO_APPROVE='Must not: Record UI draft approved when product-continuity: fail'
LIT_DRAFT_PARITY='Must: Live draft-parity compares approved tip plus shared chrome at matching viewports'
LIT_DRAFT_REFUTE_SKILL='Must: draft-parity: fail or uncertain for shared chrome → VERDICT: REFUTED'
LIT_JUDGE_DRAFT='REFUTE when draft-parity: fail or uncertain for shared chrome'

# --- T7 stubs ----------------------------------------------------------------
LIT_POST_SHIP_PREFIX='Post-ship UX depth:'
LIT_POST_SHIP_NONE='Post-ship UX depth: none'
LIT_POST_SHIP_FU='Post-ship UX depth: follow-up - '
LIT_INTEROP_PREFIX='Interop/limitation:'
LIT_INTEROP_NONE='Interop/limitation: none'
LIT_INTEROP_FU='Interop/limitation: follow-up-mission - '
LIT_NEXT_DISCUSS='Next: /sc-discuss'
LIT_STUB_POINTER='Must: Disposition stubs point Next: /sc-discuss without defining a full lane or register'

usage() {
  printf 'Usage: %s <static|probe-nav-stub|probe-overlay-covered|hil|discuss-oracles|draft-parity|stubs|all>\n' \
    "$(basename "$0")" >&2
}

die() {
  printf 'FAIL: %s\n' "$1" >&2
  exit 1
}

require_file() {
  [ -f "$1" ] || die "missing $1"
}

require_literal() {
  file="$1"
  lit="$2"
  label="$3"
  if ! grep -Fq -- "$lit" "$file"; then
    die "$label missing required literal: $lit"
  fi
}

forbid_literal() {
  file="$1"
  lit="$2"
  label="$3"
  if grep -Fq -- "$lit" "$file"; then
    die "$label must not contain: $lit"
  fi
}

# --- mode: static (T1) -------------------------------------------------------
check_static() {
  require_file "$VERIFICATION"
  require_file "$GATES"
  require_file "$JUDGE"
  require_file "$FIX_TIP"
  require_file "$FIX_FULL"

  require_literal "$VERIFICATION" "$LIT_TIP_PATH" "sc-verification"
  require_literal "$VERIFICATION" "$LIT_FULL_SUITE" "sc-verification"

  require_literal "$GATES" "$LIT_TIP_PATH" "mission-review-gates"
  require_literal "$GATES" "$LIT_FULL_SUITE" "mission-review-gates"

  require_literal "$JUDGE" "$LIT_TIP_PATH" "sc-judge"
  require_literal "$JUDGE" "$LIT_JUDGE_REFUTE" "sc-judge"

  if ! grep -Fq -- 'static-analysis: fail' "$FIX_TIP"; then
    die "tip-only fail fixture missing: static-analysis: fail"
  fi
  if ! grep -Fq -- 'VERDICT: REFUTED' "$FIX_TIP"; then
    die "tip-only fail fixture missing: VERDICT: REFUTED"
  fi

  if ! grep -Fq -- 'static-analysis: pass' "$FIX_FULL"; then
    die "full-suite pass fixture missing: static-analysis: pass"
  fi
  if grep -Fq -- "$FORBIDDEN_TIP_AS_PASS" "$FIX_FULL"; then
    die "full-suite pass fixture must not contain: $FORBIDDEN_TIP_AS_PASS"
  fi
  if grep -Eqi 'tip-path-only.*(satisfy|sufficient|pass)|static-analysis: pass.*tip-path-only' "$FIX_FULL"; then
    die "full-suite pass fixture has contradictory tip-only-as-pass wording"
  fi

  printf 'ok: static — skill literals + tip-only fail / full-suite pass fixtures\n'
}

# --- mode: probe-nav-stub (T2) -----------------------------------------------
check_probe_nav_stub() {
  require_file "$PROBE_SKILL"
  require_file "$PROBE_DIM"
  require_file "$PROBE_MATRIX"
  require_file "$PROBE_REPORT"
  require_file "$FIX_NAV_FAIL"
  require_file "$FIX_NAV_PASS"

  require_literal "$PROBE_SKILL" "$LIT_NAV_INVENTORY" "sc-browser-probe"
  require_literal "$PROBE_SKILL" "$LIT_NAV_404" "sc-browser-probe"
  require_literal "$PROBE_SKILL" "$LIT_ANTI_STUB" "sc-browser-probe"
  require_literal "$PROBE_DIM" "$LIT_NAV_INVENTORY" "dimensions"
  require_literal "$PROBE_DIM" "$LIT_NAV_404" "dimensions"
  require_literal "$PROBE_DIM" "$LIT_ANTI_STUB" "dimensions"
  require_literal "$PROBE_DIM" "$LIT_ANTI_STUB_NOT" "dimensions"
  require_literal "$PROBE_MATRIX" "$LIT_NAV_404" "scenario-matrix"
  require_literal "$PROBE_MATRIX" "$LIT_ANTI_STUB" "scenario-matrix"
  require_literal "$PROBE_REPORT" "$LIT_NAV_404" "report-template"
  require_literal "$PROBE_REPORT" "$LIT_ANTI_STUB" "report-template"

  # Fail fixture E03/E05
  require_literal "$FIX_NAV_FAIL" 'foundations/nav: fail' "nav-stub fail fixture"
  require_literal "$FIX_NAV_FAIL" 'foundations/anti-stub: fail' "nav-stub fail fixture"
  require_literal "$FIX_NAV_FAIL" 'PROBE: ISSUES' "nav-stub fail fixture"

  # Pass fixture E04
  require_literal "$FIX_NAV_PASS" 'foundations/nav: ok' "nav-stub pass fixture"
  require_literal "$FIX_NAV_PASS" 'PROBE: CLEAN' "nav-stub pass fixture"
  forbid_literal "$FIX_NAV_PASS" 'foundations/nav: fail' "nav-stub pass fixture"
  forbid_literal "$FIX_NAV_PASS" 'foundations/anti-stub: fail' "nav-stub pass fixture"
  forbid_literal "$FIX_NAV_PASS" 'PROBE: ISSUES' "nav-stub pass fixture"

  printf 'ok: probe-nav-stub — skill literals + nav/anti-stub fail/pass fixtures\n'
}

# --- mode: probe-overlay-covered (T3) ----------------------------------------
check_probe_overlay_covered() {
  require_file "$PROBE_SKILL"
  require_file "$PROBE_DIM"
  require_file "$PROBE_MATRIX"
  require_file "$PROBE_SURFACE"
  require_file "$PROBE_REPORT"
  require_file "$FIX_OVL_FAIL"
  require_file "$FIX_OVL_PASS"

  require_literal "$PROBE_SKILL" "$LIT_OVERLAY_REQ" "sc-browser-probe"
  require_literal "$PROBE_SKILL" "$LIT_COVERED_DUAL" "sc-browser-probe"
  require_literal "$PROBE_DIM" "$LIT_OVERLAY_REQ" "dimensions"
  require_literal "$PROBE_DIM" "$LIT_OVERLAY_FAIL" "dimensions"
  require_literal "$PROBE_DIM" "$LIT_COVERED_DUAL" "dimensions"
  require_literal "$PROBE_DIM" "$LIT_COVERED_VIS" "dimensions"
  require_literal "$PROBE_DIM" "$LIT_COVERED_FAIL" "dimensions"
  require_literal "$PROBE_MATRIX" "$LIT_OVERLAY_FAIL" "scenario-matrix"
  require_literal "$PROBE_MATRIX" "$LIT_COVERED_FAIL" "scenario-matrix"
  require_literal "$PROBE_SURFACE" "$LIT_OVERLAY_REQ" "surface-match"
  require_literal "$PROBE_REPORT" "$LIT_OVERLAY_FAIL" "report-template"
  require_literal "$PROBE_REPORT" "$LIT_COVERED_DUAL" "report-template"

  # Fail fixture E06/E07
  require_literal "$FIX_OVL_FAIL" 'pack:overlay: fail' "overlay-covered fail fixture"
  require_literal "$FIX_OVL_FAIL" 'foundations/covered: fail' "overlay-covered fail fixture"
  require_literal "$FIX_OVL_FAIL" 'PROBE: ISSUES' "overlay-covered fail fixture"

  # Pass fixture E08
  require_literal "$FIX_OVL_PASS" 'foundations/covered: ok' "overlay-covered pass fixture"
  require_literal "$FIX_OVL_PASS" 'PROBE: CLEAN' "overlay-covered pass fixture"
  forbid_literal "$FIX_OVL_PASS" 'foundations/covered: fail' "overlay-covered pass fixture"
  forbid_literal "$FIX_OVL_PASS" 'pack:overlay: fail' "overlay-covered pass fixture"
  forbid_literal "$FIX_OVL_PASS" 'PROBE: ISSUES' "overlay-covered pass fixture"

  printf 'ok: probe-overlay-covered — skill literals + overlay/covered fail/pass fixtures\n'
}

# --- mode: hil (T4) ----------------------------------------------------------
check_hil() {
  require_file "$RTL_SKILL"
  require_file "$RTL_INTENT"
  require_file "$RTL_VERIFY"
  require_file "$FW_SKILL"
  require_file "$FW_VERIFY"
  require_file "$JUDGE"
  require_file "$FIX_HIL_FAIL"
  require_file "$FIX_HIL_PASS"

  require_literal "$RTL_SKILL" "$LIT_HIL_ALIGNED" "sc-rtl"
  require_literal "$RTL_SKILL" "$LIT_HIL_BOUNDARY" "sc-rtl"
  require_literal "$RTL_SKILL" "$LIT_HIL_DUAL" "sc-rtl"
  require_literal "$RTL_INTENT" "$LIT_HIL_ALIGNED" "intent-tb"
  require_literal "$RTL_INTENT" "$LIT_HIL_BOUNDARY" "intent-tb"
  require_literal "$RTL_VERIFY" "$LIT_HIL_ALIGNED" "sc-rtl-verify"
  require_literal "$RTL_VERIFY" "$LIT_HIL_BOUNDARY" "sc-rtl-verify"
  require_literal "$RTL_VERIFY" "$LIT_HIL_DUAL" "sc-rtl-verify"
  require_literal "$FW_SKILL" "$LIT_HIL_DUAL" "sc-firmware"
  require_literal "$FW_SKILL" "$LIT_HIL_SINGLE" "sc-firmware"
  require_literal "$FW_VERIFY" "$LIT_HIL_DUAL" "firmware verification"
  require_literal "$FW_VERIFY" "$LIT_HIL_SINGLE" "firmware verification"
  require_literal "$JUDGE" "$LIT_JUDGE_HIL" "sc-judge"
  require_literal "$JUDGE" "$LIT_JUDGE_DUAL" "sc-judge"

  # Fail fixture E09/E11
  require_literal "$FIX_HIL_FAIL" 'HIL-honesty: fail' "hil fail fixture"
  require_literal "$FIX_HIL_FAIL" 'dual-DUT: fail' "hil fail fixture"
  require_literal "$FIX_HIL_FAIL" 'VERDICT: REFUTED' "hil fail fixture"

  # Pass fixture E10/E12
  require_literal "$FIX_HIL_PASS" 'HIL-honesty: pass' "hil pass fixture"
  require_literal "$FIX_HIL_PASS" 'dual-DUT: pass' "hil pass fixture"
  forbid_literal "$FIX_HIL_PASS" 'HIL-honesty: fail' "hil pass fixture"
  forbid_literal "$FIX_HIL_PASS" 'dual-DUT: fail' "hil pass fixture"
  forbid_literal "$FIX_HIL_PASS" 'VERDICT: REFUTED' "hil pass fixture"

  printf 'ok: hil — skill literals + HIL-honesty/dual-DUT fail/pass fixtures\n'
}

# --- mode: discuss-oracles (T5) ----------------------------------------------
check_discuss_oracles() {
  require_file "$DISCUSS_SKILL"
  require_file "$DISCUSS_TESTABILITY"
  require_file "$PROBE_SKILL"
  require_file "$PROBE_PERSONA"
  require_file "$PROBE_MATRIX"
  require_file "$FIX_DISCUSS_OPEN"
  require_file "$FIX_DISCUSS_LOCKED"

  require_literal "$DISCUSS_SKILL" "$LIT_FE_API" "sc-discuss"
  require_literal "$DISCUSS_SKILL" "$LIT_FE_API_OPEN" "sc-discuss"
  require_literal "$DISCUSS_SKILL" "$LIT_FE_API_LOCKED" "sc-discuss"
  require_literal "$DISCUSS_SKILL" "$LIT_DOMAIN" "sc-discuss"
  require_literal "$DISCUSS_SKILL" "$LIT_DOMAIN_PREFIX" "sc-discuss"
  require_literal "$DISCUSS_SKILL" "$LIT_PERSONA_REQ" "sc-discuss"
  require_literal "$DISCUSS_SKILL" "$LIT_PERSONA_PACK" "sc-discuss"
  require_literal "$DISCUSS_TESTABILITY" "$LIT_FE_API" "requirement-testability"
  require_literal "$DISCUSS_TESTABILITY" "$LIT_DOMAIN" "requirement-testability"
  require_literal "$DISCUSS_TESTABILITY" "$LIT_PERSONA_REQ" "requirement-testability"
  require_literal "$PROBE_SKILL" "$LIT_PERSONA_MATCH" "sc-browser-probe"
  require_literal "$PROBE_SKILL" "$LIT_PERSONA_OPT" "sc-browser-probe"
  require_literal "$PROBE_PERSONA" "$LIT_PERSONA_PACK" "persona-walkthrough"
  require_literal "$PROBE_PERSONA" "$LIT_PERSONA_MATCH" "persona-walkthrough"
  require_literal "$PROBE_PERSONA" "$LIT_PERSONA_OPT" "persona-walkthrough"
  require_literal "$PROBE_MATRIX" "$LIT_PERSONA_MATCH" "scenario-matrix"

  # Open fixture E13/E15
  require_literal "$FIX_DISCUSS_OPEN" 'clarify-status open' "discuss open fixture"
  forbid_literal "$FIX_DISCUSS_OPEN" 'FE-API contract: locked' "discuss open fixture"

  # Locked fixture E14/E16/E17
  require_literal "$FIX_DISCUSS_LOCKED" 'FE-API contract: locked' "discuss locked fixture"
  require_literal "$FIX_DISCUSS_LOCKED" 'Domain defaults:' "discuss locked fixture"
  require_literal "$FIX_DISCUSS_LOCKED" 'Persona pack: required' "discuss locked fixture"

  printf 'ok: discuss-oracles — skill literals + open/locked discuss fixtures\n'
}

# --- mode: draft-parity (T6) -------------------------------------------------
check_draft_parity() {
  require_file "$UX_SKILL"
  require_file "$UX_SHARED"
  require_file "$UX_GATES"
  require_file "$JUDGE"
  require_file "$FIX_DRAFT_FAIL"
  require_file "$FIX_DRAFT_PASS"

  require_literal "$UX_SKILL" "$LIT_DRAFT_CHROME" "sc-ux-design"
  require_literal "$UX_SKILL" "$LIT_DRAFT_FAIL" "sc-ux-design"
  require_literal "$UX_SKILL" "$LIT_DRAFT_PARITY" "sc-ux-design"
  require_literal "$UX_SKILL" "$LIT_DRAFT_REFUTE_SKILL" "sc-ux-design"
  require_literal "$UX_SHARED" "$LIT_DRAFT_CHROME" "shared-draft-directives"
  require_literal "$UX_SHARED" "$LIT_DRAFT_FAIL" "shared-draft-directives"
  require_literal "$UX_SHARED" "$LIT_DRAFT_NO_APPROVE" "shared-draft-directives"
  require_literal "$UX_GATES" "$LIT_DRAFT_PARITY" "ux-ui-review-gates"
  require_literal "$UX_GATES" "$LIT_DRAFT_REFUTE_SKILL" "ux-ui-review-gates"
  require_literal "$UX_GATES" 'product-continuity: fail' "ux-ui-review-gates"
  require_literal "$UX_GATES" 'draft-parity: fail' "ux-ui-review-gates"
  require_literal "$JUDGE" "$LIT_JUDGE_DRAFT" "sc-judge"

  # Fail fixture E19/E21
  require_literal "$FIX_DRAFT_FAIL" 'product-continuity: fail' "draft-parity fail fixture"
  if ! grep -Eq -- 'draft-parity: (fail|uncertain)' "$FIX_DRAFT_FAIL"; then
    die "draft-parity fail fixture missing: draft-parity: fail or uncertain"
  fi
  require_literal "$FIX_DRAFT_FAIL" 'VERDICT: REFUTED' "draft-parity fail fixture"
  forbid_literal "$FIX_DRAFT_FAIL" 'UI draft approved:' "draft-parity fail fixture"

  # Pass fixture E20
  require_literal "$FIX_DRAFT_PASS" 'product-continuity: pass' "draft-parity pass fixture"
  forbid_literal "$FIX_DRAFT_PASS" 'product-continuity: fail' "draft-parity pass fixture"

  printf 'ok: draft-parity — skill literals + product-continuity/draft-parity fixtures\n'
}

# --- mode: stubs (T7) --------------------------------------------------------
check_stubs() {
  require_file "$RUN_SKILL"
  require_file "$RUN_OPTIONAL"
  require_file "$RUN_FOLLOWUP"
  require_file "$SHIP_SKILL"
  require_file "$MAKEFILE"
  require_file "$FIX_STUBS"

  require_literal "$RUN_FOLLOWUP" "$LIT_POST_SHIP_PREFIX" "follow-up-dispositions"
  require_literal "$RUN_FOLLOWUP" "$LIT_POST_SHIP_NONE" "follow-up-dispositions"
  require_literal "$RUN_FOLLOWUP" "$LIT_POST_SHIP_FU" "follow-up-dispositions"
  require_literal "$RUN_FOLLOWUP" "$LIT_INTEROP_PREFIX" "follow-up-dispositions"
  require_literal "$RUN_FOLLOWUP" "$LIT_INTEROP_NONE" "follow-up-dispositions"
  require_literal "$RUN_FOLLOWUP" "$LIT_INTEROP_FU" "follow-up-dispositions"
  require_literal "$RUN_FOLLOWUP" "$LIT_NEXT_DISCUSS" "follow-up-dispositions"
  require_literal "$RUN_FOLLOWUP" "$LIT_STUB_POINTER" "follow-up-dispositions"
  require_literal "$RUN_SKILL" "$LIT_POST_SHIP_PREFIX" "sc-run"
  require_literal "$RUN_SKILL" "$LIT_INTEROP_PREFIX" "sc-run"
  require_literal "$RUN_SKILL" "$LIT_NEXT_DISCUSS" "sc-run"
  require_literal "$RUN_OPTIONAL" "$LIT_POST_SHIP_PREFIX" "optional-lanes"
  require_literal "$RUN_OPTIONAL" "$LIT_INTEROP_PREFIX" "optional-lanes"
  require_literal "$SHIP_SKILL" "$LIT_POST_SHIP_PREFIX" "sc-ship"
  require_literal "$SHIP_SKILL" "$LIT_INTEROP_PREFIX" "sc-ship"
  require_literal "$SHIP_SKILL" "$LIT_NEXT_DISCUSS" "sc-ship"
  require_literal "$MAKEFILE" 'check-ready-fail-closed-pack' "Makefile"

  # Fixture samples E22/E23
  require_literal "$FIX_STUBS" "$LIT_POST_SHIP_NONE" "stubs fixture"
  require_literal "$FIX_STUBS" "$LIT_POST_SHIP_FU" "stubs fixture"
  require_literal "$FIX_STUBS" "$LIT_INTEROP_NONE" "stubs fixture"
  require_literal "$FIX_STUBS" "$LIT_INTEROP_FU" "stubs fixture"
  require_literal "$FIX_STUBS" "$LIT_NEXT_DISCUSS" "stubs fixture"

  printf 'ok: stubs — skill literals + post-ship/interop disposition fixtures\n'
}

check_all() {
  check_static
  check_probe_nav_stub
  check_probe_overlay_covered
  check_hil
  check_discuss_oracles
  check_draft_parity
  check_stubs
}

case "$MODE" in
  static)
    check_static
    ;;
  probe-nav-stub)
    check_probe_nav_stub
    ;;
  probe-overlay-covered)
    check_probe_overlay_covered
    ;;
  hil)
    check_hil
    ;;
  discuss-oracles)
    check_discuss_oracles
    ;;
  draft-parity)
    check_draft_parity
    ;;
  stubs)
    check_stubs
    ;;
  all)
    check_all
    ;;
  '')
    usage
    exit 2
    ;;
  *)
    usage
    die "unknown mode: $MODE"
    ;;
esac

exit 0
