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
- Test Ideas:
  - Positive: …
  - Negative: …
  - Edge: …
  - Overlooked: …
- Implementation pitfalls: …
- Requirement Bugs: …
- Question queue: <count> parked in questions.md (blocking N / non-blocking M / researchable K)
```

Or skip:

```
Testability pass skipped: <reason>
```

## Good / Bad

- Good: measurable outcomes; machine-checkable Verify when possible; question candidates classified and parked; one blocking ask via sc-clarify at a time; structured Test Ideas (Positive/Negative/Edge/Overlooked) risk-driven when SFDIPOT/quality apply; Implementation pitfalls distinct from Requirement Bugs; usable by sc-planner / sc-tester
- Bad: inventing Verify; essay dumps instead of structured Test Ideas; dumping many questions in one user-facing turn; expertise cosplay; clearing while `Not Testable` and Verify still soft/missing

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

4. **SFDIPOT scan** (silent checklist - fuels Risks and Test Ideas; do not add new top-level `decisions.md` sections)

   Scan Product Elements; note only what applies. For each relevant area: key area example + risk + coverage idea (charter-style, concrete).

   | Area | Scan for |
   |------|----------|
   | Structure | Components, modules, layers, boundaries |
   | Function | Capabilities, workflows, business rules |
   | Data | Inputs, outputs, persistence, migrations, integrity |
   | Interfaces | APIs, UI surfaces, events, contracts |
   | Platform | OS, browser, device, runtime, deployment target |
   | Operations | Install, config, monitoring, support, recovery |
   | Time/Timing | Concurrency, latency, timeouts, scheduling, ordering |

5. **Quality criteria scan** (silent checklist - fuels Risks and Test Ideas; do not add new top-level `decisions.md` sections)

   Scan quality dimensions; note only what applies. For each relevant criterion: what it means HERE + failure to guard + check/experiment.

   | Criterion | Scan for |
   |-----------|----------|
   | Capability | Does it do what users need? |
   | Reliability | Failures, recovery, consistency |
   | Usability | Learnability, errors, accessibility |
   | Charisma | Delight, trust, polish (when product-relevant) |
   | Security | Authz, data exposure, abuse |
   | Scalability | Load, growth, resource limits |
   | Compatibility | Versions, platforms, integrations |
   | Performance | Speed, throughput, resource use |
   | Installability | Setup, upgrade, rollback |
   | Development | Testability, maintainability, operability |

6. **Fill remaining sections** (use SFDIPOT/quality scan output in Risks and Test Ideas)
   - **Notes** - observations that affect understanding
   - **Risks** - if unclear or incomplete, what breaks in build/ship; include SFDIPOT/quality-informed risks when those apply
   - **Test Ideas** - structured buckets (Positive / Negative / Edge / Overlooked). Default format per idea (compact, planner-usable): `Scenario: … | Steps: … | Expected: …`. UI/user-facing may use story form: `As a [role], when I …, then …` plus brief Notes (risk/edge/usability). Bucket mapping:
     - **Positive** - happy path
     - **Negative** - invalid input, error handling, permission denial
     - **Edge** - rare/boundary conditions
     - **Overlooked** - cases testers often miss (includes exploratory / creative paths; deep exploratory charters may also live under Strategy pass Charter ideas)
     - When UI/visual: if draft or screenshot available, extract UI elements/flows into ideas; else note in Notes or pass `Notes: No screenshot/draft - scenarios based only on textual requirement.`
   - **Implementation pitfalls** - short checklist of potential bugs/pitfalls in **implementing** the requirement (UI, data handling, error messaging, perf, security, etc.) - **distinct** from Requirement Bugs (flaws in the requirement text itself)
   - **Requirement Bugs** - flaws, contradictions, ambiguities in the requirement itself

7. **Record** - write `## Testability pass` (or skip) to `decisions.md`. Tighten `spec.md` Verify when the human confirms fixes.

## Clear rule

- Do **not** clear while Testability is `Not Testable` **and** Verify is still soft/missing
- If the human later supplies machine-checkable Verify, re-score or record a new pass; `Not Testable` alone is not a permanent block once Verify is fixed

## Must / Must not

- **Must**: Exhaust `spec.md` / `questions.md` / `decisions.md` / repo before asking
- **Must**: Park question candidates; ask via sc-clarify one at a time
- **Must not**: Invent Verify from Test Ideas
- **Must not**: Expertise cosplay ("as a QA expert…")
- **Must not**: Replace sc-clarify or mission brief
