# Installation

Spacecraft is installed as Cursor project configuration plus a local CLI. Install it into each project where you want the mission workflow available.

## Prerequisites

- Cursor
- Git
- `curl`
- macOS or Linux
- Node.js 18 or newer for the CLI (`cli/spacecraft.mjs`) and companion installs (`make install-machine`)

## User layer vs Project layer

Spacecraft installs in two layers:

- **User layer** (`make install-global` or `make install-machine`, once per machine): agents, skills, MCP config, the CLI, and global safety hooks. It also generates `~/.cursor/spacecraft/USER-RULES.txt` from six User-layer sources (`000-spacecraft`, `026-intent-coach`, `027-th-en-hil`, `050-style`, `100-conventions`, `200-workflow`). Paste that file's contents into Cursor Settings -> Rules -> User Rules so those policies apply in every workspace (Cursor does not load them from a repo checkout alone). After each User-layer regen, re-paste the updated file - human step; no Settings API automation. **Agent chat language:** English default for technical substance; Thai for HIL questions, short status, and handoff summaries (no dual `ไทย:` / `EN:` blocks).
- **Project layer** (`./bootstrap.sh` or `make install-project`, once per repo): the domain/glob rules `300`-`620`, agents, skills, project hooks (including `session-start`), and a merged `.cursor/mcp.json`. It never copies the User-layer rules - those stay User layer only, so installing into many projects never re-duplicates them.

Run the User layer install once per machine; each Project layer install is independent and repeatable.

## New PC / install-machine

For a fresh machine, install the User layer plus companion CLIs in one step:

```sh
git clone https://github.com/xiivthx/spacecraft.git
cd spacecraft
make install-machine
```

From an existing checkout you can also run `scripts/install-machine.sh` directly (`make install-machine` links the Node CLI first).

`make install-machine` invokes `scripts/install-machine.sh`, which:

1. Clones or updates Spacecraft into a durable directory (`~/.local/share/spacecraft` by default; override with `SPACECRAFT_INSTALL_DIR`).
2. Runs `make install-global` from that clone (User layer: agents, skills, MCP, CLI, hooks, `USER-RULES.txt`).
3. Installs three companion tools and wires them to Cursor:
   - **caveman** - official installer (requires Node.js 18+)
   - **rtk** - official installer, then `rtk init -g --agent cursor`
   - **codegraph** - official installer, then `codegraph install --target=cursor --yes`
4. Prints a tokless-like **Tools** status box (checkmark or cross per tool, version when available).

Accepted flags: `--agents cursor` and `--yes` for explicit non-interactive Cursor wiring. v1 has no interactive agent picker; the script is non-interactive by default.

Companion installs **soft-fail**: a failed caveman, rtk, or codegraph step prints a warning and the script continues. The spacecraft User layer completes before companions run; the script exits non-zero only if the clone or User-layer install fails.

**Does not:** install the Cursor app, or run Project-layer bootstrap on arbitrary repos. Per-repo setup stays `./bootstrap.sh` or `make install-project`.

After install:

- Ensure `~/.local/bin` is on your `PATH`.
- Paste `~/.cursor/spacecraft/USER-RULES.txt` into Cursor Settings -> Rules -> User Rules (re-paste after every User-layer regen).
- Restart Cursor to pick up skills, hooks, and companion wiring.
- Per project: run `codegraph init` in each repo you want indexed.

Re-running `make install-machine` updates the durable clone in place and refreshes the User layer and companions. Re-paste `USER-RULES.txt` afterward so Cursor picks up the regenerated text.

## Install with bootstrap

From a Spacecraft checkout, pass the target project directory:

```sh
./bootstrap.sh /path/to/project
```

To bootstrap the current directory:

```sh
./bootstrap.sh
```

