---
name: sc-run
description: "AFK roadmap runner: after human clarify, auto-run incomplete missions to ready. Invoke as /sc-run <roadmap-id>."
disable-model-invocation: true
---

# sc-run

## Goal

After human clarify, AFK incomplete roadmap missions to `ready` so a human can check and then `/sc-ship`.

## Output

Missions at `state=ready` (or stop on blocked / clarify). Handoff: **AFK done. Human check, then /sc-ship.** Never merge/push/tag.

## Good / Bad

- Good: one branch per mission; plan has acceptance+verify; evidence real; review recorded; soft assumptions in `decisions.md`
- Bad: shipping; mid-loop non-blocking questions; stacking many missions on one branch

## Verify

`spacecraft validate --strict` per mission; `map next` until `All missions complete.` or blocked tip.

## Arguments

`$ARGUMENTS` = roadmap id (optional if `spacecraft map current` is set).

```
spacecraft map use <roadmap-id>
/sc-run <roadmap-id>
# or /sc-run  (uses map current)
```

## HIL vs AFK

| Phase | Who | Action |
|-------|-----|--------|
| Clarify | Human | Answer blocking questions; `spacecraft clarify-status clear` |
| Run | AI | Loop missions to `ready` |
| Check | Human | Review ready work |
| Ship | Human | Explicit `/sc-ship` only |

## Pre-flight

1. Resolve roadmap: `$ARGUMENTS` → `map use`; else `map current` (fail if missing).
2. If any incomplete mission has `clarify-status` `open` → stop; do not AFK.
3. Do not ask non-blocking questions; record assumptions in `decisions.md`.

## AFK loop

```
spacecraft map next <roadmap-id>
```

Stop when: `All missions complete.` (print handoff), tip `blocked`, or hard blocker (set clarify open / blocked).

### Per incomplete mission

1. Parse id from `map next` (`M…:` prefix); `spacecraft use <id>`.
2. One branch per mission: if not on `feat/<id>/…`, checkout from main and `spacecraft bind-branch <id>`.
3. Artifacts: minimal `spec.md` if empty; else Task(`sc-planner`) → `plan.json`; `set-state planned` then `in_progress`.
4. Build each pending task: Task(`sc-tester`) → Task(`sc-coder`) → `spacecraft evidence`; mark done. Delegate product code.
5. Review: Task(`sc-reviewer`); write `review.md` / `review.json` (`status: ready`, releaseReadiness ready); `validate --strict`; `set-state ready`.
6. Continue loop.

## Rules

- Never call `/sc-ship`, merge, push, or tag.
- One feature branch per mission id.
- Prefer `decisions.md` assumptions over mid-AFK questions.
