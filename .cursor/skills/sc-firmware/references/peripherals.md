# Firmware peripheral recipes

On-demand register/driver examples for GPIO, UART, SPI, I2C, LCD, QSPI, LF/UHF, and DMA. Critical Must/Must-not stay in `.cursor/rules/610-firmware-peripherals.mdc`. Load this file when implementing peripheral drivers.

MCU-wide: BSP pin defines, CRC-checked bounded frames, SPI/I2C discipline, ISR enqueue-only, no GPIO busy-wait. STM32 LL/HAL snippets below are STM32-target examples. nRF uses nrfx (see target-nrf52840.md). On LTDC / DMA2D / Cortex-M7 targets, apply the LCD and D-cache notes in those sections; do not copy them onto nRF52840 or STM32L412.

## GPIO

STM32-target example (LL).

```c
#define LED_PORT  GPIOA
#define LED_PIN   LL_GPIO_PIN_5

void led_init(void) {
    LL_GPIO_SetPinMode(LED_PORT, LED_PIN, LL_GPIO_MODE_OUTPUT);
    LL_GPIO_SetPinOutputType(LED_PORT, LED_PIN, LL_GPIO_OUTPUT_PUSHPULL);
    LL_GPIO_SetPinSpeed(LED_PORT, LED_PIN, LL_GPIO_SPEED_FREQ_LOW);
}

void led_on(void)  { LED_PORT->BSRR = LED_PIN; }
void led_off(void) { LED_PORT->BSRR = (LED_PIN << 16); }
void led_toggle(void) { LED_PORT->ODR ^= LED_PIN; }
```

- Always use BSP defines, never raw port/pin numbers
- Debounce: 5-20ms in software, or RC filter in hardware
- Pull-up/down: enable internal when input - floats = noise = random wakeups

## Cycle-accurate GPIO / RF bitbang

Main-loop "step once per call" DROPS RF cycles when the loop cannot keep up. For bitbang TX: spin-wait each RF cycle index in one burst (or drive from a hardware timer) until the frame completes - do not yield mid-frame to a slow main loop.

Emit a proof counter (e.g. pause-high cycle count) on UART for HIL so host tests can assert physical work, not only app state.

**Must not:** EXTI on every edge of a continuous multi-MHz / subcarrier line (IRQ storm). Prefer TIM capture, DMA, or a polled sample window.

## UART / Serial Protocol

Frame format is MCU-wide. Tick helper in the snippet is an STM32-target example.

```
Frame format: [SYNC] [LEN] [TYPE] [PAYLOAD...] [CRC16]
              1B      1B     1B     0-255B      2B
```

```c
typedef struct {
    uint8_t  rx_buf[256];
    uint8_t  rx_head, rx_tail;
    uint32_t last_byte_ticks;
} uart_rx_t;

void uart_rx_isr(uint8_t byte) {
    if (ticks_since(rx.last_byte_ticks) > FRAME_TIMEOUT_MS) {
        rx_head = 0;
    }
    rx_buf[rx_head++] = byte;
    rx.last_byte_ticks = HAL_GetTick();
    if (rx_head >= rx_buf[rx_head - 1] + 4) {
        protocol_dispatch(rx_buf, rx_head);
    }
}
```

- Timeout: 100ms inter-byte = reset parser
- Max payload: bounded, reject oversized frames
- State machine per protocol instance - never process in ISR

## SPI Bus

STM32-target example (HAL). nRF uses nrfx (see target-nrf52840.md).

```
Modes: 0 (CPOL=0,CPHA=0) - most common, 1, 2, 3
Speed: 100kHz-50MHz (peripheral-dependent)
Pins: SCK, MOSI, MISO, NSS (CS)
```

