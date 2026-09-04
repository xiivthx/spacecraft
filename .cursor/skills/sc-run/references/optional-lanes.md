# Companion lanes (shared firewall)

Shared SoT for `/sc-run` companion lanes. Each lane skill owns Goal, Output prefix, Workflow, and unique Must-nots.

Commander **Must** invoke each lane at its gate below and leave a greppable disposition. Silence is forbidden — `skipped: <reason>` counts as invoked.

## Lanes

| Skill | Disposition prefix | When Commander Must invoke |
|---|---|---|
| `sc-loop` | `Loop watch:` | AFK `/sc-run` start (arm or skip); again on stop (`3-cycle` \| `timebox` \| `blocked`) or ready handoff **before `sc-judge`** (disarm if armed → `stopped: ready-handoff`) |
| `sc-post-ready-drain` | `Post-ready drain:` | After `set-state ready`, before human `/sc-ship` handoff |
| `sc-split-to-prs` | `Split-to-prs:` | Same post-ready window as drain (after drain disposition); fat-diff split or explicit skip |

## Firewall (all lanes)

- Lane success is **never** `AUTH` / `VERIFIED` / ready / ship authority.
- Spacecraft owns pass/fail (`ready` only on judge `VERIFIED` + empty findings; ship via AUTH + `/sc-ship`).
- When a lane was **armed** and Spacecraft hits stop (`3-cycle` | `timebox` | `blocked`) or **ready** handoff: **disarm** per that skill (do not leave armed into post-ready / ship wait).
- Disposition shape (exact strings live in each skill Output): `ran` | `skipped: <reason>` | stopped variants where that skill defines them (`stopped: <reason>`, `stopped: ready-handoff`).

## Pointers

- Per-lane workflow: `.cursor/skills/<skill>/SKILL.md`
- Upstream Cursor skills: reference only; do not copy bodies into this repo
- Post-ship / interop stubs (not lanes): `follow-up-dispositions.md` - `Post-ship UX depth:`, `Interop/limitation:`
