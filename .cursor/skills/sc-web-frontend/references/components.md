> Consult when: building new React components, choosing between composition patterns, or deciding data flow strategy.

# Component patterns

Operational patterns for React components in a TypeScript + Vite + Tailwind CSS stack.

## Component-first (Must under `/sc-run`)

For visual slices ported from an approved draft:

1. Inventory reusable chrome in `[data-draft-surface]` (Button, Field, Banner, EmptyState, …).
2. Reuse or upgrade primitives under **`components/ui/`** (per-app catalog). Add a missing primitive before the page.
3. Compose the feature/layout/page from those primitives; wire the route last.

Do **not** start by painting a whole page of one-off markup and extracting later unless the markup is truly single-use (~10 lines). Do **not** install daisyUI / MUI / similar without user approval - grow a small in-project catalog that matches `DESIGN.md` + the draft.

### shadcn MCP (approved primitive source)

Pack `frontend` merges the official shadcn MCP (`npx shadcn@latest mcp`) into `.cursor/mcp.json`. Copy-into-project primitives under `components/ui/` are allowed; full UI kits are not.

When adding missing chrome:

1. Prefer MCP **search/list** over inventing markup from scratch.
2. **Add** the matching registry item into the app (needs `components.json`).
3. **Retoken** immediately to house `DESIGN.md` + approved draft surface (colors, radius, type, spacing). Do not ship default shadcn slate/zinc/neutral theme.
4. Draft / `DESIGN.md` remain visual SoT - registry is structure and a11y behavior, not look.

If MCP is disconnected, hand-port the primitive from the draft; still do not freestyle a second kit.

### `components/ui/` vs feature folders

| Location | Owns | Example |
|----------|------|---------|
| `components/ui/` | Look primitives shared across features in **this** app | `Button`, `TextField`, `Banner`, `EmptyState` |
| `components/<feature>/` | Feature screens, shells, domain widgets | `TripDaysPage`, `SignInForm` |

Keep the catalog **per-app** until a second product app needs the same kit - then promote to a shared package (user decision). Do not invent a cross-repo design-system package by default.

### Extract-after-3 (feature abstractions)

- **Draft-named chrome** (button, field, banner, empty) may land in `components/ui/` on **first** port when the draft shows it as product chrome.
- **Feature-only** patterns (trip wizard steps, portal nav shape) stay in `components/<feature>/` until the same abstraction appears **3 times** - then extract (still prefer feature folder or `ui/` only if truly look-primitive).

### Storybook (optional catalog)

- If the project **already has** Storybook (or the user asks to add a light one): every new or upgraded `components/ui/*` primitive gets a story (variants + key states). Use it as a human catalog, like browsing a component gallery.
- Storybook is a **review aid**, not a ship gate - Vitest + draft-parity + project verify remain authoritative.
- Do not add Storybook to a project unless the user asks.

## Component structure

Prefer function components with typed props. Keep components focused - one responsibility per component.

```tsx
// components/ui/Button.tsx
interface ButtonProps {
  children: React.ReactNode;
  onClick?: () => void;
  variant?: "primary" | "ghost";
  disabled?: boolean;
}

export function Button({
  children,
  onClick,
  variant = "primary",
  disabled = false
}: ButtonProps) {
  return (
    <button
      type="button"
      className={variant === "primary" ? "btn-primary" : "btn-ghost"}
      onClick={onClick}
      disabled={disabled}
    >
      {children}
    </button>
  );
}
```

Key conventions:
- Props interface named `<Component>Props`, defined above the component
- Export the component as a named export, not default
- Co-locate styles with markup using Tailwind utility classes (tokens from `DESIGN.md` / draft)
- Interactive elements must be keyboard accessible

## Data flow

**Props down, callbacks up.** Pass data through props. Communicate upward through callback props.

**Avoid premature state management.** Start with local `useState`. Lift state up when siblings need it. Introduce context when props pass through >2 intermediate components. Add a state library only when context performance becomes a bottleneck.

```tsx
// Local state - default choice
const [isOpen, setIsOpen] = useState(false);

// Lifted state - when siblings share
// Parent holds state, passes down through props

// Context - when deeply nested
const ThemeContext = createContext<Theme>("light");
```

## Composition over inheritance

Favor composition patterns:
- **Compound components** - related components that share implicit state (e.g., `<Tabs>`, `<Tabs.Panel>`)
- **Render props** - when the parent controls what a child renders
- **Children prop** - the default composition mechanism. Prefer this over render props when possible.
- **Slots via props** - pass ReactNode props for header, footer, or custom sections

```tsx
// Compound component pattern
function Card({ children }: { children: React.ReactNode }) {
  return <div className="rounded-lg border">{children}</div>;
}
Card.Header = function CardHeader({ children }: { children: React.ReactNode }) {
  return <div className="border-b px-4 py-3">{children}</div>;
};
Card.Body = function CardBody({ children }: { children: React.ReactNode }) {
  return <div className="p-4">{children}</div>;
};
```

## Error and loading states

Every data-dependent component handles three states:
- **Loading** - show a skeleton or spinner while data resolves
- **Error** - show an error message with a retry action
- **Empty** - show a meaningful empty state, not just blank space

Prefer `components/ui` primitives (`Banner`, `EmptyState`, spinner) when they exist:

```tsx
if (isLoading) return <Spinner />;
if (error) return <Banner tone="danger" action={{ label: "Retry", onClick: refetch }}>{error.message}</Banner>;
if (!data.length) return <EmptyState message="No items yet." />;
return <ItemList items={data} />;
```

## When NOT to create a component

- Single-use markup under ~10 lines - inline it
- Pure layout wrappers - use Tailwind's container/spacing utilities instead
- Over-abstraction of **feature** widgets - wait until the pattern appears 3 times before extracting (see extract-after-3 above)
- Do not skip `components/ui` for chrome the draft already names as reusable product controls

## Spacecraft integration

- Slice order: inventory → `ui/` primitive(s) → compose page → route
- Components correspond to `plan.json` acceptance checks - one feature slice per check
- Capture evidence with `spacecraft evidence "<name>:component" -- npx vitest run`
- Record component decisions in `decisions.md` when choosing between composition patterns or adding Storybook
- Reference mission `spec.md` for the feature scope - don't build beyond it
