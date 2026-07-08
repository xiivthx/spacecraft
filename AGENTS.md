# Spacecraft

Local-first mission-control harness for OpenCode-driven development.

## Structure

```
.opencode/          # Agent config, skill permissions
  skills/           # sc-* skills (git, plan, map, verify, etc.)
scripts/            # Go CLI: scripts/spacecraft
  src/              # main.go, types.go, internal/
.space/             # Mission state
  missions/         # spec.md, plan.json, evidence.jsonl per mission
  archive/          # Shipped missions
tests/              # Node integration tests
```

## Build & test

```sh
make build          # Go binary
make test           # Go + Node tests
```

## Conventions

### Branches
`feat/<id>/<title>` — one branch per feature. Never write on `main`.

### Commits
Conventional Commits. Target 1–3 per branch, max 5. Squash WIP before merge. Rebase on `main` → verify → `git merge --no-ff`.

### Mission ids
`M07FYB5W5` — compact sortable (prefix + base36 ms since 2026-01-01).

### Mission artifacts
- `spec.md` — what and why
- `plan.json` — ≤7 tasks, each with verify + evidence
- `evidence.jsonl` — no evidence = not done
- `map.json` — project survey before planning (in outputs/)

### Lifecycle
`/sc-start → clarify → design(if UI) → /sc-map → /sc-plan → /sc-git → work → verify → review → /sc-ship`

No implementation before `spec.md` + `plan.json`. `/sc-flow` loops `work → verify → commit`.

### Evidence
```sh
scripts/spacecraft evidence "<label>" -- <command>
```

### Release
Rebase → verify → `git merge --no-ff` → tag → delete branch → archive.

## Entry points
| File | Role |
|------|------|
| `scripts/spacecraft` | CLI |
| `AGENTS.md` | Project conventions (this file) |
| `PERSONA.md` | Commander behavior, handoff, release rules |
| `opencode.json` | Agent config, permissions, models |
| `DESIGN.md` | UI/visual design discipline |

Skill details: `.opencode/skills/sc-*/SKILL.md`
