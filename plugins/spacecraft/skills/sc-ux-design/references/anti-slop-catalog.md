# Anti-Slop Catalog

All 46 patterns from [impeccable.style/slop](https://impeccable.style/slop), organized by category with detection method and fix guidance.

Detection methods:
- **CLI** - Deterministic via `npx impeccable detect` on files
- **CLI opt-in** - Deterministic but off by default; enable with `--gpt` or `--gemini`
- **Browser** - Deterministic but needs real browser layout (Playwright or extension)
- **LLM only** - Heuristic review by agent during design critique

## Visual Details (7 rules)

| # | Rule | Type | Detection | Fix |
|---|------|------|-----------|-----|
| 1 | **Border accent on rounded element** - Thick accent border on a rounded card clashes with the border-radius | AI slop | CLI | Remove border or border-radius |
| 2 | **Glassmorphism everywhere** - Blur effects, glass cards, glow borders used as decoration rather than layering solution | AI slop | LLM only | Use glass only when it solves a real z-index/layering problem |
| 3 | **Side-tab accent border** - Thick colored border on one side of a card. Most recognizable AI UI tell | AI slop | CLI | Use subtler accent or remove entirely |
| 4 | **Hairline border with wide shadow** - 1px hairline paired with wide diffuse shadow. Commit to one: edge or elevation | AI slop | CLI opt-in (`--gpt`) | Pick defined edge OR soft shadow, not both |
| 5 | **Repeating-gradient stripes** - Repeating-gradient used as surface decoration | AI slop | CLI opt-in (`--gpt`) | Use deliberate texture or leave surface plain |
| 6 | **Extreme border-radius on cards** - 24px+ radius on small cards rounds everything into same blob. Cards: 12–16px max | AI slop | LLM only | Cap card radius at 12–16px; full pill for tags/buttons only |
| 7 | **Amateurish hand-drawn SVG** - Hand-coded SVG illustrations read as doodles, not whimsy | AI slop | LLM only | Ship real assets or no illustration |

## Typography (10 rules)

| # | Rule | Type | Detection | Fix |
|---|------|------|-----------|-----|
| 8 | **Flat type hierarchy** - Font sizes too close together, no visual hierarchy | AI slop | CLI | Fewer sizes, ≥1.25 ratio between steps |
| 9 | **Icon tile stacked above heading** - Rounded-square icon container above heading: universal AI feature-card template | AI slop | CLI | Side-by-side icon and heading, or icon in flow without container |
| 10 | **Italic serif display headline** - Oversized italic serif as hero headline | AI slop | CLI | Set roman, or use non-serif display face. Editorial context may justify |
| 11 | **Hero eyebrow / pill chip** - Tiny uppercase letter-spaced label above hero headline, or pill chip variant | AI slop | CLI | Drop eyebrow, fold kicker into headline, or use breadcrumb |
| 12 | **Repeated section kicker labels** - Tiny uppercase tracked labels above every section heading | AI slop | CLI | Replace with structure, artifacts, imagery, or brand system |
| 13 | **Oversized hero headline** - Full-sentence headline at display size dominates viewport, no room above fold | AI slop | CLI | Set long headlines smaller, or tighten copy to 1–2 words |
| 14 | **Crushed letter spacing** - Letter-spacing tighter than legibility threshold | AI slop | CLI | Tighten display type optically, not destructively |
| 15 | **Overused font** - Inter, Geist, Space Grotesk, Instrument Serif: every AI wave converges on same faces | AI slop | CLI | Choose a distinctive face with personality |
| 16 | **Single font for everything** - Only one font family across entire page, no heading/body contrast | AI slop | CLI | Pair distinctive display font with refined body font |
| 17 | **All-caps body text** - Long passages in uppercase destroy word-shape recognition | Quality | CLI | Reserve uppercase for short labels and headings |

## Color & Contrast (5 rules)

| # | Rule | Type | Detection | Fix |
|---|------|------|-----------|-----|
| 18 | **AI color palette** - Purple/violet gradients and cyan-on-dark are most recognizable AI tells | AI slop | CLI | Choose distinctive, intentional palette |
| 19 | **Dark mode with glowing accents** - Dark backgrounds with colored box-shadow glows | AI slop | CLI | Use subtle purposeful lighting, or skip dark theme |
| 20 | **Gradient text** - Decorative gradient on headings and metrics, kills scannability | AI slop | CLI | Use solid colors for text |
| 21 | **Gray text on colored background** - Washed out, hard to read | Quality | CLI | Use darker shade of background color or white/near-white |
| 22 | **Cream / beige palette** - Warm cream or beige background, default "tasteful" AI surface | AI slop | CLI | Choose background from deliberate palette, not safe warm off-white |

## Layout & Space (8 rules)

| # | Rule | Type | Detection | Fix |
|---|------|------|-----------|-----|
| 23 | **Hero metric layout** - Big number, small label, 3 supporting stats, gradient accent. Used everywhere, trusted nowhere | AI slop | LLM only | Use only when metrics are real and significant |
| 24 | **Identical card grids** - Same-sized cards with icon+heading+text repeated endlessly | AI slop | LLM only | Vary layouts, use asymmetric grids, highlight items differently |
| 25 | **Monotonous spacing** - Same spacing value everywhere, no rhythm | AI slop | CLI | Tight groupings for related items, generous separation between sections |
| 26 | **Nested cards** - Cards inside cards, excessive depth and visual noise | AI slop | CLI | Flatten hierarchy: spacing, typography, dividers instead of nesting |
| 27 | **Numbered section markers (01 / 02 / 03)** - Display markers as editorial scaffold | AI slop | CLI | Numbers only when section IS a sequence; otherwise drop |
| 28 | **Line length too long** - Text lines wider than ~80 chars, eye loses place | Quality | Browser | Max-width 65–75ch on text containers |
| 29 | **Content overflowing its container** - Content spills past box bounds | Quality | Browser | Let text wrap, constrain widths, or give deliberate scroll affordance |
| 30 | **Positioned child clipped by overflow container** - overflow:hidden wraps positioned child, clips menus/tooltips | Quality | Browser | Let overflow be visible, or move positioned layer outside clip |

## Motion (3 rules)

| # | Rule | Type | Detection | Fix |
|---|------|------|-----------|-----|
| 31 | **Bounce or elastic easing** - Dialog springs in with overshoot, card bounces | AI slop | CLI | Ease-out-quart/quint/expo for UI; spring for physical objects only |
| 32 | **Layout property animation** - Animating width/height/padding/margin causes layout thrash | Quality | CLI | Use transform and opacity; grid-template-rows for height |
| 33 | **Image hover transform** - Scaling or rotating image on hover, recurring AI signature | AI slop | CLI opt-in (`--gpt`) | Let imagery sit still, or use subtler purposeful interaction |

## Copy (4 rules)

| # | Rule | Type | Detection | Fix |
|---|------|------|-----------|-----|
| 34 | **Em-dash overuse** - Multiple em-dashes in body copy, AI cadence tell | AI slop | CLI | Commas, colons, periods, or parentheses instead |
| 35 | **Marketing buzzword** - Streamline, empower, supercharge, world-class, enterprise-grade | AI slop | CLI | Specific verb+noun saying what the product literally does |
| 36 | **Aphoristic-cadence copy** - Sections landing on short rebuttal or manufactured-contrast aphorism | AI slop | CLI | Once is fine; repeated pattern is the tell. Write naturally |
| 37 | **Theater framing copy** - Dismissing something as "theater" is recurring AI copy tic | AI slop | CLI opt-in (`--gpt`) | Say plainly what the thing does or does not do |

## Imagery (1 rule)

| # | Rule | Type | Detection | Fix |
|---|------|------|-----------|-----|
| 38 | **Broken or placeholder image** - `<img>` with empty/missing src or placeholder values | Quality | CLI | Use real images, generated assets, or remove the tag |

## General Quality (8 rules)

| # | Rule | Type | Detection | Fix |
|---|------|------|-----------|-----|
| 39 | **Cramped padding** - Text too close to container edge | Quality | Browser | ≥8px, ideally 12–16px padding inside bordered/colored containers |
| 40 | **Body text touching viewport edge** - Body paragraphs flush against viewport edge, no padding | Quality | Browser | Container with ≥16px horizontal padding, or max-width with mx-auto |
| 41 | **Justified text** - Justified text without hyphenation creates rivers of whitespace | Quality | CLI | `text-align: left` for body; `hyphens: auto` if must justify |
| 42 | **Low contrast text** - Text fails WCAG AA (4.5:1 body, 3:1 large) | Quality | CLI | Increase contrast between text and background |
| 43 | **Skipped heading level** - h1 then h3 with no h2 breaks document outline | Quality | CLI | No heading level skipping |
| 44 | **Tight line height** - Line-height below 1.3x font size | Quality | CLI | 1.5–1.7 for body text |
| 45 | **Tiny body text** - Body text below 12px | Quality | CLI | ≥14px for body, 16px ideal |
| 46 | **Wide letter spacing on body text** - Letter-spacing >0.05em on body disrupts reading | Quality | CLI | Reserve wide tracking for short uppercase labels only |

## Summary

| Category | Rules | Deterministic (CLI) | Browser-only | LLM-only | Opt-in |
|----------|-------|---------------------|--------------|----------|--------|
| Visual Details | 7 | 2 | 0 | 3 | 2 |
| Typography | 10 | 10 | 0 | 0 | 0 |
| Color & Contrast | 5 | 5 | 0 | 0 | 0 |
| Layout & Space | 8 | 3 | 3 | 2 | 0 |
| Motion | 3 | 2 | 0 | 0 | 1 |
| Copy | 4 | 3 | 0 | 0 | 1 |
| Imagery | 1 | 1 | 0 | 0 | 0 |
| General Quality | 8 | 6 | 2 | 0 | 0 |
| **Total** | **46** | **32** | **5** | **5** | **4** |

**Detection coverage**: 41 rules have deterministic detection (32 CLI + 5 Browser + 4 CLI opt-in). 5 rules are LLM-only heuristic guidelines. Source: [impeccable.style/slop](https://impeccable.style/slop), retrieved 2026-07-10.
