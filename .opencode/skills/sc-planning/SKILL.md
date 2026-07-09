---
name: sc-planning
description: Convert a mission spec into a small executable plan with verifiable tasks. Activate on /sc-plan, task planning, or spec decomposition.
license: MIT
compatibility: opencode
metadata:
  version: 1.0
  audience: xiivthx
---

# sc-planning

Convert a mission spec into a small executable plan. Output is `plan.json` with ≤7 verifiable tasks, each with acceptance checks and evidence requirements.

## When to use

Activate when the user asks to:

- **Plan next steps / "create a plan" / "/sc-plan"** — explicit planning
- **Break the spec into tasks** — task decomposition from `spec.md`
- **Scope work before implementation** — pre-`/sc-build` task definition

## Workflow

Use this exact sequence unless the user specifies otherwise:

1. **Resolve mission** — `scripts/spacecraft resolve --json`. Block if safety ≠ `safe`.

2. **Read inputs** — Before producing `plan.json`, read:
   - `spec.md` — what needs to be built
   - `questions.md` — any open blocking questions
   - `decisions.md` — recorded choices and assumptions
   - `outputs/map.json` — project structure survey (if present, see Map integration below)
   - If a blocking clarification question is open, stop — route to sc-clarify.

3. **Decompose into tasks** — Break the spec into ≤7 small, verifiable tasks. Each task:
   - `id` — use the mission's compact sortable ID scheme (`T1`, `T2`, ... or match existing task numbering in the plan)
   - `title` — imperative, specific (e.g., "Add health check endpoint" not "Implement health")
   - `status` — start all as `pending`
   - `files` — exact file paths when known. List only files directly touched. Use map.json touchpoints if available.
   - `acceptance` — 1–3 concrete checks per task. Verifiable statements, not abstract goals.
   - `verify` — exact command or description of verification step (e.g., `npm test`, `curl localhost:3000/healthz`)
   - `evidence` — `scripts/spacecraft evidence "<label>" -- <command>`

4. **Write plan.json** — Produce `.space/missions/<mission-id>/plan.json`:
   ```json
   {
     "planName": "<short-descriptive-name>",
     "missionId": "<mission-id>",
     "tasks": [ ... ]
   }
   ```
   Use `scripts/spacecraft missions` to confirm the mission-id if uncertain.

5. **Verify** — Before claiming done: no task is vague, every acceptance check is testable, every file path is real (check with `ls` or glob), ≤7 tasks.

### Map integration

When `outputs/map.json` exists, use it to scope accurately:
- **Touchpoints** — Scope task `files` to files identified as direct touchpoints. Cross-reference with spec intent to avoid missing critical paths.
- **Shared dependencies** — Files with >3 consumers are red-zone. Flag these in task acceptance checks (require extra review or broader test coverage).
- **Layers** — A change in one layer (e.g., skill SKILL.md) may need corresponding updates in another (e.g., agent permission files, docs). Use layer tags from map.json to catch these side effects.

If `map.json` is missing, proceed without it — it's optional input, not a hard gate.

### Edge cases

- **>7 tasks needed** — Split into sub-plans or defer lower-priority tasks. Record the split in `decisions.md`.
- **Blocking question open** — Stop and route to sc-clarify. Do not produce `plan.json` with hidden assumptions.
- **File paths uncertain** — Use map.json or inspect the repo. If still uncertain, note it in task `files` as `"<discover-during-implementation>"`.
- **Spec is incomplete** — Flag gaps in `decisions.md`. Plan only what's specified.
- **Task depends on another task** — Document in task description. Process them in dependency order during `/sc-build`.

## Rules

- **Must**: Resolve mission before planning.
- **Must**: Read `spec.md`, `questions.md`, `decisions.md`, and `map.json` (if present) before writing `plan.json`.
- **Must**: Stop if a blocking clarification is open — route to sc-clarify.
- **Must**: ≤7 tasks per plan.
- **Must**: Each task has `id`, `title`, `status`, `files`, `acceptance`, `verify`, `evidence`.
- **Must**: Every acceptance check is verifiable (can a reviewer confirm it?).
- **Must**: File paths are real — verify with `ls` or glob before writing.
- **Must not**: Use vague tasks like "improve code" or "add features". Be specific.
- **Must not**: Fill gray areas with hidden assumptions. Record assumptions explicitly.
- **Must not**: Create broad architecture plans unless the spec requires it.

## Out of scope

- Design or UI work — use sc-design
- Implementation — use /sc-build
- Verification — use sc-verification
- Clarification — use sc-clarify

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
      "evidence": "scripts/spacecraft evidence \"health-endpoint\" -- npm test -- --testPathPattern server.test.ts"
    }
  ]
}
```

## Checklist

- [ ] Mission resolved
- [ ] `spec.md`, `questions.md`, `decisions.md`, `map.json` (if present) read
- [ ] No blocking clarification open
- [ ] Plan has ≤7 tasks
- [ ] Each task has all 7 required fields
- [ ] Every acceptance check is verifiable
- [ ] File paths verified real (not guessed)
- [ ] Assumptions recorded in `decisions.md` if any

## Research auto-trigger

When planning tasks that involve unfamiliar APIs, dependencies, or frameworks, run `spacecraft research "<query>"` to verify approach and version compatibility before committing to the plan.

---

## References

- `spec.md` — mission specification (input)
- `decisions.md` — recorded choices and assumptions
- `outputs/map.json` — project structure survey (optional input)
- `scripts/spacecraft missions` — list missions and confirm IDs
