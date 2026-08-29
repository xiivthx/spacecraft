# Impeccable × `/sc-discuss` integration contract

**Status: enacted in prompts** - live SoT is `sc-ux-design` + `impeccable-orchestration` + `sc-designer`. Pilot checklist below is optional human verification, not a claim of completed pilot.

Related live SoT: `docs/ux-ui-review.md`, `.cursor/skills/sc-ux-design/references/impeccable-orchestration.md`, `.cursor/agents/sc-designer.md`, `.cursor/skills/sc-ux-design/SKILL.md`, `.cursor/skills/sc-discuss/SKILL.md`.

## Purpose

**Impeccable is the primary UX/UI craft engine.** Use the **full** Impeccable command catalog by fitness (`init`, `shape`, new-work describe, `polish`, `critique`, `audit`, `live`, plus refine/enhance/fix family) - invoker may be Commander, `sc-designer` Next, or human slash. Spacecraft keeps mission gates and the port visual SoT (approved draft HTML). `sc-designer` routes commands + owns port gates.

Canonical command map: `.cursor/skills/sc-ux-design/references/impeccable-orchestration.md`.

Spacecraft already uses `npx impeccable detect` for anti-slop CLI. Detect alone does not replace `critique` or `audit`.

## Locked HIL decisions

| # | Decision |
|---|----------|
| 1 | **Brief:** `/impeccable shape` **replaces** the sc-ux-design 6-dimension brief entirely on the Impeccable path. |
| 2 | **Craft critique:** `/impeccable critique` default; finish-reviewer when approved comp / craft-critical. |
| 3 | **Paths:** `PRODUCT.md` / `.impeccable/` / package `DESIGN.md` at the **UI package**. **Must** gitignore `.impeccable/` (`templates/gitignore`). |
| 4 | **Draft polish:** opt-in Operate; default on Persuade / craft-critical. |
| 5 | **House DESIGN.md:** default **keep house**; update only on `DESIGN conflict: update house`. |
| 6 | **Primacy:** Impeccable = craft owner (full command catalog by fitness); invoker-agnostic; `sc-designer` = router + port gates; default visual discuss `Impeccable path: active`. |

## Artifact map

| Concern | Spacecraft SoT | Impeccable SoT | Rule |
|---------|----------------|----------------|------|
| Product intent | `spec.md`, `decisions.md` | `PRODUCT.md` at UI package | Spec/decisions win for mission |
| House visual system | package `DESIGN.md` | same + sidecar | One house file; default keep house |
| Mission visual SoT | approved draft HTML + `UI draft approved:` | comps under `.impeccable/mocks/` (gitignored) | Port SoT = draft HTML only |
| Design brief | *(leaned out when path active)* | `/impeccable shape` | Shape only; record `Impeccable brief approved:` |
| Layout choice | bake-off HTML (2-3) | comps as generation refs | Record winner/skip; comps ≠ substitutes |
| Scenario / checklist | scenario matrix + `UX checklist:` | not owned | Spacecraft-only |
| Critique | Task(`sc-designer`) port gates first | critique / finish-reviewer | Order: port → craft → human HIL |
| Orchestration | `sc-designer` Next + phase | Impeccable slash commands | Designer does not edit; Commander runs commands |
| Run port | approved draft | optional polish/audit | No redesign in run |

## Command sequence

Canonical step table: `.cursor/skills/sc-ux-design/references/impeccable-orchestration.md`. Shape → comps → bake-off HTML → scenarios → `sc-designer` port → Impeccable craft → `UI draft approved`; keep house unless `update house`.

## `/sc-run` boundary

No shape / new-work / bake-off / redesign. Port from approved draft. Optional polish/audit with draft parity. `sc-designer` run-port = live-product + draft-parity.

## Hard rules

### Must

- Impeccable primary for craft on visual missions (unless path skipped).
- Single port SoT = approved draft HTML.
- `sc-designer` before craft gate before human HIL.
- Gitignore `.impeccable/` at UI package.

### Must not

- Comps as port SoT; parallel 6-dimension brief when path active; silent bake-off skip; redesign in run; soft-pass port ↔ craft; commit `.impeccable/`; default-update house without `update house`.

## Pilot checklist

- [ ] UI package gitignores `.impeccable/`
- [ ] `Impeccable path: active` recorded
- [ ] Shape brief approved; no parallel sc-ux brief
- [ ] `Impeccable direction:` recorded when comps used
- [ ] Bake-off HTML winner/skip recorded
- [ ] Scenario + checklist + ladder complete
- [ ] `sc-designer` port pass then Impeccable craft pass (or human waive)
- [ ] `UI draft approved:` → run ports that file only
- [ ] House unchanged unless `update house`
