# Mission review gates

Authoritative harness reference for evidence, scope, tests, and acceptance review in spacecraft missions. Harness process only - not a product feature.

## Goal / when to use

Use this protocol whenever mission quality (evidence, validate, scope, tests, acceptance, security/perf when in scope, authorization) could influence pass/fail for `/sc-run` review, `sc-reviewer` findings, `sc-judge` ready, or ship readiness.

- **`/sc-run`** - deterministic pre-review before `Task(sc-reviewer)`; evidence, validate, scope, tests, acceptance, security/perf when in scope
- **Review** - `sc-reviewer` consumes per-dimension verdicts for all missions (mission dimensions always)
- **Judge** - `sc-judge` hunts false completion, weakened tests, scope drift, and unauthorized action on ready claims

Missions with **visual UI** also apply `.cursor/skills/sc-ux-design/references/ux-ui-review-gates.md` (sibling). Apply **both** when both scopes apply. Mission review covers every mission; UX gates add visual dimensions.

Do not replace `sc-verification`, `sc-tdd`, `sc-judge`, or existing defect-finding schema. This raises the bar by encoding how to run them before taste.

## The five gates

1. **Deterministic first** - Prefer `spacecraft validate --strict`, evidence re-run, done-task evidence checks, project-documented read-only audit scripts, and composition contracts before LLM critique. Reserve taste for what machines cannot assert directly (acceptance nuance, security heuristics, perf impact judgment).

2. **Narrow questions per dimension** - Never one blob question ("is this good?"). Separate evidence freshness, validate strict, scope vs plan, test quality, acceptance behavior, security (when in scope), performance (when in scope), and unauthorized action. One pass/fail verdict per dimension with a short reason.

3. **Pass/fail (not scores)** - Verdicts are `pass`, `fail`, or `uncertain` per dimension. Critique notes may say "uncertain" for humans. **fail-closed for ready/ship:** `uncertain` on a required dimension is treated as `fail` (critical finding, `REFUTED`, blocked ready) - never soft-pass or `VERIFIED` on uncertain mission-review claims.

4. **Human calibration** - When validating a new rubric or changing gate wording, sample representative missions and compare agent findings to human labels. Disagreements mean unclear criteria - tighten the dimension question or machine check; document in `decisions.md` or harness docs. Not a labeling product.

5. **Recheck on change** - When model, task shape, or criteria change, re-run the gates. Prior "passed" may drift. Do not inherit old review approvals without fresh deterministic pre-review and evidence.

## Dimension table

Required dimensions depend on mission scope. Mark **required** when the scope column applies.

