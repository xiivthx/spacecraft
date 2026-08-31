# FPGA intent

On-demand FPGA-first DUT recipes: synchronous reset, BRAM/DSP inference, control sets, vendor primitives, and IOB. Critical Must/Must-not stay in `.cursor/rules/700-rtl.mdc`. Synthesizable subset, sequential versus combo, and no-latches live in `core.md`. Load this file for FPGA DUT. One FPGA-intent file; no per-part vendor cookbooks and no board SKUs.

## Reset

FPGA RTL uses synchronous reset (house default in 700): default active-high `rst`; convert board active-low at the boundary. Reset inside `always_ff @(posedge clk)` with `if (rst)`. Do not present async-assert as FPGA law. ASIC reset difference lives in `intent-asic.md` (coming stub); it is not this pack's FPGA Must.

## BRAM, DSP, and control sets

Write memories and arithmetic so synthesis infers BRAM and DSP (AMD UG901/UG949 class): a synthesizable array plus address/`we` for RAM; registered multiply-accumulate for DSP. Keep control sets small - share clock-enable and reset across related flops. Unique per-flop enables fragment packing.

## Vendor primitives and IOB

Infer first. Instantiating vendor primitives is the exception, used when SPEC names the block or inference cannot meet a documented resource or timing need. Register at the I/O boundary and pack IOB flops for pad timing. Pad timing constraints live in SDC (rule 710), not as the only source of truth in RTL comments.

## Related

- Shared DUT core: `core.md` - synthesizable subset
- Glob Musts: `.cursor/rules/700-rtl.mdc`
- ASIC reset (coming stub): `intent-asic.md`
