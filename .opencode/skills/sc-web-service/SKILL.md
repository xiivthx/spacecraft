---
name: sc-web-service
description: Build a lean local web service from scratch under mission control
license: MIT
compatibility: opencode
metadata:
  version: 1.0
  audience: xiivthx
---

# sc-web-service

Build a lean local web service from scratch under mission control.

## When to use

Activate when the user asks to:

- Create a new web service or API
- Scaffold a local backend project
- Add a health check or version endpoint

## Workflow

Use this exact sequence unless the user specifies otherwise:

1. **Choose stack** — If no stack is specified, default to Node.js + TypeScript + Fastify + Vitest. Get user approval before installing product dependencies.
2. **Scaffold minimal milestone** — Create package scripts: `dev`, `test`, `build`. Include `GET /healthz` returning `{ "ok": true }` and `GET /version` returning service metadata.
3. **Verify** — Ensure tests pass and build passes before moving on.
4. **Iterate** — Prefer small vertical slices. Keep the web service separate from harness logic.

## Rules

- **Must**: Use only when developing a web service.
- **Must**: If no stack is specified, default to Node.js + TypeScript + Fastify + Vitest.
- **Must**: The first milestone should be minimal:
  - package scripts: `dev`, `test`, `build`
  - `GET /healthz` returns `{ "ok": true }`
  - `GET /version` returns service metadata
  - tests pass
  - build passes
- **Ask before**: Installing product dependencies is allowed only after user approval.
- **Must not**: Add database, auth, Docker, deployment, queue, frontend, or observability stack unless explicitly requested.
- **Must**: Prefer small vertical slices.
- **Must**: Keep the web service separate from harness logic.
- **Must**: If a web service includes pages or UI, use sc-design before planning UI tasks.
- **Must**: If mood, theme, or art direction is unclear for a web UI, use sc-design before planning UI tasks.
- **Must**: Default first UI surface, when requested, should be minimal:
  - one home/status page
  - one health/status section
  - one version/build metadata section
  - no fake metrics
  - no marketing fluff
  - no generic AI landing-page hero
- **Must**: Keep backend service logic separate from visual components.
- **Must not**: Add frontend complexity unless explicitly requested.

## Out of scope

This skill does NOT handle:

- UI design or frontend architecture — use sc-design
- Database, auth, deployment, or observability — ask the user before adding these
- General mission work — use sc-mission

## Output format

```
package.json scripts:
  dev — development server
  test — test runner
  build — production build

API endpoints:
  GET /healthz -> { "ok": true }
  GET /version -> { "version": "...", "build": "..." }
```

## Checklist

Before claiming the web service is ready:

- [ ] Stack chosen (default: Node.js + TypeScript + Fastify + Vitest)
- [ ] Minimal milestone: `dev`, `test`, `build` scripts work
- [ ] `GET /healthz` and `GET /version` endpoints respond
- [ ] Tests pass
- [ ] Build passes
- [ ] No unapproved dependencies added

---

## References

- `scripts/spacecraft --help` — spacecraft CLI reference
- Fastify documentation (when using default stack)
