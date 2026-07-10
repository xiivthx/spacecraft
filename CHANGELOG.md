# Changelog

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
