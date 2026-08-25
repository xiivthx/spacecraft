# Design principles

Always-apply checklist for draft HTML generation, create/port critique, and `sc-designer` review under `sc-ux-design`.

These are spacecraft house rules for **how** surfaces are judged - not a default aesthetic. Tokens, type pairing, and mood come from project `DESIGN.md` (when present) and the approved brief.

Severity:

| Level | Meaning |
|-------|---------|
| **Must** | Gap is critical. Can block draft approval or UI ready when the principle is in scope. |
| **Should** | Apply when it fits the brief and surface. Gap is important polish - **not** fail-closed alone. If deliberately skipped, note why in draft chrome notes. |

## When to load

Load this file when:

- Writing or tightening a design brief
- Running layout bake-off (structure choice)
- Polishing a winning draft for approval
- Running `sc-designer` critique (create or port)
- Checking draft ↔ live parity for visual continuity

## Surface type (persuade vs operate)

| Surface | Bias |
|---------|------|
| **Persuade** (marketing, landing, hero, pitch) | Strong first-viewport hierarchy; calm symmetry often fits; one clear CTA path; sparse supporting copy |
| **Operate** (app shell, dashboard, tools, settings) | Task density and scan paths over marketing drama; focused asymmetry OK when the focal region is obvious; persistent nav and status over ornamental balance |

Do not apply persuade composition defaults to operate surfaces (or the reverse) unless the brief says so.

## How designer reports

| Finding | Report as |
|---------|-----------|
| **Must** gap (in scope) | `critical` - can block draft approval / UI ready |
| **Should** miss | `important` - polish; does **not** fail-closed alone |
| **Should** miss and brief **explicitly required** it | Treat as `critical` (brief made it Must for this mission) |

Record deliberate **Should** skips in `[data-draft-notes]` (one line each).

---

## Always apply (Must)

One actionable line each:

- **Visual Hierarchy** - One primary focal point per viewport region; size, weight, contrast, and position must make the next action obvious without hunting.
- **Mobile Friendliness & Responsiveness** - Usable at all Responsive ladder presets (375 / 768 / 1280 / 1536); structure adapts per `shared-draft-directives.md` - not a squeezed desktop clone.
- **Content Prioritization** - Lead with the job-to-be-done and primary action; defer secondary chrome, metadata, and decoration below or aside.
- **Simplicity & Whitespace** - Prefer fewer regions and structured negative space; remove competing accents, nested cards, and ornamental chrome that do not serve the task.
- **Consistency** - Same component roles, spacing steps, and interaction patterns across the surface and parent shell (when product context exists).
- **Accessibility & Readability** - Readable measure and contrast; visible focus; labeled controls; respect `prefers-reduced-motion` when motion is present.
- **Navigation** - User always knows where they are and how to move next (or exit); collapse/drawer/persistent treatment must match the active ladder preset.
- **Color Palette** - Use only tokens named in `DESIGN.md` and/or the approved brief; do not invent competing hues or "improve" the palette silently.
- **Typography System** - Use the brief / `DESIGN.md` type roles and scale only; clear display vs body vs meta; no ad-hoc font stacks.
- **Wireframe (structure, not artifact)** - Honor the brief's layout structure and the bake-off winner's structure choice. This is **not** a separate gray-box wireframe step before chrome - bake-off and approval drafts show real primary-surface chrome.

---

## Prefer when they fit (Should)

Tools, not dogma. Skip when the brief or surface type conflicts; note the skip.

- **Golden Ratio** - Prefer a modular type/spacing scale in roughly the 1.25–1.618 range for rhythm. Do **not** force every region, panel, or split to φ.
- **Centric Symmetry** - Prefer on persuade / hero / marketing first viewports. Operate / dashboard may use focused asymmetry when the focal region stays clear.
- **Fibonacci sequence** - Optional spacing or scale steps when they clarify rhythm. Not decoration for its own sake.

---

## Assemble reminder

```
shared-draft-directives.md  →  DESIGN.md (if present)  →  design-principles.md  →  brief / user content
```

On conflict with anti-slop catalog defaults, defer to `references/anti-slop-catalog.md` unless `decisions.md` records an explicit exception.
