---
name: sc-build
description: "Implements planned mission tasks when the specification and plan gates are complete."
disable-model-invocation: true
---

Use sc-mission, sc-clarify, sc-git, sc-tdd, sc-solid, sc-learn, sc-verification, sc-web-frontend, sc-web-backend, and sc-database.
Resolve the mission. Block if unsafe.

## Pre-flight checks

Read mission.json, spec.md, plan.json. If spec.md or plan.json is missing, stop - tell user to run /sc-start or /sc-plan.

Read questions.md and decisions.md when present. If blocking clarification remains open, stop. Do not implement product code until clarification status is clear or explicitly deferred. If the user explicitly chooses to defer a decision, record the deferral in decisions.md and keep implementation limited to unaffected tasks.

Run:
```
spacecraft git-info
```

### Git safety

sc-git hygiene checks auto-trigger silently during pre-flight: `git-info`, branch check, dirty-tree check, `.gitignore` check. These block only on violation - no separate git command needed. The sc-git skill handles all policy; the build command calls it silently.

- If the workspace is not a git worktree, stop before editing product files unless the user has explicitly accepted no-git implementation risk for this mission in decisions.md.
- If git exists, inspect dirty state before edits. Work with user changes; do not revert unrelated changes.
- For large, risky, or multi-session slices, prefer a separate branch or git worktree.

### Dependency check

Before code or dependency changes, check official current docs/registry/releases for direct dependencies and framework APIs. Use latest stable direct versions unless a deep dependency, ecosystem pin, or explicit user instruction says otherwise. Record source/version/date when it affects implementation.

When versions or APIs are uncertain, run `spacecraft research "<package> latest version"` before installing. Do not guess dependency versions.

## Per-task loop

**Orchestrator rule**: You are an orchestrator, not an implementer. You MUST use the Cursor subagent delegation to invoke `/sc-coder` and `/sc-tester` for all source code writing and testing. You must NOT use Write, Edit, or Bash to create or modify source files (`.go`, `.ts`, `.tsx`, `.js`, `.css`, etc.) or test files. Your write operations are limited to: updating plan.json, staging and committing artifacts, and capturing evidence. Violating this rule bypasses the zero-trust review chain.

Start from `$ARGUMENTS` task if given, otherwise the first non-completed task in plan.json. For each task:

### 1. TDD cycle (Plan → Red → Green → Verify)

For each acceptance check (one at a time):

a. **Triage** - Is this test-worthy? If the test would be a trivial tautology (simple getter, config mapping, boilerplate), skip TDD. Delegate the implementation to `/sc-coder` with clear instructions. Record the skip. Otherwise proceed with the cycle.

b. **Plan** - What seam are we testing? What is the expected behavior and independent expected value? Confirm the plan before writing code. No test at an unplanned seam.

c. **Red** - MUST delegate to `/sc-tester` via Cursor subagent delegation: write exactly ONE failing test for the planned acceptance check. /sc-tester must verify the test fails before returning. If the test passes without implementation, it is not testing the right thing - reject and re-write.

d. **Green** - MUST delegate to `/sc-coder` via Cursor subagent delegation: write the minimum production code to make the single failing test pass. /sc-coder must not change unrelated code, add features beyond the test, or modify other tests.

e. **Verify** - MUST delegate to `/sc-tester` via Cursor subagent delegation: after the test passes, /sc-tester runs the full test suite for the affected package, captures evidence with `spacecraft evidence "<label>" -- <command>`, runs `spacecraft validate`, and confirms all task acceptance criteria.

Repeat a–e for the next acceptance check until all checks for the task are covered.

### 2. Refactor (post-feature)

After ALL acceptance checks for the task pass, refactor with the full picture:

- Extract helpers, improve names, remove duplication (Rule of Three), simplify logic.
- Remove accidental debug code, dead code, noisy formatting churn.
- The existing tests protect you - refactor with confidence.
- If the refactor is broad (>20 lines touched beyond test-targeted code), defer to a separate task.
- Make at most two short passes.

### 3. Functional test gate

After refactor: run the FULL test suite (unit + integration + functional). All old tests must pass alongside new tests. If anything breaks, fix the refactor - not the old tests. Capture evidence:

```
spacecraft evidence "<task>-functional" -- <full-test-suite-command>
```

### 4. Update plan

Update plan.json with evidence ids and mark the task `"done"` only when all acceptance checks are satisfied.

### 5. Checkpoint commit

Create a checkpoint Conventional Commit on the non-main work branch: inspect dirty/untracked files, ensure `.gitignore` is current, stage only task-related safe files and mission artifacts.

### 6. Continue

Continue to the next task while the same mission remains resolved, the chat context is still useful, and no gate blocks.

### UI implementation

If implementing UI:
- Read DESIGN.md and use sc-design
- If UI mood/theme/art direction is not recorded in decisions.md or plan.json, stop and recommend /sc-design before implementation
- Implement only the next planned UI slice
- Do not invent a new visual language when DESIGN.md exists
- Prefer CSS custom properties and local component styles
- Do not add broad styling frameworks unless explicitly requested

## Research auto-trigger

When versions or APIs are uncertain, run `spacecraft research "<package> latest version"` before installing. Do not guess dependency versions. This applies to direct dependencies, framework APIs, and breaking-change migrations.

## Hard stop gates

- Resolver conflict or ambiguity
- Open blocking clarification
- Missing spec.md or plan.json
- UI implementation without recorded design direction
- Dependency/API choice needing current official docs not yet checked
- Write attempt on `main`
- Commander writing source code directly (violation of orchestrator rule)
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
- If a hard stop gate triggers, stop the loop, record the issue via sc-learn in `.space/missions/<id>/issues.md`, and hand off with state summary.
- Checkpoint commits are local rollback points. Before `/sc-ship`, squash/fixup them into logical final commits.
- When all tasks are done, apply the **review gate**:
  1. Invoke /sc-reviewer as a read-only subagent to review the full diff, evidence, and code quality. The reviewer checks: all tasks done, evidence covers every acceptance check, tests pass fresh, no critical findings, cross-reference integrity.
  2. Record findings in `review.md` and `review.json`.
  3. Fix any issues the reviewer flags. Re-run verification after fixes.
  4. Apply the **Kalama Sutta gate**:
     - "Are all tasks actually done?" - every task in plan.json marked done
     - "Does evidence cover every acceptance check?" - cross-check evidence.jsonl against plan.json
     - "Are tests passing or am I trusting cached results?" - re-run test suite fresh
     - "Did I skip any verification?" - no acceptance check without evidence
     - "Would an adversary agree this is complete?" - no undone tasks, no missing evidence
- If all gates pass: `spacecraft set-state ready`. If blocked: `spacecraft set-state blocked`.

End with evidence ids, checkpoint commit hash when created, issues/lessons discovered (record in `.space/missions/<id>/issues.md` and `learned.md` via sc-learn), next command, and session advice.
