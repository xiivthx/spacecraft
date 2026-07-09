# System Issues

Audit date: 2026-07-09. Covers all 12 skills, 9 commands, Go CLI, and docs.

---

## Skills (`.opencode/skills/sc-*/SKILL.md`)

### Fixed

| # | Skill | Issue | Fix |
|---|-------|-------|-----|
| S1 | 8 skills | Missing trigger phrases in description | Added trigger phrases to sc-clarify, sc-design, sc-git, sc-map, sc-mission, sc-planning, sc-verification, sc-web-service |
| S2 | 6 skills | Missing Research auto-trigger | Added to sc-mission, sc-verification, sc-solid, sc-design, sc-git, sc-tdd |
| S3 | sc-web-service | Workflow score 2 — no CLI commands, no edge cases | Rewrote: exact scaffold commands, npm commands, 5 edge cases |
| S4 | sc-planning | Workflow score 2 — no plan.json creation guidance, no task ID format | Rewrote: task ID guidance, plan.json commands, 5 edge cases, map integration |
| S5 | sc-clarify | No backtick commands, output format skeleton, no edge cases | Rewrote: question format, 5 edge cases (no response, deferred, contradiction) |
| S6 | sc-mission | Thin checklist (5 items), no edge cases, no research trigger | +Edge cases, +2 checklist items, +research trigger |
| S7 | sc-verification | No edge cases (81 lines), no research trigger | +5 edge cases (failed evidence, validation fail, manual checks, no plan, stale evidence), +research trigger |
| S8 | sc-tdd | Missing testing-strategy.md reference, no edge cases, no research trigger | +Edge cases, +research trigger, +testing-strategy.md ref |
| S9 | sc-solid | No research trigger, SOLID scan underspecified | +Research auto-trigger |
| S10 | sc-creator | Trimmed from 195→152, restored to 201 for operational precision | Final: 201 lines, all phases detailed |
| S11 | sc-debug | Description 328 chars (exceeded 200 limit) | Trimmed to 181 chars |
| S12 | All 12 skills | sc-tdd/sc-solid consolidation — content split across skills | Consolidated testing into sc-tdd (single source), removed tdd.md/testing.md from sc-solid |
| S13 | sc-solid, sc-tdd | Mutual cross-references (loop) and agent names in skill content | Removed loops, removed agent names, rewrote Out of scope sections |

### Open — minor polish

| # | Skill | Issue | Priority |
|---|-------|-------|----------|
| S14 | sc-debug | Rules section partially duplicated (General rules and Operating rules overlap) | Low |
| S15 | sc-debug | 260 lines — exceeds own length recommendation | Low (acceptable for its complexity) |
| S16 | sc-design | Rules section is 56% of file (100 lines) — many rules are taste-based, not verifiable | Low |
| S17 | sc-design | Thai-first rules are locale-specific — should be in a reference file | Low |
| S18 | sc-map | 363 lines — 3-4x recommended. Inline plan.json schema is 65 lines — should be a reference file | Low |
| S19 | sc-map | Phase 2 says "LLM analysis" but then says "commander performs this analysis" — slightly contradictory | Low |
| S20 | sc-map | References section mentions "Understand-Anything" and "Graphiti" — external tools not in codebase | Low |
| S21 | sc-mission | 103 lines — thin for a meta-skill that drives the entire lifecycle | Low |
| S22 | sc-solid | "Classes < 50 lines, methods < 10 lines" is an arbitrary threshold with no justification | Low |
| S23 | sc-git | Minor redundancy between rules sub-sections and post-merge/review-gate checklists | Low |
| S24 | sc-git | Bump policy complexity lives only in rules, not in workflow step 8 | Low |
| S25 | sc-tdd | "After all checks pass → Move to review" is abrupt — no connection to `/sc-review` trigger | Low |
| S26 | sc-creator | 201 lines — exceeds own 80-150 line recommendation. References section is flat (no descriptions) | Low |

---

## Commands (`.opencode/commands/*.md`)

### Fixed

