---
name: sc-database
description: "Design schemas, write migrations, optimize queries, and manage indexes. PostgreSQL by default. Activate on \"set up PostgreSQL\", \"create database schema\", \"write a migration\", \"optimize query\", or \"design database indexes\"."
---

# sc-database

Design and manage databases under mission control. Default engine: PostgreSQL. Schema, migrations, indexing, query optimization.

## When to use

PostgreSQL setup; schema/tables; migrations; query optimize / indexes; mission DB tasks.

## Workflow

1. **Resolve** - `spacecraft resolve`; conflict → `spacecraft use <selector>`.
2. **Engine** - Default PostgreSQL; else match project. Detail: `references/<engine>.md`.
3. **Schema** - Entities/relationships from `spec.md`; 1NF→2NF→3NF (denormalize only with documented reason); domain types; constraints (`NOT NULL`, `UNIQUE`, `CHECK`, `FOREIGN KEY`) at schema level.
4. **Migration** - `references/migration-patterns.md`: versioned up/down; never edit committed migrations; seeds separate.
5. **Indexes** - PK/FK; frequent `WHERE` / `ORDER BY`. Do not index every column.
6. **Verify** - `spacecraft evidence "<label>" -- <migration-test-command>`; apply + rollback clean; `EXPLAIN ANALYZE` for index claims.

### Edge cases

Non-PG → `references/<engine>.md`. Existing DB → incremental migrations. Regression → `EXPLAIN ANALYZE` / partial indexes. Zero-downtime → expand/contract pattern.

## Rules

- **Must**: Resolve mission with `spacecraft resolve` before mutating work. On conflict/ambiguity use `spacecraft use <selector>`.
- **Must**: Default to PostgreSQL when no engine is specified.
- **Must**: Every schema change goes through a versioned migration file.
- **Must**: Migrations must be reversible - every `up` has a corresponding `down`.
- **Must**: Use `EXPLAIN ANALYZE` to verify index usage before claiming query optimization.
- **Must**: Apply constraints at the database level (NOT NULL, UNIQUE, CHECK, FOREIGN KEY). The database is the last line of defense.
- **Must not**: Modify a committed migration. Write a new migration instead.
- **Must not**: Add indexes without measuring. Verify with `EXPLAIN ANALYZE` that the index is used.
- **Must not**: Store secrets or connection strings in migration files.

## Out of scope

App data access / ORMs / API endpoints · system architecture · UI · TDD (`sc-tdd`)

## Output format

Engine · tables/relationships · normalization (or documented denorm) · up/down migration · indexes with reason · `EXPLAIN ANALYZE` summary · evidence label.

## Checklist

Resolved + engine · schema + 3NF (or documented exception) · reversible migrations · indexes measured · evidence · no secrets in migrations.

## References

- `references/postgresql.md` - types, index types, EXPLAIN
- `references/migration-patterns.md` - up/down, seeding, versioning, rollback
