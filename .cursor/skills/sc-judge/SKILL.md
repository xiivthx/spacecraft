---
name: sc-judge
description: "Adversarial prove gate before ready. Treat completion claims as claims; re-run evidence; diff scope vs plan; hunt weakened tests / false completion / unauthorized action. Verdict VERIFIED | VERIFIED WITH CAVEATS | REFUTED."
---

# sc-judge

## Goal

Adversarial prove gate before `ready`: re-observe claimed completion. Never trust a report, summary, or "done" claim alone. Preserve `/sc-discuss` → `/sc-run` → `/sc-ship`; judge proves readiness - it does not replace the lifecycle.

## Output

A single verdict for the mission (or scoped claim under review):

```
VERDICT: VERIFIED | VERIFIED WITH CAVEATS | REFUTED
```

Plus: fresh evidence ids, scope-diff notes, hunt findings, and (when not `VERIFIED`) the caveats or refutation reasons. Cursor-native skill only - no Claude-plugin dependency.

## Good / Bad

- Good: treat every completion / "done" / "ready" claim as a claim; re-run claimed evidence commands and record fresh observation; diff actual change scope vs `plan.json` / spec acceptance; hunt weakened tests, false completion, unauthorized action; emit exactly one of the three verdicts; block `ready` on `REFUTED`
- Bad: believing a prior report or evidence.jsonl line without re-running; inventing evidence; soft-shipping past `REFUTED`; replacing discuss/run/ship with a judge-only flow; expanding hunt into product redesign or trap-eval suites

## Verify

Before accepting a judge pass as proof:

```
spacecraft evidence "judge-<mission-id>" -- <re-run of claimed commands>
```

Confirm verdict string is exactly one of `VERIFIED` | `VERIFIED WITH CAVEATS` | `REFUTED`, hunt notes cover the three targets, and scope-diff vs plan/spec is recorded.

## When to use

Activate when:

- A mission claims build complete or moves toward `ready`
- Review / `/sc-run` needs an adversarial prove step before `set-state ready`
- A completion, "done", or "ready" claim must be re-observed
- Commander or reviewer asks to judge evidence and scope against plan

## Workflow

Use this exact sequence unless the user specifies otherwise:

1. **Resolve the mission** - Run `spacecraft resolve`. On conflict or ambiguity, use `spacecraft use <selector>`.
2. **Collect claims** - Read the completion claim(s): task notes, evidence labels, review draft, or "ready" request. Treat each as a claim, not proof.
3. **Re-run claimed evidence** - For every command cited as proving acceptance, re-run it via `spacecraft evidence "<label>-judge" -- <command>`. Record the fresh observation. Never reuse a stale evidence line as the sole proof.
4. **Diff scope vs plan** - Compare the actual change set (diff / files touched) to `plan.json` tasks and `spec.md` acceptance. Flag work outside plan, missing acceptance coverage, or plan items marked done without matching fresh evidence.
5. **Hunt** - Actively search for:
   - **weakened tests** - assertions removed, skipped, loosened, or replaced with tautologies so GREEN is cheap
   - **false completion** - "done"/"ready" claimed while acceptance fails, evidence is missing/stale, scope does not match plan, `issues.md` still has `Status: open` (any severity), related/regression/consequence findings parked instead of fixed, or `review.json` still has findings (including warnings)
   - **unauthorized action** - outward push/deploy/publish/send (or similar) without quoted `AUTH:` and user authorization; ship/merge without `/sc-ship` gates
6. **Verdict** - Emit exactly one of:
   - `VERIFIED` - fresh evidence passes; scope matches plan/spec; hunt clean; 0 open mission issues; 0 review findings
   - `VERIFIED WITH CAVEATS` - prove holds with non-defect notes already recorded in `decisions.md` (e.g. known manual follow-up outside product). **Never** for warnings, open `issues.md`, review findings, or acceptance gaps - those are `REFUTED`
   - `REFUTED` - material hunt hit, failed re-run, scope/acceptance mismatch, open mission issues, or leftover review findings
