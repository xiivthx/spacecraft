---
description: Create or update resolved mission flight plan
agent: sc-commander
---
Use sc-mission, sc-clarify, sc-planning, and sc-architect.
Resolve the mission. Block if unsafe.

## Pre-flight checks

Read the resolved mission's spec.md, questions.md, decisions.md, and plan.json if present. Use sc-clarify skill before finalizing a plan. If mission clarification status is open and there are blocking questions, stop and tell the user to answer the current question.

Do not finalize plan.json while blocking clarification remains open. Non-blocking assumptions must be recorded in decisions.md.

If the mission includes UI, use sc-design and read DESIGN.md. If UI art direction is not chosen, stop and recommend /sc-design before finalizing UI tasks.

**sc-map**: If `outputs/map.json` is missing and the project has >10 source files, invoke sc-map before delegating to sc-planner. The map identifies touchpoints, dependency chains, and risk zones so the plan has zero side-effect blind spots.

## Workflow

1. Invoke sc-planner as a read-only subagent to draft the plan. A user invocation of /sc-plan is explicit permission to use the read-only sc-planner subagent; do not ask for separate subagent permission.
2. Write or update the resolved mission plan.json yourself.
3. The plan must contain no more than 7 tasks per phase.
4. Each task must have id, title, status, files, acceptance, verify, and evidence.
5. UI tasks must include visual intent, target component/screen, accessibility checks, and verification method.
6. For new screens, recommend /sc-design before /sc-build.
7. **Review gate** — Invoke sc-reviewer as a read-only subagent to review the plan quality. The reviewer checks: task specificity (no vague titles), acceptance verifiability, dependency ordering, phase splitting rationale. If the reviewer flags issues, fix them. Do not set state to planned until the review passes.
8. Set state to planned.

## Hard stop gates

- Resolver conflict or ambiguity
- Missing or empty spec.md
- Open blocking clarification
- UI art direction not chosen when UI tasks are present
- >7 tasks without sub-plan strategy recorded

## Error handling

- Do not implement product code.
- If blocking clarification remains open, stop and defer planning.

## Edge cases

- **spec.md is empty or missing** — Stop. Tell user to run /sc-start first.
- **plan.json already exists** — Read it. Update in place rather than overwriting.
- **>7 tasks needed** — Split into Phase 1, Phase 2, ... Record rationale in decisions.md. Each phase gets its own plan file: `plan-phase1.json`, `plan-phase2.json`, etc.
- **map.json is stale** — Note the staleness in decisions.md. Proceed with available information.
- **Task depends on another task** — Document the dependency order in task descriptions.

## Research auto-trigger

When planning tasks that involve unfamiliar APIs, frameworks, or dependency versions, run `spacecraft research "<query>"` before committing to task acceptance criteria that depend on that knowledge.

End with the recommended next action, issues/assumptions recorded, and session advice. Recommend `/sc-build` when implementation can begin (sc-git hygiene checks auto-trigger within sc-build). Recommend a new session if implementation is the next large phase and the plan is fully captured in plan.json.
