---
description: Create or update resolved mission flight plan
agent: sc-commander
---
Use sc-mission, sc-clarify, and sc-planning.
Resolve the mission. Block if unsafe.

## Pre-flight checks

Read the resolved mission's spec.md, questions.md, decisions.md, and plan.json if present. Use sc-clarify skill before finalizing a plan. If mission clarification status is open and there are blocking questions, stop and tell the user to answer the current question.

Do not finalize plan.json while blocking clarification remains open. Non-blocking assumptions must be recorded in decisions.md.

If the mission includes UI, use sc-design and read DESIGN.md. If UI art direction is not chosen, stop and recommend /sc-design before finalizing UI tasks.

## Workflow

1. Invoke sc-planner as a read-only subagent to draft the plan. A user invocation of /sc-plan is explicit permission to use the read-only sc-planner subagent; do not ask for separate subagent permission.
2. Write or update the resolved mission plan.json yourself.
3. The plan must contain no more than 7 tasks.
4. Each task must have id, title, status, files, acceptance, verify, and evidence.
5. UI tasks must include visual intent, target component/screen, accessibility checks, and verification method.
6. For new screens, recommend /sc-design before /sc-build.
7. Set state to planned.

## Error handling

- Do not implement product code.
- If blocking clarification remains open, stop and defer planning.

## Edge cases

- **spec.md is empty or missing** — Stop. Tell user to run /sc-start first.
- **plan.json already exists** — Read it. Update in place rather than overwriting.
- **>7 tasks needed** — Split into sub-plans. Record rationale in decisions.md.
- **map.json is stale** — Note the staleness in decisions.md. Proceed with available information.
- **Task depends on another task** — Document the dependency order in task descriptions.

## Research auto-trigger

When planning tasks that involve unfamiliar APIs, frameworks, or dependency versions, run `spacecraft research "<query>"` before committing to task acceptance criteria that depend on that knowledge.

End with the recommended next action and session advice. Recommend `/sc-git` then `/sc-build` when implementation can begin. Recommend a new session if implementation is the next large phase and the plan is fully captured in plan.json.
