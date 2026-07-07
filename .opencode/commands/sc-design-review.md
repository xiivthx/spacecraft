---
description: Run read-only design critique for current UI work
agent: sc-commander
---
Use sc-mission and sc-design.
Run:
node scripts/spacecraft.mjs resolve --json
If resolver safety is not `safe` or no mission is selected, stop before design review. Show the conflict/candidates and tell the user to run `node scripts/spacecraft.mjs missions` then `node scripts/spacecraft.mjs use <number|id|title>`, or set `SPACECRAFT_MISSION=<mission-id>` for one command.
Treat `.space/current` as fallback state, not sole authority.
Read DESIGN.md, the resolved mission's spec.md, plan.json, decisions.md, design artifacts if present, evidence.jsonl if present, and git diff.
Invoke sc-designer as a read-only subagent.
A user invocation of /sc-design-review is explicit permission to use the read-only sc-designer subagent; do not ask for separate subagent permission.
The subagent must not edit files.
Ask the subagent to review:
- hierarchy
- layout
- typography
- spacing
- color use
- interaction states
- accessibility
- responsiveness
- anti-slop checklist
- Feynman clarity: plain-language explanation, labeled visuals, obvious gain/tradeoff, and no unnecessary jargon
- visual economy: HTML was used only when visual comparison materially helped the decision
- reference fit: implementation learns from selected references without cloning them
If UI art direction was selected, check whether the implementation matches it.
If selected references exist in decisions.md or design/references.md, compare the UI against them for hierarchy, layout rhythm, density, palette discipline, art direction, and interaction feel.
If reviewing design options or artifacts, flag same-y option sets where choices share the same skeleton and differ only by palette, labels, or microcopy.
If reviewing a design artifact for a Thai or multilingual mission, check that user-facing copy is Thai-first with simple English support and is easy for the user to choose from.
If reviewing an HTML artifact, flag verbose copy, theory dumps, option cards with more than 3 visible bullets, paragraphs longer than 2 short lines, unlabeled visuals, decorative visuals that do not explain the choice, captions that could fit every option, and artifacts that should have been a normal chat question.
When visual review needs a browser and no app-specific dev server is running, use:
node .opencode/skills/sc-design/scripts/serve-html.mjs <artifact-or-dir> --open
If reviewing design decisions, check that independent configs were recorded separately enough to allow mixing, for example layout choice, palette choice, typography choice, art/3D choice, and motion choice.
After the subagent responds, record a concise design review section in review.md.
If review.json exists, add design findings to review.json using severity:
critical, important, minor.
Critical design findings block /sc-ship.
Do not implement fixes in this command unless the user explicitly asks.
End with the recommended next action and session advice. Suggest /sc-polish for small UI cleanup or /sc-review when design findings are clear.
