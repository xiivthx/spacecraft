---
name: sc-browser-probe
description: "Live browser probe after implement/fix or sc-run UI/workflow. Inventories the running product, sweeps Foundations + matched packs, then AFK fix-loops (find → Task fix → re-probe) until CLEAN or 3-cycle handback. Optional throughput/perf estimate. Use when user asks browser test, full-system probe, post-run live sweep, tags/sec or rate feasibility."
---

# sc-browser-probe

## Goal

Escape net after implement/fix or inside `/sc-run` when UI or multi-step workflow exists. Live product sweep: inventory surfaces, score Foundations always, deep-sweep matched packs only. Catch UX/UI/functional/workflow/responsive escapes unit/plan verify missed. Not happy-path smoke. **AFK fix-loop by default:** every finding (critical / important / minor) fixed and re-probed until `PROBE: CLEAN` or 3-cycle / timebox / blocked stop. Optional throughput/perf when asked. Does not replace sc-verification or sc-judge.

## Output

Report structure: [references/report-template.md](references/report-template.md).

Verdict line (required): `PROBE: CLEAN` | `PROBE: ISSUES` | `PROBE: PARTIAL` | `PROBE: BLOCKED`

## Good / Bad

- Good: inventory before sweep; every in-scope foundation + matched pack scored; interaction not screenshot-only; severity + repro; fix every finding then re-probe; stop on CLEAN or same-issue 3-cycle / timebox / blocked; measured perf or explicit "not measured"
- Bad: report-only when runnable product exists; one happy path; `CLEAN` with any finding or deferred required packs; walking unmatched packs / catalog; browser / MCP / chat as `ready` / `VERIFIED` / `AUTH` / ship; Commander product edits; skipping setup

## Verify

- Foundations all scored when in scope (`references/dimensions.md`)
- Every inventoried pack scored or `n/a` / `deferred:` (`references/surface-match.md`)
- Scenario minimums from `references/scenario-matrix.md` (plus extra buckets when that pack is present)
- Every finding has repro steps (or `blocked:` reason)
- Fix-loop until no findings remain, or stop reason recorded (`3-cycle:` / `timebox:` / `blocked:`)
- `PROBE: CLEAN` forbidden if any finding remains or any required coverage row is `deferred`
- Perf claims only with measured data or explicit estimate assumptions
- Spacecraft owns disposition: `PROBE: CLEAN` | `PROBE: ISSUES` | `PROBE: PARTIAL` | `PROBE: BLOCKED` (not MCP/chat alone)
- Browser / MCP / chat success Must not authorize `ready` / `VERIFIED` / `AUTH` / ship (sc-run may still require `PROBE: CLEAN` before ready when probe is in scope; probe alone never grants those)
- Task-delegate product fixes only (`sc-coder` / `sc-tester` / `sc-firmware` / `sc-rtl`); no Commander product edits

## When to use

- After implement/fix when user asks to probe / browser-test
- Inside `/sc-run`: Task(`sc-browser-probe`) after UI recheck / fix pass when UI or multi-step workflow touched - AFK, not recommend-only
- Full-system or feature live sweep; throughput / rate feasibility (perf module)

## When not to use

- Pure docs/prompt/config with no running UI
- Replacing `spacecraft evidence` / `sc-verification` / `sc-judge`
- Visual draft-parity gate alone → `sc-ux-design` Tier 3

## Workflow

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

- URL from user, or start the app (product origin, not API-only port when UI exists)
- Note auth, seed fixtures, feature flags
- Default viewports: 375 / 768 / 1280 (add 1536 if multi-region layout)
- Timebox: default 25m `feature:<name>`, 40m `full` (user override wins) - covers sweep + fix-loop
- Scope: `full` | `feature:<name>` from user call or sc-run

Foundations first, then packs by risk (forms / auth / destructive / money first). Timebox mid-sweep → `PROBE: PARTIAL` + deferred packs; do not claim `CLEAN`. Timebox mid-fix-loop → `PROBE: ISSUES` or `PARTIAL`, list remaining, hand back. No runnable URL → `PROBE: BLOCKED` (no fix-loop).

