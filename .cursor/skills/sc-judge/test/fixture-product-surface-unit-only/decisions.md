# Fixture decisions (product-surface unit-only)

Planted mission notes for sc-judge product-surface smoke. UI task claims user-visible acceptance with unit-only verify.

## Testability pass

- Notes: fixture for Overlooked Test Idea path - user-visible UI acceptance + unit-only verify
- Risks: ready claimed without product-surface marker
- Test Ideas:
  - Overlooked: overlooked-unit-only-ui - Scenario: "UI task claims user-visible acceptance with unit-only verify (no product-surface marker)" | Steps: inspect plan verify/acceptance for verify.product|browser|curl|composition | Expected: judge VERDICT REFUTED
- Implementation pitfalls: treating unit test mention as product-surface proof
- Requirement Bugs: none
- Question queue: 0 parked
