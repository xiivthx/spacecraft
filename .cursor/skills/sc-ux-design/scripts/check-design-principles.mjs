#!/usr/bin/env node
/**
 * check-design-principles.mjs — hard-fail aesthetic calculation checker.
 *
 * Checks:
 * 1. Root DESIGN.md YAML spacing + typography fontSize consecutive ratios
 *    in [1.2, 1.618] (±0.001) OR exact 2.0 doubling (±0.001).
 * 2. Aesthetic fixture φ primary split (≈1.618 ± 0.08) and DCM within ±8%.
 *
 * Exit 0 only when all pass; else exit 1 with failures on stderr.
 */

import { readFileSync, existsSync } from "node:fs";
import { dirname, resolve } from "node:path";
import { fileURLToPath } from "node:url";

const __dirname = dirname(fileURLToPath(import.meta.url));
const SKILL_DIR = resolve(__dirname, "..");
const FIXTURE = resolve(SKILL_DIR, "test", "fixture-design-principles.html");

/** Resolve repo-root DESIGN.md from cwd or by walking up from this script. */
function findDesignMd() {
  const fromCwd = resolve(process.cwd(), "DESIGN.md");
  if (existsSync(fromCwd)) return fromCwd;
  let dir = __dirname;
  for (let i = 0; i < 8; i++) {
    const candidate = resolve(dir, "DESIGN.md");
    if (existsSync(candidate)) return candidate;
    const parent = resolve(dir, "..");
    if (parent === dir) break;
    dir = parent;
  }
  return fromCwd;
}

const DESIGN_MD = findDesignMd();

const RATIO_MIN = 1.2;
const RATIO_MAX = 1.618;
const RATIO_TOL = 0.001;
const DOUBLING = 2.0;
const PHI = 1.618;
const PHI_TOL = 0.08;
const DCM_FRAC = 0.08;

const failures = [];

function fail(msg) {
  failures.push(msg);
}

function parsePx(raw) {
  const m = String(raw).trim().match(/^(\d+(?:\.\d+)?)\s*px$/i);
  return m ? Number(m[1]) : null;
}

/** Extract YAML front matter between first --- fences. */
function extractFrontMatter(text) {
  const m = text.match(/^---\r?\n([\s\S]*?)\r?\n---/);
  if (!m) return null;
  return m[1];
}

/**
 * Collect spacing px values and typography fontSize px values from
 * indentation-based YAML (no external parser).
 */
