# Mission review gates

Authoritative harness reference for evidence, scope, tests, and acceptance review in spacecraft missions. Harness process only - not a product feature.

## Goal / when to use

Use this protocol whenever mission quality (evidence, validate, scope, tests, acceptance, security/perf when in scope, authorization) could influence pass/fail for `/sc-run` review, `sc-reviewer` findings, `sc-judge` ready, or ship readiness.

- **`/sc-run`** - deterministic pre-review before `Task(sc-reviewer)`; evidence, validate, scope, tests, acceptance, security/perf when in scope; Cursor `bugbot` + `security-review` ingest into `review.json`
- **Review** - `sc-reviewer` adds only mission-dimension gaps Cursor does not cover; consumes per-dimension verdicts for all missions (mission dimensions always)
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
| **approved-scenarios** | required on product path | `approved-scenarios.md` has freeze footer (`Approved-scenarios: frozen-from-contract` or `frozen-by-human`) or `Approved-scenarios skipped: docs/prose-only`; `spacecraft freeze-check` passes when freeze is machine-required (Gates ≥ M9G7IHV3); frozen expected literals not silently edited | Missing freeze when required; thawed oracles; scenarios invent values not in design-contract/spec; postdated-freeze or freeze-drift |
| **static-analysis** | required on product path | `evidence.jsonl` has `static-…` with **Full package/project static suite required (0 warning / 0 error)** when a project static tool runs **or** `Static-analysis skipped: no project static tool` / `Static-analysis waived: <reason>`; fix or waive failures. Tip-path-only lint/typecheck MUST NOT satisfy static-analysis for ready | Tip-path-only as sole static proof; ignored lint/typecheck without skip/waive; open warnings/errors without waive |
| **diff-coverage** | required on product path | `diff-cov-…` evidence showing touched executable **line and branch ≥90%** (sanity band 90–95%) **or** `Diff-coverage skipped: no project coverage tool` / `Diff-coverage waived: <reason>`; never global 95–100% as the bar | Missing attempt; below 90% line+branch without waive; tautology padding for coverage |
| **mutation** | required disposition on product path | In scope when any of: greppable `Mutation: required`, pack id `quality`, or greppable `Mutation: high-risk` (SoT: `docs/mission-artifacts.md`). Then `mutation-…` evidence (**>80%** scoped, or project higher bar) when tool present; else greppable `Mutation skipped: not in scope` / `Mutation skipped: no project mutation tool` / `Mutation waived: <reason>`. Ordinary missions: `Mutation skipped: not in scope` is valid | Silent omit; in-scope without evidence or skip; score below target without waive; inventing mutator installs without ask |
| **pbt** | required disposition on product path | **100%** of design-contract **core-logic** modules (branching business rules / pure domain / state machines) have `pbt-…` evidence (invariants + generators via project-existing `fast-check` / Hypothesis / equivalent) **or** greppable `Pbt skipped: no project pbt tool` / `Pbt skipped: not core logic` / `Pbt waived: <reason>` | Silent omit on core-logic; missing disposition without skip/waive; inventing PBT lib install mid-mission |
| **security-when-in-scope** | when auth/API/secrets/deps touched | Evaluate declared **SEC** Verify bars from `spec.md` (tool + evidence label; `NFR source:` when present). Default hard bar when discuss set one: `no new critical/high SAST findings vs baseline`. Commander captures read-only project checks / SAST as machine evidence when available (ordering step 5). Cursor `Task(security-review)` once in ready-path ordering step 8 (parallel with bugbot). `Task(sc-security)` only on Cursor subagent failure / unavailable skip fallback or explicit on-demand heuristic | Heuristic gaps; patterns machines miss; missing declared SEC bar when in scope. **Boundary:** no dynamic CVE tools - do not require them here |
| **perf-when-in-scope** | when perf-touched paths / hot paths | Evaluate declared **PERF** Verify bars or recorded-skip debt when perf in scope. Default relative bar when discuss set one: `no p95 regression >10% vs baseline <evidence-id>`. Measure or documented baseline when `sc-performance` applies; existing benchmark in `evidence.jsonl`; accept greppable `<Gate> skipped: no tool` debt when no tool | `measure-first` / unclear impact without evidence; invented bars; skip without debt line |
| **unauthorized-action** | required | No outward push/deploy/publish/send without quoted `AUTH:`; no ship without `/sc-ship` gates | Outward action without authorization; merge/tag without lifecycle gates |

**security-when-in-scope fail-closed:** When auth/API/secrets/deps are in scope, uncertain or skipped with no machine evidence and no `security-review` pass (and no greppable `Sc-security fallback: pass` when Cursor failed or on-demand) ⇒ treat as `fail` (critical / `REFUTED`). Fail-closed when security is in scope and both `security-review` and fallback `sc-security` fail. When security is in scope, greppable `Cursor review skipped:` is valid **only** after a greppable `Sc-security fallback: pass` (optionally also `Sc-security fallback: findings drained`; optional evidence label `sc-security-…`) **or** SEC machine-evidence pass; otherwise fail / `REFUTED` - skip alone does not satisfy security-when-in-scope. Declared SEC Verify bars (tool + evidence label) are required when security is in scope - do not invent a new review dimension.

