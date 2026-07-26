# Spacecraft

Spacecraft is a Cursor-native mission-control harness for AI-driven software development. It combines always-on project rules, specialized agents, reusable skills, safety hooks, optional MCP integration, and a local CLI for traceable mission state.

## What it provides

- Mission workflow from scope and planning through implementation, verification, review, and shipping
- Local mission artifacts and evidence under `.space/`
- Cursor-native rules, agents, skills, hooks, and MCP configuration under `.cursor/`
- Git safety with feature branches, Conventional Commits, and an explicit ship gate
- Specialized support for application development, testing, design, architecture, and embedded firmware

## Requirements

- Cursor
- Git
- Go 1.21 or newer when building the CLI from source
- macOS or Linux

## Installation

Spacecraft installs in two layers: a **User layer** (`make install-global`, once per machine) for agents, skills, MCP, the CLI, and global safety hooks - plus a generated `~/.cursor/spacecraft/USER-RULES.txt` you paste into Settings -> Rules -> User Rules so the `alwaysApply` rules apply in every workspace; and a **Project layer** (`./bootstrap.sh` or `make install-project`, once per repo) for the domain rules (`300`-`620`), agents, skills, and project hooks.

Install into a project with the bootstrap script:

```sh
./bootstrap.sh /path/to/project
```

When working from a clone of this repository, build and install with:

```sh
make install
```

For the User layer (once per machine):

```sh
make install-global
```

See the [installation guide](docs/installation.md) for setup and verification details.

## Quick start

Open the project in Cursor. User-facing slash skills are `/sc-discuss`, `/sc-run`, and `/sc-ship`:

```
/sc-discuss
/sc-run <roadmap-id>   # multi-mission
/sc-run                # mission-only when Sizing: single|phases (or map current)
/sc-ship
```

Flow:

1. `/sc-discuss` - size (`Sizing: single|phases|roadmap`), clarify, brainstorm, decide; for visual UI/FE approve draft HTML; clear clarify-status. Prefer a **new session** for run.
2. `/sc-run <roadmap-id>` AFKs incomplete roadmap missions to `ready` (or `/sc-run` mission-only for single/phases). Jigsaw plan → per-acceptance RED-GREEN via agents → combine/refactor → review. Visual UI/FE: requires draft already approved in discuss (not for `*-data` / `*-functional` seams); then screenshots/visual + functional recheck before ready. `/sc-ship` squashes AFK checkpoints to ≤5 commits before merge.
3. Human checks the ready work.
4. `/sc-ship` validates and closes out only when explicitly requested.

Roadmap selection helpers:

```sh
spacecraft map use <roadmap-id>   # set current roadmap
spacecraft map current            # print current roadmap id
spacecraft map next <roadmap-id>  # next incomplete mission on named roadmap
```

Skills live under `.cursor/skills/`. User-facing slash skills are `/sc-discuss`, `/sc-run`, `/sc-ship`, and `/sc-quick`. Spacecraft does not use `.cursor/commands/`.

## Cursor modes

Spacecraft lanes map to Cursor modes. Source of truth: `.cursor/rules/200-workflow.mdc`.

| User intent | Spacecraft lane | Cursor mode / action |
|---|---|---|
| Ask / clarify / brainstorm / visual draft | Discuss | Agent + `/sc-discuss` |
| Roadmap implement | Mission | Agent + `/sc-run` (after discuss clear) |
| Bug hunt | Debug | Cursor Debug Mode (no slash skill) |
| Ship | Ship | Agent + `/sc-ship` (hooks gate git) |
| Small edit / commit | Quick | Agent + `/sc-quick` (no mission; still INTENT/AUTH/TWINS/3-cycle) |

## Agents

Cursor discovers eight specialized agents in `.cursor/agents/`:

- `sc-coder` - implements production code
- `sc-tester` - writes tests and captures verification evidence
- `sc-planner` - converts mission specs into executable plans
- `sc-reviewer` - reviews changes, evidence, and release readiness
- `sc-designer` - reviews UI and visual design quality
- `sc-adviser` - advises on complex architecture and logic
- `sc-firmware` - implements STM32 firmware and embedded C
- `sc-writer` - writes and edits docs, prompts, messages, and other non-code prose

The always-on Spacecraft rules act as Commander and route work to these agents.

## CLI

Run the repository binary as `./spacecraft`, or use `spacecraft` after installation.

