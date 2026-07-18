# How we instruct agents

Spacecraft follows Cursor’s rules/skills split and Wharton Prompting Science: clarity beats prompting tricks.

## Layers

| Layer | Load | Use for |
|---|---|---|
| `.cursor/rules/` always-on | Every turn | Hard Never/Always, lanes, prompting standard |
| `.cursor/rules/` with globs | When matching files | Domain constraints (firmware, security, DB) |
| `.cursor/skills/*/SKILL.md` | On demand | Workflows and runbooks |
| `.cursor/agents/*.md` | Subagent Task | Job contracts with checkable outputs |

Keep always-on short. Put long detail in skills or `references/`.

## Spec Contract

Every agent and new skill should state:

1. **Goal** - why / how the result is used next
2. **Output** - format and handshake
3. **Good vs Bad** - success bar for this task
4. **Verify** - how correctness is checked

## Clarity gate

If any of the four are unclear: research repo/spec/plan/docs first; ask the human for preferences or unverifiable bars; never invent Verify. Mid `/sc-run`, soft gaps go to `decisions.md`; hard gaps stop for clarify. Ship requires machine-checkable Verify.

## Do not use

- Threats, tips, or “important to my career” framing
- Expertise cosplay (“senior expert with 20 years”) for accuracy
- Forced chain-of-thought on reasoning models (waste tokens; little gain)

Role names (Commander, Coder, Tester) are routing contracts (permissions + output shape), not expertise claims.
