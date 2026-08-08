# Mission brief

`/sc-discuss` exit gate: after blocking Q&A and (when visual) draft approval, **before** `clarify-status clear`.

Present what the human **must know** about the mission as four bullet blocks: **Goal**, **Will do**, **Impact**, **Extra**. They **Accept**, **Adjust**, or **Reject**. Not a quiz - do not make the human guess; every stake lands as a plain bullet.

Brief closes skim-and-run holes; it does **not** replace sc-clarify. Empty Spec Mirror slots or soft Verify → return to clarify, do not invent completeness in the brief.

## Procedure

1. **Spec Mirror** - From `spec.md`, fill five slots (one short phrase each). Any empty or soft slot that blocks safe AFK → stop; ask via sc-clarify; do not pose the brief yet.

   | Goal | Output | Good vs Bad | Verify | Out of scope |

2. **Stake coverage map** - Tick only slots with skim-and-run risk (wrong Goal, missed Verify, silent scope creep, unowned tradeoff). Prefer fewer bullets under Extra; never invent filler to hit a count.

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

3. **Present** - Chat shows **only** the shape below. Fill all four sections from Spec Mirror + ticked stakes. Prefer short bullets; dump nothing material.

   ```
   ### Mission brief

   **Goal**
   - …

   **Will do**
   - …

   **Impact**
   - …

   **Extra**
   - …

   **Decision:** Accept | Adjust | Reject
   ```

   - **Goal** - what success looks like (from Spec Mirror Goal + Output + Good vs Bad when they bind the outcome). One or more bullets; not an essay.
   - **Will do** - concrete work AFK will perform (surfaces, files/APIs when they matter, sizing seam when `phases` / `roadmap`).
   - **Impact** - what changes for the product/process after ship; who or what feels it; what stays the same when that closes a skim risk.
   - **Extra** - everything else the human must own before AFK: Verify bar, out of scope, tradeoffs, Wrong-if (pre-mortem), human-owned policy/look/sizing, decision-digest lines from `decisions.md` that bind AFK. Put each ticked stake here as its own bullet. Never hide a stake to keep the section short.
   - **Pre-mortem** - when the mission is not trivial (new behavior, soft Verify risk, sizing split, or multi-concern), include **at least one** Wrong-if bullet under **Extra**. Skip pre-mortem only when every stake is machine-checkable and uncontested.
   - **Decision digest** - prefer Extra bullets drawn from 2–3 behavior-changing lines in `decisions.md` when they bind AFK; never harness trivia.

4. **Decide** - Wait for Accept / Adjust / Reject. Bad bullet (ambiguous or harness trivia) → rewrite or drop; do not punish the human for agent noise.

5. **Optional teach-back** - After Accept, may ask once: summarize this mission in one sentence in your own words. Mismatch → treat as Adjust (`Mission brief: adjust - teach-back mismatch`); update spec/decisions; re-brief. Do not block Accept on teach-back unless the human opts in or a prior Accept looked hollow.

6. **Record**
   | Human | Agent | `decisions.md` |
   |-------|-------|----------------|
   | **Accept** | Clear only if other exit gates hold | `Mission brief: accepted` |
   | **Adjust** | Do not clear; update spec/decisions; re-brief | `Mission brief: adjust - <summary>` |
   | **Reject** | Do not clear; stop AFK for this mission | `Mission brief: rejected - <reason>` |
   | Skip (step 2) | Proceed without brief body | `Mission brief: skipped - no material gaps` |
   | Human skip | Record reason; do not invent sections | `Mission brief: skipped - <reason>` |

Never clear while a posed brief awaits a decision (unless skip recorded).

## Bullet craft

| Do | Don't |
|----|-------|
| Spec / requirement / process / result / tradeoff / out-of-scope / approved look / Wrong-if | Harness trivia (`clarify-status`, slash skills, branches, evidence labels) |
| Stake the human owns before AFK | Filler bullets to "have a brief" |
| Plain bullets under Goal / Will do / Impact / Extra | Hide stakes; quiz the human first |
| Plain chat labels: Goal / Will do / Impact / Extra | Brand labels in chat (`Feynman`, `5W1H`, `Quiz`, `Stake map`, `I/Q/A`) |
| Extra bullets only for ticked stake-map slots (plus Verify / out-of-scope when they bind) | Extra Q&A or I/Q/A cards inside the brief |

### Bullet style (required two beats when a stake needs both)

1. **Plain** - one short clause a smart friend would get
2. **Technical** - exact Verify bar, limits, out-of-scope, file/API names when they matter

Same bullet may carry both beats after a hyphen, or split into two bullets. Do not turn briefs into speeches. Fuller narrative rewrite → sc-writer `prose-rhythm.md`.

Harness missions: frame as **product behavior**, never "which markdown line".

### Pre-mortem bullet (Wrong-if)

When required by Procedure step 3, under **Extra**:

```
- Wrong-if: <plain failure mode>. Detect: <symptom + how we would catch it>.
```

### Hollow signals (rewrite or return to clarify)

- Goal / Will do / Impact empty or vague while Extra is long
- Extra missing Verify, limits, or out-of-scope when those bind AFK
- No Wrong-if when pre-mortem is required
- Bullet is harness trivia
- Spec Mirror still has an empty/soft slot but brief is posed anyway

## Example

Mission: hard ≤7 tasks per plan phase.

```
### Mission brief

**Goal**
- Plans stay at most seven jigsaw tasks per phase - a hard wall, not a vibe
- Over the wall means phases or a vertical multi-mission split, not soft exceptions

**Will do**
- Enforce the ≤7 task cap in planning guidance and sizing protocol
- Reject soft "prefer ≤7" and 8-9 exception bands

**Impact**
- Oversized single-phase plans stop shipping
- Discuss must size `phases` (`plan-phaseN.json`) or roadmap seams before AFK continues

**Extra**
- Verify: count jigsaw tasks in `plan.json` / `plan-phaseN.json` - more than seven without recorded `Sizing: phases|roadmap` is wrong
- Out of scope: cross-feature layer waterfalls and `*-ux` roadmap seams
- Wrong-if: AFK treats the cap as a vibe; oversized plans ship and discuss sizing never runs. Detect by task count without a recorded split
- Sizing when multi-seam: `*-data` → `*-functional` → `*-ui` as needed; UX draft stays inside the ui mission discuss

**Decision:** Accept | Adjust | Reject
```

### Optional sizing bullets (when split or multi-concern)

When `Sizing:` is `phases` or `roadmap`, put Human-owned / Tradeoff stakes under **Extra** (and mirror the chosen path under **Will do** when it binds the work):

```
**Will do**
- Size as <single | phases | roadmap <id>> before AFK

**Extra**
- Tradeoff: this requirement touches more than one seam (data / functional / ui); one mission may choke the ≤7 wall or block data work behind UI draft HIL
- Human-owned: UX draft stays inside the ui mission discuss - not a fourth seam
```

## Must not

- Mid-clarify or `/sc-run`
- Invent sections when the skip above applies
- Clear while a posed brief is undecided (unless skip recorded)
- Present stakes only after the human "takes a quiz"
- Add blocking questions inside the brief (return to sc-clarify)
- Pose a brief while Spec Mirror still has an empty/soft AFK-blocking slot
- Force a fixed bullet count or brand-label the stake map in chat
- Reintroduce Information / Question / Answer cards in the posed brief
