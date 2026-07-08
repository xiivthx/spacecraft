# Spacecraft Routing Table

## Command → Agent → Subagent → Skill → Permission

### Legend

- **agent**: The primary agent dispatched by the command (always `sc-commander`)
- **subagent**: Read-only or write-capable subagents the command may invoke
- **skill**: Skills the command loads via the `Use:` frontmatter
- **permission**: The `task.permission` and `skill.permission` entries on `sc-commander` that authorize subagents

---

## Commands

| Command | agent | subagent (task) | skill (Use:) | permission |
|---------|-------|-----------------|--------------|------------|
| `/sc-start` | sc-commander | — | sc-mission, sc-clarify | — |
| `/sc-clarify` | sc-commander | — | sc-mission, sc-clarify | — |
| `/sc-design` | sc-commander | sc-designer (read-only) | sc-mission, sc-clarify, sc-design | task: sc-designer → allow; skill: sc-design → allow |
| `/sc-plan` | sc-commander | sc-planner (read-only) | sc-mission, sc-clarify, sc-planning | task: sc-planner → allow; skill: sc-planning → allow |
| `/sc-git` | sc-commander | — | sc-mission, sc-git | skill: sc-git → allow |
| `/sc-work` | sc-commander | sc-coder (write), sc-tester (write) | sc-mission, sc-clarify, sc-git, sc-verification | task: sc-coder → allow, sc-tester → allow; skill: sc-git → allow, sc-verification → allow |
| `/sc-verify` | sc-commander | — | sc-verification | skill: sc-verification → allow |
| `/sc-flow` | sc-commander | sc-coder (write), sc-tester (write) | sc-mission, sc-git, sc-verification | task: sc-coder → allow, sc-tester → allow; skill: sc-git → allow, sc-verification → allow |
| `/sc-design-review` | sc-commander | sc-designer (read-only) | sc-mission, sc-design | task: sc-designer → allow; skill: sc-design → allow |
| `/sc-polish` | sc-commander | sc-designer (read-only) | sc-mission, sc-design | task: sc-designer → allow; skill: sc-design → allow |
| `/sc-review` | sc-commander | sc-reviewer (read-only), sc-designer (read-only, optional) | sc-mission, sc-verification, sc-git | task: sc-reviewer → allow, sc-designer → allow; skill: sc-git → allow, sc-verification → allow |
| `/sc-status` | sc-commander | — | — | — |
| `/sc-ship` | sc-commander | — | sc-mission, sc-verification, sc-git | skill: sc-git → allow, sc-verification → allow |

---

## Subagent → Skill → Permission

| Subagent | mode | skill.permission (allows) | bash.permission (notable) |
|----------|------|---------------------------|---------------------------|
| sc-commander | primary | `sc-*` (all skills) | `scripts/spacecraft *`, git read, rtk |
| sc-coder | write | sc-implementation | git status/diff, ls, rg |
| sc-tester | write | sc-testing, sc-verification | git status/diff, ls, rg, go/npm/pytest, rtk |
| sc-designer | read-only | sc-mission, sc-design, sc-web-service | git status/diff/log, rg, ls, find |
| sc-planner | read-only | sc-mission, sc-planning | git status/diff/log, rg, ls |
| sc-reviewer | read-only | sc-mission, sc-git, sc-verification | git status/diff/log |

---

## Agent Hierarchy

```
sc-commander (primary)
├── sc-planner    (read-only,  /sc-plan)
├── sc-designer   (read-only,  /sc-design, /sc-design-review, /sc-polish)
├── sc-reviewer   (read-only,  /sc-review)
├── sc-coder      (write,      /sc-work, /sc-flow)
└── sc-tester     (write,      /sc-work, /sc-flow)
```

- `sc-commander` is the only primary agent. All slash commands dispatch to it.
- Read-only subagents (`sc-planner`, `sc-designer`, `sc-reviewer`) are invoked for planning, design, and review phases.
- Write-capable subagents (`sc-coder`, `sc-tester`) are invoked during implementation (`/sc-work`, `/sc-flow`) for TDD pair-programming.

---

## Skill References

| Skill | File | Used By |
|-------|------|---------|
| sc-mission | `.opencode/skills/sc-mission/` | All commands (except /sc-verify, /sc-status) |
| sc-clarify | `.opencode/skills/sc-clarify/` | /sc-start, /sc-clarify, /sc-design, /sc-plan, /sc-work |
| sc-design | `.opencode/skills/sc-design/` | /sc-design, /sc-design-review, /sc-polish |
| sc-planning | `.opencode/skills/sc-planning/` | /sc-plan |
| sc-git | `.opencode/skills/sc-git/` | /sc-git, /sc-work, /sc-flow, /sc-review, /sc-ship |
| sc-verification | `.opencode/skills/sc-verification/` | /sc-verify, /sc-work, /sc-flow, /sc-review, /sc-ship |
| sc-implementation | `.opencode/skills/sc-implementation/` | sc-coder |
| sc-testing | `.opencode/skills/sc-testing/` | sc-tester |
| sc-web-service | `.opencode/skills/sc-web-service/` | sc-designer |

---

## Permission Flow

1. A user runs a slash command → OpenCode dispatches to `sc-commander`
2. `sc-commander` loads skills listed in the command's `Use:` frontmatter
3. If the command invokes a subagent, `sc-commander` checks its own `task.permission` block:
   - `"sc-coder": allow` → authorized for /sc-work, /sc-flow
   - `"sc-tester": allow` → authorized for /sc-work, /sc-flow
   - `"sc-planner": allow` → authorized for /sc-plan
   - `"sc-designer": allow` → authorized for /sc-design, /sc-design-review, /sc-polish
   - `"sc-reviewer": allow` → authorized for /sc-review
4. Subagents load their own `skill.permission` block to access skills
5. All agents check `bash.permission` for shell command authorization

*Last updated: 2026-07-08*
