# Cortex-M ISA recipes

On-demand Cortex-M ISA: M4 vs M7, D-cache, MPU, and DMA coherency. Critical Must/Must-not stay in `.cursor/rules/600-firmware.mdc`. Load this file when the core is Cortex-M. Do not impose Cortex-M7 D-cache rules on M4.

## Cortex-M4 vs Cortex-M7

**Cortex-M4** (L412, nRF52840): typically no data cache. CPU and DMA share the same SRAM view. DMA coherency is not the F7 problem. MPU, when present, is for execute-never / privilege, not a cacheable-vs-DMA split.

**Cortex-M7** (F746, H723ZG): I-cache + D-cache. MUST clean/invalidate D-cache around DMA. MPU regions for cacheable vs DMA buffers vs peripherals.

## D-cache and DMA (Cortex-M7 only)

Enable after clock config, before any DMA or peripheral init. Line size is 32 bytes.

```c
SCB_EnableICache();
SCB_EnableDCache();

// Flush before DMA writes to memory (ADC, SPI RX, I2C RX)
void cache_clean_invalidate(void* addr, uint32_t size) {
    uint32_t start = (uint32_t)addr & ~0x1F;  // 32-byte cache line
    uint32_t end   = ((uint32_t)addr + size + 31) & ~0x1F;
    for (uint32_t i = start; i < end; i += 32) {
        SCB->DCCIMVAC = i;  // clean + invalidate by MVA
    }
    __DSB();
    __ISB();
}

// Clean before DMA reads from memory (SPI TX, I2C TX)
void cache_clean(void* addr, uint32_t size) {
    uint32_t start = (uint32_t)addr & ~0x1F;
    uint32_t end   = ((uint32_t)addr + size + 31) & ~0x1F;
    for (uint32_t i = start; i < end; i += 32) {
        SCB->DCCMVAC = i;  // clean by MVA
    }
    __DSB();
}
```

- **DMA source buffer**: `cache_clean()` before DMA TX
- **DMA dest buffer**: `cache_clean_invalidate()` after DMA RX, before CPU read
- **ISR-shared DMA buffer**: clean/invalidate before the main loop reads
- **Display consumer**: write-through MPU or `cache_clean()` before reload; board map is the vendor target

Skip this section on M4.

## MPU (Cortex-M7)

Split cacheable code/data from DMA buffers and MMIO. Bases and sizes are board-specific (vendor target). Typical regions:

- Flash / code: normal, cacheable
- DMA buffers: normal, non-cacheable (or 32-byte aligned plus clean/invalidate)
- External RAM / framebuffer when present: write-through, or clean before the consumer
- Peripherals: device memory, strongly-ordered

```c
void mpu_init(void) {
    // Region 0: Flash - normal, cacheable, not shareable
    mpu_config_region(0, 0x08000000,
        MPU_RASR_ATTR_CACHEABLE_WB | MPU_RASR_SIZE_1MB);

    // Region 1: DMA SRAM - normal, non-cacheable (base/size from board)
    mpu_config_region(1, dma_sram_base,
        MPU_RASR_ATTR_NON_CACHEABLE | dma_sram_size);

    // Region 2: Peripherals - device, strongly-ordered
    mpu_config_region(2, 0x40000000,
        MPU_RASR_ATTR_DEVICE | MPU_RASR_SIZE_512MB);

    SCB->SHCSR |= SCB_SHCSR_MEMFAULTENA_Msk;
    MPU->CTRL  |= MPU_CTRL_ENABLE_Msk;
}
```

Add a write-through region when the board has external RAM used as a framebuffer. Do not copy an F7 SDRAM map onto M4 or onto an M7 Nucleo with no external RAM.

## Related

- Rule: `.cursor/rules/600-firmware.mdc` - Must / Must-not invariants
- Core: `core.md`
- STM32 boards: `target-stm32.md`
- Peripherals: `peripherals.md`
- Verification: `verification.md`
