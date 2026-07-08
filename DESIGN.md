# Design System: Orbital Console

## North Star

Orbital Console is the default design language for local-first web interfaces. It should feel precise, calm, technical, slightly cinematic, and intentionally sparse. The interface should be usable before it is impressive.

The UI must not become a generic SaaS landing page, decorative toy, crypto dashboard, or AI-generated template. It should earn attention through hierarchy, spacing, typography, rules, and useful information architecture.

## Principles

- Usefulness before spectacle.
- Evidence over ornament.
- Sparse but not empty.
- Structured negative space.
- One accent earns attention.
- No generic AI SaaS visual tropes.

## Visual defaults

```css
:root {
  --sc-bg: #0e1116;
  --sc-bg-deep: #080a0d;
  --sc-surface: #151a21;
  --sc-surface-raised: #1b222b;
  --sc-rule: #2a3340;
  --sc-rule-strong: #3a4656;
  --sc-text: #f3ead7;
  --sc-text-muted: #b9ad98;
  --sc-text-faint: #776f63;
  --sc-accent: #f6b44b;
  --sc-accent-strong: #ffcb6b;
  --sc-cyan: #62d6cf;
  --sc-danger: #ff6b5f;
  --sc-success: #7bd88f;
  --sc-radius-sm: 4px;
  --sc-radius-md: 8px;
  --sc-radius-lg: 14px;
  --sc-radius-pill: 999px;
}
```

## Typography

- Use a strong sans or humanist sans for body text.
- Do not automatically default to Inter or system fonts unless the mission asks for a neutral enterprise UI.
- Use monospace only for metadata, commands, IDs, timestamps, and code-like details.
- Use display typography sparingly for page-level statements.

## Layout and components

- Prefer left-aligned operator layouts over centered marketing shells.
- Use thin rules and table-like alignment for structured information.
- Avoid nested cards, fake metrics, meaningless badges, and broad decorative gradients.
- Use card-like containers only for concrete tools, repeated items, or detail panels.
- Empty states should explain what is absent and offer one useful next action.
- Forms need visible labels, helper text, validation, and strong focus states.
- Buttons use stable dimensions, low radius, and visible hover/focus/active states.

## Motion

- Motion must clarify state, orientation, or feedback.
- Respect reduced motion.
- Do not add decorative animation without purpose.

## Design artifacts

- When original design exploration is weak or generic, scout references before inventing more options.
- Split references by job: layout/template, mood/art, and interaction/motion.
- Use references to calibrate taste and structure, not to clone.
- Treat HTML design artifacts as decision aids, not essays.
- Use normal text questions when the decision can be made clearly in words.
- Create HTML only when side-by-side visual comparison materially helps the user choose.
- Use Feynman-style explanation: plain language, one familiar analogy when useful, labeled visuals, and clear gain/tradeoff.
- Ask one design config question per artifact whenever possible.
- Keep visible copy short: one main sentence, lists of 3 bullets or fewer, and no long theory paragraphs.
- Put design theory behind simple Thai explanations. The user should not need to know terms like hierarchy, IA, affordance, or motion grammar to choose.
- Make visuals teach the difference. Label the important parts and keep unchanged parts visually quiet.
- Remove any caption or adjective that could fit every option.

## Accessibility

- Maintain strong contrast for body text, muted text, controls, and state indicators.
- Provide visible keyboard focus for links, inputs, buttons, tabs, and command controls.
- Use semantic headings in order.
- Do not rely on color alone to convey status.
- Keep line length readable, especially in inspection panels and long-form content.

## Implementation notes

- Prefer CSS custom properties.
- Prefer plain CSS or the project's existing styling approach.
- Do not introduce a styling framework unless explicitly requested.
- Keep design tokens close to the app entry CSS.
- UI code must remain boring and maintainable even if the design feels distinctive.
- Keep backend service logic separate from visual components.
