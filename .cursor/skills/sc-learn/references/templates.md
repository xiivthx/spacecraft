# Mission Tracking File Templates

> Copy when creating `.space/missions/<id>/issues.md`, `solved.md`, `learned.md`.

## issues.md

```markdown
# Issues - <mission-title>

> Mission: <mission-id>. Ready/ship require 0 open.
> Policy: sc-learn. Drain: sc-run.

---
```

### Issue entry

```markdown
### <short title>
- **Date**: YYYY-MM-DD
- **Severity**: critical | important | minor
- **Class**: regression | consequence | related | unrelated | preexisting
- **Status**: open
- **Source**: <task id, review finding, or discovery context>
- **Description**: <what was found>
- **Impact**: <what it affects>
```

Filed (unrelated/preexisting only):

```markdown
- **Status**: filed
- **GitHub**: <#N or URL>
```

---

## solved.md

```markdown
# Solved - <mission-title>

> Mission: <mission-id>. Migrated to .space/trust/ on ship (local, gitignored).

---
```

```markdown
### <same title as original issue>
- **Date found**: YYYY-MM-DD
- **Date solved**: YYYY-MM-DD
- **Severity**: critical | important | minor
- **Solution**: <how it was fixed>
- **Verification**: <how the fix was verified>
- **Commit**: <hash or reference>
```

---

## learned.md

```markdown
# Lessons Learned - <mission-title>

> Mission: <mission-id>. Migrated to .space/trust/ on ship (local, gitignored).

---
```

```markdown
### <lesson title>
- **Date**: YYYY-MM-DD
- **Context**: <what triggered this insight>
- **Lesson**: <the principle or pattern learned>
- **Application**: <how to apply this in future missions>
```
