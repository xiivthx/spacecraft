---
description: Prepare final Spacecraft delivery summary
agent: sc-commander
---
Use sc-mission, sc-verification, and sc-git.
Run:
node scripts/spacecraft.mjs resolve --json
If resolver safety is not `safe` or no mission is selected, stop before release closeout. Show the conflict/candidates and tell the user to run `node scripts/spacecraft.mjs missions` then `node scripts/spacecraft.mjs use <number|id|title>`, or set `SPACECRAFT_MISSION=<mission-id>` for one command.
Treat `.space/current` as fallback state, not sole authority.
Read the resolved mission's spec.md, plan.json, evidence.jsonl, review.md, review.json, questions.md, decisions.md, and git diff when git is available.
Run:
node scripts/spacecraft.mjs validate
Run:
node scripts/spacecraft.mjs git-info
Run:
node scripts/spacecraft.mjs closeout-check
Treat "ship", "release", "merge", "finish mission", and "close branch" as release closeout requests.
Do not treat ordinary session handoff or "continue in a new session" as release closeout.
Only close/ship/merge if:
- blocking clarification questions are resolved
- acceptance checks have evidence
- important verification commands have passing evidence
- review.json status is ready
- there are no critical findings
- mutating product work has a rollback story: git base sha, branch/worktree, or an explicit no-git risk acceptance recorded in decisions.md
- product work was not committed or edited directly on main
- work branch has been rebased on latest main or has a documented local-only exception
- final branch history has 5 or fewer commits unless explicitly justified
- final commit messages follow Conventional Commits
- .gitignore is current and no unsafe files are staged/tracked accidentally
- tests, verification, and validation pass after the latest rebase and before merge into main
- merge plan uses `git merge --no-ff <branch>`
- version bump is complete or explicitly deferred with rationale
- changelog and short spec/release note are updated when product behavior changed
- tag plan exists for the bumped version
- if UI files changed, review.md or review.json includes a design review result
- UI work has no unresolved critical design findings
- if UI files changed, art direction decisions are recorded in decisions.md or explicitly deferred
If any gate fails, block closeout. List exact missing actions and next command.
If all gates pass, prepare merge to main with sc-git:
- final commits are logical Conventional Commits
- branch is rebased on latest main
- verification is rerun after rebase
- merge uses `git merge --no-ff <branch>`
- version tag is created after merge when version was bumped
- merged local branch is deleted unless user asks to keep it
- no push unless explicitly requested
Produce concise final summary:
- Mission id
- What changed
- Evidence ids
- Review status
- Git branch or rollback status
- Suggested Conventional Commit message if the user intends to commit
- Design evidence or manual visual verification notes when applicable
- Important confirmed decisions and assumptions when relevant
- Known limitations
- Suggested next step
Then set state to shipped if appropriate.
Do not git push unless the user explicitly asks.
Suggested commit messages must follow Conventional Commits:
<type>[optional scope]: <description>
End with session advice. Usually recommend a new session after a shipped mission or major phase boundary, with `/sc-status` as pickup.
