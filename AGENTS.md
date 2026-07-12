# Spacecraft

Local-first mission-control harness for OpenCode-driven development.

> **Read together with [`PERSONA.md`](./PERSONA.md)** — lane decisions, commander behavior, self-review, and release rules live there. Always check both files.

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

### Development lanes

Commander auto-detects intent and routes to the appropriate lane. No user decision required.

| Lane | Intent | Entry | Workflow |
|------|--------|-------|----------|
| 💬 **Advisory** | ask, talk, consult, research | default | direct response |
| 🚀 **Mission** | add, build, implement, feature | `/sc-start` | full flow — below |
| 🔧 **Debug** | fix, debug, diagnose, error | `/sc-debug` | 5-step discipline |
| ⚡ **Quick** | human edits, config, small fix | `/sc-quick` | branch → self-review → ship |

#### Advisory lane (default)
- Questions, discussion, research, consultation — no mutating work
- No mission created, no git operations
- Research auto-trigger still active (`spacecraft research`)

#### Mission lane
`/sc-start → /sc-design(if UI) → /sc-plan → /sc-build → /sc-ship`

- Full artifacts: `spec.md`, `plan.json`, `evidence.jsonl`, `review.md`, `review.json`
- Commander auto-handles clarification, mapping, verification
- No implementation before `spec.md` + `plan.json`
- `/sc-build` loops per task: implement → verify → checkpoint commit
- **Zero trust**: sc-reviewer reviews plan inside `/sc-plan`, diff + evidence inside `/sc-build`. `/sc-review` is a standalone manual command — not part of the pipeline.

#### Debug lane
`/sc-debug` — five-step discipline: reproduce → trace fail path → falsify hypothesis → cross-reference → post-mortem

- Scoped to fix/diagnose. No feature scope creep
- Evidence captured within debug workflow

#### Quick lane
`/sc-quick → branch → commit freely → fast self-review → report ready → wait for explicit /sc-ship`

- For: prompt tweaks, config, docs, small fixes — where full flow is overhead
- Skips: `spec.md`, `plan.json`, TDD build, formal review, evidence capture
- Keeps: git safety, Conventional Commits, changelog, versioning, `--no-ff` merge
- Commander self-review directly — no subagent, no review artifacts

### Evidence
```sh
scripts/spacecraft evidence "<label>" -- <command>
```

### Release
Merge to main only on explicit `/sc-ship` or user release command. Never auto-detect.
Rebase → verify → `git merge --no-ff` → tag → delete branch → archive.

## Additional skills

| Skill | Purpose |
|-------|---------|
| `sc-performance` | Performance review — N+1 detection, memory leaks, bundle size, render optimization |
| `sc-security` | Static security review — OWASP detection, secrets, injection patterns, manifest scanning |

## Entry points
| File | Role |
|------|------|
| `scripts/spacecraft` | CLI |
| `AGENTS.md` | Project conventions (this file) |
| `PERSONA.md` | Commander persona, lane detection, session handoff, release rules — **always read with AGENTS.md** |
| `opencode.json` | Agent config, permissions, models |
| `DESIGN.md` | UI/visual design discipline |

Skill details: `.opencode/skills/sc-*/SKILL.md`
