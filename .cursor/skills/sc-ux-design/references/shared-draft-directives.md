# Shared draft directives

Always-on layer for draft HTML prompt assembly under `sc-ux-design`. Load this first; then project `DESIGN.md` when present; then the approved design brief / user content tail.

These directives are spacecraft house rules for **how** drafts are built - not a default aesthetic. Typography, palette, and layout personality come from `DESIGN.md` (when present) and the approved brief.

**Layout bake-off:** When generating bake-off candidates (2–3 layouts before winner polish), still follow Tech + scaffold below, but the **full scenario matrix** may be deferred until the winning draft is polished for approval. Candidates must show distinct page structures and real primary-surface chrome (not wireframe boxes). Approval drafts must meet the full scenario matrix rules.

## Tech

- Output **standalone HTML** suitable for local preview (`serve-html.mjs`). Prefer self-contained markup (inline `<style>` or a single embedded stylesheet). Do not require a product app build, CDN framework stack, or React tree for the draft itself.
- Root element (`html` or `body`) must carry `data-draft="true"`.
- Use the **draft scaffold** below (chrome outside / framed surface inside). Missing scaffold = incomplete draft.
- Prefer CSS custom properties for tokens named in the brief / `DESIGN.md` on the **surface** (production look), not on scaffold chrome.
- Motion: respect `prefers-reduced-motion: reduce`; keep UI transitions short; no bounce/elastic easing or decorative-only motion.
- Do not vendor third-party skill catalogs, marketing surface templates, or foreign-language prompt blobs into the draft pipeline.

## Draft scaffold (Must) - production frame vs explanation

Humans must see at a glance what is **production UI** vs **annotation**. Use this fixed split:

| Region | Attribute | Port to product? | Contents |
|--------|-----------|------------------|----------|
| Scaffold chrome | `data-draft-chrome` | **No** | DRAFT banner, viewport toggles, prose notes, legends, scenario switcher, borrow/conflict callouts |
| Device frame | `data-draft-frame` | **No** (preview chrome only) | Visible border/bezel around the surface; width set by viewport preset |
| Production surface | `data-draft-surface` | **Yes** - visual SoT | Real product layout, tokens, and component chrome only |

### Required structure

```html
<html lang="en" data-draft="true">
<body>
  <div data-draft-chrome>
    <p class="sc-draft-banner">DRAFT - Not Final</p>
    <div class="sc-draft-viewports" role="toolbar" aria-label="Viewport presets">
      <button type="button" data-viewport-set="mobile" aria-pressed="false">Mobile 375</button>
      <button type="button" data-viewport-set="tablet" aria-pressed="false">Tablet 768</button>
      <button type="button" data-viewport-set="desktop" aria-pressed="true">Desktop 1280</button>
      <button type="button" data-viewport-set="widescreen" aria-pressed="false">Widescreen 1536</button>
    </div>
    <!-- Notes, legends, scenario switcher - explanation only -->
    <aside data-draft-notes>…</aside>
    <nav data-draft-scenarios aria-label="Scenario matrix">…</nav>
  </div>

  <div data-draft-stage>
    <div data-draft-frame data-viewport="desktop" style="width:1280px">
      <div data-draft-surface>
        <!-- ONLY portable production UI + data-state panels -->
      </div>
    </div>
  </div>
</body>
</html>
```

### Scaffold rules

- Put **all** explanatory copy (why this state, token callouts, "this is a draft", designer notes) in `[data-draft-chrome]` / `[data-draft-notes]` - **outside** the frame.
- Put **only** production-looking UI inside `[data-draft-surface]`. No long prose tutorials inside the surface.
- `[data-draft-frame]` must show a clear visible border (device bezel). Label it visually as the production preview (e.g. small chrome caption "Production surface" outside or on the bezel - caption stays chrome, not surface).
- Scenario **switcher** (buttons/links listing the applicable states for that draft) lives in chrome. Scenario **panels** (`data-state="…"`) live inside the surface. Switching scenarios must not require resizing the browser window.
- Do not style scaffold chrome with product tokens in a way that could be mistaken for UI - keep chrome visually distinct (neutral meta UI).

### Viewport presets (Must)

Include working toggles that resize **only** `[data-draft-frame]` (not the whole browser necessarily):

| Preset | `data-viewport` | Frame width |
|--------|-----------------|-------------|
| mobile | `mobile` | `375px` |
| tablet | `tablet` | `768px` |
| desktop | `desktop` | `1280px` |
| widescreen | `widescreen` | `1536px` |

