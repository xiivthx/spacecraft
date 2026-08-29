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

Shared always-on INTENT / AUTH / TWINS and 3-cycle stop (`.cursor/rules/000-spacecraft.mdc`, `200-workflow.mdc`). AUTH does not bypass `/sc-ship` or `SPACECRAFT_SHIP=1`.

- Before `ready`: Cursor `bugbot` / `security-review` primary; Spacecraft supplementary (`docs/mission-review.md`); `sc-judge` adversarial prove - ready only on `VERIFIED`.
- Ready proof: `evidence.jsonl` + empty `review.json` findings + `validate --strict` + judge hunts.
- Optional post-ready (before `/sc-ship`): `sc-post-ready-drain` (git-primary; PR / autopilot optional); `sc-split-to-prs` (human plan approval; quoted `AUTH:` before outward push/PR publish; never merge; never required for ship; re-ready if work branch mutates after `ready`).
- Firewalls: never soft-pass ready/ship from Cursor chat alone; opened split PRs are never ready or ship proof. AUTH + `/sc-ship` remain merge/tag authority.

## Overnight `/sc-run`

After `/sc-discuss` clear, AFK `/sc-run` avoids mid-HIL except hard blocks. Stop on `3-cycle` | `timebox` | `blocked` → write `.space/missions/<id>/handback.md`.

- Optional watch: `sc-loop` (disarm on stop or ready handoff).
- Optional Slack HIL: `sc-automate-slack` (Slack notify / reply may cue resume under `/sc-run` gates; Slack resume ≠ AUTH; never required for ready/ship; disarm when armed).
- Detail: `sc-run`, `sc-loop`, `sc-automate-slack`. No overnight runner CLI.

## Mission canvas artifacts (`/sc-run`)

Optional via `sc-canvas-sot` for plan | findings | evidence human-check emits. Disposition: `Canvas-sot: ran` | `Canvas-sot: skipped: <reason>`. Never required for ready/ship.

Live path: `~/.cursor/projects/<workspace>/canvases/<missionId>-<kind>.canvas.tsx` (`kind` ∈ `plan` | `findings` | `evidence`). Record matching `decisions.md` lines with absolute paths when emitted.

Firewalls: canvas ≠ AUTH / `VERIFIED` / ready / ship. Do not put canvases under `.space/` or repo `.cursor/`. Do not replace mission brief or draft HTML with a canvas.

## Goals-mirror (optional)

Optional via `sc-goal-roadmap` for multi-mission roadmaps (`Sizing: roadmap`). Disposition: `Goal-roadmap: ran` | `Goal-roadmap: skipped: <reason>`. Never required for ready/ship. Goals / Goal complete ≠ AUTH / `VERIFIED` / ready / ship. SoT remains `mission.json` + roadmap JSON.

## Avoid

- Threats, tips, or career-stakes framing
- Expertise cosplay (including STORM lenses as personas instead of jobs)
- Forced chain-of-thought on reasoning models
- Tombstones after cuts - "formerly", "no longer", "removed", "deprecated in favor of", named absences. Rewrite survivors as the **current** product only (Cut hygiene in `.cursor/rules/000-spacecraft.mdc`).

Role names (Commander, Coder, Tester) are routing contracts, not expertise claims.

## On-demand decision jobs

Same class: decision jobs that write greppable `decisions.md` / craft artifacts - not expertise cosplay, not tutors, not always-on gates unless noted.

- **Lens pass** - five jobs → `## Lens pass` or `Lens pass skipped:` - `.cursor/skills/sc-discuss/references/lens-pass.md`; Tier 3 → `sc-storm`
- **Testability / Strategy / RCRCRC** - `## Testability pass` / `## Strategy pass` / `## RCRCRC pass` or skips - `requirement-testability.md`, `htsm-strategy.md`, `rcrcrc-impact.md`
- **Test Ideas** - Positive / Negative / Edge / Overlooked + Implementation pitfalls (impl bugs ≠ Requirement Bugs)
- **SFDIPOT** - coverage of existing tests vs requirement - `sfdipot-coverage.md` (on-demand)
- **Test data design** - Positive/Negative/Boundary/Exploratory/Security rows → Test Ideas buckets - `test-data-design.md` (on-demand)
- **Oracles** - FEW HICCUPPS judgment - `test-oracles.md` (on-demand)
- **Defect findings** - impact-first craft in `review.json` / run summary - `sc-run/references/defect-finding.md`
- **Code walkthrough** - `sc-solid/references/code-walkthrough.md` (on-demand)
- **Prose craft** - rhythm / narrative context - `sc-writer/references/prose-rhythm.md`, `narrative-context.md` (extract craft only; threats/cosplay/forced CoT are anti-patterns)
- **Prompt-refine** - diagnose→rewrite for Spec Contract fidelity - `sc-writer/references/prompt-refine.md`
