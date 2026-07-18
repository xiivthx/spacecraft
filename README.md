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

Install into a project with the bootstrap script:

```sh
./bootstrap.sh /path/to/project
```

When working from a clone of this repository, build and install with:

```sh
make install
```

See the [installation guide](docs/installation.md) for setup and verification details.

## Quick start

Open the project in Cursor and run the workflow skills in chat:

```text
/sc-start
/sc-plan
/sc-build
/sc-ship
```

These are explicit Cursor skills stored under `.cursor/skills/`. Spacecraft does not use `.cursor/commands/`.

The core lifecycle is:

1. `/sc-start` creates a feature branch and mission artifacts.
2. `/sc-plan` turns the mission spec into a small, verifiable plan.
3. `/sc-build` delegates implementation and testing, then captures evidence.
4. `/sc-ship` validates and closes out the mission only when explicitly requested.

Additional workflow skills include `/sc-design`, `/sc-review`, `/sc-quick`, `/sc-research`, `/sc-resume`, and `/sc-debug`.

## Cursor modes

Spacecraft lanes map to Cursor modes. Source of truth: `.cursor/rules/200-workflow.mdc`.

| User intent | Spacecraft lane | Cursor mode / action |
|---|---|---|
| Ask / explain | Advisory | Ask Mode (or Agent with no writes) |
| Spec / design / plan | Mission (pre-build) | Plan Mode + `/sc-plan` |
| Implement | Mission (build) | Agent + Task(`sc-coder` / `sc-tester`) |
| Bug hunt | Debug | Debug Mode + `/sc-debug` |
| Formal review | Review | Agent Review / Task(`sc-reviewer`) + `/sc-review` |
| Ship | Ship | Agent + `/sc-ship` (hooks gate git) |
| Small edit / commit | Quick | Agent (no full mission gates) |

## Agents

Cursor discovers seven specialized agents in `.cursor/agents/`:

- `sc-coder` - implements production code
- `sc-tester` - writes tests and captures verification evidence
- `sc-planner` - converts mission specs into executable plans
- `sc-reviewer` - reviews changes, evidence, and release readiness
- `sc-designer` - reviews UI and visual design quality
- `sc-adviser` - advises on complex architecture and logic
- `sc-firmware` - implements STM32 firmware and embedded C

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
| `spacecraft roadmap <new\|add\|rm\|ls\|show\|next\|archive> [...]` | Manage roadmaps, alias: `map` |
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
  agents/                seven specialized Cursor agents
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

Spacecraft work belongs on `feat/<mission-id>/<title>`, not directly on `main`. Shipping is never inferred. `/sc-ship` runs only after an explicit request to merge or release, validates the mission, and applies the repository's release gates.

Before claiming build complete, prefer `spacecraft validate --strict`. Before merge, run `spacecraft closeout-check` (or `ship-check`). With `SPACECRAFT_SHIP=1`, the Cursor ship hook re-runs closeout before allowing `git merge` / `git push` / `git tag`.

Local gate (Go tests + hook unit tests):

```sh
make gate
```

On Cursor `sessionStart`, `.cursor/hooks/session-start.sh` prints `spacecraft status` (or `No active spacecraft mission.`) so the agent gets mission context.

Project behavior and policy are defined by the always-on files in `.cursor/rules/`.
