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
- During `/sc-run` after visual implementation - Tier 3 live product recheck + draft-parity (not draft discovery)

## Workflow

### Prompt assembly (draft generation)

When generating draft HTML, assemble the generation prompt in this **fixed order**:

1. **Shared draft directives** - always load `references/shared-draft-directives.md` (tech + fidelity + scenario matrix + anti-slop alignment).
2. **Design system** - when `DESIGN.md` exists at the project root (or package root for the UI), load it next as the house look / tokens / personality. If missing, skip this layer; the design brief phase must produce a candidate `DESIGN.md`.
3. **Brief / content tail** - append the approved design brief, `spec.md` feature/state requirements, and any user-supplied copy or constraints last so they bind the earlier layers.

Do not reverse or interleave these layers. Shared directives set how drafts are built; `DESIGN.md` is the project look SoT when present; the brief/content tail is the mission-specific layout, states, and copy for this draft.

### Design brief (forced checkpoint - `/sc-discuss`)

Before clearing discuss on visual work (and before any UI implementation code):

1. **Product context (Must before brief on brownfield):** Record in `decisions.md`:
   - `Product context: <routes + shell/layout file paths + screenshot paths>` when editing an existing app, **or**
   - `Product context skipped: greenfield` when there is no parent product shell.
   When not greenfield: read the parent shell/layout and nearby page patterns before drafting. Bake-off candidates **Must** include existing app chrome (nav, sidebar, shell) when editing in-app screens - not a floating marketing shell on operator pages.
2. **Read `DESIGN.md`** when present. Treat it as the default look (tokens, type, mood, principles). Do not ask the human to pick an art-direction pack - there are no packs.
3. **Reference extract (Must when images/refs supplied):** When the human supplies mood boards, screenshots, URLs, or other visual references, run `references/reference-extract.md` **before** the brief. Produce `.space/missions/<id>/design/refs/extract.md`; record `Reference extract: design/refs/extract.md` in `decisions.md`. Brief **Must** cite extract rows. Record human-confirmed `Reference borrow: <scope>`. **Must not** enter layout bake-off when `Reference borrow:` is set but the extract artifact is missing.
4. **References (image / text):** after extract (when refs supplied), fold cues into the brief per extract rows and borrow scope. Never silent-clone full chrome from a reference. Require an explicit **borrow scope** (exactly one):
   - `mood` - atmosphere, density, motion feel only
   - `tokens` - mood + color / type / spacing
   - `layout` - tokens + primary page structure
   - `chrome` - layout + component look (buttons, tables, empty/error)
   Default when the human wants "that vibe" without listing components: `mood` or `tokens`. Record `Reference borrow: <scope>` (and source path/URL if given) in `decisions.md` after extract when refs were supplied.
5. **House conflict:** if proposed art/style conflicts with existing `DESIGN.md`, stop and ask once - do not silently override. Record exactly one outcome in `decisions.md`:
   - `DESIGN conflict: mission exception` - this screen may diverge; leave `DESIGN.md` unchanged
   - `DESIGN conflict: update house` - edit `DESIGN.md` to the new SoT, then align the brief
   - `DESIGN conflict: keep house` - reject the new direction; brief stays on `DESIGN.md`
   Explicit user choice wins over `DESIGN.md` only after A or B is recorded. Anti-slop catalog still wins unless a catalog exception is also recorded.
6. **Produce a design brief** covering 6 dimensions (align to `DESIGN.md` when it exists and conflict outcome is not mission exception / update-pending; invent a candidate system when `DESIGN.md` is missing):
   - **Product metaphor and mood** - e.g., "studio dashboard", "reading room"
   - **Typography direction** - display + body pairing with rationale
   - **Color palette** - 3–5 tokens: bg, surface, text, accent, danger
   - **Layout structure** - first screen wireframe description
   - **Motion intent** - subtle / standard / none
   - **Spacing scale** - 4pt or 8pt base
   Include borrow scope, extract citations (when refs supplied), product context summary (when brownfield), and conflict outcome lines when applicable.
7. **Present the brief for user approval**. No implementation code until explicitly approved. When `DESIGN.md` was missing, the brief includes a candidate `DESIGN.md` for approval (seeded from references within the chosen borrow scope when present).

### Draft preview (`/sc-discuss`)

