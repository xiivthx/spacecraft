---
name: sc-design
description: Shape, critique, and polish visual/UI design using DESIGN.md. Activate on /sc-design, UI, layout, styling, or visual design requests.
license: MIT
compatibility: opencode
metadata:
  version: 1.0
  audience: xiivthx
---

# sc-design

Shape, critique, and polish visual/UI design using the project DESIGN.md.

## When to use

Activate when the user asks to:

- Design or modify web UI
- Choose design direction, layout, palette, or art treatment
- Review or critique existing UI
- Scout design references
- Create design HTML artifacts for decision-making

## Workflow

Use this exact sequence unless the user specifies otherwise:

1. **Read DESIGN.md** — Before designing or modifying UI, read `DESIGN.md`. If missing, create a minimal one using Orbital Console.
2. **Clarify intent** — Before shaping UI direction, use sc-clarify if design intent is ambiguous. Ask one design question at a time with a recommended answer.
3. **Configure design** — Follow the default config order: product metaphor and mood, primary user journey, first screen layout, navigation, interaction, color palette, typography, art treatment, 3D depth, transition/animation, density/spacing, key states/accessibility. Ask one question at a time. Show 2–4 focused options per config. Record in `decisions.md`.
4. **Create artifact (if needed)** — Create HTML under `.space/missions/<id>/design/` only when visual comparison materially helps. Use Thai-first, simple-English labels. Apply Feynman clarity gate.
5. **Reference scout (if needed)** — Browse public sources. Separate by purpose: layout/template, mood/art, interaction/motion. Present 5–9 candidates in 2–3 directions.
6. **Polish and review** — Before ship, the polish phase handles small UI fixes and read-only design critique against DESIGN.md.

## Rules

> **Reading guide**: Rules marked **Must** / **Must not** are verifiable constraints. Rules marked **Prefer** are taste-based guidelines — sound defaults, not gates.

### General

- **Must**: Read `DESIGN.md` before designing or modifying UI.
- **Must**: Before shaping UI direction, use sc-clarify if design intent is ambiguous.
- **Must**: Ask one design question at a time.
- **Must**: Include a recommended answer.
- **Must not**: Invent brand, mood, audience, or visual metaphor silently when it materially changes the product.
- **Must**: Record design choices in `decisions.md`.
- **Must**: Make the main design image clear before planning or implementation.
- **Must not**: Add CSS frameworks unless explicitly requested.
- **Must**: Use accessible semantic HTML and visible focus states.

### Design config

- **Must**: Default to asking a normal text question. Do not create an HTML artifact for every design config.
- **Must**: Create HTML only when seeing options materially improves the decision or the user explicitly asks for a visual artifact.
- **Must**: If mood, tone, theme, color, composition, 3D, transition, animation, or art direction is easy to answer in words, ask in chat instead of making HTML.
- **Must**: Use a quick visual-economy gate before creating HTML: would one short question resolve this faster? Is the difference spatial, visual, motion-based, or hard to imagine from text? Will the artifact make the decision clearer, not just prettier? Will the user need to compare visible differences side by side? If mostly no, ask normally.
- **Must**: Text-first configs usually include: product metaphor and mood (if simply named), primary user journey, navigation model, interaction rules, density and pacing intent, key states, accessibility constraints, risks, and priorities.
- **Must**: Visual-helpful configs may include: first screen layout and information architecture, color palette, typography feel, art treatment and imagery, 3D depth, transition and animation behavior.

### Reference scouting

- **Must**: Use reference scouting when the design feels weak, generic, hard to imagine, or the user asks for references.
- **Must**: Reference scouting usually happens before deep design config.
- **Must**: Reference review happens after UI exists, using the chosen references to check quality and drift.
- **Must not**: Create a separate command for reference scouting unless the user explicitly asks. Keep it inside sc-design.
- **Must**: When scouting references, browse current public sources if needed and separate them by purpose: layout/template, mood/art, interaction/motion.
- **Must**: Prefer 5 to 9 candidate references grouped into 2 or 3 directions.
- **Must**: For each reference, record: source name and URL, what it is useful for (layout, mood, palette, typography, art, interaction, or motion), what to borrow as a pattern, what not to copy, why it fits or does not fit.
- **Must**: Treat references as calibration, not source material to clone.
- **Must not**: Copy third-party screenshots, illustrations, layouts, copy, or brand identity into product UI.
- **Must not**: Hotlink copyrighted images into local artifacts. Link to sources instead, or use simple original diagrams.

