> Consult when: designing REST routes, adding Fastify schema validation, structuring route handlers with services, or implementing authentication.

# API patterns

Operational patterns for building APIs with Node.js + TypeScript + Fastify.

## Route structure

Separate route registration from business logic. Routes handle HTTP concerns; services handle domain logic.

```typescript
// src/routes/users.ts
import { FastifyInstance, FastifyRequest, FastifyReply } from "fastify";
import { createUser, UserInput } from "../services/users";

interface CreateUserBody {
  name: string;
  email: string;
}

export async function userRoutes(app: FastifyInstance) {
  app.post<{ Body: CreateUserBody }>("/users", {
    schema: {
      body: {
        type: "object",
        required: ["name", "email"],
        properties: {
          name: { type: "string", minLength: 1 },
          email: { type: "string", format: "email" },
        },
      },
      response: {
        201: {
          type: "object",
          properties: {
            id: { type: "string" },
            name: { type: "string" },
            email: { type: "string" },
          },
        },
      },
    },
    handler: async (request, reply) => {
      const user = await createUser(request.body);
      return reply.status(201).send(user);
    },
  });
}
```

Key conventions:
- Route files export an async function receiving the Fastify instance
- Schema validation on every route — body, params, query, response
- Handlers delegate to service functions, never contain business logic directly
- Typed generics on route options: `app.post<{ Body: T }>`

## Health and version endpoints

Every new service starts with these endpoints:

```typescript
app.get("/healthz", async () => ({ ok: true }));
app.get("/version", async () => ({
  version: process.env.npm_package_version || "0.0.0",
  build: process.env.BUILD_ID || "dev",
}));
```

## Error handling

Throw domain errors with status codes; let Fastify format via `setErrorHandler`.

```typescript
class NotFoundError extends Error { statusCode = 404; }
class ValidationError extends Error { statusCode = 400; }

app.setErrorHandler((error, request, reply) => {
  const statusCode = error.statusCode || 500;
  reply.status(statusCode).send({
    error: statusCode >= 500 ? "Internal Server Error" : error.message,
    statusCode,
  });
});
```

## Middleware and plugins

Fastify plugins encapsulate cross-cutting concerns via `fastify-plugin`.

```typescript
async function authPlugin(app: FastifyInstance) {
  app.decorateRequest("userId", null);
  app.addHook("onRequest", async (request) => {
    const token = request.headers.authorization?.replace("Bearer ", "");
    if (token) request.userId = verifyToken(token).sub;
  });
}
export default fp(authPlugin);
```

Common plugins: auth (bearer/session/API keys), CORS (`@fastify/cors`), rate limiting (`@fastify/rate-limit`), logging (built-in `logger: true`).

## REST conventions

| Method | Path | Purpose | Success code |
|--------|------|---------|-------------|
| `GET` | `/items` | List resources | 200 |
| `GET` | `/items/:id` | Get one resource | 200 |
| `POST` | `/items` | Create resource | 201 |
| `PUT` | `/items/:id` | Replace resource | 200 |
| `PATCH` | `/items/:id` | Partial update | 200 |
| `DELETE` | `/items/:id` | Remove resource | 204 |

Use plural nouns. Nest related resources: `/users/:userId/orders`. Keep nesting ≤2 levels deep.

## Spacecraft integration

- Routes correspond to `plan.json` acceptance checks — one slice per endpoint group
- Capture evidence: `scripts/spacecraft evidence "<endpoint>:api" -- npm test`
- Scaffold phase: `scripts/spacecraft evidence "web:endpoints" -- npm test`
- Environment variables documented in `decisions.md` when adding new config
