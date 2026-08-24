---
name: sc-clarify
description: "Blocking-question protocol used inside /sc-discuss. Activate when resolving ambiguous requirements; prefer /sc-discuss as the session entrypoint. Mid-ask escapes via natural language: re-pitch on confusion, research (sc-search then sc-storm), visualize (bake-off or state table)."
---

# sc-clarify

Blocking-question protocol for `/sc-discuss`. Map blocking decisions as a **design tree**, then ask the open **frontier** in rounds (up to **max 3** independent questions per turn). Record answers. Never proceed with hidden assumptions on blocking classes. Prefer `/sc-discuss` as the human slash entrypoint for ask/clarify/brainstorm/decide sessions. Mid-ask unlocks use natural language (**re-pitch** · **research** · **visualize**) - no new slash skills.

## Goal

Resolve blocking ambiguity so `/sc-discuss` can clear: every blocking decision is settled or explicitly deferred; facts and **clear technical winners** are researched (and auto-picked when evidence is strong) by the agent; user decides Verify / architecture / scope and unclear technical trade-offs.

## Output

Updated `questions.md` (Open / Answered) and `decisions.md` (choices, auto-picks with evidence summaries, true soft assumptions, explicit deferrals). Chat rounds use the Chat ask format (short or rich): **Thai content** in field bodies with greppable **English labels**. `questions.md` / `decisions.md` stay English. Done when the blocking frontier is empty (settled, auto-picked, or deferred).

## Good / Bad

- Good: exhaust context before any ask; research then optional measure before asking technical / performance / library ambiguities; **auto-pick** clear technical winners with greppable/citable evidence (large and clear win); classify blocking / non-blocking / researchable; design tree with prereqs → dependents; ask only the frontier (≤3 independent Qs); serial ask when B depends on A; Verify / architecture / in-out scope soft gaps stay on the frontier until settled or explicitly deferred; facts via read/research/sub-agent; record answers and decisions; Chat ask uses Thai content + English labels (short vs rich); on mid-ask confusion, run Re-pitch on confusion or mid-ask escape (research / visualize) without new slash commands
- Bad: questionnaire dumps; asking look-uppable facts; asking dependent Qs in the same round; asking clearly-won technical / performance / library choices; **auto-picking** Verify / architecture fork / in-out scope (even if "obvious"); silently assuming Verify / architecture / scope into `decisions.md` during discuss; implementing or finalizing draft while blocking frontier items remain open; auto-triggering mid-AFK except hard blockers; English-only Chat ask bodies; dual `ไทย:`/`EN:` Chat blocks; throwaway HTML for mid-ask visualize; new slash skills (`wait-what` / `prototype` / `research`)

## Verify

Blocking frontier empty (all blocking Open items settled, auto-picked, or deferred in `questions.md` / `decisions.md`). No silent assumption recorded for Verify, architecture fork, or in/out scope during `/sc-discuss`. Greppable artifacts only - no new CLI validator.

## When to use

Activate inside `/sc-discuss` (or when ambiguity blocks plan/build) for:

- **Clarify scope, behavior, or direction** - ambiguous requirements
- **Resolve decisions before planning** - blocking questions for planning phase
- **Record decisions from a discussion** - capture choices in `decisions.md`
- **Mid-ask unblock** - human stuck mid-round (confused wording, needs facts, or cannot picture state) → Re-pitch on confusion or mid-ask escape

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
   - For **technical / performance / library** ambiguity: exhaust research (files, sc-search/docs) and, when uncertainty is about performance or measurable trade-offs, run a cheap local check/bench or use existing test evidence when practical.
   - Dispatch research / read / sub-agent for facts. **If the answer exists in any of these sources, do not ask the user.**

3. **Classify** - Categorize each ambiguity:
   - **blocking** - cannot safely plan, design, or implement without user input (includes gaps that affect Goal / Output / Verify, an architecture fork, or in/out scope unless the user explicitly defers). Technical / performance / library items stay here only until settled by user ask, **auto-pick**, or explicit deferral.
   - **non-blocking** - true soft gap (not Verify / architecture / scope); can proceed with an explicit assumption written to `decisions.md`. After a successful **auto-pick**, the technical item is settled (record and treat as closed - not an open frontier ask).
   - **researchable** - answer by reading files, code, or sc-search (WebSearch/WebFetch); never wait on the user for look-uppable facts