After design brief approval, before `/sc-run` / real implementation:

1. **Context fidelity (Must before bake-off):** Record in `decisions.md` or draft chrome notes:
   ```
   Context fidelity: DESIGN.md | shell:<path> | extract:<path> | product-shot:<path>
   ```
   Omit absent paths. Greenfield may omit `shell:` and `product-shot:`. Draft generation must load every listed source before bake-off candidates.
2. **Layout bake-off (before detail polish):** Generate **2–3** standalone HTML layout candidates under `.space/missions/<id>/design/drafts/` so the human can pick structure before deep scenario/chrome polish. Name them `<name>-draft-v1-<layout-id>.html` (e.g. `list`, `board`, `split`). Each candidate must use the draft scaffold (`data-draft`, chrome outside framed surface) and brief tokens; show the **primary** happy-path surface with real component chrome (not wireframe boxes only). Full scenario matrix may wait until a winner is chosen. Serve the drafts dir; present candidates side-by-side (Cursor IDE browser or default browser). Record the pick in `decisions.md` as `Layout bake-off winner: <draft-file>` (or `Layout bake-off skipped: <reason>` when layout is already forced by `DESIGN.md`, borrow scope `layout`|`chrome`, or an explicit user-named structure). Do not skip bake-off silently.

3. **Responsive structure (Must):** Viewport toggles that resize the frame alone are **not** enough at **any** preset. The surface **Must** target all four presets via `[data-draft-frame][data-viewport="…"]` (and/or media) with **size-appropriate organization** at each width:

   | Preset | Width | Expectation (multi-region UI) |
   |--------|-------|-------------------------------|
   | mobile | 375 | Single column / stacked; nav as drawer/bottom/collapsed; no side-by-side dense tool chrome |
   | tablet | 768 | Intermediate organization - not identical to mobile squeeze and not identical to full desktop; often 1.5-2 column or condensed shell |
   | desktop | 1280 | Full multi-region / persistent nav / multi-column as brief requires |
   | widescreen | 1536 | Use extra width deliberately (wider content measure, optional extra column/panel, or constrained max-width + calm margins) - **Must not** be a stretched clone of desktop with empty dead space or unreadably wide lines |

   Adjacent presets **Must not** be pixel-squeezed copies of each other when the UI has multi-region chrome - each step shows intentional adaptation (structure, density, nav treatment, column count, and/or content measure). Pairwise 375-vs-1280 alone is **insufficient** as the gate. **Blocking:** any preset that is only a horizontally squeezed version of another; any preset unusable/overflowing; widescreen = stretched desktop with no measure control when content is text-dense. Bake-off: candidates **Must** demonstrate the responsive ladder (spot-check all four), not only mobile vs desktop. Designer/approval: critique **Responsive ladder** across all four presets. Single-column content pages may keep one column but **Must** still adapt density/spacing/nav at each of the four presets (not identical chrome at all widths) - document `Responsive: single-column - density/nav adapt only` in chrome notes when intentional. Optional `decisions.md` line: `Responsive ladder: mobile=<note>; tablet=<note>; desktop=<note>; widescreen=<note>`. See `references/shared-draft-directives.md` for CSS patterns.
4. **Dimension lock (iteration):** After the winner is chosen (or bake-off skipped), refine **one visual dimension per human round** - exactly one of: `typography` | `color` | `layout` | `motion` | `spacing` | `chrome`. Lock the others. Do not change type + color + layout in the same pass (whack-a-mole). Prefer feedback like "list diffs vs reference/draft; focus: <dimension>" over "make it look better." Optional: keep refs under `design/refs/` labeled by dimension. Record the active lock when useful: `Dimension lock: <dimension>`.

5. **Polish the winning draft** with a full **scenario matrix** - enough for the human to judge look, structure, and state coverage. Not a wireframe-only sketch. Assemble the draft prompt per **Prompt assembly** above (shared directives → `DESIGN.md` when present → brief/content tail). Do not generate draft HTML until the design brief is approved. Filename stays versioned (`<name>-draft-vN.html` or keep the bake-off winner name and bump `vN` on major edits).

