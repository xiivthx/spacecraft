---
name: sc-loop
description: "Optional overnight/AFK /sc-run watch via Cursor /loop (CI/jobs). Never required for ready/ship. On Spacecraft stop: disarm loop + handback.md."
disable-model-invocation: true
---

# sc-loop

## Goal

Optional overnight/AFK `/sc-run` watch via Cursor `/loop` (CI/jobs). Humans or agents may arm it. Not a `/sc-discuss` HIL substitute.

## Output

Greppable disposition (exactly one):

- `Loop watch: ran`
- `Loop watch: stopped: <reason>`
- `Loop watch: skipped: <reason>`

## When to use

Overnight/AFK `/sc-run` when watching CI/jobs helps. Else `Loop watch: skipped: <reason>`.

## Workflow

1. **Optional arm** - invoke Cursor `/loop` via `~/.cursor/skills-cursor/loop/SKILL.md` (reference only). Emit `Loop watch: ran` when armed.
2. **Tick work** - on wake, run watch prompt (status, CI, jobs). Ticks only resume mid-lifecycle under existing `/sc-run` gates. Must not `set-state ready` or progress ship.
3. **Spacecraft stop** - on `3-cycle` | `timebox` | `blocked`, when armed: **disarm** (cloud unsubscribe and/or local PID kill) + write `.space/missions/<id>/handback.md` → `Loop watch: stopped: <reason>`.
4. **Ready handoff** - when armed at ready: **disarm** → `Loop watch: stopped: ready-handoff` (or `stopped: <reason>` with ready-handoff as example). Do not leave loop armed into post-ready drain / ship wait.
5. **Stop without prior loop** - no watch armed → `Loop watch: skipped: <reason>`; still write `handback.md` per `/sc-run`.

Shared firewall: [../sc-run/references/optional-lanes.md](../sc-run/references/optional-lanes.md)

## Verify

- Disposition: `Loop watch: ran` | `Loop watch: stopped: <reason>` | `Loop watch: skipped: <reason>`
- Stop or ready handoff with prior watch: disarm; on stop also `handback.md`
- Upstream only via `~/.cursor/skills-cursor/loop/SKILL.md`

## Must not

- Leave Cursor loop armed after Spacecraft stop or ready handoff
- Treat loop ticks, CI green, or watch success as ready/ship/`VERIFIED`/AUTH
- `set-state ready` or progress ship from ticks
- Substitute for `/sc-discuss` HIL
- Copy upstream `/loop` body into this repo
