---
name: sc-discuss
description: "Pre-build HIL: clarify, brainstorm, decide, and approve visual draft before implement. Invoke as /sc-discuss. Mid-ask unblock via natural language (re-pitch / research / visualize) - details in sc-clarify."
disable-model-invocation: true
---

# sc-discuss

## Goal

Find what we want before implement: clear `spec.md`, decisions, answered questions, and (when visual) an approved draft HTML. Exit with `clarify-status clear` so a new session can `/sc-run`.

## Output

Mission ready to build: solid `spec.md`; `questions.md` / `decisions.md` updated; visual: `UI draft approved: <file>` (or non-visual skip); mission brief accepted or skip recorded; `spacecraft clarify-status clear`. Handoff by sizing (see Exit). Never plan AFK, RED-GREEN, product code, or ship.

## Good / Bad

- Good: sharp Goal / Output / Good-Bad / Verify; sizing recorded (`Sizing: …`); lens pass or skip recorded before clear; testability pass or skip recorded before clear (structured Test Ideas Positive/Negative/Edge/Overlooked + Implementation pitfalls when testability runs); on-demand test-data design via `test-data-design.md` when variable-level fixtures matter (not a clear gate); on-demand oracle evaluation via `test-oracles.md` when problem judgment needs grounding (not a clear gate); strategy pass or skip recorded before clear; RCRCRC when two requirement versions exist; frontier rounds via sc-clarify (≤3 independent blocking questions per turn; serial when dependent); Verify / architecture / scope soft gaps stay on the open frontier until settled or explicitly deferred; true soft gaps → `decisions.md`; visual brief + product context + reference extract (when refs) + context fidelity + layout bake-off (or skip) + responsive ladder (all four presets) + scenario-complete draft with designer gate before human; dimension-locked polish; mission brief via Spec Mirror + stake map + Goal / Will do / Impact / Extra bullets (plain + technical; pre-mortem Wrong-if under Extra when non-trivial) then Accept/Adjust/Reject before clear; mid-ask unblock via natural language (re-pitch / research / visualize) per sc-clarify - no new slash commands
- Bad: implementing; writing `plan.json` AFK; shipping; skipping bake-off silently; polishing type+color+layout in one pass; serving unreviewed or scenario-incomplete draft HTML; clearing while draft unapproved or mission brief undecided; clearing while Testability is `Not Testable` and Verify soft/missing; dumping the full testability queue or more than a frontier round (≤3 independent) in one turn; quizzing the human instead of presenting the brief; hollow briefs (vague Goal/Will do/Impact, no Wrong-if when required, Spec Mirror soft/empty yet posed); cross-feature layer waterfalls or `*-ux` roadmap seams; replacing mission brief (Accept/Adjust/Reject chat HIL) with a canvas; using canvas as draft HTML / visual SoT; inventing mid-ask slash skills; throwaway HTML for mid-ask visualize

## Verify

Human confirms spec (and draft when visual), then Accepts (or skip) the mission brief:

```
spacecraft clarify-status clear
# decisions.md contains "Sizing: single" OR "Sizing: phases" OR "Sizing: roadmap <id>"
# visual: decisions.md contains "UI draft approved: <draft-file>" OR "UI draft skipped:"
# visual (when draft approved): "Layout bake-off winner: …" OR "Layout bake-off skipped: …"
# visual draft includes surface-relevant scenario matrix (per shared-draft-directives) when approved
# decisions.md contains "Mission brief: accepted" OR "Mission brief: skipped - <reason>"
# soft-pass: ("## Lens pass" OR "Lens pass skipped:") AND ("## Testability pass" OR "Testability pass skipped:") AND ("## Strategy pass" OR "Strategy pass skipped:") AND (when two requirement versions: "## RCRCRC pass" OR "RCRCRC pass skipped:")
#   OR (eligible fast path) "Discuss path: fast" stands in for lens/testability/strategy/RCRCRC soft-pass
```

