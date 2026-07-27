# Changelog

## 0.42.3 - 2026-07-27

- Chore: route high-volume agents to Cursor Models pool - `sc-writer` → Composer 2.5 Fast, `sc-coder` → Grok 4.5; keep tester/planner/gates on Other Models

## 0.42.2 - 2026-07-27

- Docs: encode clean commenting doctrine - why-not-what, best comment is none, provenance ban in style/clean-code/sc-coder Bad, plus trust-seed lesson
- Fix: ship gate closeout resolves workspace via hook payload (`workspace_roots` / cwd) so global hooks do not fail when process cwd is `~/.cursor`

## 0.42.1 - 2026-07-27

- Docs: encode optional roadmap `*-integrate` tip in mission-sizing (Must-when, skip, owns, must-not) with discuss/run/planning/workflow pointers and trust-seed lesson that per-tip combine is not whole-map reconcile

## 0.42.0 - 2026-07-27

- Feat: split Cursor install into User vs Project layers - `install-global` generates `~/.cursor/spacecraft/USER-RULES.txt` and merges absolute-path safety hooks; `install-project` installs domain rules `300`-`620` + `session-start` without re-copying `alwaysApply`
- Docs: document User Rules paste into Settings → Rules → User Rules and the User vs Project install split

## 0.41.14 - 2026-07-27

