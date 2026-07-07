---
description: Run the resolved Spacecraft mission workflow until a real gate blocks it
agent: sc-commander
---
Use sc-mission, sc-git, and sc-verification.
Run:
scripts/spacecraft resolve --json
Run:
scripts/spacecraft flow
If resolver safety is not `safe` or `flow` reports blockers, stop before writing. Show the blockers and exact next action.
Treat `.space/current` as fallback state, not sole authority.

Purpose: reduce unnecessary HIL inside one chat. Continue through the safe loop without asking the user to type each command:
`/sc-work Txx -> /sc-verify Txx -> checkpoint commit -> next task`.

Loop rules:
- Use `$ARGUMENTS` as an optional task selector or mode hint. If no task is given, use the first non-completed task in plan.json.
- Follow the `/sc-work` contract for exactly one task at a time: inspect spec.md, plan.json, questions.md, decisions.md, git state, and touched files before editing.
- Follow the `/sc-verify` contract immediately after each task: capture focused evidence with `scripts/spacecraft evidence "<label>" -- <command>`, then run `scripts/spacecraft validate`.
- Capture failures too. If verification fails, stop, summarize the failure evidence, and recommend the repair step.
- After passing verification for a task, update plan.json with evidence ids and mark the task completed only when acceptance checks are satisfied.
- If git is a valid non-main work branch, inspect dirty/untracked files, ensure `.gitignore` is current, stage only task-related safe files and mission artifacts, then create a checkpoint Conventional Commit.
- Do not checkpoint on `main`, detached HEAD, unsafe dirty state, missing git, failed verification, unrelated user changes, or uncertain files.
- Checkpoint commits are local rollback points. Before `/sc-ship`, squash/fixup them into logical final commits.
- Then continue to the next task while the same mission remains resolved, the chat context is still useful, and no gate blocks.

Hard stop gates:
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

Do not push, rebase, merge, tag, delete branches, create worktrees, or run release closeout unless the user explicitly asks or invokes `/sc-ship`.
End with evidence ids, checkpoint commit hash when created, next command, and session advice.
