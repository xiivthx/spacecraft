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

`/sc-start` · `/sc-clarify` · `/sc-design` · `/sc-plan` · `/sc-git` · `/sc-build` · `/sc-review` · `/sc-resume` · `/sc-ship`

## CLI commands

| Command | What it does |
|---|---|
| `scripts/spacecraft init` | Create `.space/` structure |
| `scripts/spacecraft new "<title>"` | Create and select a mission |
| `scripts/spacecraft missions` | List all missions |
| `scripts/spacecraft use <selector>` | Select mission by number/id/title |
| `scripts/spacecraft resolve [--json]` | Resolve active mission |
| `scripts/spacecraft status` | Print resolved mission state |
| `scripts/spacecraft flow` | Print workflow readiness |
| `scripts/spacecraft bind-branch` | Record current branch on mission |
| `scripts/spacecraft git-info` | Print git status |
| `scripts/spacecraft git-suggest [type] [slug]` | Suggest branch/commit names |
| `scripts/spacecraft evidence "<label>" -- <cmd>` | Capture verification evidence |
| `scripts/spacecraft validate` | Validate mission artifacts |
| `scripts/spacecraft closeout-check` | Check release readiness |
| `scripts/spacecraft archive [selector]` | Compact shipped mission |
| `scripts/spacecraft set-state <state>` | Set mission state |
| `scripts/spacecraft clarify-status <status>` | Set clarification status |

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

`/sc-start` → `/sc-clarify` → `/sc-design` (if UI) → `/sc-plan` → `/sc-git` → `/sc-build` → `/sc-review` → `/sc-ship`

`/sc-build` without args loops `build → verify → checkpoint commit` until blocked.

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
