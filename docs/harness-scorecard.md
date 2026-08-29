# Harness quality scorecard

Maintainer-facing health check for the Spacecraft **model × process** stack. It reports whether existing harness fixtures still pass - not base-model IQ, and not a public leaderboard.

v1 is a thin, fail-closed wrapper over gates the repo already runs. Green means the install path, false-completion traps, judge skill smokes, and process-grammar checks are intact on this checkout.

## Required dimensions

Each line below is a dimension id, its source of truth (SoT), and what pass/fail means. Dimension ids are frozen: `install-smoke`, `false-completion`, `judge-skill`, `process-grammar`.

| Dimension | SoT | Pass | Fail |
|-----------|-----|------|------|
| `install-smoke` | `scripts/smoke.sh` with `PROJECT` = repo root | Smoke exits 0 against this checkout | Smoke non-zero or not invoked with repo-root project dir |
| `false-completion` | `scripts/check-judge-break.sh` | Judge-break exits 0 (false-completion traps still catch) | Judge-break non-zero |
| `judge-skill` | all `.cursor/skills/sc-judge/test/*smoke*.sh` | Every matching smoke exits 0 | Any matching smoke non-zero or missing when expected |
| `process-grammar` | all `scripts/check-workflow-*.sh`, `scripts/check-sc-planning-*.sh`, `scripts/check-sc-planner-*.sh` | Every matching check exits 0 | Any matching check non-zero |

## Measure method

Thin runner: `scripts/harness-scorecard.sh` (Make target `test-harness-scorecard` wires it into `make test`).

- Prints greppable lines: `SCORECARD <dimension-id> <pass|fail>`
- Exit 0 iff **all** required dimensions pass
- Fail-closed: any required child fail, missing SoT, or forced fail ⇒ non-zero exit
- Wrap and report only; do not reimplement smoke, judge-break, or the check scripts

## Out of scope (v1)

Not required for scorecard green. Deferred or kept as separate Make/test legs:

- golden multi-trial / pass@k / pass^k
- LM process rubrics
- HarnessAudit mid-path packs
- sc-eval revive
- dashboard UI
- Antigravity as a *required* scorecard dimension (may stay a separate `make test` leg)

## Advisory fixtures

These may still run elsewhere in `make test`. They are **not** required for scorecard exit:

- config-smoke
- antigravity-smoke
- cli/test

## Related

- Runner: `scripts/harness-scorecard.sh`
- Install smoke: [installation.md](./installation.md)
- Mission review process: [mission-review.md](./mission-review.md)
- Antigravity (separate surface): [antigravity.md](./antigravity.md)
