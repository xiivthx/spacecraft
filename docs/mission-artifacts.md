# Mission artifact schemas

Canonical shapes for `.space/` artifacts. Skills point here; do not paste full schemas into every prompt.

## mission.json

```json
{
  "id": "<mission-id>",
  "title": "<title>",
  "state": "active|planned|in_progress|ready|blocked|shipped",
  "branches": ["<branch-name>"],
  "createdAt": "<iso-date>",
  "pickup": {
    "phase": "<discuss|run|ship|quick|…>",
    "next": "<one-liner>",
    "blockers": ["<optional>"],
    "updatedAt": "<iso-date>"
  }
}
```

Optional `pickup` for handoffs (`spacecraft status`). Not a ship gate.

## plan.json

Each `acceptance[]` item = one RED-GREEN cycle. ≤7 tasks per phase; overflow → `plan-phaseN.json` (`Sizing: phases`).

```json
{
  "planName": "<short-name>",
  "missionId": "<mission-id>",
  "tasks": [
    {
      "id": "T1",
      "title": "<imperative slice>",
      "status": "pending|in_progress|done|blocked",
      "dependsOn": [],
      "files": ["<path>"],
      "acceptance": ["<one check>"],
      "verify": "<command>",
      "evidence": ["<label>"]
    }
  ]
}
```

## evidence.jsonl

One JSON object per line. `output` = real captured stdout (never fabricated). Optional `outputHash` = SHA-256 of full raw.

```json
{"label":"<label>","command":"<command>","output":"<output>","outputHash":"<hex>","ts":"<iso>"}
```

When raw > 65536 bytes: truncate JSONL `output` with `\n...[truncated]`; full raw under `evidence-raw/`; set `outputTruncated`, `outputBytes`, `outputRawPath`. `validate` hashes the sidecar when truncated.

## Roadmap (`.space/roadmaps/<id>.json`)

```json
{
  "id": "<kebab-slug>",
  "title": "<title>",
  "description": "<description>",
  "missions": [{"id": "<mission-id>", "description": "<brief>"}],
  "issues": [{"number": 1, "title": "...", "url": "https://...", "state": "open", "phase": "mvp|polish|ship"}],
  "createdAt": "<iso>",
  "updatedAt": "<iso>"
}
```

Manage with `spacecraft map`.

## design-contract.md

After plan, before product RED/GREEN. Sections: Scope · Modules · Data structures · Public seams · Edge matrix (expected literals) · Out of scope. Footer: `Design-contract: complete`. Docs/prose-only → skip line in **Outcome-gate** below (no file).

## approved-scenarios.md

After design-contract, before product RED/GREEN. Freeze Edge matrix + spec examples. Footer: `Approved-scenarios: frozen-from-contract` or `frozen-by-human`. Must not thaw frozen literals without Commander + `Scenario oracle change:`. Docs/prose skip when design-contract skip applies.

### Test freeze

Gates ≥ `M9G7IHV3`: `spacecraft freeze` / `freeze-check` when freeze footer present (exempt on `Approved-scenarios skipped:`). Checks: `freeze-drift`, `postdated-freeze`; cross-model critic line required. Gates ≥ `M9G7II1F` + `.space/config.json` `criticFamily` → match or fail (`configured-but-skipped` / `critic-family-mismatch`).

## Outcome-gate skip / waive grammar (SoT)

Do not redefine these prefixes elsewhere.

| Gate | Greppable line(s) |
|------|-------------------|
| design-contract | `Design-contract skipped: docs/prose-only` |
| approved-scenarios | `Approved-scenarios skipped: docs/prose-only` |
| static-analysis | `Static-analysis skipped: no project static tool` · `Static-analysis waived: <reason>` |
| diff-coverage | `Diff-coverage skipped: no project coverage tool` · `Diff-coverage waived: <reason>` |
| mutation | `Mutation skipped: no project mutation tool` · `Mutation skipped: not in scope` · `Mutation waived: <reason>` |
| PBT | `Pbt skipped: no project pbt tool` · `Pbt skipped: not core logic` · `Pbt waived: <reason>` |
| NFR no tool | `<Gate> skipped: no tool` (pair with debt line; never invent pass) |
| characterization | `Characterization waived: <reason>` |

## NFR

```
NFR source: user | measured-baseline | default(<why>)
```

Relative bars vs baseline evidence id (`relative-bar`). Unknown preference-bound bars → sc-clarify.

## Quality debt

Debt (not silent green): `Mutation skipped: no tool` / `<Gate> skipped: no tool` / `Characterization waived:`. `Mutation skipped: not in scope` is valid, not debt. Ceiling: `Debt ceiling: 3`. Grammar after Gates `M9G7IHON`.

## Outcome evidence labels

| Gate | Label prefix | Bar |
|------|--------------|-----|
| Static | `static-` | 0 warning / 0 error when tool runs |
| Diff coverage | `diff-cov-` | touched line+branch ≥90% (band 90–95%); no global 95–100% chase |
| Mutation | `mutation-` | In scope: `Mutation: required` \| pack `quality` \| `Mutation: high-risk`; score >80% or skip/waive |
| PBT | `pbt-` | 100% of design-contract core-logic modules, or skip/waive |
| Freeze | `freeze` | sha256 manifest; `freeze-check` at combine/ready |

## Other

- `spec.md` — what/why
- `decisions.md` / `questions.md` — decisions / opens
- `review.json` — findings; ready only on empty findings + judge `VERIFIED`. Overview: [review.md](./review.md). SoT: `mission-review-gates.md` / `ux-ui-review-gates.md`

Quick + Mission: INTENT / AUTH / TWINS / 3-cycle. Judge evidence labels include `judge`. AUTH does not bypass ship hooks.
