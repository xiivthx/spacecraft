# Spacecraft Skills

OpenCode skill suite for structured, evidence-driven development. Available globally via `~/.config/opencode/`.

## Conventions

### Git
- **Never write on `main`.** Create a feature branch.
- **Conventional Commits.** Target 1–3 per branch, max 5. Squash WIP before merge.
- **Rebase on `main`** → verify → `git merge --no-ff`.

### Evidence
No evidence = not done. Every implementation task must include concrete verification output (command output, test results, screenshots).

### Zero trust
All AI output must be reviewed by an independent subagent before it becomes truth. No output becomes authoritative without passing a review gate.

## Development lanes

Commander auto-detects intent and routes to the appropriate lane.

| Lane | Intent | Entry | Workflow |
|------|--------|-------|----------|
| Advisory | ask, talk, consult, research | default | direct response |
| Mission | add, build, implement, feature | `/sc-start` | full flow |
| Debug | fix, debug, diagnose, error | `/sc-debug` | 5-step discipline |
| Quick | human edits, config, small fix | `/sc-quick` | branch → self-review → ship |

## Available commands

| Command | Purpose |
|---------|---------|
| `/sc-start` | Begin a new mission (feature). Creates spec, branch, and mission artifacts. |
| `/sc-design` | UI/visual design review and critique. |
| `/sc-plan` | Convert spec into executable plan (≤7 tasks). |
| `/sc-build` | Implement tasks one-by-one with verification and checkpoint commits. |
| `/sc-review` | Formal review of diff, evidence, and release readiness. |
| `/sc-ship` | Merge to main, tag, archive mission. |
| `/sc-quick` | Fast lane for small changes: branch, commit, self-review, ship. |
| `/sc-resume` | Resume an active mission with full context handoff. |

## Available skills

| Skill | Purpose |
|-------|---------|
| `sc-architect` | Design system architecture, write ADRs, C4 diagrams, and tradeoff analysis. |
| `sc-clarify` | Resolve ambiguity through focused user clarification. |
| `sc-creator` | Create new Spacecraft skills from datasources. |
| `sc-database` | Design schemas, write migrations, optimize queries, and manage indexes (PostgreSQL default). |
| `sc-debug` | Five-step debugging discipline (reproduce → trace → falsify → cross-reference → post-mortem). |
| `sc-design` | Shape, critique, and polish UI/visual design. |
| `sc-git` | Git safety, branching, Conventional Commits, no-ff merge, versioning. |
| `sc-learn` | Capture mission knowledge, issues, and lessons learned. |
| `sc-map` | Survey project structure before planning. |
| `sc-mission` | Manage mission artifacts and lifecycle. |
| `sc-planning` | Convert spec into small executable plan with verifiable tasks. |
| `sc-solid` | SOLID principles and code quality discipline. |
| `sc-tdd` | Test-driven development discipline (Plan-Red-Green-Verify-Refactor-Review). |
| `sc-verification` | Capture fresh command evidence before claiming work complete. |
| `sc-web-backend` | Build server APIs with Node.js, TypeScript, Fastify, and Vitest. |
| `sc-web-frontend` | Build browser UI with React, TypeScript, Vite, Tailwind CSS, and Vitest. |
| `sc-web-service` | Build lean local web service from scratch under mission control. |

## Per-project setup

Each project using spacecraft skills should have its own `AGENTS.md` (project conventions) and `PERSONA.md` (commander persona) at the project root.

Skill details: `.opencode/skills/sc-*/SKILL.md`
