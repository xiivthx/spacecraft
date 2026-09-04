---
name: sc-verification
description: "Capture fresh command evidence before claiming work is complete. Activate after task implementation, verify step, or when evidence is needed."
---

# sc-verification

Capture fresh command evidence before claiming work is complete.

## When to use

Task/mission verify; evidence capture; validation before done/review/ship claims.

## Workflow

1. **Resolve** - `spacecraft resolve`; conflict → `spacecraft use <selector>`.
2. **Capture** - `spacecraft evidence "<label>" -- <command>` per acceptance check.
3. **Validate** - `spacecraft validate` after capture (not-doc-drift / not-10X-validate). Prefer `validate --strict` before build-complete claims (`--strict`: `exitCode` on every entry, ≥1 entry, matching evidence per done plan task).
4. **Map** - Evidence ids in summaries → each `plan.json` acceptance.

### Edge cases

- **Evidence command fails** - Capture the failure as evidence. Do not skip. Fix the issue and re-capture by appending a new evidence entry with the same or a clearer label. To discard a bad line, delete that line from `evidence.jsonl` manually (there is no overwrite flag).
- **Validation fails** - `spacecraft validate` (or `validate --strict`) returns non-zero. Check which acceptance criteria are unmet. Fix before claiming done.
- **Check cannot be automated** - State why in the evidence label. Mark as `manual`. Document the manual verification steps.
- **No plan.json exists** - Cannot map evidence to acceptance checks. Ask user to create a plan first.
- **Evidence already captured for this check** - Re-run and append fresh evidence. Never reuse stale evidence.
- **Large output** - Oversized capture may truncate the JSONL `output` (marker `\n...[truncated]`) and write full raw under mission `evidence-raw/` (`outputTruncated`, `outputBytes`, `outputRawPath`). `outputHash` is SHA-256 of the full raw; the terminal still prints the full output. See `docs/mission-artifacts.md`.

## Rules

- **Must**: No done/pass/verified/ready claim without evidence.
- **Must**: Resolve the mission with `spacecraft resolve`; `.space/current` is fallback state, not sole authority. On conflict/ambiguity use `spacecraft use <selector>`.
- **Must**: Use `spacecraft evidence "<label>" -- <command>`.
- **Must**: Capture failures too.
- **Must**: Map acceptance checks to evidence ids in final summaries.
- **Must**: If a check cannot be automated, state why and mark it manual.
- **Must**: Prefer focused verification first, then broader build/test checks before shipping.
- **Must** (mission product path before ready): static-analysis, diff-coverage, mutation, and PBT disposition - evidence labels (`static-…`, `diff-cov-…`, `mutation-…`, `pbt-…`) **or** greppable skip/waive in `decisions.md` (`docs/mission-artifacts.md`).
  - Static (when project tool runs): **Full package/project static suite required (0 warning / 0 error)** (else skip/waive). Tip-path-only lint/typecheck MUST NOT satisfy static-analysis for ready.
  - Diff coverage: touched executable **line and branch ≥90%** (sanity band 90-95%); never global 95-100%.
  - Mutation: in scope if `Mutation: required`, pack id `quality`, or `Mutation: high-risk`; then **>80%** scoped (or project higher bar); else `Mutation skipped: not in scope`.
  - PBT: **100%** of design-contract **core-logic** modules need invariants + generators (project `fast-check` / Hypothesis / equivalent) with `pbt-…`, or `Pbt skipped: no project pbt tool` / `Pbt skipped: not core logic` / `Pbt waived: <reason>`. **Must not** invent PBT lib install mid-mission.
- **Must**: Evidence must demonstrate functional correctness, not just configuration validity.
  - **Weak**: evidence that echoes the config back (e.g., "PASS: model set to X")
  - **Strong**: evidence that exercises actual behavior (e.g., "PASS: model X produces correct output for test case Y")
  - Prefer functional proof. If only config validation is possible, explicitly state why.
- **Must**: When claiming UI/workflow/user-visible behavior, `verify` and/or acceptance text Must include a product-surface marker among `verify.product` | `browser` | `curl` | `composition`; unit-only verify is insufficient.
- **Must**: Before claiming verification passed, self-audit: "Did I verify behavior or just read config? Did I cover edge cases?"
- **Must**: After defect fixes, require `TWINS:` - project-wide search for the same construct / twin occurrences before claiming done.
- **Must**: After **3 failed fix-verify cycles**, stop and hand back to human. Do not keep looping.

## Out of scope

Test runner execution · code/design review (Task(`sc-reviewer`) / sc-ux-design / Task(`sc-designer`)) · full closeout readiness (Task(`sc-reviewer`))

## Output format

```
spacecraft evidence "<label>" -- <command>
```

Appends to mission `evidence.jsonl` (label, output, timestamp, status).

## Checklist

Resolved · evidence per acceptance (incl. failures) · `validate --strict` before build-complete · ids mapped in summary. Full Musts above.

## References

- `spacecraft evidence --help` / `spacecraft validate --help` (`--strict` for ship/build claims)
- `docs/mission-artifacts.md` - evidence schema; outcome-gate skip/waive SoT
