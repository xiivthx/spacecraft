---
description: Shape UI direction for the resolved Spacecraft mission
agent: sc-commander
---
Use sc-mission, sc-clarify, and sc-design.
Run:
scripts/spacecraft resolve --json
If resolver safety is not `safe` or no mission is selected, stop before writing design artifacts. Show the conflict/candidates and tell the user to run `scripts/spacecraft missions` then `scripts/spacecraft use <number|id|title>`, or set `SPACECRAFT_MISSION=<mission-id>` for one command.
Treat `.space/current` as fallback state, not sole authority.
Read DESIGN.md and the resolved mission's spec.md, questions.md, decisions.md, and plan.json if present.
If there is no active mission, explain that /sc-start should be used first.
If DESIGN.md is missing, create it using the Orbital Console defaults.

Goal: make the main design image clear before planning or implementation.

First inspect repo context and mission artifacts. Do not ask questions that can be answered from files.

If design intent has blocking ambiguity, ask exactly one simple question and stop. Include:
- why it matters
- your recommended answer
- what happens if the recommendation is accepted

If the current design direction is weak, generic, hard to imagine, or the user asks for references, run a reference scout before deeper design configuration.

Reference scout workflow:
- browse current public sources if needed
- collect 5 to 9 candidate references
- group them into 2 or 3 reference directions
- separate reference purpose:
  - layout/template references, such as Bootstrap examples, template galleries, or real app/site layouts
  - mood/art references, such as Pinterest boards, editorial moodboards, photography, art direction, or place/product imagery
  - interaction references, such as real products, component examples, or motion examples
- for each reference, record source name, URL, useful part, what to borrow as a pattern, what not to copy, and why it fits or does not fit the mission
- ask one text-first question: which reference direction should guide the next config?
- include your recommendation
- record selected references and rejected references in decisions.md

Use references as calibration, not source material to clone. Do not copy third-party screenshots, illustrations, layouts, copy, brand identity, or exact compositions into product UI.

If reference comparison is mostly links and rationale, keep it in chat or create `.space/missions/<id>/design/references.md`.
Create a lightweight reference-board HTML only when side-by-side visual comparison is necessary. Do not hotlink copyrighted images; link to sources or draw simple original diagrams that explain the pattern.

Use a step-by-step design configuration workflow. Do not ask the user to choose one complete design package when the choice can be split into independent config decisions.

Default config order:
1. product metaphor and mood
2. primary user journey
3. first screen layout and information architecture
4. navigation model
5. interaction model
6. color palette
7. typography direction
8. art treatment and imagery
9. 3D depth, if useful
10. transition and animation behavior
11. density, spacing, and pacing
12. key states and accessibility constraints

Ask only one config question at a time. For each question:
- show 2 to 4 focused options for that config only
- include your recommendation
- explain what changes if the recommendation is accepted
- allow mixing, such as layout A with palette B
- record the chosen config in decisions.md before moving to the next config

Default to a normal text question. Do not create HTML for every config.

Before creating HTML, run the visual-economy gate:
- would one short question resolve this faster?
- is the difference spatial, visual, motion-based, or hard to imagine from text?
- will the artifact make the decision clearer, not just prettier?
- does the user need side-by-side visible comparison?
- if mostly no, ask in chat and record the decision instead.

Use text-first questions for configs that can be chosen clearly in words, such as:
- product metaphor and mood, if the options can be named simply
- primary user journey
- navigation model
- interaction rules
- density and pacing intent
- states, accessibility constraints, risks, and priorities

Use HTML only when the visual comparison matters, such as:
- first screen layout and information architecture
- color palette
- typography feel
- art treatment and imagery
- 3D depth
- transition and animation behavior

If a config choice is mostly visual and easier to judge by seeing it, you may create a static HTML artifact under:
.space/missions/<id>/design/

Before creating visual options for a config, define the divergence axes for that config. Options must differ materially in that config, not just in wording.

When exploring full design directions, options must differ materially across at least 4 of these dimensions:
- product metaphor
- primary user journey
- information architecture
- first screen layout
- navigation model
- interaction model
- visual language
- art treatment
- motion/transition model
- density and pacing

