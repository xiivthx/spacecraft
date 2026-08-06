# HTSM strategy pass (slim)

`/sc-discuss` soft gate: risk-driven test strategy snapshot for high-risk missions. Decision job - not QA persona / expertise cosplay. Grounded in HTSM (Project Environment, Product Elements/SFDIPOT, Quality Criteria, Techniques) but **slim default** - not the full 10-section Strategy Template.

## Goal

Produce a concise, risk-driven test strategy snapshot so planning and RED-GREEN cover what matters without blind gaps or generic "test everything."

## Output

Chat summary plus greppable artifact in `decisions.md`:

```
## Strategy pass
- Mission & success: …
- Scope in: … | Scope out: …
- Top risks (L x I): …
- SFDIPOT highlights: …
- Quality → risk → coverage: …
- Charter ideas (3–7): …
- Assumptions: …  # or none
- Open questions: …  # park detail in questions.md; ask via sc-clarify one at a time
```

Or skip:

```
Strategy pass skipped: <reason>
```

## Good / Bad

- Good: concrete charter-style ideas; top risks with L×I; SFDIPOT/quality highlights tied to coverage; assumptions explicit; open questions parked
- Bad: expertise cosplay; always-on full 10-section essay; vague "test everything"; inventing Verify from charters; replacing testability, RCRCRC, lens, or sc-clarify

## Verify

`decisions.md` has `## Strategy pass` OR `Strategy pass skipped:`. Strategy incompleteness does **not** block clear the way `Not Testable` + soft Verify does. Greppable markers only - no CLI validator required.

## When required (any one)

- Greenfield / new product surface
- Multi-platform / multi-device / multi-browser matrix matters
- Security, PII, compliance, or authz sensitivity
- Critical external integrations or stated SLOs
- Human asks for test strategy / HTSM / risk-driven strategy

## When skip

- Routine bug with clear Verify
- Single-path feature, low blast radius, no triggers above
- `/sc-quick`
- Testability skipped and nothing else raises risk (still prefer explicit skip reason)

Record skip in `decisions.md`:

```
Strategy pass skipped: <reason>
```

## Procedure (slim)

1. **Read** - `spec.md` + `decisions.md` (incl. `## Testability pass` when present). If application context missing/trivial → respond ONLY: `Application context missing or insufficient - please provide: product summary, target users, key workflows, platforms, key constraints, dates, and critical integrations.` Keep clarify open; record skip or leave open per discuss norms.

2. **Partial info** - proceed with explicit Assumptions; park Open Questions in `questions.md`; ask only high-impact via sc-clarify (one at a time).

3. **Fill slim template** - prefer concrete charter-style ideas over generics. Top risks use L x I (Likelihood x Impact).

4. **Techniques vocabulary** - may appear inside Charter ideas / Quality coverage (exploratory charters, boundary, state/flow, API checks, property-based, perf, security probes, a11y, compat, chaos, data integrity, observability) - do **not** require a Techniques-by-Area matrix by default.

5. **Full Strategy Template** (10 sections: Project Environment detail, Techniques matrix, Environments strategy, Risk Register owners, Reporting/Exit, 12-step implementation) **only if human explicitly asks** - then write optional mission artifact `.space/missions/<id>/test-strategy.md` and still keep a slim `## Strategy pass` summary in `decisions.md`.

## Clear rule

- Soft gate: missing Strategy pass must record skip before clear when discuss runs; Strategy incompleteness does **not** block clear the way `Not Testable` + soft Verify does.
- Do not invent Verify from charters.

## Must / Must not

- **Must**: record `## Strategy pass` OR `Strategy pass skipped:` before discuss clear
- **Must not**: expertise cosplay; always-on full 10-section essay; replace testability, RCRCRC, lens, or sc-clarify; invent Verify