6. **Every draft MUST include**: `data-draft="true"`; scaffold with `[data-draft-chrome]` (banner, notes, scenario switcher, viewport toggles) **outside** a visible `[data-draft-frame]` that wraps `[data-draft-surface]` (production UI only); versioned filename; CSS custom properties for brief tokens on the surface; after bake-off (or skip), a **surface-relevant** scenario matrix of `data-state` panels inside the surface (happy path + failure/degraded the surface can enter; `loading` when async; `empty`/`few`/`many` when the surface presents a variable-length collection; plus feature/behavior surfaces from `spec.md`). Each panel shows real component chrome - not layout boxes only. Bake-off candidates may defer full matrix until the winner polish step.

7. **Designer gate (required):** Task(`sc-designer`) on the draft presented for approval (winner after bake-off, or sole draft when skipped). Commander applies critical and important fixes (`sc-designer` is readonly). Missing scenario states (on approval candidate), missing scaffold/frame, non-portable chrome, product-continuity gaps (brownfield), missing extract when borrow is set, or frame-resize-only / squeeze-only responsive structure at **any** preset = critical. Check all four viewport presets (375 / 768 / 1280 / 1536) and **Responsive ladder** - size-appropriate organization at each preset, not pairwise mobile-vs-desktop only - before and after fixes. Do not serve or present the **approval** draft to the human until this gate passes. Bake-off candidates may be shown for layout pick after a lighter scaffold/viewport/responsive-ladder sanity check; full designer gate still required before `UI draft approved`.

8. **Serve for review**: `node .cursor/skills/sc-ux-design/scripts/serve-html.mjs .space/missions/<id>/design/drafts/ --open`

9. **Under `/sc-discuss`**: iterate (draft → designer → fix → human) with dimension lock until the human likes it. On approval, record `UI draft approved: <draft-file>` in `decisions.md` only when the scenario matrix is complete **and** `Layout bake-off winner:` or `Layout bake-off skipped:` is recorded; then `spacecraft clarify-status clear`. Do not start RED-GREEN until that record exists.

10. **Under `/sc-run`**: do **not** invent or iterate draft HTML, and do **not** run layout bake-off. Port look from the approved draft. If approval is missing, stop and recommend `/sc-discuss`. Layout discovery belongs in discuss only.

11. **Iterate** in discuss until approved (max 3 human rounds after bake-off winner - if still unapproved, escalate to user for direction). Each new draft version re-runs the designer gate before human HIL.
12. **Before human approval**: exercise all four viewport presets (mobile 375, tablet 768, desktop 1280, widescreen 1536) via the draft toggles. Confirm the **Responsive ladder** at each preset - intentional, size-appropriate organization; adjacent presets not pixel-squeezed copies; widescreen uses extra width deliberately (not stretched desktop with dead space or unreadably wide lines). If layout breaks, overflows, or any preset is frame-resize/squeeze-only, fix before asking for approval.

### DESIGN.md integration

1. **Read `DESIGN.md`** before any UI work.
2. **If missing**, the design brief phase produces a candidate `DESIGN.md` for approval (from brief + reference borrow within scope).
3. After `UI draft approved`, **update `DESIGN.md` from the approved draft** (tokens, type, spacing) so product CSS matches the draft - unless `DESIGN conflict: mission exception` was recorded (then leave house `DESIGN.md` alone; mission look lives in the approved draft only). On `DESIGN conflict: update house`, sync `DESIGN.md` to the new SoT from the approved draft. Post-approval, the approved draft owns look for port; `DESIGN.md` is the house token doc, not a license to freestyle chrome.

### Anti-slop detection & visual verification

Run after implementation:

**Step 0 - Draft parity** (before detection): Capture **paired** evidence, then compare. Serve/open the approved draft HTML and screenshot **`[data-draft-surface]`** (ignore `[data-draft-chrome]` / frame bezel) at the same viewports used for live (375 / 768 / 1280, + 1536 when multi-region). Capture matching live product screenshots (Tier 3). Record **both** path sets in evidence / `decisions.md`. Side-by-side LLM/browser compare draft vs live for tokens, layout, component chrome, and applicable scenario states. Also verify: color palette matches (no drift), typography pairing is intact, layout structure matches the surface, component chrome matches (buttons/inputs/tables/empty/error - not layout-only), motion intent is respected, and each draft `data-state` has a corresponding product state or test. Flag layout-only match with different chrome as blocking drift. Missing pair ⇒ draft-parity fail/uncertain. Fix before running detectors.