- Default active preset: `desktop` (or `mobile` when the brief is mobile-first - record choice in notes).
- Frame must scroll internally if content is taller than the viewport; stage may center the frame.
- Surface content must remain usable at **all four** widths (no horizontal overflow, no edge-flush body text, no broken columns). Check all four before human HIL - not mobile alone.
- Minimal inline script is allowed for viewport toggles and scenario switching only.

### Responsive layout (Must)

Changing `[data-draft-frame]` width **Must** visibly change surface **structure** for multi-region UIs at **every** preset - not only scale or squeeze the same desktop grid into a narrower frame.

**Responsive ladder (all four presets):**

| Preset | Width | Expectation (multi-region UI) |
|--------|-------|-------------------------------|
| mobile | 375 | Single column / stacked; nav as drawer/bottom/collapsed; no side-by-side dense tool chrome |
| tablet | 768 | Intermediate organization - not identical to mobile squeeze and not identical to full desktop; often 1.5-2 column or condensed shell |
| desktop | 1280 | Full multi-region / persistent nav / multi-column as brief requires |
| widescreen | 1536 | Use extra width deliberately (wider content measure, optional extra column/panel, or constrained max-width + calm margins) - **Must not** be a stretched clone of desktop with empty dead space or unreadably wide lines |

- Prefer CSS `@media` on the surface keyed to frame widths. The viewport toggle script should set `data-viewport` on `[data-draft-frame]` (and optionally on `[data-draft-surface]`) so CSS can target `[data-draft-frame][data-viewport="mobile"]` etc. and style descendants.
- Viewport toggle changing frame width alone is never enough at **any** preset.
- Adjacent presets **Must not** be pixel-squeezed copies of each other when the UI has multi-region chrome - each step shows intentional adaptation (structure, density, nav treatment, column count, and/or content measure).
- Pairwise 375-vs-1280 alone is **insufficient** as the gate.
- Horizontal squeeze of one preset into another without reflow = incomplete draft.
- Check all four presets before human HIL.
- Single-column pages may keep one column but **Must** still adapt density, spacing, and nav treatment at **each** of the four presets (not identical chrome at all widths). Record in chrome notes: `Responsive: single-column - density/nav adapt only`.

Example CSS pattern (adapt selectors to brief layout):

```css
/* mobile: stack, drawer nav */
[data-draft-frame][data-viewport="mobile"] [data-draft-surface] .app-shell { flex-direction: column; }
[data-draft-frame][data-viewport="mobile"] [data-draft-surface] .app-sidebar { display: none; }
[data-draft-frame][data-viewport="mobile"] [data-draft-surface] .app-nav-drawer { display: block; }

/* tablet: condensed shell, intermediate columns */
[data-draft-frame][data-viewport="tablet"] [data-draft-surface] .app-shell { flex-direction: row; }
[data-draft-frame][data-viewport="tablet"] [data-draft-surface] .app-sidebar { width: 64px; }
[data-draft-frame][data-viewport="tablet"] [data-draft-surface] .app-nav-drawer { display: none; }
[data-draft-frame][data-viewport="tablet"] [data-draft-surface] .app-main { grid-template-columns: 1fr; }

/* desktop: full multi-region */
[data-draft-frame][data-viewport="desktop"] [data-draft-surface] .app-shell { flex-direction: row; }
[data-draft-frame][data-viewport="desktop"] [data-draft-surface] .app-sidebar { width: 240px; display: block; }
[data-draft-frame][data-viewport="desktop"] [data-draft-surface] .app-main { grid-template-columns: 1fr 1fr; }

/* widescreen: deliberate extra width - not stretched desktop */
[data-draft-frame][data-viewport="widescreen"] [data-draft-surface] .app-main { max-width: 72rem; margin-inline: auto; }
[data-draft-frame][data-viewport="widescreen"] [data-draft-surface] .app-main { grid-template-columns: 1fr 1fr 320px; }
```

Example toggle behavior (adapt as needed):