| Dimension | Run (ready path) | Machine-checkable first | LLM / human critique |
|-----------|------------------|-------------------------|----------------------|
| **evidence-fresh** | required | Done tasks have matching `evidence.jsonl` lines; re-run `spacecraft evidence` / claimed verify commands for acceptance claims | Stale or missing labels; hand-written output; claims without fresh observation |
| **validate-strict** | required | `spacecraft validate --strict` passes before reviewer taste and before `ready` | Schema/artifact drift not caught by validate |
| **scope-vs-plan** | required | Diff file list vs `plan.json` tasks; acceptance strings vs `spec.md` | Silent extras; missing acceptance coverage; done tasks without matching work |
| **test-quality** | required when tests exist / TDD path | No removed/skipped/tautology assertions; composition contracts when FE/BE apply | Weakened expectations; GREEN without behavioral assertion |
| **acceptance-behavior** | required | Evidence proves behavior (not config-only); re-run acceptance verify commands | Wrong behavior with green tests; config-only proof |
| **approved-scenarios** | required on product path | `approved-scenarios.md` has freeze footer (`Approved-scenarios: frozen-from-contract` or `frozen-by-human`) or `Approved-scenarios skipped: docs/prose-only`; frozen expected literals not silently edited | Missing freeze; thawed oracles; scenarios invent values not in design-contract/spec |
| **static-analysis** | required on product path | `evidence.jsonl` has `static-…` label with **0 warning / 0 error** when a project static tool runs **or** greppable `Static-analysis skipped: no project static tool` / `Static-analysis waived: <reason>` in `decisions.md`; failures fixed or waived | Lint/typecheck ignored with no skip/waive; warnings/errors left open without waive |
| **diff-coverage** | required on product path | `diff-cov-…` evidence showing touched executable **line and branch ≥90%** (sanity band 90–95%) **or** `Diff-coverage skipped: no project coverage tool` / `Diff-coverage waived: <reason>`; never global 95–100% as the bar | Missing attempt; below 90% line+branch without waive; tautology padding for coverage |
| **mutation** | required disposition on product path | `mutation-…` evidence (**>80%** scoped, or project higher bar) when in scope + tool present; else greppable `Mutation skipped: not in scope` / `Mutation skipped: no project mutation tool` / `Mutation waived: <reason>` / opt-in via `Mutation: required` | Silent omit; in-scope without evidence or skip; score below target without waive; inventing mutator installs without ask |
| **pbt** | required disposition on product path | **100%** of design-contract **core-logic** modules (branching business rules / pure domain / state machines) have `pbt-…` evidence (invariants + generators via project-existing `fast-check` / Hypothesis / equivalent) **or** greppable `Pbt skipped: no project pbt tool` / `Pbt skipped: not core logic` / `Pbt waived: <reason>` | Silent omit on core-logic; missing disposition without skip/waive; inventing PBT lib install mid-mission |
| **security-when-in-scope** | when auth/API/secrets/deps touched | Commander captures read-only project checks as evidence when available (lint, typecheck, project-documented audit scripts); then heuristic `Task(sc-security)` scan | Heuristic gaps; patterns machines miss. **Boundary:** `sc-security` forbids dynamic CVE tools - do not require them here |
| **perf-when-in-scope** | when perf-touched paths / hot paths | Measure or documented baseline when `sc-performance` applies; existing benchmark in `evidence.jsonl` | `measure-first` / unclear impact without evidence |
| **unauthorized-action** | required | No outward push/deploy/publish/send without quoted `AUTH:`; no ship without `/sc-ship` gates | Outward action without authorization; merge/tag without lifecycle gates |

**security-when-in-scope fail-closed:** When auth/API/secrets/deps are in scope, uncertain or skipped with no machine evidence and no heuristic `sc-security` pass ⇒ treat as `fail` (critical / `REFUTED`).

**perf-when-in-scope fail-closed:** When perf paths are in scope, `measure-first` unclear impact without evidence ⇒ `uncertain` ⇒ fail-closed for ready.

## Deterministic-before-review ordering

Commander runs this layer in `/sc-run` **before** `Task(sc-reviewer)`:

1. `spacecraft validate --strict`
2. Confirm every `plan.json` task marked `done` has matching `evidence.jsonl` entries for its `evidence` labels
3. Confirm approved-scenarios freeze footer or docs/prose skip; confirm static-analysis, diff-coverage, mutation, and PBT (`pbt-…` or `Pbt skipped`/`Pbt waived`) evidence labels or skip/waive lines (`docs/mission-artifacts.md`)
4. Re-run or spot-check claimed verify commands when acceptance is in doubt
5. When security in scope: capture read-only project checks as evidence (lint, typecheck, documented audit scripts); then `Task(sc-security)` heuristic scan (no dynamic CVE tools per `sc-security` skill)
6. When performance in scope: capture measurement or documented baseline per `sc-performance`; flag unclear hot-path impact without evidence
7. When UI or multi-step workflow touched: Task(`sc-browser-probe`) to `PROBE: CLEAN` (skip when no runnable UI/workflow surface)
8. Only then: `Task(sc-reviewer)` (+ `Task(sc-designer)` / UX gates when visual UI)

## Post-review canvas handoff

After `review.json` is written, Commander may emit canvases for human check under `~/.cursor/projects/<workspace>/canvases/`:

