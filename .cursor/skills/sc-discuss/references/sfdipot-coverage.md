# SFDIPOT coverage review

On-demand procedure: evaluate quality and completeness of **existing tests** against a requirement using SFDIPOT. Decision job - not expertise cosplay. **Not** an always-on discuss soft gate.

## Goal

Assess whether current tests adequately cover the requirement across Structure, Function, Data, Interfaces, Platform, Operations, and Time - and surface blind spots without inventing Verify.

## Output

Chat summary plus greppable artifact in `decisions.md`:

```
## SFDIPOT coverage
- Requirement summary: …
- Current test list summary: …
- Structure: Coverage: Covered | Partially Covered | Missing | Observations: … | Suggestions: …
- Function: Coverage: … | Observations: … | Suggestions: …
- Data: Coverage: … | Observations: … | Suggestions: …
- Interfaces: Coverage: … | Observations: … | Suggestions: …
- Platform: Coverage: … | Observations: … | Suggestions: …
- Operations: Coverage: … | Observations: … | Suggestions: …
- Time: Coverage: … | Observations: … | Suggestions: …
- Overall assessment:
  - Strengths: …
  - Weaknesses: …
  - Priority recommendations: …
```

Or skip:

```
SFDIPOT coverage skipped: <reason>
```

## Good / Bad

- Good: concrete coverage ratings per SFDIPOT area; observations tied to named tests or gaps; suggestions actionable for planner/tester; distinguishes blind spots from false-completion claims
- Bad: expertise cosplay; inventing Verify from suggestions; replacing testability pass, strategy pass, or RCRCRC; blocking discuss clear on coverage gaps alone

## Verify

When run: `decisions.md` has `## SFDIPOT coverage` OR `SFDIPOT coverage skipped:`. Does **not** block discuss clear. Greppable markers only - no CLI validator required.

## When to use

- Human asks for SFDIPOT coverage review, test suite review, or blind-spot hunt against requirement
- Mid-run suite review before expanding acceptance or claiming coverage complete
- sc-judge optional hunt aid when `## Testability pass` Test Ideas or `## Strategy pass` Charter ideas exist

## When skip

- No requirement text available
- No current test list (files, suite summary, or pasted cases)
- Routine bug with clear repro and single targeted test
- Discuss clear not blocked by missing this pass - record skip if invoked but inputs missing

## Error handling

- Requirement missing → respond exactly: `Requirement missing - please provide requirement details.`
- Test list missing → respond exactly: `Test list missing - please provide the current test cases.`

## Procedure

1. **Read** - requirement (`spec.md` or human paste). If missing → error phrase above. Current tests (suite files, pasted list, or evidence of what exists). If missing → error phrase above.

2. **Requirement summary** - one tight paragraph: what must be true, key behaviors, constraints.

3. **Current test list summary** - enumerate or summarize existing tests/cases and what each claims to cover.

4. **SFDIPOT analysis** - for each area below, rate **Coverage** (`Covered` | `Partially Covered` | `Missing`), **Observations** (what tests address or miss), **Suggestions** (concrete gaps - not invented Verify):

   | Area | Scan for |
   |------|----------|
   | Structure | Components, modules, layers, boundaries |
   | Function | Capabilities, workflows, business rules |
   | Data | Inputs, outputs, persistence, migrations, integrity |
   | Interfaces | APIs, UI surfaces, events, contracts |
   | Platform | OS, browser, device, runtime, deployment target |
   | Operations | Install, config, monitoring, support, recovery |
   | Time | Concurrency, latency, timeouts, scheduling, ordering |

5. **Overall assessment**
   - **Strengths** - well-covered areas
   - **Weaknesses** - systematic gaps
   - **Priority recommendations** - ordered follow-ups for planner/tester (still no invented Verify)

6. **Record** - write `## SFDIPOT coverage` to `decisions.md` (or skip message).

## Clear rule

- **Soft** - does **not** block discuss clear. Missing SFDIPOT coverage alone is not a clarify blocker.
- **sc-judge** - blind-spot suggestions alone do **not** justify `REFUTED`. `REFUTED` only when claimed acceptance or Test Idea coverage maps to a **Missing** SFDIPOT area that was asserted done **without** fresh evidence (false completion).

## Must / Must not

- **Must**: Ground analysis in actual requirement text and actual test list
- **Must not**: Expertise cosplay ("as a QA expert…")
- **Must not**: Invent Verify from suggestions or coverage gaps
- **Must not**: Replace `requirement-testability.md`, `htsm-strategy.md`, or `rcrcrc-impact.md`
- **Must not**: Add always-on discuss gate behavior
