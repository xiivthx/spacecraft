---
name: sc-browser-probe
description: "Live browser probe after implement/fix or sc-run UI/workflow. Inventories the running product, sweeps Foundations + matched packs, then AFK fix-loops (find → Task fix → re-probe) until CLEAN or 3-cycle handback. Optional throughput/perf estimate. Use when user asks browser test, full-system probe, post-run live sweep, tags/sec or rate feasibility."
---

# sc-browser-probe

## Goal

Escape net after implement/fix or inside `/sc-run` when UI or multi-step workflow exists. A **live product sweep** of the running product in a real browser: inventory surfaces, score Foundations always, then deep-sweep only matched packs. Catch UX, UI, functional, workflow, and responsive escapes that unit/plan verify missed. Not a happy-path smoke. **AFK fix-loop by default:** every finding (critical / important / minor) is fixed and re-probed until `PROBE: CLEAN` or the 3-cycle stop. Optionally measure or estimate throughput when asked. Does not replace sc-verification or sc-judge.

## Output

1. **Setup** - URL/env, auth, seed data, viewports
2. **Inventory** - surfaces seen (ids from `.cursor/skills/sc-ux-design/references/checklists/README.md` via `references/surface-match.md`)
3. **Scenario table** - id, input, steps, result (`pass` | `fail` | `blocked`), notes
4. **Coverage table** - each in-scope foundation area and matched pack: `ok` | `fail` | `n/a` | `deferred`
5. **Findings** - severity `critical` | `important` | `minor` with repro (empty when CLEAN)
6. **Fix-loop log** - rounds, what fixed, re-probe result (omit when first pass CLEAN)
7. **Perf** (only if requested or product has a rate target) - see `references/perf-probe.md`
8. **Verdict line** - `PROBE: CLEAN` | `PROBE: ISSUES` | `PROBE: PARTIAL` | `PROBE: BLOCKED`

## Good / Bad

- Good: inventory before sweep; every in-scope foundation + matched pack scored; interaction not screenshot-only; pack items from `sc-ux-design/references/checklists/` via `surface-match.md`; severity + repro; fix every finding then re-probe; stop on CLEAN or same-issue 3-cycle; timebox; measured perf or explicit "not measured"
- Bad: report-only when runnable product exists; one happy path; treating invalid / empty / loading / error as the same state; `CLEAN` with any finding or deferred required packs; walking unmatched packs; treating designer draft-parity as this job; claiming `ready` / replacing sc-judge; skipping setup; infinite fix without 3-cycle stop

## Verify

- Foundations all scored when in scope (`references/dimensions.md`)
- Every inventoried pack scored or `n/a` / `deferred:` (`references/surface-match.md`)
- Scenario minimums from `references/scenario-matrix.md` (plus extra buckets when that pack is present)
- Every finding has repro steps (or `blocked:` reason)
- Fix-loop ran until no findings remain, or stop reason recorded (`3-cycle:` / `timebox:` / `blocked:`)
- `PROBE: CLEAN` forbidden if any finding remains or any required coverage row is `deferred`
- Perf claims only with measured data or explicit estimate assumptions
- Does **not** alone authorize mission `ready` (sc-run may require CLEAN before ready when probe is in scope)

## When to use

- After implement/fix when user asks to probe / browser-test
- Inside `/sc-run`: Task(`sc-browser-probe`) after UI recheck / fix pass when UI or multi-step workflow touched - AFK, not a handoff recommend
- Full-system or feature live sweep via browser
- Throughput / rate feasibility questions (perf module)

## When not to use

- Pure docs/prompt/config with no running UI
- Replacing `spacecraft evidence` / `sc-verification` / `sc-judge`
- Visual draft-parity gate alone → use `sc-ux-design` Tier 3

## Workflow

Copy and track:

```
Probe progress:
- [ ] 1. Resolve target
- [ ] 2. Inventory + scenario matrix
- [ ] 3. Run scenarios (browser)
- [ ] 4. Foundations + matched packs
- [ ] 5. Report findings + verdict
- [ ] 6. Fix-loop (if any finding)
- [ ] 7. Perf probe (if in scope)
```

