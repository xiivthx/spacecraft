---
name: sc-search
description: Quick internet search with 3-tier escalation for resolving stuck issues, gray areas, and stale knowledge. Activates on unfamiliar errors, deprecated APIs, dependency uncertainty, or technical gray areas.
license: MIT
compatibility: opencode
metadata:
  version: 1.0
  audience: xiivthx
---

# sc-search

Quick internet search with structured escalation. Find answers for stuck issues, gray areas, and stale knowledge — without leaving the development flow.

## When to use

Activate when the Commander encounters any of these triggers:

- **Unfamiliar error or stack trace** — error message or crash that isn't obvious from context
- **Deprecated or unfamiliar API** — method, function, or pattern that looks wrong or outdated
- **Dependency version or compatibility uncertainty** — unsure about latest version, breaking changes, or peer dependencies
- **Technical gray area** — knowledge is stale or missing, best practices may have changed

## Workflow

Use this exact escalation sequence. Never skip a tier unless a shortcut applies (see Tier shortcuts below).

### Tier 1 — Quick search (`google_search`)

First attempt. Fast lookup for direct answers.
- Use `google_search` with a targeted query (error message, API name, package + version).
- Goal: find a relevant page (official docs, GitHub issue, Stack Overflow, blog post).
- Timeout: ~5s. If no clear answer, escalate.

### Tier 2 — Deep read (`webfetch`)

When Tier 1 returns a specific relevant URL.
- Use `webfetch` to fetch the page content.
- Extract the concrete answer (version number, fix, migration path, correct syntax).
- Timeout: ~10s. If the page doesn't resolve the question, escalate.
- Alternative: use `ctx_fetch_and_index` to fetch + index, then `ctx_search` to query — useful when you need to cross-reference multiple sources or re-query the same page.

### Tier 3 — Systematic research (`spacecraft research`)

When Tiers 1–2 are insufficient, or the topic matches a known scope.
- Run `spacecraft research "<query>"` with relevant flags.
- If the topic matches a known scope (React, Go, PostgreSQL, Tailwind, etc.), use `--scope`.
- For deep investigation, use `--deep true`.
- Timeout: ~10s. `spacecraft research` uses Brave Search for scoped, versioned results.

### Tier 4 — Ask user

When all tiers fail to resolve the question.
- Ask exactly one question.
- Include: the question, what you tried (tiers 1–3), and a recommended next step.
- Do not proceed with implementation while the question is open.

## Rules

- **Must**: Always try Tier 1 first. Only escalate when the answer is still unclear.
- **Must**: Never search for the same thing twice in one session. Cache results in session context.
- **Must**: Record findings in mission context (decisions.md, questions.md, or session notes). Don't just stash them.
- **Must**: If all tiers fail, ask the user exactly one question with the context gathered so far.
- **Must not**: Skip tiers in the general case — even obvious answers deserve a quick search for confirmation. The only exceptions are the shortcuts in the table below.
- **Must not**: Use this skill for systematic literature reviews or multi-source synthesis — use `/sc-research` or `spacecraft research` directly.
- **Must not**: Use for casual browsing or curiosity — only for blocking technical questions.

## Tier shortcuts

When the topic is immediately recognizable, shortcut to the appropriate tier:

| Scenario | Start at |
|----------|----------|
| Exact error message | Tier 1 (google_search) |
| Known package, need latest version | Tier 2 (webfetch npm/pypi/crates.io, or ctx_fetch_and_index → ctx_search) |
| Framework best practice (React, Go, etc.) | Tier 3 (spacecraft research --scope) |
| "How do I...?" with no specific docs | Tier 1 → escalate as needed |

## Out of scope

- Systematic literature review — use `/sc-research` or `spacecraft research` directly
- Debugging discipline — use sc-debug
- Requirement clarification — use sc-clarify
- Architectural decisions — use sc-architect

## References

- `spacecraft research --help` — research subcommand flags and usage
- `google_search` tool — OpenCode built-in web search
- `webfetch` tool — OpenCode built-in page fetcher
- `decisions.md` — record findings that affect mission direction
- `questions.md` — record open questions escalated to the user
