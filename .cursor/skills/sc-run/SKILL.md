---
name: sc-run
description: "AFK roadmap runner: after human clarify, auto-run incomplete missions to ready. User-facing with /sc-ship only. Invoke as /sc-run <roadmap-id>."
disable-model-invocation: true
---

# sc-run

Human clarifies gray areas first. Then AFK through the roadmap. Stop for human check. Never ship - that is `/sc-ship` only.

## Arguments

`$ARGUMENTS` = roadmap id (optional if `spacecraft map current` is set).

```
spacecraft map use <roadmap-id>   # optional persist
/sc-run <roadmap-id>
# or
/sc-run                            # uses map current
```

## HIL vs AFK

| Phase | Who | Action |
|-------|-----|--------|
| Clarify | Human | Answer blocking questions; `spacecraft clarify-status clear` |
| Run | AI (this skill) | Loop missions to `ready` |
| Check | Human | Review ready work |
| Ship | Human | Explicit `/sc-ship` only |

## Pre-flight

1. Resolve roadmap:
   - If `$ARGUMENTS` non-empty → use it and `spacecraft map use $ARGUMENTS`
   - Else → `spacecraft map current` (fail if missing)
2. Clarify gate: for each incomplete mission on the roadmap (or the tip from `map next`), if `clarify-status` is `open`, **stop**. Tell the human to finish clarification. Do not AFK.
3. Do not ask new non-blocking questions. Record assumptions in `decisions.md`.

## AFK loop

Repeat until stop:

```
spacecraft map next <roadmap-id>
```

### Stop conditions

- Output `All missions complete.` → print: **AFK done. Human check, then /sc-ship.** Stop.
- Tip has `state=blocked` → print blocked id; stop for human.
- Hard blocker (missing secret, impossible acceptance) → set clarify open or blocked; stop.

### Per incomplete mission

1. Parse id from `map next` line (`M…:` prefix).
2. `spacecraft use <id>`
3. **One branch per mission** (never stack many missions on one branch):
   - If not already on `feat/<id>/…`, create/checkout `feat/<id>/<slug>` from main and `spacecraft bind-branch <id>`
4. Ensure artifacts:
   - If no `spec.md` body: write minimal spec from roadmap mission description + decisions
   - If no plan / empty tasks: Task(`sc-planner`) → write `plan.json`; `spacecraft set-state planned` then `in_progress`
5. Build: for each pending plan task, Task(`sc-tester`) then Task(`sc-coder`) then evidence via `spacecraft evidence`; mark task done. Use sc-tdd / sc-verification / sc-git as detail skills (auto). Delegate product code - do not write it in Commander.
6. Review: Task(`sc-reviewer`); write `review.md` / `review.json` with `status: ready` and releaseReadiness objects (`changelog`/`specNote` status ready when preparing for later ship). `spacecraft validate --strict`. `spacecraft set-state ready`.
7. Continue loop (next mission).

## Rules

- **Must**: Only user-facing slash skills for lifecycle are `/sc-run` and `/sc-ship`.
- **Must not**: Call `/sc-ship`, merge, push, or tag.
- **Must not**: Mid-loop clarify unless hard blocker.
- **Must**: One feature branch per mission id.
- **Must**: Prefer assumptions in `decisions.md` over asking the human during AFK.

## Legacy

Old slash skills (`/sc-start`, `/sc-plan`, `/sc-build`, …) live under `.deleted/skills/` and are not Cursor-discovered. Restore only if needed.
