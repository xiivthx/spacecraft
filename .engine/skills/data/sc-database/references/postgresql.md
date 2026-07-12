> Consult when: working with PostgreSQL specifically — choosing column types, creating indexes, or analyzing query performance with EXPLAIN.

# PostgreSQL engine reference

PostgreSQL-specific operational patterns. Currently at version 18 (stable, 2026). For universal database practices (normalization, schema design), see the parent SKILL.md.

## Column types

Prefer domain-appropriate types. Do not default to `text` or `integer` for everything.

| Use case | Type | Notes |
|----------|------|-------|
| Auto-incrementing ID | `SERIAL` or `BIGSERIAL` | Or `GENERATED ALWAYS AS IDENTITY` (SQL standard) |
| UUID primary key | `UUID` with `gen_random_uuid()` | Better for distributed systems |
| Short text (≤255) | `VARCHAR(n)` or `TEXT` with `CHECK` | `TEXT` has no performance penalty in PG |
| Boolean | `BOOLEAN` | Not `INTEGER` with 0/1 |
| Date/time | `TIMESTAMPTZ` | Always use timezone-aware. `TIMESTAMP` is ambiguous. |
| Money | `NUMERIC(19,4)` | Not `FLOAT` — floating point loses precision |
| JSON | `JSONB` | For semi-structured data. Not `JSON` — `JSONB` is indexed and faster for reads. |
| Enum | Native `CREATE TYPE` or `VARCHAR` with `CHECK` | Native enums are harder to modify; prefer CHECK for evolvable enums |
| Array | `TEXT[]`, `INTEGER[]` | Use sparingly. Normalized tables are usually better. |

## Index types

| Type | Use case | Example |
|------|----------|---------|
| **B-tree** (default) | Equality, range, sorting | `CREATE INDEX idx_users_email ON users (email);` |
| **Hash** | Equality only (no range) | Faster than B-tree for `=` but rarely needed |
| **GIN** | Full-text search, JSONB, arrays | `CREATE INDEX idx_posts_content ON posts USING GIN (to_tsvector('english', content));` |
| **GiST** | Geometric data, full-text | Spatial queries, trigram matching |
| **BRIN** | Very large tables, sequential data | `CREATE INDEX idx_events_time ON events USING BRIN (created_at);` — much smaller index |

**Partial indexes**: Index only a subset of rows to save space and write overhead.

```sql
-- Only index active users
CREATE INDEX idx_active_users ON users (email) WHERE deleted_at IS NULL;
```

**Composite indexes**: Column order matters — put equality columns first, range columns last.

```sql
-- Good: user_id (equality) before created_at (range)
CREATE INDEX idx_orders_user_time ON orders (user_id, created_at DESC);

-- Bad: range first — B-tree can't use created_at efficiently before user_id
CREATE INDEX idx_orders_time_user ON orders (created_at, user_id);
```

## EXPLAIN ANALYZE

Always verify query performance and index usage with `EXPLAIN ANALYZE`. Key indicators:

- **Seq Scan** → full table scan. May need an index if the table is large and the query is frequent.
- **Index Scan** → reading from index, then fetching rows from heap. Good for selective queries.
- **Index Only Scan** → reading from index without touching heap. Best — all columns in the index.
- **Nested Loop** → join strategy. Fast for small result sets; slow for large ones.
- **Hash Join** → build hash table then probe. Good for medium/large joins.
- **Merge Join** → both inputs sorted. Fastest for large sorted joins.

```sql
EXPLAIN ANALYZE
SELECT u.name, o.total
FROM users u
JOIN orders o ON u.id = o.user_id
WHERE u.created_at > now() - interval '30 days'
ORDER BY o.total DESC
LIMIT 50;
```

## Connection pooling

PostgreSQL uses process-per-connection. Do not open a connection per request. Use a pool:

- **PgBouncer** — external pool (production, multiple services)
- **Application pool** — `pg-pool` (Node.js), `HikariCP` (Java), `SQLAlchemy pool` (Python)
- **Serverless** — `@neondatabase/serverless` for edge/worker environments

Pool sizing: start with `(CPU cores * 2) + effective_disk_spindle_count`. Monitor `pg_stat_activity` and adjust.

## Spacecraft integration

- PostgreSQL engine choice is an architectural decision — record in `decisions.md`
- Migration files and index definitions are product code — stage alongside schema
- EXPLAIN ANALYZE output is evidence of query performance — capture with `scripts/spacecraft evidence`
- Connection configuration (pool size, credentials) goes in environment variables, never in code
