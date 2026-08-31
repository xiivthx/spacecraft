---
name: sc-debug
description: "RCA session opener: reproduce, observe the fail path, falsify, ledger, then fix. Invoke as /sc-debug. Activate on bug hunt, diagnose, stack trace, broken/failing/throwing, or unknown root cause. Not for obvious one-line fixes (/sc-quick) or spec holes (/sc-discuss)."
disable-model-invocation: true
---

# sc-debug

## Goal

Open a debug session, pick **one pack**, prove root cause with a repro artifact, then hand off. Does not merge, tag, or ship. Small proven fixes → `/sc-quick`. Spec holes / multi-module / design-contract change → `/sc-discuss`. Mid `/sc-run` → apply this spine silently; do not start a second lifecycle.

## Output

```
Pack: software | hardware-mcu | hardware-fpga | visual
Repro: <artifact path or command> | blocked: no-repro
Fail path: <where it breaks>
Cause: <one sentence> | open
Exit: /sc-quick | /sc-discuss | stay-run | blocked: <reason>
```

Handshake: `done` | `blocked: <reason>` | `needs-input: <question>`. Thai HIL for Debug Mode switch and exit lane.

## Good / Bad

- Good: pack before tactics; reproducible signal before hypotheses; disproof first; ledger every run; characterization stays; Task-delegate the fix; probes tagged and removed
- Bad: recite a mantra at the user; propose a fix with no repro; software debugger-first on MCU timing; Cursor Debug Mode on RTL/visual look; pixel-diff as RCA; Commander product edits; merge/tag from this skill

## Verify

Greppable `Pack:` + `Repro:` (or `blocked: no-repro`) + `Exit:`. No ship. Characterization or evidence command exists when a fix landed. Probes gone or tagged for cleanup.

## When to use

- `/sc-debug` session open
- User reports a bug, paste of stack/log, "broken / throwing / failing", diagnose / investigate / root cause
- `/sc-run` fix loop when the failure cause is unknown (silent spine; 3-cycle still hands back)

Skip: obvious one-file fix with known cause (`/sc-quick`); requirement ambiguity (`/sc-discuss`); taste-only UI (`/sc-discuss` + designer).

## Workflow

1. **Pack first** - Classify before any tactic. Load **one** ref. Mixed look+behavior → split: behavior = software, look = visual.
   - **software** - app / API / CLI / tests / process race / git regression
   - **hardware-mcu** - MCU firmware, IRQ, DMA, peripherals, HIL
   - **hardware-fpga** - SystemVerilog, sim, ILA, timing at speed
   - **visual** - rendered look vs approved draft / DESIGN.md; functional click-path is software
2. **Reproduce** - Runnable artifact (test, curl, CLI, sim, HIL, screenshot pair). Flake: raise rate before debug. No repro → `blocked: no-repro`; ask env, dump, or instrument permission. Do not hypothesise.
3. **Simplify** - Drop irrelevant input / steps / commits (`git bisect` on regressions). Flip one knob at a time.
4. **Observe fail path** - Follow the loaded pack ladder. Escalate only when the prior tactic cannot reach the bug.
5. **Falsify** - 3-5 ranked hypotheses. Walk each end-to-end. Run the **disproof** first. Ledger: what changed, what happened, ruled in/out. A hypothesis that contradicts any prior breadcrumb is wrong or incomplete.
6. **Fix** - After cause holds against the ledger: characterization/RED via Task(`sc-tester`) when tests apply; impl via Task(`sc-coder` / `sc-firmware` / `sc-rtl` / `sc-designer`). Commander does not write product code or tests. TWINS after defect fixes. 3 failed fix-verify → human.
7. **Exit** - Remove probes. Then:
   - small obvious fix, no mission needed → `/sc-quick`
   - spec / scope / look-contract hole → `/sc-discuss`
   - already in `/sc-run` → stay-run (evidence + continue)
   - no repro / no bench / no oracle → `blocked`

Apply steps 2-6 **in order**. Catch a proposed fix without a repro → return to step 2. Carry the spine silently - do not quote it back.

Cursor Debug Mode (Shift+Tab) is a **software-pack** runtime-log tactic only. Commander cannot switch that mode; ask in Thai when the software ladder needs it.

## Must / Must not

- **Must**: Detect pack before tactics; stop without a repro; disproof before committing to a cause; Task-delegate fixes; 3-cycle handback; INTENT before behavior-changing edits
- **Must not**: Merge/tag/push; recite the spine; use software halt/printf-first on timing HW; use Cursor Debug Mode for MCU/FPGA/visual look; treat screenshot diff as root cause; thaw design-contract oracles

## References

- [software.md](references/software.md) - tests, debugger, knobs, Cursor Debug Mode, bisect
- [hardware.md](references/hardware.md) - MCU observe-without-halt; FPGA sim/waves/ILA
- [visual.md](references/visual.md) - look vs behavior; draft-parity; computed style
- `/sc-quick`, `/sc-discuss`, `/sc-run`, sc-tdd, sc-verification, sc-rtl, sc-firmware, sc-ux-design, sc-browser-probe
