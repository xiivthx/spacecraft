# Clean Code

> Consult when: naming anything, structuring a class, or wondering if code is readable. Apply consistently — inconsistent naming is worse than imperfect naming.

## Naming (priority order)

| # | Rule | Good | Bad |
|---|------|------|-----|
| 1 | **Consistency** | `getUser`, `getOrder`, `getProduct` | `getUser`, `fetchCustomer`, `retrieveClient` |
| 2 | **Specificity** | `OrderRepository`, `validatePayment` | `DataManager`, `processInfo` |
| 3 | **Domain language** | `activeCustomers` | `filteredUserArray` |
| 4 | **Brevity** | `activeUsers` | `listOfAllActiveUsersInTheSystem` |
| 5 | **Searchability** | `fetchOrderSummary()` | `data` |
| 6 | **Pronounceability** | `timestamp` | `genymdhms` |

**Must not**: use `data`, `info`, `manager`, `handler`, `processor`, `utils` — these mean nothing.

---

## Value Objects (mandatory)

Wrap every domain concept. Never pass raw primitives for domain types.

```typescript
// MUST: wrap these
class UserId   { constructor(readonly value: string) {} }
class Email    { constructor(readonly value: string) { /* validate */ } }
class Money    { constructor(readonly amount: number, readonly currency: string) {} }

// NEVER:
function createOrder(userId: string, email: string, amount: number) {} // primitive obsession
function createOrder(userId: UserId, email: Email, amount: Money) {}   // correct
```

---

## Structure Rules

- **One indentation level** per method. Extract if deeper.
- **No `else`** — use early return, guard clause, or polymorphism.
  ```typescript
  // BAD
  if (user.isPremium) { return 20; } else { return 0; }
  // GOOD
  if (user.isPremium) return 20; return 0;
  ```
- **One dot per line** (Law of Demeter). `order.customer.address.city` → `order.shippingCity()`.
- **Classes < 50 lines, methods < 10 lines, files < 100 lines.** If larger, split.
- **Max 2 instance variables** per class. Compose with smaller objects.
- **No getters/setters that expose raw state.** Objects have behavior, not data bags. `account.withdraw(amount)`, not `account.setBalance(account.getBalance() - amount)`.
- **`Object.hasOwn(obj, key)`** for untrusted string checks — never `in` operator (matches prototype).

---

## Object Calisthenics (9 Rules)

Apply strictly during implementation, relax slightly in review after justification.

1. One level of indentation per method
2. Don't use `else`
3. Wrap all primitives and strings
4. First-class collections (collection-only classes)
5. One dot per line
6. Don't abbreviate
7. Keep entities small (<50 lines class, <10 method)
8. Max 2 instance variables per class
9. No getters/setters — behavior over data

---

## Comments

**Only explain WHY, not WHAT.** Code explains what. Comments explain business reasons, non-obvious decisions, or warnings.

```typescript
// BAD: redundant
// Add 1 to counter
counter++;

// GOOD: explains why
// Compensate for 0-based indexing in legacy API
counter++;
```

Prefer renaming over commenting. If you need a comment to explain what code does, the name is wrong.

---

## Spacecraft integration

- Match existing codebase naming conventions found in the mission's code. Read existing files before naming.
- Value objects become part of the domain layer — they should be in the `domain/` folder per the architecture reference.
- Calisthenics rules 8 and 9 (max 2 ivars, no getters/setters) are aspirational — flag violations in review but don't block small pragmatic classes.
