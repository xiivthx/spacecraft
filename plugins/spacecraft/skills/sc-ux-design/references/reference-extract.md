# Reference extract (on-demand gate)

Run when the human supplies reference image(s), screenshots, or URLs during `/sc-discuss` visual work. Produces a structured extract artifact **before** the design brief so borrow scope and layout cues are mechanical, not improvised from memory.

## When to run

- Human supplies mood boards, screenshots, product shots, competitor URLs, or other visual references.
- **Must** run before the design brief when any reference is supplied.
- **Must not** enter layout bake-off when `Reference borrow:` is recorded in `decisions.md` but `design/refs/extract.md` is missing.

## Output artifact

Write `.space/missions/<id>/design/refs/extract.md` (create `design/refs/` when needed). One row or section per reference source.

### Required fields (per source)

| Field | What to capture |
|-------|-----------------|
| **Source** | Path under `design/refs/` or URL |
| **Tokens** | Color, type, spacing cues visible in the ref (hex/names, font families, density rhythm) |
| **Layout regions** | Header, nav, main, aside, footer; column count; denseness (airy vs compact) |
| **Chrome inventory** | Buttons, fields, tables, cards, empty/error patterns, badges, tabs - what is visible |
| **Responsive cues** | Per visible breakpoint: mobile (375), tablet (768), desktop (1280), widescreen (1536) - stack vs side-by-side, drawer vs persistent nav, column count, content measure, density; note when ref shows only a subset |
| **Borrow recommendation** | Proposed scope: `mood` \| `tokens` \| `layout` \| `chrome` (human still picks exactly one) |

When multiple refs are supplied, group by source. Note conflicts between refs in a short **Conflicts** subsection.

### Example extract (abbreviated)

```markdown
# Reference extract

## Ref 1 - design/refs/dashboard-shot.png

| Field | Notes |
|-------|-------|
| Source | design/refs/dashboard-shot.png |
| Tokens | Dark bg `#0f1117`, surface `#1a1d27`, accent cyan `#22d3ee`, 14px system sans, 8px spacing grid |
| Layout regions | Persistent left sidebar (~240px), top bar with search, main 2-column (list + detail) |
| Chrome inventory | Icon nav items, pill tabs, data table with zebra rows, inline empty state with CTA |
| Responsive cues | Not visible in ref (desktop-only shot) |
| Borrow recommendation | `layout` |

## Conflicts

None.
```

## Workflow

1. **Collect** - Save human-supplied images under `design/refs/` when not already on disk; record URLs as-is.
2. **Extract** - Read each source mechanically; fill the table fields from what is visible. Do not invent product claims, metrics, or features not shown in the ref or stated in the brief/spec.
3. **Recommend borrow** - Propose exactly one borrow scope per ref (or one consolidated recommendation when refs agree). Human confirms with `Reference borrow: <scope>` in `decisions.md`.
4. **Record** - Add `Reference extract: design/refs/extract.md` to `decisions.md`.
5. **Brief** - Design brief **Must** cite relevant extract rows (tokens, layout regions, chrome) when forming the 6 dimensions. Borrow scope in the brief must match the human-confirmed `Reference borrow:` line.

## Must

- Produce `extract.md` before the design brief when references are supplied.
- Cite extract rows in the design brief.
- Record `Reference extract: design/refs/extract.md` in `decisions.md`.
- Honor the human-confirmed borrow scope only - extract proposes; human picks.

## Must not

- Silent full-chrome clone of a reference (borrow scope still governs what may be ported).
- Invent product claims, brand lines, metrics, or UI states not visible in the ref or supplied in brief/spec.
- Skip extract and jump to bake-off when `Reference borrow:` is set.
- Enter bake-off without the extract artifact when references were supplied.

## Handoff to brief and bake-off

- Brief **Reference borrow** line must align with `decisions.md` and cite extract rows.
- Before bake-off, record **Context fidelity** (see `sc-ux-design` SKILL.md): `Context fidelity: DESIGN.md | shell:<path> | extract:design/refs/extract.md | product-shot:<path>` (omit absent paths).
- Bake-off candidates must respect borrow scope and product context (shell/nav when brownfield).
