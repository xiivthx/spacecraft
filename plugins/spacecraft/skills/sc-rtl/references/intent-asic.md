# ASIC intent (stub)

Foundry-class ASIC target is out of scope for this FPGA pack. This file is a coming stub, not a foundry recipe. Load it only to bound the gap. A coming pack or later mission owns ASIC flow; do not treat this pack as tapeout-ready.

## Reset

Do not treat FPGA synchronous reset as ASIC law. ASIC often uses async-assert / sync-deassert. That async-assert policy differs from the FPGA default in `.cursor/rules/700-rtl.mdc`. Do not apply FPGA sync-reset as a universal digital-IC Must.

## Not in this pack

No UPF low-power, no tapeout DRC/LVS, no foundry decks. Those stay out of this pack.

## Related

- FPGA reset (house default): `intent-fpga.md`
- Glob Musts: `.cursor/rules/700-rtl.mdc`
