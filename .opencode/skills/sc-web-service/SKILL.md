---
name: sc-web-service
description: Build a lean local web service from scratch under mission control. Activate on "build web service", "create API", "scaffold server", or new web project.
license: MIT
compatibility: opencode
metadata:
  version: 1.0
  audience: xiivthx
---

# sc-web-service

Build a lean local web service from scratch under mission control. Default stack: Node.js + TypeScript + Fastify + Vitest. Alternative stacks require user approval.

## When to use

Activate when the user asks to:

- **"Build a web service" / "create an API" / "scaffold a server"** — new web project
- **"Add a health endpoint" / "add version endpoint"** — minimal milestone additions
- When a mission task requires a local backend service

## Workflow

Use this exact sequence unless the user specifies otherwise:

1. **Resolve mission** — `scripts/spacecraft resolve --json`. Block if safety ≠ `safe`.

2. **Choose stack** — Default: Node.js + TypeScript + Fastify + Vitest. If the user specifies a different stack, adapt accordingly. Run `spacecraft research "fastify v5 typescript setup"` before scaffolding if versions are uncertain. Ask before installing any `npm` dependencies — list the packages and wait for approval.

3. **Scaffold** — Create the project with these exact files:
   - `package.json` — scripts: `dev`, `test`, `build`, `start`
   - `tsconfig.json` — strict mode, ESNext target
   - `src/server.ts` — Fastify instance with `GET /healthz` and `GET /version`
   - `src/server.test.ts` — tests for both endpoints
   - `.gitignore` — `node_modules/`, `dist/`, `.env`

   ```bash
   npm install    # after user approves dependencies
   npm test       # verify tests pass
   npm run build  # verify build passes
   ```

4. **Verify** — `scripts/spacecraft evidence "web:endpoints" -- npm test`. Both endpoints must return 200. Build must succeed.

5. **Iterate** — Prefer small vertical slices. Keep the web service separate from harness logic. Each new endpoint gets a test first (see sc-tdd).

### Edge cases

- **User specifies a stack** — Use that stack. Still require minimal milestone with health/version endpoints.
- **Project already exists** — Don't re-scaffold. Add endpoints to existing structure.
- **Tests fail** — Fix before proceeding. Never skip verification.
- **User wants database/auth/Docker** — These are out of scope. Remind the user this skill is for lean services only.

## Rules

- **Must**: Resolve mission before mutating work.
- **Must**: Default to Node.js + TypeScript + Fastify + Vitest when no stack is specified.
- **Must**: First milestone is always: `dev`/`test`/`build` scripts, `GET /healthz`, `GET /version`, passing tests, passing build.
- **Must**: Verify with `scripts/spacecraft evidence` after each milestone.
- **Must**: Prefer small vertical slices over broad horizontal scaffolding.
- **Ask before**: Installing any npm dependencies. List packages first, get approval.
- **Must not**: Add database, auth, Docker, deployment, queues, frontend, or observability unless explicitly requested.
- **Must not**: Add fake metrics, marketing copy, or AI-generated landing pages to status endpoints.

## Out of scope

- UI design or frontend architecture — use sc-design
- Database, auth, deployment, observability — separate concerns, ask before adding
- General mission workflow — use sc-mission
- TDD workflow — use sc-tdd

## Output format

```
Stack: Node.js + TypeScript + Fastify + Vitest
Scaffold:
  package.json ✓ (dev, test, build, start)
  tsconfig.json ✓ (strict, ESNext)
  src/server.ts ✓ (GET /healthz, GET /version)
  src/server.test.ts ✓ (2 endpoint tests)
Verify:
  npm test → PASS
  npm run build → PASS
Evidence: web:endpoints
```

## Checklist

- [ ] Mission resolved, branch created
- [ ] Stack chosen, dependencies approved by user
- [ ] `package.json` has `dev`, `test`, `build`, `start` scripts
- [ ] `GET /healthz` returns `{ "ok": true }` with 200
- [ ] `GET /version` returns `{ "version": "...", "build": "..." }` with 200
- [ ] Tests pass (`npm test`)
- [ ] Build passes (`npm run build`)
- [ ] Evidence captured with `scripts/spacecraft evidence`
- [ ] No unapproved dependencies

## Research auto-trigger

When the default stack version is uncertain or the user specifies an unfamiliar framework, run `spacecraft research "<framework> setup guide"` before scaffolding.

---

## References

- `scripts/spacecraft --help` — mission resolver and evidence capture
- Fastify v5 docs: `https://fastify.dev/docs/latest/` (check before scaffolding for API changes)
- Vitest docs: `https://vitest.dev/` (check for config syntax)
