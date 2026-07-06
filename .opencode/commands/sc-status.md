---
description: Show current Spacecraft mission status
agent: sc-commander
---
Run:
node scripts/spacecraft.mjs status
Then summarize the current mission state, git/rollback status, blockers, next recommended command, and whether to continue this chat or start a new session.
If the mission includes UI but art direction is not selected, recommend /sc-design.
If implementation is next, recommend /sc-git before /sc-work.
If implementation is next and the workspace is not a git worktree, recommend initializing git or explicitly accepting no-git implementation risk before /sc-work.