```c
void spi_write_reg(uint8_t addr, uint8_t value) {
    CS_LOW();
    HAL_SPI_Transmit(&hspi1, &addr, 1, HAL_MAX_DELAY);
    HAL_SPI_Transmit(&hspi1, &value, 1, HAL_MAX_DELAY);
    CS_HIGH();
}

uint8_t spi_read_reg(uint8_t addr) {
    uint8_t tx = addr | 0x80;
    uint8_t rx = 0;
    CS_LOW();
    HAL_SPI_Transmit(&hspi1, &tx, 1, HAL_MAX_DELAY);
    HAL_SPI_Receive(&hspi1, &rx, 1, HAL_MAX_DELAY);
    CS_HIGH();
    return rx;
}

typedef struct {
    SPI_HandleTypeDef* hspi;
    DMA_HandleTypeDef* hdma_tx;
    DMA_HandleTypeDef* hdma_rx;
    uint8_t*          tx_buf;
    uint8_t*          rx_buf;
    volatile bool     busy;
    volatile bool     error;
    volatile bool     rx_complete;
} spi_dma_t;

void spi_dma_transfer(spi_dma_t* dev, uint8_t* tx, uint8_t* rx, uint16_t len) {
    dev->busy = true;
    if (rx) {
        HAL_SPI_TransmitReceive_DMA(dev->hspi, tx, rx, len);
    } else {
        HAL_SPI_Transmit_DMA(dev->hspi, tx, len);
    }
}

void spi_dma_full_duplex(spi_dma_t* dev, uint8_t* tx, uint8_t* rx, uint16_t len) {
    dev->rx_complete = false;
    HAL_SPI_TransmitReceive_DMA(dev->hspi, tx, rx, len);
}

void HAL_SPI_TxCpltCallback(SPI_HandleTypeDef* hspi) {
    spi_dev.busy = false;
    if (spi_dev.on_tx_done) spi_dev.on_tx_done();
}
void HAL_SPI_ErrorCallback(SPI_HandleTypeDef* hspi) {
    spi_dev.error = true;
    spi_dev.busy = false;
}
```

- Chip select timing: CS low → 1μs → transfer → 1μs → CS high
- Multi-byte: keep CS low across entire burst, raise only at end
- Clock polarity/phase: read datasheet - wrong mode = silent failure
- Speed: start at 1MHz, test, increase. Max = min(CPU/2, peripheral max)
- DMA priority: SPI > UART (bulk data), round-robin if multiple SPI ports
- CS pin: hardware NSS preferred; if GPIO, use BSRR for atomic toggle

## I2C Bus

STM32-target example (HAL).

```
Speed: Standard 100kHz, Fast 400kHz, Fast+ 1MHz
Pins: SCL, SDA (open-drain - requires external pull-up 2.2-4.7kΩ)
Addressing: 7-bit (most common) or 10-bit
```

