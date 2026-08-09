# Probe dimensions

Silent sweep after (or during) scenarios. File a finding only with repro - no essay dumps.

## UX

- Feedback on action (loading, success, error) visible and timely
- Destructive / irreversible actions confirm or are clearly marked
- Empty and error states explain next step
- Focus / keyboard path usable for primary flows

## UI

- Primary action discoverable; hierarchy clear
- Copy matches outcome (no misleading labels)
- Overlap, cut-off text, low contrast that blocks use → finding
- Component chrome consistent on the same surface (not draft-parity - that is sc-ux-design)

## Functional

- Happy path completes with correct observable result
- Invalid / empty / boundary inputs rejected or handled without crash
- State persists across expected refresh/navigation when required
- No silent data loss; errors reachable from UI, not only console

## Workflow

- Multi-step flows can complete end-to-end without dead ends
- Back / cancel / retry do not leave corrupt mid-state
- Auth / session expiry recoverable or clearly messaged
- Create → use (or equivalent composition) works on product origin

## Responsive

- 375 / 768 / 1280: primary path usable (no trapped scroll, untappable controls)
- No horizontal overflow that hides primary CTA
- Touch targets usable on 375 for primary actions

## Out of scope here

- Draft-parity / token match → `sc-ux-design`
- Unit assertion quality → `sc-tdd` / `sc-judge`
- Security deep hunt → `sc-security` when in mission scope
