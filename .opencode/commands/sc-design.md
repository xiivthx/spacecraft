---
description: Shape UI direction for the resolved Spacecraft mission
agent: sc-commander
---
Use sc-mission, sc-clarify, and sc-design.
Load sc-design skill for the full workflow; this command provides dispatch only.
Resolve the mission. Block if unsafe.
Read DESIGN.md and the resolved mission's spec.md, questions.md, decisions.md, and plan.json if present.
If no active mission, say to use /sc-start first.
If DESIGN.md is missing, create it using Orbital Console defaults.

Goal: clear the main design image before planning or implementation.

If design intent has blocking ambiguity, ask exactly one question and stop — why it matters, your recommendation, what happens if accepted.

When ready, invoke sc-designer as a read-only subagent to shape the UI direction.
A user invocation of /sc-design is explicit permission; do not ask for separate permission.

Then update spec.md and/or plan.json with a concise UI section:
target screen/component, user goal, chosen mood/tone, visual metaphor, info hierarchy, layout structure, design primitives, visual constraints, palette/typography/art/3D/transition/animation constraints when relevant, accessibility checks, verification method.
Record chosen direction, rejected directions, assumptions, and open design risks in decisions.md.
When enough configs are chosen, synthesize into one design brief instead of keeping earlier options as packages.

Do not implement product code, UI code, or add dependencies.
Set mission state to draft/planned depending on progress.
End with recommended next action and session advice. Prefer /sc-plan next.
