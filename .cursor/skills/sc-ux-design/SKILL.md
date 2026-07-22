---
name: sc-ux-design
description: "UI quality: draft HTML under /sc-discuss; anti-slop and visual verify under /sc-run. Activate on slop check, draft preview, visual verify, or UI quality review."
disable-model-invocation: true
---

# sc-ux-design

UI quality companion: design brief + draft HTML under `/sc-discuss`; anti-slop and browser visual verification under `/sc-run` / sc-web-frontend. Draft discovery is discuss-owned; AFK run only builds against an approved draft.

## When to use

Activate on:
- **"Check for slop" / "anti-slop" / "slop audit"** - run anti-slop detection
- **"Preview draft" / "create draft" / "draft HTML"** - under `/sc-discuss` before implement
- **"Visual verify" / "visual test" / "browser check"** - playwright-cli / Cursor IDE browser visual verification (post-build)
- **"UI quality check"** - comprehensive UX quality review
- During `/sc-discuss` for visual UI/FE - design brief + draft checkpoint
- During `/sc-run` after visual implementation - Tier 3 recheck (not draft discovery)

## Workflow

### Prompt assembly (draft generation)

When generating draft HTML, assemble the generation prompt in this **fixed order**:

1. **Shared draft directives** - always load `references/shared-draft-directives.md` (tech + fidelity + anti-slop alignment).
2. **Pack body** - when an art-direction pack is selected, load that pack under `references/art-directions/<pack>/` (iron rules + locked layout/section pool). Skip this layer when the choice is none / custom brief only.
3. **Brief / content tail** - append the approved design brief and any user-supplied copy or constraints last so they bind the earlier layers.

Do not reverse or interleave these layers. Shared directives set how drafts are built; packs (optional) set structural personality; the brief/content tail is the mission-specific source of truth for tokens and copy.

### Design brief (forced checkpoint - `/sc-discuss`)

Before clearing discuss on visual work (and before any UI implementation code):

1. **Produce a design brief** covering 6 dimensions:
   - **Product metaphor and mood** - e.g., "studio dashboard", "reading room"
   - **Typography direction** - display + body pairing with rationale
   - **Color palette** - 3–5 tokens: bg, surface, text, accent, danger
   - **Layout structure** - first screen wireframe description
   - **Motion intent** - subtle / standard / none
   - **Spacing scale** - 4pt or 8pt base

2. **Present the brief for user approval**. No implementation code until explicitly approved.

### Art-direction pack selection (`/sc-discuss`)

Before draft HTML generation, require explicit **pack selection** among:

- `swiss-grid`
- `editorial`
- `none - custom brief only`

Record the choice in `decisions.md` (e.g. `Art-direction pack: swiss-grid` or `Art-direction pack: none - custom brief only`). Human or explicit brief choice only - no silent keyword auto-matcher. Skip the pack body layer in Prompt assembly when the choice is `none - custom brief only`.

### Draft preview (`/sc-discuss`)

After design brief approval and pack selection, before `/sc-run` / real implementation:

1. **Generate a standalone HTML draft** under `.space/missions/<id>/design/drafts/` that shows **layout, style tokens (colors/type/spacing), and key components** - enough for the human to judge look and structure. Not a wireframe-only sketch. Assemble the draft prompt per **Prompt assembly** above (shared directives → pack body when selected → brief/content tail). Do not generate draft HTML until pack selection is recorded.

2. **Every draft MUST include**: visible "DRAFT - Not Final" banner, `data-draft="true"` on root element, versioned filename (`<name>-draft-v1.html`).

3. **Designer gate (required):** Task(`sc-designer`) on the draft. Commander applies critical and important fixes (`sc-designer` is readonly). Check 375px before and after fixes. Do not serve or present the draft to the human until this gate passes.

4. **Serve for review**: `node .cursor/skills/sc-ux-design/scripts/serve-html.mjs .space/missions/<id>/design/drafts/ --open`

5. **Under `/sc-discuss`**: iterate (draft → designer → fix → human) until the human likes it. On approval, record `UI draft approved: <draft-file>` in `decisions.md`, then `spacecraft clarify-status clear`. Do not start RED-GREEN until that record exists.

6. **Under `/sc-run`**: do **not** invent or iterate draft HTML. If approval is missing, stop and recommend `/sc-discuss`.

7. **Iterate** in discuss until approved (max 3 human rounds - if still unapproved, escalate to user for direction). Each new draft version re-runs the designer gate before human HIL.
8. **Before human approval**: check the draft at 375px viewport width. If layout breaks at mobile, fix before asking for approval.

### DESIGN.md integration

1. **Read `DESIGN.md`** before any UI work. Apply its tokens.
2. **If missing**, the design brief phase produces a candidate `DESIGN.md` for approval.
3. All drafts and implementation must reference `DESIGN.md` as source of truth.

### Anti-slop detection & visual verification

Run after implementation:

**Step 0 - Brief trace** (before detection): Cross-check the implementation against the approved design brief. Verify: color palette matches (no drift), typography pairing is intact, layout structure follows the wireframe, motion intent is respected. Flag and fix any drift before running detectors.

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
2. **`playwright-cli`** — primary real-browser visual / interact (`open` → `snapshot` / `screenshot`; resize 375/768/1280)
3. **Cursor IDE browser** (`cursor-ide-browser` MCP) — fallback when `playwright-cli` cannot run

