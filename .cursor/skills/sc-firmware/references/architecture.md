# Firmware architecture recipes

On-demand recipes for STM32 board setup, CubeMX, cache/MPU, linker, layering, and coding patterns. Critical Must/Must-not stay in `.cursor/rules/600-firmware.mdc`. Load this file when implementing firmware architecture.

## Target: STM32F746NG-Discovery (Cortex-M7)

The F7 series adds L1 cache, SDRAM controller, LTDC, DMA2D, and double-precision FPU. The Discovery kit has 4.3" 480×272 LCD with capacitive touch.

Key differences from F4 (Cortex-M4):
- **L1 Cache**: 4KB I-cache + 4KB D-cache - must flush/invalidate before/after DMA
- **SDRAM**: 128 Mbit external - framebuffer, heap, large buffers live here
- **LTDC**: hardware LCD controller - replaces SPI bit-bang for display
- **DMA2D**: Chrom-Art 2D accelerator - hardware blit, fill, blend, color convert
- **QSPI**: Quad-SPI Flash 128 Mbit - memory-mapped for execute-in-place or asset storage
- **MPU**: must configure regions: Flash (cacheable), SRAM (non-cacheable), SDRAM (write-through), peripherals (strongly-ordered)

## CubeMX2 Setup (STM32F746NG-Discovery)

Essential peripherals to enable in `.ioc`:
```
RCC:   HSE (25 MHz external), PLL → 216 MHz SYSCLK, 108 MHz APB2, 54 MHz APB1
SDRAM: 128 Mbit, 16-bit, 2 banks, 4K refresh
LTDC:  480×272 RGB565, pixel clock 9 MHz, HSYNC 41, VSYNC 10, back porch 13/2, front porch 32/2
DMA2D: enable for 2D acceleration
FMC:   Bank 2 (SDRAM), Bank 1 (QSPI)
CRC:   enable for CRC16/32 hardware
```

## Cache Management (Cortex-M7)

```c
// Enable after SystemClock_Config(), before any DMA or peripheral init
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

Recipe notes:
- **DMA source buffer**: `cache_clean()` before DMA TX
- **DMA dest buffer**: `cache_clean_invalidate()` after DMA RX, before reading
- **Framebuffer in SDRAM**: mark as write-through (MPU) or clean manually before LTDC refresh
- **Cache coherency in ISR**: DMA buffer touched by ISR → clean/invalidate before main loop reads

## MPU Configuration

```c
void mpu_init(void) {
    // Region 0: Flash (0x0800_0000) - normal memory, cacheable, not shareable
    mpu_config_region(0, 0x08000000,
        MPU_RASR_ATTR_CACHEABLE_WB | MPU_RASR_SIZE_1MB);

    // Region 1: SRAM1 (0x2001_0000) - normal, non-cacheable (DMA buffers live here)
    mpu_config_region(1, 0x20010000,
        MPU_RASR_ATTR_NON_CACHEABLE | MPU_RASR_SIZE_256KB);

    // Region 2: SDRAM (0xC000_0000) - normal, write-through (framebuffer)
    mpu_config_region(2, 0xC0000000,
        MPU_RASR_ATTR_CACHEABLE_WT | MPU_RASR_SIZE_8MB);

    // Region 3: Peripherals (0x4000_0000) - device memory, strongly-ordered
    mpu_config_region(3, 0x40000000,
        MPU_RASR_ATTR_DEVICE | MPU_RASR_SIZE_512MB);

    SCB->SHCSR |= SCB_SHCSR_MEMFAULTENA_Msk;
    MPU->CTRL  |= MPU_CTRL_ENABLE_Msk;
}
```

## SDRAM Initialization

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

## Linker Script for SDRAM

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

## Memory Map (STM32F746NG-Discovery)

```
0x0800_0000  Flash (1MB)
0x1FF0_0000  System memory (bootloader)
0x2000_0000  SRAM1 (320KB)
0x6000_0000  FMC Bank 1 (QSPI)
0x9000_0000  QSPI memory-mapped (16MB)
0xC000_0000  FMC Bank 2 - SDRAM (8MB)
0xE000_0000  Cortex-M7 system control
```

## Project Layering (top to bottom)

```
app/            Application logic (state machines, protocols, business rules)
  ├── tasks/    FreeRTOS tasks (if RTOS used)
  └── modules/  Feature modules (protocol, display, comms)
hal_if/         Hardware Abstraction Interface (your own wrapper)
  ├── gpio.h    Never include STM32 headers in app/
  ├── uart.h
  ├── spi.h
  └── timer.h
drivers/        Peripheral drivers (thin wrappers around HAL)
  ├── lcd/      LCD driver (ILI9341, ST7789, etc)
  ├── radio/    LF/UHF transceiver driver
  └── sensor/   Sensor drivers
