# Live surface packs

House-owned live packs for `/sc-browser-probe`. Consult URLs are optional reading - **house items are the score source**.

Do **not** copy, quote, or vendor external checklist item text into the repo or the probe report. Record the consult URL only.

## When to use

After inventory of the **running product**. Skip packs that are not visible in this scope.

**Must not** walk the whole catalog. Inventory first. One product can match several packs.

## Match

| Id | When | Consult |
|----|------|---------|
| `login` | Sign-in / session start | https://www.checklist.design/web-app/login |
| `sign-up` | Create account | https://www.checklist.design/website/sign-up |
| `form-submit` | Primary job is submit (create/save/contact) | https://www.checklist.design/flows/submitting-a-form |
| `input-error` | Validation-heavy form | https://www.checklist.design/flows/showing-input-error |
| `empty-state` | Collection first-run / no rows | https://www.checklist.design/web-app/empty-state |
| `settings` | Preferences / account settings | https://www.checklist.design/web-app/settings |
| `search-results` | Search or filter results | https://www.checklist.design/web-app/search-results |
| `multi-step` | Wizard / stepped form | https://www.checklist.design/web-app/multi-step-form |
| `save-dirty` | Edit + save existing record | https://www.checklist.design/flows/saving-changes |
| `filter` | Collection filters | https://www.checklist.design/flows/filtering-items |
| `upload` | File / media upload | https://www.checklist.design/flows/uploading-media |
| `auth-recovery` | Forgot / reset password | https://www.checklist.design/flows/resetting-password |
| `destructive` | Delete account or irreversible destroy | https://www.checklist.design/flows/deleting-account |
| `notifications` | In-app notification list/inbox | https://www.checklist.design/web-app/notifications |
| `onboarding` | First-run setup | https://www.checklist.design/web-app/onboarding |
| `detail` | Single item detail | https://www.checklist.design/web-app/single-item-detail |
| `not-found` | Unknown URL / 404 | https://www.checklist.design/website/404 |
| `nav-tabs` | Persistent tab/bottom nav | https://www.checklist.design/mobile/tab-bar-navigation |
| `overlay` | Modal, drawer, or action sheet in the flow | https://www.checklist.design/design-system/modal |
| `table` | Data table | https://www.checklist.design/design-system/table |

Discuss ids (`login` … `search-results`) overlap `sc-ux-design/references/surface-checklist.md` by id only. That file is draft/designer gated (one id). This file is a live multi-pack sweep.

## Severity

Per applicable house item, emit `ok` | `fail` | `n/a` | `deferred`:

- Missing **(state)** → finding **critical**
- Missing chrome/path → finding **important**
- `n/a` when the product lacks the capability
- `deferred` when the item was not reached (timebox or blocked mid-sweep)

**State-like** items are marked `(state)`.

## Packs

### login

- Identity and password fields with visible labels
- Show/hide when the password is masked
- Submit labeled as sign-in
- Path to recover access when password auth exists (chrome; `n/a` otherwise)
- Path to sign-up when the product has it (chrome; `n/a` otherwise)
- Auth failure visible without leaking whether the identifier exists, unless the spec requires otherwise (state)
- In-progress state that blocks double submit (state)

### sign-up

- Surface purpose is create-account, not sign-in
- Required fields with visible labels
- Path back to login (chrome)
- Password rules visible before fail (chrome)
- In-progress after submit (state)
- Success or explicit next step (state)
- Field-level and submit-level failures visible (state)

### form-submit

- Primary submit control with an action-specific label
- In-progress after submit; control cannot double-fire (state)
- Success confirmation (state)
- Submit-level failure (state)
- Field-level invalid shown at the field (state)

### input-error

- Do not show field errors while the user is still typing the first value
- After blur or submit, error at the field (state)
- Error says how to fix
- First invalid focused or scrolled into view (chrome)
- Valid fields stay filled (chrome)
- Error clears or updates on re-attempt (state)

### empty-state

- Explains why empty (zero-data vs no-matches) (state)
- Next action when create or filter is possible (chrome; `n/a` if neither)
- Distinct from loading and from error (state)
- Not a blank hole

### settings

- Related settings grouped; each control has a visible label
- Unsaved changes visible before leave or save (chrome)
- Save success and save-failure (state)
- Destructive actions confirm (`n/a` if none)
- Loaded values match persisted state

### search-results

- Query visible and editable
- Results, no-matches, and error are distinct (state)
- Count or an explicit no-matches message (chrome)
- Loading distinct from empty (state)
- A result is actionable

### multi-step

- Progress (which step / how many)
- Step title
- Validate before next (state)
- Back keeps later-step data
- Save/resume (`n/a` if product has no draft)
- Review-before-final (`n/a` if a single confirm is enough)

### save-dirty

- Edit affordance
- Save inert until dirty
- Save goes in-progress (state)
- Saved confirmation (state)
- Leave-with-dirty warned or blocked (chrome)

### filter

- Control near the collection
- Options usable
- Active filters visible
- Clear one and clear all
- Result count (optional chrome)
- Filtered-empty explains adjust/clear (state)

### upload

- Empty drop/click target
- Limits (type/size) stated before attempt
- Progress (state)
- Success and fail with how to fix (state)
- Remove/retry action
- Multi-file layout if multi allowed (`n/a` if single)

### auth-recovery

- Reset entry next to the password field
- Identity field to send the reset
- Sent confirmation without account-enumeration when the spec requires that
- New-password field + rules
- Success then path to login (state)

### destructive

- Reachable without hunting
- Consequences explained before confirm
- Confirm is explicit
- Success confirmation (state)
- Cancel leaves account/data intact

### notifications

- List readable
- Read vs unread (state)
- Empty list distinct (state)
- Mark-all (`n/a` if none)
- Item action if the product has one

### onboarding

- Progress if multi-step
- First action clear
- Can skip or finish into the product (`n/a` if forced)
- Completion acknowledgement (state)

### detail

- Item identity clear
- Primary actions reachable
- Missing/unknown id is error, not blank (state)
- Back returns to list/context

### not-found

- Page states it is unknown/missing (state)
- Recovery link (home or search)
- Recovery lands on a real product surface
- Not an uncaught crash

### nav-tabs

- 3-5 destinations typical
- Icon+label or a clear label
- Active destination distinct
- Target tappable on 375
- Stays put on section roots

### overlay

- Title
- Primary action
- Close / Esc
- Focus moves in and restores
- Viewport fit on 375
- Backdrop does not cover the close control

### table

- Headers
- Row actions if rows are actionable
- Sort/filter (`n/a` if none)
- Overflow strategy on 375 (stack / scroll / cards - not clipped unreadably)
- Pagination (`n/a` if short list)

## Must / Must not

- **Must**: Inventory first; score only matched packs
- **Must**: Score each applicable house item `ok` | `fail` | `n/a` | `deferred`
- **Must**: File a finding on `fail` with repro
- **Must not**: Walk packs that are not in inventory
- **Must not**: Copy external checklist wording
- **Must not**: Use the consult site, Figma plugin, or its AI review as a pass/fail source
- **Must not**: Treat this as the discuss one-id draft gate (`surface-checklist.md`)

## Related

- `../SKILL.md` - probe workflow and verdict
- `dimensions.md` - Foundations (always)
- `scenario-matrix.md` - extra buckets when a pack is present
- `sc-ux-design/references/surface-checklist.md` - discuss/designer one-id (not this job)
