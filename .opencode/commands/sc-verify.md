---
description: Capture verification evidence for the current mission
agent: sc-commander
---
Use sc-verification.
Read plan.json and acceptance checks.
Run focused verification using:
node scripts/spacecraft.mjs evidence "<label>" -- <command>
Capture evidence for tests, build, typecheck, or lint as appropriate.
Capture failures too.
Update plan.json task evidence references when a check maps to a task.
Run:
node scripts/spacecraft.mjs validate
Use rtk for noisy verification commands when available. Use raw output or `rtk proxy` passthrough when exact evidence is needed.
Do not say verified without evidence ids.
End with next action and session advice. Suggest /sc-design-review if UI changed, otherwise /sc-review or /sc-work.
