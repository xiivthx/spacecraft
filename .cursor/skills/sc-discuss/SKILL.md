---
name: sc-discuss
description: "Pre-build HIL: clarify, brainstorm, decide, and approve visual draft before implement. Invoke as /sc-discuss. Mid-ask unblock via natural language (re-pitch / research / visualize) - details in sc-clarify."
disable-model-invocation: true
---

# sc-discuss

## Goal

Find what we want before implement: clear `spec.md`, decisions, answered questions, and (when visual) an approved draft HTML. Exit with `clarify-status clear` so a new session can `/sc-run`.

## Output

Mission ready to build: solid `spec.md`; `questions.md` / `decisions.md` updated; visual: `UI draft approved: <file>` (or non-visual skip); mission brief accepted or skip recorded; `spacecraft clarify-status clear`. Handoff by sizing (see Verify / Exit). Never plan AFK, RED-GREEN, product code, or ship.

## Good / Bad

**Good**
- Sharp Goal / Output / Good-Bad / Verify; `Sizing: …` recorded
- Soft-pass: `## Lens pass` / `Lens pass skipped:`, `## Testability pass` / `Testability pass skipped:`, `## Strategy pass` / `Strategy pass skipped:`, and when two requirement versions `## RCRCRC pass` / `RCRCRC pass skipped:` - or eligible `Discuss path: fast`
- Testability when run: structured Test Ideas + Implementation pitfalls; on-demand `test-data-design.md` / `test-oracles.md` (not clear gates)
- Frontier rounds via sc-clarify (≤3 independent; serial when dependent); Verify / architecture / scope stay on frontier until settled or explicitly deferred; true soft → `decisions.md`
- Visual: product context + reference extract (when refs) + bake-off (or skip) + responsive ladder + scenario-complete draft; designer port + Impeccable craft before human; dimension-locked polish; then `UI draft approved:`
- Mission brief via Spec Mirror + stake map + Goal / Will do / Impact / Extra (Wrong-if when required); Accept / Adjust / Reject before clear
- Mid-ask via natural language per sc-clarify - no new slash commands

**Bad**
- Implementing; writing `plan.json` AFK; shipping
- Skipping bake-off silently; multi-dimension polish in one pass; unreviewed or scenario-incomplete draft HTML
- Clearing while draft unapproved, mission brief undecided, or Testability `Not Testable` + Verify soft/missing
- Dumping full testability queue or >1 frontier round (≤3) in one turn; quizzing instead of presenting the brief; hollow briefs
- Cross-feature layer waterfalls or `*-ux` roadmap seams; canvas as mission brief or draft HTML / visual SoT
- Inventing mid-ask slash skills; throwaway HTML for mid-ask visualize

## Verify

Human confirms spec (and draft when visual), then Accepts (or skip) the mission brief. All of the following must hold before clear:

```
spacecraft clarify-status clear
# decisions.md contains "Sizing: single" OR "Sizing: phases" OR "Sizing: roadmap <id>"
# when Sizing: roadmap: every seam decisions.md contains "Roadmap contract: locked <file>" OR "Contract lock deferred: exploratory, skeleton-first" (or grandfather "Roadmap contract skipped: pre-M1 map")
# visual: decisions.md contains "UI draft approved: <draft-file>" OR "UI draft skipped:"
# visual (when draft approved): "Layout bake-off winner: …" OR "Layout bake-off skipped: …"
# visual draft includes surface-relevant scenario matrix (per shared-draft-directives) when approved
# decisions.md contains "Mission brief: accepted" OR "Mission brief: skipped - <reason>"
# soft-pass: ("## Lens pass" OR "Lens pass skipped:") AND ("## Testability pass" OR "Testability pass skipped:") AND ("## Strategy pass" OR "Strategy pass skipped:") AND (when two requirement versions: "## RCRCRC pass" OR "RCRCRC pass skipped:")
#   OR (eligible fast path) "Discuss path: fast" stands in for lens/testability/strategy/RCRCRC soft-pass
```

Also: no open blocking questions; Verify present; do not clear while Testability is `Not Testable` and Verify soft/missing.

