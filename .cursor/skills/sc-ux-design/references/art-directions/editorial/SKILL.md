---
name: editorial
description: "Magazine-like art direction: serif/sans pairing, clear hierarchy, locked section/layout pool. Optional draft pack under sc-ux-design."
---

# Editorial art direction

Magazine hierarchy for discuss draft HTML: expressive serif display, refined sans for UI chrome and short labels, locked section structures. Not a broadsheet/newspaper skin and not the swiss-grid hairline identity.

Load this pack **after** `references/shared-draft-directives.md` and **before** the approved brief / content tail.

## Authority

This pack **defers** to `references/anti-slop-catalog.md` on conflict: anti-slop-catalog is authoritative (wins). Do not ship catalog-banned house style - no Inter-as-sole-font (or other overused faces alone), no purple/cream defaults, no soft-shadow + hairline as elevation recipe - even if a layout below feels "editorial."

## Iron rules

These iron rules are mandatory for every draft that selects this pack.

1. **Magazine hierarchy** - Few type sizes with a clear ratio (≥1.25 between steps). Display headline, dek/subhead, body, and captions must read as distinct levels. No flat same-size stacks.
2. **Serif / sans pairing** - Serif for long reading and primary display; sans for nav, labels, captions, and compact UI. Never a single family for the whole page. Never use the catalog's overused faces as the house pair.
3. **Locked section pool** - Pick structure only from the **named layouts** below. Do not invent a new page architecture each draft run. Compose a page from one primary layout plus optional supporting layouts from the same pool.
4. **Content fidelity** - Brief and user copy are source of truth. Do not invent metrics, pull quotes, bylines, or brand lines the brief did not supply. If a layout needs a quote or caption and none was given, use a neutral placeholder labeled as such - or choose a layout that does not require it.
5. **Not broadsheet-as-default** - Magazine rhythm and hierarchy, not dense multi-column newspaper with hairline rules and zero radius as the only look (that identity belongs to `swiss-grid`). Soft radii and quiet rules are fine when they serve reading; do not mimic swiss-grid chrome.
6. **Anti-slop deference** - No hero eyebrow / pill chips, no repeating uppercase section kickers, no decorative `01 / 02 / 03` scaffolds, no identical icon+heading card grids, no purple house palette, no beige / warm off-white surface default. Cards only when they hold a real interaction.

## Locked layout / section pool

Use these **named layouts** only. Name the chosen layout(s) in the draft comment or a `data-layout` attribute when helpful.

### 1. Cover Spread

Full-bleed (or edge-to-edge) visual plane behind or beside a short masthead wordmark, one display headline, one dek, one CTA group. First viewport stays lean - no stats strip, schedule, or promo cluster.

### 2. Feature Opener

Article start: serif headline, sans dek, byline/meta row, then a lead paragraph that is visually heavier than body. Optional single dominant image under the opener - not a card collage.

### 3. Reading Columns

Long-form body in a constrained measure (about 65-75ch). One or two columns max. Drop cap or bold lead sentence allowed; justified text only with hyphenation - prefer left-aligned body.

### 4. Pull Quote Band

Full-width band interrupting the reading flow: oversized serif quote, sans attribution. One quote per band; no floating badge chrome on imagery.

### 5. Image Caption Strip

Asymmetric photo (or color field stand-in) with a caption/credit rail beside or below. Caption in sans; no icon tile above a heading.

### 6. Contents Rail

Editorial table of contents or chapter list: sans labels, serif titles, quiet separators. Not a dashboard card grid. No numbered display markers unless the brief is literally a sequence.

### 7. Interview Stack

Q in sans (short, medium weight); A in serif (reading size). Vertical rhythm with generous separation between pairs - not nested cards.

### 8. Colophon Close

Closing credits, next-story teaser, or small-type legal/meta. Compact sans; no second hero.

## Type and color defaults (pack-level, overridable by brief)

- **Suggested pair** (system stacks; brief may replace): display/body serif - `"Iowan Old Style", "Palatino Linotype", Palatino, Georgia, serif`; UI/caption sans - `"Avenir Next", Avenir, "Segoe UI", Candara, sans-serif`.
- **Suggested tokens** when brief is silent: cool paper `--bg: #F5F6F8`, ink `--text: #1A1F24`, muted `--muted: #5A6570`, accent `--accent: #0B6E4F` (forest). Do **not** use purple/violet or warm off-white / beige paper as pack defaults.
- Spacing: generous between sections, tighter within a related group (anti-monotonous spacing).

## Example

`example.html` in this folder is a throwaway reference demonstrating pool layouts. It is not product UI.
