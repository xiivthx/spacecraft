---
name: sc-quick
description: "No-mission fast lane for manual edits/fixes/docs: branch → verify → commit → ship. Invoke as /sc-quick."
disable-model-invocation: true
---

# sc-quick

## Goal

Ship small, obvious changes **without a mission**: no `spacecraft new`, no spec/plan/TDD/formal review/closeout-check. Still use git safety, Conventional Commits, self-review, and an explicit ship step.

## Output

Change on `main` via `--no-ff` merge + annotated tag (or blocked with exact missing actions). No mission archive. No `set-state shipped`.

## Good / Bad

- Good: docs tweak, diagram fix, config typo, prompt edit, one-file bug with obvious fix, “just commit my manual edits”
- Bad: new product feature with unknown scope, multi-module refactor, UI needing design HIL, API contracts, anything needing evidence/review → use `/sc-discuss` then `/sc-run`

## Verify

Self-review of `git diff` + project tests when code changed (or note docs-only skip). After ship: on `main`, tag exists, work branch deleted.

## Arguments

```
/sc-quick
/sc-quick <short description>
/sc-quick ship
```

`$ARGUMENTS` = optional title/description, or `ship` to finish an in-progress quick branch.

## When to use

Use when **no mission is appropriate**:

- Manual edits the human already made (or asks you to make)
- Docs / diagrams / changelog-only / config / prompts
- Tiny fixes where mission ceremony costs more than the change

Do **not** invent a mission stub. If `spacecraft resolve` finds a mission, prefer that mission’s `/sc-run` / `/sc-ship` path unless the user explicitly wants quick lane anyway.

## Pre-flight

```
spacecraft git-info
git status -sb
git branch --show-current
```

- Not a git worktree → stop (unless user accepts no-git risk in chat).
- On `main` with mutations pending → create work branch before editing/committing.
- Unrelated dirty files → do not stage; warn.
- Use sc-git for Conventional Commits, `.gitignore`, staging hygiene.
- **Skip** `spacecraft resolve` / `use` / `new` / `closeout-check` / `validate --strict`.

## Workflow

### 1. Branch (no mission id)

From latest `main`:

```
git checkout main
git pull   # if remote exists
git checkout -b <type>/<short-title>
```

`<type>`: `docs` | `fix` | `chore` | `style` | rarely `feat` (only if truly tiny).  
**No** `M…` mission id segment. Example: `docs/gpio-mapping-real-pins`.

### 2. Edit / adopt changes

- Implement the requested edit, **or** take the user’s existing dirty files.
- INTENT before behavior-changing code (same as Mission).
- Prefer `git add <paths>` for only this change.

### 3. Verify

| Change kind | Action |
|-------------|--------|
| Code / tests | Run project test command; fix failures before ship |
| Docs / HTML / markdown only | Visual or content check; note “docs-only — no suite” in ship summary |
| Mixed | Run suite for touched packages |

No `evidence.jsonl` required. Optional: paste key command output in the ship summary.

### 4. Fast self-review

Inspect `git diff` / `git diff --cached`:

- No secrets, credentials, or `.env`
- No unrelated files / debug noise
- Matches user request
- TWINS search after defect fixes

Fix and recommit if needed.

### 5. Commits

Target **1–2** commits (max 3):

1. Implementation: `docs:` / `fix:` / `chore:` …
2. **Separate** release notes: `chore: release notes for vX.Y.Z` (CHANGELOG; tag = next patch unless user sets otherwise)

Body: `-` bullets, lowercase start, no mission ids.

### 6. Ship (explicit only)

Ship only when the user says **ship** / **merge** / `/sc-quick ship` (or equivalent). AUTH with a **quoted** user phrase.

```
git rebase main          # or origin/main when available
# re-run tests if code changed
# ensure CHANGELOG commit exists

SPACECRAFT_SHIP=1 SPACECRAFT_QUICK=1 git merge --no-ff <type>/<short-title> -m "merge: <short title>"
SPACECRAFT_SHIP=1 SPACECRAFT_QUICK=1 git tag -a vX.Y.Z -m "vX.Y.Z"
git branch -d <type>/<short-title>
unset SPACECRAFT_SHIP SPACECRAFT_QUICK
```

Bump policy: docs/chore → next **patch** tag; tiny fix → patch; do not invent major/minor without user ask.

**Do not push** unless the user explicitly asks. Push still needs AUTH + `SPACECRAFT_SHIP=1 SPACECRAFT_QUICK=1`.

Full gate lists: `references/ship-gates.md`.

## Hard stops

- Scope too large / unclear → hand off to `/sc-discuss`
- Write attempt on `main` for product commits
- Secrets or unsafe staging
- Self-review finds critical issues
- User did not authorize ship
- Mission already in flight for the same work → prefer `/sc-ship`

## Errors

- Never use `closeout-check` as a blocker on this lane (hook skips it when `SPACECRAFT_QUICK=1`)
- Never create `.space/missions/…` for quick work unless the user asks for a mission
- After merge you are on `main` — next mutation needs a new branch

## Research

Unfamiliar tooling/APIs: use sc-search before committing. Fast ≠ skip-research.

## Summary format

```
Lane: quick (no mission)
Branch: <type>/<title> → merged --no-ff → deleted
Commits: …
Tag: vX.Y.Z
Files: …
Verify: <tests or docs-only>
AUTH: "<quoted user ship phrase>"
Next: push? / done
```