```c
typedef struct {
    I2C_HandleTypeDef* hi2c;
    uint8_t  addr;
    uint32_t timeout_ms;
    volatile bool busy;
    volatile bool error;
    uint8_t  tx_buf[32];
    uint8_t  rx_buf[32];
} i2c_dev_t;

uint8_t i2c_reg_read(i2c_dev_t* dev, uint8_t reg) {
    uint8_t val = 0;
    HAL_I2C_Mem_Read(dev->hi2c, dev->addr << 1, reg,
                     I2C_MEMADD_SIZE_8BIT, &val, 1, dev->timeout_ms);
    return val;
}

bool i2c_burst_read(i2c_dev_t* dev, uint8_t reg, uint8_t* buf, uint16_t len) {
    return HAL_I2C_Mem_Read(dev->hi2c, dev->addr << 1, reg,
                            I2C_MEMADD_SIZE_8BIT, buf, len, dev->timeout_ms) == HAL_OK;
}

void i2c_read_dma(i2c_dev_t* dev, uint8_t reg, uint16_t len) {
    dev->busy = true;
    HAL_I2C_Mem_Read_DMA(dev->hi2c, dev->addr << 1, reg,
                         I2C_MEMADD_SIZE_8BIT, dev->rx_buf, len);
}

bool i2c_write_regs(i2c_dev_t* dev, uint8_t reg, uint8_t* data, uint16_t len) {
    return HAL_I2C_Mem_Write(dev->hi2c, dev->addr << 1, reg,
                             I2C_MEMADD_SIZE_8BIT, data, len, dev->timeout_ms) == HAL_OK;
}

void HAL_I2C_MemTxCpltCallback(I2C_HandleTypeDef* hi2c) {
    i2c_dev.busy = false;
}
void HAL_I2C_ErrorCallback(I2C_HandleTypeDef* hi2c) {
    uint32_t err = HAL_I2C_GetError(hi2c);
    if (err & HAL_I2C_ERROR_AF)  { /* NACK */ }
    if (err & HAL_I2C_ERROR_BERR) { /* bus error */ }
    if (err & HAL_I2C_ERROR_ARLO) { /* arbitration lost */ }
    i2c_dev.error = true;
    i2c_dev.busy = false;
}

void i2c_scan(I2C_HandleTypeDef* hi2c) {
    for (uint8_t addr = 1; addr < 127; addr++) {
        if (HAL_I2C_IsDeviceReady(hi2c, addr << 1, 1, 10) == HAL_OK) {
            printf("Device at 0x%02X\n", addr);
        }
    }
}
```

- **Pull-up resistors**: mandatory - 2.2k-4.7kΩ on both SCL and SDA
- **Bus capacitance**: max 400pF at 100kHz, 200pF at 400kHz
- **NACK handling**: retry with backoff when device busy
- **Clock stretching**: enable in I2C init if device needs it
- **Multi-master**: prefer single master unless necessary
- **Timeout**: always set per-device timeout
- **DMA for bulk**: EEPROM pages, sensor FIFO dumps
- **Bus recovery**: if SDA stuck low → toggle SCL 9 times; reset peripheral after

### Common I2C Devices

| Device | Typical addr | Notes |
|---|---|---|
| EEPROM (24Cxx) | 0x50-0x57 | 5ms write cycle, poll ACK after write |
| Temperature (LM75, TMP102) | 0x48-0x4F | 12-bit, conversion time 30-300ms |
| OLED (SSD1306) | 0x3C/0x3D | 128x64, 1KB framebuffer |
| RTC (DS3231) | 0x68 | BCD format, temperature-compensated |
| GPIO expander (MCP23017) | 0x20-0x27 | 16 GPIO, 2 register banks |
| Accelerometer (MPU6050) | 0x68/0x69 | 6-axis, FIFO buffer, interrupt pin |

## LCD Display (LTDC / DMA2D targets)

On LTDC / DMA2D / Cortex-M7 targets only. STM32-target example below is STM32F746-Discovery. Do not apply these notes to nRF52840 or STM32L412.

```
Layer: app/ → hal_if/lcd.h → drivers/lcd/ltdc_stm32f7.c → bsp/board.c

hal_if/lcd.h:
  void lcd_init(void);
  void lcd_set_pixel(uint16_t x, uint16_t y, uint16_t color);
  void lcd_fill(uint16_t color);
  void lcd_draw_rect(uint16_t x, uint16_t y, uint16_t w, uint16_t h, uint16_t color);
  void lcd_draw_bitmap(uint16_t x, uint16_t y, uint8_t* data, uint16_t w, uint16_t h);
  void lcd_flush(void);
```

