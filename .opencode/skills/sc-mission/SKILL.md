---
name: sc-mission
description: Manage mission artifacts and lifecycle for local OpenCode development
license: MIT
compatibility: opencode
metadata:
  version: 1.0
  audience: xiivthx
---

# sc-mission

Manage mission artifacts and lifecycle for local OpenCode development.

## When to use

Activate when the user asks to:

- Resolve, list, or select a mission
- Create mission artifacts (spec, plan, evidence)
- Check mission state or lifecycle
- Route ambiguity to clarification
- Handle session handoff or release closeout

## Workflow

Use this exact sequence unless the user specifies otherwise:

1. **Resolve mission** — Run `scripts/spacecraft resolve [selector] [--json]`, `status`, or `missions`; `.space/current` is fallback state, not sole authority.
2. **Read artifacts** — Read `mission.json`, `spec.md`, `questions.md`, `decisions.md`, `plan.json`, design artifacts, `evidence.jsonl`, and `review.json` when available.
3. **Route ambiguity** — If intent, scope, or acceptance criteria is ambiguous, route to sc-clarify before proceeding.
4. **Enforce lifecycle** — Follow: mission -> clarify -> spec -> visual design if needed -> plan -> build -> verify -> review -> ship. Use `/sc-build` to repeat build -> verify -> checkpoint commit for successive tasks.
5. **Release or handoff** — On ship intent, run release closeout. On session end, give handoff summary.

## Rules

- **Must**: Resolve the active mission with `scripts/spacecraft resolve [selector] [--json]`, `status`, or `missions`; `.space/current` is fallback state, not sole authority.
- **Must**: Resolver priority is explicit selector or `SPACECRAFT_MISSION`, session binding, branch mission id, branch metadata, `.space/current`, then single active mission.
- **Must**: Strong signal conflicts or ambiguous active missions block mission writes until the user selects with `scripts/spacecraft use <number|id|title>` or an explicit selector.
- **Must**: Users may choose by list number, mission id, exact title, or unique title substring; do not expect the user to know a mission id.
- **Must**: New mission and evidence ids are compact sortable ids with no hyphen, such as `M07FYB5W5`; legacy `M-YYYYMMDD-HHmmss` ids remain valid.
- **Must**: Read the resolved mission's `mission.json`, `spec.md`, `questions.md`, `decisions.md`, `plan.json`, design artifacts, `evidence.jsonl`, and `review.json` when available.
- **Must**: sc-mission owns lifecycle but must route ambiguity to sc-clarify.
- **Must**: Enforce order: mission -> clarify -> spec -> visual design if needed -> plan -> build -> verify -> review -> ship. `/sc-build` may repeat build -> verify -> checkpoint commit for successive tasks until a gate blocks.
- **Must not**: Skip clarification when user intent, scope, acceptance criteria, or visual design direction is materially ambiguous.
- **Must**: If clear mutating work is requested and no suitable mission or branch exists, create the mission and non-main branch without another blocking question when policy permits it.
- **Must not**: Implement if spec or plan is missing.
- **Must**: Use sc-git for git safety, release branching, commits, rebasing, merge, version bump, changelog/spec notes, and tagging.
- **Must**: Mutating work must satisfy sc-git gates before implementation and before ship.
- **Must**: Treat stop-chat/close-session/end-session/new-session requests as session handoff unless release intent is explicit.
- **Must**: If "close session" is ambiguous and work appears ready, recommend `/sc-ship`; do not merge automatically.
- **Must**: Treat ship/release/merge/finish-mission/close-branch requests as release closeout prep. Block closeout when gates are incomplete.
- **Must**: After successful release closeout, archive shipped mission artifacts under `.space/archive/` unless the user asks to keep the full live mission folder.
- **Must**: Keep mission artifacts small and human-readable.
- **Must**: Prefer explicit evidence over narrative claims.
- **Must**: End each session with a recommended next action and session advice: continue this chat for small adjacent steps, or start a new session when the phase changed, the thread is context-heavy, or mission artifacts are sufficient for handoff.

## Out of scope

This skill does NOT handle:

- Git operations, branching, or commits — use sc-git
- Planning tasks — use sc-planning
- Design or UI direction — use sc-design
- Evidence capture or verification — use sc-verification
- Clarification questions — use sc-clarify (sc-mission routes to it)

## Output format

```
Mission artifacts follow the standard layout:
.space/missions/<id>/
  mission.json     # core metadata
  spec.md          # spec
  questions.md     # open/answered questions
  decisions.md     # confirmed choices and assumptions
  plan.json        # task plan
  evidence.jsonl   # evidence entries
  review.md        # review narrative
  review.json      # review findings and readiness
  design/          # design artifacts
  outputs/         # task outputs
```

## Checklist

Before claiming the mission lifecycle is handled:

- [ ] Mission resolved and confirmed safe before writes
- [ ] All relevant artifacts read before decisions
- [ ] Ambiguity routed to sc-clarify when blocking
- [ ] Lifecycle enforced: clarify -> spec -> visual design -> plan -> work -> verify -> review -> ship
- [ ] Session handoff or release closeout handled correctly per user intent

---

## References

- `scripts/spacecraft --help` — spacecraft CLI reference
- `scripts/spacecraft resolve --help` — resolver subcommand
- `scripts/spacecraft use --help` — mission selection
- `scripts/spacecraft evidence --help` — evidence subcommand
- `scripts/spacecraft closeout-check --help` — closeout verification
