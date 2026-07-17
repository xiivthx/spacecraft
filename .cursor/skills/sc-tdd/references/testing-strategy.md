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

## Common Mistakes

| Mistake | Fix |
|---------|-----|
| Testing implementation (brittle) | Test behavior through public interface |
| Too many mocks (tests prove nothing) | Use fakes, real objects when possible |
| Shared state between tests (flaky) | Isolate - fresh state per test |
| Testing trivial code | Focus on logic, edge cases, invariants |
| Slow test suite | Keep unit tests fast. Optimize integration tests. |

---

## Spacecraft integration

- Unit tests for domain, integration for boundaries.
- Production code must satisfy the pyramid distribution. No unit tests = not done.
- Evidence must include test type and layer.
- Evidence label convention: `"unit:domain:Money.add"`, `"integration:infra:PostgresOrderRepo"`.
