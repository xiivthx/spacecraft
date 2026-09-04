# RTL verify signoff (load on demand)

## Pre-ready checklist (FPGA / RISC-V IP)

Adapt thresholds to SPEC; do not invent foundry signoff.

- [ ] Lint clean (Verilator/Verible as available) + documented waivers
- [ ] Synthesis/elab clean (Yosys `check` / vendor synth): no unintended latches, multi-drivers, undriven outputs
- [ ] Self-checking sim regression green (assert/scoreboard/signature)
- [ ] System/integration sim + software image when in scope (software owns DUT bring-up) before board HIL
- [ ] Coverage vs written plan (or skip line) - line% alone insufficient
- [ ] CDC: sync structures + constraints reviewed (`710`; `intent-cdc.md`)
- [ ] ISA filter when CPU: ACT and/or Spike/Sail compare for SPEC subset
- [ ] Formal on critical props / riscv-formal when enabled (`intent-formal.md`; L4)
- [ ] PnR timing meets Fmax; resource within budget
- [ ] Physical HIL evidence when FPGA target in scope; scenario↔observable map when project mapping doc exists
- [ ] Each gate captured in `evidence.jsonl` via `spacecraft evidence`

## Tool map (open-source default)

| Need | Tools |
|------|--------|
| Lint | Verilator `--lint-only`, Verible |
| Sim | Verilator, Icarus, cocotb |
| Synth / struct | Yosys |
| Formal | SymbiYosys (`sby`), riscv-formal |
| PnR / timing | nextpnr, icetime / vendor STA |
| ISA | Spike, Sail, riscv-arch-test |

## Evidence label examples

```
spacecraft evidence "rtl-lint" -- make lint
spacecraft evidence "rtl-sim-smoke" -- make sim
spacecraft evidence "rtl-isa-rv32i" -- make isa-rv32i
spacecraft evidence "rtl-formal-rvfi" -- make formal
spacecraft evidence "fpga-timing" -- make -C rtl/fpga
```

## sc-rtl references

Implement recipes (not this gate): `core.md`, `intent-fpga.md`, `intent-tb.md`, `intent-cdc.md`, `intent-formal.md`. CPU/ISA: `arch.md`.
