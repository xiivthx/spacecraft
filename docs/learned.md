# Lessons Learned

> Aggregated from shipped missions. Source of truth for internal research (`spacecraft research`).
> Each entry includes mission context for traceability.

---

## Solved

Specific issues that were identified, fixed, and verified during missions.

| Mission | Date | Problem | Solution | Evidence |
|---------|------|---------|----------|----------|
| M07IMLU48 | 2026-07-09 | 20 documented issues across skills, commands, agents, docs — all pre-existing from audit | All 20 resolved via 9 tasks; S15 (sc-debug length) initially deferred, then trimmed to 199 lines | E07IMTO4G–E07IP8RUW |
| M07IMLU48 | 2026-07-09 | TDD cycle missing Plan and Refactor steps; no triage gate; no phase splitting for >7 tasks | Added Plan→Red→Green→Verify→Refactor→Review; Skip TDD when rules; Phase 1/2/N with separate plan files | E07INHF6H |
| M07IMLU48 | 2026-07-09 | Stale evidence for 5 tasks after pipeline changes; 14 commits exceeded branch limit | Recaptured fresh evidence; squashed to 1 commit | E07IP81C5, E07IP86UT, E07IP86W7 |
|---------|------|---------|----------|----------|
| M07IDBF29 | 2026-07-09 | Duplicate helper functions across Go packages (fileExists, readJSON, writeJSON, taskIsComplete, blockingFindings) | Consolidated into util/fs.go and mission package | E07IEJDSL, E07IEMBWU |
| M07IDBF29 | 2026-07-09 | copyTextFile silently fails (returns bool, callers ignore) | Changed to return error, handle at call sites | E07IEJDSL |
| M07IDBF29 | 2026-07-09 | researchCmd goto anti-pattern and deep variable collision | Replaced with structured control flow, renamed variable | E07IERLO6 |
| M07IDBF29 | 2026-07-09 | Non-standard exit codes in checkDepsCmd (2=update,1=error) and resolveCmd | Normalized to Unix conventions (0=success,1=error) | E07IEV9VK |
| M07IDBF29 | 2026-07-09 | lookupPackage nil return with no error log; evidenceCmd exit code propagation breaks pipelines | Added stderr warning log; exit 0 always, print subprocess code to stdout | E07IEW9Q6 |
| M07IDBF29 | 2026-07-09 | Dead empty if-block in workflow.go; missing sections in sc-quick.md/sc-resume.md; fragile inline Python in sc-resume.md | Removed dead code; added Pre-flight/Error Handling/Hard Stop Gates sections; replaced Python with Go/shell approach | E07IEXXYT, E07IEOJKN |
| M07IDBF29 | 2026-07-09 | Closeout checker gate statuses are keyword-enforced — specNote rejected "not-applicable"; only specific keywords allowed | Used "deferred" instead; check `defaultReleaseGateStatuses` before populating review.json releaseReadiness | E07IFL9C2 |
| M07IDBF29 | 2026-07-09 | review.json findings missing `id`/`summary` fields cause "unnamed" blocking findings | Populated `id` and `summary` per Finding struct schema for proper blocking-findings identification | E07IFL9C2 |
| M07IG6R17 | 2026-07-09 | Duplicate ID normalization in archive.go (normalizeMissionIdSimple) — third copy of slug logic | Unified with util.NormalizeMissionId (G3) | E07IGRR4U |
| M07IG6R17 | 2026-07-09 | Config package only supports root-derived paths — no way to override individual paths for testing | Added ConfigOption functional options: WithSpaceDir, WithMissionsDir, WithArchiveDir, WithCurrentFile (G15) | E07IGUZF7 |
| M07IG6R17 | 2026-07-09 | 22 unused backward-compat type aliases in types.go | Removed all aliases (G16) | E07IGWVI1 |
| M07IG6R17 | 2026-07-09 | GitInfo and GitInfoData — overlapping structs with different fields, JSON consumers see different shapes | Consolidated into GitInfoData (additive superset — preserves all old JSON fields) (G17) | E07IH12NO |
| M07IG6R17 | 2026-07-09 | NoopRunner variable named as noop but runs real OS commands (gitutil/git.go) | Renamed to DefaultRunner to reflect actual behavior (G18) | E07IH2MZG |
| M07IG6R17 | 2026-07-09 | Manual flag parsing loop in researchCmd and checkDepsCmd — flag package supports dashes | Replaced with flag.NewFlagSet (G19) | E07IHF9J4 |
| M07IG6R17 | 2026-07-09 | Commander session handoff recommends auto-trigger skills as user slash commands | Clarified auto-trigger vs user commands in PERSONA.md (D7) | E07IHJBKL |
| M07JSKJRB | 2026-07-10 | 3 Node integration tests in resolver suite failed because test workspaces lacked the `.opencode/commands/` directory required by `validateNextCommand` — CLI initialization set the path but test harness never created it | Added `ensureCommandsDir` helper to populate `.opencode/commands/` in test workspaces; updated expected clarification format from `/sc-clarify` to `(clarify)` | E07JWWW0N |
| M07JZ2OPD | 2026-07-10 | `.opencode/.gitignore` blanket-ignored `package.json`, `package-lock.json`, `.gitignore` across all skills — sc-ux-design needed its own package.json for Playwright but files were invisible to git | Added negation patterns (`!skills/sc-ux-design/package.json`, etc.) to `.opencode/.gitignore`; verified with `git check-ignore` | E07K0ASLZ |
| M07MTPHTR | 2026-07-12 | T4/T7 evidence truncated — multi-command chains (`echo && grep && grep`) only captured first command output; verify passed because echo succeeded | Re-captured each grep as separate `spacecraft evidence` call (E07MW9P28, E07MW9PM6, E07MW9Q6Z for T4; E07MW9X07, E07MW9XKQ, E07MW9YJW for T7) | E07MW9P28, E07MW9X07 |
| M07MTPHTR | 2026-07-12 | Phantom reference files in sc-performance SKILL.md — 5 local files listed in References section that didn't exist | Replaced with external URLs matching sc-security pattern | E07MW9X07 |
| M07MTPHTR | 2026-07-12 | AGENTS.md skill table ordering — sc-security before sc-performance (should be alphabetical) | Swapped rows to correct alphabetical order | E07MUVRW3 |
| M07MTPHTR | 2026-07-12 | SPACECRAFT.md table had duplicate rows and incorrect alphabetical ordering for sc-performance/sc-security | Moved to correct positions, removed duplicates | E07MUVRW3 |
| M07MTPHTR | 2026-07-12 | AC11 evidence was echo claim instead of actual review output | Replaced with reference to actual review.json (E07MWA5YH) | E07MWA5YH |

