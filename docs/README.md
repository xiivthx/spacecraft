# Product docs

Tracked product Source of Truth. Prefer short current pages. Cold-start: read `docs/` then `.space/`.

## Layout

```text
docs/
  README.md              # this map
  installation.md        # install User + Project layers
  prompting.md           # how we instruct agents
  mission-artifacts.md   # mission schemas + outcome gates
  review.md              # mission + UX review overview
  impeccable.md          # Impeccable × discuss contract
  harness.md             # maintainer scorecard
  conventions/           # engineering norms
  demos/                 # landing + console HTML
  specs/                 # durable product contracts (add when needed)
  architecture/decisions/# ADRs (add when needed)
  runbooks/              # ops playbooks (add when needed)
```

Create folders when you have real content.

## Belong / do not

| Belong | Do not |
|--------|--------|
| Specs, ADRs, conventions, runbooks | Secrets, scratch notes, build caches |
| Lasting product contracts | Full mission discuss dumps |

At ship, promote durable contracts into `specs/` or an ADR.

## Start

1. [installation.md](./installation.md) — install
2. [mission-artifacts.md](./mission-artifacts.md) — mission shapes
3. [prompting.md](./prompting.md) — agent instruction
4. [conventions/](./conventions/) — day-to-day norms
