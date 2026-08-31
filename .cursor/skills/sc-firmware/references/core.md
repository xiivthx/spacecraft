# Firmware core recipes

On-demand shared MCU core: house layering, Power-of-Ten-class risk rules, and FDIR handshake. Critical Must/Must-not stay in `.cursor/rules/600-firmware.mdc`. Load this file for every MCU firmware task. OS class refs and target refs load from SKILL routing; this file is the shared core.

## House layering

```
app/        Application logic (state machines, protocols, business rules)
hal_if/     Portable hardware interface (no vendor headers in app/)
drivers/    Thin wrappers around vendor HAL or SDK
bsp/        Board pins, clocks, startup
```

Dependency: `app/` -> `hal_if/` -> `drivers/` / `bsp/`. `app/` includes only `hal_if/`. Wrap vendor HAL behind `hal_if/`. Never `app/` -> `drivers/` directly. No circular includes.

cFS, F Prime, and Zephyr: match the consuming project. Do not force those frameworks as the house default.

## Power-of-Ten-class invariants

Risk rules only (not a NASA JPL PDF paste):

- no recursion (stack overflow on small MCUs)
- no heap after init in hot paths or ISR; no `malloc` in the loop
- bounded loops (every loop has a proven upper bound)
- checked returns on every non-void call; validate parameters at the boundary
- restrict preprocessor to include guards and simple macros; prefer `inline`
- static analysis for flight-grade claims (warnings-as-errors plus a checker before the claim)
- fixed-width types from `stdint.h` (`uint8_t`, `uint16_t`, `uint32_t`)

## FDIR handshake

failing assert / WDT / HardFault -> documented safe state recipe

Assertions stay enabled when claiming FDIR. Do not strip `assert` from the flight image when the claim is FDIR. Each failure class (assert, WDT timeout, HardFault) has a documented safe state: outputs held or deasserted, optional reset, optional degrade mode, and a log of PC/LR to retained RAM when the core allows it. Kick the WDT in the main loop, not in ISR.

## Coding patterns

MCU-wide. Board clocks, memory maps, and vendor init live in target refs.

**Types.** `stdint.h` only for widths. `volatile` for MMIO and ISR-shared data. `static` for file-local symbols. `const` data in Flash.

**ISR.** Keep it short: set a flag or wake a task; do the work in task or main-loop context. Shared variables: `volatile` plus atomic access or a critical section. Never `printf`, delay, or blocking take in ISR.

```c
void irq_handler(void) {
    uint8_t byte = hw_read_rx();
    rx_buffer_put(&rx_ring, byte);
    task_wake_from_isr(rx_task);
}
```

**Memory.** No heap after init; use static buffers or a pool filled at startup. Review the linker script for stack size. `.bss` should match the map file.

**Macros.** Parenthesize. Hardware addresses and bit masks only; no function-like control flow in the preprocessor.

## Related

- Rule: `.cursor/rules/600-firmware.mdc` - Must / Must-not invariants
- Peripherals: `peripherals.md`
- Verification: `verification.md`
- OS class and target recipes: load from SKILL routing
