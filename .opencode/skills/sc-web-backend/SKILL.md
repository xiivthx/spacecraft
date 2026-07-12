---
name: sc-web-backend
description: Build server APIs and services with Node.js, TypeScript, Fastify, and Vitest. Activate on "create an API endpoint", "build a web service", "scaffold a server", "add authentication", or "design REST API".
license: MIT
compatibility: opencode
metadata:
  version: 1.0
  audience: xiivthx
---

# sc-web-backend

Build server-side APIs and services under mission control. Default stack: Node.js + TypeScript + Fastify + Vitest.

## When to use

Activate when the user asks to:

- **"Create an API endpoint" / "build a web service"** — new server features
- **"Scaffold a server" / "set up Fastify"** — project initialization
- **"Add authentication" / "design REST API"** — API architecture and security
- **"Add a health endpoint" / "add version endpoint"** — minimal service milestones
- When a mission task requires server-side implementation

## Workflow

Use this exact sequence unless the user specifies otherwise:

1. **Resolve mission** — `scripts/spacecraft resolve --json`. Block if safety ≠ `safe`.

2. **Choose stack** — Default: Node.js + TypeScript + Fastify + Vitest. If the user specifies a different stack, adapt accordingly. Run `spacecraft research "fastify v5 typescript setup"` before scaffolding if versions are uncertain.

3. **Scaffold (new project)** — If no project exists:
   - `package.json` — scripts: `dev`, `test`, `build`, `start`
   - `tsconfig.json` — strict mode, ESNext target
   - `src/server.ts` — Fastify instance with `GET /healthz` returning `{ "ok": true }` and `GET /version` returning `{ "version": "...", "build": "..." }`
   - `src/server.test.ts` — tests for both endpoints
   - `.gitignore` — `node_modules/`, `dist/`, `.env`
   Wait for user approval before installing npm dependencies.

4. **Build by slice** — Implement one vertical feature slice at a time:
   - Route definition with Fastify schema validation
   - Handler logic with typed request/response
   - Service layer (business logic, separated from route handlers)
   - Test for the endpoint behavior via Vitest + `fastify.inject()`
   Prefer small, focused routes. Extract shared patterns to `references/api-patterns.md`.

5. **Verify** — `scripts/spacecraft evidence "<label>" -- npm test`. Tests must pass. Build must succeed (`npm run build`). Run the full test suite after each slice.

### Edge cases

- **Project already exists** — Don't re-scaffold. Add routes to the existing structure.
- **User requests auth** — Implement via Fastify plugin. Use bearer tokens or session cookies. Hash passwords with bcrypt. Store secrets in environment variables.
- **Database access needed** — This is a separate concern. Use the project's existing data layer. For schema/migration questions, defer to the appropriate concern.
- **GraphQL requested** — Adapt route layer to resolvers. Same principles: typed schemas, tested behavior.
- **Build or tests fail** — Fix before proceeding. Never skip verification.

## Rules

- **Must**: Resolve mission with `scripts/spacecraft resolve --json` before mutating work.
- **Must**: Default to Node.js + TypeScript + Fastify + Vitest when no stack is specified.
- **Must**: First milestone for new projects: `dev`/`test`/`build` scripts, `GET /healthz`, `GET /version`, passing tests, passing build.
- **Must**: Verify with `scripts/spacecraft evidence` after each milestone.
- **Must**: Prefer small vertical slices over broad horizontal scaffolding.
- **Must**: Every route handler returns typed responses. Use Fastify schema validation for request bodies and params.
- **Must**: Environment variables for secrets and configuration. Never hardcode credentials.
- **Ask before**: Installing any npm dependencies. List packages first, get approval.
- **Must not**: Add database schema, migrations, or query optimization — separate concern.
- **Must not**: Add deployment config, Docker, CI/CD, or observability unless explicitly requested.

## Reviewer checklist

Use this checklist when reviewing backend code:

- [ ] **Fastify plugin lifecycle violations**
  - Plugins registered after `server.listen()`
  - Async plugin registered without `await register()`
  - Plugin options not typed
  - Duplicate route registrations across plugins
- [ ] **Missing input validation**
  - Route handler without Fastify schema validation (`schema: { body: ..., params: ..., querystring: ... }`)
  - Trusting raw `request.body` without validation
  - Missing response schema for typed output
  - Using `any` or `as` casts to bypass type checking
- [ ] **Unhandled promise rejections**
  - `async` route handlers without try/catch
  - Promise chains without `.catch()`
  - `void` used to suppress unhandled rejection warnings
  - Error handler not registered on Fastify instance
- [ ] **Route handler anti-patterns**
  - Business logic in route handler instead of service layer
  - Direct database access from handler
  - Synchronous blocking calls in async handlers
  - Response sent before async operation completes
  - Mixing REST and RPC patterns inconsistently

## Out of scope

- UI design or frontend architecture — separate concern
- Database schema, migrations, or query optimization — separate concern
- System architecture decisions or ADR writing — separate concern
- TDD discipline — use sc-tdd for test-first workflow
- Browser-side code — separate concern

## Output format

```
Stack: Node.js + TypeScript + Fastify + Vitest
Scaffold (new project):
  package.json ✓ (dev, test, build, start)
  tsconfig.json ✓ (strict, ESNext)
  src/server.ts ✓ (GET /healthz, GET /version)
  src/server.test.ts ✓ (endpoint tests)
Route: <method> <path>
  Schema: Fastify JSON Schema (request body, params, response)
  Handler: typed request → service → typed response
  Test: fastify.inject() → status + body assertions
Verify:
  npm test → PASS
  npm run build → PASS
Evidence: <label>
```

## Checklist

Before claiming backend work done:

- [ ] Mission resolved, branch created
- [ ] Stack confirmed: Node.js + TypeScript + Fastify + Vitest (or approved alternative)
- [ ] Routes use Fastify schema validation for request bodies and params
- [ ] Handler logic separated from route registration
- [ ] Secrets in environment variables, never hardcoded
- [ ] Tests cover success, validation error, and server error cases
- [ ] All tests pass (`npm test`)
- [ ] Build passes (`npm run build`)
- [ ] Evidence captured with `scripts/spacecraft evidence`
- [ ] No unapproved dependencies

## Research auto-trigger

When the default stack version is uncertain or the user specifies an unfamiliar framework, run `spacecraft research "<framework> setup guide"` before scaffolding. Backend frameworks evolve — verify current APIs.

---

## References

- `references/api-patterns.md` — REST route design, Fastify schema validation, error handling, middleware
- `references/testing.md` — Vitest + Fastify inject patterns, integration testing, mocking external services
