---
description: Fast lane — branch, commit, ship without full spec/plan/build/review flow
agent: sc-commander
---
Use sc-mission and sc-git.
Start a quick mission for: $ARGUMENTS

## When to use

Use `/sc-quick` for small, straightforward changes that don't need detailed spec, task planning, TDD cycles, or formal review. Examples:
- Prompt tweaks, config changes, documentation
- Small fixes with obvious scope
- Personal tooling adjustments
- Anything where full flow overhead exceeds the change itself

Do NOT use for: new features with unknown scope, multi-file refactors, UI changes, API integrations, or anything needing design review.

## Pre-flight Checks

Resolve the mission. Block if unsafe.

Run:
```
scripts/spacecraft git-info
```

### Git safety

- If the workspace is not a git worktree, stop before editing product files unless the user has explicitly accepted no-git implementation risk for this mission in decisions.md.
- If git exists, inspect dirty state before edits. Work with user changes; do not revert unrelated changes.
- Use sc-git for branch naming, branch creation, checkpoint commits, .gitignore hygiene, and staging safety checks.

## Workflow

### 1. Mission stub

Run:
```
scripts/spacecraft new "$ARGUMENTS"
```

Then immediately:
- Set clarification to clear: `scripts/spacecraft clarify-status clear`
- Record in decisions.md: "Fast lane — skipping spec.md, plan.json, evidence, and formal review."
- Do NOT write spec.md or plan.json stubs
- Set state to planned: `scripts/spacecraft set-state planned`

### 2. Branch

Create a non-main work branch from latest `main` using sc-git naming:
```
scripts/spacecraft git-suggest feat
```

Create and switch to the branch. Bind it:
```
scripts/spacecraft bind-branch
```

### 3. Commit freely

- No TDD cycle required — implement directly
- No acceptance checks required
- No evidence capture required
- Conventional Commits still expected
- Checkpoint commits allowed on work branch
- `.gitignore` must stay current; no secrets staged

### 4. Fast self-review

Before ship, commander performs a lightweight self-review (see PERSONA.md fast self-review section):
- Inspect `git diff` — check for secrets, debug code, unrelated edits, dead code
- Run nearest cheap test if available
- Verify the change does what was intended
- No subagent reviewer — commander reviews directly
- No review.md or review.json required

If self-review finds issues, fix them and recommit.

### 5. Ship

When ready to ship, do release closeout with streamlined gates:

**Required (same as normal flow):**
- Version bump (or defer with rationale in decisions.md)
- Changelog update (mandatory — never defer)
- Rebase work branch on latest `main`
- Verify after rebase (run tests if available)
- `git merge --no-ff <branch>` into `main`
- Annotated tag: `git tag -a v<version> -m "v<version>"`
- Delete merged local branch
- `scripts/spacecraft archive` to compact shipped artifacts

**Skipped (fast lane only):**
- No evidence.jsonl requirement
- No review.md or review.json requirement
- No `scripts/spacecraft closeout-check` (requires evidence + review gates)
- No sc-reviewer subagent

**Ship checklist:**
- [ ] Version bumped (or deferred with rationale)
- [ ] Changelog updated
- [ ] Rebased on latest main
- [ ] Tests pass after rebase (if applicable)
- [ ] `git merge --no-ff` completed
- [ ] Tag created
- [ ] Branch deleted
- [ ] Mission archived

Produce a summary: mission id, what changed, git branch/merge info, suggested commit message, known limitations, next step.

## Research auto-trigger

When quick-lane changes touch unfamiliar tooling, configuration, or dependency APIs, run `spacecraft research "<topic>"` before committing. Fast lane is not skip-research lane.

## Hard Stop Gates

- Resolver conflict or ambiguity
- Write attempt on `main`
- Dirty/untracked files that cannot be safely attributed to the current task
- Unsafe files or secrets before staging
- Self-review finding critical issues
- Release actions requiring `/sc-ship`
- Context is too heavy for safe continuation; give handoff instead

## Error handling

- Do not push unless explicitly asked
- Do not write on `main` — always create branch first
- If git worktree is dirty with unrelated changes, warn before staging
- If self-review finds issues, fix before ship
- Conventional Commits format: `<type>: <description>`
- After post-merge cleanup, you are on `main` — next mutation requires new branch

End with session advice. Recommend new session after shipped mission.