The HTML artifact should be dependency-free and directly openable in a browser. Prefer focused config comparison screens over all-in-one direction packages. Examples:
- layout-only comparison with neutral palette
- palette-only comparison using the same layout
- motion-only storyboard using the chosen layout and palette
- art/3D comparison using the chosen layout and palette

When creating or referencing an HTML artifact, also provide an easy preview command using the bundled local server:
node .opencode/skills/sc-design/scripts/serve-html.mjs .space/missions/<id>/design/<artifact>.html --open

If the user wants to browse all resolved mission design artifacts, use:
node .opencode/skills/sc-design/scripts/serve-html.mjs --open

Treat each HTML artifact as a decision aid, not a design essay. Apply a Feynman explanation pass before showing it:
- explain the choice in plain Thai first
- use short simple-English labels only when they help the user answer
- explain the concept like teaching a smart friend who does not know design theory
- use a familiar analogy when useful
- show the idea with a labeled visual
- state what the user gains and what they give up
- rewrite jargon into everyday words

Keep artifact copy compact:
- one artifact answers one design config question
- one main sentence explains what the user is choosing
- no more than 3 bullets in any visible list
- no paragraph longer than 2 short lines
- no theory dump in the main reading path
- put extra rationale in a small "เหตุผลสั้น ๆ (Why)" block only when it helps

Visuals must teach the difference:
- label the important parts directly in the visual
- make the visual demonstrate the current config, not a decorative fake screen
- keep fixed parts quiet so changed parts stand out
- avoid abstract visuals unless the labels make the decision clear

Before showing an artifact, run the clarity gate:
- can the user understand the question in 10 seconds?
- can the user tell how options differ without long reading?
- does each visual match its explanation?
- would any caption fit every option? If yes, rewrite it.
- did you remove adjective-only claims that do not point to concrete UI behavior?

For a focused config artifact, include:
- the one config being decided
- 2 to 4 options for that config
- a tiny "what stays fixed" note only if needed
- one recommended option

Each option card should stay compact:
- "แปลว่าอะไร (Meaning)"
- "เหมาะเมื่อ (Best when)"
- "แลกกับอะไร (Tradeoff)"

Only create a full direction comparison when the early concept itself is unclear. If creating full directions, clearly label which parts are mixable.

Use Thai-first, simple-English support in the artifact when the user uses Thai or the mission context is Thai/multilingual:
- write main headings and explanations in Thai
- add short simple-English labels in parentheses where helpful, such as "แผนที่ก่อน (Map-first)"
- keep English terms short and familiar
- explain design terms in plain Thai instead of relying on theory words
- make option names easy to say back in chat
- avoid long English-only paragraphs
- prefer simple section labels such as "เลือกอะไร (Decision)", "ดูตรงนี้ (Look here)", "เหมาะเมื่อ (Best when)", "แลกกับอะไร (Tradeoff)", and "คำแนะนำ (Pick)"

Do not produce options that share the same screen skeleton with swapped colors, labels, or microcopy unless the current question is specifically about palette.
At least one option should be intentionally outside the default Orbital Console comfort zone while still respecting mission constraints.
Include a similarity audit in the artifact:
- what is genuinely different across options
- what is intentionally shared
- why the options are not just theme variations
If the options are visually or structurally too similar, discard them and create a new set before showing the user.

Use clear option labels so the user can answer with an option name or number.
Do not require the user to understand design theory.
Do not lock design direction until required configs are chosen or explicitly deferred.

Invoke sc-designer as a read-only subagent to shape the UI direction.
A user invocation of /sc-design is explicit permission to use the read-only sc-designer subagent; do not ask for separate subagent permission.
Then update the resolved mission spec.md and/or plan.json with a concise UI section covering:
- target screen or component
- user goal
- chosen mood and tone
- theme and visual metaphor
- information hierarchy
- layout structure
- design primitives
- visual constraints
- palette, typography, art, 3D, transition, and animation constraints when relevant
- accessibility checks
- verification method
Record chosen direction, rejected directions, assumptions, and open design risks in decisions.md.
When the required configs are chosen, synthesize them into one design brief instead of treating any earlier option as an inseparable package.
Do not implement product code in this command.
Do not implement UI code.
Do not add dependencies.
Set or keep mission state as specified/planned depending on current progress.
End with the recommended next action and session advice. Prefer /sc-plan next when design direction is clear.
