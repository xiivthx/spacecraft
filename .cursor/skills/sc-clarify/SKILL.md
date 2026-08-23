---
name: sc-clarify
description: "Blocking-question protocol used inside /sc-discuss. Activate when resolving ambiguous requirements; prefer /sc-discuss as the session entrypoint."
---

# sc-clarify

Blocking-question protocol for `/sc-discuss`. Map blocking decisions as a **design tree**, then ask the open **frontier** in rounds (up to **max 3** independent questions per turn). Record answers. Never proceed with hidden assumptions on blocking classes. Prefer `/sc-discuss` as the human slash entrypoint for ask/clarify/brainstorm/decide sessions.

## Goal

Resolve blocking ambiguity so `/sc-discuss` can clear: every blocking decision is settled or explicitly deferred; facts are researched by the agent; user only decides.

## Output

Updated `questions.md` (Open / Answered) and `decisions.md` (choices, true soft assumptions, explicit deferrals). Chat rounds use the Chat ask format (short or rich, English labels) below. Done when the blocking frontier is empty (settled or deferred).

## Good / Bad

- Good: exhaust context before any ask; classify blocking / non-blocking / researchable; design tree with prereqs → dependents; ask only the frontier (≤3 independent Qs); serial ask when B depends on A; Verify / architecture / in-out scope soft gaps stay on the frontier until settled or explicitly deferred; facts via read/research/sub-agent; record answers and decisions; chat asks use English Chat ask format (short vs rich)
- Bad: questionnaire dumps; asking look-uppable facts; asking dependent Qs in the same round; silently assuming Verify / architecture / scope into `decisions.md` during discuss; implementing or finalizing draft while blocking frontier items remain open; auto-triggering mid-AFK except hard blockers

## Verify

Blocking frontier empty (all blocking Open items settled or deferred in `questions.md` / `decisions.md`). No silent assumption recorded for Verify, architecture fork, or in/out scope during `/sc-discuss`. Greppable artifacts only - no new CLI validator.

## When to use

Activate inside `/sc-discuss` (or when ambiguity blocks plan/build) for:

- **Clarify scope, behavior, or direction** - ambiguous requirements
- **Resolve decisions before planning** - blocking questions for planning phase
- **Record decisions from a discussion** - capture choices in `decisions.md`

Mission brief presents Goal / Will do / Impact / Extra bullets and does not quiz - grilling lives here.

## Workflow

Use this exact sequence unless the user specifies otherwise:

1. **Resolve mission** - `spacecraft resolve`. On conflict or ambiguity, use `spacecraft use <selector>`.

2. **Inspect context** - Before asking the user ANY question, exhaust these sources:
   - Mission `spec.md` - is the answer already stated?
   - `questions.md` - was this already asked?
   - `decisions.md` - was a decision already recorded?
   - Repo files - can the answer be found by reading code?
   - If the question is about ecosystem conventions or API usage, use sc-search (WebSearch/WebFetch) first.
   - Dispatch research / read / sub-agent for facts. **If the answer exists in any of these sources, do not ask the user.**

3. **Classify** - Categorize each ambiguity:
   - **blocking** - cannot safely plan, design, or implement without user input (includes gaps that affect Goal / Output / Verify, an architecture fork, or in/out scope unless the user explicitly defers)
   - **non-blocking** - true soft gap (not Verify / architecture / scope); can proceed with an explicit assumption written to `decisions.md`
   - **researchable** - answer by reading files, code, or sc-search (WebSearch/WebFetch); never wait on the user for look-uppable facts

4. **Design tree** - Map blocking decisions as a tree (prerequisites → dependents). Document open nodes under `questions.md` Open (brief `depends-on:` notes allowed). The **frontier** is the set of blocking decisions whose prerequisites are settled.

5. **Ask a frontier round** - If the frontier is non-empty: ask every independent frontier question in one round, capped at **max 3**. A question that still depends on another open question belongs to a later round - ask the prerequisite only. Present directly in chat using the **Chat ask format** below (short or rich). Number `Q1`…`Qn` for the round (use `Q1` even for a single-question round).

### Chat ask format

**Language:** The entire ask block is English only. Paths, commands, product names, and tech terms stay as-is. Skill instructions on disk stay US English.

**Short vs rich:**
- **Rich** - question has explicit choices (A/B/…) **or** is a heavy blocking class (Verify / architecture fork / in-out scope). Requires Feynman plain-explain, context, and per-choice pros/cons.
- **Short** - simple yes/no or single-path questions. Skip Feynman / context / trade-offs.

**Short template** (all labels English):
```
**Q1 - <short title>:** <question>

**Why it matters:** <one sentence>
**Recommendation:** <suggested answer + brief rationale>
**If accepted:** <what happens next>
```

**Rich template** (all labels English):
```
**Q1 - <short title>**

**Plain explain:** <Feynman: what the problem is, why we need a decision, what the choice will change - plain language, short paragraphs or bullets; no jargon dump>
**Context:** <what we already know from spec/decisions/repo; what is still ambiguous; what is at stake>
**Question:** <question>
- A) ...
- B) ...

**Trade-offs:**
- **A)** Pros: … | Cons: …
- **B)** Pros: … | Cons: …

**Why it matters:** <one sentence>
**Recommendation:** <suggested answer + brief rationale>
**If accepted:** <what happens next>
```

Keep trade-offs tight: 1-2 bullets each side per choice. Field order and bold labels above are greppable - do not rename.

