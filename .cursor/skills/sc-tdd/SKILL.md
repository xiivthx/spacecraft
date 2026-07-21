---
name: sc-tdd
description: "Test-driven development discipline. Activate on \"TDD\", \"Plan-Red-Green-Verify-Refactor-Review\", \"test-first\", or integration tests. Skip TDD for trivial code - just code and review."
---

# sc-tdd

Plan → Red → Green → Verify → Repeat. Then Refactor → Review when the feature is complete. Plan what to test before writing. Skip TDD when the test would be a trivial tautology - just code and review. Every cycle produces one vertical slice.

Under `/sc-run`, the Commander drives each `plan.json` acceptance as one cycle via Task(`sc-tester`) then Task(`sc-coder`), auto-commits after RED and GREEN, then runs a combine/refactor + functional gate before review.

## Principles

Distilled from common testing anti-patterns. Apply silently; surface violations.

| Principle | Anti-pattern defended |
|-----------|----------------------|
| **Test behavior, not implementation.** Test via public interfaces. A test that checks a struct field will break on any refactor. A test that checks `loginAsGuest()` survives. | Testing internal implementation |
| **Test critical code, not everything.** 20% of code causes 80% of bugs. Chasing 100% coverage wastes time on trivial paths. | Coverage obsession |
| **Every production bug gets a test.** A bug that shipped once should never ship twice. | Not converting bugs to tests |
| **Tests must be deterministic.** Flaky tests destroy trust in the whole suite. Fix or quarantine flaky tests immediately. | Flaky/slow tests |
| **TDD is a tool, not a religion.** Skip tests for trivial code, experiments, and spikes. Write tests before, during, or after - what matters is that they exist and verify behavior. | TDD as religion |
| **DRY applies to tests too.** Centralize fixtures, share helpers, avoid copy-paste. Test code is production code. | Test code as second-class |
| **Read the test framework docs first.** Know what parameterized tests, fixtures, mocks, and test categorization your framework offers before writing utilities. | Writing tests without docs |
| **Match test type to the code.** Business logic → unit tests. External systems → integration tests. Both are needed - they catch different things. | Wrong kind of tests |

## When to use

Activate on these triggers:

- Writing a failing test before production code exists
- Implementing the minimum code to pass a specific failing test
- User explicitly says "TDD", "test-first", "red-green-refactor"
- Adding integration tests that verify behavior through public interfaces
- A `plan.json` task requires test-first acceptance checks

## Workflow

### Triage: should this be TDD?

Before starting any cycle, ask: *"Would a test for this be a trivial tautology?"* If yes - skip TDD. Just code and review. TDD is a tool, not a religion - see Principles.

**Skip TDD when:**
- Simple getters/setters, pass-through wrappers, one-line delegations
- Pure config mappings (env → config object)
- Boilerplate that the framework generates or enforces
- **Struct-constructor tests**: creating a struct, then checking its own fields with zero transformation, marshaling, or serialization in between. The test is `assert(x == x)` - the type definition IS the spec. Review the struct directly.
- Any test that would be `expect(add(1,2)).toBe(3)` - the implementation IS the spec

**Use TDD when:**
- Behavior has branching logic, edge cases, or error states
- The correct output is not trivially obvious from the input
- The implementation could plausibly be wrong

Record skipped-TDD decisions for a task in the task output or plan.json notes.

### Per acceptance check (one cycle)

1. **Plan** - Before writing anything, decide:
   - What seam are we testing? (public boundary)
   - What is the expected behavior? (from spec, literal, worked example)
   - What is the independent expected value? (never recompute it the way the code does)
   - Is the test non-trivial? (re-apply triage gate)
   Write the seam and expected behavior down. No test at an unconfirmed, unplanned seam.

2. **Red** - Write exactly one failing test. Must verify behavior through the public interface. Verify the test fails before proceeding. If the test passes without implementation, it is not testing the right thing - reject and re-write.

3. **Green** - Write the absolute minimum production code to pass the failing test. No speculative features. No refactoring. No anticipating future tests.

4. **Verify** - Run the test suite for the affected package. Capture evidence: `spacecraft evidence "<label>" -- <test command>`. Confirm the task's acceptance check is satisfied.

5. **Repeat** - Next acceptance check. One plan, one test, one implementation per cycle. Each cycle is a tracer bullet informed by the last.

### After all checks pass for a feature

1. **Refactor** - Now you have the full picture. Extract helpers, improve names, remove duplication (Rule of Three), simplify logic. The tests protect you - refactor with confidence. Under `/sc-run`, Commander auto-commits this step (`refactor:`).

2. **Functional test gate** - Run the full test suite (unit + integration + functional). All old tests must pass alongside new tests. If anything breaks, fix the refactor, not the old tests. Capture evidence: `spacecraft evidence "<label>-functional" -- <full-test-suite>`. Auto-commit if the gate adds tests or fixes (`test:` / `fix:`).

3. **Review** - Self-review the diff. Then move to formal review for code review, design review (if UI), and release readiness before shipping.

## Rules

### Test quality (non-negotiable)

- **Must**: Test behavior through public interfaces only. Tests survive refactors because they don't couple to internal structure. A test that checks `loginAsGuest()` survives field changes; a test that checks `customer.type == 0` breaks on every refactor.
- **Must**: Expected values from independent source - literal, worked example, spec. Never recompute expected value the same way the code does.
- **Must not**: Test private methods, mock internal collaborators, verify through side channels (e.g., querying DB instead of using the interface).
- **Must**: Tests must be deterministic. No sleep-based waits, no random seeds without pinning, no order-dependent test state. A test that sometimes fails undermines the entire suite.
- **Must**: Every production bug gets a test. PBCNT (Percent of Bugs that Create New Tests) target: 100%. A bug that shipped once should never ship twice.
- **Must**: Read the test framework documentation before writing tests. Know parameterized tests, fixtures, setup/teardown, test categorization, and mock capabilities.

