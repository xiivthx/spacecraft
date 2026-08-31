# Software pack

App, API, CLI, tests, in-process races, "worked yesterday" regressions. Load only after `Pack: software`.

## Repro

Failing test, curl script, CLI invocation, or replay harness. Target a fast (1-5 s) deterministic signal. Pin time, seed RNG, freeze network when those axes matter.

Flaky: raise the rate (loop, parallel, stress, tighter windows). ~50% flake is debuggable; ~1% is not. Do not proceed on a single unreproducible hit.

## Observe (escalate only when the prior fails)

1. **Debugger** - Attach and step to the failure. Skip when halt would hide a race (then treat as Heisenbug: go to knobs / logs / Cursor Debug Mode).
2. **Source trace + knobs** - Walk the path. List flags, env, feature toggles, input shape, timing, build options. Flip **one** at a time.
3. **Instrumentation** - printf / logs with a unique prefix (e.g. `[DBG-a4f2]`) so cleanup is one grep. Or ask Thai HIL: switch Cursor Debug Mode (Shift+Tab) when Agent cannot reach runtime state. Debug Mode = hypothesise, log to Cursor debug server, human repro, read logs, targeted fix, strip probes. Capability stays Cursor; this pack still requires a repro artifact first.

## Isolate

- Regression vs `main` → `git bisect` (delta on history)
- Large failing input → minimize until each remaining piece is necessary
- Passing vs failing run → the delta is the candidate cause, not the latest stack frame alone

## Falsify / ledger

Same spine as the parent skill. Do not promote a cause that fails any prior run.

## Fix

RED characterization (Task `sc-tester`) then GREEN (Task `sc-coder`). Coder does not edit tests. `spacecraft evidence` when a mission is active. Strip `[DBG-…]` probes before `/sc-quick`.

## Out

Unfamiliar API / stale docs while stuck → sc-search (does not replace this pack). Spec unclear → `/sc-discuss`.