- Chore: ship sizing-prompt-cleanup closeout (findings #10-#15 verify evidence; plan-phaseN naming, sizing-outcome metaphor, discuss lane pointer, artifacts, trust-seed, `*-api` alias already released in 0.41.10)

## 0.41.13 - 2026-07-27

- Chore: ship sizing-run-handoff closeout (findings #1 #6 verify evidence; mission-only `/sc-run` for single|phases and non-ui draft skip already released in 0.41.10)

## 0.41.12 - 2026-07-26

- Chore: ship sizing-decision-tree closeout (findings #4 #7 #8 #9 verify evidence; always-size, phases ownership, resize protocol, hard ≥4 already released in 0.41.10)

## 0.41.11 - 2026-07-26

- Chore: ship sizing-map-protocol closeout (findings #2 #3 #5 verify evidence; map create protocol already released in 0.41.10)

## 0.41.10 - 2026-07-26

- Docs: discuss owns multi-mission map create/add; planning must hand off (no planning-owned roadmap creator)
- Docs: harden mission-sizing - always-size, ordered stub map protocol, resize protocol, mission-only `/sc-run`, non-ui draft skips, hard ≥4 single-seam split

## 0.41.9 - 2026-07-26

- Feat: approved draft HTML is visual source of truth - port chrome from draft; require scenario matrix (empty/error/few/many) before `UI draft approved`; draft-parity UI recheck and judge draft-drift hunt
- Feat: vertical multi-mission sizing playbook - `/sc-discuss` sizes single vs phases vs roadmap seams (`*-data` → `*-functional` → `*-ui`); ban `*-ux` seam and cross-feature layer waterfalls

## 0.41.8 - 2026-07-25

- Feat: new `sc-writer` agent for docs/prompts/messages/non-code prose; Commander, `sc-tdd`, `/sc-run`, `/sc-discuss`, `sc-coder`, and `sc-tester` route docs/prose work to `sc-writer` instead of `sc-coder`; README and installation list eight agents
- Chore: install script refactor - `bootstrap.sh` becomes an executable with an exit-trap fix; extract `scripts/install-binary.sh`; add `scripts/global-sync.sh` for agents+skills; add `scripts/test-install.sh`; `scripts/smoke.sh` checks `.space/roadmaps`; Makefile slimmed

## 0.41.7 - 2026-07-25

- Feat: ready/ship clean-slate bar - 0 open `issues.md` items and 0 review findings (including warnings); related/regression/consequence must be fixed in `/sc-run`; ship blocks instead of migrating open debt to GitHub; `closeout-check` enforces both gates
- Feat: `/sc-run` **issues drain loop** after plan+combine(+UI) - queue open issues, act by class+severity (must-fix related; file-or-fix unrelated), append new findings and loop until 0 open, then review/judge; review hits re-enter drain
- Docs: lifecycle prompt simplify - canonical discuss → run → ship in `200-workflow`; sc-learn owns issues policy; standardized handoffs; sc-judge caveats = non-defect only; lifecycle stays Skills not Commands
- Docs: `/sc-discuss` mission brief replaces comprehension quiz - Information/Question/Answer (Feynman + technical); human Accept/Adjust/Reject before clear (`references/mission-brief.md`)
- Docs: drop `docs/learned.md`; local source of trust is `.space/trust/` (gitignored with `.space/`); tracked seed in `sc-learn/references/trust-seed/`; `spacecraft init` seeds empty trust files; ship migrates only high-signal lessons/solved (no diary noise)
- Docs: sc-learn / sc-run / sc-ship / sc-judge / workflow rules updated for issue class triage and clean ready

## 0.41.6 - 2026-07-25

- Docs: sc-tdd composition-test rules (create→use / join→claim / auth→mutate) — mocked seams alone are insufficient; production composition bugs need chain contracts; mirror in testing-strategy, mocking, sc-web-backend and sc-web-frontend testing refs

## 0.41.5 - 2026-07-23

- Chore: retier Task agent models for plan/do/check diversity - planner Opus, doers Sonnet 5 high, reviewer GPT-5.6 Sol (cross-family check); adviser stays Opus; designer stays Sonnet

## 0.41.4 - 2026-07-23

- Fix: `spacecraft new --help`/`-h` and `spacecraft map new --help`/`-h` print usage and exit 0 instead of creating a mission/roadmap titled `--help`

## 0.41.3 - 2026-07-23

- Feat: after ship, recommend next roadmap mission - `spacecraft archive` prints next incomplete mission + `/sc-discuss` hint; `/sc-ship` and `/sc-quick` surface `map next` when a roadmap still has work

## 0.41.2 - 2026-07-23

- Chore: pin spacecraft Task agents to third-party models (Opus for reviewer/adviser; Sonnet for coder/tester/planner/designer/firmware) so Other Models usage is drawn instead of Auto inherit

## 0.41.1 - 2026-07-22

- Docs: Tier 3 browser matrix — `playwright-cli` primary, Cursor IDE browser fallback; remove system Chrome headless and browser-use/CDP from visual gates (`sc-ux-design`, `sc-run`, `sc-web-frontend`)
- Fix: ship/main hooks honor documented `SPACECRAFT_SHIP=1` / `SPACECRAFT_QUICK=1` command prefixes (Cursor hook process does not inherit agent shell env)

## 0.41.0 - 2026-07-22

- Feat: draft art-direction packs for `sc-ux-design` (`swiss-grid`, `editorial`) with shared prompt assembly, locked layout pools, and discuss pack selection (no silent auto-matcher)
- Feat: `sc-designer` pack-fidelity critique when a pack is selected; anti-slop catalog remains authoritative on conflict

## 0.40.2 - 2026-07-22

- Feat: always-on US English prompt coach (`.cursor/rules/025-english-coach.mdc`) - natural rewrite + short notes, or cheer when already correct
- Docs: `/sc-discuss` comprehension quiz before clear (customer-lens 5W1H on spec/result; Feynman keys; skip if nothing material; harness trivia banned; procedure owned by `references/comprehension-quiz.md`)

## 0.40.1 - 2026-07-22

- Docs: require Task(`sc-designer`) critique and Commander critical/important fixes before serving visual draft HTML to the human in `/sc-discuss`

## 0.40.0 - 2026-07-22

- Feat: restore `/sc-quick` as no-mission fast lane (branch → verify → commit → ship; no mission stub/closeout)
- Feat: ship hook skips `closeout-check` when `SPACECRAFT_SHIP=1` and `SPACECRAFT_QUICK=1`
- Feat: add `/sc-diagram` for standalone interactive HTML block diagrams (click-to-trace nets)
- Docs: wire Quick lane in README, workflow, and Commander rules; sc-git notes quick branch naming and closeout skip
- Docs: point sc-architect interactive wiring HTML to `sc-diagram`; list `sc-diagram` as explicit-only skill
- Test: hooks_test covers QUICK bypass and SHIP-required-with-QUICK

## 0.39.2 - 2026-07-22

- Docs: do not put mission ids in git commit subjects or bodies (merge commits and AFK checkpoints included); keep ids on the work branch / mission artifacts only

## 0.39.1 - 2026-07-21

- Docs: keep ≤7 jigsaw tasks per phase as a hard Must; phase (`plan-phaseN.json`) and roadmap (`spacecraft map`) as escape hatches; reject soft prefer ≤7 and 8-9 exception bands
- Docs: AFK TDD triage skips docs/prose wording-only acceptances (direct write + evidence; no phrase-echo RED harnesses)
- Test: policy check scripts for hard Must, escape hatches, and soft/8-9 rejection wording

## 0.39.0 - 2026-07-21

- Feat: always-on INTENT / AUTH / TWINS and 3-cycle stop for Quick and Mission agents
- Feat: `sc-judge` adversarial prove gate before `ready`; block on `REFUTED`
- Test: sc-judge false-completion and pass smoke fixtures
- Docs: prompting, mission-artifacts, README, and hooks_test assert AUTH does not bypass `SPACECRAFT_SHIP=1`

## 0.38.0 - 2026-07-21

- Feat: `/sc-discuss` pre-build HIL (clarify, decisions, visual draft approval); `/sc-run` is AFK implement only
- Feat: visual draft HTML gate moves from mid-`/sc-run` stop into `/sc-discuss`; run requires `UI draft approved` record
- Docs: workflow rules, README, installation, prompting, and related skills/agents updated for discuss → run → ship

## 0.37.3 - 2026-07-21

- Fix: prefer `.space/current` (`spacecraft use`) over branch `feat/<id>/` in lean CLI resolve/status/flow and evidence without `--mission` (#38 regression)
- Docs: sc-mission resolver priority is explicit selector / `--mission`, then `.space/current`, then branch

## 0.37.2 - 2026-07-21

- Fix: restore SHA-256 `outputHash` on evidence capture and validate integrity for the lean CLI (#37)
- Docs: document optional `outputHash` in mission-artifacts evidence schema

## 0.37.1 - 2026-07-21

- Docs: strip mission id from work branch name immediately before `/sc-ship` merge (`feat/<id>/<title>` → `feat/<title>`)
- Docs: sc-ship, sc-git, conventions, and README document work vs merge branch naming

## 0.37.0 - 2026-07-19

- Feat: UI HTML-first HIL gate - draft HTML (layout/style/components) must be approved before visual FE implementation under `/sc-run`
- Feat: after UI build, require screenshots/visual verify plus functional test evidence before `ready`
- Docs: sc-web-frontend, sc-ux-design, sc-designer, and workflow describe draft stop + visual recheck

## 0.36.0 - 2026-07-19

- Feat: `/sc-run` AFK build uses jigsaw plan tasks and per-acceptance RED-GREEN via Task agents with auto checkpoint commits
- Feat: `/sc-ship` squashes AFK checkpoints to ≤5 Conventional Commits before merge
- Docs: planner/coder/tester prompts and sc-planning/sc-tdd/sc-git describe jigsaw + cycle discipline

## 0.35.1 - 2026-07-18

- Fix: `make install-global` copies `sc-*` skills into `~/.cursor/skills` (including `/sc-run` and `/sc-ship`)
- Docs: installation guide covers global skills install

## 0.35.0 - 2026-07-18

- Feat: user-facing slash UX is `/sc-run` (AFK to ready) and `/sc-ship` only
- Feat: roadmap AFK helpers - `spacecraft map use` / `current` / `next`
- Feat: Cursor lane-to-mode bridge in workflow rules and README
- Feat: hard Cursor hooks - deny mutating git on main; deny merge/push/tag unless `SPACECRAFT_SHIP=1` and closeout-check passes
- Feat: `validate --strict` requires evidence `exitCode` and passing labels for every done plan task
- Feat: harden `closeout-check` review.json gates (status ready, no critical/blocksShip, releaseReadiness objects)
- Feat: add `sc-firmware` skill; polish agent proactive descriptions
- Fix: remove `research` and `check-deps` CLI stubs; `set-state`/`state` auto-resolve when mission-id omitted
- Fix: strip skills/rules/agents to lean CLI truth
- Docs: trim redundant prompt fluff; Goal/Output/Good/Bad/Verify without duplicated clarity gates
- Docs: README and mission-artifact docs for AFK/ship flow

## 0.34.1 - 2026-07-17

- Docs: replace the retired docs issue registry with GitHub Issues in sc-learn

## 0.34.0 - 2026-07-17

- Feat: complete Cursor-native spacecraft migration - rules, agents, skills, hooks, MCP
- Feat: restore mission CLI with aliases, evidence exit codes, and validation
- Feat: safe project/global install with MCP and hooks merge (no ~/.cursorrules clobber)
- Fix: installer no longer overwrites existing Cursor hooks (#50)
- Docs: rewrite README and installation guide for Cursor runtime
- Ci: restore make test and GitHub Actions verification

## 0.33.1 - 2026-07-17

- Feat: add schema versioning for Mission/Plan/Review artifacts
- Fix: surface hook errors instead of silently discarding them
- Fix: add subagent stall detection and retry guardrails
- Fix: add task dependency enforcement for subagent workflows
- Fix: SHA-256 integrity verification for evidence files
- Fix: replace review.json self-reported release gates with machine-readable checks
- Fix: decouple severity from ship-blocking gate
- Fix: enforce state machine transitions with prerequisite validation
- Fix: require automated verification in quick lane
- Fix: validate branch naming convention in closeout-check
- Fix: prioritize .space/current over branch signals in resolver
- Fix: align sc-planner model with sc-coder to prevent plan-code drift
- Fix: grant sc-reviewer bash access to verify evidence
- Fix: conditionally load DESIGN.md only for UI work
- Docs: remove explanatory filler from prompts

## 0.33.0 - 2026-07-16

- Feat: enrich roadmap missions with descriptions - add MissionEntry data model with ID and Description fields, backward compatible with old string array format

## 0.32.1 - 2026-07-16

- Docs: add parallel UI review rule for sc-designer + sc-llm-vision

## 0.32.0 - 2026-07-16

- Feat: suggest next roadmap mission after ship - print next unshipped mission with /sc-start hint

## 0.31.14 - 2026-07-15

- Docs: update sc-ship workflow - migrate issues to GitHub, not docs/issues.md
- Docs: combine knowledge migration + changelog + version bump into one commit

## 0.31.13 - 2026-07-15

- Refactor: move sc-debug from skill to command for explicit debug lane entry
- Docs: simplify sc-debug wording (remove jargon, British spellings)

## 0.31.12 - 2026-07-15

- Fix: archive readiness allows quick-fix missions without formal plan/review - evidence is sufficient proof
- Fix: archive summary accurately reflects which artifacts are kept

## 0.31.11 - 2026-07-15

- Fix: evidence capture runs from user CWD instead of spacecraft root (remove cmd.Dir override)
- Fix: require explicit set-state shipped and evidence capture for GitHub issue closing (#35)
- Docs: sc-git/sc-mission add post-merge checklist for issue closing verification

## 0.31.10 - 2026-07-15

- Feat: add criteria decomposition to sc-review (specification, evidence, quality, integrity)
- Feat: add optional Criterion field to Finding struct for traceability

## 0.31.9 - 2026-07-15

- Feat: auto-archive mission on set-state shipped
- Feat: auto-clean roadmap when all missions shipped

## 0.31.8 - 2026-07-14

- Fix: closeout-check now blocks merges without CHANGELOG.md update
- Fix: git-info warns when dirty changes exist on main branch
- Fix: git-info warns when CHANGELOG.md is missing from branch commits

## 0.31.7 - 2026-07-14

- Fix: clear stale .space/current in resolver (#31)
- Fix: add archive hint after set-state shipped (#28)

## 0.31.6 - 2026-07-14

- Fix: add dirty state check for non-main branches in workflow (#18)
- Fix: create bootstrap eval labelled examples (#7)

## 0.31.5 - 2026-07-14

- Chore: remove sc-map skill (redundant with codegraph) (M07Q5XKDZ)

## 0.31.4 - 2026-07-14

- Feat: auto-close GitHub issues with extended patterns, multi-source scanning, pre-check, and --no-close-issues flag (M07Q5XKB2)

## 0.31.3 - 2026-07-14

- Feat: add CI workflow with test & coverage gate (75% threshold) (M07Q5XHEI)

## 0.31.2 - 2026-07-14

- Feat: add AI review checklist to spec.md template (M07PFFNEB)

## 0.31.1 - 2026-07-14

- Fix: add archive-aware loadMission helper for roadmap commands to recognize shipped missions (M07Q480PR)
- Docs: add bug radar rule to auto-create GitHub issues for discovered problems

## 0.31.0 - 2026-07-14

- Feat: add sc-memory skill wrapping ctx_search and ctx_index with spacecraft conventions (M07PYRGLG)
- Feat: wire ctx_index hooks into sc-mission, sc-learn auto-indexing mission artifacts and lessons (M07PYRGLG)
- Feat: add ctx_search queries to sc-resume for prior mission context in handoff (M07PYRGLG)

## 0.30.2 - 2026-07-14

- Fix: accept string shorthand for ReleaseGate JSON unmarshaling to prevent archive blocking (M07PYOT0G)

## 0.30.1 - 2026-07-14

- Refactor: inject behavioral directives into all 7 sc-* agent prompts for tighter communication discipline (M07PXR32X)

## 0.30.0 - 2026-07-14

- Feat: add deploy hooks (before:deploy, after:deploy) and --ci flag for archive command (M07PFFNAV)
- Docs: enhance spec template with Edge Cases, Error Handling, Integration Points, Test Plan sections (M07PFFNAV)
- Docs: add CI/CD examples (GitHub Actions, GitLab CI, hook config) (M07PFFNAV)
- Cancelled: sandbox, model routing, AI review, MCP integration (over-engineering, no real value)

## 0.29.0 - 2026-07-14

- Feat: auto-close GitHub issues on ship when spec.md/decisions.md contains "fixes #N" or "closes #N" (M07PSM4N3)
- Refactor: sc-creator skill now supports creating skills, agents, and commands from templates (M07PSM4N3)
- Docs: evaluate sc-map vs codegraph - keep both as complementary tools (M07PSM4N3)
- Refactor: move templates from docs/templates/ to sc-creator skill directory (M07PSM4N3)

## 0.28.0 - 2026-07-14

- Fix: clear .space/current or set next roadmap mission on archive to prevent resolver conflicts (M07PFFJGY)

## 0.27.0 - 2026-07-14

- Docs: expand README.md with Overview, Installation, Quick Start, and 5 Usage sections (M07PFFIY3)
- Docs: add comprehensive installation guide (docs/installation.md) covering Homebrew, source build, and binary download
- Tooling: add dev and docs targets to Makefile for developer workflow
- Tooling: add developer automation scripts (setup.sh, test.sh, coverage.sh) with cross-platform support
- Docs: fix broken SPEC.md references in README.md

## 0.26.12 - 2026-07-14

- Test: increase coverage from 52.9% to 85.4% — util 98.6%, config 100%, mission 88.4%, eval 93.1%, resolver 89.9%, state 100%, gitutil 100%, hooks 90.8%, main.go 81.0% (M07PFFIGL)
- Test: add 425 tests across 16 packages, no test exceeds 5s
- Test: add CLI lifecycle tests (init, new, current, resolve, missions, use, bind-branch, status, set-state, clarify-status, evidence, exec, validate, closeout, archive)
- Test: add CLI utility tests (research, eval, check-deps, traces, cost, git-info, git-suggest, workflow, roadmap)
- Test: add eval package tests (runner, types, init, deterministic, rubric, lmjudge)
- Test: add mission package tests (store, model, archive, closeout)
- Test: add foundation package tests (util/fs, util/slug, config)
- Test: add critical package tests (resolver, state, gitutil, hooks)

## 0.26.11 - 2026-07-14

- Fix: remove duplicated MarshalJSON in eval/init.go - callers use single public function (M07PFFHYJ)
- Fix: eval scorers resolve file content before scoring - hallucination/response-quality checks read actual output, not path strings (M07PFFHYJ)
- Fix: filter eval-type entries from scoring - prevents feedback loop where eval results influence their own scores (M07PFFHYJ)

## 0.26.10 - 2026-07-14

- Feat: add issue tracking to roadmap CLI — issues stored in roadmap JSON, displayed grouped by phase with `roadmap show` (M07PFMSUD)

## 0.26.9 - 2026-07-14

- Feat: add roadmap feature — CLI commands for multi-mission long-term work (new, add, remove, show, list, continue, archive) with `.space/roadmaps/` JSON store and derived lifecycle states (M07PDUTCZ)
- Feat: add roadmap lifecycle validation — active/done/archived state transitions, duplicate detection, cross-roadmap conflict guard
- Feat: add large-scope auto-detect rule in PERSONA.md — commander suggests roadmap when task count exceeds 7 or scope is multi-phase

## 0.26.8 - 2026-07-14

- Feat: split sc-commander into commander + sc-adviser agent for architectural design escalation
- Feat: add sc-adviser read-only subagent with off-hours protocol and architect-tasks leave mechanism

## 0.26.7 - 2026-07-14

- Docs: add issues roadmap with 24 open issues across 6 phases, sorted by priority and dependency order

## 0.26.6 - 2026-07-14

- Fix: replace stale `.opencode/skills/` references with `.engine/skills/` (20 occurrences, 8 files)

## 0.26.5 - 2026-07-13

- Docs: add Standards section to AGENTS.md (em dash ban, CHANGELOG immutability, quality-first, E2E bug repro, pixel-perfect UI, lint/test excellence)
- Docs: add Bug fixes, Quality over cost, and Engineering excellence values to PERSONA.md; replace em dashes with plain dashes throughout

## 0.26.4 - 2026-07-13

- Fix: engine plugin self-locates via `import.meta.url` instead of plugin context, removes sync-plugin target and script

## 0.26.3 - 2026-07-13

- Refactor: engine plugin uses `directory` from plugin context instead of hardcoded paths (M07O1KTGQ)
- Feat: add `make sync-plugin` target and sync helper script

## 0.26.2 - 2026-07-13

- Docs: restructure PERSONA.md into 9 components (Identity, Values, Communication, Expertise, Boundaries, Workflow, Tool Usage, Memory Policy, Examples) — each rule lives in exactly one authoritative file
- Docs: deduplicate lane behavior, release rules, and research auto-trigger from AGENTS.md into PERSONA.md; AGENTS.md now cross-references PERSONA.md for those domains
- Docs: infuse Feynman (first-principles clarity), Zen (simplicity as discipline, non-attachment), and hacker ethos (source-first, trace the wire) into commander persona values
- Fix: enforce changelog and version bump as mandatory in sc-ship — never deferrable, hard gate
- Docs: trim prompt redundancy across AGENTS.md/PERSONA.md

## 0.26.1 - 2026-07-13

- Fix: harden sc-ship process against squash, force-add, and placeholder evidence (M07NXO1XE) — sc-commander constraints against squash-merge and `git add -f`; sc-ship fork-point confirmation before rebase; evidence command rejects empty-stdout placeholders and supports `--force` for stale entry cleanup; sc-git rebase-target mismatch detection and tag-only-after-merge policy; sc-verification cleanup guidance for failed evidence

## 0.26.0 - 2026-07-12

- Feat: add lifecycle hook system (M07N95XAN) — user-defined shell commands in `.space/hooks.json` that fire on mission events (created, state.changed, evidence.appended, validated, shipped, archived, wildcard `*`); supports blocking/non-blocking modes, configurable timeouts, stdout/stderr streaming with `[label]` prefix, and SIGINT forwarding

## 0.25.0 - 2026-07-12

- Feat: add traces and cost CLI subcommands for observability tracking — `spacecraft traces <id>` prints execution trace table with timestamp, latency, and tokens; `spacecraft cost --all` shows aggregate token usage and estimated cost per mission (M07N6P7I4)

## 0.24.0 - 2026-07-12

- Feat: add eval framework — `spacecraft eval <id>` runs structured evaluation (deterministic checks, rubric scoring, LM judge) against mission evidence, `eval init <id>` scaffolds eval directory (M07N361SC)
- Feat: add sc-eval skill — agent-facing instructions for creating labelled eval examples, writing rubrics, and running eval suites
- Feat: add eval coverage gate to sc-ship — blocks merge when coverage below configured threshold (default 0.8)
- Feat: add evidence.jsonl eval-type entries — eval results written back to mission evidence with `"type":"eval"` for downstream traceability
- Docs: fix sc-security SKILL.md trim, README/sc-review permission updates, docs/issues.md resolution from M07MTPHTR

## 0.23.0 - 2026-07-12

- Feat: add `make install` for global Spacecraft setup — symlinks CLI to `~/.local/bin/`, writes `~/.config/opencode/opencode.jsonc` with absolute paths
- Feat: add `make uninstall` to remove global symlink
- Refactor: move `opencode.json.example` to `docs/opencode.md` as config template with field reference
- Refactor: rename `SPACECRAFT.md` to `README.md` for repo discoverability
- Refactor: rename `INSTALL.md` to `docs/how-to-install.md`
- Docs: update file layout, skill paths, and slash command list to match actual `.engine/` structure

## 0.22.0 - 2026-07-12

- Feat: make .engine an opencode plugin — add `.opencode/plugins/engine.js` for auto-loading skills, commands, and persona context (M07N1L422)
- Feat: add `docs/INSTALL.md` with setup, troubleshooting, and git-backed install examples
- Refactor: consolidate engine structure under `.engine/` — move scripts, docs, config; remove 21 skill symlinks
- Refactor: simplify Makefile to 4 targets (build, test, clean, lint)
- Chore: clean `opencode.json` (remove commented-out models), remove old resolver tests

## 0.21.1 - 2026-07-12

- Chore: reorganize skills into category folders (core/quality/design/web/data/meta) with flat symlinks for OpenCode discovery
- Chore: move agents, commands, skills from .opencode/ to .engine/ — source of truth consolidated
- Chore: remove project-local .opencode/ symlink; skills served from global ~/.config/opencode/skills/ only

## 0.21.0 - 2026-07-12

- Feat: optimize agent models for task profiles — sc-reviewer `glm-5.2` → `deepseek-v4-pro` (reasoning), sc-planner `flash-free` → `qwen3.7-plus` (structured JSON), sc-tester `flash-free` → `deepseek-v4-flash` (Go Flash). Coder stays `kimi-k2.7-code`, commander stays `deepseek-v4-pro`, designer stays `glm-5.2`. All agents get reasoning effort variants (M07MTPHTR)
- Feat: add sc-security skill — static security review (OWASP top 10, hardcoded secrets, SQL/command injection, manifest scanning) for sc-reviewer read-only use
- Feat: add sc-performance skill — performance review (N+1 detection, memory leaks, bundle size, React re-render anti-patterns) for sc-reviewer read-only use
- Feat: expand sc-reviewer edge cases from 3 to 7 — add false-green tests, unaddressed prior findings (regression), huge diffs (>500 lines), conflicting evidence
- Feat: add research-request output pattern to sc-reviewer — reviewer emits "research needed: <query>" finding; commander executes spacecraft research
- Feat: add reviewer-facing sections to sc-tdd (false-green heuristics), sc-web-backend (Fastify lifecycle, validation, promises), sc-web-frontend (React hooks, useEffect, Tailwind, keys)
- Feat: register sc-security and sc-performance in AGENTS.md, SPACECRAFT.md routing tables, and sc-reviewer agent permissions
- Feat: add sc-llm-vision skill — LLM-driven visual design review via Gemini + agy CLI for screenshot-based UI quality inspection
- Docs: migrate M07MTPHTR knowledge — 6 solved issues, 4 deferred minors, 3 general lessons

## 0.20.0 - 2026-07-11

- Feat: add sc-search skill — auto-triggered 3-tier search escalation (google_search → webfetch → spacecraft research → ask user) for stuck issues, gray areas, and stale knowledge (M07KGTNR0)
- Feat: add /sc-research command — user-invoked systematic research via spacecraft research CLI (Brave Search, scoped docs, deep analysis)
- Docs: register sc-search in sc-commander.md auto-trigger skills and AGENTS.md skill table
- Docs: register /sc-research in AGENTS.md commands table
- Docs: update PERSONA.md Research auto-trigger section to reference sc-search skill

## 0.19.3 - 2026-07-10

- Fix: standardize task status on `done` — remove `completed` synonym from `TaskIsComplete()`, `taskIsOpen()`, and all test fixtures. Release gates drop `complete`/`completed` (keep `done`).
- Fix: add explicit `/sc-ship` merge gate — PERSONA.md, AGENTS.md, .opencode/AGENTS.md now block auto-merge. Quick lane ends with "report ready, wait for /sc-ship."

## 0.19.2 - 2026-07-10

- Fix: commander persona now delegates product code to sc-coder and tests to sc-tester — aligns persona with command-file delegation model. Added constraint preventing direct implementation.

## 0.19.1 - 2026-07-10

- Feat: `make install` now symlinks `scripts/spacecraft` binary to `~/.local/bin/spacecraft` — callable from any workspace. Symlink auto-updates on rebuild. `make uninstall` and `clean-global` also remove the binary symlink.

## 0.19.0 - 2026-07-10

- Feat: add sc-localize skill — culturally-aware bilingual copy review; catches literal translations that break in the target culture (M07K689IV)
- Docs: th-TH locale reference (collocations, register, UI conventions, anti-patterns), universal localization rules
- Docs: register sc-localize in SPACECRAFT.md auto-triggers and skill references table

## 0.18.0 - 2026-07-10

- Feat: add sc-pathfinder skill — chart a shared map of investigation tickets for work too large for one session; resolve tickets one at a time until the destination is clear (M07K4S68B)
- Feat: register sc-pathfinder in SPACECRAFT.md skill references table (explicit invocation only, not auto-triggered)

## 0.17.0 - 2026-07-10

- Feat: add sc-ux-design skill — anti-slop UI control (46 impeccable.style patterns), HTML draft preview with design briefs, animation/transition guidelines, and Playwright-based visual verification (M07JZ2OPD)
- Feat: register sc-ux-design in AGENTS.md skill table, agent permissions (sc-coder, sc-designer, sc-reviewer), and SPACECRAFT.md references
- Chore: add negation patterns to `.opencode/.gitignore` for skill-local package.json files (sc-ux-design Playwright dependency)

## 0.16.2 - 2026-07-10

- Feat: add 4 web developer skills — sc-web-frontend (React/TypeScript/Vite/Tailwind), sc-web-backend (Node.js/Fastify), sc-architect (ADR/C4), sc-database (PostgreSQL/migrations) (M07JSKJRB)
- Feat: register new skills in AGENTS.md skill table, agent permissions, command Use: lines, and SPACECRAFT.md routing tables
- Remove: sc-web-service skill — patterns absorbed by sc-web-backend
- Fix: stale Node test expectations for workflow recommendations

## 0.16.1 - 2026-07-10

- Fix: `make install` detects local `.opencode/` and warns about double-load collision (FORCE=1 to override)

## 0.16.0 - 2026-07-10

- Feat: `make install/uninstall` — symlink spacecraft skills/agents/commands to global OpenCode config
- Feat: generic `.opencode/AGENTS.md` for global consumption (no project-specific references)
- Feat: merge only agent definitions into global `opencode.json` (preserves plugins/providers)
- Feat: `make clean-global` — full reset target for removing all spacecraft symlinks
- Feat: `make fix-opencode` — strip project-specific keys from previously merged global config
- Chore: update agent model configs (sc-commander variant medium, sc-planner/sc-tester flash-free)

## 0.15.0 - 2026-07-10

- Feat: add 4 new built-in compact filters — go vet, npm test, docker ps, curl (M07J009BS)
- Feat: export filter config DSL (.space/compact/filters.json) — composable pipeline stages (include, exclude, dedup, truncate, stripPrefix) for user-defined compact rules without Go code (M07J009BS)
- Feat: integrate compact with evidence — `spacecraft evidence --compact` saves compacted output alongside raw (M07J009BS)
- Feat: compact EvidenceEntry model gains optional `compact` field for evidence JSONL entries

## 0.14.5 - 2026-07-09

- Feat: add `spacecraft compact` command — token-optimized output filtering for LLM context (M07IYSMYA)
- Feat: compact supports git status/diff/log, go test/build, ls, cat filters, plus generic dedup+truncation
- Feat: `--tee` flag saves full unfiltered output on non-zero exit for LLM fallback

## 0.14.4 - 2026-07-09

- Docs: fix 13 skill issues (sc-debug dedup, sc-design reading guide, sc-map schema extract, sc-mission expansion, sc-solid justification, sc-git redundancy, sc-tdd trigger, sc-creator trim) (S14–S26, M07IMLU48)
- Docs: add Hard Stop Gates and Research auto-trigger to all 9 command files (C10–C11, M07IMLU48)
- Docs: auto-trigger sc-git checks silently from sc-build; remove `/sc-git` command (C12, M07IMLU48)
- Docs: polish 3 agent files — brevity examples, locale notes, auto-trigger clarity (A7–A9, M07IMLU48)
- Docs: sync SPACECRAFT.md CLI table with main.go — add workflow alias (D6, M07IMLU48)
- Refactor: process improvements — Plan→Red→Green→Verify→Refactor→Review TDD cycle; triage gate for trivial tests; phase splitting for >7 tasks; struct-constructor test defense (T9, M07IMLU48)

## 0.14.3 - 2026-07-09

- Fix: unify archive.go ID normalization with `util.NormalizeMissionId` (G3, M07IG6R17)
- Feat: add `ConfigOption` functional options for path overrides in config package (G15, M07IG6R17)
- Refactor: remove 22 unused backward-compat type aliases from package main (G16, M07IG6R17)
- Refactor: consolidate `GitInfo` into `GitInfoData` in mission model — ResolveOutput now carries full git state (G17, M07IG6R17)
- Fix: rename `NoopRunner` → `DefaultRunner` to reflect actual behavior (G18, M07IG6R17)
- Refactor: replace manual flag parsing with `flag.FlagSet` in researchCmd and checkDepsCmd (G19, M07IG6R17)
- Docs: fix Commander session handoff to distinguish auto-trigger skills from user slash commands (D7, M07IG6R17)

## 0.14.2 - 2026-07-09

- Docs: clarify Solved vs Lessons distinction in `sc-learn/SKILL.md` and `docs/learned.md` — Lessons are general principles, Solved are specific issues

## 0.14.1 - 2026-07-09

- Fix: consolidate duplicated Go helper functions — `fileExists`, `readJSON`, `writeJSON` across archive/mission/state now use `util/fs.go` (G4, M07IDBF29)
- Fix: deduplicate `taskIsComplete` and `blockingFindings` between archive and closeout into mission package (G5, G6, M07IDBF29)
- Fix: correct `copyTextFile` to return error instead of silently failing with bool (G12, M07IDBF29)
- Fix: normalize `checkDepsCmd` and `resolveCmd` exit codes to Unix conventions (G8, G9, M07IDBF29)
- Fix: add fallback error log for `lookupPackage` nil return and stop `evidenceCmd` exit code propagation (G10, G11, M07IDBF29)
- Refactor: remove `goto` statements from `researchCmd` and rename `deep` variable collision (G7, G13, M07IDBF29)
- Fix: remove dead empty if-block in `workflow.go` (G14, M07IDBF29)
- Docs: add Pre-flight Checks and Hard Stop Gates to `sc-quick.md` (C7, M07IDBF29)
- Docs: replace fragile inline Python script in `sc-resume.md` and add missing sections (C8, C9, M07IDBF29)
- Docs: add D7 commander skill-as-command recommendation issue to `docs/issues.md`

## 0.14.0 - 2026-07-09

- Feat: add sc-tdd skill — test-driven development discipline with seam identification, anti-pattern detection, mocking guidelines (M07I7FY1P)
- Feat: add sc-solid skill — SOLID principles, clean code, complexity management, architecture (7 references) (M07I7FY1P)
- Feat: add sc-creator skill — 6-phase skill creation workflow (gather, create, wire, polish, register, verify) (M07I7FY1P)
- Feat: add sc-learn skill — mission knowledge capture and migration to global docs (issues, solved, learned) (M07I7FY1P)
- Feat: enhance all 6 agent files — Role/Context/Constraints sections, edge cases, skill awareness
- Feat: enhance 9 command files — Hard Stop Gates, Error Handling, edge cases, Research auto-triggers
- Feat: workflow command validation — verify Next command exists in .opencode/commands/ before proceeding
- Feat: system-wide quality review — 71 findings tracked (33 fixed) across skills, commands, agents, Go CLI
- Fix: nowISO() in research/deep.go returns real UTC timestamp instead of empty string
- Refactor: remove mutual cross-references and agent names from skill content
- Refactor: consolidate testing content into sc-tdd as single source of truth
- Docs: add SPACECRAFT.md CLI commands table, docs/issues.md, docs/learned.md
- Docs: add trigger phrases to all 12 skill descriptions, Research auto-triggers to 10/12 skills

## 0.13.0 - 2026-07-09

- Feat: add `/sc-quick` fast lane command — branch → commit freely → fast self-review → ship, skipping spec/plan/TDD/formal review (M07I88T8A)
- Feat: formalize 4 development lanes (Advisory, Mission, Debug, Quick) with commander auto-detection in AGENTS.md and PERSONA.md (M07I8KQYA)
- Feat: add fast self-review checklist to PERSONA.md — commander performs directly, no subagent required
- Docs: cross-reference AGENTS.md and PERSONA.md with "always read both" banners
- Chore: simplify agent bash permissions — wildcard allow with deny-list for dangerous ops

## 0.12.1 - 2026-07-09

- Chore: expand Commander shell command permissions via rtk proxy allow rules — git fetch, git merge-base, git rev-list, git rev-parse, git cat-file, git show, go test, nlm, python3 (M07HN2B5J)

## 0.12.0 - 2026-07-08

- Feat: add `spacecraft research <query>` subcommand — internet research via Brave Search API with config-driven search scopes and package registry lookups (Go, npm, PyPI, crates.io)
- Feat: add `spacecraft check-deps` subcommand — project-wide dependency freshness audit with concurrent registry lookups
- Feat: add `--deep` flag for AI-powered page analysis via browser-use (default) and notebooklm-mcp-cli
- Feat: add Commander auto-trigger to invoke research automatically in planning, implementation, debugging, and clarification lanes
- Feat: add configurable search scopes (.space/scopes.json) with built-in defaults for react, tailwindcss, nextjs, go, rust, postgresql, npm, pypi

## 0.11.0 - 2026-07-08

- Feat: add sc-debug skill with 5-step debug mantra (reproduce → trace fail path → falsify hypothesis → breadcrumb ledger → post-mortem), anti-rationalization guards, and cross-references to sc-verification, sc-git, sc-clarify
- Feat: add question-intent pattern to sc-reviewer agent — ask whether the change should exist before line-by-line review, consider simpler alternatives

## 0.10.0 - 2026-07-08

- Feat: merge sc-design, sc-polish, sc-design-review into single /sc-design with phase detection (design vs polish based on mission state)
- Feat: merge sc-work, sc-flow into /sc-build — single implementation loop with TDD cycle, self-review, checkpoint commits
- Feat: add /sc-resume command for session handoff pickup with live state injection via !`command` syntax
- Feat: add make status and make resume CLI convenience targets
- Refactor: remove phantom commands /sc-verify and /sc-status — never had command files
- Refactor: simplify state machine to 6 states (draft, planned, built, ready, shipped, blocked) — remove verifying and reviewing
- Refactor: add Kalama Sutta zero-trust self-audit gates to state transitions in /sc-build and /sc-review
- Command count: 12 → 8

## 0.9.0 - 2026-07-08

- Feat: add sc-map skill for project structure survey before planning — 3-phase survey (deterministic discovery + LLM semantic analysis) produces map.json with touchpoints, dependencies, risk zones, and layer classification
- Feat: sc-planning reads map.json as optional input for comprehensive task scoping
- Refactor: distribute SPEC.md content across AGENTS.md and PERSONA.md; lean AGENTS.md to project conventions

## 0.8.0 - 2026-07-08

- Feat: embed proactive rigor into process — selection decisions must enumerate ≥2 alternatives, self-audit before completion, evidence must prove functional correctness
- Feat: add evidence quality rules to sc-verification skill — distinguish config-echo (weak) from functional proof (strong)
- Close IS-01: root cause of AI laziness addressed via process gates

## 0.7.0 - 2026-07-08

- Feat: configure Spacecraft agents with OpenCode Go cost-aware model matrix — 4 distinct models across 6 agents (deepseek-v4-pro, deepseek-v4-flash, glm-5.2, kimi-k2.7-code)
- Feat: add per-agent colors and task permissions for sc-coder and sc-tester

## 0.6.6 - 2026-07-08

- Fix: accept `"done"` and `"cancelled"` task statuses in closeout-check and archive readiness, matching workflow.go `taskIsOpen` semantics
- Fix: sc-flow command delegates to sc-coder (implement) and sc-tester (test+verify) instead of spawning a nested sc-commander

## 0.6.5 - 2026-07-08

- Refactor: move sc-verify, sc-clarify, sc-status from commands to auto-triggered skills (12 → 9 commands)
- Feat: sc-commander agent auto-triggers sc-verification after every task, sc-clarify on ambiguity, sc-mission at session start
- Fix: harden changelog merge guard — changelog mandatory before merge, never deferrable in review/ship/git gates

## 0.6.4 - 2026-07-08

- Refactor: lean 13 → 12 command files — standardized frontmatter (`subtask: true` on sc-flow, sc-ship, sc-work), resolver preamble, and structured sections (Pre-flight/Workflow/Error handling) per `docs/templates/command.md`
- Refactor: merge sc-design-review.md into sc-review.md (design-review checks now part of review workflow)
- Refactor: de-duplicate git policy — sc-git.md sole source; sc-work, sc-ship, sc-review, sc-flow reference it

## 0.6.3 - 2026-07-08

- Docs: move skill and command templates from `.space/` to `docs/templates/` and remove `-template` suffix
- Docs: create `docs/templates/agent.md` from `.opencode/agents/` real-world patterns
- Fix: add changelog-before-merge guard and separate-commit pattern in sc-git workflow, commits rule, and review gate

## 0.6.2 - 2026-07-08

- Fix: add post-merge checklist in sc-git skill and guard against direct main edits

## 0.6.1 - 2026-07-08

- Refactor: realign all 7 SKILL.md files to `.space/skill-template.md` format (consistent frontmatter with name, description, license, compatibility, metadata.version, metadata.audience; template section layout: When to use, Workflow, Rules, Out of scope, Output format, Checklist, References).
- Disambiguate "design" → "visual design" in AGENTS.md, SPEC.md, sc-clarify, sc-mission, and sc-web-service skill files.
- Remove redundant AGENTS.md and DESIGN.md references from skill-level prompts (centralized in AGENTS.md).
- All original skill rules preserved; only structure and missing sections changed.

## 0.6.0 - 2026-07-08

- Add `sc-coder`/`sc-tester` to commander `task.permission` block and `sc-verification` to `sc-work.md` `Use:` frontmatter for `/sc-work` and `/sc-flow` TDD loops.
- Create centralized routing table at `docs/routing.md` documenting command → agent → subagent → skill → permission.
- Add `sc-web-service` to coder `skill.permission` block and open bash CLI permissions for tester evidence capture and reviewer resolve/status/validate.
- All 5 agent-architecture-review issues resolved. Verified with grep-based acceptance checks.

## 0.5.3 - 2026-07-07

- Refactor `scripts/src` from flat `package main` monolith (12 files, 2900 lines) into 10 internal/ Go packages with interfaces, dependency injection, and error returns.
- Split packages: `config`, `id`, `util`, `gitutil`, `mission/model`, `mission/store`, `resolver`, `state`, `workflow`, `closeout`, `archive`.
- All business logic returns errors; single `os.Exit(1)` in `main()`.
- ~140 Go unit tests + 23 Node integration tests pass. No behavior changes.

## 0.5.2 - 2026-07-07

- Migrate package.json scripts to GNU Makefile with `make test`, `make build`, `make help`, and `make sc-*` targets.
- Update docs and agent permissions to reference `make` instead of `npm run`.
- Remove package.json (Node.js `type:module` setting was unused).

## 0.5.1 - 2026-07-07

- Rewrite `scripts/spacecraft.mjs` (Node.js) as a zero-dependency Go binary (`scripts/spacecraft`) with sub-5ms cold-start.
- All 18 subcommands behave identically; JSON output (`--json` flags) is byte-compatible.
- Invocation changes from `node scripts/spacecraft.mjs <cmd>` to `scripts/spacecraft <cmd>`.
- Update all docs, skills, and command files to reference the Go binary.
- Update branch naming convention to `<type>/<id>/<title>`.
- Refine Conventional Commits guidance: no scopes by default, lowercase bullet-point bodies.

## 0.5.0 - 2026-07-07

- Add `/sc-flow` workflow runner guidance and `flow` helper status for repeated work, verify, checkpoint loops.
- Add compact sortable mission and evidence ids without hyphens while preserving legacy mission id resolution.
- Add shipped mission archive compaction with `node scripts/spacecraft.mjs archive`.
- Add English-only root `SPEC.md` for the project-level Spacecraft contract.

## 0.4.0 - 2026-07-07

- Add active mission resolver docs and prompts for multi-mission sessions, including title/number selection.

## 0.3.0 - 2026-07-07

- Add session handoff versus release closeout policy so chat handoff does not trigger merges.
- Add automatic mission and branch start guidance for clear mutating work.
- Move commander persona guidance to `PERSONA.md`.
- Require fresh official dependency/version checks before direct dependency and current API work.
- Add rtk shell-output proxy policy and deny passthrough for push, sudo, and destructive commands.
- Add `sc:closeout` and harden `closeout-check` release gates, evidence validation, and evidence id reservation.

## 0.2.0 - 2026-07-07

- Add a bounded lightweight self-review/self-test loop to `/sc-work`.
- Treat required read-only subagents as part of their slash command contract for `/sc-plan`, `/sc-design`, `/sc-design-review`, `/sc-polish`, and `/sc-review`.
- Document the subagent invocation model and review workflow in Spacecraft docs and agent rules.
