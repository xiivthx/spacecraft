---
name: sc-git
description: Enforce Spacecraft git safety, branching, hygiene, Conventional Commits, verification before main, no-ff merge, version bump, changelog/spec update, and release tag policy.
license: MIT
---
- Treat git as the rollback and release boundary for mutating work.
- Discovery, clarification, design, planning, and review may run without git.
- Before git work, resolve mission with `scripts/spacecraft resolve --json`. If safety ≠ `safe`, block until user selects with `scripts/spacecraft use <number|id|title>` or `SPACECRAFT_MISSION`.
- Before committing/merging/releasing, run `scripts/spacecraft git-info`.
- Never write product changes on `main`. If on `main` when mutation is requested, create a work branch.
- Do not auto-run `git init`, create worktrees, rebase, merge, tag, or push unless asked.
- Use rtk for noisy git output when available.

## Branching

- One branch = one feature/fix/scoped change. Branch from latest `main`.
- Pattern: `<type>/<id>/<title>` — e.g. `feat/m07fp1l7z/go-rewrite`.
- `release/v<major>.<minor>.<patch>` only for release prep.
- If a branch needs >5 final commits, split the feature first.

## Commits

- Target 1–3 final commits per branch, max 5 unless justified in decisions.md.
- Frequent WIP/checkpoint commits OK on work branch only. After a passing `/sc-flow` task, may create a checkpoint commit.
- Before merge, squash/fixup into logical Conventional Commits. Do not squash unrelated changes.
- Never stage unrelated user changes. Prefer `git add <specific-files>`.
- Subject: `<type>: <description>` — imperative, lowercase, ~72 chars.
- Body: `-` bullets, lowercase start, no period end. Include mission/evidence ids.
- Types: `feat`, `fix`, `docs`, `refactor`, `test`, `build`, `ci`, `chore`, `perf`, `style`, `revert`.
- Breaking: `!` after type or `BREAKING CHANGE:` footer.

## Ignore Hygiene

- Before staging, inspect untracked/modified files. Update `.gitignore` for new build outputs, caches, logs, env files, deps, or machine files.
- Use stack-aware patterns (Node, Rust, Python, etc.) when those tools are present.
- Keep source, tests, migrations, lockfiles, configs, and release notes tracked.
- If unsure about a file, ask before staging. Before release, verify no secrets are tracked.

## Rebase & Merge

- Before merge, rebase on latest `main`. Reverify after rebase.
- Merge only with `--no-ff`. No fast-forward or squash-merge into `main`.
- Do not rebase/rewrite `main`. Resolve conflicts on work branch, then reverify.
- After merge, delete local branch unless asked to keep.

## Release Prep

- User "ship/release/merge/finish/close branch" → release closeout.
- "stop/close/end session" → session handoff (no merge/tag/branch delete).
- If ambiguous, recommend `/sc-ship`; don't auto-merge.
- Run `scripts/spacecraft closeout-check` before claiming readiness. It requires `review.json.releaseReadiness` object entries (string/boolean gates invalid).
- If any gate is incomplete, block and list missing actions.
- Bump version by impact: breaking=major, feature=minor, fix/patch=patch, docs/chore=no bump unless part of a release (record in decisions.md).
- Update changelog and spec/release note when behavior changed.

## Tagging

- After no-ff merge, create annotated tag: `git tag -a v<major>.<minor>.<patch> -m "v<major>.<minor>.<patch>"`. Do not tag before merge exists. Do not push unless asked.

## Review Gate

Before shipping or merging, check:
- not on `main` while editing product files
- branch based/rebased on latest `main`
- ≤5 final commits (or exception in decisions.md)
- final commits use Conventional Commits
- `.gitignore` current; no secrets/deps/builds staged
- tests/verification pass after latest rebase
- version + changelog/spec update present or deferred with rationale
- merge plan uses `--no-ff`; tag plan exists

## Dependency Freshness

- Before adding/updating direct deps, check official docs/registry/releases. Prefer latest stable. Record source, version, date when it affects implementation.
- Exceptions: transitive deps, ecosystem pins, lockfile constraints, security advisories, or explicit instruction.
