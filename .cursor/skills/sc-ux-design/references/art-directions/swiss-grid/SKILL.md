---
name: swiss-grid
description: "Art-direction pack: strict modular grid, hairline rules, zero border-radius. Locked layout pool for draft HTML under sc-ux-design."
---

# swiss-grid

Optional art-direction pack for draft HTML. Load after `shared-draft-directives.md` when the discuss brief selects **swiss-grid**. Pick structure from the locked layout pool below - do not invent a new page architecture each run.

## Authority

This pack **defers** to `references/anti-slop-catalog.md` on conflict: anti-slop-catalog is authoritative (wins). Do not ship catalog-banned house style - no Inter-as-sole-font (or other overused AI faces alone), no purple/cream defaults, no soft-shadow + hairline as elevation recipe.

## iron rules

1. **Strict grid** - Compose on an explicit modular grid (prefer 12 columns or a documented unit). Align edges to the grid; no free-floating clusters that ignore column gutters.
2. **Hairlines** - Use 1px (or hairline) borders / rules for structure and separation. No soft multi-layer shadows as the elevation recipe; edges define modules, not blur.
3. **Zero border-radius** - All surfaces, controls, and media frames use `border-radius: 0`. No pills, rounded cards, or soft blobs.
4. **Content fidelity** - Copy, names, metrics, and brand lines come only from the approved brief / user content. Do not invent claims to fill empty grid cells; leave cells sparse or use brief-supplied placeholders.
5. **Type** - Prefer a grotesque / neo-grotesque pairing with clear display vs body contrast. Do not default to overused AI grotesques. Brief tokens override pack suggestions.
6. **Palette** - Cool neutrals and one decisive accent from the brief. Do not default to warm cream / beige paper or purple-indigo gradients.

## Locked layout pool

Agents **must** pick one named layout from this pool for the primary page composition (or combine at most two when the brief needs a secondary band). Names are stable identifiers for discuss notes and designer critique.

### Named layouts

- **poster-masthead** - Full-width oversized display title (or brand wordmark) flush to the top grid band; hairline under the masthead; multi-column body on the remaining rows.
- **split-rule** - Two equal columns separated by a continuous vertical hairline; shared top/bottom horizontal rules lock the frame.
- **index-rail** - Narrow left rail (index, meta, or nav list) and wide main field, divided by a vertical hairline; rail stays flush to the outer margin.
- **quad-matrix** - 2x2 equal modules with crossing hairlines; each cell holds one content block aligned to the same inset.
- **triptych-band** - Three equal columns with vertical hairlines; optional full-width hairline header band above the trio.
- **asymmetric-twelve** - 12-unit row with an 8+4 or 9+3 split (or mirror); vertical hairline on the column break; optional stacked bands that keep the same ratio.
- **baseline-stack** - Full-width horizontal bands stacked with hairline rules between; each band is one content row on the same column grid.
- **corner-flag** - Small top-leading meta / flag cell and a large adjacent content field (classic poster asymmetry); outer frame hairline optional.

## Usage notes

- Demonstrate the chosen layout with real brief content; empty decorative chrome is not a substitute for structure.
- Cards are not required. Prefer grid cells, rules, and type hierarchy over nested card chrome.
- Reference: `example.html` in this folder (throwaway DRAFT mock - not product UI).
