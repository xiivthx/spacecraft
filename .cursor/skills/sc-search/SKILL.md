---
name: sc-search
description: "Quick internet search with 3-tier escalation for stuck issues, gray areas, and stale knowledge. Mission-affecting external facts use Fact-check (see references/fact-check.md)."
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

### Fact-check

Before settling a mission-affecting external fact: `references/fact-check.md`. Emit `Fact-check: corroborated` | `contested: <id>` | `skipped: <reason>`. One `Task(sc-fact-check)` only when that SoT requires it (claim block ≤5). Contested → no auto-pick / implement on that claim.

### Tier 4 - Ask user

One question: what you tried (tiers 1–3) + recommended next step. Do not implement while open.

## Rules

- **Must**: Always try Tier 1 first. Only escalate when the answer is still unclear.
- **Must**: Never search for the same thing twice in one session. Cache results in session context.
- **Must**: Record findings in mission context (decisions.md, questions.md, or session notes). Don't just stash them.
- **Must**: If all tiers fail, ask the user exactly one question with the context gathered so far.
- **Must**: When stating mission-affecting external facts from this skill, leave `Fact-check:` per `references/fact-check.md`.
- **Must not**: Skip tiers in the general case - even obvious answers deserve a quick search for confirmation. The only exceptions are the shortcuts in the table below.
- **Must not**: Use this skill for long systematic literature reviews - use sc-storm for that.
- **Must not**: Use for casual browsing or curiosity - only for blocking technical questions.
- **Must not**: Pass chat or mission trees into `Task(sc-fact-check)`; treat `Fact-check:` as ready/ship authority.

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
- Debugging discipline - `/sc-debug`
- Requirement clarification - use `/sc-discuss` / sc-clarify
- Architectural decisions - use sc-architect

## References

- `WebSearch` tool - Cursor built-in web search
- `WebFetch` tool - Cursor built-in page fetcher
- `references/fact-check.md` - claim cross-check + disposition
- `.cursor/agents/sc-fact-check.md` - critic agent
- sc-storm - open-domain / strategy research for discuss
- `decisions.md` - record findings that affect mission direction
- `questions.md` - record open questions escalated to the user
