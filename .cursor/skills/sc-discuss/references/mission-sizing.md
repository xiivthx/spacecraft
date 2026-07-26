# Mission sizing (vertical seams)

Canonical playbook for `/sc-discuss`: size work by **feature load**, not by global layers. Prefer one vertical mission; when large, split into **feature + seam** missions. Never invent a cross-feature waterfall (`all-db` → `all-api` → `all-ui`).

## Concern checklist (4 - coverage only)

Ask which concerns the requirement needs (omit absent ones):

| Concern | Meaning |
|---------|---------|
| **UX** | Design brief + scenario-complete draft (discuss-owned; visual SoT). Lives in the `*-ui` mission's `/sc-discuss` (or the single vertical mission) - **not** a roadmap seam. |
| **UI** | Port implementation from the approved draft |
| **Functional** | Behavior / API / domain logic + Verify |
| **Database** | Schema / migrations / queries |

## Mission seams when splitting (3 only)

| Seam | Mission title pattern | Owns |
|------|----------------------|------|
| **data** | `<feature>-data` | schema, migrations, queries |
| **functional** | `<feature>-functional` | API / domain / Verify |
| **ui** | `<feature>-ui` | UX draft HIL in discuss + port implementation |

**Must** use `*-functional` for new maps. `*-api` is allowed **only** when that title already exists in the repo/team convention - do not invent `*-api` on a fresh split.

**Must not** create `*-ux` as a roadmap seam.

**Order default:** `*-data` → `*-functional` → `*-ui`. Omit a seam when that concern is absent.

## Decision tree

**Always size** on every `/sc-discuss` (not only large asks). Default `Sizing: single` when the work fits one mission.

1. **Single mission (default)** - One primary user capability; estimated ≤7 jigsaw tasks after discuss; fold needed checklist concerns into that mission (discuss → plan → run). Draft still required when visual.
2. **Same-mission phases** - Still one capability / one ship story, but >7 tasks and seams are **not** independently shippable → record `Sizing: phases` (+ phase count / rationale). **Discuss owns the phases decision;** `/sc-run` planning may write `plan-phaseN.json` (never a multi-mission map).
3. **Multi-mission roadmap (Must when)** any of:
   - Estimated >7 tasks **and** at least two of the **3 seams** can ship/verify independently, **or**
   - Any single seam alone would need **≥4** jigsaw tasks (hard Must, not "roughly"), **or**
   - Visual draft HIL would block unrelated data/functional work if kept in one mission

Then create the map under **Map creation (discuss only)** below. Discuss the current tip only; later seams wait for post-ship handoff (`Next: /sc-discuss <id>`).

## Map creation (discuss only)

`spacecraft map` is queue plumbing for `/sc-run` - not a planning ceremony. Only `/sc-discuss` (this playbook) may create or resize a multi-mission map.

**Ordered steps (Must):**

1. Record intended `Sizing: roadmap <id>` (+ seams / rationale) - will copy onto every seam in step 5.
2. For each needed seam title: `spacecraft new "<feature>-data|functional|ui"` (create stubs **before** `map add`).
3. `spacecraft map new "<feature-or-roadmap-title>"` **unless** the human explicitly approved reusing an existing roadmap id. Never silently append to unrelated `map current`.
4. `spacecraft map add <roadmap-id> <mission-id> --desc "<seam title>"` for each stub, in order `data` → `functional` → `ui` (omit absent).
5. **Stub `decisions.md` on every seam** (tip and later) with the same:
   ```
   Sizing: roadmap <roadmap-id>
   Sizing seams: <feature>-data, <feature>-functional, <feature>-ui
   Sizing rationale: <one line>
   ```
   Later seams may add `Mission brief: skipped - stub seam; discuss after prior tip ships` until their discuss turn.
6. `spacecraft map use <roadmap-id>`; `spacecraft use <tip-id>`; discuss **current tip only** (full spec / draft / brief). Leave later seams for post-ship handoff.

**Must not** call `map new` / `map add` from `/sc-run` planning (`sc-planning` / `sc-planner`). If scope needs a new or resized multi-mission split mid-plan → stop and apply **Resize protocol** below via `/sc-discuss`.

## Resize protocol (mid-plan / mid-run)

When planning or AFK discovers the mission must become multi-mission (or the map must change):

1. Stop AFK / planning. Do not `map new` from planning.
2. Hand to `/sc-discuss` with the in-progress mission id.
3. Discuss chooses one fate for the current mission (record in `decisions.md`):
   - **keep-as-seam** - retitle/repurpose current mission as one seam (usually the tip); create sibling stubs + map for the rest
   - **supersede** - archive or abandon current as wrong shape; create fresh seam stubs + new map
   - **re-tip** - current stays on map but is no longer tip; reorder / `map add` siblings; discuss the new tip
4. Re-record `Sizing: roadmap <id>` on every involved mission; clear only the tip after its discuss gates.

## Handoff by sizing

After tip `clarify-status clear`:

| `Sizing:` | Next |
|-----------|------|
| `roadmap <id>` | **Spec clear. New session: `/sc-run <id>`.** (`map use` that id) |
| `single` or `phases` | **Spec clear. New session: `/sc-run`.** (mission-only AFK on resolved current mission; no multi-mission map required) |

## Visual draft by seam

- Draft HIL required for `*-ui` and for `Sizing: single` / `phases` when the mission is visual UI/FE.
- `*-data` and `*-functional` tips: **non-visual** - record `UI draft skipped: non-visual seam (<data|functional>)` in that mission's `decisions.md`. `/sc-run` must not demand a draft for those tips.

## Record

In each involved mission `decisions.md` (at map create for roadmap; at discuss for single/phases):

```
Sizing: single | phases | roadmap <roadmap-id>
Sizing seams: <feature>-data, <feature>-functional, <feature>-ui   # omit absent; roadmap only
Sizing rationale: <one line>
```

For `phases`, also record: `Sizing phases: <N> - <one-line rationale>`.

Roadmap description should state feature name, seam order, and why split.

## Must not

- Global layer roadmaps across unrelated features
- Fourth `*-ux` roadmap seam
- Soft "prefer ≤7" or 8-9 exception bands (≤7 remains a hard Must per phase)
- Soft "roughly ≥4" for single-seam split (use hard ≥4)
- Inventing roadmap create/resize mid-`/sc-run` or mid-plan without returning to `/sc-discuss`
- Planning-owned "create a roadmap" as an escape hatch (multi-mission is a discuss sizing outcome)
- Silent reuse of `map current` for a different feature
- `map add` without prior `spacecraft new` stubs
- Leaving later seam `decisions.md` without `Sizing:` at map create
