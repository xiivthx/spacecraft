---
name: sc-quick
description: "No-mission fast lane for manual edits/fixes/docs: branch → verify → commit → ship in one pass. Invoke as /sc-quick."
disable-model-invocation: true
---

# sc-quick

## Goal

Ship small, obvious changes **without a mission**: no `spacecraft new`, no spec/plan/TDD/formal review/closeout-check. Still use git safety, Conventional Commits, and self-review. Invoking `/sc-quick` runs the full lane through local merge + tag - no separate ship step.

## Output

Change on `main` via `--no-ff` merge + annotated tag (or blocked with exact missing actions). No mission archive. No `set-state shipped`.

## Good / Bad

- Good: docs tweak, diagram fix, config typo, prompt edit, one-file bug with obvious fix, one-off polish/optimize/refactor/adjust, “just commit my manual edits”
- Bad: new product feature with unknown scope, multi-module refactor, UI needing design HIL, API contracts, batched polish (≥3 `.space/polish-backlog.md` items or explicit batch discuss), anything needing evidence/review → use `/sc-discuss` then `/sc-run`

## Verify

Self-review of `git diff` + project tests when code changed (or note docs-only skip). After ship: on `main`, tag exists, work branch deleted.

## Arguments

```
/sc-quick
/sc-quick <short description>
```

`$ARGUMENTS` = optional title/description for the work branch / summary.

## When to use

Use when **no mission is appropriate**:

- Manual edits the human already made (or asks you to make)
- Docs / diagrams / changelog-only / config / prompts
- Tiny fixes where mission ceremony costs more than the change
- One-off tiny polish / optimize / refactor / adjust (no mission ceremony)

### Polish: one-off vs batch

- **One-off** → stay on this lane.
- **Batched** (≥3 items in `.space/polish-backlog.md` OR explicit batch discuss) → stop; hand off to `/sc-discuss` then `/sc-run`. Mission path: backlog + evidence + baseline do-not-break; max ≤5 files/batch; never unbounded improve.

Do **not** invent a mission stub. If `spacecraft resolve` finds a mission, prefer that mission’s `/sc-run` / `/sc-ship` path unless the user explicitly wants quick lane anyway.

## Pre-flight

```
git status -sb
git branch --show-current
git rev-parse --abbrev-ref HEAD
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

If already on a quick work branch with the intended commits, skip recreate and continue from verify → ship.

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

### 6. Ship (automatic)

Invoking `/sc-quick` **is** the ship authorization for this lane. After verify + commits succeed, merge and tag in the same run - do **not** wait for a second “ship” / “merge” message.

AUTH: quote the user’s `/sc-quick` invocation (or equivalent quick-lane request).

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

**Do not push** unless the user explicitly asks. Push needs AUTH + `SPACECRAFT_SHIP=1 SPACECRAFT_QUICK=1` **and** human approval in the ship hook ask prompt.

Full gate lists: `references/ship-gates.md`.

## Hard stops

- Scope too large / unclear → hand off to `/sc-discuss`
- Batched polish (≥3 `.space/polish-backlog.md` items or explicit batch discuss) → `/sc-discuss` then `/sc-run` (≤5 files/batch; never unbounded improve)
- Write attempt on `main` for product commits
- Secrets or unsafe staging
- Self-review finds critical issues
- Mission already in flight for the same work → prefer `/sc-ship`

## Errors

- Never use `closeout-check` as a blocker on this lane (hook skips it when `SPACECRAFT_QUICK=1`)
- Never create `.space/missions/…` for quick work unless the user asks for a mission
- After merge you are on `main` — next mutation needs a new branch

## Research

Unfamiliar tooling/APIs: use sc-search before committing. Fast ≠ skip-research.

## Summary format

After successful merge/tag, if `spacecraft map current` succeeds and `spacecraft map next <rid>` is not `All missions complete.`, handoff: **Next: new session → /sc-discuss <id>**. Otherwise `Next: push? / done`. When a resolved mission exists and the handoff is useful, set or update optional `mission.json` `pickup` (`phase`, `next` one-liner, `updatedAt`) so `spacecraft status` / session-start shows Pickup. Not a ship gate (quick has no closeout).

```
Lane: quick (no mission)
Branch: <type>/<title> → merged --no-ff → deleted
Commits: …
Tag: vX.Y.Z
Files: …
Verify: <tests or docs-only>
AUTH: "<quoted /sc-quick invocation>"
Next: new session → /sc-discuss <id>   # when roadmap has next
# or: Next: push? / done               # no current roadmap / all complete
```
