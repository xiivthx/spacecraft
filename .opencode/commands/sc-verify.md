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
Do not say verified without evidence ids.
End with the recommended next action and session advice. Suggest /sc-design-review if UI changed, otherwise /sc-review or /sc-work for remaining tasks.
