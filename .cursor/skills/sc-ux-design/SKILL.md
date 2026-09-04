---
name: sc-ux-design
description: "UI quality: draft HTML under /sc-discuss; anti-slop and visual verify under /sc-run. Activate on slop check, draft preview, visual verify, or UI quality review."
disable-model-invocation: true
---

# sc-ux-design

UI quality companion: **Impeccable-primary** craft under `/sc-discuss`; draft HTML remains port visual SoT; anti-slop and browser visual verification under `/sc-run` / sc-web-frontend. `sc-designer` orchestrates Impeccable and owns Spacecraft port gates. AFK run **ports** look from an approved draft, with scenario coverage required before approval.

## When to use

Activate on:
- **"Check for slop" / "anti-slop" / "slop audit"** - run anti-slop detection
- **"Preview draft" / "create draft" / "draft HTML"** - under `/sc-discuss` before implement
- **"Visual verify" / "visual test" / "browser check"** - playwright-cli / Cursor IDE browser visual verification (post-build)
- **"UI quality check"** - comprehensive UX quality review
- During `/sc-discuss` for visual UI/FE - design brief + draft checkpoint
- During `/sc-run` after visual implementation - Tier 3 live product recheck + draft-parity (not draft discovery)

## Procedure

**Discuss / draft:** `references/impeccable-orchestration.md` + `references/shared-draft-directives.md` (also load `references/design-principles.md`; consult `references/layout-patterns.md` when naming bake-off structures).

**Prompt assembly (fixed order):** shared-draft-directives → `DESIGN.md` when present → design-principles → brief/content tail (`/impeccable shape` when `Impeccable path: active`; else 6-dimension brief) + `spec.md` Musts. Do not reverse or interleave.

Default on visual work: **`Impeccable path: active`**. Legacy 6-dimension brief only when **`Impeccable path: skipped: <reason>`**.

## Spacecraft gates (record in `decisions.md`)

These stay inline - orchestration refs do not replace them.

1. **Product context:** `Product context: <routes + shell/layout file paths + screenshot paths>` or `Product context skipped: greenfield`. Also record `Impeccable path: active` or `Impeccable path: skipped: <reason>`. Brownfield: read parent shell/layout and nearby pages before craft; bake-off candidates **Must** include existing app chrome for in-app screens.

2. **Surface checklist:** Match one primary id via `references/checklists/README.md` + `references/surface-checklist.md`. Record `UX checklist: <id>` or `UX checklist: none - <reason>` before bake-off. Read that one checklist file; fold applicable `- [ ]` items into `spec.md` Must and the approval draft. **Must not** walk the catalog or use an external checklist site as a spacecraft gate.

3. **Reference extract / borrow (when images/refs supplied):** Run `references/reference-extract.md` before brief → `design/refs/extract.md`; record `Reference extract: design/refs/extract.md` and exactly one `Reference borrow: mood | tokens | layout | chrome`. Brief **Must** cite extract rows. **Must not** enter bake-off when `Reference borrow:` is set but extract is missing. Never silent-clone full chrome beyond borrow scope.

4. **DESIGN conflict:** Read package `DESIGN.md` when present (default look). If style conflicts, ask once; record exactly one of `DESIGN conflict: mission exception | update house | keep house` (default **keep house** when file exists). Anti-slop catalog still wins unless a catalog exception is also recorded. **Must not** ask the human to pick an art-direction pack.

5. **Brief**
   - **Path active:** `/impeccable shape <surface>` is the only design brief. Ensure UI-package `.impeccable/` is gitignored; `/impeccable init` when `PRODUCT.md` missing on new/replacement world. Human approval → `Impeccable brief approved: …`. Name surface type (`persuade` | `operate`). **Do not** author a parallel 6-dimension brief.
   - **Path skipped only:** Produce a 6-dimension design brief (align to `DESIGN.md` when present): Product metaphor and mood; Typography direction; Color palette (3-5 tokens: bg, surface, text, accent, danger); Layout structure; Motion intent; Spacing scale (4pt or 8pt). Include borrow scope, extract citations (when refs), product context summary (when brownfield), `UX checklist:` id or none, conflict outcome lines; candidate `DESIGN.md` when house file was missing. Present for user approval before implementation code.

6. **Context fidelity (before bake-off):** `Context fidelity: DESIGN.md | shell:<path> | extract:<path> | product-shot:<path>` (omit absent; greenfield may omit shell/product-shot). Draft generation must load every listed source.

