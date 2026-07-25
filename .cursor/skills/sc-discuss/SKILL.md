---
name: sc-discuss
description: "Pre-build HIL: clarify, brainstorm, decide, and approve visual draft before implement. Invoke as /sc-discuss."
disable-model-invocation: true
---

# sc-discuss

## Goal

Find what we want before implement: clear `spec.md`, decisions, answered questions, and (when visual) an approved draft HTML. Exit with `clarify-status clear` so a new session can `/sc-run`.

## Output

Mission ready to build: solid `spec.md`; `questions.md` / `decisions.md` updated; visual: `UI draft approved: <file>`; mission brief accepted or skip recorded; `spacecraft clarify-status clear`. Handoff: **Spec clear. New session: /sc-run.** Never plan AFK, RED-GREEN, product code, or ship.

## Good / Bad

- Good: sharp Goal / Output / Good-Bad / Verify; one blocking question at a time; soft gaps → `decisions.md`; visual brief + draft with designer gate before human; mission brief (I/Q/A) then Accept/Adjust/Reject before clear
- Bad: implementing; writing `plan.json` AFK; shipping; serving unreviewed draft HTML; clearing while draft unapproved or mission brief undecided; quizzing the human instead of presenting Answers

## Verify

Human confirms spec (and draft when visual), then Accepts (or skip) the mission brief:

```
spacecraft clarify-status clear
# visual: decisions.md contains "UI draft approved: <draft-file>"
# decisions.md contains "Mission brief: accepted" OR "Mission brief: skipped - <reason>"
```

## Arguments

```
/sc-discuss
/sc-discuss <mission-id|title>
```

## Lifecycle

Canonical: `.cursor/rules/200-workflow.mdc` - this skill is discuss HIL only. Next: new session `/sc-run`.

## Pre-flight

1. Resolve: `$ARGUMENTS` → `spacecraft use`; else `spacecraft resolve`. `spacecraft new` only if user wants mutating work and none exists.
2. Read `spec.md`, `questions.md`, `decisions.md`, drafts when present.
3. Do not start `/sc-run` build, product code, or ship.

## Discuss loop

```
resolve → inspect → classify gaps → talk / ask / decide → (visual: brief + draft → designer → fix → human HIL) → mission brief → clear → handoff
```

### Spec and decisions

1. Ensure `spec.md` has Goal, Output, Good vs Bad, Verify (machine-checkable where possible). Skim `.space/trust/lessons.md` before inventing process (sc-learn: seed if missing).
2. Blocking ambiguity: sc-clarify (one question at a time); record in `questions.md` / `decisions.md`.
3. Soft gaps → assumptions in `decisions.md` (do not block clear alone).
4. Deep architecture: optional Task(`sc-adviser`). Keep `clarify-status open` while blockers or unapproved visual draft remain.

### Visual design (when UI/FE)

Detect from intent / `spec.md`. If visual:

1. sc-ux-design design brief (6 dimensions); human approval.
2. **Pack selection before draft HTML:** `swiss-grid`, `editorial`, or `none - custom brief only`. Record in `decisions.md`. Human or explicit brief choice only - no silent auto-matcher.
3. Generate draft HTML under `.space/missions/<id>/design/drafts/` (not wireframe-only). Check 375px. Follow sc-ux-design prompt assembly.
4. **Designer gate before human:** Task(`sc-designer`); Commander applies critical/important fixes; re-check 375px. Do not present draft until this passes.
5. Serve via `serve-html.mjs`; iterate (draft → designer → fix → human) until approved (max 3 human rounds). Each new draft re-runs designer gate.
6. On approval: `UI draft approved: <draft-file>` in `decisions.md`.
7. Skip draft only for non-visual FE; record skip reason.

### Mission brief (before clear)

Follow `references/mission-brief.md`. Present Information / Question / Answer cards (Feynman + technical); human **Accept | Adjust | Reject**.

- Accept → record `Mission brief: accepted` (then clear if other gates hold)
- Adjust → record `Mission brief: adjust - <summary>`; update spec/decisions; re-brief; do not clear
- Reject → record `Mission brief: rejected - <reason>`; do not clear
- Skip → `Mission brief: skipped - <reason>`

Never clear while a posed brief awaits a decision (unless skip recorded).

### Exit

1. No open blocking questions; Verify present; visual approved or skip recorded; mission brief accepted or skip recorded.
2. `spacecraft clarify-status clear`.
3. Handoff: **Spec clear. New session: /sc-run.**

## Rules

- Never `/sc-run` build, `/sc-ship`, merge, push, tag, or product implementation/tests (draft HTML only).
- Never present draft before designer gate + critical/important fixes.
- Never clear while mission brief undecided (unless skip recorded).
- Prefer `spec.md` / `decisions.md` / `questions.md` over chat-only memory.
- One mission focus per discuss session.

## Specialist skills

| Concern | Where |
|---------|--------|
| Blocking questions | sc-clarify |
| Draft HTML / visual-verify | sc-ux-design |
| Draft critique | Task(`sc-designer`) |
| Mission brief | `references/mission-brief.md` |
| Architecture | Task(`sc-adviser`) |
| Plan / TDD / evidence | `/sc-run` |

## References

- sc-clarify, sc-ux-design, `/sc-run`, `/sc-ship`
- `references/mission-brief.md`
