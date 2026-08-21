# Prompt refine

On-demand procedure: diagnose and rewrite agent/skill/rule prompts for fidelity to the Spec Contract. Decision job - not expertise cosplay. **Not** always-on. **Not** a tutor.

## Goal

Diagnose draft prompt text against house Spec Contract fidelity, then rewrite only what fails or needs caution - without changing gates, policy, or runtime behavior.

## Output

```
## Prompt refine
### Diagnosis
Goal/structure: …
Ratings:
1. Task Fidelity - Pass|Caution|Fail - …
2. Clarity and Specificity - Pass|Caution|Fail - …
3. Context Utilisation - Pass|Caution|Fail - …
4. Accuracy and Verifiability - Pass|Caution|Fail - …
5. Tone and Persona Consistency - Pass|Caution|Fail - …
6. Error Handling - Pass|Caution|Fail - …
7. Resource Efficiency - Pass|Caution|Fail - …
Triggers: … (or none)
### Before/After
… (≤2 lines micro-example, or one-sentence N/A rationale)
### Revised
``` … revised prompt … ```
```

Or:

```
Prompt refine skipped: <reason>
```

Revised prompts for agents/skills keep **Goal**, **Output**, **Good vs Bad**, and **Verify** when editing agent/skill text.

## Good / Bad

- Good: Pass/Caution/Fail ratings with one-line notes; rewrite only Caution/Fail items; preserve purpose, scope, persona/routing contract; numbered-step Procedure; brevity + clarity; explicit trigger resolution when marked; Spec Contract Verify present and never invented; US English; ASCII hyphen-minus only
- Bad: expertise cosplay; threats/tips/career-stakes framing; forced chain-of-thought; Unicode em dash; changing what gates/rules *do*; inventing Verify; applying lyrical rhythm craft (that is `prose-rhythm.md`); questionnaire dumps; light paraphrase without diagnosis

## Verify

When run: `## Prompt refine` with Diagnosis + Before/After + Revised (or fenced revised) OR `Prompt refine skipped:` with reason. Does **not** block discuss clear or ship.

## When to use

- Human asks to refine or diagnose a draft prompt
- `Task(sc-writer)` for agent/skill/rule prompt fidelity rewrite
- Before shipping new always-on rule text that needs a compression check

## When skip

- Narrative engagement rewrite → `prose-rhythm.md`
- Thin narrative context → `narrative-context.md`
- Product code or tests
- Changing gate behavior (stop and report instead)

## Procedure

### Phase 1 - Rapid Diagnosis

1. **Goal/structure** - one short paragraph: draft prompt's goal and structure.
2. **Rate** each criterion Pass | Caution | Fail with one-line note:
   1. Task Fidelity (Goal matches author intent)
   2. Clarity and Specificity
   3. Context Utilisation (inputs/siblings/house Avoid used)
   4. Accuracy and Verifiability (Spec Contract Verify present; never invent Verify)
   5. Tone and Persona Consistency (routing job, not cosplay; US English; ASCII hyphen-minus)
   6. Error Handling (edge cases / blocked / needs-input)
   7. Resource Efficiency (tokens/latency; always-on vs on-demand; no filler)
3. **Triggers** - mark any that apply: Context Preservation | Intent Refinement | Error Prevention

### Phase 2 - Precision Rewrite

1. Improve only Caution/Fail items.
2. Preserve purpose, scope, persona/routing contract.
3. Prefer numbered-step Procedure.
4. Brevity + clarity.
5. If a trigger was marked, explicitly show how addressed.

## Must / Must not

- **Must**: Diagnose before rewrite
- **Must**: Rate all seven criteria with one-line notes
- **Must**: Rewrite only Caution/Fail items
- **Must**: Keep Goal, Output, Good vs Bad, Verify on agent/skill text
- **Must**: Stop and report if wording would change gate or policy behavior
- **Must not**: Expertise cosplay ("as a Senior Prompt Architect…")
- **Must not**: Threats, tips, or career-stakes framing
- **Must not**: Forced chain-of-thought on reasoning models
- **Must not**: Unicode em dash - ASCII hyphen-minus only
- **Must not**: Change what gates, rules, or checks *do* while editing wording
- **Must not**: Invent Verify
- **Must not**: Apply lyrical rhythm craft here (use `prose-rhythm.md` for narrative)
- **Must not**: Questionnaire dumps

## Related

- Spec Contract: `docs/prompting.md`
- Rhythm rewrite (narrative only): `prose-rhythm.md`
- Narrative context harvest: `narrative-context.md`
- Agent: `.cursor/agents/sc-writer.md`
