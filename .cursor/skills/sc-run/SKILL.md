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

- Good: clarify clear first; jigsaw plan tasks; per-acceptance RED→GREEN via Task; auto checkpoint commits; refactor+integration after slices; evidence real; review recorded
- Bad: shipping; mid-loop non-blocking questions; stacking many missions on one branch; bulk implement without RED; skipping checkpoint commits

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
| Run | AI | Loop missions to `ready` (RED-GREEN + checkpoints) |
| Check | Human | Review ready work |
| Ship | Human | Explicit `/sc-ship` only (squashes checkpoints → ≤5 commits) |

## Pre-flight

1. Resolve roadmap: `$ARGUMENTS` → `map use`; else `map current` (fail if missing).
2. If any incomplete mission has `clarify-status` `open` → stop; do not AFK. Requirements must be clear first.
3. Do not ask non-blocking questions; record assumptions in `decisions.md`.

## AFK loop

```
spacecraft map next <roadmap-id>
```

Stop when: `All missions complete.` (print handoff), tip `blocked`, or hard blocker (set clarify open / blocked).

### Per incomplete mission

1. Parse id from `map next` (`M…:` prefix); `spacecraft use <id>`.
2. One branch per mission: if not on `feat/<id>/…`, checkout from main and `spacecraft bind-branch <id>`.
3. Artifacts: minimal `spec.md` if empty; else Task(`sc-planner`) → jigsaw `plan.json`; `set-state planned` then `in_progress`.
4. **Build (atomic RED-GREEN)** - for each pending plan task, for each `acceptance[]` check:
   1. Triage via sc-tdd (skip tautologies; record skip in task notes / `decisions.md`).
   2. **RED**: Task(`sc-tester`) for that single acceptance. Auto checkpoint commit (`test: …`).
   3. **GREEN**: Task(`sc-coder`) / Task(`sc-firmware`) minimum code. `spacecraft evidence`. Auto checkpoint commit (`feat:` / `fix:`).
   4. Mark task `done` only after all its acceptances pass.
5. **Combine**: after all plan tasks done - refactor for cohesion; run unit + integration/functional suite; `spacecraft evidence` for the full gate. Auto checkpoint commit (`refactor:` / `test:`).
6. Review: Task(`sc-reviewer`); write `review.md` / `review.json` (`status: ready`, releaseReadiness ready); `validate --strict`; `set-state ready`.
7. Continue loop. Do not squash here - `/sc-ship` squashes checkpoints to ≤5 Conventional Commits.

## Checkpoint commits (mandatory during AFK)

Commander auto-commits on the work branch after every RED, every GREEN, and after the combine/refactor gate. Use Conventional Commit subjects; body may note `wip checkpoint` + mission id + acceptance. See sc-git §Checkpoint commits. Never push.

## Rules

- Never call `/sc-ship`, merge, push, or tag.
- One feature branch per mission id.
- Prefer `decisions.md` assumptions over mid-AFK questions.
- Delegate product code/tests via Task - Commander orchestrates only.
- Do not start AFK while clarify is `open`.