4. **Design tree** - Map blocking decisions as a tree (prerequisites → dependents). Document open nodes under `questions.md` Open (brief `depends-on:` notes allowed). The **frontier** is the set of blocking decisions whose prerequisites are settled.

5. **Auto-pick (technical)** - Before asking on a **technical / performance / library** frontier item, run research → optional measure → then auto-pick **or** ask:

   **Evidence sources:** repo facts + sc-search/docs + simple local bench or existing tests when useful.

   **Evidence bar (auto-pick only when the win is large and clear):** e.g. clear order-of-magnitude advantage, or unanimous house convention / already locked in `decisions.md` / `spec.md`. If unsure or the gap is modest → research/test first; if still unclear → ask with Chat ask format (rich when choices).

   - **Auto-pick:** One option clearly wins by a large margin with greppable/citable evidence → write the choice + evidence summary to `decisions.md` (and Answered in `questions.md` if it was Open). Do **not** ask the user. Cite why (evidence bar). Recompute the frontier.
   - **Ask:** Evidence is weak, conflicting, or the advantage is modest → proceed to step 6. Recommendation may still point at the better option.
   - **Never auto-pick:** Verify, architecture fork, or in/out scope - even if "obvious." Those stay on the blocking frontier until settled or explicitly deferred.

   **Micro-example:** Auto-pick - "Repo and existing tests already standardize on Vitest; house convention unanimous → record Vitest in `decisions.md` with evidence; do not ask." Ask - "Bench shows ~1.2x and docs conflict → ask with Chat ask format."

6. **Ask a frontier round** - If the frontier is non-empty after auto-pick: ask every independent frontier question in one round, capped at **max 3**. A question that still depends on another open question belongs to a later round - ask the prerequisite only. Present directly in chat using the **Chat ask format** below (short or rich). Number `Q1`…`Qn` for the round (use `Q1` even for a single-question round).

### Chat ask format

**Language:** Chat ask field bodies are **Thai content**; bold **English labels** stay greppable (`Plain explain`, `Context`, `Question`, `Trade-offs`, `Why it matters`, `Recommendation`, `If accepted`). Paths, commands, product names, and tech terms stay as-is. No dual `ไทย:`/`EN:` Chat blocks. `questions.md` / `decisions.md` stay English. Skill instructions on disk stay US English.

**Short vs rich:**
- **Rich** - question has explicit choices (A/B/…) **or** is a heavy blocking class (Verify / architecture fork / in-out scope). Requires Plain explain, Context, and per-choice Trade-offs.
- **Short** - simple yes/no or single-path questions. Skip Plain explain / Context / Trade-offs.

**Short template** (English labels; Thai bodies):
```
**Q1 - <short title>:** <คำถาม>

**Why it matters:** <หนึ่งประโยค>
**Recommendation:** <คำแนะนำ + เหตุผลสั้น>
**If accepted:** <ขั้นถัดไป>
```

**Rich template** (English labels; Thai bodies):
```
**Q1 - <short title>**

**Plain explain:** <ปัญหาคืออะไร ทำไมต้องตัดสินใจ เลือกแล้วเปลี่ยนอะไร - ภาษาง่าย สั้น>
**Context:** <รู้อะไรแล้วจาก spec/decisions/repo; อะไรยังคลุมเครือ; stake คืออะไร>
**Question:** <คำถาม>
- A) ...
- B) ...

**Trade-offs:**
- **A)** Pros: … | Cons: …
- **B)** Pros: … | Cons: …

**Why it matters:** <หนึ่งประโยค>
**Recommendation:** <คำแนะนำ + เหตุผลสั้น>
**If accepted:** <ขั้นถัดไป>
```

Keep trade-offs tight: 1-2 bullets each side per choice. Field order and bold labels above are greppable - do not rename.

