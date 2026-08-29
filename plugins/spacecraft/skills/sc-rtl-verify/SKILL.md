---
name: sc-rtl-verify
description: >-
  FPGA/RISC-V RTL verification plan and gates: lint, self-checking sim, ISA/ACT,
  formal (riscv-formal/RVFI), coverage, STA evidence. Activate on testbench,
  cocotb, Verilator regression, formal, riscv-arch-test, timing signoff, or
  when claiming RTL/CPU ready. Use proactively with fpga pack verify work.
---

# sc-rtl-verify

Verification discipline for SystemVerilog / FPGA / RISC-V. Complements `sc-rtl` (implement). TB via Task(`sc-tester`); RTL fixes via Task(`sc-rtl`/`sc-coder`).

## When to use

TB / cocotb / formal / ISA regression; before ready/ship when RTL/CPU in scope; after RTL change needing evidence beyond elaborate; proactively on fpga-pack verify.
## Workflow

1. **Resolve mission** - `spacecraft resolve` / `use`.
2. **Read plan** - acceptance in `spec.md` / `plan.json`; note ISA subset and FPGA target.
3. **Pick layers** (run cheap first; skip only with documented reason):

| Layer | Gate | Typical cmd / artifact |
|-------|------|------------------------|
| L0 Lint | Verilator lint clean (waivers documented) | `verilator --lint-only …` or `make lint` |
| L1 Struct | Yosys elaborates; no unintended latch/multi-driver | `yosys -p '…; proc; check'` or `make synth` |
| L2 Sim | Self-checking regression exit 0 | `make sim` / Verilator binary / cocotb |
| L3 ISA | SPEC subset vs Spike/Sail or ACT ELFs | project `make isa` / arch-test runner |
| L4 Formal | Critical props or riscv-formal when RVFI exists | `sby` / `make formal` |
| L5 Impl | PnR meets Fmax + resource budget | `make` in `rtl/fpga` / nextpnr report |

4. **Evidence** - each done acceptance: `spacecraft evidence "<label>" -- <cmd>`.
5. **Disposition** lint/CDC/static hits: **fix-in-block** | **waive-system** (e.g. reset sync at top) | **monitor** (FPGA fanout until timing fails). Re-run gate after fix.
6. **Ready bar** - do not claim CPU/RTL ready on L0 alone. Minimum for production-minded FPGA core: L0+L1+L2 green; L3 when CPU ISA in scope; L4 when formal tooling in repo; L5 when FPGA target in scope.

### Observe-first (HW / sim bugs)

`$display` / waves / UART before theorizing FSM. Same contract as `sc-rtl` debug section.

### RISC-V notes

- **riscv-arch-test (ACT):** minimal filter for ISA compatibility - not a substitute for directed/random DV ([riscv-arch-test](https://github.com/riscv/riscv-arch-test)).
- **riscv-formal:** optional RVFI port + SymbiYosys checkers when enabling formal ([YosysHQ/riscv-formal](https://github.com/YosysHQ/riscv-formal)).
- Match Sail/Spike config to DUT extensions (SPEC may be RV32IMC only).

### Coverage

Functional/cover properties beat raw line %. Record plan + measured numbers or `Coverage skipped: <reason>`.

## Must / Must not

- **Must**: Prefer Makefile/CI; evidence every verify acceptance
- **Must**: Self-checking TB (assert/scoreboard/signature) - no passive pass
- **Must**: Consult rules `700` / `710` / `720` when editing matching globs
- **Must not**: Sign off on lint-only; invent EDA installs mid-mission
- **Must not**: Treat ACT pass as "fully verified CPU"

## Out of scope

- Writing production RTL - `sc-rtl`
- STM32 firmware - `sc-firmware`
- Foundry tapeout DRC/LVS decks - record out-of-scope until ASIC freeze
- TDD ceremony - `sc-tdd` (still use RED TB before GREEN RTL when plan says TDD)

## Checklist

- [ ] Layers selected vs SPEC (L0–L5)
- [ ] L0+L1+L2 evidence present for RTL claims
- [ ] L3 if CPU ISA in scope (or skip line)
- [ ] L4/L5 if tools/target in scope (or skip line)
- [ ] Waivers greppable; timing report if FPGA

## References

- [references/signoff.md](references/signoff.md) - compact signoff checklist + tool map
- Rules: `700-rtl.mdc`, `710-rtl-timing.mdc`, `720-rtl-verify.mdc`
- Skill: `sc-rtl` (implement)
