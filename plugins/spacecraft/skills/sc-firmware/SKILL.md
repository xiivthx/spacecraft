---
name: sc-firmware
description: "Embedded system engineer for MCU firmware: embedded C, HAL/LL drivers, peripherals, and firmware protocols. STM32 is one vendor. Activate on STM32, nRF, HAL, CubeMX, firmware, embedded C, peripheral. Use proactively for MCU firmware. Do not activate on FPGA, SystemVerilog, RTL, or digital IC."
---

# sc-firmware

This skill is the embedded system engineer for MCU firmware (not digital IC / FPGA). Confirm OS class and target, then load matching references. Delegate production C to Task(`sc-firmware`). Commander does not write firmware.

Progressive disclosure: keep this file loaded. Shared invariants → `references/core.md`. OS class → `references/os-bare-metal.md` or `references/os-rtos.md`. ISA → `references/target-cortex-m.md`. Vendor/board → `references/target-stm32.md` or `references/target-nrf52840.md`. Glob rules `600`/`610`/`620` when editing matching files. Project SPEC / mission `spec.md` wins for OS class and target; this skill owns routing, house layering, and the verify loop.

## When to use

- **Firmware / embedded C / MCU** - application or driver work. "Embedded" here means MCU firmware, not FPGA.
- **STM32 / HAL / CubeMX** - STM32-target init, generated-code wrappers, BSP. CubeMX/HAL are STM32-target recipes, not the always-on identity.
- **nRF / nRF52840** - Nordic target path; do not force CubeMX for nRF.
- **GPIO / IRQ / UART / SPI / DMA** - peripheral bring-up or protocol
- **HIL / host unit test for firmware** - firmware verification path
- Mission build task that touches MCU firmware; **use proactively** (do not wait for an explicit skill name)

Do not treat FPGA, SystemVerilog, RTL, or digital IC as this skill - those use `sc-rtl`.

## Workflow

1. **Resolve mission** - `spacecraft resolve`. Conflict/ambiguity → `spacecraft use <selector>`.
2. **Confirm OS class and target** - Confirm OS class (bare-metal | RTOS) and target from mission `spec.md` or project SPEC before loading those references. If target omitted: handshake `needs-input` (no silent F746 default).
3. **Load matching references** - Always `references/core.md`. Then OS: `references/os-bare-metal.md` or `references/os-rtos.md`. ISA: `references/target-cortex-m.md`. Vendor/board: `references/target-stm32.md` or `references/target-nrf52840.md`.
4. **Scope layers** - Surgical edits in `app/`, `hal_if/`, `drivers/`, `bsp/`. House layering: `app/` → `hal_if/` → `drivers/` / `bsp/`. STM32 CubeMX: never edit generated `MX_*` bodies; wrap them. Glob rules `600`/`610`/`620` when editing matching files.
5. **Delegate** - Task(`sc-firmware`) for production C. Commander does not write firmware. Pair with sc-tester / rule 620 when the plan requires TDD.
6. **Verify** - `spacecraft evidence "<label>" -- <host-or-hil-test-command>`. Prefer host unit tests for logic; target/HIL for hardware-dependent paths per `620-firmware-testing.mdc` and `references/verification.md`.

### Edge cases

- **No failing test yet** - Stop. Red before green when the mission uses TDD.
- **Cache/DMA when Cortex-M7 target** - Clean/invalidate D-cache before/after DMA.
- **LTDC target** - No CPU pixel loops for LCD; use DMA2D.
- **STM32 CubeMX regenerate** - Review `git diff` before accepting; keep custom code outside `MX_*`.

## Rules

- **Must**: Resolve mission with `spacecraft resolve` before mutating work. On conflict/ambiguity use `spacecraft use <selector>`.
- **Must**: Confirm OS class and target from mission `spec.md` or project SPEC before loading OS/target references.
- **Must**: Handshake `needs-input` when target is omitted; do not pick a board.
- **Must**: Consult rules `600`, `610`, and `620` before firmware changes in their domains; load matching `references/*.md` after OS class and target are confirmed.
- **Must**: Delegate production firmware writes to Task(`sc-firmware`), not Commander.
- **Must**: Capture evidence with `spacecraft evidence` for verify steps.
- **Must**: Keep ISR short; no blocking, delay, or printf in ISR.
- **Must not**: Edit generated `MX_*` functions in place on STM32 CubeMX targets - wrap in BSP/HAL interface layers.
- **Must not**: Use dynamic allocation after init in hot paths or ISR.
- **Must not**: Skip cache clean/invalidate around DMA on Cortex-M7.
- **Must not**: Force CubeMX or STM32 HAL on nRF targets.

## Out of scope

- FPGA / SystemVerilog / RTL / digital IC - use sc-rtl
- Linux / Yocto / userspace process model - out of this skill; future pack `iot`
- Application-layer web/API work - use sc-web-frontend / sc-web-backend
- Pure architecture ADRs without firmware edits - use sc-architect / sc-adviser
- UI design for non-embedded surfaces - use sc-ux-design / Task(sc-designer)
- TDD process mechanics - use sc-tdd (still apply rule 620 for firmware test shapes)

## Output format

```
OS class: bare-metal | RTOS
Target: <MCU / board>
Rules consulted: 600 | 610 | 620
References loaded: core | os-bare-metal | os-rtos | target-cortex-m | target-stm32 | target-nrf52840
Scope:
  Files: <paths>
  Layer: app | hal_if | drivers | bsp
Delegate: Task(sc-firmware)
Verify:
  Command: <test command>
  Evidence: <label>
```

## Checklist

Before claiming firmware work done:

- [ ] Mission resolved
- [ ] OS class and target confirmed from spec before loading references
- [ ] Target omitted → handshake `needs-input`
- [ ] Matching references loaded; rules 600/610/620 consulted as needed
- [ ] House layers `app/`, `hal_if/`, `drivers/`, `bsp/` respected
- [ ] Implementation delegated to Task(`sc-firmware`)
- [ ] STM32 CubeMX: no direct `MX_*` body edits; wrappers used
- [ ] Cache/DMA and ISR constraints respected when Cortex-M7 / LTDC target
- [ ] Tests run; evidence captured with `spacecraft evidence`
- [ ] Scope limited to active plan task files

## References

- [references/core.md](references/core.md) - shared MCU invariants (every MCU task)
- [references/os-bare-metal.md](references/os-bare-metal.md) - bare-metal OS class
- [references/os-rtos.md](references/os-rtos.md) - RTOS OS class
- [references/target-cortex-m.md](references/target-cortex-m.md) - Cortex-M ISA (M4 vs M7)
- [references/target-stm32.md](references/target-stm32.md) - STM32 vendor/board recipes
- [references/target-nrf52840.md](references/target-nrf52840.md) - nRF52840 vendor/board recipes
- `.cursor/rules/600-firmware.mdc` - firmware architecture Must / Must-not
- `.cursor/rules/610-firmware-peripherals.mdc` - peripheral Must / Must-not
- `.cursor/rules/620-firmware-testing.mdc` - host / target / HIL Must / Must-not
- `references/peripherals.md` - peripheral and protocol examples (on-demand)
- `references/verification.md` - Unity, target, HIL, CI recipes (on-demand)
- `.cursor/agents/sc-firmware.md` - write-capable firmware agent
