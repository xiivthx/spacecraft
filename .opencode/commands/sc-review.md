---
description: Review current mission diff and evidence
agent: sc-commander
---
Use sc-mission, sc-verification, and sc-git.
Read current mission spec.md, plan.json, evidence.jsonl, review.json, and git diff when git is available.
Run:
node scripts/spacecraft.mjs git-info
Review sc-git readiness:
- no product edits were made directly on main
- final commits are planned to be 5 or fewer
- final commit messages follow Conventional Commits
- .gitignore is current and no unsafe files are staged/tracked accidentally
- branch will be rebased on latest main before merge
- tests, verification, and validation will be rerun after rebase before merge
- merge plan is no-ff
- version bump, changelog/spec note, and tag plan exist when shipping
Invoke sc-reviewer as a read-only subagent.
The reviewer must not edit files.
If UI files changed, recommend /sc-design-review or invoke sc-designer read-only if appropriate.
Critical design findings block shipping the same way critical code findings do.
After the subagent returns findings, record the review in:
- review.md
- review.json
Set state to reviewing while working.
If review status is ready, set state to ready.
If review status is blocked, set state to blocked.
Do not implement fixes in the same command unless the user explicitly asks.
End with the recommended next action and session advice. Suggest /sc-ship only when ready; suggest a new session for larger fix cycles.
