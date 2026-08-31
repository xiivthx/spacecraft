# RTOS OS class

On-demand recipe when OS class is RTOS. Load after `core.md`.

Do not force FreeRTOS, Zephyr, or cFS OSAL. Match the consuming project. cFS, F Prime, and Zephyr: match the consuming project; do not force those frameworks as the house default.

House layering stays `app/` -> `hal_if/` -> `drivers/` / `bsp/` unless the consuming project already uses those frameworks.

## Tasks and queues

Work lives in tasks. Queues (or the project's equivalent) carry events between ISR, producer tasks, and consumer tasks. Bound every queue. Prefer messages over shared mutable state.

## Interrupts

ISR defers to a task: set a flag, give from ISR, or post to a queue from ISR. Do the work in task context. Never take a blocking lock in ISR. Keep ISR short per `core.md`.

## Timing

RTOS tick plus hardware timers. Kick the WDT from a trusted task, not from ISR. Use the project's delay and timeout APIs in task context only.

## Ownership

One task owns each peripheral unless the project's RTOS already defines a driver model. When using house layering, wrap vendor HAL behind `hal_if/`. Synchronization (mutex or equivalent) only for data that two tasks must share; never take it in ISR.

## Related

- Shared core: `core.md`
- Rule: `.cursor/rules/600-firmware.mdc`
