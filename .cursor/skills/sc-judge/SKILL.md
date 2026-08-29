---
name: sc-judge
description: "Adversarial prove gate before ready. Treat completion claims as claims; re-run evidence; diff scope vs plan; hunt weakened tests / false completion / unauthorized action. Verdict VERIFIED | REFUTED only."
---

# sc-judge

## Goal

Adversarial prove gate before `ready`: re-observe claimed completion. Never trust a report, summary, or "done" claim alone. Judge proves readiness - it does not replace `/sc-discuss` → `/sc-run` → `/sc-ship`.

## Output

```
VERDICT: VERIFIED | REFUTED
```

Plus: fresh evidence ids, scope-diff notes, hunt findings, per-finding dissent labels, cross-model critic disposition, and (when `REFUTED`) remediation for `/sc-run`. Cursor-native only.

## Good / Bad

- Good: treat every completion / "done" / "ready" claim as a claim; re-run claimed evidence; diff scope vs `plan.json` / spec; hunt oracle-tamper, AUTH gaps, leftover `review.json` findings (any severity / any `source`), and greppable `Cursor review:` / `Cursor review skipped:` when Cursor review ran or is required; when disposition claims `ran`, require corroboration (`cursor-review-…` evidence or greppable `Cursor ingest: session`); emit exactly one verdict; allow `ready` **only** on `VERIFIED`; require A/B/C outcome-gate disposition per `docs/mission-artifacts.md`; judge from **structured-lines-only** zero-context input; label every finding `AGREE` / `DISAGREE_EVIDENCE` / `DISAGREE_CONCERN` with cited code; emit `Cross-model critic:` or skip disposition (never silent)
- Bad: trusting a prior report or `evidence.jsonl` line without re-run; inventing evidence; soft-shipping past `REFUTED`; third verdicts or caveat soft-pass; re-walking the full mission-review dimension table; treating optional canvases as `VERIFIED` / ready proof; free-text builder rationale as judge input; silent cross-model critic omission; missing Cursor-review disposition when review ran/required; accepting `Cursor review: … ran` without corroboration

## Verify

```
spacecraft evidence "judge-<mission-id>" -- <re-run of claimed commands>
```

Confirm verdict is exactly `VERIFIED` | `REFUTED`, hunts cover the ready-path targets below, scope-diff is recorded, dissent labels + cross-model disposition are present, and `Ready: allowed` only on `VERIFIED`.

## When to use

- Mission claims build complete or moves toward `ready`
- Review / `/sc-run` needs an adversarial prove step before `set-state ready`
- A completion, "done", or "ready" claim must be re-observed

## Workflow

