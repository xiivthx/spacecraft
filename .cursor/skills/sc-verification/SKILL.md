---
name: sc-verification
description: "Capture fresh command evidence before claiming work is complete. Activate after task implementation, verify step, or when evidence is needed."
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

1. **Resolve the mission** - Run `spacecraft resolve`. On conflict or ambiguity, use `spacecraft use <selector>`.
2. **Capture evidence** - Run `spacecraft evidence "<label>" -- <command>` for each acceptance check.
3. **Validate** - Run `spacecraft validate` after evidence capture. Prefer `spacecraft validate --strict` before claiming a build task or mission build is complete (`--strict` requires `exitCode` on every evidence entry, ≥1 evidence entry, and matching evidence for each done plan task).
4. **Map to acceptance** - Record evidence ids in final summaries, referencing each acceptance check from `plan.json`.

### Edge cases

- **Evidence command fails** - Capture the failure as evidence. Do not skip. Fix the issue and re-capture by appending a new evidence entry with the same or a clearer label. To discard a bad line, delete that line from `evidence.jsonl` manually (there is no overwrite flag).
- **Validation fails** - `spacecraft validate` (or `validate --strict`) returns non-zero. Check which acceptance criteria are unmet. Fix before claiming done.
- **Check cannot be automated** - State why in the evidence label. Mark as `manual`. Document the manual verification steps.
- **No plan.json exists** - Cannot map evidence to acceptance checks. Ask user to create a plan first.
- **Evidence already captured for this check** - Re-run and append fresh evidence. Never reuse stale evidence.

## Rules

- **Must**: No done/pass/verified/ready claim without evidence.
- **Must**: Resolve the mission with `spacecraft resolve`; `.space/current` is fallback state, not sole authority. On conflict/ambiguity use `spacecraft use <selector>`.
- **Must**: Use `spacecraft evidence "<label>" -- <command>`.
- **Must**: Capture failures too.
- **Must**: Map acceptance checks to evidence ids in final summaries.
- **Must**: If a check cannot be automated, state why and mark it manual.
- **Must**: Prefer focused verification first, then broader build/test checks before shipping.
- **Must**: Evidence must demonstrate functional correctness, not just configuration validity.
  - **Weak**: evidence that echoes the config back (e.g., "PASS: model set to X")
  - **Strong**: evidence that exercises actual behavior (e.g., "PASS: model X produces correct output for test case Y")
  - Prefer functional proof. If only config validation is possible, explicitly state why.
- **Must**: Before claiming verification passed, self-audit: "Did I verify behavior or just read config? Did I cover edge cases?"
- **Must**: After defect fixes, require `TWINS:` - project-wide search for the same construct / twin occurrences before claiming done.
- **Must**: After **3 failed fix-verify cycles**, stop and hand back to human. Do not keep looping.

## Out of scope

This skill does NOT handle:

- Automated test execution - use the project's test runner instead
- Code review or design critique - use sc-review or sc-design
- Release readiness verification - use sc-review for full closeout checks

## Output format

```
spacecraft evidence "<label>" -- <command>
```

Evidence entries are appended to `evidence.jsonl` in the mission directory. Each entry is a JSON object with label, command output, timestamp, and status.

## Checklist

Before claiming verification passed:

- [ ] Mission resolved with `spacecraft resolve` (on conflict/ambiguity: `spacecraft use <selector>`)
- [ ] Evidence captured for every acceptance check
- [ ] Failures captured as evidence too
- [ ] Validation passed with `spacecraft validate --strict` before claiming build complete
- [ ] Evidence ids mapped to acceptance checks in summary

---

## References

- `spacecraft evidence --help` - evidence subcommand reference
- `spacecraft validate --help` - validation reference (`--strict` for ship/build claims)
- `evidence.jsonl` in the active mission directory