### 1. Resolve target

- URL from user, or start the app and use the product origin (not API-only port when UI exists)
- Note auth, seed fixtures, feature flags
- Default viewports: 375 / 768 / 1280 (add 1536 if multi-region layout)
- Timebox: default 25m `feature:<name>`, 40m `full` (user override wins) - covers sweep + fix-loop
- Scope: `full` | `feature:<name>` from user call or sc-run

Foundations first, then packs by risk (forms / auth / destructive / money first). Timebox cut mid-sweep → `PROBE: PARTIAL` + list deferred packs. Do not claim `CLEAN`. Timebox mid-fix-loop → keep `PROBE: ISSUES` or `PARTIAL`, list remaining findings, hand back.

If no runnable URL → `PROBE: BLOCKED` and stop (no fix-loop).

### 2. Inventory + scenario matrix

Inventory visible routes/surfaces; resolve ids via `.cursor/skills/sc-ux-design/references/checklists/README.md` and `references/surface-match.md`. One product may match several packs. Do not walk packs that are not in inventory. Do not walk the catalog.

Then follow `references/scenario-matrix.md`. Seed from (in order):

1. User-supplied examples
2. Spec / fixtures / Test Ideas in `decisions.md` when present
3. Generated diversity (short/long/empty/invalid/boundary) - label as agent-generated (`gen:`)

Minimum when scope allows: happy + empty/invalid + boundary/long + one mobile pass. Add extra buckets only when that pack is in inventory.

### 3. Run scenarios

Browser matrix (same preference as spacecraft UI work):

1. **Chrome DevTools MCP** (Antigravity native: `navigate_page`, `click`, `fill`, `type_text`, `take_screenshot`, `list_console_messages`, `list_network_requests`, `lighthouse_audit`, `resize_page`, `evaluate_script`)
2. **`playwright-cli`** - standalone CLI (`open` → interact → `snapshot` / `screenshot`)
3. **Cursor IDE browser** (`cursor-ide-browser` MCP) - Cursor fallback

For each scenario: execute steps, record result, capture screenshot path on fail. Prefer interaction over passive screenshots.

### 4. Dimension sweep

1. Confirm inventory
2. Score **Foundations** always (`references/dimensions.md`)
3. Score **matched packs** only (`references/surface-match.md`)
4. File a finding only on `fail` with repro

Mark each check `ok` | `fail` | `n/a` (capability absent) | `deferred` (not reached).

### 5. Report (each sweep round)

Use this template:

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

Severity:

| Level | Meaning |
|-------|---------|
| critical | Blocks core path or data loss / security-shaped; missing (state) on a pack item |
| important | Wrong behavior or major UX/workflow friction; missing chrome/path on a pack item |
| minor | Polish, copy, non-blocking layout |

Verdict:

- `PROBE: BLOCKED` - setup failed
- `PROBE: ISSUES` - any finding remains (critical / important / minor) after sweep or after stop; wins over PARTIAL; still list deferred rows
- `PROBE: PARTIAL` - required foundation/pack not scanned (timebox or blocked mid-sweep) and no findings yet
- `PROBE: CLEAN` - in-scope sweep finished, no required row `deferred`, **zero findings** (critical / important / minor)

### 6. Fix-loop (default; AFK)

When any finding exists and URL is runnable:

```
while findings nonempty AND not stopped:
  1. Prefer critical → important → minor
  2. Task(`sc-tester`) when a failing test is the right gate; else Task(`sc-coder`/`sc-firmware`/`sc-rtl`) with repro + requiredFix
  3. Commander Must not write product code/tests - only Task-delegate
  4. Re-run focused scenarios + failed coverage rows (full re-sweep when unclear)
  5. Same root issue fails verify **3** times → stop; hand human with repro + attempts
  6. Timebox exhausted → stop; report remaining findings
```

Exit:

- Zero findings + no required `deferred` → `PROBE: CLEAN`
- Stopped with findings left → `PROBE: ISSUES` (or `PARTIAL` if coverage still deferred and findings empty)
- Setup dead → `PROBE: BLOCKED`

