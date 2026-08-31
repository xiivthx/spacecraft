# Visual pack

Rendered look is wrong while behavior may be fine. Oracle is approved draft / `DESIGN.md` / live product - not `expect()`. Load only after `Pack: visual`.

Click-handler / data / routing bugs are **software**. Taste-only ("make it prettier") is `/sc-discuss` + designer, not this pack.

## Classify

| Symptom | Pack |
|---|---|
| Wrong data, dead click, bad API, failed assert | software |
| Right DOM/action, wrong paint (color, overlap, shift, hidden by stacking) | visual |
| Look and behavior both wrong | split: software first if the flow is false, then visual vs oracle |

Look vs behavior conflict on the spec → `/sc-discuss` (Commander does not pick).

## Repro

Paired screenshots at the same viewport (375 / 768 / 1280, + 1536 when multi-region), theme, and seed data. Record URL, route, and state. Draft present → draft-surface vs live (sc-ux-design Step 0). No oracle (no draft, no DESIGN.md, no accepted screenshot) → `blocked: no-repro` or escalate `/sc-discuss`.

Stabilize animation/fonts before treating a one-frame diff as a bug.

## Observe

1. **Side-by-side** - draft vs live or before vs after. Pixel diff is a **detector**, not RCA (anti-alias, font, GPU noise).
2. **Box + stack** - computed style, margin/padding, overflow, z-index, stacking context, transform/opacity creating a context.
3. **Live pass** - sc-browser-probe / playwright-cli / Cursor IDE browser. Interaction, not screenshot-only. `PROBE:` is not `ready` / `VERIFIED` / AUTH.
4. **Designer** - Task(`sc-designer`) when the oracle is look and the delta is craft, after the computed-style fail path is known.

Do not use Cursor Debug Mode. Do not "fix CSS" by guessing a file.

## Isolate

One viewport, one surface, one token. Token-level change that paints many screens is still one cause - confirm with a second surface. Responsive: if only one preset breaks, the ladder (not the token) is the suspect.

## Fix

Surgical CSS/layout via Task(`sc-coder`) when the oracle is already approved. Changing the look contract (new world, new tokens, new layout idea) → `/sc-discuss`, not `/sc-quick`. Re-check the paired screenshots. TWINS for the same construct (z-index, overflow, token).

## Out

Draft-parity / live-product gates for `/sc-run` stay in sc-ux-design. This pack only finds why this screen diverges **now**.
