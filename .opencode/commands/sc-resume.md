---
description: Resume active mission after session handoff
agent: sc-commander
---

Use sc-mission.
Resolve the mission. If no mission resolves, say so and stop — do not create a new mission.

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
!`scripts/spacecraft resolve --json 2>/dev/null | python3 -c "
import json,sys,os
d=json.load(sys.stdin)
mid=d.get('selected',{}).get('id') or d.get('currentMissionId')
if mid:
    p=os.path.join('.space','missions',mid,'evidence.jsonl')
    if os.path.exists(p):
        lines=[l.strip() for l in open(p) if l.strip()]
        if lines:
            print(lines[-1])
        else:
            print('(no evidence)')
    else:
        print('(no evidence file)')
else:
    print('(no mission)')
" 2>/dev/null || echo "(could not read evidence)"`

## Handoff resume

Based on the live state above, present a concise handoff resume:

1. **Mission**: ID, title, state, how resolved (branch/session/explicit)
2. **Git**: branch, HEAD short sha, clean or dirty (list dirty files if ≤5)
3. **Clarification**: status + blocking question count
4. **Progress**: tasks completed/total, next task (ID + title)
5. **Evidence**: count, last entry (label + date + exit code)
6. **Review**: status + unresolved finding count
7. **Blockers**: list if any
8. **Next action**: exact slash command to pick up work (e.g. `/sc-build T03`, `/sc-review`, `/sc-plan`)
9. **Session advice**: continue this chat or start new session, with brief reason

## Constraints

- Do NOT start implementing, designing, planning, or mutating anything.
- Do NOT ask the user what they want to do — the resume output IS the answer.
- This command is strictly read-only.
- If git is dirty and the state is not `draft`, flag it prominently.
- If no mission resolves, say: "No active mission. Start one with `/sc-start <title>`."

End with the **next action** line only — no pleasantries, no filler.
