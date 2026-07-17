---
name: sc-design
description: "Shapes UI direction when a mission needs design exploration, critique, or polish against DESIGN.md."
disable-model-invocation: true
---

Use sc-mission, sc-clarify, sc-design, and sc-web-frontend.
Resolve the mission. Block if unsafe.

## Pre-flight checks

Read DESIGN.md and the resolved mission's spec.md, questions.md, decisions.md, and plan.json if present. If no active mission, say to use /sc-start first. If DESIGN.md is missing, create it using Orbital Console defaults.

## Phase detection

Check mission state and design decisions to determine the phase:

| Condition | Phase | Behavior |
|-----------|-------|----------|
| State is `draft`, or no design decisions exist | **Design** | Shape UI direction, record decisions, update spec/plan |
| State is `planned` or later, AND design decisions exist | **Polish** | Small UI fixes before ship |

If ambiguous (e.g. `planned` but no design decisions), ask the user which phase and stop. Do not assume.

---

## Design phase

Goal: clear the main design image before planning or implementation.

### Workflow

1. If design intent has blocking ambiguity, ask exactly one question and stop - why it matters, your recommendation, what happens if accepted.
2. When ready, invoke /sc-designer as a read-only subagent to shape the UI direction. 
3. Update spec.md and/or plan.json with a concise UI section: target screen/component, user goal, chosen mood/tone, visual metaphor, info hierarchy, layout structure, design primitives, visual constraints, palette/typography/art/3D/transition/animation constraints when relevant, accessibility checks, verification method.
4. Record chosen direction, rejected directions, assumptions, and open design risks in decisions.md.
5. When enough configs are chosen, synthesize into one design brief instead of keeping earlier options as packages.
6. **Review gate** - Invoke /sc-reviewer as a read-only subagent to review the design decisions against DESIGN.md. The reviewer checks: anti-slop, option diversity (not same-y), Feynman clarity, Thai-first where applicable, art direction consistency. If the reviewer flags issues, fix them.
7. Set mission state to draft/planned depending on progress.

### Constraints

- Do not implement product code, UI code, or add dependencies.

End with recommended next action and session advice. Prefer /sc-plan next.

---

## Polish phase

Goal: small, low-risk UI fixes before shipping.

### Pre-flight checks

Read DESIGN.md, the resolved mission's spec.md, plan.json, review.md, review.json, and git diff. If the mission has no UI changes, say so and stop.

### Workflow

1. Invoke /sc-designer as a read-only subagent to identify focused polish items and critique against DESIGN.md. 
2. Implement only small, low-risk polish changes that improve:
   - spacing rhythm
   - typography hierarchy
   - color consistency
   - focus/hover/active states
   - empty/loading/error states
   - accessible labels and semantics
   - removal of generic AI-template patterns
3. After polishing, tell the user to run sc-verification and /sc-review.

### Constraints

- Do not add dependencies.
- Do not redesign the whole app.
- Do not change backend behavior.
- Do not claim the UI is ready without verification and design review.

End with session advice. Prefer continuing this chat for immediate verification, unless the thread is context-heavy.

## Research auto-trigger

When design decisions involve unfamiliar UI patterns, accessibility standards, or CSS framework capabilities, run `spacecraft research "<topic>"` before committing to a design direction.

## Hard stop gates

- Resolver conflict or ambiguity
- DESIGN.md missing and cannot be created from defaults
- Ambiguous phase (planned state but no design decisions)
- Blocking clarification open
- UI art direction not chosen when UI tasks are present

## Error handling

- Do not implement product code, UI code, or add dependencies.
- If /sc-designer subagent returns critical design blockers, stop and report - do not proceed to polish.
- If DESIGN.md was modified externally during the session, re-read it before making design decisions.
- Do not assume product or design direction silently. Ask if ambiguous.
