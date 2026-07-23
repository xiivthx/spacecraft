---
name: sc-ship
description: "Close out a ready mission on explicit /sc-ship: validate, merge --no-ff, archive."
disable-model-invocation: true
---

## Goal

Close out a `ready` mission only after explicit human `/sc-ship`: validate, merge `--no-ff`, archive.

## Output

Shipped mission on `main` (or blocked with exact missing gates). Never infer ship from AFK or handoff.

## Good / Bad

- Good: AFK checkpoints squashed to ≤5 Conventional Commits; `validate --strict` + `closeout-check` pass; CHANGELOG + version bump committed; branch stripped to `<type>/<title>` before merge; `SPACECRAFT_SHIP=1` for gated git
- Bad: shipping with unsquashed WIP checkpoints; shipping without evidence/review; merging while still named `<type>/<id>/<title>`; deferring changelog

## Verify

`spacecraft validate --strict` and `spacecraft closeout-check` (`ship-check`) both exit 0 before merge.

Use sc-mission, sc-verification, sc-git, sc-learn. Resolve the mission. Block if unsafe.

## Pre-flight

Read resolved mission artifacts + git diff when available. Run in order:

```
spacecraft validate --strict
spacecraft git-info
git log main..HEAD --oneline -- CHANGELOG.md | grep -q . || { echo "FAIL: CHANGELOG.md not updated"; exit 1; }
spacecraft closeout-check
```

Changelog + version bump are mandatory before merge (sc-git tags next patch even for chores).

With `SPACECRAFT_SHIP=1`, the ship hook re-runs closeout before allow. Preflight does not replace that.

## Workflow

Treat ship/release/merge/finish-mission/close-branch as closeout. Do not treat session handoff as ship.

### 1. Gates

Only merge if: clarify clear; acceptance has evidence; `review.json` status `ready`; no critical/`blocksShip` findings; sc-git hygiene OK; CHANGELOG + version bump committed; tag plan set; UI changes have design review + art direction in `decisions.md` when applicable.

If any fail: block with exact missing actions.

### 2. Squash AFK checkpoints

Before release commit / merge, rewrite the work-branch history into ≤5 logical Conventional Commits (target 1–3 + separate changelog commit):

1. Inspect `git log main..HEAD --oneline`. If already ≤5 logical commits and no `wip checkpoint` noise, skip.
2. Soft-reset or equivalent non-interactive squash/fixup of RED/GREEN/refactor checkpoints into coherent `feat:` / `fix:` / `test:` / `refactor:` commits. Prefer `git reset --soft $(git merge-base HEAD main)` then re-commit logical groups - never interactive rebase (`-i`).
3. Re-run mission verify / tests after squash; capture fresh evidence if history rewrite dropped nothing but reaffirm suite green.
4. Record squash summary in the ship handoff (how many checkpoints → how many final commits).

Do not push rewritten history unless the user explicitly asks.

### 3. Release commit

Before merge, one commit: sc-learn migration + CHANGELOG + version bump.

1. Read `issues.md` / `solved.md` / `learned.md`; open GitHub issues for unresolved; migrate solved/learned to `docs/learned.md`; mark migrated (do not delete).
2. Update CHANGELOG.md.
3. Bump version.
4. Commit.

### 4. Merge

- Finish all gates, release commit, and closeout while still on `<type>/<id>/<title>` (resolver needs the id segment).
- Rebase on latest `main` after confirming fork point; reverify.
- **Strip mission id from the branch name** immediately before merge:
  - If current branch matches `<type>/<id>/<title>` and `<id>` is a mission id (`M…`), rename to `<type>/<title>`:
    `git branch -m <type>/<title>`
  - Example: `feat/M07FP1L7Z/go-rewrite` → `feat/go-rewrite`
  - Skip rename if already `<type>/<title>` or not a mission-id pattern.
  - Optional: `spacecraft bind-branch <id>` after rename so `mission.json.branches` stays accurate.
  - If the old name was pushed: delete the remote old name after merge (or when cleaning up); do not block merge on remote rename.
- Merge `--no-ff` only (no squash), using the stripped name.
- Tag annotated `vX.Y.Z`; delete local branch unless asked to keep.
- `spacecraft archive` unless user keeps the live mission folder.
- No push unless explicitly requested.

#### SPACECRAFT_SHIP

Hooks deny `git merge` / `push` / `tag` unless `SPACECRAFT_SHIP=1` and closeout passes. Prefer per-command env:

```
# after strip: feat/<title> (not feat/<id>/<title>)
SPACECRAFT_SHIP=1 git merge --no-ff feat/<title>
SPACECRAFT_SHIP=1 git push origin main
SPACECRAFT_SHIP=1 git tag -a vX.Y.Z -m "vX.Y.Z"
```

Unset after ship ops. Never set for ordinary git work.

### 5. Summary

Mission id, what changed, evidence, review status, git/tag state, limitations, next step. Then `set-state shipped` and archive if not already.

After archive, surface any CLI next-mission lines in the ship summary. If archive was silent, still run:

```
spacecraft map current
spacecraft map next <roadmap-id>
```

When current roadmap exists and `map next` is not `All missions complete.`, recommend that mission:

`Next: new session → /sc-discuss <id> (then /sc-run)`

Do not auto-start discuss or run. If no current roadmap or all complete: recommend a new session; pickup via `spacecraft status`.

## Hard stops

Resolver conflict; clarify open; missing evidence; review not ready; critical findings; sc-git fail; no CHANGELOG/version commit; UI without design review.

## Errors

- No push unless user asks.
- Conventional Commits: `<type>: <description>`; body `- ` bullets, lowercase start. Do not put mission ids in commit subjects or bodies (including merge commits).
- On gate fail: list exact missing actions.
