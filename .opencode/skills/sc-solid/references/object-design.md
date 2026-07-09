# Object-Oriented Design

> Consult when: designing a new class, deciding between inheritance and composition, or choosing between value object and entity.

## Responsibility-Driven Design

Objects are defined by their **responsibilities**, not their data.

For every class, ask:
1. "What does it **know**?" (data it holds)
2. "What does it **do**?" (behaviors it provides)
3. "What does it **decide**?" (logic it owns)

Then verify:
1. "What pattern is this?" — Which stereotype? (see below)
2. "Is it doing too much?" — Check calisthenics limits.

If you can't answer clearly, the class needs refactoring.

---

## Object Stereotypes

Every class fits one (rarely two). If it fits three, split it.

| Stereotype | Role | Example |
|------------|------|---------|
| **Information Holder** | Holds data, minimal behavior | `User`, `Product`, `Address` |
| **Structurer** | Maintains relationships | `OrderItems`, `UserGroup` |
| **Service Provider** | Performs work, stateless | `PaymentProcessor`, `EmailSender` |
| **Coordinator** | Orchestrates workflow | `CheckoutService` |
| **Controller** | Makes decisions, delegates | `OrderController` |
| **Interfacer** | Transforms between systems | `UserAdapter`, `DatabaseMapper` |

---

## Value Objects vs Entities

| | Value Object | Entity |
|---|---|---|
| **Identity** | None — defined by attributes | Has identity (survives attribute changes) |
| **Mutability** | Immutable | Mutable via methods |
| **Equality** | By value | By identity |
| **Examples** | `Email`, `Money`, `DateRange` | `User`, `Order`, `Product` |

```typescript
// Value Object
class Money {
  constructor(readonly amount: number, readonly currency: string) {}
  add(other: Money): Money {
    if (this.currency !== other.currency) throw new CurrencyMismatch();
    return new Money(this.amount + other.amount, this.currency);
  }
}

// Entity
class User {
  constructor(readonly id: UserId, private email: Email, private name: string) {}
  changeEmail(newEmail: Email): void { this.email = newEmail; } // same user, new email
  equals(other: User): boolean { return this.id.equals(other.id); }
}
```

---

## Tell, Don't Ask

Command objects to do work. Don't interrogate them and do the work yourself.

```typescript
// BAD: Ask, then decide, then act
if (account.balance >= amount) {
  account.balance -= amount;
}

// GOOD: Tell the object
const result = account.withdraw(amount);
```

The object that owns the data should own the behavior.

---

## Design by Contract

Every method has:
- **Preconditions** — what must be true before calling
- **Postconditions** — what will be true after calling
- **Invariants** — what is always true about the object

```typescript
class BankAccount {
  // INVARIANT: balance >= 0
  // PRECONDITION: amount > 0
  // POSTCONDITION: balance decreased by amount, or error
  withdraw(amount: Money): Result {
    if (amount.isNegativeOrZero()) return Result.fail("Invalid amount");
    if (this.balance.lessThan(amount)) return Result.fail("Insufficient funds");
    this.balance = this.balance.subtract(amount);
    return Result.ok();
  }
}
```

---

## Composition Over Inheritance

**Default to composition.** Use inheritance only for true "is-a" relationships or Template Method pattern.

```typescript
// PREFER: Composition
class User {
  constructor(private discountPolicy: DiscountPolicy) {}
  getDiscount() { return this.discountPolicy.calculate(); }
}
new User(new PremiumDiscount());
new User(new NoDiscount()); // pluggable

// AVOID: Inheritance for code reuse
class PremiumUser extends User { getDiscount() { return 20; } }
```

---

## Aggregates

A cluster of objects treated as one unit. One root. External code only references the root. Root enforces invariants for the whole cluster.

```typescript
// Order is the aggregate root
class Order {
  private items: OrderItem[] = [];
  addItem(product: Product, qty: number) { ... }  // through root
  // External code never does: order.items.push(...)
}
```

---

## Spacecraft integration

- Value objects live in domain layer. Entities too. Services and coordinators in application layer.
- When designing for testability, prefer composition — it makes dependency injection trivial.
- The "Tell, Don't Ask" rule aligns with public-interface-only testing. Both enforce behavior verification over state inspection.
