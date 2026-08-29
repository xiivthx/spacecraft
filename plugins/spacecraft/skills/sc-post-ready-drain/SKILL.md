---
name: sc-post-ready-drain
description: "Optional git-primary post-ready drain after mission ready and before /sc-ship. Resolve conflicts vs main, run scoped local verify. Cursor autopilot only when an open PR already exists. Never merge."
disable-model-invocation: true
---

# sc-post-ready-drain

## Goal

Optional **git-primary** post-ready drain after mission `ready` and before `/sc-ship`. Humans may run it; it is **never required** for `/sc-ship`. Spacecraft AUTH + `/sc-ship` remain merge/tag authority. Drain success is not ship authority.

## Output

Greppable disposition (exactly one):

- `Post-ready drain: ran`
- `Post-ready drain: skipped: <reason>`

## Good / Bad

- Good: git-only path with no GitHub PR; conflict vs latest `main` then scoped local verify; optional Cursor `autopilot` only when an open PR already exists; re-ready after any post-ready commit
- Bad: requiring a PR or `gh`; merging or soft-passing ship; shipping with new commits after `ready` without re-ready; treating drain success as ship authority

## When to use

After `/sc-run` reaches `ready`, before human `/sc-ship`. Skip when the human ships without a drain (`Post-ready drain: skipped: <reason>`).

## Workflow

1. **Resolve conflicts vs latest `main`** - merge or rebase the work branch against latest `main`; resolve conflicts locally. Git-primary; no GitHub PR required. If this step creates any new commit, stop for **re-ready** before claiming drain success or shipping (do not emit `Post-ready drain: ran` yet).
2. **Scoped local verify** - run the mission's scoped local verify/tests. When claiming drain success, include `spacecraft validate --strict`.
3. **Optional Cursor autopilot** - only when an **open PR** already exists. Invoke by reading `~/.cursor/skills-cursor/autopilot/SKILL.md` (do not copy its body into this repo). Skip when no open PR. Autopilot commits also force **re-ready**.
4. **Report** - emit `Post-ready drain: ran` only when no unre-reviewed post-ready commits remain on the work branch; otherwise `Post-ready drain: skipped: <reason>` (e.g. awaiting re-ready).

## Boundaries

- **Never merge.** Must not merge PRs or local main from this lane.
- **Must not enable auto-merge.**
- **Must not mark-draft-ready.**
- Must not treat drain success as ship authority; AUTH + `/sc-ship` own merge/tag.
- This lane is **never required** for `/sc-ship`.
- Must not claim `Post-ready drain: ran` while HEAD has commits after the last `ready` that have not completed review→judge→`ready` again.

## Re-ready

Any new commit on the work branch after `ready` forces **re-ready**: re-enter `/sc-run` review→judge→`ready` again before `/sc-ship`. A local commit is enough (push not required). Machine closeout binding of ready HEAD is out of scope for this thin skill; policy + disposition above are the gate.

## Verify

- Disposition line present: `Post-ready drain: ran` or `Post-ready drain: skipped: <reason>`
- When `ran`: conflicts vs `main` addressed; scoped verify including `spacecraft validate --strict` when claiming success
- No merge / auto-merge / mark-draft-ready from this lane