Do not park findings for a later recommend handoff.

### 7. Perf probe (optional)

Run only when:

- User asked rate/throughput/feasibility, **or**
- Product documents a tags/sec (or similar) target

Follow `references/perf-probe.md`. Keep perf separate from functional `PROBE:` unless user makes rate a hard pass bar.

## Call template

```
/sc-browser-probe
target: <url | "start app">
scope: full | feature:<name>
examples: <optional comma list>
perf: none | <question e.g. "is 20 tags/sec feasible?">
timebox: 25m | 40m | <override>
```

Defaults: 25m `feature:<name>`, 40m `full`. User override wins. Fix-loop is always on when findings exist.

### Task from sc-run

```
Task(sc-browser-probe)
target: <live product URL>
scope: feature:<name> | full
mission: <id>
examples: <optional>
```

Agent follows this skill end-to-end (sweep + fix-loop). Return final report + verdict. sc-run Must not set `ready` while in-scope probe is not `PROBE: CLEAN` (hand human on `ISSUES` / `PARTIAL` / `BLOCKED` after stop).

## Lifecycle placement

```
sc-run build → evidence (sc-verification)
            → UI recheck (when visual)
            → Task(sc-browser-probe) when UI/workflow  [AFK fix-loop → CLEAN]
            → review + sc-judge
            → ready → human → sc-ship
```

After small fix outside mission: run this skill on the touched surface when user asks.

## Must / Must not

- **Must**: Real product URL; record setup before scenarios
- **Must**: Inventory before sweep; score Foundations + every matched pack (or `n/a` / `deferred:`)
- **Must**: Diverse scenarios per matrix minimums when in scope, plus extra buckets when that pack is present
- **Must**: Findings include severity + repro (or blocked reason)
- **Must**: Fix-loop every finding (including minor) until CLEAN, or stop with `3-cycle:` / `timebox:` / `blocked:`
- **Must**: Task-delegate product fixes (`sc-coder` / `sc-tester` / `sc-firmware` / `sc-rtl`); no Commander product edits
- **Must**: When invoked from `/sc-run` on UI/workflow work, run AFK via Task(`sc-browser-probe`) - no recommend-only handoff
- **Must not**: Treat quality-tier / mutation out-of-scope as a probe waive — UI/workflow still requires `PROBE: CLEAN` (non-UI skip only)
- **Must not**: Expertise cosplay ("as a QA expert…")
- **Must not**: Invent throughput numbers
- **Must not**: Use an external checklist site as pass/fail; walk packs that are not in inventory
- **Must not**: Claim `CLEAN` when any finding remains or any required coverage row is `deferred`
- **Must not**: Treat designer draft-parity as this job
- **Must not**: Replace sc-verification or sc-judge; absorb into judge/verify as the sole gate; alone allow `ready`
- **Must not**: Hardcode product-specific fixtures into this skill core
- **Must not**: End with "recommend probe later" when a runnable UI/workflow target exists

## References

- [dimensions.md](references/dimensions.md) - Foundations sweep (always)
- [surface-match.md](references/surface-match.md) - inventory ids; load matched files from `.cursor/skills/sc-ux-design/references/checklists/`
- [scenario-matrix.md](references/scenario-matrix.md) - how to build cases
- [perf-probe.md](references/perf-probe.md) - throughput / machine estimate
- [examples.md](examples.md) - call examples (including NFC-shaped)

## Checklist

- [ ] Target URL runnable
- [ ] Inventory recorded
- [ ] Scenario minimums met (plus pack extras when present, or scoped skip noted)
- [ ] Foundations all scored when in scope
- [ ] Every inventoried pack scored or `n/a` / `deferred:`
- [ ] Dimension `fail`s filed as findings with repro
- [ ] Fix-loop until zero findings or stop reason recorded
- [ ] Verdict line emitted (`CLEAN` forbidden if findings remain or required row `deferred`)
- [ ] Perf only if in scope; measured or assumptions explicit
- [ ] No ready/ship claim from probe alone
