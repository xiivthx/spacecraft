---
name: sc-tdd
description: "Test-driven development discipline. Activate on \"TDD\", \"Plan-Red-Green-Verify-Refactor-Review\", \"test-first\", or integration tests. Skip TDD for trivial code - just code and review."
---

# sc-tdd

## Goal

Design-contract → Plan → Red → Green → Verify → Repeat; then Refactor → Review. Under `/sc-run`, triage each acceptance; skip tautologies or run TDD. Outcome disposition SoT: `docs/mission-artifacts.md`.

## Workflow

1. **Triage** - Would a test be a trivial tautology? If yes → skip TDD.
   - **Skip:** getters/pass-throughs; config maps; boilerplate; struct-constructor tautology; docs/prose/wording-only. Record skip. Under `/sc-run`: Task(`sc-writer`) for docs/prose; Task(`sc-coder`) for other tautologies → evidence → checkpoint. No phrase-echo RED.
   - **TDD:** branching, edges, errors, non-obvious output, verify paths that can fail for real reasons.
2. **Before first product RED/GREEN (`/sc-run`)** - `Design-contract: complete` or skipped docs/prose; `Approved-scenarios:` freeze or skipped docs/prose (`docs/mission-artifacts.md`). Missing on behavioral work → stop; write via sc-planning. Frozen literals = RED oracle.
3. **Plan** - seam (public boundary), expected from Edge matrix / frozen scenario / spec, independent expected value. One condition per test.
4. **Red** - Task(`sc-tester`); one acceptance → one RED; deep assert via public interface. Passes without impl → rewrite. No shallow asserts.
5. **Green** - Task(`sc-coder`/`sc-firmware`/`sc-rtl`); minimum production code. **Coder Must not edit tests.**
6. **Verify** - package suite; `spacecraft evidence "<label>" -- <cmd>`. Repeat next acceptance.
7. **After checks** - Refactor; functional suite evidence; outcome A/B/C/PBT disposition via `docs/mission-artifacts.md` only; then Review.

## Principles (silent; surface violations)

- Public-interface behavior · critical paths not coverage religion · bug → test · deterministic
- Skip tautologies · design-contract before first AFK test · coder ≠ tests · composition over mocked seams

## Must / Must not

- **Must**: Triage every cycle; red before green when TDD; independent oracles; one condition per test; design-contract + approved-scenarios (or docs/prose skip) before product RED/GREEN; coder Must not edit tests; disposition before ready (`docs/mission-artifacts.md`).
- **Must not**: Phrase-echo RED for docs/prose; shallow asserts; mock internal collaborators; chase global 95-100% coverage; recompute expected like impl; invent skip-prefix strings here; bake developer home paths, usernames, or machine-local absolute paths into fixtures or oracles; edit production docs or code solely so a grep/string oracle passes - change the oracle or the spec instead.
- After defect fixes: `TWINS:` project-wide. After **3** failed fix-verify → human.

## Output

```
Design-contract: complete | skipped: docs/prose-only
Approved-scenarios: frozen-from-contract | frozen-by-human | skipped: docs/prose-only
Triage: TDD / Skip
Cycle: <acceptance> — Plan / Red / Green / Evidence
--- All checks pass ---
  Outcomes: see docs/mission-artifacts.md | Functional: PASSES
```

## Verify

Evidence labels present for RED/GREEN/functional; coder did not edit tests; design-contract + approved-scenarios complete or docs/prose skip before product cycles; outcome disposition recorded per `docs/mission-artifacts.md`.

## References

- `docs/mission-artifacts.md` - design-contract, approved-scenarios, outcome-gate SoT
- `references/examples.md`, `references/mocking.md`, `references/testing-strategy.md`
