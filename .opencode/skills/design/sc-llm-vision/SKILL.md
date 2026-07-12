---
name: sc-llm-vision
description: Use LLM vision models (Gemini via agy CLI) to review UI screenshots for design quality, pixel precision, and visual issues. Activate on visual review, screenshot audit, design QA, pixel check, or UI quality inspection.
license: MIT
compatibility: opencode
metadata:
  version: 1.0
  audience: xiivthx
---

# sc-llm-vision

Use multi-modal LLMs (Gemini 3.5 Flash via `agy` CLI) to perform deep visual design review of UI screenshots. Unlike DOM-based automated QA (which checks computed styles, overflow, accessibility), LLM vision review catches human-perceptible issues: layout feel, color harmony, composition, polish, anti-slop, art direction, and brand consistency.

## Prerequisites

- `agy` CLI installed and `~/.local/bin/agy` available on PATH
- `agy models` should list vision-capable models (e.g., Gemini 3.5 Flash)
- Screenshots must exist on disk (typically from Playwright capture scripts)

## When to use

Activate on:

- **"Review screenshots visually" / "visual QA" / "design review"** — run LLM vision review
- **"Check pixel precision" / "pixel-perfect audit"** — detailed per-screenshot review
- **"Find smallest issues" / "obsessive review"** — fovea (per-screenshot) mode
- After capturing screenshots via Playwright — before claiming complete

## Models

```bash
agy models
# Gemini 3.5 Flash (Medium) — recommended for vision review (best detail/accuracy)
# Gemini 3.5 Flash (High)    — faster but less thorough
# Gemini 3.1 Pro (High)      — slower, architectural analysis
```

Default: `Gemini 3.5 Flash (Medium)`.

## Two Review Modes

### 1. Survey Mode (initial scan, broad)

Send 20-80 screenshots in one batch. Good for first-pass triage — finding systemic patterns and obvious issues.

- **Batch**: 20-80 screenshots
- **Detail**: Surface patterns, ~1-2 findings/screenshot
- **Use when**: First pass on a new app, finding the biggest problems

### 2. Fovea Mode (detail review, per-route)

Send 1-2 routes (4-8 screenshots) per batch. Gemini inspects each screenshot at pixel level. Finds tiny issues the survey misses.

- **Batch**: 4-8 screenshots (1-2 routes × 4 viewports)
- **Detail**: Pixel-level, ~5-10 findings/screenshot, 20-35 per route
- **Use when**: Primary review mode — this is where real value comes from

**Rate limit**: Wait ~30 minutes between **agy** batches to avoid Antigravity usage limits. Running batches back-to-back causes timeouts.

## Workflow

### 0. Prerequisites

```bash
which agy && agy models
```

Confirm Gemini models available. The `agy` CLI uses the Antigravity platform.

### 1. Capture Screenshots

Use Playwright to capture full-page screenshots at 4 viewports (1440, 1024, 768, 320). Save to `.space/screenshots/v{version}/`. Naming: `{num}_{group}_{scenario}_{viewport}.png`.

### 2. Survey Mode (First Pass)

```bash
agy --model="Gemini 3.5 Flash (Medium)" --prompt \
"You are a senior visual/product designer. Review {N} screenshots across {M} routes.

DESIGN REFERENCE: [paste key tokens: colors, typography, layout, avoid-patterns]

FILES: [list all screenshot filenames grouped by route]

For each route evaluate: layout, typography, color, responsive, polish, interaction.

OUTPUT JSON array of findings grouped by route." \
--dangerously-skip-permissions \
--add-dir /path/to/screenshots 2>&1 | tee .space/screenshots/v{version}/_gemini-survey.txt
```

### 3. Fovea Mode (Per-Route Detail)

```bash
agy --model="Gemini 3.5 Flash (Medium)" --prompt \
"You are a senior visual/product designer. Review exactly 4 screenshots of the {ROUTE NAME} at 4 viewports.

DESIGN REFERENCE: [full tokens: colors, sizes, layout rules, component specs]

FILES:
{route}_desktop.png (1440) | {route}_small-desk.png (1024)
{route}_tablet.png (768)   | {route}_mobile.png (320)

For EACH screenshot individually, find EVERY issue, no matter how small:
1. PIXEL PRECISION: alignment, margins, padding, uneven gaps, misaligned elements
2. TYPOGRAPHY: font sizes vs spec, line-heights, weights, overflow, clipping
3. COLOR: exact token matches, contrast, unauthorized colors
4. COMPOSITION: layout execution, hierarchy, right-rail presence
5. RESPONSIVE: breakpoint adaptation, collapse behavior, hidden elements
6. POLISH: anything unfinished, placeholder-ish, dev-artifact, sloppy
7. ANTI-SLOP: decorative blobs, nested cards, fake metrics, AI-template feel
8. ART DIRECTION: premium feel, on-brand, professional warmth

OUTPUT JSON array only." \
--dangerously-skip-permissions \
--add-dir /path/to/screenshots 2>&1 | tee .space/screenshots/v{version}/_gemini-fovea-{route}.txt
```

**Repeat** for each route, waiting 30 minutes between batches.

### 4. Extract and Categorize

Parse JSON output. Categorize by severity:

