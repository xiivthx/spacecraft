# Spacecraft

Local-first mission-control harness for OpenCode-driven development.

> **Read together with [`PERSONA.md`](./PERSONA.md)** — lane decisions, commander behavior, self-review, and release rules live there. Always check both files.

## Structure

```
.engine/            # OpenCode engine config, skills, and conventions
  agents/           # sc-* agent definitions (commander, coder, reviewer, etc.)
  skills/           # sc-* skills in categorized subdirs (core/, quality/, design/, etc.)
  opencode.json     # Agent config, permissions, models
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

### Standards
- Never use the em dash ("-"). Use plain dash instead.
- Never manually modify CHANGELOG.md or any auto-generated files.
- Prefer quality, simplicity, robustness, scalability, and long-term maintainability over development cost.
- Bug fixes: reproduce in E2E setting as close to end-user experience as possible before fixing.
- E2E testing: be picky about UI, obsessed with pixel perfection. Fix anything that looks off, even if unrelated.
- Engineering excellence: lint, test failures, test flakiness - fix even if not caused by current work.
- Bug radar: when you spot incorrect workflow, a silent logic bug, or a broken assumption during any lane (even advisory), auto-create a GitHub issue. Don't let discovered problems rot in chat history.

### Mission ids
`M07FYB5W5` — compact sortable (prefix + base36 ms since 2026-01-01).

### Mission artifacts
- `spec.md` — what and why
- `plan.json` — ≤7 tasks, each with verify + evidence
- `evidence.jsonl` — no evidence = not done
- `map.json` — project survey before planning (in outputs/)

## Available commands

| Command | Purpose |
|---------|---------|
| `/sc-start` | Begin a new mission (feature). Creates spec, branch, and mission artifacts. |
| `/sc-design` | UI/visual design review and critique. |
| `/sc-plan` | Convert spec into executable plan (≤7 tasks). |
| `/sc-build` | Implement tasks one-by-one with verification and checkpoint commits. |
| `/sc-review` | Formal review of diff, evidence, and release readiness. |
| `/sc-ship` | Merge to main, tag, archive mission. |
| `/sc-quick` | Fast lane for small changes: branch, commit, self-review, report ready. Ship only on explicit `/sc-ship`. |
| `/sc-research` | Run systematic research via spacecraft research CLI (Brave Search, scoped docs, deep analysis). |
| `/sc-resume` | Resume an active mission with full context handoff. |

### Development lanes

Commander auto-detects intent and routes to the appropriate lane without user input. Full lane behavior, decision flow, self-review, release rules, and session handoff: see [`PERSONA.md`](./PERSONA.md).

### Available agents

| Agent | Role |
|-------|------|
| `sc-commander` | Primary agent — mission orchestration, lane detection, subagent delegation |
| `sc-coder` | Write-capable — implements production code |
| `sc-tester` | Write-capable — writes tests and captures evidence (TDD) |
| `sc-planner` | Read-only — converts spec into executable plan |
| `sc-reviewer` | Read-only — reviews diffs, evidence, and release readiness |
| `sc-designer` | Read-only — shapes UI direction and critiques visual design |
| `sc-adviser` | Read-only — complex system design and deep logic restructuring, invoked on escalation |

### Evidence
```sh
scripts/spacecraft evidence "<label>" -- <command>
```

## Available skills

| Skill | Purpose |
|-------|---------|
| `sc-architect` | Design system architecture, write ADRs, C4 diagrams, and tradeoff analysis |
| `sc-clarify` | Resolve ambiguity through focused user clarification |
| `sc-creator` | Create new Spacecraft skills from datasources |
| `sc-database` | Design schemas, write migrations, optimize queries, and manage indexes (PostgreSQL default) |
| `sc-debug` | Five-step debugging discipline (reproduce → trace → falsify → cross-reference → post-mortem) |
| `sc-design` | Shape, critique, and polish UI/visual design |
| `sc-git` | Git safety, branching, Conventional Commits, no-ff merge, versioning |
| `sc-learn` | Capture mission knowledge, issues, and lessons learned |
| `sc-llm-vision` | Use LLM vision models (Gemini via agy CLI) to review UI screenshots |
| `sc-localize` | Review bilingual copy for cultural fit |
| `sc-map` | Survey project structure before planning |
| `sc-memory` | Wraps ctx_search and ctx_index with spacecraft conventions for structured cross-mission memory |
| `sc-mission` | Manage mission artifacts and lifecycle |
| `sc-pathfinder` | Chart a map of tickets for large, multi-session work |
| `sc-performance` | Performance review — N+1 detection, memory leaks, bundle size, render optimization |
| `sc-planning` | Convert spec into small executable plan with verifiable tasks |
| `sc-search` | Quick internet search with 3-tier escalation for stuck issues, gray areas, and stale knowledge |
| `sc-security` | Static security review — OWASP detection, secrets, injection patterns, manifest scanning |
| `sc-solid` | SOLID principles and code quality discipline |
| `sc-tdd` | Test-driven development discipline (Plan-Red-Green-Verify-Refactor-Review) |
| `sc-ux-design` | UI quality control: anti-slop enforcement, draft previews, animation quality, visual verification |
| `sc-verification` | Capture fresh command evidence before claiming work complete |
| `sc-web-backend` | Build server APIs with Node.js, TypeScript, Fastify, and Vitest |
| `sc-web-frontend` | Build browser UI with React, TypeScript, Vite, Tailwind CSS, and Vitest |

## Entry points

| File | Role |
|------|------|
| `scripts/spacecraft` | CLI |
| `.engine/AGENTS.md` | Project conventions (this file) |
| `.engine/PERSONA.md` | Commander persona, lane detection, session handoff, release rules — **always read with AGENTS.md** |
| `.engine/opencode.json` | Agent config, permissions, models |
| `.engine/DESIGN.md` | UI/visual design discipline |

Skill details: `.engine/skills/*/sc-*/SKILL.md`

## Research auto-trigger

When encountering gray areas, outdated knowledge, or uncertainty, invoke the search escalation via `sc-search` skill. Full escalation tiers, trigger conditions, and examples: see [`PERSONA.md`](./PERSONA.md) §4 Expertise — Research.
