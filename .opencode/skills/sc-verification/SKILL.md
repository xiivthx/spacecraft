---
name: sc-verification
description: Capture fresh command evidence before claiming Spacecraft work is complete
license: MIT
compatibility: opencode
---
- No done/pass/verified/ready claim without evidence.
- Resolve the mission with `scripts/spacecraft resolve --json`; `.space/current` is fallback state, not sole authority.
- If resolver safety is not `safe`, stop before evidence capture and ask the user to choose with `scripts/spacecraft use <number|id|title>` or an explicit `SPACECRAFT_MISSION`.
- Use `scripts/spacecraft evidence "<label>" -- <command>`.
- Capture failures too.
- Map acceptance checks to evidence ids in final summaries.
- If a check cannot be automated, state why and mark it manual.
- Prefer focused verification first, then broader build/test checks before shipping.
- Use rtk for noisy verification commands when available. Use raw output or `rtk proxy` passthrough when exact evidence is needed.
