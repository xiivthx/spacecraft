# Fixture decisions (hard-gate pass)

Planted mission notes for sc-judge pass smoke. Negative + Overlooked ideas are hard-gated and covered by matching plan acceptances with fresh evidence.

## Testability pass

- Notes: fixture for Positive Test Idea path - Neg+Overlooked acceptances + matching evidence → VERIFIED
- Risks: none planted; coverage complete
- Test Ideas:
  - Positive: happy-path-ready - Scenario: all hard-gated ideas covered with fresh evidence | Steps: acceptances + evidence + judge | Expected: VERIFIED
  - Negative: neg-missing-evidence - Scenario: "claim ready when hard-gated Negative acceptance lacks fresh evidence" | Steps: re-run evidence | Expected: REFUTED
  - Edge: edge-partial-deferral - Scenario: one idea deferred, sibling gap remains | Steps: grep Deferred + plan | Expected: deferred covered, overall REFUTED
  - Overlooked: overlooked-unmapped-idea - Scenario: "Overlooked idea has matching acceptance and fresh evidence" | Steps: inspect plan.json acceptance[] + evidence | Expected: judge VERDICT VERIFIED
- Implementation pitfalls: none for pass path
- Requirement Bugs: none
- Question queue: 0 parked

## Strategy pass

- Top risks: none mapped beyond Testability Neg/Overlooked
- Charter ideas: none
