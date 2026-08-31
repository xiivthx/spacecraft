# Formal / SVA intent

On-demand DUT-versus-TB ownership for formal and SVA. Load this file when adding assertions, properties, or formal collateral. The verification gate is `sc-rtl-verify` L4. Do not merge verify into `sc-rtl`.

## DUT-versus-TB ownership

**DUT.** May carry synthesizable assertions that do not break synthesis:

- Immediate `assert` on protocol or illegal states when the synth/sim flow treats them as synthesizable (or strips them with a documented `SYNTHESIS` guard)
- Concurrent `assert property` / SVA only when the project flow supports that construct in DUT. If the construct is sim-only or breaks synth, it belongs in TB, not DUT

DUT assertions are design intent. They must not introduce `$display`, `#` delays, `initial`, or other TB-only constructs (`intent-tb.md`).

**TB.** Owns stimulus and non-synthesizable verification collateral:

- Stimulus, `$display`, `$fatal` in checkers, `#` delays, `initial`
- Cover groups and cover properties used as coverage (not DUT)
- Bind files that attach SVA to DUT without editing DUT source
- riscv-formal wrappers / RVFI checkers
- SymbiYosys `sby` scripts, formal assumptions, and cut-points

Write TB via Task(`sc-tester`) with `sc-rtl-verify`. Formal run evidence is L4 (`sby` / `make formal`).

## Related

- TB firewall: `intent-tb.md`
- Gate: skill `sc-rtl-verify` L4
- Rule: `.cursor/rules/720-rtl-verify.mdc`
