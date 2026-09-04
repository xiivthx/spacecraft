# How we instruct agents

Clarity over tricks. Keep always-on rules short; detail lives in skills / `references/`.

## Layers

| Layer | Load | Use |
|-------|------|-----|
| Always-on rules | Every turn | Hard Never/Always, lanes |
| Glob rules | Matching files | Domain constraints |
| Skills | On demand | Workflows (`/sc-*` skills; not `.cursor/commands/`) |
| Agents | Task | Goal / Inputs / Ban / Handshake |

Lifecycle: `.cursor/rules/200-workflow.mdc` + slash skills.

## Spec Contract

1. **Goal** · 2. **Output** · 3. **Good vs Bad** · 4. **Verify**

Unclear → research, then `/sc-discuss` for preferences; never invent Verify. Mid-run soft gaps → `decisions.md`; hard → `/sc-discuss`. Ship needs machine-checkable Verify.

Discuss grill: re-pitch / research (`sc-search` → `sc-storm`) / visualize. Chat asks: Thai bodies, English labels. Detail: `sc-clarify`.

## Gates

INTENT / AUTH / TWINS / 3-cycle (always-on). AUTH does not bypass `/sc-ship` or `SPACECRAFT_SHIP=1`.

- Before `ready`: `bugbot` + `security-review`; Spacecraft supplementary ([review.md](./review.md)); `sc-judge` → `VERIFIED` only; greppable `Loop watch:`
- Ready proof: `evidence.jsonl` + empty findings + `validate --strict` + judge
- After `ready`: `Post-ready drain:` then `Split-to-prs:` (`ran` or `skipped:`) before ship — `sc-run/references/optional-lanes.md`
- Overnight: `sc-loop` disposition at `/sc-run` start; stop → `handback.md`

## On-demand jobs

- **Lens pass** — `sc-discuss/references/lens-pass.md`; Tier 3 → `sc-storm`
- **Fact-check** — `Fact-check: corroborated` \| `contested:` \| `skipped:` — `sc-search/references/fact-check.md`; agent `sc-fact-check`
- Other discuss/test craft jobs live under skill `references/` (testability, strategy, defect-finding, prompt-refine, …)

## Avoid

Threats / career stakes · expertise cosplay · forced chain-of-thought · tombstones after cuts (cut hygiene in `000-spacecraft.mdc`). Role names are routing contracts, not expertise claims.