**Rich micro-example** (copy shape, not content):
```
**Q1 - Choose API framework**

**Plain explain:** We need to pick the HTTP API library for the server. The choice affects startup speed, TypeScript support, and how familiar the team is when debugging.
**Context:** Spec already locks Node.js + TypeScript, but not the framework. Repo has no server package yet. Stake is routing/plugin lock-in before planning.
**Question:** Use Fastify or Express?
- A) Fastify
- B) Express

**Trade-offs:**
- **A)** Pros: fast, TS-first | Cons: team may know Express better
- **B)** Pros: more docs/examples | Cons: slower, messier typing

**Why it matters:** Wrong pick forces a routing rewrite during implement.
**Recommendation:** A) Fastify - fits the TS stack and startup target
**If accepted:** Record in decisions.md, then ask the next frontier item (or clear if frontier is empty)
```

6. **Record** - After the user answers:
   - Record each question and answer in `questions.md` under `### Answered`
   - Record confirmed choices in `decisions.md`
   - If the user accepted a recommendation, note: "Recommendation accepted: <outcome>"
   - If the user defers: record explicit deferral in `decisions.md` (`Deferred: <question>. Proceeding without.` or equivalent) and close that Open node
   - Non-blocking (true soft): write the assumption to `decisions.md` without asking
   - Recompute the frontier; repeat step 5 until the blocking frontier is empty

### Soft gaps in `/sc-discuss`

- Gaps that affect **Goal / Output / Verify**, an **architecture fork**, or **in/out scope** → treat as **blocking frontier** items (or require **explicit deferral** recorded). Do **not** clear those gray areas by silent assumption into `decisions.md` during discuss.
- **True soft** gaps (not those classes) may still go to `decisions.md` as assumptions.
- AFK `/sc-run` soft→`decisions.md` behavior is unchanged (see AFK mode).

### Edge cases

- **User doesn't respond** - Do not proceed with implementation. Keep frontier questions open. If session ends, hand off with the open frontier.
- **Answer raises new questions** - Classify, attach to the design tree, ask only the new frontier (≤3 independent).
- **Dependent chain** - If B depends on A and A is open, ask A only this round.
- **Answer contradicts spec** - Update `spec.md` to reflect the decision. The user's answer is authoritative.
- **User defers decision** - Record the deferral in `decisions.md`. Only continue past a deferred blocker when the deferral is explicit; do not invent settlement.

## AFK mode

During `/sc-run` (mission `in_progress` with clarify-status clear): **do not auto-trigger**. Only ask if the Commander hits a true hard blocker (missing secret, impossible acceptance). Otherwise write assumptions to `decisions.md`. Prefer hand off to `/sc-discuss` over a long Q&A mid-AFK. Do not expand the AFK ask surface beyond hard blockers.

## Rules

- **Must**: Exhaust all context sources before asking the user. Never ask a question answerable from files or research.
- **Must**: Classify every ambiguity as blocking, non-blocking, or researchable.
- **Must**: Maintain a design tree; ask only the open frontier, max 3 independent questions per turn; serial when dependent.
- **Must**: During `/sc-discuss`, keep Verify / architecture / scope soft gaps on the frontier until settled or explicitly deferred - do not silently assume them.
- **Must**: Every asked question uses the Chat ask format (short or rich). Short: Question + Why it matters + Recommendation + If accepted. Rich (choices or Verify / architecture / in-out scope): also Plain explain + Context + Trade-offs (Pros/Cons per choice). Entire ask block in English.
- **Must not**: Dump mutually dependent questions in one round, or exceed 3 questions per round.
- **Must not**: Implement, plan, or finalize visual draft while a blocking frontier item is open (unless explicitly deferred).
- **Must not**: Auto-trigger mid-AFK (`/sc-run`) unless a hard blocker.
- **Must**: Record answered questions in `questions.md`. Record decisions and deferrals in `decisions.md`.
- **Must**: Prefer user clarity over agent cleverness. If the user's answer seems suboptimal, state your concern once and accept their decision.

## Out of scope

- Session entry / brainstorm / visual draft ownership - use `/sc-discuss`
- Planning - use sc-planning (via `/sc-run` for roadmap work)
- Visual draft HTML - `/sc-discuss` + sc-ux-design; required critique via Task(`sc-designer`) before human HIL
- Git operations - use sc-git
- Implementation - use `/sc-run` (AFK) or agents; ship with `/sc-ship` only

After `/sc-discuss` clears clarify:
- `Sizing: roadmap <id>` → handoff **Spec clear. New session: `/sc-run <id>`.**
- `Sizing: single|phases` → handoff **Spec clear. New session: `/sc-run`.** (mission-only)

## Output format

```
### Open
- **Q: <question>** - blocking | awaiting response | depends-on: <other Q or none>

### Answered
- **2026-07-09** - **Q: Should we use Fastify or Express?**
  Recommendation accepted: Fastify (better TS support, faster startup)
  Source: user accepted recommendation
```

## Checklist

- [ ] Mission resolved, context inspected
- [ ] Research auto-trigger checked via sc-search (if applicable)
- [ ] Ambiguity classified (blocking / non-blocking / researchable)
- [ ] Design tree updated; frontier identified
- [ ] Frontier round asked (≤3 independent; serial if dependent)
- [ ] Each question uses Chat ask format (short or rich; English labels; rich adds Plain explain + Context + Trade-offs)
- [ ] Answer recorded in `questions.md`
- [ ] Decision or explicit deferral recorded in `decisions.md`
- [ ] Verify / architecture / scope gaps not silently assumed during discuss
- [ ] Blocking frontier empty before planning or implementation
- [ ] After clarify clear: handoff `Spec clear. New session: /sc-run.`

## References

- `questions.md` - open and answered questions per mission
- `decisions.md` - confirmed choices, assumptions, and deferrals
- sc-search - WebSearch/WebFetch escalation for researchable questions
- `/sc-discuss` - pre-build HIL session (owns clarify + visual draft)
- `/sc-run` - AFK roadmap runner after clarify is clear
- `/sc-ship` - explicit human-only ship
