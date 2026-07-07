---
description: Resolve the next blocking Spacecraft mission question
agent: sc-commander
---
Use sc-mission and sc-clarify.
Load sc-clarify skill for the detailed workflow; this command provides dispatch only.
Resolve the mission. Block if unsafe.
Read the resolved mission's mission.json, spec.md, questions.md, decisions.md, and plan.json if present.
Create questions.md or decisions.md with standard headings if either is missing.
Inspect repo files if they can answer the ambiguity.
If the ambiguity is visual mood/theme/layout/art/3D/transition/animation, prefer /sc-design — it can create HTML comparisons when seeing is clearer than text.
Do not implement product code.
Do not create or modify product files.
Ask exactly one blocking question at a time if one is open in questions.md. Include why it matters, your recommendation, and what happens if accepted.
Record answered questions in questions.md and decisions.md.
When no blocking questions remain, run:
scripts/spacecraft clarify-status clear
If unavailable, update mission.json directly.
Summarize confirmed decisions briefly.
Recommend next command: /sc-design, /sc-plan, or /sc-work only if already planned.
Recommend whether to continue this chat or start a new session.
