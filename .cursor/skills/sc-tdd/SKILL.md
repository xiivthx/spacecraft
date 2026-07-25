---
name: sc-tdd
description: "Test-driven development discipline. Activate on \"TDD\", \"Plan-Red-Green-Verify-Refactor-Review\", \"test-first\", or integration tests. Skip TDD for trivial code - just code and review."
---

# sc-tdd

Plan → Red → Green → Verify → Repeat. Then Refactor → Review when the feature is complete. Plan what to test before writing. Skip TDD when the test would be a trivial tautology - just code and review. Every cycle produces one vertical slice.

Under `/sc-run`, the Commander triages each `plan.json` acceptance first. If TDD applies: Task(`sc-tester`) then Task(`sc-coder`), checkpoint after RED and GREEN. If triage skips: direct write - Task(`sc-writer`) for docs/prose/wording-only, Task(`sc-coder`) for other tautologies - evidence via the task verify command, one checkpoint (no RED harness). Then combine/refactor + functional gate before review.

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
| **Composition over mocked seams.** A create→use or multi-step user path needs at least one test that walks the real chain. Isolated unit tests with happy mocks will not catch credential, session, or wiring bugs between steps. | Mocked-seam false confidence |

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
- **Docs / prose / wording-only**: acceptance is "file must contain phrase X" (skills, agents, rules, markdown policy). The text *is* the spec - a RED harness that greps for the phrase you will write is a tautology. Direct write + `spacecraft evidence` with the task `verify` command (e.g. `rg`).

**Use TDD when:**
- Behavior has branching logic, edge cases, or error states
- The correct output is not trivially obvious from the input
- The implementation could plausibly be wrong
- A verify script / CLI / runtime path can fail for reasons other than "missing the sentence we plan to paste"

Record skipped-TDD decisions for a task in the task output, `plan.json` notes, or `decisions.md`. On skip under `/sc-run`: do **not** invent phrase-harness scripts; `sc-writer` writes docs/prose, `sc-coder` writes other tautology fixes; evidence runs task `verify`.

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
- **Must**: When the bug was a **composition** failure (step A’s output unusable by step B, or env/proxy/client wiring), the new test must exercise that chain — not only re-mock the failing seam.
- **Must**: Read the test framework documentation before writing tests. Know parameterized tests, fixtures, setup/teardown, test categorization, and mock capabilities.

### Triage

- **Must**: Apply the triage gate before every cycle. If the test would be a trivial tautology → skip TDD, code directly, evidence + review.
- **Must**: Record skipped-TDD decisions per task. Don't silently skip - make the call explicit.
- **Must not**: Invent RED harnesses for docs/prose/wording-only acceptances (phrase-echo greps).
- **Must not**: Skip TDD for anything with branching logic, edge cases, error states, or non-obvious output.

### Loop discipline

- **Must**: Plan before red. No test without a confirmed seam and expected behavior written down.
- **Must**: Red before green when TDD applies. No production code without a failing test **or** an explicit triage skip.
- **Must**: One slice per cycle. TDD path: one acceptance → one test → one implementation. Skip path: one acceptance → direct write → evidence. Plan tasks are jigsaw pieces; acceptances are cycles inside them.
- **Must**: Refactor after all acceptance checks pass - not mid-cycle. You need the full picture.
- **Must** (AFK): Checkpoint after RED and GREEN when TDD applies; after the direct-write+evidence step when triage skips; and after post-feature refactor (Commander; see sc-git).
- **Must**: Run functional test suite after refactor. Old tests must pass alongside new tests.
- **Must not**: Horizontal slice - bulk tests before bulk implementation.
- **Must**: After defect fixes, emit `TWINS:` - project-wide search for the same construct / twin occurrences before claiming done.
- **Must**: After **3 failed fix-verify cycles**, stop and hand back to human. Do not keep looping (3-cycle stop).

### Mocking

