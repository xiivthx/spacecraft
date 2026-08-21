# Spacecraft hard contract (Antigravity)

Rules are context; hooks are enforcement. Depth lives in skills - keep this file short.

## Hard gates (hooks)

- Secrets / force-push / catastrophic rm / main mutate / ship without `SPACECRAFT_SHIP=1`: `hooks/safety-check.mjs`
- `git push` is **denied in-agent** (no auto-push). Human pushes after AUTH.

## Soft contract

- **AUTH:** Quoted user authorization before outward push/deploy/publish/send.
- **INTENT:** Class + intended behavior before behavior-changing edits.
- **Commander:** No product code/tests - Task-delegate (`sc-coder` / `sc-tester` / `sc-firmware`; prose → `sc-writer`).
- **Lanes:** `/sc-discuss` → `/sc-run` → human check → `/sc-ship`. Small edits: `/sc-quick`.
- **SoT:** explicit user > approved draft + spec > DESIGN.md > process rules > evidence > code.
- **Language:** English technical substance; Thai for HIL / status / handoff.

See plugin skills for UX, TDD, judge, and domain encyclopedias.
