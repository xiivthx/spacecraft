---
name: sc-research
description: "Explicit systematic research via /sc-research. Thin wrapper around sc-search tiers for deeper multi-source investigation."
disable-model-invocation: true
---

Use sc-mission, sc-verification, and sc-search.
Run systematic research for: $ARGUMENTS

## When to use

Use `/sc-research <query>` for systematic investigation when:
- Quick auto-triggered sc-search is insufficient
- You need multi-source synthesis across docs, issues, and changelogs
- You want structured, saved research notes for a mission

Not for quick lookups - use sc-search auto-trigger for that.

## Workflow

### 1. Resolve mission

```
spacecraft resolve
```
On conflict or ambiguity, use `spacecraft use <selector>`.

### 2. Delegate to sc-search tiers

Follow sc-search escalation for `$ARGUMENTS`:
1. **Tier 1** - `WebSearch` with a targeted query
2. **Tier 2** - `WebFetch` on the best URLs (or `ctx_fetch_and_index` → `ctx_search`)
3. **Tier 3** - Deeper multi-source synthesis: additional `WebSearch`/`WebFetch` passes across official docs, release notes, and issue trackers; cross-check versions and contradictions

### 3. Save results

Save notes to `.space/missions/<id>/research/`:
```
mkdir -p .space/missions/<id>/research/
```

If no active mission, save to `.space/research/` instead.

### 4. Present results

Format and present results inline. Include:
- Query used
- Key findings (summarized)
- Source URLs
- Open contradictions or remaining unknowns

### 5. Record evidence

```
spacecraft evidence "research-<query-slug>" -- echo "Research complete: <one-line summary>"
```

## Hard stop gates

- Resolver conflict or ambiguity (resolve via `spacecraft resolve`; on conflict use `spacecraft use <selector>`)
- Query is empty

## Error handling

- If the query is empty, ask the user for a query
- If search returns no useful results, report and suggest an alternative query
- Do not implement product code from this command - research only

## Edge cases

- **No active mission** - Run without mission binding; save to `.space/research/` instead
- **Research directory exists** - Append with timestamp to avoid overwrite
- **Network error** - Report the error; retry Tier 1–2 or ask the user

End with research summary, evidence id, and session advice.
