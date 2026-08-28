# Impeccable orchestration (Spacecraft UX)

Authoritative workflow when Impeccable is the **primary** UX/UI craft engine inside Spacecraft missions. Loaded by `sc-designer` and by `/sc-discuss` / `sc-ux-design` on visual work.

Contract overview: `docs/impeccable-discuss-integration.md`.

## Primacy

| Owner | Owns |
|-------|------|
| **Impeccable** | Full craft command set (see Command catalog) |
| **Spacecraft** | Mission gates: product context, UX checklist, extract/borrow, bake-off HTML SoT, scenario matrix, responsive ladder, `UI draft approved`, `/sc-run` port, live-product / draft-parity |
| **Invoker** | Whoever runs UX work (Commander, Task(`sc-designer`) Next, or human slash) **Must** pick the matching Impeccable command below — not invent a parallel craft path |

Port visual SoT remains approved draft HTML under `.space/missions/<id>/design/drafts/`. Impeccable comps are craft north stars only. Invoker identity does not matter; **command fitness** does.

## Default path

On visual UI/FE discuss: record `Impeccable path: active` unless human records `Impeccable path: skipped: <reason>` (legacy sc-ux 6-dimension brief path).

## UI package roots

Resolve `<ui-package>` as the package that owns the surface (not monorepo root unless that package is the root).

| Artifact | Location |
|----------|----------|
| `PRODUCT.md`, `DESIGN.md` | `<ui-package>/` |
| `.impeccable/` | `<ui-package>/.impeccable/` — **Must** be gitignored |
| Draft HTML | `.space/missions/<id>/design/drafts/` |
| Refs | `.space/missions/<id>/design/refs/` |

## Command catalog (use the right slash)

Map intent → Impeccable. Prefer the **primary** command; add **follow-ons** when the finding/need matches. Do not skip a listed primary when its trigger fires on an active path.

### Core lifecycle (user-facing defaults)

| Intent | Command | When in Spacecraft |
|--------|---------|-------------------|
| Start a project / capture product truth | `/impeccable init` | Missing `PRODUCT.md` at UI package; new product or first visual mission in package |
| Plan UX before code | `/impeccable shape <surface>` | Discuss brief — **only** brief when path active (replaces sc-ux 6-dimension brief) |
| Build a new surface | `/impeccable <describe surface>` (new-work) | After shape; new/replacement composition or visual world; comps / direction lock |
| Polish a page | `/impeccable polish <target>` | Winner draft or live UI after structure is right; Persuade default-on; Operate opt-in / before ship |
| Find design issues | `/impeccable critique <target>` | Default craft gate before `UI draft approved`; also mid-iterate when human asks for design review |
| Check implementation | `/impeccable audit <target>` | After `/sc-run` port or on live product: a11y / responsive / technical QC alongside Spacecraft live-product |
| Iterate in the browser | `/impeccable live` | Human wants variant exploration on a running surface; discuss polish rounds or post-port tweak **without** reopening art direction |

### Build / system

| Intent | Command | When |
|--------|---------|------|
| Record house from shipped UI | `/impeccable document` | Only if `DESIGN conflict: update house` or house `DESIGN.md` missing after approval |
| Pull tokens/components into system | `/impeccable extract <target>` | After stable port; promoting patterns into house (needs update-house or explicit extract ask) |

### Refine / enhance / fix (on demand)

Use when findings or human ask match — usually **after** shape + direction, on draft or live target. Still honor dimension lock in discuss (one of typography | color | layout | motion | spacing | chrome per human round).

| Intent | Command |
|--------|---------|
| Too safe / bland | `/impeccable bolder <target>` |
| Too loud | `/impeccable quieter <target>` |
| Too complex | `/impeccable distill <target>` |
| Errors / i18n / edge cases | `/impeccable harden <target>` |
| First-run / empty / activation | `/impeccable onboard <target>` |
| Motion | `/impeccable animate <target>` |
| Weak color | `/impeccable colorize <target>` |
| Type hierarchy | `/impeccable typeset <target>` |
| Spacing / rhythm | `/impeccable layout <target>` |
| Personality | `/impeccable delight <target>` |
| Ambitious effects | `/impeccable overdrive <target>` |
| Copy / labels / errors | `/impeccable clarify <target>` |
| Device / breakpoint fit | `/impeccable adapt <target>` |
| UI performance | `/impeccable optimize <target>` |

