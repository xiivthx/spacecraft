---
name: sc-database
description: "Design schemas, write migrations, optimize queries, and manage indexes. PostgreSQL by default. Activate on \"set up PostgreSQL\", \"create database schema\", \"write a migration\", \"optimize query\", or \"design database indexes\"."
---

# sc-database

Design and manage databases under mission control. Universal practices with PostgreSQL as the default engine. Schema design, normalization, migrations, indexing, and query optimization.

## When to use

Activate when the user asks to:

- **"Set up PostgreSQL" / "connect to database"** - database initialization
- **"Create database schema" / "design tables"** - schema design and normalization
- **"Write a migration"** - schema versioning and change management
- **"Optimize query" / "design database indexes"** - performance tuning
- When a mission task requires database design or management

## Workflow

Use this exact sequence unless the user specifies otherwise:

1. **Resolve mission** - `spacecraft resolve --json`. Block if safety ≠ `safe`.

2. **Choose engine** - Default: PostgreSQL. If the project already uses a different engine (MySQL, SQLite), match it. Engine-specific details are in `references/<engine>.md`.

3. **Design schema** - Before writing DDL:
   - Identify entities and relationships from the mission `spec.md`
   - Apply normalization: 1NF (atomic columns) → 2NF (no partial dependencies) → 3NF (no transitive dependencies). Denormalize only with a documented reason.
   - Choose column types: prefer domain-appropriate types over generic `text`/`integer`
   - Add constraints: `NOT NULL`, `UNIQUE`, `CHECK`, `FOREIGN KEY` at schema level

4. **Write migration** - Use `references/migration-patterns.md` for strategy:
   - Up migration: creates tables, columns, indexes, constraints
   - Down migration: reverses the up migration (drops in reverse order)
   - Seed data: reference data inserted via separate seed files
   Version migrations sequentially. Never modify a committed migration - write a new one.

5. **Index strategy** - Add indexes for:
   - Primary keys (automatic)
   - Foreign keys (join performance)
   - Columns in `WHERE` clauses of frequent queries
   - Columns in `ORDER BY` of sorted queries
   Do not index every column - each index costs write performance and storage.

6. **Verify** - `spacecraft evidence "<label>" -- <migration-test-command>`. Migrations must apply cleanly and roll back cleanly. Query tests must demonstrate index usage (`EXPLAIN ANALYZE`).

### Edge cases

- **Non-PostgreSQL engine** - Load the relevant `references/<engine>.md`. Universal practices (normalization, migration patterns, indexing strategy) still apply.
- **Existing database** - Don't redesign from scratch. Add migrations incrementally. Document denormalization decisions in `decisions.md`.
- **Performance regression** - Check query plans with `EXPLAIN ANALYZE`. Validate index usage. Consider partial indexes for filtered queries.
- **Zero-downtime migration** - For production: add columns with defaults, deploy code that handles both old and new schema, then drop old columns in a separate migration.

## Rules

- **Must**: Resolve mission with `spacecraft resolve --json` before mutating work.
- **Must**: Default to PostgreSQL when no engine is specified.
- **Must**: Every schema change goes through a versioned migration file.
- **Must**: Migrations must be reversible - every `up` has a corresponding `down`.
- **Must**: Use `EXPLAIN ANALYZE` to verify index usage before claiming query optimization.
- **Must**: Apply constraints at the database level (NOT NULL, UNIQUE, CHECK, FOREIGN KEY). The database is the last line of defense.
- **Must not**: Modify a committed migration. Write a new migration instead.
- **Must not**: Add indexes without measuring. Verify with `EXPLAIN ANALYZE` that the index is used.
- **Must not**: Store secrets or connection strings in migration files.

## Out of scope

- Application-layer data access (repositories, ORMs, API endpoints) - separate concern
- System architecture decisions - separate concern
- UI design or frontend architecture - separate concern
- TDD discipline - use sc-tdd

## Output format

```
Engine: PostgreSQL (default)
Schema:
  Tables: <count>
  Relationships: <diagram or list>
  Normalization: 3NF (or documented denormalization)
Migration:
  Up: create <objects>
  Down: drop <objects> (reverse order)
Indexes:
  - <table>.<column>: <type> - <reason>
Verify:
  EXPLAIN ANALYZE <query> → <plan summary>
Evidence: <label>
```

## Checklist

Before claiming database work done:

- [ ] Mission resolved, engine chosen
- [ ] Schema designed: entities, relationships, column types, constraints
- [ ] Normalization applied (3NF target, documented exceptions)
- [ ] Migration files written: up migration, down migration, reversible
- [ ] Indexes added for foreign keys and query WHERE/ORDER BY columns
- [ ] `EXPLAIN ANALYZE` confirms index usage for performance-critical queries
- [ ] Migrations apply and roll back cleanly
- [ ] Evidence captured with `spacecraft evidence`
- [ ] No secrets or connection strings in migration files

## References

- `references/postgresql.md` - PostgreSQL-specific engine details, types, index types, EXPLAIN
- `references/migration-patterns.md` - Up/down migrations, seeding, versioning, rollback strategies
