# Shared draft directives

Always-on layer for draft HTML prompt assembly under `sc-ux-design`. Load this first; then project `DESIGN.md` when present; then the approved design brief / user content tail.

These directives are spacecraft house rules for **how** drafts are built - not a default aesthetic. Typography, palette, and layout personality come from `DESIGN.md` (when present) and the approved brief.

## Tech

- Output **standalone HTML** suitable for local preview (`serve-html.mjs`). Prefer self-contained markup (inline `<style>` or a single embedded stylesheet). Do not require a product app build, CDN framework stack, or React tree for the draft itself.
- Root element must carry `data-draft="true"`. Include a visible **DRAFT - Not Final** banner.
- Use semantic structure (`header`, `main`, `section`, `nav`, headings in order). Prefer CSS custom properties for tokens named in the brief / `DESIGN.md`.
- Support a usable first pass at **375px** and desktop widths; avoid horizontal overflow and edge-flush body text.
- Motion: respect `prefers-reduced-motion: reduce`; keep UI transitions short; no bounce/elastic easing or decorative-only motion.
- Do not vendor third-party skill catalogs, marketing surface templates, or foreign-language prompt blobs into the draft pipeline.

## Fidelity

- Treat project **`DESIGN.md`** (when present) as the house look SoT for tokens, type pairing, mood, and principles. Do not invent a competing design system unless `decisions.md` records `DESIGN conflict: mission exception` or `update house`.
- Treat the approved **design brief** and any user-supplied copy, names, and constraints as source of truth for mission-specific layout, states, and copy during draft generation. Do not invent product claims, metrics, schedules, or brand lines the brief did not supply.
- **References:** honor the recorded borrow scope only (`mood` | `tokens` | `layout` | `chrome`). Do not silent-clone full chrome from a reference when scope is narrower. Default vibe-only borrows stop at `mood` or `tokens`.
- Follow the brief's layout structure and type/color directions; when `DESIGN.md` exists, brief tokens must align with it unless `DESIGN conflict: mission exception` or `update house` is recorded.
- Tokens in the draft (bg, surface, text, accent, danger; type pairing; spacing base) must match the brief / effective house. Flag drift rather than "improving" the palette silently.
- After human approval, the draft is the **visual source of truth** for `/sc-run`: implementers **port** structure, tokens, spacing, type, and component chrome from the approved draft. Do not rebuild a second look that only "matches the brief."
- Behavior, Verify, and acceptance remain owned by `spec.md`. If draft look and spec behavior conflict, stop and return to `/sc-discuss` - do not freestyle.

## Scenario matrix (Must)

Every visual draft must include a visible **Scenario matrix** (nav or section list) with labeled panels using `data-state="<name>"` for each primary UI surface in scope. Required states (minimum):

- `empty` - empty data
- `error` - error / failure handling
- `few` - few data
- `many` - many / dense data
- feature / behavior surfaces called out in `spec.md` (key interactive states as static panels when live interaction is not draftable)

Include `loading` / pending when the spec implies async work. Each panel must show **real component chrome** (buttons, inputs, tables, empty states, error banners) - not layout boxes only. Missing required states block designer pass, human serve, and `UI draft approved`.

## Anti-slop alignment

- Treat `references/anti-slop-catalog.md` as the pattern authority for draft generation and critique. Do not ship known AI-slop tells as the default look.
- **Must not** establish spacecraft house defaults that the catalog flags: overused single typefaces as the only face, purple/violet gradient palettes, soft diffuse shadow paired with hairline borders as the house elevation recipe, or safe warm off-white / beige paper as the default surface.
- Prefer distinctive, brief-driven type pairing and a deliberate palette. Cards and chrome stay subordinate to content; avoid nested cards, hero eyebrows/pill chips, gradient text, and identical icon+heading card grids unless the brief explicitly requires them.
- On conflict between `DESIGN.md` / brief suggestion and the anti-slop catalog, **defer to the catalog** unless the human recorded an explicit exception in `decisions.md`.

## Assemble reminder

```
shared-draft-directives.md  →  DESIGN.md (if present)  →  brief / user content
```
