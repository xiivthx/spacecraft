---
name: sc-verification
description: Capture fresh command evidence before claiming Spacecraft work is complete
license: MIT
compatibility: opencode
---
- No done/pass/verified/ready claim without evidence.
- Use `node scripts/spacecraft.mjs evidence "<label>" -- <command>`.
- Capture failures too.
- Map acceptance checks to evidence ids in final summaries.
- If a check cannot be automated, state why and mark it manual.
- Prefer focused verification first, then broader build/test checks before shipping.
