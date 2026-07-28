---
name: sc-ux-design
description: "UI quality: draft HTML under /sc-discuss; anti-slop and visual verify under /sc-run. Activate on slop check, draft preview, visual verify, or UI quality review."
disable-model-invocation: true
---

# sc-ux-design

UI quality companion: design brief + draft HTML under `/sc-discuss`; anti-slop and browser visual verification under `/sc-run` / sc-web-frontend. Draft discovery is discuss-owned; AFK run **ports** look from an approved draft (visual source of truth), with scenario coverage required before approval.

## When to use

Activate on:
- **"Check for slop" / "anti-slop" / "slop audit"** - run anti-slop detection
- **"Preview draft" / "create draft" / "draft HTML"** - under `/sc-discuss` before implement
- **"Visual verify" / "visual test" / "browser check"** - playwright-cli / Cursor IDE browser visual verification (post-build)
- **"UI quality check"** - comprehensive UX quality review
- During `/sc-discuss` for visual UI/FE - design brief + draft checkpoint
- During `/sc-run` after visual implementation - Tier 3 recheck + draft-parity (not draft discovery)

## Workflow

### Prompt assembly (draft generation)

When generating draft HTML, assemble the generation prompt in this **fixed order**:

1. **Shared draft directives** - always load `references/shared-draft-directives.md` (tech + fidelity + scenario matrix + anti-slop alignment).
2. **Design system** - when `DESIGN.md` exists at the project root (or package root for the UI), load it next as the house look / tokens / personality. If missing, skip this layer; the design brief phase must produce a candidate `DESIGN.md`.
3. **Brief / content tail** - append the approved design brief, `spec.md` feature/state requirements, and any user-supplied copy or constraints last so they bind the earlier layers.

Do not reverse or interleave these layers. Shared directives set how drafts are built; `DESIGN.md` is the project look SoT when present; the brief/content tail is the mission-specific layout, states, and copy for this draft.

### Design brief (forced checkpoint - `/sc-discuss`)

Before clearing discuss on visual work (and before any UI implementation code):

1. **Read `DESIGN.md`** when present. Treat it as the default look (tokens, type, mood, principles). Do not ask the human to pick an art-direction pack - there are no packs.
2. **References (image / text):** when the human supplies mood boards, screenshots, URLs, or prose references, extract cues into the brief. Never silent-clone full chrome from a reference. Require an explicit **borrow scope** (exactly one):
   - `mood` - atmosphere, density, motion feel only
   - `tokens` - mood + color / type / spacing
   - `layout` - tokens + primary page structure
   - `chrome` - layout + component look (buttons, tables, empty/error)
   Default when the human wants "that vibe" without listing components: `mood` or `tokens`. Record `Reference borrow: <scope>` (and source path/URL if given) in `decisions.md`.
3. **House conflict:** if proposed art/style conflicts with existing `DESIGN.md`, stop and ask once - do not silently override. Record exactly one outcome in `decisions.md`:
   - `DESIGN conflict: mission exception` - this screen may diverge; leave `DESIGN.md` unchanged
   - `DESIGN conflict: update house` - edit `DESIGN.md` to the new SoT, then align the brief
   - `DESIGN conflict: keep house` - reject the new direction; brief stays on `DESIGN.md`
   Explicit user choice wins over `DESIGN.md` only after A or B is recorded. Anti-slop catalog still wins unless a catalog exception is also recorded.
4. **Produce a design brief** covering 6 dimensions (align to `DESIGN.md` when it exists and conflict outcome is not mission exception / update-pending; invent a candidate system when `DESIGN.md` is missing):
   - **Product metaphor and mood** - e.g., "studio dashboard", "reading room"
   - **Typography direction** - display + body pairing with rationale
   - **Color palette** - 3–5 tokens: bg, surface, text, accent, danger
   - **Layout structure** - first screen wireframe description
   - **Motion intent** - subtle / standard / none
   - **Spacing scale** - 4pt or 8pt base
   Include borrow scope and conflict outcome lines when applicable.
