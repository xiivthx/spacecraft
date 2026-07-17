#!/usr/bin/env node
/**
 * visual-verify.mjs - Tier 3 browser visual verification for sc-ux-design.
 *
 * Usage: node scripts/visual-verify.mjs <html-file>
 * Output: JSON to stdout with breakpoint results, issues, and screenshot paths.
 */

import { chromium } from '@playwright/test';
import { writeFileSync, mkdirSync, existsSync } from 'fs';
import { resolve, dirname, basename, extname } from 'path';
import { fileURLToPath } from 'url';

const __dirname = dirname(fileURLToPath(import.meta.url));
const SKILL_DIR = resolve(__dirname, '..');
const SCREENSHOTS_DIR = resolve(SKILL_DIR, 'test', 'screenshots');

const BREAKPOINTS = [
  { name: 'mobile', width: 375, height: 812 },
  { name: 'tablet', width: 768, height: 1024 },
  { name: 'desktop', width: 1280, height: 900 },
];

function usage() {
  console.error('Usage: node scripts/visual-verify.mjs <html-file>');
  process.exit(1);
}

async function detectIssues(page) {
  const issues = [];

  // Horizontal overflow: elements wider than viewport
  const overflowing = await page.evaluate(() => {
    const results = [];
    const vw = window.innerWidth;
    document.querySelectorAll('*').forEach((el) => {
      const rect = el.getBoundingClientRect();
      if (rect.right > vw + 2 || rect.left < -2) {
        results.push({
          selector: el.tagName + (el.className ? '.' + el.className.split(' ')[0] : ''),
          kind: 'horizontal-overflow',
          detail: `Element extends beyond viewport (right: ${Math.round(rect.right)}px, vw: ${vw}px)`,
        });
      }
    });
    return results;
  });
  overflowing.forEach((i) => issues.push({ ...i, severity: 'error' }));

  // Cramped padding: elements with very small padding in bordered/colored containers
  const cramped = await page.evaluate(() => {
    const results = [];
    document.querySelectorAll('*').forEach((el) => {
      const style = window.getComputedStyle(el);
      const padTop = parseFloat(style.paddingTop);
      const padBottom = parseFloat(style.paddingBottom);
      const padLeft = parseFloat(style.paddingLeft);
      const bg = style.backgroundColor;
      const border = style.borderWidth !== '0px' || style.borderTopWidth !== '0px';
      const hasVisibleBg = bg && bg !== 'rgba(0, 0, 0, 0)' && bg !== 'transparent';
      if (
        (hasVisibleBg || border) &&
        el.textContent.trim().length > 10 &&
        (padTop < 8 || padBottom < 8 || padLeft < 8) &&
        (padTop > 0 || padBottom > 0 || padLeft > 0)
      ) {
        const minPad = Math.min(padTop, padBottom, padLeft);
        results.push({
          selector: el.tagName + (el.className ? '.' + el.className.split(' ')[0] : ''),
          kind: 'cramped-padding',
          detail: `Padding too small: min ${minPad}px (recommend ≥8px)`,
        });
      }
    });
    return results;
  });
  cramped.forEach((i) => issues.push({ ...i, severity: 'warning' }));

  // Body text touching viewport edge: body/paragraphs flush against edge
  const edgeText = await page.evaluate(() => {
    const results = [];
    document.querySelectorAll('p, h1, h2, h3, h4, h5, h6, span, div').forEach((el) => {
      if (el.textContent.trim().length < 20) return;
      const style = window.getComputedStyle(el);
      const marginLeft = parseFloat(style.marginLeft);
      const paddingLeft = parseFloat(style.paddingLeft);
      const left = el.getBoundingClientRect().left;
      const right = el.getBoundingClientRect().right;
      const vw = window.innerWidth;
      if (left <= 4 && marginLeft <= 4 && paddingLeft <= 4) {
        results.push({
          selector: el.tagName + (el.className ? '.' + el.className.split(' ')[0] : ''),
          kind: 'text-touching-viewport-edge',
          detail: `Content flush against left viewport edge (left: ${Math.round(left)}px)`,
        });
      }
      if (right >= vw - 4) {
        // check that it isn't full-width on purpose (like hero sections)
        const parentRect = el.parentElement?.getBoundingClientRect();
        if (parentRect && parentRect.left <= 4) {
          results.push({
            selector: el.tagName + (el.className ? '.' + el.className.split(' ')[0] : ''),
            kind: 'text-touching-viewport-edge',
            detail: `Content flush against right viewport edge (right: ${Math.round(right)}px, vw: ${vw}px)`,
          });
        }
      }
    });
    return results;
  });
  edgeText.forEach((i) => issues.push({ ...i, severity: 'warning' }));

  // Gradient detection: purple/blue gradients (AI palette tell)
  const gradients = await page.evaluate(() => {
    const results = [];
    document.querySelectorAll('*').forEach((el) => {
      const bg = window.getComputedStyle(el).backgroundImage;
      if (bg && bg.includes('gradient')) {
        // Check for purple-blue-ish gradients
        const lower = bg.toLowerCase();
        if (
          (lower.includes('7c3aed') || lower.includes('3b82f6') || lower.includes('purple') || lower.includes('blue') || lower.includes('violet')) ||
          lower.includes('135deg') || lower.includes('to right') || lower.includes('to bottom')
        ) {
          results.push({
            selector: el.tagName + (el.className ? '.' + el.className.split(' ')[0] : ''),
            kind: 'ai-gradient',
            detail: 'Purple/blue gradient detected - AI slop tell',
          });
        }
      }
    });
    return results;
  });
  gradients.forEach((i) => issues.push({ ...i, severity: 'warning' }));

  // Nested cards: card-like elements inside card-like elements
  const nestedCards = await page.evaluate(() => {
    const results = [];
    const isCard = (el) => {
      const style = window.getComputedStyle(el);
      const hasShadow = style.boxShadow !== 'none';
      const hasRadius = parseFloat(style.borderRadius) > 4;
      const hasBg = style.backgroundColor !== 'rgba(0, 0, 0, 0)' && style.backgroundColor !== 'transparent';
      return hasShadow || (hasRadius && hasBg);
    };
    // Find card elements with card children
    const allCards = Array.from(document.querySelectorAll('*')).filter(isCard);
    for (let i = 0; i < allCards.length; i++) {
      const card = allCards[i];
      const children = Array.from(card.children).filter(isCard);
      if (children.length > 0) {
        results.push({
          selector: card.tagName + (card.className ? '.' + card.className.split(' ')[0] : ''),
          kind: 'nested-cards',
          detail: `Card contains ${children.length} nested card(s) - flatten hierarchy`,
        });
      }
    }
    return results;
  });
  nestedCards.forEach((i) => issues.push({ ...i, severity: 'warning' }));

  // Cream/beige background detection
  const creamBg = await page.evaluate(() => {
    const results = [];
    const bg = window.getComputedStyle(document.body).backgroundColor;
    // Check for warm cream/beige colors
    const match = bg.match(/rgba?\((\d+),\s*(\d+),\s*(\d+)/);
    if (match) {
      const r = parseInt(match[1]), g = parseInt(match[2]), b = parseInt(match[3]);
      // Cream/beige: warm off-white, higher red, lower blue
      if (r > 220 && g > 210 && b > 190 && r >= b + 10 && (r - b) < 60) {
        results.push({
          selector: 'body',
          kind: 'cream-palette',
          detail: `Warm cream/beige background (rgb(${r},${g},${b})) - AI "tasteful" default`,
        });
      }
    }
    return results;
  });
  creamBg.forEach((i) => issues.push({ ...i, severity: 'info' }));

  return issues;
}

async function main() {
  const args = process.argv.slice(2);
  if (args.length < 1) usage();

  const htmlPath = resolve(args[0]);
  if (!existsSync(htmlPath)) {
    console.error(`File not found: ${htmlPath}`);
    process.exit(1);
  }

  const baseName = basename(htmlPath, extname(htmlPath));
  mkdirSync(SCREENSHOTS_DIR, { recursive: true });

  const browser = await chromium.launch({ headless: true });
  const results = [];

  try {
    for (const bp of BREAKPOINTS) {
      const context = await browser.newContext({
        viewport: { width: bp.width, height: bp.height },
      });
      const page = await context.newPage();
      await page.goto(`file://${htmlPath}`, { waitUntil: 'networkidle' });

      const screenshotPath = resolve(SCREENSHOTS_DIR, `${baseName}-${bp.name}-${bp.width}.png`);
      await page.screenshot({ path: screenshotPath, fullPage: true });

      const issues = await detectIssues(page);

      results.push({
        breakpoint: bp.width,
        issues,
        screenshots: [screenshotPath],
      });

      await context.close();
    }
  } finally {
    await browser.close();
  }

  const report = {
    file: htmlPath,
    breakpoints: BREAKPOINTS.map((b) => b.width),
    results,
  };

  console.log(JSON.stringify(report, null, 2));
}

main().catch((err) => {
  console.error('Visual verification failed:', err.message);
  process.exit(1);
});
