---
name: sc-tdd
description: "Test-driven development discipline. Activate on \"TDD\", \"Plan-Red-Green-Verify-Refactor-Review\", \"test-first\", or integration tests. Skip TDD for trivial code - just code and review."
---

# sc-tdd

Design-contract → Plan (per acceptance) → Red → Green → Verify → Repeat. Then Refactor → Review when the feature is complete. Under `/sc-run`, mission `design-contract.md` lands **before** first product RED/GREEN - see `docs/mission-artifacts.md`. Skip TDD when the test would be a trivial tautology - just code and review.

Under `/sc-run`, Commander triages each `plan.json` acceptance. If TDD applies: Task(`sc-tester`) then Task(`sc-coder`); one checkpoint per plan task after acceptances done. Both read `design-contract.md` + `approved-scenarios.md` when present; expected values from Edge matrix / frozen scenarios / spec literals - never recomputed like the implementation. Coder **Must not** edit tests. If triage skips: Task(`sc-writer`) for docs/prose/wording-only, Task(`sc-coder`) for other tautologies → evidence via task verify → one checkpoint. Then combine/refactor + functional gate before review.

**AFK note:** Seeing red is hygiene, not proof the test fails for the right reason. Prefer independent oracles (design-contract + approved-scenarios frozen literals; mutation disposition later). Deep asserts ground expected values in those oracles - not impl recomputation. Do not spend tokens on process theater beyond **one acceptance → one RED**.

## Principles

Apply silently; surface violations.

| Principle | Defends against |
|-----------|-----------------|
| Test behavior via public interfaces | Implementation-coupled tests |
| Test critical code, not everything | Coverage obsession / 95–100% global bars |
| Every production bug gets a test | Bugs that ship twice |
| Tests must be deterministic | Flaky suite trust collapse |
| TDD is a tool - skip tautologies | TDD as religion |
| Design-contract before first AFK test | Design-by-first-test |
| Coder does not own tests | Impl agent rewriting oracles |
| Composition over mocked seams | Happy-mock false confidence |

## Triage

Before each cycle: would a test be a trivial tautology? If yes → skip TDD.

**Skip:** getters/pass-throughs; pure config maps; framework boilerplate; struct-constructor `assert(x == x)`; `expect(add(1,2)).toBe(3)` where impl is the spec; **docs/prose/wording-only** (phrase-echo RED harnesses).

**Use TDD:** branching, edges, error states, non-obvious output, verify paths that can fail for reasons other than "missing the sentence we will paste".

Record skips in task output / `plan.json` notes / `decisions.md`. On skip under `/sc-run`: no invented phrase-harness scripts.

## Before first product cycle (`/sc-run`)

Confirm `Design-contract: complete` or `Design-contract skipped: docs/prose-only`, then `Approved-scenarios:` freeze footer or `Approved-scenarios skipped: docs/prose-only` (`docs/mission-artifacts.md`). Missing for behavioral work → stop; write via sc-planning. Frozen literals are the RED oracle.

## Per acceptance (one cycle)

1. **Plan** - seam (public boundary), expected behavior (Edge matrix / frozen scenario / spec), independent expected value, non-trivial triage. Write seam + expected down.
2. **Red** - **one acceptance → one RED**; **one condition per test**; deep assert via public interface (oracles = design-contract + approved-scenarios). If it passes without impl → rewrite. No bulk-suite-before-GREEN. No shallow asserts (`expect(true)`, typeof-only, status-only, phrase-echo).
3. **Green** - minimum production code. **Must not** edit tests.
4. **Verify** - suite for affected package; `spacecraft evidence "<label>" -- <cmd>`.
5. **Repeat** - next acceptance.

## After all checks pass

