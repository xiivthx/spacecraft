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

- Good: `validate --strict` + `closeout-check` pass; CHANGELOG + version bump committed; `SPACECRAFT_SHIP=1` for gated git
- Bad: shipping without evidence/review; deferring changelog

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

### 2. Release commit

Before merge, one commit: sc-learn migration + CHANGELOG + version bump.

1. Read `issues.md` / `solved.md` / `learned.md`; open GitHub issues for unresolved; migrate solved/learned to `docs/learned.md`; mark migrated (do not delete).
2. Update CHANGELOG.md.
3. Bump version.
4. Commit.

### 3. Merge

- Rebase on latest `main` after confirming fork point.
- Merge `--no-ff` only (no squash).
- Tag annotated `vX.Y.Z`; delete local branch unless asked to keep.
- `spacecraft archive` unless user keeps the live mission folder.
- No push unless explicitly requested.

#### SPACECRAFT_SHIP

Hooks deny `git merge` / `push` / `tag` unless `SPACECRAFT_SHIP=1` and closeout passes. Prefer per-command env:

```
SPACECRAFT_SHIP=1 git merge --no-ff feat/<id>/<title>
SPACECRAFT_SHIP=1 git push origin main
SPACECRAFT_SHIP=1 git tag -a vX.Y.Z -m "vX.Y.Z"
```

Unset after ship ops. Never set for ordinary git work.

### 4. Summary

Mission id, what changed, evidence, review status, git/tag state, limitations, next step. Then `set-state shipped` and archive if not already.

## Hard stops

Resolver conflict; clarify open; missing evidence; review not ready; critical findings; sc-git fail; no CHANGELOG/version commit; UI without design review.

## Errors

- No push unless user asks.
- Conventional Commits: `<type>: <description>`; body `- ` bullets, lowercase start.
- On gate fail: list exact missing actions.

After ship, recommend a new session; pickup via `spacecraft status`.
