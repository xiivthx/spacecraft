---
name: sc-architect
description: "Design system architecture, write ADRs, and analyze tradeoffs using C4 model and design patterns. Activate on \"design the system architecture\", \"write an ADR\", \"choose between microservices and monolith\", or \"architectural decision\"."
---

# sc-architect

Design system architecture under mission control. Universal architecture patterns with domain-specific references. Covers ADR format, C4 diagrams, design patterns, and tradeoff analysis.

## When to use

Activate when the user asks to:

- **"Design the system architecture" / "create C4 diagram"** - architecture documentation
- **"Write an ADR" / "architectural decision"** - decision records and rationale
- **"Choose between microservices and monolith"** - architectural tradeoff analysis
- **"What pattern should I use for"** - design pattern selection and justification
- When a mission requires architectural planning or tradeoff documentation

## Workflow

Use this exact sequence unless the user specifies otherwise:

1. **Resolve mission** - `spacecraft resolve`. On conflict or ambiguity, use `spacecraft use <selector>`.

2. **Understand context** - Read the mission `spec.md`, existing `decisions.md`, and any architecture artifacts. Identify the decision scope and constraints.

3. **Analyze tradeoffs** - For each architectural decision:
   - Enumerate ≥2 alternatives with pros and cons
   - Evaluate against non-functional requirements (scalability, maintainability, cost)
   - Select the best option with explicit rationale
   Document the analysis in `decisions.md`.

4. **Write ADR** - Use the template from `references/adr-templates.md`. Minimum sections:
   - **Title** - descriptive and searchable
   - **Status** - proposed, accepted, deprecated, superseded
   - **Context** - what problem does this decision address?
   - **Decision** - what was chosen and why
   - **Consequences** - what becomes easier, harder, or constrained

5. **Diagram** - Use C4 model levels as appropriate:
   - **Level 1 (System Context)** - the system and its users/external systems
   - **Level 2 (Container)** - applications, data stores, message queues
   - **Level 3 (Component)** - major structural building blocks within a container
   - **Level 4 (Code)** - class-level detail (only when critical)
   - For **interactive HTML** block/wiring diagrams (click-to-trace nets), use `sc-diagram` instead of ASCII/C4 prose alone.

6. **Verify** - `spacecraft evidence "<label>" -- echo "Architecture decision documented"`. ADRs and diagrams are manual artifacts. Evidence is the existence of the documented decision.

### Edge cases

- **Domain-specific patterns** - Load the relevant `references/` file (e.g., `web.md` for web architecture, add others for mobile, embedded, data pipelines)
- **Reversible decision** - Still write an ADR. Distinguish between one-way and two-way door decisions. One-way doors require deeper analysis.
- **No clear winner** - Document the deadlock, pick the simplest option, note conditions that would change the decision.
- **Team disagreement** - Record all positions in the ADR's "Alternatives considered" section.

## Rules

- **Must**: Resolve mission with `spacecraft resolve` before documenting architecture. On conflict/ambiguity use `spacecraft use <selector>`.
- **Must**: Record non-trivial architectural decisions in ADRs. If it affects system structure, technology choice, or cross-cutting concern, write it down.
- **Must**: Enumerate ≥2 alternatives with pros and cons for each significant decision.
- **Must**: Use C4 model for diagrams. Start at Level 1; drill down only where needed.
- **Must**: Store ADRs in `decisions.md` within the mission directory, or a dedicated `adr/` directory for project-level decisions.
- **Must**: Link ADRs to mission acceptance checks when they inform implementation.
- **Must not**: Design for hypothetical futures (YAGNI). Architecture serves current mission needs.
- **Must not**: Microservices by default. Start simple; distribute only when needed.

## Out of scope

- API-level design within an established architecture - separate concern
- Database schema design, migrations, indexing - separate concern
- UI design or frontend architecture - separate concern
- Interactive HTML block/wiring diagrams - use `sc-diagram`
- Code-level implementation patterns - separate concern
- TDD discipline - use sc-tdd

## Output format

```
ADR: <title>
Status: proposed | accepted | deprecated | superseded
Context: <problem statement>
Decision: <what and why>
Consequences: <tradeoffs>
Alternatives considered:
  - <option>: pros/cons
  - <option>: pros/cons
C4 level: System Context | Container | Component | Code
Diagram: <ascii or description>
```

## Checklist

Before claiming architecture work done:

- [ ] Mission resolved, spec.md and constraints understood
- [ ] ≥2 alternatives considered for each significant decision
- [ ] ADR written with all 5 sections (Title, Status, Context, Decision, Consequences)
- [ ] C4 diagram at appropriate level (≥Level 1)
- [ ] Decision linked to mission acceptance checks where applicable
- [ ] No speculative architecture (YAGNI applied)
- [ ] Domain-specific references consulted when applicable

## References

- `references/adr-templates.md` - ADR templates, decision frameworks, tradeoff documentation
- `references/web.md` - Web-specific architecture: SPA vs SSR, API gateway, CDN strategy, scaling
