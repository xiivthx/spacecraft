# Call examples

## Generic full-system

```
/sc-browser-probe
target: http://localhost:5173
scope: full
examples: empty form, valid sample, max-length string
perf: none
timebox: 40m
```

Report includes Inventory + Coverage. Fix-loop runs if any finding. Example snippet:

```
### Inventory
- Surfaces: form-submit, upload

### Coverage
| id | result | notes |
|----|--------|-------|
| foundations/ux | ok | |
| foundations/ui | ok | |
| foundations/functional | ok | |
| foundations/workflow | ok | |
| foundations/responsive | ok | |
| pack:form-submit | ok | 5/5 |
| pack:upload | deferred | timebox |

PROBE: PARTIAL
```

## Feature-scoped after fix

```
/sc-browser-probe
target: start app
scope: feature:checkout
examples: expired card then valid card
perf: none
timebox: 25m
```

## Persona walkthrough (opt-in)

```
/sc-browser-probe
target: http://localhost:5173
scope: feature:onboarding
persona: on
examples: first-time signup, keyboard-only tab through form
perf: none
timebox: 25m
```

Coverage includes `pack:persona-walkthrough` with archetype ids in notes. Journeys come from inventory/spec - not from the pack file. Findings use critical/important/minor only.

## Timebox PARTIAL

Sweep finished Foundations + `form-submit`. `upload` inventoried but not reached.

```
/sc-browser-probe
target: http://localhost:5173
scope: full
timebox: 40m
```

Coverage: Foundations `ok`; `pack:form-submit` `ok`; `pack:upload` `deferred: timebox`. No findings.

`PROBE: PARTIAL` - do not claim `CLEAN`.

## NFC / sim-shaped (product-specific - call site only)

```
/sc-browser-probe
target: <nfc roll sim URL>
scope: feature:nfc-roll
examples: short url, long url, trc sample
perf: is 20 tags/sec feasible? per-scenario rate + recommend speed / machine estimate
timebox: 25m
```

Do not bake NFC fixtures into SKILL.md core - pass them as `examples:` / `perf:`.

## From sc-run (AFK Task)

When `/sc-run` touches UI or multi-step workflow, after fix pass:

```
Task(sc-browser-probe)
target: <live product URL>
scope: feature:<name>
mission: <id>
```

Agent sweeps + fix-loops until `PROBE: CLEAN` (or stops on 3-cycle / timebox). sc-run does not ready until CLEAN (or probe skipped for non-UI work). Not a handoff recommend line.
