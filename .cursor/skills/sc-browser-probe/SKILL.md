---
name: sc-browser-probe
description: "Live browser probe after implement/fix or sc-run green. Inventories the running product, always sweeps Foundations, then deep-sweeps matched surface packs (UX/UI/functional/workflow/responsive). Optional throughput/perf estimate. Use when user asks browser test, full-system probe, post-run smoke, tags/sec or rate feasibility."
---

# sc-browser-probe

## Goal

Recommend-only escape net after implement/fix or `/sc-run` green for UI/workflow. A **live product sweep** of the running product in a real browser: inventory surfaces, score Foundations always, then deep-sweep only matched packs. Catch UX, UI, functional, workflow, and responsive escapes that unit/plan verify missed. Not a happy-path smoke. Optionally measure or estimate throughput when asked. Does not replace sc-verification or sc-judge.

## Output

1. **Setup** - URL/env, auth, seed data, viewports
2. **Inventory** - surfaces seen (pack ids from `references/surface-match.md`)
3. **Scenario table** - id, input, steps, result (`pass` | `fail` | `blocked`), notes
4. **Coverage table** - each in-scope foundation area and matched pack: `ok` | `fail` | `n/a` | `deferred`
5. **Findings** - severity `critical` | `important` | `minor` with repro
6. **Perf** (only if requested or product has a rate target) - see `references/perf-probe.md`
7. **Verdict line** - `PROBE: CLEAN` | `PROBE: ISSUES` | `PROBE: PARTIAL` | `PROBE: BLOCKED`

## Good / Bad

- Good: inventory before sweep; every in-scope foundation + matched pack scored; interaction not screenshot-only; consult site never used as pass/fail source; severity + repro; timebox; measured perf or explicit "not measured"
- Bad: one happy path; treating invalid / empty / loading / error as the same state; `CLEAN` with deferred required packs; copying external checklist text; treating designer draft-parity as this job; claiming `ready` / replacing sc-judge; skipping setup

## Verify

- Foundations all scored when in scope (`references/dimensions.md`)
- Every inventoried pack scored or `n/a` / `deferred:` (`references/surface-match.md`)
- Scenario minimums from `references/scenario-matrix.md` (plus extra buckets when that pack is present)
- Every finding has repro steps (or `blocked:` reason)
- `PROBE: CLEAN` forbidden if any required coverage row is `deferred`
- Perf claims only with measured data or explicit estimate assumptions
- Does **not** alone authorize mission `ready`

## When to use

- After implement/fix when user asks to probe / browser-test
- After `/sc-run` green: recommend-only handoff (not mandatory): `Recommend: /sc-browser-probe` when UI or multi-step workflow exists
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
- [ ] 6. Perf probe (if in scope)
```

### 1. Resolve target

- URL from user, or start the app and use the product origin (not API-only port when UI exists)
- Note auth, seed fixtures, feature flags
- Default viewports: 375 / 768 / 1280 (add 1536 if multi-region layout)
- Timebox: default 25m `feature:<name>`, 40m `full` (user override wins)
- Scope: `full` | `feature:<name>` from user call

Foundations first, then packs by risk (forms / auth / destructive / money first). Timebox cut → `PROBE: PARTIAL` + list deferred packs. Do not claim `CLEAN`.

If no runnable URL → `PROBE: BLOCKED` and stop.

### 2. Inventory + scenario matrix

Inventory visible routes/surfaces using ids in `references/surface-match.md`. One product may match several packs. Do not walk packs that are not in inventory. Do not walk a page-type catalog.

Then follow `references/scenario-matrix.md`. Seed from (in order):

1. User-supplied examples
2. Spec / fixtures / Test Ideas in `decisions.md` when present
3. Generated diversity (short/long/empty/invalid/boundary) - label as agent-generated (`gen:`)

Minimum when scope allows: happy + empty/invalid + boundary/long + one mobile pass. Add extra buckets only when that pack is in inventory.

### 3. Run scenarios

Browser matrix (same preference as spacecraft UI work):

1. **`playwright-cli`** - primary (`open` → interact → `snapshot` / `screenshot`)
2. **Cursor IDE browser** (`cursor-ide-browser` MCP) - fallback

For each scenario: execute steps, record result, capture screenshot path on fail. Prefer interaction over passive screenshots.

### 4. Dimension sweep

1. Confirm inventory
2. Score **Foundations** always (`references/dimensions.md`)
3. Score **matched packs** only (`references/surface-match.md`)
4. File a finding only on `fail` with repro

Mark each check `ok` | `fail` | `n/a` (capability absent) | `deferred` (not reached).

### 5. Report

Use this template:

```markdown
## Browser probe
- Target: <url>
- Scope: <full | feature:…>
- Timebox: <n>m
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

Pack rows: list each house item as `ok` | `fail` | `n/a` in notes or a sub-list. `fail` → finding.

### Findings
- **critical** - <title> - repro: …
- **important** - …
- **minor** - …

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
- `PROBE: ISSUES` - any critical/important (wins over PARTIAL; still list deferred rows)
- `PROBE: PARTIAL` - required foundation/pack not scanned (timebox or blocked mid-sweep) and no critical/important yet
- `PROBE: CLEAN` - in-scope sweep finished, no required row `deferred`, no critical/important (minors listed)

Hand findings back to `/sc-run` fix or `/sc-quick` - do not set mission `ready` from this skill.

### 6. Perf probe (optional)

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

Defaults: 25m `feature:<name>`, 40m `full`. User override wins.

## Lifecycle placement

```
sc-run build → evidence (sc-verification)
            → [recommend] sc-browser-probe
            → review + sc-judge
            → ready → human → sc-ship
```

After small fix: probe the touched surface only when user asks or handoff recommends.

## Must / Must not

- **Must**: Real product URL; record setup before scenarios
- **Must**: Inventory before sweep; score Foundations + every matched pack (or `n/a` / `deferred:`)
- **Must**: Diverse scenarios per matrix minimums when in scope, plus extra buckets when that pack is present
- **Must**: Findings include severity + repro (or blocked reason)
- **Must**: Stay recommend-only escape net - recommend (not force) after UI/workflow `/sc-run` green via `Recommend: /sc-browser-probe`
- **Must not**: Expertise cosplay ("as a QA expert…")
- **Must not**: Invent throughput numbers
- **Must not**: Copy external checklist text; use a consult site as pass/fail
- **Must not**: Claim `CLEAN` when any required coverage row is `deferred`
- **Must not**: Treat designer draft-parity as this job
- **Must not**: Replace sc-verification or sc-judge; absorb into judge/verify as the sole gate; alone allow `ready`
- **Must not**: Hardcode product-specific fixtures into this skill core

## References

- [dimensions.md](references/dimensions.md) - Foundations sweep (always)
- [surface-match.md](references/surface-match.md) - inventory ids + matched packs
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
- [ ] Verdict line emitted (`CLEAN` forbidden if any required row is `deferred`)
- [ ] Perf only if in scope; measured or assumptions explicit
- [ ] No ready/ship claim from probe alone
