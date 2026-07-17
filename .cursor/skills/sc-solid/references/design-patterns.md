# Design Patterns

> Consult when: a design problem matches a known pattern. Do NOT consult when looking for a place to apply a pattern - patterns solve problems, they don't create them.

## The Golden Rule

> Let patterns emerge from refactoring. Don't force them upfront.

A pattern should: solve a problem you HAVE, simplify the code, and be understood by the team. If any of these is false, don't use it.

---

## When a Pattern Helps

| Problem | Pattern | Why |
|---------|---------|-----|
| Multiple variants, chosen at runtime | **Strategy** | Encapsulate interchangeable algorithms |
| One instance needed globally | **Singleton** (rare - prefer DI) | Shared resource |
| Complex object construction | **Builder** | Step-by-step, readable |
| Incompatible interfaces | **Adapter** | Third-party/library integration |
| Add behavior without modifying | **Decorator** | Wrapping with extra behavior |
| Undo/redo, command queue | **Command** | Encapsulate actions as objects |
| Notify multiple listeners | **Observer** | Pub/sub, event systems |
| Tree structure, uniform treatment | **Composite** | Part-whole hierarchy |
| Common algorithm, varying steps | **Template Method** | Skeleton with hooks |

---

## Patterns to Avoid Forcing

| Pattern | The Problem with Early Application |
|---------|-----------------------------------|
| **Singleton** | Usually DI is better. Singleton makes testing hard and hides dependencies. |
| **Factory** | If construction is simple, just use `new`. Add factory when construction logic gets complex. |
| **Abstract Factory** | Wait until you have ≥2 factory families. Never start here. |

---

## Anti-Pattern Warnings

| Anti-Pattern | Sign |
|--------------|------|
| **Golden Hammer** | Using the same pattern for everything |
| **God Object** | One class doing everything |
| **Over-engineering** | More abstraction layers than concrete code |
| **Pattern cargo cult** | Applying a pattern because "that's how it's done" |

---

## Spacecraft integration

- Code that uses a pattern should have a one-line comment naming the pattern (e.g., `// Strategy: pricing algorithm`).
- Apply patterns during the refactor phase after tests pass. Never introduce patterns during initial implementation.
- If a pattern adds more lines than it saves, skip it. KISS over pattern purity.
