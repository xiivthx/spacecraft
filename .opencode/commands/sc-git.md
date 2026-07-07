---
description: Prepare or review Spacecraft git branch, commit, release, merge, and tag plan
agent: sc-commander
---
Use sc-mission and sc-git.
Run:
node scripts/spacecraft.mjs resolve --json
If resolver safety is not `safe` or no mission is selected, stop before branch, commit, merge, or release prep. Show the conflict/candidates and tell the user to run `node scripts/spacecraft.mjs missions` then `node scripts/spacecraft.mjs use <number|id|title>`, or set `SPACECRAFT_MISSION=<mission-id>` for one command.
Treat `.space/current` as fallback state, not sole authority.
Read the resolved mission's mission.json, plan.json, decisions.md, evidence.jsonl, and git state.
Run:
node scripts/spacecraft.mjs git-info

If the user asks for a branch suggestion, run:
node scripts/spacecraft.mjs git-suggest $ARGUMENTS

Apply Spacecraft release branching:
- never write product changes directly on main
- one branch per feature, fix, issue, or tightly scoped change
- branch from latest main
- use branch names like `<type>/<issue-or-mission>-<slug>`
- when clear mutating work has no non-main branch, create/switch to the branch without another blocking question
- keep `.gitignore` current before staging, committing, or merging
- do not allow secrets, local env files, private data, caches, logs, dependency folders, build outputs, or machine-specific files into git/public artifacts
- agent may commit frequently only on a valid non-main work branch
- plan final commits before implementation
- final branch history should have 1 to 3 commits and should not exceed 5 unless justified
- squash/fixup checkpoint commits into logical commits before merge
- rebase work branch on latest main before merge
- test, verify, and validate after rebase and before merge into main
- merge into main only with `git merge --no-ff <branch>`
- bump version before merge unless explicitly deferred with rationale
- update changelog and short spec/release note before merge when behavior changed
- create version tag after the no-ff merge into main
- after successful merge to main, delete the merged local branch unless the user asks to keep it
- if the user asks to ship/release/merge/finish mission/close branch, run release closeout prep; block if any gate is incomplete and list exact missing actions
- if the user asks only to stop this chat or continue in a new session while work is unfinished, do session handoff instead of release closeout

Use Conventional Commits.
Do not push unless the user explicitly asks.
Use rtk for noisy git/status/diff/log output when available; never use it to bypass denied operations.
End with next action and session advice. If implementation gates are ready, recommend `/sc-flow`.
