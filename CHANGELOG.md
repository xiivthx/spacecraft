# Changelog

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
