# Mission sizing (vertical seams)

Canonical playbook for `/sc-discuss`: size work by **feature load**, not by global layers. Prefer one vertical mission; when large, split into **feature + seam** missions. Never invent a cross-feature waterfall (`all-db` → `all-api` → `all-ui`).

## Concern checklist (5 - coverage only)

Ask which concerns the requirement needs (omit absent ones):

| Concern | Meaning |
|---------|---------|
| **UX** | Design brief + scenario-complete draft (discuss-owned; visual SoT). Lives in the `*-ui` mission's `/sc-discuss` (or the single vertical mission) - **not** a roadmap seam. |
| **UI** | Port implementation from the approved draft |
| **Functional** | Behavior / API / domain logic + Verify |
| **Database** | Schema / migrations / queries |
| **Quality** | Testbench / security / performance (coverage-only). Fold measurable SEC/PERF bars into the owning seam's `Verify` when in scope. Dedicated `*-security` / `*-perf` tips only when quality **is** the deliverable - not routine coverage. |

## Mission seams when splitting (3 only)

| Seam | Mission title pattern | Owns |
|------|----------------------|------|
| **data** | `<feature>-data` | schema, migrations, queries |
| **functional** | `<feature>-functional` | API / domain / Verify |
| **ui** | `<feature>-ui` | UX draft HIL in discuss + port implementation |

**Must** use `*-functional` for new maps. `*-api` is allowed **only** when that title already exists in the repo/team convention - do not invent `*-api` on a fresh split.

**Must not** create `*-ux` as a roadmap seam.

## Optional integrate tip (fourth, non-visual)

`*-integrate` is not a feature seam and not `*-ux` - it is an optional fourth roadmap tip, added only after the last feature seam, to reconcile a multi-seam map once every feature seam has shipped.

| Tip | Mission title pattern | Owns |
|-----|----------------------|------|
| **integrate** (optional) | `<feature>-integrate` | **contract-conformance** + **drift drain** (cross-seam E2E Verify against locked `roadmap-contract.md`; dead code/path removal; route/contract alignment; cross-mission issue drain) |

**Must-when** (add `*-integrate` only if one holds):
- The map has more than one feature seam and a later seam's contract could invalidate an earlier seam's assumptions
- `*-ui` follows `*-functional` and the map needs a final cross-seam pass after `*-ui` ships
- Discuss records known drift, or a signal that integrate is needed, at map create

