---
name: sc-planning
description: "Convert a mission spec into a small executable plan with verifiable tasks. Activate during /sc-run planning, Task(sc-planner), task planning, or spec decomposition."
---

# sc-planning

## Goal

Turn `spec.md` into a jigsaw `plan.json` with ≤7 verifiable tasks per phase so `/sc-run` can execute per-acceptance RED-GREEN cycles via sc-tester / sc-coder.

## Output

Writable `plan.json` (schema in `docs/mission-artifacts.md`). Each task needs acceptance + verify + evidence. Each acceptance item is one RED-GREEN cycle.

## Good / Bad

- Good: atomic jigsaw slices; concrete acceptance (1-3 per task); exact verify; no hidden assumptions
- Bad: one coarse "implement feature" task; vague titles; missing verify; filling gray areas silently

## Verify

Every acceptance is testable and maps to one cycle; ≤7 tasks per phase; file paths real; no open blocking clarify.

## When to use

Activate when:

- `/sc-run` needs a plan (or Task `sc-planner`)
- User asks to plan / break the spec into tasks
- Scope work before implementation

When scope exceeds 7 tasks, split into Phase 1, Phase 2, ... each with its own plan.json.

## Workflow

Use this exact sequence unless the user specifies otherwise:

1. **Resolve mission** - `spacecraft resolve`. On conflict or ambiguity, use `spacecraft use <selector>`.

2. **Read inputs** - Before producing `plan.json`, read:
   - `spec.md` - what needs to be built
   - `questions.md` - any open blocking questions
   - `decisions.md` - recorded choices and assumptions
   - `outputs/map.json` - project structure survey (if present, see Map integration below)
   - If a blocking clarification question is open, stop - route to `/sc-discuss` / sc-clarify.

3. **Decompose into jigsaw tasks** - Break the feature into atomic behavioral vertical slices (puzzle pieces that combine into the full capability). ≤7 tasks per phase. When scope demands more:
   - Split into Phase 1, Phase 2, ... each with its own `plan.json` (`plan-phase1.json`, `plan-phase2.json`, ...)
   - Phase 1 covers the highest-priority, blocking, or foundational work
   - Each phase is independently buildable and verifiable
   - Record the split rationale in `decisions.md`

   **Jigsaw rule:** prefer one plan task per independently testable slice (form fields, remember-me, submit, API, errors, theme), not one task for the whole feature. `/sc-run` runs RED then GREEN once per acceptance item.

   Each task:
   - `id` - use the mission's compact sortable ID scheme (`T1`, `T2`, ... or match existing task numbering in the plan)
   - `title` - imperative, names the slice (e.g., "Bind username and password inputs with validation" not "Implement login")
    - `status` - start all as `pending`
    - `dependsOn` - optional array of task IDs this task depends on. Build loop skips tasks whose deps aren't `done`.
    - `files` - exact file paths when known. List only files directly touched. Use map.json touchpoints if available.
    - `acceptance` - 1–3 concrete checks per task. **Each item = one RED-GREEN cycle.** Split the task if you need more than 3.
    - `verify` - exact command or description of verification step (e.g., `npm test`, `curl localhost:3000/healthz`)
    - `evidence` - `spacecraft evidence "<label>" -- <command>`

4. **Write plan.json** - Produce `.space/missions/<mission-id>/plan.json`:
   ```json
   {
     "planName": "<short-descriptive-name>",
     "missionId": "<mission-id>",
     "tasks": [ ... ]
   }
   ```
   Use `spacecraft missions` to confirm the mission-id if uncertain.

5. **Verify** - Before claiming done: no task is vague, every acceptance check is testable, every file path is real (check with `ls` or glob), ≤7 tasks per phase.

### Map integration

When `outputs/map.json` exists, use it to scope accurately:
- **Touchpoints** - Scope task `files` to files identified as direct touchpoints. Cross-reference with spec intent to avoid missing critical paths.
- **Shared dependencies** - Files with >3 consumers are red-zone. Flag these in task acceptance checks (require extra review or broader test coverage).
- **Layers** - A change in one layer (e.g., skill SKILL.md) may need corresponding updates in another (e.g., agent permission files, docs). Use layer tags from map.json to catch these side effects.

