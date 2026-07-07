---
description: Resolve the next blocking Spacecraft mission question
agent: sc-commander
---
Use sc-mission and sc-clarify.
Run:
node scripts/spacecraft.mjs resolve --json
If resolver safety is not `safe` or no mission is selected, stop before writing clarification artifacts. Show the conflict/candidates and tell the user to run `node scripts/spacecraft.mjs missions` then `node scripts/spacecraft.mjs use <number|id|title>`, or set `SPACECRAFT_MISSION=<mission-id>` for one command.
Treat `.space/current` as fallback state, not sole authority.
Read the resolved mission's mission.json, spec.md, questions.md, decisions.md, and plan.json if present.
Create questions.md or decisions.md with the standard headings if either is missing.
Inspect repo files if they can answer the ambiguity.
If the ambiguity is visual mood, theme, layout feel, art direction, 3D, transition, or animation, prefer /sc-design because it can create HTML comparison artifacts when seeing options is clearer than a text question.
Do not implement product code.
Do not create or modify product files.
If there is an unanswered blocking question already in questions.md, ask that one question again with:
- why it matters
- your recommended answer
- what happens if the recommendation is accepted
If the user has answered a previous question in the conversation, record the answer in questions.md and decisions.md.
Then determine whether another blocking question remains.
If yes, ask exactly one next blocking question and stop.
If no blocking questions remain, set clarification status to clear using:
node scripts/spacecraft.mjs clarify-status clear
If the helper is unavailable, update mission.json directly.
Summarize the current confirmed decisions briefly.
Recommend the next command:
- /sc-design if UI/design direction is needed
- /sc-plan if ready to plan
- /sc-work only if already planned
Also recommend whether to continue this chat or start a new session. Prefer continuing this chat when the user should answer the current question.
