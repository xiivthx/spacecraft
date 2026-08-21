# Fixture decisions (hard-gate omit)

Planted mission notes for sc-judge omit smoke. Overlooked idea is hard-gated and intentionally uncovered.

## Testability pass

- Notes: fixture for Negative Test Idea path - Overlooked present, plan omits coverage
- Risks: ready claimed without hard-gate coverage
- Test Ideas:
  - Positive: happy-path-ready - Scenario: all hard-gated ideas covered | Steps: acceptances + evidence | Expected: VERIFIED
  - Negative: neg-false-ready - Scenario: claim ready with failing evidence | Steps: re-run evidence | Expected: REFUTED
  - Edge: edge-partial-deferral - Scenario: one idea deferred, sibling gap remains | Steps: grep Deferred + plan | Expected: deferred covered, overall REFUTED
  - Overlooked: overlooked-orphan-claim - Scenario: "claim ready while Overlooked idea has no acceptance or Deferred line" | Steps: inspect plan.json acceptance[] and decisions.md for Deferred test idea: | Expected: judge VERDICT REFUTED
- Implementation pitfalls: soft SFDIPOT language treating Overlooked as optional aid
- Requirement Bugs: none
- Question queue: 0 parked

## Strategy pass

- Top risks: none mapped beyond Testability Overlooked
- Charter ideas: none
