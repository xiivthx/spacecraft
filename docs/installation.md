# Installation

Cursor project config + local CLI. Two layers; three front doors.

## Prerequisites

Cursor · Git · `curl` · macOS/Linux · Node.js 18+ (`cli/spacecraft.mjs`, `make install-machine`)

## Front doors

| Command | Role |
|---------|------|
| `make install` | User layer + CLI + smoke (default checkout) |
| `spacecraft setup` | Project packs (per repo); alias `make install-project` |
| `make install-machine` | New PC: durable clone + User layer + companions |

**User layer** (`make install` / `install-global` / `install-machine`): agents, lean-core skills, MCP, CLI, global safety hooks, short `~/.cursor/spacecraft/USER-RULES.txt` CORE from `010-hard-contract`. Paste into Settings → Rules → User Rules after regen.

**Project layer** (`spacecraft setup` / `./bootstrap.sh` / `make install-project`): alwaysApply hard-contract, pack-selected domain skills/rules, `session-start` + safety hooks, merged MCP. Never copies agents or lean-core (those stay under `~/.cursor`).

| Layer | Role |
|-------|------|
| Hooks | Hard deny/ask (secrets, destructive, main-write, ship + push ask) |
| Always-on `010-hard-contract` | Soft contract every chat |
| Skills / glob rules | On demand |

## New PC

```sh
git clone https://github.com/xiivthx/spacecraft.git
cd spacecraft
make install-machine
```

Runs User layer, then companions (caveman, rtk, codegraph) with Cursor wiring; prints a **Tools** status box. Companion steps soft-fail (warn + continue). Ensure `~/.local/bin` on `PATH`; re-paste `USER-RULES.txt`; restart Cursor; per repo `codegraph init` when indexing.

## Checkout / project

```sh
make install                          # User + CLI + smoke
spacecraft setup --packs frontend,quality
# or: ./bootstrap.sh /path/to/project
```

`make install-global` alone refreshes `~/.cursor` without smoke. Lean reconcile is **destructive**: it **prunes** spacecraft-managed domain encyclopedia skills under `~/.cursor/skills` outside the lean allowlist; unrelated files stay.

Advanced (not default): `SPACECRAFT_SKILL_PROFILE=full` or `make install FULL=1` installs domain encyclopedias into User layer.

First `.space` create (`init` / bootstrap / ensure): `git init` if needed; starter `.gitignore` with `.space/`; may soft-run `codegraph init`. Does not seed product `docs/`.

## Packs

```sh
spacecraft setup
spacecraft setup --packs frontend,quality
spacecraft setup --reconfigure --packs quality
```

TTY without profile → interactive (`quality` pre-checked). Non-TTY needs `--packs` or `SPACECRAFT_PACKS`. Profile present → silent reconcile. Selectable: `frontend`, `backend`, `database`, `embedded`, `quality`, `fpga`. Coming (listed, not installable): `iot`, `pcb`, `management`. The `frontend` pack merges the official shadcn MCP into the project's `.cursor/mcp.json` (removed on reconfigure without `frontend`; unrelated user MCP servers stay).

## Verify

```sh
test -d .cursor/rules && test -f .cursor/hooks.json && test -d .space
spacecraft help   # or ./spacecraft help
```

Restart Cursor; confirm `/sc-discuss`, `/sc-run`, `/sc-ship`, `/sc-quick`, `/sc-debug` and User-layer agents under `~/.cursor/agents/`.

## Docs vs `.space/`

- **`docs/`** — tracked product SoT (commit)
- **`.space/`** — gitignored mission runtime (do not commit)

Promote durable contracts at ship into `docs/specs/` or an ADR.
