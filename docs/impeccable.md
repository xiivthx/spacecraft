# Impeccable × discuss

Impeccable is the primary UX/UI craft engine. Spacecraft owns mission gates and the port visual SoT (approved draft HTML). `sc-designer` routes commands and owns port gates.

Orchestration SoT: `.cursor/skills/sc-ux-design/references/impeccable-orchestration.md`

## Locked decisions

| # | Decision |
|---|----------|
| 1 | Brief: `/impeccable shape` replaces the 6-dimension brief when path active |
| 2 | Craft: `/impeccable critique` default; finish-reviewer when approved comp / craft-critical |
| 3 | Paths: `PRODUCT.md` / `.impeccable/` / package `DESIGN.md` at UI package; gitignore `.impeccable/` |
| 4 | Draft polish: opt-in Operate; default on Persuade / craft-critical |
| 5 | House `DESIGN.md`: keep house unless `DESIGN conflict: update house` |
| 6 | Default visual discuss: `Impeccable path: active` |

## Hard rules

**Must:** Impeccable primary for craft (unless path skipped); port SoT = approved draft HTML; `sc-designer` port before craft before human HIL; gitignore `.impeccable/`.

**Must not:** comps as port SoT; parallel 6-dimension brief when path active; silent bake-off skip; redesign in `/sc-run`; soft-pass port ↔ craft; commit `.impeccable/`; update house without `update house`.

## `/sc-run`

No shape / new-work / bake-off / redesign. Port from approved draft. Optional polish/audit with draft parity.