5. **Present the brief for user approval**. No implementation code until explicitly approved. When `DESIGN.md` was missing, the brief includes a candidate `DESIGN.md` for approval (seeded from references within the chosen borrow scope when present).

### Draft preview (`/sc-discuss`)

After design brief approval, before `/sc-run` / real implementation:

1. **Generate a standalone HTML draft** under `.space/missions/<id>/design/drafts/` that shows **layout, style tokens (colors/type/spacing), key components, and a full scenario matrix** - enough for the human to judge look, structure, and state coverage. Not a wireframe-only sketch. Assemble the draft prompt per **Prompt assembly** above (shared directives → `DESIGN.md` when present → brief/content tail). Do not generate draft HTML until the design brief is approved.

2. **Every draft MUST include**: visible "DRAFT - Not Final" banner, `data-draft="true"` on root element, versioned filename (`<name>-draft-v1.html`), CSS custom properties for brief tokens, and a visible **Scenario matrix** with `data-state="<name>"` panels covering at least: empty, error, few, many, plus feature/behavior surfaces from `spec.md` (and loading when async is implied). Each panel shows real component chrome - not layout boxes only.

3. **Designer gate (required):** Task(`sc-designer`) on the draft. Commander applies critical and important fixes (`sc-designer` is readonly). Missing scenario states or non-portable chrome = critical. Check 375px before and after fixes. Do not serve or present the draft to the human until this gate passes.

4. **Serve for review**: `node .cursor/skills/sc-ux-design/scripts/serve-html.mjs .space/missions/<id>/design/drafts/ --open`

5. **Under `/sc-discuss`**: iterate (draft → designer → fix → human) until the human likes it. On approval, record `UI draft approved: <draft-file>` in `decisions.md` only when the scenario matrix is complete; then `spacecraft clarify-status clear`. Do not start RED-GREEN until that record exists.

6. **Under `/sc-run`**: do **not** invent or iterate draft HTML. Port look from the approved draft. If approval is missing, stop and recommend `/sc-discuss`.

7. **Iterate** in discuss until approved (max 3 human rounds - if still unapproved, escalate to user for direction). Each new draft version re-runs the designer gate before human HIL.
8. **Before human approval**: check the draft at 375px viewport width. If layout breaks at mobile, fix before asking for approval.

### DESIGN.md integration

1. **Read `DESIGN.md`** before any UI work.
2. **If missing**, the design brief phase produces a candidate `DESIGN.md` for approval (from brief + reference borrow within scope).
3. After `UI draft approved`, **update `DESIGN.md` from the approved draft** (tokens, type, spacing) so product CSS matches the draft - unless `DESIGN conflict: mission exception` was recorded (then leave house `DESIGN.md` alone; mission look lives in the approved draft only). On `DESIGN conflict: update house`, sync `DESIGN.md` to the new SoT from the approved draft. Post-approval, the approved draft owns look for port; `DESIGN.md` is the house token doc, not a license to freestyle chrome.

### Anti-slop detection & visual verification

Run after implementation:

**Step 0 - Draft parity** (before detection): Cross-check the implementation against the **approved draft HTML** (and brief tokens embedded there). Verify: color palette matches (no drift), typography pairing is intact, layout structure matches the draft, component chrome matches (buttons/inputs/tables/empty/error - not layout-only), motion intent is respected, and each draft `data-state` has a corresponding product state or test. Flag layout-only match with different chrome as blocking drift. Fix before running detectors.

**Tier 1 - CLI** (36 rules + 5 browser-rendered via headless):
`npx impeccable detect <html-file>` - catches 41 patterns, 5 need browser rendering.

**Tier 2 - LLM-only** (5 rules): Review each with a concrete heuristic:
- **Glassmorphism**: Does any blur/glass effect serve a real layering purpose? If purely decorative → flag.
- **Extreme border-radius**: Are any cards or inputs rounded past 16px? If tags/buttons only → acceptable. If cards → flag.
- **Amateurish SVG**: Are there hand-coded SVG illustrations? If not production-quality vector assets → flag.
- **Hero metric layout**: Big number + small label + three supporting stats in a row? If not real data → flag.
- **Identical card grids**: Same-sized cards repeated with icon + heading + text? If no differentiation → flag.

