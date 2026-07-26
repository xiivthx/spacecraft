# Installation

Spacecraft is installed as Cursor project configuration plus a local CLI. Install it into each project where you want the mission workflow available.

## Prerequisites

- Cursor
- Git
- `curl`
- macOS or Linux
- Go 1.21 or newer when building the CLI from source

## User layer vs Project layer

Spacecraft installs in two layers:

- **User layer** (`make install-global`, once per machine): agents, skills, MCP config, the CLI, and global safety hooks. It also generates `~/.cursor/spacecraft/USER-RULES.txt` from the five `alwaysApply` rules (`000-spacecraft`, `025-english-coach`, `050-style`, `100-conventions`, `200-workflow`). Paste that file's contents into Cursor Settings -> Rules -> User Rules once - that is how the `alwaysApply` rules take effect in every workspace, since Cursor does not read a repo's `alwaysApply: true` rules outside that repo.
- **Project layer** (`./bootstrap.sh` or `make install-project`, once per repo): the domain/glob rules `300`-`620`, agents, skills, project hooks (including `session-start`), and a merged `.cursor/mcp.json`. It never copies the `alwaysApply` rules - those stay User layer only, so installing into many projects never re-duplicates them.

Run the User layer install once per machine; each Project layer install is independent and repeatable.

## Install with bootstrap

From a Spacecraft checkout, pass the target project directory:

```sh
./bootstrap.sh /path/to/project
```

To bootstrap the current directory:

```sh
./bootstrap.sh
```

The bootstrap installer prepares project-local `.cursor/` and `.space/` content and installs the repository CLI when a compatible prebuilt binary is available. This is the Project layer only - see [User layer vs Project layer](#user-layer-vs-project-layer) for the one-time global setup.

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

`make install` builds the Go CLI from `cmd/spacecraft/` and installs Spacecraft for use from Cursor and your shell. Ensure `~/.local/bin` is on `PATH` if your shell cannot find `spacecraft`:

```sh
export PATH="$HOME/.local/bin:$PATH"
```

For the User layer - Cursor-wide agents, slash skills (`/sc-discuss`, `/sc-run`, `/sc-ship`, `/sc-quick`, and other `sc-*` skills), and the `alwaysApply` rules:

```sh
make install-global
```

That copies `~/.cursor/agents/sc-*.md` and `~/.cursor/skills/sc-*/`, merges MCP into `~/.cursor/mcp.json`, links the CLI, installs global safety hooks (`check-main-write`, `check-ship-commands`) into `~/.cursor/hooks.json`, and generates `~/.cursor/spacecraft/USER-RULES.txt` from the five `alwaysApply` rules. Paste that file's contents into Cursor Settings -> Rules -> User Rules once so Commander, workflow, English coaching, style, and conventions apply in every workspace. Restart Cursor afterward. Unrelated skills and hooks (for example personal ones) are left alone.

For the Project layer in another repo, either run `./bootstrap.sh /path/to/project` or, from this checkout:

```sh
make install-project PROJECT=/path/to/project
```

Both install the domain/glob rules (`300`-`620`), agents, skills, and project hooks (including `session-start`) - never the `alwaysApply` rules, which are User layer only.

To build without installing:

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

Confirm the CLI:

```sh
spacecraft help
```

When using the repository binary directly:

```sh
./spacecraft help
```

The help output should begin with `Spacecraft mission helper` and list the mission, evidence, validation, and roadmap commands.

## Verify Cursor discovery

After restarting Cursor:

1. Open the installed project.
2. Confirm `/sc-discuss`, `/sc-run`, `/sc-ship`, and `/sc-quick` are available as skills.
3. Confirm the eight agents are discoverable: `sc-coder`, `sc-tester`, `sc-planner`, `sc-reviewer`, `sc-designer`, `sc-adviser`, `sc-firmware`, and `sc-writer`.
4. Approve the project MCP server if Cursor asks for confirmation.

Workflow prompts are Agent Skills under `.cursor/skills/` (explicit `/` via `disable-model-invocation: true`). Do not migrate them to `.cursor/commands/` - Cursor's direction is Commands → Skills (`/migrate-to-skills`).

## Start a project

If the target does not have mission state yet:

```sh
spacecraft init
```

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
  trust/                   # local source of trust (lessons.md, solved.md); seeded by init
```

Trust is not committed. Tracked seed: `.cursor/skills/sc-learn/references/trust-seed/`. Agents read `.space/trust/lessons.md` before inventing process.

The Spacecraft repository also contains the CLI source at `cmd/spacecraft/`.
