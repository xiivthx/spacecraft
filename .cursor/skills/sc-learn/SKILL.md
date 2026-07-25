---
name: sc-learn
description: "Capture mission knowledge: issues, solutions, and lessons learned. Activate during ship migration, when recording findings during a mission, or on \"lesson learned\", \"what did we learn\", \"capture knowledge\". Also activated during implementation to record issues as discovered."
---

# sc-learn

Mission knowledge: issues, solutions, lessons. **Sole owner of issues triage policy.** Ready/ship require **0 open** issues and empty review findings.

**Local source of trust** (gitignored with the rest of `.space/`):

```
.space/trust/lessons.md   # must-read for agents
.space/trust/solved.md    # project-specific fixes
```

Tracked seed only: `references/trust-seed/` (copy into `.space/trust/` when missing). Never write trust to `docs/`.

## When to use

- Record / solve / file a finding during a mission
- Reflect ("what did we learn?")
- Ship migration: confirm 0 open; migrate solved/learned into `.space/trust/` before version bump
- Ensure trust exists (seed if missing) before discuss/plan/run

## Ensure trust (before read or migrate)

1. If `.space/trust/lessons.md` is missing: `mkdir -p .space/trust` and copy from `references/trust-seed/lessons.md` and `solved.md` (or create empty tables with the headers in those seeds).
2. Agents **must read** `.space/trust/lessons.md` before inventing process in `/sc-discuss` or `/sc-run` plan/build. Skim for matching lessons; do not invent conflicting process.

## Issues policy (source of truth)

**Two-rule summary** (for always-on rules / sc-run):

1. **Mission-caused** (`related` / `regression` / `consequence`): **must fix** in `/sc-run` drain → `solved.md`. Never file-away.
2. **Not mission-caused** (`unrelated` / `preexisting`): **file** during run (`gh issue create` → `filed`) unless **suite-breaking**, on **touched path**, or **trivial** one-shot fix - then fix → `solved.md`.

Ready/ship: 0 `Status: open` (any severity, including minor). Drain loop mechanics: **sc-run**.

### Full drain matrix

| Class | Severity | Drain action |
|-------|----------|--------------|
| `regression` / `consequence` / `related` | any | **Must fix** → `solved.md`. Never file-away. |
| `unrelated` / `preexisting` | critical / important | **Fix** if suite-breaking, on touched path, or small/local; else **file** → `filed`. |
| `unrelated` / `preexisting` | minor | **File** by default; **fix** only if trivial one-shot. |

No soft "prefer fix when uncertain." Suite-breaking or touched-path ⇒ fix; else file.

During drain: new findings **append** and stay in the loop until 0 open. Review/judge hits become new open entries and re-enter drain (sc-run).

## Workflow

### During a mission

1. **Create tracking files** when mission moves to `planned` (or earlier): `issues.md`, `solved.md`, `learned.md` under `.space/missions/<id>/`. Templates: `references/templates.md`. After each write, ctx_index with `sc-memory/<mission-id>/<type>` (best-effort).

2. **Classify and record** when a finding appears:
   ```
   ### <short title>
   - **Date**: YYYY-MM-DD
   - **Severity**: critical | important | minor
   - **Class**: regression | consequence | related | unrelated | preexisting
   - **Status**: open
   - **Source**: <task id, review finding, or discovery context>
   - **Description**: <what was found>
   - **Impact**: <what it affects>
   ```

3. **Mark as solved** - move from `issues.md` to `solved.md`:
   ```
   ### <same title>
   - **Date found**: YYYY-MM-DD
   - **Date solved**: YYYY-MM-DD
   - **Severity**: critical | important | minor
   - **Solution**: <how it was fixed>
   - **Verification**: <how the fix was verified>
   - **Commit**: <hash or reference>
   ```

4. **Mark as filed** (unrelated/preexisting only):
   ```
   - **Status**: filed
   - **GitHub**: <#N or URL>
   ```

5. **Record lessons** (general principles, not specific bugs) in mission `learned.md`. Solved = project-specific fix; Lessons = transferable principle.

### Clean-slate gate

Before `set-state ready` and before ship release:

- `issues.md`: **0** `Status: open` (`spacecraft closeout-check`)
- `review.json` `findings`: **empty**
- Do not leave open items for ship to turn into GitHub Issues

### During ship (migration)

1. If any issue is `open` → **block ship**. Fix or file in `/sc-run`, then re-ready.
2. Ensure `.space/trust/` exists (seed if missing).
3. **Quality bar (must)** - do not dump mission noise into trust:
   - **Lessons:** only general, reusable principles. Skip duplicates of rows already in `.space/trust/lessons.md`. Skip mission-diary / one-off UI taste / harness trivia. Prefer ≤1 new lesson per ship unless clearly distinct.
   - **Solved:** only regressions or fixes others would hit again. Skip nits, typos, and issues already covered by an existing lesson. Prefer empty append over pad.
   - Reword lessons to world-wide principles before append; keep rows short (one line each cell).
4. Migrate qualifying mission `solved.md` → `.space/trust/solved.md`: `| <mission-id> | <date> | <problem> | <solution> | <evidence> |`
5. Migrate qualifying mission `learned.md` → `.space/trust/lessons.md`: `| <date> | <lesson> | <why it matters> | <mission-id> |`
6. Continue version bump + changelog. Do not delete mission tracking files (archive with mission). Trust stays local/gitignored - do not commit `.space/`.

## Rules

- **Must**: Create tracking files at `planned` (or earlier); record findings as found; classify each.
- **Must**: Follow the issues policy matrix above; clear all `open` before ready/ship.
- **Must**: Migrate only high-signal solved/learned into `.space/trust/` at ship; reword lessons; skip duplicates and nits.
- **Must**: Read `.space/trust/lessons.md` before inventing process (seed if missing). Keep trust files short.
- **Must not**: Ship or ready with open issues; create GitHub Issues at ship for open debt; write trust into `docs/`; pad trust with mission diary noise.
- **Must**: GitHub Issues = global registry for debt filed during run.
- **Must**: ctx_index after mission tracking writes (best-effort).

## Output format

```
Mission: <mission-id>
Issues: 0 open (clean)
Solved: [N] items → .space/trust/solved.md
Lessons: [N] items → .space/trust/lessons.md
Migration complete. Ready for version bump.
```

## Checklist

- [ ] Tracking files exist
- [ ] Zero open issues (closeout-check)
- [ ] Trust dir present (seeded if needed)
- [ ] Solved + lessons migrated to `.space/trust/`
- [ ] Proceed to version bump / changelog

## References

- `references/templates.md` - mission tracking templates
- `references/trust-seed/` - tracked seed for local trust
- `.space/trust/` - local source of trust (gitignored)
- sc-run - issues drain loop
- sc-memory - ctx_index conventions
