> Consult when: applying Tailwind CSS classes, building responsive layouts, ensuring accessibility in visual design, or choosing between utility and component styles.

# Styling patterns

Operational patterns for Tailwind CSS in a React + Vite stack. Co-located utilities, responsive design, and accessibility-first styling.

## Utility-first convention

Style elements directly in JSX using Tailwind utility classes. Do not create separate CSS files for components. Do not use inline styles except for truly dynamic values (e.g., a color from API data).

```tsx
// PREFERRED - utility classes in className
<button className="rounded-lg bg-blue-600 px-4 py-2 text-white hover:bg-blue-700">
  Submit
</button>

// AVOID - separate CSS files for component styles
// AVOID - inline style objects
<button style={{ borderRadius: 8, backgroundColor: "#2563eb" }}>Submit</button>
```

## Responsive design

Tailwind is mobile-first. Base classes apply to all screens. Prefix with breakpoint to override at larger sizes.

| Breakpoint | Min width | Typical use |
|---|---|---|
| `sm` | 640px | Large phones, small tablets |
| `md` | 768px | Tablets |
| `lg` | 1024px | Small laptops |
| `xl` | 1280px | Desktops |
| `2xl` | 1536px | Large displays |

```tsx
// Mobile-first: single column → two columns on md+
<div className="grid grid-cols-1 md:grid-cols-2 gap-4">
  <Card />
  <Card />
</div>

// Hide on mobile, show on desktop
<nav className="hidden md:flex">...</nav>

// Full width on mobile, constrained on desktop
<main className="w-full max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
```

## Accessibility styling

Tailwind provides accessibility utilities that should be part of every component:

- **Focus rings** - `focus:ring-2 focus:ring-blue-500 focus:outline-none` on interactive elements
- **Forced colors** - `forced-colors:outline` for high-contrast mode support
- **Screen reader only** - `sr-only` for visually hidden text that screen readers announce
- **Reduced motion** - `motion-reduce:transition-none` for users who prefer reduced motion

```tsx
// Accessible interactive element
<button
  className="rounded bg-blue-600 px-4 py-2 text-white 
             hover:bg-blue-700 focus:ring-2 focus:ring-blue-500 
             focus:outline-none motion-reduce:transition-none"
  aria-label="Close dialog"
>
  <span aria-hidden="true">×</span>
  <span className="sr-only">Close</span>
</button>
```

## Extracting shared patterns

Wait until a set of classes repeats 3 times before extracting:

```tsx
// After 3+ usages, extract to a component or apply directive
// Option A: Extract to a component
function PrimaryButton({ children, ...props }: ButtonHTMLAttributes<HTMLButtonElement>) {
  return (
    <button
      className="rounded-lg bg-blue-600 px-4 py-2 text-white hover:bg-blue-700 focus:ring-2 focus:ring-blue-500 focus:outline-none"
      {...props}
    >
      {children}
    </button>
  );
}

// Option B: Use @apply in a single global CSS layer (only for truly repeated patterns)
@layer components {
  .btn-primary {
    @apply rounded-lg bg-blue-600 px-4 py-2 text-white hover:bg-blue-700;
  }
}
```

Prefer Option A (component extraction) over Option B (CSS `@apply`). Components are explicit, composable, and type-safe.

## Dark mode

Use the `dark:` variant for dark mode support. Toggle via class strategy (`darkMode: "class"` in Tailwind config).

```tsx
<div className="bg-white text-gray-900 dark:bg-gray-900 dark:text-gray-100">
  <h1 className="text-2xl font-bold text-gray-900 dark:text-white">Dashboard</h1>
</div>
```

## When NOT to add styling

- Pure layout wrappers - use the parent's grid/flex classes
- Spacing between components - use parent's `gap` or `space-y-*` rather than margin on each child
- One-off overrides - keep them inline; don't extract a pattern for a single use

## Page layout

For page skeletons (single/two/three column, app shell, split, hero, card grid) and responsive collapse rules, see `references/layout.md`. Design-time bake-off selection lives in `sc-ux-design/references/layout-patterns.md`.

## Spacecraft integration

- Styles are verified as part of component evidence - `spacecraft evidence "<name>:component" -- npx vitest run`
- Visual regression is out of scope for automated evidence; note manual visual checks in evidence labels
- Design direction should be recorded in `decisions.md` before implementing styled components
- Tailwind config changes (theme extension, plugins) are infrastructure; stage them alongside component code
