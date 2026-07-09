---
name: sc-git
description: Enforce git safety, branching, Conventional Commits, no-ff merge, versioning, and release tags. Activate on /sc-git, commit, branch, merge, release, or ship.
license: MIT
compatibility: opencode
metadata:
  version: 1.0
  audience: xiivthx
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

1. **Resolve mission** — Before git work, resolve mission with `scripts/spacecraft resolve --json`. If safety ≠ `safe`, block until user selects.
2. **Check git state** — Run `scripts/spacecraft git-info` before committing/merging/releasing.
3. **Branch** — Create a non-main work branch from latest `main` before mutating. Never write product changes on `main`.
4. **Commit (implementation)** — Use Conventional Commits. Target 1–3 final commits per branch, max 5 unless justified in decisions.md.
5. **Commit (release notes)** — Add version bump + changelog update as a **separate commit** in the work branch before merge (type `chore:` or `docs:`). Never defer after merge.
6. **Verify** — After latest rebase, reverify. Run `scripts/spacecraft closeout-check` before claiming release readiness.
7. **Rebase & Merge** — Before merge, rebase on latest `main`. Merge with `--no-ff` only.
8. **Tag** — After no-ff merge, create annotated tag. Do not tag before merge exists.

## Rules

### General

- **Must**: Treat git as the rollback and release boundary for mutating work.
- **May**: Discovery, clarification, design, planning, and review may run without git.
- **Must**: Before git work, resolve mission with `scripts/spacecraft resolve --json`. If safety ≠ `safe`, block until user selects with `scripts/spacecraft use <number|id|title>` or `SPACECRAFT_MISSION`.
- **Must**: Before committing/merging/releasing, run `scripts/spacecraft git-info`.
- **Must not**: Write product changes on `main`. If on `main` when mutation is requested, create a work branch.
- **Must not**: Auto-run `git init`, create worktrees, rebase, merge, tag, or push unless asked.

### Branching

- **Must**: One branch = one feature/fix/scoped change. Branch from latest `main`.
- **Must**: Pattern: `<type>/<id>/<title>` — e.g. `feat/m07fp1l7z/go-rewrite`.
- **Must**: `release/v<major>.<minor>.<patch>` only for release prep.
- **Must**: If a branch needs >5 final commits, split the feature first.

### Commits

- **Must**: Target 1–3 final commits per branch, max 5 unless justified in decisions.md.
- **Must**: Separate version bump + changelog update into its own `chore:` or `docs:` commit — do not bundle with implementation commit.
- **May**: Frequent WIP/checkpoint commits OK on work branch only. After a passing `/sc-build` task, may create a checkpoint commit.
- **Must**: Before merge, squash/fixup into logical Conventional Commits. Do not squash unrelated changes.
- **Must not**: Stage unrelated user changes. Prefer `git add <specific-files>`.
- **Must**: Subject: `<type>: <description>` — imperative, lowercase, ~72 chars.
- **Must**: Body: `-` bullets, lowercase start, no period end. Include mission/evidence ids.
- **Must**: Types: `feat`, `fix`, `docs`, `refactor`, `test`, `build`, `ci`, `chore`, `perf`, `style`, `revert`.
- **Must**: Breaking: `!` after type or `BREAKING CHANGE:` footer.

### Ignore Hygiene

- **Must**: Before staging, inspect untracked/modified files. Update `.gitignore` for new build outputs, caches, logs, env files, deps, or machine files.
- **Must**: Use stack-aware patterns (Node, Rust, Python, etc.) when those tools are present.
- **Must**: Keep source, tests, migrations, lockfiles, configs, and release notes tracked.
- **Must**: If unsure about a file, ask before staging. Before release, verify no secrets are tracked.

### Rebase & Merge

- **Must**: Before merge, rebase on latest `main`. Reverify after rebase.
- **Must**: Merge only with `--no-ff`. No fast-forward or squash-merge into `main`.
- **Must not**: Rebase/rewrite `main`. Resolve conflicts on work branch, then reverify.
- **Must**: After merge, delete local branch unless asked to keep.

### Release Prep

