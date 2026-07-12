---
name: sc-performance
description: >
  Performance review discipline. Activate on N+1 detection, memory leak check, bundle size analysis, render optimization, or performance bottleneck review.
license: MIT
compatibility: opencode
metadata:
  version: 1.0
  audience: xiivthx
---

# sc-performance

Find expensive work before it ships. Scan code, plans, and evidence for performance anti-patterns. Report findings grouped by category. Fix when the current task scope allows; flag otherwise.

## When to use

Activate on these triggers:

- User asks for a performance review or optimization
- Detecting N+1 queries, lazy loading loops, or missing eager loading
- Investigating memory leaks or unbounded growth
- Reviewing bundle size or dependency imports
- React re-render issues: memo, useMemo/useCallback, context cascades
- Synchronous blocking work, hot path allocations, or excessive DOM churn

## Workflow

Use this sequence during implementation or review:

1. **Scope scan** — Read the current task's `plan.json`, `spec.md`, and `evidence.jsonl`. Only flag issues tied to changed code or accepted risk.
2. **Static pattern scan** — Walk the diff for the categories below. Look at call sites, not just declarations.
3. **Quantify if possible** — Prefer measurements over guesses. If a benchmark, bundle analyzer output, or memory profile exists, reference it in the finding.
4. **Classify and report** — Group every finding by category and severity.
5. **Propose a fix** — Give one concrete remediation per finding. Do not block on hypothetical optimizations outside the task scope.

## Rules

### Must

- **N+1 detection**: Flag any loop that triggers lazy-loaded relations, repeated queries, or per-iteration network calls.
- **Memory leak check**: Flag unbounded collections, uncleared timers, unsubscribed listeners, WebSocket leaks, closure captures of large objects, and accumulating global state.
- **Bundle size analysis**: Flag whole-library imports that could be tree-shaken, missing dynamic imports, unoptimized assets, and duplicate dependencies.
- **Render optimization**: Flag inline function/object props, missing memoization on expensive pure components, context value churn, state lifted too high, and unstable `key={index}`.
- **Hot path discipline**: Flag synchronous blocking work on the main thread, excessive DOM manipulation, deep clones, JSON.parse/stringify in render, unnecessary allocations, and uncached regex.
- Report findings with file, line, and a concrete example snippet.

### Must not

- Optimize without a measured or clearly observable cost.
- Add premature abstraction in the name of performance.
- Change unrelated code outside the task's scope.
- Ignore platform-level evidence when `evidence.jsonl` already contains a benchmark.

### Prefer

- Eager loading, batch queries, or data-loader patterns over per-item fetches.
- `React.memo`, `useMemo`, and `useCallback` on genuinely expensive subtrees, not every component.
- Dynamic imports for routes or large rarely-used features.
- Specific function imports over default imports of entire libraries.
- Caching computed values outside hot loops.

## Out of scope

- Functional correctness — separate concern
- Security vulnerabilities — separate concern
- Accessibility review — separate concern
- Infrastructure scaling, CDN tuning, database server configuration
- Code style and SOLID violations — use the appropriate discipline

## Output format

```
Performance scan: <scope>
Severity: [critical / high / medium / low / note]

N+1 / query:
  <file:line> — <pattern> — <fix>

Memory:
  <file:line> — <pattern> — <fix>

Bundle:
  <file:line> — <pattern> — <fix>

Render:
  <file:line> — <pattern> — <fix>

Hot path:
  <file:line> — <pattern> — <fix>

Verdict: [clean / fix-now / fix-in-follow-up / measure-first]
```

## Checklist

Before claiming performance work done:

- [ ] N+1 loops and missing eager loading checked
- [ ] Timers, listeners, subscriptions, and unbounded collections checked for leaks
- [ ] Imports reviewed for tree-shaking and dynamic-import opportunities
- [ ] React memoization, context values, keys, and state placement reviewed
- [ ] Hot paths checked for blocking work, deep clones, JSON serialization, and allocations
- [ ] Findings reference `plan.json` scope or `spec.md` acceptance criteria
- [ ] Any measured evidence from `evidence.jsonl` is cited
- [ ] Verdict and severity assigned for each finding

## References

- `references/performance-patterns.md` — N+1 fixes, batching, eager loading, data loaders
- `references/memory-leak-patterns.md` — cleanup patterns for timers, listeners, subscriptions, closures
- `references/bundle-analysis.md` — tree-shaking, dynamic imports, dependency duplication
- `references/react-render-patterns.md` — memoization, context splitting, stable keys, state colocation
- `references/hot-path-patterns.md` — avoiding blocking work, allocations, and serialization in tight loops
