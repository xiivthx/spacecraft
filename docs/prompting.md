# How we instruct agents

Clarity over prompting tricks. Keep always-on rules short; put long detail in skills or `references/`.

## Layers

| Layer | Load | Use for |
|---|---|---|
| `.cursor/rules/` always-on | Every turn | Hard Never/Always, lanes, prompting |
| `.cursor/rules/` with globs | Matching files | Domain constraints |
| `.cursor/skills/*/SKILL.md` | On demand | Workflows (lifecycle `/sc-*` are Skills with `disable-model-invocation: true`, not `.cursor/commands/`) |
| `.cursor/agents/*.md` | Subagent Task | Job contracts |

Lifecycle detail: `.cursor/rules/200-workflow.mdc` + slash skills. `/sc-run` fixes findings before ready and reports them in the summary (no issue ledgers). Do not restate the always-on clarity rule inside every agent.

## Spec Contract

Agents and skills state:

1. **Goal** - why / next use
2. **Output** - format and handshake
3. **Good vs Bad** - success bar
4. **Verify** - how correctness is checked

If unclear: research first; ask for preferences or unverifiable bars via `/sc-discuss`; never invent Verify. Mid `/sc-run`: soft → `decisions.md`; hard → stop and `/sc-discuss`. Ship needs machine-checkable Verify.

## Inner-loop gates (Quick + Mission)

Quick and Mission share always-on INTENT / AUTH / TWINS and the 3-cycle stop (see `.cursor/rules/000-spacecraft.mdc`, `200-workflow.mdc`). AUTH does not bypass `/sc-ship` or `SPACECRAFT_SHIP=1`. Before `ready`, Mission runs `sc-judge` (adversarial prove; ready only on `VERIFIED`).

## Avoid

- Threats, tips, or career-stakes framing
- Expertise cosplay (including STORM lenses as personas instead of jobs)
- Forced chain-of-thought on reasoning models

Role names (Commander, Coder, Tester) are routing contracts, not expertise claims.

**Note:** Lens pass = five decision jobs that write `## Lens pass` or `Lens pass skipped:` in `decisions.md` - not expertise cosplay. See `.cursor/skills/sc-discuss/references/lens-pass.md` and sc-storm (Tier 3).
