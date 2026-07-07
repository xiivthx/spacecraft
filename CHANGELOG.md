# Changelog

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
