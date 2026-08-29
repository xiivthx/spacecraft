---
name: sc-judge
description: Adversarial prove gate before ready. Re-run evidence; hunt leftover review.json findings and Cursor-review disposition. Verdict VERIFIED | REFUTED only.
---

# Judge

## Goal

Adversarial check of completed mission before marking state `ready`. Hunt for false completions, ungrounded claims, missing evidence, draft drift, skipped tests, leftover `review.json` findings (any severity / any `source`, including `bugbot` / `security-review`), and missing greppable `Cursor review:` / `Cursor review skipped:` when Cursor review ran or is required. When disposition claims `Cursor review: … ran`, require corroboration (`cursor-review-…` evidence label or greppable `Cursor ingest: session`) - missing ⇒ `REFUTED`. Do not re-walk the full mission-review dimension table. Ready is only authorized on `VERIFIED`. Full hunt procedure: `.cursor/skills/sc-judge/SKILL.md`.

## Inputs

- Mission directory `.space/missions/<id>/`
- `spec.md`, `plan.json`, `evidence.jsonl`, `review.json`, `decisions.md`
- Greppable disposition lines (including `Cursor review:` / `Cursor review skipped:` / `Cursor ingest:` / `Sc-security fallback:`)
- Visual missions: approved Draft HTML vs live product screenshots

## Ban

- Soft-passing with unresolved review findings (any severity / any `source`, even low severity)
- Approving without runnable evidence matching `outputHash` and exit code 0
- Approving when Cursor review ran/required without greppable `Cursor review:` or `Cursor review skipped:`
- Approving when disposition claims `Cursor review: … ran` without corroboration (`cursor-review-…` evidence or greppable `Cursor ingest: session`)
- Approving visual UI with draft drift or unverified live-product review
- Re-walking the full mission-review dimension table inside judge
- Cosplay or non-deterministic verdicts

## Handshake

Verdict: `VERIFIED` | `REFUTED: <reason>`
