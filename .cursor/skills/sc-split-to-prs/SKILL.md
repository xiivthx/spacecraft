---
name: sc-split-to-prs
description: "Optional post-ready / pre-ship fat-diff hygiene via Cursor split-to-prs. Human plan approval required. Must not replace discuss sizing or merge pull requests. Never ready/ship authority."
disable-model-invocation: true
---

# sc-split-to-prs

## Goal

Optional **post-ready / pre-ship** fat-diff hygiene via Cursor `split-to-prs`. Humans may run it when the work branch diff is large; it is **never required** for `ready` or `/sc-ship`. Spacecraft AUTH + `/sc-ship` remain merge/tag SoT. Split success is not ready or ship authority.

## Output

Greppable disposition (exactly one):

- `Split-to-prs: ran`
- `Split-to-prs: skipped: <reason>`

## Good / Bad

- Good: optional `split-to-prs` after `ready`, before `/sc-ship`, for fat diffs only; human-approved split plan first; quoted `AUTH:` before any outward push or PR publish; sizing seams / roadmap contract untouched; thin skill + upstream by reference; disposition greppable; re-ready after post-ready work-branch mutation; AUTH + ready/ship unchanged
- Bad: auto-split without human plan approval; push or PR publish without quoted `AUTH:`; using split mid-discuss/run as a substitute for mission sizing / map resize; merging pull requests or turning on GitHub automatic merge from this lane; soft-passing ready/ship because PRs opened; claiming `Split-to-prs: ran` with unre-reviewed post-ready commits; requiring split before every ship; copying full upstream body into the repo

## When to use

After `/sc-run` reaches `ready`, before human `/sc-ship`, when the human judges the diff fat. Skip otherwise (`Split-to-prs: skipped: <reason>`). Not a mid-discuss/run sizing substitute.

## Workflow

1. **Human plan approval** - the human Must approve the split plan before any branch/PR work.
2. **Quoted AUTH** - state quoted `AUTH:` before any outward push or PR publish.
3. **Invoke upstream** - read and follow `~/.cursor/skills-cursor/split-to-prs/SKILL.md` (do not copy its body into this repo).
4. **Re-ready if mutated** - if split creates any new commit on the mission work branch after `ready`, stop for **re-ready** before claiming split success or shipping (do not emit `Split-to-prs: ran` yet).
5. **Report** - emit `Split-to-prs: ran` only when no unre-reviewed post-ready commits remain on the mission work branch; otherwise `Split-to-prs: skipped: <reason>` (e.g. awaiting re-ready).

## Boundaries

- Must not replace `/sc-discuss` sizing. Sizing seams stay Spacecraft discuss-owned. Must not treat a fat branch as roadmap resize. Must not replace mid-run map resize.
- Must not merge pull requests from this lane. Must not turn on GitHub automatic merge. Must not merge via GitHub merge controls from this lane.
- Opened PRs are never ready or ship proof. Must not treat split chat output as ready/ship authority, `VERIFIED`, or skip AUTH.
- Quoted `AUTH:` is required before any outward push or PR publish from this lane.
- Lane is **review hygiene** only - not a mid-discuss/run sizing substitute.
- This lane is **never required** for `/sc-ship`.
- Must not claim `Split-to-prs: ran` while HEAD has commits after the last `ready` that have not completed review→judge→`ready` again.

## Re-ready

If split mutates the mission work branch after `ready`, Must **re-ready**: re-enter `/sc-run` review→judge→`ready` again before claiming success or `/sc-ship`. Do not emit `Split-to-prs: ran` while unre-reviewed post-ready commits remain on the mission work branch. A local commit is enough (push not required).

## Verify

- Disposition line present: `Split-to-prs: ran` or `Split-to-prs: skipped: <reason>`
- Human approved the split plan before branch/PR work when `ran`
- Quoted `AUTH:` present before any outward push or PR publish when those occur
- Discuss sizing / roadmap seams untouched; no merge of pull requests from this lane
- No unre-reviewed post-ready commits on the mission work branch when claiming `ran`
