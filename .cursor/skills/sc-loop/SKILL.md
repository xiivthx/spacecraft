---
name: sc-loop
description: "Companion overnight/AFK /sc-run watch via Cursor /loop (CI/jobs). Must disposition at /sc-run start (ran or skipped). On Spacecraft stop or ready: disarm when armed + handback.md. Never ready/ship authority."
disable-model-invocation: true
---

# sc-loop

## Goal

Companion overnight/AFK `/sc-run` watch via Cursor `/loop` (CI/jobs). Commander **Must** leave a greppable disposition at `/sc-run` start (`ran` or `skipped:`). Not a `/sc-discuss` HIL substitute. Watch success is never ready/ship/`VERIFIED`/AUTH.

## Output

Greppable disposition (exactly one active line at a time; silence forbidden):

- `Loop watch: ran`
- `Loop watch: stopped: <reason>`
- `Loop watch: skipped: <reason>`

## When to use

**Must** invoke at every `/sc-run` start. Arm when overnight/AFK CI/job watch helps → `Loop watch: ran`. Else → `Loop watch: skipped: <reason>`.

## Workflow

1. **Start disposition** - at `/sc-run` step 1: either arm Cursor `/loop` via `~/.cursor/skills-cursor/loop/SKILL.md` (reference only) → `Loop watch: ran`, or skip → `Loop watch: skipped: <reason>`. Silence forbidden.
2. **Tick work** - on wake, run watch prompt (status, CI, jobs). Ticks only resume mid-lifecycle under existing `/sc-run` gates. Must not `set-state ready` or progress ship.
3. **Spacecraft stop** - on `3-cycle` | `timebox` | `blocked`, when armed: **disarm** (cloud unsubscribe and/or local PID kill) + write `.space/missions/<id>/handback.md` → `Loop watch: stopped: <reason>`.
4. **Ready handoff** - when armed toward ready: **disarm** → `Loop watch: stopped: ready-handoff` **before** Commander invokes `sc-judge` (judge requires `stopped:` whenever prior `ran`). Do not leave loop armed into judge / post-ready drain / ship wait.
5. **Stop without prior arm** - no watch armed → keep/emit `Loop watch: skipped: <reason>`; still write `handback.md` per `/sc-run`.

Shared firewall: [../sc-run/references/optional-lanes.md](../sc-run/references/optional-lanes.md)

## Verify

- Disposition present: `Loop watch: ran` | `Loop watch: stopped: <reason>` | `Loop watch: skipped: <reason>`
- Stop or ready handoff with prior watch: disarm; on stop also `handback.md`
- Upstream only via `~/.cursor/skills-cursor/loop/SKILL.md`

## Must not

- Omit disposition at `/sc-run` start (silence)
- Leave Cursor loop armed after Spacecraft stop or ready handoff
- Treat loop ticks, CI green, or watch success as ready/ship/`VERIFIED`/AUTH
- `set-state ready` or progress ship from ticks
- Substitute for `/sc-discuss` HIL
- Copy upstream `/loop` body into this repo
