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

Then: `spacecraft map new` (or use current map) and add only the needed `<feature>-data` / `<feature>-functional` / `<feature>-ui` missions. Discuss the current tip only; later seams wait for post-ship handoff (`Next: /sc-discuss <id>`).

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
- Inventing roadmap resize mid-`/sc-run` without returning to `/sc-discuss`
