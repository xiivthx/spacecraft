# Bare-metal OS class

On-demand recipe when OS class is bare-metal. Load after `core.md`. No RTOS primitives.

House layering stays `app/` -> `hal_if/` -> `drivers/` / `bsp/`. Shared invariants stay in `core.md`.

## Scheduling

Main loop plus ISR. Cooperative or superloop: poll flags, run one slice of work, return to the loop. No tasks, mutexes, or RTOS queues.

```c
for (;;) {
    wdt_kick();
    if (rx_flag) { rx_flag = 0; protocol_poll(); }
    app_tick();
}
```

## Interrupts

ISR sets a flag or writes a lock-free buffer. Work runs in the main loop. Keep ISR short per `core.md`.

## Timing

Hardware timers and the WDT. Kick the WDT in the main loop. Timeouts are timer compare or tick counts, not an RTOS delay.

## Ownership

One context owns each peripheral: the main loop plus that peripheral's ISR. Do not invent task-style split ownership.

## Related

- Shared core: `core.md`
- Rule: `.cursor/rules/600-firmware.mdc`
