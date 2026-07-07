---
description: Create or update resolved mission flight plan
agent: sc-commander
---
Use sc-mission, sc-clarify, and sc-planning.
Resolve the mission. Block if unsafe.
Read the resolved mission's spec.md, questions.md, decisions.md, and plan.json if present.
Use sc-clarify before finalizing a plan.
If mission clarification status is open and there are blocking questions, stop and tell the user to run /sc-clarify or answer the current question.
Do not finalize plan.json while blocking clarification remains open.
Non-blocking assumptions must be recorded in decisions.md.
If the mission includes UI, use sc-design and read DESIGN.md.
If UI art direction is not chosen, stop and recommend /sc-design before finalizing UI tasks.
Invoke sc-planner as a read-only subagent to draft the plan.
A user invocation of /sc-plan is explicit permission to use the read-only sc-planner subagent; do not ask for separate subagent permission.
Then write or update the resolved mission plan.json yourself.
The plan must contain no more than 7 tasks.
Each task must have id, title, status, files, acceptance, verify, and evidence.
UI tasks must include visual intent, target component/screen, accessibility checks, and verification method.
For new screens, recommend /sc-design before /sc-work.
Set state to planned.
Do not implement product code.
End with the recommended next action and session advice. Recommend `/sc-git` then `/sc-flow` when implementation can begin. Recommend a new session if implementation is the next large phase and the plan is fully captured in plan.json.
