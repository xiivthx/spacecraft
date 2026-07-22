> Consult when emitting or fixing an interactive block-diagram HTML file.

# HTML template contract

Standalone file. Inline `<style>` and `<script>`. No build step. Core highlight must work offline (no CDN required for click-trace).

## Document skeleton

```html
<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8" />
  <meta name="viewport" content="width=device-width, initial-scale=1" />
  <title><!-- diagram title --></title>
  <style>/* tokens + layout + highlight states */</style>
</head>
<body>
  <header class="chrome">
    <h1><!-- title --></h1>
    <p class="hint">Click a wire or port to highlight its net. Click empty space to clear.</p>
  </header>
  <aside class="legend" hidden><!-- optional: net color / name list --></aside>
  <main class="stage">
    <svg id="diagram" viewBox="0 0 W H" role="img" aria-label="…">
      <g id="wires"></g>
      <g id="blocks"></g>
      <g id="labels"></g>
    </svg>
  </main>
  <script>/* highlight */</script>
</body>
</html>
```

Paint order inside SVG: **wires → blocks → labels** so wires sit behind boards and chips.

## CSS tokens (dark default)

```css
:root {
  --bg: #0f1419;
  --surface: #1a222c;
  --text: #e7ecf1;
  --muted: #8b98a8;
  --wire: #5a6a7a;
  --wire-hi: #7dd3fc;
  --accent: #38bdf8;
}
body { background: var(--bg); color: var(--text); }
.wire { stroke: var(--wire); stroke-width: 2; fill: none; opacity: 0.45; }
.wire.is-hi, .port.is-hi, .label.is-hi { opacity: 1; }
.wire.is-hi { stroke: var(--wire-hi); stroke-width: 3; }
.dim { opacity: 0.2; }
```

## Data attributes

| Attribute | On | Purpose |
|-----------|-----|---------|
| `data-net="NET_ID"` | wire path, port, related label | Click-trace key (stable, uppercase slug) |
| `data-block="BLOCK_ID"` | block group | Module identity |
| `data-port="PORT_ID"` | pad/port marker | Endpoint identity |

One net id per logical signal (e.g. `SPI_SCLK`, `GND`). Power/ground may be shared nets.

## Wire and port markup

```html
<g id="wires">
  <path class="wire" data-net="SPI_SCLK" d="M …" />
</g>
<g id="blocks">
  <g class="block" data-block="MCU">
    <rect … />
    <circle class="port" data-block="MCU" data-port="D13" data-net="SPI_SCLK" cx="…" cy="…" r="4" />
  </g>
</g>
<g id="labels">
  <text class="label" data-net="SPI_SCLK" x="…" y="…">SCLK</text>
</g>
```

## Highlight JS contract

```js
(function () {
  const svg = document.getElementById("diagram");
  const nodes = () => svg.querySelectorAll("[data-net]");

  function clear() {
    nodes().forEach((el) => {
      el.classList.remove("is-hi", "dim");
    });
  }

  function highlight(net) {
    nodes().forEach((el) => {
      const match = el.getAttribute("data-net") === net;
      el.classList.toggle("is-hi", match);
      el.classList.toggle("dim", !match);
    });
  }

  svg.addEventListener("click", (e) => {
    const t = e.target.closest("[data-net]");
    if (!t || !svg.contains(t)) {
      clear();
      return;
    }
    highlight(t.getAttribute("data-net"));
  });

  svg.addEventListener("mouseover", (e) => {
    const t = e.target.closest("[data-net]");
    if (!t) return;
    // optional: soft preview without locking selection
  });
})();
```

Minimum behavior:

1. Click element with `data-net` → highlight all matching; dim others.
2. Click empty diagram chrome/background → clear.
3. Keyboard optional; mouse sufficient for v1.

## Blocks

Prefer rounded rects for modules; keep port circles on the **edge** facing the peer. Title text inside the block, not over ports.

## Legend (optional)

List nets with swatches. Clicking a legend row should call the same `highlight(net)` path.

## Accessibility

- `role="img"` + `aria-label` on SVG
- Visible hint in header for click-to-trace
- Do not rely on color alone: `is-hi` also increases stroke width / opacity