**Tier 1 - CLI** (36 rules + 5 browser-rendered via headless):
`npx impeccable detect <html-file>` - catches 41 patterns, 5 need browser rendering.

**Tier 2 - LLM-only** (5 rules): Review each with a concrete heuristic:
- **Glassmorphism**: Does any blur/glass effect serve a real layering purpose? If purely decorative → flag.
- **Extreme border-radius**: Are any cards or inputs rounded past 16px? If tags/buttons only → acceptable. If cards → flag.
- **Amateurish SVG**: Are there hand-coded SVG illustrations? If not production-quality vector assets → flag.
- **Hero metric layout**: Big number + small label + three supporting stats in a row? If not real data → flag.
- **Identical card grids**: Same-sized cards repeated with icon + heading + text? If no differentiation → flag.

**Tier 3 - Live product visual check** (**required** after visual UI implementation):

Target the **running product URL** (start the app; open real product routes). Draft HTML serve alone is not live product review and does not satisfy **live-product** for ready.

Canonical browser matrix (do not expand):
1. **Vitest + happy-dom** — behavior only (not visual pixels)
2. **`playwright-cli`** — primary real-browser visual / interact on the product URL (`open` → `snapshot` / `screenshot`; resize 375/768/1280, + 1536 when multi-region); also capture draft-surface shots at those same viewports; side-by-side LLM/browser compare draft vs live
3. **Cursor IDE browser** (`cursor-ide-browser` MCP) — fallback when `playwright-cli` cannot run

Optional scripted audit (same Playwright family):  
`node .cursor/skills/sc-ux-design/scripts/visual-verify.mjs <product-url>`  
(3 viewports, overflow/clip audits, JSON report). Install: `cd .cursor/skills/sc-ux-design && npm install`. Prefer a product URL here when claiming live-product.

- Capture **paired** draft-surface + live screenshot paths in evidence / `decisions.md`; Task(`sc-designer`) live critique (**live-product** + draft-parity) with both image sets plus live URL; fix blocking visual issues, live-product gaps, and draft-parity gaps before `ready`.
- **live-product** and **draft-parity** (paired evidence + side-by-side compare) pass required before claiming UI ready (fail-closed; `uncertain` or missing pair blocks).
- **Do not use** system Chrome headless or browser-use/CDP for Tier 3 (removed from the official matrix).

## Rules

### Anti-slop

- **Authority**: When `DESIGN.md` / brief guidance and `references/anti-slop-catalog.md` disagree, anti-slop-catalog is authoritative unless the human explicitly approved an exception in `decisions.md`.
- **Must**: Run `npx impeccable detect` on all HTML output before claiming UI work is complete. Fix all CLI-detected violations before shipping.
- **Must not**: Use any pattern flagged in `references/anti-slop-catalog.md` (purple-blue gradients, glassmorphism, nested cards, side-tab borders, cream/beige palettes, gradient text, hero eyebrows) without explicit user approval. Document intentional exceptions in `decisions.md`.
- **Must not**: Use Inter/Geist/Space Grotesk as sole font without deliberate pairing.

### Design brief

- **Must**: Record `Product context: …` or `Product context skipped: greenfield` before brief on brownfield work.
- **Must**: When brownfield, read parent shell/layout and nearby page patterns before drafting; bake-off candidates include existing app chrome for in-app screens.
- **Must**: Read `DESIGN.md` when present and use it as the default look before drafting.
- **Must**: When references are supplied, run `references/reference-extract.md` before the brief; produce `design/refs/extract.md`; record `Reference extract: design/refs/extract.md` and human-confirmed `Reference borrow:` in `decisions.md`; brief must cite extract rows.
- **Must**: When references are supplied, record exactly one borrow scope (`mood` | `tokens` | `layout` | `chrome`) in the brief and `decisions.md`.
- **Must**: When proposed style conflicts with `DESIGN.md`, ask once and record `DESIGN conflict: mission exception | update house | keep house` before drafting.
- **Must**: Produce a 6-dimension design brief before writing UI implementation code.
- **Must**: Obtain explicit user approval on the brief before proceeding.
- **Must not**: Enter layout bake-off when `Reference borrow:` is set but `design/refs/extract.md` is missing.
- **Must not**: Silent-clone full chrome from a reference image or text when borrow scope is narrower.
- **Must not**: Ask the human to pick an art-direction pack (packs removed - brief + `DESIGN.md` only).
- **Must not**: Skip the brief checkpoint even for "quick" UI changes that affect visual design.

