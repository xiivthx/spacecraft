# How we instruct agents

Clarity over prompting tricks. Keep always-on rules short; put long detail in skills or `references/`.

## Layers

| Layer | Load | Use for |
|---|---|---|
| `.cursor/rules/` always-on | Every turn | Hard Never/Always, lanes, prompting |
| `.cursor/rules/` with globs | Matching files | Domain constraints |
| `.cursor/skills/*/SKILL.md` | On demand | Workflows |
| `.cursor/agents/*.md` | Subagent Task | Job contracts |

Do not restate the always-on clarity rule inside every agent.

## Spec Contract

Agents and skills state:

1. **Goal** - why / next use
2. **Output** - format and handshake
3. **Good vs Bad** - success bar
4. **Verify** - how correctness is checked

If unclear: research first; ask for preferences or unverifiable bars via `/sc-discuss`; never invent Verify. Mid `/sc-run`: soft → `decisions.md`; hard → stop and `/sc-discuss`. Ship needs machine-checkable Verify.

## Avoid

- Threats, tips, or career-stakes framing
- Expertise cosplay
- Forced chain-of-thought on reasoning models

Role names (Commander, Coder, Tester) are routing contracts, not expertise claims.
