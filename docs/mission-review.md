# Mission review gates

Spacecraft harness process for evidence, scope, tests, and acceptance review on every mission. This is not a product feature - it is how agents and humans decide when those dimensions are good enough for `ready` and ship.

## Five gates

1. **Deterministic first** - `spacecraft validate --strict`, evidence re-run, done-task evidence checks, and project read-only audits before LLM taste.
2. **Narrow per dimension** - Separate pass/fail questions (evidence, validate, scope, tests, acceptance, security, perf, unauthorized action) - not one "is this good?" blob.
3. **Pass/fail, fail-closed** - Verdicts are `pass`, `fail`, or `uncertain`. For `ready` and ship, **`uncertain` counts as fail** - critical finding or `REFUTED`, never soft-pass.
4. **Human calibration** - When rubric wording changes, sample missions and compare agent vs human labels; disagreements mean unclear criteria.
5. **Recheck on change** - Model, task shape, or criteria change ⇒ re-run gates; old passes do not carry forward without fresh evidence.

## Where it applies

| Phase | Use |
|-------|-----|
| `/sc-run` | Deterministic pre-review before `sc-reviewer`; validate, evidence, scope, tests |
| Review | `sc-reviewer` per-dimension verdicts (all missions) |
| Judge | `sc-judge` hunt for false completion, weakened tests, unauthorized action |

Missions with visual UI also use UX/UI review gates (sibling) - apply both when both scopes apply.

## Source of truth for agents

Full dimension table, machine vs critique split, deterministic ordering, and output format:

`.cursor/skills/sc-run/references/mission-review-gates.md`

Related: `sc-run`, `sc-verification`, `sc-judge`, `sc-tdd`, `sc-security` (heuristic only - no dynamic CVE tools), `sc-performance`, `ux-ui-review-gates.md` (visual sibling).

## Mission canvas milestones

`/sc-run` emits Cursor Canvas artifacts at plan, post-review findings-or-skip, and evidence (before `sc-judge` on the ready path). Live files live only under Cursor-managed `canvases/`:

`~/.cursor/projects/<workspace>/canvases/<missionId>-<kind>.canvas.tsx` where `kind` ∈ `plan` | `findings` | `evidence`.

Do **not** put `.canvas.tsx` under mission `.space/` or repo `.cursor/` (IDE detect is managed `canvases/` only).

Ready-path order after suite clean + deterministic pre-review: `sc-reviewer` → findings-or-skip canvas → evidence canvas → `sc-judge` → `set-state ready` on `VERIFIED`.

Greppable `decisions.md` lines (absolute paths; chat + `decisions.md` include absolute markdown links when a canvas exists):

| Milestone | Line |
|-----------|------|
| After plan | `Canvas plan: ` + absolute path to `<missionId>-plan.canvas.tsx` |
| After review (nonempty findings) | `Canvas findings: ` + absolute path to `<missionId>-findings.canvas.tsx` |
| After review (empty findings) | `Canvas findings skipped: empty` |
| Before `sc-judge` (ready path) | `Canvas evidence: ` + absolute path to `<missionId>-evidence.canvas.tsx` |

Ready is blocked when applicable lines are missing (chat-only link without the decisions line fails). Gate is file existence + decisions lines only - do not inspect canvas TSX/JSON shape. Mid-build short chat dumps are not required as canvases.

**Must not:** replace mission brief (Accept/Adjust/Reject chat HIL) with a canvas; use canvas as draft HTML / visual SoT.
