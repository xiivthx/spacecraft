> Consult when: shipping a quick-lane mission, checking hard-stop gates, or handling git errors after self-review.

# Quick lane ship gates

## Required (same as normal flow)

- Version bump (or defer with rationale in decisions.md)
- Changelog update (mandatory - never defer)
- Rebase work branch on latest `main`
- Verify after rebase (run tests if available)
- `git merge --no-ff <branch>` into `main`
- Annotated tag: `git tag -a v<version> -m "v<version>"`
- Delete merged local branch
- `spacecraft archive` to compact shipped artifacts

## Skipped (fast lane only)

- No evidence.jsonl requirement
- No review.md or review.json requirement
- No `spacecraft closeout-check` (requires evidence + review gates)
- No /sc-reviewer subagent

## Ship checklist

- [ ] Version bumped (or deferred with rationale)
- [ ] Changelog updated
- [ ] Rebased on latest main
- [ ] Tests pass after rebase (if applicable)
- [ ] `git merge --no-ff` completed
- [ ] Tag created
- [ ] Branch deleted
- [ ] Mission archived

## Hard stop gates

- Resolver conflict or ambiguity
- Write attempt on `main`
- Dirty/untracked files that cannot be safely attributed to the current task
- Unsafe files or secrets before staging
- Self-review finding critical issues
- Release actions requiring `/sc-ship`
- Context is too heavy for safe continuation; give handoff instead

## Error handling

- Do not push unless explicitly asked
- Do not write on `main` - always create branch first
- If git worktree is dirty with unrelated changes, warn before staging
- If self-review finds issues, fix before ship
- Conventional Commits format: `<type>: <description>`
- After post-merge cleanup, you are on `main` - next mutation requires new branch

## Spacecraft integration

Quick missions still resolve via `spacecraft resolve` / branch binding. Archive with `spacecraft archive` after ship. Record skipped formal gates in `decisions.md`.
