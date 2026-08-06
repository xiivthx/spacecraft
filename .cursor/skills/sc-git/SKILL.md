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
2. **Check git state** - Run `spacecraft git-info` before committing/merging/releasing.
3. **Branch** - Create a non-main work branch from latest `main` before mutating. Never write product changes on `main`.
4. **Commit (AFK checkpoints)** - During `/sc-run` build, auto-commit after every RED, GREEN, triage-skip direct-write+evidence, and post-feature refactor (see §Checkpoint commits). These are WIP on the work branch only.
5. **Squash before ship** - On `/sc-ship`, squash/fixup checkpoints into 1–3 logical Conventional Commits (max 5) before merge. See sc-ship. `/sc-quick` keeps 1–3 commits (no mission squash ceremony).
6. **Commit (release notes)** - Add changelog update as a **separate commit** in the work branch before merge (type `chore:` or `docs:`). Never defer after merge.
7. **Verify** - After latest rebase, reverify. For missions, run `spacecraft closeout-check` before claiming release readiness. `/sc-quick` skips closeout (hook uses `SPACECRAFT_QUICK=1`).
8. **Rebase & Merge** - Before merge, rebase on latest `main`. Merge with `--no-ff` only. Mission: `SPACECRAFT_SHIP=1`. Quick: `SPACECRAFT_SHIP=1 SPACECRAFT_QUICK=1`.
9. **Tag** - After no-ff merge, create annotated tag. Follow bump policy from Rules §Release Prep: breaking=major, feature=minor, fix=patch, docs/chore=still tag as next patch. Do not tag before merge exists.

## Rules

### General

- **Must**: Treat git as the rollback and release boundary for mutating work.
- **May**: Discovery, clarification, design, planning, and review may run without git.
- **Must**: Before git work on a mission, resolve via `spacecraft resolve`. On conflict or ambiguity, use `spacecraft use <number|id|title>` or `SPACECRAFT_MISSION`. Skip resolve for `/sc-quick`.
- **Must**: Before committing/merging/releasing, run `spacecraft git-info`.
- **Must not**: Write product changes on `main`. If on `main` when mutation is requested, create a work branch.
- **Must not**: Auto-run `git init`, create worktrees, rebase, merge, tag, or push unless asked.
- **Must**: Before outward actions (push, deploy, publish, send), state `AUTH:` with a **quoted** user authorization from the conversation. Mission local ship still requires `/sc-ship` + hooks + `SPACECRAFT_SHIP=1`. `/sc-quick` authorizes local merge/tag in that lane (hooks + `SPACECRAFT_SHIP=1 SPACECRAFT_QUICK=1`); push still needs separate AUTH. AUTH does not bypass ship hooks.

### Branching

- **Must**: One branch = one feature/fix/scoped change. Branch from latest `main`.
- **Must**: Pattern: `<type>/<id>/<title>` for missions - e.g. `feat/m07fp1l7z/go-rewrite`. `/sc-quick` uses `<type>/<title>` with no mission id.
- **Must**: `release/v<major>.<minor>.<patch>` only for release prep.
- **Must**: If a branch needs >5 final commits, split the feature first.

### Commits

- **Must**: Target 1–3 **final** commits per branch after ship squash, max 5 unless justified in decisions.md.
- **Must**: Separate version bump + changelog update into its own `chore:` or `docs:` commit - do not bundle with implementation commit.
- **Must**: During `/sc-run` AFK build, auto-create checkpoint commits (see §Checkpoint commits).
- **Must**: Before `/sc-ship` merge, squash/fixup checkpoints into logical Conventional Commits (≤5). Do not squash unrelated changes.
- **Must not**: Stage unrelated user changes. Prefer `git add <specific-files>`.
- **Must**: Subject: `<type>: <description>` - imperative, lowercase, ~72 chars.
- **Must**: Body: `-` bullets, lowercase start, no period end. Describe what changed; evidence lives in `evidence.jsonl`, not the commit.
- **Must not**: Put mission ids (`M…`) in commit subjects or bodies (including merge commits and AFK checkpoints). Keep ids on the work branch name / mission artifacts only.
- **Must**: Types: `feat`, `fix`, `docs`, `refactor`, `test`, `build`, `ci`, `chore`, `perf`, `style`, `revert`.
- **Must**: Breaking: `!` after type or `BREAKING CHANGE:` footer.

### Checkpoint commits

Used by `/sc-run` on the work branch. Auto-commit; never push.

| Step | When | Suggested type |
|------|------|----------------|
| RED | Failing test for one acceptance is committed | `test:` |
| GREEN | Minimal code passes that acceptance + evidence captured | `feat:` or `fix:` |
| Skip | Triage skip (tautology / docs-prose): direct write + evidence; no RED harness | `docs:` / `feat:` / `fix:` |
| Combine | Post-feature refactor and/or integration/functional gate | `refactor:` / `test:` |

