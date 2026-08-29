---
name: sc-loop
description: "Optional overnight/AFK /sc-run watch via Cursor /loop (CI/jobs). Never required for ready/ship. On Spacecraft stop: disarm loop + handback.md."
disable-model-invocation: true
---

# sc-loop

## Goal

Optional overnight/AFK `/sc-run` watch via Cursor `/loop` (CI/jobs). Humans or agents may arm it; it is **never required** for ready or ship. Spacecraft owns stop and ready/ship. This lane is **not** a `/sc-discuss` HIL substitute.

## Output

Greppable disposition (exactly one):

- `Loop watch: ran`
- `Loop watch: stopped: <reason>`
- `Loop watch: skipped: <reason>`

## Good / Bad

- Good: optional `/loop` watch during overnight/AFK `/sc-run` only; Spacecraft stop → disarm loop + `handback.md`; ready handoff → disarm first; upstream `/loop` by reference; disposition greppable
- Bad: requiring `/loop` for every `/sc-run`; leaving loop armed after stop or ready handoff; using this lane as discuss HIL; treating loop ticks, CI green, or watch success as ready/ship/`VERIFIED`/AUTH authority; ticks calling `set-state ready` or progressing ship

## When to use

During overnight/AFK `/sc-run` when watching CI/jobs helps. Skip when the human runs `/sc-run` without a watch (`Loop watch: skipped: <reason>`).

## Workflow

1. **Optional arm** - when watching CI/jobs during overnight/AFK `/sc-run`, invoke Cursor `/loop` by reading `~/.cursor/skills-cursor/loop/SKILL.md` (do not copy its body into this repo). Emit `Loop watch: ran` when a watch was armed.
2. **Tick work** - on each wake, run the watch prompt (status, CI, jobs). Loop ticks only resume mid-lifecycle work under existing `/sc-run` gates. Must not `set-state ready` or progress ship. Loop ticks are never ready/ship/`VERIFIED`/AUTH authority.
3. **Spacecraft stop** - on `3-cycle` | `timebox` | `blocked`, Must both:
   1. **Disarm** the Cursor loop: cloud unsubscribe and/or local PID kill (follow upstream stop steps).
   2. **Write** `.space/missions/<id>/handback.md` with stop reason + remaining work cue.
   Then emit `Loop watch: stopped: <reason>`.
4. **Ready handoff** - when overnight/AFK `/sc-run` reaches **ready** handoff and a watch was armed, Must **disarm** (cloud unsubscribe and/or local PID kill) and emit `Loop watch: stopped: ready-handoff` (or `Loop watch: stopped: <reason>` with ready-handoff as example reason). Do not leave the loop armed into post-ready drain / ship wait.
5. **Stop without prior loop** - if stop fires and no watch was armed → `Loop watch: skipped: <reason>` (e.g. no prior loop); still write `handback.md` per `/sc-run`.

## Boundaries

- This lane is **never required** for ready or ship.
- Must not substitute for `/sc-discuss` HIL.
- Must not leave a Cursor loop armed after Spacecraft stop or ready handoff.
- Loop ticks Must not `set-state ready` or progress ship; ticks only resume mid-lifecycle work under existing `/sc-run` gates.
- Must not treat loop ticks, CI green, or watch success as ready/ship/`VERIFIED`/AUTH authority.
- Must not copy upstream `/loop` skill body into this repo; reference only.

## Verify

- Disposition line present: `Loop watch: ran` | `Loop watch: stopped: <reason>` | `Loop watch: skipped: <reason>`
- On Spacecraft stop or ready handoff with a prior watch: cloud unsubscribe and/or local PID kill; on stop also `handback.md`
- Upstream invoked only via `~/.cursor/skills-cursor/loop/SKILL.md`
