---
name: sc-rtl
description: >-
  Multi-cycle RISC-V SystemVerilog RTL: FSM, memory handshakes, coding
  conventions, FPGA timing, lint/synth gates, observe-first HW debug.
  Activate on .sv, SystemVerilog, RTL, CPU FSM, imem/dmem, Verilator, Yosys,
  nextpnr, or FPGA pack work. Use proactively for RTL implementation.
---

# sc-rtl

Implement and maintain SystemVerilog RTL under mission control (multi-cycle RISC-V CPU defaults). Delegate production `.sv` writes to Task(`sc-rtl`) when available, else Task(`sc-coder`) with this skill active. Commander does not write RTL.

Progressive disclosure: keep this file loaded. Read [references/arch.md](references/arch.md) only for ISA lists, cycle counts, or module tree. Verify gates → skill `sc-rtl-verify`. Glob rules `700`/`710`/`720` apply when editing matching files.

**SoT conflict:** if project `docs/SPECIFICATION.md` (or mission `spec.md`) narrows ISA/FPGA/paths, SPEC wins for product scope. This skill still owns SV conventions, debug discipline, and quality loop.

## When to use

- **SystemVerilog / `.sv` / RTL** - module edit, new block, decode, FSM
- **CPU architecture** - RV32I/M/A/C/F/Zicsr, multi-cycle states, mem interface
- **Verilator lint / Yosys / nextpnr / FPGA synth** - gates after SV change
- **HW debug** - FSM stuck, handshake hang, wrong ALU/mem result
- Mission build task that touches RTL; **use proactively** for FPGA-pack RTL work

## Workflow

1. **Resolve mission** - `spacecraft resolve`. Conflict/ambiguity → `spacecraft use <selector>`.
2. **Confirm scope** - Match ISA/target from `spec.md` / project SPEC. Prefer surgical module edits.
3. **Apply hard rules** (below) before writing SV.
4. **Delegate** - Task(`sc-rtl`) or Task(`sc-coder`) for production SV. Pair with sc-tester when plan requires TDD (cocotb/Verilator TB as project defines).
5. **Verify** - lint then synth (or project Makefile). Prefer `sc-rtl-verify` layer table + `spacecraft evidence`. HW bugs: observe-first (below) before abstract reasoning.

### Hard rules

1. Every `.sv`: start `` `default_nettype none ``, end `` `default_nettype wire ``.
2. Sync reset only inside RTL. Default active-high `rst`. Convert board active-low at boundary.
3. Sequential: `always_ff @(posedge clk)` with `if (rst)` inside.
4. Reset **control** bits (`valid`/`pending`), not payload-only regs. Refresh payload when asserting control; ignore payload when control low.
5. `snake_case`; purpose prefixes (`imem_`, `dmem_`, `alu_`). Keep RISC-V names (`rs1`, `rs2`, `rd`, `funct3`).
6. Prefer `case` for mutually exclusive FSM/decode. Nested `if` only when priority is real.
7. Fmax first: register at natural boundaries; multi-cycle stage large cones. Exempt `*_ready`-style returns; document other unregistered exemptions.
8. No latches. Staging = FFs. `x0` hardwired zero in hardware.
9. After SV change: lint, then synth. No done claim without fresh evidence.

### Architecture (compact)

Multi-cycle **non-pipelined** RV32IMACF_Zicsr defaults. External imem/dmem. Variable latency via ready/valid. `instr_complete`: 1-cycle pulse at finish.

| State | Hex | Role |
|-------|-----|------|
| S_BOOT | 0x0 | Wait boot after reset |
| S_FETCH | 0x1 | `imem_req`, wait `imem_ready` |
| S_DECODE | 0x2 | Decode, read regs |
| S_EXECUTE | 0x3 | ALU |
| S_MEM_ADDR | 0x4 | Load/store address |
| S_MEM_READ | 0x5 | `dmem` read, wait ready |
| S_MEM_WRITE | 0x6 | `dmem` write, wait ready |
| S_WRITEBACK | 0x7 | Write rd |
| S_BRANCH | 0x8 | Branch compare + PC |
| S_CSR | 0x9 | CSR op |
| S_HALT | 0xA | ECALL/EBREAK |
| S_ATOMIC_RMW | 0xB | Atomic RMW |

**imem:** `imem_req`, `imem_ready`, `imem_addr`, `imem_data`  
**dmem:** `dmem_req`, `dmem_ready`, `dmem_addr`, `dmem_wdata`, `dmem_rdata`, `dmem_we`, `dmem_re`, `dmem_size`

### Debug (non-negotiable)

Observe first; reason from evidence. `$display` values + FSM transitions; cycle logs. No assumed wires / predicted FSM / guessed timing.

### AI RTL quality loop

1. Lint (static analysis if available) before long sim debates.
2. Disposition findings — do not blind-fix-all:
   - **Fix in block** — unread/readback gaps, input regs, ITE→case for exclusive FSM
   - **Waive / system** — reset sync deassert at top-level controller; document; no sync buried in reusable IP
   - **Monitor** — FPGA fanout; act only if post-impl timing fails
3. Re-lint/re-synth. Close when gates clean or waivers documented.

### Verify commands

Prefer repo Makefile/CI. Else (adjust if tree is `src/rtl/`):

```bash
find rtl -name '*.sv' -exec verilator --lint-only --Wno-MULTITOP {} +
(cd rtl/fpga && make)
```

## Must / Must not

- **Must**: Resolve mission before mutating RTL work
- **Must**: Delegate production `.sv` to Task(`sc-rtl`) / Task(`sc-coder`)
- **Must**: Capture evidence with `spacecraft evidence` for verify steps
- **Must**: Observe-first on HW bugs (`$display` / sim) before claiming root cause
- **Must not**: Skip `default_nettype` guards on new `.sv`
- **Must not**: Async reset inside project RTL modules (sync `rst` default)
- **Must not**: Invent signal behavior without sim evidence

## Out of scope

- STM32 / HAL / CubeMX - sc-firmware
- App web/API - sc-web-frontend / sc-web-backend
- Pure ADRs without RTL edits - sc-architect / sc-adviser
- TDD process mechanics - sc-tdd

## Checklist

- [ ] Mission resolved
- [ ] SPEC/ISA scope confirmed
- [ ] Hard rules applied
- [ ] Delegated Task(`sc-rtl`/`sc-coder`)
- [ ] Lint + synth (or project verify) + evidence

## References

- [references/arch.md](references/arch.md) - hierarchy, cycles, ISA inventory, disposition cheat sheet
- Skill `sc-rtl-verify` - lint/sim/ISA/formal/STA gates
- Rules `700-rtl.mdc`, `710-rtl-timing.mdc`, `720-rtl-verify.mdc`
- `.cursor/agents/sc-rtl.md` - write-capable RTL agent
