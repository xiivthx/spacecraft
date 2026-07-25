---
name: sc-run
description: "AFK roadmap runner: after /sc-discuss clear, auto-run incomplete missions to ready. Invoke as /sc-run <roadmap-id>."
disable-model-invocation: true
---

# sc-run

## Goal

After `/sc-discuss` clear, AFK incomplete roadmap missions to `ready` for human check, then `/sc-ship`.

## Output

Missions at `state=ready` (or stop on blocked / clarify / missing draft approval). Handoff: **Ready. Human check, then /sc-ship.** Never merge/push/tag.

## Good / Bad

- Good: discuss clear first; plan → build → combine(+UI) → issues drain → review/judge; 0 open issues; empty findings; real evidence
- Bad: shipping; discuss work mid-AFK; skip drain; ready with open issues or findings; ready without `sc-judge` or after `REFUTED`; invent phrase-echo RED for docs/prose

## Verify

`spacecraft validate --strict` per mission; `map next` until `All missions complete.` or blocked tip.

## Arguments

`$ARGUMENTS` = roadmap id (optional if `spacecraft map current` is set).

```
spacecraft map use <roadmap-id>
/sc-run <roadmap-id>
# or /sc-run  (uses map current)
```

## Lifecycle

Canonical: `.cursor/rules/200-workflow.mdc` - discuss (HIL) → run (AFK) → human check → ship.

## Pre-flight

1. Resolve roadmap: `$ARGUMENTS` → `map use`; else `map current` (fail if missing).
2. Incomplete mission with `clarify-status` `open` → stop; `/sc-discuss`.
3. Visual UI/FE without `UI draft approved: …` (and no non-visual skip) → stop; `/sc-discuss`.
4. Soft gaps → `decisions.md`. No design-brief / draft-HTML discovery here.

## AFK loop

```
spacecraft map next <roadmap-id>
```

Stop when: `All missions complete.`, tip `blocked`, hard clarify, or missing draft approval.

### Per incomplete mission

1. Parse id from `map next` (`M…:`); `spacecraft use <id>`.
2. Branch `feat/<id>/…` from main; `spacecraft bind-branch <id>` if needed.
3. Spec already discuss-ready; ensure `.space/trust/lessons.md` (sc-learn seed if missing) and skim before planning; Task(`sc-planner`) → `plan.json`; `set-state planned` then `in_progress`. Ensure mission `issues.md` / `solved.md` / `learned.md` (sc-learn).
4. **Build (per acceptance)** - for each pending task / `acceptance[]`:
   1. Triage via sc-tdd; record `skip: <reason>` when tautology (docs/prose/wording-only).
   2. **TDD:** RED Task(`sc-tester`) → checkpoint → GREEN Task(`sc-coder`/`sc-firmware`) → evidence → checkpoint.
   3. **Skip:** direct write - Task(`sc-writer`) for docs/prose/wording-only, Task(`sc-coder`) for other tautologies - → evidence with task `verify` → one checkpoint. No phrase-harness scripts.
   4. **Findings mid-build:**
      - Blocks **current** acceptance → fix now; record → `solved.md`.
      - Else → append `open` to `issues.md` (class + severity); continue plan. Drain after combine(+UI).
   5. Mark task `done` only when all its acceptances pass.
5. **Combine:** refactor; full suite; evidence. Append new failures to `issues.md`. Checkpoint.
6. **UI recheck (visual UI/FE):** visual (sc-ux-design Tier 3: `playwright-cli` primary) + functional suite; append failures to `issues.md`. No ready yet.
7. **Issues drain** - until 0 `Status: open` (policy: **sc-learn**):
   1. Queue open entries (critical → important → minor; mission-caused before unrelated).
   2. Act per sc-learn matrix (must-fix mission-caused; file-or-fix unrelated).
   3. Verify each action; fixed → `solved.md`; filed → `filed` + GitHub URL.
   4. New problems → append `open`; continue loop.
   5. When empty, re-run suite (+ UI if UI); reopen → drain again.
   6. Same issue fails fix-verify **3** times or hard blocker → stop to human.
8. **Review + sc-judge** (only after 0 open):
   1. Task(`sc-reviewer`); UI also Task(`sc-designer`).
   2. Run sc-judge; evidence label including `judge`.
   3. Findings or `REFUTED` → open entries in `issues.md` → back to step 7.
   4. Clean: `review.json` (`status: ready`, empty `findings`); `validate --strict`; `set-state ready`.
9. Handoff: **Ready. Human check, then /sc-ship.** Continue `map next`. Squash is `/sc-ship` only.

```mermaid
flowchart TD
  A["plan build"] --> B["findings: fix if blocks acceptance else append"]
  B --> C["combine + UI"]
  C --> D["issues drain"]
  D --> E{"open?"}
  E -->|yes| F["sc-learn: fix or file"]
  F --> G{"new finding?"}
  G -->|yes| H["append open"]
  H --> E
  G -->|no| E
  E -->|no| I["re-run suite"]
  I -->|fails| H
  I -->|clean| J["review + judge"]
  J -->|findings or REFUTED| H
  J -->|clean| K["ready → HIL → ship"]
```

## Checkpoint commits

Auto-commit after every RED, GREEN, skip+evidence, combine, and material drain fix. Conventional Commits; no mission id in subject/body. Never push. See sc-git.

## Rules

- Never `/sc-ship`, merge, push, tag; never discuss/draft work mid-AFK.
- One feature branch per mission id; Task-delegate product code/tests.
- No ready without sc-judge, with `REFUTED`, with open issues, or with review findings.
- Must drain after plan+combine(+UI) until 0 open (sc-learn policy).
- After 3 failed fix-verify on the same issue → human.

## References

- `/sc-discuss`, `/sc-ship`
- sc-learn - issues policy
- sc-judge - ready prove gate
- sc-ux-design - post-build visual QC