- **Must**: User "ship/release/merge/finish/close branch" → release closeout.
- **Must**: "stop/close/end session" → session handoff (no merge/tag/branch delete).
- **Must**: If ambiguous, recommend `/sc-ship`; don't auto-merge.
- **Must**: Run `scripts/spacecraft closeout-check` before claiming readiness. It requires `review.json.releaseReadiness` object entries (string/boolean gates invalid).
- **Must**: If any gate is incomplete, block and list missing actions.
- **Must**: Bump version by impact: breaking=major, feature=minor, fix/patch=patch, docs/chore=no bump unless part of a release (record in decisions.md).
- **Must**: Update changelog and spec/release note when behavior changed.

### Tagging

- **Must**: After EVERY no-ff merge to `main`, create an annotated tag: `git tag -a v<major>.<minor>.<patch> -m "v<major>.<minor>.<patch>"`. Tagging is mandatory after every merge regardless of bump policy.
- **Must**: When bump policy says "no bump" (docs/chore), still tag as the next patch version (e.g., v0.6.0 → v0.6.1). The tag tracks changes to main; the bump policy only controls whether the version number reflects feature/breaking scope.
- **Must**: Tag immediately after merge, before any other operation. Do not pause or defer.
- **Must**: Do not tag before merge exists. Do not push unless asked.

### Post-merge

After the no-ff merge completes, immediately:

- [ ] Create annotated tag (`git tag -a v<major>.<minor>.<patch> -m "v<major>.<minor>.<patch>"`)
- [ ] Verify tag exists (`git tag -l 'v*' | tail -3`)
- [ ] Delete merged local branch (`git branch -d <branch>`)
- [ ] Run `scripts/spacecraft archive` to compact shipped mission artifacts

**Important**: After post-merge cleanup, you are on `main`. Any further edits — even small fixes or documentation — MUST first create a new non-main branch. Creating a branch is automatic and non-optional whenever mutation is requested while on `main`. Do not edit, commit, or stage anything on `main` directly.

### Review Gate

Before shipping or merging, check:

- [ ] not on `main` while editing product files (this applies post-merge too — always branch first)
- [ ] branch based/rebased on latest `main`
- [ ] ≤5 final commits (or exception in decisions.md)
- [ ] final commits use Conventional Commits
- [ ] `.gitignore` current; no secrets/deps/builds staged
- [ ] tests/verification pass after latest rebase
- [ ] version bump + changelog/spec update in a **separate commit** (`chore:`/`docs:`) in work branch (not deferred after merge)
- [ ] merge plan uses `--no-ff`; tag plan exists

### Dependency Freshness

- **Must**: Before adding/updating direct deps, check official docs/registry/releases. Prefer latest stable. Record source, version, date when it affects implementation.
- **May**: Exceptions: transitive deps, ecosystem pins, lockfile constraints, security advisories, or explicit instruction.

## Out of scope

This skill does NOT handle:

- Mission lifecycle or planning — use sc-mission, sc-planning
- Evidence capture or verification — use sc-verification
- UI design — use sc-design
- Direct user clarification — use sc-clarify

## Output format

```
Branch pattern: <type>/<id>/<title>
Commit subject: <type>: <description> (~72 chars)
Merge: git merge --no-ff <branch>
Tag: git tag -a v<major>.<minor>.<patch> -m "v<major>.<minor>.<patch>"
```

## Checklist

Before claiming git work is done:

- [ ] Not on `main` for product changes
- [ ] Branch named with `<type>/<id>/<title>` pattern
- [ ] `.gitignore` current; no secrets staged
- [ ] Conventional Commits used
- [ ] Verification passed after latest rebase
- [ ] Closeout check passes before release claim

## Research auto-trigger

When git operations involve unfamiliar flags, rebase conflict resolution strategies, or tag/signing conventions, run `spacecraft research "git <topic>"` before executing. Git mistakes are hard to undo — verify the command before running it.

---

## References

- `scripts/spacecraft resolve --help` — resolver subcommand
- `scripts/spacecraft git-info` — git state inspection
- `scripts/spacecraft closeout-check --help` — closeout verification
