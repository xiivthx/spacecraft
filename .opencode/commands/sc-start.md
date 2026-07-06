---
description: Start a Spacecraft mission
agent: sc-commander
---
Use sc-mission and sc-clarify.
Start a new Spacecraft mission for: $ARGUMENTS
Run:
node scripts/spacecraft.mjs new "$ARGUMENTS"
The helper records git base sha when the workspace is a git worktree.
Draft only a minimal initial spec.md from the user request.
Create or update questions.md and decisions.md.
Inspect available repo context if useful.
Identify gray areas before planning or implementation.
If there is a blocking ambiguity, ask exactly one question and stop.
Include your recommended answer.
If no blocking ambiguity exists, record assumptions in decisions.md and set clarification status to clear.
Set state to specified only when the mission has enough clarity for the initial spec.
Do not implement product code.
Do not create a detailed plan.
Do not run /sc-design implicitly.
Do not assume product or design direction silently.
End with the recommended next action and whether the user should continue this chat or start a new session.
