---
description: Resolve the next blocking Spacecraft mission question
agent: sc-commander
---
Use sc-mission and sc-clarify.
Resolve the mission. Block if unsafe.

## Pre-flight checks

Read the resolved mission's mission.json, spec.md, questions.md, decisions.md, and plan.json if present. Create questions.md or decisions.md with standard headings if either is missing. Inspect repo files if they can answer the ambiguity.

If the ambiguity is visual mood/theme/layout/art/3D/transition/animation, prefer /sc-design — it can create HTML comparisons when seeing is clearer than text.

## Workflow

1. Ask exactly one blocking question at a time if one is open in questions.md. Include why it matters, your recommendation, and what happens if accepted.
2. Record answered questions in questions.md and decisions.md.
3. When no blocking questions remain, run:
   ```
   scripts/spacecraft clarify-status clear
   ```
   If unavailable, update mission.json directly.
4. Summarize confirmed decisions briefly.
5. Recommend next command: /sc-design, /sc-plan, or /sc-work only if already planned.
6. Recommend whether to continue this chat or start a new session.

## Error handling

- Do not implement product code.
- Do not create or modify product files.
- If the clarify-status command is unavailable, update mission.json manually.

End with recommended next action and session advice.
