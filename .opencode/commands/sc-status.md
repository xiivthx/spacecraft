---
description: Show resolved Spacecraft mission status
agent: sc-commander
---
Run:
node scripts/spacecraft.mjs status
Treat the helper's `Selected by`, `Mission safety`, `Conflicts`, and `Candidates` output as authoritative.
Then summarize the resolved mission state, git/rollback status, resolver safety, blockers, next recommended command, and whether to continue this chat or start a new session.
If resolver safety is not safe, do not recommend write/verify/review/ship work. Tell the user to choose with `node scripts/spacecraft.mjs missions` then `node scripts/spacecraft.mjs use <number|id|title>`, or use `SPACECRAFT_MISSION=<mission-id>` for one command.
If the mission includes UI but art direction is not selected, recommend /sc-design.
If implementation is next, recommend /sc-git before /sc-flow.
If implementation is next and no suitable mission or non-main branch exists, create or switch to one when the user's mutating intent is clear; otherwise recommend /sc-git.
If the user asks to stop this chat, close/end the session, or continue in a new session, give session handoff: state, blockers, git status, and pickup command. Do not run release closeout unless they ask to ship/release/merge/finish mission/close branch. If work appears ready, recommend /sc-ship.
