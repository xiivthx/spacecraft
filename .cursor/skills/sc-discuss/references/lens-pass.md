# Lens pass (STORM Tier 0-3)

Five STORM lenses as **decision jobs** that produce greppable artifacts in `decisions.md` - not personas, not lens-named agents, not expertise cosplay.

| Lens | Job |
|------|-----|
| Practitioner | What ships daily; ops and maintenance reality |
| Academic | Known patterns, literature, formal tradeoff framing |
| Skeptic | Failure modes, hidden assumptions, what we might be wrong about |
| Economist | Cost, scope, ongoing tax, build vs defer |
| Historian | Prior decisions in repo history and mission `decisions.md` (`none found` allowed) |

## When required (any one)

- Architecture or multi-option tradeoff with material fork
- Product or policy preference not settled by spec alone
- Soft or preference-bound Verify (human taste, bar unclear)
- Sizing may explode (roadmap split, phases, or jigsaw count uncertainty)

## When skip

- Routine bugfix with clear Verify
- `/sc-quick` (no mission lens gate)
- Verify already machine-checkable and uncontested
- No material decision this discuss session

Record skip in `decisions.md`:

```
Lens pass skipped: <reason>
```

## Shared `decisions.md` template

Use exactly this shape when a lens pass runs (one `## Lens pass` block per topic):

```
## Lens pass (<topic>)
- Practitioner: …
- Academic: …
- Skeptic: …
- Economist: …
- Historian: …  # repo history / decisions.md; none found allowed
- Synthesis: <one path + what we reject>
- Lens tier used: 0|1|2|3
```

Greppable markers: `## Lens pass` OR `Lens pass skipped:`

## Tier 0 - Discuss checklist (Commander)

Commander fills the five bullets + Synthesis in discuss (no subagent required). Record in `decisions.md` with `Lens tier used: 0`. Default for preference-bound product/process choices that do not need deep architecture.

## Tier 1 - sc-adviser

When escalation triggers fire, Task(`sc-adviser`) Output includes a lens pass **before** the single recommendation. Template: this file. Record copy in `decisions.md` when discuss owns the decision; adviser chat may hold the pass first. `Lens tier used: 1`.

## Tier 2 - Parallel readonly Tasks (Commander)

High-stakes fork: Commander may launch **2-3** readonly Tasks with lens-scoped prompts.

- Default trio: Skeptic, Economist, Practitioner (not five by default)
- **Never** five agents by default
- Synthesize into one `## Lens pass` in `decisions.md`
- Prefer Task(`sc-adviser`) or focused prompts - no lens-named agents
- `Lens tier used: 2`

## Tier 3 - sc-storm

Open-domain / strategy / systematic multi-source research feeding discuss only. Skill: `sc-storm` (not a lifecycle slash). Writes `research-brief.md` under the mission dir and copies Synthesis into `decisions.md` with `Lens tier used: 3`.

## Hard rules

- **Must not** invent Verify from lenses - preference-bound bars → sc-clarify
- **Not used** in `/sc-run` build loop
- **sc-judge** at ready stays Skeptic-only adversarial prove - no full STORM at ready
- **Must not** spawn agents named after lenses
- One Synthesis path - lenses inform; they do not vote as separate owners

## Must not

- Always-on five lenses every mission
- STORM inside RED-GREEN
- Revive a retired research lifecycle slash as peer to discuss/run/ship
- Expertise cosplay ("as a historian I feel…") without an artifact line
