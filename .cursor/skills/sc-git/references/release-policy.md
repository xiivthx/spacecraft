# Git release policy (closeout, tagging, checkpoints)

On-demand tables and long policy for `/sc-run` checkpoints, closeout, tagging, and post-merge. Inline Must for Conventional Commits and ship gates stay in `../SKILL.md`. Load this file when running AFK checkpoints, ship squash, or release closeout.

## Checkpoint commits

Used by `/sc-run` on the work branch. Auto-commit; never push.

| Step | When | Suggested type |
|------|------|----------------|
| Task | One Conventional Commit per plan task after that task's acceptances are done (TDD RED+GREEN complete, or triage-skip direct-write+evidence) - not per RED and per GREEN | `feat:` / `fix:` / `test:` / `docs:` |
| Combine | Post-feature refactor and/or integration/functional gate | `refactor:` / `test:` |
| Fix | Material fix during fix pass | `fix:` |

- **Must**: One Conventional Commit per plan task after that task's acceptances are done. Combine and material-fix checkpoints may remain.
- **Must**: Subject stays Conventional Commits; body may include `- wip checkpoint`, task id, acceptance summary (and `skip: <reason>` when triage skipped). Do not include the mission id.
- **Must not**: Invent RED `test:` checkpoints for triage-skip / docs-prose wording-only acceptances - one `docs:` / `feat:` / `fix:` checkpoint is enough.
- **Must not**: Treat checkpoint count as the final ship commit budget - squash at `/sc-ship`.
- **Must not**: Checkpoint-commit unrelated user dirty files.

## Release Prep / closeout

- **Must**: User "ship/release/merge/finish/close branch" → release closeout.
- **Must**: "stop/close/end session" → session handoff (no merge/tag/branch delete).
- **Must**: If ambiguous, recommend ship; don't auto-merge.
- **Must**: For missions, run `spacecraft closeout-check` (alias: `ship-check`) before claiming readiness. The CLI machine-enforces: required artifacts present; mission state `ready` or `shipped`; clarify-status not `open`; evidence.jsonl non-empty with `label`/`command`/`output`/`ts`/`exitCode`; `review.json` status `ready` with **empty** `findings` (0 errors, 0 warnings); `releaseReadiness.changelog` and `releaseReadiness.specNote` objects with status `ready` (string/boolean gates invalid); at least one commit touching `CHANGELOG.md` since `main`/`origin/main`. `/sc-quick` skips closeout; ship with `SPACECRAFT_SHIP=1 SPACECRAFT_QUICK=1`.
- **Must**: If any gate is incomplete, block and list missing actions.
- **Must**: Cursor ship hooks re-run closeout when `SPACECRAFT_SHIP=1` before allowing merge/push/tag.
- **Must**: Bump version by impact: breaking=major, feature=minor, fix/patch=patch, docs/chore=no bump unless part of a release (record in decisions.md).
- **Must**: Update changelog and spec/release note when behavior changed.

## Tagging

- **Must**: After EVERY no-ff merge to `main`, create an annotated tag: `git tag -a v<major>.<minor>.<patch> -m "v<major>.<minor>.<patch>"`. Tagging is mandatory after every merge regardless of bump policy.
- **Must**: When bump policy says "no bump" (docs/chore), still tag as the next patch version (e.g., v0.6.0 → v0.6.1). The tag tracks changes to main; the bump policy only controls whether the version number reflects feature/breaking scope.
- **Must**: Create tag ONLY AFTER merge is confirmed clean (`git merge --no-ff` completes successfully). Never create tag before merge, even if merge retries are expected.
- **Must**: If merge is reverted/redone, delete the premature tag first before retrying. Verify tag was created only once on final clean merge.
- **Must**: Do not push unless asked.

## Post-merge

After the no-ff merge completes, immediately execute tagging and release-prep cleanup:

- Tag, verify tag, delete branch.
- **Must**: Run `spacecraft set-state shipped` to trigger archive and close GitHub issues referenced in artifacts.
- **Must**: Capture evidence of issue closing output (e.g., "Issues: X closed, Y already closed").
- After cleanup you are on `main` - any further mutation requires a new non-main branch.

## Review Gate (ship)

Before shipping or merging: rebased on `main`, ≤5 commits, Conventional Commits, `.gitignore` current, version bump + changelog in separate commit, `--no-ff` merge plan, tagging after merge.

## Related

- Skill: `../SKILL.md` - Conventional Commits Must and ship gates
- `spacecraft closeout-check --help`