| # | Command | Issue | Fix |
|---|---------|-------|-----|
| C1 | sc-design | Missing Hard Stop Gates and Error Handling | +Hard Stop Gates, +Error Handling (subagent blockers, DESIGN.md drift) |
| C2 | sc-git | Missing Hard Stop Gates and Error Handling, no edge cases | +Hard Stop Gates, +Error Handling (detached HEAD, dirty state, name collision) |
| C3 | sc-plan | No hard stop gates, no edge cases, no Research auto-trigger | +Edge cases (empty spec, stale map, task deps), +Research trigger |
| C4 | sc-review | No hard stop gates, no edge cases | +Edge cases (corrupt review.json, missing evidence, no subagent findings) |
| C5 | sc-start | No hard stop gates, no edge cases | +Edge cases (duplicate mission, no-git warning, empty args) |
| C6 | sc-build | No Research auto-trigger in dependency check | +Research auto-trigger reference |

### Open

| # | Command | Issue | Priority |
|---|---------|-------|----------|
| C7 | sc-quick | Missing Pre-flight Checks and Hard Stop Gates sections | Medium |
| C8 | sc-resume | Uses inline Python script to read last evidence entry — fragile across environments | Medium |
| C9 | sc-resume | Missing Pre-flight Checks and Error Handling sections | Medium |
| C10 | All 9 | Only 3/9 have all 4 sections (Pre-flight, Hard Stop Gates, Error Handling, Resolver Gate) | Low |
| C11 | All 9 | Research auto-trigger is only referenced in 3/9 commands | Low |

---

## Agents (`.opencode/agents/*.md`)

### Fixed

| # | Agent | Issue | Fix |
|---|-------|-------|-----|
| A1 | sc-planner | 21 lines, no structured sections, no map.json integration, no task limits, no plan.json schema | +Role/Context/Constraints/Output Format sections, +map.json integration, +≤7 rule, +plan.json schema |
| A2 | sc-coder | No SOLID/TDD awareness, no edge cases, no dependency research trigger | +sc-solid/sc-tdd guidance, +3 edge cases (no test, multiple checks, breaking tests), +research trigger |
| A3 | sc-tester | No evidence command, no seams concept, no edge cases | +exact evidence command syntax, +seams guidance, +4 edge cases (false pass, no acceptance checks, unknown framework, suite failure) |
| A4 | sc-commander | Only 1 constraint, missing sc-debug/sc-map/research auto-triggers | +4 constraints (skip gates, main writes, handoff merge, multi-question), +3 auto-triggers |
| A5 | sc-reviewer | No SOLID/code quality checks, no Kalama Sutta gate, no edge cases | +sc-solid integration, +Kalama Sutta gate, +3 edge cases |
| A6 | sc-designer | No structured sections (bare body text), no edge cases | +Role/Context/Constraints/Edge cases sections, +3 edge cases |

### Open

| # | Agent | Issue | Priority |
|---|-------|-------|----------|
| A7 | sc-commander | Auto-triggers section doesn't mention sc-tdd or sc-solid — though they load via sc-* wildcard | Low |
| A8 | sc-coder | "caveman-style brevity" instruction could be more specific about expected format | Low |
| A9 | sc-designer | Some rules are extremely locale-specific (Thai-first) — could be a reference, not agent body | Low |

---

## Go CLI (`scripts/src/`)

### Fixed

| # | File | Issue | Fix |
|---|------|-------|-----|
| G1 | research/deep.go | `nowISO()` always returned empty string — all deep analysis results lost their timestamp | Added `time` import, returns `time.Now().UTC().Format(time.RFC3339)` |
| G2 | resolver/resolver.go | Duplicate ID normalization — same regex as util/slug.go | Confirmed intentional (comment: "moved from util.go to avoid circular dep"). No fix needed. |
| G3 | archive/archive.go | Third copy of ID normalization (`normalizeMissionIdSimple`) | Open — could unify with util but has slightly different behavior |

### Open — critical

| # | File | Issue | Priority |
|---|------|-------|----------|
| G4 | archive, mission, state, util | Duplicate helper functions across 3+ packages: `fileExists`/`exists`, `readJSON`, `writeJSON` each implemented independently. `util/fs.go` has shared versions but they're unused by any internal package. | 🔴 High |
| G5 | archive, closeout | Duplicate `taskIsComplete` function — identical in `archive/archive.go:333` and `closeout/closeout.go:244` | 🔴 High |
| G6 | archive, closeout | Duplicate `blockingFindings`/`blockingReviewFindings` — nearly identical in `archive/archive.go:354` and `closeout/closeout.go:228` | 🔴 High |

