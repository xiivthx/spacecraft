# Mission artifact schemas

Canonical shapes for `.space/` artifacts. Always-on rules point here; do not paste full schemas into every prompt.

## mission.json

```json
{
  "id": "<mission-id>",
  "title": "<title>",
  "state": "active|planned|in_progress|ready|blocked|shipped",
  "branches": ["<branch-name>"],
  "createdAt": "<iso-date>",
  "pickup": {
    "phase": "<discuss|run|ship|quick|…>",
    "next": "<one-liner for the next session>",
    "blockers": ["<optional blocker>"],
    "updatedAt": "<iso-date>"
  }
}
```

Optional `pickup`: set or update on discuss / run / ship / quick handoffs so `spacecraft status` (and session-start) can print `Pickup: <next>`. Fields: `phase`, `next` (required for the status line), optional `blockers`, and `updatedAt`. Omit `pickup` when there is nothing useful to resume. Not a closeout or ship gate.

## plan.json

Jigsaw tasks: each task is one behavioral slice of the feature. Each `acceptance[]` item is one RED-GREEN cycle under `/sc-run`.

```json
{
  "planName": "<short-name>",
  "missionId": "<mission-id>",
  "tasks": [
    {
      "id": "T1",
      "title": "<imperative jigsaw slice>",
      "status": "pending|in_progress|done|blocked",
      "dependsOn": [],
      "files": ["<path>"],
      "acceptance": ["<one check per RED-GREEN cycle>"],
      "verify": "<command>",
      "evidence": ["<label>"]
    }
  ]
}
```

When `Sizing: phases` (discuss-owned), same-mission overflow uses sequential files `plan-phase1.json`, `plan-phase2.json`, … (or `plan-phaseN.json` naming). Each file matches the `plan.json` shape above; ≤7 tasks per phase. Default active plan for build is the current incomplete phase file (Commander tracks which phase is active in `decisions.md` or by first pending phase).

## evidence.jsonl

One JSON object per line:

```json
{"label":"<label>","command":"<command>","output":"<output>","outputHash":"<hex>","ts":"<iso-timestamp>"}
```

`output` must be actual captured command output - never fabricated.

`outputHash` is optional: lowercase hex SHA-256 of the **full** raw command output. Omitted via `omitempty`; entries without it remain valid (backward compatible).

When raw output exceeds the capture limit (65536 bytes), `spacecraft evidence` keeps the full raw under the mission `evidence-raw/` sidecar, truncates the JSONL `output` field (prefix plus trailing marker `\n...[truncated]`), and records:

```json
{"label":"<label>","command":"<command>","output":"<prefix>\\n...[truncated]","outputHash":"<hex of full raw>","outputTruncated":true,"outputBytes":<full-raw-byte-length>,"outputRawPath":"evidence-raw/<ts>-<safe-label>.log","ts":"<iso-timestamp>"}
```

- `outputTruncated` - `true` when the JSONL `output` is truncated
- `outputBytes` - byte length of the full raw
- `outputRawPath` - path relative to the mission dir under `evidence-raw/`
- `outputHash` - SHA-256 of the full raw (not of the truncated JSONL string)

`spacecraft validate` hashes the sidecar when `outputTruncated` is true (fails if the sidecar is missing or mismatches); otherwise it hashes `output` when `outputHash` is present. The terminal still prints the full raw on capture.

## Roadmap (`.space/roadmaps/<id>.json`)

```json
{
  "id": "<kebab-slug>",
  "title": "<title>",
  "description": "<description>",
  "missions": [
    {"id": "<mission-id>", "description": "<brief>"}
  ],
  "issues": [
    {"number": 1, "title": "...", "url": "https://...", "state": "open", "phase": "mvp|polish|ship"}
  ],
  "createdAt": "<iso-date>",
  "updatedAt": "<iso-date>"
}
```

Roadmap ID: lowercase kebab-case from title. Manage with `spacecraft roadmap` (`map` alias).

## design-contract.md

Mission-scoped build contract written after `plan.json` and **before** AFK RED/GREEN. Path: `.space/missions/<mission-id>/design-contract.md`.

Purpose: force whole-feature analysis (modules, data shapes, public seams, edge literals) so tests and code share one frozen oracle - not design-by-first-test.

Required sections (headings may match exactly or clearly synonym):

