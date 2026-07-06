---
description: Prepare final Spacecraft delivery summary
agent: sc-commander
---
Use sc-mission, sc-verification, and sc-git.
Read spec.md, plan.json, evidence.jsonl, review.md, review.json, questions.md, decisions.md, and git diff when git is available.
Run:
node scripts/spacecraft.mjs validate
Run:
node scripts/spacecraft.mjs git-info
Only ship if:
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
Produce a concise final summary:
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
Do not git push.
Do not create commits unless the user explicitly asks.
Suggested commit messages must follow Conventional Commits:
<type>[optional scope]: <description>
Do not auto-merge.
Do not create tags unless the user explicitly asks.
End with session advice. Usually recommend starting a new session after a shipped mission or major phase boundary, with `/sc-status` as the pickup command.
