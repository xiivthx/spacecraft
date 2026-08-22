# Spacecraft

Spacecraft is a Cursor-native mission-control harness for AI-driven software development. It combines always-on project rules, specialized agents, reusable skills, safety hooks, optional MCP integration, and a local CLI for traceable mission state.

## What it provides

- Mission workflow from scope and planning through implementation, verification, review, and shipping
- Local mission artifacts and evidence under `.space/` (gitignored; first-use ensure may `git init`, write starter `.gitignore` from `templates/gitignore`, and soft-run `codegraph init` when no `.codegraph/` index exists - warn and continue on missing binary or failure)
- Cursor-native rules, agents, skills, hooks, and MCP configuration under `.cursor/`
- Git safety with feature branches, Conventional Commits, and an explicit ship gate
- Specialized support for application development, testing, design, architecture, and embedded firmware

## Requirements

- Cursor
- Git
- Node.js 18 or newer for the CLI (`cli/spacecraft.mjs`) and companion installs (`make install-machine`)
- macOS or Linux

## Installation

Spacecraft installs in two layers: a **User layer** (once per machine) for agents, lean-core skills, MCP, the CLI, and global safety hooks - plus a short `~/.cursor/spacecraft/USER-RULES.txt` CORE (`010-hard-contract`). Paste into Settings -> Rules -> User Rules after regen. And a **Project layer** (`./bootstrap.sh`, `make install-project`, or `spacecraft setup`) for alwaysApply hard-contract, **pack-selected** domain skills/rules, `session-start`, and **safety hooks** (secrets / destructive / main-write / ship+push-ask) so cloud agents get the same gates. Lean-core lifecycle skills and agents stay User layer (`~/.cursor`).

**Enforcement map:** hooks = hard; `010-hard-contract` = always-on soft; skills/glob rules = on demand. Markdown alone does not block `.env` reads or push.

**New PC** - User layer plus companion tools (caveman, rtk, codegraph) with Cursor wiring:

```sh
git clone https://github.com/xiivthx/spacecraft.git
cd spacecraft
make install-machine
```

User layer only (default **lean**: lifecycle + process skills):

```sh
make install-global
```

Lean reconcile prunes spacecraft-managed domain encyclopedia skills under `~/.cursor/skills` that sit outside the lean allowlist; unrelated files under `~/.cursor` stay put. Opt in to domain encyclopedias with `--full` via `SPACECRAFT_SKILL_PROFILE=full` or `make install-global FULL=1`:

```sh
SPACECRAFT_SKILL_PROFILE=full make install-global
# or: make install-global FULL=1
```

Project layer (per repo) installs **selected** domain packs locally - no User `--full` required. Choose packs with `spacecraft setup` (interactive default: **quality**; coming packs such as `iot`/`fpga`/`pcb`/`management` are listed but not installable). Non-TTY needs `--packs` or `SPACECRAFT_PACKS`, or fails. Existing `.cursor/spacecraft-profile.json` → silent reconcile; change packs with `--reconfigure`. User-layer lean/full is unchanged; lean-core stays out of the project layer.

```sh
./bootstrap.sh /path/to/project
# or: spacecraft setup --packs frontend,quality
```

Bootstrap / `install-cursor` first `.space` create also ensures git (init if needed), starter `.gitignore` with `.space/`, and may soft-run `codegraph init` when no `.codegraph/` index exists (warn and continue on missing binary or failure).

When working from a clone of this repository, build and install with:

```sh
make install
```

See the [installation guide](docs/installation.md) (project pack setup, lean vs full) and [Antigravity guide](docs/antigravity.md) for companions, Tools status output, and verification details.

### Antigravity Installation

Spacecraft fully supports **Google Antigravity** via its plugin and project scaffold system:

```sh
# Global Antigravity plugin (~/.gemini/config/plugins/spacecraft)
make install-antigravity

# Project layer (.agents/ + GEMINI.md)
./bootstrap.sh --antigravity /path/to/project
# or: make install-antigravity-project PROJECT=/path/to/project
```

## Quick start

Open the project in Cursor. User-facing slash skills are `/sc-discuss`, `/sc-run`, and `/sc-ship`:

```
/sc-discuss
/sc-run <roadmap-id>   # multi-mission
/sc-run                # mission-only when Sizing: single|phases (or map current)
/sc-ship
```

Flow:

1. `/sc-discuss` - size (`Sizing: single|phases|roadmap`), lens pass or skip (`## Lens pass` / `Lens pass skipped:`), testability pass or skip (`## Testability pass` / `Testability pass skipped:`), strategy pass or skip (`## Strategy pass` / `Strategy pass skipped:`), RCRCRC when two requirement versions (`## RCRCRC pass` / skip), clarify, brainstorm, decide; for visual UI/FE approve draft HTML; clear clarify-status. Prefer a **new session** for run.
2. `/sc-run <roadmap-id>` AFKs incomplete roadmap missions to `ready` (or `/sc-run` mission-only for single/phases). Jigsaw plan → per-acceptance RED-GREEN via agents → combine/refactor → review. Visual UI/FE: requires draft already approved in discuss (not for `*-data` / `*-functional` seams); then live product review on the running app URL (Tier 3 + designer **live-product**) plus draft-parity (paired draft-surface vs live screenshots, side-by-side) and functional recheck before ready. `/sc-ship` squashes AFK checkpoints to ≤5 commits before merge.
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

