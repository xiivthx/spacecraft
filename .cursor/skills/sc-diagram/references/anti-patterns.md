> Consult when reviewing layout quality of an interactive block diagram.

# Anti-patterns

## Z-order / occlusion

- **Wires on top of blocks** - always paint `#wires` before `#blocks` / `#labels`.
- **Labels covering pads** - place label chips beside the port, not centered on the pad text.
- **Opaque chips over wire endpoints** - leave a clear gap at the pad; wire should meet the port circle.

## Routing

- **Hairball crossings** - prefer orthogonal stubs then a shared trunk; reduce diagonal crossings through text.
- **Wires through board titles** - route in gutters between blocks.
- **Identical paths stacked** - offset parallel nets by a few px so each remains clickable.

## Highlight / opacity

- **Wires at full opacity always** - default soft (`~0.4–0.5`); only highlighted net goes bright.
- **Highlight without dimming others** - dim non-matching nets so the active path is obvious.
- **Color-only highlight** - also bump stroke width / opacity for colorblind safety.

## Labels and density

- **Tiny unreadable type** - body labels ≥12px CSS px equivalent in the viewBox scale.
- **Duplicate conflicting names** - one SoT name per net; Arduino alias vs connector pin as secondary chip, not a second net id.
- **Legend missing when >8 nets** - add a legend once net count gets crowded.

## Color soup

- **Rainbow per wire by default** - one muted wire color; reserve strong hues for highlight or bus families (SPI / UART / power) if needed.
- **Purple-glow AI slop** - stick to the dark token set in `html-template.md`; no neon purple gradients.

## Structure

- **Multiple HTML files for one diagram** - v1 is one standalone file.
- **CDN-required for click-trace** - core JS must run offline.
- **Invented nets to "look complete"** - stop and ask; incomplete SoT beats fake wires.

## Interaction

- **Hover-only with no click lock** - click must set a sticky highlight; empty click clears.
- **Hit targets too thin** - stroke/hit area ≥ clickable; invisible wider hit path OK if needed.

## Checklist before "done"

- [ ] Wires behind blocks
- [ ] Soft default wire opacity + strong highlight + dim others
- [ ] Every wire/port/label in a net shares the same `data-net`
- [ ] No label covering pad text
- [ ] Click empty space clears
- [ ] Readable at default zoom; scrollable on narrow viewports