```c
__attribute__((section(".lcd_fb"))) uint16_t lcd_fb[2][480*272];

void lcd_fill(uint16_t color) {
    DMA2D->CR = 0;
    DMA2D->OPFCCR = DMA2D_OUTPUT_RGB565;
    DMA2D->OOR = 0;
    DMA2D->OMAR = (uint32_t)&lcd_fb[back_buf][0];
    DMA2D->NLR = (272 << 16) | 480;
    DMA2D->OCOLR = color;
    DMA2D->CR = DMA2D_CR_START;
    while (DMA2D->CR & DMA2D_CR_START);
}

void lcd_draw_bitmap(uint16_t x, uint16_t y, uint8_t* src, uint16_t w, uint16_t h) {
    DMA2D->CR = DMA2D_CR_MODE_M2M_PFC;
    DMA2D->FGPFCCR = DMA2D_INPUT_ARGB8888;
    DMA2D->OPFCCR = DMA2D_OUTPUT_RGB565;
    DMA2D->FGMAR = (uint32_t)src;
    DMA2D->OMAR = (uint32_t)&lcd_fb[back_buf][y * 480 + x];
    DMA2D->FGOR = 0;
    DMA2D->OOR = 480 - w;
    DMA2D->NLR = (h << 16) | w;
    DMA2D->CR = DMA2D_CR_START | DMA2D_CR_MODE;
    while (DMA2D->CR & DMA2D_CR_START);
}

void lcd_blend_layer(uint8_t* fg, uint8_t* bg, uint16_t alpha) {
    DMA2D->CR = DMA2D_CR_MODE_M2M_BLEND;
    DMA2D->FGPFCCR = DMA2D_INPUT_ARGB8888 | (alpha << 24);
    DMA2D->BGPFCCR = DMA2D_INPUT_ARGB8888;
    DMA2D->OPFCCR = DMA2D_OUTPUT_RGB565;
    DMA2D->FGMAR = (uint32_t)fg;
    DMA2D->BGMAR = (uint32_t)bg;
    DMA2D->OMAR = (uint32_t)lcd_fb[back_buf];
    DMA2D->NLR = (272 << 16) | 480;
    DMA2D->CR = DMA2D_CR_START;
    while (DMA2D->CR & DMA2D_CR_START);
}

void LTDC_IRQHandler(void) {
    if (LTDC->ISR & LTDC_ISR_LIF) {
        LTDC->ISR = LTDC_ISR_LIF;
        LTDC->L2CR ^= LTDC_L2CR_VIS;
        front_buf ^= 1; back_buf ^= 1;
        LTDC->L2CFBAR = (uint32_t)lcd_fb[front_buf];
    }
}
```

On LTDC / DMA2D / Cortex-M7 targets:

- **Double-buffered**: draw to back buffer, LTDC reads front, swap on vblank
- **DMA2D for fill/blit/blend**: color convert - never CPU pixel loops
- **SDRAM bandwidth**: LTDC @ 9MHz pixel clock ≈ 18 MB/s - leave headroom for DMA2D
- **Layer 1**: background. **Layer 2**: foreground. Blend with alpha
- **Pixel format**: RGB565 native; ARGB8888 for images with DMA2D PFC
- **Cache**: framebuffer write-through SDRAM (MPU) or `cache_clean()` before LTDC reload
- **Touch**: FT5336 I2C at 0x38 - poll at 20-30Hz in low-priority task

## QSPI Flash (STM32-target example - 128Mbit on Discovery)

```c
void qspi_init(void) {
    QSPI_CommandTypeDef cmd = {
        .InstructionMode   = QSPI_INSTRUCTION_1_LINE,
        .AddressMode       = QSPI_ADDRESS_4_LINES,
        .DataMode          = QSPI_DATA_4_LINES,
        .DummyCycles       = 6,
        .SIOOMode          = QSPI_SIOO_INST_ONLY_FIRST_CMD,
    };
    HAL_QSPI_Init(&hqspi);
    QSPI->CR |= QSPI_CR_TCEN;
    HAL_QSPI_MemoryMapped(&hqspi, &cmd, NULL);
}

__attribute__((section(".qspi"))) const uint8_t splash_image[480*272*2];

void qspi_write_page(uint32_t addr, uint8_t* data, uint16_t len) {
    qspi_write_enable();
    qspi_sector_erase(addr);
    while (qspi_busy());
    qspi_write_enable();
    qspi_page_program(addr, data, len);
    while (qspi_busy());
}
```

