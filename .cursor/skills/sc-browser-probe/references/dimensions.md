# Probe dimensions

Silent sweep after (or during) scenarios. File a finding only with repro - no essay dumps.

## How to sweep

1. Inventory visible routes/surfaces (ids from `surface-match.md`)
2. Score **Foundations** always
3. Score **matched packs** only
4. Mark each check `ok` | `fail` | `n/a` (capability absent) | `deferred` (not reached)
5. File a finding only on `fail` with repro

## Foundations (always)

### UX

- After click/submit: visible in-progress, then success or error (not silent)
- Destructive / irreversible: confirm or clearly marked; cancel leaves state intact
- Empty, no-results, load-error, and loading are distinct when those states exist
- Primary flow keyboard-completable; focus visible; overlays do not trap focus with no escape
- Unsaved edits: leave/refresh warns or save affordance is obvious
- Errors say what failed and how to fix; field errors sit with the field

### UI

- Primary action discoverable; hierarchy clear
- Visible labels on controls (not placeholder-only)
- Action copy matches outcome
- Overlap, clip, truncated critical text, or contrast that blocks reading → finding
- Same-surface chrome consistent (not draft-parity - that is sc-ux-design)
- Interactive states present when used: default, hover or press, focus, disabled
- Overlays (modal / drawer / menu): titled, closable (Esc or explicit close), stay in viewport
- Toast / banner readable; not the only channel for submit failure; dismissible if it persists

### Functional

- Happy path completes with correct observable result
- Empty / invalid / boundary inputs handled; no crash
- Double-submit does not create duplicates
- Required state persists across refresh/navigation
- Errors reachable in UI, not console-only
- Unknown route: product 404 or recovery, not a blank crash
- Format constraints (email, number, …) do not wipe valid sibling fields

### Workflow

- Multi-step: progress visible; back keeps data; validate before next; no dead end
- Back / cancel / retry do not leave corrupt mid-state
- Auth / session expiry recoverable or clearly messaged
- Create → use (composition) works on product origin
- Filter / search: query visible and editable; active filters clearable; filtered-empty ≠ zero-data
- Upload (if present): limits stated before attempt; progress; success and fail; bad type/size rejected in UI
- Save: inert or disabled until dirty; in-progress; saved confirmation

### Responsive + live a11y

- 375 / 768 / 1280 usable (add 1536 if multi-region)
- No horizontal overflow that hides primary CTA
- Touch targets usable on 375 for primary actions
- 200% zoom: primary path still usable
- If motion is essential to meaning and reduced-motion is detectable, essential info still available without it

## Out of scope here

- Draft-parity / token match → `sc-ux-design`
- Unit assertion quality → `sc-tdd` / `sc-judge`
- Security deep hunt → `sc-security` when in mission scope
- Consult site / Figma AI as a gate source
