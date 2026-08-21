# RCRCRC impact pass

`/sc-discuss` (and re-plan under `/sc-run`) delta-only gate: analyze testing impact between two requirement versions using the RCRCRC heuristic. Decision job - not expertise cosplay.

## Goal

Turn a requirement rewrite into prioritized testing focus so planning and RED-GREEN cover what changed without blind regression gaps.

## Output

Chat summary plus greppable artifact in `decisions.md`:

```
## RCRCRC pass
- Key differences: …
- Recent - Observation: … | Testing Focus: …
- Core - Observation: … | Testing Focus: …
- Risky - Observation: … | Testing Focus: …
- Configuration Sensitive - Observation: … | Testing Focus: …
- Repaired - Observation: … | Testing Focus: …
- Chronic - Observation: … | Testing Focus: …
- Testing Priorities:
  1. …
  2. …
  3. …
```

Or skip when only one version exists:

```
RCRCRC pass skipped: Need both existing and updated requirements to perform RCRCRC analysis.
```

## Good / Bad

- Good: concrete diffs between versions; each RCRCRC category has Observation + Testing Focus; top-3 priorities usable by sc-planning / sc-tester
- Bad: running with a single version; inventing history for Chronic; expertise cosplay; vague "test everything"

## Verify

When two versions were available: `decisions.md` has `## RCRCRC pass`. When only one version: `RCRCRC pass skipped:` with the need-both message (or equivalent reason). Greppable markers only - no CLI validator required.

## When required

Existing **and** updated requirement are both available, e.g.:

- Human pastes two versions
- Mid-mission `spec.md` rewrite with prior text recoverable (git / decisions note / backup paste)
- Discuss replaces a previously cleared behavior bar

## When skip

Only one requirement version - respond exactly:

`Need both existing and updated requirements to perform RCRCRC analysis.`

Record `RCRCRC pass skipped: Need both existing and updated requirements to perform RCRCRC analysis.`

If both versions are vague or incomplete:

`Requirements unclear - please provide more detail.`

Keep clarify open; do not invent Verify or Chronic history.

## Procedure

1. **Requirement Comparison** - Summarize key differences between existing and updated requirement.

2. **RCRCRC Analysis** - For each category, provide:
   - **Observation** - what changed or matters
   - **Testing Focus** - what to test or explore

   | Category | Meaning |
   |----------|---------|
   | Recent | New areas of code or functionality |
   | Core | Essential functions that must not regress |
   | Risky | Features/areas more prone to defects |
   | Configuration Sensitive | Env, settings, integrations |
   | Repaired | Defect fixes / modifications that may create new issues |
   | Chronic | Areas with a history of frequent problems (`none found` allowed - do not invent) |

3. **Testing Priorities** - Short prioritized list (top 3) of what to test first and why. Write into `## RCRCRC pass`.

4. **Hand off** - sc-planning prefers these priorities when ordering jigsaw tasks / acceptance; sc-tester may use Testing Focus when choosing the RED scenario for the active acceptance (still one acceptance → one test).

## Must / Must not

- **Must**: Require both versions before analysis
- **Must**: Allow `none found` for Chronic without fabricating history
- **Must not**: Invent Verify from Testing Focus
- **Must not**: Expertise cosplay or RST persona framing
- **Must not**: Replace sc-clarify, lens pass, or mission brief
