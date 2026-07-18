---
name: sc-ux-design
description: "UI quality control: anti-slop enforcement, HTML draft previews, animation quality rules, browser visual verification. Activate on slop check, draft preview, visual verify, or UI quality review."
disable-model-invocation: true
---

# sc-ux-design

Quality-control companion for UI implementation: anti-slop enforcement, draft previews before code, animation quality, browser visual verification. Works alongside sc-design - sc-design owns design direction; sc-ux-design enforces quality during implementation.

## When to use

Activate on:
- **"Check for slop" / "anti-slop" / "slop audit"** - run anti-slop detection
- **"Preview draft" / "create draft" / "draft HTML"** - pre-implementation draft workflow
- **"Visual verify" / "visual test" / "browser check"** - Playwright visual verification
- **"UI quality check"** - comprehensive UX quality review
- Before UI implementation begins - design brief checkpoint

## Workflow

### Design brief (forced checkpoint)

Before writing any UI implementation code:

1. **Produce a design brief** covering 6 dimensions:
   - **Product metaphor and mood** - e.g., "studio dashboard", "reading room"
   - **Typography direction** - display + body pairing with rationale
   - **Color palette** - 3–5 tokens: bg, surface, text, accent, danger
   - **Layout structure** - first screen wireframe description
   - **Motion intent** - subtle / standard / none
   - **Spacing scale** - 4pt or 8pt base

2. **Present the brief for user approval**. No implementation code until explicitly approved.

### Draft preview

After design brief approval, before real implementation:

1. **Generate a standalone HTML draft** under `.space/missions/<id>/design/drafts/` showing layout structure, typography, color scheme, component arrangement, and spacing rhythm.

2. **Every draft MUST include**: visible "DRAFT - Not Final" banner, `data-draft="true"` on root element, versioned filename (`<name>-draft-v1.html`).

3. **Serve for review**: `node .cursor/skills/sc-ux-design/scripts/serve-html.mjs .space/missions/<id>/design/drafts/ --open`

4. **Iterate** until approved (max 3 rounds - if still unapproved, escalate to user for direction). Only then begin real implementation.
5. **Before approval**: check the draft at 375px viewport width. If layout breaks at mobile, fix before asking for approval.

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

**Tier 3 - Browser visual check** (optional, needs Playwright):
`node .cursor/skills/sc-ux-design/scripts/visual-verify.mjs <html-file>`
- 3 viewports (375/768/1280px), full-page screenshots
- Audits: horizontal overflow, clipped content, text touching viewport edge, cramped padding
- JSON report: `breakpoint`, `issues` (selector, kind, severity), `screenshots`
- Install: `cd .cursor/skills/sc-ux-design && npm install`
- Skill functions without Tier 3 if Playwright is unavailable.

## Rules

### Anti-slop

- **Must**: Run `npx impeccable detect` on all HTML output before claiming UI work is complete. Fix all CLI-detected violations before shipping.
- **Must not**: Use any pattern flagged in `references/anti-slop-catalog.md` (purple-blue gradients, glassmorphism, nested cards, side-tab borders, cream/beige palettes, gradient text, hero eyebrows) without explicit user approval. Document intentional exceptions in `decisions.md`.
- **Must not**: Use Inter/Geist/Space Grotesk as sole font without deliberate pairing.

### Design brief

- **Must**: Produce a 6-dimension design brief before writing UI implementation code.
- **Must**: Obtain explicit user approval on the brief before proceeding.
- **Must not**: Skip the brief checkpoint even for "quick" UI changes that affect visual design.
- **Must**: After implementation, cross-check the output against the approved brief (palette, typography, layout, motion). Fix any drift before running anti-slop detection.

### Draft preview

- **Must**: Generate a draft HTML preview and obtain user approval before writing implementation code.
- **Must**: Every draft HTML file includes: visible "DRAFT - Not Final" banner, `data-draft="true"` on root element, versioned filename.
- **Must**: Check draft layout at 375px viewport width before asking for approval.
- **Must not**: Use draft HTML as the basis for production implementation. Drafts are throwaway mockups.
- **Must**: After 3 draft rounds without approval, escalate to the user for direction instead of iterating indefinitely.

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

- Design direction and visual identity - use sc-design
- Full accessibility audit (WCAG, ARIA, keyboard nav) - use sc-review
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

- [ ] Design brief produced and approved (6 dimensions)
- [ ] Draft preview generated, responsive-checked at 375px, and approved (max 3 rounds)
- [ ] `DESIGN.md` read and applied (or candidate generated)
- [ ] Implementation cross-checked against approved design brief (no palette/typo/layout/motion drift)
- [ ] `npx impeccable detect` run - zero unfixed violations
- [ ] 5 LLM-only patterns reviewed with concrete heuristics (glassmorphism, extreme radius, amateur SVG, hero metrics, identical grids)
- [ ] Animation: durations in range, easing rules followed, reduced-motion respected
- [ ] No banned fonts (Inter/Geist/Space Grotesk) without deliberate pairing
- [ ] Tier 3 visual verification run if Playwright available

## References

- `references/anti-slop-catalog.md` - all 46 impeccable.style patterns with detection methods and fixes
- `references/animation-guidelines.md` - duration standards, easing rules, reduced-motion, anti-patterns
- `scripts/serve-html.mjs` - local HTML draft preview server
- `scripts/visual-verify.mjs` - Playwright browser visual verification script
- `test/fixture-slop.html` - test fixture with known slop patterns for script validation
- [impeccable.style/slop](https://impeccable.style/slop) - source catalog (2026-07-10)
