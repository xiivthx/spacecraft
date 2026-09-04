# Fact-check

SoT for mission-affecting **external** claims (API, version, compat, docs, open-domain synthesis).

## Disposition

When this protocol runs, emit exactly one greppable line (silence forbidden):

- `Fact-check: corroborated`
- `Fact-check: contested: <claim-id>`
- `Fact-check: skipped: <reason>`

Contested → escalate (more search or ask); no auto-pick. Cite + short note in `decisions.md`.

## When

After `sc-search` / `sc-storm` (or clarify web auto-pick) before treating an external fact as settled.

**Skip second agent** (still emit `Fact-check: skipped: …` when research stated the fact):

| Reason | Line |
|--------|------|
| Path cite in repo | `Fact-check: skipped: repo-local` |
| Tier 1+2 agree on one primary; no contradiction | `Fact-check: skipped: single-primary` |
| Non-blocking | `Fact-check: skipped: non-blocking` |
| Preference / Verify / architecture | `Fact-check: skipped: preference-or-gate` |

**Must `Task(sc-fact-check)`** when sources disagree, a version/API/compat pin drives implement or auto-pick, or the human asks for fact-check.

## Token budget

| Cap | Limit |
|-----|-------|
| Claims | ≤5 |
| Per claim | 1 line + 1 cite (URL or `path:lines`) |
| Critic input | claim block only |
| Critic Tasks | one per batch |
| Extra fetches in critic | ≤2, only if cite missing or conflict |

Optional `NOTE <id>: <≤12 words>`.

## Claim block

```
CLAIM <id>: <one-line fact>
CITE <id>: <url-or-path>
```

## Critic → disposition

Critic returns per claim: `OK <id>` | `CONTEST <id>: <≤12 words>` | `INSUFFICIENT <id>: <≤12 words>`.

- all `OK` → `Fact-check: corroborated`
- any `CONTEST` / `INSUFFICIENT` → `Fact-check: contested: <id>` (first id; list rest in `decisions.md`)

## Shape

1. Primary builds claim block from sources.
2. One `Task(sc-fact-check)` with claim block only (readonly; different model family).
3. Primary writes disposition; escalate if contested.

## Must not

- Treat disposition as `AUTH` / `VERIFIED` / ready / ship
- Pass chat, full `spec.md`, or full research brief to the critic
- Spawn one critic Task per claim
- Soft-pass contested claims
