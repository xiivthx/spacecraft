---
description: Show current Spacecraft mission status
agent: sc-commander
---
Run:
node scripts/spacecraft.mjs status
Then summarize the current mission state, git/rollback status, blockers, next recommended command, and whether to continue this chat or start a new session.
If the mission includes UI but art direction is not selected, recommend /sc-design.
If implementation is next, recommend /sc-git before /sc-work.
If implementation is next and no suitable mission or non-main branch exists, create or switch to one when the user's mutating intent is clear; otherwise recommend /sc-git.
If the user asks to stop this chat, close/end the session, or continue in a new session, give session handoff: state, blockers, git status, and pickup command. Do not run release closeout unless they ask to ship/release/merge/finish mission/close branch. If work appears ready, recommend /sc-ship.
