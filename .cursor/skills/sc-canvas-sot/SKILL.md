---
name: sc-canvas-sot
description: "Optional canvas-SoT for /sc-run plan|findings|evidence human-check emits. Upstream canvas by reference. Never required for ready/ship. Canvas ≠ AUTH / VERIFIED / ready / ship authority."
disable-model-invocation: true
---

# sc-canvas-sot

## Goal

Optional canvas-SoT for `/sc-run` plan | findings | evidence human-check emits. Upstream Cursor canvas is emit-pattern SoT by reference; mission artifacts and judge remain pass/fail SoT.

## Output

Greppable disposition (exactly one):

- `Canvas-sot: ran`
- `Canvas-sot: skipped: <reason>`

When a canvas is emitted, keep greppable `decisions.md` lines valid:

- `Canvas plan: ` + absolute path
- `Canvas findings: ` + absolute path
- `Canvas findings skipped: empty`
- `Canvas evidence: ` + absolute path

## When to use

During `/sc-run` when optional plan | findings | evidence canvas helps human check. Else `Canvas-sot: skipped: <reason>`. Never required for every `/sc-run`.

## Workflow

1. **Optional emit** - invoke canvas via `~/.cursor/skills-cursor/canvas/SKILL.md` (reference only). Write under managed path only: `~/.cursor/projects/<workspace>/canvases/<missionId>-<kind>.canvas.tsx` (`kind` ∈ `plan` | `findings` | `evidence`).
2. **Record lines** - when emitted, write matching greppable lines in `decisions.md` (absolute path; chat and `decisions.md` include absolute markdown links): `Canvas plan:`, `Canvas findings:` (or `Canvas findings skipped: empty`), `Canvas evidence:`.
3. **Disposition** - `Canvas-sot: ran` when at least one emit ran; else `Canvas-sot: skipped: <reason>`. Missing canvas files/lines do not block judge or ready.

Shared firewall: [../sc-run/references/optional-lanes.md](../sc-run/references/optional-lanes.md)

## Verify

- Disposition: `Canvas-sot: ran` | `Canvas-sot: skipped: <reason>`
- Managed `canvases/` path only; greppable Canvas plan/findings/evidence lines valid when emitted
- Upstream only via `~/.cursor/skills-cursor/canvas/SKILL.md`

## Must not

- Treat canvas as AUTH / `VERIFIED` / ready / ship authority
- Commit `.canvas.tsx` under mission `.space/` or repo `.cursor/` as tip deliverable
- Use canvas as discuss draft HTML / mission-brief SoT
- Copy upstream canvas body; add sample `.canvas.tsx` under this skill dir
