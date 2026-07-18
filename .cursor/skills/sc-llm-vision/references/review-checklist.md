> Consult when: categorizing findings, filing GitHub issues, merging with DOM QA, or auditing anti-patterns.

# Review checklist

## Severity

| Severity | Criteria |
|----------|----------|
| **critical** | Crash, content missing, page unusable, security/brand leak |
| **important** | Design system violated, layout broken, accessibility blocked, data inconsistency |
| **minor** | Pixel misalignment, wrong color token, typography off-spec, polish |
| **info** | Nice-to-have improvements, taste suggestions |

Deduplicate: same issue across routes → one systematic issue, not N per-route issues.

## Create GitHub issues

Group related findings by theme and route:

```bash
cat << 'BODY' | gh issue create --title "...(route+theme)..." --body-file -
## Source
Gemini 3.5 Flash fovea review - {N} screenshots

## Description
...

## Evidence
Screenshots: {filenames}

## Suggested Fix
...

## Verification Criteria
...
BODY
```

Grouping rules:

- One route's desktop issues → one issue
- One route's mobile/responsive issues → one issue
- Systematic patterns across routes → one cross-cutting issue
- Pre-existing coverage → reference existing issue instead of duplicating

## Merge with automated QA

Combine LLM vision findings with Playwright DOM QA (overflow, tokens, hit areas, accessibility). LLM finds visual feel, composition, brand fit that DOM cannot.

## Evidence layout

```
.space/screenshots/v{version}/
  _gemini-survey.txt
  _gemini-fovea-{route}.txt
  _gemini-review-prompt.md
  _merged-design-report.json   # optional
```

## Anti-patterns

- Don't send all screenshots in one batch
- Don't omit DESIGN.md from the prompt
- Don't skip fovea for complex routes
- Don't create 1 issue per finding - group by route + theme
- Don't run without `--dangerously-skip-permissions`
- Don't run batches back-to-back
- Don't trust the model on data correctness - verify counts independently

## Integration with /sc-review

```
/sc-review
  ├── Playwright DOM QA
  ├── Gemini Survey (1 batch)
  ├── Gemini Fovea (N batches, 30min between)
  └── Merge → dedupe → GitHub issues → review.json
```

## Spacecraft integration

Link issues in mission `issues.md` when present. Fold release-blocking findings into `review.json` severity groups before ship.
