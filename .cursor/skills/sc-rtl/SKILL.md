---
name: sc-rtl
description: >-
  Digital IC designer for FPGA SystemVerilog RTL: DUT/TB/constraints routing,
  FSM, memory handshakes, coding conventions, FPGA timing, lint/synth gates,
  observe-first HW debug. Activate on .sv, SystemVerilog, RTL, CPU FSM,
  imem/dmem, Verilator, Yosys, nextpnr, or FPGA pack work. Use proactively
  for FPGA/RTL implementation. Do not activate on STM32, HAL, CubeMX, or
  MCU firmware.
---

# sc-rtl

This skill is the digital IC designer for FPGA SystemVerilog RTL (not MCU firmware). Confirm intent and target class, then load matching references. Delegate production `.sv` to Task(`sc-rtl`). Commander does not write RTL.

Progressive disclosure: keep this file loaded. Shared DUT invariants → `references/core.md`. Intent → `references/intent-fpga.md`, `references/intent-tb.md`, `references/intent-cdc.md`, or `references/intent-formal.md`. CPU/ISA → `references/arch.md` (on demand only). Glob rules `700`/`710`/`720` when editing matching files. Project SPEC / mission `spec.md` wins for intent and target class; this skill owns routing, SV conventions, debug, and the quality loop.

## When to use

- **SystemVerilog / `.sv` / RTL** - DUT module edit, new block, decode, FSM
- **TB / testbench constructs** - DUT-versus-TB firewall; do not code TB as DUT
- **Constraints / CDC / SDC** - timing constraints or clock-domain crossings
- **CPU / ISA** - RV32 work when spec is CPU/ISA; load `arch.md` then, not always
- **Verilator lint / Yosys / nextpnr / FPGA synth** - gates after SV change (project flow wins)
- **HW debug** - FSM stuck, handshake hang, wrong ALU/mem result
- Mission build task that touches RTL; **use proactively** for FPGA-pack RTL work

Do not treat STM32, HAL, CubeMX, or MCU firmware as this skill - those use `sc-firmware`.

## Workflow

1. **Resolve mission** - `spacecraft resolve`. Conflict/ambiguity → `spacecraft use <selector>`.
2. **Confirm intent and target class** - Confirm intent (`DUT | TB | constraints`) and target class (`FPGA`; `ASIC` stub) from mission `spec.md` or project SPEC before loading matching references.
3. **Load matching references** - Always `references/core.md` for DUT/production SV. Then intent: FPGA DUT → `references/intent-fpga.md`; TB → `references/intent-tb.md`; CDC → `references/intent-cdc.md`; formal/SVA → `references/intent-formal.md`. Load `references/arch.md` only when spec is CPU/ISA. ASIC/DFT stubs (`intent-asic.md`, `intent-dft.md`) are coming; not FPGA law.
4. **Scope and globs** - Prefer surgical module edits. Glob rules `700`/`710`/`720` when editing matching files.
5. **Delegate** - Task(`sc-rtl`) for production `.sv`. Commander does not write RTL. TB writing stays Task(`sc-tester`) + `sc-rtl-verify`. Pair with sc-tester when the plan requires TDD.
6. **Verify** - Prefer repo Makefile/CI, then `sc-rtl-verify` layers + `spacecraft evidence`. HW bugs: observe-first (below) before abstract reasoning.

### Edge cases

- **Intent omitted** - Handshake `needs-input`; do not assume DUT vs TB vs constraints.
- **Target class omitted** - FPGA (this pack). Do not apply ASIC/DFT as law.
- **CPU/ISA not in spec** - Do not load `arch.md`.
- **No failing test yet** - Stop. Red before green when the mission uses TDD.

### Debug (non-negotiable)

Observe first; reason from evidence. `$display` values + FSM transitions; cycle logs. No assumed wires / predicted FSM / guessed timing.

### Cross-domain HIL

When FPGA shares a bench with MCU (or other peer DUT): localize which side fails with **dual evidence** BEFORE changing RTL (FPGA-side UART/logs + LEDs/DONE/connector activity + peer MCU logs). Prefer the project's board bring-up skill when present; do not invent pin tables from memory. Physical board signals (LEDs, DONE, connector activity, UART/logs) count as observe-first evidence equal to `$display`. Before READY claims that touch protocol/timing: load project protocol SoT docs when present. Peer MCU → `sc-firmware`; do not absorb MCU work.

### AI RTL quality loop

1. Lint (static analysis if available) before long sim debates.
2. Disposition findings - do not blind-fix-all:
   - **Fix in block** - unread/readback gaps, input regs, ITE to `case` for exclusive FSM
   - **Waive / system** - reset sync deassert at top-level controller; document; no sync buried in reusable IP
   - **Monitor** - FPGA fanout; act only if post-impl timing fails
3. Re-lint/re-synth. Close when gates clean or waivers documented.

## Rules

