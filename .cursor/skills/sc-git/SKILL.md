---
name: sc-git
description: "Enforce git safety, branching, Conventional Commits, no-ff merge, versioning, and release tags. Activate on commit, branch, merge, release, or ship."
---

# sc-git

Enforce git safety, branching, hygiene, Conventional Commits, verification before main, no-ff merge, version bump, changelog/spec update, and release tag policy.

## When to use

Activate when the user asks to:

- Set up or check git branch, commit, or merge
- Run release closeout or prepare for ship
- Verify git hygiene before staging, committing, or merging
- Check git safety before mutating work

## Workflow

Use this exact sequence unless the user specifies otherwise:

1. **Resolve mission** - Before git work on a mission, resolve via `spacecraft resolve`. On conflict or ambiguity, use `spacecraft use <selector>`. Skip for `/sc-quick` (no mission).
2. **Check git state** - Run `git status`, `git rev-parse`, and related plain git checks before committing/merging/releasing.
3. **Branch** - Create a non-main work branch from latest `main` before mutating. Never write product changes on `main`.
4. **Commit (AFK checkpoints)** - During `/sc-run` build, auto-commit one Conventional Commit per plan task after that task's acceptances are done; also after combine/refactor and material fixes. Details: `references/release-policy.md` §Checkpoint commits. These are WIP on the work branch only.
5. **Squash before ship** - On `/sc-ship`, squash/fixup checkpoints into 1–3 logical Conventional Commits (max 5) before merge. See sc-ship. `/sc-quick` keeps 1–3 commits (no mission squash ceremony).
6. **Commit (release notes)** - Add changelog update as a **separate commit** in the work branch before merge (type `chore:` or `docs:`). Never defer after merge.
7. **Verify** - After latest rebase, reverify. For missions, run `spacecraft closeout-check` before claiming release readiness. `/sc-quick` skips closeout (hook uses `SPACECRAFT_QUICK=1`).
8. **Rebase & Merge** - Before merge, rebase on latest `main`. Merge with `--no-ff` only. Mission: `SPACECRAFT_SHIP=1`. Quick: `SPACECRAFT_SHIP=1 SPACECRAFT_QUICK=1`.
9. **Tag** - After no-ff merge, create annotated tag. Bump policy and tagging tables: `references/release-policy.md`. Do not tag before merge exists.

## Rules

### General

- **Must**: Treat git as the rollback and release boundary for mutating work.
- **May**: Discovery, clarification, design, planning, and review may run without git.
- **Must**: Before git work on a mission, resolve via `spacecraft resolve`. On conflict or ambiguity, use `spacecraft use <number|id|title>` or `SPACECRAFT_MISSION`. Skip resolve for `/sc-quick`.
- **Must**: Before committing/merging/releasing, inspect git state with plain git (`git status`, `git rev-parse`, etc.).
- **Must**: Before the first write, confirm the path is under `git rev-parse --show-toplevel` or a path the user named.
- **Must not**: Write product changes on `main`. If on `main` when mutation is requested, create a work branch.
- **Must not**: Ad-hoc `git init` (agents inventing repo setup). **May**: spacecraft `ensureProjectReady`, `spacecraft init`, and bootstrap / `install-cursor` run `git init` when the project is not yet a git repo.
- **Must not**: Auto-create worktrees, rebase, merge, tag, or push unless asked.
- **Must**: Before outward actions (push, deploy, publish, send), state `AUTH:` with a **quoted** user authorization from the conversation. Mission local ship still requires `/sc-ship` + hooks + `SPACECRAFT_SHIP=1`. `/sc-quick` authorizes local merge/tag in that lane (hooks + `SPACECRAFT_SHIP=1 SPACECRAFT_QUICK=1`). Push still needs separate AUTH **and** human approval in the ship hook (`permission: ask`). AUTH does not bypass ship hooks. Force-push is hook-denied.

### Branching

- **Must**: One branch = one feature/fix/scoped change. Branch from latest `main`.
- **Must**: Pattern: `<type>/<id>/<title>` for missions - e.g. `feat/m07fp1l7z/go-rewrite`. `/sc-quick` uses `<type>/<title>` with no mission id.
- **Must**: `release/v<major>.<minor>.<patch>` only for release prep.
- **Must**: If a branch needs >5 final commits, split the feature first. Before adding a 4th commit, if `git log main..HEAD` already mixes unrelated chores, stop and split - do not pile on.

### Commits (Conventional Commits)

