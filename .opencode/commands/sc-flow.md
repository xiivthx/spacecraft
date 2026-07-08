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

Continue through the safe loop without asking the user to type each command:
`/sc-work Txx -> /sc-verify Txx -> checkpoint commit -> next task`.

1. Use `$ARGUMENTS` as an optional task selector or mode hint. If no task is given, use the first non-completed task in plan.json.
2. Follow the `/sc-work` contract for exactly one task at a time: inspect spec.md, plan.json, questions.md, decisions.md, git state, and touched files before editing.
3. Follow the `/sc-verify` contract immediately after each task: capture focused evidence with `scripts/spacecraft evidence "<label>" -- <command>`, then run `scripts/spacecraft validate`.
4. Capture failures too. If verification fails, stop, summarize the failure evidence, and recommend the repair step.
5. After passing verification for a task, update plan.json with evidence ids and mark the task completed only when acceptance checks are satisfied.
6. Use sc-git for checkpoint commits on valid non-main work branch only. Inspect dirty/untracked files, ensure `.gitignore` is current, stage only task-related safe files and mission artifacts, then create a checkpoint Conventional Commit.
7. Then continue to the next task while the same mission remains resolved, the chat context is still useful, and no gate blocks.

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