1. **Resolve** - `spacecraft resolve` (or `spacecraft use <selector>` on conflict).
2. **Collect claims (zero-context)** - Input is **structured-lines-only**: `spec.md` + design-contract + diff + evidence labels + extracted `decisions.md` lines only (grep prefixes: `Scenario oracle change:`, `Roadmap contract:`, `Cursor review:`, `Cursor review skipped:`, `Cursor ingest:`, `Sc-security fallback:`, `Loop watch:`, debt/disposition lines such as `Mutation skipped:`, `Pbt skipped:`, `Characterization waived:`, `Debt ceiling:`). Treat each as a claim. **Never** free-text builder rationale, narrative task notes, or unprefixed `decisions.md` prose. Optional canvases aid human check only - not claims to prove; not a `VERIFIED` / ready gate.
3. **Re-run claimed evidence** - For every command cited as proving acceptance, re-run via `spacecraft evidence "<label>-judge" -- <command>`. Never reuse a stale `evidence.jsonl` line as sole proof.
4. **Diff scope vs plan** - Compare the change set to `plan.json` tasks and `spec.md` acceptance. Flag out-of-plan work, missing acceptance coverage, or "done" without matching fresh evidence.
5. **Hunt (ready path)** - Search these targets only (do **not** re-walk the full mission-review dimension table; reviewer already applied `mission-review-gates` / UX gates):
   - **oracle-tamper** - assertions removed/skipped/loosened/tautologies; frozen scenario expected literals edited to force green; coder edited tests without Commander + `decisions.md` note; `Scenario oracle change:` missing when oracle changed
   - **false completion** - "done"/"ready" while acceptance fails, evidence missing/stale, scope mismatch, or defects left unfixed
   - **AUTH** - outward push/deploy/publish/send without quoted `AUTH:`; ship/merge without `/sc-ship` gates
   - **leftover findings** - any severity in `review.json` (`critical` / `important` / `minor`) and any `source` (including `bugbot` / `security-review` / `sc-reviewer`) ⇒ `REFUTED`. Empty findings only; do not re-score mission-review dimensions here
   - **Cursor review disposition** - when Cursor review ran or is required on the ready path, require greppable `Cursor review:` or `Cursor review skipped:` (same grammar as sc-run, e.g. `Cursor review: bugbot+security-review ran`); missing disposition ⇒ `REFUTED`. When disposition claims `Cursor review: … ran`, require **corroboration**: an evidence label matching `cursor-review-…` **or** greppable `Cursor ingest: session` in `decisions.md`; missing corroboration ⇒ `REFUTED` (disposition theater). When security is in scope and disposition is `Cursor review skipped:`, require greppable `Sc-security fallback: pass` (optionally also `Sc-security fallback: findings drained`; optional evidence label `sc-security-…`) or SEC machine-evidence pass; skip alone ⇒ `REFUTED`
   - **Loop watch disposition** - when greppable `Loop watch: ran` is present → require greppable `Loop watch: stopped:` before ready; `Loop watch: skipped:` alone after `ran` without any `stopped:` ⇒ `REFUTED`; coexisting `skipped:` does not REFUTE when `stopped:` is also present. When no `ran` → `Loop watch: skipped:` OK / hunt N/A. Armed loop claimed after ready without disarm ⇒ `REFUTED`
   - **outcome-gate disposition** - product ready needs `Design-contract: complete` or `Design-contract skipped: docs/prose-only`; `Approved-scenarios:` freeze footer or `Approved-scenarios skipped: docs/prose-only`; static / diff-cov / mutation / PBT disposition (evidence labels or greppable lines such as `Static-analysis skipped:…` / `Mutation skipped:…` / `Pbt skipped:…`). Exact prefixes and bars: `docs/mission-artifacts.md` - static **0 warning / 0 error** when tool runs; diff-cov touched **line and branch ≥90%** when measured; mutation in scope when any of `Mutation: required` | pack `quality` | `Mutation: high-risk`, then **>80%** scoped when tool present (else `Mutation skipped: not in scope` valid); PBT **100%** of design-contract **core-logic** modules (`pbt-…` invariants + generators via project-existing `fast-check` / Hypothesis / equivalent, or `Pbt skipped: no project pbt tool` / `Pbt skipped: not core logic` / `Pbt waived: <reason>`). Missing disposition without skip/waive ⇒ `REFUTED`. Inventing PBT lib install mid-mission ⇒ `REFUTED`. Global 95–100% coverage as success bar ⇒ `REFUTED`
   - **hard-gated Test Ideas** (when present) - Neg/Overlooked (+ Strategy Top risk/Charter when mapped) without matching `acceptance[]` and without `Deferred test idea: <id> - <reason>` ⇒ `REFUTED`; claimed done without fresh evidence ⇒ `REFUTED`
   - **product-surface miss** (when UI/workflow claimed) - need `verify.product` | `browser` | `curl` | `composition`; unit-only ⇒ `REFUTED`
   - **draft drift (visual UI only)** - when `UI draft approved:` is recorded, REFUTE clear chrome divergence, missing paired draft+live screenshot evidence, or fail/uncertain draft-parity / live-product. Point to `ux-ui-review-gates.md`; do not re-score every visual dimension here
6. **Dissent (per finding)** - For each review / hunt finding, emit `AGREE` | `DISAGREE_EVIDENCE` | `DISAGREE_CONCERN` with cited code (path + lines or greppable symbol). Every `DISAGREE_*` resolves to a fresh evidence id **or** a greppable `decisions.md` line. Bare disagreement without cite + resolution ⇒ `REFUTED`.
7. **Cross-model critic disposition** - Required output line (best-effort; **not** a Must): `Cross-model critic: <family>` or `Cross-model critic skipped: no second family configured`. Never silent omission.
8. **Verdict**
   - `VERIFIED` - fresh evidence passes; scope matches; hunts clean; **0** review findings (any `source`); Cursor-review disposition present when required (and corroboration when disposition claims `ran`); dissent trail complete; cross-model disposition present. Ready allowed (subject to other gates).
   - `REFUTED` - any material hunt hit, failed re-run, scope/acceptance mismatch, leftover findings, missing Cursor-review disposition when required, missing ran-corroboration, incomplete dissent resolution, or failed verify. Ready blocked.
