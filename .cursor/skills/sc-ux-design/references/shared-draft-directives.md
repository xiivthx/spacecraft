# Shared draft directives

Always-on layer for draft HTML prompt assembly under `sc-ux-design`. Load this first; then an optional art-direction pack body; then the approved design brief / user content tail.

These directives are spacecraft house rules for **how** drafts are built - not a default aesthetic. Typography, palette, and layout personality come from the brief and (when selected) the pack.

## Tech

- Output **standalone HTML** suitable for local preview (`serve-html.mjs`). Prefer self-contained markup (inline `<style>` or a single embedded stylesheet). Do not require a product app build, CDN framework stack, or React tree for the draft itself.
- Root element must carry `data-draft="true"`. Include a visible **DRAFT - Not Final** banner.
- Use semantic structure (`header`, `main`, `section`, `nav`, headings in order). Prefer CSS custom properties for tokens named in the brief.
- Support a usable first pass at **375px** and desktop widths; avoid horizontal overflow and edge-flush body text.
- Motion: respect `prefers-reduced-motion: reduce`; keep UI transitions short; no bounce/elastic easing or decorative-only motion.
- Do not vendor third-party skill catalogs, marketing surface templates, or foreign-language prompt blobs into the draft pipeline.

## Fidelity

- Treat the approved **design brief** and any user-supplied copy, names, and constraints as source of truth for tokens and copy during draft generation. Do not invent product claims, metrics, schedules, or brand lines the brief did not supply.
- When an art-direction pack is selected, pick structure from that pack's **locked layout / section pool** - do not freestyle a new page architecture each run.
- When no pack is selected (`none - custom brief only`), follow the brief's layout and type/color directions only.
- Tokens in the draft (bg, surface, text, accent, danger; type pairing; spacing base) must match the brief. Flag drift rather than "improving" the palette silently.
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
- On conflict between a pack suggestion and the anti-slop catalog, **defer to the catalog** (pack fidelity is refined in later workflow steps; catalog wins on pattern bans).

## Assemble reminder

```
shared-draft-directives.md  →  art-directions/<pack>/ (if selected)  →  brief / user content
```