### HTML artifacts

- **Must**: HTML design artifacts belong under `.space/missions/<id>/design/`.
- **Must**: Use the artifact to show options, not to implement product UI.
- **Must**: To preview HTML artifacts, use: `node .engine/skills/sc-design/scripts/serve-html.mjs [artifact-or-dir] --open`.
- **Must**: If no artifact path is provided, the preview script serves the resolved mission's `design/` folder.
- **Must**: After creating an HTML artifact, include the exact preview command in the response.
- **Must**: Treat HTML artifacts as decision aids, not design essays.
- **Must**: Apply a Feynman explanation pass before showing any HTML artifact: name the choice in plain Thai first (with simple-English label when useful), explain as if teaching a smart friend, use one familiar analogy (map, notebook, counter, timeline, tray, or studio), show the idea with a labeled visual that makes the explanation obvious, say what the user gains and gives up, rewrite jargon until it uses everyday words.
- **Must**: Keep user-facing artifact copy short: one question per artifact, one main sentence, no more than 3 bullets in any visible list, no paragraph longer than 2 short lines, no theory dump in the main reading path. Move extra rationale into a small "เหตุผลสั้น ๆ (Why)" or omit it.
- **Must**: Every option card should answer: what this means in plain language, when to choose it, what changes on screen, what risk or tradeoff matters.
- **Must**: Visuals must teach the difference: label the important parts directly in the visual, make the visual demonstrate the config being chosen, keep fixed elements visually quiet, avoid abstract composition unless labels make the decision clear.
- **Must**: Before showing an artifact, run the clarity gate: can the user understand the question in 10 seconds? Can the user tell how options differ without reading long text? Does each visual match the explanation beside it? Would any caption fit every option (if yes, rewrite)? Did you remove adjectives that do not point to concrete UI behavior?

### Options and divergence

- **Must**: Use a step-by-step design configuration workflow. Do not force the user to choose one complete design package when decisions can be mixed.
- **Must**: Default config order: product metaphor and mood, primary user journey, first screen layout and information architecture, navigation model, interaction model, color palette, typography direction, art treatment and imagery, 3D depth (if useful), transition and animation behavior, density/spacing/pacing, key states and accessibility constraints.
- **Must**: Ask one config question at a time.
- **Must**: For each config question: show 2 to 4 focused options for that config only, include your recommendation, explain what changes if accepted, allow mixing (layout A with palette B), record in decisions.md.
- **Must**: For HTML option cards, prefer compact structure: "แปลว่าอะไร (Meaning)", "เหมาะเมื่อ (Best when)", "แลกกับอะไร (Tradeoff)".
- **Must**: Prefer focused config artifacts over all-in-one direction packages: layout-only comparison with neutral palette, palette-only comparison using same layout, motion-only storyboard using chosen layout and palette, art/3D comparison using chosen layout and palette.
- **Must**: Only create full direction comparisons when the early concept itself is unclear.
- **Must**: When required configs are chosen, synthesize them into one design brief instead of treating any earlier option as an inseparable package.
- **Must**: Before creating options, define divergence axes. Options must differ materially across at least 4 dimensions from: product metaphor, primary user journey, information architecture, first screen layout, navigation model, interaction model, visual language, art treatment, motion/transition model, density and pacing.
- **Must not**: Produce options that share the same screen skeleton with swapped colors, labels, or microcopy unless the current question is specifically about palette.
- **Must**: At least one option should be intentionally outside the default Orbital Console comfort zone while still respecting mission constraints.
- **Must**: Include a similarity audit in any HTML comparison artifact: what is genuinely different, what is intentionally shared, why the options are not just theme variations.
- **Must**: If options are visually or structurally too similar, discard them and create a new set before showing the user.

