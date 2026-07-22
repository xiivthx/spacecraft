> Consult when: shipping a `/sc-quick` no-mission change, checking hard stops, or handling git errors after self-review.

# Quick lane ship gates (no mission)

## Required

- Work branch `<type>/<short-title>` (no mission id)
- Self-review of diff completed
- Tests run when code changed (or docs-only noted)
- CHANGELOG.md updated (next patch heading unless user sets otherwise)
- Rebase on latest `main` (warn if fork point mismatch)
- Re-verify after rebase when code changed
- `SPACECRAFT_SHIP=1 SPACECRAFT_QUICK=1 git merge --no-ff <branch>`
- Annotated tag `vX.Y.Z`
- Delete merged local branch
- Unset `SPACECRAFT_SHIP` and `SPACECRAFT_QUICK`

## Skipped (no-mission quick only)

- No `spacecraft new` / mission artifacts
- No `spec.md` / `plan.json` / `evidence.jsonl` / `review.json`
- No `spacecraft validate --strict`
- No `spacecraft closeout-check` (hook skips when `SPACECRAFT_QUICK=1`)
- No `/sc-ship` mission closeout / archive / `set-state shipped`
- No reviewer / judge subagents

## Ship checklist

- [ ] Branch not `main`
- [ ] Only related paths staged
- [ ] Conventional Commits (≤3 commits; changelog separate)
- [ ] CHANGELOG committed (heading matches tag)
- [ ] Rebased on main
- [ ] Tests green or docs-only noted
- [ ] Merge `--no-ff` with both env flags
- [ ] Tag created
- [ ] Branch deleted
- [ ] Env unset

## Hard stops

- Scope needs discuss/run
- Dirty unrelated files staged
- Secrets in diff
- Ship without user authorization
- `SPACECRAFT_SHIP=1` without `SPACECRAFT_QUICK=1` on a no-mission branch (closeout will fail — set both)

## Error handling

- No push unless asked
- No product commits on `main`
- If closeout runs by mistake (QUICK unset), unset and retry with both flags
- After merge, next edit needs a new branch

## Hook contract

```
SPACECRAFT_SHIP=1 SPACECRAFT_QUICK=1 git merge|push|tag …
```

- `SPACECRAFT_SHIP=1` alone → still runs `closeout-check` (mission ship)
- Both flags → allow ship git without closeout (quick lane)
- Neither → deny merge/push/tag
