# STM32 target recipes

On-demand STM32 board, CubeMX, clock, memory, and vendor SDK recipes. Critical Must/Must-not stay in `.cursor/rules/600-firmware.mdc`. CubeMX/HAL/LL belong here; shared core 600 does not require CubeMX. Load this file when the vendor is STM32. Cortex-M4 vs Cortex-M7 D-cache / MPU / DMA: `target-cortex-m.md`.

## CubeMX, HAL, and LL

Wrap generated `MX_*` in `bsp/` / `hal_if/`. Do not edit `MX_*` bodies in place. On regenerate, review `git diff` before accepting.

Prefer LL on hot paths; HAL is fine for `MX_*` init, USB/Ethernet/SDIO, and bring-up:

```c
HAL_GPIO_WritePin(GPIOA, GPIO_PIN_5, GPIO_PIN_SET);  // bloated
LL_GPIO_SetOutputPin(GPIOA, LL_GPIO_PIN_5);          // lean
```

Always wrap in `hal_if/` - never expose a HAL handle to `app/`:

```c
// hal_if/gpio.h
void gpio_set(uint8_t pin);
bool gpio_read(uint8_t pin);
```

## STM32F746NG-Discovery

Cortex-M7. Discovery kit: 4.3" 480x272 LCD with capacitive touch. F7 adds L1 cache, SDRAM, LTDC, DMA2D, and double-precision FPU. Apply M7 cache/MPU/DMA from `target-cortex-m.md`.

- **L1 cache**: 4KB I-cache + 4KB D-cache - clean/invalidate around DMA
- **SDRAM**: 128 Mbit external - framebuffer, large buffers
- **LTDC**: hardware LCD controller - not SPI bit-bang
- **DMA2D**: Chrom-Art blit, fill, blend, color convert - no CPU pixel loops
- **QSPI**: 128 Mbit, memory-mapped XIP or assets
- **MPU**: Flash cacheable, SRAM non-cacheable (DMA), SDRAM write-through, peripherals strongly-ordered

### CubeMX2 (216 MHz)

Essential peripherals in `.ioc`:

```
RCC:   HSE (25 MHz external), PLL -> 216 MHz SYSCLK, 108 MHz APB2, 54 MHz APB1
SDRAM: 128 Mbit, 16-bit, 2 banks, 4K refresh
LTDC:  480x272 RGB565, pixel clock 9 MHz, HSYNC 41, VSYNC 10, back porch 13/2, front porch 32/2
DMA2D: enable for 2D acceleration
FMC:   Bank 2 (SDRAM), Bank 1 (QSPI)
CRC:   enable for CRC16/32 hardware
```

### MPU regions (this board)

```c
void mpu_init(void) {
    mpu_config_region(0, 0x08000000,
        MPU_RASR_ATTR_CACHEABLE_WB | MPU_RASR_SIZE_1MB);     // Flash
    mpu_config_region(1, 0x20010000,
        MPU_RASR_ATTR_NON_CACHEABLE | MPU_RASR_SIZE_256KB);  // SRAM1, DMA buffers
    mpu_config_region(2, 0xC0000000,
        MPU_RASR_ATTR_CACHEABLE_WT | MPU_RASR_SIZE_8MB);     // SDRAM framebuffer
    mpu_config_region(3, 0x40000000,
        MPU_RASR_ATTR_DEVICE | MPU_RASR_SIZE_512MB);         // peripherals
    SCB->SHCSR |= SCB_SHCSR_MEMFAULTENA_Msk;
    MPU->CTRL  |= MPU_CTRL_ENABLE_Msk;
}
```

### SDRAM init

```c
void sdram_init(void) {
    SDRAM_HandleTypeDef hsdram1;
    FMC_SDRAM_TimingTypeDef timing = {
        .LoadToActiveDelay    = 2,
        .ExitSelfRefreshDelay = 7,
        .SelfRefreshTime      = 4,
        .RowCycleDelay        = 7,
        .WriteRecoveryTime    = 2,
        .RPDelay              = 2,
        .RCDDelay             = 2,
    };
    HAL_SDRAM_Init(&hsdram1, &timing);

    uint32_t* sdram = (uint32_t*)0xC0000000;
    sdram[0] = 0xDEADBEEF;
    if (sdram[0] != 0xDEADBEEF) { Error_Handler(); }
    sdram[0] = 0;
}
```