Handoff by sizing: roadmap → `/sc-run <id>`; single|phases → `/sc-run` (mission-only).
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
5. If roadmap: follow **Map creation (discuss only)** ordered steps in `references/mission-sizing.md` (`spacecraft new` stubs → `map new` unless human approved reuse → `map add` → stub `Sizing:` on every seam → discuss tip only). Never leave map create/resize to planning.

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

When the human is stuck mid-frontier (confused wording, needs facts, or cannot picture state), unlock via **natural language** - **no new slash** commands. Owner of Re-pitch / mid-ask escape details: **sc-clarify**.

| Cue | Route |
|-----|-------|
| Re-pitch | sc-clarify **Re-pitch on confusion** |
| Research | `sc-search` first → escalate `sc-storm` for open-domain/strategy |
| Visualize | Existing layout **bake-off** / draft (UI) or chat state/example table (non-UI) |

Do not add `/wait-what`, `/prototype`, or `/research` skills. After research or visualize, return to the open frontier round.

### Visual design (when UI/FE)

Detect from intent / `spec.md`. If visual:

1. **Product context:** Record `Product context: <routes + shell/layout file paths + screenshot paths>` or `Product context skipped: greenfield` in `decisions.md`. When brownfield, read parent shell/layout and nearby pages before brief; bake-off candidates must include existing app chrome for in-app screens.
2. **Surface checklist:** Match one primary id via `.cursor/skills/sc-ux-design/references/checklists/README.md` and `sc-ux-design/references/surface-checklist.md`. Record `UX checklist: <id>` or `UX checklist: none - <reason>` in `decisions.md` before bake-off. Read that one file under `.cursor/skills/sc-ux-design/references/checklists/` and fold applicable `- [ ]` items into `spec.md` Must and the approval draft.
3. Read project `DESIGN.md` when present (house look SoT). **Reference extract when images/refs supplied:** run `sc-ux-design/references/reference-extract.md` before brief → `design/refs/extract.md`; record `Reference extract: design/refs/extract.md` and human-confirmed `Reference borrow:` (`mood` | `tokens` | `layout` | `chrome`). Brief must cite extract rows. **Must not** enter bake-off when `Reference borrow:` is set but extract is missing. If proposed style conflicts with `DESIGN.md`, ask once and record `DESIGN conflict: mission exception | update house | keep house`. Then sc-ux-design design brief (6 dimensions, aligned to effective house); human approval. No art-direction pack question - packs removed.
4. **Context fidelity before bake-off:** Record `Context fidelity: DESIGN.md | shell:<path> | extract:<path> | product-shot:<path>` (omit absent; greenfield may omit shell/product-shot).
5. **Layout bake-off:** After brief approval, generate **2–3** HTML layout candidates under `.space/missions/<id>/design/drafts/` (`<name>-draft-v1-<layout-id>.html`) with scaffold + brief tokens + primary surface chrome + **Responsive ladder** across all four presets (375 / 768 / 1280 / 1536) for multi-region UIs. Serve and let the human pick. Record `Layout bake-off winner: <file>` or `Layout bake-off skipped: <reason>` (forced by `DESIGN.md`, borrow `layout`|`chrome`, or user-named structure). Do not skip silently. Full scenario matrix may wait until the winner.
6. **Polish winner** (sc-ux-design): full draft with scaffold (chrome outside framed surface); viewport toggles (375 / 768 / 1280 / 1536) with responsive CSS (not frame-resize-only at any preset); **surface-relevant scenario matrix** per `spec.md` + primary surface shape (happy path + failure/degraded the surface can enter; `loading` when async; `empty`/`few`/`many` when variable-length collection; plus feature/behavior surfaces from `spec.md`); applicable `- [ ]` items from the recorded id's file under `.cursor/skills/sc-ux-design/references/checklists/` as real chrome; real component chrome. Check all four viewports and **Responsive ladder** at each preset. Prompt assembly: `shared directives` → `DESIGN.md` → `design-principles.md` → brief. Honor borrow scope - no silent full-chrome clone.
7. **Dimension lock while iterating:** change only one of `typography` | `color` | `layout` | `motion` | `spacing` | `chrome` per human round. Prefer "list diffs vs draft/reference; focus: <dimension>" over "make it look better."
8. **Designer gate before approval HIL:** Task(`sc-designer`) on the approval candidate; check scenario coverage + surface checklist + port readiness + scaffold/frame + product continuity + reference extract + **Responsive ladder**; Commander applies critical/important fixes; re-check all four viewports and size-appropriate organization at each preset. Do not present the approval draft until this passes. Missing required states, missing production frame, frame-resize-only / squeeze-only responsive at any preset, or product-continuity gaps = critical - do not serve.
9. Serve via `serve-html.mjs`; iterate (draft → designer → fix → human) with dimension lock until approved (max 3 human rounds after bake-off). Each new draft re-runs designer gate.
10. On approval: record `UI draft approved: <draft-file>` in `decisions.md` **only if** the scenario matrix is complete, `UX checklist:` (id or none) is recorded, and bake-off winner/skip is recorded. Incomplete states → refuse approval; iterate draft.
11. Skip draft for non-visual FE, or for `*-data` / `*-functional` / `*-integrate` seams: record `UI draft skipped: non-visual seam (<data|functional|integrate>)` or other skip reason (e.g. `UI draft skipped: non-visual seam (integrate)`). Bake-off not required when draft is skipped.
12. Tell the human: approved draft is the **visual source of truth** for `/sc-run` (port look; do not freestyle chrome; no layout bake-off in run). Prefer draft surface chrome that maps cleanly to reusable `components/ui` primitives (button, field, banner, empty) so `/sc-run` can build component-first.

