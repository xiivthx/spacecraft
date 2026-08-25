---
version: alpha
name: Orbital Console
description: Default design language for local-first web interfaces - precise, calm, technical, slightly cinematic, intentionally sparse.
colors:
  bg: "#0e1116"
  bg-deep: "#080a0d"
  surface: "#151a21"
  surface-raised: "#1b222b"
  rule: "#2a3340"
  rule-strong: "#3a4656"
  text: "#f3ead7"
  text-muted: "#b9ad98"
  text-faint: "#776f63"
  accent: "#f6b44b"
  accent-strong: "#ffcb6b"
  cyan: "#62d6cf"
  danger: "#ff6b5f"
  success: "#7bd88f"
typography:
  body:
    fontFamily: IBM Plex Sans
    fontSize: 16px
    fontWeight: 400
    lineHeight: 1.5
  display:
    fontFamily: IBM Plex Sans
    fontSize: 32px
    fontWeight: 600
    lineHeight: 1.15
    letterSpacing: -0.02em
  meta:
    fontFamily: IBM Plex Mono
    fontSize: 12px
    fontWeight: 500
    lineHeight: 1.4
rounded:
  sm: 4px
  md: 8px
  lg: 14px
  pill: 999px
spacing:
  xs: 4px
  sm: 8px
  md: 16px
  lg: 24px
  xl: 32px
  xxl: 48px
components:
  button-primary:
    backgroundColor: "{colors.accent}"
    textColor: "{colors.bg-deep}"
    rounded: "{rounded.sm}"
    padding: 12px
  button-primary-hover:
    backgroundColor: "{colors.accent-strong}"
  input:
    backgroundColor: "{colors.surface}"
    textColor: "{colors.text}"
    rounded: "{rounded.sm}"
    padding: 12px
---

# Design System: Orbital Console

## Overview

Orbital Console is the default design language for local-first web interfaces. It should feel precise, calm, technical, slightly cinematic, and intentionally sparse. The interface should be usable before it is impressive.

The UI must not become a generic SaaS landing page, decorative toy, crypto dashboard, or AI-generated template. It should earn attention through hierarchy, spacing, typography, rules, and useful information architecture.

### Principles

- Usefulness before spectacle.
- Evidence over ornament.
- Sparse but not empty.
- Structured negative space.
- One accent earns attention.
- No generic AI SaaS visual tropes.

**Always apply (Must)** - hierarchy, Responsive ladder (375 / 768 / 1280 / 1536), content priority, whitespace, consistency, a11y/readability, navigation, palette and type from this file / brief only, and brief/bake-off structure (real chrome - not a gray-box wireframe artifact).

**Prefer when they fit (Should)** - modular scale near 1.25–1.618 (Golden Ratio as rhythm, not every region to φ); centric symmetry on persuade/hero; focused asymmetry OK on operate/dashboard; Fibonacci steps only when they clarify spacing rhythm.

Full checklist and reporting (Must → critical; Should → important unless the brief required it): `.cursor/skills/sc-ux-design/references/design-principles.md`.

### CSS custom properties

Map YAML tokens into app CSS (prefer these names at the entry stylesheet):

```css
:root {
  --sc-bg: #0e1116;
  --sc-bg-deep: #080a0d;
  --sc-surface: #151a21;
  --sc-surface-raised: #1b222b;
  --sc-rule: #2a3340;
  --sc-rule-strong: #3a4656;
  --sc-text: #f3ead7;
  --sc-text-muted: #b9ad98;
  --sc-text-faint: #776f63;
  --sc-accent: #f6b44b;
  --sc-accent-strong: #ffcb6b;
  --sc-cyan: #62d6cf;
  --sc-danger: #ff6b5f;
  --sc-success: #7bd88f;
  --sc-radius-sm: 4px;
  --sc-radius-md: 8px;
  --sc-radius-lg: 14px;
  --sc-radius-pill: 999px;
}
```

## Colors