1. **Findings:** nonempty `findings` → optional `<missionId>-findings.canvas.tsx`, `Canvas findings: ` + absolute path in `decisions.md`, and an absolute markdown link in chat. Empty `findings` → optional `Canvas findings skipped: empty` (no findings canvas file).
2. **Evidence:** optional `<missionId>-evidence.canvas.tsx`; `Canvas evidence: ` + absolute path; absolute markdown link in chat (and `decisions.md`).

Then run `sc-judge`. Missing canvas files or `decisions.md` lines do not block `sc-judge` or ready. Ready proof is `evidence.jsonl` + empty `review.json` findings + `validate --strict` + judge `VERIFIED`.

Do not put canvases under mission `.space/` or repo `.cursor/`. Do not inspect canvas TSX/JSON shape. See `/sc-run` Optional canvas (human check).

## Verdict mapping

Per dimension, emit exactly one of:

| Verdict | Meaning | Ready / ship |
|---------|---------|--------------|
| `pass` | Dimension satisfied with evidence | Allows progress on that dimension |
| `fail` | Dimension not satisfied | Blocks ready on that dimension |
| `uncertain` | Cannot decide with available evidence | **fail-closed:** treat as `fail` for `ready`, ship, and `sc-judge` |

Rules:

- No 1-5 scores or numeric rubrics.
- Critique output may still label `uncertain` in notes for human follow-up.
- `sc-judge`: `uncertain` on any required mission-review dimension ⇒ `REFUTED` (note `uncertain` in hunt reasons). No third verdict string.
- `sc-reviewer`: `uncertain` on a required dimension ⇒ critical finding; `status: blocked`.

## Output snippet (reviewers)

Short per-dimension lines only:

```
Mission review (mission-review-gates):
- evidence-fresh: pass | fail | uncertain - <one line reason>
- validate-strict: pass | fail | uncertain - <reason>
- scope-vs-plan: pass | fail | uncertain - <reason>
- test-quality: pass | fail | uncertain | n/a - <reason>
- acceptance-behavior: pass | fail | uncertain - <reason>
- security-when-in-scope: pass | fail | uncertain | n/a - <reason>
- perf-when-in-scope: pass | fail | uncertain | n/a - <reason>
- unauthorized-action: pass | fail | uncertain - <reason>
```

Omit `n/a` dimensions or mark `n/a` when out of scope. Required dimensions for the current scope must not be `uncertain` on the ready path.

When visual UI also applies, add the UX snippet from `ux-ui-review-gates.md` separately.

## Human calibration (gate 4)

Lightweight harness process:

1. Pick 3-5 representative missions (pass, fail, edge: missing evidence, weakened tests, scope drift).
2. Run deterministic pre-review + `sc-reviewer` with this table; record per-dimension verdicts.
3. Compare to human labels. Mismatches ⇒ tighten the dimension question or add a machine check.
4. Record outcome in `decisions.md` when rubric wording changes (e.g. `Mission review rubric calibrated: <date> - <note>`).

## Recheck on change (gate 5)

Re-run relevant gates when any of these change:

- Subagent model or reviewer prompt
- Task shape (new acceptance, security/perf scope added)
- Dimension table or deterministic pre-review ordering

Prior review approval does not grandfather later runs without fresh validate, evidence, and judge re-run.

## Cross-links

- `.cursor/skills/sc-run/SKILL.md` - AFK loop; deterministic pre-review before reviewer
- `.cursor/skills/sc-verification/SKILL.md` - fresh evidence capture mechanics
- `.cursor/skills/sc-judge/SKILL.md` - adversarial prove gate; `VERIFIED` | `REFUTED` only
- `.cursor/skills/sc-run/references/defect-finding.md` - actionable findings schema
- `.cursor/skills/sc-tdd/SKILL.md` - test-first discipline; test-quality dimension
- `.cursor/skills/sc-security/SKILL.md` - heuristic static scan; **no dynamic CVE tools**
- `.cursor/skills/sc-performance/SKILL.md` - measure-first; hot-path discipline
- `.cursor/skills/sc-ux-design/references/ux-ui-review-gates.md` - sibling for visual UI
- `docs/mission-review.md` - short human-facing overview