### Draft preview

- **Must**: Record `Context fidelity: …` before bake-off (omit absent paths; greenfield may omit shell/product-shot).
- **Must**: Run a **layout bake-off** of 2–3 HTML candidates (distinct page structures, brief tokens, primary surface chrome) before deep scenario polish; record `Layout bake-off winner: <file>` or `Layout bake-off skipped: <reason>` in `decisions.md`.
- **Must**: Use responsive CSS so **all four presets** (375 / 768 / 1280 / 1536) show size-appropriate organization for multi-region UIs - not frame-resize-only or squeeze-only at any preset. Adjacent presets must show intentional adaptation (structure, density, nav, column count, and/or content measure). Document `Responsive: single-column - density/nav adapt only` when single-column is intentional (still adapt density/spacing/nav at each preset).
- **Must**: After bake-off (or skip), polish the winning draft with layout + style tokens + key components + **scenario matrix**; obtain user approval before writing implementation code.
- **Must**: During draft iteration, apply **dimension lock** - change only one of `typography` | `color` | `layout` | `motion` | `spacing` | `chrome` per human round; do not restyle multiple dimensions in one pass.
- **Must**: Every approval-candidate draft HTML file includes: `data-draft="true"`; scaffold with `[data-draft-chrome]` outside a framed `[data-draft-surface]`; viewport toggles for 375 / 768 / 1280 / 1536; versioned filename; and surface-relevant `data-state` panels **inside the surface** (happy path + failure/degraded the surface can enter; `loading` when implied; `empty`/`few`/`many` when variable-length collection; plus spec feature/behavior surfaces).
- **Must**: Keep explanatory copy outside the frame; port only `[data-draft-surface]`.
- **Must**: Run Task(`sc-designer`) and apply critical/important fixes before serving or presenting the **approval** draft to the human.
- **Must**: Check draft layout at all four viewport presets before asking for approval; confirm **Responsive ladder** - size-appropriate organization at each preset, not pairwise mobile-vs-desktop only.
- **Must**: Own draft discovery (including bake-off) under `/sc-discuss`; `/sc-run` requires `UI draft approved: …` already recorded (or non-visual skip) and must not invent layouts.
- **Must not**: Serve or present raw/unreviewed approval-candidate draft HTML to the human.
- **Must not**: Record `UI draft approved` when required scenario states are missing, scaffold/frame/surface split is missing, bake-off winner/skip is missing, or the responsive ladder fails (adjacent presets pixel-squeezed, any preset unusable/overflowing, widescreen stretched desktop with no measure control) when multi-region UI is implied.
- **Must not**: Skip bake-off silently when layout is still open; use an explicit skip line when forced.
- **Must not**: Rely on horizontal squeeze of one preset's layout into another without reflow at **any** width - incomplete draft.
- **Must not**: Prompt with vague "make it look better" - prefer dimension-scoped diffs vs reference or draft. Vague aesthetic asks are gated by always-on `026-intent-coach.mdc` (ask intent first; then propose).
- **Must**: Treat `[data-draft-surface]` in the approved draft as the **visual source of truth** for production implementation - port structure, tokens, spacing, type, and component chrome; do not invent a second visual system; do not port scaffold chrome.
- **Must**: After 3 human draft rounds without approval (post bake-off), escalate to the user for direction instead of iterating indefinitely.

### Review gates

- **Must**: Apply `.cursor/skills/sc-ux-design/references/ux-ui-review-gates.md` (five gates) for discuss designer gate, run visual QC, and review/judge visual ready.
- **Must**: Run deterministic tiers (Step 0, `npx impeccable detect`, scenario `data-state` / mapping, scripted audits) before LLM taste on each dimension.
- **Must**: Per-dimension `pass` | `fail` | `uncertain` + short reason - no numeric scores. **`uncertain` on a required dimension blocks `ready`** (fail-closed; treat as fail for approval and judge).

### Visual recheck

