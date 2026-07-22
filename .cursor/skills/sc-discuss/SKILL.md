---
name: sc-discuss
description: "Pre-build HIL: clarify, brainstorm, decide, and approve visual draft before implement. Invoke as /sc-discuss."
disable-model-invocation: true
---

# sc-discuss

## Goal

Find what we want before implement: clear `spec.md`, recorded decisions, answered questions, and (when visual) an approved draft HTML. Exit with `clarify-status clear` so a new session can `/sc-run`.

## Output

Mission ready to build: `spec.md` solid; `questions.md` / `decisions.md` updated; visual missions have `UI draft approved: <file>` in `decisions.md`; `spacecraft clarify-status clear`. Handoff: **Spec clear. New session: /sc-run.** Never plan AFK, RED-GREEN, product code, or ship.

## Good / Bad

- Good: talk until Goal / Output / Good-Bad / Verify are sharp; one blocking question at a time (sc-clarify protocol); soft gaps → `decisions.md`; visual brief + draft HTML; Task(`sc-designer`) then Commander fixes before human sees draft; iterate until human likes it; new session for `/sc-run`
- Bad: implementing; writing `plan.json` AFK; checkpoint commits; mid-discuss shipping; showing raw/unreviewed draft HTML to the human; clearing clarify while draft unapproved on visual work; stacking roadmap AFK in this session

## Verify

Human confirms the spec (and draft when visual). Then:

```
spacecraft clarify-status clear
# visual: decisions.md contains "UI draft approved: <draft-file>"
```

## Arguments

`$ARGUMENTS` = mission selector or roadmap mission id (optional if already resolved).

```
/sc-discuss
/sc-discuss <mission-id|title>
```

## HIL vs AFK

| Phase | Who | Action |
|-------|-----|--------|
| Discuss | Human + AI | Spec, brainstorm, decisions, questions, visual draft |
| Run | AI (`/sc-run`) | AFK plan → RED-GREEN → review → ready |
| Check | Human | Review ready work |
| Ship | Human | Explicit `/sc-ship` only |

## Pre-flight

1. Resolve mission: `$ARGUMENTS` → `spacecraft use`; else `spacecraft resolve`. Create with `spacecraft new` only if user wants mutating work and none exists.
2. Read `spec.md`, `questions.md`, `decisions.md`, design drafts when present.
3. Do not start `/sc-run` build, product code, or ship from this skill.

## Discuss loop

```
resolve → inspect → classify gaps → talk / ask / decide → (visual: brief + draft → designer → fix → human HIL) → clear → handoff
```

### Spec and decisions

1. Ensure `spec.md` has Goal, Output, Good vs Bad, Verify (machine-checkable where possible). Prefer editing the spec over chat-only agreements.
2. Use sc-clarify protocol for blocking ambiguity: exhaust files/research first; ask exactly one blocking question at a time (question + why + recommendation + if-accepted); record in `questions.md` / `decisions.md`.
3. Soft / non-blocking gaps: write explicit assumptions to `decisions.md` (do not block clear on soft gaps alone).
4. Deep architecture tradeoffs: optional Task(`sc-adviser`). Visual draft critique: required Task(`sc-designer`) before human HIL (see Visual design).
5. Keep `spacecraft clarify-status open` while blocking questions or unapproved visual draft remain.

### Visual design (when UI/FE surface)

Detect from intent / `spec.md` (layout, style, pages, components, design). If visual:

1. Follow sc-ux-design design brief (6 dimensions); get human approval.
2. Generate standalone draft HTML under `.space/missions/<id>/design/drafts/` (layout, style tokens, key components - not wireframe-only). Check 375px.
3. **Designer gate (required before human sees draft):** Task(`sc-designer`) on the draft. Commander applies all critical and important fixes to the draft HTML (`sc-designer` is readonly). Re-check 375px after fixes. Do not serve or present the draft to the human until this gate passes.
4. Serve via `serve-html.mjs` and present the cleaned draft for human review.
5. Iterate (draft → designer → fix → human) until the human likes it (max 3 human rounds, then escalate for direction). Each new draft version re-runs the designer gate before human HIL.
6. On approval: record `UI draft approved: <draft-file>` in `decisions.md`. Optionally note art direction / DESIGN.md updates.
7. Skip draft only for non-visual FE (pure logic/hooks); record skip reason in `decisions.md`.

Non-visual missions skip the draft path.

### Exit

1. No open blocking questions.
2. Spec is implementable (Verify present).
3. Visual: approved draft recorded (or non-visual skip recorded).
4. `spacecraft clarify-status clear`.
5. Recommend **new session** `/sc-run` (roadmap id if applicable). Do not continue into AFK build in this chat unless the user explicitly overrides.

## Rules

- Never call `/sc-run` build steps, `/sc-ship`, merge, push, or tag from this skill.
- Never write product implementation or tests; drafts are throwaway HTML only.
- Never serve or present visual draft HTML to the human before Task(`sc-designer`) and Commander fixes for critical/important findings.
- Prefer recording over memory: `spec.md`, `decisions.md`, `questions.md` are the handoff.
- Ask/clarify/pathfinder-style work lives here; do not invent separate slash commands for them.
- One mission focus per discuss session (roadmap selection via `map use` is fine; AFK loop is `/sc-run`).

## Specialist skills (not slash phases)

| Concern | Where |
|---------|--------|
| One-question clarify protocol | sc-clarify (used inside discuss) |
| Draft HTML + anti-slop / visual-verify scripts | sc-ux-design |
| Required draft critique before human HIL | Task(`sc-designer`) then Commander fixes |
| Deep architecture advice | Task(`sc-adviser`) |
| Plan / TDD / evidence | `/sc-run` + sc-planning / sc-tdd / sc-verification |

## References

- sc-clarify - blocking question protocol
- sc-ux-design - brief, draft HTML, post-build visual QC
- `/sc-run` - AFK implement after clarify clear
- `/sc-ship` - explicit ship only
