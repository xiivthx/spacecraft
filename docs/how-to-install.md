# Installing Spacecraft for OpenCode

Spacecraft is a local-first mission-control harness for OpenCode-driven development. It provides mission lifecycle management (spec → plan → build → review → ship), TDD discipline, subagent orchestration, and evidence-based verification gates.

## Prerequisites

- [OpenCode.ai](https://opencode.ai) installed
- [Go](https://go.dev/dl/) 1.26+ (for the mission-control CLI)
- Git

## Installation

### Local (per-workspace)

```sh
make build
```

This compiles the Go helper binary (`scripts/spacecraft`) used for mission state, evidence capture, git branching, research, and resolver operations.

Spacecraft loads as a local OpenCode plugin. On startup, OpenCode will:

1. Load `.opencode/plugins/engine.js` — registers skills paths, commands, and injects commander persona
2. Read `opencode.json` — registers agent definitions (sc-commander, sc-coder, sc-tester, etc.)
3. Inject PERSONA.md + AGENTS.md + DESIGN.md context on session start

### Global (all workspaces)

```sh
make install
```

This:
- Builds the CLI
- Symlinks `spacecraft` to `~/.local/bin/`
- Writes global config to `~/.config/opencode/opencode.jsonc` with absolute paths

Restart OpenCode after install.

## Verify

Ask the commander:

```
/sc-resume
```

Or verify skills are discovered:

```
load skill sc-clarify
```

## Manual alternative

If the plugin doesn't auto-load, check your `opencode.json` includes the project path:

```json
{
  "plugin": ["./.opencode/plugins/engine.js"]
}
```

Both regular and git-backed installs work. To use a specific tag:

```json
{
  "plugin": ["spacecraft@git+https://github.com/anomalyco/spacecraft.git#v0.21.1"]
}
```

## Updating

```sh
git pull origin main
make build
```

New skills, commands, and persona updates take effect on next OpenCode restart.

## Troubleshooting

### CLI missing

```sh
make build
```

Verify: `scripts/spacecraft help`

### Plugin not loading

Check logs for engine.js errors. Verify `.engine/` and `.opencode/` directories exist at workspace root.

### Skill not found

Verify the skill file exists:
```sh
ls .engine/skills/<category>/sc-<name>/SKILL.md
```

### Commands not showing

Ensure `.engine/commands/sc-*.md` files exist and the plugin's config hook ran without errors.

## Getting Help

- Report issues: https://github.com/anomalyco/spacecraft/issues