**Rich micro-example** (copy shape, not content):
```
**Q1 - Choose API framework**

**Plain explain:** ต้องเลือก HTTP library สำหรับเซิร์ฟเวอร์ การเลือกมีผลต่อความเร็วตอนสตาร์ท การซัพพอร์ต TypeScript และความคุ้นของทีมตอนดีบัก
**Context:** Spec ล็อก Node.js + TypeScript แล้ว แต่ยังไม่ล็อก framework ในรีโปยังไม่มี server package Stake คือ binding เส้นทาง/plugin ก่อนวางแผน
**Question:** ใช้ Fastify หรือ Express?
- A) Fastify
- B) Express

**Trade-offs:**
- **A)** Pros: เร็ว, TS-first | Cons: ทีมอาจคุ้น Express มากกว่า
- **B)** Pros: เอกสาร/ตัวอย่างเยอะ | Cons: ช้ากว่า พิมพ์ยากกว่า

**Why it matters:** เลือกผิดต้อง rewrite routing ตอน implement
**Recommendation:** A) Fastify - เข้ากับสแต็ก TS และเป้าสตาร์ท
**If accepted:** บันทึกใน decisions.md แล้วถาม frontier ข้อถัดไป (หรือ clear ถ้า frontier ว่าง)
```

### Re-pitch on confusion

When the human signals confusion mid-round (cues under mid-ask escape), or 027 Auto-Clarity drops caveman and triggers re-pitch:

1. STE-lite restatement of the **current frontier** questions only - shorter, plainer, same decisions at stake.
2. Use ubiquitous language from mission `spec.md` / `decisions.md` (not repo `CONTEXT.md`).
3. Re-present with Chat ask format (Thai content, English labels). Do not restart the mission or invent new frontier items.

### mid-ask escape

Natural-language unlock mid-grill - **no new slash** skills (`wait-what` / `prototype` / `research`). Match cue → escape:

| Escape | Cue examples | Agent action |
|--------|--------------|--------------|
| Re-pitch | อธิบายใหม่ / งง / wait what | Run **Re-pitch on confusion** |
| Research | หาข้อมูลก่อน / research แล้วถามใหม่ | `sc-search` first → escalate `sc-storm` only for open-domain/strategy → record findings → re-ask frontier with better options. No `/research` slash. |
| Visualize | นึกภาพไม่ออก / โชว์ state | **UI:** point to existing sc-discuss / sc-ux-design bake-off or draft. **Non-UI:** short chat explanation + state/example table. **No throwaway HTML files.** |

After research or visualize, return to the open frontier round (re-ask if options improved).

7. **Record** - After the user answers (or after an auto-pick):
   - Record each question and answer in `questions.md` under `### Answered`
   - Record confirmed choices in `decisions.md`
   - Auto-pick: note evidence summary and source (e.g. "Auto-pick: <choice>. Evidence: <bar citation>")
   - If the user accepted a recommendation, note: "Recommendation accepted: <outcome>"
   - If the user defers: record explicit deferral in `decisions.md` (`Deferred: <question>. Proceeding without.` or equivalent) and close that Open node
   - Non-blocking (true soft): write the assumption to `decisions.md` without asking
   - Recompute the frontier; repeat steps 5-6 until the blocking frontier is empty

### Soft gaps in `/sc-discuss`

- Gaps that affect **Goal / Output / Verify**, an **architecture fork**, or **in/out scope** → treat as **blocking frontier** items (or require **explicit deferral** recorded). Do **not** clear those gray areas by silent assumption or **auto-pick** into `decisions.md` during discuss.
- **True soft** gaps (not those classes) may still go to `decisions.md` as assumptions.
- AFK `/sc-run` soft→`decisions.md` behavior is unchanged (see AFK mode).

### Edge cases

- **User doesn't respond** - Do not proceed with implementation. Keep frontier questions open. If session ends, hand off with the open frontier.
- **Answer raises new questions** - Classify, attach to the design tree, ask only the new frontier (≤3 independent).
- **Dependent chain** - If B depends on A and A is open, ask A only this round.
- **Answer contradicts spec** - Update `spec.md` to reflect the decision. The user's answer is authoritative.
- **User defers decision** - Record the deferral in `decisions.md`. Only continue past a deferred blocker when the deferral is explicit; do not invent settlement.
- **Mid-ask stuck** - Treat confusion / need-facts / cannot-picture cues as mid-ask escape (re-pitch · research · visualize); do not abandon the frontier round.

## AFK mode

During `/sc-run` (mission `in_progress` with clarify-status clear): **do not auto-trigger**. Only ask if the Commander hits a true hard blocker (missing secret, impossible acceptance). Otherwise write assumptions to `decisions.md`. Prefer hand off to `/sc-discuss` over a long Q&A mid-AFK. Do not expand the AFK ask surface beyond hard blockers.

## Rules