`craft` is deprecated → treat as new-work. Pin/hooks/doctor are tooling, not mission craft steps (run when human asks or hook status requires).

### Anti-patterns

- Do **not** use `shape` / new-work / `live` redesign during `/sc-run` after `UI draft approved` (unless human reopens discuss).
- Do **not** substitute sc-ux 6-dimension brief for `shape` when path active.
- Do **not** skip `critique` (or finish-reviewer when triggered) before approval HIL on path active.
- Do **not** treat `npx impeccable detect` alone as full `critique` or `audit`.

## Discuss sequence (path active)

Invoker executes; `sc-designer` returns the next **exact** Impeccable/Spacecraft command when invoked.

| Step | Owner | Action | `decisions.md` |
|------|-------|--------|----------------|
| 0 | SC | Opt-in (default active) | `Impeccable path: active` |
| 1 | SC | Product context + UX checklist | `Product context:` / skip; `UX checklist:` |
| 2 | IMP | `/impeccable init` if no `PRODUCT.md` | — |
| 3 | SC | Reference extract + borrow when refs supplied | `Reference extract:`; `Reference borrow:` |
| 3b | SC | Context fidelity | `Context fidelity: …` |
| 4 | IMP | `/impeccable shape <surface>` | `Impeccable brief approved: …` |
| 5 | IMP | `/impeccable <describe>` new-work / comps as needed | `Impeccable direction: …` |
| 6 | SC | House `DESIGN.md` read/conflict | `DESIGN conflict:` (default keep house if file exists) |
| 7 | SC | Bake-off 2–3 draft HTML (comps = refs, not substitutes) | `Layout bake-off winner:` / skip |
| 8 | BOTH | Scenario matrix + ladder (SC Must). Then Impeccable refine as needed: `polish` (Persuade default / Operate opt-in), else targeted `typeset` / `layout` / `clarify` / `harden` / `onboard` / `live` | `Impeccable draft polish: on \| skipped` (+ optional `Impeccable refine: <cmd>`) |
| 9a | SC | Task(`sc-designer`) **port gate** first | — |
| 9b | IMP | `/impeccable critique <draft>` default; finish-reviewer when approved comp / craft-critical | `Impeccable craft gate:`; `Impeccable craft:` |
| 10 | SC | Human HIL | `UI draft approved: <file>` |
| 11 | IMP/SC | `/impeccable document` only if update-house or missing house | default keep house |

**Gate order:** port → craft (`critique`) → human HIL. Never soft-pass port ↔ craft.

## Shape replaces Spacecraft brief

When path active: shape is the only brief; prompt assembly tail = shape brief + `spec.md` Musts + checklist. Path skipped: legacy 6-dimension brief.

## Run boundary

| Allowed | Forbidden |
|---------|-----------|
| `/impeccable polish` / `audit` / `adapt` / `optimize` / `clarify` / `harden` on **ported** UI if draft parity preserved | `shape`, new-work, bake-off, `live` art-direction reopen, `bolder`/`quieter` that change the approved world without discuss |
| Task(`sc-designer`) run-port (live-product + draft-parity) | Treating comps as port SoT |

Record optional `Impeccable run assist: polish|audit|…` in `decisions.md` when used.

## `sc-designer` invocation map

| Phase | Designer job |
|-------|----------------|
| Mid-discuss | Orchestrate: missing records + **Next = exact `/impeccable …` or SC step**; may list **Also consider** one secondary command from the catalog |
| Bake-off | Scaffold + responsive ladder only |
| Approval-port | Full port dimensions; craft pending → Next = `/impeccable critique` (or finish-reviewer rule) |
| Approval-craft-check | Verify craft record; if failed findings need refine → Next = matching refine cmd then re-critique |
| `/sc-run` visual | live-product + draft-parity; Next = `/impeccable audit` and/or `polish` when technical/craft gaps remain **without** redesign |

When Next is an Impeccable command, write it as a copy-pasteable slash line (e.g. `/impeccable polish .space/missions/<id>/design/drafts/<file>.html`).

## Craft waive

`Impeccable craft: waived: <reason>` only with explicit human reason. Designer must not waive. Missing craft on path active without waive = critical before approval HIL.

## Finish-reviewer trigger

Use finish-reviewer when `Impeccable direction:` points at an approved comp under `.impeccable/mocks/` or surface is Persuade / replacement world / craft-critical redesign. Otherwise `/impeccable critique` suffices.