### Mission brief (before clear)

Follow `references/mission-brief.md`. Spec Mirror → stake coverage map → Goal / Will do / Impact / Extra bullets (plain + technical; pre-mortem Wrong-if under Extra when non-trivial); human **Accept | Adjust | Reject**. Optional teach-back after Accept. Empty/soft Spec Mirror slots that block AFK → return to sc-clarify, do not invent brief completeness.

- Accept → record `Mission brief: accepted` (then clear if other gates hold)
- Adjust → record `Mission brief: adjust - <summary>`; update spec/decisions; re-brief; do not clear
- Reject → record `Mission brief: rejected - <reason>`; do not clear
- Skip → `Mission brief: skipped - <reason>`

Never clear while a posed brief awaits a decision (unless skip recorded).

**Must not:** replace the mission brief (Accept/Adjust/Reject chat HIL) with a Cursor Canvas; do not use canvas as draft HTML / visual SoT (draft stays HTML under discuss; canvases are `/sc-run` plan/findings/evidence milestones only).

### Exit

1. No open blocking questions; Verify present; `Sizing: …` recorded; (`## Lens pass` or `Lens pass skipped:`) and (`## Testability pass` or `Testability pass skipped:`) and (`## Strategy pass` or `Strategy pass skipped:`) and (when two requirement versions, `## RCRCRC pass` or `RCRCRC pass skipped:`) — or, when eligible, `Discuss path: fast` stands in for those soft-pass lines; visual approved or skip recorded; mission brief accepted or skip recorded. Do not clear while Testability is `Not Testable` and Verify soft/missing.
2. `spacecraft clarify-status clear`.
3. Handoff by sizing (`references/mission-sizing.md`):
   - `Sizing: roadmap <id>` → **Spec clear. New session: `/sc-run <id>`.**
   - `Sizing: single` or `phases` → **Spec clear. New session: `/sc-run`.** (mission-only AFK on resolved current mission)
