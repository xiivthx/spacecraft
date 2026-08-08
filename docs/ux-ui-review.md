# UX/UI review gates

Spacecraft harness process for missions with visual UI. This is not a product feature - it is how agents and humans decide when draft HTML and implemented UI are good enough to approve or ship.

## Five gates

1. **Deterministic first** - Tests, CLI (`npx impeccable detect`), scenario `data-state` checks, and evidence re-run before LLM taste.
2. **Narrow per dimension** - Separate pass/fail questions (parity, ladder, slop, a11y, motion, continuity) - not one "is this good?" blob.
3. **Pass/fail, fail-closed** - Verdicts are `pass`, `fail`, or `uncertain`. For `ready` and ship, **`uncertain` counts as fail** - critical finding or `REFUTED`, never soft-pass.
4. **Human calibration** - When rubric wording changes, sample cases and compare agent vs human labels; disagreements mean unclear criteria.
5. **Recheck on change** - Model, task shape, or criteria change ⇒ re-run gates; old passes do not carry forward without fresh evidence.

## Where it applies

| Phase | Use |
|-------|-----|
| `/sc-discuss` | Designer gate before `UI draft approved` |
| `/sc-run` | Draft-parity (Step 0), anti-slop, Tier 3 live product review on the running product URL |
| Review / judge | `sc-reviewer` findings; `sc-judge` draft drift hunt |

**Live product review** is required before visual `ready`: start the app, open real product routes, capture screenshots, and pass the **live-product** dimension (fail-closed). Draft HTML serve alone does not satisfy that gate. **Draft parity** also requires paired draft-surface + live screenshots at matching viewports and a side-by-side compare before ready. Human browser check stays after ready.

## Source of truth for agents

Full dimension table, machine vs critique split, and output format:

`.cursor/skills/sc-ux-design/references/ux-ui-review-gates.md`

Related: `sc-designer` agent, `sc-ux-design` skill, `sc-judge` draft drift, `anti-slop-catalog.md`, `shared-draft-directives.md`, `mission-review-gates.md` (sibling for every mission).
