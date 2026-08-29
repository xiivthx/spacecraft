# Optional lanes (shared firewall)

Shared SoT for optional `/sc-run` / post-ready lanes. Each lane skill owns Goal, Output prefix, Workflow, and unique Must-nots.

## Lanes

| Skill | Disposition prefix |
|---|---|
| `sc-loop` | `Loop watch:` |
| `sc-automate-slack` | `Automate-Slack:` |
| `sc-canvas-sot` | `Canvas-sot:` |
| `sc-goal-roadmap` | `Goal-roadmap:` |
| `sc-post-ready-drain` | `Post-ready drain:` |
| `sc-split-to-prs` | `Split-to-prs:` |

## Firewall (all lanes)

- **Never required** for ready or ship.
- Must not treat lane success, chat, notify, ticks, canvas, Goals, drain, or split as `AUTH` / `VERIFIED` / ready / ship authority.
- Spacecraft owns pass/fail (`ready` only on judge `VERIFIED` + empty findings; ship via AUTH + `/sc-ship`).
- When a lane was **armed** and Spacecraft hits stop (`3-cycle` | `timebox` | `blocked`) or **ready** handoff: **disarm** per that skill (do not leave armed into post-ready / ship wait).
- Disposition shape (exact strings live in each skill Output): `ran` | `skipped: <reason>` | stopped variants where that skill defines them (`stopped: <reason>`, `stopped: ready-handoff`, plus lane-specific forms such as `resumed`).

## Pointers

- Per-lane workflow: `.cursor/skills/<skill>/SKILL.md`
- Upstream Cursor skills: reference only; do not copy bodies into this repo
