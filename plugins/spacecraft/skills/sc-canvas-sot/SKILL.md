---
name: sc-canvas-sot
description: "Optional canvas-SoT for /sc-run plan|findings|evidence human-check emits. Upstream canvas by reference. Never required for ready/ship. Canvas ≠ AUTH / VERIFIED / ready / ship authority."
disable-model-invocation: true
---

# sc-canvas-sot

## Goal

Optional canvas-SoT for `/sc-run` plan | findings | evidence human-check emits. Humans or agents may emit; this lane is **never required** for ready or ship, and never required for every `/sc-run`. Upstream Cursor canvas is the emit-pattern SoT by reference; mission artifacts and judge remain pass/fail SoT. Must not treat canvas as AUTH / `VERIFIED` / ready / ship authority.

## Output

Greppable disposition (exactly one):

- `Canvas-sot: ran`
- `Canvas-sot: skipped: <reason>`

When a canvas is emitted, keep greppable `decisions.md` lines valid:

- `Canvas plan: ` + absolute path
- `Canvas findings: ` + absolute path
- `Canvas findings skipped: empty`
- `Canvas evidence: ` + absolute path

## Good / Bad

- Good: optional emit for human check during `/sc-run`; upstream canvas by reference; managed `canvases/` path only; disposition greppable; Canvas plan/findings/evidence lines remain valid; never required for ready/ship or every `/sc-run`
- Bad: requiring this lane for every `/sc-run` or for ready/ship; tip-deliverable under mission `.space/` or repo `.cursor/`; canvas as discuss draft HTML / mission-brief SoT; copying upstream skill body; sample fixtures under this skill dir
- Bad (firewall): Must not treat canvas as AUTH / `VERIFIED` / ready / ship authority

## When to use

During `/sc-run` when an optional plan | findings | evidence canvas helps human check. Skip when no emit is wanted (`Canvas-sot: skipped: <reason>`).

## Workflow

1. **Optional emit** - when a plan | findings | evidence human-check canvas is wanted, invoke Cursor canvas by reading `~/.cursor/skills-cursor/canvas/SKILL.md` (do not copy its body into this repo). Write under managed path only: `~/.cursor/projects/<workspace>/canvases/<missionId>-<kind>.canvas.tsx` (`kind` ∈ `plan` | `findings` | `evidence`).
2. **Record lines** - when emitted, write matching greppable lines in `decisions.md` (absolute path; chat and `decisions.md` include absolute markdown links): `Canvas plan:`, `Canvas findings:` (or `Canvas findings skipped: empty` when findings are empty), `Canvas evidence:`.
3. **Disposition** - emit `Canvas-sot: ran` when at least one canvas emit ran for this mission; otherwise `Canvas-sot: skipped: <reason>`.
4. **Skip** - if the lane is not used → `Canvas-sot: skipped: <reason>`. Missing canvas files or lines do not block judge or ready.

## Boundaries

- This lane is **never required** for ready or ship, and never required for every `/sc-run`.
- Must not treat canvas as AUTH / `VERIFIED` / ready / ship authority.
- Must not commit `.canvas.tsx` under mission `.space/` or repo `.cursor/` as tip deliverable; live files live only under managed `canvases/`.
- Must not use canvas as discuss draft HTML / mission-brief SoT (brief stays Accept/Adjust/Reject chat HIL; draft stays HTML).
- Must not copy upstream canvas skill body into this repo; reference only: `~/.cursor/skills-cursor/canvas/SKILL.md`.
- Must not add sample `.canvas.tsx` under this skill directory.
- Canvas files and Canvas plan/findings/evidence lines are optional aids for human check - not ready or `VERIFIED` proof.

## Verify

- Disposition line present: `Canvas-sot: ran` | `Canvas-sot: skipped: <reason>`
- Firewall greppable: Must not treat canvas as AUTH / `VERIFIED` / ready / ship; never required for ready/ship; managed path only; no tip-deliverable under mission `.space/` or repo `.cursor/`
- Upstream invoked only via `~/.cursor/skills-cursor/canvas/SKILL.md`
- Greppable Canvas plan / findings / findings skipped / evidence lines remain valid when emitted
