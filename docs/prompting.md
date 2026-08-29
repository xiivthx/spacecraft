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

Quick and Mission share always-on INTENT / AUTH / TWINS and the 3-cycle stop (see `.cursor/rules/000-spacecraft.mdc`, `200-workflow.mdc`). AUTH does not bypass `/sc-ship` or `SPACECRAFT_SHIP=1`. Before `ready`, Mission uses Cursor `bugbot` / `security-review` as primary defect/security surfaces (Spacecraft supplementary - `docs/mission-review.md`) and runs `sc-judge` (adversarial prove; ready only on `VERIFIED`). Ready proof is `evidence.jsonl` + empty `review.json` findings + `validate --strict` (mission artifacts and evidence; not-doc-drift / not-10X-validate) + judge hunts. After `ready` and before `/sc-ship`, humans may run an optional post-ready drain via `sc-post-ready-drain` (git-primary; GitHub PR and Cursor `autopilot` are optional add-ons, not mandatory). After `ready` and before `/sc-ship`, humans may also run optional fat-diff hygiene via `sc-split-to-prs` (Cursor `split-to-prs`; human plan approval; quoted `AUTH:` before any outward push or PR publish; Must not replace `/sc-discuss` sizing or mid-run map resize; Must not merge; never required for ship; re-ready if the mission work branch mutates after `ready`). AUTH + `/sc-ship` remain merge/tag authority; capability → Cursor primary; orchestration / pass-fail → Spacecraft gate (same rules). Never soft-pass ready/ship from Cursor chat tables alone. Opened split PRs are never ready or ship proof.

## Overnight `/sc-run`

After `/sc-discuss` clear, AFK `/sc-run` avoids mid-HIL except hard blocks. Stop on `3-cycle` | `timebox` | `blocked`. On those stops write `.space/missions/<id>/handback.md` with stop reason + remaining work cue. Optional overnight/AFK watch via `sc-loop` / Cursor `/loop` (CI/jobs); on stop or ready handoff, disarm the loop (and write `handback.md` on stop). Optional Automations+Slack HIL AFK via `sc-automate-slack`: Slack notify on handback/needs-HIL; Slack reply may cue resume under `/sc-run` gates; Slack resume ≠ AUTH; never required for ready/ship; on stop or ready handoff, disarm Automations lane when armed. Detail: `.cursor/skills/sc-run/SKILL.md`, `.cursor/skills/sc-loop/SKILL.md`, `.cursor/skills/sc-automate-slack/SKILL.md`. No overnight runner CLI.

## Mission canvas artifacts (`/sc-run`)

Optional canvas-SoT via `sc-canvas-sot` for `/sc-run` plan | findings | evidence human-check emits. Never required for ready/ship or every `/sc-run`. Disposition: `Canvas-sot: ran` | `Canvas-sot: skipped: <reason>`. Detail: `.cursor/skills/sc-canvas-sot/SKILL.md` (upstream canvas by reference).

Live files under managed `canvases/` only:

`~/.cursor/projects/<workspace>/canvases/<missionId>-<kind>.canvas.tsx` (`kind` ∈ `plan` | `findings` | `evidence`).

When emitted, record matching lines in `decisions.md`: `Canvas plan: `, `Canvas findings: ` (or `Canvas findings skipped: empty`), `Canvas evidence:` - each with an absolute path; chat and `decisions.md` include absolute markdown links. Canvas files and those decisions lines are optional aids for human check - not ready or `VERIFIED` proof. Must not treat canvas as AUTH / `VERIFIED` / ready / ship authority. Do **not** put canvases under mission `.space/` or repo `.cursor/`. Do **not** replace mission brief or draft HTML / visual SoT with a canvas (brief stays Accept/Adjust/Reject chat HIL; draft stays HTML).

## Goals-mirror (optional)

Optional Goals-mirror via `sc-goal-roadmap` for multi-mission roadmaps (`Sizing: roadmap`). Disposition: `Goal-roadmap: ran` | `Goal-roadmap: skipped: <reason>`. Never required for ready/ship or every roadmap discuss/run. Must not treat Goals / Goal complete as AUTH / `VERIFIED` / ready / ship authority. Spacecraft `mission.json` + roadmap JSON remain SoT. Detail: `.cursor/skills/sc-goal-roadmap/SKILL.md`.

## Avoid

- Threats, tips, or career-stakes framing
- Expertise cosplay (including STORM lenses as personas instead of jobs)
- Forced chain-of-thought on reasoning models
- Tombstones after cuts - "formerly", "no longer", "removed", "deprecated in favor of", named absences. Rewrite survivors as the **current** product only (Cut hygiene in `.cursor/rules/000-spacecraft.mdc`).

Role names (Commander, Coder, Tester) are routing contracts, not expertise claims.

**Note:** Lens pass = five decision jobs that write `## Lens pass` or `Lens pass skipped:` in `decisions.md` - not expertise cosplay. See `.cursor/skills/sc-discuss/references/lens-pass.md` and sc-storm (Tier 3). Testability pass, Strategy pass, and RCRCRC pass are the same class of discuss decision jobs (`## Testability pass` / `## Strategy pass` / `## RCRCRC pass` or skips) - not QA personas. See `requirement-testability.md`, `htsm-strategy.md`, and `rcrcrc-impact.md`. Test Ideas use structured buckets (Positive / Negative / Edge / Overlooked) plus Implementation pitfalls (impl bugs, distinct from Requirement Bugs). SFDIPOT coverage review of **existing tests** vs requirement is on-demand via `sfdipot-coverage.md` - not an always-on discuss gate. On-demand test data design (`test-data-design.md`) produces variable-level rows (Positive/Negative/Boundary/Exploratory/Security) that map to Test Ideas buckets - not an always-on gate. On-demand oracle evaluation (`test-oracles.md`, FEW HICCUPPS) grounds problem judgment for observations - not expertise cosplay, not a tutor. Defect findings use `sc-run/references/defect-finding.md` (impact-first craft in `review.json` / run summary - not a ledger); domain security/performance/database scans map checklist priority to the same house severity before filing. On-demand code walkthrough via `sc-solid/references/code-walkthrough.md` - not expertise cosplay. On-demand prose craft (`sc-writer/references/prose-rhythm.md`, `narrative-context.md`) is the same class of decision job - rhythm rewrite and narrative context harvest, not expertise cosplay, not a tutor; source prompts that demonstrated threats, cosplay, and forced CoT are anti-patterns - extract craft only. On-demand prompt-refine (`sc-writer/references/prompt-refine.md`) is the same class of decision job - diagnose→rewrite for agent/skill/rule prompt fidelity to the Spec Contract, not expertise cosplay; extract craft only.
