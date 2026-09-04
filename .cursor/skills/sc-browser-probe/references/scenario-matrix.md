# Scenario matrix

Build a small diverse set. Prefer category coverage over many near-duplicate happies.

## Buckets (minimum when in scope)

| Bucket | Intent | Example seeds |
|--------|--------|----------------|
| Happy | Primary success path | typical valid input |
| Empty / missing | Required fields blank | submit with no data |
| Invalid | Wrong shape / type | malformed id, bad URL |
| Boundary / long | Limits and overflow | max length, very long string |
| Mobile | Same happy (or fail) at 375 | resize then repeat S1 |

Add when risk warrants: auth expiry mid-flow, slow network (if controllable), concurrent tab.

## Extra buckets (when pack present)

Do not require these on every probe. Required only when that pack is in inventory.

| Extra bucket | When | Intent |
|--------------|------|--------|
| dirty-save | `save-dirty` / `settings` | edit, attempt leave, save |
| filter-empty | `filter` / `search-results` | filters or query that yield zero |
| upload-reject | `upload` | disallowed type or oversized |
| not-found | `full` scope or routing touched | unknown URL |
| keyboard-overlay | `overlay` / modal in flow | open, Tab, Esc, close |
| double-submit | `form-submit` / `login` | mash submit |
| blur-validate | `input-error` / `form-submit` | type invalid, blur, fix |
| persona-step | `persona-walkthrough` matched | one critical journey step per active archetype (label `persona:<id>`) |

## Seeding order

1. User `examples:` list
2. Mission `decisions.md` Test Ideas / Test data design rows
3. Spec fixtures
4. Agent-generated (label `gen:` in notes)

## Table shape

| id | bucket | input | steps | expected |
|----|--------|-------|-------|----------|
| S1 | happy | … | 1. … 2. … | … |
| S2 | invalid | … | … | error visible; no crash |
| S3 | gen:filter-empty | … | … | filtered-empty ≠ zero-data |

## Rules

- **Must**: At least one row per in-scope bucket above
- **Must**: Extra rows only when that pack is in inventory
- **Must**: Record `pass` / `fail` / `blocked` per row after run
- **Must**: Label agent-generated rows `gen:`
- **Must not**: Ten tiny variants of the same happy path
- **Must not**: Skip invalid/boundary on forms or parsers without noting `scoped skip: <reason>`
