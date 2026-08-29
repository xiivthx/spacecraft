---
name: sc-search
description: "Quick internet search with 3-tier escalation for resolving stuck issues, gray areas, and stale knowledge. Activates on unfamiliar errors, deprecated APIs, dependency uncertainty, or technical gray areas."
---

# sc-search

Quick internet search with 3-tier escalation for stuck issues, gray areas, and stale knowledge.

## When to use

Unfamiliar error/stack; deprecated or unfamiliar API; dependency version/compat uncertainty; technical gray area where knowledge may be stale.

## Workflow

Never skip a tier unless a shortcut applies.

### Tier 1 - Quick search (`WebSearch`)

Targeted query (error, API, package + version). ~5s. Escalate if unclear.

### Tier 2 - Deep read (`WebFetch`)

Fetch a Tier-1 URL; extract concrete answer. ~10s. Optional: `ctx_fetch_and_index` → `ctx_search` for multi-query of one page.

### Tier 3 - Multi-source synthesis

Refined WebSearch + 2–4 WebFetch (docs, release notes, issues). Prefer primary docs; note contradictions. Open-domain strategy for discuss → **sc-storm**, not this skill.

### Tier 4 - Ask user

One question: what you tried (tiers 1–3) + recommended next step. Do not implement while open.

## Rules

- **Must**: Always try Tier 1 first. Only escalate when the answer is still unclear.
- **Must**: Never search for the same thing twice in one session. Cache results in session context.
- **Must**: Record findings in mission context (decisions.md, questions.md, or session notes). Don't just stash them.
- **Must**: If all tiers fail, ask the user exactly one question with the context gathered so far.
- **Must not**: Skip tiers in the general case - even obvious answers deserve a quick search for confirmation. The only exceptions are the shortcuts in the table below.
- **Must not**: Use this skill for long systematic literature reviews - use sc-storm for that.
- **Must not**: Use for casual browsing or curiosity - only for blocking technical questions.

## Tier shortcuts

When the topic is immediately recognizable, shortcut to the appropriate tier:

| Scenario | Start at |
|----------|----------|
| Exact error message | Tier 1 (WebSearch) |
| Known package, need latest version | Tier 2 (WebFetch npm/pypi/crates.io, or ctx_fetch_and_index → ctx_search) |
| Framework best practice (React, Go, etc.) | Tier 3 (multi-source WebSearch/WebFetch) |
| "How do I...?" with no specific docs | Tier 1 → escalate as needed |

## Out of scope

- Systematic literature review sessions - use sc-storm
- Debugging discipline - Cursor Debug Mode
- Requirement clarification - use `/sc-discuss` / sc-clarify
- Architectural decisions - use sc-architect

## References

- `WebSearch` tool - Cursor built-in web search
- `WebFetch` tool - Cursor built-in page fetcher
- sc-storm - Tier 3 systematic research feeding discuss (open-domain / strategy; not API gray areas)
- `decisions.md` - record findings that affect mission direction
- `questions.md` - record open questions escalated to the user