**Skip:** when the last tip's own combine already covers end-to-end behavior and there is no drift signal, do **not** add `*-integrate` - record `Integrate tip skipped: <reason>` in that mission's (or map-tip's) `decisions.md`. When integrate is skipped, the **last seam's combine** checks **contract-conformance** against the locked roadmap contract.

**Owns:** contract-conformance + drift drain; cross-seam E2E Verify vs locked contract; dead code/path removal; route/contract alignment; cross-mission issue drain.

**Must not:**
- Add new features
- Reopen UX art direction (visual SoT stays with the `*-ui` draft)
- Silently invent or resize the map mid-`/sc-run` - integrate is a discuss sizing decision like any other seam
- Rewrite a shipped tip's spec without a recorded delta in `decisions.md`

Per-mission combine (each tip's own `/sc-run` combine step) still runs as normal; `*-integrate` does not replace it - integrate is the cross-seam pass that runs after every feature seam has already combined on its own.

**Order default:** `*-data` → `*-functional` → `*-ui` → optional `*-integrate` when present. Integrate never blocks UI work - it always comes after the last feature seam. Omit a seam (or the integrate tip) when that concern is absent.

## Task granularity (planning pointer)

This file owns mission sizing (single / phases / roadmap). Task-level split heuristics - units, when to split a task, task shape, no wall-clock time - live in `sc-planning` **Split formula** (workflow step 3). Planning obeys `Sizing:` recorded here; do not duplicate this decision tree in planning prompts.

## Decision tree

**Always size** on every `/sc-discuss` (not only large asks). Default `Sizing: single` when the work fits one mission.

1. **Single mission (default)** - One primary user capability; estimated ≤7 jigsaw tasks after discuss; fold needed checklist concerns into that mission (discuss → plan → run). Draft still required when visual.
2. **Same-mission phases** - Still one capability / one ship story, but >7 tasks and seams are **not** independently shippable → record `Sizing: phases` (+ phase count / rationale). **Discuss owns the phases decision;** `/sc-run` planning may write `plan-phaseN.json` (never a multi-mission map).
3. **Multi-mission roadmap (Must when)** any of:
   - Estimated >7 tasks **and** at least two of the **3 seams** can ship/verify independently, **or**
   - Any single seam alone would need **≥4** jigsaw tasks (hard Must, not "roughly"), **or**
   - Visual draft HIL would block unrelated data/functional work if kept in one mission

When Must-when holds: **auto-split** (auto-apply roadmap sizing). Heuristic: **can-split** ⇒ work is **too big for one** mission. **Must not ask** one-vs-many / one mission vs many — apply the split; do not quiz single vs multi when Must-when is clear. Ask only when independent-shipability stays genuinely ambiguous after this playbook.

Then create the map under **Map creation (discuss only)** below. Discuss the current tip only; later seams wait for post-ship handoff (`Next: /sc-discuss <id>`).

## Map creation (discuss only)

`spacecraft map` is queue plumbing for `/sc-run` - not a planning ceremony. Only `/sc-discuss` (this playbook) may create or resize a multi-mission map.

**Ordered steps (Must):**

1. Record intended `Sizing: roadmap <id>` (+ seams / rationale) - will copy onto every seam in step 7.
2. For each needed seam title: `spacecraft new "<feature>-data|functional|ui"` (create stubs **before** `map add`).
3. `spacecraft map new "<feature-or-roadmap-title>"` **unless** the human explicitly approved reusing an existing roadmap id. Never silently append to unrelated `map current`.
4. `spacecraft map add <roadmap-id> <mission-id> --desc "<seam title>"` for each stub, in order `data` → `functional` → `ui` → optional `integrate` when Must-when holds (omit absent). The `integrate` stub still needs its own `spacecraft new "<feature>-integrate"` before `map add`, same as any other seam - discuss-owned only.
5. **Contract lock (Must on every multi-seam roadmap):** produce `roadmap-contract.md` under the roadmap or tip mission - cross-seam schema/API shapes + **interface-level scenario** skeleton **only** (not full visual draft, not per-mission hash oracles). Human-approve before any seam builds on the shapes. Record `Roadmap contract: locked <file>` on every seam's `decisions.md`. **Exempt:** `Sizing: single|phases` never pays contract cost. **Sanctioned skip only:** `Contract lock deferred: exploratory, skeleton-first` (record on every seam; still discuss-owned).
6. **Wireframe (when map has a `*-ui` seam):** at map-create, produce a quick lo-fi HTML wireframe (house: `serve-html.mjs`) naming surface regions + a **region→data-shape** mapping table bound to the locked contract. One fast HIL approval. Locks **structure + data flow**, **not** look (manga: name → draft → final). Visual draft approval **stays** at `*-ui` discuss and **Must** **conform** to the locked roadmap contract; record conformance in the `*-ui` mission's `decisions.md` (e.g. `Wireframe conformance: ok - <contract-file>`). Do **not** lock visual taste at map-create.
7. **Stub `decisions.md` on every seam** (tip and later) with the same:
   ```
   Sizing: roadmap <roadmap-id>
   Sizing seams: <feature>-data, <feature>-functional, <feature>-ui
   Sizing rationale: <one line>
   Roadmap contract: locked <path-to-roadmap-contract.md>
   ```
   (or the sanctioned `Contract lock deferred: exploratory, skeleton-first` line instead of the locked line). Later seams may add `Mission brief: skipped - stub seam; discuss after prior tip ships` until their discuss turn.
8. `spacecraft map use <roadmap-id>`; `spacecraft use <tip-id>`; discuss **current tip only** (full spec / draft / brief). Leave later seams for post-ship handoff.

**Must not** call `map new` / `map add` from `/sc-run` planning (`sc-planning` / `sc-planner`). If scope needs a new or resized multi-mission split mid-plan → stop and apply **Resize protocol** below via `/sc-discuss`.

## Roadmap contract freeze points

- **Map contract** (`roadmap-contract.md`): **interface-level scenario** skeleton + cross-seam shapes only. Not the per-mission behavior-oracle hash.
- **Per-mission** `Approved-scenarios:` remains the hash anchor for AFK freeze (see `docs/mission-artifacts.md`). Do not treat map-contract interface skeletons as frozen expected literals for product RED/GREEN.

## Roadmap contract re-lock protocol

When a locked cross-seam shape is wrong mid-map (never silent amend, never route-around):

1. **Stop** AFK / planning on seams that depend on the stale contract.
2. Hand to `/sc-discuss` with the roadmap / tip id.
3. Record the delta in `decisions.md` (what changed and why).
4. Stamp: `Roadmap contract: re-locked v<N> <file> supersedes v<N-1>: <reason>`
5. Resume only after the re-locked file is human-approved and every active seam records the new lock line.

## Gates version (grandfathering)

Record `Gates version: <governing-mission-id>` when a map or tip adopts the contract-lock / quality-gate grammar from a governing tip (e.g. after that tip ships).

- **New gates apply only after the governing mission ships** - do not retro-gate in-flight maps.
- In-flight maps created before the governing tip may record: `Roadmap contract skipped: pre-M1 map` (or equivalent pre-governing waiver) and continue without contract lock until their next discuss resize.
## Resize protocol (mid-plan / mid-run)

When planning or AFK discovers the mission must become multi-mission (or the map must change):

1. Stop AFK / planning. Do not `map new` from planning.
2. Hand to `/sc-discuss` with the in-progress mission id.
3. Discuss chooses one fate for the current mission (record in `decisions.md`):
   - **keep-as-seam** - retitle/repurpose current mission as one seam (usually the tip); create sibling stubs + map for the rest
   - **supersede** - archive or abandon current as wrong shape; create fresh seam stubs + new map
   - **re-tip** - current stays on map as a non-tip node; reorder / `map add` siblings; discuss the new tip
4. Re-record `Sizing: roadmap <id>` on every involved mission; clear only the tip after its discuss gates.
5. Discuss may also `map add` a `*-integrate` tip to an existing map when drift appears after map create (e.g. a later seam's contract diverged from an earlier one). Still discuss-only - never invent an integrate tip mid-`/sc-run` or from planning.

## Handoff by sizing

After tip `clarify-status clear`:

| `Sizing:` | Next |
|-----------|------|
| `roadmap <id>` | **Spec clear. New session: `/sc-run <id>`.** (`map use` that id) |
| `single` or `phases` | **Spec clear. New session: `/sc-run`.** (mission-only AFK on resolved current mission; no multi-mission map required) |

## Visual draft by seam

- Draft HIL required for `*-ui` and for `Sizing: single` / `phases` when the mission is visual UI/FE.
- On roadmap maps with a locked contract: `*-ui` discuss draft **Must** **conform** to the locked `roadmap-contract.md` (structure + data flow from the map-create wireframe). Record conformance in that `*-ui` `decisions.md`. Look/taste stays at `*-ui` discuss - never locked at map-create.
- `*-data`, `*-functional`, and `*-integrate` tips: **non-visual** - record `UI draft skipped: non-visual seam (<data|functional|integrate>)` in that mission's `decisions.md`. `/sc-run` must not demand a draft for those tips.
- Integrate tip example: `UI draft skipped: non-visual seam (integrate)`.

## Record

In each involved mission `decisions.md` (at map create for roadmap; at discuss for single/phases):

```
Sizing: single | phases | roadmap <roadmap-id>
Sizing seams: <feature>-data, <feature>-functional, <feature>-ui   # omit absent; roadmap only
Sizing rationale: <one line>
Roadmap contract: locked <path-to-roadmap-contract.md>   # roadmap only; or sanctioned skip / grandfather line
```

For `phases`, also record: `Sizing phases: <N> - <one-line rationale>`.

Roadmap description should state feature name, seam order, and why split.

Sanctioned roadmap contract lines (pick one greppable disposition when roadmap):

- `Roadmap contract: locked <file>`
- `Contract lock deferred: exploratory, skeleton-first`
- `Roadmap contract skipped: pre-M1 map` (grandfather / pre-governing waiver only)
- `Roadmap contract: re-locked v<N> <file> supersedes v<N-1>: <reason>` (after re-lock protocol)
## Must not

- Global layer roadmaps across unrelated features
- Fourth `*-ux` roadmap seam - the only sanctioned optional fourth tip is `*-integrate` reconcile, never `*-ux`
- Soft "prefer ≤7" or 8-9 exception bands (≤7 remains a hard Must per phase)
- Soft "roughly ≥4" for single-seam split (use hard ≥4)
- Inventing roadmap create/resize mid-`/sc-run` or mid-plan without returning to `/sc-discuss`
- Planning-owned "create a roadmap" as an escape hatch (multi-mission is a discuss sizing outcome)
- Silent reuse of `map current` for a different feature
- `map add` without prior `spacecraft new` stubs
- Leaving later seam `decisions.md` without `Sizing:` at map create
- Multi-seam roadmap without `Roadmap contract: locked <file>` or a sanctioned skip/grandfather line
- Silent amend of a locked `roadmap-contract.md` (use re-lock protocol)
- Locking full visual draft / look at map-create (wireframe = structure + data flow only)
- Inventing dedicated `*-security` / `*-perf` tips when Quality is only coverage, not the deliverable
- Retroactively applying new contract gates to in-flight maps without `Gates version:` / grandfather waiver