7. **Layout bake-off:** Generate **2-3** HTML candidates under `.space/missions/<id>/design/drafts/` (`<name>-draft-v1-<layout-id>.html`) with scaffold + primary surface chrome + Responsive ladder across all four presets (375 / 768 / 1280 / 1536) for multi-region UIs. Serve; human picks. Record `Layout bake-off winner: <file>` or `Layout bake-off skipped: <reason>`. Do not skip silently. Full scenario matrix may wait until the winner.

8. **Polish + scenario matrix + dimension lock:** After winner (or skip), full draft with scaffold; viewport toggles; **surface-relevant scenario matrix** per shared-draft-directives + `spec.md` + checklist chrome. **Dimension lock:** one of `typography` | `color` | `layout` | `motion` | `spacing` | `chrome` per human round. Impeccable draft polish: default on Persuade / craft-critical; opt-in Operate (`Impeccable draft polish: on | skipped`).

9. **Designer port + Impeccable craft (before approval HIL):** Task(`sc-designer`, model `gemini-3.8-flash-high`) port gate first (scaffold, scenarios, checklist, ladder, continuity, extract, port readiness) - fail-closed. Then Impeccable craft (path active): `/impeccable critique` default, or finish-reviewer when approved comp / craft-critical; record `Impeccable craft gate:` and `Impeccable craft: pass | waived: <reason>`. Path skipped: designer port gate alone. Do not present approval draft until both pass (or human waive craft).

10. **`UI draft approved:`** Record `UI draft approved: <draft-file>` **only if** scenario matrix complete, `UX checklist:` recorded, bake-off winner/skip recorded, shape brief approved (path active), and port+craft gates pass (or craft waived). Default keep house `DESIGN.md` unless `DESIGN conflict: update house`. Skip draft for non-visual FE or `*-data` / `*-functional` / `*-integrate`: `UI draft skipped: …` (bake-off not required). Serve via `scripts/serve-html.mjs`; max 3 human rounds after bake-off. Under `/sc-run`: port look from approved draft - do not invent draft HTML or run bake-off/shape.

### Responsive ladder (blocking)

All four presets via `[data-draft-frame][data-viewport="…"]` (and/or media): mobile 375, tablet 768, desktop 1280, widescreen 1536. Adjacent presets **Must not** be pixel-squeezed copies; widescreen **Must not** be stretched desktop with no measure control. Document `Responsive: single-column - density/nav adapt only` when intentional. Optional: `Responsive ladder: mobile=<note>; tablet=<note>; desktop=<note>; widescreen=<note>`. Patterns in shared-draft-directives.

### Every approval-candidate draft Must include

`data-draft="true"`; `[data-draft-chrome]` outside framed `[data-draft-surface]`; viewport toggles; versioned filename; CSS custom properties for brief tokens; surface-relevant `data-state` panels inside the surface (happy path + failure/degraded; `loading` when async; `empty`/`few`/`many` when variable-length collection; + `spec.md` features).

## DESIGN.md integration

Read before UI work. Missing → candidate from brief phase. After `UI draft approved`, default keep house; sync only on `DESIGN conflict: update house`. Mission exception → leave house alone. Approved draft owns look for port.

## Anti-slop & visual verification (`/sc-run`)

**Step 0 - Draft parity:** Paired screenshots of approved `[data-draft-surface]` vs live product at matching viewports (375 / 768 / 1280, + 1536 when multi-region); record both path sets; side-by-side compare. Missing pair ⇒ draft-parity fail/uncertain.

**Tier 1 - CLI:** `npx impeccable detect <html-file>` - fix all CLI violations before claiming complete.

**Tier 2 - LLM-only:** glassmorphism purpose; extreme border-radius (>16px on cards); amateurish SVG; hero metric layout without real data; identical card grids.

**Tier 3 - Live product (required after visual UI):** Running product URL via `playwright-cli` (preferred) or Cursor IDE browser (fallback). Optional: `scripts/visual-verify.mjs <product-url>`. Task(`sc-designer`) live critique (**live-product** + draft-parity) with both image sets + live URL. Fail-closed for ready. **Must not** use system Chrome headless or browser-use/CDP. Pair with functional tests via `spacecraft evidence`.

Review protocol: `references/ux-ui-review-gates.md` (five gates); per-dimension `pass` | `fail` | `uncertain` - `uncertain` blocks ready. Deterministic tiers before LLM taste.

## Rules

### Path

