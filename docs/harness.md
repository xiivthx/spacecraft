# Harness scorecard

Maintainer check that install smoke, false-completion traps, judge smokes, and process-grammar scripts still pass on this checkout.

| Dimension | SoT | Pass |
|-----------|-----|------|
| `install-smoke` | `scripts/smoke.sh` with `PROJECT` = repo root | exit 0 |
| `false-completion` | `scripts/check-judge-break.sh` | exit 0 |
| `judge-skill` | `.cursor/skills/sc-judge/test/*smoke*.sh` | all exit 0 |
| `process-grammar` | `scripts/check-workflow-*.sh`, `check-sc-planning-*.sh`, `check-sc-planner-*.sh` | all exit 0 |

Runner: `scripts/harness-scorecard.sh` (`make test-harness-scorecard`). Prints `SCORECARD <id> <pass|fail>`; exit 0 iff all required dimensions pass.

Related: [installation.md](./installation.md), [review.md](./review.md)