| Severity | Criteria |
|----------|----------|
| **critical** | Crash, content missing, page unusable, security/brand leak |
| **important** | Design system violated, layout broken, accessibility blocked, data inconsistency |
| **minor** | Pixel misalignment, wrong color token, typography off-spec, polish |
| **info** | Nice-to-have improvements, taste suggestions |

Deduplicate: same issue appearing across multiple routes → 1 systematic issue, not N per-route issues.

### 5. Create GitHub Issues

Group related findings by theme and route:

```bash
cat << 'BODY' | gh issue create --title "...(route+theme)..." --body-file -
## Source
Gemini 3.5 Flash fovea review — {N} screenshots

## Description
...

## Evidence
Screenshots: {filenames}

## Suggested Fix
...

## Verification Criteria
...
BODY
```

**Grouping rules**:
- One route's desktop issues → one issue
- One route's mobile/responsive issues → one issue
- Systematic patterns across routes (e.g., "orange color everywhere") → one cross-cutting issue
- Pre-existing coverage → reference the existing issue instead of duplicating

### 6. Merge with Automated QA

Combine LLM vision findings with Playwright DOM QA (overflow, tokens, hit areas, accessibility). LLM finds what DOM can't — visual feel, composition, brand fit.

### 7. Record Evidence

```
.space/screenshots/v{version}/
  _gemini-survey.txt              # survey output
  _gemini-fovea-{route}.txt       # per-route fovea outputs
  _gemini-review-prompt.md        # prompt copy for audit
  _merged-design-report.json      # compiled findings (optional)
```

## Prompt Engineering

### Always Include DESIGN.md Context

The model must know what "correct" looks like. Include at minimum:

- **Color tokens**: exact hex values with semantic roles
- **Typography scale**: font sizes, line heights, weights, families
- **Layout rules**: grid columns, sidebar widths, breakpoints
- **Component rules**: button sizes, card radii, border styles
- **Anti-patterns**: what to flag (blobs, nested cards, gradients)

### Output Format

Always request JSON array:

```
OUTPUT JSON array only. Each finding:
{
  "screenshot": "filename.png",
  "severity": "critical|important|minor|info",
  "category": "layout|typography|color|interaction|responsive|polish|brand|accessibility",
  "finding": "one-line",
  "element": "which UI element",
  "suggestion": "actionable fix"
}
```

### Be Exhaustive, Not Polite

> "Find EVERY issue no matter how small. Be obsessive. Compare against high-end references."

The model defaults to polite reviews — push toward critical scrutiny.

## Batch Size Guide

| Mode | Screenshots/Batch | Detail/Shot | Wait | Use Case |
|------|-------------------|-------------|------|----------|
| Survey | 20-80 | ~1-2 | None | First-pass triage |
| Fovea | 4-8 | ~5-10 | 30min | Per-route deep review |
| Mixed | 12-16 | ~3-5 | 30min | Quick section review |

**Never batch all screenshots at once** — the model skims and misses most issues.

## Rate Limiting

- **agy (Gemini via Antigravity)**: wait ~30 minutes between fovea batches
- **Symptom of hitting limit**: agent starts tool-calling (search, file reads) instead of analyzing images
- **Timeout handling**: set `--print-timeout` to 5-10 minutes per batch. If batch times out with no output, retry once — if still failing, skip that route or wait longer.

## Anti-Patterns

- Don't send all screenshots in one batch — model skims, misses pixel issues
- Don't omit DESIGN.md from prompt — model can't find violations without knowing the spec
- Don't skip fovea mode for complex routes (itinerary, map, overview, detail planner)
- Don't create 1 issue per finding — group by route + theme
- Don't run without `--dangerously-skip-permissions` — agy prompts for each tool use
- Don't run batches back-to-back — rate limits cause timeouts
- Don't trust the model's judgment on data correctness — verify settlement counts, member counts, currency values independently

## Integration Example

```
/sc-review
  ├── Playwright DOM QA → overflow, tokens, hit area, a11y report
  ├── Gemini Survey (1 batch, 40-80 shots) → systemic patterns
  ├── Gemini Fovea (N batches, 4-8 shots each) → route-level pixel issues
  │     Wait 30min between batches
  └── Merge findings → deduplicate → GitHub issues → review.json
```

Commander orchestrates. This skill is the LLM vision review mechanic.

## Output Format

```
Batch: landing-page (4 screenshots)
Model: Gemini 3.5 Flash (Medium)
Findings: 28 (3 critical, 8 important, 12 minor, 5 info)
Output: .space/screenshots/v0.7.2/_gemini-fovea-landing.txt
GitHub issues: #88, #89

Time: 2min | Wait: 30min until next batch
```

## Checklist

- [ ] `agy` CLI available and authenticated
- [ ] Screenshots exist at expected paths
- [ ] DESIGN.md content included in every prompt
- [ ] Survey mode completed for systemic patterns
- [ ] Fovea mode completed for all priority routes
- [ ] JSON output captured to `_gemini-fovea-{route}.txt` files
- [ ] Findings deduplicated across routes
- [ ] GitHub issues created with severity, evidence, fix suggestions
- [ ] 30-minute waits observed between fovea batches
- [ ] All findings merged with automated QA into unified review

---

## References

- `agy --help` — agy CLI reference
- `agy models` — available vision models
- `DESIGN.md` — project design system (provide in prompt context)
- Google Antigravity CLI docs: `~/.gemini/antigravity-cli/`
