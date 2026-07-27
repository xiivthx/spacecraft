---
name: sc-storm
description: "Tier 3 systematic multi-source research for open-domain strategy questions feeding /sc-discuss. Use for literature-style synthesis, not stuck-API gray areas (those stay sc-search)."
---

# sc-storm

## Goal

Run STORM Tier 3: gather sources, draft five lens notes, debate contradictions, and produce greppable discuss artifacts - without inventing Verify or touching `/sc-run` build.

## Output

Under `.space/missions/<id>/`:

1. `research-brief.md` - question, sources, lens notes, contradictions, synthesis
2. `decisions.md` - one `## Lens pass (<topic>)` block copied from template in `references/lens-pass.md` (path: `.cursor/skills/sc-discuss/references/lens-pass.md`), ending with `Lens tier used: 3`

Handshake: `done` | `blocked: <reason>` | `needs-input: <question>`

## Good / Bad

- Good: framed question; primary sources preferred; five lens bullets are jobs not personas; one Synthesis path; feeds discuss only; Historian skims `.space/trust/lessons.md`
- Bad: lifecycle slash peer to discuss/run/ship/quick; always-on; product code; inventing Verify; stuck API / deprecation lookups (use sc-search); full STORM at ready (sc-judge stays Skeptic-only)

## Verify

```
test -f .space/missions/<id>/research-brief.md
rg -q '## Lens pass' .space/missions/<id>/decisions.md
rg -q 'Lens tier used: 3' .space/missions/<id>/decisions.md
```

Commander confirms Synthesis is one path and Verify bar still comes from spec + sc-clarify, not from lenses.

## When to use

- Open-domain or strategy research before implement
- Systematic multi-source synthesis (competitors, patterns, policy landscape)
- Discuss escalation when sc-search tiers are insufficient **and** the question is not a blocking technical gray area

**Not** for: unfamiliar errors, deprecated APIs, version pins - **sc-search**. **Not** for: `/sc-run` RED-GREEN.

## Workflow

1. **Resolve mission** - `spacecraft resolve` or explicit id; read `spec.md`, `decisions.md`.
2. **Frame question** - one decision-shaped question; if Verify is preference-bound, flag for sc-clarify - do not invent bars.
3. **Gather sources** - `WebSearch` / `WebFetch` (sc-search-style escalation); prefer official docs, papers, primary posts; note contradictions.
4. **Draft five lens notes** - Practitioner, Academic, Skeptic, Economist, Historian (skim `.space/trust/lessons.md`).
5. **Debate contradictions** - reconcile or state what remains contested; still one Synthesis path.
6. **Write** `research-brief.md` in the mission dir (question, sources list, lens notes, synthesis).
7. **Copy** Synthesis + five bullets into `decisions.md` as `## Lens pass (<topic>)` per `references/lens-pass.md`; set `Lens tier used: 3`.

## Out of scope

- `/sc-run` build, product code, tests
- Inventing Verify
- Always-on use
- `/sc-research` slash (retired - this skill replaces that pointer)

## References

- `.cursor/skills/sc-discuss/references/lens-pass.md` - shared template and tier rules
- sc-search - technical gray areas and API stuck states
- sc-clarify - preference-bound Verify
- `/sc-discuss` - consumer of lens pass artifacts
