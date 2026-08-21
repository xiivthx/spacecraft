# Testing Strategy

> Consult when: deciding what kind of test to write, or designing test doubles for infrastructure boundaries.

## The Testing Pyramid

```
       /\       E2E - critical paths only. Slow, brittle.
      /  \      
     /----\     Integration - component boundaries. Medium speed.
    /      \    
   /--------\   Unit - most tests. Fast, isolated, comprehensive.
  /          \  
```

Distribution: ~70% unit, ~20% integration, ~10% E2E.

---

## Test Types

### Unit Tests
One class/function, isolated. No external deps (use fakes/stubs). **Most tests.**

```typescript
it('calculates total from items', () => {
  const order = new Order();
  order.addItem({ price: 100 });
  order.addItem({ price: 50 });
  expect(order.total()).toBe(150);
});
```

### Integration Tests
Multiple components together. Test boundaries between modules.

### E2E Tests
Full system, user perspective. Critical paths only. Use sparingly.

---

## AAA - Arrange, Act, Assert

Structure every test:

```typescript
it('applies 20% discount for premium users', () => {
  // ARRANGE
  const user = new User({ premium: true });
  const cart = new Cart(user);
  cart.add({ price: 100 });
  // ACT
  const total = cart.checkout();
  // ASSERT
  expect(total).toBe(80);
});
```

---

## Deep assert

**One condition per test.** **One acceptance → one RED.** Expected values from design-contract Edge matrix + approved-scenarios frozen literals - never recomputed like the implementation.

**Banned shallow patterns:** `expect(true)`, typeof-only, status-only, phrase-echo. Hollow green ≠ behavioral proof.

**Error paths:** assert **ErrorCode** or error **instance** (project typed equivalent OK). Message-substring / HTTP status alone = shallow → rewrite.

See sc-tester Ban / Deep assert; sc-tdd Rules.

---

## Test Doubles

| Double | Purpose | When |
|--------|---------|------|
| **Dummy** | Placeholder, never used | Satisfy constructor params |
| **Stub** | Returns predefined values | Control test inputs |
| **Spy** | Records calls | Verify side effects |
| **Mock** | Verifies expected interactions | Contract verification |
| **Fake** | Working simplified impl | In-memory repo, test DB |

**Rule:** Fake > Stub > Mock. Prefer real implementations when fast enough. Mock only at system boundaries (APIs, payment, email).

---

## Testing by Layer

- **Domain** (most - unit): business rules, value objects, entities. No mocks needed.
- **Application** (some - integration): use case orchestration. Mock infrastructure.
- **Infrastructure** (few - integration): real database/API adapters. Use test DB.

---

## Contract Tests

Verify that all implementations of an interface satisfy its contract:

```typescript
function testRepoContract(makeRepo: () => UserRepo) {
  it('returns null for missing entity', async () => {
    expect(await makeRepo().findById('none')).toBeNull();
  });
}
testRepoContract(() => new InMemoryUserRepo());
testRepoContract(() => new PostgresUserRepo(testDb));
```

---

## Composition tests

**Composition** = a user-visible path that mints something in step A and consumes it in step B (create→use, join→claim, auth→mutate, invite→accept).

Isolated unit tests with mocked `fetch` / happy fixtures prove each seam in isolation. They **will not** catch:

- Credentials returned by create that cannot unlock the next step
- Flags/state that mark an entity “ready” without a usable secret
- Empty public API base / reverse-proxy / same-origin rewrite miswiring (“could not reach server”)
- UI that only spies on the first call while the real chain is broken

**Required when shipping such a feature:**

1. At least one **API or inject-style contract** that walks the chain with values returned by earlier steps (not hard-coded fixtures that skip the mint).
2. If the client uses relative `/api` or a dev proxy, at least one test or config guard that empty/default base still targets the documented same-origin path, and local/dev config stays aligned.
3. A production composition bug is closed only when the composition test is added in the same change — not by another mocked unit test alone.

```typescript
it('create credentials unlock the follow-up step', async () => {
  const created = await app.inject({ method: 'POST', url: '/resources', payload: { name: 'x' } });
  expect(created.statusCode).toBe(201);
  const { id, secret } = created.json();

  const join = await app.inject({
    method: 'POST',
    url: '/sessions',
    payload: { resourceId: id, secret },
  });
  expect(join.statusCode).toBe(200);

  const finish = await app.inject({
    method: 'POST',
    url: `/resources/${id}/claim`,
    payload: { sessionToken: join.json().token, secret },
  });
  expect(finish.statusCode).toBe(200);
});
```

Evidence labels: prefer `integration:composition:<flow-name>` so review can see the chain was exercised.

---

## Test data diversity

When designing fixtures, prefer category diversity (positive / negative / boundary / exploratory / security-shaped) over many near-duplicate happy paths. See `sc-discuss/references/test-data-design.md`.

---

## Property-based testing (PBT)

**Core logic** = design-contract modules with branching business rules / pure domain / state machines - not chrome, docs, or thin adapters.

**Product path Must:** **100%** of those core-logic modules get property-based invariants + generators, using a project-existing lib (`fast-check` / Hypothesis / equivalent). Capture evidence with label prefix `pbt-…`.

**Disposition when PBT does not run** (exact greppable lines; SoT in `docs/mission-artifacts.md`):

- `Pbt skipped: no project pbt tool`
- `Pbt skipped: not core logic`
- `Pbt waived: <reason>`

**Must not** invent installing a PBT library mid-mission unless the human asked. Missing `pbt-…` without skip/waive ⇒ judge-REFUTE material.

Prefer invariants that encode domain laws (round-trip, idempotence, monotonicity, state-machine transitions) over example-only cases.

---

## Common Mistakes

| Mistake | Fix |
|---------|-----|
| Testing implementation (brittle) | Test behavior through public interface |
| Too many mocks (tests prove nothing) | Use fakes, real objects when possible |
| Shared state between tests (flaky) | Isolate - fresh state per test |
| Testing trivial code | Focus on logic, edge cases, invariants |
| Slow test suite | Keep unit tests fast. Optimize integration tests. |
| Mocked seam only for multi-step flows | Add a composition contract across create→use / join→claim |
| Asserting only create response shape | Also assert the next step succeeds with returned credentials |
| Shallow assert (`expect(true)`, typeof/status/phrase-echo) | Deep assert from design-contract + approved-scenarios oracles |
| Error path checks status/message only | Assert ErrorCode or error instance (typed equivalent OK) |
| Multi-condition / multi-acceptance RED | One condition per test; one acceptance → one RED |
| Core-logic module ships with example tests only | Add property-based invariants + generators (`fast-check` / Hypothesis / project equivalent) or greppable `Pbt skipped`/`Pbt waived` |
| Inventing a PBT lib install mid-mission | Use project-existing tool only; else `Pbt skipped: no project pbt tool` |

---

## Spacecraft integration

- Unit tests for domain, integration for boundaries.
- Production code must satisfy the pyramid distribution. No unit tests = not done.
- Evidence must include test type and layer.
- Evidence label convention: `"unit:domain:Money.add"`, `"integration:infra:PostgresOrderRepo"`.
- Core-logic modules: property-based evidence `pbt-…` or greppable PBT skip/waive per `docs/mission-artifacts.md`.