Cursor discovers specialized agents from the User layer (`~/.cursor/agents/` after `make install-global` / `make install-machine`):

- `sc-coder` - implements production code
- `sc-tester` - writes tests and captures verification evidence
- `sc-planner` - converts mission specs into executable plans
- `sc-reviewer` - reviews changes, evidence, and release readiness
- `sc-designer` - reviews UI and visual design quality
- `sc-adviser` - advises on complex architecture and logic
- `sc-firmware` - implements STM32 firmware and embedded C
- `sc-writer` - writes and edits docs, prompts, messages, and other non-code prose
- `sc-browser-probe` - live browser sweep + AFK find→fix→re-probe until CLEAN

The always-on Spacecraft rules act as Commander and route work to these agents.

## CLI

The CLI is Node (`cli/spacecraft.mjs`). Run the checkout link as `./spacecraft`, or use `spacecraft` after installation.

| Command | Purpose |
|---|---|
| `spacecraft init` | Initialize `.space/` (git init if needed; starter `.gitignore` with `.space/` on first create; may soft-run `codegraph init` when no `.codegraph/` index) |
| `spacecraft setup [--packs a,b] [--reconfigure]` | Project pack selection; writes `.cursor/spacecraft-profile.json`; selective install + prune (`SPACECRAFT_PACKS` when flag absent) |
| `spacecraft new <title>` | Create a mission with a generated ID |
| `spacecraft missions` | List missions |
| `spacecraft use <number\|id\|title>` | Select the current mission |
| `spacecraft current` | Print the current mission ID |
| `spacecraft resolve [selector]` | Resolve a mission from a selector or branch |
| `spacecraft status` | Show mission status |
| `spacecraft flow` | Show the resolved mission workflow snapshot |
| `spacecraft bind-branch [selector]` | Bind the current branch to a mission |
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
  agents/                specialized Cursor agents (coder, tester, planner, reviewer, designer, adviser, firmware, writer, browser-probe)
  skills/                workflow and domain skills
  mcp.json               project MCP server configuration
  hooks.json             Cursor hook registration
  hooks/                 hook scripts
.space/
  missions/<id>/         spec, plan, decisions, evidence, and review artifacts
  archive/               shipped mission archives
  roadmaps/              multi-mission roadmaps
cli/                     Node CLI entry (spacecraft.mjs) and tests
spacecraft               repository CLI link (cli/spacecraft.mjs)
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

No-mission small edits use `/sc-quick` on branch `<type>/<title>` (no mission id). One `/sc-quick` run goes branch → verify → commit → local merge/tag; use `SPACECRAFT_SHIP=1 SPACECRAFT_QUICK=1` so the hook skips `closeout-check`. Push still needs an explicit ask.

Before claiming mission build complete, prefer `spacecraft validate --strict`. Before mission merge, run `spacecraft closeout-check` (or `ship-check`). With `SPACECRAFT_SHIP=1` alone, the Cursor ship hook re-runs closeout before allowing `git merge` / `git push` / `git tag`. With both `SPACECRAFT_SHIP=1` and `SPACECRAFT_QUICK=1`, closeout is skipped (quick lane only).

Local gate (Node CLI tests + hook unit tests):

```sh
make gate
```

On Cursor `sessionStart`, `.cursor/hooks/session-start.sh` prints `spacecraft status` (or `No active spacecraft mission.`) so the agent gets mission context.

## Lean profile

User-facing slash skills: **`/sc-discuss`**, **`/sc-run`**, **`/sc-ship`**, and **`/sc-quick`**.

- **HIL discuss:** `/sc-discuss` - clarify, decide, approve visual draft HTML
- **AFK run:** `/sc-run` loops `map next` until missions are `ready` or blocked; build is atomic RED-GREEN with auto checkpoint commits; UI missions require prior draft approval and live product review (running URL + **live-product**) with paired draft-parity compare and functional evidence
- **HIL ship:** final check + `/sc-ship`
- **Quick (no mission):** `/sc-quick` - manual edits/fixes/docs; branch → verify → commit → local ship in one pass (no mission artifacts or closeout; push still explicit)
- **Active detail skills** under `.cursor/skills/` support agents (mission, planning, tdd, git, domains, sc-storm, …)
- **Explicit-only** (not auto-invoked): `sc-solid`, `sc-security`, `sc-performance`, `sc-ux-design`, `sc-diagram` - glob rules still apply. `sc-storm` activates on open-domain / strategy research feeding discuss.

Project behavior and policy are defined by the always-on files in `.cursor/rules/`.

## How we instruct agents

Clarity over tricks. Agents use Goal, Output, Good vs Bad, Verify. If unclear, research then ask - never invent Verify. Details: [docs/prompting.md](docs/prompting.md). Artifact schemas: [docs/mission-artifacts.md](docs/mission-artifacts.md).
