# Surface UX checklist adapter

House-owned completeness items for common visual surfaces. Use during `/sc-discuss` (spec + draft) and `/sc-run` / `sc-designer` (verify). Consult URLs are optional reading with attribution - **house items are the gate source**.

Do **not** copy, quote, or vendor external checklist item text into this repo, `spec.md`, drafts, or `decisions.md`. Record the consult URL only.

## When to use

Visual UI/FE discuss and run. Skip when `UI draft skipped:` / non-visual.

**Must** record exactly one line in `decisions.md` before layout bake-off (and before `UI draft approved`):

```
UX checklist: <id>
UX checklist: none - <reason>
```

`<id>` is one primary surface from the catalog below. Pick the main user job, not every widget on the page.

**Must not** walk the whole catalog. **Must not** treat a consult URL as a pass/fail source.

## Match

| Id | Match when the primary surface is… |
|----|-------------------------------------|
| `login` | Sign-in / session start |
| `sign-up` | Create-account / register |
| `form-submit` | A form whose primary job is submit (create, save, contact) and no more specific id fits |
| `input-error` | Validation-heavy form where invalid input is the main risk (prefer `form-submit` when both fit) |
| `empty-state` | Collection (list, table, inbox, grid) whose first-run / no-rows treatment matters |
| `settings` | Preferences / account / notification settings |
| `search-results` | Search or filter results |

Use `none` for read-only marketing, charts-only dashboards, or surfaces outside this table. One id only.

## Workflow

1. **Match** - From `spec.md` / brief, pick one id or `none`.
2. **Record** - `UX checklist: <id>` or `UX checklist: none - <reason>` in `decisions.md`. Optional: `UX checklist consult: <url>` (the row's Consult URL).
3. **Fold** - Copy applicable **house items** into `spec.md` Must and into draft `data-state` / chrome. Mark an item `n/a` when the product lacks that capability (e.g. no sign-up → skip login's create-account path).
4. **Draft** - Approval candidate shows each applicable item as real chrome (not a note-only promise).
5. **Verify** - `sc-designer` scores the **surface-checklist** dimension against house items (discuss: draft; run: live product).

## Severity

Per applicable house item, emit `present` | `missing` | `n/a`:

- **State-like** items (in-progress, success, submit/field failure, empty vs loading) - `missing` = **critical** (same bar as scenario coverage).
- **Chrome/path** items (labels, recovery links, grouping, result count) - `missing` = **important**.
- Visual discuss with no `UX checklist:` line before approval - **important**; Commander **Must** add `id` or `none` before `UI draft approved`.
- Bake-off candidates: do not score this dimension.

## Catalog

Each item is house-owned. **State-like** items are marked `(state)`.

### login

Consult: https://www.checklist.design/website/login-page

- Identity field (email or username) with a visible label
- Password field; offer show/hide when the value is masked
- Submit control labeled as sign-in, not a generic Submit
- Path to recover access when a password can be forgotten (chrome)
- Path to create an account when sign-up exists (chrome; `n/a` otherwise)
- Auth failure visible without leaking whether the identifier exists, unless the spec requires otherwise (state)
- In-progress state that blocks double submit (state)

### sign-up

Consult: https://www.checklist.design/website/sign-up

- Surface purpose is create-account, not sign-in
- Required identity and credential fields with visible labels
- Path back to login when accounts already exist (chrome)
- Password rules visible before submit, not only after failure (chrome)
- In-progress state after submit (state)
- Success or explicit next step after create (state)
- Field-level and submit-level failures visible (state)

### form-submit

Consult: https://www.checklist.design/flows/submitting-a-form

- Primary submit control with an action-specific label
- In-progress state after submit; control cannot double-fire (state)
- Success confirmation (inline, toast, or next step) (state)
- Submit-level failure message (state)
- Field-level invalid input shown at the field (state)

### input-error

Consult: https://www.checklist.design/flows/showing-input-error

- Error text adjacent to the invalid field (state)
- Error says what failed and how to fix it
- First invalid field is focused or scrolled into view on submit (chrome)
- Valid fields stay filled after a failed submit (chrome)
- Error is grouped with its control (visible grouping or name/description association)

### empty-state

Consult: https://www.checklist.design/web-app/empty-state

- Explains why the surface is empty (no data vs no matches) (state)
- Next action when the user can create an item or change filters (chrome; `n/a` if neither is possible)
- Distinct from loading and from error (state)
- Not a blank hole or broken layout

### settings

Consult: https://www.checklist.design/web-app/settings

- Related settings grouped; each control has a visible label
- Unsaved changes are visible before leave or save (chrome)
- Save success and save-failure feedback (state)
- Destructive actions confirm before they run (chrome; `n/a` if none)
- Loaded values match persisted state

### search-results

Consult: https://www.checklist.design/web-app/search-results

- Query is visible and editable
- Results, no-matches, and error are distinct states (state)
- Result count or an explicit no-matches message after a query runs (chrome)
- Loading is distinct from empty (state)
- A result is openable or otherwise actionable

## Designer output

When `UX checklist: <id>` is set, list each house item:

```
surface-checklist (<id>):
- <item>: present | missing | n/a - <one line>
```

Then the gates snippet verdict (`pass` | `fail` | `uncertain` | `n/a`). `fail` if any applicable item is `missing`. `n/a` when the line is `none` or the work is non-visual.

## Must / Must not

- **Must**: Record `UX checklist: <id>` or `UX checklist: none - <reason>` before bake-off on visual discuss.
- **Must**: Fold applicable house items into spec Must and the approval draft.
- **Must**: Score house items in `sc-designer` when an id is recorded (discuss draft; run live product).
- **Must not**: Copy external checklist wording into repo artifacts.
- **Must not**: Use the consult site, Figma plugin, or its AI review as a spacecraft gate.
- **Must not**: Require more than one id, or score bake-off candidates.

## Related

- `SKILL.md` - discuss record + fold; run verify
- `ux-ui-review-gates.md` - **surface-checklist** dimension
- `.cursor/agents/sc-designer.md` - critique
- `shared-draft-directives.md` - scenario matrix still owns generic states