### Linker (SDRAM + LCD fb)

```ld
MEMORY {
  FLASH    (rx)  : ORIGIN = 0x08000000, LENGTH = 1024K
  RAM      (rwx) : ORIGIN = 0x20000000, LENGTH = 320K
  SDRAM    (rwx) : ORIGIN = 0xC0000000, LENGTH = 8M
  QSPI     (rx)  : ORIGIN = 0x90000000, LENGTH = 16M
}

SECTIONS {
  .sdram_section (NOLOAD): {
    *(.sdram)
    *(.sdram*)
  } > SDRAM

  .lcd_fb (NOLOAD): {
    *(.lcd_fb)
  } > SDRAM
}
```

```c
__attribute__((section(".lcd_fb"))) uint16_t lcd_fb[480][272];
__attribute__((section(".sdram"))) uint8_t asset_pool[4 * 1024 * 1024];
```

### Memory map

```
0x0800_0000  Flash (1MB)
0x1FF0_0000  System memory (bootloader)
0x2000_0000  SRAM1 (320KB)
0x6000_0000  FMC Bank 1 (QSPI)
0x9000_0000  QSPI memory-mapped (16MB)
0xC000_0000  FMC Bank 2 - SDRAM (8MB)
0xE000_0000  Cortex-M7 system control
```

## NUCLEO-H723ZG

STM32H723ZGT6. Cortex-M7 up to 550 MHz. Nucleo-144 with ST-LINK. CubeMX/HAL/LL in scope. No onboard LCD like Discovery. Apply M7 cache/MPU/DMA from `target-cortex-m.md`.

**Clock.** VOS0 for 550 MHz; set Flash wait states for that scale. Nucleo HSE is typically 8 MHz from ST-LINK MCO, not the Discovery 25 MHz crystal. PLL1 for SYSCLK; separate PLLs for USB/kernel clocks. Do not paste the F746 216 MHz CubeMX tree.

**Memory.** Split SRAM: ITCM (code), DTCM (CPU-tight data; DMA does not share this view), AXI SRAM (DMA buffers). About 1 MB Flash, ~564 KB SRAM total. No onboard SDRAM - do not reuse `0xC0000000` or the Discovery linker. MPU: cacheable Flash/ITCM, non-cacheable DMA SRAM (AXI), peripherals strongly-ordered.

**Debug.** ST-LINK/V3 on the Nucleo-144. SWD + SWO. Disconnect ST-LINK (CN2) for an external probe. High SYSCLK: connect-under-reset if attach fails.

**Vendor SDK.** H7 HAL pack, not F7. In CubeMX enable I-cache and D-cache on Cortex_M7. Pick DMA1/2, BDMA, or MDMA for the bus that can reach the buffer - DTCM is the usual trap. Ethernet/USB on this Nucleo need their own clock domains and PHY/USB setup.

## NUCLEO-L412KB

STM32L412KBU6. Cortex-M4 80 MHz. Nucleo-32. 128 KB flash / 40 KB SRAM. CubeMX/HAL/LL in scope. No F7 D-cache, no LTDC. Do not copy F746 cache or LCD Musts onto L412.

**Clock.** Range 1 for 80 MHz. MSI or 8 MHz HSE from ST-LINK MCO. No 216 MHz PLL, no LTDC pixel clock.

**Memory.** 40 KB SRAM is the whole budget - static buffers, no SDRAM heap, no 480x272 framebuffer. DMA is CPU-coherent; skip `cache_clean` / D-cache enable. MPU, if used, is protection only.

**Debug.** Nucleo-32 ST-LINK/V2-1. SWD on the ST-LINK; the UFQFPN32 is pin-starved. Disconnect ST-LINK solder bridges for an external probe.

**Vendor SDK.** L4 HAL pack. CubeMX on this board: GPIO/UART/SPI/I2C/timers as needed. Do not enable LTDC, DMA2D, FMC SDRAM, or F7 cache options.

## Related

- Rule: `.cursor/rules/600-firmware.mdc` - Must / Must-not invariants
- Core: `core.md`
- Cortex-M ISA: `target-cortex-m.md`
- Peripherals: `peripherals.md`
- Verification: `verification.md`
