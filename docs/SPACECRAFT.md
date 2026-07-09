# Spacecraft

Spacecraft is a lean, local-first OpenCode harness for mission-driven software development.

**Persona:** `PERSONA.md` · **Rules:** `AGENTS.md` · **Spec:** `SPEC.md` · **Design:** `DESIGN.md`

## Quickstart

```sh
make build                          # build the Go helper binary
scripts/spacecraft new "My mission" # create a mission
scripts/spacecraft missions         # list missions
scripts/spacecraft use 1            # select a mission
scripts/spacecraft status           # show resolved mission state
```

## Slash commands

`/sc-start` · `/sc-design` · `/sc-plan` · `/sc-build` · `/sc-review` · `/sc-quick` · `/sc-resume` · `/sc-ship`

## CLI commands

| Command | What it does |
|---|---|
| `scripts/spacecraft archive [selector]` | Compact shipped mission |
| `scripts/spacecraft bind-branch` | Record current branch on mission |
| `scripts/spacecraft check-deps [path]` | Check dependency freshness |
| `scripts/spacecraft clarify-status <status>` | Set clarification status |
| `scripts/spacecraft closeout-check` | Check release readiness |
| `scripts/spacecraft current` | Print currently active mission |
| `scripts/spacecraft evidence "<label>" -- <cmd>` | Capture verification evidence |
| `scripts/spacecraft flow` | Print workflow readiness (alias: `workflow`) |
| `scripts/spacecraft git-info` | Print git status |
| `scripts/spacecraft git-suggest [type] [slug]` | Suggest branch/commit names |
| `scripts/spacecraft help` | Show CLI help |
| `scripts/spacecraft init` | Create `.space/` structure |
| `scripts/spacecraft missions` | List all missions |
| `scripts/spacecraft new "<title>"` | Create and select a mission |
| `scripts/spacecraft research "<query>"` | Search web, registries, and analyze |
| `scripts/spacecraft resolve [--json]` | Resolve active mission |
| `scripts/spacecraft set-state <state>` | Set mission state |
| `scripts/spacecraft status` | Print resolved mission state |
| `scripts/spacecraft use <selector>` | Select mission by number/id/title |
| `scripts/spacecraft validate` | Validate mission artifacts |

## File layout

```
.space/
  current                      fallback mission id
  sessions/                    session→mission bindings
  missions/<id>/               mission artifacts + outputs
    mission.json, spec.md, plan.json, evidence.jsonl
    questions.md, decisions.md, review.md, review.json
    design/, outputs/
  archive/<id>/                compact shipped missions
.opencode/agents/              agent configs (sc-commander, etc.)
.opencode/commands/            slash command prompts
.opencode/skills/              reusable skills
scripts/spacecraft             Go binary helper
```

## Mission lifecycle (summary)

`/sc-start → /sc-design(if UI) → /sc-plan → /sc-build → /sc-ship`

Commander auto-handles clarification, mapping, git hygiene, review, and verification within these steps. sc-reviewer reviews plans and diffs. `/sc-review` is a standalone manual command — not in the pipeline.

See `SPEC.md` §Workflow and §Verification for full gate rules.

## Git

- One non-main branch per feature/fix (`<type>/<id>/<title>`)
- Never write product code on main
- Final branch: 1–3 Conventional Commits (max 5)
- Before merge: rebase on main → verify → bump version → update changelog
- Merge: `git merge --no-ff <branch>`
- After merge: annotated version tag → delete merged branch
- No push without explicit request

See `AGENTS.md` §Git And Release Branching and skill `sc-git` for full policy.

## Session handoff vs release closeout

- **Handoff**: stop mid-work, preserve state, no merge. Pickup: `/sc-resume` then next command.
- **Closeout**: ship/merge/finish mission. Gates: evidence, review, git, release gates all pass.

Default to handoff. Only closeout on explicit ship/release/merge intent.

## Routing Table

### Command → Agent → Subagent → Skill → Permission

**Legend**

- **agent**: The primary agent dispatched by the command (always `sc-commander`)
- **subagent**: Read-only or write-capable subagents the command may invoke
- **skill**: Skills the command loads via the `Use:` frontmatter
- **permission**: The `task.permission` and `skill.permission` entries on `sc-commander` that authorize subagents