### Exit clear checklist

Before clear, Confirm: soft-pass lines present - `Lens pass skipped:` / `Testability pass skipped:` / `Strategy pass skipped:` / `RCRCRC pass skipped:` (when two versions) **or** eligible `Discuss path: fast`. Handoff by sizing (`references/mission-sizing.md`):
- `Sizing: roadmap <id>` → **Spec clear. New session: `/sc-run <id>`.**
- `Sizing: single` or `phases` → **Spec clear. New session: `/sc-run`.** (mission-only)

On handoff, set or update optional `mission.json` `pickup` (`phase`, `next` one-liner, `updatedAt`) so `spacecraft status` / session-start shows Pickup. Not a clear or closeout gate.

## Arguments

```
/sc-discuss
/sc-discuss <mission-id|title>
```

## Lifecycle

Canonical: `.cursor/rules/200-workflow.mdc` - this skill is discuss HIL only. Next: handoff by sizing (`/sc-run <roadmap-id>` or mission-only `/sc-run`).

## Pre-flight

1. Resolve: `$ARGUMENTS` → `spacecraft use`; else `spacecraft resolve`. `spacecraft new` only if user wants mutating work and none exists.
2. Read `spec.md`, `questions.md`, `decisions.md`, drafts when present.
3. Do not start `/sc-run` build, product code, or ship.

## Discuss loop

```
resolve → inspect → sizing gate → lens-pass gate → testability soft gate → strategy soft gate → (RCRCRC when two versions) → classify gaps → talk / ask / decide → (visual: brief → bake-off → polish+scenarios → designer → dimension-locked fix → human HIL) → mission brief → clear → handoff
```

### Sizing gate

**Always** apply `references/mission-sizing.md` (default `Sizing: single` when work fits one mission) before deep clarify or draft:

1. Classify checklist concerns present: UX / UI / functional / database.
2. Rough jigsaw count for a vertical slice.
3. Choose `single` | `phases` | `roadmap` per the playbook (3 feature seams when splitting: `*-data` → `*-functional` → `*-ui`, plus an optional `*-integrate` tip after the last feature seam when Must-when holds). **Must** use `*-functional` on new maps; **Must not** add a `*-ux` seam. When roadmap Must-when fires: auto-split per that playbook (**Must not ask** one-vs-many).
4. Record in `decisions.md`: `Sizing: single | phases | roadmap <id>` (+ seams/rationale when roadmap; + `Sizing phases: N - …` when phases).
5. If roadmap: follow **Map creation (discuss only)** ordered steps in `references/mission-sizing.md` (`spacecraft new` stubs → `map new` unless human approved reuse → `map add` → contract lock / wireframe when required → stub `Sizing:` + contract disposition on every seam → discuss tip only). Never leave map create/resize to planning.

### Soft-gate load rule

For lens / testability / strategy / RCRCRC: when triggers clearly do **not** fire, record the skip line in `decisions.md` and continue - **do not** load the full reference file only to skip. Load the matching `references/*.md` only when a trigger fires (or the human asks for that pass).

### Discuss fast path

Eligible when all hold: `Sizing: single`, non-visual, Verify present, empty blocking frontier. Then stamp `Discuss path: fast` in `decisions.md` — that marker stands in for lens / testability / strategy / RCRCRC soft-pass clear lines (satisfies soft-pass without per-pass skip paragraphs). Legacy `## Lens pass` / `Lens pass skipped:` (and Testability / Strategy / RCRCRC equivalents) remain valid without the fast marker.

**Must not** stamp `Discuss path: fast` when sizing is roadmap or phases; when Verify is soft or missing; or when visual draft is required. Marker alone is not enough if ineligible.

Fast path still requires solid `spec.md`, `Sizing: single`, UI draft skipped or N/A, mission brief accepted or skipped, then clear — brief is not skipped by fast path alone.

### Lens-pass gate

