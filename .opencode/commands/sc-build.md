---
description: Implement mission tasks in a continuous loop until a gate blocks
agent: sc-commander
subtask: true
---

Use sc-mission, sc-clarify, sc-git, sc-tdd, sc-solid, and sc-verification.
Resolve the mission. Block if unsafe.

## Pre-flight checks

Read mission.json, spec.md, plan.json. If spec.md or plan.json is missing, stop — tell user to run /sc-start or /sc-plan.

Read questions.md and decisions.md when present. If blocking clarification remains open, stop. Do not implement product code until clarification status is clear or explicitly deferred. If the user explicitly chooses to defer a decision, record the deferral in decisions.md and keep implementation limited to unaffected tasks.

Run:
```
scripts/spacecraft git-info
```

### Git safety

- If the workspace is not a git worktree, stop before editing product files unless the user has explicitly accepted no-git implementation risk for this mission in decisions.md.
- If git exists, inspect dirty state before edits. Work with user changes; do not revert unrelated changes.
- For large, risky, or multi-session slices, prefer a separate branch or git worktree.
- Use sc-git for branch naming, branch creation, checkpoint commits, .gitignore hygiene, and staging safety checks.

### Dependency check

Before code or dependency changes, check official current docs/registry/releases for direct dependencies and framework APIs. Use latest stable direct versions unless a deep dependency, ecosystem pin, or explicit user instruction says otherwise. Record source/version/date when it affects implementation.

## Per-task loop

Start from `$ARGUMENTS` task if given, otherwise the first non-completed task in plan.json. For each task:

### 1. Atomic TDD cycle

For each acceptance check (one at a time):

a. **Red** — delegate to `sc-tester`: write exactly ONE failing test for the current acceptance check. sc-tester must verify the test fails before returning. If the test passes without implementation, it is not testing the right thing — reject and re-write.

b. **Green** — delegate to `sc-coder`: write the minimum production code to make the single failing test pass. sc-coder must not change unrelated code, add features beyond the test, or modify other tests.

c. **Verify** — delegate to `sc-tester`: after the test passes, sc-tester runs the full test suite for the affected package, captures evidence with `scripts/spacecraft evidence "<label>" -- <command>`, runs `scripts/spacecraft validate`, and confirms all task acceptance criteria.

Repeat a–c for the next acceptance check until all checks for the task are covered.

### 2. Self-review / refactor

Commander inspects the touched diff:
- Make at most two short passes.
- Check the change against the current task and acceptance checks.
- Remove accidental debug code, unrelated edits, dead code, and noisy formatting churn.
- Extract helpers, improve names, remove duplication, simplify logic.
- Check obvious edge cases, error/loading states, accessibility labels, and mobile layout when relevant.
- Run the nearest cheap focused self-test when practical.
- If the pass finds small, low-risk issues, fix them and re-run the test.
- If the refactor is broad (>20 lines touched beyond the test-targeted code), defer to a separate task.

Keep this lightweight. Do not invoke sc-reviewer, do not write review.md or review.json.

### 3. Update plan

Update plan.json with evidence ids and mark the task `"done"` only when all acceptance checks are satisfied.

### 4. Checkpoint commit

Create a checkpoint Conventional Commit on the non-main work branch: inspect dirty/untracked files, ensure `.gitignore` is current, stage only task-related safe files and mission artifacts.

### 5. Continue

Continue to the next task while the same mission remains resolved, the chat context is still useful, and no gate blocks.

### UI implementation

If implementing UI:
- Read DESIGN.md and use sc-design
- If UI mood/theme/art direction is not recorded in decisions.md or plan.json, stop and recommend /sc-design before implementation
- Implement only the next planned UI slice
- Do not invent a new visual language when DESIGN.md exists
- Prefer CSS custom properties and local component styles
- Do not add broad styling frameworks unless explicitly requested

## Hard stop gates

- Resolver conflict or ambiguity
- Open blocking clarification
- Missing spec.md or plan.json
- UI implementation without recorded design direction
- Dependency/API choice needing current official docs not yet checked
- Write attempt on `main`
- Dirty/untracked files that cannot be safely attributed to the current task
- Unsafe files or secrets before staging
- Failed verification or validation
- Critical review/design finding
- Release actions requiring `/sc-ship`
- Context is too heavy for safe continuation; give handoff instead

## Error handling

- No git worktree → stop unless no-git risk accepted in decisions.md.
- Blocking clarification → stop and use sc-clarify skill.
- Missing spec.md or plan.json → stop and tell user to run /sc-start or /sc-plan.
- If self-review stops with known remaining risk, state that risk and recommend the next concrete command.
- If a useful self-test is skipped, say why and recommend the next concrete verification command.
- Do not push, rebase, merge, tag, delete branches, create worktrees, or run release closeout unless the user explicitly asks or invokes `/sc-ship`.
- If a hard stop gate triggers, stop the loop and hand off with state summary.
- Checkpoint commits are local rollback points. Before `/sc-ship`, squash/fixup them into logical final commits.
- When all tasks are done, apply the **Kalama Sutta gate** before claiming `built`:
  1. "Are all tasks actually done?" — every task in plan.json marked done
  2. "Does evidence cover every acceptance check?" — cross-check evidence.jsonl against plan.json acceptance arrays
  3. "Are tests passing or am I trusting cached results?" — re-run test suite fresh
  4. "Did I skip any verification?" — no acceptance check without evidence
  5. "Would an adversary agree this is complete?" — no undone tasks, no missing evidence
- If all 5 pass: `scripts/spacecraft set-state built`

End with evidence ids, checkpoint commit hash when created, next command, and session advice.
