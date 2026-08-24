/**
 * spacecraft context — T1–T4 (M9C0ZW13).
 *
 * Oracles: design-contract Public seams + Edge matrix; approved-scenarios frozen.
 *
 * Pbt skipped: not core logic
 */
import assert from 'node:assert/strict';
import { spawnSync } from 'node:child_process';
import {
  chmodSync,
  mkdirSync,
  mkdtempSync,
  readFileSync,
  rmSync,
  symlinkSync,
  writeFileSync,
} from 'node:fs';
import os from 'node:os';
import path from 'node:path';
import test from 'node:test';
import { fileURLToPath } from 'node:url';

const repoRoot = path.resolve(path.dirname(fileURLToPath(import.meta.url)), '../..');
const entryPath = path.join(repoRoot, 'cli', 'spacecraft.mjs');
const sessionStartHook = path.join(repoRoot, '.cursor', 'hooks', 'session-start.sh');

const HEADING_README = '## docs/README.md';
const HEADING_STATUS = '## spacecraft status';
/** design-contract section 4 — lessons when `.space/trust/lessons.md` exists. */
const HEADING_LESSONS = '## .space/trust/lessons.md';
/** design-contract: lessons body = split(`\\n`) indices 0..19 before budget trim. */
const LESSONS_LINE_CAP = 20;
const NO_MISSION_LINE = 'No selected mission.';
/** Hook fallback when spacecraft binary is absent (session-start.sh). */
const HOOK_ABSENT_FALLBACK = 'No active spacecraft mission.';
/** design-contract truncation marker (substring `truncated` required). */
const TRUNC_MARKER = '\n...[truncated]';
const ENV_BUDGET = 'SPACECRAFT_CONTEXT_BUDGET';
/** Oversized fixture body so assembled pack exceeds T3 budget caps. */
const OVERSIZE_README = `${'X'.repeat(800)}\n`;

/** Inject command must be context (design-contract Hook / plan T4 verify). */
const HOOK_CONTEXT_INJECT_RE = /(?:\.\/spacecraft|spacecraft)\s+context\b/;
/** status-alone inject line (plan verify must not match after GREEN). */
const HOOK_STATUS_INJECT_RE = /(?:\.\/spacecraft|spacecraft)\s+status\s*$/m;

function runContextHelp() {
  return spawnSync(process.execPath, [entryPath, 'context', '--help'], {
    encoding: 'utf8',
    cwd: repoRoot,
  });
}

function spaceRoot() {
  const dir = mkdtempSync(path.join(os.tmpdir(), 'spacecraft-context-'));
  mkdirSync(path.join(dir, '.space', 'missions'), { recursive: true });
  mkdirSync(path.join(dir, '.space', 'roadmaps'), { recursive: true });
  return dir;
}

function cleanup(dir) {
  rmSync(dir, { recursive: true, force: true });
}

function writeMissionWithPickup(root, id, pickupNext) {
  const dir = path.join(root, '.space', 'missions', id);
  mkdirSync(dir, { recursive: true });
  writeFileSync(
    path.join(dir, 'mission.json'),
    `${JSON.stringify(
      {
        id,
        title: 'Context Pack Mission',
        state: 'active',
        createdAt: '2026-01-01T00:00:00Z',
        branches: [],
        pickup: {
          phase: 'run',
          next: pickupNext,
          updatedAt: '2026-08-10T00:00:00Z',
        },
      },
      null,
      2,
    )}\n`,
  );
  writeFileSync(path.join(dir, 'spec.md'), '# Spec\n');
  writeFileSync(
    path.join(dir, 'plan.json'),
    `${JSON.stringify({ planName: 'test', missionId: id, tasks: [] }, null, 2)}\n`,
  );
  writeFileSync(path.join(dir, 'evidence.jsonl'), '');
  writeFileSync(path.join(root, '.space', 'current'), `${id}\n`);
}

/**
 * @param {string} cwd
 * @param {string[]} [extraArgs]
 * @param {NodeJS.ProcessEnv} [envOverrides] merged over process.env when provided
 */
