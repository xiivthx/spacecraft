---
name: sc-goal-roadmap
description: "Optional Goals-mirror for multi-mission roadmaps (Sizing: roadmap). Upstream CreateGoal/UpdateGoal by reference. Never required for ready/ship. Goals ≠ AUTH / VERIFIED / ready / ship authority."
disable-model-invocation: true
---

# sc-goal-roadmap

## Goal

Optional Goals-mirror for multi-mission roadmaps (`Sizing: roadmap`). One Goal per roadmap (not per tip). Spacecraft `mission.json` + roadmap JSON remain SoT.

## Output

Greppable disposition (exactly one):

- `Goal-roadmap: ran`
- `Goal-roadmap: skipped: <reason>`

## When to use

Multi-mission roadmap work when optional human-visibility Goals-mirror helps. Skip single-mission (Goals-mirror N/A) or when unwanted → `Goal-roadmap: skipped: <reason>`. Never required for every roadmap discuss/run.

## Workflow

1. **Optional propose** - may propose a Goal objective from the active roadmap (one Goal per roadmap, not per tip).
2. **CreateGoal** - only after human ask/confirm.
3. **UpdateGoal complete** - only when last tip on that map ships AND a Goal for that roadmap exists; else skip. No Goal created → skip complete; single-mission → N/A skip.
4. **Disposition** - `Goal-roadmap: ran` when CreateGoal and/or UpdateGoal ran; else `Goal-roadmap: skipped: <reason>`.

Shared firewall: [../sc-run/references/optional-lanes.md](../sc-run/references/optional-lanes.md)

## Verify

- Disposition: `Goal-roadmap: ran` | `Goal-roadmap: skipped: <reason>`
- CreateGoal only after human ask/confirm; one Goal per roadmap; UpdateGoal complete only when last tip ships and Goal exists
- Spacecraft `mission.json` + roadmap JSON remain SoT

## Must not

- Treat Goals / Goal complete as AUTH / `VERIFIED` / ready / ship authority
- CreateGoal without human ask/confirm
- Invent tip-by-tip progress (API has none)
- Replace `mission.json` / roadmap JSON with Goals as SoT
