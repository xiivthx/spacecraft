# Test oracles

On-demand procedure: evaluate whether an observation is a problem using FEW HICCUPPS oracles. Decision job - not expertise cosplay. **Not** an always-on discuss soft gate. **Not** a teaching tutor.

## Goal

Ground problem judgment: name which oracles fire for a concrete observation, state why it is a problem, and suggest follow-up exploration - without inventing Verify or adopting RST tutor persona.

## Output

Chat summary plus greppable artifact in `decisions.md` (discuss) or equivalent block in review/finding notes (run):

```
## Oracle evaluation
- Observation: …
- Oracles fired: …  # mnemonic letter + short name, only those that apply
- Problem statement: …  # why this is a problem (oracle-grounded)
- Confidence: high | medium | low
- Follow-up ideas: …  # 1-3 exploratory/charter or Test Ideas bucket hints; do not invent Verify
- Open questions: …  # or none
```

Or skip:

```
Oracle evaluation skipped: <reason>
```

## Good / Bad

- Good: concrete observation; only fired oracles listed; problem statement usable by defect-finding `issue`/`impact`; follow-ups map to Test Ideas buckets or Strategy charter ideas; soft gate only
- Bad: expertise cosplay / RST persona; teaching Intro→Story→Why→Tools→Learning format; inventing Verify; always-on gate; listing all 11 oracles every time; essay dumps

## Verify

When run: `decisions.md` has `## Oracle evaluation` OR `Oracle evaluation skipped:` (discuss), or equivalent greppable block in review/finding notes (run). Does **not** block discuss clear. **sc-judge**: oracle gaps alone do **not** justify `REFUTED`.

## When to use

- Human asks "is this a bug?", oracle evaluation, or why something is wrong
- Mid-discuss when Overlooked Test Ideas need "why it's a problem"
- Mid-run / review when a finding's problem judgment needs grounding before filing defect-finding
- Exploring a charter outcome ("we saw X - problem?")

## When skip

- Clear acceptances already machine-checkable
- No observation
- Routine RED-GREEN with explicit expected literals

## Error handling

- **No observation/behavior** → respond exactly: `Please provide a specific observation or behavior (what you saw, where, and what you expected) so I can evaluate oracles.`
- **Observation too vague** → respond exactly: `Observation unclear - please provide more detail about what happened and what seemed wrong.`

## Procedure

1. **Read** - observation (what happened, where, what seemed wrong or unexpected). If missing → exact error phrase above. If vague → exact error phrase above.

2. **Scan FEW HICCUPPS** (silent checklist - list only oracles that fire):

   | Letter | Oracle | Means |
   |--------|--------|-------|
   | F | Familiar problems | Matches known failure patterns / bug families |
   | E | Explainability | Behavior hard to explain / unjustified |
   | W | World | Inconsistent with how the real world works |
   | H | History | Inconsistent with prior versions / past behavior |
   | I | Image | Damages brand / trust / professional image |
   | C | Comparable products | Worse than reasonable comparable products |
   | C | Claims | Contradicts documented claims / specs / marketing / help |
   | U | User expectations | Violates reasonable user wants / mental model |
   | P | Product | Inconsistent with other parts of the same product |
   | P | Purpose | Fails the product's purpose / mission for the user |
   | S | Statutes | Violates laws, standards, contracts, accessibility regs |

3. **Problem statement** - one concise oracle-grounded why (usable for defect-finding `issue` / `impact`).

4. **Confidence** - high | medium | low based on observation clarity and oracle strength.

5. **Follow-up ideas** - 1-3 hints only: exploratory charter angles, or Test Ideas bucket (Positive / Negative / Edge / Overlooked). Optional reflection from observation - not learner homework.

6. **Open questions** - park in `questions.md` when needed, or `none`.

7. **Record** - write `## Oracle evaluation` (or skip) to `decisions.md` or finding notes. Do not invent Verify.

## Heuristics pointer

Oracles answer **"is this a problem?"** Heuristics answer **"where/how to look"** - use existing refs; do not duplicate them here:

- SFDIPOT + quality criteria: `requirement-testability.md`
- Strategy / charters: `htsm-strategy.md`
- Requirement delta: `rcrcrc-impact.md`
- Test data categories: `test-data-design.md`
- SFDIPOT coverage of existing tests: `sfdipot-coverage.md`
- Defect filing after judgment: `sc-run/references/defect-finding.md`

## Clear rule

- **Soft** - does **not** block discuss clear.
- **sc-judge** - oracle evaluation gaps alone do **not** justify `REFUTED`.

## Must / Must not

- **Must**: Ground in a specific observation
- **Must**: List only oracles that fire (letter + short name)
- **Must**: Use exact error strings for missing/vague observations
- **Must not**: Expertise cosplay ("as a testing expert…")
- **Must not**: RST tutor persona or 6-step teaching format
- **Must not**: Invent Verify from oracle evaluation
- **Must not**: Replace `requirement-testability.md`, `htsm-strategy.md`, `rcrcrc-impact.md`, `sfdipot-coverage.md`, `test-data-design.md`, or `defect-finding.md`
- **Must not**: Add always-on discuss gate behavior
- **Must not**: Unicode em dash - ASCII hyphen-minus only

## Related

- Defect findings after judgment: `sc-run/references/defect-finding.md`
- Test Ideas buckets: `requirement-testability.md`
- Charter follow-up: `htsm-strategy.md`
- Inspiration (RST oracles framing): James Bach, Michael Bolton, Maaike Brinkhof - FEW HICCUPPS mnemonic; not required reading
