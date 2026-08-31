# Hardware pack

Observer effect is the law. Halt, blocking printf, and heavy probes change timing and can hide the bug. Load MCU **or** FPGA after pack detect - never apply the software debugger-first ladder here.

No board, probe, dump, or sim license → `blocked: no-repro`. Do not invent waveforms.

Domain depth: `sc-firmware` / `sc-rtl` / `sc-rtl-verify`. This file only orders RCA.

---

## hardware-mcu

MCU firmware, IRQ, DMA, peripherals, HIL. Delegate fixes to Task(`sc-firmware`).

### Repro

Host unit test for logic. Target/HIL for pins, IRQ, DMA, clocks. Record stimulus (input, timing, power). ISR bugs need a trigger that is not "set a breakpoint in the ISR".

### Observe (timing-sensitive: this order)

1. **Truth on the wire** - scope, logic analyzer, GPIO toggle. Least disturbance to the CPU path.
2. **Non-halt trace** - SWO / ITM / RTT / ring buffer. No printf in ISR. No blocking UART in a 1 kHz loop.
3. **Halt GDB / OpenOCD** - last, for hardfault and logic that is not timing-sensitive. Halt masks races.
4. **UART printf** - last resort. Blocking I/O is milliseconds; it creates Heisenbugs.

Logic-only (no timing): host test + GDB is fine. Cache/DMA/ISR are separate hypotheses from "algorithm wrong". Cortex-M7 DMA: clean/invalidate D-cache is a first-class knob.

### Isolate

Reproduce on host to drop silicon. If only on target: electrical vs firmware vs ISR vs DMA vs clock. One axis per experiment.

### Fix

Task(`sc-firmware`). Keep ISR short. Strip SWO/printf probes or leave a documented trace plan. Evidence: host or HIL command via `spacecraft evidence` when a mission is active.

---

## hardware-fpga

SystemVerilog DUT, sim, on-chip debug, "fails at speed". Delegate RTL to Task(`sc-rtl`); TB to Task(`sc-tester`) + `sc-rtl-verify`. Observe-first: no guessed wires or predicted FSM.

### Repro

Self-checking sim first (assert / scoreboard / signature). Silicon-only / at-speed: ILA/SignalTap at the **native** clock with a trigger that matches the failure.

### Observe (cheap → expensive)

1. Lint / elaborate (unintended latch, undriven net)
2. Sim waves at the handshake / FSM boundary (`$display` tagged; unique prefix)
3. Formal / SVA when the property already exists in-repo
4. On-chip ILA - few probes, shallow depth, same clock domain as the signals. Timing closure of the probe is part of the evidence.

Functional fail in sim ≠ timing fail at Fmax. CDC is its own hypothesis. STA fail is not a logic bug until the path is identified.

### Isolate

Unit the block in sim before the full SoC. If only on board: I/O voltage, clock, reset sync, CDC, then DUT logic.

### Fix

Task(`sc-rtl`). Do not code TB as DUT. Re-run the sim (and ILA if silicon) after the fix. Remove or gate ILA for ship unless the design keeps a debug mux on purpose (document).

## Must not (both)

- Software pack halt/printf-first on a timing bug
- Cursor Debug Mode as the HW observer
- Claim root cause from RTL reading without sim or capture
