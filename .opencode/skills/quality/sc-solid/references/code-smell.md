# Code Smells

> Consult when: reviewing a diff, doing self-review during `/sc-build`, or feeling that code is "fighting you." Not all smells need fixing — confirm with a test before refactoring.

## The 7 That Matter Most

### 1. Long Method

**Smell:** Method > 10 lines, doing multiple things.
**Fix:** Extract Method. Each extracted method does exactly one thing.
**Detect:** Can you name the method without "and"?

```typescript
// SMELL
function processOrder(order) {
  if (!order.items.length) throw Error();  // validation
  let total = 0;
  for (const item of order.items) total += item.price * item.quantity; // calculation
  db.orders.insert({ ...order, total });   // persistence
  email.send(order.customer.email);         // notification
}

// FIX
function processOrder(order) {
  validate(order);
  const total = calculateTotal(order);
  save(order, total);
  notify(order);
}
```

### 2. Large Class

**Smell:** Class > 50 lines or multiple unrelated methods.
**Fix:** Extract Class by responsibility.
**Detect:** Can you describe the class without "and"?

### 3. Feature Envy

**Smell:** Method uses another object's data more than its own.
**Fix:** Move Method to the envied class.
**Detect:** Does the method access `other.x`, `other.y`, `other.z` more than `this.*`?

```typescript
// SMELL: Order envies Customer's data
class Order {
  shippingCost(customer) {
    if (customer.country === 'US') return customer.state === 'CA' ? 10 : 15;
    return 25;
  }
}
// FIX: Move to Customer
class Customer {
  shippingCost() { /* same logic, but on owner */ }
}
```

### 4. Primitive Obsession

**Smell:** Raw strings/numbers for domain concepts (email, money, ID).
**Fix:** Wrap in value objects (see `clean-code.md`).
**Detect:** Do you see validation of the same primitive scattered across files?

### 5. Switch Statements

**Smell:** Type-code switching repeated across the codebase.
**Fix:** Replace with Polymorphism (Strategy pattern).
**Detect:** Same `switch(type)` or `if(type === ...)` chain appears in >1 place.

### 6. Inappropriate Intimacy

**Smell:** Classes reaching into each other's internals.
**Fix:** Tell, Don't Ask. Move the behavior to the data owner.
**Detect:** Does class A access private/protected members of class B?

### 7. Speculative Generality

**Smell:** "Just in case" abstractions, unused parameters, hook points for future needs.
**Fix:** Delete. YAGNI.
**Detect:** Is there a caller for this abstraction? If not, remove it.

---

## Quick Smell → Fix Table

| Smell | Fix |
|-------|-----|
| Long Method | Extract Method |
| Large Class | Extract Class |
| Long Parameter List (>3) | Introduce Parameter Object |
| Feature Envy | Move Method |
| Primitive Obsession | Wrap in Value Object |
| Switch Statements | Replace with Polymorphism |
| Data Clumps | Extract Class |
| Shotgun Surgery | Move Method/Field together |
| Speculative Generality | Delete (YAGNI) |
| Dead Code | Delete |

---

## Remediation Workflow

1. **Confirm** — Is this actually causing pain? Not all smells need fixing.
2. **Cover** — Ensure test coverage exists for the affected code.
3. **Refactor** — Small step, keep tests green.
4. **Commit** — One refactoring per checkpoint commit.

---

## Spacecraft integration

- Flag smells in code review. Critical smells (Primitive Obsession on domain types, Switch Statements on core logic) block ship.
- Fix smells before marking a `plan.json` task done.
- During self-review, scan for the 7 smells. If fixing would exceed 20 lines, create a follow-up task instead.
