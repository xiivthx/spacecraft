# Fixture decisions (hard-gate deferral edge)

Planted mission notes for sc-judge deferral-edge smoke. overlooked-x is deferred; sibling hard-gated idea stays uncovered.

## Testability pass

- Notes: fixture for Edge Test Idea path - Deferred covers overlooked-x while sibling hard-gate gap remains
- Risks: deferral treated as blanket VERIFIED for all hard-gated ideas
- Test Ideas:
  - Positive: happy-path-ready - Scenario: all hard-gated ideas covered | Steps: acceptances + evidence | Expected: VERIFIED
  - Negative: neg-sibling-uncovered - Scenario: "sibling Negative remains without acceptance or Deferred line while overlooked-x is deferred" | Steps: inspect plan.json acceptance[] and decisions.md | Expected: judge VERDICT REFUTED overall
  - Edge: edge-partial-deferral - Scenario: Deferred covers overlooked-x, sibling gap remains | Steps: grep Deferred + plan | Expected: deferred covered, overall REFUTED
  - Overlooked: overlooked-x - Scenario: "hard-gated Overlooked deferred out of size with no acceptance" | Steps: grep Deferred test idea: overlooked-x | Expected: that idea covered by deferral
- Implementation pitfalls: treating one Deferred line as coverage for every hard-gated idea
- Requirement Bugs: none
- Question queue: 0 parked

## Deferral

Deferred test idea: overlooked-x - out of size

## Strategy pass

- Top risks: none mapped beyond Testability hard-gated set
- Charter ideas: none
