> Consult when: building new React components, choosing between composition patterns, or deciding data flow strategy.

# Component patterns

Operational patterns for React components in a TypeScript + Vite + Tailwind CSS stack.

## Component structure

Prefer function components with typed props. Keep components focused — one responsibility per component.

```tsx
// src/components/UserCard.tsx
interface UserCardProps {
  name: string;
  email: string;
  onSelect: (email: string) => void;
}

export function UserCard({ name, email, onSelect }: UserCardProps) {
  return (
    <article
      className="rounded-lg border p-4 hover:shadow-md"
      onClick={() => onSelect(email)}
      role="button"
      tabIndex={0}
      onKeyDown={(e) => e.key === "Enter" && onSelect(email)}
    >
      <h3 className="text-lg font-semibold">{name}</h3>
      <p className="text-sm text-gray-600">{email}</p>
    </article>
  );
}
```

Key conventions:
- Props interface named `<Component>Props`, defined above the component
- Export the component as a named export, not default
- Co-locate styles with markup using Tailwind utility classes
- Interactive elements must be keyboard accessible

## Data flow

**Props down, callbacks up.** Pass data through props. Communicate upward through callback props.

**Avoid premature state management.** Start with local `useState`. Lift state up when siblings need it. Introduce context when props pass through >2 intermediate components. Add a state library only when context performance becomes a bottleneck.

```tsx
// Local state — default choice
const [isOpen, setIsOpen] = useState(false);

// Lifted state — when siblings share
// Parent holds state, passes down through props

// Context — when deeply nested
const ThemeContext = createContext<Theme>("light");
```

## Composition over inheritance

Favor composition patterns:
- **Compound components** — related components that share implicit state (e.g., `<Tabs>`, `<Tabs.Panel>`)
- **Render props** — when the parent controls what a child renders
- **Children prop** — the default composition mechanism. Prefer this over render props when possible.
- **Slots via props** — pass ReactNode props for header, footer, or custom sections

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
- **Loading** — show a skeleton or spinner while data resolves
- **Error** — show an error message with a retry action
- **Empty** — show a meaningful empty state, not just blank space

```tsx
if (isLoading) return <Skeleton />;
if (error) return <Error message={error.message} onRetry={refetch} />;
if (!data.length) return <EmptyState message="No items yet." />;
return <ItemList items={data} />;
```

## When NOT to create a component

- Single-use markup under ~10 lines — inline it
- Pure layout wrappers — use Tailwind's container/spacing utilities instead
- Over-abstraction — wait until the pattern appears 3 times before extracting

## Spacecraft integration

- Components correspond to `plan.json` acceptance checks — one slice per check
- Capture evidence with `scripts/spacecraft evidence "<name>:component" -- npx vitest run`
- Record component decisions in `decisions.md` when choosing between composition patterns
- Reference mission `spec.md` for the feature scope — don't build beyond it
