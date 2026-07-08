---
name: sc-debug
description: Five-step debugging discipline — reproduce, trace the fail path, falsify the hypothesis, cross-reference every breadcrumb, document the post-mortem. Recite the mantra block verbatim at the start of any debugging session, then apply the five steps in order before proposing any fix. Trigger on /sc-debug and proactively whenever debugging starts — user reports a bug, says something is broken/throwing/failing, asks to debug/diagnose/investigate an issue, or pastes a stack trace or error log.
license: MIT
compatibility: opencode
metadata:
  version: 1.0
  audience: xiivthx
---

# sc-debug

Five-step debugging discipline for root-cause analysis. Engineer-level rigor, scientist-level falsification. Recite verbatim, apply in order, refuse to ship without evidence. Follows the Spacecraft skill format and cross-references existing sc-* skills without duplicating them.

## When to use

Activate when the user asks to:

- Debug or diagnose a bug or unexpected behavior
- Investigate a failure, error, stack trace, or broken test
- Find the root cause of an issue
- Document a fix with a structured post-mortem after debugging

## Workflow

Use this exact sequence unless the user specifies otherwise:

### Recite the mantra — verbatim, as the first thing in your first response

> **Mantra:**
> 1. **First is reproducibility.** Can the issue be reproduced reliably?
> 2. **Know the fail path.** Debugger first; then source trace + knob enumeration; then in-code instrumentation.
> 3. **Question your hypothesis.** What would disprove it?
> 4. **Every run is a breadcrumb.** Cross-reference all of them.
> 5. **Document the post-mortem.** Repro, root cause, fix, validation — all four, or refuse to draft.

Then begin work.

---

### 1. Reproduce reliably

Build a runnable repro before anything else.

- **Reliable repro** → capture the exact steps, inputs, and environment as a runnable artifact: failing test, curl script, CLI invocation, replay harness. Use `sc-verification` to capture repro evidence.
- **Flaky repro** → the bug is not yet debuggable. Raise the rate first: loop the trigger, parallelise, add stress, narrow timing windows, inject sleeps. 50% flake is debuggable; 1% is not.
- **No repro at all** → stop. Say so explicitly. Ask the user for env access, captured artifacts (HAR, log dump, core), or permission to instrument. Do **not** proceed to hypothesise.

Target: a fast (1–5 s), deterministic pass/fail signal. Pin time, seed the RNG, freeze network, isolate filesystem.

**Use `sc-git`** to create an investigation branch from latest `main` before any diagnostic code changes.

---

### 2. Know the fail path

Once reproducible, find *where* the code breaks and *what stops it from breaking*. The differential narrows the search. Try in this order — escalate only when the prior tactic fails.

