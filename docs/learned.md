# Lessons Learned

> Aggregated from shipped missions. Source of truth for internal research (`spacecraft research`).
> Each entry includes mission context for traceability.

---

## Solved

Specific issues that were identified, fixed, and verified during missions.

| Mission | Date | Problem | Solution | Evidence |
|---------|------|---------|----------|----------|
| M07IDBF29 | 2026-07-09 | Duplicate helper functions across Go packages (fileExists, readJSON, writeJSON, taskIsComplete, blockingFindings) | Consolidated into util/fs.go and mission package | E07IEJDSL, E07IEMBWU |
| M07IDBF29 | 2026-07-09 | copyTextFile silently fails (returns bool, callers ignore) | Changed to return error, handle at call sites | E07IEJDSL |
| M07IDBF29 | 2026-07-09 | researchCmd goto anti-pattern and deep variable collision | Replaced with structured control flow, renamed variable | E07IERLO6 |
| M07IDBF29 | 2026-07-09 | Non-standard exit codes in checkDepsCmd (2=update,1=error) and resolveCmd | Normalized to Unix conventions (0=success,1=error) | E07IEV9VK |
| M07IDBF29 | 2026-07-09 | lookupPackage nil return with no error log; evidenceCmd exit code propagation breaks pipelines | Added stderr warning log; exit 0 always, print subprocess code to stdout | E07IEW9Q6 |
| M07IDBF29 | 2026-07-09 | Dead empty if-block in workflow.go; missing sections in sc-quick.md/sc-resume.md; fragile inline Python in sc-resume.md | Removed dead code; added Pre-flight/Error Handling/Hard Stop Gates sections; replaced Python with Go/shell approach | E07IEXXYT, E07IEOJKN |
| M07IDBF29 | 2026-07-09 | Closeout checker gate statuses are keyword-enforced — specNote rejected "not-applicable"; only specific keywords allowed | Used "deferred" instead; check `defaultReleaseGateStatuses` before populating review.json releaseReadiness | E07IFL9C2 |
| M07IDBF29 | 2026-07-09 | review.json findings missing `id`/`summary` fields cause "unnamed" blocking findings | Populated `id` and `summary` per Finding struct schema for proper blocking-findings identification | E07IFL9C2 |

---

## Lessons

General principles and transferable insights — applicable beyond this codebase. Emerged from mission work but framed as world-wide solutions.

| Mission | Date | Lesson | Why it matters |
|---------|------|--------|----------------|
| M07IDBF29 | 2026-07-09 | Always verify pre-existing test failures against the base commit before treating them as blocking regressions | CI failures are not always introduced by the current change. A failing test that also fails on the base commit is pre-existing noise — document it, don't block the release on it |
| M07IDBF29 | 2026-07-09 | Resolved findings should not block release — severity downgrade prevents false blocking | A finding marked "critical" still blocks even after resolution. When a finding is resolved, either downgrade severity or mark `blocksShip: false` so the release gate reflects reality, not history |

---

## Patterns

Recurring solutions that proved effective across missions.

| Pattern | First Seen (Mission) | Description | When to Apply |
|---------|---------------------|-------------|---------------|
<!-- Emerge from repeated lessons; deduplicate manually -->
