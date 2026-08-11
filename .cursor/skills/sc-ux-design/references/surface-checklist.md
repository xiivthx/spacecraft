# Surface UX checklist adapter

Discuss / designer process. Item SoT: `references/checklists/` (index: `references/checklists/README.md`).

## When to use

Visual UI/FE discuss and run. Skip when `UI draft skipped:` / non-visual.

**Must** record exactly one line in `decisions.md` before layout bake-off (and before `UI draft approved`):

```
UX checklist: <id>
UX checklist: none - <reason>
```

`<id>` is a README alias or `category/slug`. Pick the main user job, not every widget.

**Must not** walk the catalog. **Must not** record more than one id.

## Workflow

1. **Match** - Open `references/checklists/README.md`. One id or `none`.
2. **Record** - `UX checklist: <id>` or `UX checklist: none - <reason>`.
3. **Read** - That one file under `references/checklists/` (skip when `none`).
4. **Fold** - Applicable `- [ ]` titles into `spec.md` Must and draft chrome. `n/a` if the product lacks that capability. Tips are hints, not gates.
5. **Draft** - Approval candidate shows each applicable item as real chrome.
6. **Verify** - `sc-designer` scores those items (discuss: draft; run: live product).

## Score

Per README. Discuss/designer: `present` | `missing` | `n/a`.

- `(state)` `missing` = **critical**
- chrome/path `missing` = **important**
- No `UX checklist:` line = **important** (Commander adds id or none before `UI draft approved`)
- Bake-off candidates: do not score

## Designer output

```
surface-checklist (<id>):
- <item>: present | missing | n/a - <one line>
```

Then the gates snippet (`pass` | `fail` | `uncertain` | `n/a`). `fail` if any applicable item is `missing`. `n/a` when the line is `none` or the work is non-visual.

## Must / Must not

- **Must**: Record `UX checklist:` before bake-off; Read that one file; fold and score `- [ ]` items
- **Must not**: Walk the catalog or require more than one id
- **Must not**: Use an external checklist site or its AI review as a gate
- **Must not**: Score bake-off candidates

## Related

- `checklists/README.md` - item SoT
- `SKILL.md` / `ux-ui-review-gates.md` / `.cursor/agents/sc-designer.md`
- `shared-draft-directives.md` - scenario matrix owns generic states
- `sc-browser-probe/references/surface-match.md` - live multi-pack
