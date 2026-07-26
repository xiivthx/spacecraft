# Mission brief

`/sc-discuss` exit gate: after blocking Q&A and (when visual) draft approval, **before** `clarify-status clear`.

Present what the human **must know** about the mission as Information → Question → Answer cards (Feynman + technical). They **Accept**, **Adjust**, or **Reject**. Not a quiz - do not make the human guess; the Answer is in the brief.

## Procedure

1. **Probe** - From `spec.md` + product lines in `decisions.md`, list risks of skim-and-run (wrong Goal, missed Verify, silent scope creep). If none material → record skip and stop:
   `Mission brief: skipped - no material gaps`
2. **Present** - 1–5 cards (prefer fewer). Chat shows **only**:
   ```
   ### Mission brief

   1.
   **Information:** …
   **Question:** …
   **Answer:** …

   2. …

   **Decision:** Accept | Adjust | Reject
   ```
3. **Decide** - Wait for Accept / Adjust / Reject. Bad card (ambiguous or harness trivia) → rewrite or drop; do not punish the human for agent noise.
4. **Record**
   | Human | Agent | `decisions.md` |
   |-------|-------|----------------|
   | **Accept** | Clear only if other exit gates hold | `Mission brief: accepted` |
   | **Adjust** | Do not clear; update spec/decisions; re-brief | `Mission brief: adjust - <summary>` |
   | **Reject** | Do not clear; stop AFK for this mission | `Mission brief: rejected - <reason>` |
   | Skip (step 1) | Proceed without cards | `Mission brief: skipped - no material gaps` |
   | Human skip | Record reason; do not invent cards | `Mission brief: skipped - <reason>` |

Never clear while a posed brief awaits a decision (unless skip recorded).

## Card craft

| Do | Don't |
|----|-------|
| Spec / requirement / process / result / tradeoff / out-of-scope / approved look | Harness trivia (`clarify-status`, slash skills, branches, evidence labels) |
| Stake the human owns before AFK | Filler cards to "have a brief" |
| Answer in the brief (Feynman + technical) | Hide the answer; quiz the human first |
| Plain chat labels: Information / Question / Answer | Brand labels in chat (`Feynman`, `5W1H`, `Quiz`) |

### Answer style (required two beats)

1. **Feynman** - one plain sentence a smart friend would get
2. **Technical** - exact Verify bar, limits, out-of-scope, file/API names when they matter

Harness missions: frame as **product behavior**, never "which markdown line".

## Example

Mission: hard ≤7 tasks per plan phase.

```
### Mission brief

1.
**Information:** Plans blow up when "a few more tasks" sneak in after the hard cap.
**Question:** What is the hard task limit per plan phase after this ships?
**Answer:** At most seven tasks - a wall, not a vibe. Soft "prefer ≤7" and 8-9 exception bands are rejected; split with plan-phaseN or a vertical roadmap of feature seams (`*-data` → `*-functional` → `*-ui`) per mission-sizing - never a cross-feature layer waterfall or a `*-ux` seam.

**Decision:** Accept | Adjust | Reject
```

### Optional sizing card (when split or multi-concern)

When `Sizing:` is `phases` or `roadmap`, include a brief card so the human owns the split:

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