| Command | agent | subagent (task) | skill (Use:) | permission |
|---------|-------|-----------------|--------------|------------|
| `/sc-start` | sc-commander | — | sc-mission, sc-clarify | — |
| `/sc-design` | sc-commander | sc-designer (read-only) | sc-mission, sc-clarify, sc-design | task: sc-designer → allow; skill: sc-design → allow |
| `/sc-plan` | sc-commander | sc-planner (read-only) | sc-mission, sc-clarify, sc-planning | task: sc-planner → allow; skill: sc-planning → allow |
| `/sc-build` | sc-commander | sc-coder (write), sc-tester (write) | sc-mission, sc-clarify, sc-git, sc-tdd, sc-solid, sc-verification | task: sc-coder → allow, sc-tester → allow; skill: sc-git → allow, sc-tdd → allow, sc-solid → allow, sc-verification → allow |
| `/sc-resume` | sc-commander | — | sc-mission | — |
| `/sc-review` | sc-commander | sc-reviewer (read-only), sc-designer (read-only, optional) | sc-mission, sc-verification, sc-git, sc-solid | task: sc-reviewer → allow, sc-designer → allow; skill: sc-git → allow, sc-solid → allow, sc-verification → allow |
| `/sc-quick` | sc-commander | — | sc-mission, sc-git | skill: sc-git → allow |
| `/sc-ship` | sc-commander | — | sc-mission, sc-verification, sc-git, sc-learn | skill: sc-git → allow, sc-verification → allow, sc-learn → allow |

Commander auto-triggers: sc-clarify (on ambiguity), sc-mission (session start), sc-map (before planning), sc-debug (on error/stack trace), sc-verification (after task implementation), Research (gray areas).

### Subagent → Skill → Permission

| Subagent | mode | skill.permission (allows) | bash.permission (notable) |
|----------|------|---------------------------|---------------------------|
| sc-commander | primary | `sc-*` (all skills) | R/W |
| sc-coder | write | sc-solid, sc-tdd, sc-web-service | R/W |
| sc-tester | write | sc-solid, sc-tdd, sc-verification | R/W |
| sc-designer | read-only | sc-mission, sc-design, sc-web-service | RO |
| sc-planner | read-only | sc-mission, sc-planning, sc-solid | RO |
| sc-reviewer | read-only | sc-mission, sc-git, sc-solid, sc-verification | RO |

### Agent Hierarchy

```
sc-commander (primary)
├── sc-planner    (read-only,  /sc-plan)
├── sc-designer   (read-only,  /sc-design)
├── sc-reviewer   (read-only,  /sc-review)
├── sc-coder      (write,      /sc-build)
└── sc-tester     (write,      /sc-build)
```

- `sc-commander` is the only primary agent. All slash commands dispatch to it.
- Read-only subagents (`sc-planner`, `sc-designer`, `sc-reviewer`) are invoked for planning, design, and review phases.
- Write-capable subagents (`sc-coder`, `sc-tester`) are invoked during `/sc-build` for TDD pair-programming.

### Skill References

| Skill | File | Used By |
|-------|------|---------|
| sc-mission | `.opencode/skills/sc-mission/` | All commands |
| sc-clarify | `.opencode/skills/sc-clarify/` | /sc-start, /sc-design, /sc-plan, /sc-build (auto-triggered on ambiguity) |
| sc-design | `.opencode/skills/sc-design/` | /sc-design |
| sc-planning | `.opencode/skills/sc-planning/` | /sc-plan |
| sc-git | `.opencode/skills/sc-git/` | sc-build, sc-quick, sc-review, sc-ship (auto-triggered silently within sc-build) |
| sc-verification | `.opencode/skills/sc-verification/` | /sc-build, /sc-review, /sc-ship (auto-triggered after task implementation) |
| sc-debug | `.opencode/skills/sc-debug/` | Commander auto-trigger (error/stack trace/debug request) |
| sc-creator | `.opencode/skills/sc-creator/` | Commander (skill creation workflow) |
| sc-map | `.opencode/skills/sc-map/` | Commander auto-trigger (before /sc-plan when map.json missing) |
| sc-learn | `.opencode/skills/sc-learn/` | /sc-ship, Commander (knowledge capture and migration) |
| sc-tdd | `.opencode/skills/sc-tdd/` | /sc-build, sc-tester, sc-coder (TDD red-green-refactor) |
| sc-solid | `.opencode/skills/sc-solid/` | /sc-build, /sc-review, sc-coder, sc-tester, sc-planner, sc-reviewer (SOLID, clean code, architecture) |
| sc-web-service | `.opencode/skills/sc-web-service/` | sc-coder, sc-designer |

### Permission Flow

1. A user runs a slash command → OpenCode dispatches to `sc-commander`
2. `sc-commander` loads skills listed in the command's `Use:` frontmatter
3. If the command invokes a subagent, `sc-commander` checks its own `task.permission` block:
   - `"sc-coder": allow` → authorized for /sc-build
   - `"sc-tester": allow` → authorized for /sc-build
   - `"sc-planner": allow` → authorized for /sc-plan
   - `"sc-designer": allow` → authorized for /sc-design
   - `"sc-reviewer": allow` → authorized for /sc-review
4. Subagents load their own `skill.permission` block to access skills
5. All agents check `bash.permission` for shell command authorization
