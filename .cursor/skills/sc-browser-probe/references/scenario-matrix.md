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

Add when risk warrants: auth expiry mid-flow, double-submit, slow network (if controllable), concurrent tab.

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

## Rules

- **Must**: At least one row per in-scope bucket above
- **Must**: Record `pass` / `fail` / `blocked` per row after run
- **Must not**: Ten tiny variants of the same happy path
- **Must not**: Skip invalid/boundary on forms or parsers without noting `scoped skip: <reason>`
