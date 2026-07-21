---
name: sc-run
description: "AFK roadmap runner: after /sc-discuss clear, auto-run incomplete missions to ready. Invoke as /sc-run <roadmap-id>."
disable-model-invocation: true
---

# sc-run

## Goal

After `/sc-discuss` (clarify clear + approved visual draft when required), AFK incomplete roadmap missions to `ready` so a human can check and then `/sc-ship`.

## Output

Missions at `state=ready` (or stop on blocked / clarify / missing draft approval). Handoff: **AFK done. Human check, then /sc-ship.** Never merge/push/tag.

## Good / Bad

- Good: discuss clear first; visual draft already approved in discuss; jigsaw plan tasks; per-acceptance RED→GREEN via Task; auto checkpoint commits; visual + functional recheck for UI; evidence real; review recorded; `sc-judge` before ready; block `set-state ready` on `REFUTED`
- Bad: shipping; clarifying or iterating draft HTML in this session; implementing UI without approved draft record; mid-loop non-blocking questions; stacking many missions on one branch; bulk implement without RED; skipping checkpoint commits; claiming UI ready without screenshots/visual + functional evidence; setting ready without `sc-judge` or after `REFUTED`

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
| Discuss | Human (`/sc-discuss`) | Spec, decisions, visual draft approval; `clarify-status clear` |
| Run | AI | Loop missions to `ready` (RED-GREEN + checkpoints) |
| Check | Human | Review ready work |
| Ship | Human | Explicit `/sc-ship` only (squashes checkpoints → ≤5 commits) |

## Pre-flight

1. Resolve roadmap: `$ARGUMENTS` → `map use`; else `map current` (fail if missing).
2. If any incomplete mission has `clarify-status` `open` → stop; recommend `/sc-discuss`. Do not AFK.
3. If any incomplete mission is visual UI/FE and `decisions.md` lacks `UI draft approved: …` (and no recorded non-visual skip) → stop; recommend `/sc-discuss`.
4. Do not ask non-blocking questions; record assumptions in `decisions.md`. Do not run design-brief / draft-HTML discovery here.

## AFK loop

```
spacecraft map next <roadmap-id>
```

Stop when: `All missions complete.` (print handoff), tip `blocked`, clarify opened by hard blocker, or missing draft approval discovered.

### Per incomplete mission

1. Parse id from `map next` (`M…:` prefix); `spacecraft use <id>`.
2. One branch per mission: if not on `feat/<id>/…`, checkout from main and `spacecraft bind-branch <id>`.
3. Artifacts: `spec.md` must already be discuss-ready; Task(`sc-planner`) → jigsaw `plan.json`; `set-state planned` then `in_progress`.
4. **Build (atomic RED-GREEN)** - for each pending plan task, for each `acceptance[]` check:
   1. Triage via sc-tdd (skip tautologies; record skip in task notes / `decisions.md`).
   2. **RED**: Task(`sc-tester`) for that single acceptance. Auto checkpoint commit (`test: …`).
   3. **GREEN**: Task(`sc-coder`) / Task(`sc-firmware`) minimum code. `spacecraft evidence`. Auto checkpoint commit (`feat:` / `fix:`).
   4. Mark task `done` only after all its acceptances pass.
5. **Combine**: after all plan tasks done - refactor for cohesion; run unit + integration/functional suite; `spacecraft evidence` for the full gate. Auto checkpoint commit (`refactor:` / `test:`).
6. **UI visual + functional recheck (when visual UI/FE)** - before review:
   1. Visual: sc-ux-design Tier 3 (`visual-verify.mjs`) when Playwright available; else browser screenshots. Capture screenshot paths in evidence / `decisions.md`. Cross-check against approved draft / design brief.
   2. Functional: Vitest/RTL or project functional suite via `spacecraft evidence`.
   3. Fix issues found; do not set `ready` without both.
7. Review + `sc-judge` ready gate:
   1. Task(`sc-reviewer`); for UI also Task(`sc-designer`) when visual.
   2. Run `sc-judge` (`.cursor/skills/sc-judge/SKILL.md`) - adversarial prove before ready.
   3. Capture judge evidence via `spacecraft evidence` with a label including `judge` (e.g. `judge-<mission-id>`).
   4. If verdict is `REFUTED`: do **not** `set-state ready`; leave blocked / fix and re-judge. Handshake blocked.
   5. Only when judge is not `REFUTED`: write `review.md` / `review.json` (`status: ready`, releaseReadiness ready); `validate --strict`; `set-state ready`.
8. Continue loop. Do not squash here - `/sc-ship` squashes checkpoints to ≤5 Conventional Commits.

## Checkpoint commits (mandatory during AFK)

Commander auto-commits on the work branch after every RED, every GREEN, and after the combine/refactor gate. Use Conventional Commit subjects; body may note `wip checkpoint` + mission id + acceptance. See sc-git §Checkpoint commits. Never push.

## Rules

- Never call `/sc-ship`, merge, push, or tag.
- Never run `/sc-discuss` work (clarify Q&A, design brief, draft HTML iteration) inside AFK - stop and hand off to `/sc-discuss` if discovery is still needed.
- One feature branch per mission id.
- Prefer `decisions.md` assumptions over mid-AFK questions (hard blockers only: missing secret, impossible acceptance).
- Delegate product code/tests via Task - Commander orchestrates only.
- Do not start AFK while clarify is `open` or visual draft approval is missing.
- Do not implement visual UI/FE without an approved draft HTML record (from discuss) or recorded non-visual skip.
- Do not `set-state ready` without `sc-judge`, or when verdict is `REFUTED`. Capture judge evidence (label including `judge`).

## References

- `/sc-discuss` - clarify, decisions, visual draft approval
- `/sc-ship` - explicit ship only
- sc-ux-design - post-build visual QC (brief/draft owned by discuss)
- sc-judge - adversarial prove gate before ready; block on `REFUTED`
