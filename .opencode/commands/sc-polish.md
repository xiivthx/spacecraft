---
description: Polish UI changes against DESIGN.md before shipping
agent: sc-commander
---
Use sc-mission and sc-design.
Run:
scripts/spacecraft resolve --json
If resolver safety is not `safe` or no mission is selected, stop before polishing. Show the conflict/candidates and tell the user to run `scripts/spacecraft missions` then `scripts/spacecraft use <number|id|title>`, or set `SPACECRAFT_MISSION=<mission-id>` for one command.
Treat `.space/current` as fallback state, not sole authority.
Read DESIGN.md, the resolved mission's spec.md, plan.json, review.md, review.json, and git diff.
If the mission has no UI changes, say so and stop.
First invoke sc-designer as a read-only subagent to identify focused polish items.
A user invocation of /sc-polish is explicit permission to use the read-only sc-designer subagent; do not ask for separate subagent permission.
Then implement only small, low-risk polish changes that improve:
- spacing rhythm
- typography hierarchy
- color consistency
- focus/hover/active states
- empty/loading/error states
- accessible labels and semantics
- removal of generic AI-template patterns
Do not add dependencies.
Do not redesign the whole app.
Do not change backend behavior.
After polishing, tell the user to run /sc-verify and /sc-design-review.
Do not claim the UI is ready without verification and design review.
End with session advice. Prefer continuing this chat for immediate verification, unless the thread is context-heavy.