- **Must**: Target 1–3 **final** commits per branch after ship squash, max 5 unless justified in decisions.md.
- **Must**: Separate version bump + changelog update into its own `chore:` or `docs:` commit - do not bundle with implementation commit.
- **Must**: During `/sc-run` AFK build, auto-create checkpoint commits (see `references/release-policy.md` §Checkpoint commits).
- **Must**: Before `/sc-ship` merge, squash/fixup checkpoints into logical Conventional Commits (≤5). Do not squash unrelated changes.
- **Must not**: Stage unrelated user changes. Prefer `git add <specific-files>`.
- **Must**: Subject: `<type>: <description>` - imperative, lowercase, ~72 chars.
- **Must**: Body: `-` bullets, lowercase start, no period end. Describe what changed; evidence lives in `evidence.jsonl`, not the commit.
- **Must not**: Put mission ids (`M…`) in commit subjects or bodies (including merge commits and AFK checkpoints). Keep ids on the work branch name / mission artifacts only.
- **Must**: Types: `feat`, `fix`, `docs`, `refactor`, `test`, `build`, `ci`, `chore`, `perf`, `style`, `revert`.
- **Must**: Breaking: `!` after type or `BREAKING CHANGE:` footer.

### Ignore Hygiene

- **Must**: Before staging, inspect untracked/modified files. Update `.gitignore` for new build outputs, caches, logs, env files, deps, or machine files.
- **Must**: Use stack-aware patterns (Node, Rust, Python, etc.) when those tools are present.
- **Must**: Keep source, tests, migrations, lockfiles, configs, and release notes tracked.
- **Must**: If unsure about a file, ask before staging. Before release, verify no secrets are tracked.

### Rebase & Merge (ship gates)

- **Must**: Before merge, rebase on latest `main`. Reverify after rebase.
- **Must**: Identify fork point with `git log --oneline main..HEAD | head -1` before rebase. If `main` has advanced beyond the expected base, warn: "Rebase target mismatch: main HEAD differs from fork point. Confirm correct base before rebase."
- **Must**: Immediately before merge, strip the mission id from the work branch when it matches `<type>/<id>/<title>` and `<id>` is a mission id (`M…`): `git branch -m <type>/<title>`. Finish closeout/evidence while the id is still in the name. See sc-ship §Merge.
- **Must**: Merge only with `--no-ff`. No fast-forward or squash-merge into `main`.
- **Must not**: Rebase/rewrite `main`. Resolve conflicts on work branch, then reverify.
- **Must**: After merge, delete local branch unless asked to keep.
- **Must**: Mission merge/tag uses `SPACECRAFT_SHIP=1` (and `SPACECRAFT_QUICK=1` for `/sc-quick`). Hooks enforce closeout unless quick.

### Release / tag (summary)

- **Must**: After every no-ff merge to `main`, create an annotated tag `v<major>.<minor>.<patch>` (see `references/release-policy.md` §Tagging).
- **Must**: For missions, `spacecraft closeout-check` before claiming readiness (quick skips). Full closeout field bars: `references/release-policy.md` §Release Prep.
- **Must**: After merge: tag → `spacecraft set-state shipped` → capture issue-close evidence. Details: `references/release-policy.md` §Post-merge.
- **Must not**: Tag before merge succeeds; push unless asked.

### Dependency Freshness

- **Must**: Before adding/updating direct deps, check official docs/registry/releases. Prefer latest stable. Record source, version, date when it affects implementation.
- **May**: Exceptions: transitive deps, ecosystem pins, lockfile constraints, security advisories, or explicit instruction.

## Out of scope

This skill does NOT handle:

- Mission lifecycle or planning - use sc-mission, sc-planning
- Evidence capture or verification - use sc-verification
- Direct user clarification - use `/sc-discuss` / sc-clarify
- UI design - draft under `/sc-discuss`; QC via sc-ux-design

## Output format

```
Branch pattern (work): <type>/<id>/<title>
Branch pattern (merge): <type>/<title> after stripping id
Commit subject: <type>: <description> (~72 chars)
Merge: git merge --no-ff <type>/<title>
Tag: git tag -a v<major>.<minor>.<patch> -m "v<major>.<minor>.<patch>"
```

## Checklist

Before claiming git work is done:

- [ ] Not on `main` for product changes
- [ ] Branch named `<type>/<id>/<title>` while working; stripped to `<type>/<title>` immediately before merge
- [ ] `.gitignore` current; no secrets staged
- [ ] Conventional Commits used
- [ ] AFK checkpoints present during `/sc-run`; squashed to ≤5 before ship merge
- [ ] Verification passed after latest rebase
- [ ] Closeout check passes before release claim
- [ ] After merge: ran `set-state shipped`, captured evidence of issue closing

## References

- `references/release-policy.md` - checkpoint table, closeout bars, tagging, post-merge (on-demand)
- `spacecraft resolve --help` - resolver subcommand
- `spacecraft closeout-check --help` - closeout verification
