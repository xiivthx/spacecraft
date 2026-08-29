---
name: sc-goal-roadmap
description: "Optional Goals-mirror for multi-mission roadmaps (Sizing: roadmap). Upstream CreateGoal/UpdateGoal by reference. Never required for ready/ship. Goals ≠ AUTH / VERIFIED / ready / ship authority."
disable-model-invocation: true
---

# sc-goal-roadmap

## Goal

Optional Goals-mirror for multi-mission roadmaps (`Sizing: roadmap`). Agent may optionally propose a Goal objective from the active roadmap; invoke `CreateGoal` only after human ask/confirm. One Goal per roadmap (not per tip). Spacecraft `mission.json` + roadmap JSON remain SoT. This lane is **never required** for ready or ship, and never required for every roadmap discuss/run. Must not treat Goals / Goal complete as AUTH / `VERIFIED` / ready / ship authority.

## Output

Greppable disposition (exactly one):

- `Goal-roadmap: ran`
- `Goal-roadmap: skipped: <reason>`

## Good / Bad

- Good: optional Goals-mirror for multi-mission roadmaps; propose Goal objective from active roadmap then CreateGoal only after human ask/confirm; one Goal per roadmap (not per tip); UpdateGoal complete only when last tip on that map ships AND a Goal for that roadmap exists; Spacecraft `mission.json` + roadmap JSON remain SoT; disposition greppable; never required for ready/ship or every roadmap discuss/run
- Bad: requiring this lane for every roadmap discuss/run or for ready/ship; one Goal per tip; inventing tip-by-tip progress; CreateGoal without human ask/confirm; replacing `mission.json` / roadmap JSON with Goals as SoT
- Bad (firewall): Must not treat Goals / Goal complete as AUTH / `VERIFIED` / ready / ship authority

## When to use

During multi-mission roadmap work (`Sizing: roadmap`) when optional human-visibility Goals-mirror helps. Skip for single-mission work (Goals-mirror N/A) or when no Goal mirror is wanted → `Goal-roadmap: skipped: <reason>`.

## Workflow

1. **Optional propose** - agent may optionally propose a Goal objective from the active roadmap (one Goal per roadmap, not per tip).
2. **CreateGoal** - invoke CreateGoal only after human ask/confirm. Must not CreateGoal without human request.
3. **UpdateGoal complete** - may UpdateGoal complete only when the last tip on that map ships AND a Goal for that roadmap exists; otherwise skip. Edge: roadmap with no Goal created → skip complete; last tip ship with Goal → may UpdateGoal complete; single-mission work → Goals-mirror N/A (skip).
4. **Disposition** - emit `Goal-roadmap: ran` when CreateGoal and/or UpdateGoal for this roadmap mirror ran; otherwise `Goal-roadmap: skipped: <reason>`.

## Boundaries

- This lane is **never required** for ready or ship, and never required for every roadmap discuss/run.
- Must not treat Goals / Goal complete as AUTH / `VERIFIED` / ready / ship authority.
- Spacecraft `mission.json` + roadmap JSON remain SoT; Goals are optional human-visibility only.
- Must not invent tip-by-tip progress (API has none).
- Must not CreateGoal without human ask/confirm.
- Goals and Goal complete are optional aids for human visibility - not ready or `VERIFIED` proof.

## Verify

- Disposition line present: `Goal-roadmap: ran` | `Goal-roadmap: skipped: <reason>`
- Firewall greppable: Must not treat Goals as AUTH / `VERIFIED` / ready / ship; never required for ready/ship
- CreateGoal only after human ask/confirm; one Goal per roadmap; UpdateGoal complete only when last tip ships and Goal exists
- Spacecraft `mission.json` + roadmap JSON remain SoT
