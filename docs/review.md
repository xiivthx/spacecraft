# Review gates

Harness process for `ready` / ship — not a product feature. Fail-closed: `uncertain` counts as fail.

## Shared idea

1. Deterministic checks before LLM taste
2. Narrow per-dimension pass/fail
3. `pass` | `fail` | `uncertain` (uncertain → fail for ready/ship)
4. Recalibrate when rubric wording changes
5. Re-run after model / task / criteria change

## Mission

SoT: `.cursor/skills/sc-run/references/mission-review-gates.md`

Ready path: suite clean → Cursor `bugbot` + `security-review` → `sc-reviewer` → `sc-judge` → `ready` on `VERIFIED`. Proof: `evidence.jsonl` + empty `review.json` findings + `validate --strict` + judge (incl. `Loop watch:`). After ready: `Post-ready drain:` then `Split-to-prs:` before `/sc-ship` (see `sc-run/references/optional-lanes.md`).

## UX / UI

SoT: `.cursor/skills/sc-ux-design/references/ux-ui-review-gates.md`

Applies when visual UI is in scope (`/sc-discuss` draft gate; `/sc-run` draft-parity + live-product). Impeccable craft: [impeccable.md](./impeccable.md).