1. **Refactor** - extract, name, DRY (Rule of Three). Tests protect you.
2. **Functional gate** - full suite; `spacecraft evidence "<label>-functional" -- <full suite>`.
3. **Outcome gates (A/B/C + PBT)** - disposition required before ready (`docs/mission-artifacts.md`):
   - **A** - design-contract + tester/coder split + independent literals + coder cannot edit tests
   - **B** - frozen approved-scenarios; static-analysis evidence (**0 warning / 0 error** when tool runs) or `Static-analysis skipped/waived:…`; diff coverage **line and branch ≥90%** touched when measured, or skip/waive
   - **C** - mutation in scope when any of `Mutation: required` | pack `quality` | `Mutation: high-risk` (`docs/mission-artifacts.md`); then evidence (**>80%** scoped when tool; kill **boundary/operator** mutants) or `Mutation skipped:…` / waive; ordinary missions: `Mutation skipped: not in scope` valid
   - **PBT** - **100%** of design-contract **core-logic** modules (branching business rules / pure domain / state machines) need invariants + generators via project-existing lib (`fast-check` / Hypothesis / equivalent) and `pbt-…` evidence; else greppable `Pbt skipped: no project pbt tool` / `Pbt skipped: not core logic` / `Pbt waived: <reason>`. **Must not** invent PBT lib install mid-mission.
4. **Review** - then formal review / release readiness.

## Rules

- **Must**: Public-interface behavior tests; independent expected values; deterministic tests; triage before every cycle; record skips.
- **Must**: **One condition per test**; **one acceptance → one RED**. Deep asserts grounded in design-contract + approved-scenarios oracles (Edge matrix / frozen expected literals) - never recomputed like the implementation.
- **Must** (`/sc-run` product): design-contract + approved-scenarios complete or docs/prose skip before RED/GREEN; static / diff-cov / mutation / PBT disposition before ready.
- **Must** (PBT, product path): **100%** of design-contract **core-logic** modules carry property-based invariants + generators (`fast-check` / Hypothesis / equivalent already in the project) with `pbt-…` evidence, or a greppable `Pbt skipped:…` / `Pbt waived:…` line. Missing disposition ⇒ judge-REFUTE material. **Must not** invent PBT lib install mid-mission.
- **Must**: Red before green when TDD applies; one implementation slice per cycle; refactor after all checks; functional suite after refactor.
- **Must**: Coder Must not edit tests; oracle/scenario changes via Commander + `decisions.md` / `Scenario oracle change:`.
- **Must not**: Invent phrase-echo RED for docs/prose; mock internal collaborators; chase global 95–100% coverage; claim regression quality from red-green alone when mutation is in scope (`Mutation: required` | pack `quality` | `Mutation: high-risk`) with no disposition.
- **Must not**: Shallow asserts (`expect(true)`, typeof-only, status-only, phrase-echo). Error-path asserts **Must** use **ErrorCode** or error **instance** (typed equivalent OK) - see tester Ban / Deep assert.
- **Must**: Mock at system boundaries only. Composition paths (create→use, join→claim, auth→mutate) need a chain contract, not only mocked seams - see `references/testing-strategy.md`.
- **Must**: After defect fixes, `TWINS:` project-wide search. After **3** failed fix-verify cycles → human.

## Assessing test quality (review)

Red flags: assertions on private/internal state; tautological expected values; happy-path-only when branches exist; mock interaction-only asserts; mocked seam without composition coverage; shallow asserts; multi-condition tests. Fix: public behavior, independent literals from design-contract + approved-scenarios, one condition per test, edge paths, real chain contracts.

## Anti-patterns (short)

Implementation-coupled · tautological · struct-constructor · shallow assert · multi-condition · horizontal bulk impl · missing design-contract · coder edits tests · silent mutation omit · silent PBT omit on core-logic · mocked-seam only.

## Output format

```
Design-contract: complete | skipped: docs/prose-only
Approved-scenarios: frozen-from-contract | frozen-by-human | skipped: docs/prose-only
Triage: TDD / Skip
Cycle: <acceptance>
  Plan / Red / Green / Evidence
--- All checks pass ---
  Static | Diff-coverage | Mutation | PBT: evidence-* | skipped/waived
  Functional gate: PASSES
```

## References

- `docs/mission-artifacts.md` - design-contract, approved-scenarios, outcome-gate skip/waive SoT
- `references/examples.md`, `references/mocking.md`, `references/testing-strategy.md`
