---
name: sc-clarify
description: Resolve mission ambiguity through focused user clarification before planning, visual design, or implementation
license: MIT
compatibility: opencode
metadata:
  version: 1.0
  audience: xiivthx
---

# sc-clarify

Resolve mission ambiguity through focused user clarification before planning, visual design, or implementation.

## When to use

Activate when the user asks to:

- Clarify mission scope, product behavior, or visual design direction
- Resolve ambiguity before planning or implementation
- Record decisions and assumptions from a discussion

## Workflow

Use this exact sequence unless the user specifies otherwise:

1. **Inspect available context** — Check resolved mission `spec.md`, `plan.json` if present, package/project files, and existing source files. If a question can be answered by inspecting the repo, inspect the repo instead of asking the user.
2. **Classify ambiguity** — Determine if the ambiguity is:
   - **blocking**: cannot safely plan/design/implement without user input
   - **non-blocking**: can proceed with an explicit assumption
   - **researchable**: answer by reading files/code
3. **Ask one question** — Ask exactly one blocking question at a time. Each question must include: the question, why it matters, your recommended answer, and what will happen if the user accepts the recommendation.
4. **Record answers** — Record answered questions in `questions.md`. Record confirmed choices and assumptions in `decisions.md`. If the user says to proceed with your recommendation, record that as a confirmed decision.

## Rules

- **Must**: Use sc-clarify when a mission, product behavior, or visual design direction has meaningful ambiguity.
- **Must not**: Use sc-clarify to create bureaucracy for obvious small tasks.
- **Must**: First inspect available context before asking:
  - resolved mission `spec.md`
  - `plan.json` if present
  - package/project files when relevant
  - existing source files when relevant
- **Must**: If a question can be answered by inspecting the repo, inspect the repo instead of asking the user.
- **Must**: Classify ambiguity:
  - blocking: cannot safely plan/design/implement without user input
  - non-blocking: can proceed with an explicit assumption
  - researchable: answer by reading files/code
- **Must**: Ask exactly one blocking question at a time.
- **Must**: Every question must include:
  - the question
  - why it matters
  - your recommended answer
  - what will happen if the user accepts the recommendation
- **Must not**: Ask multiple questions in one message.
- **Must not**: Implement while a blocking clarification question is open.
- **Must not**: Finalize plan or visual design direction while blocking clarification remains open.
- **Must**: Record answered questions in `questions.md`.
- **Must**: Record confirmed choices and assumptions in `decisions.md`.
- **Must**: If the user says to proceed with your recommendation, record that as a confirmed decision.
- **Must**: If ambiguity is low-risk, write it as an assumption and continue.
- **Must**: Prefer user clarity over agent cleverness.

## Out of scope

This skill does NOT handle:

- Planning — use sc-planning
- Visual design direction — use sc-design
- Git operations — use sc-git
- Implementation — use sc-coder or /sc-build

## Output format

```
### Open
- **question** — status (blocking/non-blocking)

### Answered
- **date** — **Q: <question>**
  Recommendation accepted: <outcome>
  Source: <source>
```

## Checklist

Before claiming clarification is resolved:

- [ ] Available context inspected before asking
- [ ] Ambiguity classified as blocking, non-blocking, or researchable
- [ ] One blocking question asked at a time (if needed)
- [ ] Answered questions recorded in `questions.md`
- [ ] Confirmed choices and assumptions recorded in `decisions.md`
- [ ] No blocking question remains open before implementation

## Research auto-trigger

Before asking the user about ecosystem conventions, API usage, or framework-specific practices, invoke `spacecraft research <query>` to check current documentation. Only escalate to the user when research doesn't resolve the ambiguity.

---

## References

- `questions.md` — open and answered questions
- `decisions.md` — confirmed choices and assumptions
