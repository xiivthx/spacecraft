---
name: sc-clarify
description: "Blocking-question protocol used inside /sc-discuss. Activate when resolving ambiguous requirements; prefer /sc-discuss as the session entrypoint."
---

# sc-clarify

Blocking-question protocol for `/sc-discuss`. Ask exactly one blocking question at a time. Record answers. Never proceed with hidden assumptions. Prefer `/sc-discuss` as the human slash entrypoint for ask/clarify/brainstorm/decide sessions.

## When to use

Activate inside `/sc-discuss` (or when ambiguity blocks plan/build) for:

- **Clarify scope, behavior, or direction** - ambiguous requirements
- **Resolve a decision before planning** - blocking questions for planning phase
- **Record decisions from a discussion** - capture choices in `decisions.md`

## Workflow

Use this exact sequence unless the user specifies otherwise:

1. **Resolve mission** - `spacecraft resolve`. On conflict or ambiguity, use `spacecraft use <selector>`.

2. **Inspect context** - Before asking the user ANY question, exhaust these sources:
   - Mission `spec.md` - is the answer already stated?
   - `questions.md` - was this already asked?
   - `decisions.md` - was a decision already recorded?
   - Repo files - can the answer be found by reading code?
   - If the question is about ecosystem conventions or API usage, use sc-search (WebSearch/WebFetch) first.
   - **If the answer exists in any of these sources, do not ask the user.**

3. **Classify** - Categorize the ambiguity:
   - **blocking** - cannot safely plan, design, or implement without user input
   - **non-blocking** - can proceed with an explicit assumption written to `decisions.md`
   - **researchable** - answer by reading files, code, or sc-search (WebSearch/WebFetch)

4. **Ask one question** - If blocking: ask exactly one question. Present it directly in the chat. Format:
   ```
   **Question:** <one clear sentence>
   **Why it matters:** <one sentence on impact>
   **Recommendation:** <your suggested answer with brief rationale>
   **If accepted:** <what happens next - one sentence>
   ```

5. **Record** - After the user answers:
   - Record the question and answer in `questions.md` under `### Answered`
   - Record confirmed choices in `decisions.md`
   - If the user accepted your recommendation, note: "Recommendation accepted: <outcome>"
   - If non-blocking: write the assumption directly to `decisions.md`, no need to ask

### Edge cases

- **User doesn't respond** - Do not proceed with implementation. Keep the question open. If session ends, hand off with the open question.
- **Answer raises new questions** - Classify the new question. Ask one at a time. Record each answer before the next.
- **Multiple ambiguities found** - Classify all of them. Ask only the most blocking one first. Note the others.
- **Answer contradicts spec** - Update `spec.md` to reflect the decision. The user's answer is authoritative.
- **User defers decision** - Record the deferral in `decisions.md` with: "Deferred: <question>. Proceeding without." Only proceed if the ambiguity is non-blocking.

## AFK mode

During `/sc-run` (mission `in_progress` with clarify-status clear): **do not auto-trigger**. Only ask if the Commander hits a true hard blocker (missing secret, impossible acceptance). Otherwise write assumptions to `decisions.md`. Prefer hand off to `/sc-discuss` over a long Q&A mid-AFK.

## Rules

- **Must**: Exhaust all context sources before asking the user. Never ask a question answerable from files or research.
- **Must**: Classify every ambiguity as blocking, non-blocking, or researchable.
- **Must**: Ask exactly one blocking question at a time.
- **Must**: Every question includes: the question, why it matters, a recommendation, and what happens if accepted.
- **Must not**: Ask multiple questions in one message.
- **Must not**: Implement, plan, or finalize visual draft while a blocking question is open.
- **Must not**: Auto-trigger mid-AFK (`/sc-run`) unless a hard blocker.
- **Must**: Record answered questions in `questions.md`. Record decisions in `decisions.md`.
- **Must**: Prefer user clarity over agent cleverness. If the user's answer seems suboptimal, state your concern once and accept their decision.

## Out of scope

- Session entry / brainstorm / visual draft ownership - use `/sc-discuss`
- Planning - use sc-planning (via `/sc-run` for roadmap work)
- Visual draft HTML - `/sc-discuss` + sc-ux-design; critique via Task(`sc-designer`)
- Git operations - use sc-git
- Implementation - use `/sc-run` (AFK) or agents; ship with `/sc-ship` only

After `/sc-discuss` clears clarify for roadmap work, recommend a **new session** `/sc-run <roadmap-id>`.

## Output format

```
### Open
- **Q: <question>** - blocking | awaiting response

### Answered
- **2026-07-09** - **Q: Should we use Fastify or Express?**
  Recommendation accepted: Fastify (better TS support, faster startup)
  Source: user accepted recommendation
```

## Checklist

- [ ] Mission resolved, context inspected
- [ ] Research auto-trigger checked via sc-search (if applicable)
- [ ] Ambiguity classified (blocking / non-blocking / researchable)
- [ ] One blocking question asked at a time (if needed)
- [ ] Question includes: question + why + recommendation + what-if-accepted
- [ ] Answer recorded in `questions.md`
- [ ] Decision recorded in `decisions.md`
- [ ] No blocking question remains open before planning or implementation
- [ ] After clarify clear on roadmap work: recommend new session `/sc-run`

## References

- `questions.md` - open and answered questions per mission
- `decisions.md` - confirmed choices, assumptions, and deferrals
- sc-search - WebSearch/WebFetch escalation for researchable questions
- `/sc-discuss` - pre-build HIL session (owns clarify + visual draft)
- `/sc-run` - AFK roadmap runner after clarify is clear
- `/sc-ship` - explicit human-only ship
