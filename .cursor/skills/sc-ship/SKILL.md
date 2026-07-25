---
name: sc-ship
description: "Close out a ready mission on explicit /sc-ship: validate, merge --no-ff, archive."
disable-model-invocation: true
---

# sc-ship

## Goal

Close out a `ready` mission only after explicit human `/sc-ship`: validate, merge `--no-ff`, archive.

## Output

Shipped mission on `main` (or blocked with exact missing gates). Never infer ship from AFK or handoff.

## Good / Bad

- Good: squash checkpoints ≤5; `validate --strict` + `closeout-check` pass; CHANGELOG + version bump; strip branch to `<type>/<title>`; `SPACECRAFT_SHIP=1`
- Bad: ship without ready/closeout; open `issues.md` or review findings; create GitHub Issues at ship for parked debt; merge while still `<type>/<id>/<title>`

## Verify

```
spacecraft validate --strict
spacecraft closeout-check
```

Both exit 0 before merge. Closeout enforces 0 open issues and empty review findings.

## Lifecycle

Canonical: `.cursor/rules/200-workflow.mdc`. This skill is ship only after human check of `ready` work.

## Pre-flight

```
spacecraft validate --strict
spacecraft git-info
git log main..HEAD --oneline -- CHANGELOG.md | grep -q . || { echo "FAIL: CHANGELOG.md not updated"; exit 1; }
spacecraft closeout-check
```

Changelog + version bump mandatory before merge. Ship hook re-runs closeout when `SPACECRAFT_SHIP=1`.

## Workflow

### 1. Squash AFK checkpoints

Rewrite work-branch history into ≤5 logical Conventional Commits (target 1–3 + changelog commit):

1. `git log main..HEAD --oneline` - skip if already clean (≤5, no `wip checkpoint`).
2. Soft-reset or equivalent non-interactive squash (`git reset --soft $(git merge-base HEAD main)` then re-commit). Never `rebase -i`.
3. Re-run verify/tests; fresh evidence if needed.
4. Note squash summary in handoff.

Do not push rewritten history unless user asks.

### 2. Release commit

One commit: sc-learn migration + CHANGELOG + version bump.

1. If any `issues.md` entry is `open` → **block** (file/fix in `/sc-run`, re-ready). Migrate solved/learned only (sc-learn).
2. Update CHANGELOG.md; bump version; commit.

### 3. Merge

- Finish gates while still on `<type>/<id>/<title>` (resolver needs id).
- Rebase on latest `main`; reverify.
- Strip mission id: `git branch -m <type>/<title>` (e.g. `feat/M…/go-rewrite` → `feat/go-rewrite`).
- Optional: `spacecraft bind-branch <id>` after rename.
- `SPACECRAFT_SHIP=1 git merge --no-ff <type>/<title>`
- Annotated tag `vX.Y.Z`; delete local branch unless asked to keep; `spacecraft archive` unless kept live.
- No push unless user asks.

```
SPACECRAFT_SHIP=1 git merge --no-ff feat/<title>
SPACECRAFT_SHIP=1 git push origin main
SPACECRAFT_SHIP=1 git tag -a vX.Y.Z -m "vX.Y.Z"
```

Unset env after ship ops.

### 4. Summary

Mission id, what changed, evidence, git/tag state, next step. `set-state shipped` + archive if needed.

```
spacecraft map current
spacecraft map next <roadmap-id>
```

If roadmap has work: **Next: new session → /sc-discuss <id>**. Do not auto-start discuss/run. Else: new session; `spacecraft status`.

## Hard stops

Any `closeout-check` or `validate --strict` failure; clarify open; sc-git fail; no CHANGELOG/version commit; UI without design review when required. List exact missing actions.

## Errors

- No push unless user asks.
- Conventional Commits; no mission ids in commit subjects/bodies.
