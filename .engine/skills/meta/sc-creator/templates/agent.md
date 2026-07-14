---
description: Brief description of the agent's role and capability
mode: subagent  # primary | subagent
temperature: 0.1  # 0.1 for deterministic, 0.2 for balanced
permission:
  edit: deny  # deny | allow | ask
  external_directory: deny
  bash:
    "*": ask  # default permission
    "git status*": allow
    "git diff*": allow
    "git log*": allow
    "ls*": allow
    "rg *": allow
    "scripts/spacecraft *": allow
    # Add command-specific allow/deny rules here
  task:
    "*": deny
    # "sc-coder": allow  # uncomment if this agent can delegate to others
  skill:
    "*": deny
    "sc-*": allow  # uncomment to allow all spacecraft skills
    # Add specific skill allows here
---

## Role & Identity

You are a [role name].
Your primary goal is to [mission statement: what this agent does].

## Context & Guidelines

When handling tasks, you must follow these rules:

- [Rule 1: core behavior]
- [Rule 2: what to read or rely on]
- [Rule 3: how to communicate]
- [Rule 4: what to avoid]

## Constraints

Do NOT:

- [Constraint 1]
- [Constraint 2]
- [Constraint 3]

## Output Format

[Describe the expected output format or give an example.]

```
{
  "status": "ok" | "error",
  "message": "Concise status or result"
}
```
