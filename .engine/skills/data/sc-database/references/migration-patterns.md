> Consult when: writing database migrations, setting up a migration framework, designing seed data, or planning rollback strategies.

# Migration patterns

Operational patterns for database schema versioning, seeding, and rollback. Framework-agnostic principles with examples.

## Migration structure

Every migration has an `up` and `down` direction. The system tracks which migrations have been applied.

```
migrations/
├── 001_create_users.sql
├── 001_create_users.down.sql
├── 002_add_email_index.sql
├── 002_add_email_index.down.sql
├── 003_create_orders.sql
└── 003_create_orders.down.sql
```

Or with a migration framework:

```typescript
// migrations/001_create_users.ts (Knex example)
export async function up(knex: Knex): Promise<void> {
  await knex.schema.createTable("users", (t) => {
    t.uuid("id").primary().defaultTo(knex.raw("gen_random_uuid()"));
    t.string("name").notNullable();
    t.string("email").notNullable().unique();
    t.timestamps(true, true); // created_at, updated_at
  });
}

export async function down(knex: Knex): Promise<void> {
  await knex.schema.dropTable("users");
}
```

## Migration rules

1. **Sequential versioning** — Number migrations sequentially. Never renumber or reorder committed migrations.
2. **Never modify a committed migration** — Write a new migration to alter the schema.
3. **Every up has a down** — `down` reverses exactly what `up` did, in reverse order.
4. **Idempotent when possible** — Use `IF EXISTS` / `IF NOT EXISTS` to avoid errors on re-run.
5. **No data in schema migrations** — Schema changes and data changes are separate. Use seed files for reference data.

```sql
-- 001_create_users.up.sql
CREATE TABLE IF NOT EXISTS users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name TEXT NOT NULL,
  email TEXT NOT NULL UNIQUE,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- 001_create_users.down.sql
DROP TABLE IF EXISTS users;
```

## Seeding

Seed files populate reference data and development/test fixtures. Separate from schema migrations — seeds can be re-run safely.

```typescript
// seeds/01_roles.ts
export async function seed(knex: Knex): Promise<void> {
  await knex("roles").insert([
    { name: "admin" },
    { name: "editor" },
    { name: "viewer" },
  ]).onConflict("name").ignore(); // Safe to re-run
}
```

**Seed conventions**:
- Seeds are additive — they insert or upsert, never truncate
- Use `ON CONFLICT ... DO NOTHING` or upsert to make seeds re-runnable
- Seeds run after all migrations are applied
- Production seeds are only for reference data (roles, categories, feature flags), never sample/test data

## Rollback strategy

| Scenario | Strategy |
|----------|----------|
| **Dev environment** | Roll back the migration, fix, re-apply |
| **Staging** | Roll back last migration, fix, re-apply (data is disposable) |
| **Production** | Write a forward-only fix migration. Never roll back a migration that has been in production — rolling back may destroy data. |
| **Data-destructive change** | Expand first: add new column, deploy code that writes both old and new, backfill, deploy code that reads new only, drop old column |

**Production safety**: The column rename dance:
```sql
-- Step 1: Add new column (migration N)
ALTER TABLE users ADD COLUMN full_name TEXT;

-- Step 2: Deploy code that writes to both name and full_name

-- Step 3: Backfill (migration N+1 - data migration)
UPDATE users SET full_name = name WHERE full_name IS NULL;

-- Step 4: Deploy code that reads full_name only

-- Step 5: Drop old column (migration N+2)
ALTER TABLE users DROP COLUMN name;
```

## Migration frameworks

| Framework | Language | Notes |
|-----------|----------|-------|
| **Knex** | Node.js | Query builder + migrations + seeds |
| **Drizzle Kit** | Node.js | Type-safe, generate from schema, push mode for dev |
| **Prisma Migrate** | Node.js | Declarative schema → generated migrations |
| **Alembic** | Python | SQLAlchemy migrations |
| **Flyway** | Java/Kotlin | SQL-first, Java API |
| **golang-migrate** | Go | SQL file-based, minimal |

## Spacecraft integration

- Migration files are product code — stage and commit them alongside schema changes
- Evidence: `scripts/spacecraft evidence "migration:<name>" -- <migrate-up && migrate-down command>`
- Migration framework choice is a one-way-door architectural decision — record in `decisions.md`
- Seed data for tests is part of the test setup; reference data seeds are product code
- Rollback decisions for production are architectural decisions — write an ADR
