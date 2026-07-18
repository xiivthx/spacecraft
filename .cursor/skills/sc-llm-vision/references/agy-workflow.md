> Consult when: running survey or fovea batches with `agy`, choosing models, or hitting rate limits.

# agy workflow

## Models

```bash
agy models
# Gemini 3.5 Flash (Medium) - recommended (best detail/accuracy)
# Gemini 3.5 Flash (High)    - faster, less thorough
# Gemini 3.1 Pro (High)      - slower, architectural analysis
```

Default: `Gemini 3.5 Flash (Medium)`.

## Two review modes

### Survey (initial scan)

- **Batch**: 20-80 screenshots
- **Detail**: ~1-2 findings/screenshot, systemic patterns
- **Use when**: First pass on a new app

### Fovea (per-route detail)

- **Batch**: 4-8 screenshots (1-2 routes × 4 viewports)
- **Detail**: Pixel-level, ~5-10 findings/screenshot
- **Use when**: Primary review mode - where real value comes from

**Rate limit**: Wait ~30 minutes between **agy** batches. Back-to-back batches cause timeouts.

## Survey prompt shape

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

## Fovea prompt shape

```bash
agy --model="Gemini 3.5 Flash (Medium)" --prompt \
"You are a senior visual/product designer. Review exactly 4 screenshots of the {ROUTE NAME} at 4 viewports.

DESIGN REFERENCE: [full tokens: colors, sizes, layout rules, component specs]

FILES:
{route}_desktop.png (1440) | {route}_small-desk.png (1024)
{route}_tablet.png (768)   | {route}_mobile.png (320)

For EACH screenshot individually, find EVERY issue, no matter how small:
1. PIXEL PRECISION: alignment, margins, padding, uneven gaps
2. TYPOGRAPHY: sizes vs spec, line-heights, weights, overflow
3. COLOR: token matches, contrast, unauthorized colors
4. COMPOSITION: layout execution, hierarchy
5. RESPONSIVE: breakpoint adaptation, collapse behavior
6. POLISH: unfinished, placeholder-ish, sloppy
7. ANTI-SLOP: blobs, nested cards, fake metrics, AI-template feel
8. ART DIRECTION: premium feel, on-brand

OUTPUT JSON array only." \
--dangerously-skip-permissions \
--add-dir /path/to/screenshots 2>&1 | tee .space/screenshots/v{version}/_gemini-fovea-{route}.txt
```

Repeat per route; wait 30 minutes between batches.

## Prompt engineering

Always include DESIGN.md context: color tokens, typography scale, layout rules, component rules, anti-patterns.

Request JSON array findings:

```
{
  "screenshot": "filename.png",
  "severity": "critical|important|minor|info",
  "category": "layout|typography|color|interaction|responsive|polish|brand|accessibility",
  "finding": "one-line",
  "element": "which UI element",
  "suggestion": "actionable fix"
}
```

Push exhaustive scrutiny: "Find EVERY issue no matter how small."

## Batch size guide

| Mode | Screenshots/Batch | Detail/Shot | Wait | Use Case |
|------|-------------------|-------------|------|----------|
| Survey | 20-80 | ~1-2 | None | First-pass triage |
| Fovea | 4-8 | ~5-10 | 30min | Per-route deep review |
| Mixed | 12-16 | ~3-5 | 30min | Quick section review |

Never batch all screenshots at once.

## Rate limiting

- Wait ~30 minutes between fovea batches
- Symptom of hitting limit: agent tool-calls (search, file reads) instead of analyzing images
- Timeout: set `--print-timeout` to 5-10 minutes; retry once, then skip route or wait longer
- Always pass `--dangerously-skip-permissions` or agy prompts per tool use

## Spacecraft integration

Store outputs under `.space/screenshots/v{version}/`. Reference findings from mission `review.json` / GitHub issues. Capture verification with `spacecraft evidence` when the plan requires it.
