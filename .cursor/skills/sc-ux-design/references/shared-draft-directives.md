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
- Scenario **switcher** (buttons/links listing empty/error/few/many/…) lives in chrome. Scenario **panels** (`data-state="…"`) live inside the surface. Switching scenarios must not require resizing the browser window.
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

Example toggle behavior (adapt as needed):

```js
const widths = { mobile: 375, tablet: 768, desktop: 1280, widescreen: 1536 };
document.querySelectorAll("[data-viewport-set]").forEach((btn) => {
  btn.addEventListener("click", () => {
    const key = btn.getAttribute("data-viewport-set");
    const frame = document.querySelector("[data-draft-frame]");
    frame.style.width = widths[key] + "px";
    frame.setAttribute("data-viewport", key);
    document.querySelectorAll("[data-viewport-set]").forEach((b) => {
      b.setAttribute("aria-pressed", String(b === btn));
    });
  });
});
```

## Fidelity

- Treat project **`DESIGN.md`** (when present) as the house look SoT for tokens, type pairing, mood, and principles. Do not invent a competing design system unless `decisions.md` records `DESIGN conflict: mission exception` or `update house`.
- Treat the approved **design brief** and any user-supplied copy, names, and constraints as source of truth for mission-specific layout, states, and copy during draft generation. Do not invent product claims, metrics, schedules, or brand lines the brief did not supply.
- **References:** honor the recorded borrow scope only (`mood` | `tokens` | `layout` | `chrome`). Do not silent-clone full chrome from a reference when scope is narrower. Default vibe-only borrows stop at `mood` or `tokens`.
- Follow the brief's layout structure and type/color directions; when `DESIGN.md` exists, brief tokens must align with it unless `DESIGN conflict: mission exception` or `update house` is recorded.
- Tokens in the **surface** (bg, surface, text, accent, danger; type pairing; spacing base) must match the brief / effective house. Flag drift rather than "improving" the palette silently.
- After human approval, **`[data-draft-surface]`** is the **visual source of truth** for `/sc-run`: implementers **port** structure, tokens, spacing, type, and component chrome from the surface only. Never port `[data-draft-chrome]` / frame bezel / viewport toolbar.
- Prefer surface chrome that maps to reusable product primitives (button, field, banner, empty, table) so `/sc-run` can upgrade or add `components/ui` first, then compose the page - not one-off page-only markup for shared controls.
- Behavior, Verify, and acceptance remain owned by `spec.md`. If draft look and spec behavior conflict, stop and return to `/sc-discuss` - do not freestyle.

## Scenario matrix (Must)

Every visual draft must include a visible **Scenario matrix** with labeled panels using `data-state="<name>"` for each primary UI surface in scope. Required states (minimum):

- `empty` - empty data
- `error` - error / failure handling
- `few` - few data
- `many` - many / dense data
- feature / behavior surfaces called out in `spec.md` (key interactive states as static panels when live interaction is not draftable)

Include `loading` / pending when the spec implies async work. Each panel must show **real component chrome** (buttons, inputs, tables, empty states, error banners) - not layout boxes only - and must live **inside** `[data-draft-surface]`. Missing required states block designer pass, human serve, and `UI draft approved`.

## Anti-slop alignment

- Treat `references/anti-slop-catalog.md` as the pattern authority for draft generation and critique. Do not ship known AI-slop tells as the default look.
- **Must not** establish spacecraft house defaults that the catalog flags: overused single typefaces as the only face, purple/violet gradient palettes, soft diffuse shadow paired with hairline borders as the house elevation recipe, or safe warm off-white / beige paper as the default surface.
- Prefer distinctive, brief-driven type pairing and a deliberate palette. Cards and chrome stay subordinate to content; avoid nested cards, hero eyebrows/pill chips, gradient text, and identical icon+heading card grids unless the brief explicitly requires them.
- On conflict between `DESIGN.md` / brief suggestion and the anti-slop catalog, **defer to the catalog** unless the human recorded an explicit exception in `decisions.md`.

## Assemble reminder

```
shared-draft-directives.md  →  DESIGN.md (if present)  →  brief / user content
```
