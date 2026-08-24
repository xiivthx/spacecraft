/**
 * Docs ↔ mission drift report (`spacecraft drift`).
 * First-ship: four rules; skip-vs-find degrade; default exit 0; --strict non-zero-on-findings.
 */

import { existsSync, readdirSync, readFileSync, statSync } from 'node:fs';
import path from 'node:path';
import { missionDir, resolveActive } from './resolve.mjs';

/** Explicit docs/specs/* or docs/epics/* path tokens in mission prose. */
const DOCS_REF_RE = /docs\/(?:specs|epics)\/[A-Za-z0-9._/-]+/g;

/** Docs map claim: conventions/ or docs/conventions mentioned. */
const CONVENTIONS_CLAIM_RE = /(?:docs\/)?conventions\//;

/** ATX heading whose title starts with Goal or Verify. */
const GOAL_HEADING_RE = /^#{1,6}\s+Goal\b/m;
const VERIFY_HEADING_RE = /^#{1,6}\s+Verify\b/m;

/** Optional orphan trees (Rule 4). */
const ORPHAN_KINDS = ['epics', 'specs'];

function hasHelpFlag(args) {
  return args.some((a) => a === '--help' || a === '-h');
}

function hasStrictFlag(args) {
  return args.some((a) => a === '--strict');
}

function printDriftHelp() {
  console.log('Usage: spacecraft drift [--strict]');
  console.log('');
  console.log('Report docs ↔ mission drift on stdout (read-only; does not write docs/).');
  console.log('');
  console.log('First-ship rules:');
  console.log(
    '  1. Goal/Verify heading: when clarify-status is clear, current mission spec.md',
  );
  console.log('     must have a Goal heading and a Verify heading.');
  console.log(
    '  2. Broken refs: mission spec.md / decisions.md citing missing docs/specs/*',
  );
  console.log('     or docs/epics/* paths.');
  console.log(
    '  3. Conventions claim: docs/conventions/ missing while docs map / context claims',
  );
  console.log('     conventions.');
  console.log(
    '  4. Orphan: file under a non-empty docs/epics/ or docs/specs/ tree never',
  );
  console.log('     referenced from current mission spec.md / decisions.md.');
  console.log('');
  console.log('Skip vs finding:');
  console.log(
    '  Absent optional docs trees (e.g. docs/epics, docs/specs) emit an explicit',
  );
  console.log('  skip line for that rule. Skips are not findings.');
  console.log('');
  console.log('Exit:');
  console.log('  default exit 0 for a completed report (including findings).');
  console.log('  --strict exits non-zero when there is at least one finding');
  console.log('  (skips alone do not fail --strict).');
  console.log('');
  console.log('Options:');
  console.log('  --strict                 Non-zero exit on findings');
  console.log('  --help, -h               Show this help');
}

function readUtf8(filePath) {
  try {
    return readFileSync(filePath, 'utf8');
  } catch {
    return null;
  }
}

function readClarifyStatus(mDir) {
  const raw = readUtf8(path.join(mDir, 'clarify-status'));
  return raw === null ? '' : raw.trim();
}

/**
 * Rule 1: clear + missing Goal or Verify ATX heading → finding.
 * @param {string} mDir
 * @returns {string[]}
 */
function ruleClearGoalVerify(mDir) {
  if (readClarifyStatus(mDir) !== 'clear') return [];
  const spec = readUtf8(path.join(mDir, 'spec.md'));
  if (spec === null) {
    return [
      'finding: clarify-status clear but spec.md missing Goal/Verify heading (spec.md absent)',
    ];
  }
  const hasGoal = GOAL_HEADING_RE.test(spec);
  const hasVerify = VERIFY_HEADING_RE.test(spec);
  if (hasGoal && hasVerify) return [];
  const missing = [];
  if (!hasGoal) missing.push('Goal');
  if (!hasVerify) missing.push('Verify');
  return [
    `finding: clarify-status clear but spec.md missing Goal/Verify heading (${missing.join(', ')})`,
  ];
}

/**
 * Collect unique docs/specs|epics path refs from mission text files.
 * @param {string} mDir
 * @returns {string[]}
 */
function collectDocsRefs(mDir) {
  const texts = [];
  for (const name of ['spec.md', 'decisions.md']) {
    const body = readUtf8(path.join(mDir, name));
    if (body !== null) texts.push(body);
  }
  const found = new Set();
  for (const text of texts) {
    DOCS_REF_RE.lastIndex = 0;
    let m;
    while ((m = DOCS_REF_RE.exec(text)) !== null) {
      // Strip trailing punctuation that often follows a path in prose.
      const ref = m[0].replace(/[.,;:!?)]+$/, '');
      if (ref) found.add(ref);
    }
  }
  return [...found];
}