### Triage

- **Must**: Apply the triage gate before every cycle. If the test would be a trivial tautology → skip TDD, code directly, review.
- **Must**: Record skipped-TDD decisions per task. Don't silently skip - make the call explicit.
- **Must not**: Skip TDD for anything with branching logic, edge cases, error states, or non-obvious output.

### Loop discipline

- **Must**: Plan before red. No test without a confirmed seam and expected behavior written down.
- **Must**: Red before green. No production code without a failing test (or explicit triage skip).
- **Must**: One slice per cycle. One acceptance → one test → one implementation. Plan tasks are jigsaw pieces; acceptances are cycles inside them.
- **Must**: Refactor after all acceptance checks pass - not mid-cycle. You need the full picture.
- **Must** (AFK): Checkpoint-commit after RED, GREEN, and post-feature refactor (Commander; see sc-git).
- **Must**: Run functional test suite after refactor. Old tests must pass alongside new tests.
- **Must not**: Horizontal slice - bulk tests before bulk implementation.
- **Must**: After defect fixes, emit `TWINS:` - project-wide search for the same construct / twin occurrences before claiming done.
- **Must**: After **3 failed fix-verify cycles**, stop and hand back to human. Do not keep looping (3-cycle stop).

### Mocking

- **Must**: Mock at system boundaries only - external APIs, payment, email, time/randomness.
- **Must not**: Mock your own classes, internal collaborators, or anything under your control.
- **Prefer**: Dependency injection and SDK-style interfaces over generic fetchers. See `references/mocking.md`.

## Assessing test quality during review

> These are review-time heuristics - distinct from the authoring-time rules above.

Use these during code review to spot tests that pass green but don't actually verify behavior.

1. **Assertions on implementation detail**
   - Red flag: `expect(component.state.field).toBe(...)`, `expect(service.internalCounter).toBe(1)`, or any assertion on private fields, internal state, or non-public seams.
   - Fix: Rewrite assertions against public behavior or user-observable outcomes.

2. **Tautological assertions**
   - Red flag: `expect(result).toBe(a + b)` where `result = add(a, b)`; expected value computed the same way as the implementation.
   - Fix: Use an independent expected value - literal, worked example, or spec.

3. **Missed edge paths**
   - Red flag: Code contains `if`, `try/catch`, guards, or null/undefined handling, but tests only exercise the happy path.
   - Fix: Add tests for `else`, `catch`, empty input, null/undefined, and boundary values.

4. **Mock over-verification**
   - Red flag: Test mocks internal collaborators and verifies only mock interactions (`expect(mock.fn).toHaveBeenCalledWith(...)`); no assertion on actual output.
   - Fix: Assert on real behavior through public interfaces; mock only external system boundaries.

### Anti-patterns

| Type | Tell | Fix |
|------|------|-----|
| **Implementation-coupled** | Test breaks on refactor but behavior unchanged | Rewrite against public interface |
| **Tautological** | `expect(add(a,b)).toBe(a+b)` - can't disagree with code | Use a literal expected value |
| **Struct-constructor** | Creates struct, checks its own fields with zero transformation - `assert(x == x)` | Delete. Review struct type directly. No test needed. |
| **Horizontal slicing** | All tests written before any implementation | Restart with vertical slices |

## Out of scope

- Code quality principles and architecture - separate concern
- Review workflow - handled at review stage
- Evidence capture - handled by verification pipeline
- Git operations - handled by git infrastructure
- Mission lifecycle - handled by mission management
- UI testing - handled by design infrastructure

## Output format

```
Triage: TDD (non-trivial behavior) / Skip (trivial - coding directly)
Seams: [list of confirmed public boundaries]
Cycle 1/3: <acceptance check description>
  Plan: <seam + expected behavior>
  Red: test "<name>" - FAILS (expected)
  Green: minimal impl in <file:line> - PASSES
  Evidence: <label>

--- All checks pass ---
  Refactor: <what was improved>
  Functional test gate: <full suite> - PASSES
  Review: formal review
```

## Checklist

- [ ] Triage applied per acceptance check (TDD or skip - recorded)
- [ ] Plan done before red: seam confirmed, expected behavior written
- [ ] Every TDD test written before its implementation
- [ ] Tests verify behavior through public interfaces only
- [ ] No implementation-coupled or tautological tests
- [ ] Vertical slices only - one plan → test → impl per cycle
- [ ] Refactor done after all checks pass (not mid-cycle)
- [ ] Functional test suite passes after refactor
- [ ] Mocks only at system boundaries
- [ ] Evidence captured for each passing test suite (unit + functional)

### Edge cases

- **Test passes without implementation** - The test is not testing the right thing. Reject and re-write.
- **Test framework unfamiliar** - Use sc-search (WebSearch/WebFetch) for `"<framework> assertion API"` before writing tests. Wrong assertions produce false confidence.
- **No acceptance checks in plan.json** - Cannot verify against acceptance criteria. Ask for a plan first.

## References

- `references/examples.md` - good and bad test patterns with detection rules
- `references/mocking.md` - when to mock, dependency injection, SDK-style interfaces
- `references/testing-strategy.md` - testing pyramid, test types, AAA pattern, test doubles
