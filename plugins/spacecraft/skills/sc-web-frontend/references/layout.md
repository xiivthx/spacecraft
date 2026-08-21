> Consult when: choosing or implementing page structure, porting draft layout regions to React + Tailwind, or collapsing multi-region chrome across breakpoints.
> Design / bake-off selection: see `sc-ux-design/references/layout-patterns.md`.

# Page layout patterns

Operational page-structure patterns for React + TypeScript + Vite + Tailwind. Pick **one primary structure per page**; nest secondary patterns inside `main` only when the task needs them (e.g. sidebar shell + card grid body).

## Taxonomy

| Class | Patterns | Meaning |
|---|---|---|
| **Page structure** | Single column, two column, three column, sidebar left/right, split screen, card grid, hero + CTA | Real DOM / CSS skeleton |
| **Gaze / scan** | F-pattern, Z-pattern | How people scan - not a separate CSS layout. Shape copy and hierarchy; do not "draw an F" as chrome |

F-pattern (NN/g): common default scan on walls of unformatted text. Prefer subheads, bullets, and bold keywords (layer-cake) so users do not skip right-side content. Z-pattern: marketing flow for sparse above-the-fold landings (logo/nav → value → proof → CTA).

## Decision tree

```
Page job?
├─ Long read / narrative → Single (± two-col TOC)
├─ Sell / promote → Hero + CTA (± Z above the fold)
├─ Many scannable items → Card grid
├─ Dense comparable rows → Table (not a card grid)
├─ Pick item then inspect often → Split / master-detail
├─ Multi-destination app chrome → Sidebar left (+ optional three-col body)
└─ Article + contextual meta → Two-col or sidebar right
```

Do not stack every pattern on one screen. Match structure to the page job, then order attention: first look → primary action → next step.

## Pattern cards

### 1. Single column

- **Use:** Articles, long docs, narrative landings.
- **Skeleton:** Centered measure (`max-w-prose` / `max-w-3xl`), clear vertical rhythm.
- **Avoid:** Dashboards that need simultaneous modules.
- **Responsive:** Already one column; keep readable measure on widescreen.
- **Tailwind sketch:**

```tsx
<main className="mx-auto w-full max-w-prose px-4 py-8 sm:px-6">
  {/* headings + body */}
</main>
```

### 2. Two column (main + sidebar)

- **Use:** Blog, docs, knowledge - primary content + TOC / related.
- **Skeleton:** Main flexible; aside ~240px. Keep document order: main before aside when stacking.
- **Avoid:** Empty or duplicate-of-nav sidebars.
- **Responsive:** `md+` two columns; mobile stack with main first.
- **Tailwind sketch:**

```tsx
<div className="mx-auto grid max-w-6xl grid-cols-1 gap-8 px-4 md:grid-cols-[minmax(0,1fr)_240px]">
  <main>…</main>
  <aside>…</aside>
</div>
```

### 3. Three column

- **Use:** Portals, dense dashboards, mail-style nav | list | detail.
- **Skeleton:** CSS Grid + named areas (holy grail or Outlook 3-pane).
- **Avoid:** Narrow viewports without a one-pane-at-a-time flow.
- **Responsive:** Mobile = one pane (list → detail navigation); tablet may keep two panes.
- **Must:** `minmax(0, 1fr)` on flexible tracks to prevent phantom scrollbars from wide tables.

```tsx
<div className="grid min-h-0 grid-cols-1 lg:grid-cols-[200px_minmax(0,1fr)_minmax(0,320px)]">
  <nav>…</nav>
  <section className="min-w-0 overflow-auto">…</section>
  <aside className="min-w-0 overflow-auto">…</aside>
</div>
```

### 4. Sidebar left (app shell)

- **Use:** Admin, ERP, CRM, SaaS with many destinations or nested nav.
- **Skeleton:** Persistent left nav + topbar (context / account) + scrollable main. Sidebar ~240–280px; icon rail ~64px.
- **Avoid:** Marketing sites or apps with only 2–3 top-level routes (prefer top nav).
- **Responsive ladder:**
  - Desktop (`lg+`): full sidebar
  - Tablet (`md`): icon rail (same nav data)
  - Mobile: drawer / sheet triggered from topbar
