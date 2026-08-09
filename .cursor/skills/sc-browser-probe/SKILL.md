---
name: sc-browser-probe
description: "Exploratory browser probe after implement/fix or sc-run green. Finds UX/UI/functional/workflow/responsive issues via real-browser scenarios; optional throughput/perf estimate. Use when user asks browser test, full-system probe, post-run smoke, tags/sec or rate feasibility."
---

# sc-browser-probe

## Goal

Recommend-only escape net after implement/fix or `/sc-run` green for UI/workflow. Exercise the **running product** in a real browser. Surface potential issues across UX, UI, functional, workflow, and responsive design. Optionally measure or estimate throughput when asked. Catch escapes that unit/plan verify missed. Does not replace sc-verification or sc-judge.

## Output

1. **Setup** - URL/env, auth, seed data, viewports
2. **Scenario table** - id, input, steps, result (`pass` | `fail` | `blocked`), notes
3. **Findings** - severity `critical` | `important` | `minor` with repro
4. **Perf** (only if requested or product has a rate target) - see `references/perf-probe.md`
5. **Verdict line** - `PROBE: CLEAN` | `PROBE: ISSUES` | `PROBE: BLOCKED`

## Good / Bad

- Good: real browser on live URL; diverse inputs; severity + repro; timebox; measured perf or explicit "not measured"; product-surface proof
- Bad: expertise cosplay; inventing metrics; domain hardcode in core steps; one happy path only; claiming `ready` / replacing sc-judge; skipping setup; screenshot-only with no interaction

## Verify

- ≥1 scenario each bucket when in scope: happy, empty/invalid, boundary/long-input, mobile viewport (see `references/scenario-matrix.md`)
- Every finding has repro steps (or `blocked:` reason)
- Perf claims only with measured data or explicit estimate assumptions
- Does **not** alone authorize mission `ready`

## When to use

- After implement/fix when user asks to probe / browser-test
- After `/sc-run` green: recommend-only handoff (not mandatory): `Recommend: /sc-browser-probe` when UI or multi-step workflow exists
- Full-system or feature smoke via browser
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
- [ ] 2. Build scenario matrix
- [ ] 3. Run scenarios (browser)
- [ ] 4. Dimension sweep
- [ ] 5. Report findings + verdict
- [ ] 6. Perf probe (if in scope)
```

### 1. Resolve target

- URL from user, or start the app and use the product origin (not API-only port when UI exists)
- Note auth, seed fixtures, feature flags
- Default viewports: 375 / 768 / 1280 (add 1536 if multi-region layout)
- Timebox: default 20m unless user sets otherwise
- Scope: `full` | `feature:<name>` from user call

If no runnable URL → `PROBE: BLOCKED` and stop.

### 2. Build scenario matrix

Follow `references/scenario-matrix.md`. Seed from (in order):

1. User-supplied examples
2. Spec / fixtures / Test Ideas in `decisions.md` when present
3. Generated diversity (short/long/empty/invalid/boundary) - label as agent-generated

Minimum when scope allows: happy + empty/invalid + boundary/long + one mobile pass.

### 3. Run scenarios

Browser matrix (same preference as spacecraft UI work):

1. **`playwright-cli`** - primary (`open` → interact → `snapshot` / `screenshot`)
2. **Cursor IDE browser** (`cursor-ide-browser` MCP) - fallback

For each scenario: execute steps, record result, capture screenshot path on fail. Prefer interaction over passive screenshots.

### 4. Dimension sweep

After scenarios (or interleaved), scan open gaps per `references/dimensions.md`. Do not essay - only file findings with repro.

### 5. Report

Use this template:

```markdown
## Browser probe
- Target: <url>
- Scope: <full | feature:…>
- Timebox: <n>m
- PROBE: CLEAN | ISSUES | BLOCKED

### Setup
- Env / auth / seed: …

### Scenarios
| id | input | result | notes |
|----|-------|--------|-------|
| S1 | … | pass/fail/blocked | … |

### Findings
- **critical** - <title> — repro: …
- **important** - …
- **minor** - …

### Perf
none | see section below
```

Severity:

| Level | Meaning |
|-------|---------|
| critical | Blocks core path or data loss / security-shaped |
| important | Wrong behavior or major UX/workflow friction |
| minor | Polish, copy, non-blocking layout |

`PROBE: CLEAN` only if no critical/important (minors may remain, list them).  
`PROBE: ISSUES` if any critical/important.  
`PROBE: BLOCKED` if setup failed.

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
timebox: 20m
```

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
- **Must**: Diverse scenarios per matrix minimums when in scope
- **Must**: Findings include severity + repro (or blocked reason)
- **Must**: Stay recommend-only escape net - recommend (not force) after UI/workflow `/sc-run` green via `Recommend: /sc-browser-probe`
- **Must not**: Expertise cosplay ("as a QA expert…")
- **Must not**: Invent throughput numbers
- **Must not**: Replace sc-verification or sc-judge; absorb into judge/verify as the sole gate; alone allow `ready`
- **Must not**: Hardcode product-specific fixtures into this skill core

## References

- [dimensions.md](references/dimensions.md) - UX/UI/functional/workflow/responsive sweep
- [scenario-matrix.md](references/scenario-matrix.md) - how to build cases
- [perf-probe.md](references/perf-probe.md) - throughput / machine estimate
- [examples.md](examples.md) - call examples (including NFC-shaped)

## Checklist

- [ ] Target URL runnable
- [ ] Scenario minimums met (or scoped skip noted)
- [ ] Dimension gaps filed as findings or N/A
- [ ] Verdict line emitted
- [ ] Perf only if in scope; measured or assumptions explicit
- [ ] No ready/ship claim from probe alone