### 2. Inventory + scenario matrix

Inventory visible routes/surfaces; resolve ids via `.cursor/skills/sc-ux-design/references/checklists/README.md` and `references/surface-match.md`. Match several packs when needed. **Inventory before catalog** - do not walk unmatched packs or the full catalog.

Then follow `references/scenario-matrix.md`. Seed order: (1) user examples (2) spec / fixtures / Test Ideas in `decisions.md` (3) generated diversity labeled `gen:`.

Minimum when scope allows: happy + empty/invalid + boundary/long + one mobile pass. Extra buckets only when that pack is in inventory.

### 3. Run scenarios

Browser matrix (Cursor-first; sc-ux-design Tier 3 visual matrix unchanged):

1. **Cursor IDE browser** (`cursor-ide-browser` MCP) - primary
2. **`playwright-cli`** - fallback
3. **Chrome DevTools MCP** - last fallback

For each scenario: execute steps, record result, screenshot path on fail. Prefer interaction over passive screenshots.

### 4. Dimension sweep

1. Confirm inventory
2. Score **Foundations** always (`references/dimensions.md`)
3. Score **matched packs** only (`references/surface-match.md`)
4. File a finding only on `fail` with repro

Mark each check `ok` | `fail` | `n/a` | `deferred`.

### 5. Report

Fill [references/report-template.md](references/report-template.md) each sweep round (severity + verdict tables live there).

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

Exit: zero findings + no required `deferred` → `PROBE: CLEAN`; stopped with findings → `PROBE: ISSUES` (or `PARTIAL` if coverage deferred and findings empty); setup dead → `PROBE: BLOCKED`. Do not park findings for a later recommend handoff.

### 7. Perf probe (optional)

Only when user asked rate/throughput/feasibility **or** product documents a tags/sec (or similar) target. Follow `references/perf-probe.md`. Keep perf separate from functional `PROBE:` unless user makes rate a hard pass bar.

## Call template

```
/sc-browser-probe
target: <url | "start app">
scope: full | feature:<name>
examples: <optional comma list>
perf: none | <question e.g. "is 20 tags/sec feasible?">
timebox: 25m | 40m | <override>
```

Defaults: 25m `feature:<name>`, 40m `full`. Fix-loop always on when findings exist.

### Task from sc-run

```
Task(sc-browser-probe)
target: <live product URL>
scope: feature:<name> | full
mission: <id>
examples: <optional>
```

Return final report + verdict. sc-run Must not set `ready` while in-scope probe is not `PROBE: CLEAN`.

## Lifecycle placement

```
sc-run build → evidence (sc-verification)
            → UI recheck (when visual)
            → Task(sc-browser-probe) when UI/workflow  [AFK fix-loop → CLEAN]
            → review + sc-judge
            → ready → human → sc-ship
```

## Checklist

See **Verify** (SoT). Short track:

- [ ] Target runnable; inventory before catalog
- [ ] Foundations + matched packs scored; findings have repro
- [ ] Fix-loop to CLEAN or `3-cycle:` / `timebox:` / `blocked:`
- [ ] Verdict line; `CLEAN` forbidden with findings or required `deferred`
- [ ] No ready/`VERIFIED`/`AUTH`/ship from browser/MCP/chat or probe alone; Task-delegate fixes only

## References

- [report-template.md](references/report-template.md) - report markdown + severity/verdict
- [dimensions.md](references/dimensions.md) - Foundations sweep (always)
- [surface-match.md](references/surface-match.md) - inventory ids; load matched files from `.cursor/skills/sc-ux-design/references/checklists/`
- [scenario-matrix.md](references/scenario-matrix.md) - how to build cases
- [perf-probe.md](references/perf-probe.md) - throughput / machine estimate
- [examples.md](examples.md) - call examples
