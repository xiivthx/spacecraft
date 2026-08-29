---
name: sc-architect
description: "Design system architecture, write ADRs, and analyze tradeoffs using C4 model and design patterns. Activate on \"design the system architecture\", \"write an ADR\", \"choose between microservices and monolith\", or \"architectural decision\"."
---

# sc-architect

Design system architecture under mission control: ADRs, C4, patterns, tradeoffs. Domain detail in `references/`.

## When to use

Architecture docs / C4; ADR / architectural decision; stack tradeoffs (e.g. microservices vs monolith); pattern selection; mission architecture planning.

## Workflow

1. **Resolve** - `spacecraft resolve`; conflict → `spacecraft use <selector>`.
2. **Context** - Read `spec.md`, `decisions.md`, existing architecture artifacts; scope + constraints.
3. **Tradeoffs** - ≥2 alternatives with pros/cons; score vs NFRs; pick with rationale; record in `decisions.md`.
4. **ADR** - Template + required sections: `references/adr-templates.md` (Title, Status, Context, Decision, Consequences).
5. **Diagram** - C4 L1→L4 as needed (drill down only where required). Interactive HTML block/wiring → `sc-diagram`.
6. **Verify** - `spacecraft evidence "<label>" -- echo "Architecture decision documented"` (existence of the documented decision).

### Edge cases

- Domain refs: load `references/<domain>.md` (e.g. `web.md`).
- One-way doors need deeper analysis; still write an ADR for reversible choices.
- No clear winner → simplest option + change conditions. Disagreement → all positions under Alternatives considered.

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

API design inside an established architecture · DB schema/migrations · UI/frontend · interactive HTML diagrams (`sc-diagram`) · code-level patterns · TDD (`sc-tdd`)

## Output format

ADR body + sections: `references/adr-templates.md`. Handshake note: C4 level used + diagram path/description; Alternatives considered (≥2).

## Checklist

Mission resolved · ≥2 alternatives · ADR 5 sections · C4 ≥L1 · linked to acceptance when relevant · YAGNI · domain refs when needed.

## References

- `references/adr-templates.md` - ADR templates, decision frameworks, tradeoff documentation
- `references/web.md` - Web-specific architecture: SPA vs SSR, API gateway, CDN strategy, scaling