- **Must**: Exhaust all context sources before asking the user. Never ask a question answerable from files or research.
- **Must**: Before asking on a technical / performance / library ambiguity: exhaust research (files, sc-search/docs) and, when uncertainty is measurable, run a cheap local check/bench or use existing test evidence when practical.
- **Must**: Auto-pick a technical / performance / library choice only when the win is **large and clear** (order-of-magnitude advantage, or unanimous house convention / already locked in decisions/spec); record choice + evidence summary in `decisions.md` (and Answered in `questions.md` if Open). Do not ask.
- **Must**: Classify every ambiguity as blocking, non-blocking, or researchable.
- **Must**: Maintain a design tree; ask only the open frontier, max 3 independent questions per turn; serial when dependent.
- **Must**: During `/sc-discuss`, keep Verify / architecture / scope soft gaps on the frontier until settled or explicitly deferred - do not silently assume them.
- **Must**: Every asked question uses the Chat ask format (short or rich). Short: Question + Why it matters + Recommendation + If accepted. Rich (choices or Verify / architecture / in-out scope): also Plain explain + Context + Trade-offs (Pros/Cons per choice). Field bodies are Thai content; bold labels are English.
- **Must**: On mid-ask confusion or escape cues, run Re-pitch on confusion or mid-ask escape (research via sc-search then sc-storm when needed; visualize via bake-off/draft or state/example table) - no new slash skills.
- **Must not**: Auto-pick Verify, architecture fork, or in/out scope - even if "obvious."
- **Must not**: Dump mutually dependent questions in one round, or exceed 3 questions per round.
- **Must not**: Implement, plan, or finalize visual draft while a blocking frontier item is open (unless explicitly deferred).
- **Must not**: Auto-trigger mid-AFK (`/sc-run`) unless a hard blocker.
- **Must not**: Dual `ไทย:`/`EN:` Chat blocks; English-only Chat ask bodies; throwaway HTML files for mid-ask visualize.
- **Must**: Record answered questions in `questions.md`. Record decisions, auto-picks, and deferrals in `decisions.md` (English).
- **Must**: Prefer user clarity over agent cleverness. If the user's answer seems suboptimal, state your concern once and accept their decision.

## Out of scope

- Session entry / brainstorm / visual draft ownership - use `/sc-discuss`
- Planning - use sc-planning (via `/sc-run` for roadmap work)
- Visual draft HTML - `/sc-discuss` + sc-ux-design; required critique via Task(`sc-designer`) before human HIL
- Throwaway HTML for mid-ask visualize - use bake-off/draft pointer or chat state/example table instead
- New slash skills for mid-ask (`wait-what` / `prototype` / `research`) - embed escapes here
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

- **2026-08-24** - **Q: Vitest or Jest?**
  Auto-pick: Vitest. Evidence: unanimous house convention in repo + existing tests
  Source: auto-pick (large and clear)
```

## Checklist

- [ ] Mission resolved, context inspected
- [ ] Research auto-trigger checked via sc-search (if applicable)
- [ ] Technical / performance / library: research (+ optional cheap measure) before ask
- [ ] Ambiguity classified (blocking / non-blocking / researchable)
- [ ] Design tree updated; frontier identified
- [ ] Auto-pick applied only for large-and-clear technical wins; never for Verify / architecture / scope
- [ ] Frontier round asked (≤3 independent; serial if dependent) for remaining open items
- [ ] Each question uses Chat ask format (short or rich; Thai content + English labels; rich adds Plain explain + Context + Trade-offs)
- [ ] Mid-ask confusion handled via Re-pitch on confusion or mid-ask escape when cued
- [ ] Answer or auto-pick recorded in `questions.md`
- [ ] Decision, auto-pick evidence, or explicit deferral recorded in `decisions.md` (English)
- [ ] Verify / architecture / scope gaps not silently assumed or auto-picked during discuss
- [ ] Blocking frontier empty before planning or implementation
- [ ] After clarify clear: handoff `Spec clear. New session: /sc-run.`

## References

- `questions.md` - open and answered questions per mission (English)
- `decisions.md` - confirmed choices, auto-picks, assumptions, and deferrals (English)
- sc-search - WebSearch/WebFetch escalation for researchable questions and mid-ask research escape
- sc-storm - open-domain/strategy research escalate after sc-search (mid-ask research escape)
- `/sc-discuss` - pre-build HIL session (owns clarify + visual draft / bake-off)
- `/sc-run` - AFK roadmap runner after clarify is clear
- `/sc-ship` - explicit human-only ship