- QSPI is **read-only** in production - write only during firmware update
- Memory-mapped reads are cached (via MPU)
- Erase before write: 4KB sector erase (~45ms), 256B page program (~0.5ms)
- Wear leveling: ≈ 100K erase cycles - not for frequently-updated data
- Asset pipeline: PNG → `xxd -i` → link to `.qspi` section

## LF Protocol (125-134 kHz)

```
Carrier: 125 kHz (LF)
Modulation: ASK (100% depth typical)
Data rate: 1-8 kbps
Encoding: Manchester, Biphase, or NRZ
```

```c
typedef enum {
    LF_IDLE, LF_CARRIER_ON, LF_BIT_SAMPLE, LF_FRAME_END
} lf_rx_state_t;

void lf_timer_capture_isr(uint16_t capture_val) {
    uint16_t width = capture_val - last_capture;
    last_capture = capture_val;

    if (width > 6 && width < 10) {
        lf_shift_bit(1);
    } else if (width > 10 && width < 16) {
        lf_shift_bit(0);
    } else {
        lf_reset();
    }
}
```

Architecture:
- HW timer capture → ISR measures pulse width → ring buffer → task processes bits
- Error handling: preamble detection, CRC16, bit timing tolerance ±25%
- Modulation: PWM or OOK via GPIO + timer, or dedicated LF driver IC

## UHF Protocol (433/868/915 MHz)

```
Carrier: 433/868/915 MHz
Modulation: FSK, GFSK, OOK
Typical IC: CC1101, SX1276, AX5043
```

```c
typedef struct {
    radio_state_t state;
    uint8_t tx_power;
    uint32_t frequency;
    uint32_t bitrate;
    uint8_t sync_word[4];
    uint8_t tx_buf[64];
    uint8_t tx_len;
    uint8_t rx_buf[64];
    volatile bool packet_ready;
} uhf_radio_t;

void uhf_tx_packet(uhf_radio_t* radio, uint8_t* data, uint8_t len) {
    // 1. Set TX mode - SPI write register
    // 2. Write FIFO - SPI burst via DMA
    // 3. Issue TX command - SPI strobe
    // 4. Wait for GDO0 pin IRQ (packet sent)
}
```

- Always wait for GDO pin (packet ready) - never poll register
- Preamble: ≥4 bytes for AGC settling
- CRC: hardware CRC on CC1101, software CRC16-CCITT for custom protocol
- Retry: max 3 retransmits with exponential backoff (50ms, 100ms, 200ms)
- RSSI/LQI: reject packets below threshold

## DMA Best Practices

Circular mode is MCU-wide for continuous streams. On LTDC / DMA2D / Cortex-M7 targets: flush/invalidate D-cache around DMA. STM32-target example (HAL ADC callbacks):

```c
#define DMA_BUF_SIZE 256
static uint16_t dma_buf[2][DMA_BUF_SIZE];

void HAL_ADC_ConvHalfCpltCallback(ADC_HandleTypeDef* hadc) {
    process_adc_data(dma_buf[0], DMA_BUF_SIZE);
}
void HAL_ADC_ConvCpltCallback(ADC_HandleTypeDef* hadc) {
    process_adc_data(dma_buf[1], DMA_BUF_SIZE);
}
```

- DMA priority: highest for audio/ADC, medium for SPI/UART, low for memory
- Cache coherency: on Cortex-M7 targets, flush/invalidate D-cache before/after DMA
- Circular mode: for continuous streaming - never stop/restart mid-stream
- Error handling: DMA transfer error ISR → log error, reset peripheral

## Related

- Rule: `.cursor/rules/610-firmware-peripherals.mdc` - Must / Must-not invariants
- Core: `core.md`
- Target refs: `target-cortex-m.md`, `target-stm32.md`, `target-nrf52840.md`
- Verification: `verification.md`