**Tier 3 - Browser visual check** (**required** after visual UI implementation):

Canonical browser matrix (do not expand):
1. **Vitest + happy-dom** — behavior only (not visual pixels)
2. **`playwright-cli`** — primary real-browser visual / interact (`open` → `snapshot` / `screenshot`; resize 375/768/1280); compare app screenshots to approved draft states (side-by-side LLM/browser review)
3. **Cursor IDE browser** (`cursor-ide-browser` MCP) — fallback when `playwright-cli` cannot run

Optional scripted audit (same Playwright family):  
`node .cursor/skills/sc-ux-design/scripts/visual-verify.mjs <html-file-or-url>`  
(3 viewports, overflow/clip audits, JSON report). Install: `cd .cursor/skills/sc-ux-design && npm install`.

- Capture screenshot paths in evidence / `decisions.md`; fix blocking visual issues and draft-parity gaps before `ready`.
- **Do not use** system Chrome headless or browser-use/CDP for Tier 3 (removed from the official matrix).

## Rules

### Anti-slop

- **Authority**: When `DESIGN.md` / brief guidance and `references/anti-slop-catalog.md` disagree, anti-slop-catalog is authoritative unless the human explicitly approved an exception in `decisions.md`.
- **Must**: Run `npx impeccable detect` on all HTML output before claiming UI work is complete. Fix all CLI-detected violations before shipping.
- **Must not**: Use any pattern flagged in `references/anti-slop-catalog.md` (purple-blue gradients, glassmorphism, nested cards, side-tab borders, cream/beige palettes, gradient text, hero eyebrows) without explicit user approval. Document intentional exceptions in `decisions.md`.
- **Must not**: Use Inter/Geist/Space Grotesk as sole font without deliberate pairing.

### Design brief

- **Must**: Read `DESIGN.md` when present and use it as the default look before drafting.
- **Must**: When references are supplied, record exactly one borrow scope (`mood` | `tokens` | `layout` | `chrome`) in the brief and `decisions.md`.
- **Must**: When proposed style conflicts with `DESIGN.md`, ask once and record `DESIGN conflict: mission exception | update house | keep house` before drafting.
- **Must**: Produce a 6-dimension design brief before writing UI implementation code.
- **Must**: Obtain explicit user approval on the brief before proceeding.
- **Must not**: Silent-clone full chrome from a reference image or text when borrow scope is narrower.
- **Must not**: Ask the human to pick an art-direction pack (packs removed - brief + `DESIGN.md` only).
- **Must not**: Skip the brief checkpoint even for "quick" UI changes that affect visual design.

### Draft preview

- **Must**: Generate a draft HTML preview showing layout + style tokens + key components + **scenario matrix**; obtain user approval before writing implementation code.
- **Must**: Every draft HTML file includes: visible "DRAFT - Not Final" banner, `data-draft="true"` on root element, versioned filename, and `data-state` panels for empty, error, few, many, plus spec feature/behavior surfaces (loading when implied).
- **Must**: Run Task(`sc-designer`) and apply critical/important fixes before serving or presenting any draft to the human.
- **Must**: Check draft layout at 375px viewport width before asking for approval.
- **Must**: Own draft discovery under `/sc-discuss`; `/sc-run` requires `UI draft approved: …` already recorded (or non-visual skip).
- **Must not**: Serve or present raw/unreviewed draft HTML to the human.
- **Must not**: Record `UI draft approved` when required scenario states are missing.
- **Must**: Treat the approved draft as the **visual source of truth** for production implementation - port structure, tokens, spacing, type, and component chrome; do not invent a second visual system.
- **Must**: After 3 human draft rounds without approval, escalate to the user for direction instead of iterating indefinitely.

### Visual recheck

- **Must**: After visual UI implementation, run Step 0 draft-parity against the approved draft, then Tier 3 with `playwright-cli` (preferred) or Cursor IDE browser (fallback); record screenshot paths.
- **Must**: Flag layout-only match with different chrome, or missing draft states in the product, as blocking issues.
- **Must**: Pair visual recheck with functional tests (Vitest/RTL or project suite) via `spacecraft evidence` before claiming UI ready.
- **Must not**: Use system Chrome headless or browser-use/CDP as the visual gate.

