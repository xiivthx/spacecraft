---
name: sc-diagram
description: "Build standalone interactive HTML block diagrams with click-to-trace nets. Use when the user asks for a block diagram, wiring diagram, interactive diagram, architecture blocks HTML, or /sc-diagram."
disable-model-invocation: true
---

# sc-diagram

## Goal

Produce a **single-file** interactive HTML block diagram (boards, modules, buses, signal nets) with hover/click net highlighting. No build step. No Mermaid-in-markdown. Not product UI (use `sc-ux-design`).

## Output

`docs/<name>.html` (or path the user names): standalone HTML with inline CSS/JS, SVG canvas, legend optional. Link from the nearest README when one exists.

## Good / Bad

- Good: bring-up interconnect, module block map, pin/net wiring overview, signal-path explainer
- Bad: pixel-perfect PCB, schematic capture, live app screens, C4/ADR prose (use `sc-architect`), mission ceremony

## Verify

Open the HTML in a browser. Click each net: matching wires + endpoints highlight. Labels readable; pads/text not covered. Horizontal scroll OK on narrow viewports. No invented nets vs agreed SoT.

## Arguments

```
/sc-diagram
/sc-diagram <short title or path>
```

`$ARGUMENTS` = optional diagram title or output path hint.

## When to use

Activate on: "block diagram", "interactive diagram", "wiring diagram", "architecture blocks HTML", `/sc-diagram`.

Prefer this over ASCII/Mermaid when the human needs **click-to-trace** nets. Prefer `sc-architect` for C4/ADR text decisions.

## Pre-flight

1. Confirm output path (default `docs/<slug>.html`).
2. Collect **source of truth** - never invent nets:
   - Blocks/modules (name, role)
   - Ports/pads (stable ids)
   - Nets (name → endpoints)
3. If SoT missing or ambiguous → ask before drawing.

## Workflow

### 1. Inventory HIL

Propose a short list for confirm:

```
Blocks: …
Nets: …
Open questions: …
```

Do not emit HTML until the inventory is accepted (or user says "just draw it" with enough SoT).

### 2. Emit HTML

Follow `references/html-template.md`:

- Dark interactive default
- SVG wires **behind** blocks (paint order / z-index)
- Soft wire opacity; highlight brightens net + endpoints
- Stable `data-net` / `data-block` / `data-port` attributes
- Label chips beside pads, not on top of text
- Inline CSS/JS only (one file)

### 3. Layout pass

Apply `references/anti-patterns.md`: reduce wire crossings, fix label collisions, keep legend clear.

### 4. Review

Serve or open the file. Iterate until click-trace and labels pass Verify.

### 5. Docs link

If a nearby README exists, add one link to the diagram.

## Rules

- **Must**: Ask when nets/ports are unknown - do not invent interconnect.
- **Must**: One standalone HTML file; no bundler, no CDN-required runtime for core highlight.
- **Must**: Use `data-net` on every wire and related port/label so one click highlights the full net.
- **Must**: Paint wires behind blocks; soft default opacity, strong highlight.
- **Must not**: Replace `sc-architect` C4/ADR workflow for architecture decisions.
- **Must not**: Build React/Vite apps or print A4 twins in v1 (print twin is out of scope).
- **Must not**: Cover pad text with wires or opaque chips.

## Hard stops

- No SoT and user will not provide one
- Request is full product UI → `sc-ux-design`
- Request is ADR/C4 only → `sc-architect`
- Scope becomes a multi-page app

## Summary format

```
Diagram: docs/<name>.html
Blocks: N
Nets: N
SoT: <source>
Verify: click-trace OK | issues…
Next: iterate layout / link README / done
```

## References

- `references/html-template.md` - HTML/SVG skeleton + highlight JS contract
- `references/anti-patterns.md` - overlap, z-order, labels, color soup
