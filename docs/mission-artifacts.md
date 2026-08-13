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

## Outcome-gate skip / waive grammar (SoT)

Lifecycle skills (`sc-run`, `sc-judge`, `sc-tdd`, `sc-verification`, `sc-planning`, `sc-mission`) cite this section for greppable `decisions.md` line prefixes - do not redefine the strings elsewhere. When each gate applies and what evidence to capture: `design-contract.md` / `approved-scenarios.md` above and **Outcome evidence labels** / **Mutation** below.

| Gate | Greppable line(s) |
|------|-------------------|
| design-contract | `Design-contract skipped: docs/prose-only` |
| approved-scenarios | `Approved-scenarios skipped: docs/prose-only` |
| static-analysis | `Static-analysis skipped: no project static tool` · `Static-analysis waived: <reason>` (rare; reason required) |
| diff-coverage | `Diff-coverage skipped: no project coverage tool` · `Diff-coverage waived: <reason>` (e.g. docs-only touch, generated files only) |
| mutation | `Mutation skipped: no project mutation tool` · `Mutation skipped: not in scope` · `Mutation waived: <reason>` (rare; reason required) |

## Outcome evidence labels (layer B)

Capture with `spacecraft evidence` when the gate runs:

| Gate | Typical label prefix | When |
|------|----------------------|------|
| Static analysis | `static-` | combine and/or deterministic pre-review - project lint/typecheck/format-check if present |
| Diff coverage | `diff-cov-` | combine or pre-review - coverage on **touched** lines/files when a project coverage tool exists |
| Mutation testing | `mutation-` | combine or pre-review when mutation is **in scope** (see below) and a project mutation tool exists |

Skip / waive lines: exact prefixes in **Outcome-gate skip / waive grammar (SoT)** above.

**Must not** use global line coverage 95–100% as a hard gate. Diff coverage target when measured: **≥80% of touched executable lines** in the mission diff (exclude pure docs/generated unless the mission is about those files). Below target → fix tests or waive with reason - do not pad tautologies.

### Mutation (phase C)

**In scope when any hold:**

1. `decisions.md` has greppable `Mutation: required`, **or**
2. Mission product path touches logic-heavy code (branching business rules, pure domain/calc/parse pipelines) per `design-contract.md` Modules / Data structures - not docs/prose-only, not pure styling/chrome-only

**Out of scope:** docs/prose-only plans; UI look-only with no behavioral logic; no mutation tool in the project.

When in scope and a mutation tool exists: run it scoped to touched packages/files when the tool allows; capture `spacecraft evidence "mutation-…" -- <cmd>`. Default target: **≥70% mutation score** on that scope (or the project's documented higher bar). Below target → strengthen behavioral tests (not tautologies) or `Mutation waived: <reason>`.

When in scope but no tool: `Mutation skipped: no project mutation tool` (does not invent installing a mutator mid-mission unless the human asked).

When not in scope: `Mutation skipped: not in scope` (one line; required so ready is unambiguous).

Do not treat mutation as a substitute for approved-scenarios / design-contract oracles.

## Other

- `spec.md` - what and why (free-form markdown)
- `decisions.md` / `questions.md` - assumptions and blocking questions
- `review.md` / `review.json` - release readiness review output (after `sc-judge`; ready only when verdict is `VERIFIED` and findings empty). Mission review: `.cursor/skills/sc-run/references/mission-review-gates.md` (overview `docs/mission-review.md`). Visual UI: UX/UI review gates (`.cursor/skills/sc-ux-design/references/ux-ui-review-gates.md`; overview `docs/ux-ui-review.md`)

## Inner-loop / judge

Quick and Mission both apply INTENT / AUTH / TWINS and the 3-cycle stop. Capture judge evidence with a label that includes `judge` (e.g. `judge-<mission-id>` or `judge-pass-validate`). `output` must be real command stdout - never hand-written. AUTH does not bypass ship hooks or `SPACECRAFT_SHIP=1`. Findings are fixed in `/sc-run` and listed in the run/ship summary.