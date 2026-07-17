---
name: sc-learn
description: "Capture mission knowledge: issues, solutions, and lessons learned. Activate during ship migration, when recording findings during a mission, or on \"lesson learned\", \"what did we learn\", \"capture knowledge\". Also activated during implementation to record issues as discovered."
---

# sc-learn

Capture knowledge from missions: track issues found, solutions applied, and lessons learned. During ship, create GitHub Issues for unresolved items and migrate solved/learned items to `docs/learned.md` for internal research reuse.

## When to use

Activate when the user asks to:

- **"Record this issue" / "track this bug" / "note this finding"** - during a mission
- **"Mark as solved" / "this is fixed"** - after resolving an issue
- **"What did we learn?" / "lesson learned" / "capture knowledge"** - reflection
- During ship migration - create GitHub Issues for open items; migrate solved/learned to `docs/learned.md` before version bump
- During implementation - record issues and lessons as they are discovered

## Workflow

### During a mission

1. **Create tracking files** - When a mission moves to `planned` state, ensure these files exist in `.space/missions/<id>/`:
   - `issues.md` - open issues found during this mission
   - `solved.md` - issues resolved during this mission
   - `learned.md` - lessons extracted from this mission

   If they don't exist, create them from the templates below.

   After each write to these files, ctx_index them with source label `sc-memory/<mission-id>/<type>` (best-effort, non-blocking). See sc-memory for conventions.

2. **Record issues** - When a bug, gap, or finding is discovered:
   ```
   ### <short title>
   - **Date**: YYYY-MM-DD
   - **Severity**: critical | important | minor
   - **Status**: open
   - **Source**: <task id, review finding, or discovery context>
   - **Description**: <what was found>
   - **Impact**: <what it affects>
   ```

3. **Mark as solved** - When an issue is resolved, move it from `issues.md` to `solved.md`:
   ```
   ### <same title>
   - **Date found**: YYYY-MM-DD
   - **Date solved**: YYYY-MM-DD
   - **Severity**: critical | important | minor
   - **Solution**: <how it was fixed>
   - **Verification**: <how the fix was verified>
   - **Commit**: <hash or reference>
   ```

4. **Record lessons** - When a general principle or transferable insight emerges (NOT a specific issue - those go to solved.md):
    ```
    ### <lesson title>
    - **Date**: YYYY-MM-DD
    - **Context**: <what triggered this insight>
    - **Lesson**: <the general principle - framed as world-wide solution, not project-specific>
    - **Application**: <how this applies beyond the current mission>
    ```
    
    **Distinction**: Solved = specific bugs fixed in this project. Lessons = general truths reusable anywhere. A closeout checker quirk is solved; "verify pre-existing failures before blocking" is a lesson.

### During ship (migration)

Before the version bump and changelog commit, run this migration:

1. **Read mission files** - Load `.space/missions/<id>/issues.md`, `solved.md`, `learned.md`.

2. **Migrate unresolved issues** - For each issue in `issues.md` with status `open`:
   - Create a GitHub Issue with mission context: `gh issue create --title "<title>" --body "<body>" --label bug`
   - Body must include: mission id/title, date, severity, description, impact
   - Update the mission `issues.md` to mark each as `status: migrated` and record the GitHub issue URL/number
   - If the mission belongs to a roadmap, append the issue to the roadmap `issues` array

3. **Migrate solved items** - For each entry in `solved.md` (specific issues fixed):
    - Append to `docs/learned.md` under the `## Solved` table with mission context
    - Format: `| <mission-id> | <date> | <problem summary> | <solution summary> | <evidence> |`

4. **Migrate lessons** - For each entry in `learned.md` (general principles, not project-specific):
    - Reword from mission context to general principle before migrating - strip project-specific details, keep the transferable insight
    - Append to `docs/learned.md` under the `## Lessons` table
    - Format: `| <mission-id> | <date> | <lesson - general principle> | <why it matters - world-wide relevance> |`

5. **Proceed with ship** - After migration, continue with version bump and changelog as normal.

### Templates

See `references/templates.md` for the full file templates. Copy them when creating mission tracking files.

## Rules

- **Must**: Create tracking files when mission moves to `planned` state (or earlier if issues arise).
- **Must**: Record issues as they are found - don't batch at the end.
- **Must**: Move issues from `issues.md` to `solved.md` when resolved during the mission.
- **Must**: Distinguish Solved (specific issues fixed in this project) from Lessons (general principles, transferable to any codebase). If an insight only makes sense in the context of this specific tool, it's a solved issue, not a lesson.
- **Must**: During migration, reword lesson entries from project-specific context into general principles before writing to `docs/learned.md`.
- **Must**: During ship, create GitHub Issues for unresolved mission issues before the version bump commit.
- **Must**: During ship, migrate solved and learned items to `docs/learned.md` before the version bump commit.
- **Must not**: Ship with unresolved issues still only in the mission folder - they must be promoted to GitHub Issues.
- **Must**: Use GitHub Issues as the sole global issue registry.
- **Must not**: Delete mission tracking files after migration - archive them with the mission.
- **Must**: After writing to issues.md, solved.md, or learned.md, ctx_index the file with source label `sc-memory/<mission-id>/<type>` (best-effort, non-blocking -- warn on failure). See sc-memory for label format and conventions.

## Out of scope

- Git operations - use sc-git
- Mission lifecycle - use sc-mission
- Code quality issues - use sc-solid for detection, sc-learn for tracking
- Automated lesson extraction from git diffs - manual reflection only (for now)

## Output format

```
Mission: <mission-id>
Issues: [N] open → GitHub Issues created (#N, #M, ...)
Solved: [N] items → migrated to docs/learned.md
Lessons: [N] items → migrated to docs/learned.md
Migration complete. Ready for version bump.
```

## Checklist

Before claiming knowledge migration is done:

- [ ] Mission tracking files exist (issues.md, solved.md, learned.md)
- [ ] All open issues created as GitHub Issues with mission context
- [ ] All solved items migrated to `docs/learned.md` with mission context
- [ ] All lessons migrated to `docs/learned.md` with mission context
- [ ] Mission files marked as `status: migrated` (not deleted)
- [ ] Proceed to version bump and changelog after migration

---

## References

- `references/templates.md` - mission tracking file templates (issues.md, solved.md, learned.md)
- GitHub Issues - global issue registry (via `gh issue create`)
- `docs/learned.md` - global lessons learned (aggregation target, internal research source)
- sc-memory - cross-mission memory conventions (ctx_index/ctx_search wrapping)
- `.space/missions/<id>/issues.md` - per-mission issue tracking
- `.space/missions/<id>/solved.md` - per-mission resolved items
- `.space/missions/<id>/learned.md` - per-mission lessons
