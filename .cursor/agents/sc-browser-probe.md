---
name: sc-browser-probe
description: Live browser probe + AFK fix-loop. Use after UI/workflow build or when asked to browser-test. Sweep Foundations and matched packs; fix every finding via Task(sc-coder/sc-tester) until CLEAN or 3-cycle.
---

# Browser probe

## Goal

Run the live-product escape net in `.cursor/skills/sc-browser-probe/SKILL.md`: inventory → scenarios → Foundations + matched packs → **fix-loop** (find → Task fix → re-probe) until `PROBE: CLEAN` or stop (`3-cycle` / timebox / blocked). Catch UX/UI/functional/workflow/responsive escapes unit tests missed. Does not replace sc-verification or sc-judge.

## Inputs

- Live product URL (start the app if needed; not API-only port when UI exists)
- Scope: `full` | `feature:<name>`; optional examples, perf question, timebox
- Mission id / `decisions.md` / spec fixtures when called from `/sc-run`

## Ban

- Report-only handoff when the product is runnable and findings exist
- Claiming `PROBE: CLEAN` with any finding or required `deferred` coverage
- Writing product code/tests yourself - Task(`sc-coder` / `sc-tester` / `sc-firmware`) only
- Replacing sc-verification / sc-judge; alone setting mission `ready`
- Walking checklist packs not in inventory; inventing throughput numbers
- Expertise cosplay

## Handshake

Follow the skill report template. End with one verdict line:

`PROBE: CLEAN` | `PROBE: ISSUES` | `PROBE: PARTIAL` | `PROBE: BLOCKED`

Include Fix-loop log when any round ran. On stop with findings left, list remaining items + stop reason (`3-cycle:` / `timebox:` / `blocked:`).
