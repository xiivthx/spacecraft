# Persona walkthrough pack

Optional live pack for `sc-browser-probe`. Cognitive walkthrough by archetype **jobs** - not STORM lenses, not discuss cosplay, not a 1-5 scorecard.

Coverage id: `pack:persona-walkthrough`

## When to match

**Must match** when any:

- Call has `persona: on`
- User asks multi-persona, cognitive walkthrough, or persona probe
- `decisions.md` has `Persona pack: required`

**Must not** auto-match every probe. Skip when unmatched; no coverage row required.

## Procedure

1. Pick **3-5** archetypes from the table below (drop `domain` unless spec names a specialist role).
2. Fill matrix: archetype id, goal (from inventory/spec), key constraint. No biography essays.
3. Choose journeys from **inventory + spec** only - do not invent product-specific journeys in this file.
4. For each archetype × critical journey step, answer the four cognitive questions. Any `no` → finding with repro.
5. Map severity; file under Findings. No numeric scores.

## Archetypes (jobs)

| id | Job focus |
|----|-----------|
| novice | Onboarding clarity, jargon, error recovery, cognitive load |
| power | Speed, keyboard, bulk/filter/export, density |
| a11y | Contrast, focus, keyboard-only, zoom 200%, SR labels |
| mobile | Thumb reach, glance ≤3s, flaky net, one-hand use |
| privacy | Storage transparency, export/purge, trackers, session isolation |
| domain | Domain terms/units/handoff fidelity - only when spec requires |

## Cognitive questions (per step)

1. Will the user try the right effect?
2. Will they notice the correct action is available?
3. Will they associate that action with the intended outcome?
4. After acting, will they see clear progress feedback?

## Severity map

| Walkthrough | Probe finding |
|-------------|---------------|
| Blocks task; no workaround | **critical** |
| Severe confusion / high abandon risk | **important** |
| Hesitation, fatigue, suboptimal path | **minor** |
| Polish / microcopy only | **minor** (or skip if Foundations already covers) |

Do not invent a parallel S4-S1 ladder in the report - use critical / important / minor only.

## Coverage score

Mark `pack:persona-walkthrough` as:

- `ok` - matrix run; every `no` filed as a finding (or zero `no`)
- `fail` - unused (findings carry fail); prefer `ok` + findings
- `n/a` - pack not matched
- `deferred` - matched but not reached (timebox / blocked)

Notes must include: archetype ids used + journey ids + step count.

## Must / Must not

- **Must**: Journeys from inventory/spec; findings have repro; pass/fail via findings only
- **Must not**: STORM lens names as personas; expertise monologues; 1-5 scores; MoSCoW essays in report (priority is finding severity)
- **Must not**: Replace Foundations, surface catalog packs, draft-parity, or sc-judge
- **Must not**: Bake product-domain scenarios into this reference

## Related

- `../SKILL.md` - call `persona: off \| on`
- `surface-match.md` - catalog packs stay separate
- `scenario-matrix.md` - extra bucket when this pack matched
- `report-template.md` - coverage row + findings
- `.cursor/skills/sc-discuss/references/lens-pass.md` - decision jobs; do not conflate
