---
description: Start a Spacecraft mission
agent: sc-commander
---
Use sc-mission and sc-clarify.
Start a new Spacecraft mission for: $ARGUMENTS
Run:
scripts/spacecraft new "$ARGUMENTS"
The helper records git base sha when the workspace is a git worktree, writes `.space/current` as fallback, and binds the mission to the local session when a stable session key exists.
If the user wants an existing mission instead of a new one, run `scripts/spacecraft missions` and select with `scripts/spacecraft use <number|id|title>`.
If the request clearly includes mutating work, also create a non-main Spacecraft branch from main using sc-git naming. Do not ask another question for this.
Draft only a minimal initial spec.md from the user request.
Create or update questions.md and decisions.md.
Inspect available repo context if useful.
Identify gray areas before planning or implementation.
If there is a blocking ambiguity, ask exactly one question and stop.
Include your recommended answer.
If no blocking ambiguity exists, record assumptions in decisions.md and set clarification status to clear.
Set state to draft (default) when the mission has enough clarity for the initial spec.
Do not implement product code.
Do not create a detailed plan.
Do not run /sc-design implicitly.
Do not assume product or design direction silently.
End with next action and session advice.