### Open — medium

| # | File | Issue | Priority |
|---|------|-------|----------|
| G7 | main.go:922 | `researchCmd` uses `goto` — control flow anti-pattern. If code before the goto changes, behavior may be surprising. | 🟡 Medium |
| G8 | main.go | `checkDepsCmd` non-standard exit codes: 2 = updates available, 1 = error. Conflicts with Unix convention (1 = general error). Callers checking `$?` may misinterpret. | 🟡 Medium |
| G9 | main.go | `resolveCmd` returns inconsistent exit codes: 0 when no mission resolved but resolver succeeds, 1 only when "No mission resolved." is printed. Other commands use `requireResolved` which calls `os.Exit(1)`. | 🟡 Medium |
| G10 | main.go:1107 | `lookupPackage` silently returns `nil, ""` if no registry matches — caller checks `if pkg != nil` but there's no fallback error log for the nil case. | 🟡 Medium |
| G11 | main.go:704 | `evidenceCmd` always propagates subprocess exit code — if the evidence command fails (intentionally for capturing failures), the CLI exits non-zero. Could break shell pipelines. | 🟡 Medium |
| G12 | archive/archive.go:369 | `copyTextFile` silently fails — returns `bool` but callers ignore it. If spec.md or decisions.md fail to copy, archive is incomplete with no error. | 🟡 Medium |
| G13 | main.go:843 | `deep` variable name collision — used both as the `--deep` flag value and throughout the function for different purposes. | 🟡 Medium |
| G14 | workflow/workflow.go:249 | Dead code — empty if-block with comment "Check dirty state isn't from res". No TODO or issue reference. | 🟡 Medium |

### Open — low

| # | File | Issue | Priority |
|---|------|-------|----------|
| G15 | config/config.go | `NewConfig(root)` only — no ability to override individual paths. All paths derived from `root` alone. Testing forced to use temp directories. | 🟢 Low |
| G16 | types.go | 33 lines of backward-compat type aliases — could be removed now that refactoring is complete. | 🟢 Low |
| G17 | mission/model.go | `GitInfoData` and `GitInfo` — similar but different structures. `GitInfoData` has `Available`, `Dirty`, `DirtyFiles`; `GitInfo` only has `Branch`, `Sha`, `IsRepo`. | 🟢 Low |
| G18 | gitutil/git.go:70 | `NoopRunner` is misleading — declared as `var NoopRunner CommandRunner = OSCommandRunner{}` but `OSCommandRunner` actually runs real OS commands. No actual noop implementation exists. | 🟢 Low |
| G19 | main.go | `check-deps` and `research` flag parsing uses manual loop instead of `flag` package. Comment says "Go's flag package doesn't support dashes" — incorrect; `flag.String()` does support dashes. | 🟢 Low |

---

## Documentation (`docs/SPACECRAFT.md`)

### Fixed

| # | Issue | Fix |
|---|-------|-----|
| D1 | `sc-implementation` referenced in skill table but no directory exists on disk | Removed stale reference |
| D2 | `sc-testing` referenced in subagent table and skill table but no directory exists | Removed stale reference |
| D3 | `/sc-quick` command exists on disk but not listed in slash commands or routing table | Added to both |
| D4 | `sc-web-service` Used By column said only "sc-designer" but sc-coder also has it | Updated to "sc-coder, sc-designer" |
| D5 | CLI command table missing `current`, `research`, `check-deps`, `help` | Added all 4, restored dropped `archive`/`bind-branch`, alphabetized |

### Open

| # | Issue | Priority |
|---|-------|----------|
| D6 | CLI command table still missing some commands from code (minor flags-only variants) | Low |

---

## Summary

| Domain | Total | Fixed | Open |
|--------|-------|-------|------|
| Skills | 26 | 13 | 13 |
| Commands | 11 | 6 | 5 |
| Agents | 9 | 6 | 3 |
| Go CLI | 19 | 3 | 16 |
| Docs | 6 | 5 | 1 |
| **Total** | **71** | **33** | **38** |

### Open high-priority items

1. **G4**: Duplicate helper functions across packages (fileExists, readJSON, writeJSON)
2. **G5**: Duplicate `taskIsComplete` in archive and closeout
3. **G6**: Duplicate `blockingFindings` in archive and closeout