### Thai-first artifacts

> **Locale-dependent**: These rules apply only when the user works in Thai or the mission context is Thai/multilingual. For other locales, adapt labels and language to the user's working language.

- **Must**: Use Thai-first, simple-English support in HTML artifacts when the user uses Thai or the mission context is Thai/multilingual: main headings and explanations in Thai, short simple-English labels in parentheses where helpful, plain Thai explanations for design theory terms, option names that are easy to say back in chat, no long English-only paragraphs.
- **Must**: In Thai-first artifacts, prefer simple labels such as: "เลือกอะไร (Decision)", "ดูตรงนี้ (Look here)", "เหมาะเมื่อ (Best when)", "แลกกับอะไร (Tradeoff)", "คำแนะนำ (Pick)".

### Full direction exploration

- **Must**: For full direction exploration, present 3 to 5 distinct options covering: mood and tone, product metaphor, palette, typography direction, composition and layout, primary journey and navigation model, UX emphasis, art treatment, 3D possibility (only if useful), transition and animation principles, risks and anti-patterns.
- **Must**: Each option must include: name, intended feeling, best-fit use case, the simplest explanation, what changes in the UI, key color/layout/motion/accessibility notes only when they affect the decision, what would be built first.
- **Must**: Include a recommendation and explain why.
- **Must not**: Lock art direction until the user chooses an option or explicitly accepts the recommendation.

### DESIGN.md

- **Must**: If `DESIGN.md` is missing, create a minimal one using Orbital Console.
- **Must**: Treat `DESIGN.md` as the source of truth for look and feel.
- **Must not**: Copy third-party designs.
- **Must not**: Imitate a named brand exactly.
- **Must**: Use references only to understand patterns and quality bars.

### Anti-patterns

- **Must not**: Use generic AI slop: no generic SaaS hero defaults, no purple-blue gradient default, no card grids as the automatic answer, no nested cards, no meaningless badges, no fake metrics, no stock illustration placeholders, no magic-wand/sparkle AI cliches.

### Pre-implementation shaping

- **Must**: Before implementation, shape the UI: user goal, primary screen, information hierarchy, layout structure, components needed, states needed, acceptance checks.
- **Must**: Prefer small, implementable vertical slices.
- **Must**: For each UI task, define: target screen/component, visual intent, interaction behavior, accessibility requirement, verification method.

### Review

- **Must**: During review, check: hierarchy, clarity, consistency with `DESIGN.md`, accessibility, responsive behavior, unnecessary ornament, generic AI-template smell.
- **Must**: If browser/screenshot tooling does not exist, use the bundled HTML preview server for manual visual verification when useful.
- **Must not**: Add screenshot tooling unless explicitly requested.

## Out of scope

This skill does NOT handle:

- Product implementation — use sc-coder or the build command
- Git operations or release — use sc-git
- Mission planning — use sc-planning
- Evidence capture — use sc-verification
- Clarification routing — use sc-clarify

## Output format

Design decisions are recorded in `decisions.md`. Design artifacts go under `.space/missions/<id>/design/`. The ultimate design reference is `DESIGN.md`.

```
Design artifact layout:
.space/missions/<id>/design/
  <config-or-option-name>.html   # decision aid HTML
  references.md                   # scouted references (optional)
```

## Checklist

Before claiming design work is ready:

- [ ] `DESIGN.md` read (or created if missing)
- [ ] Design intent clarified with user if ambiguous
- [ ] Design choices recorded in `decisions.md`
- [ ] Configs asked one at a time, 2–4 options each
- [ ] Options differ materially (not just palette swaps)
- [ ] No AI slop: no generic hero, purple-blue gradient, card grid default, fake metrics
- [ ] HTML artifacts pass clarity gate
- [ ] Thai-first labels used when user works in Thai
- [ ] Art direction not locked without user choice

## References

- `DESIGN.md` — project design direction
- Preview server: `node ./scripts/serve-html.mjs [artifact-or-dir] --open`
