# nRF52840 target recipes

On-demand vendor/board recipes for the nRF52840 DK. Load when the mission target is nRF52840. Shared MCU invariants stay in `core.md`. Cortex-M4 vs M7 cache, MPU, and DMA: `target-cortex-m.md`.

## Board: nRF52840 DK PCA10056

Full DK (PCA10056), not the USB dongle. Cortex-M4F at ~64 MHz, 1 MB flash, 256 KB RAM, onboard J-Link. Board capabilities: BLE, 802.15.4, and USB. The nRF52840 path does not require CubeMX/HAL.

## Vendor SDK

nRF Connect SDK, nRF5 SDK, or nrfx - not STM32 CubeMX. This path does not depend on STM32 CubeMX. Wrap nrfx (or the chosen Nordic SDK) in `drivers/` / `bsp/`. House layering stays `app/` -> `hal_if/` -> `drivers/` / `bsp/`. Match a consuming Zephyr or nRF Connect SDK tree when the project already uses it; do not force Zephyr as the house default.

## Gotchas

**Clock.** Typical HFCLK source is the 32 MHz HFXO. Start HFXO before the radio. LFCLK is usually the 32.768 kHz LFXO for RTC.

**Memory map.** Flash at `0x00000000` (1 MB). SRAM at `0x20000000` (256 KB). Review the linker script for stack and `.bss`.

**Debug.** Onboard J-Link over SWD.

**Vendor SDK.** Stay on one Nordic path (nRF Connect SDK / nRF5 SDK / nrfx) for this board.

## ISA

nRF52840 is Cortex-M4F. No F7 D-cache DMA ritual. M4 vs M7: `target-cortex-m.md`.

## Related

- Core: `core.md`
- ISA: `target-cortex-m.md`
- Rule: `.cursor/rules/600-firmware.mdc`
