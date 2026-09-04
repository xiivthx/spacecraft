# Live surface packs

Live multi-pack adapter for `/sc-browser-probe`. Item SoT: `.cursor/skills/sc-ux-design/references/checklists/` (index: `README.md` there). Foundations stay in `dimensions.md`. Probe-owned optional pack: `persona-walkthrough.md` (not in the UX checklist catalog).

## When to use

After inventory of the **running product**. Skip packs that are not visible.

**Must not** walk the catalog. Inventory first. One product can match several packs.

Discuss is one primary id (`sc-ux-design/references/surface-checklist.md`). This file is the live multi-pack sweep.

## Match

1. Inventory visible surfaces.
2. Resolve each hit via README aliases or `category/slug`.
3. **Read only those files.**
4. Also Read a `design-system/` file when that chrome is on screen: modal, drawer, table, button, input, toast, loading. Do not load every `design-system/` file.
5. If `persona-walkthrough` match rules fire (`references/persona-walkthrough.md`), Read that file and score `pack:persona-walkthrough`. Do not treat it as a checklist catalog id.

## Score

Per catalog README. Probe: `ok` | `fail` | `n/a` | `deferred`.

- `(state)` `fail` → finding **critical**
- chrome/path `fail` → finding **important**
- `n/a` when the product lacks the capability
- `deferred` when not reached (timebox or blocked)

Persona pack: score per `persona-walkthrough.md` (findings carry severity; no 1-5 scores).

Tips are hints, not gates. File a finding on `fail` with repro. Do not invent a second item list.

## Must / Must not

- **Must**: Inventory first; score only matched files; finding + repro on `fail`
- Must: When dialog/modal/drawer inventoried, overlay pack is required
- **Must not**: Walk unmatched packs or every `design-system/` file
- **Must not**: Use an external checklist site or its AI review as pass/fail
- **Must not**: Treat this as the discuss one-id draft gate
- **Must not**: Auto-match `persona-walkthrough` on every probe

## Related

- `.cursor/skills/sc-ux-design/references/checklists/README.md`
- `persona-walkthrough.md` - optional cognitive walkthrough pack
- `../SKILL.md` / `dimensions.md` / `scenario-matrix.md`
- `sc-ux-design/references/surface-checklist.md` - discuss one-id
