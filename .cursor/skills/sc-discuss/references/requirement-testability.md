# Requirement testability pass

`/sc-discuss` soft gate (lens-style): stress-test `spec.md` for testability before deep plan/build. Decision job - not expertise cosplay.

## Goal

Make the requirement testable enough that `/sc-run` can turn Verify into acceptance checks without inventing bars.

## Output

Chat summary plus greppable artifact in `decisions.md`:

```
## Testability pass
- Testability: Testable | Not Testable
- Fixes needed: <concrete Verify/acceptance changes, or none>
- Notes: …
- Risks: …
- Test Ideas: …
- Requirement Bugs: …
- Question queue: <count> parked in questions.md (blocking N / non-blocking M / researchable K)
```

Or skip:

```
Testability pass skipped: <reason>
```

## Good / Bad

- Good: measurable outcomes; machine-checkable Verify when possible; question candidates classified and parked; one blocking ask via sc-clarify at a time; Test Ideas usable by sc-planner / sc-tester
- Bad: inventing Verify; dumping many questions in one user-facing turn; expertise cosplay; clearing while `Not Testable` and Verify still soft/missing

## Verify

`decisions.md` has `## Testability pass` OR `Testability pass skipped:`. If Testability is `Not Testable` and Verify is still soft/missing → keep `clarify-status open`. Greppable markers only - no CLI validator required.

## When required (any one)

- Soft or missing Verify
- New feature with behavioral uncertainty
- Human asks for requirement / testability review
- Mission brief probe finds skim-and-run risk on Verify

## When skip

- Verify already machine-checkable and uncontested
- Routine bug with clear repro
- Non-behavior docs-only work

Record skip in `decisions.md`:

```
Testability pass skipped: <reason>
```

## Procedure

1. **Read** - `spec.md` (Goal, Output, Good vs Bad, Verify). If the requirement is missing, empty, or too vague to score → respond exactly: `Requirement unclear - please provide more detail.` Keep clarify open. Do not invent Verify.

2. **Score testability**
   - **Testable** - outcomes are observable; Verify (or clear path to Verify) is measurable; acceptance can become RED-GREEN cycles
   - **Not Testable** - list specific fixes (clearer acceptance, measurable outcomes, edge cases, out-of-scope)

3. **Question queue** (do not dump as the user-facing turn)
   - List clarifying questions you would ask product
   - Classify each: **blocking** | **non-blocking** | **researchable**
   - Park under `questions.md` (Open)
   - Ask only the most blocking via sc-clarify (one at a time: Question / Why it matters / Recommendation / If accepted)
   - Non-blocking → assumption in `decisions.md`; researchable → read code / sc-search first

4. **Fill remaining sections**
   - **Notes** - observations that affect understanding
   - **Risks** - if unclear or incomplete, what breaks in build/ship
   - **Test Ideas** - high-level scenarios (happy path, edges, failure); keep them behavioral
   - **Requirement Bugs** - flaws, contradictions, ambiguities in the requirement itself

5. **Record** - write `## Testability pass` (or skip) to `decisions.md`. Tighten `spec.md` Verify when the human confirms fixes.

## Clear rule

- Do **not** clear while Testability is `Not Testable` **and** Verify is still soft/missing
- If the human later supplies machine-checkable Verify, re-score or record a new pass; `Not Testable` alone is not a permanent block once Verify is fixed

## Must / Must not

- **Must**: Exhaust `spec.md` / `questions.md` / `decisions.md` / repo before asking
- **Must**: Park question candidates; ask via sc-clarify one at a time
- **Must not**: Invent Verify from Test Ideas
- **Must not**: Expertise cosplay ("as a QA expert…")
- **Must not**: Replace sc-clarify or mission brief
