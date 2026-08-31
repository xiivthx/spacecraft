# DUT versus TB firewall

On-demand DUT-versus-TB firewall. Simulator-only constructs do not belong in synthesizable DUT (Cummings class: simulator-only facts versus synth-only facts). Synthesizable subset lives in `core.md`. Load this file when intent is TB or when a DUT edit might leak testbench constructs. This is a firewall, not a UVM cookbook.

## TB-only (forbidden in DUT)

`# delays`, `initial`, and `$display` are TB-only constructs forbidden in DUT coding.

- `# delays` stay in the testbench; DUT must not wait on simulation time
- DUT must not use `initial` for reset or stimulus
- `$display` stays TB-side (observe-first prints, not synthesizable RTL)

Delays and `$display` stay out of synthesizable DUT.

## Who writes TB

TB writing stays Task(`sc-tester`) with `sc-rtl-verify`. Do not write testbenches as DUT work. Lint, self-checking sim, ISA, formal, and STA gates belong to `sc-rtl-verify`.

## Related

- Shared DUT core: `core.md` - synthesizable subset
- Skill `sc-rtl-verify` - verification gates
