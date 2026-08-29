---
name: sc-post-ready-drain
description: "Optional git-primary post-ready drain after mission ready and before /sc-ship. Resolve conflicts vs main, run scoped local verify. Cursor autopilot only when an open PR already exists. Never merge into main or PRs."
disable-model-invocation: true
---

# sc-post-ready-drain

## Goal

Optional **git-primary** post-ready drain after mission `ready` and before `/sc-ship`. Spacecraft AUTH + `/sc-ship` remain merge/tag authority. Drain success is not ship authority.

## Output

Greppable disposition (exactly one):

- `Post-ready drain: ran`
- `Post-ready drain: skipped: <reason>`

## When to use

After `/sc-run` reaches `ready`, before human `/sc-ship`. Else `Post-ready drain: skipped: <reason>`.

## Workflow

1. **Resolve conflicts vs latest `main`** - update the **work branch** from latest `main` (prefer `git rebase main`, or `git merge main` into the work branch only). Sync *onto* the work branch - not a merge into `main`. Git-primary; no GitHub PR required. Any new commit → stop for **re-ready** (do not emit `ran` yet).
2. **Scoped local verify** - mission scoped verify/tests; when claiming success include `spacecraft validate --strict`.
3. **Optional Cursor autopilot** - only when an **open PR** already exists; invoke `~/.cursor/skills-cursor/autopilot/SKILL.md` (reference only). Skip when no open PR. Autopilot commits also force **re-ready**.
4. **Report** - `Post-ready drain: ran` only when no unre-reviewed post-ready commits remain; else `Post-ready drain: skipped: <reason>` (e.g. awaiting re-ready).

**Re-ready:** any new commit on the work branch after `ready` forces re-enter `/sc-run` review→judge→`ready` before `/sc-ship`. Local commit enough (push not required).

Shared firewall: [../sc-run/references/optional-lanes.md](../sc-run/references/optional-lanes.md)

## Verify

- Disposition: `Post-ready drain: ran` | `Post-ready drain: skipped: <reason>`
- When `ran`: conflicts vs `main` addressed; scoped verify including `spacecraft validate --strict` when claiming success
- No merge into `main` / PR merge / auto-merge / mark-draft-ready from this lane

## Must not

- Merge into `main` or merge PRs (work-branch sync from `main` in step 1 is allowed)
- Enable auto-merge or mark-draft-ready
- Treat drain success as ship authority
- Claim `Post-ready drain: ran` while unre-reviewed post-ready commits remain
