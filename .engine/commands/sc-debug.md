---
description: Debug and diagnose bugs with five-step discipline
agent: sc-commander
---
Use sc-verification, sc-git, sc-clarify.
Debug: $ARGUMENTS

## Pre-flight checks

If the user is reporting a bug, error, or unexpected behavior, enter debug mode. State the five constraints once in your first response, then apply in order.

## Five steps

State these once, first response:

1. **Reproduce first.** No fix without a reliable repro.
2. **Know the fail path.** Debugger first; then source trace + config parameters; then in-code instrumentation.
3. **Test before accepting.** What would disprove the hypothesis? Run counter-example first.
4. **Every run is a breadcrumb.** Cross-reference all experiments.
5. **Post-mortem or refuse to draft.** Repro, root cause, fix, validation - all four or stop.

---

### 1. Reproduce reliably

Build a runnable repro before anything else.
- **Reliable repro** - capture exact steps, inputs, environment as a runnable artifact (failing test, curl, CLI). Use `sc-verification` for evidence.
- **Flaky repro** - raise the rate first: loop the trigger, run in parallel, add stress. 50% flake is debuggable; 1% is not.
- **No repro at all** - stop. Ask the user for env access, captured artifacts (HAR, log dump, core), or permission to instrument. Do not proceed to guess.

Target: a fast (1-5 s), deterministic pass/fail signal.

**Use `sc-git`** to create an investigation branch from latest `main` before any diagnostic code changes.

---

### 2. Know the fail path

Find *where* the code breaks and *what stops it from breaking*. Escalate only when the prior tactic fails.

1. **Attach a debugger.** Step to the failure site. One breakpoint beats ten logs.
2. **Source trace + config parameters.** Trace the code path end-to-end. List every config parameter: flags, env vars, feature toggles, branch conditions, input shape, timing, concurrency. Flip one at a time.
3. **In-code instrumentation.** `printf` / log statements at the suspected fail site. Tag every probe with a unique prefix (e.g. `[DBG-a4f2]`) so cleanup is a single grep.

---

### 3. Test the hypothesis

Examine a candidate root cause **before** testing it.

- Does it explain the symptom end-to-end?
- What is the cleanest **counter-example**? Run it first.
- Generate 3-5 ranked hypotheses, not one.

---

### 4. Breadcrumb ledger

Maintain a running ledger of every experiment. Each entry: what changed, what happened, what it ruled in or out.

| Field | Description |
|-------|-------------|
| `timestamp` | ISO 8601 |
| `hypothesis` | What was believed before the run |
| `command` | Exact command or action taken |
| `result` | PASS, FAIL, or specific observation |
| `rules_in` | What this run confirms |
| `rules_out` | What this run eliminates |

- When a new hypothesis surfaces, walk the ledger. Does it hold for **every** prior observation?
- If any past run contradicts it, refine or discard.
- Update the ledger after every run.

---

### 5. Post-mortem

**Refuse to draft without ALL four mandatory inputs:**

- [ ] **Repro steps** - deterministic repro exists (runnable)
- [ ] **Root cause** - actual mechanism identified (not a hypothesis)
- [ ] **Fix** - PR / commit / branch pointer
- [ ] **Validation** - original repro now passes

#### Structure

1. **Summary** - What broke, what fixed it, one sentence.
2. **Symptom** - What was observed.
3. **Root cause** - Actual bug mechanism with file paths, commit SHAs.
4. **Why it produced the symptom** - Link root cause to symptom.
5. **Fix** - What changed, why it addresses root cause.
6. **How it was found** - Debugging path, rejected hypotheses, confirming experiment.
7. **Why it slipped through** - CI gap, untested code, review miss. Blameless.
8. **Validation** - Original repro passes, test names, before/after numbers.
9. **Action items** - Concrete next-steps. If none: "None."

Use `sc-git` for commit/PR pointers. Use `sc-verification` for validation evidence.

---

## Anti-rationalization guards

These excuses are always wrong. When you catch yourself thinking one, return to the step you're skipping:

| Excuse | Correct action |
|--------|----------------|
| "This is an obvious one-liner" | Return to step 1. Reproduce first. |
| "I've seen this pattern before" | Treat as hypothesis. Trace (step 2), then test (step 3). |
| "A repro would take too long" | Build the repro. If truly too long, say so and ask for help. |
| "I'll add a ledger later" | Start the ledger now. |
| "The debugger isn't available" | Follow step 2 escalation: trace + config parameters, then instrumentation. |
| "I don't need all four inputs" | All four mandatory. Refuse to draft. |

## Hard stop gates

- No reliable reproduction exists
- Hypothesis formed before fail path traced
- Conclusion accepted before testing attempted
- Post-mortem drafted without all four mandatory inputs
- Ambiguity about bug behavior blocks diagnosis - route to `sc-clarify`

End with post-mortem (if fix validated) or next debug step.
