---
name: sc-research
description: "Runs systematic research when a mission needs current documentation, web evidence, or deeper technical analysis."
disable-model-invocation: true
---

Use sc-mission, sc-verification, and sc-search.
Run systematic research for: $ARGUMENTS

## When to use

Use `/sc-research <query>` for systematic investigation when:
- Quick search (sc-search) is insufficient
- You need scoped, versioned results from Brave Search
- You need deep analysis (browser-use or NotebookLM)
- You want structured, saved research output for a mission

Not for quick lookups - use sc-search auto-trigger for that.

## Pre-flight checks

Resolve the mission. Block if unsafe.

Run:
```
spacecraft git-info
```

### API key check

`spacecraft research` requires `SPACECRAFT_BRAVE_API_KEY` in the environment. Check before running:

```
spacecraft research "test" --json --no-save 2>&1
```

If the output contains "SPACECRAFT_BRAVE_API_KEY", warn the user:
```
SPACECRAFT_BRAVE_API_KEY is not set. Set it in your environment to use Brave Search.
export SPACECRAFT_BRAVE_API_KEY="your-key"
```
Exit cleanly with exit code 1. Do not proceed without the key.

## Workflow

### 1. Resolve mission

```
spacecraft resolve --json
```
Block if safety ≠ `safe`.

### 2. Parse user flags

Extract flags from `$ARGUMENTS`:
- `--scope <name>` - force a scope (react, tailwindcss, nextjs, storybook, postgresql, go, rust)
- `--deep true|nlm` - deep research (true=browser-use, nlm=notebooklm)
- `--results <n>` - number of results (default 5)
- `--timeout <d>` - request timeout (default 10s)
- `--json` - JSON output
- `--no-save` - skip persistence

The query is everything in `$ARGUMENTS` that isn't a flag.

### 3. Run research

```
spacecraft research "<query>" [user-provided flags]
```

### 4. Save results

Save results to `.space/missions/<id>/research/`:
```
mkdir -p .space/missions/<id>/research/
```

If using `--no-save`, skip persistence but still present results.

### 5. Present results

Format and present results inline. Include:
- Query and scope used
- Number of results
- Key findings (summarized)
- Source URLs

### 6. Record evidence

```
spacecraft evidence "research-<query-slug>" -- spacecraft research "<query>" [flags]
```

## Research auto-trigger

This command wraps `spacecraft research`, which uses Brave Search for scoped, versioned documentation searches. For quick lookups during implementation, prefer the sc-search skill (auto-triggered). Use `/sc-research` when you need systematic, saved research output.

## Hard stop gates

- Resolver conflict or ambiguity
- Missing `SPACECRAFT_BRAVE_API_KEY`
- Query is empty or only flags

## Error handling

- If `SPACECRAFT_BRAVE_API_KEY` is not set, warn and exit code 1
- If the query is empty, ask the user for a query
- If research returns no results, report and suggest alternative query
- Do not implement product code from this command - it's research-only

## Edge cases

- **No active mission** - Run without mission binding; save to `.space/research/` instead
- **Research directory exists** - Append with timestamp to avoid overwrite
- **User provides invalid scope** - List valid scopes and ask to retry
- **Network error** - Report the error, suggest retry or fallback to sc-search

End with research summary, evidence id, and session advice.