7. **Ready gate** - If `REFUTED`, block `ready`. Do not soften to ship or ready.

### Edge cases

- **No claimed evidence commands** - `REFUTED` (or refuse `VERIFIED`).
- **Evidence re-run fails** - Capture as fresh evidence; `REFUTED` until fixed and re-judged.
- **Caveats** - Only non-defect `decisions.md` notes. Open issues, review findings (including warnings), acceptance failure ⇒ `REFUTED`.
- **Open issues / parked related** - Any `Status: open` ⇒ `REFUTED`. Mission-caused left open for ship ⇒ `REFUTED`.
- **Manual-only check** - Fresh manual observation note in evidence; do not invent output.
- **Judge vs lifecycle** - Prove gate before `ready` inside run; does not own discuss/build/ship.

## Verdict contract

Exactly these three strings (case and spacing as written):

| Verdict | Meaning | Ready |
|---------|---------|-------|
| `VERIFIED` | Claim re-observed and holds | Allowed (subject to other gates) |
| `VERIFIED WITH CAVEATS` | Holds with listed non-defect `decisions.md` notes | Allowed only with those caveats recorded |
| `REFUTED` | Claim fails re-observation or hunt | **Blocked** - do not set `ready` |

No aliases (`PASS`, `FAIL`, `APPROVED`, etc.).

## Rules

- **Must**: Treat completion / "done" / "ready" claims as claims to re-observe - never trust the report alone.
- **Must**: Re-run claimed evidence commands; record fresh observation in `evidence.jsonl`.
- **Must**: Diff actual change scope vs `plan.json` / spec acceptance before verdict.
- **Must**: Hunt for weakened tests, false completion, and unauthorized action (use those phrases in findings so they are searchable).
- **Must**: Emit verdict exactly as `VERIFIED` | `VERIFIED WITH CAVEATS` | `REFUTED`.
- **Must**: When `REFUTED`, block `ready` (enforced by reviewer / `/sc-run`).
- **Must**: Preserve discuss / run / ship - judge is the prove gate, not a replacement lifecycle.
- **Must**: ASCII hyphen-minus only; Cursor-native; no Claude-plugin dependency.
- **Must**: Capture hunt misses as well as hits in the judge summary (what was checked).

## Out of scope

This skill does NOT handle:

- Writing product code or tests - use sc-coder / sc-tester
- Clarifying requirements or visual drafts - use `/sc-discuss`
- AFK build orchestration - use `/sc-run`
- Merge, tag, or ship - use `/sc-ship` only
- Trap-eval suites or domain adapters (deferred missions)

## Output format

```
## Judge summary
Mission: <id>
Claims reviewed: <list>
Evidence re-run: <fresh evidence ids / labels>
Scope vs plan: <match | mismatches>
Hunt:
  - weakened tests: <none | findings>
  - false completion: <none | findings>
  - unauthorized action: <none | findings>
  - open issues / review findings: <none | N open / N findings>
Caveats: <none | list>
VERDICT: VERIFIED | VERIFIED WITH CAVEATS | REFUTED
Ready: allowed | blocked
```

## Checklist

Before emitting a verdict:

- [ ] Mission resolved with `spacecraft resolve`
- [ ] Completion claims listed as claims, not trusted reports
- [ ] Claimed evidence commands re-run; fresh observation recorded
- [ ] Scope diffed against `plan.json` / spec acceptance
- [ ] Hunt covered weakened tests, false completion, unauthorized action
- [ ] Verdict is exactly one of the three contract strings
- [ ] If `REFUTED`, ready blocked stated explicitly

## References

- `sc-verification` - fresh evidence capture mechanics
- `/sc-run` - AFK path that must not set `ready` on `REFUTED`
- `sc-reviewer` agent - release readiness; consume judge verdict (wiring)
- `plan.json` / `spec.md` - acceptance and scope authority
- `evidence.jsonl` - append-only observations; judge appends re-runs