- **Must**: One checkpoint per RED, per GREEN, per triage-skip direct-write+evidence, and after the combine/refactor gate.
- **Must**: Subject stays Conventional Commits; body may include `- wip checkpoint`, task id, acceptance summary (and `skip: <reason>` when triage skipped). Do not include the mission id.
- **Must not**: Invent RED `test:` checkpoints for triage-skip / docs-prose wording-only acceptances.
- **Must not**: Treat checkpoint count as the final ship commit budget - squash at `/sc-ship`.
- **Must not**: Checkpoint-commit unrelated user dirty files.

### Ignore Hygiene

- **Must**: Before staging, inspect untracked/modified files. Update `.gitignore` for new build outputs, caches, logs, env files, deps, or machine files.
- **Must**: Use stack-aware patterns (Node, Rust, Python, etc.) when those tools are present.
- **Must**: Keep source, tests, migrations, lockfiles, configs, and release notes tracked.
- **Must**: If unsure about a file, ask before staging. Before release, verify no secrets are tracked.

### Rebase & Merge

- **Must**: Before merge, rebase on latest `main`. Reverify after rebase.
- **Must**: Identify fork point with `git log --oneline main..HEAD | head -1` before rebase. If `main` has advanced beyond the expected base, warn: "Rebase target mismatch: main HEAD differs from fork point. Confirm correct base before rebase."
- **Must**: Immediately before merge, strip the mission id from the work branch when it matches `<type>/<id>/<title>` and `<id>` is a mission id (`M…`): `git branch -m <type>/<title>`. Finish closeout/evidence while the id is still in the name. See sc-ship §Merge.
- **Must**: Merge only with `--no-ff`. No fast-forward or squash-merge into `main`.
- **Must not**: Rebase/rewrite `main`. Resolve conflicts on work branch, then reverify.
- **Must**: After merge, delete local branch unless asked to keep.

### Release Prep

- **Must**: User "ship/release/merge/finish/close branch" → release closeout.
- **Must**: "stop/close/end session" → session handoff (no merge/tag/branch delete).
- **Must**: If ambiguous, recommend ship; don't auto-merge.
- **Must**: For missions, run `spacecraft closeout-check` (alias: `ship-check`) before claiming readiness. The CLI machine-enforces: required artifacts present; mission state `ready` or `shipped`; clarify-status not `open`; evidence.jsonl non-empty with `label`/`command`/`output`/`ts`/`exitCode`; `review.json` status `ready` with **empty** `findings` (0 errors, 0 warnings); `releaseReadiness.changelog` and `releaseReadiness.specNote` objects with status `ready` (string/boolean gates invalid); at least one commit touching `CHANGELOG.md` since `main`/`origin/main`. `/sc-quick` skips closeout; ship with `SPACECRAFT_SHIP=1 SPACECRAFT_QUICK=1`.
- **Must**: If any gate is incomplete, block and list missing actions.
- **Must**: Cursor ship hooks re-run closeout when `SPACECRAFT_SHIP=1` before allowing merge/push/tag.
- **Must**: Bump version by impact: breaking=major, feature=minor, fix/patch=patch, docs/chore=no bump unless part of a release (record in decisions.md).
- **Must**: Update changelog and spec/release note when behavior changed.

### Tagging

- **Must**: After EVERY no-ff merge to `main`, create an annotated tag: `git tag -a v<major>.<minor>.<patch> -m "v<major>.<minor>.<patch>"`. Tagging is mandatory after every merge regardless of bump policy.
- **Must**: When bump policy says "no bump" (docs/chore), still tag as the next patch version (e.g., v0.6.0 → v0.6.1). The tag tracks changes to main; the bump policy only controls whether the version number reflects feature/breaking scope.
- **Must**: Create tag ONLY AFTER merge is confirmed clean (`git merge --no-ff` completes successfully). Never create tag before merge, even if merge retries are expected.
- **Must**: If merge is reverted/redone, delete the premature tag first before retrying. Verify tag was created only once on final clean merge.
- **Must**: Do not push unless asked.

### Post-merge

After the no-ff merge completes, immediately execute the mandatory steps from Rules §Tagging and §Release Prep:
- Tag, verify tag, delete branch.
- **Must**: Run `spacecraft set-state shipped` to trigger archive and close GitHub issues referenced in artifacts.
- **Must**: Capture evidence of issue closing output (e.g., "Issues: X closed, Y already closed").
- After cleanup you are on `main` - any further mutation requires a new non-main branch.

### Review Gate

Before shipping or merging: verify all Rules §Rebase & Merge, §Release Prep, §Tagging gates pass. Key checks: rebased on `main`, ≤5 commits, Conventional Commits, `.gitignore` current, version bump + changelog in separate commit, `--no-ff` merge plan.

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

- `spacecraft resolve --help` - resolver subcommand
- `spacecraft git-info` - git state inspection
- `spacecraft closeout-check --help` - closeout verification
