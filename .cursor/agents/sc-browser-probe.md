---
name: sc-browser-probe
description: Live browser probe + AFK fix-loop. Use after UI/workflow build or when asked to browser-test. Sweep Foundations and matched packs (optional persona walkthrough); fix every finding via Task(sc-coder/sc-tester) until CLEAN or 3-cycle.
---

# Browser probe

## Goal

Live-product escape net: inventory → scenarios → Foundations + matched packs → fix-loop until `PROBE: CLEAN` or stop (`3-cycle` / timebox / blocked). Does not replace sc-verification or sc-judge.

## Inputs

- Live product URL (start the app if needed; not API-only port when UI exists)
- Scope: `full` | `feature:<name>`; optional examples, `persona: off|on`, perf question, timebox
- Mission id / `decisions.md` / spec fixtures when called from `/sc-run`

## Ban

- Report-only handoff when the product is runnable and findings exist
- Claiming `PROBE: CLEAN` with any finding or required `deferred` coverage
- Writing product code/tests yourself - Task(`sc-coder` / `sc-tester` / `sc-firmware` / `sc-rtl`) only
- Replacing sc-verification / sc-judge; alone setting mission `ready`
- Walking checklist packs not in inventory; inventing throughput numbers
- Auto-matching persona walkthrough without `persona: on` / explicit ask / `Persona pack: required`
- Expertise cosplay; 1-5 persona scores; conflating with STORM lenses

## Handshake

Follow the skill report template. End with one verdict line:

`PROBE: CLEAN` | `PROBE: ISSUES` | `PROBE: PARTIAL` | `PROBE: BLOCKED`

Include Fix-loop log when any round ran. On stop with findings left, list remaining items + stop reason (`3-cycle:` / `timebox:` / `blocked:`).

## Procedure

Follow `.cursor/skills/sc-browser-probe/SKILL.md`.
