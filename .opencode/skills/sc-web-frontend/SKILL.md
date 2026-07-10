---
name: sc-web-frontend
description: Build and maintain browser UI with React, TypeScript, Vite, Tailwind CSS, and Vitest. Activate on "build a React component", "create a frontend page", "style the UI", "add responsive layout", or "write frontend tests".
license: MIT
compatibility: opencode
metadata:
  version: 1.0
  audience: xiivthx
---

# sc-web-frontend

Build and maintain browser-side user interfaces under mission control. Default stack: React + TypeScript + Vite + Tailwind CSS + Vitest.

## When to use

Activate when the user asks to:

- **"Build a React component" / "create a frontend page"** — new UI features
- **"Style the UI" / "add responsive layout"** — visual and layout changes
- **"Write frontend tests"** — component and integration testing
- When a mission task requires browser-side implementation

## Workflow

Use this exact sequence unless the user specifies otherwise:

1. **Resolve mission** — `scripts/spacecraft resolve --json`. Block if safety ≠ `safe`.

2. **Choose stack** — Default: React + TypeScript + Vite + Tailwind CSS + Vitest. If the project already has a frontend stack, match it. Run `spacecraft research "react latest hooks api"` before using unfamiliar APIs.

3. **Build by slice** — Implement one vertical feature slice at a time:
   - Component with its styles (co-located or Tailwind classes)
   - TypeScript types for props, state, and data
   - Test for the component's behavior via Vitest + React Testing Library
   - Wire into the app's routing or parent component
   Prefer small, focused components. Extract shared patterns to `references/components.md` patterns.

4. **Verify** — `scripts/spacecraft evidence "<label>" -- npx vitest run`. Tests must pass before claiming done. Run the full test suite after each slice to catch regressions.

5. **Iterate** — Each acceptance check in `plan.json` drives one slice. Keep the UI working at every checkpoint.

### Edge cases

- **User specifies a different stack** — Adapt patterns. Still require component tests and build verification.
- **Project has no frontend yet** — Scaffold with `npm create vite@latest`, install Tailwind CSS, set up Vitest.
- **Design direction missing** — Stop and recommend running the design workflow before implementing UI.
- **Accessibility concern** — Check against Tailwind's accessibility utilities and React Testing Library's accessibility queries.
- **Build fails** — Fix before proceeding. Never skip build verification.

## Rules

- **Must**: Resolve mission with `scripts/spacecraft resolve --json` before mutating work.
- **Must**: Default to React + TypeScript + Vite + Tailwind CSS + Vitest when no stack is specified.
- **Must**: Verify with `scripts/spacecraft evidence` after each implementation slice.
- **Must**: Prefer small vertical slices over broad horizontal scaffolding.
- **Must**: Component tests cover behavior via public interfaces — render output and user interactions, not internal state.
- **Must**: TypeScript strict mode. All props, state, and event handlers typed.
- **Must**: Styles scoped to component or Tailwind utility classes. No global CSS churn.
- **Must not**: Add state management (Redux, Zustand) unless the component tree exceeds simple prop drilling.
- **Must not**: Add router unless the feature requires multiple views or URL state.
- **Must not**: Install UI libraries or component frameworks without user approval.

## Out of scope

- API design, server logic, or backend architecture — separate concern
- Database schema, migrations, or query optimization — separate concern
- System architecture decisions or ADR writing — separate concern
- UI design direction and visual critique — use sc-design before code when needed
- TDD discipline — use sc-tdd for test-first workflow

## Output format

```
Stack: React + TypeScript + Vite + Tailwind CSS + Vitest
Component: <name>
  Props: <typed interface>
  Styles: Tailwind classes (co-located in JSX)
  Test: vitest + @testing-library/react
Verify:
  npx vitest run → PASS
  npm run build → PASS
Evidence: <label>
```

## Checklist

Before claiming frontend work done:

- [ ] Mission resolved, branch created
- [ ] Stack confirmed: React + TypeScript + Vite + Tailwind CSS + Vitest (or approved alternative)
- [ ] Components typed with TypeScript interfaces
- [ ] Styles applied via Tailwind utility classes, scoped to component
- [ ] Component tests pass (`npx vitest run`)
- [ ] Build passes (`npm run build`)
- [ ] Evidence captured with `scripts/spacecraft evidence`
- [ ] No unapproved dependencies
- [ ] Accessibility: semantic HTML, focus management, aria labels where needed

## Research auto-trigger

When using unfamiliar React APIs, Tailwind CSS utilities, or Vitest features, run `spacecraft research "<topic> latest docs"` before implementing. Frontend tooling evolves rapidly — verify current APIs.

---

## References

- `references/components.md` — React component patterns, TypeScript typing, data flow, composition
- `references/testing.md` — Vitest + React Testing Library patterns, query strategies, mocking
- `references/styling.md` — Tailwind CSS conventions, responsive design, accessibility utilities