After sizing, when triggers fire (architecture fork, policy preference, soft Verify, sizing explosion risk), apply `references/lens-pass.md`. Otherwise record `Lens pass skipped: <reason>` without loading the ref. Tier 0: Commander checklist. Tier 1: Task(`sc-adviser`). Tier 2: 2-3 readonly Tasks (default Skeptic, Economist, Practitioner) → one `## Lens pass`. Tier 3 open-domain: sc-storm (not sc-search gray areas).

### Testability soft gate

After lens-pass, when triggers fire (soft/missing Verify, new feature with behavioral uncertainty, human asks for requirement review, mission brief probe finds Verify skim risk), apply `references/requirement-testability.md`. Otherwise record `Testability pass skipped: <reason>` without loading the ref. Park candidates in `questions.md`; ask via sc-clarify frontier rounds (≤3 independent; serial when dependent). Do not clear while Testability is `Not Testable` and Verify is still soft/missing.

Fail-closed discuss oracles (before `clarify-status clear`):

- Must: Browser-to-API seam — named error-envelope fields + FE agreement with OpenAPI or shared schema before clarify-status clear
- Must: Missing FE-API lock → clarify-status open
- Stamp greppable `FE-API contract: locked` when seam locked
- Must: Currency, locale, or AI seed in tip → Domain defaults: with concrete expected value before clear
- Stamp greppable `Domain defaults:` with concrete `key=value` pairs
- Must: Persona pack: required when persona explicitly enabled
- Stamp greppable `Persona pack: required` when persona enabled

### Strategy soft gate

After testability, when triggers fire (greenfield, multi-platform matrix, security/PII/compliance, critical integrations/SLOs, human asks for test strategy), apply `references/htsm-strategy.md`. Otherwise record `Strategy pass skipped: <reason>` without loading the ref. Strategy incompleteness does not block clear the way `Not Testable` + soft Verify does. Do not invent Verify from charters.

### RCRCRC (when two versions)

When existing and updated requirements are both available, apply `references/rcrcrc-impact.md` and record `## RCRCRC pass`. If only one version (or no delta): record `RCRCRC pass skipped: …` without loading the ref.

### Spec and decisions

1. Ensure `spec.md` has Goal, Output, Good vs Bad, Verify (machine-checkable where possible).
2. Blocking ambiguity: sc-clarify frontier rounds (≤3 independent; serial when dependent); record in `questions.md` / `decisions.md`.
3. Soft gaps: Verify / architecture / in-out scope → open frontier until settled or explicitly deferred (do not silently assume). True soft gaps → assumptions in `decisions.md` (do not block clear alone).
4. Lens / testability / strategy / RCRCRC: record `## … pass` (load ref when trigger fires) **or** the matching skip line without loading the ref when triggers clearly do not fire.
5. Deep architecture: Task(`sc-adviser`) with Tier 1 lens pass; high-stakes may use Tier 2; open-domain → sc-storm (Tier 3). Keep `clarify-status open` while blockers or unapproved visual draft remain.

### Mid-ask unblock

When the human is stuck mid-frontier (confused wording, needs facts, or cannot picture state), unlock via **natural language** - **no new slash** commands. Owner of Re-pitch / mid-ask escape details: **sc-clarify**. Cue → route: Re-pitch → sc-clarify Re-pitch on confusion; Research → `sc-search` then `sc-storm` for open-domain/strategy; Visualize → existing bake-off / draft (UI) or chat state/example table (non-UI). Do not add `/wait-what`, `/prototype`, or `/research` skills.

### Visual design (when UI/FE)

Detect from intent / `spec.md`. If visual:

- Follow `.cursor/skills/sc-ux-design/references/impeccable-orchestration.md`
- Task(`sc-designer`) owns port gates
- Default `Impeccable path: active` unless `Impeccable path: skipped: <reason>`
- Path skipped → legacy 6-dimension brief in sc-ux-design

Record Spacecraft gates in `decisions.md` (details in sc-ux-design): `Product context:` / `Product context skipped: greenfield`; `UX checklist:`; `Reference extract:` / `Reference borrow:` when refs; `DESIGN conflict:` when style conflicts; `Context fidelity:`; `Layout bake-off winner:` or `Layout bake-off skipped:`; surface-relevant scenario matrix; Responsive ladder (all four presets); `Impeccable craft: pass | waived:`; then `UI draft approved: <draft-file>` (or `UI draft skipped:` for non-visual FE / `*-data` / `*-functional` / `*-integrate`). Bake-off not required when draft is skipped. Approved draft is visual SoT for `/sc-run`.

