/**
 * Budgeted cold-start context pack (`spacecraft context`).
 */

import { existsSync, readdirSync, readFileSync, statSync } from 'node:fs';
import path from 'node:path';
import { formatStatus } from './mission.mjs';

const DEFAULT_BUDGET = 4096;
const ENV_BUDGET = 'SPACECRAFT_CONTEXT_BUDGET';
const TRUNC_MARKER = '\n...[truncated]';
const LESSONS_LINE_CAP = 20;

function hasHelpFlag(args) {
  return args.some((a) => a === '--help' || a === '-h');
}

function printContextHelp() {
  console.log('Usage: spacecraft context [--budget <n>]');
  console.log('');
  console.log('Print a budgeted cold-start context pack.');
  console.log('Read order: docs/ (product) then .space/ (local mission/runtime).');
  console.log('Missing paths are omitted; exit 0 for normal missing-path cases.');
  console.log('');
  console.log('Options:');
  console.log(`  --budget <n>              Cap pack length in chars (default ${DEFAULT_BUDGET})`);
  console.log('  --help, -h                Show this help');
  console.log('');
  console.log(
    `${ENV_BUDGET}=<n>  Same as --budget when the flag is absent (default ${DEFAULT_BUDGET})`,
  );
}

/** Section: `## <heading>\n<body>\n` */
function section(heading, body) {
  return `## ${heading}\n${body}\n`;
}

function readUtf8IfFile(filePath) {
  try {
    if (!existsSync(filePath) || !statSync(filePath).isFile()) return null;
    return readFileSync(filePath, 'utf8');
  } catch {
    return null;
  }
}

/** Relative paths under conventionsDir, POSIX separators, lex sort. */
function listConventionRels(conventionsDir) {
  const rels = [];
  function walk(absDir, relBase) {
    let entries;
    try {
      entries = readdirSync(absDir, { withFileTypes: true });
    } catch {
      return;
    }
    for (const ent of entries) {
      const rel = relBase ? `${relBase}/${ent.name}` : ent.name;
      const abs = path.join(absDir, ent.name);
      if (ent.isDirectory()) {
        walk(abs, rel);
      } else if (ent.isFile()) {
        rels.push(rel);
      }
    }
  }
  walk(conventionsDir, '');
  rels.sort((a, b) => (a < b ? -1 : a > b ? 1 : 0));
  return rels;
}

function topLines(text, maxLines) {
  const parts = text.split('\n');
  return parts.slice(0, maxLines).join('\n');
}

function assemblePack(projectRoot, spaceDir, mid) {
  const parts = [];

  const readmePath = path.join(projectRoot, 'docs', 'README.md');
  const readmeBody = readUtf8IfFile(readmePath);
  if (readmeBody !== null) {
    parts.push(section('docs/README.md', readmeBody));
  }

  const conventionsDir = path.join(projectRoot, 'docs', 'conventions');
  if (existsSync(conventionsDir)) {
    for (const rel of listConventionRels(conventionsDir)) {
      const body = readUtf8IfFile(path.join(conventionsDir, ...rel.split('/')));
      if (body === null) continue;
      parts.push(section(`docs/conventions/${rel}`, body));
    }
  }

  const status = formatStatus(spaceDir, mid);
  if (status.ok) {
    parts.push(section('spacecraft status', status.text));
  }

  const lessonsPath = path.join(spaceDir, 'trust', 'lessons.md');
  const lessonsRaw = readUtf8IfFile(lessonsPath);
  if (lessonsRaw !== null) {
    parts.push(section('.space/trust/lessons.md', topLines(lessonsRaw, LESSONS_LINE_CAP)));
  }

  return parts.join('');
}

function applyBudget(assembled, budget) {
  if (assembled.length <= budget) return assembled;
  if (budget < TRUNC_MARKER.length) return TRUNC_MARKER.slice(0, budget);
  return assembled.slice(0, budget - TRUNC_MARKER.length) + TRUNC_MARKER;
}

/** Valid budget token: integer string with n >= 1. */
function parseBudgetToken(raw) {
  if (typeof raw !== 'string' || raw === '' || !/^\d+$/.test(raw)) return null;
  const n = Number(raw);
  if (!Number.isInteger(n) || n < 1) return null;
  return n;
}

/**
 * Precedence: --budget > SPACECRAFT_CONTEXT_BUDGET > default 4096.
 * @param {string[]} args
 * @param {NodeJS.ProcessEnv} [env]
 * @returns {{ ok: true, budget: number } | { ok: false, message: string }}
 */
function resolveBudget(args, env = process.env) {
  const flagIdx = args.indexOf('--budget');
  if (flagIdx !== -1) {
    const next = args[flagIdx + 1];
    if (next == null || next.startsWith('-')) {
      return { ok: false, message: 'spacecraft context: invalid budget (missing value)' };
    }
    const n = parseBudgetToken(next);
    if (n === null) {
      return { ok: false, message: 'spacecraft context: invalid budget' };
    }
    return { ok: true, budget: n };
  }

  if (env[ENV_BUDGET] !== undefined) {
    const n = parseBudgetToken(env[ENV_BUDGET]);
    if (n === null) {
      return { ok: false, message: 'spacecraft context: invalid budget' };
    }
    return { ok: true, budget: n };
  }

  return { ok: true, budget: DEFAULT_BUDGET };
}

/**
 * CLI entry for `spacecraft context`.
 * @param {string[]} args
 * @param {string} spaceDir
 * @param {string} [_cwd]
 * @param {string|null} [mid]
 * @returns {number} exit code
 */
export function contextCmd(args, spaceDir, _cwd, mid) {
  if (hasHelpFlag(args)) {
    printContextHelp();
    return 0;
  }

  const resolved = resolveBudget(args);
  if (!resolved.ok) {
    console.error(resolved.message);
    return 2;
  }

  const projectRoot = path.dirname(spaceDir);
  const assembled = assemblePack(projectRoot, spaceDir, mid ?? null);
  const output = applyBudget(assembled, resolved.budget);
  process.stdout.write(output);
  return 0;
}