1. `## Scope` - in-scope behavior from `spec.md` (short)
2. `## Modules` - modules / packages / files and responsibilities (scoped to this mission; not a product-wide architecture essay)
3. `## Data structures` - types, records, invariants
4. `## Public seams` - interfaces tests will call (APIs, CLI, exported functions)
5. `## Edge matrix` - table or list: case, input, **expected literal** (from spec / worked example - never "same as implementation"), acceptance id/ref when known
6. `## Out of scope` - Wrong-if / non-goals that bind build

Footer line (greppable): `Design-contract: complete`

**Skip (docs/prose-only plans):** when every pending acceptance is docs/prose/wording-only tautology, do not write the file; append the design-contract skip line from **Outcome-gate skip / waive grammar (SoT)** below to `decisions.md` instead.

`/sc-run` Must not start product RED/GREEN until `design-contract.md` exists with `Design-contract: complete`, or the skip line is present.

## approved-scenarios.md

Frozen behavior oracles for AFK confidence (Approved Scenarios pattern). Path: `.space/missions/<mission-id>/approved-scenarios.md`.

Write **after** `design-contract.md` and **before** product RED/GREEN. Derive rows from the Edge matrix + spec worked examples - do not invent new expected values.

Minimum shape:

1. Short provenance line (e.g. `Source: design-contract Edge matrix + spec worked examples`)
2. Table or list with: `id`, acceptance ref, input/steps, **expected literal**, `status` (`frozen` for rows that bind build)
3. Footer (greppable): `Approved-scenarios: frozen-from-contract` (AFK freeze from contract) **or** `Approved-scenarios: frozen-by-human` when a human explicitly froze/re-froze at discuss or ready check

**Rules:**

- Agents **Must not** edit frozen expected literals to make tests pass. Oracle changes need Commander + greppable `Scenario oracle change: <id> - <reason>` in `decisions.md` (and usually human confirm at ready).
- Product RED/GREEN **Must not** start until this file exists with a freeze footer, or skip applies.
- **Skip (docs/prose-only):** when the design-contract skip applies, also append the approved-scenarios skip line (no file required). Exact prefixes: **Outcome-gate skip / waive grammar (SoT)** below.

This is not a separate product test runner - executable proof stays in the project test suite / verify commands mapped to scenario ids via acceptance.

### Test freeze (machine-checked D4)

Product missions with `Gates version: M9G7IHV3` (or later in the ordered registry) get machine-checked test freeze via `spacecraft freeze` and `spacecraft freeze-check`. Pre-tip missions without a valid `Gates version:` line are grandfathered (no new freeze gates).

**Capture:** `spacecraft freeze [--mission <id>] <path...>` — one append-only evidence event with label `freeze`. The `output` field holds JSON:

```json
{"kind":"freeze-manifest","files":[{"path":"<repo-relative>","sha256":"<hex>"}]}
```

Paths are explicit (no glob), repo-relative, forward slashes. When the manifest JSON exceeds 65536 bytes, the CLI truncates `output`, sets `outputTruncated` + `outputRawPath` sidecar under `evidence-raw/` (same pattern as `spacecraft evidence`).

**Required-ness:**

- **Machine-required** when `approved-scenarios.md` has a frozen footer (`Approved-scenarios: frozen-from-contract` or `frozen-by-human`) and Gates version ≥ M9G7IHV3.
- **Exempt** when `Approved-scenarios skipped:` appears in `decisions.md` or `approved-scenarios.md`.
- **Skill-gated** when neither frozen footer nor skip line exists (no new machine net).

**Checks (`spacecraft freeze-check`, also enforced in `closeout-check`):**

- `freeze-drift` — manifest file hash differs from disk without greppable `Scenario oracle change:` in `decisions.md` (re-freeze after oracle change).
- `postdated-freeze` — latest `freeze` evidence line appears after any `test-…` or `test-run` evidence line in the append-only log (retroactive freeze).
- `Cross-model critic:` or `Cross-model critic skipped:` required in `decisions.md` when Gates version ≥ M9G7IHV3 (`silent-cross-model-critic` at closeout).

Run `spacecraft freeze-check` at combine/ready before product combine proceeds; exit 1 blocks ship until drift/postdate/missing freeze is resolved.

## Outcome-gate skip / waive grammar (SoT)

