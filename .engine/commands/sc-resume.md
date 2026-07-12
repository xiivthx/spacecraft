---
description: Resume active mission after session handoff
agent: sc-commander
---

Use sc-mission.
Resolve the mission. If no mission resolves, say so and stop — do not create a new mission.

## Pre-flight Checks

Run:
```
scripts/spacecraft resolve --json
```

If resolver safety is not `safe` or no mission is selected, state the issue and stop. Run `scripts/spacecraft missions` and `scripts/spacecraft use <number|id|title>` to resolve.

Run:
```
scripts/spacecraft git-info
```

If git is dirty and the mission state is not `draft`, flag it prominently in the resume output.

## Live state

Mission status:
!`scripts/spacecraft status`

Workflow:
!`scripts/spacecraft flow`

Git:
!`scripts/spacecraft git-info`

Last commit:
!`git log -1 --oneline --no-decorate 2>/dev/null || echo "(no commits)"`

Dirty files (if any):
!`git diff --stat 2>/dev/null; git diff --cached --stat 2>/dev/null || echo "(clean or not a repo)"`

Last evidence entry:
!`mid=$(scripts/spacecraft resolve --json 2>/dev/null | sed -n 's/.*"currentMissionId":"\([^"]*\)".*/\1/p'); [ -n "$mid" ] && [ -f ".space/missions/$mid/evidence.jsonl" ] && tail -1 ".space/missions/$mid/evidence.jsonl" || echo "(no evidence)"`

## Handoff resume

Based on the live state above, present a concise handoff resume:

1. **Mission**: ID, title, state, how resolved (branch/session/explicit)
2. **Git**: branch, HEAD short sha, clean or dirty (list dirty files if ≤5)
3. **Clarification**: status + blocking question count
4. **Progress**: tasks completed/total, next task (ID + title)
5. **Evidence**: count, last entry (label + date + exit code)
6. **Review**: status + unresolved finding count
7. **Blockers**: list if any
8. **Next action**: if the workflow `Next:` field is a slash command (starts with `/sc-`), quote it exactly as the pickup command. If it is a parenthesized status like `(clarify)` or `(shipped)`, explain what the Commander should do: for `(clarify)` — sc-clarify skill will auto-trigger, just tell the user to continue the session; for `(shipped)` — mission is complete, nothing to do.
9. **Session advice**: continue this chat or start new session, with brief reason

## Constraints

- Do NOT start implementing, designing, planning, or mutating anything.
- Do NOT ask the user what they want to do — the resume output IS the answer.
- This command is strictly read-only.
- If git is dirty and the state is not `draft`, flag it prominently.
- If no mission resolves, say: "No active mission. Start one with `/sc-start <title>`."

## Research auto-trigger

sc-resume is read-only — no research trigger needed during resume. Research decisions belong to the commands that follow (sc-plan, sc-build, etc.).

## Hard stop gates

- No mission resolves
- Resolver conflict or ambiguity
- Git dirty with non-draft state (warn prominently before proceeding with mutations)

## Error handling

- No mission resolves → stop and tell user: "No active mission. Start one with `/sc-start <title>`."
- Resolver conflict or ambiguity → display candidates and tell user to run `scripts/spacecraft missions` then `scripts/spacecraft use <number|id|title>`
- Git not available → still display mission state but flag git as unavailable
- Dirty workspace with non-draft mission → warn prominently in resume output
- Cannot read evidence → display "(no evidence)" and continue

End with the **next action** line only — no pleasantries, no filler.
