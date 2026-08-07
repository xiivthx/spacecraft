# Mission brief

`/sc-discuss` exit gate: after blocking Q&A and (when visual) draft approval, **before** `clarify-status clear`.

Present what the human **must know** about the mission as Information → Question → Answer cards (Feynman + technical). They **Accept**, **Adjust**, or **Reject**. Not a quiz - do not make the human guess; the Answer is in the brief.

Brief closes skim-and-run holes; it does **not** replace sc-clarify. Empty Spec Mirror slots or soft Verify → return to clarify, do not invent completeness in the brief.

## Procedure

1. **Spec Mirror** - From `spec.md`, fill five slots (one short phrase each). Any empty or soft slot that blocks safe AFK → stop; ask via sc-clarify; do not pose the brief yet.

   | Goal | Output | Good vs Bad | Verify | Out of scope |

2. **Stake coverage map** - Tick only slots with skim-and-run risk (wrong Goal, missed Verify, silent scope creep, unowned tradeoff). Prefer fewer cards; never invent filler to hit a count.

   | Slot | Risk if skimmed |
   |------|-----------------|
   | Goal | AFK builds the wrong change |
   | Done / Verify | "Done" is uncheckable or taste-only |
   | In vs Out | Scope creeps or cuts a must-have |
   | Tradeoff | Human did not own the fork |
   | Wrong-if | Failure mode surfaces after merge |
   | Human-owned | Sizing, look, or policy left implicit |

   If **no** material risks after Spec Mirror is solid → record skip and stop:
   `Mission brief: skipped - no material gaps`

3. **Present** - Chat shows **only** the shape below. Lead with one-breath. Emit 1–5 I/Q/A cards for ticked stakes only (prefer fewer).

   ```
   ### Mission brief

   **In one breath:** …

   1.
   **Information:** …
   **Question:** …
   **Answer:** …

   2. …

   **Decision:** Accept | Adjust | Reject
   ```

   - **In one breath** - one sentence: what this mission changes and when it is done. Not an essay; not a substitute for cards.
   - **Pre-mortem** - when the mission is not trivial (new behavior, soft Verify risk, sizing split, or multi-concern), include **at least one** Wrong-if card (see Card craft). Skip pre-mortem only when every stake is machine-checkable and uncontested.
   - **Decision digest** - prefer Information drawn from 2–3 behavior-changing lines in `decisions.md` when they bind AFK; never harness trivia.

4. **Decide** - Wait for Accept / Adjust / Reject. Bad card (ambiguous or harness trivia) → rewrite or drop; do not punish the human for agent noise.

5. **Optional teach-back** - After Accept, may ask once: summarize this mission in one sentence in your own words. Mismatch → treat as Adjust (`Mission brief: adjust - teach-back mismatch`); update spec/decisions; re-brief. Do not block Accept on teach-back unless the human opts in or a prior Accept looked hollow.

6. **Record**
   | Human | Agent | `decisions.md` |
   |-------|-------|----------------|
   | **Accept** | Clear only if other exit gates hold | `Mission brief: accepted` |
   | **Adjust** | Do not clear; update spec/decisions; re-brief | `Mission brief: adjust - <summary>` |
   | **Reject** | Do not clear; stop AFK for this mission | `Mission brief: rejected - <reason>` |
   | Skip (step 2) | Proceed without cards | `Mission brief: skipped - no material gaps` |
   | Human skip | Record reason; do not invent cards | `Mission brief: skipped - <reason>` |

Never clear while a posed brief awaits a decision (unless skip recorded).

## Card craft

| Do | Don't |
|----|-------|
| Spec / requirement / process / result / tradeoff / out-of-scope / approved look / Wrong-if | Harness trivia (`clarify-status`, slash skills, branches, evidence labels) |
| Stake the human owns before AFK | Filler cards to "have a brief" |
| Answer in the brief (Feynman + technical) | Hide the answer; quiz the human first |
| Plain chat labels: Information / Question / Answer | Brand labels in chat (`Feynman`, `5W1H`, `Quiz`, `Stake map`) |
| Cards only for ticked stake-map slots | Extra Q&A inside the brief |

### Answer style (required two beats)

1. **Feynman** - one plain sentence a smart friend would get; prefer natural rhythm (not a stack of same-length fragments) - still one sentence, not an essay
2. **Technical** - exact Verify bar, limits, out-of-scope, file/API names when they matter

Do not turn briefs into speeches. Fuller narrative rewrite → sc-writer `prose-rhythm.md`.

Harness missions: frame as **product behavior**, never "which markdown line".

### Pre-mortem card (Wrong-if)

When required by Procedure step 3:

```
**Information:** If AFK ships the wrong reading, the failure shows up after merge.
**Question:** What is the most important failure mode to own before AFK?
**Answer:** <Feynman one sentence>. <Technical: symptom + how we would detect it>.
```

### Hollow signals (rewrite or return to clarify)

- Answer has Feynman only - missing Verify, limits, or out-of-scope
- No Wrong-if when pre-mortem is required
- Card is harness trivia
- Spec Mirror still has an empty/soft slot but brief is posed anyway

## Example

Mission: hard ≤7 tasks per plan phase.

```
### Mission brief

**In one breath:** Plans may not grow past seven tasks per phase; overage forces phases or a vertical multi-mission split, not soft exceptions.

1.
**Information:** Plans blow up when "a few more tasks" sneak in after the hard cap.
**Question:** What is the hard task limit per plan phase after this ships?
**Answer:** At most seven tasks - a wall, not a vibe. Soft "prefer ≤7" and 8-9 exception bands are rejected; discuss sizes `phases` (`plan-phaseN.json`) or a vertical multi-mission of feature seams (`*-data` → `*-functional` → `*-ui`) per mission-sizing - never a cross-feature layer waterfall or a `*-ux` seam.

2.
**Information:** If AFK treats the cap as a vibe, oversized plans ship and discuss sizing never runs.
**Question:** What is the most important failure mode to own before AFK?
**Answer:** Soft caps and silent 8-9 "just this once" bands. Detect by counting jigsaw tasks in `plan.json` / `plan-phaseN.json` - more than seven without a recorded `Sizing: phases|roadmap` split is wrong.

**Decision:** Accept | Adjust | Reject
```

### Optional sizing card (when split or multi-concern)

When `Sizing:` is `phases` or `roadmap`, include a brief card so the human owns the split (Human-owned / Tradeoff stake):

```
**Information:** This requirement touches more than one seam (data / functional / ui), so one mission may choke the ≤7 task wall or block data work behind UI draft HIL.
**Question:** How are we sizing this before AFK?
**Answer:** <single | phases | roadmap <id> with <feature>-data → <feature>-functional → <feature>-ui as needed>. UX draft stays inside the ui mission discuss - not a fourth seam.
```

## Must not

- Mid-clarify or `/sc-run`
- Invent cards when the skip above applies
- Clear while a posed brief is undecided (unless skip recorded)
- Present Answer only after the human "takes a quiz"
- Add blocking questions inside the brief (return to sc-clarify)
- Pose a brief while Spec Mirror still has an empty/soft AFK-blocking slot
- Force five cards or brand-label the stake map in chat
