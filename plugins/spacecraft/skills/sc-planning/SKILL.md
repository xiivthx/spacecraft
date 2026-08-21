---
name: sc-planning
description: "Convert a mission spec into a small executable plan with verifiable tasks. Activate during /sc-run planning, Task(sc-planner), task planning, or spec decomposition."
---

# sc-planning

## Goal

Turn `spec.md` into a jigsaw `plan.json` with ≤7 verifiable tasks per phase, then mission `design-contract.md` + `approved-scenarios.md`, so `/sc-run` can execute per-acceptance RED-GREEN against a frozen oracle.

## Output

Writable `plan.json`, `design-contract.md`, and `approved-scenarios.md` (schemas in `docs/mission-artifacts.md`). Each task: acceptance + verify + evidence. Each acceptance = one RED-GREEN cycle. Design-contract holds modules, seams, edge literals; approved-scenarios freezes them before build.

## Good / Bad

- Good: atomic jigsaw slices; concrete acceptance (1-3 per task); exact verify; hard-gated Test Ideas covered or deferred; design-contract + scenarios complete or docs/prose skip
- Bad: one coarse "implement feature" task; vague titles; missing verify; silent gray-area fills; inventing Verify; `map new` from planning

## Verify

Every acceptance testable and one cycle; ≤7 tasks per phase; file paths real; no open blocking clarify; each hard-gated Test Idea in `acceptance[]` or `Deferred test idea: <id> - <reason>`; design-contract + approved-scenarios complete/frozen or skip lines per `docs/mission-artifacts.md`.

## When to use

`/sc-run` planning / Task(`sc-planner`); user asks to plan or decompose. Over 7 tasks → same-mission `plan-phaseN.json` when `Sizing: phases`, else hand `/sc-discuss` + mission-sizing (planning never creates maps).

## Workflow

1. **Resolve** - `spacecraft resolve` / `use`.
2. **Read** - `spec.md`, `questions.md`, `decisions.md` (Test Ideas / Strategy Top risks / Charter / RCRCRC / Test data when present), optional `outputs/map.json`. Stop if blocking clarify open, or Testability `Not Testable` with soft Verify.
3. **Decompose** - ≤7 jigsaw tasks per phase as a **hard Must** (not preference-only; reject soft prefer ≤7; reject any 8-9 exception band). Prefer one independently testable behavioral slice per task (1-3 acceptance, one verify + evidence, real `dependsOn` only). Split when acceptance would exceed 3, verify surface differs, happy vs material error need separate cycles, or look vs behavior conflict → `/sc-discuss`. Over 7 with one ship story and `Sizing: phases` → same-mission `plan-phaseN.json`. Independent seams → hand multi-mission to `/sc-discuss` + mission-sizing. **Must not** invent `*-ux` seams, cross-feature waterfalls, or `*-integrate` tips mid-plan; **Must not** call `spacecraft map new` / `map add`. Mid-run multi-mission need → `/sc-discuss` Resize protocol.
4. **Write `plan.json`** - `planName`, `missionId`, `tasks[]` with `id`, `title`, `status`, `files`, `acceptance`, `verify`, `evidence` (+ optional `dependsOn`). When claiming UI/workflow behavior, include a product-surface marker among `verify.product` | `browser` | `curl` | `composition`.
5. **Write `design-contract.md`** - Scope, Modules, Data structures, Public seams, Edge matrix with expected literals, Out of scope; footer `Design-contract: complete`. Docs/prose-only: append `Design-contract skipped: docs/prose-only` (`docs/mission-artifacts.md`).
6. **Write `approved-scenarios.md`** - freeze Edge matrix + worked examples; footer `Approved-scenarios: frozen-from-contract` (or `frozen-by-human`). Docs/prose-only: append `Approved-scenarios skipped: docs/prose-only`. Do not invent expected values.
7. **Verify** - no vague tasks; paths real; ≤7; contract + scenarios complete/frozen or skipped.

### Hard-gated coverage

When Test Ideas / Strategy Top risks / Charter exist: cover each hard-gated idea (all Neg + Overlooked; plus Top risk/Charter when mapped or listed) in `acceptance[]` **or** `Deferred test idea: <id> - <reason>`. Positive / Edge prefer-only unless also Top risk/Charter.

### Map integration (optional)

When `outputs/map.json` exists, prefer touchpoints for `files`, flag high-consumer shared deps, and catch cross-layer side effects. Missing map is fine.

## Rules

- **Must**: Resolve; read inputs; stop on blocking clarify or soft Verify + `Not Testable`.
- **Must**: ≤7 tasks per phase as a hard Must (not preference-only); each task fully fielded; every acceptance verifiable as one cycle.
- **Must**: Hard-gated ideas in acceptance or Deferred line; product-surface marker when UI/workflow claimed.
- **Must**: After plan, write design-contract + approved-scenarios (or docs/prose skip lines per `docs/mission-artifacts.md`) before product build.
- **Must not**: Soft prefer ≤7; any 8-9 exception band; soft-prefer-only for hard-gated ideas; call `map new` / invent maps/seams; invent scenario literals; unbounded architecture essays in the contract.

## Out of scope

Draft UI - `/sc-discuss` + sc-ux-design. Implementation - sc-coder / sc-firmware. Evidence - sc-verification. Clarify - `/sc-discuss`.

## References

- `docs/mission-artifacts.md` - plan + design-contract + approved-scenarios schemas; outcome-gate skip/waive SoT
- sc-discuss `references/mission-sizing.md` - single / phases / roadmap (discuss owns map)
