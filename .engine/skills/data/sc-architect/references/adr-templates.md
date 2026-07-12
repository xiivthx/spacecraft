> Consult when: writing an ADR, structuring a decision record, or documenting tradeoff analysis.

# ADR templates

Architecture Decision Records: lightweight, immutable documents that capture the context and rationale behind architectural choices.

## Standard ADR template

```markdown
# ADR-<NNN>: <descriptive title>

**Status**: proposed | accepted | deprecated | superseded by [ADR-XXX]
**Date**: YYYY-MM-DD
**Deciders**: <names or roles>

## Context

What problem are we solving? What constraints exist?
Describe the forces at play: technical, business, timeline, team.
If this decision replaces or is informed by a previous ADR, reference it.

## Decision

What did we choose? Be specific — name the technology, pattern, or approach.
Why this option over the alternatives? Provide explicit rationale.

## Consequences

What becomes easier? What becomes harder?
What new constraints does this introduce?
What follow-up actions or ADRs are needed?

## Alternatives considered

### Alternative 1: <name>
- **Pros**: ...
- **Cons**: ...
- **Why rejected**: ...

### Alternative 2: <name>
- **Pros**: ...
- **Cons**: ...
- **Why rejected**: ...
```

## ADR lifecycle

```
Proposed → Accepted → (maybe) Deprecated → Superseded
```

- **Proposed** — under discussion; not yet binding
- **Accepted** — agreed and active; implementation follows this decision
- **Deprecated** — no longer followed but not yet replaced
- **Superseded** — replaced by a newer ADR; reference the superseding ADR

## When to write an ADR

Write an ADR when:
- The decision affects system structure (monolith vs microservices, layered vs hexagonal)
- The decision constrains future choices (language, framework, database)
- The decision has non-obvious tradeoffs (consistency vs availability, simplicity vs performance)
- The decision is a one-way door — hard to reverse (data model, auth protocol, API versioning)
- Team members have expressed conflicting preferences

Skip an ADR when:
- The decision is trivially reversible (library choice within same category)
- The decision is forced by existing constraints (must use company-wide auth provider)
- The decision has zero architectural impact (formatting, linting rules)

## Tradeoff documentation

For every significant architectural choice, document the tradeoff in a structured format:

```
Decision: <what was chosen>
Drivers: <business or technical requirements pushing this choice>
Upside: <what improves — speed, simplicity, maintainability>
Downside: <what degrades — flexibility, performance, operational complexity>
Reversal cost: low | medium | high | one-way door
Revisit trigger: <what would cause us to reconsider>
```

## Decision frameworks

**Two-way door vs one-way door** (Bezos):
- Two-way door: easy to reverse. Decide quickly with ~70% information.
- One-way door: hard/expensive to reverse. Slow down, get more data, consider more alternatives.

**Build vs buy vs adopt**:
Track record → Community health → Integration complexity → Total cost (not just license) → Lock-in risk

**Consistency models**:
Strong consistency → easier reasoning, harder scaling. Eventual consistency → harder reasoning, easier scaling. Match the model to the use case, not the habit.

## Spacecraft integration

- Store ADRs in the mission `decisions.md` for mission-scoped decisions
- For project-level decisions, create `adr/` directory with numbered ADR files
- Link each ADR to relevant `plan.json` tasks and `spec.md` acceptance criteria
- Architecture decisions that affect implementation are part of mission evidence — `scripts/spacecraft evidence "adr:<title>" -- echo "ADR documented"`
