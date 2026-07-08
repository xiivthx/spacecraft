---
name: sc-planning
description: Convert a mission spec into a small executable plan with verifiable tasks
license: MIT
compatibility: opencode
metadata:
  version: 1.0
  audience: xiivthx
---

# sc-planning

Convert a mission spec into a small executable plan with verifiable tasks.

## When to use

Activate when the user asks to:

- Plan next steps from a spec
- Create or update a `plan.json`
- Break a mission into tasks
- Scope work before implementation

## Workflow

Use this exact sequence unless the user specifies otherwise:

1. **Read inputs** — Check the mission's `spec.md`, `questions.md`, and `decisions.md` before producing `plan.json`.
2. **Identify tasks** — Break the spec into ≤7 small, verifiable tasks. Each task needs `id`, `title`, `status`, `files`, `acceptance`, `verify`, and `evidence`.
3. **Write plan** — Produce or update `plan.json` with status values: `pending`, `in-progress`, `done`, `blocked`.
4. **Verify** — Ensure the plan is complete before the mission moves to implementation.

## Rules

- **Must**: Create `plan.json`-ready output.
- **Must**: Before producing `plan.json`, check `questions.md` and `decisions.md`.
- **Must not**: Fill gray areas with hidden assumptions.
- **Must**: Record low-risk assumptions explicitly.
- **Must**: Blocking product/design decisions require sc-clarify.
- **Must**: Keep tasks ≤ 7.
- **Must**: Each task needs `id`, `title`, `status`, `files`, `acceptance`, `verify`, and `evidence`.
- **Must**: Use status values: `pending`, `in-progress`, `done`, `blocked`.
- **Must not**: Use vague tasks.
- **Must**: Include exact files when known.
- **Must**: Prefer focused tests and build checks.
- **Must not**: Create broad architecture plans unless the mission explicitly requires it.

## Out of scope

This skill does NOT handle:

- Design or UI work — use sc-design instead
- Implementation — use sc-work or sc-coder
- Verification — use sc-verification

## Output format

```
{
  "planName": "<short-name>",
  "missionId": "<mission-id>",
  "tasks": [
    {
      "id": "<task-id>",
      "title": "<short description>",
      "status": "pending|in-progress|done|blocked",
      "files": ["<path1>", "<path2>"],
      "acceptance": ["<check1>", "<check2>"],
      "verify": "<verification-command-or-description>",
      "evidence": "<evidence-capture-command>"
    }
  ]
}
```

## Checklist

Before claiming the plan is ready:

- [ ] Read `spec.md`, `questions.md`, and `decisions.md`
- [ ] Plan has ≤7 tasks
- [ ] Each task has `id`, `title`, `status`, `files`, `acceptance`, `verify`, `evidence`
- [ ] No vague or unverifiable tasks
- [ ] Assumptions recorded if low-risk

---

## References

- `.space/skill-template.md` — section template reference
- `scripts/spacecraft plan --help` — plan subcommand reference (if available)
