---
name: sc-start
description: "Starts a mission when the user begins scoped feature work or selects an existing mission."
disable-model-invocation: true
---

Use sc-mission and sc-clarify.
Start a new mission for: $ARGUMENTS

## Pre-flight checks

If the user wants an existing mission instead of a new one, run `spacecraft missions` and select with `spacecraft use <number|id|title>`.

## Workflow

1. Run:
   ```
   spacecraft new "$ARGUMENTS"
   ```
   The helper records git base sha when the workspace is a git worktree, writes `.space/current` as fallback, and binds the mission to the local session when a stable session key exists.

2. If the request clearly includes mutating work, also create a non-main branch from main using sc-git naming. Do not ask another question for this.
3. Draft only a minimal initial spec.md from the user request.
4. **Review gate** - Invoke /sc-reviewer as a read-only subagent to review the spec. The reviewer checks: scope is clear, acceptance criteria defined, no hidden assumptions, spec matches user intent. If the reviewer flags gaps, fix them before proceeding.
5. Create or update questions.md and decisions.md.
6. Inspect available repo context if useful.
7. Identify gray areas before planning or implementation.
8. If there is a blocking ambiguity, ask exactly one question and stop. Include your recommended answer.
9. If no blocking ambiguity exists, record assumptions in decisions.md and set clarification status to clear.
10. Keep state `active` (default) when the mission has enough clarity for the initial spec. Use `spacecraft set-state [mission-id] <new-state>` when a later transition is needed.

## Research auto-trigger

When mission scope involves unfamiliar tools, frameworks, or APIs, use sc-search (WebSearch/WebFetch) for `"<topic>"` before drafting the initial spec. A well-informed spec eliminates downstream replanning.

## Hard stop gates

- Resolver conflict or ambiguity
- Empty or vague mission title
- Duplicate mission detected without user acknowledge
- No-git workspace without user accepting risk

## Error handling

- Do not implement product code.
- Do not create a detailed plan.
- Do not run /sc-design implicitly.
- Do not assume product or design direction silently.

## Edge cases

- **Mission already exists for this request** - Use `spacecraft missions` to list. Offer to select existing rather than create duplicate.
- **Not in a git worktree** - The mission is created but warn: "No-git implementation risk." Record the warning in decisions.md if user accepts.
- **$ARGUMENTS is empty or vague** - Ask for a one-line mission description. Do not create a mission from an empty title.
- **Multiple sessions active** - The new mission binds to the current session. Warn if another session has an active mission.

End with next action, assumptions recorded in decisions.md, and session advice. Recommend `/sc-design` if the mission has UI, otherwise `/sc-plan`.