4. On handoff, set or update optional `mission.json` `pickup` (`phase`, `next` one-liner, `updatedAt`) so `spacecraft status` / session-start shows Pickup. Not a clear or closeout gate.

## Rules

- Never `/sc-run` build, `/sc-ship`, merge, push, tag, or product implementation/tests (draft HTML only).
- Never present approval draft before designer gate + critical/important fixes.
- Never record `UI draft approved` when required scenario states are missing, when bake-off winner/skip is missing on a visual draft path, when the responsive ladder fails (adjacent presets pixel-squeezed, any preset unusable/overflowing, widescreen stretched desktop with no measure control) without documented single-column exception, or when `Reference borrow:` is set but extract artifact is missing.
- Never skip layout bake-off silently when layout is still open.
- Never restyle multiple visual dimensions in one human draft round (dimension lock).
- Never clear while mission brief undecided (unless skip recorded).
- Never replace mission brief (Accept/Adjust/Reject chat HIL) with a canvas; never use canvas as draft HTML / visual SoT.
- Never clear while Testability is `Not Testable` and Verify is still soft/missing.
- Never dump the full testability question queue in one user-facing turn - park in `questions.md` and ask via sc-clarify frontier rounds (≤3 independent; serial when dependent).
- Never create `*-ux` roadmap seams or cross-feature layer waterfalls (see `references/mission-sizing.md`).
- Soft gates: when triggers clearly do not fire, record skip lines without loading full reference files.
- Prefer `spec.md` / `decisions.md` / `questions.md` over chat-only memory.
- Never invent mid-ask slash skills (`wait-what` / `prototype` / `research`); route mid-ask to sc-clarify (re-pitch / research via sc-search then sc-storm / visualize via bake-off or chat state table).
- One mission focus per discuss session.

## Specialist skills

| Concern | Where |
|---------|--------|
| Blocking questions / mid-ask escapes | sc-clarify (Re-pitch / mid-ask escape owner); research: sc-search → sc-storm; visualize: bake-off or chat state table |
| Mission sizing / roadmap split | `references/mission-sizing.md` |
| Requirement testability | `references/requirement-testability.md` |
| SFDIPOT coverage review (existing tests vs requirement) | `references/sfdipot-coverage.md` (on-demand; not a soft gate) |
| Test data design (variable-level fixtures) | `references/test-data-design.md` (on-demand; not a soft gate) |
| Test oracles (FEW HICCUPPS) | `references/test-oracles.md` (on-demand; not a soft gate) |
| HTSM strategy (slim) | `references/htsm-strategy.md` |
| Requirement delta / RCRCRC | `references/rcrcrc-impact.md` |
| Draft HTML / visual-verify | sc-ux-design |
| Reference extract | `sc-ux-design/references/reference-extract.md` |
| Surface checklist | `.cursor/skills/sc-ux-design/references/checklists/README.md` (item SoT) + `sc-ux-design/references/surface-checklist.md` (one-id adapter) |
| Draft critique | Task(`sc-designer`) |
| Mission brief | `references/mission-brief.md` |
| Architecture | Task(`sc-adviser`) |
| Lens pass / STORM | `references/lens-pass.md`; optional sc-storm (Tier 3) |
| Prompt / docs / spec wording | optional Task(`sc-writer`) |
| Plan / TDD / evidence | `/sc-run` |

## References

- sc-clarify (mid-ask Re-pitch / escape owner), sc-ux-design, `/sc-run`, `/sc-ship`
- `references/mission-sizing.md`
- `references/mission-brief.md`
- `references/lens-pass.md`
- `references/requirement-testability.md`
- `references/sfdipot-coverage.md`
- `references/test-data-design.md`
- `references/test-oracles.md`
- `references/htsm-strategy.md`
- `references/rcrcrc-impact.md`
- `.cursor/skills/sc-ux-design/references/checklists/README.md`
- `sc-ux-design/references/surface-checklist.md`
- sc-search → sc-storm (mid-ask research escape)
