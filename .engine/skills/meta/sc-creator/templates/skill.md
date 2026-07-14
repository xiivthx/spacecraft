---
name: your-skill-name
description: >
  One precise sentence: what this skill does + when the agent should activate it.
  Include 2-3 trigger phrases the agent can match against.
  Max 200 chars. This is the ONLY thing the agent sees before loading the full skill.
license: MIT
compatibility: opencode
metadata:
  version: 1.0
  audience: xiivthx
---

# Skill Name

One-line purpose. The agent reads this after loading.

## When to use

Activate when the user asks to:

- **[Trigger pattern]** — e.g. "create a board", "new project"
- **[Trigger pattern]** — e.g. "review this code", "check for issues"
- When the task involves X, Y, or Z

## Workflow

Use this exact sequence unless the user specifies otherwise:

1. **Step one** — What to do first.
   `exact command or reference`
2. **Step two** — What to do next.
   Expected output / format
3. **Step three** — Final step.
   - Sub-bullets for variations
   - Edge cases to check

## Rules

- **Must**: thing the agent must always do
- **Must not**: thing the agent must never do (stronger than "avoid")
- **Ask before**: things requiring user approval

## Out of scope

This skill does NOT handle:

- Task A — use [other-skill] instead
- Task B — ask the user for direction
- Task C — out of scope entirely

## Output format

```
Expected output shape.
[type] description
[type] description
```

## Checklist

Before claiming done:

- [ ] Criterion 1
- [ ] Criterion 2
- [ ] Criterion 3
- (Add/remove as needed)

---

## References

Load details on demand from:

- `references/rules.md` — detailed conventions
- `references/examples.md` — worked examples
- `scripts/` — executables (if any)
