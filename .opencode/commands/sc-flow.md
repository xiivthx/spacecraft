---
description: Run the resolved Spacecraft mission workflow until a real gate blocks it
agent: sc-commander
subtask: true
---
Use sc-mission, sc-git, and sc-verification.
Resolve the mission. Block if unsafe.

## Pre-flight checks

Run:
```
scripts/spacecraft workflow
```

Verify mission state is planned, branch is non-main, and git is clean.

## Workflow

Continue through the safe loop without asking the user to type each command.
For each task, use **atomic TDD cycles** (Red → Green → Refactor) with write-capable subagents.
Break acceptance checks into the smallest testable units and iterate one at a time.

### Per-task loop

1. Use `$ARGUMENTS` as an optional task selector or mode hint. If no task is given, use the first non-completed task in plan.json.

2. **Atomic TDD cycle** — for each acceptance check (one at a time):

   a. **Red — delegate to `sc-tester`**: write exactly ONE failing test for the current acceptance check. sc-tester must verify the test fails before returning. If the test passes without implementation, it is not testing the right thing — reject and re-write.

   b. **Green — delegate to `sc-coder`**: write the minimum production code to make the single failing test pass. sc-coder must not change unrelated code, add features beyond the test, or modify other tests.

   c. **Refactor — commander reviews the diff**: after Green passes, commander inspects the touched diff for cleanup opportunities without changing behavior — extract helpers, improve names, remove duplication, simplify logic. If cleanup is needed, apply it and re-run the test. If the refactor is broad (>20 lines touched beyond the test-targeted code), defer to a separate task.

   d. Repeat a–c for the next acceptance check until all checks for the task are covered.

3. **Final verification — delegate to `sc-tester`**: after all acceptance checks pass, sc-tester runs the full test suite for the affected package, captures evidence with `scripts/spacecraft evidence "<label>" -- <command>`, runs `scripts/spacecraft validate`, and confirms all task acceptance criteria.

4. Update plan.json with evidence ids and mark the task `"done"` only when all acceptance checks are satisfied.

5. Create a checkpoint Conventional Commit on the non-main work branch: inspect dirty/untracked files, ensure `.gitignore` is current, stage only task-related safe files and mission artifacts.

6. Continue to the next task while the same mission remains resolved, the chat context is still useful, and no gate blocks.

## Hard stop gates

- resolver conflict or ambiguity
- open blocking clarification
- missing spec.md or plan.json
- UI implementation without recorded design direction
- dependency/API choice needing current official docs not yet checked
- write attempt on `main`
- dirty/untracked files that cannot be safely attributed to the current task
- unsafe files or secrets before staging
- failed verification or validation
- critical review/design finding
- release actions requiring `/sc-ship`
- context is too heavy for safe continuation; give handoff instead

## Error handling

- Do not push, rebase, merge, tag, delete branches, create worktrees, or run release closeout unless the user explicitly asks or invokes `/sc-ship`.
- If a hard stop gate triggers, stop the loop and hand off with state summary.
- Checkpoint commits are local rollback points. Before `/sc-ship`, squash/fixup them into logical final commits.

End with evidence ids, checkpoint commit hash when created, next command, and session advice.
