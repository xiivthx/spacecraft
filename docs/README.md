# Product documentation

Tracked source of truth for product behavior, design choices, and operating knowledge.
Keep this tree readable by any engineer who joins the repo.

## Layout

```text
docs/
  README.md                 # this map
  vision.md                 # product vision and non-goals (add when ready)
  epics/                    # outcomes and delivery slices (add as needed)
  specs/                    # durable product contracts and behavior
  conventions/              # shared engineering conventions
  architecture/
    decisions/              # architecture decision records (ADRs)
  runbooks/                 # operational how-tos and incident playbooks
```

Create folders when you have real content. Empty placeholder directories are not required.

## What belongs here

| Area | Use for |
|------|---------|
| `vision.md` | Why the product exists, who it serves, and what is out of scope |
| `epics/` | Multi-step outcomes and how work is sliced for delivery |
| `specs/` | Lasting contracts: APIs, flows, acceptance rules, invariants |
| `conventions/` | Naming, code style, review norms, shared practices |
| `architecture/decisions/` | Significant design choices and the tradeoffs behind them |
| `runbooks/` | Deploy, recover, rotate, and operate the system in production |

Update `docs/` when product truth changes. Prefer short, current pages over historical dumps.

## What does not belong here

- Secrets, credentials, or environment-specific private values
- Scratch notes, local experiments, or one-off investigation dumps
- Generated build artifacts or tool caches
- Content that only exists to satisfy an automation pipeline

Put durable product rules in `specs/` or an ADR under `architecture/decisions/`. Leave temporary working notes out of this tree.

## Start here

1. Read `vision.md` when present, then relevant specs and ADRs.
2. Follow `conventions/` for day-to-day engineering norms.
3. Use `runbooks/` when operating or recovering the system.
