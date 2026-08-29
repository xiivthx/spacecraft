---
name: sc-solid
description: "Code quality discipline. Activate on SOLID, clean code, refactoring, architecture decisions, or code review."
disable-model-invocation: true
---

# sc-solid

Design code for changeability, testability, and readability. Apply principles silently; surface only violations. Do not recite theory - act on it.

## When to use

Activate on production code change, diff/self-review, module-boundary planning, design alternatives, smell fix.
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

On-demand only (not every silent scan): when a human pastes a snippet to analyze/explain, or asks for a beginner-friendly walkthrough, follow `references/code-walkthrough.md`. Routine diff review stays on the checklist below.

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
- **Must**: Classes < 50 lines, methods < 10 lines.
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

TDD / test writing · evidence capture · git · mission lifecycle · UI design

## Output format

```
SOLID scan: SRP=ok, OCP=ok, LSP=n/a, ISP=ok, DIP=violation (explain)
Smells: [none / list with fix]
Complexity: KISS=ok, YAGNI=ok, DRY=deferred (dup #1)
Verdict: [clean / fix-before-commit / needs-refactor-task]
```

## Checklist

SOLID clean · no new smells (or deferred) · domain primitives wrapped · deps inward · KISS/YAGNI · Rule of Three · naming matches codebase.

## References

- `references/solid-principles.md` · `clean-code.md` · `code-smell.md` · `complexity.md` · `architecture.md` · `object-design.md` · `design-patterns.md`
- `references/code-walkthrough.md` - on-demand snippet analysis (not silent SOLID default)
