---
name: sc-clarify
description: Resolve mission ambiguity through focused user clarification before planning, design, or implementation
license: MIT
compatibility: opencode
---
- Use sc-clarify when a mission, product behavior, or design direction has meaningful ambiguity.
- Do not use sc-clarify to create bureaucracy for obvious small tasks.
- First inspect available context before asking:
  - `AGENTS.md`
  - `DESIGN.md` if UI/design is involved
  - resolved mission `spec.md`
  - `plan.json` if present
  - package/project files when relevant
  - existing source files when relevant
- If a question can be answered by inspecting the repo, inspect the repo instead of asking the user.
- Classify ambiguity:
  - blocking: cannot safely plan/design/implement without user input
  - non-blocking: can proceed with an explicit assumption
  - researchable: answer by reading files/code
- Ask exactly one blocking question at a time.
- Every question must include:
  - the question
  - why it matters
  - your recommended answer
  - what will happen if the user accepts the recommendation
- Do not ask multiple questions in one message.
- Do not implement while a blocking clarification question is open.
- Do not finalize plan or design direction while blocking clarification remains open.
- Record answered questions in `questions.md`.
- Record confirmed choices and assumptions in `decisions.md`.
- If the user says to proceed with your recommendation, record that as a confirmed decision.
- If ambiguity is low-risk, write it as an assumption and continue.
- Prefer user clarity over agent cleverness.
