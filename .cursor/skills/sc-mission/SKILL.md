---
name: sc-mission
description: "Manage mission artifacts and lifecycle. Activate on /sc-discuss, /sc-run mission work, mission creation, status check, or lifecycle management."
---

# sc-mission

Manage mission artifacts and lifecycle for local Cursor development.

## When to use

Resolve/list/select a mission; create artifacts; check state; route ambiguity to `/sc-discuss`; session handoff or release closeout prep.

## Workflow

1. **Resolve** - `spacecraft resolve` / `status` / `missions`; `.space/current` is fallback, not sole authority. Conflict → `spacecraft use <selector>`. Detect lane (Discuss, Mission, Debug, Quick) from intent.
2. **Read** - `mission.json`, `spec.md`, `questions.md`, `decisions.md`, `plan.json`, `design-contract.md` / `approved-scenarios.md` when present, design artifacts, `evidence.jsonl`, `review.json`.
3. **Route ambiguity** - intent/scope/acceptance unclear → `/sc-discuss` before `/sc-run`. Apply mission sizing (`sc-discuss/references/mission-sizing.md`) before bind/run; discuss owns `spacecraft map` create/add and Resize.
4. **Enforce lifecycle** - discuss → `/sc-run` (plan → design-contract → approved-scenarios → build → verify → review → judge) → ship. AFK requires design-contract + approved-scenarios (or docs/prose skips) before product RED/GREEN; combine needs static + diff-cov + mutation disposition (evidence or skip/waive per `docs/mission-artifacts.md`).
5. **Release or handoff** - ship intent → `/sc-ship`. Session end → handoff summary. Prefer new session when phase changes.

### Edge cases

- No mission + mutating work → `spacecraft new "<title>"` + work branch; note in `decisions.md`.
- Multiple matches → show candidates; `spacecraft use`.
- Mid-session end → summarize task, blockers, dirty git, next pickup; do not merge/tag/delete branches.
- Invalid state transition or unparseable `plan.json` / `evidence.jsonl` → block; ask regenerate/restore.
- Plan without spec → treat as corrupt; spec is SoT.

## Lifecycle states

| State | Meaning | Next |
|-------|---------|------|
| `active` | Created | `planned`, `blocked` |
| `planned` | Plan ready | `in_progress`, `blocked` |
| `in_progress` | Building | `ready`, `blocked` |
| `ready` | Review + judge passed | `shipped`, `blocked` |
| `blocked` | Gated | `active`, `in_progress` |
| `shipped` | Merged/archived | Terminal |

`spacecraft set-state [mission-id] <state>` enforces transitions.

## Rules

- **Must**: Resolve via CLI priority (explicit selector → `.space/current` → branch `feat/<id>/…`); block writes on strong conflicts until `use`.
- **Must**: Read artifacts before mutating; route ambiguity to `/sc-discuss`; enforce discuss → plan → design-contract → approved-scenarios → build → verify → review → ship.
- **Must**: AFK product path has design-contract + approved-scenarios (or skips) before RED/GREEN; static / diff-cov / mutation disposition before ready (`docs/mission-artifacts.md`).
- **Must**: Use sc-git for git safety; mutating work satisfies sc-git gates; ship/release requests are closeout prep (block if gates incomplete).
- **Must**: After merge, `set-state shipped` (archive under `.space/archive/` unless asked otherwise); capture evidence; prefer evidence over narrative claims.
- **Must not**: Skip clarify when materially ambiguous; implement without spec/plan; write product on `main`.
- **Must**: End sessions with next action; prefer new chat when phase changes or context is heavy.

## Out of scope

sc-git · sc-planning · sc-ux-design / sc-designer · sc-verification · `/sc-discuss` clarify/draft · Debug Mode · Task(`sc-reviewer`)

## Layout

```
.space/missions/<id>/
  mission.json  spec.md  questions.md  decisions.md
  plan.json  design-contract.md  approved-scenarios.md
  evidence.jsonl  review.json  design/  outputs/
```

Skip/waive prefixes: `docs/mission-artifacts.md` (Outcome-gate skip / waive grammar).

## References

- `spacecraft resolve|use|evidence|closeout-check --help`
- `docs/mission-artifacts.md` - schemas + outcome-gate skip/waive SoT
