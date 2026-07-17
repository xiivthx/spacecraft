---
name: sc-memory
description: "Wraps ctx_search and ctx_index with spacecraft conventions for structured memory across missions. Activate on cross-mission recall, lesson lookup, session context recovery, or \"memory\", \"what do we know about\", \"prior mission context\", \"search memory\", \"recall\"."
---

# sc-memory

Wraps context-mode ctx_search and ctx_index with spacecraft conventions. No new storage layer -- this is conventions and lifecycle integration over the existing FTS5 knowledge base.

## When to use

Activate when the user asks to:

- **"What do we know about X?" / "search memory" / "recall"** -- cross-mission knowledge lookup
- **"Prior mission context" / "what happened last session"** -- session handoff recovery
- **"Record this for later" / "remember this"** -- index a fact or artifact
- Any agent or skill that needs cross-mission recall or structured memory access

## Conventions

### Source label format

All ctx_index calls use this label format:

```
sc-memory/<mission-id>/<type>
```

Where `<type>` is one of: `spec`, `plan`, `decisions`, `questions`, `issues`, `solved`, `learned`, `review`, `evidence`.

Example: `sc-memory/M07PYRGLG/spec`

### Querying

Use ctx_search with these conventions:

```sh
# Scoped to a mission
ctx_search(queries: ["<search terms>"], source: "sc-memory/<mission-id>")

# Scoped to a type across all missions
ctx_search(queries: ["<search terms>"], source: "sc-memory/spec")

# Timeline-sorted for session recall (replaces non-existent ctx_memory)
ctx_search(queries: ["last user prompt", "open blockers", "active skills"], sort: "timeline")
```

Timeline-sorted ctx_search is the canonical replacement for ctx_memory -- ctx_memory does not exist as a tool. Timeline sort surfaces auto-captured session events (decisions, errors, blockers, plans, user prompts, rejected approaches) alongside indexed content.

### Memory pipeline

```
Write artifact → ctx_index(source: "sc-memory/<mission-id>/<type>") → ctx_search retrieves it
```

## Degradation

When context-mode is unavailable:

1. Try ctx_search / ctx_index first
2. On failure: warn once ("context-mode unavailable, falling back to file reads")
3. Read the corresponding markdown file directly from `.space/missions/<id>/`
4. Continue -- missing memory is not a blocker

## Integration points

| Integration | What | Where |
|-------------|------|-------|
| sc-mission | ctx_index spec/plan/decisions/questions on create/update | Workflow step 2 |
| sc-learn | ctx_index issues.md/solved.md/learned.md after writes | After each write step |
| sc-resume | ctx_search for prior mission context | Pre-flight checks, handoff resume |

## Rules

- **Must**: Use source label format `sc-memory/<mission-id>/<type>` for all ctx_index calls.
- **Must**: ctx_index is best-effort -- warn on failure, never block the mission lifecycle.
- **Must**: ctx_search with `sort: "timeline"` replaces ctx_memory for session recall.
- **Must**: Fall back to file reads when context-mode tools fail.
- **Must not**: Create new storage layers -- ctx_index/ctx_search are the single source of truth.
- **Must not**: Block mission progress on memory operations.

## Out of scope

- New storage layer (`sc-memory` is conventions only)
- Custom summarizer or retrieval engine
- Replacing context-mode -- this is wiring, not building
- Agent memory without context-mode present

## Checklist

Before claiming memory operations are correct:

- [ ] Source labels use `sc-memory/<mission-id>/<type>` format
- [ ] ctx_index calls are best-effort with warn-on-failure
- [ ] ctx_search queries use mission-scoped sources when appropriate
- [ ] Timeline sort used for session recall (not guessing a non-existent tool)
- [ ] Degradation path documented: fall back to file reads
- [ ] No new storage layer introduced

---

## References

- context-mode tools: `ctx_index`, `ctx_search`, `ctx_execute`, `ctx_memory`
- `.space/missions/<id>/` -- mission artifact files (fallback when ctx unavailable)
- sc-learn -- lesson and issue capture (indexes via sc-memory conventions)
- sc-mission -- artifact lifecycle (indexes via sc-memory conventions)
- sc-resume -- session handoff (queries via sc-memory conventions)
