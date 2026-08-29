# Judge-break fixtures

Known-bad mission packs that **must** be rejected by deterministic closeout predicates.

Run: `scripts/check-judge-break.sh [repo-root] [spacecraft-binary]`

These are not LLM trap-eval suites. They prove the exit gate fails on bad disk state.

## Packs

| Pack | Leak class | Typical `mustContain` |
| --- | --- | --- |
| `empty-evidence/` | No evidence captured | `no evidence captured` |
| `false-completion/` | Done plan + empty evidence | `no evidence captured` |
| `review-findings/` | Non-empty review findings (Cursor-sourced leftover with `source: bugbot` must fail closeout) | `review finding` |
| `false-consensus/` | VERIFIED/ready without per-finding `AGREE\|DISAGREE_*` dissent | `false-consensus` |
| `charitable-reviewer/` | Free-text `builderRationale` (not structured-lines-only) | `charitable-reviewer` |
| `silent-mutation-skip/` | `Mutation: required\|high-risk` without disposition or `mutation-` evidence | `silent-mutation-skip` |
| `retroactive-oracle-change/` | Frozen scenarios thawed/edited without `Scenario oracle change:` | `retroactive-oracle-change` |
| `freeze-postdate/` | Test-run evidence before freeze event (retroactive freeze) | `postdated-freeze` |
| `freeze-drift-without-oracle-line/` | Frozen file hash drift without `Scenario oracle change:` | `freeze-drift` |
| `silent-cross-model-critic/` | Gates ≥ M9G7IHV3 without `Cross-model critic:` or skip line | `silent-cross-model-critic` |