function runContext(cwd, extraArgs = [], envOverrides) {
  const result = spawnSync(process.execPath, [entryPath, 'context', ...extraArgs], {
    cwd,
    encoding: 'utf8',
    env: envOverrides ? { ...process.env, ...envOverrides } : process.env,
  });
  return {
    stdout: result.stdout ?? '',
    stderr: result.stderr ?? '',
    code: result.status ?? 1,
  };
}

/**
 * @param {string} cwd
 * @param {NodeJS.ProcessEnv} [envOverrides]
 */
function runSessionStartHook(cwd, envOverrides) {
  const result = spawnSync('sh', [sessionStartHook], {
    cwd,
    encoding: 'utf8',
    env: envOverrides ? { ...process.env, ...envOverrides } : process.env,
  });
  return {
    stdout: result.stdout ?? '',
    stderr: result.stderr ?? '',
    code: result.status ?? 1,
  };
}

/**
 * T1 acceptance 2: context --help exits 0 and documents read order, default
 * budget 4096, --budget, and SPACECRAFT_CONTEXT_BUDGET (design-contract Public seams).
 */
test('spacecraft context --help exits 0 and documents pack budget seams', () => {
  const help = runContextHelp();
  const helpOut = `${help.stdout ?? ''}${help.stderr ?? ''}`;

  assert.equal(
    help.status,
    0,
    `context --help must exit 0\nstderr=${help.stderr}\nstdout=${help.stdout}`,
  );

  // Read order: docs/ then .space/ (Public seams Help text Must mention)
  const docsIdx = helpOut.indexOf('docs/');
  const spaceIdx = helpOut.indexOf('.space/');
  assert.ok(docsIdx !== -1, `context --help must mention docs/\n${helpOut}`);
  assert.ok(spaceIdx !== -1, `context --help must mention .space/\n${helpOut}`);
  assert.ok(
    docsIdx < spaceIdx,
    `context --help must document read order docs/ then .space/\n${helpOut}`,
  );

  assert.match(
    helpOut,
    /4096/,
    `context --help must document default budget 4096\n${helpOut}`,
  );
  assert.match(
    helpOut,
    /--budget/,
    `context --help must document --budget\n${helpOut}`,
  );
  assert.match(
    helpOut,
    /SPACECRAFT_CONTEXT_BUDGET/,
    `context --help must document SPACECRAFT_CONTEXT_BUDGET\n${helpOut}`,
  );
});

/**
 * T2 pos-order-happy: docs README + conventions + mission status/Pickup →
 * README heading before status; conventions headings before status; exit 0.
 */
test('pos-order-happy: docs README and conventions headings appear before status', () => {
  const dir = spaceRoot();
  const id = 'M9CTX001';
  try {
    writeMissionWithPickup(dir, id, 'continue T2 pack order');
    mkdirSync(path.join(dir, 'docs', 'conventions'), { recursive: true });
    writeFileSync(path.join(dir, 'docs', 'README.md'), '# Product README\n');
    writeFileSync(path.join(dir, 'docs', 'conventions', 'naming.md'), '# Naming\n');

    const res = runContext(dir);
    assert.equal(res.code, 0, `context exit=${res.code}\nstderr=${res.stderr}\nstdout=${res.stdout}`);

    const out = res.stdout;
    const readmeIdx = out.indexOf(HEADING_README);
    const convIdx = out.indexOf('## docs/conventions/naming.md');
    const statusIdx = out.indexOf(HEADING_STATUS);

    assert.ok(readmeIdx !== -1, `stdout must include ${HEADING_README}\n${out}`);
    assert.ok(convIdx !== -1, `stdout must include ## docs/conventions/naming.md\n${out}`);
    assert.ok(statusIdx !== -1, `stdout must include ${HEADING_STATUS}\n${out}`);
    assert.ok(
      readmeIdx < statusIdx,
      `${HEADING_README} must appear before ${HEADING_STATUS}\n${out}`,
    );
    assert.ok(
      convIdx < statusIdx,
      'conventions headings must appear before ## spacecraft status\n' + out,
    );
  } finally {
    cleanup(dir);
  }
});

