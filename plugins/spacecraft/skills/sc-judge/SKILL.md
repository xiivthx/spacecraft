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

Plus: fresh evidence ids, scope-diff notes, hunt findings, and (when `REFUTED`) remediation for `/sc-run`. Cursor-native only.

## Good / Bad

- Good: treat every completion / "done" / "ready" claim as a claim; re-run claimed evidence; diff scope vs `plan.json` / spec; hunt oracle-tamper, AUTH gaps, leftover `review.json` findings; emit exactly one verdict; allow `ready` **only** on `VERIFIED`; require A/B/C outcome-gate disposition per `docs/mission-artifacts.md`
- Bad: trusting a prior report or `evidence.jsonl` line without re-run; inventing evidence; soft-shipping past `REFUTED`; third verdicts or caveat soft-pass; re-walking the full reviewer dimension table; treating optional canvases as `VERIFIED` / ready proof

## Verify

```
spacecraft evidence "judge-<mission-id>" -- <re-run of claimed commands>
```

Confirm verdict is exactly `VERIFIED` | `REFUTED`, hunts cover the ready-path targets below, scope-diff is recorded, and `Ready: allowed` only on `VERIFIED`.

## When to use

- Mission claims build complete or moves toward `ready`
- Review / `/sc-run` needs an adversarial prove step before `set-state ready`
- A completion, "done", or "ready" claim must be re-observed

## Workflow

1. **Resolve** - `spacecraft resolve` (or `spacecraft use <selector>` on conflict).
2. **Collect claims** - Task notes, evidence labels, review draft, or "ready" request. Treat each as a claim. Optional canvases aid human check only - not claims to prove; not a `VERIFIED` / ready gate.
3. **Re-run claimed evidence** - For every command cited as proving acceptance, re-run via `spacecraft evidence "<label>-judge" -- <command>`. Never reuse a stale `evidence.jsonl` line as sole proof.
4. **Diff scope vs plan** - Compare the change set to `plan.json` tasks and `spec.md` acceptance. Flag out-of-plan work, missing acceptance coverage, or "done" without matching fresh evidence.
5. **Hunt (ready path)** - Search these targets only (do **not** re-walk the full reviewer dimension table; reviewer already applied `mission-review-gates` / UX gates):
   - **oracle-tamper** - assertions removed/skipped/loosened/tautologies; frozen scenario expected literals edited to force green; coder edited tests without Commander + `decisions.md` note; `Scenario oracle change:` missing when oracle changed
   - **false completion** - "done"/"ready" while acceptance fails, evidence missing/stale, scope mismatch, or defects left unfixed
   - **AUTH** - outward push/deploy/publish/send without quoted `AUTH:`; ship/merge without `/sc-ship` gates
   - **leftover findings** - any severity in `review.json` (critical / important / minor) ⇒ `REFUTED`
   - **outcome-gate disposition** - product ready needs `Design-contract: complete` or `Design-contract skipped: docs/prose-only`; `Approved-scenarios:` freeze footer or `Approved-scenarios skipped: docs/prose-only`; static / diff-cov / mutation disposition (evidence labels or greppable lines such as `Static-analysis skipped:…` / `Mutation skipped:…`). Exact prefixes: `docs/mission-artifacts.md`. Missing disposition without skip/waive ⇒ `REFUTED`. Global 95–100% coverage as success bar ⇒ `REFUTED`
   - **hard-gated Test Ideas** (when present) - Neg/Overlooked (+ Strategy Top risk/Charter when mapped) without matching `acceptance[]` and without `Deferred test idea: <id> - <reason>` ⇒ `REFUTED`; claimed done without fresh evidence ⇒ `REFUTED`
   - **product-surface miss** (when UI/workflow claimed) - need `verify.product` | `browser` | `curl` | `composition`; unit-only ⇒ `REFUTED`
   - **draft drift (visual UI only)** - when `UI draft approved:` is recorded, REFUTE clear chrome divergence, missing paired draft+live screenshot evidence, or fail/uncertain draft-parity / live-product. Point to `ux-ui-review-gates.md`; do not re-score every visual dimension here
6. **Verdict**
   - `VERIFIED` - fresh evidence passes; scope matches; hunts clean; **0** review findings. Ready allowed (subject to other gates).
   - `REFUTED` - any material hunt hit, failed re-run, scope/acceptance mismatch, leftover findings, or failed verify. Ready blocked.
7. **Ready gate** - Allow `ready` **only** on `VERIFIED`. On `REFUTED`, block ready and list remediation for `/sc-run` → fix → re-review → re-judge. No caveat / soft-pass.

### Edge cases

No claimed evidence / re-run fails → `REFUTED`. Non-defect `decisions.md` notes may sit beside `VERIFIED`; unfinished follow-up → `REFUTED`. Uncertain hard-gate or visual ready → fail-closed `REFUTED`. Manual-only → fresh observation note; do not invent output.

## Verdict contract

| Verdict | Ready |
|---------|-------|
| `VERIFIED` | **Allowed** |
| `REFUTED` | **Blocked** - fix → re-judge |

No aliases (`PASS`, `FAIL`, `VERIFIED WITH CAVEATS`, etc.).

## Rules

- **Must**: Re-run claimed evidence into `evidence.jsonl`; diff scope vs plan/spec; hunt oracle-tamper, AUTH, leftover `review.json` findings, A/B/C disposition per `docs/mission-artifacts.md`.
- **Must**: Emit `VERIFIED` | `REFUTED` only; allow `ready` only on `VERIFIED`; prove from evidence + scope + hunts + empty findings + `validate --strict` - not canvas.
- **Must not**: Re-walk the full reviewer dimension table; soft-pass `REFUTED`; invent evidence.
- **Must**: ASCII hyphen-minus; Cursor-native.

## Judge-break fixtures

`make test-judge-break` (or `scripts/check-judge-break.sh`) - known-bad packs under `references/judge-break/`. Run before claiming sc-judge or closeout predicate changes are safe.

## Out of scope

Product code/tests · clarify/drafts · AFK build · ship · trap-eval / LLM-as-judge transcript scoring.

## Output format

```
## Judge summary
Mission: <id> | Claims: <list>
Evidence re-run: <fresh ids> | Scope vs plan: <match | mismatches>
Hunt: oracle-tamper / AUTH / leftover findings / outcome-gate / hard-gated Test Ideas
Remediation (REFUTED): <none | list>
VERDICT: VERIFIED | REFUTED | Ready: allowed | blocked
```

## References

- `sc-verification`, `/sc-run`, `plan.json` / `spec.md`, `evidence.jsonl` / `review.json`
- `docs/mission-artifacts.md` - outcome-gate skip/waive SoT (`Design-contract`, `Approved-scenarios`, `Static-analysis`, `Mutation skipped`, …)
- `.cursor/skills/sc-run/references/mission-review-gates.md` - reviewer applies; judge does not re-walk
- `.cursor/skills/sc-ux-design/references/ux-ui-review-gates.md` - visual fail-closed when UI draft approved
- `references/judge-break/` - known-bad fixtures
