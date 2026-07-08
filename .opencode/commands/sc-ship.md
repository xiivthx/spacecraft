---
description: Prepare final delivery summary
agent: sc-commander
subtask: true
---
Use sc-mission, sc-verification, and sc-git.
Resolve the mission. Block if unsafe.

## Pre-flight checks

Read the resolved mission's spec.md, plan.json, evidence.jsonl, review.md, review.json, questions.md, decisions.md, and git diff when git is available.

Run:
```
scripts/spacecraft validate
scripts/spacecraft git-info
scripts/spacecraft closeout-check
```

## Workflow

Treat "ship", "release", "merge", "finish mission", and "close branch" as release closeout requests. Do not treat ordinary session handoff or "continue in a new session" as release closeout.

### 1. Check closeout gates

Only close/ship/merge if:
- blocking clarification questions are resolved
- acceptance checks have evidence
- important verification commands have passing evidence
- review.json status is ready
- there are no critical findings
- sc-git gates pass: branch hygiene, commit style, rebase status, merge plan
- changelog updated with this merge's changes (mandatory — never defer)
- version bump complete or explicitly deferred with rationale
- tag plan exists for the bumped version
- if UI files changed, review.md or review.json includes a design review result
- UI work has no unresolved critical design findings
- if UI files changed, art direction decisions are recorded in decisions.md or explicitly deferred

If any gate fails, block closeout. List exact missing actions and next command.

### 2. Prepare merge

If all gates pass, use sc-git to prepare merge to main:
- rebase, verify, merge with `--no-ff`, tag, branch cleanup
- compact shipped mission artifacts with `scripts/spacecraft archive` unless the user asks to keep the full live mission folder
- no push unless explicitly requested

### 3. Produce summary

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

Then set state to shipped if appropriate. After state is shipped and release closeout is complete, run `scripts/spacecraft archive` to move the mission from `.space/missions/` to `.space/archive/` with compact durable artifacts, unless the user asks to keep the full live mission folder.

## Error handling

- Do not git push unless the user explicitly asks.
- Suggested commit messages must follow Conventional Commits: `<type>: <description>` — no scope by default; body uses `- ` bullet points with lowercase first character.
- If gates fail, block closeout with exact missing actions listed.

End with session advice. Usually recommend a new session after a shipped mission or major phase boundary, with sc-mission status as pickup.
