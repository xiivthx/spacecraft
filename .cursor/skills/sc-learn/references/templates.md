# Mission Tracking File Templates

> Copy these templates when creating `.space/missions/<id>/issues.md`, `solved.md`, `learned.md`.

## issues.md

```markdown
# Issues - <mission-title>

> Mission: <mission-id>. Issues found during development.
> Created as GitHub Issues on ship if unresolved.

---
```

### Issue entry format

```markdown
### <short title>
- **Date**: YYYY-MM-DD
- **Severity**: critical | important | minor
- **Status**: open
- **Source**: <task id, review finding, or discovery context>
- **Description**: <what was found>
- **Impact**: <what it affects>
```

---

## solved.md

```markdown
# Solved - <mission-title>

> Mission: <mission-id>. Issues resolved during development.
> Migrated to docs/learned.md on ship.

---
```

### Solved entry format

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

> Mission: <mission-id>. Principles and patterns discovered.
> Migrated to docs/learned.md on ship.

---
```

### Lesson entry format

```markdown
### <lesson title>
- **Date**: YYYY-MM-DD
- **Context**: <what triggered this insight>
- **Lesson**: <the principle or pattern learned>
- **Application**: <how to apply this in future missions>
```
