# OpenCode Command Template

Command files go in `.opencode/commands/<name>.md`. Run with `/<name>` in TUI.

Based on real-world usage from top OpenCode repos (Cluster444/agentic, CommandKenobi, opencode-commands, opencode-template).

---

## Anatomy

Every command has two parts:

1. **Frontmatter** (YAML `---`) — metadata
2. **Template body** (plain text/Markdown) — the prompt sent to the LLM

```yaml
---
description: One-line shown in TUI autocomplete
agent: build          # optional, defaults to current agent
model: provider/model # optional, overrides model
subtask: true         # optional, forces subagent invocation
---
```

---

## Minimal templates

### Simple (no args)

```yaml
---
description: Run tests with coverage
agent: build
---
Run the full test suite with coverage report and show any failures.
Focus on the failing tests and suggest fixes.
```

### With arguments (`$ARGUMENTS` or `$1`, `$2`)

```yaml
---
description: Create a new React component with TypeScript
agent: build
---
Create a new React component named $ARGUMENTS with TypeScript support.
Include proper typing and basic structure.
```

```yaml
---
description: Create a file with content in a directory
---
Create a file named $1 in the directory $2 with the following content:
$3
```

### With shell output injection

```yaml
---
description: Analyze test coverage
---
Here are the current test results:
!`npm test`

Based on these results, suggest improvements to increase coverage.
```

```yaml
---
description: Review recent git changes
---
Recent git commits:
!`git log --oneline -10`

Review these changes and suggest any improvements.
```

### With file references

```yaml
---
description: Review component for performance issues
---
Review the component in @src/components/Button.tsx.
Check for performance issues and suggest improvements.
```

---

## Production templates (from real repos)

### Git commit with user approval

```yaml
---
description: Create a conventional commit from all changes with user approval
agent: build
---
You are executing the `/commit-all` command. Follow this workflow:

## Pre-flight Checks

1. **Verify Git Configuration**
   - Run: `git config user.name` — if empty, abort
   - Run: `git config user.email` — if empty, abort

2. **Check Current Branch**
   - Run: `git branch --show-current`
   - If `main`, `master`, `develop`, or `release/*`:
     Warn user and offer: Commit Anyway / Create Branch / Cancel

## Workflow

### 1. Read All Changes (parallel)
- `git status`
- `git diff --cached --stat`
- `git diff --cached`
- `git diff --stat`
- `git diff`
- `git log --oneline -5`

### 2. Analyze Changes
- Which files are staged vs unstaged
- Full diff
- Recent patterns

### 3. Generate Commit Message
- Format: `type(scope): subject`
- Types: feat, fix, docs, style, refactor, test, chore, perf
- Subject: under 72 chars, imperative mood, no period

### 4. Ask user (use question tool)
- Present: Accept / Suggest Again / Cancel

### 5. Execute
- On accept: `git add . && git commit -m "MESSAGE"`

## Error Handling
- Nothing to commit → message
- Hook failed → show error + suggest fix
- Staging failed → check permissions
```

### Plan generation (ticket → plan)

```yaml
---
description: Create implementation plan from ticket + research
agent: plan
subtask: true
---
# Implementation Plan

## Step 1: Context & Research
1. Read ALL mentioned files fully (no limit/offset)
2. Spawn parallel research tasks:
   - codebase-locator — find relevant files
   - codebase-analyzer — understand current implementation
   - thoughts-locator — find existing docs
3. Read ALL files identified by research
4. Present informed analysis + focused questions

## Step 2: Plan Structure
Get user alignment on:
- Phases and ordering
- Approach/design options

## Step 3: Write Plan
Write to `thoughts/plans/<name>.md` with sections:
- Overview
- Current State Analysis
- What We're NOT Doing
- Per-phase: Changes Required + Success Criteria
- Testing Strategy
- Performance Considerations

## Rules
- Be skeptical — verify with code, don't assume
- Interactive — get buy-in at each step
- No open questions in final plan
- Separate automated vs manual verification
```

### PR review

```yaml
---
description: Review pull request changes and suggest improvements
agent: plan
---
You are reviewing changes for a pull request.

## Process

1. **Read Changes**
   - `git diff main...HEAD`
   - `git log main...HEAD --oneline`
   - List changed files

2. **Analyze**
   - Logic correctness
   - Edge case handling
   - Error handling
   - Test coverage
   - Performance impact
   - Security concerns
   - Code style consistency
   - Documentation completeness

3. **Report**
   ```
   ## Summary
   [overall assessment]

   ## What's good
   - [specific strength]

   ## Issues
   - [P0] Critical: [description + suggestion]
   - [P1] Major: [description + suggestion]
   - [P2] Minor: [description + suggestion]

   ## Questions
   - [unresolved question]
   ```
```

---

## Best practices from real usage

| Practice | Why | Source |
|---|---|---|
| `description` is required | Shows in TUI autocomplete | Docs |
| Use `$ARGUMENTS` not inline | Cleaner positional handling | Docs |
| Pre-flight checks first | Catch issues early | CommandKenobi, Agentic |
| Use question tool for approval | User stays in control | CommandKenobi |
| Spawn sub-tasks for research | Parallel context gathering | Agentic (plan, research) |
| Separate automated vs manual verification | Clear handoff | Agentic |
| Prefer `agent: build` for mutating commands | Safety | All real repos |
| Use `subtask: true` for long commands | Don't pollute main context | opencode-commands |
| Include error handling section | Robustness | CommandKenobi |
| Write output to files not memory | Persistence, reviewability | Agentic |

## Writing guidelines

- **Be specific**: refer to exact file paths, commands, patterns
- **Be structured**: use headings, lists, code blocks
- **Be interactive**: break into steps with user checkpoints
- **Be safe**: check branch, check git state, ask before destructive ops
- **Be verifiable**: include success criteria the LLM can check

## References

- Docs: https://opencode.ai/docs/commands/
- Agentic workflow: https://github.com/Cluster444/agentic
- CommandKenobi: https://github.com/obiwancenobi/CommandKenobi
- opencode-commands: https://github.com/nguyenngothuong/opencode-commands
- opencode-template: https://github.com/julianromli/opencode-template
- oh-my-opencode: https://github.com/code-yeongyu/oh-my-opencode
