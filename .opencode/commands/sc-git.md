---
description: Prepare or review Spacecraft git branch, commit, release, merge, and tag plan
agent: sc-commander
---
Use sc-mission and sc-git.
Read .space/current, mission.json, plan.json, decisions.md, evidence.jsonl, and git state.
Run:
node scripts/spacecraft.mjs git-info

If the user asks for a branch suggestion, run:
node scripts/spacecraft.mjs git-suggest $ARGUMENTS

Apply Spacecraft release branching:
- never write product changes directly on main
- one branch per feature, fix, issue, or tightly scoped change
- branch from latest main
- use branch names like `<type>/<issue-or-mission>-<slug>`
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

Use Conventional Commits for final commit suggestions.
Do not create branches, commits, rebases, merges, tags, or pushes unless the user explicitly asks.
End with the recommended next action and session advice.
