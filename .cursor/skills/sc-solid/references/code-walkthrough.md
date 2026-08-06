# Code walkthrough

On-demand procedure when reviewing a pasted snippet, explaining code to a non-developer, or deepening a code-quality pass beyond the silent SOLID scan. Decision job - not expertise cosplay.

## Goal

Give a structured, concise read of a code snippet: compliance, risks, and optional plain-language explanation - without duplicating sc-solid reference docs or reciting theory on every edit.

## Output

Compact headings matching the seven steps below. **Beginner Explanation** last - only when triggered (human asks, or mission-brief Feynman Answers).

```
## Compliance
…

## Consistency
…

## Security
…

## Maintainability
…

## Redundant / unreachable
…

## Test ideas
…

## Beginner explanation
…  # only when triggered
```

## Good / Bad

- Good: ordered checks; pointers to existing refs instead of duplicating; concise bullets; exact error strings for bad input; Beginner Explanation only when asked
- Bad: expertise cosplay; essay-length output; Beginner Explanation on every silent SOLID scan; duplicating `solid-principles.md` / `clean-code.md` verbatim; inventing behavior not in the snippet

## Verify

Output uses the compact template (≤7 sections). Missing/invalid snippet returns exact fixed strings below. Commander or human confirms headings match procedure and refs are pointers not copies.

## When to use

- Human pastes a snippet to analyze or explain
- Human asks for beginner-friendly walkthrough
- Deepening a code-quality pass beyond routine silent SOLID on a focused snippet

## When skip

- Routine silent SOLID scan on every edit (sc-solid main workflow - no walkthrough output)
- Diff review already covered by `sc-reviewer` / defect findings
- Snippet absent or unusable (return error string; do not proceed)

## Error handling

- **No snippet** → respond exactly: `No code snippet was found. Please provide one so I can analyze it.`
- **Incomplete/invalid snippet** → respond exactly: `The code seems incomplete/invalid. Could you provide the full snippet?`

## Procedure (ordered)

1. **Compliance** - Does the snippet meet project rules and acceptance intent? Check `.cursor/rules/` globs on touched paths; mission `spec.md` when mission-scoped. Note gaps only - do not restate full rule text.

2. **Consistency** - Naming, patterns, and style vs surrounding codebase. See `references/clean-code.md` and `references/code-smell.md` for detection - cite violations, not theory.

3. **Security** - Obvious risks in the snippet (injection, secrets, unsafe deserialization, auth gaps). See `.cursor/rules/300-security.mdc` - pointer only; flag concrete lines.

4. **Maintainability** - SOLID / complexity / architecture fit. See `references/solid-principles.md`, `references/complexity.md`, `references/architecture.md` - flag SRP/DIP and accidental complexity.

5. **Redundant / unreachable** - Dead code, duplicate logic, unreachable branches. See `references/code-smell.md` (dead code, duplication).

6. **Test ideas** - 2-5 concrete checks mapped to buckets: **Positive** (happy path), **Negative** (invalid input / failure), **Edge** (boundaries). Align with sc-tdd discipline when writing tests - do not invent Verify bars.

7. **Beginner explanation** - **Only when triggered**: step-by-step plain language (what each part does, data flow, outcome). Skip on silent SOLID scans.

## Must / Must not

- **Must**: Keep each step 2-5 lines max in output
- **Must**: Point to existing refs; do not duplicate long excerpts
- **Must**: Use exact error strings for missing/invalid snippets
- **Must not**: Expertise cosplay ("as a senior engineer…")
- **Must not**: Beginner Explanation on every routine scan
- **Must not**: Unicode em dash - ASCII hyphen-minus only

## Related

- Silent SOLID workflow: `sc-solid/SKILL.md` (On every code change)
- Defect findings after review: `sc-run/references/defect-finding.md`
- Test buckets: `requirement-testability.md` (discuss), sc-tdd
