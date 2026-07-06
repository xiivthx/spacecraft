---
name: sc-web-service
description: Build a lean local web service from scratch under Spacecraft mission control
license: MIT
compatibility: opencode
---
- Use only when developing a web service.
- If no stack is specified, default to Node.js + TypeScript + Fastify + Vitest.
- The first milestone should be minimal:
  - package scripts: `dev`, `test`, `build`
  - `GET /healthz` returns `{ "ok": true }`
  - `GET /version` returns service metadata
  - tests pass
  - build passes
- Installing product dependencies is allowed only after user approval.
- Do not add database, auth, Docker, deployment, queue, frontend, or observability stack unless explicitly requested.
- Prefer small vertical slices.
- Keep the web service separate from Spacecraft harness logic.
- If a web service includes pages or UI, use sc-design and `DESIGN.md`.
- If mood, theme, or art direction is unclear for a web UI, use sc-design before planning UI tasks.
- Default first UI surface, when requested, should be minimal:
  - one home/status page
  - one health/status section
  - one version/build metadata section
  - no fake metrics
  - no marketing fluff
  - no generic AI landing-page hero
- Keep backend service logic separate from visual components.
- Do not add frontend complexity unless explicitly requested.
