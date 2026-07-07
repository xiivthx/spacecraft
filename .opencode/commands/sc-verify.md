---
description: Capture verification evidence for the resolved mission
agent: sc-commander
---
Use sc-verification.
Run:
scripts/spacecraft resolve --json
If resolver safety is not `safe` or no mission is selected, stop before capturing evidence. Show the conflict/candidates and tell the user to run `scripts/spacecraft missions` then `scripts/spacecraft use <number|id|title>`, or set `SPACECRAFT_MISSION=<mission-id>` for one command.
Treat `.space/current` as fallback state, not sole authority.
Read the resolved mission's plan.json and acceptance checks.
Run focused verification using:
scripts/spacecraft evidence "<label>" -- <command>
Capture evidence for tests, build, typecheck, or lint as appropriate.
Capture failures too.
Update plan.json task evidence references when a check maps to a task.
Run:
scripts/spacecraft validate
Use rtk for noisy verification commands when available. Use raw output or `rtk proxy` passthrough when exact evidence is needed.
Do not say verified without evidence ids.
End with next action and session advice. Suggest /sc-design-review if UI changed, otherwise /sc-flow when tasks remain, or /sc-review when implementation is complete.