- **Must**: Record `Impeccable path: active` (default) or `Impeccable path: skipped: <reason>` on visual work; follow `references/impeccable-orchestration.md` when active.
- **Must**: When path active, obtain `/impeccable shape` brief approval (`Impeccable brief approved:`) before bake-off - do not author a parallel 6-dimension brief.
- **Must**: When path skipped, produce a 6-dimension design brief before UI implementation code.
- **Must**: Ensure UI-package `.impeccable/` is gitignored when path active.
- **Must**: Obtain explicit user approval on the brief before proceeding.

### Spacecraft gates

- **Must**: Record Product context / greenfield skip; UX checklist id or none; Context fidelity before bake-off; Layout bake-off winner or skip; Reference extract + borrow when refs supplied.
- **Must**: When style conflicts with `DESIGN.md`, record `DESIGN conflict: mission exception | update house | keep house` before drafting.
- **Must**: Run layout bake-off of 2-3 candidates before deep polish unless explicit skip; dimension lock one dimension per human round.
- **Must**: Approval draft includes scaffold + surface-relevant scenario matrix; Responsive ladder at all four presets; Task(`sc-designer`) port then Impeccable craft (path active) or port alone (path skipped).
- **Must not**: Enter bake-off when `Reference borrow:` is set but extract missing; skip bake-off silently; serve unreviewed approval draft; record `UI draft approved` when scenarios / checklist / bake-off / ladder / shape (path active) / craft pass-or-waive fail.
- **Must not**: Replace draft HTML / visual SoT with a Cursor Canvas.
- **Must**: Own draft discovery under `/sc-discuss`; `/sc-run` requires `UI draft approved: …` (or skip) already recorded.
- **Must**: Treat `[data-draft-surface]` as visual SoT for port; after 3 human rounds without approval post bake-off, escalate.

### Anti-slop / recheck / DESIGN.md / animation

- **Authority**: When `DESIGN.md` / brief and `references/anti-slop-catalog.md` disagree, anti-slop-catalog wins unless human exception in `decisions.md`.
- **Must**: Run `npx impeccable detect` on HTML before claiming UI complete; no banned catalog patterns or Inter/Geist/Space Grotesk as sole font without deliberate pairing.
- **Must**: After visual implementation: Tier 3 live-product + draft-parity pass (paired evidence) before ready; apply `references/ux-ui-review-gates.md`.
- Must: Brownfield UI draft requires recorded shell, header, and navigation chrome
- Must: Tip-only draft omitting shared chrome → product-continuity: fail
- Must: Live draft-parity compares approved tip plus shared chrome at matching viewports
- Must: draft-parity: fail or uncertain for shared chrome → VERDICT: REFUTED
- **Must**: Read/keep house `DESIGN.md` (prefer ≤~200 lines); sync from draft only on update-house.
- **Must**: Respect `prefers-reduced-motion`; micro-interactions 150-300ms, complex ≤400ms; no bounce/elastic, width/height animation, decorative-only, or linear easing for discrete UI transitions.

## Out of scope

- UX orchestrate + port/craft gates before human HIL - Task(`sc-designer`) + Impeccable
- Full accessibility audit; CSS framework / component library selection; product implementation; mission planning; automated pixel-diff tooling

## Output format

### Design brief (path skipped / legacy shape)
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
- **UX checklist**: [<id> | none - <reason>]
- **DESIGN conflict**: [none|mission exception|update house|keep house]
```

### Visual verification (JSON)
```json
{"file": "<html>", "breakpoints": [375,768,1280], "results": [
  {"breakpoint": 375, "issues": [{"selector": ".card", "kind": "overflow", "severity": "error"}], "screenshots": ["...png"]}
]}
```

## Checklist

Before claiming UI ready: Product context; UX checklist; Reference extract when refs; Context fidelity; `UI draft approved:`; bake-off winner/skip; dimension lock used; scaffold + scenario matrix; Responsive ladder all four; designer port (+ craft when path active); DESIGN.md disposition; port from `[data-draft-surface]` + Step 0 draft-parity; `npx impeccable detect` clean; Tier 2 LLM patterns; animation rules; Tier 3 live-product + paired draft-parity; Task(`sc-designer`) live critique pass; functional tests with `spacecraft evidence`.

## References

- `references/impeccable-orchestration.md`, `references/shared-draft-directives.md`, `references/design-principles.md`, `references/layout-patterns.md`
- `references/checklists/README.md`, `references/surface-checklist.md`, `references/reference-extract.md`
- `references/ux-ui-review-gates.md`, `references/anti-slop-catalog.md`, `references/animation-guidelines.md`
- `scripts/serve-html.mjs`, `scripts/visual-verify.mjs`, `test/fixture-slop.html`
- Project `DESIGN.md`; [impeccable.style/slop](https://impeccable.style/slop)
