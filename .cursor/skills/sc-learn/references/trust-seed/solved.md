# Trust solved

Tracked seed for `.space/trust/solved.md`. Prefer **empty tables** here; append only high-signal regressions at ship (not every nit).

| Mission | Date | Problem | Solution | Evidence |
|---------|------|---------|----------|----------|
| M07SMYDMZ | 2026-07-21 | Resolver preferred branch over `spacecraft use` (`.space/current`) | Prefer `.space/current` in `resolveActive`; evidence without `--mission` uses it | validate --strict |
| M07SML0VX | 2026-07-21 | Evidence rewrite dropped `outputHash` integrity | Optional SHA-256 `outputHash` on write; validate when present | ship-reverify |