```js
const widths = { mobile: 375, tablet: 768, desktop: 1280, widescreen: 1536 };
document.querySelectorAll("[data-viewport-set]").forEach((btn) => {
  btn.addEventListener("click", () => {
    const key = btn.getAttribute("data-viewport-set");
    const frame = document.querySelector("[data-draft-frame]");
    const surface = document.querySelector("[data-draft-surface]");
    frame.style.width = widths[key] + "px";
    frame.setAttribute("data-viewport", key);
    if (surface) surface.setAttribute("data-viewport", key);
    document.querySelectorAll("[data-viewport-set]").forEach((b) => {
      b.setAttribute("aria-pressed", String(b === btn));
    });
  });
});
```

## Fidelity

- Treat project **`DESIGN.md`** (when present) as the house look SoT for tokens, type pairing, mood, and principles. Do not invent a competing design system unless `decisions.md` records `DESIGN conflict: mission exception` or `update house`.
- Treat the approved **design brief** and any user-supplied copy, names, and constraints as source of truth for mission-specific layout, states, and copy during draft generation. Do not invent product claims, metrics, schedules, or brand lines the brief did not supply.
- **References:** honor the recorded borrow scope only (`mood` | `tokens` | `layout` | `chrome`). When `design/refs/extract.md` exists, treat extract rows as the mechanical source for reference cues (see `references/reference-extract.md`); do not silent-clone full chrome from a reference when scope is narrower. Default vibe-only borrows stop at `mood` or `tokens`.
- **Product context:** when `decisions.md` records `Product context:` (not greenfield skip), draft surfaces must reflect the parent shell, nav, and nearby page patterns from the listed paths - not a disconnected marketing shell on in-app screens.
- Follow the brief's layout structure and type/color directions; when `DESIGN.md` exists, brief tokens must align with it unless `DESIGN conflict: mission exception` or `update house` is recorded.
- Tokens in the **surface** (bg, surface, text, accent, danger; type pairing; spacing base) must match the brief / effective house. Flag drift rather than "improving" the palette silently.
- After human approval, **`[data-draft-surface]`** is the **visual source of truth** for `/sc-run`: implementers **port** structure, tokens, spacing, type, and component chrome from the surface only. Never port `[data-draft-chrome]` / frame bezel / viewport toolbar.
- Prefer surface chrome that maps to reusable product primitives (button, field, banner, empty, table) so `/sc-run` can upgrade or add `components/ui` first, then compose the page - not one-off page-only markup for shared controls.
- Behavior, Verify, and acceptance remain owned by `spec.md`. If draft look and spec behavior conflict, stop and return to `/sc-discuss` - do not freestyle.

## Scenario matrix (Must)

Every visual draft must include a visible **Scenario matrix** with labeled `data-state="<name>"` panels for states the primary surface can enter per `spec.md` and surface shape.

Include:
- Happy path (`default` or equivalent)
- Failure and degraded states the surface can enter (e.g. `error`, `reduced-motion` when the spec calls for them)
- `loading` when async work is implied
- `empty`, `few`, `many` when the surface presents a variable-length collection (list, table, feed, or item cards)
- Feature and behavior surfaces from `spec.md`

When `decisions.md` records `UX checklist: <id>`, load that id's file under `.cursor/skills/sc-ux-design/references/checklists/` (README + `references/surface-checklist.md`) and show applicable `- [ ]` items as real chrome or extra `data-state` panels. State-like items stay in the matrix.

Chrome notes: `Scenario matrix: <states>` (optional short note when collection density states do not apply).

Gate: missing an applicable state = critical (designer pass, human serve, and `UI draft approved` require the applicable set).

Each panel must show **real component chrome** (buttons, inputs, tables, empty states, error banners) - not layout boxes only - and must live **inside** `[data-draft-surface]`.

## Anti-slop alignment

- Treat `references/anti-slop-catalog.md` as the pattern authority for draft generation and critique. Do not ship known AI-slop tells as the default look.
- **Must not** establish spacecraft house defaults that the catalog flags: overused single typefaces as the only face, purple/violet gradient palettes, soft diffuse shadow paired with hairline borders as the house elevation recipe, or safe warm off-white / beige paper as the default surface.
- Prefer distinctive, brief-driven type pairing and a deliberate palette. Cards and chrome stay subordinate to content; avoid nested cards, hero eyebrows/pill chips, gradient text, and identical icon+heading card grids unless the brief explicitly requires them.
- On conflict between `DESIGN.md` / brief suggestion and the anti-slop catalog, **defer to the catalog** unless the human recorded an explicit exception in `decisions.md`.

## Assemble reminder

```
shared-draft-directives.md  →  DESIGN.md (if present)  →  brief / user content
```