Lifecycle skills (`sc-run`, `sc-judge`, `sc-tdd`, `sc-verification`, `sc-planning`, `sc-mission`) cite this section for greppable `decisions.md` line prefixes - do not redefine the strings elsewhere. When each gate applies and what evidence to capture: `design-contract.md` / `approved-scenarios.md` above and **Outcome evidence labels** / **Mutation** below.

| Gate | Greppable line(s) |
|------|-------------------|
| design-contract | `Design-contract skipped: docs/prose-only` |
| approved-scenarios | `Approved-scenarios skipped: docs/prose-only` |
| static-analysis | `Static-analysis skipped: no project static tool` · `Static-analysis waived: <reason>` (rare; reason required) |
| diff-coverage | `Diff-coverage skipped: no project coverage tool` · `Diff-coverage waived: <reason>` (e.g. docs-only touch, generated files only) |
| mutation | `Mutation skipped: no project mutation tool` · `Mutation skipped: not in scope` · `Mutation waived: <reason>` (rare; reason required) |
| PBT | `Pbt skipped: no project pbt tool` · `Pbt skipped: not core logic` · `Pbt waived: <reason>` |
| quality / NFR (no tool) | `<Gate> skipped: no tool` — quality debt when a declared SEC/PERF (or other NFR) gate cannot run because no project tool exists; pair with a debt line in `decisions.md` (never invent a passing bar). Debt-facing mutation synonym: `Mutation skipped: no tool` (same class as SoT `Mutation skipped: no project mutation tool` when mutation is in scope but no tool) |
| characterization | `Characterization waived: <reason>` — recorded quality debt on refactor/optimize when behavior-preservation helpers are unavailable (see **Quality debt and characterization** below) |

## NFR provenance and relative bars

Quality / NFR Verify lines (security, performance, and similar) carry provenance and stay measurable.

**Provenance tag (greppable):**

```
NFR source: user | measured-baseline | default(<why>)
```

- `user` - human-stated bar
- `measured-baseline` - bar derived from captured baseline evidence
- `default(<why>)` - harness default with a short why (e.g. `default(SAST critical/high vs baseline)`)

Unknown preference-bound bars → sc-clarify blocking question; do not invent.

**Relative-bar pattern:** compare against a named baseline evidence id, not an absolute invented number. Greppable label: `relative-bar`. Example:

```
no p95 regression >10% vs baseline <evidence-id>
```

**No-tool skip line:** when the gate is in scope but no project tool exists, record:

```
<Gate> skipped: no tool
```

Example: `PERF skipped: no tool`. Each quality Verify line names tool + evidence label, or becomes this debt pattern - never checkbox theater. Mutation in-scope with no tool still uses the SoT skip `Mutation skipped: no project mutation tool`; debt accounting also greps the synonym `Mutation skipped: no tool` (see below).

## Quality debt and characterization

### Visible quality debt (D5-rule)

These skips are **visible quality debt**, not silent green:

- `Mutation skipped: no tool` (debt-facing synonym; SoT formal skip remains `Mutation skipped: no project mutation tool` when mutation is in scope but no tool)
- Equivalent quality / NFR no-tool skips: `<Gate> skipped: no tool` (SEC, PERF, and similar)
- `Characterization waived: <reason>` (refactor/optimize behavior-preservation waiver)

`Mutation skipped: not in scope` is a valid ready disposition, **not** quality debt.

Ship closeout **Must** list open quality debts for the mission (greppable skip/waive lines above that still count as debt).

**Debt ceiling (greppable):**

```
Debt ceiling: 3
```

N=3 open quality debts ⇒ the next mission **Must** drain quality debt or be a quality mission / fold into `*-integrate`. Record the ceiling line in `decisions.md` (or ship closeout) when the tip's grammar governs.

### Characterization / behavior preservation (C3)

Refactor and optimize missions carry a behavior-preservation bar scaled to blast radius (characterization / golden / baseline do-not-break evidence preferred). Characterization helpers are deferred to the quality-tooling seam; until then, explicit waiver is allowed as recorded debt:

```
Characterization waived: <reason>
```

Reason required. Counts toward the debt ceiling like other quality debt.

### Grandfathering

New quality-debt and characterization grammar applies only after tip `Gates version: M9G7IHON` ships. Do not retro-gate missions cleared before that ship. In-flight / pre-tip missions keep prior dispositions until their next discuss that adopts the tip.
## Outcome evidence labels (layer B)

