# Browser probe report template

Use each sweep round:

```markdown
## Browser probe
- Target: <url>
- Scope: <full | feature:…>
- Timebox: <n>m
- Round: <n>
- PROBE: CLEAN | ISSUES | PARTIAL | BLOCKED

### Setup
- Env / auth / seed: …

### Inventory
- Surfaces: <pack ids or none>

### Scenarios
| id | input | result | notes |
|----|-------|--------|-------|
| S1 | … | pass/fail/blocked | … |

### Coverage
| id | result | notes |
|----|--------|-------|
| foundations/ux | ok | |
| pack:<id> | fail | item scores; see findings |
| pack:<id> | deferred | timebox |

Pack rows: list each `- [ ]` item from the loaded catalog file as `ok` | `fail` | `n/a` in notes or a sub-list. `fail` → finding.

### Findings
- **critical** - <title> - repro: …
- **important** - …
- **minor** - …

### Fix-loop
none | Round <n>: fixed <titles> → re-probe …

### Perf
none | see section below
```

## Severity

| Level | Meaning |
|-------|---------|
| critical | Blocks core path or data loss / security-shaped; missing (state) on a pack item |
| important | Wrong behavior or major UX/workflow friction; missing chrome/path on a pack item |
| minor | Polish, copy, non-blocking layout |

## Verdict

- `PROBE: BLOCKED` - setup failed
- `PROBE: ISSUES` - any finding remains (critical / important / minor) after sweep or after stop; wins over PARTIAL; still list deferred rows
- `PROBE: PARTIAL` - required foundation/pack not scanned (timebox or blocked mid-sweep) and no findings yet
- `PROBE: CLEAN` - in-scope sweep finished, no required row `deferred`, **zero findings** (critical / important / minor)
