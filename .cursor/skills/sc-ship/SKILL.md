---
name: sc-ship
description: "Prepares and executes mission delivery when review gates pass and the user explicitly requests shipping."
disable-model-invocation: true
---

Use sc-mission, sc-verification, sc-git, and sc-learn.
Resolve the mission. Block if unsafe.

## Pre-flight checks

Read the resolved mission's spec.md, plan.json, evidence.jsonl, review.md, review.json, questions.md, decisions.md, and git diff when git is available.

Run these in order. Each must pass before the next:

```
spacecraft validate
spacecraft git-info
# Hard gate: changelog and version bump. Both mandatory - never defer.
# Checks that the work branch has a commit touching CHANGELOG.md since fork.
git log main..HEAD --oneline | grep -q 'CHANGELOG' || { echo "FAIL: CHANGELOG.md not updated. Add a version bump + changelog commit before merge."; exit 1; }
spacecraft closeout-check
```

Changelog update and version bump are **never deferrable**. sc-git requires even docs/chore changes to tag the next patch version. Do not proceed with merge until these commits exist on the work branch.

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
- changelog updated with this merge's changes (mandatory - never defer)
- version bump committed (mandatory - never defer. sc-git tags next patch even for chores)
- tag plan for the bumped version
- if UI files changed, review.md or review.json includes a design review result
- UI work has no unresolved critical design findings
- if UI files changed, art direction decisions are recorded in decisions.md

If any gate fails, block closeout. List exact missing actions and next command.

### 2. Migrate mission knowledge and prepare release commit

Before merge, combine knowledge migration with version bump and changelog in one commit:

1. Use sc-learn to preserve what was learned:
   - Read `.space/missions/<id>/issues.md`, `solved.md`, `learned.md`.
   - For unresolved issues: create GitHub Issues with mission context.
   - Migrate solved items and lessons to `docs/learned.md` with mission context.
   - Mark mission files as migrated. Do not delete them - they archive with the mission.

2. Update CHANGELOG.md with this merge's changes.

3. Bump version (sc-git tags next patch even for chores).

4. Commit all three together: knowledge migration + changelog + version bump.

### 3. Prepare merge

If all gates pass, use sc-git to prepare merge to main:
- confirm fork point and rebase target: `git log --oneline main..HEAD | head -1` to identify fork point, then rebase on `main` only after confirming `main` HEAD matches expected base
- never squash - merge with `--no-ff` preserving all granular commits
- tag, branch cleanup
- compact shipped mission artifacts with `spacecraft archive` unless the user asks to keep the full live mission folder
- no push unless explicitly requested

#### SPACECRAFT_SHIP gate (required)

Cursor hooks deny `git merge`, `git push`, and `git tag` unless `SPACECRAFT_SHIP=1`. Set it only for those gated ship commands, then unset immediately:

```
# Prefer per-command env (auto-clears):
SPACECRAFT_SHIP=1 git merge --no-ff feat/<id>/<title>
SPACECRAFT_SHIP=1 git push origin main
SPACECRAFT_SHIP=1 git tag vX.Y.Z

# Or export for a short block, then unset:
export SPACECRAFT_SHIP=1
git merge --no-ff feat/<id>/<title>
git push origin main
git tag vX.Y.Z
unset SPACECRAFT_SHIP
```

Never leave `SPACECRAFT_SHIP=1` exported in the shell after ship ops. Do not set it for ordinary non-ship git work.

### 4. Produce summary

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

Then set state to shipped if appropriate. After state is shipped and release closeout is complete, run `spacecraft archive` to move the mission from `.space/missions/` to `.space/archive/` with compact durable artifacts, unless the user asks to keep the full live mission folder.

## Research auto-trigger

sc-ship gates are verification gates - research should have happened during sc-plan and sc-build. If a gate check reveals unexpected behavior or version conflicts, use sc-search (WebSearch/WebFetch) for `"<topic>"` before blocking.

## Hard stop gates

- Resolver conflict or ambiguity
- Blocking clarification unresolved
- Missing evidence for any acceptance check
- review.json status not `ready`
- Critical findings unresolved
- sc-git gates fail (branch hygiene, commit style, rebase, merge plan)
- Changelog not updated (mandatory - never defer)
- Version bump not committed (mandatory - never defer)
- UI changes without design review recorded

## Error handling

- Do not git push unless the user explicitly asks.
- Suggested commit messages must follow Conventional Commits: `<type>: <description>` - no scope by default; body uses `- ` bullet points with lowercase first character.
- If gates fail, block closeout with exact missing actions listed.

End with session advice. Usually recommend a new session after a shipped mission or major phase boundary, with sc-mission status as pickup.