- **Must**: After visual UI implementation, start the app and run Tier 3 against the **running product URL** with `playwright-cli` (preferred) or Cursor IDE browser (fallback); capture live screenshots; serve/open the approved draft and capture `[data-draft-surface]` shots at the same viewports; record **both** path sets; side-by-side LLM/browser compare; then Step 0 draft-parity; then Task(`sc-designer`) for **live-product** + draft-parity with both image sets plus live URL.
- **Must**: Require **live-product** pass and **draft-parity** pass (paired draft-surface + live screenshots + side-by-side compare) before claiming UI ready. Draft HTML serve alone does not satisfy live product review.
- **Must**: Flag layout-only match with different chrome, or missing draft states in the product, as blocking issues. Missing paired screenshot evidence ⇒ draft-parity fail/uncertain (fail-closed for ready).
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
- **Must**: Keep house `DESIGN.md` focused (prefer ≤~200 lines of durable tokens/rules) - specificity beats sprawl; evolve the file when design changes, then regenerate UI from it.

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
- **Reference extract**: [design/refs/extract.md or none]
- **Product context**: [routes + shell paths + screenshots or greenfield skip]
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

- [ ] Product context recorded (`Product context: …` or `Product context skipped: greenfield`)
- [ ] Reference extract produced when refs supplied (`Reference extract: design/refs/extract.md`)
- [ ] Context fidelity recorded before bake-off
- [ ] Design brief + draft approved in `/sc-discuss` (`UI draft approved: …` in `decisions.md`)
- [ ] Layout bake-off recorded (`Layout bake-off winner: …` or `Layout bake-off skipped: …`) before approval
- [ ] Draft iteration used dimension lock (one of typography/color/layout/motion/spacing/chrome per round)
- [ ] Approved draft uses scaffold (chrome outside framed `[data-draft-surface]`) + surface-relevant scenario matrix (happy path + failure/degraded the surface can enter; loading when implied; `empty`/`few`/`many` when collection; + spec features)
- [ ] Draft checked at viewport presets 375 / 768 / 1280 / 1536; **Responsive ladder** verified at all four (size-appropriate organization; no squeeze-only adjacent presets; widescreen deliberate)
- [ ] Draft passed Task(`sc-designer`) + critical/important fixes before human HIL
- [ ] `DESIGN.md` synced from approved draft and applied
- [ ] Implementation **ported** from `[data-draft-surface]` only (tokens, layout, chrome) - Step 0 draft-parity passed
- [ ] Each draft `data-state` mapped to product UI and/or tests
- [ ] `npx impeccable detect` run - zero unfixed violations
- [ ] 5 LLM-only patterns reviewed with concrete heuristics (glassmorphism, extreme radius, amateur SVG, hero metrics, identical grids)
- [ ] Animation: durations in range, easing rules followed, reduced-motion respected
- [ ] No banned fonts (Inter/Geist/Space Grotesk) without deliberate pairing
- [ ] Tier 3 live product verification on running product URL via `playwright-cli` or Cursor IDE browser; live screenshot paths recorded
- [ ] Draft-surface screenshots captured at matching viewports; both path sets recorded; side-by-side draft vs live compare passed (draft-parity)
- [ ] Task(`sc-designer`) live critique: **live-product** + draft-parity pass with paired image sets (fail-closed)
- [ ] Functional tests passed with `spacecraft evidence`

## References

- `references/shared-draft-directives.md` - always-on draft prompt layer (tech, fidelity, responsive structure, scenario matrix, anti-slop alignment)
- `references/reference-extract.md` - on-demand gate when human supplies reference images/screenshots/URLs
- Project `DESIGN.md` - house look / tokens (loaded after shared directives when present)
- `references/ux-ui-review-gates.md` - five-gate UX/UI review protocol (dimensions, fail-closed verdicts, calibration)
- `references/anti-slop-catalog.md` - all 46 impeccable.style patterns with detection methods and fixes
- `references/animation-guidelines.md` - duration standards, easing rules, reduced-motion, anti-patterns
- `scripts/serve-html.mjs` - local HTML draft preview server
- `scripts/visual-verify.mjs` - Playwright browser visual verification script
- `test/fixture-slop.html` - test fixture with known slop patterns for script validation
- [impeccable.style/slop](https://impeccable.style/slop) - source catalog (2026-07-10)