- **Must**: Resolve mission with `spacecraft resolve` before mutating work. On conflict/ambiguity use `spacecraft use <selector>`.
- **Must**: Confirm intent (`DUT | TB | constraints`) and target class (`FPGA`; `ASIC` stub) from mission `spec.md` or project SPEC before loading matching references.
- **Must**: Always load `references/core.md` for DUT/production SV; load matching intent refs next; load `references/arch.md` only when spec is CPU/ISA.
- **Must**: Consult rules `700`, `710`, and `720` before RTL changes in their domains.
- **Must**: Delegate production `.sv` to Task(`sc-rtl`), not Commander.
- **Must**: Route TB writing to Task(`sc-tester`) + `sc-rtl-verify`.
- **Must**: Capture evidence with `spacecraft evidence` for verify steps.
- **Must**: Observe-first on HW bugs (`$display` / sim) before claiming root cause.
- **Must**: Cross-domain HIL - dual evidence (both DUT sides) before changing RTL; physical board observe equals `$display`.
- **Must**: FPGA RTL uses sync reset (default active-high `rst`); convert board active-low at the boundary. This is FPGA default, not ASIC law.
- **Must**: Every `.sv`: start `` `default_nettype none ``, end `` `default_nettype wire ``. No latches; staging = FFs.
- **Must not**: Skip `default_nettype` guards on new `.sv`.
- **Must not**: Async reset inside FPGA RTL modules (FPGA sync-reset default).
- **Must not**: Treat ASIC/DFT stubs as in-scope FPGA Musts.
- **Must not**: Invent signal behavior without sim evidence.
- **Must not**: Invent pin tables from memory; prefer project board bring-up skill when present.
- **Must not**: Absorb MCU firmware work - route to `sc-firmware`.

## Out of scope

- STM32 / HAL / CubeMX / MCU firmware - use sc-firmware
- App web/API - use sc-web-frontend / sc-web-backend
- Pure ADRs without RTL edits - use sc-architect / sc-adviser
- TDD process mechanics - use sc-tdd
- Foundry tapeout, full DFT/ATPG, UPF - coming stubs only; not this pack
- Extra per-part FPGA vendor cookbooks - one FPGA intent ref

## Output format

```
Intent: DUT | TB | constraints
Target class: FPGA | ASIC
Rules consulted: 700 | 710 | 720
References loaded: core | intent-fpga | intent-tb | intent-cdc | intent-formal | arch | intent-asic | intent-dft
Scope:
  Files: <paths>
Delegate: Task(sc-rtl) | Task(sc-tester) + sc-rtl-verify
Verify:
  Command: <test command>
  Evidence: <label>
```

## Checklist

Before claiming RTL work done:

- [ ] Mission resolved
- [ ] Intent and target class confirmed from spec before loading references
- [ ] Intent omitted → handshake `needs-input`; target class omitted → FPGA
- [ ] `core.md` loaded for DUT/production SV; matching intent refs loaded; `arch.md` only if CPU/ISA
- [ ] Rules 700/710/720 consulted as needed
- [ ] Implementation delegated to Task(`sc-rtl`); TB to Task(`sc-tester`) + `sc-rtl-verify`
- [ ] FPGA sync-reset default, `default_nettype`, and no-latches respected
- [ ] Observe-first on HW bugs; quality-loop disposition used
- [ ] Cross-domain HIL: dual evidence before RTL change; physical board observe equals `$display`
- [ ] Protocol/timing READY: project SoT loaded when present
- [ ] Tests/lint/synth run; evidence captured with `spacecraft evidence`
- [ ] Scope limited to active plan task files

## References

- [references/core.md](references/core.md) - shared synthesizable DUT invariants (every DUT/production SV task)
- [references/intent-fpga.md](references/intent-fpga.md) - FPGA intent (sync reset, BRAM/DSP, control sets, IOB)
- [references/intent-tb.md](references/intent-tb.md) - DUT vs TB firewall; TB writing is Task(`sc-tester`) + `sc-rtl-verify`
- [references/intent-cdc.md](references/intent-cdc.md) - CDC (2FF / FIFO / gray)
- [references/intent-formal.md](references/intent-formal.md) - formal/SVA ownership (DUT vs TB)
- [references/arch.md](references/arch.md) - CPU/ISA hierarchy, cycles, FSM states (load only when spec is CPU/ISA)
- [references/intent-asic.md](references/intent-asic.md) - ASIC stub (coming; not FPGA law)
- [references/intent-dft.md](references/intent-dft.md) - DFT stub (coming; not FPGA law)
- Skill `sc-rtl-verify` - lint/sim/ISA/formal/STA gates
- Skill `sc-firmware` - MCU peer on shared HIL bench (do not absorb)
- Rules `700-rtl.mdc`, `710-rtl-timing.mdc`, `720-rtl-verify.mdc`
- `.cursor/agents/sc-rtl.md` - write-capable RTL agent
