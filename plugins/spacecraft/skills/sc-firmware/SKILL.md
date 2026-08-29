---
name: sc-firmware
description: "Implement STM32 embedded C, HAL/LL drivers, peripherals, CubeMX, and firmware protocols. Activate on STM32, HAL, firmware, embedded C, peripheral, CubeMX. Use proactively for embedded/firmware implementation."
---

# sc-firmware

Implement and maintain STM32 ARM Cortex-M firmware under mission control. Architecture, HAL/LL drivers, peripherals, CubeMX2 layout, and host/target verification. Delegate production writes to Task(`sc-firmware`), not Commander.

## When to use

Activate when the user asks to:

- **"STM32" / "firmware" / "embedded C"** - MCU application or driver work
- **"HAL" / "LL driver" / "CubeMX"** - peripheral init, generated code wrappers, BSP
- **"GPIO / IRQ / LCD / UART / SPI / DMA"** - peripheral bring-up or protocol
- **"HIL" / "host unit test for firmware"** - firmware verification path
- When a mission task requires embedded or firmware implementation
- **Use proactively** for embedded/firmware implementation (do not wait for an explicit skill name)

## Workflow

Use this exact sequence unless the user specifies otherwise:

1. **Resolve mission** - `spacecraft resolve`. On conflict or ambiguity, use `spacecraft use <selector>`.

2. **Consult firmware rules** - Read before editing (invariants only; load recipes on demand):
   - `.cursor/rules/600-firmware.mdc` - architecture Must / Must-not
   - `.cursor/rules/610-firmware-peripherals.mdc` - peripheral Must / Must-not
   - `.cursor/rules/620-firmware-testing.mdc` - test Must / Must-not

3. **Load recipes as needed** (on-demand, not always):
   - `references/architecture.md` - board, CubeMX, MPU, linker, cache helpers, layering
   - `references/peripherals.md` - GPIO/UART/SPI/I2C/LCD/QSPI/LF/UHF/DMA examples
   - `references/verification.md` - Unity/Ceedling, target, HIL, CI

4. **Confirm target and scope** - Match board/MCU from `spec.md` and existing CubeMX2 layout. Prefer surgical changes inside `app/`, `hal_if/`, `drivers/`, `bsp/` - never edit generated `MX_*` bodies; wrap them.

5. **Delegate implementation** - Task(`sc-firmware`) for production C writes. Commander does not write firmware source. Pair with sc-tester / rule 620 for failing tests and evidence when the plan requires TDD.

6. **Verify** - `spacecraft evidence "<label>" -- <host-or-hil-test-command>`. Prefer host unit tests for logic; target/HIL for hardware-dependent paths per `620-firmware-testing.mdc` and `references/verification.md`.

### Edge cases

- **No failing test yet** - Stop. Red before green when the mission uses TDD.
- **Cache/DMA on F7** - Clean/invalidate D-cache before/after DMA; framebuffer write-through SDRAM.
- **CubeMX regenerate** - Review `git diff` before accepting; keep custom code outside `MX_*`.

## Rules

- **Must**: Resolve mission with `spacecraft resolve` before mutating work. On conflict/ambiguity use `spacecraft use <selector>`.
- **Must**: Consult rules `600`, `610`, and `620` before firmware changes in their domains; load matching `references/*.md` when implementing recipes.
- **Must**: Delegate production firmware writes to Task(`sc-firmware`), not Commander.
- **Must**: Capture evidence with `spacecraft evidence` for verify steps.
- **Must**: Keep ISR short; no blocking, delay, or printf in ISR.
- **Must not**: Edit generated `MX_*` functions in place - wrap in BSP/HAL interface layers.
- **Must not**: Use dynamic allocation after init in hot paths or ISR.
- **Must not**: Skip cache clean/invalidate around DMA on Cortex-M7.

## Out of scope

- Application-layer web/API work - use sc-web-frontend / sc-web-backend
- Pure architecture ADRs without firmware edits - use sc-architect / sc-adviser
- UI design for non-embedded surfaces - use sc-ux-design / Task(sc-designer)
- TDD process mechanics - use sc-tdd (still apply rule 620 for firmware test shapes)

## Output format

```
Target: <MCU / board>
Rules consulted: 600 | 610 | 620
Recipes loaded: architecture | peripherals | verification | none
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
- [ ] Rules 600/610/620 consulted as needed; recipes loaded when implementing
- [ ] Implementation delegated to Task(`sc-firmware`)
- [ ] No direct `MX_*` body edits; wrappers used
- [ ] Cache/DMA and ISR constraints respected on F7 when applicable
- [ ] Tests run; evidence captured with `spacecraft evidence`
- [ ] Scope limited to active plan task files

## References

- `.cursor/rules/600-firmware.mdc` - firmware architecture Must / Must-not
- `.cursor/rules/610-firmware-peripherals.mdc` - peripheral Must / Must-not
- `.cursor/rules/620-firmware-testing.mdc` - host / target / HIL Must / Must-not
- `references/architecture.md` - board, CubeMX, MPU, linker, cache, layering recipes (on-demand)
- `references/peripherals.md` - peripheral and protocol examples (on-demand)
- `references/verification.md` - Unity, target, HIL, CI recipes (on-demand)
- `.cursor/agents/sc-firmware.md` - write-capable firmware agent
