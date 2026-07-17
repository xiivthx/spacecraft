> Consult when: writing API endpoint tests, using Fastify inject, mocking external services, or setting up integration test infrastructure.

# Testing patterns

Operational patterns for testing Fastify APIs with Vitest.

## Test structure

Use `fastify.inject()` for end-to-end HTTP testing without starting a real server. Tests are fast and deterministic.

```typescript
import { describe, it, expect, beforeAll, afterAll } from "vitest";
import Fastify from "fastify";
import { userRoutes } from "../routes/users";

describe("POST /users", () => {
  const app = Fastify();

  beforeAll(async () => {
    await app.register(userRoutes);
    await app.ready();
  });

  afterAll(async () => {
    await app.close();
  });

  it("creates a user and returns 201", async () => {
    const response = await app.inject({
      method: "POST",
      url: "/users",
      payload: { name: "Alice", email: "alice@example.com" },
    });

    expect(response.statusCode).toBe(201);
    const body = response.json();
    expect(body).toMatchObject({
      name: "Alice",
      email: "alice@example.com",
    });
    expect(body.id).toBeDefined();
  });

  it("returns 400 when email is missing", async () => {
    const response = await app.inject({
      method: "POST",
      url: "/users",
      payload: { name: "Alice" },
    });

    expect(response.statusCode).toBe(400);
    expect(response.json()).toMatchObject({
      error: expect.any(String),
      statusCode: 400,
    });
  });
});
```

## Test isolation

Each test file builds its own Fastify instance. Do not share instances across test files - shared state causes non-deterministic failures.

```typescript
// Each describe block gets a fresh app
describe("GET /users", () => {
  const app = Fastify();
  beforeAll(async () => {
    await app.register(userRoutes);
    await app.ready();
  });
  afterAll(async () => app.close());
  // tests...
});
```

## Mocking

Mock at system boundaries. Do not mock Fastify or your own services.

- **Database** - Use a test database or mock the repository interface
- **External APIs** - Mock HTTP calls with `vi.mock()` or `nock`
- **Time** - `vi.useFakeTimers()` for time-dependent logic

```typescript
vi.mock("../services/users", () => ({ createUser: vi.fn() }));
import { createUser } from "../services/users";

it("delegates to service layer", async () => {
  vi.mocked(createUser).mockResolvedValue({ id: "1", name: "Alice", email: "alice@example.com" });
  const response = await app.inject({
    method: "POST", url: "/users",
    payload: { name: "Alice", email: "alice@example.com" },
  });
  expect(createUser).toHaveBeenCalledWith({ name: "Alice", email: "alice@example.com" });
  expect(response.statusCode).toBe(201);
});
```

## Test cases checklist

Every route handler test covers:
1. **Success** - valid input, expected output, correct status code
2. **Validation errors** - missing required fields, invalid formats
3. **Not found** - resource doesn't exist (GET/PUT/DELETE by id)
4. **Auth errors** - missing or invalid tokens (protected routes)
5. **Server errors** - service layer throws, returns 500

## Integration tests

For workflows spanning multiple endpoints, write integration tests that call the real chain:

```typescript
it("creates and retrieves a user", async () => {
  // Create
  const create = await app.inject({
    method: "POST",
    url: "/users",
    payload: { name: "Bob", email: "bob@example.com" },
  });
  expect(create.statusCode).toBe(201);
  const { id } = create.json();

  // Retrieve
  const get = await app.inject({
    method: "GET",
    url: `/users/${id}`,
  });
  expect(get.statusCode).toBe(200);
  expect(get.json()).toMatchObject({ name: "Bob" });
});
```

## Spacecraft integration

- Run tests per endpoint: `spacecraft evidence "<endpoint>:test" -- npx vitest run src/routes/<name>.test.ts`
- Full suite: `spacecraft evidence "<task>-functional" -- npm test`
- Integration tests count as functional evidence - they verify end-to-end behavior
- Test files are product code - stage alongside routes in checkpoint commits
