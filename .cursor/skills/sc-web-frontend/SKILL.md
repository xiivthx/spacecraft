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
   - Open the approved draft; confirm the surface-relevant scenario matrix (applicable `data-state` panels per draft + `spec.md`) is present
   - If missing, stop and recommend `/sc-discuss` - do not invent draft mid-build
   - Skip only for non-visual FE (pure logic/hooks, no UI surface); record skip in `decisions.md`

4. **Port from draft (visual SoT)** - Before coding chrome, open the approved draft HTML and port **`[data-draft-surface]` only** (ignore `[data-draft-chrome]` / frame bezel / viewport toolbar). Sync `DESIGN.md` tokens from the surface. Implement by **porting** structure, tokens, spacing, type, and component chrome - do not invent a second look that only "matches the brief." Map each draft `data-state` to real app states and tests. Behavior/Verify stay owned by `spec.md`; look/behavior conflict → stop and `/sc-discuss`.

5. **Build by slice (component-first)** - Implement one vertical **feature** slice at a time (RED-GREEN under `/sc-run`), but inside each slice build **primitives before the page**:
   1. **Inventory chrome** in `[data-draft-surface]` for this slice (buttons, fields, banners, empty states, tables, …).
   2. **Reuse or upgrade** matching primitives under the project's `components/ui/` (or equivalent). If missing, **add** the primitive first - typed props, styles ported from the draft, Vitest + RTL behavior test.
   3. **Compose** the feature/layout/page from those primitives (and existing feature components). Wire into routing or parent last.
   Prefer small, focused components. Keep feature-only abstractions co-located under `components/<feature>/`. See `references/components.md` for `ui/` vs extract-after-3 and optional Storybook.

6. **Verify (functional)** - `spacecraft evidence "<label>" -- npx vitest run` (or project functional suite). Tests must pass before claiming done. Run the full test suite after each slice to catch regressions.

7. **Verify (visual / draft-parity)** - After visual UI work: sc-ux-design Step 0 draft-parity (side-by-side vs approved draft **surface** - tokens, layout, chrome, states) then Tier 3 via `playwright-cli` (preferred) or Cursor IDE browser (fallback); optional `visual-verify.mjs`. Layout-only match with different chrome, or missing draft states, is blocking. Record screenshot paths in evidence / `decisions.md`. Fix before ready. Do not use system Chrome headless or browser-use/CDP.

8. **Iterate** - Each acceptance check in `plan.json` drives one slice. Keep the UI working at every checkpoint.

### Edge cases

- **User specifies a different stack** - Adapt patterns. Still require component tests and build verification.
- **Project has no frontend yet** - Scaffold with `npm create vite@latest`, install Tailwind CSS, set up Vitest.
- **Design direction missing / draft not approved** - Stop. Recommend `/sc-discuss` + sc-ux-design draft HIL; do not implement UI.
- **Draft present but freestyle temptation** - Port chrome from the approved draft. Do not rebuild a "cleaner" Tailwind look that only matches layout.
- **Accessibility concern** - Check against Tailwind's accessibility utilities and React Testing Library's accessibility queries.
- **Build fails** - Fix before proceeding. Never skip build verification.

## Rules

- **Must**: Resolve mission with `spacecraft resolve` before mutating work. On conflict/ambiguity use `spacecraft use <selector>`.
- **Must**: Default to React + TypeScript + Vite + Tailwind CSS + Vitest when no stack is specified.
- **Must**: For visual UI work, require approved draft HTML from `/sc-discuss` (sc-ux-design) before writing product UI code.
- **Must**: Treat `[data-draft-surface]` in the approved draft as the **visual source of truth** - port structure, tokens, spacing, type, and component chrome; sync `DESIGN.md` from the surface; never port scaffold chrome.
- **Must not**: Freestyle alternate buttons/inputs/tables/empty/error chrome that only vaguely match the draft layout.
- **Must**: Map each draft `data-state` (surface-relevant set from the approved draft + `spec.md`) to product UI and/or tests.
- **Must**: After visual UI implementation, pass Step 0 draft-parity then capture visual verification (`playwright-cli` or Cursor IDE browser; optional `visual-verify.mjs`) and functional test evidence before claiming done.
- **Must not**: Use system Chrome headless or browser-use/CDP for visual verification.
- **Must**: Verify with `spacecraft evidence` after each implementation slice.
- **Must**: Prefer small vertical **feature** slices over broad horizontal scaffolding; within each slice, build or upgrade `components/ui` primitives before composing the page.
- **Must**: Prefer reusing the project's per-app `components/ui/` catalog for draft-named chrome; do not install third-party UI kits (daisyUI, MUI, …) or create a cross-repo design-system package without user approval.
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
Port: tokens/layout/chrome from draft surface; states mapped: <applicable data-state list>
UI primitives: reused|upgraded|added under components/ui: <list>
Page: composed from primitives + feature components
  Props: <typed interface>
  Styles: Tailwind classes ported from draft (co-located in JSX)
  Test: vitest + @testing-library/react
Verify:
  npx vitest run → PASS
  npm run build → PASS
  draft-parity + visual-verify / screenshots → PASS
Evidence: <label>
```

## Checklist

Before claiming frontend work done:

- [ ] Mission resolved, branch created
- [ ] Stack confirmed: React + TypeScript + Vite + Tailwind CSS + Vitest (or approved alternative)
- [ ] Draft HTML approved (or non-visual skip recorded in `decisions.md`)
- [ ] Look ported from approved draft surface (not freestyled from brief alone); `DESIGN.md` synced
- [ ] Draft scenario states mapped to product UI and/or tests
- [ ] Slice built component-first: `components/ui` primitives reused/upgraded/added before page compose
- [ ] Components typed with TypeScript interfaces
- [ ] Styles applied via Tailwind utility classes, scoped to component
- [ ] If the project has Storybook, new `ui/*` primitives have stories (catalog aid - not a ship gate)
- [ ] Component / functional tests pass (`npx vitest run` or project suite)
- [ ] Draft-parity + visual recheck: `playwright-cli` or Cursor IDE browser screenshots captured
- [ ] Build passes (`npm run build`)
- [ ] Evidence captured with `spacecraft evidence`
- [ ] No unapproved dependencies
- [ ] Accessibility: semantic HTML, focus management, aria labels where needed

## References

- `references/components.md` - React component patterns, TypeScript typing, data flow, composition
- `references/testing.md` - Vitest + React Testing Library patterns, query strategies, mocking
- `references/styling.md` - Tailwind CSS conventions, responsive design, accessibility utilities
