# Installation

Spacecraft is installed as Cursor project configuration plus a local CLI. Install it into each project where you want the mission workflow available.

## Prerequisites

- Cursor
- Git
- `curl`
- macOS or Linux
- Go 1.21 or newer when building the CLI from source

## Install with bootstrap

From a Spacecraft checkout, pass the target project directory:

```sh
./bootstrap.sh /path/to/project
```

To bootstrap the current directory:

```sh
./bootstrap.sh
```

The bootstrap installer prepares project-local `.cursor/` and `.space/` content and installs the repository CLI when a compatible prebuilt binary is available.

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

The help output should begin with `Spacecraft mission helper` and list the mission, evidence, validation, research, dependency, and roadmap commands.

## Verify Cursor discovery

After restarting Cursor:

1. Open the installed project.
2. Confirm `/sc-start`, `/sc-plan`, `/sc-build`, and `/sc-ship` are available as skills.
3. Confirm the seven agents are discoverable: `sc-coder`, `sc-tester`, `sc-planner`, `sc-reviewer`, `sc-designer`, `sc-adviser`, and `sc-firmware`.
4. Approve the project MCP server if Cursor asks for confirmation.

Workflow prompts are skills under `.cursor/skills/`. No `.cursor/commands/` directory is required.

## Start a project

If the target does not have mission state yet:

```sh
spacecraft init
```

Then begin in Cursor:

```text
/sc-start
/sc-plan
/sc-build
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
.space/
  missions/
  archive/
  roadmaps/
```

The Spacecraft repository also contains the CLI source at `cmd/spacecraft/`.