- **Must**: Mock at system boundaries only - external APIs, payment, email, time/randomness.
- **Must not**: Mock your own classes, internal collaborators, or anything under your control.
- **Must not**: Ship a UI or client path that only has a mocked seam (e.g. form spies on `startX`) without a companion **composition** contract or integration test for the real credentials/responses that path will receive.
- **Prefer**: Dependency injection and SDK-style interfaces over generic fetchers. See `references/mocking.md`.

### Composition paths (create → use, join → claim, auth → mutate)

These are non-negotiable for features that mint credentials, sessions, tokens, or IDs in one step and consume them in another:

1. **Backend/API contract must span the user path**, not only the create or mint response. Example shape: create resource → start session/join with returned credentials → complete the next step (claim, login, mutate) with those same credentials. Assert success on the final step.
2. **Do not leave inconsistent state** that the UI will treat as “ready” (e.g. marked claimed/active) without a usable secret or follow-up path — unless a composition test proves the UI path still works.
3. **Client wiring**: if the browser talks to the API via empty/public base URL, reverse proxy, or same-origin rewrite, keep at least one test that empty/default base builds same-origin `/api/...` (or documented relative) URLs, and keep local/dev config (Makefile, env templates) aligned. Cross-origin “could not reach server” failures are composition bugs.
4. **Before claiming done**: run the composition contract(s). Prefer also a smoke through the **app origin** (dev UI host), not only the API port.
5. **Production composition bugs**: fix + add the composition test in the same change. A mocked unit test alone does not close the bug.

See `references/testing-strategy.md` (Composition tests).

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

5. **Mocked seam without composition coverage**
   - Red flag: Feature ships with UI/unit tests that spy or mock the first call (`startJoin`, `createX`) while the real create→use / join→claim / auth→mutate chain has no API or integration contract.
   - Fix: Add a composition contract that uses credentials returned by step A in step B; do not close the bug with another happy mock.

### Anti-patterns

| Type | Tell | Fix |
|------|------|-----|
| **Implementation-coupled** | Test breaks on refactor but behavior unchanged | Rewrite against public interface |
| **Tautological** | `expect(add(a,b)).toBe(a+b)` - can't disagree with code | Use a literal expected value |
| **Struct-constructor** | Creates struct, checks its own fields with zero transformation - `assert(x == x)` | Delete. Review struct type directly. No test needed. |
| **Horizontal slicing** | All tests written before any implementation | Restart with vertical slices |
| **Mocked-seam only** | Create/join/auth UI covered only by mocked fetch; real credential handoff untested | Add composition contract across the real HTTP (or inject) chain |

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
- [ ] Docs/prose/wording-only acceptances skipped (no phrase-echo RED harness)
- [ ] Plan done before red when TDD applies: seam confirmed, expected behavior written
- [ ] Every TDD test written before its implementation
- [ ] Tests verify behavior through public interfaces only
- [ ] No implementation-coupled or tautological tests
- [ ] Vertical slices only - TDD: plan → test → impl; skip: write → evidence
- [ ] Refactor done after all checks pass (not mid-cycle)
- [ ] Functional test suite passes after refactor
- [ ] Mocks only at system boundaries
- [ ] Composition paths (create→use, join→claim, auth→mutate) have a chain contract, not only mocked seams
- [ ] Evidence captured for each passing acceptance (unit/verify + functional)

### Edge cases

- **Test passes without implementation** - The test is not testing the right thing. Reject and re-write.
- **Test framework unfamiliar** - Use sc-search (WebSearch/WebFetch) for `"<framework> assertion API"` before writing tests. Wrong assertions produce false confidence.
- **No acceptance checks in plan.json** - Cannot verify against acceptance criteria. Ask for a plan first.

## References

- `references/examples.md` - good and bad test patterns with detection rules
- `references/mocking.md` - when to mock, dependency injection, SDK-style interfaces
- `references/testing-strategy.md` - testing pyramid, test types, AAA pattern, test doubles
