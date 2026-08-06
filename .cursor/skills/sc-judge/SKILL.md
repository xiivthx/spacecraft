---
name: sc-judge
description: "Adversarial prove gate before ready. Treat completion claims as claims; re-run evidence; diff scope vs plan; hunt weakened tests / false completion / unauthorized action. Verdict VERIFIED | REFUTED only."
---

# sc-judge

## Goal

Adversarial prove gate before `ready`: re-observe claimed completion. Never trust a report, summary, or "done" claim alone. Preserve `/sc-discuss` → `/sc-run` → `/sc-ship`; judge proves readiness - it does not replace the lifecycle.

## Output

A single verdict for the mission (or scoped claim under review):

```
VERDICT: VERIFIED | REFUTED
```

Plus: fresh evidence ids, scope-diff notes, hunt findings, and (when `REFUTED`) the refutation reasons plus a remediation list for `/sc-run` to fix. Cursor-native skill only - no Claude-plugin dependency.

## Good / Bad

- Good: treat every completion / "done" / "ready" claim as a claim; re-run claimed evidence commands and record fresh observation; diff actual change scope vs `plan.json` / spec acceptance; hunt weakened tests, false completion, unauthorized action; emit exactly one of the two verdicts; allow `ready` **only** on `VERIFIED`; on `REFUTED` emit a fix plan and block ready until re-judged
- Bad: believing a prior report or evidence.jsonl line without re-running; inventing evidence; soft-shipping past `REFUTED`; third verdicts or caveat soft-pass; replacing discuss/run/ship with a judge-only flow; expanding hunt into product redesign or trap-eval suites

## Verify

Before accepting a judge pass as proof:

```
spacecraft evidence "judge-<mission-id>" -- <re-run of claimed commands>
```

Confirm verdict string is exactly one of `VERIFIED` | `REFUTED`, hunt notes cover the hunt targets, scope-diff vs plan/spec is recorded, and `Ready: allowed` only when verdict is `VERIFIED`.

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
   - **false completion** - "done"/"ready" claimed while acceptance fails, evidence is missing/stale, scope does not match plan, defects left unfixed, or `review.json` still has findings (including minor / warnings)
   - **unauthorized action** - outward push/deploy/publish/send (or similar) without quoted `AUTH:` and user authorization; ship/merge without `/sc-ship` gates
   - **draft drift (visual UI)** - when `UI draft approved: …` is recorded, treat "matches draft" / visual ready as a claim: REFUTE if product chrome clearly diverges from the approved draft (layout-only match) or required draft scenario states (`empty` / `error` / `few` / `many` / spec features) were never implemented or tested
   - **SFDIPOT blind spots (optional aid)** - when `## Testability pass` Test Ideas or `## Strategy pass` Charter ideas exist, may use `sc-discuss/references/sfdipot-coverage.md` to hunt coverage gaps. Blind-spot suggestions alone do **not** justify `REFUTED`. `REFUTED` only when claimed acceptance or Test Idea coverage maps to a **Missing** SFDIPOT area asserted done without fresh evidence (false completion)
   - **Test data gaps (optional aid)** - missing `## Test data design` rows alone do **not** justify `REFUTED`
6. **Verdict** - Emit exactly one of:
   - `VERIFIED` - perfect bar: fresh evidence passes; scope matches plan/spec; hunt clean; **0** review findings (critical / important / minor - no warning band). Ready allowed (subject to other gates).
   - `REFUTED` - any gap: material hunt hit, failed re-run, scope/acceptance mismatch, leftover review findings (any severity), or failed verify. Ready blocked.
7. **Ready gate** - Allow `ready` **only** on `VERIFIED`. On `REFUTED`, block `ready` and emit a remediation list so `/sc-run` can fix → re-review → re-judge. Do not soften to ship or ready. No caveat / soft-pass verdict.

### Edge cases

