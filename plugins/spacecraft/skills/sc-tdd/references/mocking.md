# When to Mock

Mock at **system boundaries** only:

- External APIs (payment, email, etc.)
- Databases (sometimes - prefer test DB)
- Time/randomness
- File system (sometimes)

Don't mock:

- Your own classes/modules
- Internal collaborators
- Anything you control

## Designing for Mockability

At system boundaries, design interfaces that are easy to mock:

**1. Use dependency injection**

Pass external dependencies in rather than creating them internally:

```typescript
// Easy to mock
function processPayment(order, paymentClient) {
  return paymentClient.charge(order.total);
}

// Hard to mock
function processPayment(order) {
  const client = new StripeClient(process.env.STRIPE_KEY);
  return client.charge(order.total);
}
```

**2. Prefer SDK-style interfaces over generic fetchers**

Create specific functions for each external operation instead of one generic function with conditional logic:

```typescript
// GOOD: Each function is independently mockable
const api = {
  getUser: (id) => fetch(`/users/${id}`),
  getOrders: (userId) => fetch(`/users/${userId}/orders`),
  createOrder: (data) => fetch('/orders', { method: 'POST', body: data }),
};

// BAD: Mocking requires conditional logic inside the mock
const api = {
  fetch: (endpoint, options) => fetch(endpoint, options),
};
```

The SDK approach means:
- Each mock returns one specific shape
- No conditional logic in test setup
- Easier to see which endpoints a test exercises
- Type safety per endpoint

## Mocked seams vs composition

Mocking `fetch` (or spying on `startJoin` / `createX`) in a UI unit test is fine for DOM and wiring **inside that component**. It is **not** sufficient coverage for:

- Token / password / ID handoff from create (or mint) to the next API call
- Session or claim flows that depend on server-side state set at create time
- Dev proxy / empty public base URL / same-origin `/api` rewrite behavior

For those, add a composition contract (real HTTP inject or API test app) that uses values returned by earlier steps. Do not close a production composition bug with only another happy mock. See `testing-strategy.md` (Composition tests) and sc-tdd SKILL (Composition paths).
