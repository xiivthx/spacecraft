---
description: Prepare or review git branch, commit, release, merge, and tag plan
agent: sc-commander
---
Use sc-mission and sc-git.
Resolve the mission. Block if unsafe.

## Pre-flight checks

Read the resolved mission's mission.json, plan.json, decisions.md, evidence.jsonl, and git state.

Run:
```
scripts/spacecraft git-info
```

## Workflow

### 1. Branch management

Apply release branching:
- never write product changes directly on main
- one branch per feature, fix, issue, or tightly scoped change
- branch from latest main
- use branch names like `<type>/<id>/<title>`
- when clear mutating work has no non-main branch, create/switch to the branch without another blocking question
- if the user asks for a branch suggestion, run `scripts/spacecraft git-suggest $ARGUMENTS`

### 2. Staging and hygiene

- keep `.gitignore` current before staging, committing, or merging
- do not allow secrets, local env files, private data, caches, logs, dependency folders, build outputs, or machine-specific files into git/public artifacts

### 3. Commits

- agent may commit frequently only on a valid non-main work branch
- plan final commits before implementation
- final branch history should have 1 to 3 commits and should not exceed 5 unless justified
- squash/fixup checkpoint commits into logical commits before merge
- use Conventional Commits: `<type>: <description>` subject, no scope by default. Body uses `- ` bullet points with lowercase first character.

### 4. Merge and release

- rebase work branch on latest main before merge
- test, verify, and validate after rebase and before merge into main
- merge into main only with `git merge --no-ff <branch>`
- bump version before merge unless explicitly deferred with strong rationale
- update changelog before merge — always required, never defer
- update short spec/release note before merge when behavior changed
- create version tag after the no-ff merge into main
- after successful merge to main, delete the merged local branch unless the user asks to keep it

### 5. Session handling

- if the user asks to ship/release/merge/finish mission/close branch, run release closeout prep; block if any gate is incomplete and list exact missing actions
- if the user asks only to stop this chat or continue in a new session while work is unfinished, do session handoff instead of release closeout

## Error handling

- Do not push unless the user explicitly asks.
- Do not create worktrees, force-push, or run `git init` unless asked.
- Detached HEAD state: refuse mutating work until a branch is created.
- Dirty state blocking branch creation: warn user about unrelated changes before creating branch.
- Branch name collision: if the suggested branch name exists, append a short suffix.

## Hard stop gates

- Write attempt on `main`
- Detached HEAD during mutating work
- Uncommitted changes that conflict with the operation
- Secrets or local env files detected before staging
- `.gitignore` stale — build outputs, caches, or dependency folders unstaged
- Verification failure after rebase

End with next action and session advice. If implementation gates are ready, recommend `/sc-build`.
