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
- branch is ready to rebase on latest main before merge
- tests, verification, and validation have a rerun plan after rebase
- merge plan is no-ff
- version bump, changelog/spec note, and tag plan exist when shipping
- branch cleanup will happen after successful merge to main unless explicitly kept
If the diff is intended to ship, record release readiness in review.json:
- releaseReadiness.version
- releaseReadiness.changelog
- releaseReadiness.specNote
- releaseReadiness.tagPlan
- releaseReadiness.postRebaseVerification
Use structured objects with `status`, plus `rationale` for any deferred gate. Do not use string or boolean releaseReadiness values.
Invoke sc-reviewer as a read-only subagent.
A user invocation of /sc-review is explicit permission to use the read-only sc-reviewer subagent; do not ask for separate subagent permission.
The reviewer must not edit files.
If UI files changed, recommend /sc-design-review or invoke sc-designer read-only if appropriate.
A user invocation of /sc-review is also sufficient permission for a focused read-only sc-designer sidecar when UI changes need design-risk triage; do not ask for separate subagent permission.
Critical design findings block shipping the same way critical code findings do.
After the subagent returns findings, record the review in:
- review.md
- review.json
Set state to reviewing while working.
If review status is ready, set state to ready.
If review status is blocked, set state to blocked.
Do not implement fixes in the same command unless the user explicitly asks.
End with next action and session advice. Suggest /sc-ship only when ready.
