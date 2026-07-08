---
description: Show resolved Spacecraft mission status
agent: sc-commander
---
Use sc-mission and sc-verification.
Resolve the mission. Block if unsafe.

## Workflow

1. Run:
   ```
   scripts/spacecraft status
   ```
   Treat the helper's `Selected by`, `Mission safety`, `Conflicts`, and `Candidates` output as authoritative.
2. Summarize the resolved mission state, git/rollback status, resolver safety, blockers, next recommended command, and whether to continue this chat or start a new session.

Next command recommendations:
- If resolver safety is not safe, tell the user to choose with `scripts/spacecraft missions` then `scripts/spacecraft use <number|id|title>`, or use `SPACECRAFT_MISSION=<mission-id>` for one command.
- If the mission includes UI but art direction is not selected, recommend /sc-design.
- If implementation is next, recommend /sc-git before /sc-flow.
- If implementation is next and no suitable mission or non-main branch exists, create or switch to one when the user's mutating intent is clear; otherwise recommend /sc-git.

## Error handling

- If resolver safety is not safe, do not recommend write/verify/review/ship work.
- If the user asks to stop this chat, close/end the session, or continue in a new session, give session handoff: state, blockers, git status, and pickup command. Do not run release closeout unless they ask to ship/release/merge/finish mission/close branch. If work appears ready, recommend /sc-ship.
