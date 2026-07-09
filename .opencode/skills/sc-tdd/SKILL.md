---
name: sc-tdd
description: >
  Test-driven development discipline. Activate on "TDD", "red-green-refactor", "test-first", or integration tests.
license: MIT
compatibility: opencode
metadata:
  version: 1.0
  audience: xiivthx
---

# sc-tdd

Red → Green → Repeat. Write the test first, then the minimum code to pass. Refactoring happens after the loop, during review. Every cycle produces one vertical slice: one seam, one test, one passing implementation.

## When to use

Activate on these triggers:

- Writing a failing test before production code exists
- Implementing the minimum code to pass a specific failing test
- User explicitly says "TDD", "test-first", "red-green-refactor"
- Adding integration tests that verify behavior through public interfaces
- A `plan.json` task requires test-first acceptance checks

## Workflow

### Per acceptance check (one cycle)

1. **Seams** — Confirm the public boundary under test. Ask: *"What's the public interface we're testing?"* Write it down. No test at an unconfirmed seam.

2. **Red** — Write exactly one failing test. Must verify behavior through the public interface. Expected values from an independent source: a literal, a worked example, the spec. Never recompute the expected value the way the code computes it. Verify the test fails before proceeding.

3. **Green** — Write the absolute minimum production code to pass the failing test. No speculative features. No refactoring. No anticipating future tests. Just enough to make the test pass.

4. **Verify** — Run the test suite for the affected package. Capture evidence: `scripts/spacecraft evidence "<label>" -- <test command>`. Confirm the task's acceptance check is satisfied.

5. **Repeat** — Next acceptance check. One seam, one test, one implementation. Each cycle is a tracer bullet informed by the last.

### After all checks pass

Move to review. Refactoring belongs here, not in the red-green loop: extract helpers, improve names, remove duplication (Rule of Three).

## Rules

### Test quality (non-negotiable)

- **Must**: Test behavior through public interfaces only. Tests survive refactors because they don't couple to internal structure.
- **Must**: Expected values from independent source — literal, worked example, spec. Never recompute expected value the same way the code does.
- **Must not**: Test private methods, mock internal collaborators, verify through side channels (e.g., querying DB instead of using the interface).

### Seams

- **Must**: Confirm seams before writing. Write the seam list down. No test without confirmed seam.
- **Must not**: Test against internals. A seam is where observable behavior crosses a boundary.

### Loop discipline

- **Must**: Red before green. No production code without a failing test.
- **Must**: One slice per cycle. One seam → one test → one implementation.
- **Must not**: Horizontal slice — bulk tests before bulk implementation. Tests written without implementation feedback verify imaginary behavior.

### Mocking

- **Must**: Mock at system boundaries only — external APIs, payment, email, time/randomness.
- **Must not**: Mock your own classes, internal collaborators, or anything under your control.
- **Prefer**: Dependency injection and SDK-style interfaces over generic fetchers. See `references/mocking.md`.

### Anti-patterns

| Type | Tell | Fix |
|------|------|-----|
| **Implementation-coupled** | Test breaks on refactor but behavior unchanged | Rewrite against public interface |
| **Tautological** | `expect(add(a,b)).toBe(a+b)` — can't disagree with code | Use a literal expected value |
| **Horizontal slicing** | All tests written before any implementation | Restart with vertical slices |

## Out of scope

- Code quality principles and architecture — separate concern
- Review workflow — handled at review stage
- Evidence capture — handled by verification pipeline
- Git operations — handled by git infrastructure
- Mission lifecycle — handled by mission management
- UI testing — handled by design infrastructure

## Output format

```
Seams: [list of confirmed public boundaries]
Cycle 1/3: <acceptance check description>
  Red: test "<name>" — FAILS (expected)
  Green: minimal impl in <file:line> — PASSES
  Evidence: <label>
```

## Checklist

- [ ] Seams confirmed (written down, approved)
- [ ] Every test written before its implementation
- [ ] Tests verify behavior through public interfaces only
- [ ] No implementation-coupled or tautological tests
- [ ] Vertical slices only — one seam per cycle
- [ ] All tests pass with minimal code (no speculative features)
- [ ] Mocks only at system boundaries
- [ ] Evidence captured for each passing test suite

---

## References

- `references/examples.md` — good and bad test patterns with detection rules
- `references/mocking.md` — when to mock, dependency injection, SDK-style interfaces
- `references/testing-strategy.md` — testing pyramid, test types, AAA, doubles, contracts
