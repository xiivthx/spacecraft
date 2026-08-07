# UX/UI review gates

Authoritative harness reference for visual UI quality review in spacecraft missions. Harness process only - not a product feature.

## Goal / when to use

Use this protocol whenever subjective UI quality could influence pass/fail for discuss approval, run draft-parity, visual QC, `sc-designer` critique, `sc-reviewer` findings, or `sc-judge` ready.

- **`/sc-discuss`** - designer gate on draft HTML before human HIL (`UI draft approved`)
- **`/sc-run`** - Step 0 draft-parity, anti-slop tiers, Tier 3 browser visual check before `ready`
- **Review / judge** - `sc-reviewer` consumes structured designer findings; `sc-judge` hunts draft drift on visual ready claims

Do not replace draft-parity, HIL, or existing `sc-designer` dimensions. This raises the bar by encoding how to run them.

## The five gates

1. **Deterministic first** - Prefer ordinary tests, CLI, format checks, tool calls, `npx impeccable detect`, scenario `data-state` presence, evidence re-run, and scripted overflow audits before LLM taste. Reserve critique for what machines cannot assert directly (hierarchy, clutter, subtle parity, motion feel).

2. **Narrow questions per dimension** - Never one blob question ("is this good?"). Separate draft/scaffold readiness, scenario coverage, responsive ladder, draft parity, anti-slop, accessibility blockers, motion intent, and product continuity. One pass/fail verdict per dimension with a short reason.

3. **Pass/fail (not scores)** - Verdicts are `pass`, `fail`, or `uncertain` per dimension. Critique notes may say "uncertain" for humans. **fail-closed for ready/ship:** `uncertain` on a required dimension is treated as `fail` (critical finding, `REFUTED`, blocked ready) - never soft-pass or `VERIFIED` on uncertain visual claims.

4. **Human calibration** - When validating a new UI rubric or changing gate wording, sample representative cases and compare agent findings to human labels. Disagreements mean unclear criteria - tighten the dimension question or machine check; document in `decisions.md` or harness docs. Not a labeling product.

5. **Recheck on change** - When model, task shape, or criteria change, re-run the gates. Prior "passed" may drift. Do not inherit old visual approvals without fresh evidence.

## Dimension table

Required dimensions depend on phase (discuss vs run). Mark **required** when the phase column applies.

| Dimension | Discuss (approval draft) | Run (implementation) | Machine-checkable first | LLM / human critique |
|-----------|--------------------------|----------------------|-------------------------|----------------------|
| **Draft / scaffold readiness** | required | n/a | `data-draft`, `[data-draft-chrome]` outside `[data-draft-frame]` / `[data-draft-surface]`; viewport toggles present; `Layout bake-off winner:` or skip in `decisions.md` | Scaffold clarity, port readiness, notes outside surface |
| **Scenario coverage** | required (approval candidate) | required | `data-state` panels for `empty`, `error`, `few`, `many` + spec features; loading when async implied; product mapping or tests per state | Real component chrome vs layout boxes; missing implied states |
| **Responsive ladder** | required (multi-region UI) | required when UI ships | Four presets 375 / 768 / 1280 / 1536; frame resize works; overflow/clip scripts at breakpoints | Size-appropriate organization; not pixel-squeezed adjacent presets; widescreen measure control |
| **Draft parity** | n/a | required | CSS token variables vs draft; `data-state` mapping; Step 0 checklist; screenshot paths in evidence | Layout-only match with different chrome; subtle spacing/type drift |
| **Anti-slop / catalog** | recommended | required | `npx impeccable detect` (Tier 1 CLI + browser-rendered rules) | Tier 2 heuristics (glassmorphism, extreme radius, amateur SVG, hero metrics, identical grids) |
| **Accessibility blockers** | when in scope | when in scope | Obvious missing labels on form controls in HTML; `prefers-reduced-motion` respected in CSS | Contrast/focus/keyboard gaps when visually obvious in draft or product |
| **Motion intent vs draft** | when motion in brief | required when motion in brief | `prefers-reduced-motion`; duration/easing within `animation-guidelines.md` | Motion feel matches brief/draft intent |
| **Product continuity** | required when brownfield | required when brownfield | `Product context:` paths recorded; shell files exist | Draft/product reflects parent shell/nav patterns; no floating marketing shell on in-app pages |

Bake-off candidates (pre-winner): scaffold + responsive ladder required; full scenario matrix may defer until approval polish (see `sc-designer`).

## Verdict mapping

Per dimension, emit exactly one of:

| Verdict | Meaning | Ready / ship |
|---------|---------|--------------|
| `pass` | Dimension satisfied with evidence | Allows progress on that dimension |
| `fail` | Dimension not satisfied | Blocks approval / ready on that dimension |
| `uncertain` | Cannot decide with available evidence | **fail-closed:** treat as `fail` for approval, `ready`, and `sc-judge` |

Rules:

- No 1-5 scores or numeric rubrics.
- Critique output may still label `uncertain` in notes for human follow-up.
- `sc-judge`: uncertain draft-parity or visual ready claim ⇒ `REFUTED` (note `uncertain` in hunt reasons). No third verdict string.
- `sc-reviewer`: `uncertain` on a required UI dimension ⇒ critical finding; `status: blocked`.

## Output snippet (reviewers / designers)

Short per-dimension lines only:

```
UX/UI review (ux-ui-review-gates):
- draft-scaffold: pass | fail | uncertain - <one line reason>
- scenario-coverage: pass | fail | uncertain - <reason>
- responsive-ladder: pass | fail | uncertain - <reason>
- draft-parity: pass | fail | uncertain - <reason>  (run only)
- anti-slop: pass | fail | uncertain - <reason>
- a11y-blockers: pass | fail | uncertain | n/a - <reason>
- motion-intent: pass | fail | uncertain | n/a - <reason>
- product-continuity: pass | fail | uncertain | n/a - <reason>
```

Omit `n/a` dimensions or mark `n/a` when out of scope. Required dimensions for the current phase must not be `uncertain` on the ready path.

## Human calibration (gate 4)

Lightweight harness process:

1. Pick 3-5 representative drafts or implementation diffs (pass, fail, edge).
2. Run `sc-designer` (or manual review) with this table; record per-dimension verdicts.
3. Compare to human labels. Mismatches ⇒ tighten the dimension question or add a machine check.
4. Record outcome in `decisions.md` when rubric wording changes (e.g. `UX rubric calibrated: <date> - <note>`).

## Recheck on change (gate 5)

Re-run relevant gates when any of these change:

- Subagent model or reviewer/designer prompt
- Task shape (new breakpoints, new states, brownfield context added)
- Dimension table or anti-slop catalog rules

Prior `UI draft approved` or visual evidence does not grandfather later runs without fresh Step 0 / Tier 3 / judge re-run.

## Cross-links

- `.cursor/skills/sc-ux-design/SKILL.md` - discuss draft workflow, Step 0, anti-slop tiers
- `.cursor/agents/sc-designer.md` - discuss and run critique dimensions
- `.cursor/skills/sc-judge/SKILL.md` - draft drift hunt, `VERIFIED` | `REFUTED` only
- `references/anti-slop-catalog.md` - deterministic slop patterns and fixes
- `references/shared-draft-directives.md` - draft scaffold, scenario matrix, responsive structure
- `docs/ux-ui-review.md` - short human-facing overview
