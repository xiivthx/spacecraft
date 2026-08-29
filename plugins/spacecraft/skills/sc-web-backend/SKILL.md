---
name: sc-web-backend
description: "Build server APIs and services with Node.js, TypeScript, Fastify, and Vitest. Activate on \"create an API endpoint\", \"build a web service\", \"scaffold a server\", \"add authentication\", or \"design REST API\"."
---

# sc-web-backend

Build server APIs/services under mission control. Default stack: Node.js + TypeScript + Fastify + Vitest.

## When to use

API endpoint / web service; Fastify scaffold; auth / REST design; health/version milestones; mission server tasks.

## Workflow

1. **Resolve** - `spacecraft resolve`; conflict → `spacecraft use <selector>`.
2. **Stack** - Default above; adapt if user specifies. Uncertain versions → sc-search (`fastify v5 typescript setup`).
3. **Scaffold (new)** - `package.json` (`dev`/`test`/`build`/`start`); strict `tsconfig`; `src/server.ts` with `GET /healthz` + `GET /version`; matching tests; `.gitignore`. **Ask before** npm install.
4. **Slice** - Route + schema validation → typed handler → service layer → Vitest + `fastify.inject()`. Patterns: `references/api-patterns.md`.
5. **Verify** - `spacecraft evidence "<label>" -- npm test`; `npm run build` after each slice.

### Edge cases

Existing project → add routes, no re-scaffold. Auth → Fastify plugin; hash passwords; secrets in env. DB → existing data layer / defer schema to sc-database. GraphQL → typed resolvers, same test bar. Failures → fix before continue.

## Rules

- **Must**: Resolve mission with `spacecraft resolve` before mutating work. On conflict/ambiguity use `spacecraft use <selector>`.
- **Must**: Default to Node.js + TypeScript + Fastify + Vitest when no stack is specified.
- **Must**: First milestone for new projects: `dev`/`test`/`build` scripts, `GET /healthz`, `GET /version`, passing tests, passing build.
- **Must**: Verify with `spacecraft evidence` after each milestone.
- **Must**: Prefer small vertical slices over broad horizontal scaffolding.
- **Must**: Every route handler returns typed responses. Use Fastify schema validation for request bodies and params.
- **Must**: Environment variables for secrets and configuration. Never hardcode credentials.
- **Ask before**: Installing any npm dependencies. List packages first, get approval.
- **Must not**: Add database schema, migrations, or query optimization - separate concern.
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

UI/frontend · DB schema/migrations · ADRs (`sc-architect`) · TDD ceremony (`sc-tdd`) · browser code

## Output format

Stack · scaffold status (if new) · route method/path + schema + service + inject test · `npm test` / `npm run build` · evidence label.

## Checklist

Resolved · stack confirmed · schema validation · handler ≠ registration · secrets in env · tests cover success/validation/server error · test+build pass · evidence · no unapproved deps.

## References

- `references/api-patterns.md` - REST, Fastify schema, errors, middleware
- `references/testing.md` - Vitest + inject, integration, mocking
