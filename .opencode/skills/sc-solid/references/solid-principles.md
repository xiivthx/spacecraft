# SOLID Principles

> Consult when: designing a new class, reviewing a diff, or choosing module boundaries. Apply all 5 — do not cherry-pick.

## Quick Reference

| Principle | Test | Violation Sign |
|-----------|------|----------------|
| **S**RP | One reason to change? | Described with "and" |
| **O**CP | Add via new class, not edit? | `if/switch` on type codes |
| **L**SP | Subtypes substitutable? | Type-checking in calling code |
| **I**SP | No unused methods? | Empty method bodies, `throw new Error("unused")` |
| **D**IP | Depend on interface, not concrete? | `new Concrete()` inside business logic |

---

## S — Single Responsibility

**One reason to change.** If describing the class requires "and," split it.

```typescript
// VIOLATION: persistence + presentation + business logic
class Order {
  calculateTotal() { ... }
  saveToDatabase() { ... }     // second reason to change
  generateInvoice() { ... }    // third reason to change
}

// FIX: one responsibility per class
class Order { calculateTotal() { ... } }
class OrderRepository { save(order) { ... } }
class InvoiceGenerator { generate(order) { ... } }
```

**Detection:** Would different stakeholders request changes to different parts? Would a database migration force edits to this class?

---

## O — Open/Closed

**Open for extension, closed for modification.** Add behavior by adding code, not editing existing code.

```typescript
// VIOLATION: must edit to add new behavior
function calculateShipping(method: string, value: number) {
  if (method === 'standard') return value < 50 ? 5 : 0;
  if (method === 'express') return 15;
  // new method = new if-branch = edit existing code
}

// FIX: strategy pattern
interface ShippingMethod {
  cost(orderValue: number): number;
}
class StandardShipping implements ShippingMethod { cost(v) { return v < 50 ? 5 : 0; } }
class ExpressShipping implements ShippingMethod { cost(v) { return 15; } }
```

**Detection:** Do you add `if`/`switch` branches to add features? That's OCP violation.

---

## L — Liskov Substitution

**Subtypes must not break the parent contract.** If `Bird.fly()` exists, `Penguin extends Bird` that throws on `fly()` is an LSP violation.

```typescript
// VIOLATION: subtype changes semantics
class Discount {
  getDiscount(): number { return 0; }  // contract: non-negative
}
class Surcharge extends Discount {
  getDiscount(): number { return -10; } // breaks contract!
}
```

**Detection:** Do you find yourself writing `if (x instanceof SpecificSubtype)` checks? The hierarchy is wrong.

---

## I — Interface Segregation

**No client forced to depend on methods it doesn't use.** Split fat interfaces.

```typescript
// VIOLATION: fat interface
interface WarehouseDevice {
  printLabel(id: string): void;
  scanBarcode(): string;
  packageItem(id: string): void;
}
class BasicPrinter implements WarehouseDevice {
  printLabel(id) { ... }
  scanBarcode() { throw new Error("unsupported"); } // forced!
  packageItem(id) { throw new Error("unsupported"); }
}

// FIX: segregated
interface LabelPrinter { printLabel(id: string): void; }
interface BarcodeScanner { scanBarcode(): string; }
```

**Detection:** `throw new Error("Not implemented")` or empty methods = ISP violation.

---

## D — Dependency Inversion

**Depend on abstractions, not concretions.** High-level policy must not know about low-level details. Inject interfaces.

```typescript
// VIOLATION: hardwired dependency
class OrderService {
  private email = new SendGridClient(process.env.KEY); // locked in
}

// FIX: inject interface
interface EmailSender { send(to: string, body: string): void; }
class OrderService {
  constructor(private email: EmailSender) {}
}
// Now: test with MockEmailSender, production with SendGridSender
```

**Detection:** Can you swap the database/API without changing business logic? If not, DIP is violated.

---

## Spacecraft integration

- Run a SOLID scan on every new class before marking a `plan.json` task done.
- SRP and DIP violations are the most common and costly — prioritize these in review.
- Use SOLID to decide module boundaries during `plan.json` task decomposition.
- Record complex SOLID trade-offs in `decisions.md` under the mission.
