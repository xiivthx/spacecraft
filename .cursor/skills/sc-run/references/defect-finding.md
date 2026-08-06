# Defect finding craft

On-demand craft for writing defect findings during `/sc-run` fix/review and when `sc-reviewer` emits findings. Decision job - not expertise cosplay. **Not** an always-on discuss gate. **Not** a bug tracker.

## Goal

Produce actionable, impact-clear findings in `review.json` and run summaries so `/sc-run` can fix and re-verify without inventing repro, version, or environment details.

## Output

Chat summary plus findings in `review.json` (or run summary when mid-build, before review):

```json
{
  "findings": [
    {
      "severity": "critical",
      "title": "Checkout fails on expired card",
      "file": "src/checkout/handler.ts",
      "issue": "POST /checkout returns 500 when card expiry is in the past.\nUsers cannot complete purchase with a valid new card after one expired attempt.",
      "repro": [
        "Submit checkout with expired card",
        "Replace with valid card and resubmit"
      ],
      "impact": "Purchase blocked after first expired-card attempt",
      "businessRisk": "Lost revenue on recoverable payment errors",
      "requiredFix": "Clear stale card state; return 4xx for expired card",
      "retest": [
        "Expired card then valid card completes checkout",
        "Evidence: checkout-expired-card",
        "Negative: expired card alone returns 400"
      ],
      "notes": "Screenshot: evidence/checkout-500.png",
      "environment": "unspecified",
      "version": "unspecified",
      "reproducible": "evidence: checkout-expired-card"
    }
  ]
}
```

**Compact minor** (optional fields omitted when trivial):

```json
{
  "severity": "minor",
  "file": "src/ui/badge.tsx",
  "issue": "Badge uses 11px font; house minimum is 12px.",
  "requiredFix": "Set font-size to 12px per DESIGN.md"
}
```

## Good / Bad

- Good: impact-first title (≤12 words); 2-3 line issue (problem + why); clear `requiredFix`; 2-3 `retest` ideas for critical/important; evidence-backed `reproducible`; environment/version from mission Platform / SFDIPOT / evidence or `unspecified`; business risk one line when critical/important; visual UI notes point to evidence path or "attach / evidence"
- Bad: expertise cosplay; inventing repro/version/env; default browser to Edge; inventing "Yes, Thrice" repro counts; similar bug stories unless known from repo/history; issue ledgers (`issues.md`); essay dumps; adopting High/Medium/Low/Lowest severity labels

## Verify

Findings appear in `review.json` or run summary with greppable `severity` + `requiredFix`. Critical/important include `title`, user `impact`, and 2-3 `retest` entries. No invented optional fields. Commander validates via review drain before `ready`.

## When to use

- `sc-reviewer` emits or expands findings (especially critical/important)
- `/sc-run` fix pass or findings mid-build when recording defects for summary
- Human asks for structured defect write-up during review

## When skip

- Routine minor style nits (compact `file` + `issue` + `requiredFix` enough)
- Finding already fully proven by evidence label (omit redundant `repro`)
- Discuss-phase ambiguity (use sc-clarify / `questions.md` instead)

## Finding schema

| Field | Required | Notes |
|---|---|---|
| `severity` | yes | `critical` \| `important` \| `minor` |
| `title` | critical/important | Impact-first, ≤12 words |
| `file` | yes | Path or scope |
| `issue` | yes | 2-3 lines: problem + why it matters |
| `repro` | when needed | Steps array; omit if evidence already proves |
| `impact` | critical/important | User-facing effect |
| `businessRisk` | optional | One short line when critical/important; omit trivial minor |
| `requiredFix` | yes | Concrete remediation |
| `retest` | critical/important | 2-3 local verify ideas (evidence labels or manual checks) |
| `notes` | optional | Screenshots → evidence path or "attach / evidence" |
| `environment` | optional | Mission Platform / SFDIPOT / evidence; default `unspecified` |
| `version` | optional | From evidence or tag; default `unspecified` |
| `reproducible` | optional | Evidence label (e.g. `evidence: label`) or manual count (e.g. `manual: 3/3`) |

## Severity mapping

| Impact signal | Severity |
|---|---|
| Data loss, security breach, crash, blocked core flow, wrong money/auth | `critical` |
| Degraded UX, wrong non-critical behavior, missing error handling on touched path | `important` |
| Cosmetic, docs typo, minor inconsistency off critical path | `minor` |

External High → `critical`; Medium → `important`; Low/Lowest → `minor`. Do not adopt High/Medium/Low/Lowest in house findings.

## Clarify-before-assume

When bug description is ambiguous (repro, version, environment, impact):

- **Mid-discuss** - sc-clarify or park in `questions.md`
- **Mid-run soft** - record assumption in `decisions.md`
- **Mid-run hard** - stop for `/sc-discuss`
- **Never** invent repro steps, version, or environment

## Must / Must not

- **Must**: Use house severity enum (`critical` \| `important` \| `minor`)
- **Must**: Ground optional fields in evidence, mission context, or `unspecified`
- **Must not**: Expertise cosplay ("as a QA expert…")
- **Must not**: Invent details (repro, version, env, similar stories)
- **Must not**: Revive issue ledgers; findings stay in `review.json` / run summary
- **Must not**: Unicode em dash - ASCII hyphen-minus only

## Related

- Used by `sc-reviewer` + `/sc-run` fix/summary
- `retest` ideas: borrow from Test Ideas buckets (Positive / Negative / Edge) in `decisions.md` when present
- SFDIPOT / Platform context: `sfdipot-coverage.md`, mission `spec.md`
- Optional oracle grounding for `issue`/`impact` "why": `sc-discuss/references/test-oracles.md`; optional `notes` may cite fired oracle letters - never invent oracles without observation

### Domain scans

Domain rules and skills (`300-security`, `400-performance`, `500-database`, `sc-security`, `sc-performance`) map checklist priorities to house severity before filing review findings - see each Finding Format.

| Domain | Checklist → House |
|---|---|
| Security | Critical/High → critical; Medium → important; Low → minor; informational → scan note only |
| Performance | Critical → critical; High → important; Medium → important (or minor if off hot path); Low/note → minor or scan note |
| Database | Critical → critical; High → important; Medium → minor |
