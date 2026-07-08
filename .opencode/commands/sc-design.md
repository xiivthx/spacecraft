---
description: Shape UI direction for the resolved Spacecraft mission
agent: sc-commander
---
Use sc-mission, sc-clarify, and sc-design.
Resolve the mission. Block if unsafe.

## Pre-flight checks

Read DESIGN.md and the resolved mission's spec.md, questions.md, decisions.md, and plan.json if present. If no active mission, say to use /sc-start first. If DESIGN.md is missing, create it using Orbital Console defaults.

Goal: clear the main design image before planning or implementation.

## Workflow

1. If design intent has blocking ambiguity, ask exactly one question and stop — why it matters, your recommendation, what happens if accepted.
2. When ready, invoke sc-designer as a read-only subagent to shape the UI direction. A user invocation of /sc-design is explicit permission; do not ask for separate permission.
3. Update spec.md and/or plan.json with a concise UI section: target screen/component, user goal, chosen mood/tone, visual metaphor, info hierarchy, layout structure, design primitives, visual constraints, palette/typography/art/3D/transition/animation constraints when relevant, accessibility checks, verification method.
4. Record chosen direction, rejected directions, assumptions, and open design risks in decisions.md.
5. When enough configs are chosen, synthesize into one design brief instead of keeping earlier options as packages.
6. Set mission state to draft/planned depending on progress.

## Error handling

- Do not implement product code, UI code, or add dependencies.
- If no active mission, tell user to run /sc-start first.

End with recommended next action and session advice. Prefer /sc-plan next.
