---
name: sc-run
description: "AFK roadmap runner: after human clarify, auto-run incomplete missions to ready. Invoke as /sc-run <roadmap-id>."
disable-model-invocation: true
---

# sc-run

## Goal

After human clarify, AFK incomplete roadmap missions to `ready` so a human can check and then `/sc-ship`.

## Output

Missions at `state=ready` (or stop on blocked / clarify / UI draft pending). Handoff: **AFK done. Human check, then /sc-ship.** Never merge/push/tag.

## Good / Bad

- Good: clarify clear first; UI draft HTML approved before visual FE build; jigsaw plan tasks; per-acceptance RED→GREEN via Task; auto checkpoint commits; visual + functional recheck for UI; evidence real; review recorded
- Bad: shipping; implementing UI without approved draft HTML; mid-loop non-blocking questions; stacking many missions on one branch; bulk implement without RED; skipping checkpoint commits; claiming UI ready without screenshots/visual + functional evidence

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
| UI draft | Human | Approve standalone draft HTML (layout/style/components) when mission is visual UI/FE |
| Run | AI | Loop missions to `ready` (RED-GREEN + checkpoints) after draft approval when required |
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

Stop when: `All missions complete.` (print handoff), tip `blocked`, UI draft awaiting approval, or hard blocker (set clarify open / blocked).

### Per incomplete mission

1. Parse id from `map next` (`M…:` prefix); `spacecraft use <id>`.
2. One branch per mission: if not on `feat/<id>/…`, checkout from main and `spacecraft bind-branch <id>`.
3. Artifacts: minimal `spec.md` if empty; else Task(`sc-planner`) → jigsaw `plan.json`; `set-state planned` then `in_progress`.
4. **UI draft HIL (when visual UI/FE)** - detect from `spec.md` / `plan.json` (layout, style, components, pages, design). If visual work:
   1. Follow sc-ux-design: design brief → standalone draft HTML under `.space/missions/<id>/design/drafts/` covering **layout, style tokens, and key components** (not wireframe-only).
   2. Serve draft (`serve-html.mjs`); check 375px; hand off for human approval.
   3. **Stop AFK.** Set `spacecraft clarify-status open` (or leave clear only if already approved). Record pending approval in `decisions.md`.
   4. Do **not** start RED-GREEN for visual tasks until approval is recorded in `decisions.md` (e.g. `UI draft approved: <draft-file>`) or clarify cleared after explicit user approve.
   5. Skip draft only for non-visual FE (pure logic/hooks, no UI surface); record skip reason in `decisions.md`.
5. **Build (atomic RED-GREEN)** - for each pending plan task, for each `acceptance[]` check:
   1. Triage via sc-tdd (skip tautologies; record skip in task notes / `decisions.md`).
   2. **RED**: Task(`sc-tester`) for that single acceptance. Auto checkpoint commit (`test: …`).
   3. **GREEN**: Task(`sc-coder`) / Task(`sc-firmware`) minimum code. `spacecraft evidence`. Auto checkpoint commit (`feat:` / `fix:`).
   4. Mark task `done` only after all its acceptances pass.
6. **Combine**: after all plan tasks done - refactor for cohesion; run unit + integration/functional suite; `spacecraft evidence` for the full gate. Auto checkpoint commit (`refactor:` / `test:`).
7. **UI visual + functional recheck (when visual UI/FE)** - before review:
   1. Visual: sc-ux-design Tier 3 (`visual-verify.mjs`) when Playwright available; else browser screenshots. Capture screenshot paths in evidence / `decisions.md`.
   2. Functional: Vitest/RTL or project functional suite via `spacecraft evidence`.
   3. Fix issues found; do not set `ready` without both.
8. Review: Task(`sc-reviewer`); for UI also Task(`sc-designer`) when visual; write `review.md` / `review.json` (`status: ready`, releaseReadiness ready); `validate --strict`; `set-state ready`.
9. Continue loop. Do not squash here - `/sc-ship` squashes checkpoints to ≤5 Conventional Commits.

## Checkpoint commits (mandatory during AFK)

Commander auto-commits on the work branch after every RED, every GREEN, and after the combine/refactor gate. Use Conventional Commit subjects; body may note `wip checkpoint` + mission id + acceptance. See sc-git §Checkpoint commits. Never push.

## Rules

- Never call `/sc-ship`, merge, push, or tag.
- One feature branch per mission id.
- Prefer `decisions.md` assumptions over mid-AFK questions (except UI draft approval - hard HIL stop).
- Delegate product code/tests via Task - Commander orchestrates only.
- Do not start AFK while clarify is `open`.
- Do not implement visual UI/FE without an approved draft HTML (or recorded non-visual skip).