**perf-when-in-scope fail-closed:** When perf paths are in scope, `measure-first` unclear impact without evidence ⇒ `uncertain` ⇒ fail-closed for ready. Declared PERF Verify bars or recorded-skip debt are required when perf is in scope - do not invent a new review dimension.

## Deterministic-before-review ordering

Commander runs this layer in `/sc-run` **before** `Task(sc-reviewer)`:

1. `spacecraft validate --strict`
2. Confirm every `plan.json` task marked `done` has matching `evidence.jsonl` entries for its `evidence` labels
3. Confirm approved-scenarios freeze footer or docs/prose skip; run `spacecraft freeze-check` when Gates version ≥ M9G7IHV3 and freeze is machine-required (`docs/mission-artifacts.md` **Test freeze**); confirm static-analysis, diff-coverage, mutation, and PBT (`pbt-…` or `Pbt skipped`/`Pbt waived`) evidence labels or skip/waive lines (`docs/mission-artifacts.md`)
4. Re-run or spot-check claimed verify commands when acceptance is in doubt
5. When security in scope: evaluate declared SEC Verify bars (tool + evidence label / SAST vs baseline); capture read-only project checks as **machine evidence only** (lint, typecheck, documented audit scripts / SAST). Do **not** invoke `Task(security-review)` here
6. When performance in scope: evaluate declared PERF Verify bars or recorded-skip debt; capture measurement or documented baseline per `sc-performance` (e.g. p95 relative-bar); flag unclear hot-path impact without evidence
7. When UI or multi-step workflow touched: Task(`sc-browser-probe`) to `PROBE: CLEAN` (skip when no runnable UI/workflow surface). Primary probe executor: Cursor IDE browser (`cursor-ide-browser`) via sc-browser-probe Browser matrix. Browser / MCP / chat success Must not authorize `ready` / `VERIFIED` / `AUTH` / ship.
8. Parallel Cursor `Task(bugbot)` + `Task(security-review)` - **sole** ready-path `Task(security-review)` invoke (disposition per `/sc-run`); map outputs into `review.json` via `defect-finding.md` with `source` set. On Cursor security-review failure or unavailable: `Task(sc-security)` fallback or SEC machine-evidence pass before skip is valid when security in scope - write greppable `Sc-security fallback: pass` (optionally also `Sc-security fallback: findings drained`; optional evidence label `sc-security-…`) when fallback clears the gate (no dynamic CVE tools; fail-closed when both fail)
9. Only then: `Task(sc-reviewer)` (+ `Task(sc-designer)` / UX gates when visual UI) - mission-dimension gaps only; on overlap Cursor finding wins - **remove** Spacecraft duplicate from `findings[]` (do not leave a `supersededBy` row that still counts as leftover)

## Cursor ingest and all-severity drain

After Cursor + `sc-reviewer` write `review.json`:

1. **Ingest** - every bugbot / security-review / sc-reviewer finding is a defect-finding row with `source` set when known. When corroborating a `Cursor review: … ran` disposition without a `cursor-review-…` evidence label, write greppable `Cursor ingest: session` in `decisions.md`.
2. **Dedupe** - same file + issue-family → keep Cursor row; **remove** Spacecraft duplicate from `findings[]` (`supersededBy` is not a ready exemption - row must be gone).
3. **Drain all severities** (`critical` / `important` / `minor`) - fix only when `requiredFix` is concrete (no free-form whole-tree refactor from Bugbot noise); combine keeps ordinary refactor. Loop: fix → re-run Cursor reviews → re-review → re-judge until findings empty or `3-cycle` / `timebox` handback.

Ready proof still requires empty `review.json` findings - soft-pass with open minors is forbidden.

## Post-review → judge

After `review.json` is written, Commander continues to `sc-judge`. Ready proof is `evidence.jsonl` + empty `review.json` findings + `validate --strict` + judge `VERIFIED`. Companion-lane dispositions (`Loop watch:`, `Post-ready drain:`, `Split-to-prs:`) are required by `/sc-run` / `/sc-ship` gates — never AUTH / `VERIFIED` / ship authority. See `references/optional-lanes.md`.

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
- `.cursor/skills/sc-security/SKILL.md` - fallback / on-demand heuristic static scan when Cursor `security-review` fails or is explicitly requested; **no dynamic CVE tools**
- `.cursor/skills/sc-performance/SKILL.md` - measure-first; hot-path discipline
- `.cursor/skills/sc-ux-design/references/ux-ui-review-gates.md` - sibling for visual UI
- `docs/review.md` - short human-facing overview
