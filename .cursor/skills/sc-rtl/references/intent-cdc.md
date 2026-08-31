# CDC intent (clock-domain crossing)

On-demand CDC recipes: 2FF, FIFO, and gray-code crossings. Constraint Musts stay in `.cursor/rules/710-rtl-timing.mdc`. Load this file when intent is constraints or the DUT crosses async clocks. Formal/SVA ownership is `intent-formal.md`. Timing/STA gate is `sc-rtl-verify` L5.

## Single-bit: 2FF

Use a 2FF synchronizer for a single-bit async signal into a destination clock:

- Two (or more) flops in the dest clock; the first flop samples the async input
- Mark the chain (`ASYNC_REG` / vendor keep) so tools do not merge or retime the chain away
- Source must be a level, or a pulse stretched to at least one dest-clock period; a one-cycle pulse in a faster source clock can vanish

Do not treat 2FF as a bus synchronizer.

## Multi-bit: FIFO / gray

Bare 2FF on a multi-bit bus is illegal: bits can go metastable independently and produce combinations that never existed in the source.

Use:

- **gray** pointers (or a gray-coded counter) when sampling a multi-bit count/pointer with 2FF: adjacent values differ by one bit, so a metastable sample is still a valid adjacent pointer
- Dual-clock **FIFO** (gray write/read pointers plus dual-clock RAM) for data buses and most handshake payloads
- Valid-ready with one control bit 2FF-synced and data held stable until the handshake completes (not a second 2FF on the data bus)

## Reset at the domain boundary

Sync deassert into each clock domain at the system boundary (reset controller). Do not bury per-domain reset sync inside reusable IP. Pair with rule 710 Reset.

## Constraints pairing

Pair every async crossing with a constraint policy (rule 710): `set_clock_groups -asynchronous` or scoped `set_false_path` / `set_max_delay -datapath_only`. Do not leave CDC paths fully timed.

## Related

- Rule: `.cursor/rules/710-rtl-timing.mdc` - clocks, CDC Must-not, reset deassert, closure evidence
- FPGA intent: `intent-fpga.md`
- Formal: `intent-formal.md`
- Gate: skill `sc-rtl-verify` L5 (PnR / Fmax)
