# RTL architecture reference (load on demand)

## Module hierarchy

```
top (CPU)
├── fetch_buffer      # RV32C alignment
├── decompress        # RV32C → 32-bit, combinational
├── decoder
├── alu               # RV32I + M
│   └── div_unit
├── regfile
├── csr_file
├── branch_unit       # dedicated compare (not ALU)
├── mem_interface
└── writeback_mux
```

Also: 32 FP regs `f0-f31` + FCSR when F enabled. Debug MMIO FIFO often at `0x40000000` (packet protocol) - confirm against current SPEC/memory map.

## Multi-cycle FSM states

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

## Base instruction cycles

Memory latency adds wait cycles in FETCH / MEM_READ / MEM_WRITE (e.g. 3-cycle mem → load ≈ 5 base + 3 + 3 = 11).

| Class | Base | Path |
|-------|------|------|
| R-type / I-arith / M / LUI/AUIPC / JAL/JALR | 4 | F→D→EX→WB |
| Load | 5 | F→D→MA→MR→WB |
| Store | 4 | F→D→MA→MW |
| Branch | 3 | F→D→BR |
| FENCE | 2 | F→D |
| ECALL/EBREAK | 2 | F→D→HALT |
| CSR | 4 | F→D→CSR→WB |

## ISA inventory (RV32IMACF_Zicsr ≈ 118)

Product SPEC may subset (e.g. RV32IMC only). Do not implement retired extensions against SPEC.

**RV32I (40):** ADD ADDI SUB; AND ANDI OR ORI XOR XORI; SLL SLLI SRL SRLI SRA SRAI; SLT SLTI SLTU SLTIU; BEQ BNE BLT BGE BLTU BGEU; LW LH LB LHU LBU SW SH SB; LUI AUIPC; JAL JALR; FENCE; ECALL EBREAK

**M (8):** MUL MULH MULHSU MULHU; DIV DIVU; REM REMU

**A (11):** LR.W SC.W; AMOSWAP.W AMOADD.W AMOXOR.W AMOAND.W AMOOR.W; AMOMIN.W AMOMAX.W AMOMINU.W AMOMAXU.W

**C (27):**  
Q0: C.ADDI4SPN C.LW C.SW  
Q1: C.NOP C.ADDI C.JAL C.LI C.ADDI16SP C.LUI C.SRLI C.SRAI C.ANDI C.SUB C.XOR C.OR C.AND C.J C.BEQZ C.BNEZ  
Q2: C.SLLI C.LWSP C.JR C.MV C.EBREAK C.JALR C.ADD C.SWSP

**F (26):** FADD.S FSUB.S FMUL.S FDIV.S FSQRT.S; FMADD.S FMSUB.S FNMSUB.S FNMADD.S; FMIN.S FMAX.S; FSGNJ.S FSGNJN.S FSGNJX.S; FEQ.S FLT.S FLE.S; FCVT.W.S FCVT.WU.S FCVT.S.W FCVT.S.WU; FLW FSW; FMV.X.W FMV.W.X FCLASS.S  
IEEE 754-2008; FCSR rounding + flags.

**Zicsr (6):** CSRRW CSRRS CSRRC CSRRWI CSRRSI CSRRCI

## Design decisions (why)

- Multi-cycle FSM (12 states) + variable-latency mem handshakes
- Exposed imem/dmem (TB/FPGA owns memory)
- Separate `branch_unit`; FF staging (FPGA-safe)
- Payload regs often not reset (cut reset fanout); control bits are

## Static-analysis disposition cheat sheet

| Finding class | Typical action |
|---------------|----------------|
| Async reset deassert / reset sync | System reset controller + waiver; not inside IP |
| Unregistered control inputs | Fix: register at block boundary |
| High fanout on FPGA | Monitor; fix only if timing report fails |
| Driven-but-unread / missing protocol readback | Fix: functional gap |
| Deep ITE on exclusive states | Fix: `case` for parallel select |
