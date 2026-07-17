---
name: sc-solid
description: "Code quality discipline. Activate on SOLID, clean code, refactoring, architecture decisions, or code review."
---

# sc-solid

Design code for changeability, testability, and readability. Apply principles silently; surface only violations. Do not recite theory - act on it.

## When to use

Activate on these triggers:

- Writing or modifying production code
- Reviewing a diff or doing self-review
- Planning module boundaries or task decomposition
- Choosing between design alternatives
- Detecting and fixing code smells

## Workflow

### On every code change

1. **SOLID scan** - Check the class/module against the 5 principles (see `references/solid-principles.md`). Flag any violation.
2. **Smell check** - Scan for the 7 common smells (see `references/code-smell.md`). If found, fix before commit.
3. **Complexity gate** - Is this the simplest thing that works? Apply KISS, YAGNI, DRY-with-Rule-of-Three.
4. **Naming audit** - Are names consistent, specific, searchable? See `references/clean-code.md`.
5. **Dependency check** - Do dependencies point inward (domain ← infra)? See `references/architecture.md`.

### When planning architecture

1. Choose boundaries - feature-first vertical slices, horizontal layers with dependency rule
2. Pick patterns only when the problem matches - see `references/design-patterns.md`
3. Record ADRs in `decisions.md` for non-obvious choices
4. Verify against mission `spec.md` - does the architecture serve the acceptance criteria?

### When reviewing code

1. `references/code-smell.md` - scan the diff for smells
2. `references/solid-principles.md` - flag SRP/DIP violations (most common)
3. `references/clean-code.md` - naming consistency, calisthenics rules
4. `references/complexity.md` - accidental complexity, premature abstractions
5. `references/architecture.md` - circular deps, domain depending on infra

## Rules

### SOLID (non-negotiable)

- **SRP**: One reason to change. Split if described with "and."
- **OCP**: Add behavior via new classes, not edits to existing ones.
- **LSP**: Subtypes must be substitutable without altering correctness.
- **ISP**: No client forced to depend on unused methods.
- **DIP**: Depend on interfaces, inject implementations. Never `new Concrete()` in business logic.

### Clean code

- **Must**: Wrap domain primitives in value objects (Email, Money, UserId - never raw strings/numbers).
- **Must**: No `else` when early return works.
- **Must**: One dot per line (Law of Demeter).
- **Must**: Classes < 50 lines, methods < 10 lines. (Rationale: beyond these thresholds, a class likely violates SRP. These are empirically-derived limits from clean code practice - a 50-line class can be read in one screen; a 10-line method can be understood at a glance without scrolling or mental stack.)
- **Must not**: Create abstractions before the third duplication (Rule of Three).
- **Must not**: Use `in` operator on objects with untrusted keys - use `Object.hasOwn()`.

### Complexity

- KISS: simplest solution that passes the test. No speculative architecture.
- YAGNI: no "just in case" code. Build what the mission task asks for.
- DRY: wait for 3 duplications before extracting. Wrong abstraction > duplication.

### Architecture

- Dependency rule: outer layers depend inward. Domain knows nothing of infra.
- Feature-first: organize by feature (`users/`, `orders/`), not by layer (`controllers/`, `services/`).
- Contracts: define interfaces at module boundaries. Every infra impl satisfies a domain contract.

## Out of scope

- Test writing and TDD discipline - separate concern, handled by test infrastructure
- Evidence capture - handled by verification pipeline
- Git operations and branching - handled by git infrastructure
- Mission lifecycle - handled by mission management
- UI design - handled by design infrastructure

## Output format

```
SOLID scan: SRP=ok, OCP=ok, LSP=n/a, ISP=ok, DIP=violation (explain)
Smells: [none / list with fix]
Complexity: KISS=ok, YAGNI=ok, DRY=deferred (dup #1)
Verdict: [clean / fix-before-commit / needs-refactor-task]
```

## Checklist

Before committing code:

- [ ] SOLID scan clean for all new/modified classes
- [ ] No code smells introduced (or documented as deferred)
- [ ] Domain primitives wrapped (no raw strings/numbers for concepts)
- [ ] Dependencies point inward (no domain → infra refs)
- [ ] KISS/YAGNI applied - no speculative abstractions
- [ ] Rule of Three respected - no premature extraction
- [ ] Naming consistent with existing codebase vocabulary

## References

- `references/solid-principles.md` - 5 principles with detection questions and examples
- `references/clean-code.md` - naming, calisthenics, formatting conventions
- `references/code-smell.md` - 7 common smells with before/after and fix patterns
- `references/complexity.md` - KISS, YAGNI, DRY, Rule of Three, technical debt
- `references/architecture.md` - dependency rule, feature-first, ports-adapters
- `references/object-design.md` - stereotypes, Tell Don't Ask, value objects vs entities
- `references/design-patterns.md` - when patterns help, when they hurt
