---
name: sc-verification
description: Capture fresh command evidence before claiming work is complete. Activate after task implementation, verify step, or when evidence is needed.
license: MIT
compatibility: opencode
metadata:
  version: 1.0
  audience: xiivthx
---

# sc-verification

Capture fresh command evidence before claiming work is complete.

## When to use

Activate when the user asks to:

- Verify that a task is complete
- Capture evidence for a mission task
- Run validation checks before claiming done
- Check evidence requirements before a review or ship

## Workflow

Use this exact sequence unless the user specifies otherwise:

1. **Resolve the mission** — Run `scripts/spacecraft resolve --json`. Ensure safety is `safe` before proceeding.
2. **Capture evidence** — Run `scripts/spacecraft evidence "<label>" -- <command>` for each acceptance check.
3. **Validate** — Run `scripts/spacecraft validate` after evidence capture.
4. **Map to acceptance** — Record evidence ids in final summaries, referencing each acceptance check from `plan.json`.

### Edge cases

- **Evidence command fails** — Capture the failure as evidence. Do not skip. Fix the issue and re-capture.
- **Validation fails** — `scripts/spacecraft validate` returns non-zero. Check which acceptance criteria are unmet. Fix before claiming done.
- **Check cannot be automated** — State why in the evidence label. Mark as `manual`. Document the manual verification steps.
- **No plan.json exists** — Cannot map evidence to acceptance checks. Ask user to create a plan first.
- **Evidence already captured for this check** — Re-run to ensure freshness. Never reuse stale evidence.

## Research auto-trigger

When a verification command uses an unfamiliar tool or flag, run `spacecraft research "<tool> <flag> usage"` before capturing evidence. Incorrect verification commands produce misleading evidence.

## Rules

- **Must**: No done/pass/verified/ready claim without evidence.
- **Must**: Resolve the mission with `scripts/spacecraft resolve --json`; `.space/current` is fallback state, not sole authority.
- **Must**: If resolver safety is not `safe`, stop before evidence capture and ask the user to choose with `scripts/spacecraft use <number|id|title>` or an explicit `SPACECRAFT_MISSION`.
- **Must**: Use `scripts/spacecraft evidence "<label>" -- <command>`.
- **Must**: Capture failures too.
- **Must**: Map acceptance checks to evidence ids in final summaries.
- **Must**: If a check cannot be automated, state why and mark it manual.
- **Must**: Prefer focused verification first, then broader build/test checks before shipping.
- **Must**: Evidence must demonstrate functional correctness, not just configuration validity.
  - **Weak**: evidence that echoes the config back (e.g., "PASS: model set to X")
  - **Strong**: evidence that exercises actual behavior (e.g., "PASS: model X produces correct output for test case Y")
  - Prefer functional proof. If only config validation is possible, explicitly state why.
- **Must**: Before claiming verification passed, self-audit: "Did I verify behavior or just read config? Did I cover edge cases?"

## Out of scope

This skill does NOT handle:

- Automated test execution — use the project's test runner instead
- Code review or design critique — use sc-review or sc-design
- Release readiness verification — use sc-review for full closeout checks

## Output format

```
scripts/spacecraft evidence "<label>" -- <command>
```

Evidence entries are appended to `evidence.jsonl` in the mission directory. Each entry is a JSON object with label, command output, timestamp, and status.

## Checklist

Before claiming verification passed:

- [ ] Mission resolved with `scripts/spacecraft resolve --json` (safety = `safe`)
- [ ] Evidence captured for every acceptance check
- [ ] Failures captured as evidence too
- [ ] Validation passed with `scripts/spacecraft validate`
- [ ] Evidence ids mapped to acceptance checks in summary

---

## References

- `scripts/spacecraft evidence --help` — evidence subcommand reference
- `scripts/spacecraft validate --help` — validation reference
- `evidence.jsonl` in the active mission directory
