---
name: sc-llm-vision
description: "Use LLM vision models (Gemini via agy CLI) to review UI screenshots for design quality, pixel precision, and visual issues. Activate on visual review, screenshot audit, design QA, pixel check, or UI quality inspection."
---

# sc-llm-vision

Use multi-modal LLMs (Gemini 3.5 Flash via `agy` CLI) for deep visual design review of UI screenshots. Catches human-perceptible issues DOM QA misses: layout feel, color harmony, composition, polish, anti-slop, art direction, brand consistency.

## Prerequisites

- `agy` CLI installed (`which agy && agy models`)
- Vision-capable models listed (default: `Gemini 3.5 Flash (Medium)`)
- Screenshots on disk (typically Playwright → `.space/screenshots/v{version}/`)

## When to use

- **"Review screenshots visually" / "visual QA" / "design review"**
- **"Check pixel precision" / "pixel-perfect audit"**
- **"Find smallest issues" / "obsessive review"** - fovea mode
- After Playwright capture, before claiming UI complete

## Workflow

1. **Confirm tooling** - `which agy && agy models`
2. **Capture screenshots** - 4 viewports (1440, 1024, 768, 320); naming `{num}_{group}_{scenario}_{viewport}.png`
3. **Survey mode** - 20-80 screenshots, first-pass systemic patterns (see `references/agy-workflow.md`)
4. **Fovea mode** - 4-8 screenshots per route, pixel-level review; wait ~30 min between agy batches
5. **Extract and categorize** - critical / important / minor / info; dedupe cross-route duplicates
6. **Create GitHub issues** - group by route + theme (see `references/review-checklist.md`)
7. **Merge with DOM QA** - combine with Playwright overflow/token/a11y findings
8. **Record evidence** under `.space/screenshots/v{version}/` (`_gemini-survey.txt`, `_gemini-fovea-{route}.txt`, optional merged report)

Commander orchestrates. This skill is the LLM vision mechanic. Detailed prompts, batch sizes, rate limits, and anti-patterns live in references.

## Rules

- **Must**: Include DESIGN.md tokens in every prompt (colors, type, layout, anti-patterns)
- **Must**: Prefer fovea for priority routes; survey alone is not enough for ship claims
- **Must**: Wait ~30 minutes between fovea batches (Antigravity rate limits)
- **Must**: Request JSON-only findings; group GitHub issues by route + theme
- **Must not**: Send all screenshots in one batch (model skims)
- **Must not**: Trust model judgment on data correctness - verify counts/currency independently
- **Must not**: Create one issue per finding

## Out of scope

- DOM/computed-style QA - Playwright / automated a11y tooling
- Implementing UI fixes - sc-coder / build lane
- Art direction from scratch without screenshots - sc-design / sc-designer

## Output format

```
Batch: <route> (N screenshots)
Model: Gemini 3.5 Flash (Medium)
Findings: N (critical / important / minor / info)
Output: .space/screenshots/v{version}/_gemini-fovea-{route}.txt
GitHub issues: #N, ...
Time: ~2min | Wait: 30min until next batch
```

## Checklist

- [ ] `agy` available and authenticated
- [ ] Screenshots exist at expected paths
- [ ] DESIGN.md included in every prompt
- [ ] Survey completed for systemic patterns
- [ ] Fovea completed for priority routes
- [ ] JSON captured to `_gemini-fovea-{route}.txt`
- [ ] Findings deduplicated; issues grouped by route + theme
- [ ] 30-minute waits between fovea batches
- [ ] Merged with automated QA into unified review

## References

- `references/agy-workflow.md` - survey/fovea prompts, models, batch sizes, rate limits
- `references/review-checklist.md` - severity, issue grouping, anti-patterns, evidence layout
- `DESIGN.md` - project design system (provide in prompt context)
- `agy --help` / `agy models` - CLI and model list
