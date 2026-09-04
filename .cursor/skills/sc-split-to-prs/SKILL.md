---
name: sc-split-to-prs
description: "Companion post-ready / pre-ship fat-diff hygiene via Cursor split-to-prs. Must disposition (ran or skipped). Human plan approval required. Must not replace discuss sizing or merge pull requests. Never ready/ship authority."
disable-model-invocation: true
---

# sc-split-to-prs

## Goal

Companion **post-ready / pre-ship** fat-diff hygiene via Cursor `split-to-prs`. Commander **Must** leave a greppable disposition after post-ready drain (`ran` or `skipped:`). Spacecraft AUTH + `/sc-ship` remain merge/tag SoT. Split success is not ready or ship authority.

## Output

Greppable disposition (exactly one; silence forbidden):

- `Split-to-prs: ran`
- `Split-to-prs: skipped: <reason>`

## When to use

**Must** invoke after `/sc-run` reaches `ready` (after `sc-post-ready-drain` disposition), before human `/sc-ship`. Run split when the human judges the diff fat; else `Split-to-prs: skipped: <reason>`. Not a mid-discuss/run sizing substitute.

## Workflow

1. **Human plan approval** - human Must approve the split plan before any branch/PR work.
2. **Quoted AUTH** - state quoted `AUTH:` before any outward push or PR publish.
3. **Invoke upstream** - follow `~/.cursor/skills-cursor/split-to-prs/SKILL.md` (reference only).
4. **Re-ready if mutated** - any new commit on the mission work branch after `ready` → stop for **re-ready** (do not emit `ran` yet).
5. **Report** - `Split-to-prs: ran` only when no unre-reviewed post-ready commits remain; else `Split-to-prs: skipped: <reason>` (e.g. awaiting re-ready, or diff not fat).

**Re-ready:** post-ready work-branch mutation forces re-enter `/sc-run` review→judge→`ready` before claiming success or `/sc-ship`. Local commit enough (push not required).

Shared firewall: [../sc-run/references/optional-lanes.md](../sc-run/references/optional-lanes.md)

## Verify

- Disposition: `Split-to-prs: ran` | `Split-to-prs: skipped: <reason>`
- Human approved split plan before branch/PR work when `ran`
- Quoted `AUTH:` before any outward push or PR publish when those occur
- No unre-reviewed post-ready commits when claiming `ran`

## Must not

- Omit disposition after ready (silence)
- Merge pull requests; turn on GitHub automatic merge; merge via GitHub merge controls
- Replace `/sc-discuss` sizing / roadmap resize / mid-run map resize
- Treat opened PRs or split chat as ready/ship/`VERIFIED`/AUTH skip
- Push or PR publish without quoted `AUTH:`
- Claim `Split-to-prs: ran` while unre-reviewed post-ready commits remain
- Copy full upstream body into this repo
