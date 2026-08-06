# Test data design

On-demand procedure: generate structured test **data** (variable-level values) that complements scenario-level Test Ideas. Decision job - not expertise cosplay. **Not** an always-on discuss soft gate.

## Goal

Produce diverse, non-repetitive test data rows across categories so planner/tester fixtures cover positive, negative, boundary, exploratory, and security-shaped inputs without inventing Verify.

## Output

Chat summary plus greppable artifact in `decisions.md`:

```
## Test data design
- Context / variables: …
- Distribution target: ~N rows (default 10 when human asks for a table; fewer when focused)
- Rows:
  | Row | Variable | Value | Category | Notes |
  | … | … | … | Positive|Negative|Boundary|Exploratory|Security | … |
- Bucket map: Positive→Positive; Negative→Negative; Boundary→Edge; Exploratory→Overlooked; Security→Negative or Overlooked
- Open questions: …  # or none; park detail in questions.md
```

Or skip:

```
Test data design skipped: <reason>
```

## Good / Bad

- Good: concise diverse values; clarify when variables unclear; category coverage; usable by sc-tester fixtures / sc-planner acceptance when present; maps to Test Ideas buckets
- Bad: expertise cosplay; inventing Verify; always-on gate; replacing testability/strategy/SFDIPOT; repetitive near-duplicates; dumping many clarifying Qs in one turn

## Verify

When run: `decisions.md` has `## Test data design` OR `Test data design skipped:`. Does **not** block discuss clear. Greppable only - no CLI validator required.

## When to use

- Human asks for test data / fixtures / input matrix / penetration-shaped samples for named variables
- Testability Data/Interfaces risks need concrete values before RED
- sc-tester/planner need diverse fixtures for a multi-variable acceptance
- Mid-run when building parameterized tests and variables are known

## When skip

- No variables/context
- Routine single happy-path
- Discuss clear not blocked by missing this

## Error handling

- Variables/context unclear → ask clarifying questions (1-3) before generating
- Exact phrase when nothing usable: `Variables or context unclear - please provide variable names and what the system does with them.`

## Procedure

1. **Read** - context (`spec.md`, variables list, or human paste). If unclear → ask or error phrase above.

2. **Distribute** - across categories: Positive, Negative, Boundary, Exploratory/Creative, Security/Penetration (aim for coverage of all five when generating a full table; allow focus when human scopes).

3. **Inspiration checklist** (use when relevant - do not force every domain):
   - **Paths/Files** - long names, special chars, non-existent paths
   - **Time/Date** - leap days, invalid dates, time zones, formats
   - **Numbers** - zero, max/min, scientific notation, negatives, decimals
   - **Strings** - very long, accents, CJK, delimiters, SQL/script injection, emojis

4. **Concise diverse examples** - avoid unnecessary repetition. Prefer fewer sharp rows over many near-duplicates.

5. **Record** - write `## Test data design` (or skip) to `decisions.md`. Do not invent Verify from rows.

6. **Optional** - feed noteworthy rows into Test Ideas (Edge/Negative/Overlooked) if a testability pass is also being written - scenarios stay in Test Ideas; values stay here.

## Clear rule

- **Soft** - does **not** block discuss clear.
- **sc-judge** - data-table gaps alone do **not** justify `REFUTED`.

## Must / Must not

- **Must**: Ground in named variables
- **Must not**: Expertise cosplay ("as a QA expert…")
- **Must not**: Invent Verify from data rows
- **Must not**: Replace `requirement-testability.md`, `htsm-strategy.md`, or `sfdipot-coverage.md`
- **Must not**: Add always-on discuss gate behavior
