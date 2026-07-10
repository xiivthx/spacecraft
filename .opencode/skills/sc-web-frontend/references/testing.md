> Consult when: writing component tests, choosing query strategies, mocking browser APIs, or setting up Vitest with React Testing Library.

# Testing patterns

Operational patterns for testing React components with Vitest and React Testing Library.

## Test structure

Every component test follows the AAA pattern: Arrange → Act → Assert.

```tsx
import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { describe, it, expect } from "vitest";
import { UserCard } from "./UserCard";

describe("UserCard", () => {
  it("renders user name and email", () => {
    // Arrange
    render(<UserCard name="Alice" email="alice@example.com" onSelect={vi.fn()} />);

    // Assert
    expect(screen.getByText("Alice")).toBeInTheDocument();
    expect(screen.getByText("alice@example.com")).toBeInTheDocument();
  });

  it("calls onSelect with email when clicked", async () => {
    // Arrange
    const onSelect = vi.fn();
    const user = userEvent.setup();
    render(<UserCard name="Alice" email="alice@example.com" onSelect={onSelect} />);

    // Act
    await user.click(screen.getByRole("button"));

    // Assert
    expect(onSelect).toHaveBeenCalledWith("alice@example.com");
  });
});
```

## Query priority

Use queries in this order — each successive query is less specific and more fragile:

1. **`getByRole`** — preferred. Mirrors accessibility tree. Use `name` option for text matching.
2. **`getByLabelText`** — for form fields. Pairs label with input.
3. **`getByText`** — for non-interactive text content.
4. **`getByTestId`** — last resort. Only when no other query works. Use `data-testid` sparingly.

```tsx
// Preferred
screen.getByRole("button", { name: "Submit" });
screen.getByLabelText("Email address");

// Acceptable
screen.getByText("Welcome back");

// Last resort
screen.getByTestId("user-menu");
```

## Testing behavior, not implementation

Test what the user sees and does. Never test internal state, method calls on component instances, or React internals.

```tsx
// GOOD — tests behavior
await user.click(screen.getByRole("button", { name: "Delete" }));
expect(screen.queryByText("Item 1")).not.toBeInTheDocument();

// BAD — tests implementation
expect(component.state.isOpen).toBe(true);
expect(handleClick).toHaveBeenCalled();
```

## Async testing

Use `findBy*` queries for elements that appear asynchronously. Use `waitFor` for assertions that need retry.

```tsx
it("loads and displays users", async () => {
  render(<UserList />);

  // findBy queries retry until the element appears or timeout
  expect(await screen.findByText("Alice")).toBeInTheDocument();

  // Or use waitFor for custom assertions
  await waitFor(() => {
    expect(screen.getAllByRole("listitem")).toHaveLength(3);
  });
});
```

## Mocking

Mock only at system boundaries. Do not mock React components or hooks under your control.

- **API calls** — mock `fetch` or use MSW (Mock Service Worker) for integration-like tests
- **Browser APIs** — mock `window.matchMedia` for responsive tests, `IntersectionObserver` for visibility
- **Timers** — use `vi.useFakeTimers()` for `setTimeout`/`setInterval`
- **Do not mock** — React hooks (`useState`, `useEffect`), child components, context providers

```tsx
// MSW handler for API mocking
import { http, HttpResponse } from "msw";
const handlers = [
  http.get("/api/users", () => {
    return HttpResponse.json([{ id: 1, name: "Alice" }]);
  }),
];
```

## Test file placement

Co-locate tests with components. Use `.test.tsx` suffix.

```
src/components/UserCard.tsx
src/components/UserCard.test.tsx
```

## Spacecraft integration

- Run tests per slice: `scripts/spacecraft evidence "<slice>:test" -- npx vitest run src/components/<name>.test.tsx`
- Full suite after refactor: `scripts/spacecraft evidence "<task>-functional" -- npx vitest run`
- Evidence must include test count, pass/fail, and duration
- Test files are product code — stage them alongside components in checkpoint commits