### Mission brief (before clear)

Follow `references/mission-brief.md`. Spec Mirror → stake coverage map → Goal / Will do / Impact / Extra bullets (plain + technical; pre-mortem Wrong-if under Extra when non-trivial); human **Accept | Adjust | Reject**. Optional teach-back after Accept. Empty/soft Spec Mirror slots that block AFK → return to sc-clarify, do not invent brief completeness.

- Accept → record `Mission brief: accepted` (then clear if other gates hold)
- Adjust → record `Mission brief: adjust - <summary>`; update spec/decisions; re-brief; do not clear
- Reject → record `Mission brief: rejected - <reason>`; do not clear
- Skip → `Mission brief: skipped - <reason>`

Never clear while a posed brief awaits a decision (unless skip recorded).

**Must not:** replace the mission brief (Accept/Adjust/Reject chat HIL) with a Cursor Canvas; do not use canvas as draft HTML / visual SoT (draft stays HTML under discuss; canvases are `/sc-run` plan/findings/evidence milestones only).

## Rules

- Never `/sc-run` build, `/sc-ship`, merge, push, tag, or product implementation/tests (draft HTML only).
- Never present approval draft before designer **port** gate + Impeccable **craft** gate (or human craft waive) + critical/important fixes.
- Never record `UI draft approved` when required scenario states are missing, when bake-off winner/skip is missing on a visual draft path, when the responsive ladder fails (adjacent presets pixel-squeezed, any preset unusable/overflowing, widescreen stretched desktop with no measure control) without documented single-column exception, when `Reference borrow:` is set but extract artifact is missing, when `Impeccable path: active` and shape brief is unapproved, or when path active and craft is neither `pass` nor human `waived`.
- Never skip layout bake-off silently when layout is still open.
- Never restyle multiple visual dimensions in one human draft round (dimension lock).
- Never clear while mission brief undecided (unless skip recorded).
- Never replace mission brief (Accept/Adjust/Reject chat HIL) with a canvas; never use canvas as draft HTML / visual SoT.
- Never clear while Testability is `Not Testable` and Verify is still soft/missing.
- Never dump the full testability question queue in one user-facing turn - park in `questions.md` and ask via sc-clarify frontier rounds (≤3 independent; serial when dependent).
- Never create `*-ux` roadmap seams or cross-feature layer waterfalls (see `references/mission-sizing.md`).
- Never clear `Sizing: roadmap` without `Roadmap contract: locked <file>` or sanctioned `Contract lock deferred: exploratory, skeleton-first` (or grandfather skip) on every seam.
- Soft gates: when triggers clearly do not fire, record skip lines without loading full reference files.
- Prefer `spec.md` / `decisions.md` / `questions.md` over chat-only memory.
- Never invent mid-ask slash skills (`wait-what` / `prototype` / `research`); route mid-ask to sc-clarify (re-pitch / research via sc-search then sc-storm / visualize via bake-off or chat state table).
- One mission focus per discuss session.

## References

- sc-clarify (mid-ask Re-pitch / escape owner); sc-search → sc-storm; sc-ux-design; Task(`sc-designer`); Task(`sc-adviser`); optional Task(`sc-writer`); `/sc-run`; `/sc-ship`
- `references/mission-sizing.md`, `references/mission-brief.md`, `references/lens-pass.md`, `references/requirement-testability.md`, `references/htsm-strategy.md`, `references/rcrcrc-impact.md`
- On-demand (not soft gates): `references/sfdipot-coverage.md`, `references/test-data-design.md`, `references/test-oracles.md`
- `.cursor/skills/sc-ux-design/references/impeccable-orchestration.md`, `sc-ux-design/references/reference-extract.md`, `.cursor/skills/sc-ux-design/references/checklists/README.md`, `sc-ux-design/references/surface-checklist.md`
