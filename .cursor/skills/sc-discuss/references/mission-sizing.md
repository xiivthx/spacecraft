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
| **functional** | `<feature>-functional` | API / domain / Verify (`*-api` only if the team already uses that name) |
| **ui** | `<feature>-ui` | UX draft HIL in discuss + port implementation |

**Must not** create `*-ux` as a roadmap seam.

**Order default:** `*-data` → `*-functional` → `*-ui`. Omit a seam when that concern is absent.

## Decision tree

1. **Single mission (default)** - One primary user capability; estimated ≤7 jigsaw tasks after discuss; fold needed checklist concerns into that mission (discuss → plan → run). Draft still required when visual.
2. **Same-mission phases** - Still one capability / one ship story, but >7 tasks and seams are **not** independently shippable → `plan-phaseN.json`.
3. **Multi-mission roadmap (Must when)** any of:
   - Estimated >7 tasks **and** at least two of the **3 seams** can ship/verify independently, **or**
   - Any single seam alone would need roughly ≥4 tasks, **or**
   - Visual draft HIL would block unrelated data/functional work if kept in one mission

Then create the map under **Map creation (discuss only)** below. Discuss the current tip only; later seams wait for post-ship handoff (`Next: /sc-discuss <id>`).

## Map creation (discuss only)

`spacecraft map` is queue plumbing for `/sc-run` - not a planning ceremony. Only `/sc-discuss` (this playbook) may create or resize a multi-mission map:

1. Record `Sizing: roadmap <id>` (+ seams / rationale) in `decisions.md`.
2. `spacecraft map new` (or use current map) + add only needed `<feature>-data` / `<feature>-functional` / `<feature>-ui` missions.
3. Discuss the current tip; leave later seams for post-ship handoff.

**Must not** call `map new` / `map add` from `/sc-run` planning (`sc-planning` / `sc-planner`). If scope needs a new or resized multi-mission split mid-plan → stop and hand to `/sc-discuss`.

## Record

In each involved mission `decisions.md`:

```
Sizing: single | phases | roadmap <roadmap-id>
Sizing seams: <feature>-data, <feature>-functional, <feature>-ui   # omit absent; roadmap only
Sizing rationale: <one line>
```

Roadmap description should state feature name, seam order, and why split.

## Must not

- Global layer roadmaps across unrelated features
- Fourth `*-ux` roadmap seam
- Soft "prefer ≤7" or 8-9 exception bands (≤7 remains a hard Must per phase)
- Inventing roadmap create/resize mid-`/sc-run` or mid-plan without returning to `/sc-discuss`
- Planning-owned "create a roadmap" as an escape hatch (multi-mission is a discuss sizing outcome)