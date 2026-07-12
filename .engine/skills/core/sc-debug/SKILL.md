---
name: sc-debug
description: Five-step debugging discipline. Activate on /sc-debug, errors, stack traces, bugs, or "debug/diagnose/investigate" requests. Reproduce, trace, falsify, cross-reference, post-mortem.
license: MIT
compatibility: opencode
metadata:
  version: 1.0
  audience: xiivthx
---

# sc-debug

Five-step debugging discipline for root-cause analysis. Engineer-level rigor, scientist-level falsification. Recite verbatim, apply in order, refuse to ship without evidence. Follows the skill format and cross-references existing sc-* skills without duplicating them.

## When to use

Activate when the user asks to:

- Debug or diagnose a bug or unexpected behavior
- Investigate a failure, error, stack trace, or broken test
- Find the root cause of an issue
- Document a fix with a structured post-mortem after debugging

## Workflow

### Recite the mantra — verbatim, first thing in first response

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
- **Reliable repro** → capture exact steps, inputs, environment as a runnable artifact (failing test, curl, CLI). Use `sc-verification` for evidence.
- **Flaky repro** → the bug is not yet debuggable. Raise the rate first: loop the trigger, parallelise, add stress, narrow timing windows, inject sleeps. 50% flake is debuggable; 1% is not.
- **No repro at all** → stop. Say so explicitly. Ask the user for env access, captured artifacts (HAR, log dump, core), or permission to instrument. Do **not** proceed to hypothesise.

Target: a fast (1–5 s), deterministic pass/fail signal. Pin time, seed the RNG, freeze network, isolate filesystem.

**Use `sc-git`** to create an investigation branch from latest `main` before any diagnostic code changes.

---

### 2. Know the fail path

Once reproducible, find *where* the code breaks and *what stops it from breaking*. Try in this order — escalate only when the prior tactic fails.

1. **Attach a debugger.** If the environment supports it, attach and step to the failure site. One breakpoint beats ten logs. Do this **before** turning any knobs.
2. **Source trace + knob enumeration.** If no debugger, trace the code path end-to-end and list every knob that can influence the outcome: config flags, env vars, feature toggles, branch conditions, input shape, timing, concurrency. Flip one at a time.
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

| Field | Description |
|-------|-------------|
| `timestamp` | ISO 8601 of the experiment |
| `hypothesis` | What was believed before the run |
| `command` | Exact command or action taken |
| `output_hash` | SHA-256 of output or short summary |
| `result` | PASS, FAIL, or specific observation |
| `rules_in` | What this run confirms |
| `rules_out` | What this run eliminates |

**Operating rules for the ledger:**

- When a new hypothesis surfaces, walk the ledger. Does it hold for **every** prior observation, not just the most recent?
- If any past run contradicts it, the hypothesis is wrong or incomplete — refine or discard.
- When in doubt, design the **single experiment** whose outcome makes it certain. Run that next.
- Update the ledger after every run. It is your memory across the session.
- The completed ledger is raw material for the post-mortem in step 5.

---

### 5. Post-mortem

After the fix lands and is validated, draft the canonical engineering record. For other engineers and future-you who will have forgotten in 6 months.

**Refuse to draft without ALL four mandatory inputs.** If any are missing, list what's missing and stop:

- [ ] **Repro steps** — deterministic repro exists (runnable, not just a description)
- [ ] **Root cause** — the actual mechanism is identified (not a hypothesis)
- [ ] **Fix** — PR / commit / branch pointer to the actual change
- [ ] **Validation** — the original repro now passes; the fix is confirmed

If you came in via steps 1–4, the breadcrumb ledger from step 4 is your raw material — pull from it.

#### Post-mortem structure

1. **Summary** *(mandatory)* — What broke, what fixed it, one sentence.
2. **Symptom** — What was observed (errors, logs, test failures).
3. **Root cause** *(mandatory)* — Actual bug mechanism with code identifiers, file paths, commit SHAs.
4. **Why it produced the symptom** — Link root cause to symptom.
5. **Fix** *(mandatory)* — What changed, why it addresses root cause. Link to PR/commit.
6. **How it was found** — Debugging path: repro, tools, rejected hypotheses, confirming experiment.
7. **Why it slipped through** — CI gap, latent code, review miss. Blameless.
8. **Validation** *(mandatory)* — Original repro passes, specific test names, before/after numbers.
9. **Action items / follow-ups** — Concrete next-steps. If none: "None — the fix is sufficient."

Use `sc-git` to reference commit/PR pointers in the post-mortem. Use `sc-verification` to capture validation evidence.

---

## Rules

### General

- **Must**: Recite the mantra block **once** per debug session, in your first response. Do not re-recite mid-session.
- **Must**: Recite **verbatim**. Never paraphrase, shorten, or skip lines of the recital.
- **Must**: Apply the five steps **in order**. The step order is enforced — no skipping ahead.
- **Must**: If the user says "skip the mantra" → skip the recital but still apply the five steps silently.
- **Must not**: Propose a fix before step 1 is satisfied (reliable repro exists).
- **Must not**: Form a hypothesis before step 2 has narrowed the fail path (no hypothesis before trace).
- **Must not**: Accept a conclusion before step 3 has tried to disprove the hypothesis (no conclusion before falsification).
- **Must not**: Commit to a hypothesis before step 4 confirms it against every prior breadcrumb.
- **Must not**: Draft a post-mortem before step 5's four mandatory inputs are all satisfied.
- **Must**: No fix before reproduction. If you catch yourself proposing a fix without a reliable repro, stop and return to step 1.
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

## Out of scope

This skill does NOT handle:

- Issue tracking, bug databases, CI/test runner integration
- Incident management or outage timelines
- Tool-based debugging (dap CLI, debug adapters) — this is pure workflow discipline
- Evidence capture — use `sc-verification`
- Git branch/commit hygiene — use `sc-git`
- Ambiguity resolution — use `sc-clarify`
- Code review — use sc-review
- Visual/UI design — use `sc-design`

## Output format

Session start: recite the mantra block (see Workflow §Recite the mantra). Maintain breadcrumb ledger throughout session (format below). Draft post-mortem after fix validation (mandatory sections: Summary, Root cause, Fix, Validation).

## Checklist

- [ ] Reliable repro exists and is captured as evidence (step 1)
- [ ] Fail path traced through source code or instrumentation (step 2)
- [ ] Hypothesis attempted to be falsified before acceptance (step 3)
- [ ] Breadcrumb ledger has entries for every experiment (step 4)
- [ ] Post-mortem has all four mandatory inputs (step 5)
- [ ] No fix before repro; no hypothesis before trace; no conclusion before falsification
- [ ] Cross-references to sc-verification, sc-git, sc-clarify used where appropriate
- [ ] Fix validated with sc-verification evidence before post-mortem

## Research auto-trigger

When encountering an unknown error message, unexpected framework behavior, or configuration issue during debugging, invoke `spacecraft research <query>` to search for known solutions before forming hypotheses.

---

## References

- `sc-verification` — evidence capture and validation
- `sc-git` — branch hygiene and commit references
- `sc-clarify` — ambiguity resolution during diagnosis
