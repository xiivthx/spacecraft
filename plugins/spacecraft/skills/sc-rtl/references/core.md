# Synthesizable DUT core

On-demand shared DUT core: synthesizable subset, sequential versus combo, FPGA house reset, and coding recipes. Critical Must/Must-not stay in `.cursor/rules/700-rtl.mdc`. Load this file for every DUT/production SystemVerilog task. Intent refs load from SKILL routing; this file is the shared core.

## Synthesizable subset

DUT code is the synthesizable subset. Simulator-only constructs do not belong in DUT RTL. Cummings class: simulator-only facts versus synth-only facts cause silicon mismatch (`sim≠synth`). Write DUT so simulation and synthesis agree.

TB-only constructs (`$display`, `#` delays, `initial`) and the DUT-versus-TB firewall live in `intent-tb.md`.

## Sequential vs combo

- Sequential flops: `always_ff @(posedge clk)` with reset inside the block; nonblocking (`<=`) for flops
- Combinational: `always_comb`
- no latches: full assign every branch of `always_comb`; staging is flops
- Prefer `case` / `unique case` for mutually exclusive FSM/decode; nested `if` only when priority is real

## Reset

FPGA house default is synchronous reset: default active-high `rst`; convert board active-low at the boundary. This is FPGA default, not ASIC law. Detail lives in rule 700 and `intent-fpga.md`.

## Coding patterns

DUT-wide. Vendor primitives, BRAM/DSP, and IOB live in `intent-fpga.md`. CPU/ISA tables live in `arch.md` (load only when spec is CPU/ISA).

**Names.** `snake_case`. Purpose prefixes (`imem_`, `dmem_`) when those ports exist.

**FSM.** `case` for exclusive states.

**Timing.** Register at natural Fmax boundaries; multi-cycle large cones; document unregistered exemptions (handshake `*_ready` OK).

## Debug

Observe first; reason from evidence. `$display` is TB-side. Full DUT-versus-TB firewall is `intent-tb.md`.

## Related

- Rule: `.cursor/rules/700-rtl.mdc` - Must / Must-not invariants
- FPGA intent: `intent-fpga.md`
- TB firewall: `intent-tb.md`
- CDC: `intent-cdc.md`
- Formal: `intent-formal.md`
- CPU/ISA: `arch.md` (load only when spec is CPU/ISA)
