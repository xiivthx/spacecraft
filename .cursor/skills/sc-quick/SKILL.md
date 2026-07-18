---
name: sc-quick
description: "Runs the fast lane when a small change can be safely branched, verified, committed, and shipped without a full mission flow."
disable-model-invocation: true
---

Use sc-mission and sc-git.
Start a quick mission for: $ARGUMENTS

## When to use

Use `/sc-quick` for small, straightforward changes that don't need detailed spec, task planning, TDD cycles, or formal review. Examples:
- Prompt tweaks, config changes, documentation
- Small fixes with obvious scope
- Personal tooling adjustments
- Anything where full flow overhead exceeds the change itself

Do NOT use for: new features with unknown scope, multi-file refactors, UI changes, API integrations, or anything needing design review.

## Pre-flight Checks

Resolve the mission. Block if unsafe.

Run:
```
spacecraft git-info
```

### Git safety

- If the workspace is not a git worktree, stop before editing product files unless the user has explicitly accepted no-git implementation risk for this mission in decisions.md.
- If git exists, inspect dirty state before edits. Work with user changes; do not revert unrelated changes.
- Use sc-git for branch naming, branch creation, checkpoint commits, .gitignore hygiene, and staging safety checks.

## Workflow

### 1. Mission stub

Run:
```
spacecraft new "$ARGUMENTS"
```

Then immediately:
- Set clarification to clear: `spacecraft clarify-status clear`
- Record in decisions.md: "Fast lane - skipping spec.md, plan.json, evidence, and formal review."
- Do NOT write spec.md or plan.json stubs
- Set state to planned: `spacecraft set-state planned`

### 2. Branch

Create a non-main work branch from latest `main` using sc-git naming:
```
spacecraft git-suggest feat
```

Create and switch to the branch. Bind it:
```
spacecraft bind-branch
```

### 3. Commit freely

- No TDD cycle required - implement directly
- No acceptance checks required
- Conventional Commits still expected
- Checkpoint commits allowed on work branch
- `.gitignore` must stay current; no secrets staged

### 4. Automated verification (required)

Before self-review, run the project's test suite and capture output as evidence:

```bash
# Go:    spacecraft evidence "quick lane tests" -- make test
# Node:  spacecraft evidence "quick lane tests" -- npm test
# Other: use the appropriate test command
```

If no test suite exists, skip this step but note it in decisions.md as a limitation.

### 5. Fast self-review

Before ship, commander performs a lightweight self-review (commander: `.cursor/rules/000-spacecraft.mdc`; Quick Lane: `.cursor/rules/200-workflow.mdc`):
- Inspect `git diff` - secrets, debug code, unrelated edits, dead code, noisy formatting
- Review test output from evidence - all tests must pass
- Verify the change does what was intended
- No subagent reviewer; no review.md or review.json required

If self-review finds issues, fix them and recommit.

### 6. Ship

When ready to ship, do release closeout with streamlined gates. Full required/skipped lists, ship checklist, hard-stop gates, and error handling: `references/ship-gates.md`.

Produce a summary: mission id, what changed, git branch/merge info, suggested commit message, known limitations, next step.

## Research auto-trigger

When quick-lane changes touch unfamiliar tooling, configuration, or dependency APIs, use sc-search (WebSearch/WebFetch) for `"<topic>"` before committing. Fast lane is not skip-research lane.

End with session advice. Recommend new session after shipped mission.

## References

- `references/ship-gates.md` - ship required/skipped gates, checklist, hard stops, error handling