/**
 * T2 edge-missing-docs: no docs README/conventions, no lessons → exit 0;
 * status or no-mission line; no ## docs/README.md.
 */
test('edge-missing-docs: omit docs headings; exit 0 with status or no-mission', () => {
  const dir = spaceRoot();
  try {
    const res = runContext(dir);
    assert.equal(res.code, 0, `context exit=${res.code}\nstderr=${res.stderr}\nstdout=${res.stdout}`);

    const out = res.stdout;
    assert.equal(
      out.includes(HEADING_README),
      false,
      `stdout must not include ${HEADING_README}\n${out}`,
    );

    const hasStatusHeading = out.includes(HEADING_STATUS);
    const hasNoMission = out.includes(NO_MISSION_LINE);
    const hasMissionLine = /^Mission: /m.test(out);
    assert.ok(
      hasStatusHeading || hasNoMission || hasMissionLine,
      `stdout must include status section or "${NO_MISSION_LINE}" (or Mission: line)\n${out}`,
    );
  } finally {
    cleanup(dir);
  }
});

/**
 * T2 over-conventions-order: conventions a.md / b.md → lexicographic heading order; exit 0.
 */
test('over-conventions-order: conventions headings follow lexicographic rel-path', () => {
  const dir = spaceRoot();
  try {
    mkdirSync(path.join(dir, 'docs', 'conventions'), { recursive: true });
    writeFileSync(path.join(dir, 'docs', 'conventions', 'b.md'), '# B\n');
    writeFileSync(path.join(dir, 'docs', 'conventions', 'a.md'), '# A\n');

    const res = runContext(dir);
    assert.equal(res.code, 0, `context exit=${res.code}\nstderr=${res.stderr}\nstdout=${res.stdout}`);

    const out = res.stdout;
    const headingA = '## docs/conventions/a.md';
    const headingB = '## docs/conventions/b.md';
    const aIdx = out.indexOf(headingA);
    const bIdx = out.indexOf(headingB);

    assert.ok(aIdx !== -1, `stdout must include ${headingA}\n${out}`);
    assert.ok(bIdx !== -1, `stdout must include ${headingB}\n${out}`);
    assert.ok(
      aIdx < bIdx,
      `${headingA} must appear before ${headingB}\n${out}`,
    );
  } finally {
    cleanup(dir);
  }
});

/**
 * lessons-top20 (T2/review): `.space/trust/lessons.md` with >20 lines → heading present;
 * body = first 20 lines (split on `\n`, indices 0..19) before budget trim — not all 25.
 */
test('lessons-top20: lessons body capped to first 20 lines before budget trim', () => {
  const dir = spaceRoot();
  try {
    const lessonLines = Array.from(
      { length: 25 },
      (_, i) => `LESSON_LINE_${String(i).padStart(2, '0')}`,
    );
    const lessonsRaw = `${lessonLines.join('\n')}\n`;
    mkdirSync(path.join(dir, '.space', 'trust'), { recursive: true });
    writeFileSync(path.join(dir, '.space', 'trust', 'lessons.md'), lessonsRaw);

    // Oracle: design-contract Lessons body — split on \n, indices 0..19, rejoin.
    const expectedBody = lessonsRaw.split('\n').slice(0, LESSONS_LINE_CAP).join('\n');

    const res = runContext(dir, ['--budget', '100000']);
    assert.equal(res.code, 0, `context exit=${res.code}\nstderr=${res.stderr}\nstdout=${res.stdout}`);

    const out = res.stdout;
    const headingIdx = out.indexOf(HEADING_LESSONS);
    assert.ok(headingIdx !== -1, `stdout must include ${HEADING_LESSONS}\n${out}`);

    const afterHeading = out.slice(headingIdx + HEADING_LESSONS.length);
    assert.ok(
      afterHeading.startsWith('\n'),
      `lessons heading must be followed by newline then body\n${out}`,
    );
    const afterNewline = afterHeading.slice(1);
    const nextSec = afterNewline.indexOf('\n## ');
    const rawSection = nextSec === -1 ? afterNewline : afterNewline.slice(0, nextSec);
    const body = rawSection.endsWith('\n') ? rawSection.slice(0, -1) : rawSection;
    const bodyLines = body.split('\n');

    assert.equal(
      bodyLines.length,
      LESSONS_LINE_CAP,
      `lessons body must have exactly ${LESSONS_LINE_CAP} lines (got ${bodyLines.length}; fixture had 25)\nbody=${JSON.stringify(body)}\n${out}`,
    );
    assert.equal(
      body,
      expectedBody,
      `lessons body must equal split(\\n) indices 0..19 rejoined\nexpected=${JSON.stringify(expectedBody)}\nactual=${JSON.stringify(body)}\n${out}`,
    );
    assert.equal(
      body.includes('LESSON_LINE_20'),
      false,
      `lessons body must not include line index 20+ (strict cap)\nbody=${JSON.stringify(body)}\n${out}`,
    );
  } finally {
    cleanup(dir);
  }
});

