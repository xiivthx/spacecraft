> Consult when: making web-specific architecture decisions (SPA vs SSR, API gateway, CDN), or the mission involves a web application.

# Web architecture patterns

Web-specific architecture references. Covers rendering strategies, API gateway, CDN, and scaling for browser-based applications.

## Rendering strategies

| Strategy | Description | When to use |
|----------|-------------|-------------|
| **SPA (Single Page App)** | Browser loads JS bundle; renders client-side | Interactive apps, dashboards, tools |
| **SSR (Server-Side Rendering)** | Server renders HTML per request | Content sites, SEO-critical pages, fast first paint |
| **SSG (Static Site Gen)** | Pre-render HTML at build time | Blogs, docs, marketing pages |
| **ISR (Incremental Static Regeneration)** | Rebuild static pages on demand | Hybrid: mostly static with some dynamic content |

Decision framework:
1. Is SEO critical? → SSR or SSG
2. High interactivity (dashboards, editing tools)? → SPA
3. Mostly read-only content? → SSG
4. Need both? → Hybrid (Next.js, Remix, Astro)

```ascii
Need SEO?
├── Yes → High interactivity?
│         ├── Yes → SSR with hydration (Next.js, Remix)
│         └── No  → SSG (Astro, Hugo)
└── No  → SPA (Vite + React)
```

## API gateway patterns

An API gateway sits between clients and backend services. Use when you have ≥3 backend services or need cross-cutting API concerns.

**Responsibilities**:
- **Routing** — direct requests to the correct service
- **Auth** — validate tokens before forwarding
- **Rate limiting** — protect backend services
- **Response aggregation** — combine multiple service responses
- **Protocol translation** — REST ↔ gRPC ↔ GraphQL

**Gateway options**:
- **Lightweight reverse proxy** — nginx, Caddy, Envoy (ops-managed)
- **API management platform** — Kong, Tyk, Apigee (when you need developer portals, billing)
- **BFF (Backend for Frontend)** — one gateway per client type (web, mobile, IoT)

**When NOT to use an API gateway**: Single service, no auth, no rate limiting needed. Direct client-to-service is simpler.

## CDN strategy

**Cacheable content**:
- Static assets: JS bundles, CSS, images, fonts → cache aggressively (1 year, filename hashing)
- HTML pages: cache based on freshness (SSG: long TTL, SSR: short or no cache)
- API responses: cache GET responses where data changes slowly

**Cache invalidation**:
- Filename hashing (webpack/vite chunk hashes) — no invalidation needed
- Purge by URL — for dynamic content updates
- Surrogate keys / cache tags — purge groups of related content

**Edge compute**: Run logic at CDN edge nodes for personalization, A/B testing, geo-routing without hitting origin.

## Horizontal scaling

| Pattern | Description | When to use |
|---------|-------------|-------------|
| **Stateless services** | No session affinity; any instance handles any request | Default for web APIs |
| **Sticky sessions** | Route same user to same instance | When avoiding state migration (temporary, not strategic) |
| **Read replicas** | Separate read/write database paths | Read-heavy workloads |
| **CQRS** | Separate read/write models entirely | Complex query patterns, event-sourced systems |
| **Event-driven** | Async communication via message queue | Long-running workflows, eventual consistency acceptable |

## Spacecraft integration

- Web architecture decisions belong in the mission `decisions.md` as ADRs
- Reference this file when planning or reviewing web-specific architecture
- C4 Level 1 (System Context) for web apps always includes: user browser, CDN, API gateway (if present), backend service(s), database(s)
- Document rendering strategy choice before starting frontend or backend implementation
