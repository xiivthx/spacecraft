---
name: sc-run
description: "AFK runner: after /sc-discuss clear, run roadmap queue or a single resolved mission to ready. Invoke as /sc-run [roadmap-id]."
disable-model-invocation: true
---

# sc-run

## Goal

After `/sc-discuss` clear, AFK incomplete work to `ready` for human check, then `/sc-ship`. Supports multi-mission roadmap queue **or** mission-only (`Sizing: single|phases`).

## Output

Missions at `state=ready` (or stop on blocked / clarify / missing draft approval). Handoff: **Ready. Human check, then /sc-ship.** Never merge/push/tag.

## Good / Bad

- Good: discuss clear first; plan → build → combine(+UI) → fix findings → review/judge; empty review findings; `sc-judge` `VERIFIED`; real evidence; summary lists what was fixed
- Bad: shipping; discuss work mid-AFK; parking defects for later; ready with review findings; ready without `sc-judge` or without `VERIFIED`; caveat / soft-pass ready; invent phrase-echo RED for docs/prose; demanding UI draft on `*-data` / `*-functional` / `*-integrate` seams

## Verify

`spacecraft validate --strict` per mission; roadmap mode: `map next` until `All missions complete.` or blocked tip; mission-only: that one mission `ready`.

## Arguments

`$ARGUMENTS` = roadmap id (optional).

```
spacecraft map use <roadmap-id>
/sc-run <roadmap-id>
# or /sc-run  → map current if set; else mission-only on spacecraft resolve/current
```

## Lifecycle

Canonical: `.cursor/rules/200-workflow.mdc` - discuss (HIL) → run (AFK) → human check → ship.

## Pre-flight

1. Resolve mode:
   - If `$ARGUMENTS` set → `map use` that roadmap (**roadmap mode**).
   - Else if `spacecraft map current` succeeds → **roadmap mode** on that id.
   - Else resolve mission via `spacecraft resolve` / `current` → **mission-only mode** (one mission; no map required). Fail only if neither roadmap nor mission resolves.
2. Incomplete mission with `clarify-status` `open` → stop; `/sc-discuss`.
3. Draft gate (visual only):
   - Require `UI draft approved: …` when the mission is visual UI/FE (`*-ui` title, or `Sizing: single|phases` without a non-visual skip).
   - Do **not** require draft when `decisions.md` has `UI draft skipped: non-visual seam (data|functional|integrate)` or another recorded non-visual skip. `*-integrate` tips use `UI draft skipped: non-visual seam (integrate)` - no draft demanded.
   - Missing required draft → stop; `/sc-discuss`.
4. Soft gaps → `decisions.md`. No design-brief / draft-HTML discovery here.

## AFK loop

**Roadmap mode:**

```
spacecraft map next <roadmap-id>
```

Stop when: `All missions complete.`, tip `blocked`, hard clarify, or missing draft approval.

**Mission-only mode:** run the single resolved mission through **Per incomplete mission** once; then stop (no `map next`).

### Per incomplete mission

1. Parse id from `map next` (`M…:`); `spacecraft use <id>`.
2. Branch `feat/<id>/…` from main; `spacecraft bind-branch <id>` if needed.
3. Spec already discuss-ready; Task(`sc-planner`) → `plan.json`; `set-state planned` then `in_progress`.
4. **Build (per acceptance)** - for each pending task / `acceptance[]`:
   1. Triage via sc-tdd; record `skip: <reason>` when tautology (docs/prose/wording-only).
   2. **TDD:** RED Task(`sc-tester`) → checkpoint → GREEN Task(`sc-coder`/`sc-firmware`) → evidence → checkpoint.
   3. **Skip:** direct write - Task(`sc-writer`) for docs/prose/wording-only, Task(`sc-coder`) for other tautologies - → evidence with task `verify` → one checkpoint. No phrase-harness scripts.
   4. **Findings mid-build:** fix now (especially if it blocks current acceptance or the suite). When recording/reporting defects, use impact-first craft from `references/defect-finding.md` (especially critical/important). Note each fix for the run summary.
   5. Mark task `done` only when all its acceptances pass.
5. **Combine:** refactor; full suite; evidence. Fix any new failures. Checkpoint.
6. **UI recheck (visual UI/FE):** sc-ux-design Step 0 **draft-parity** (side-by-side vs approved draft: tokens, layout, component chrome, scenario states) + Tier 3 visual (`playwright-cli` primary) + functional suite. Layout-only match with different chrome, or missing draft `data-state` coverage in the product, → fix now. No ready yet.
7. **Fix pass** - until suite (+ UI if UI) is clean:
   1. Fix mission-caused / suite-breaking / touched-path defects. For critical/important findings, follow `references/defect-finding.md` (impact-first title, user impact, 2-3 retest ideas).
   2. Unrelated preexisting that is not suite-breaking and not on touched path → note in summary only.
   3. Same issue fails fix-verify **3** times or hard blocker → stop to human.
8. **Review + sc-judge** (only after suite clean):
   1. Task(`sc-reviewer`); UI also Task(`sc-designer`).
   2. Run sc-judge; evidence label including `judge`. Verdict is binary: `VERIFIED` | `REFUTED` (no caveats).
   3. Any findings (any severity) or `REFUTED` → fix remediation now → re-review + re-judge. Do not set ready.
   4. Clean only when judge is `VERIFIED`, `review.json` has `status: ready`, and `findings` is empty: `validate --strict`; `set-state ready`.
9. Handoff: **Ready. Human check, then /sc-ship.** Include a short **Fixes** list in the summary. Continue `map next`. Squash is `/sc-ship` only.

```mermaid
flowchart TD
  A["plan build"] --> B["findings: fix now"]
  B --> C["combine + UI"]
  C --> D["fix pass"]
  D --> E{"suite clean?"}
  E -->|no| F["fix"]
  F --> D
  E -->|yes| J["review + judge"]
  J -->|findings or REFUTED| F
  J -->|VERIFIED and empty findings| K["ready → HIL → ship"]
```

## Checkpoint commits

Auto-commit after every RED, GREEN, skip+evidence, combine, and material fix. Conventional Commits; no mission id in subject/body. Never push. See sc-git.

## Rules

- Never `/sc-ship`, merge, push, tag; never discuss/draft work mid-AFK.
- One feature branch per mission id; Task-delegate product code/tests.
- Ready only after `sc-judge` `VERIFIED` and empty review findings (any severity blocks). Never ready on `REFUTED` or caveat soft-pass.
- Must fix blockers after plan+combine(+UI); report fixes in the run summary.
- After 3 failed fix-verify on the same issue → human.

## References

- `/sc-discuss`, `/sc-ship`
- sc-judge - ready prove gate
- sc-ux-design - post-build visual QC + draft-parity
- sc-web-frontend - port look from approved draft
- `references/defect-finding.md` - actionable defect findings for review/summary
