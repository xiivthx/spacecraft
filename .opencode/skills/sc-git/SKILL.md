---
name: sc-git
description: Enforce Spacecraft git safety, release branching, branch hygiene, gitignore hygiene, Conventional Commits, verification before main, version bump, changelog/spec update, no-ff merge, and release tag policy. Use before mutating work, committing, squashing, rebasing, merging, tagging, shipping, or reviewing git readiness.
license: MIT
---
- Treat git as the rollback and release boundary for mutating work.
- Discovery, clarification, design exploration, planning, and read-only review may run without git.
- Before implementation, commit, merge, or release prep, resolve the mission with `scripts/spacecraft resolve --json`; `.space/current` is fallback state, not sole authority.
- If resolver safety is not `safe`, block git-changing work until the user selects with `scripts/spacecraft use <number|id|title>` or an explicit `SPACECRAFT_MISSION`.
- Before implementation, commit, merge, or release prep, run `scripts/spacecraft git-info`.
- Never write product changes directly on `main`.
- If clear mutating work is requested and the current branch is `main`, create or switch to a non-main work branch without another blocking question when policy permits it.
- Do not auto-run `git init`, create worktrees, rebase, merge, tag, or push unless the user explicitly asks.
- Do not push unless the user explicitly asks.
- Keep `.gitignore` current before staging, committing, or merging.
- Never allow secrets, credentials, local env files, generated dependency folders, build outputs, caches, logs, private artifacts, or machine-specific files into git or public release artifacts.
- Use rtk for noisy git/status/diff/log output when available. Do not use rtk to bypass denied operations.

## Branching Model

- Use Spacecraft release branching.
- `main` is protected and release-only.
- One branch equals one feature, fix, issue, or tightly scoped change.
- Branch from the latest `main`.
- Use branch names:
  - `<type>/<id>/<title>`
  - `feat/m07fp1l7z/go-rewrite`
  - `fix/m07fp1l7z/resolve-review-findings`
  - `chore/m07fp1l7z/update-docs-and-skills`
  - `docs/m07fp1l7z/add-api-reference`
  - `refactor/m07fp1l7z/extract-mission-module`
  - `release/v<major>.<minor>.<patch>` only for release preparation work.
- Prefer mission ids over issue ids when both exist.
- Prefer a separate git worktree for large, risky, or multi-session branches.
- If a branch is expected to need more than 5 final commits, split the feature before implementation.

## Commit Plan

- Before implementation, plan the intended final commits.
- Target 1 to 3 final commits per branch.
- A branch merged into `main` should not exceed 5 final commits unless explicitly justified in `decisions.md`.
- The agent may commit frequently inside a valid non-main work branch.
- After a task has passing verification evidence, `/sc-flow` may create a local checkpoint commit before starting the next task.
- Frequent WIP/checkpoint commits are allowed only on the work branch.
- Before merge, squash/fixup/checkpoint commits into logical commits so the branch history stays reviewable.
- Do not squash unrelated logical changes together.
- Never stage unrelated user changes.
- Prefer `git add <specific-files>` over broad staging.

## Ignore Hygiene

- Before staging or committing, inspect untracked and modified files.
- Update `.gitignore` when new tools, frameworks, build outputs, caches, logs, env files, local databases, generated artifacts, or machine-specific files appear.
- Use stack-aware ignore patterns for the project, such as Node/Next, Rust, Python, databases, IDEs, OS files, test coverage, and deployment output when those tools are present.
- Keep project artifacts that should be reviewed, such as specs, plans, migrations, source files, tests, and intended docs.
- Do not hide source files, migrations, lockfiles, configs, or release notes just to make status clean.
- If unsure whether a file is safe to track, stop and ask one question before staging it.
- Before public release, check that no secrets or private data are tracked or staged.

## Conventional Commits

- Use Conventional Commits for every final commit.
- Subject format:
  `<type>: <description>`
- Avoid scopes (`(scope)`) unless the change touches a distinct subsystem and the scope meaningfully aids review.
- Common types:
  `feat`, `fix`, `docs`, `refactor`, `test`, `build`, `ci`, `chore`, `perf`, `style`, `revert`.