If `map.json` is missing, proceed without it - it's optional input, not a hard gate.

### Edge cases

- **>7 tasks needed** - Two escape hatches: (1) same-mission phase split via `plan-phaseN.json` (Phase 1 → `plan-phase1.json`, Phase 2 → `plan-phase2.json`, ...); each phase gets its own `plan-phase<N>.json`. (2) roadmap/multi-mission split via `spacecraft map`. Record the split rationale in `decisions.md`.
- **Blocking question open** - Stop and route to `/sc-discuss` / sc-clarify. Do not produce `plan.json` with hidden assumptions.
- **File paths uncertain** - Use map.json or inspect the repo. If still uncertain, note it in task `files` as `"<discover-during-implementation>"`.
- **Spec is incomplete** - Flag gaps in `decisions.md`. Plan only what's specified.
- **Task depends on another task** - Use `dependsOn: ["T01"]` field. Build loop mechanically skips tasks whose deps aren't `done`.

## Rules

- **Must**: Resolve mission before planning.
- **Must**: Read `spec.md`, `questions.md`, `decisions.md`, and `map.json` (if present) before writing `plan.json`.
- **Must**: Stop if a blocking clarification is open - route to `/sc-discuss` / sc-clarify.
- **Must**: ≤7 tasks per phase as a hard Must (not preference-only). Split into Phase 1, Phase 2, ... if scope exceeds this.
- **Must**: Each task has `id`, `title`, `status`, `files`, `acceptance`, `verify`, `evidence`.
- **Must**: Every acceptance check is verifiable (can a reviewer confirm it?).
- **Must**: File paths are real - verify with `ls` or glob before writing.
- **Must not**: Soft prefer ≤7; reject any 8-9 exception band.
- **Must not**: Use vague tasks like "improve code", "add features", or one task that swallows the whole feature.
- **Must not**: Fill gray areas with hidden assumptions. Record assumptions explicitly.
- **Must not**: Create broad architecture plans unless the spec requires it.
- **Must**: Treat each `acceptance[]` item as one RED-GREEN cycle for `/sc-run`.

## Out of scope

- Design or UI work - draft under `/sc-discuss` + sc-ux-design; critique via Task(`sc-designer`)
- Implementation - Task `sc-coder` / `sc-firmware` under `/sc-run`
- Verification - use sc-verification
- Clarification - use `/sc-discuss` (sc-clarify protocol)

## Output format

```json
{
  "planName": "add-health-endpoint",
  "missionId": "M07FYB5W5",
  "tasks": [
    {
      "id": "T1",
      "title": "Add GET /healthz endpoint returning { ok: true }",
      "status": "pending",
      "files": ["src/server.ts", "src/server.test.ts"],
      "acceptance": [
        "Endpoint responds 200 with { ok: true }",
        "Test verifies 200 response and body shape"
      ],
      "verify": "npm test -- --testPathPattern server.test.ts",
      "evidence": "spacecraft evidence \"health-endpoint\" -- npm test -- --testPathPattern server.test.ts"
    }
  ]
}
```

## Checklist

- [ ] Mission resolved
- [ ] `spec.md`, `questions.md`, `decisions.md`, `map.json` (if present) read
- [ ] No blocking clarification open
- [ ] Plan has ≤7 jigsaw tasks per phase (split if needed)
- [ ] Each task is a behavioral slice (not the whole feature)
- [ ] Each task has all required fields
- [ ] Every acceptance check is verifiable and sized as one RED-GREEN cycle
- [ ] File paths verified real (not guessed)
- [ ] Assumptions recorded in `decisions.md` if any

## References

- `spec.md` - mission specification (input)
- `decisions.md` - recorded choices and assumptions
- `outputs/map.json` - project structure survey (optional input)
- `spacecraft missions` - list missions and confirm IDs
