# RTL verify signoff (load on demand)

## Pre-ready checklist (FPGA / RISC-V IP)

Adapt thresholds to SPEC; do not invent foundry signoff.

- [ ] Lint clean (Verilator/Verible as available) + documented waivers
- [ ] Synthesis/elab clean (Yosys `check` / vendor synth): no unintended latches, multi-drivers, undriven outputs
- [ ] Self-checking sim regression green (assert/scoreboard/signature)
- [ ] Coverage vs written plan (or skip line) - line% alone insufficient
- [ ] CDC: sync structures + constraints reviewed (`710`)
- [ ] ISA filter when CPU: ACT and/or Spike/Sail compare for SPEC subset
- [ ] Formal on critical props / riscv-formal when enabled
- [ ] PnR timing meets Fmax; resource within budget
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