| Command | Purpose |
|---|---|
| `spacecraft init` | Initialize `.space/` mission state |
| `spacecraft new <title>` | Create a mission with a generated ID |
| `spacecraft missions` | List missions |
| `spacecraft use <number\|id\|title>` | Select the current mission |
| `spacecraft current` | Print the current mission ID |
| `spacecraft resolve [selector]` | Resolve a mission from a selector or branch |
| `spacecraft status` | Show mission status |
| `spacecraft flow` | Show the resolved mission workflow snapshot |
| `spacecraft bind-branch [selector]` | Bind the current branch to a mission |
| `spacecraft git-info` | Show Git worktree status |
| `spacecraft git-suggest [type] [slug]` | Suggest branch and commit conventions |
| `spacecraft set-state [mission-id] <new-state>` | Set mission state (mission-id optional when resolved from branch or current), alias: `state` |
| `spacecraft clarify-status <open\|clear\|deferred>` | Set clarification status |
| `spacecraft evidence [--mission <id>] <label> -- <command...>` | Run a command and capture evidence, alias: `evi` |
| `spacecraft validate [--strict] [mission-id]` | Validate mission artifacts and evidence, alias: `val`. `--strict` also requires `exitCode` on every evidence entry and evidence for each done plan task |
| `spacecraft closeout-check` | Check whether a mission is ready to close out, alias: `ship-check` |
| `spacecraft ship-check` | Alias for `closeout-check` |
| `spacecraft archive [selector]` | Archive a shipped mission |
| `spacecraft roadmap <new\|add\|rm\|ls\|show\|next\|archive\|use\|current> [...]` | Manage roadmaps, alias: `map`. `map use` / `map current` / `map next` support `/sc-run` |
| `spacecraft help` | Show live CLI help |

Mission states progress as:

```text
active -> planned -> in_progress -> ready -> shipped
```

`blocked` is available from any active state.

Use the CLI as the source of truth for current syntax:

```sh
./spacecraft help
```

## File layout

```text
.cursor/
  rules/                 always-on Commander, workflow, and domain rules
  agents/                eight specialized Cursor agents
  skills/                workflow and domain skills
  mcp.json               project MCP server configuration
  hooks.json             Cursor hook registration
  hooks/                 hook scripts
.space/
  missions/<id>/         spec, plan, decisions, evidence, and review artifacts
  archive/               shipped mission archives
  roadmaps/              multi-mission roadmaps
cmd/spacecraft/          Go CLI source and tests
spacecraft               repository CLI binary
bootstrap.sh             project bootstrap installer
Makefile                 build and install targets
```

## Mission artifacts

Each mission lives at `.space/missions/<id>/`. The primary files are:

- `mission.json` - identity, state, branches, and creation time
- `spec.md` - scope and intent
- `plan.json` - tasks, acceptance criteria, verification commands, and status
- `evidence.jsonl` - captured command evidence
- `decisions.md` and `questions.md` - resolved and open decisions
- `review.md` and `review.json` - formal review results

## Git and shipping

Mission work belongs on `feat/<mission-id>/<title>`, not directly on `main`. Immediately before `/sc-ship` merge, rename to `feat/<title>` (strip the mission id) so the merge commit uses the short name. Shipping is never inferred. `/sc-ship` runs only after an explicit request to merge or release, validates the mission, and applies the repository's release gates.

No-mission small edits use `/sc-quick` on branch `<type>/<title>` (no mission id). Ship with `SPACECRAFT_SHIP=1 SPACECRAFT_QUICK=1` so the hook skips `closeout-check`.

Before claiming mission build complete, prefer `spacecraft validate --strict`. Before mission merge, run `spacecraft closeout-check` (or `ship-check`). With `SPACECRAFT_SHIP=1` alone, the Cursor ship hook re-runs closeout before allowing `git merge` / `git push` / `git tag`. With both `SPACECRAFT_SHIP=1` and `SPACECRAFT_QUICK=1`, closeout is skipped (quick lane only).

Local gate (Go tests + hook unit tests):

```sh
make gate
```

On Cursor `sessionStart`, `.cursor/hooks/session-start.sh` prints `spacecraft status` (or `No active spacecraft mission.`) so the agent gets mission context.

## Lean profile

User-facing slash skills: **`/sc-discuss`**, **`/sc-run`**, **`/sc-ship`**, and **`/sc-quick`**.

- **HIL discuss:** `/sc-discuss` - clarify, decide, approve visual draft HTML
- **AFK run:** `/sc-run` loops `map next` until missions are `ready` or blocked; build is atomic RED-GREEN with auto checkpoint commits; UI missions require prior draft approval and recheck with visual + functional evidence
- **HIL ship:** final check + `/sc-ship`
- **Quick (no mission):** `/sc-quick` - manual edits/fixes/docs; branch → verify → commit → ship without mission artifacts or closeout
- **Active detail skills** under `.cursor/skills/` support agents (mission, planning, tdd, git, domains, …)
- **Explicit-only** (not auto-invoked): `sc-solid`, `sc-security`, `sc-performance`, `sc-ux-design`, `sc-diagram` - glob rules still apply

Project behavior and policy are defined by the always-on files in `.cursor/rules/`.

## How we instruct agents

Clarity over tricks. Agents use Goal, Output, Good vs Bad, Verify. If unclear, research then ask - never invent Verify. Details: [docs/prompting.md](docs/prompting.md). Artifact schemas: [docs/mission-artifacts.md](docs/mission-artifacts.md).
