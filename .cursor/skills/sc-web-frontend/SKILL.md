---
name: sc-web-frontend
description: "Build and maintain browser UI with React, TypeScript, Vite, Tailwind CSS, and Vitest. Activate on \"build a React component\", \"create a frontend page\", \"style the UI\", \"add responsive layout\", or \"write frontend tests\"."
---

# sc-web-frontend

Build and maintain browser-side user interfaces under mission control. Default stack: React + TypeScript + Vite + Tailwind CSS + Vitest.

## When to use

Activate when the user asks to:

- **"Build a React component" / "create a frontend page"** - new UI features
- **"Style the UI" / "add responsive layout"** - visual and layout changes
- **"Write frontend tests"** - component and integration testing
- When a mission task requires browser-side implementation

## Workflow

Use this exact sequence unless the user specifies otherwise:

1. **Resolve mission** - `spacecraft resolve`. On conflict or ambiguity, use `spacecraft use <selector>`.

2. **Choose stack** - Default: React + TypeScript + Vite + Tailwind CSS + Vitest. If the project already has a frontend stack, match it. Use sc-search (WebSearch/WebFetch) for `"react latest hooks api"` before using unfamiliar APIs.

3. **UI draft gate (hard)** - For any visual layout/style/component work, require `/sc-discuss` approval first (sc-ux-design brief + draft HTML):
   - `decisions.md` must contain `UI draft approved: <path>` (or recorded non-visual skip)
   - If missing, stop and recommend `/sc-discuss` - do not invent draft mid-build
   - Skip only for non-visual FE (pure logic/hooks, no UI surface); record skip in `decisions.md`

4. **Build by slice** - Implement one vertical feature slice at a time (RED-GREEN under `/sc-run`):
   - Component with its styles (co-located or Tailwind classes)
   - TypeScript types for props, state, and data
   - Test for the component's behavior via Vitest + React Testing Library
   - Wire into the app's routing or parent component
   Prefer small, focused components. Extract shared patterns to `references/components.md` patterns.

5. **Verify (functional)** - `spacecraft evidence "<label>" -- npx vitest run` (or project functional suite). Tests must pass before claiming done. Run the full test suite after each slice to catch regressions.

6. **Verify (visual)** - After visual UI work: sc-ux-design Tier 3 via `playwright-cli` (preferred) or Cursor IDE browser (fallback); optional `visual-verify.mjs`. Record screenshot paths in evidence / `decisions.md`. Fix visual issues before ready. Do not use system Chrome headless or browser-use/CDP.

7. **Iterate** - Each acceptance check in `plan.json` drives one slice. Keep the UI working at every checkpoint.

### Edge cases

- **User specifies a different stack** - Adapt patterns. Still require component tests and build verification.
- **Project has no frontend yet** - Scaffold with `npm create vite@latest`, install Tailwind CSS, set up Vitest.
- **Design direction missing / draft not approved** - Stop. Recommend `/sc-discuss` + sc-ux-design draft HIL; do not implement UI.
- **Accessibility concern** - Check against Tailwind's accessibility utilities and React Testing Library's accessibility queries.
- **Build fails** - Fix before proceeding. Never skip build verification.

## Rules

- **Must**: Resolve mission with `spacecraft resolve` before mutating work. On conflict/ambiguity use `spacecraft use <selector>`.
- **Must**: Default to React + TypeScript + Vite + Tailwind CSS + Vitest when no stack is specified.
- **Must**: For visual UI work, require approved draft HTML from `/sc-discuss` (sc-ux-design) before writing product UI code.
- **Must**: After visual UI implementation, capture visual verification (`playwright-cli` or Cursor IDE browser; optional `visual-verify.mjs`) and functional test evidence before claiming done.
- **Must not**: Use system Chrome headless or browser-use/CDP for visual verification.
- **Must**: Verify with `spacecraft evidence` after each implementation slice.
- **Must**: Prefer small vertical slices over broad horizontal scaffolding.
- **Must**: Component tests cover behavior via public interfaces - render output and user interactions, not internal state.
- **Must**: TypeScript strict mode. All props, state, and event handlers typed.
- **Must**: Styles scoped to component or Tailwind utility classes. No global CSS churn.
- **Must not**: Implement visual layout/style/components without an approved draft (unless non-visual skip is recorded).
- **Must not**: Add state management (Redux, Zustand) unless the component tree exceeds simple prop drilling.
- **Must not**: Add router unless the feature requires multiple views or URL state.
- **Must not**: Install UI libraries or component frameworks without user approval.

## Reviewer checklist

Use this checklist when reviewing frontend code:

- [ ] **React hooks rule violations**
  - Hooks called conditionally (inside `if` or loops)
  - Hooks called after early return
  - `useState` or `useEffect` inside callbacks or conditions
  - Custom hooks not prefixed with `use`
  - Hooks called from regular functions instead of components/hooks
- [ ] **useEffect overuse**
  - `useEffect` used for derived state that could be computed during render
  - Chained effects triggering cascading re-renders
  - Missing cleanup functions for subscriptions, timers, or event listeners
  - Effects with missing or incorrect dependency arrays
  - `useEffect` used where `useSyncExternalStore` or an event handler would be simpler
- [ ] **Tailwind class conflicts**
  - Conflicting utility classes on the same element (e.g., `p-4 p-6`, `text-red-500 text-blue-500`)
  - Unnecessary `!important` modifiers
  - Missing responsive variants
  - Using inline styles instead of Tailwind equivalents
  - Long unbroken class strings without grouping
- [ ] **Missing key props**
  - Mapped elements without a `key` prop
  - Using `index` as key with dynamic or reorderable lists
  - Non-stable keys (random, `Date.now()`)
  - Duplicate keys across siblings

## Out of scope

- API design, server logic, or backend architecture - separate concern
- Database schema, migrations, or query optimization - separate concern
- System architecture decisions or ADR writing - separate concern
- UI design direction and draft HIL - `/sc-discuss` + sc-ux-design (draft HTML approval) before `/sc-run` code
- TDD discipline - use sc-tdd for test-first workflow

## Output format

```
Stack: React + TypeScript + Vite + Tailwind CSS + Vitest
Draft: approved <path> | skip non-visual: <reason>
Component: <name>
  Props: <typed interface>
  Styles: Tailwind classes (co-located in JSX)
  Test: vitest + @testing-library/react
Verify:
  npx vitest run → PASS
  npm run build → PASS
  visual-verify / screenshots → PASS
Evidence: <label>
```

## Checklist

Before claiming frontend work done:

- [ ] Mission resolved, branch created
- [ ] Stack confirmed: React + TypeScript + Vite + Tailwind CSS + Vitest (or approved alternative)
- [ ] Draft HTML approved (or non-visual skip recorded in `decisions.md`)
- [ ] Components typed with TypeScript interfaces
- [ ] Styles applied via Tailwind utility classes, scoped to component
- [ ] Component / functional tests pass (`npx vitest run` or project suite)
- [ ] Visual recheck: `playwright-cli` or Cursor IDE browser screenshots captured
- [ ] Build passes (`npm run build`)
- [ ] Evidence captured with `spacecraft evidence`
- [ ] No unapproved dependencies
- [ ] Accessibility: semantic HTML, focus management, aria labels where needed

## References

- `references/components.md` - React component patterns, TypeScript typing, data flow, composition
- `references/testing.md` - Vitest + React Testing Library patterns, query strategies, mocking
- `references/styling.md` - Tailwind CSS conventions, responsive design, accessibility utilities
