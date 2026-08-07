# Mission artifact schemas

Canonical shapes for `.space/` artifacts. Always-on rules point here; do not paste full schemas into every prompt.

## mission.json

```json
{
  "id": "<mission-id>",
  "title": "<title>",
  "state": "active|planned|in_progress|ready|blocked|shipped",
  "branches": ["<branch-name>"],
  "createdAt": "<iso-date>"
}
```

## plan.json

Jigsaw tasks: each task is one behavioral slice of the feature. Each `acceptance[]` item is one RED-GREEN cycle under `/sc-run`.

```json
{
  "planName": "<short-name>",
  "missionId": "<mission-id>",
  "tasks": [
    {
      "id": "T1",
      "title": "<imperative jigsaw slice>",
      "status": "pending|in_progress|done|blocked",
      "dependsOn": [],
      "files": ["<path>"],
      "acceptance": ["<one check per RED-GREEN cycle>"],
      "verify": "<command>",
      "evidence": ["<label>"]
    }
  ]
}
```

When `Sizing: phases` (discuss-owned), same-mission overflow uses sequential files `plan-phase1.json`, `plan-phase2.json`, … (or `plan-phaseN.json` naming). Each file matches the `plan.json` shape above; ≤7 tasks per phase. Default active plan for build is the current incomplete phase file (Commander tracks which phase is active in `decisions.md` or by first pending phase).

## evidence.jsonl

One JSON object per line:

```json
{"label":"<label>","command":"<command>","output":"<output>","outputHash":"<hex>","ts":"<iso-timestamp>"}
```

`output` must be actual captured command output - never fabricated.

`outputHash` is optional: lowercase hex SHA-256 of `output`. Omitted via `omitempty`; entries without it remain valid (backward compatible).

## Roadmap (`.space/roadmaps/<id>.json`)

```json
{
  "id": "<kebab-slug>",
  "title": "<title>",
  "description": "<description>",
  "missions": [
    {"id": "<mission-id>", "description": "<brief>"}
  ],
  "issues": [
    {"number": 1, "title": "...", "url": "https://...", "state": "open", "phase": "mvp|polish|ship"}
  ],
  "createdAt": "<iso-date>",
  "updatedAt": "<iso-date>"
}
```

Roadmap ID: lowercase kebab-case from title. Manage with `spacecraft roadmap` (`map` alias).

## Other

- `spec.md` - what and why (free-form markdown)
- `decisions.md` / `questions.md` - assumptions and blocking questions
- `review.md` / `review.json` - release readiness review output (after `sc-judge`; ready only when verdict is `VERIFIED` and findings empty). Visual UI: UX/UI review gates (`.cursor/skills/sc-ux-design/references/ux-ui-review-gates.md`; overview `docs/ux-ui-review.md`)

## Inner-loop / judge

Quick and Mission both apply INTENT / AUTH / TWINS and the 3-cycle stop. Capture judge evidence with a label that includes `judge` (e.g. `judge-<mission-id>` or `judge-pass-validate`). `output` must be real command stdout - never hand-written. AUTH does not bypass ship hooks or `SPACECRAFT_SHIP=1`. Findings are fixed in `/sc-run` and listed in the run/ship summary.