/**
 * Rule 2: explicit docs/specs/* or docs/epics/* refs that do not exist → finding.
 * @param {string} cwd
 * @param {string} mDir
 * @returns {string[]}
 */
function ruleBrokenDocsRefs(cwd, mDir) {
  const findings = [];
  for (const ref of collectDocsRefs(mDir)) {
    const abs = path.join(cwd, ref);
    if (!existsSync(abs)) {
      findings.push(`finding: broken docs ref ${ref}`);
    }
  }
  return findings;
}

/**
 * Whether docs map claims conventions (docs/README.md mentions conventions/).
 * @param {string} cwd
 * @returns {boolean}
 */
function docsMapClaimsConventions(cwd) {
  const readme = readUtf8(path.join(cwd, 'docs', 'README.md'));
  if (readme === null) return false;
  return CONVENTIONS_CLAIM_RE.test(readme);
}

/**
 * Rule 3: docs map claims conventions but docs/conventions/ missing → finding.
 * @param {string} cwd
 * @returns {string[]}
 */
function ruleConventionsClaim(cwd) {
  if (!docsMapClaimsConventions(cwd)) return [];
  const convDir = path.join(cwd, 'docs', 'conventions');
  try {
    if (existsSync(convDir) && statSync(convDir).isDirectory()) return [];
  } catch {
    // treat as missing
  }
  return [
    'finding: docs/conventions/ missing while docs map claims conventions',
  ];
}

/**
 * List files under docs/<kind>/ as posix paths docs/<kind>/….
 * @param {string} absDir
 * @param {string} kind
 * @returns {string[]}
 */
function listTreeFiles(absDir, kind) {
  const out = [];
  function walk(abs, relBase) {
    let entries;
    try {
      entries = readdirSync(abs, { withFileTypes: true });
    } catch {
      return;
    }
    for (const ent of entries) {
      const rel = relBase ? `${relBase}/${ent.name}` : ent.name;
      const child = path.join(abs, ent.name);
      if (ent.isDirectory()) {
        walk(child, rel);
      } else if (ent.isFile()) {
        out.push(`docs/${kind}/${rel}`);
      }
    }
  }
  walk(absDir, '');
  out.sort((a, b) => (a < b ? -1 : a > b ? 1 : 0));
  return out;
}

/**
 * Rule 4: orphans under non-empty docs/epics|specs; skip when tree absent/empty.
 * @param {string} cwd
 * @param {string} mDir
 * @returns {{ findings: string[], skips: string[] }}
 */
function ruleOrphanAndSkip(cwd, mDir) {
  const findings = [];
  const skips = [];
  const refs = new Set(collectDocsRefs(mDir));

  for (const kind of ORPHAN_KINDS) {
    const abs = path.join(cwd, 'docs', kind);
    let isDir = false;
    try {
      isDir = existsSync(abs) && statSync(abs).isDirectory();
    } catch {
      isDir = false;
    }
    if (!isDir) {
      skips.push(`skip: docs/${kind} tree absent`);
      continue;
    }
    const files = listTreeFiles(abs, kind);
    if (files.length === 0) {
      skips.push(`skip: docs/${kind} tree empty`);
      continue;
    }
    for (const file of files) {
      if (!refs.has(file)) {
        findings.push(
          `finding: orphan ${file} not referenced from mission spec.md / decisions.md`,
        );
      }
    }
  }

  return { findings, skips };
}

/**
 * CLI entry for `spacecraft drift`.
 * @param {string[]} args
 * @param {string} spaceDir
 * @param {string} cwd
 * @param {string|null} mid
 * @returns {number} exit code
 */
export function driftCmd(args, spaceDir, cwd, mid) {
  if (hasHelpFlag(args)) {
    printDriftHelp();
    return 0;
  }

  const strict = hasStrictFlag(args);
  const id = resolveActive(spaceDir, mid);
  if (!id) {
    console.error(
      "spacecraft drift: no active mission - pass a mission via branch or 'use'",
    );
    return 1;
  }

  const mDir = missionDir(spaceDir, id);
  const orphan = ruleOrphanAndSkip(cwd, mDir);
  const findings = [
    ...ruleClearGoalVerify(mDir),
    ...ruleBrokenDocsRefs(cwd, mDir),
    ...ruleConventionsClaim(cwd),
    ...orphan.findings,
  ];
  const skips = orphan.skips;

  for (const line of findings) {
    console.log(line);
  }
  for (const line of skips) {
    console.log(line);
  }

  if (strict && findings.length >= 1) return 1;
  return 0;
}
