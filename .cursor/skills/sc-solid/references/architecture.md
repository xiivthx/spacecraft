# Software Architecture

> Consult when: scoping a new mission, decomposing a plan into tasks, or choosing where a new class lives.

## The Goal

Enable the team to **add, change, remove, test, and deploy** features with minimal friction. Architecture decisions serve development velocity, not purity.

---

## The Dependency Rule

**Dependencies point INWARD. Domain knows nothing of infrastructure.**

```
Infrastructure  ───→  Application  ───→  Domain
    (outer)              (middle)         (inner)
```

- Inner layers define interfaces (ports). Outer layers implement them (adapters).
- Domain has ZERO imports from `pg`, `express`, `fetch`, `fs`, etc.
- Invert with interfaces:

```typescript
// Domain defines the contract (inner)
interface OrderRepository {
  save(order: Order): Promise<void>;
  findById(id: OrderId): Promise<Order | null>;
}

// Infrastructure implements it (outer)
class PostgresOrderRepo implements OrderRepository { ... }

// Domain service depends on abstraction
class OrderService {
  constructor(private repo: OrderRepository) {}  // DIP applied
}
```

---

## Vertical vs Horizontal

### Feature-first (vertical) - preferred

Group by feature, not by technical layer. Changes to "users" stay in `users/`.

```
src/
  users/         # UserController, UserService, UserRepository
  orders/        # OrderController, OrderService, OrderRepository
  shared/        # Truly shared: value objects, utilities
```

### Layer-first (horizontal) - anti-pattern

Avoid this. Changes to one feature scatter across layers.

```
src/
  controllers/   # UserController, OrderController
  services/      # UserService, OrderService
  repositories/  # UserRepository, OrderRepository
```

---

## Ports & Adapters (Hexagonal)

Domain at center. Adapters around the edges.

- **Port** = interface defined by domain (`OrderRepository`, `PaymentGateway`).
- **Adapter** = implementation of a port (`PostgresOrderRepo`, `StripeGateway`).

Every external system gets an adapter. Swapping Postgres for SQLite means one new adapter, zero domain changes.

---

## The Walking Skeleton

For a new mission: build the thinnest end-to-end slice first. Touches all layers, deployable from day one, proves the architecture works. Then flesh out features.

---

## Red Flags

| Flag | Why It Matters |
|------|----------------|
| Domain importing from `pg`/`express`/`fetch` | Dependency rule violated |
| Circular imports between modules | Boundaries are wrong |
| `utils/`, `common/`, `helpers/` that grow unbounded | Dumping ground - no cohesion |
| Shared mutable state across modules | Implicit coupling, hard to test |
| Database schema driving domain model | Domain should drive schema, not reverse |

---

## Spacecraft integration

- Use the dependency rule when decomposing tasks. Each `plan.json` task should touch one layer or one feature.
- Place new files following feature-first convention. Consult existing codebase structure - don't invent a new layout.
- Flag any domain file importing from infrastructure. That's a DIP/architecture violation that blocks ship.
- Record architectural decisions as ADRs in `decisions.md` under the mission.