- **Must:** One scroll owner - chrome fixed; only `main` scrolls.

```tsx
<div className="grid h-dvh grid-cols-1 grid-rows-[auto_minmax(0,1fr)] md:grid-cols-[64px_minmax(0,1fr)] lg:grid-cols-[256px_minmax(0,1fr)]">
  <aside className="hidden md:block">…</aside>
  <header className="col-start-1 md:col-start-2">…</header>
  <main className="min-h-0 overflow-auto md:col-start-2">…</main>
</div>
```

### 5. Sidebar right

- **Use:** Article / KB contextual TOC, related, meta (LTR).
- **Differs from left:** Right rail is contextual, not primary app navigation.
- **Responsive:** Move TOC above content or to sticky jump links on mobile.
- **RTL:** Mirror rail to the start edge for the writing direction.

### 6. Split screen

- **Use:** Login/register (brand | form), compare two items, master–detail with frequent switches.
- **Skeleton:** Equal or weighted halves; optional resizable divider for detail workflows.
- **Avoid:** Rare selection that always opens full-page detail.
- **Responsive:** Single sequence on mobile (form first, or list → push detail).

```tsx
<div className="grid min-h-dvh grid-cols-1 lg:grid-cols-2">
  <section>…</section>
  <section>…</section>
</div>
```

### 7. Hero header + CTA

- **Use:** Landing / promo first viewport.
- **Budget (house):** Brand + one headline + one supporting line + CTA group + one dominant visual. No stats strips, pill chips, or cards in the hero.
- **Avoid:** Operator / admin screens inside an app shell.
- **Responsive:** Full-bleed visual plane; CTA hit target ≥44px; text must remain readable over image on mobile.

### 8. F-pattern (gaze)

- Not a CSS layout. Place critical words early in lines and near the top of the **content** area; use meaningful H2/H3; avoid unmarked prose walls.
- Prefer layer-cake scanning over hoping users follow an F.

### 9. Z-pattern (gaze)

- Sparse marketing above the fold: start top-left, cross to top-right, diagonal to lower-left, end at lower-right CTA.
- Long multi-section pages need repeated local hierarchy - not one giant Z.

### 10. Card grid

- **Use:** Products, portfolio, services, heterogeneous KPI tiles.
- **Avoid:** Dense sortable/filterable data - use a table.
- **Rules:** One concept + one primary action per card; avoid identical icon+heading grids (anti-slop).
- **Responsive:** Prefer auto-fit over hard-coded 3→2→1 when possible.

```tsx
<div className="grid grid-cols-[repeat(auto-fit,minmax(min(100%,280px),1fr))] gap-6">
  {/* cards */}
</div>
```

## Implementation rules

| Concern | Practice |
|---|---|
| Page skeleton | CSS **Grid** (`grid-template-areas` when regions are named) |
| Toolbars / nav rows | **Flex** |
| App shell overflow | Scroll `main` only; `minmax(0, 1fr)` on flexible tracks |
| Card density | `auto-fit` + `minmax(min(100%, Npx), 1fr)` |
| Breakpoints | Align with draft presets: 375 / 768 / 1280 (+ 1536 when multi-region) |
| Sidebar collapse | Full → icon rail → drawer; **one** nav data source |
| Port from draft | Map `[data-draft-surface]` regions 1:1; do not invent a second skeleton |

## Anti-patterns

- Choosing layout for "pro look" instead of task fit
- Three columns on mobile without a pane flow
- Marketing hero inside admin chrome
- Card grid for tabular data
- Treating F-pattern as a goal rather than a failure mode of unstructured text
- Multiple competing CTAs in the first viewport
- Pure layout wrapper components when parent `grid` / `flex` / `gap` suffice (see `styling.md`)

## Related

- `references/styling.md` - Tailwind utilities, breakpoints, a11y styling
- `references/components.md` - compose pages from `components/ui` primitives
- `.cursor/skills/sc-ux-design/references/layout-patterns.md` - bake-off selection and gaze guidance
- `.cursor/skills/sc-ux-design/references/anti-slop-catalog.md` - hero budget, identical card grids
- `.cursor/skills/sc-ux-design/references/shared-draft-directives.md` - responsive ladder CSS for drafts
