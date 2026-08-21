> Consult when: writing a design brief layout structure, naming layout bake-off candidates, or picking page structure before draft polish.
> Implement / port in React + Tailwind: see `sc-web-frontend/references/layout.md`.

# Layout patterns (discuss / bake-off)

Selection reference for `/sc-discuss` visual work. Use this to name bake-off candidates (`list`, `board`, `split`, `shell`, …), fill the brief **Layout structure** dimension, and keep gaze advice out of fake chrome.

## Taxonomy

| Class | Patterns | Role in discuss |
|---|---|---|
| **Page structure** | Single column, two column, three column, sidebar left/right, split screen, card grid, hero + CTA | Distinct bake-off skeletons; record winner as `Layout bake-off winner:` |
| **Gaze / scan** | F-pattern, Z-pattern | Copy and hierarchy guidance inside a structure - **not** bake-off layout-ids |

F-pattern (NN/g): default scan on unformatted prose walls - **bad** for users and businesses if important content sits right/low without signals. Mitigate with meaningful subheads, bullets, and keyword emphasis (layer-cake). Do not treat "F layout" as a goal.

Z-pattern: sparse marketing flow above the fold (start → value → proof → CTA). Complements Hero + CTA; does not replace a page skeleton on long pages.

## Decision tree

```
Page job?
├─ Long read / narrative → Single (± two-col TOC)
├─ Sell / promote → Hero + CTA (± Z above the fold)
├─ Many scannable items → Card grid
├─ Dense comparable rows → Table (prefer table bake-off id, not card grid)
├─ Pick item then inspect often → Split / master-detail
├─ Multi-destination app chrome → Sidebar left (+ optional three-col body)
└─ Article + contextual meta → Two-col or sidebar right
```

**Must:** One primary structure per page. Nest a secondary pattern only inside main (e.g. app shell + card grid body). Do not combine hero, three-column, and card grid as competing first-viewport systems.

## Pattern cards (bake-off lens)

### 1. Single column

- **Bake-off id hint:** `single` / `article`
- **Show:** Clear vertical sections, readable measure, layer-cake headings
- **Responsive:** Same column; adapt density, spacing, and nav treatment across 375 / 768 / 1280 / 1536 (not a squeezed clone)

### 2. Two column (main + sidebar)

- **Bake-off id hint:** `two-col` / `article-toc`
- **Show:** Primary content + real TOC/related chrome (not empty boxes)
- **Responsive:** Stack with **main before** aside on mobile

### 3. Three column

- **Bake-off id hint:** `three-col` / `portal`
- **Show:** Three real regions with distinct jobs (nav | list | detail, or tools | canvas | inspector)
- **Responsive:** One pane at a time on mobile; do not pixel-squeeze three columns

### 4. Sidebar left (app shell)

- **Bake-off id hint:** `shell` / `admin`
- **Show:** Persistent left nav + topbar + main; brownfield must reuse existing shell chrome
- **Responsive ladder:** full sidebar → icon rail → drawer (same destinations)
- **Widths (guidance):** sidebar ~240–280px; rail ~64px

### 5. Sidebar right

- **Bake-off id hint:** `sidebar-right` / `kb`
- **Show:** Contextual rail (TOC/meta), not primary app nav
- **Responsive:** TOC becomes top jump links or collapsible block on mobile

### 6. Split screen

- **Bake-off id hint:** `split` / `auth-split` / `compare`
- **Show:** Two equal-weight jobs (brand|form, A|B, list|detail)
- **Responsive:** Single sequence on mobile; no unusable half-width forms

### 7. Hero header + CTA

- **Bake-off id hint:** `hero` / `landing`
- **First viewport budget:** brand + one headline + one supporting sentence + CTA group + one dominant visual
- **Must not (anti-slop):** hero eyebrows/pill chips, stats strips, cards in hero, detached promo badges on media
- **Avoid:** Marketing hero on operator screens inside an app shell

### 8. F-pattern (gaze only)

- Apply when polishing article/docs **copy hierarchy** inside single/two-col structures
- Do not name a bake-off candidate `f-pattern`

### 9. Z-pattern (gaze only)

- Apply on sparse landings with Hero + CTA
- Do not name a bake-off candidate `z-pattern` unless the human explicitly asks for that label as structure shorthand

### 10. Card grid

- **Bake-off id hint:** `cards` / `catalog`
- **Show:** Heterogeneous but scannable items; one primary action per card
- **Must not:** Identical icon + heading + blurb grids (anti-slop catalog)
- **Prefer table** when items are dense, uniform, and comparison-heavy

## Bake-off checklist

Before presenting layout candidates:

- [ ] Each candidate uses a **distinct** primary structure from the page-structure class
- [ ] Candidates show real primary-surface chrome (not wireframe boxes only)
- [ ] Responsive ladder spot-checked at 375 / 768 / 1280 / 1536 (structure change, not frame-resize-only)
- [ ] Brownfield candidates include existing app chrome when editing in-app screens
- [ ] Gaze advice (F/Z) folded into copy/hierarchy - not fake layout chrome
- [ ] Winner recorded: `Layout bake-off winner: <file>` or `Layout bake-off skipped: <reason>`

## Anti-patterns

- Aesthetic pick ("looks more SaaS") without task fit
- Three-column mobile squeeze
- Hero inside admin shell
- Card grid for tabular data
- Multiple competing CTAs above the fold
- Silent skip of bake-off when structure is still open

## Related

- `references/shared-draft-directives.md` - draft scaffold, responsive ladder CSS
- `references/anti-slop-catalog.md` - hero and card-grid slop
- `references/reference-extract.md` - borrow scope `layout` | `chrome`
- `.cursor/skills/sc-web-frontend/references/layout.md` - Tailwind / Grid implementation
