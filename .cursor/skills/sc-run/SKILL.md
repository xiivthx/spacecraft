---
name: sc-run
description: "AFK runner: after /sc-discuss clear, run roadmap queue or a single resolved mission to ready. Invoke as /sc-run [roadmap-id]."
disable-model-invocation: true
---

# sc-run

## Goal

After `/sc-discuss` clear, AFK incomplete work to `ready` for human check, then `/sc-ship`. Roadmap queue or mission-only. Never merge/push/tag. TDD theory lives in sc-tdd only.

## Output

Missions at `state=ready`, or stop on `3-cycle` / `timebox` / `blocked` / clarify / missing draft → write `.space/missions/<id>/handback.md`. Handoff: **Ready. Human check, then /sc-ship.**

## Good / Bad

- Good: discuss clear → ordered lifecycle; Task-delegate; disposition via `docs/mission-artifacts.md`; `sc-judge` `VERIFIED` + empty findings before ready
- Bad: shipping; product RED/GREEN without design-contract/approved-scenarios (or docs/prose skips); coder edits tests; ready without judge/`VERIFIED`

## Verify

`spacecraft validate --strict`; matching `evidence.jsonl`; empty `review.json` findings; `sc-judge` `VERIFIED`; UI/workflow in scope → `PROBE: CLEAN` (or probe skip). Canvas ≠ ready proof.

## Pre-flight

1. Resolve: `$ARGUMENTS` → roadmap; else `map current`; else mission-only via `resolve`/`current`. Fail only if neither resolves.
2. `clarify-status` open → stop; `/sc-discuss`.
3. Visual UI/FE: require `UI draft approved:` unless `UI draft skipped:`. Missing → stop; `/sc-discuss`.
4. Soft gaps → `decisions.md`.

## AFK lifecycle (per incomplete mission)

Ordered: discuss-clear → plan → design-contract → approved-scenarios → build → combine → fix → (browser-probe when UI) → review → judge.

1. `spacecraft use <id>`; branch `feat/<id>/…`.
2. **plan** - Task(`sc-planner`) → `plan.json`; `set-state planned` then `in_progress`.
3. **design-contract** - write or `Design-contract skipped: docs/prose-only` (`docs/mission-artifacts.md`). No product RED/GREEN until complete or skip.
4. **approved-scenarios** - freeze or `Approved-scenarios skipped: docs/prose-only`. Agents Must not thaw frozen literals; oracle changes need Commander + `Scenario oracle change:`.
5. **build** - triage via sc-tdd. TDD: RED Task(`sc-tester`) → GREEN Task(`sc-coder`/`sc-firmware`/`sc-rtl`) → evidence; coder **Must not** edit tests. Skip: Task(`sc-writer`) docs/prose, else Task(`sc-coder`) → evidence. Commander Task-delegates product code/tests (role split).
6. **combine** - ordinary refactor (keep here; Bugbot/security-review drive the fix loop only); full suite; `spacecraft freeze-check` (exit 1 → handback: drift/postdate/missing freeze); static/diff-cov/mutation disposition → `docs/mission-artifacts.md` only (do not paste bars here). Checkpoint.
7. **fix** - until suite (+ UI if UI) clean. Same issue **3** times → human. Stop reasons → `handback.md`.
8. **browser-probe** (UI/workflow) - Task(`sc-browser-probe`) to `PROBE: CLEAN`; else handback. Skip when no runnable UI/workflow.
9. **review** → **judge** - `validate --strict`; parallel Task(`bugbot`) + Task(`security-review`) before or as input to Task(`sc-reviewer`) (exact Cursor prompt: `Full Repository Path: <abs repo>`; `Diff: branch changes` default; optional `Base Branch` / `Custom Instructions` / bugbot-only `Change Description` per upstream skills - invoke + thin adapter only, do not copy skill bodies); write greppable `Cursor review: bugbot+security-review ran` or `Cursor review skipped: <reason>` (unavailable Cursor types ⇒ skip disposition). When writing `Cursor review: … ran`, Must also capture `spacecraft evidence "cursor-review-…" -- …` **or** write greppable `Cursor ingest: session` in `decisions.md`. When security is in scope, `Cursor review skipped:` is valid **only** after a greppable `Sc-security fallback: pass` (optionally also `Sc-security fallback: findings drained`; optional evidence label `sc-security-…`) **or** SEC machine-evidence pass; otherwise fail / ready blocked (skip alone does not satisfy security-when-in-scope). Map bugbot + security-review (+ sc-reviewer) outputs into `review.json` via `defect-finding.md` with `source` set; Task(`sc-reviewer`); sc-judge. Drain **all severities** (`critical`/`important`/`minor`): findings/`REFUTED` → fix (only when `requiredFix` is concrete; no free-form whole-tree refactor from Bugbot noise) → re-run Cursor reviews → re-review → re-judge until findings empty or `3-cycle`/`timebox` handback. Ready only on `VERIFIED` + empty findings: `set-state ready`. Handoff; never merge/push/tag.

Overnight: stop on `3-cycle` | `timebox` | `blocked` → `handback.md`. Checkpoint commits: one Conventional Commit per plan task after acceptances; see sc-git. Never push.

## References

- sc-tdd, sc-judge, sc-browser-probe, `/sc-discuss`, `/sc-ship`, sc-git
- `docs/mission-artifacts.md` - design-contract / approved-scenarios / outcome-gate SoT
- `references/defect-finding.md`, `references/mission-review-gates.md`