The bootstrap installer prepares project-local `.cursor/` and `.space/` content and links the Node CLI (`cli/spacecraft.mjs`) when Node.js is on `PATH`, then runs project smoke checks against that link. On first `.space` create it ensures a git repo (`git init` when needed), a starter `.gitignore` from `templates/gitignore` (always ignores `.space/`), and may soft-run `codegraph init` when no index (`.codegraph/`) exists - missing binary or failure warns and continues. This is the Project layer only - see [User layer vs Project layer](#user-layer-vs-project-layer) for the one-time global setup.

You can also run the published bootstrap script from the target project:

```sh
curl -fsSL https://raw.githubusercontent.com/xiivthx/spacecraft/main/bootstrap.sh | sh
```

Restart Cursor after installation so it refreshes the project's Cursor configuration.

## Build and install from source

Clone the repository, then use the Makefile:

```sh
git clone https://github.com/xiivthx/spacecraft.git
cd spacecraft
make install
```

`make install` links the Node CLI (`cli/spacecraft.mjs`) and installs Spacecraft for use from Cursor and your shell. Ensure `~/.local/bin` is on `PATH` if your shell cannot find `spacecraft`:

```sh
export PATH="$HOME/.local/bin:$PATH"
```

For the User layer only (no companions) - Cursor-wide agents, slash skills (`/sc-discuss`, `/sc-run`, `/sc-ship`, `/sc-quick`, and other `sc-*` skills), and the six User-layer sources via `USER-RULES.txt`:

```sh
make install-global
```

Default User-layer skills are lean-core (lifecycle + process). Lean reconcile is destructive for spacecraft-managed domain packs under `~/.cursor/skills` - it prunes encyclopedia skills outside the allowlist. Unrelated files under `~/.cursor` stay put. Opt in to domain encyclopedias with `SPACECRAFT_SKILL_PROFILE=full` or `make install-global FULL=1` (documented `--full` equivalent):

```sh
SPACECRAFT_SKILL_PROFILE=full make install-global
# or: make install-global FULL=1
```

On a new PC, prefer [`make install-machine`](#new-pc--install-machine) to also install caveman, rtk, and codegraph with Cursor wiring.

That copies `~/.cursor/agents/sc-*.md` and `~/.cursor/skills/sc-*/`, merges MCP into `~/.cursor/mcp.json`, links the CLI, installs global safety hooks (`check-main-write`, `check-ship-commands`) into `~/.cursor/hooks.json`, and generates `~/.cursor/spacecraft/USER-RULES.txt` from the six User-layer sources. Paste that file's contents into Cursor Settings -> Rules -> User Rules so Commander, workflow, intent coaching, Agent chat language (English default; Thai for HIL/status/handoff), style, and conventions apply in every workspace. After each regen, re-paste the file (human step - no Settings API automation). Restart Cursor afterward. Unrelated skills and hooks (for example personal ones) are left alone.

For the Project layer in another repo, either run `./bootstrap.sh /path/to/project` or, from this checkout:

```sh
make install-project PROJECT=/path/to/project
```

Both install the domain/glob rules (`300`-`620`), agents, skills, and project hooks (including `session-start`) - never the User-layer rules, which stay User layer only.

To link the Node CLI in the checkout without a full install:

```sh
make build
./spacecraft help
```

## Verify the installation

In the target project, confirm the Cursor-native files:

```sh
test -d .cursor/rules
test -d .cursor/agents
test -d .cursor/skills
test -f .cursor/mcp.json
test -f .cursor/hooks.json
test -d .space
```

Confirm the Node CLI (linked entry or `PATH`):

```sh
spacecraft help
```

From a checkout or bootstrap target that has the `./spacecraft` link:

```sh
./spacecraft help
```

The help output should begin with `Spacecraft mission helper` and list the mission, evidence, validation, and roadmap commands. Smoke after bootstrap uses the same Node CLI path - Node.js 18+ on `PATH` is enough for that check.

## Verify Cursor discovery

After restarting Cursor:

1. Open the installed project.
2. Confirm `/sc-discuss`, `/sc-run`, `/sc-ship`, and `/sc-quick` are available as skills.
3. Confirm detail skill `sc-storm` is discoverable (Tier 3 open-domain research; not a lifecycle slash).
4. Confirm the eight agents are discoverable: `sc-coder`, `sc-tester`, `sc-planner`, `sc-reviewer`, `sc-designer`, `sc-adviser`, `sc-firmware`, and `sc-writer`.
5. Approve the project MCP server if Cursor asks for confirmation.

Workflow prompts are Agent Skills under `.cursor/skills/` (explicit `/` via `disable-model-invocation: true`). Do not migrate them to `.cursor/commands/` - Cursor's direction is Commands → Skills (`/migrate-to-skills`).

## Start a project

If the target does not have mission state yet:

```sh
spacecraft init
```

First-use ensure (`spacecraft init`, any CLI command when `.space/` is missing, and project bootstrap / `install-cursor`) runs `ensureProjectReady`: `git init` if the directory is not already a repo, scaffolds `.space/`, writes starter `.gitignore` from `templates/gitignore` (includes `.space/`), and may soft-run `codegraph init` when no index (`.codegraph/`) exists - missing binary or failure warns and continues. When `.space/` already exists, later ensure only adds a `.space/` line to `.gitignore` if missing - it does not replace the whole file.

Then begin in Cursor:

```text
/sc-discuss
/sc-run
/sc-ship
```

## Installed layout

```text
.cursor/
  rules/
  agents/
  skills/
  mcp.json
  hooks.json
  hooks/
.space/                    # fully gitignored (local state)
  missions/
  archive/
  roadmaps/
```

The Spacecraft repository CLI lives at `cli/spacecraft.mjs` (stdlib Node; zero npm CLI dependencies).