1. **Attach a debugger.** If the environment supports it, attach and step to the failure site. One breakpoint beats ten logs. Do this **before** turning any knobs.
2. **Source trace + knob enumeration.** If no debugger (or it can't reach the bug), trace the code path end-to-end and list every knob that can influence the outcome:
   - config flags, env vars, feature toggles
   - branch conditions, input shape
   - timing, concurrency, build options
   Each knob is a candidate axis to flip in the differential. Flip one at a time.
3. **In-code instrumentation.** If outside knobs can't move the failure, go inside: `printf` / log statements at the suspected fail site, dump the relevant internal state. Tag every probe with a unique prefix (e.g. `[DBG-a4f2]`) so cleanup is a single grep. Let the trace show where reality diverges from your model.

---

### 3. Falsify the hypothesis

When a candidate root cause surfaces, scrutinise it **before** testing it.

- Does it actually explain the symptom end-to-end? Walk it through.
- What is the simplest **proof**? What is the cleanest **disproof**?
- Run the **disproof first**. If the hypothesis survives, it's real. If it dies, you saved yourself from chasing a phantom.
- Generate 3–5 ranked hypotheses, not one. Single-hypothesis thinking anchors on the first plausible idea.

---

### 4. Every run is a breadcrumb

Maintain a running **ledger** of every experiment in this session. Each entry: what changed, what happened, what it ruled in or out.

#### Breadcrumb ledger format

Each entry must contain these fields:

| Field | Description |
|-------|-------------|
| `timestamp` | ISO 8601 timestamp of the experiment |
| `hypothesis` | What you believed before the run — one sentence |
| `command` | Exact command or action taken |
| `output_hash` | SHA-256 of the output (for reproducibility) or short summary |
| `result` | What happened: PASS, FAIL, or specific observation |
| `rules_in` | What this run confirms (list) |
| `rules_out` | What this run eliminates (list) |

**Operating rules for the ledger:**

- When a new hypothesis surfaces, walk the ledger. Does it hold for **every** prior observation, not just the most recent?
- If any past run contradicts it, the hypothesis is wrong or incomplete — refine or discard.
- When in doubt, design the **single experiment** whose outcome makes it certain. Run that next, instead of churning on adjacent runs.
- Update the ledger after every run. It is your memory across the session.
- The completed ledger is the raw material for the post-mortem in step 5.

---

### 5. Post-mortem

After the fix lands and is validated, draft the canonical engineering record. Written **for** other engineers (and future-you, who will have forgotten everything in 6 months). Code identifiers are first-class here.

**Refuse to draft without ALL four mandatory inputs.** If any are missing, list what's missing and stop:

- [ ] **Repro steps** — deterministic repro exists (runnable, not just a description)
- [ ] **Root cause** — the actual mechanism is identified (not a hypothesis)
- [ ] **Fix** — PR / commit / branch pointer to the actual change
- [ ] **Validation** — the original repro now passes; the fix is confirmed

If you came in via steps 1–4, the breadcrumb ledger from step 4 is your raw material — pull from it.

#### Post-mortem structure

1. **Summary** *(mandatory)* — What broke, what fixed it, in one sentence. PR number, owner.
2. **Symptom** — What was actually observed. Error messages, test failures, log lines.
3. **Root cause** *(mandatory)* — The actual bug mechanism. Code identifiers, file paths, function names, commit SHAs of the offending change. Walk the cause chain end-to-end.
4. **Why it produced the symptom** — Link the root cause to the symptom. Often non-obvious.
5. **Fix** *(mandatory)* — What changed and why it addresses the root cause. Link to PR / commit. Name prior fix attempts and what was wrong with them.
6. **How it was found** — The debugging path: what repro worked, tools used, hypotheses rejected (from breadcrumb ledger), the single experiment that confirmed the cause.
7. **Why it slipped through** — CI gap, latent code, workload gap, incomplete prior fix, review miss. Blameless — describe the gap, not the person.
8. **Validation** *(mandatory)* — How we know the fix works. Original repro now passes, specific test names, concrete before/after numbers. State coverage honestly.
9. **Action items / follow-ups** — Concrete next-steps with owners and tracking artifacts. If none, write "None — the fix is sufficient."

Use `sc-git` to reference commit/PR pointers in the post-mortem. Use `sc-verification` to capture validation evidence.

---

## Rules

### General

- **Must**: Recite the mantra block **once** per debug session, in your first response. Do not re-recite mid-session.
- **Must**: Recite **verbatim**. Never paraphrase, shorten, or skip lines of the recital.
- **Must**: Apply the five steps **in order**. The step order is enforced — no skipping ahead.
- **Must**: If the user says "skip the mantra" → skip the recital but still apply the five steps silently.
- **Must not**: Propose a fix before step 1 is satisfied (reliable repro exists).
- **Must not**: Start testing hypotheses before step 2 has narrowed the fail path.
- **Must not**: Commit to a hypothesis before step 3 has tried to disprove it.
- **Must not**: Declare a hypothesis correct until step 4 confirms it against every prior breadcrumb.
- **Must not**: Draft a post-mortem before step 5's four mandatory inputs are all satisfied.
- **Must**: If you catch yourself proposing a fix without a reliable repro, stop and return to step 1.
- **Must**: If ambiguity about bug behavior blocks diagnosis, route to `sc-clarify` before continuing.
- **Must**: Use `sc-verification` for capturing reproducible evidence (step 1) and fix validation (step 5). Do not reimplement evidence capture.
- **Must**: Use `sc-git` for creating investigation branches before diagnostic code changes and for recording commit/PR references in the post-mortem. Do not reimplement git hygiene.
- **Must**: Use `sc-clarify` when user intent, bug behavior, or expected behavior is materially ambiguous. Do not reimplement clarification workflow.
- **Must**: The mantra is a constraint **you** carry through the session — not advice to deliver back to the user.

### Anti-rationalization guards

Agents often try to skip debug steps with plausible-sounding excuses. These excuses are always wrong. When you hear yourself thinking one of these, stop and return to the step you're trying to skip:

| Excuse | Why it's wrong | Correct action |
|--------|----------------|----------------|
| "This is an obvious one-liner — I can see the fix" | You're fixing what you *think* is broken, not what's actually broken. Many "obvious" fixes paper over symptoms while leaving the root cause untouched. | Return to step 1. Reproduce first, then confirm the fix addresses the root cause, not just the symptom. |
| "I've seen this pattern before — I know what's wrong" | Pattern-matching is hypothesis generation, not diagnosis. The same symptom can have different root causes. You're skipping steps 2 and 3. | Treat pattern recognition as a hypothesis. Trace the fail path (step 2), then falsify (step 3). |
| "A repro would take too long — let me just try this fix" | Without a repro, you can't validate the fix. You're trading 5 minutes of repro-building for potentially hours of blind debugging. | Build the repro. If it truly takes too long, say so explicitly and ask for help isolating. |
| "I'll add a breadcrumb ledger later — let me just explore first" | You're losing data. Every unrecorded run is wasted effort that can't be cross-referenced. | Start the ledger now. Even a partial entry is better than none. |
| "The debugger isn't available, so I'll skip to hypothesis" | Knob enumeration (step 2.2) and in-code instrumentation (step 2.3) are available when debuggers aren't. | Follow the step 2 escalation order: debugger → source trace + knobs → instrumentation. |
| "It's a post-mortem for myself — I don't need all four inputs" | A post-mortem without root cause is a symptom log. Without repro, the next person can't verify. Without validation, it's a guess. | All four inputs are mandatory. Refuse to draft. |

## Operating rules

- Recite the mantra block **once** per debug session, in your first response. Do not re-recite mid-session.
- Recite **verbatim**. Never paraphrase, shorten, or skip lines of the recital.
- If the user says "skip the mantra" → skip the recital but still apply the five steps silently.
- Apply the five steps **in order**:
  - Do not propose a fix before step 1 is satisfied (reliable repro exists).
  - Do not start testing hypotheses before step 2 has narrowed the fail path.
  - Do not commit to a hypothesis before step 3 has tried to disprove it.
  - Do not declare a hypothesis correct until step 4 confirms it against every prior breadcrumb.
  - Do not draft a post-mortem before step 5's four mandatory inputs are all satisfied.
- If you catch yourself proposing a fix without a reliable repro, stop and return to step 1.
- The mantra is a constraint **you** carry through the session — not advice to deliver back to the user.

## Out of scope

This skill does NOT handle:

- General-purpose issue tracking or bug databases
- CI/test runner integration or automated test generation
- Customer-visible incident management or outage timelines — those need a separate incident report
- Tool-based debugging (dap CLI, debug adapters) — this is pure workflow discipline, not tool integration
- Evidence capture — use `sc-verification`
- Git branch/commit hygiene — use `sc-git`
- User ambiguity resolution — use `sc-clarify`
- Code review — use `sc-reviewer`
- Visual/UI design — use `sc-design`

## Output format

Session start (first response only):
```
> **Mantra:**
> 1. **First is reproducibility.** Can the issue be reproduced reliably?
> 2. **Know the fail path.** Debugger first; then source trace + knob enumeration; then in-code instrumentation.
> 3. **Question your hypothesis.** What would disprove it?
> 4. **Every run is a breadcrumb.** Cross-reference all of them.
> 5. **Document the post-mortem.** Repro, root cause, fix, validation — all four, or refuse to draft.
```

Breadcrumb ledger (maintained throughout session):
```
| timestamp | hypothesis | command | output_hash | result | rules_in | rules_out |
|-----------|------------|---------|-------------|--------|----------|-----------|
| ...       | ...        | ...     | ...         | ...    | ...      | ...       |
```

Post-mortem (after fix lands):
```
## Post-mortem: <issue title>

### Summary
...

### Root cause
...

### Fix
...

### Validation
...
```

## Checklist

Before claiming the debug session is complete:

- [ ] Reliable repro exists and is captured as evidence (step 1)
- [ ] Fail path is traced through source code or instrumentation (step 2)
- [ ] Hypothesis attempted to be falsified before acceptance (step 3)
- [ ] Breadcrumb ledger has entries for every experiment, cross-referenced (step 4)
- [ ] Post-mortem has all four mandatory inputs before drafting (step 5)
- [ ] No fix attempted before repro was built
- [ ] No hypothesis formed before fail path was traced
- [ ] No conclusion reached before falsification was attempted
- [ ] Cross-references to sc-verification, sc-git, sc-clarify used where appropriate
- [ ] Git investigation branch created via sc-git
- [ ] Fix validated with sc-verification evidence before post-mortem

---

## References

- `sc-verification` — evidence capture and validation
- `sc-git` — branch hygiene and commit references
- `sc-clarify` — ambiguity resolution during diagnosis