Capture with `spacecraft evidence` when the gate runs:

| Gate | Typical label prefix | When |
|------|----------------------|------|
| Static analysis | `static-` | combine and/or deterministic pre-review - project lint/typecheck/format-check if present |
| Diff coverage | `diff-cov-` | combine or pre-review - coverage on **touched** executable line and branch when a project coverage tool exists |
| Mutation testing | `mutation-` | combine or pre-review when mutation is **in scope** (see below) and a project mutation tool exists |
| Property-based (PBT) | `pbt-` | product path: **100%** of design-contract **core-logic** modules (invariants + generators via project-existing `fast-check` / Hypothesis / equivalent) |
| Test freeze | `freeze` | before RED / at per-mission freeze event — sha256 manifest of explicit test file list + `approved-scenarios.md`; re-checked at combine/ready via `spacecraft freeze-check` |

Skip / waive lines: exact prefixes in **Outcome-gate skip / waive grammar (SoT)** above.

**Static analysis:** when a project static tool runs, require **0 warning / 0 error**. No tool → skip/waive per the table above.

**Diff coverage:** when measured, touched executable **line and branch ≥90%** (sanity band **90–95%**). Exclude pure docs/generated unless the mission is about those files. **Must not** chase whole-repo global coverage, and **Must not** pad tautologies toward 100%. Below target → fix tests or waive with reason.

### Mutation (phase C)

**In scope when any hold:**

1. `decisions.md` has greppable `Mutation: required`
2. Project pack selection includes pack id `quality` (profile / `SPACECRAFT_PACKS` / setup packs — same meaning as install “quality pack”)
3. `decisions.md` has greppable `Mutation: high-risk`

**Out of scope:** ordinary missions with no quality pack and without `Mutation: required` / `Mutation: high-risk` in `decisions.md`. For those, `Mutation skipped: not in scope` is a **valid** ready disposition (one line; required so ready is unambiguous).

When in scope and a mutation tool exists: run it scoped to touched packages/files when the tool allows; capture `spacecraft evidence "mutation-…" -- <cmd>`. Default target: **>80% mutation** score on that scope (or the project's documented higher bar). Below target → strengthen behavioral tests (not tautologies) or `Mutation waived: <reason>`.

When in scope but no tool: `Mutation skipped: no project mutation tool` (does not invent installing a mutator mid-mission unless the human asked). Same debt class as `Mutation skipped: no tool` under **Quality debt and characterization**.

When not in scope: `Mutation skipped: not in scope`.

Do not treat mutation as a substitute for approved-scenarios / design-contract oracles.

### Property-based testing (PBT)

**Core logic** (design-contract modules): branching business rules / pure domain / state machines - not chrome, docs, or thin adapters.

**Must (product path):** **100%** of design-contract modules marked **core-logic** require property-based invariants + generators, using a project-existing lib (`fast-check` / Hypothesis / equivalent). Capture `spacecraft evidence "pbt-…" -- <cmd>`.

**Must not:** invent installing a PBT library mid-mission unless the human asked.

Missing `pbt-…` evidence without a greppable skip/waive disposition below is **REFUTED** material for `sc-judge`.

When not applicable, one greppable disposition line (required so ready is unambiguous):

- `Pbt skipped: no project pbt tool`
- `Pbt skipped: not core logic`
- `Pbt waived: <reason>` (rare; reason required)

## Other

- `spec.md` - what and why (free-form markdown)
- `decisions.md` / `questions.md` - assumptions and blocking questions
- `review.md` / `review.json` - release readiness review output (after `sc-judge`; ready only when verdict is `VERIFIED` and findings empty). Mission review: `.cursor/skills/sc-run/references/mission-review-gates.md` (overview `docs/mission-review.md`). Visual UI: UX/UI review gates (`.cursor/skills/sc-ux-design/references/ux-ui-review-gates.md`; overview `docs/ux-ui-review.md`)

## Inner-loop / judge

Quick and Mission both apply INTENT / AUTH / TWINS and the 3-cycle stop. Capture judge evidence with a label that includes `judge` (e.g. `judge-<mission-id>` or `judge-pass-validate`). `output` must be real command stdout - never hand-written. AUTH does not bypass ship hooks or `SPACECRAFT_SHIP=1`. Findings are fixed in `/sc-run` and listed in the run/ship summary.