bsp/            Board Support Package - pin mappings, clock config
  ├── bsp.h     Board-specific #defines (SPI1_CS_PORT, LCD_RST_PIN, etc)
  └── board.c   MX_* init functions from CubeMX
third_party/    STM32 HAL, CMSIS, FreeRTOS, middleware
```

Dependency:
- `app/` includes only `hal_if/` - never touches HAL directly
- `hal_if/` includes `drivers/` - wraps HAL in portable interfaces
- `drivers/` includes `bsp/` + CMSIS - thin, replaceable
- `bsp/` includes STM32 HAL headers
- Never: `app/` → `drivers/` (skip the interface layer)
- Never: circular includes

## HAL Usage Patterns

Prefer LL (Low-Layer) over HAL for performance-critical paths:
```c
HAL_GPIO_WritePin(GPIOA, GPIO_PIN_5, GPIO_PIN_SET);  // bloated
LL_GPIO_SetOutputPin(GPIOA, LL_GPIO_PIN_5);          // lean
```

HAL is acceptable for initialization (`MX_*`), complex peripherals (USB, Ethernet, SDIO), and prototyping.

Always wrap in your own interface:
```c
// hal_if/gpio.h - never expose HAL handle to app code
void gpio_set(uint8_t pin);
bool gpio_read(uint8_t pin);
```

## State Machine Pattern

Flat state machine for event-driven code:
```c
typedef enum {
    STATE_IDLE,
    STATE_TX_ACTIVE,
    STATE_TX_DONE,
    STATE_ERROR,
} radio_state_t;

typedef struct {
    radio_state_t state;
    uint32_t     last_event_time;
    void*        context;
} radio_fsm_t;

void radio_fsm_dispatch(radio_fsm_t* fsm, radio_event_t event) {
    switch (fsm->state) {
    case STATE_IDLE:
        if (event == EVT_START_TX) {
            radio_start_tx();
            fsm->state = STATE_TX_ACTIVE;
            fsm->last_event_time = HAL_GetTick();
        }
        break;
    case STATE_TX_ACTIVE:
        if (event == EVT_TX_COMPLETE) {
            fsm->state = STATE_TX_DONE;
        } else if (event == EVT_TIMEOUT) {
            radio_abort();
            fsm->state = STATE_ERROR;
        }
        break;
    default:
        break;
    }
}
```

Nested state machine (HSM) for complex protocols:
- Parent state handles common events (reset, error)
- Child states handle protocol-specific transitions
- Entry/exit actions at each level

## Coding Standards (MISRA C inspired)

**Types**
- Use `stdint.h`: `uint8_t`, `uint16_t`, `uint32_t` - never `int`, `short`, `long`
- `volatile` for: memory-mapped registers, variables shared with ISRs
- `static` for file-local functions and variables
- `const` for immutable data to place in Flash (`.rodata`)

**Functions**
- Max 40 lines (embedded pragmatism - HSM needs room)
- Single return at the end (MISRA 15.5)
- No recursion (stack overflow risk)
- No dynamic memory after init (`malloc` only in setup, never in loop)

**Macros**
- Use `inline` over macros for functions
- Macros only for: hardware addresses, bit manipulation, compile-time constants
- Always parenthesize: `#define LED_PIN (1 << 5)`

**Safety**
- Watchdog: enable IWDG with reasonable timeout, kick in main loop
- HardFault handler: log stack pointer, PC, LR to retained RAM
- Assert: use `assert_param()` in DEBUG, strip in release
- Stack overflow: check with watermark pattern or MPU

## Interrupt Safety Recipe

```c
void USART1_IRQHandler(void) {
    if (USART1->SR & USART_SR_RXNE) {
        uint8_t byte = USART1->DR;
        rx_buffer_put(&rx_ring, byte);  // lock-free ring buffer
        task_notify_from_isr(rx_task_handle);
    }
    // NEVER: printf, delay, HAL_Delay, semaphore take
}
```

- ISR < 10μs - set a flag, wake a task, do real work in task context
- Shared variables: `volatile` + atomic access (LDREX/STREX or critical section)
- Disabled interrupts: measure max disable time, keep < 5μs
- Priority: protocol timing (highest), UI updates (lower), debug (lowest)

## Memory Budgets

- Stack: 4KB default, 8KB for RTOS + complex state machines
- Heap: ~0 (no `malloc` in production) or fixed pool at init
- `.bss` zero-init: check map file, no surprises
- Linker script: review `.ld` for stack/heap sizes, section alignment
- Flash wear: EEPROM emulation with wear leveling for persistent state

## Related

- Rule: `.cursor/rules/600-firmware.mdc` - Must / Must-not invariants
- Peripherals: `peripherals.md`
- Verification: `verification.md`