- **No claimed evidence commands** - `REFUTED` (or refuse `VERIFIED`).
- **Evidence re-run fails** - Capture as fresh evidence; `REFUTED` until fixed and re-judged.
- **Non-defect `decisions.md` notes** - Allowed alongside `VERIFIED` as recorded decisions only. They do **not** create a third verdict. Unfinished follow-up work ⇒ `REFUTED`.
- **Any review finding** - Critical, important, or minor (including warnings) ⇒ `REFUTED` until findings are empty.
- **Manual-only check** - Fresh manual observation note in evidence; do not invent output.
- **Judge vs lifecycle** - Prove gate before `ready` inside run; does not own discuss/build/ship.

## Verdict contract

Exactly these two strings (case and spacing as written):

| Verdict | Meaning | Ready |
|---------|---------|-------|
| `VERIFIED` | Claim re-observed and holds at the perfect bar | **Allowed** (subject to other gates) |
| `REFUTED` | Claim fails re-observation or hunt | **Blocked** - fix → re-judge; do not set `ready` |

No aliases (`PASS`, `FAIL`, `APPROVED`, `VERIFIED WITH CAVEATS`, etc.).

## Rules

- **Must**: Treat completion / "done" / "ready" claims as claims to re-observe - never trust the report alone.
- **Must**: Re-run claimed evidence commands; record fresh observation in `evidence.jsonl`.
- **Must**: Diff actual change scope vs `plan.json` / spec acceptance before verdict.
- **Must**: Hunt for weakened tests, false completion, unauthorized action, and (when visual UI) draft drift (use those phrases in findings so they are searchable).
- **Must**: Emit verdict exactly as `VERIFIED` | `REFUTED`.
- **Must**: Allow `ready` only when verdict is `VERIFIED` (enforced by reviewer / `/sc-run`).
- **Must**: When `REFUTED`, block `ready`, list remediation for `/sc-run` to fix, and require re-judge after fixes.
- **Must**: Preserve discuss / run / ship - judge is the prove gate, not a replacement lifecycle.
- **Must**: ASCII hyphen-minus only; Cursor-native; no Claude-plugin dependency.
- **Must**: Capture hunt misses as well as hits in the judge summary (what was checked).

## Judge-break fixtures

Known-bad packs under `references/judge-break/` prove the ready/ship exit gate **rejects** bad disk state (empty evidence, review findings, false completion). Deterministic only - no LLM.

```
make test-judge-break
# or
scripts/check-judge-break.sh [repo-root] [spacecraft-binary]
```

Run before claiming sc-judge skill or closeout predicate changes are safe. Adding a new reject path: add a fixture pack with `expect.json` (`id`, `mustContain`) plus minimal mission files.

## Out of scope

This skill does NOT handle:

- Writing product code or tests - use sc-coder / sc-tester
- Clarifying requirements or visual drafts - use `/sc-discuss`
- AFK build orchestration - use `/sc-run`
- Merge, tag, or ship - use `/sc-ship` only
- Trap-eval suites, trajectory rubrics, or LLM-as-judge scoring of agent transcripts (deferred)

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
  - draft drift (visual UI): <n/a | none | findings>
  - review findings: <none | N findings>
Remediation (when REFUTED): <none | list for /sc-run to fix>
VERDICT: VERIFIED | REFUTED
Ready: allowed | blocked
```

## Checklist

Before emitting a verdict:

- [ ] Mission resolved with `spacecraft resolve`
- [ ] Completion claims listed as claims, not trusted reports
- [ ] Claimed evidence commands re-run; fresh observation recorded
- [ ] Scope diffed against `plan.json` / spec acceptance
- [ ] Hunt covered weakened tests, false completion, unauthorized action, and draft drift when visual UI
- [ ] Verdict is exactly one of the two contract strings
- [ ] `Ready: allowed` only when verdict is `VERIFIED`
- [ ] If `REFUTED`, ready blocked and remediation listed explicitly

## References

- `sc-verification` - fresh evidence capture mechanics
- `/sc-run` - AFK path that must set `ready` only on `VERIFIED`
- `sc-reviewer` agent - release readiness; consume judge verdict (wiring)
- `plan.json` / `spec.md` - acceptance and scope authority (behavior)
- approved draft HTML - visual look authority when `UI draft approved` is recorded
- `evidence.jsonl` - append-only observations; judge appends re-runs
- `references/judge-break/` - known-bad fixtures; `scripts/check-judge-break.sh`