9. **Ready gate** - Allow `ready` **only** on `VERIFIED`. On `REFUTED`, block ready and list remediation for `/sc-run` → fix → re-review → re-judge. No caveat / soft-pass.

### Edge cases

No claimed evidence / re-run fails → `REFUTED`. Non-defect `decisions.md` notes may sit beside `VERIFIED`; unfinished follow-up → `REFUTED`. Uncertain hard-gate or visual ready → fail-closed `REFUTED`. Manual-only → fresh observation note; do not invent output. Free-text rationale slipped into judge input → discard; use greppable lines only. `DISAGREE_*` without fresh evidence id or decisions line → `REFUTED`. Missing cross-model disposition line → `REFUTED` (emit skip line when no second family). Cursor review ran/required without greppable `Cursor review:` / `Cursor review skipped:` → `REFUTED`. Disposition claims `Cursor review: … ran` without corroboration (`cursor-review-…` evidence label or greppable `Cursor ingest: session`) → `REFUTED`. Security in scope + `Cursor review skipped:` without greppable `Sc-security fallback: pass` (or SEC machine-evidence pass) → `REFUTED`.

## Verdict contract

| Verdict | Ready |
|---------|-------|
| `VERIFIED` | **Allowed** |
| `REFUTED` | **Blocked** - fix → re-judge |

No aliases (`PASS`, `FAIL`, `VERIFIED WITH CAVEATS`, etc.).

## Rules

- **Must**: Re-run claimed evidence into `evidence.jsonl`; diff scope vs plan/spec; hunt oracle-tamper, AUTH, leftover `review.json` findings (any severity / any `source`), Cursor-review disposition when review ran/required (and corroboration when disposition claims `ran`), A/B/C + PBT disposition per `docs/mission-artifacts.md`.
- **Must**: Judge from **structured-lines-only** zero-context input (spec + design-contract + diff + evidence labels + greppable `decisions.md` prefixes only); never free-text builder rationale.
- **Must**: Per finding emit `AGREE` | `DISAGREE_EVIDENCE` | `DISAGREE_CONCERN` with cited code; every `DISAGREE_*` resolves to fresh evidence id or `decisions.md` line.
- **Must**: Emit `VERIFIED` | `REFUTED` only; allow `ready` only on `VERIFIED`; prove from evidence + scope + hunts + empty findings + `validate --strict` - not canvas.
- **Must not**: Re-walk the full mission-review dimension table; soft-pass `REFUTED`; invent evidence.
- **Must**: ASCII hyphen-minus; Cursor-native.
- Cross-model critic: required disposition line `Cross-model critic: <family>` or `Cross-model critic skipped: no second family configured` - never silent; **not** a Must (best-effort second family; skip when none configured).

## Judge-break fixtures

`make test-judge-break` (or `scripts/check-judge-break.sh`) - known-bad packs under `references/judge-break/`. Run before claiming sc-judge or closeout predicate changes are safe.

## Out of scope

Product code/tests · clarify/drafts · AFK build · ship · trap-eval / LLM-as-judge transcript scoring.

## Output format

```
## Judge summary
Mission: <id> | Claims: <list>
Input: structured-lines-only (spec | design-contract | diff | evidence labels | decisions prefixes)
Evidence re-run: <fresh ids> | Scope vs plan: <match | mismatches>
Hunt: oracle-tamper / AUTH / leftover findings / Cursor review disposition / outcome-gate / hard-gated Test Ideas
Dissent: <finding> -> AGREE | DISAGREE_EVIDENCE | DISAGREE_CONCERN + cite + (evidence id | decisions line)
Cross-model critic: <family> | Cross-model critic skipped: no second family configured
Remediation (REFUTED): <none | list>
VERDICT: VERIFIED | REFUTED | Ready: allowed | blocked
```

## References

- `sc-verification`, `/sc-run`, `plan.json` / `spec.md`, `evidence.jsonl` / `review.json`
- `docs/mission-artifacts.md` - outcome-gate skip/waive SoT (`Design-contract`, `Approved-scenarios`, `Static-analysis`, `Mutation skipped`, `Pbt skipped`, …)
- `.cursor/skills/sc-run/references/mission-review-gates.md` - reviewer applies; judge does not re-walk
- `.cursor/skills/sc-ux-design/references/ux-ui-review-gates.md` - visual fail-closed when UI draft approved
- `references/judge-break/` - known-bad fixtures