/**
 * T3 neg-budget-trunc: oversize fixtures + `--budget 200` → exit 0;
 * output.length <= 200; contains `truncated` (marker `\n...[truncated]`).
 */
test('neg-budget-trunc: --budget 200 caps output and includes truncated', () => {
  const dir = spaceRoot();
  try {
    mkdirSync(path.join(dir, 'docs'), { recursive: true });
    writeFileSync(path.join(dir, 'docs', 'README.md'), OVERSIZE_README);

    const res = runContext(dir, ['--budget', '200']);
    assert.equal(res.code, 0, `context exit=${res.code}\nstderr=${res.stderr}\nstdout=${res.stdout}`);

    const out = res.stdout;
    assert.ok(
      out.length <= 200,
      `stdout.length must be <= 200 (got ${out.length})\n${out}`,
    );
    assert.ok(
      out.includes('truncated'),
      `stdout must include truncated (marker ${JSON.stringify(TRUNC_MARKER)})\n${out}`,
    );
  } finally {
    cleanup(dir);
  }
});

/**
 * T3 edge-budget-precedence: env SPACECRAFT_CONTEXT_BUDGET=1000 + `--budget 500`
 * → effective cap 500 (flag wins).
 */
test('edge-budget-precedence: --budget 500 wins over SPACECRAFT_CONTEXT_BUDGET=1000', () => {
  const dir = spaceRoot();
  try {
    mkdirSync(path.join(dir, 'docs'), { recursive: true });
    writeFileSync(path.join(dir, 'docs', 'README.md'), OVERSIZE_README);

    const res = runContext(dir, ['--budget', '500'], {
      [ENV_BUDGET]: '1000',
    });
    assert.equal(res.code, 0, `context exit=${res.code}\nstderr=${res.stderr}\nstdout=${res.stdout}`);

    const out = res.stdout;
    assert.ok(
      out.length <= 500,
      `effective budget must be 500 (stdout.length <= 500, got ${out.length})\n${out}`,
    );
  } finally {
    cleanup(dir);
  }
});

/**
 * T3 neg-budget-invalid: `--budget 0` → non-zero; stderr invalid budget;
 * no silent default to 4096.
 */
test('neg-budget-invalid: --budget 0 exits non-zero with invalid budget stderr', () => {
  const dir = spaceRoot();
  try {
    mkdirSync(path.join(dir, 'docs'), { recursive: true });
    writeFileSync(path.join(dir, 'docs', 'README.md'), OVERSIZE_README);

    const res = runContext(dir, ['--budget', '0']);
    assert.notEqual(
      res.code,
      0,
      `context --budget 0 must exit non-zero (no silent default 4096)\nstderr=${res.stderr}\nstdout=${res.stdout}`,
    );
    assert.match(
      res.stderr,
      /spacecraft context:.*invalid budget/i,
      `stderr must match /spacecraft context:.*invalid budget/i\nstderr=${res.stderr}`,
    );
  } finally {
    cleanup(dir);
  }
});