Optional scripted audit (same Playwright family):  
`node .cursor/skills/sc-ux-design/scripts/visual-verify.mjs <html-file-or-url>`  
(3 viewports, overflow/clip audits, JSON report). Install: `cd .cursor/skills/sc-ux-design && npm install`.

- Capture screenshot paths in evidence / `decisions.md`; fix blocking visual issues before `ready`.
- **Do not use** system Chrome headless or browser-use/CDP for Tier 3 (removed from the official matrix).

## Rules

### Anti-slop

- **Authority**: When an art-direction pack and `references/anti-slop-catalog.md` disagree, anti-slop-catalog is authoritative.
- **Must**: Run `npx impeccable detect` on all HTML output before claiming UI work is complete. Fix all CLI-detected violations before shipping.
- **Must not**: Use any pattern flagged in `references/anti-slop-catalog.md` (purple-blue gradients, glassmorphism, nested cards, side-tab borders, cream/beige palettes, gradient text, hero eyebrows) without explicit user approval. Document intentional exceptions in `decisions.md`.
- **Must not**: Use Inter/Geist/Space Grotesk as sole font without deliberate pairing.

### Design brief

- **Must**: Produce a 6-dimension design brief before writing UI implementation code.
- **Must**: Obtain explicit user approval on the brief before proceeding.
- **Must not**: Skip the brief checkpoint even for "quick" UI changes that affect visual design.
- **Must**: After implementation, cross-check the output against the approved brief (palette, typography, layout, motion). Fix any drift before running anti-slop detection.

### Pack selection

- **Must**: Require explicit pack selection among `swiss-grid`, `editorial`, or `none - custom brief only` before draft HTML generation.
- **Must**: Record pack selection in `decisions.md` (human or explicit brief choice only).
- **Must not**: Use a silent keyword auto-matcher for pack selection. No silent keyword inference from the brief or copy.

### Draft preview

- **Must**: Generate a draft HTML preview showing layout + style tokens + key components; obtain user approval before writing implementation code.
- **Must**: Every draft HTML file includes: visible "DRAFT - Not Final" banner, `data-draft="true"` on root element, versioned filename.
- **Must**: Run Task(`sc-designer`) and apply critical/important fixes before serving or presenting any draft to the human.
- **Must**: Check draft layout at 375px viewport width before asking for approval.
- **Must**: Own draft discovery under `/sc-discuss`; `/sc-run` requires `UI draft approved: …` already recorded (or non-visual skip).
- **Must not**: Serve or present raw/unreviewed draft HTML to the human.
- **Must not**: Use draft HTML as the basis for production implementation. Drafts are throwaway mockups.
- **Must**: After 3 human draft rounds without approval, escalate to the user for direction instead of iterating indefinitely.

### Visual recheck

- **Must**: After visual UI implementation, run Tier 3 with `playwright-cli` (preferred) or Cursor IDE browser (fallback); record screenshot paths.
- **Must**: Pair visual recheck with functional tests (Vitest/RTL or project suite) via `spacecraft evidence` before claiming UI ready.
- **Must not**: Use system Chrome headless or browser-use/CDP as the visual gate.

### Animation

- **Must**: Respect `prefers-reduced-motion: reduce`. Disable non-essential animations when requested.
- **Must not**: Use bounce/elastic easing, width/height animation, or decorative-only animations.
- **Must not**: Use linear easing for discrete UI transitions.
- **Must**: Keep micro-interactions 150–300ms, complex transitions ≤400ms.

### DESIGN.md

- **Must**: Read and apply `DESIGN.md` before any UI implementation work.
- **Must**: Generate a candidate `DESIGN.md` when the project lacks one, during the design brief phase.

## Out of scope

This skill does NOT handle:

- Draft critique before human HIL - Task(`sc-designer`); Commander applies fixes
- Full accessibility audit (WCAG, ARIA, keyboard nav) - note gaps; escalate if blocking
- CSS framework selection or component library decisions - ask the user
- Product implementation - that's the build command's scope
- Mission planning - use sc-planning

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
- [ ] Draft passed Task(`sc-designer`) + critical/important fixes before human HIL
- [ ] `DESIGN.md` read and applied (or candidate generated during discuss)
- [ ] Implementation cross-checked against approved design brief (no palette/typo/layout/motion drift)
- [ ] `npx impeccable detect` run - zero unfixed violations
- [ ] 5 LLM-only patterns reviewed with concrete heuristics (glassmorphism, extreme radius, amateur SVG, hero metrics, identical grids)
- [ ] Animation: durations in range, easing rules followed, reduced-motion respected
- [ ] No banned fonts (Inter/Geist/Space Grotesk) without deliberate pairing
- [ ] Tier 3 visual verification via `playwright-cli` or Cursor IDE browser; paths recorded
- [ ] Functional tests passed with `spacecraft evidence`

## References

- `references/shared-draft-directives.md` - always-on draft prompt layer (tech, fidelity, anti-slop alignment)
- `references/art-directions/` - optional art-direction packs (loaded after shared directives when selected)
- `references/anti-slop-catalog.md` - all 46 impeccable.style patterns with detection methods and fixes
- `references/animation-guidelines.md` - duration standards, easing rules, reduced-motion, anti-patterns
- `scripts/serve-html.mjs` - local HTML draft preview server
- `scripts/visual-verify.mjs` - Playwright browser visual verification script
- `test/fixture-slop.html` - test fixture with known slop patterns for script validation
- [impeccable.style/slop](https://impeccable.style/slop) - source catalog (2026-07-10)