function parseDesignScales(yaml) {
  const lines = yaml.split(/\r?\n/);
  const spacing = [];
  const fontSizes = [];
  let section = null; // 'spacing' | 'typography' | null

  for (const line of lines) {
    if (/^\s*#/.test(line) || !line.trim()) continue;

    const spacingKey = line.match(/^spacing:\s*$/);
    if (spacingKey) {
      section = "spacing";
      continue;
    }
    const typoKey = line.match(/^typography:\s*$/);
    if (typoKey) {
      section = "typography";
      continue;
    }
    // Top-level key ends current section
    if (/^[a-zA-Z_][\w-]*:\s*/.test(line) && !/^\s/.test(line)) {
      if (section === "spacing" || section === "typography") {
        section = null;
      }
    }

    if (section === "spacing") {
      const sm = line.match(/^\s+[a-zA-Z_][\w-]*:\s*(.+)\s*$/);
      if (sm) {
        const px = parsePx(sm[1].replace(/["']/g, ""));
        if (px != null) spacing.push(px);
      }
      continue;
    }

    if (section === "typography") {
      const fs = line.match(/^\s+fontSize:\s*(.+)\s*$/);
      if (fs) {
        const px = parsePx(fs[1].replace(/["']/g, ""));
        if (px != null) fontSizes.push(px);
      }
    }
  }

  return {
    spacing: [...new Set(spacing)].sort((a, b) => a - b),
    fontSizes: [...new Set(fontSizes)].sort((a, b) => a - b),
  };
}

function ratioAllowed(ratio) {
  if (ratio >= RATIO_MIN - RATIO_TOL && ratio <= RATIO_MAX + RATIO_TOL) return true;
  if (Math.abs(ratio - DOUBLING) <= RATIO_TOL) return true;
  return false;
}

function checkScale(label, values) {
  if (values.length < 2) {
    fail(`${label}: need at least 2 scale steps, got ${values.length} (${values.join(", ") || "none"})`);
    return;
  }
  for (let i = 1; i < values.length; i++) {
    const smaller = values[i - 1];
    const larger = values[i];
    if (smaller <= 0) {
      fail(`${label}: non-positive step ${smaller}`);
      continue;
    }
    const ratio = larger / smaller;
    if (!ratioAllowed(ratio)) {
      fail(
        `${label}: ratio ${smaller}→${larger} = ${ratio.toFixed(4)} not in [${RATIO_MIN}, ${RATIO_MAX}] (±${RATIO_TOL}) and not doubling ${DOUBLING} (±${RATIO_TOL})`
      );
    }
  }
}

function checkDesignMd() {
  if (!existsSync(DESIGN_MD)) {
    fail(`DESIGN.md not found at ${DESIGN_MD}`);
    return;
  }
  const text = readFileSync(DESIGN_MD, "utf8");
  const yaml = extractFrontMatter(text);
  if (!yaml) {
    fail("DESIGN.md: missing YAML front matter between --- fences");
    return;
  }
  const { spacing, fontSizes } = parseDesignScales(yaml);
  checkScale("spacing", spacing);
  checkScale("typography fontSize", fontSizes);
}

/**
 * Read φ split from data-phi-a / data-phi-b attributes or --phi-a / --phi-b CSS vars.
 */
function extractPhiPair(html) {
  const dataA = html.match(/data-phi-a\s*=\s*["']([^"']+)["']/i);
  const dataB = html.match(/data-phi-b\s*=\s*["']([^"']+)["']/i);
  if (dataA && dataB) {
    const a = Number(String(dataA[1]).replace(/px$/i, ""));
    const b = Number(String(dataB[1]).replace(/px$/i, ""));
    if (Number.isFinite(a) && Number.isFinite(b) && a > 0 && b > 0) {
      return { a, b, source: "data-phi-a/b" };
    }
  }

  const varA = html.match(/--phi-a\s*:\s*([\d.]+)\s*px/i);
  const varB = html.match(/--phi-b\s*:\s*([\d.]+)\s*px/i);
  if (varA && varB) {
    const a = Number(varA[1]);
    const b = Number(varB[1]);
    if (Number.isFinite(a) && Number.isFinite(b) && a > 0 && b > 0) {
      return { a, b, source: "--phi-a/--phi-b" };
    }
  }
  return null;
}

function extractDcmJson(html) {
  const m = html.match(
    /<script[^>]*\bid=["']dcm-weights["'][^>]*type=["']application\/json["'][^>]*>([\s\S]*?)<\/script>/i
  ) || html.match(
    /<script[^>]*type=["']application\/json["'][^>]*\bid=["']dcm-weights["'][^>]*>([\s\S]*?)<\/script>/i
  );
  if (!m) return null;
  try {
    return JSON.parse(m[1].trim());
  } catch (e) {
    fail(`fixture DCM JSON parse error: ${e.message}`);
    return null;
  }
}

function checkPhiSplit(html) {
  const pair = extractPhiPair(html);
  if (!pair) {
    fail("fixture: missing data-phi-a/data-phi-b or --phi-a/--phi-b widths");
    return;
  }
  const lo = Math.min(pair.a, pair.b);
  const hi = Math.max(pair.a, pair.b);
  const ratio = hi / lo;
  if (Math.abs(ratio - PHI) > PHI_TOL) {
    fail(
      `fixture φ split (${pair.source}): max/min = ${ratio.toFixed(4)} not within ${PHI} ± ${PHI_TOL} (got ${hi}/${lo})`
    );
  }
}

function checkDcm(html) {
  const data = extractDcmJson(html);
  if (!data) {
    if (!failures.some((f) => f.includes("DCM JSON"))) {
      fail('fixture: missing <script type="application/json" id="dcm-weights">');
    }
    return;
  }
  const { width: W, height: H, rects } = data;
  if (!(W > 0) || !(H > 0) || !Array.isArray(rects) || rects.length === 0) {
    fail("fixture DCM: need positive width/height and non-empty rects");
    return;
  }

  let mass = 0;
  let mx = 0;
  let my = 0;
  for (const r of rects) {
    const w = Number(r.w);
    const h = Number(r.h);
    const x = Number(r.x);
    const y = Number(r.y);
    const weight = Number(r.weight);
    if (![w, h, x, y, weight].every(Number.isFinite) || weight < 0) {
      fail(`fixture DCM: invalid rect ${JSON.stringify(r)}`);
      return;
    }
    const cx = x + w / 2;
    const cy = y + h / 2;
    mass += weight;
    mx += weight * cx;
    my += weight * cy;
  }
  if (mass <= 0) {
    fail("fixture DCM: total weight must be > 0");
    return;
  }
  const cmX = mx / mass;
  const cmY = my / mass;
  const dx = Math.abs(cmX - W / 2) / (W / 2);
  const dy = Math.abs(cmY - H / 2) / (H / 2);
  if (dx > DCM_FRAC || dy > DCM_FRAC) {
    fail(
      `fixture DCM: center of mass offset too large (dx=${(dx * 100).toFixed(2)}%, dy=${(dy * 100).toFixed(2)}%; limit ${DCM_FRAC * 100}% of half-extent). cm=(${cmX.toFixed(2)}, ${cmY.toFixed(2)}) center=(${(W / 2).toFixed(2)}, ${(H / 2).toFixed(2)})`
    );
  }
}

function checkFixture() {
  if (!existsSync(FIXTURE)) {
    fail(`fixture not found at ${FIXTURE}`);
    return;
  }
  const html = readFileSync(FIXTURE, "utf8");
  checkPhiSplit(html);
  checkDcm(html);
}

checkDesignMd();
checkFixture();

if (failures.length > 0) {
  for (const f of failures) {
    console.error(`FAIL: ${f}`);
  }
  console.error(`${failures.length} check(s) failed`);
  process.exit(1);
}

console.log("PASS: DESIGN.md modular scale + fixture φ split + DCM");
process.exit(0);