/**
 * T3 neg-budget-invalid (env arm): SPACECRAFT_CONTEXT_BUDGET=abc without flag →
 * non-zero; stderr invalid budget; no silent default to 4096.
 */
test('neg-budget-invalid: SPACECRAFT_CONTEXT_BUDGET=abc exits non-zero with invalid budget stderr', () => {
  const dir = spaceRoot();
  try {
    mkdirSync(path.join(dir, 'docs'), { recursive: true });
    writeFileSync(path.join(dir, 'docs', 'README.md'), OVERSIZE_README);

    const res = runContext(dir, [], { [ENV_BUDGET]: 'abc' });
    assert.notEqual(
      res.code,
      0,
      `context with ${ENV_BUDGET}=abc must exit non-zero (no silent default 4096)\nstderr=${res.stderr}\nstdout=${res.stdout}`,
    );
    assert.match(
      res.stderr,
      /spacecraft context:.*invalid budget/i,
      `stderr must match /spacecraft context:.*invalid budget/i\nstderr=${res.stderr}`,
    );
  } finally {
    cleanup(dir);
  }
});

/**
 * T4 pos-hook-context: source invokes `spacecraft context` / `./spacecraft context`
 * (not status alone); with binary present, hook prints context pack and exits 0.
 */
test('pos-hook-context: session-start invokes context and prints pack', () => {
  const src = readFileSync(sessionStartHook, 'utf8');
  assert.match(
    src,
    HOOK_CONTEXT_INJECT_RE,
    `session-start.sh must invoke spacecraft context or ./spacecraft context\n${src}`,
  );
  assert.equal(
    HOOK_STATUS_INJECT_RE.test(src),
    false,
    `session-start.sh must not use status alone as inject\n${src}`,
  );

  const dir = spaceRoot();
  try {
    mkdirSync(path.join(dir, 'docs'), { recursive: true });
    writeFileSync(path.join(dir, 'docs', 'README.md'), '# Product README\n');
    symlinkSync(entryPath, path.join(dir, 'spacecraft'));
    chmodSync(path.join(dir, 'spacecraft'), 0o755);

    const res = runSessionStartHook(dir);
    assert.equal(
      res.code,
      0,
      `session-start must exit 0\nstderr=${res.stderr}\nstdout=${res.stdout}`,
    );

    const out = res.stdout;
    assert.ok(
      out.includes(HEADING_README) || out.includes(HEADING_STATUS),
      `stdout must be context pack (include ${HEADING_README} or ${HEADING_STATUS}), not status-only inject\n${out}`,
    );
    assert.equal(
      out.trim() === HOOK_ABSENT_FALLBACK,
      false,
      `with binary present, stdout must not be absent-binary fallback alone\n${out}`,
    );
  } finally {
    cleanup(dir);
  }
});

/**
 * T4: hook always exits 0; fallback message only when spacecraft binary is absent.
 */
test('session-start: exits 0 with fallback only when spacecraft binary absent', () => {
  const dir = spaceRoot();
  try {
    // No ./spacecraft in cwd; PATH stripped of any spacecraft binary.
    const res = runSessionStartHook(dir, {
      PATH: '/usr/bin:/bin',
    });
    assert.equal(
      res.code,
      0,
      `session-start must exit 0 when binary absent\nstderr=${res.stderr}\nstdout=${res.stdout}`,
    );
    assert.match(
      res.stdout,
      new RegExp(`^${HOOK_ABSENT_FALLBACK.replace(/\./g, '\\.')}\\n?$`),
      `stdout must be absent-binary fallback only\nstdout=${res.stdout}`,
    );

    const src = readFileSync(sessionStartHook, 'utf8');
    // Fallback printf must not ride on `spacecraft … ||` (design: only when binary absent).
    assert.equal(
      /\|\|\s*printf[^\n]*No active spacecraft mission/.test(src),
      false,
      `fallback printf must not follow spacecraft invoke via ||\n${src}`,
    );
  } finally {
    cleanup(dir);
  }
});