### Animation

- **Must**: Respect `prefers-reduced-motion: reduce`. Disable non-essential animations when requested.
- **Must not**: Use bounce/elastic easing, width/height animation, or decorative-only animations.
- **Must not**: Use linear easing for discrete UI transitions.
- **Must**: Keep micro-interactions 150–300ms, complex transitions ≤400ms.

### DESIGN.md

- **Must**: Read `DESIGN.md` before any UI implementation work.
- **Must**: Generate a candidate `DESIGN.md` when the project lacks one, during the design brief phase.
- **Must**: After draft approval, sync `DESIGN.md` tokens from the approved draft before or during port, except when `DESIGN conflict: mission exception` is recorded (leave house unchanged).

## Out of scope

This skill does NOT handle:

- Draft critique before human HIL - Task(`sc-designer`); Commander applies fixes
- Full accessibility audit (WCAG, ARIA, keyboard nav) - note gaps; escalate if blocking
- CSS framework selection or component library decisions - ask the user
- Product implementation - that's the build command's scope (must port from approved draft)
- Mission planning - use sc-planning
- Automated pixel-diff tooling - side-by-side LLM/browser compare in v1

## Output format

### Design brief
```
## Design Brief - [screen/component]
- **Metaphor & Mood**: [metaphor], [mood]
- **Typography**: [display] + [body] - [rationale]
- **Color Palette**: bg=[], surface=[], text=[], accent=[], danger=[]
- **Layout**: [first screen wireframe]
- **Motion**: [subtle/standard/none] - [reasoning]
- **Spacing**: [4pt/8pt] base
- **Reference borrow**: [mood|tokens|layout|chrome] - [source or none]
- **DESIGN conflict**: [none|mission exception|update house|keep house]
```

### Visual verification (JSON)
```json
{"file": "<html>", "breakpoints": [375,768,1280], "results": [
  {"breakpoint": 375, "issues": [{"selector": ".card", "kind": "overflow", "severity": "error"}], "screenshots": ["...png"]}
]}
```

## Checklist

Before claiming UI implementation is ready:

- [ ] Design brief + draft approved in `/sc-discuss` (`UI draft approved: …` in `decisions.md`)
- [ ] Approved draft includes scenario matrix (`empty`, `error`, `few`, `many`, + spec features; loading when implied)
- [ ] Draft passed Task(`sc-designer`) + critical/important fixes before human HIL
- [ ] `DESIGN.md` synced from approved draft and applied
- [ ] Implementation **ported** from approved draft (tokens, layout, chrome) - Step 0 draft-parity passed
- [ ] Each draft `data-state` mapped to product UI and/or tests
- [ ] `npx impeccable detect` run - zero unfixed violations
- [ ] 5 LLM-only patterns reviewed with concrete heuristics (glassmorphism, extreme radius, amateur SVG, hero metrics, identical grids)
- [ ] Animation: durations in range, easing rules followed, reduced-motion respected
- [ ] No banned fonts (Inter/Geist/Space Grotesk) without deliberate pairing
- [ ] Tier 3 visual verification via `playwright-cli` or Cursor IDE browser; paths recorded; side-by-side vs draft
- [ ] Functional tests passed with `spacecraft evidence`

## References

- `references/shared-draft-directives.md` - always-on draft prompt layer (tech, fidelity, scenario matrix, anti-slop alignment)
- Project `DESIGN.md` - house look / tokens (loaded after shared directives when present)
- `references/anti-slop-catalog.md` - all 46 impeccable.style patterns with detection methods and fixes
- `references/animation-guidelines.md` - duration standards, easing rules, reduced-motion, anti-patterns
- `scripts/serve-html.mjs` - local HTML draft preview server
- `scripts/visual-verify.mjs` - Playwright browser visual verification script
- `test/fixture-slop.html` - test fixture with known slop patterns for script validation
- [impeccable.style/slop](https://impeccable.style/slop) - source catalog (2026-07-10)
