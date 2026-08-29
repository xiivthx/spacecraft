# How we instruct agents

Clarity over prompting tricks. Keep always-on rules short; put long detail in skills or `references/`.

## Layers

| Layer | Load | Use for |
|---|---|---|
| `.cursor/rules/` always-on | Every turn | Hard Never/Always, lanes, prompting |
| `.cursor/rules/` with globs | Matching files | Domain constraints |
| `.cursor/skills/*/SKILL.md` | On demand | Workflows (lifecycle `/sc-*` are Skills with `disable-model-invocation: true`, not `.cursor/commands/`) |
| `.cursor/agents/*.md` | Subagent Task | Job contracts: Goal / Inputs / Ban / Handshake |

Lifecycle detail: `.cursor/rules/200-workflow.mdc` + slash skills. `/sc-run` fixes findings before ready and reports them in the summary (no issue ledgers). Do not restate always-on rules inside agents.

## Spec Contract

Agents and skills state:

1. **Goal** - why / next use
2. **Output** - format and handshake
3. **Good vs Bad** - success bar
4. **Verify** - how correctness is checked

If unclear: research first; ask for preferences or unverifiable bars via `/sc-discuss`; never invent Verify. Mid `/sc-run`: soft → `decisions.md`; hard → stop and `/sc-discuss`. Ship needs machine-checkable Verify.

Mid `/sc-discuss` grill: if the human is stuck (confused wording, needs facts, or cannot picture state), unlock in place with natural language - re-pitch, research (`sc-search` then `sc-storm` when open-domain), or visualize (existing bake-off/draft, or a short chat state table). No new slash commands. Chat asks use Thai field bodies with English labels; `questions.md` / `decisions.md` stay English. Detail: `.cursor/skills/sc-clarify/SKILL.md`.

## Inner-loop gates (Quick + Mission)

Quick and Mission share always-on INTENT / AUTH / TWINS and the 3-cycle stop (see `.cursor/rules/000-spacecraft.mdc`, `200-workflow.mdc`). AUTH does not bypass `/sc-ship` or `SPACECRAFT_SHIP=1`. Before `ready`, Mission uses Cursor `bugbot` / `security-review` as primary defect/security surfaces (Spacecraft supplementary - `docs/mission-review.md`) and runs `sc-judge` (adversarial prove; ready only on `VERIFIED`). Ready proof is `evidence.jsonl` + empty `review.json` findings + `validate --strict` (mission artifacts and evidence; not-doc-drift / not-10X-validate) + judge hunts.

## Overnight `/sc-run`

After `/sc-discuss` clear, AFK `/sc-run` avoids mid-HIL except hard blocks. Stop on `3-cycle` | `timebox` | `blocked`. On those stops write `.space/missions/<id>/handback.md` with stop reason + remaining work cue. Detail: `.cursor/skills/sc-run/SKILL.md`. No overnight runner CLI.

## Mission canvas artifacts (`/sc-run`)

`/sc-run` may emit Cursor Canvas files under managed `canvases/` for human check after ready:

`~/.cursor/projects/<workspace>/canvases/<missionId>-<kind>.canvas.tsx` (`kind` ∈ `plan` | `findings` | `evidence`).

When emitted, record matching lines in `decisions.md`: `Canvas plan: `, `Canvas findings: ` (or `Canvas findings skipped: empty`), `Canvas evidence:` - each with an absolute path; chat and `decisions.md` include absolute markdown links. Canvas files and those decisions lines are optional aids for human check - not ready or `VERIFIED` proof. Do **not** put canvases under mission `.space/` or repo `.cursor/`. Do **not** replace mission brief or draft HTML / visual SoT with a canvas (brief stays Accept/Adjust/Reject chat HIL; draft stays HTML).

## Avoid

- Threats, tips, or career-stakes framing
- Expertise cosplay (including STORM lenses as personas instead of jobs)
- Forced chain-of-thought on reasoning models
- Tombstones after cuts - "formerly", "no longer", "removed", "deprecated in favor of", named absences. Rewrite survivors as the **current** product only (Cut hygiene in `.cursor/rules/000-spacecraft.mdc`).

Role names (Commander, Coder, Tester) are routing contracts, not expertise claims.

**Note:** Lens pass = five decision jobs that write `## Lens pass` or `Lens pass skipped:` in `decisions.md` - not expertise cosplay. See `.cursor/skills/sc-discuss/references/lens-pass.md` and sc-storm (Tier 3). Testability pass, Strategy pass, and RCRCRC pass are the same class of discuss decision jobs (`## Testability pass` / `## Strategy pass` / `## RCRCRC pass` or skips) - not QA personas. See `requirement-testability.md`, `htsm-strategy.md`, and `rcrcrc-impact.md`. Test Ideas use structured buckets (Positive / Negative / Edge / Overlooked) plus Implementation pitfalls (impl bugs, distinct from Requirement Bugs). SFDIPOT coverage review of **existing tests** vs requirement is on-demand via `sfdipot-coverage.md` - not an always-on discuss gate. On-demand test data design (`test-data-design.md`) produces variable-level rows (Positive/Negative/Boundary/Exploratory/Security) that map to Test Ideas buckets - not an always-on gate. On-demand oracle evaluation (`test-oracles.md`, FEW HICCUPPS) grounds problem judgment for observations - not expertise cosplay, not a tutor. Defect findings use `sc-run/references/defect-finding.md` (impact-first craft in `review.json` / run summary - not a ledger); domain security/performance/database scans map checklist priority to the same house severity before filing. On-demand code walkthrough via `sc-solid/references/code-walkthrough.md` - not expertise cosplay. On-demand prose craft (`sc-writer/references/prose-rhythm.md`, `narrative-context.md`) is the same class of decision job - rhythm rewrite and narrative context harvest, not expertise cosplay, not a tutor; source prompts that demonstrated threats, cosplay, and forced CoT are anti-patterns - extract craft only. On-demand prompt-refine (`sc-writer/references/prompt-refine.md`) is the same class of decision job - diagnose→rewrite for agent/skill/rule prompt fidelity to the Spec Contract, not expertise cosplay; extract craft only.
