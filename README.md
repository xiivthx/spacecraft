# Spacecraft

Cursor-native mission-control harness: rules, agents, skills, safety hooks, and a local CLI for traceable mission work under `.space/`.

## Requirements

Cursor · Git · Node.js 18+ · macOS or Linux

## Install

| Command | Role |
|---------|------|
| `make install` | User layer + CLI + smoke |
| `spacecraft setup` | Project packs (per repo) |
| `make install-machine` | New PC + companions |

```sh
git clone https://github.com/xiivthx/spacecraft.git
cd spacecraft
make install
spacecraft setup --packs quality
```

Details: [docs/installation.md](docs/installation.md)

Companions (`make install-machine`): caveman, rtk, codegraph — Tools status box.

## Quick start

Slash skills: `/sc-discuss` → `/sc-run` → human check → `/sc-ship`. Small edits: `/sc-quick`. Unknown RCA: `/sc-debug`.

```
/sc-discuss
/sc-run
/sc-ship
```

1. Discuss — size, clarify, decide; approve visual draft when needed; clear clarify-status
2. Run — plan → build → review → judge `VERIFIED` → `ready`
3. Human check
4. Ship — explicit `/sc-ship` only

Roadmap helpers: `spacecraft map use|current|next`. Live CLI: `spacecraft help`.

Validate: `spacecraft validate` / `val` — Validate mission artifacts and evidence (not-doc-drift).

## Lanes

| Intent | Lane |
|--------|------|
| Clarify / draft | `/sc-discuss` |
| Implement | `/sc-run` |
| Debug RCA | `/sc-debug` |
| Small edit | `/sc-quick` |
| Ship | `/sc-ship` |

SoT: `.cursor/rules/200-workflow.mdc`

## Layout

```text
.cursor/     rules, skills, hooks, MCP
docs/        tracked product SoT
.space/      gitignored missions / evidence
cli/         Node CLI (spacecraft.mjs)
```

## Details

- [Installation](docs/installation.md)
- [Mission artifacts](docs/mission-artifacts.md)
- [Prompting](docs/prompting.md)
- [Docs map](docs/README.md)
- Agents after install: `~/.cursor/agents/` (`sc-coder`, `sc-tester`, `sc-planner`, `sc-reviewer`, `sc-designer`, `sc-adviser`, `sc-firmware`, `sc-rtl`, `sc-writer`, `sc-browser-probe`, `sc-fact-check`, …)
