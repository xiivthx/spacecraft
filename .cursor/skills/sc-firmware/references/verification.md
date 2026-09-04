# Firmware verification recipes

On-demand Unity/Ceedling, target integration, HIL, and CI recipes. Critical Must/Must-not stay in `.cursor/rules/620-firmware-testing.mdc`. Load this file when writing or running firmware tests.

## Three-Layer Testing

```
Layer 1: Host Unit Tests (PC)
  → Ceedling/Unity, CMock, CppUTest
  → Mock vendor HAL/SDK headers, mock drivers, test pure logic
  → Target: 80%+ coverage of app/ and hal_if/

Layer 2: Target Integration (on-device)
  → Flash firmware, run test suite via UART/SWO
  → Test hardware interaction: GPIO toggles, SPI transfers, IRQ latency
  → Target: every peripheral init + basic I/O verified

Layer 3: HIL (Hardware-in-the-Loop)
  → Automated test rig: inject signals, measure outputs
  → Protocol conformance: inject known-good/bad packets, verify response
  → Target: smoke test before release
```

## Host Unit Testing

```c
#include "unity.h"
#include "mock_gpio.h"
#include "mock_spi.h"

static radio_state_t state;

void setUp(void) {
    state = (radio_state_t){0};
}

void test_radio_tx_packet_valid(void) {
    spi_send_cmd_ExpectAndReturn(CMD_TX, STATUS_READY);
    gpio_set_Expect(LED_TX_PIN);

    radio_send_packet(&state, test_packet, 8);

    TEST_ASSERT_EQUAL(STATE_TX_DONE, state.current);
}

void test_radio_tx_no_ack_timeout(void) {
    spi_send_cmd_ExpectAndReturn(CMD_TX, STATUS_TIMEOUT);
    gpio_set_Expect(LED_ERR_PIN);

    radio_send_packet(&state, test_packet, 8);

    TEST_ASSERT_EQUAL(STATE_ERROR, state.current);
}
```

- Mock ALL hardware interfaces - never include real vendor HAL/SDK headers in host unit tests
- Test error paths: timeout, CRC fail, buffer overflow, protocol violations
- Test state machine transitions: every event in every state
- Test boundary values: `UINT8_MAX`, `0`, `packet_len - 1`, `packet_len + 1`

## Target Integration Testing

```c
void test_uart_loopback(void) {
    uart_init(USART2, 115200);
    uint8_t tx[] = "Hello";
    uint8_t rx[5] = {0};

    uart_send(USART2, tx, 5);
    HAL_Delay(10);
    uart_receive(USART2, rx, 5);

    TEST_ASSERT_EQUAL_MEMORY(tx, rx, 5);
}

void test_irq_latency(void) {
    DWT->CYCCNT = 0;
    trigger_irq(TIM2_IRQn);
    uint32_t cycles = DWT->CYCCNT;
    uint32_t us = cycles / (SystemCoreClock / 1000000);

    TEST_ASSERT_LESS_THAN(10, us);  // ISR must complete < 10μs
}
```

## HIL Testing

```
Test rig:
  Host PC ←UART→ Test Board ←GPIO/SPI→ Device Under Test

Test flow:
  1. Host sends command via UART to test board
  2. Test board injects signal / sends packet to DUT
  3. Test board measures DUT response (GPIO, SPI capture, protocol packet)
  4. Test board reports pass/fail back to host
```

Automate with:
- `pytest` (Python) for test orchestration
- `OpenOCD` + GDB for target flashing + debug
- `expect` / `pexpect` for UART interaction
- CI: GitHub Actions with self-hosted runner (board connected)

**Proof oracle.** Host exact-line tests + target UART lines that prove physical work (counts, timing windows, pin IDR) - not only app-state asserts.

**Dual-board bench.** FPGA + MCU: correlate both UARTs / LEDs before changing either side.

- Must: Peer DUT harness requires dual-DUT correlated evidence for HIL GREEN or ready
- Must not: Single-DUT evidence when peer harness exists

**Host green ≠ silicon green.** Passing host unit tests does not claim HIL or target integration.

## CI Pipeline

When claiming flight-grade: run static analysis (cppcheck or clang-tidy) with warnings-as-errors. That pass is separate from host unit, target integration, and HIL.

```yaml
build:
  - arm-none-eabi-gcc build → verify 0 warnings (warnings-as-errors)
  - cppcheck / clang-tidy → static analysis
  - cmocka/ceedling → host unit tests

flash:
  - openocd → flash to target board
  - uart console → verify boot message + self-test pass

hil:
  - pytest → run HIL suite
  - collect results → pass/fail report
```

## Debug Tools

- **DWT cycle counter**: `DWT->CYCCNT` - measure execution time
- **SWO trace**: `ITM_SendChar()` - printf to debugger without UART
- **FreeRTOS trace**: `uxTaskGetStackHighWaterMark()` - detect stack overflow
- **HardFault analyzer**: read stacked PC, LR, PSR - decode fault address
- **Logic analyzer**: Saleae/PulseView - verify SPI timing, protocol bit patterns

## Related

- Rule: `.cursor/rules/620-firmware-testing.mdc` - Must / Must-not invariants
- Core: `core.md`
- Target refs: `target-cortex-m.md`, `target-stm32.md`, `target-nrf52840.md`
- Peripherals: `peripherals.md`