---

## Lessons

General principles and transferable insights — applicable beyond this codebase. Emerged from mission work but framed as world-wide solutions.

| Mission | Date | Lesson | Why it matters |
|---------|------|--------|----------------|
| M07IMLU48 | 2026-07-09 | Machine-validated JSON schemas with enum-like keyword constraints silently reject natural-language status values — always check the source code's allowed statuses before populating `ReleaseGate.status` | "ready" and "approved" feel natural to humans but don't match the machine's allowlist. Check `defaultReleaseGateStatuses` in the source before populating any gate's status field |
| M07IMLU48 | 2026-07-09 | A resolved critical finding still blocks release if the severity field is not downgraded — the finding struct checks severity, not resolution status | When a review finding is resolved, either downgrade its severity from "critical" to "info" or set `blocksShip: false`. The finding documents what was found; the severity determines whether the gate passes |
| M07IMLU48 | 2026-07-09 | Version bump baseline must come from the authoritative version record (CHANGELOG), not from the most recent git tag — tags can drift from the human-facing version history | Before planning a version bump, always read CHANGELOG.md to find the latest version entry. Git tags may be stale or from a different tag sequence |
|---------|------|--------|----------------|
| M07IDBF29 | 2026-07-09 | Always verify pre-existing test failures against the base commit before treating them as blocking regressions | CI failures are not always introduced by the current change. A failing test that also fails on the base commit is pre-existing noise — document it, don't block the release on it |
| M07IDBF29 | 2026-07-09 | Resolved findings should not block release — severity downgrade prevents false blocking | A finding marked "critical" still blocks even after resolution. When a finding is resolved, either downgrade severity or mark `blocksShip: false` so the release gate reflects reality, not history |
| M07IG6R17 | 2026-07-09 | Configuration-as-code schemas with machine-enforced enums must document allowed values — `ReleaseGate.status` silently rejects descriptive words like "satisfied", "approved", "pending" because only specific keywords in the allowlist pass | Always check the source code's status keyword allowlist before populating machine-validated JSON fields. Descriptive language that feels natural to humans often doesn't match the machine's enum |
| M07IG6R17 | 2026-07-09 | Version bump and changelog commits must be separate from implementation commits — git history should distinguish "what changed" from "what was released" | Keep release note commits (`chore:` or `docs:`) as a dedicated final commit in the work branch, never bundled with feature or fix implementation. This makes rollback and audit easier |
| M07JSKJRB | 2026-07-10 | When adding filesystem-based validation to a CLI, update integration test harnesses to populate the required directory structure — otherwise integration tests silently exercise fallback code paths instead of the validation path | `validateNextCommand` checks if `.opencode/commands/<cmd>.md` exists; if the test workspace doesn't create that directory, every command falls back to `/sc-resume`. Unit tests bypass this with mocks, creating a false sense of coverage. Always verify that integration test fixtures mirror the CLI's initialization assumptions |
| M07JZ2OPD | 2026-07-10 | Gitignore negation patterns require the full directory path relative to the gitignore file — a pattern like `!sc-ux-design/package.json` won't match `skills/sc-ux-design/package.json` | `git check-ignore -v` is the authoritative test before committing — if the tool says a negation pattern is active but the file is still ignored, the path prefix is wrong. Always verify with the real command output before shipping |
| M07JZ2OPD | 2026-07-10 | When a system-level constraint (global gitignore) conflicts with a design decision (skill-local dependencies), resolve the conflict in both the constraint file (with explanatory comments) and the architecture decisions | Future skills needing Node packages, Python venvs, or build artifacts will hit the same conflict. Precedents and resolutions should be discoverable without re-inventing the solution |
| M07MTPHTR | 2026-07-12 | Multi-command evidence chains (`echo && grep && grep`) silently truncate to first command only in CLI runners that capture stdout from a single subprocess — exit code 0 from echo masks missing grep results | Any evidence capture tool that executes a single subprocess and captures its stdout will silently drop subsequent commands in a chain. Always run one check per evidence call to avoid false positives |
| M07MTPHTR | 2026-07-12 | Skill reference files listed in a References section must be verified for existence before commit — phantom refs violate sc-creator convention and break downstream consumers silently | Prefer external URLs over local file paths for non-essential references. If local files are required, add an existence check in the skill creation workflow |
| M07MTPHTR | 2026-07-12 | When a planned change touches a shared config file affecting multiple components, expanding scope to cover the whole file in one mission can reduce total churn vs. repeated incremental changes across missions | For configuration-as-code systems, evaluate coherence before scoping narrowly. A single-file change touching multiple agents may be cleaner as one mission than fragmented across several |

---

## Patterns

Recurring solutions that proved effective across missions.

| Pattern | First Seen (Mission) | Description | When to Apply |
|---------|---------------------|-------------|---------------|
<!-- Emerge from repeated lessons; deduplicate manually -->