- **bg (#0e1116):** Page canvas / app background only.
- **bg-deep (#080a0d):** Deepest wells, inset panels, primary-button text on accent.
- **surface (#151a21):** Default panels, inputs, tool surfaces.
- **surface-raised (#1b222b):** Elevated tool panels and selected rows - not decorative cards.
- **rule (#2a3340):** Hairline dividers and table rules.
- **rule-strong (#3a4656):** Stronger separators and idle control borders.
- **text (#f3ead7):** Primary body and headlines.
- **text-muted (#b9ad98):** Secondary copy and labels.
- **text-faint (#776f63):** Tertiary metadata; keep contrast usable.
- **accent (#f6b44b):** Single attention color - primary CTAs and critical highlights only.
- **accent-strong (#ffcb6b):** Accent hover / active emphasis.
- **cyan (#62d6cf):** Informational / link-adjacent cues - not a second CTA color.
- **danger (#ff6b5f):** Errors and destructive actions only.
- **success (#7bd88f):** Success and healthy status only.

## Typography

- **Body:** IBM Plex Sans at ~16px for narrative UI (or another expressive humanist/grotesque approved in brief). Do not default to Inter, Roboto, Arial, or system stacks unless a mission explicitly asks for a neutral enterprise UI.
- **Display:** IBM Plex Sans semibold for page-level statements; tighter tracking than body.
- **Meta:** IBM Plex Mono only for metadata, commands, IDs, timestamps, and code-like details.
- Prefer at most two weights per screen unless a dense operator table needs a third for status.

## Layout

- Prefer left-aligned operator layouts over centered marketing shells.
- Strict **4px / 8px** spacing rhythm (`xs`–`xxl` in YAML). Group related controls; leave structured negative space rather than filler widgets.
- First viewport of a promotional surface: brand, one headline, one short support line, one CTA group, one dominant visual - not a dashboard of stats.
- Cards are not the default container. Use card-like surfaces only for concrete tools, repeated items, or detail panels.

## Elevation & Depth

- Depth comes from **tonal layers** (`bg` → `surface` → `surface-raised`) and thin rules - not multi-layer drop shadows or glow.
- Avoid broad decorative gradients as the main visual idea; atmosphere is secondary to hierarchy and IA.

## Shapes

- Interactive controls and inputs: **sm (4px)** radius by default - engineered, not pill-happy.
- Larger panels may use **md/lg**; pills (`999px`) only when the component is explicitly a chip/tag/toggle track.
- Do not mix sharp and heavily rounded languages in the same view.

## Components

### Buttons

- Stable height/padding, low radius (`rounded.sm`), visible **default / hover / focus / active / disabled** states.
- Primary: `accent` fill, deep text; hover uses `accent-strong`. One primary CTA per screen region.
- Secondary/tertiary: rule border or quiet surface - never compete with primary accent.

### Inputs

- Visible labels, helper text, and validation.
- States: **default / hover / focus / disabled / error** (danger border + message; do not rely on color alone).
- 1px `rule-strong` border on `surface` background; strong focus ring using accent or cyan at accessible contrast.

### Lists & tables

- Thin rules and table-like alignment for structured information.
- Selected row: `surface-raised` + clear focus; avoid badge spam.

### Empty / error / loading

- Empty: explain what is absent and offer one useful next action.
- Error: danger color + plain-language recovery.
- Loading: reserved, non-decorative; respect `prefers-reduced-motion`.

### Forms

- Labels always visible (not placeholder-only).
- Group fields with `md`/`lg` spacing; keep actions aligned to the operator reading edge (usually left/start).

## Do's and Don'ts

- Do use **accent** only for the single most important action (or critical highlight) in a region.
- Don't invent new palette colors mid-component - extend `DESIGN.md` first.
- Do keep body text, muted text, and controls at strong contrast.
- Don't ship nested card stacks, fake metrics, meaningless badges, or purple-on-white / cream-serif-terracotta AI template looks.
- Do provide visible keyboard focus on links, inputs, buttons, tabs, and command controls.
- Don't rely on color alone for status.
- Do prefer CSS custom properties mapped from these tokens; don't add a styling framework unless explicitly requested.
- Don't port `[data-draft-chrome]` scaffold into product UI - port `[data-draft-surface]` only after draft approval.

## Motion

- Motion must clarify state, orientation, or feedback.
- Respect `prefers-reduced-motion`.
- Do not add decorative animation without purpose (2–3 intentional motions max on visually led surfaces).

## Accessibility

- Maintain strong contrast for body text, muted text, controls, and state indicators.
- Provide visible keyboard focus for links, inputs, buttons, tabs, and command controls.
- Use semantic headings in order.
- Do not rely on color alone to convey status.
- Keep line length readable, especially in inspection panels and long-form content.

## Design artifacts

- When original design exploration is weak or generic, scout references before inventing more options.
- Split references by job: layout/template, mood/art, and interaction/motion.
- Use references to calibrate taste and structure, not to clone. Record borrow scope (`mood` | `tokens` | `layout` | `chrome`) in mission `decisions.md`.
- Treat HTML design artifacts as decision aids. Create HTML when side-by-side comparison helps; otherwise use normal text questions.
- Ask one design config question per artifact when possible.
- Keep visible copy short: one main sentence, lists of 3 bullets or fewer.
- After `/sc-discuss` approval, sync tokens here from the approved draft unless `DESIGN conflict: mission exception` is recorded.

## Implementation notes

- Prefer CSS custom properties; keep tokens close to the app entry CSS.
- Prefer plain CSS or the project's existing styling approach.
- UI code must remain boring and maintainable even if the design feels distinctive.
- Keep backend service logic separate from visual components.
- House file stays focused (durable tokens/rules). Evolve this file when design changes, then regenerate UI from it - specificity beats sprawl.
