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
- Bad: ship without ready/closeout; review findings present; merge while still `<type>/<id>/<title>`

## Must / Must not

- **Must**: Promote only durable product contracts from mission working notes into `docs/specs` or `docs/architecture/decisions` (ADR), in human engineering language.
- **Must not**: Dump every mission discuss into `docs/`.
- **Must not**: Add AI-flavored filenames or framing into product docs/.

## Verify

```
spacecraft validate --strict
spacecraft closeout-check
```

Both exit 0 before merge. Closeout enforces empty review findings (and other closeout gates).

## Lifecycle

Canonical: `.cursor/rules/200-workflow.mdc`. This skill is ship only after human check of `ready` work.

## Pre-flight

```
spacecraft validate --strict
git status -sb
git rev-parse --abbrev-ref HEAD
git log main..HEAD --oneline -- CHANGELOG.md | grep -q . || { echo "FAIL: CHANGELOG.md not updated"; exit 1; }
spacecraft closeout-check
```

Changelog + version bump mandatory before merge. Ship hook re-runs closeout when `SPACECRAFT_SHIP=1`.

**Companion dispositions (Must before merge):** greppable `Post-ready drain: ran` | `Post-ready drain: skipped: <reason>` **and** `Split-to-prs: ran` | `Split-to-prs: skipped: <reason>` in mission `decisions.md` (or equivalent greppable log). Silence ⇒ **block** ship; invoke `/sc-run` post-ready companions or emit explicit skip. See `sc-run/references/optional-lanes.md`. Lane success is never ship authority.

Out-of-scope depth: stamp `Post-ship UX depth:` and/or `Interop/limitation:` ending `Next: /sc-discuss` (`sc-run/references/follow-up-dispositions.md`). Stubs are not ship blockers; do not invent lanes.

## Workflow

### 1. Squash AFK checkpoints

Rewrite work-branch history into ≤5 logical Conventional Commits (target 1–3 + changelog commit):

1. `git log main..HEAD --oneline` - skip if already clean (≤5, no `wip checkpoint`).
2. Soft-reset or equivalent non-interactive squash (`git reset --soft $(git merge-base HEAD main)` then re-commit). Never `rebase -i`.
3. Re-run verify/tests; fresh evidence if needed.
4. Note squash summary in handoff.

Do not push rewritten history unless user asks.

### 2. Release commit

One commit: CHANGELOG + version bump.

1. If `review.json` still has findings → **block** (fix in `/sc-run`, re-ready).
2. Update CHANGELOG.md; bump version; commit.

### 3. Merge

- Finish gates while still on `<type>/<id>/<title>` (resolver needs id).
- Rebase on latest `main`; reverify.
- Strip mission id: `git branch -m <type>/<title>` (e.g. `feat/M…/go-rewrite` → `feat/go-rewrite`).
- Optional: `spacecraft bind-branch <id>` after rename.
- `SPACECRAFT_SHIP=1 git merge --no-ff <type>/<title>`
- Annotated tag `vX.Y.Z`; delete local branch unless asked to keep; `spacecraft archive` unless kept live.
- No push unless user asks. After `SPACECRAFT_SHIP=1`, the ship hook still **asks** the human before `git push` runs.

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

On handoff (especially when roadmap has a next mission), set or update that mission's optional `mission.json` `pickup` (`phase`, `next` one-liner, `updatedAt`) so `spacecraft status` / session-start shows Pickup. Not a closeout or ship gate.

## Hard stops

Any `closeout-check` or `validate --strict` failure; clarify open; sc-git fail; no CHANGELOG/version commit; UI without design review when required; missing `Post-ready drain:` or `Split-to-prs:` disposition. List exact missing actions.

## Errors

No push unless user asks (hook still asks after `SPACECRAFT_SHIP=1`). Conventional Commits; no mission ids in subjects/bodies.