- Use imperative, lowercase descriptions.
- Keep the subject short, roughly 72 characters when practical.
- Body format:
  - Use bullet points with `- ` prefix.
  - Start each bullet with a **lowercase** letter.
  - No period at end of each bullet.
- Put mission id, evidence ids, and longer rationale in the body when useful.
- Mark breaking changes with `!` after the type or a `BREAKING CHANGE:` footer.
- Examples:
  - `feat: implement core CLI commands`
    ```
    - add Go mission CRUD, resolution, workflow, and evidence subcommands
    ```
  - `fix: resolve review findings`
    ```
    - remove tracked Go binaries
    - implement printStatus() matching mjs behavior
    - update .gitignore for scripts/src/ binary patterns
    ```
  - `chore: update docs and skills to reference Go binary`
    ```
    - replace node scripts/spacecraft.mjs with scripts/spacecraft in AGENTS.md, SPEC.md
    - update all 14 command files in .opencode/commands/
    - update 3 skills: sc-git, sc-mission, sc-verification
    ```

## Rebase And Merge

- Before merge, rebase the work branch on the latest `main`.
- Re-run verification after rebase.
- Test, verify, and validate the system before any merge into `main`.
- Merge into `main` only with `--no-ff`.
- Do not fast-forward merge into `main`.
- Do not squash-merge into `main`; squash/fixup on the branch before the no-ff merge.
- Do not rebase or rewrite `main`.
- Resolve conflicts on the work branch, then verify again.
- After a successful merge to `main`, delete the merged local branch unless the user asks to keep it.

## Release Prep

- Treat user requests to ship, release, merge, finish the mission, or close a branch as release closeout prep.
- Treat ordinary stop-chat, close-session, end-session, or continue-in-new-session requests as session handoff. Do not merge, tag, or delete branches for handoff.
- If "close session" is ambiguous and work appears ready, recommend `/sc-ship`; do not merge automatically.
- Closeout must prepare the branch for merge into `main`.
- Run `scripts/spacecraft closeout-check` before claiming closeout readiness.
- `closeout-check` requires `review.json.releaseReadiness` object entries for version, changelog, spec note, tag plan, and post-rebase verification; deferred gates need rationale. String and boolean gates are invalid.
- If any gate is incomplete, block closeout and list exact missing actions.
- Before merging into `main`, bump the project version.
- Before merging into `main`, run the mission's required tests, verification commands, and validation commands.
- Choose version bump by impact:
  - breaking change: major
  - user-visible feature: minor
  - bug fix or small compatible change: patch
  - docs/test/chore-only changes: patch only when they are part of a release; otherwise record no version bump needed in `decisions.md`
- If the version source is unclear, ask one question before changing files.
- Update the changelog before merge.
- Update a short spec/release note before merge when the mission changed product behavior.
- Keep release notes concise:
  - what changed
  - why it matters
  - verification evidence
  - known limitations

## Tagging

- After the no-ff merge into `main`, create a version tag.
- Prefer annotated tags: `git tag -a v<major>.<minor>.<patch> -m "v<major>.<minor>.<patch>"`.
- Do not create a tag before the merge commit exists on `main`.
- Do not push tags unless the user explicitly asks.

## Review Gate

Before shipping or merging, check:
- not on `main` while editing product files
- branch was based on latest `main` or rebased onto latest `main`
- branch has 5 or fewer final commits, or an explicit exception
- final commits use Conventional Commits
- `.gitignore` is current for the stack and generated outputs
- no secrets, local env files, private data, caches, logs, dependencies, or build outputs are staged/tracked accidentally
- required tests, verification, and validation pass after the latest rebase
- version bump is present or explicitly deferred with rationale
- changelog/spec note is updated when needed
- verification evidence exists after the latest rebase
- merge plan uses `--no-ff`
- tag plan exists for the bumped version

## Dependency Freshness

- Before adding or updating direct dependencies, frameworks, or generated code that depends on current APIs, check official docs, registries, or releases.
- Prefer latest stable direct dependency versions.
- Exceptions: deep transitive dependencies, ecosystem pins, lockfile constraints, security advisories, or explicit user instruction.
- Record source, version, and date in decisions or evidence when the choice affects implementation.
