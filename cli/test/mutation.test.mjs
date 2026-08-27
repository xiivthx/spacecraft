/**
 * T3–T4 — spacecraft mutation diff scope, evidence, exit codes (M9G7II1F).
 *
 * Stub env vars for coder (no real Stryker in CI):
 * - SPACECRAFT_MUTATION_STUB=pass   → score 85 (above default threshold 80)
 * - SPACECRAFT_MUTATION_STUB=fail   → score 70 (below threshold)
 * - SPACECRAFT_MUTATION_STUB=missing → mutator not installed
 * - SPACECRAFT_MUTATION_SCOPE_CAP=<n> → override default scope cap (50)
 */
import assert from 'node:assert/strict';
import { spawnSync } from 'node:child_process';
import {
  existsSync,
  mkdirSync,
  mkdtempSync,
  readFileSync,
  rmSync,
  writeFileSync,
} from 'node:fs';
import os from 'node:os';
import path from 'node:path';
import test from 'node:test';
import { fileURLToPath } from 'node:url';

const repoRoot = path.resolve(path.dirname(fileURLToPath(import.meta.url)), '../..');
const entryPath = path.join(repoRoot, 'cli', 'spacecraft.mjs');

function spaceRoot() {
  const dir = mkdtempSync(path.join(os.tmpdir(), 'spacecraft-mutation-'));
  mkdirSync(path.join(dir, '.space', 'missions'), { recursive: true });
  mkdirSync(path.join(dir, '.space', 'roadmaps'), { recursive: true });
  return dir;
}

function cleanup(dir) {
  rmSync(dir, { recursive: true, force: true });
}

function isolatedGitEnv(dir) {
  const globalCfg = path.join(dir, '.gitconfig');
  const xdgConfig = path.join(dir, '.config');
  mkdirSync(xdgConfig, { recursive: true });
  if (!existsSync(globalCfg)) writeFileSync(globalCfg, '');
  return {
    ...process.env,
    HOME: dir,
    GIT_CONFIG_GLOBAL: globalCfg,
    GIT_CONFIG_NOSYSTEM: '1',
    GIT_TERMINAL_PROMPT: '0',
    XDG_CONFIG_HOME: xdgConfig,
  };
}

function runGit(dir, args) {
  const result = spawnSync('git', args, {
    cwd: dir,
    encoding: 'utf8',
    env: isolatedGitEnv(dir),
  });
  if (result.status !== 0) {
    throw new Error(`git ${args.join(' ')} failed: ${result.stderr || result.stdout}`);
  }
}

function initGitRepo(dir, branch = 'main') {
  runGit(dir, ['-c', 'init.defaultBranch=main', 'init', '--template=']);
  runGit(dir, ['config', 'user.email', 'test@example.com']);
  runGit(dir, ['config', 'user.name', 'Test']);
  writeFileSync(path.join(dir, '.gitkeep'), '');
  runGit(dir, ['add', '.gitkeep']);
  runGit(dir, ['commit', '-m', 'chore: init']);
  if (branch && branch !== 'main' && branch !== 'master') {
    runGit(dir, ['checkout', '-b', branch]);
  }
}

function writeMission(root, id, { branches = [] } = {}) {
  const dir = path.join(root, '.space', 'missions', id);
  mkdirSync(dir, { recursive: true });
  writeFileSync(
    path.join(dir, 'mission.json'),
    `${JSON.stringify(
      {
        id,
        title: 'Mutation Test Mission',
        state: 'active',
        createdAt: '2026-01-01T00:00:00Z',
        branches,
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
  return dir;
}

function writeCurrent(root, id) {
  writeFileSync(path.join(root, '.space', 'current'), `${id}\n`);
}

function readEvidenceLines(root, id) {
  const data = readFileSync(
    path.join(root, '.space', 'missions', id, 'evidence.jsonl'),
    'utf8',
  );
  return data
    .split('\n')
    .map((l) => l.trim())
    .filter(Boolean)
    .map((l) => JSON.parse(l));
}

function mutationLines(lines) {
  return lines.filter((e) => typeof e.label === 'string' && /^mutation-/.test(e.label));
}

function parseMutationSummary(entry) {
  assert.match(entry.label, /^mutation-/);
  const summary = JSON.parse(entry.output);
  assert.equal(typeof summary.files, 'object');
  assert.ok(Array.isArray(summary.files));
  assert.equal(typeof summary.score, 'number');
  assert.equal(typeof summary.threshold, 'number');
  assert.equal(typeof summary.pass, 'boolean');
  return summary;
}

function runCLI(dir, args, envExtra = {}) {
  const result = spawnSync(process.execPath, [entryPath, ...args], {
    cwd: dir,
    encoding: 'utf8',
    env: { ...process.env, ...envExtra },
  });
  return {
    stdout: result.stdout ?? '',
    stderr: result.stderr ?? '',
    code: result.status ?? 1,
  };
}

function combined(res) {
  return `${res.stdout}${res.stderr}`;
}

function assertNotStub(res, label) {
  assert.doesNotMatch(
    combined(res),
    /not implemented/i,
    `${label} must be real dispatch, not stub\n${combined(res)}`,
  );
  assert.doesNotMatch(
    combined(res),
    /unknown command/i,
    `${label} must dispatch, not unknown\n${combined(res)}`,
  );
}

/**
 * Git fixture: main + feature branch with bound branch; optional mutable file changes.
 */
function writeMutationGitFixture(root, id, { changedFiles = 0, branchSuffix = 'mutation' } = {}) {
  const branch = `feat/${id}/${branchSuffix}`;
  initGitRepo(root, branch);
  writeMission(root, id, { branches: [branch] });
  writeCurrent(root, id);

  mkdirSync(path.join(root, 'fixtures', 'mutation'), { recursive: true });
  for (let i = 0; i < changedFiles; i++) {
    const rel = `fixtures/mutation/file-${i}.mjs`;
    writeFileSync(path.join(root, rel), `export const n${i} = ${i};\n`);
    runGit(root, ['add', rel]);
  }
  if (changedFiles > 0) {
    runGit(root, ['commit', '-m', `test: add ${changedFiles} mutable files`]);
  }
  return { branch };
}

function runMutation(root, id, envExtra = {}) {
  return runCLI(root, ['mutation', '--mission', id], envExtra);
}

// --- T3: diff scope ---

test('edge-empty-scope: zero mutable files yields scope-empty marker and exit 0', () => {
  const dir = spaceRoot();
  const id = 'M9G7MUT01';
  try {
    writeMutationGitFixture(dir, id, { changedFiles: 0 });

    const res = runMutation(dir, id, { SPACECRAFT_MUTATION_STUB: 'pass' });
    assertNotStub(res, 'mutation');
    assert.equal(res.code, 0, `empty scope must exit 0\n${combined(res)}`);
    assert.match(combined(res), /scope-empty/);

    const lines = readEvidenceLines(dir, id);
    const mut = mutationLines(lines);
    if (mut.length > 0) {
      const summary = parseMutationSummary(mut[0]);
      assert.ok(
        !('score' in summary) || summary.score === undefined,
        'empty scope must not fabricate numeric score',
      );
    }
  } finally {
    cleanup(dir);
  }
});

test('edge-oversize-scope: scope exceeding cap truncates files and records marker', () => {
  const dir = spaceRoot();
  const id = 'M9G7MUT02';
  try {
    writeMutationGitFixture(dir, id, { changedFiles: 5 });

    const res = runMutation(dir, id, {
      SPACECRAFT_MUTATION_STUB: 'pass',
      SPACECRAFT_MUTATION_SCOPE_CAP: '3',
    });
    assertNotStub(res, 'mutation');
    assert.equal(res.code, 0, `oversize scope run\n${combined(res)}`);

    const mut = mutationLines(readEvidenceLines(dir, id));
    assert.equal(mut.length, 1, 'expected one mutation evidence line');
    const summary = parseMutationSummary(mut[0]);
    assert.ok(summary.files.length <= 3, `files truncated to cap, got ${summary.files.length}`);
    assert.match(combined(res), /truncat/i);
  } finally {
    cleanup(dir);
  }
});

test('over-no-merge-base: mission without bound branch exits non-zero naming merge-base', () => {
  const dir = spaceRoot();
  const id = 'M9G7MUT03';
  try {
    initGitRepo(dir, 'main');
    writeMission(dir, id, { branches: [] });
    writeCurrent(dir, id);
    mkdirSync(path.join(dir, 'fixtures'), { recursive: true });
    writeFileSync(path.join(dir, 'fixtures', 'orphan.mjs'), 'export const x = 1;\n');
    runGit(dir, ['add', 'fixtures/orphan.mjs']);
    runGit(dir, ['commit', '-m', 'test: orphan change']);

    const res = runMutation(dir, id, { SPACECRAFT_MUTATION_STUB: 'pass' });
    assertNotStub(res, 'mutation');
    assert.notEqual(res.code, 0, `missing merge-base must fail\n${combined(res)}`);
    assert.match(combined(res), /merge-base/i);
    assert.equal(mutationLines(readEvidenceLines(dir, id)).length, 0);
  } finally {
    cleanup(dir);
  }
});

// --- T4: evidence and exit codes ---

test('pos-mutation-pass: at/above threshold writes mutation-<scope> evidence and exits 0', () => {
  const dir = spaceRoot();
  const id = 'M9G7MUT04';
  try {
    writeMutationGitFixture(dir, id, { changedFiles: 2 });

    const before = readEvidenceLines(dir, id).length;
    const res = runMutation(dir, id, { SPACECRAFT_MUTATION_STUB: 'pass' });
    assertNotStub(res, 'mutation');
    assert.equal(res.code, 0, `pass run\n${combined(res)}`);

    const after = readEvidenceLines(dir, id);
    assert.equal(after.length, before + 1, 'exactly one evidence line appended');
    const mut = mutationLines(after);
    assert.equal(mut.length, 1);
    assert.match(mut[0].label, /^mutation-/);

    const summary = parseMutationSummary(mut[0]);
    assert.equal(summary.threshold, 80);
    assert.ok(summary.score >= summary.threshold);
    assert.equal(summary.pass, true);
    assert.ok(summary.files.length > 0);
  } finally {
    cleanup(dir);
  }
});

test('neg-mutation-fail: below-threshold stub writes pass:false and exits non-zero', () => {
  const dir = spaceRoot();
  const id = 'M9G7MUT05';
  try {
    writeMutationGitFixture(dir, id, { changedFiles: 1 });

    const res = runMutation(dir, id, { SPACECRAFT_MUTATION_STUB: 'fail' });
    assertNotStub(res, 'mutation');
    assert.notEqual(res.code, 0, `below threshold must fail\n${combined(res)}`);

    const mut = mutationLines(readEvidenceLines(dir, id));
    assert.equal(mut.length, 1);
    const summary = parseMutationSummary(mut[0]);
    assert.equal(summary.pass, false);
    assert.ok(summary.score < summary.threshold);
  } finally {
    cleanup(dir);
  }
});

test('neg-mutation-fail: missing mutator exits non-zero with install hint and no evidence', () => {
  const dir = spaceRoot();
  const id = 'M9G7MUT06';
  try {
    writeMutationGitFixture(dir, id, { changedFiles: 1 });

    const before = readEvidenceLines(dir, id).length;
    const res = runMutation(dir, id, { SPACECRAFT_MUTATION_STUB: 'missing' });
    assertNotStub(res, 'mutation');
    assert.notEqual(res.code, 0, `missing tool must fail\n${combined(res)}`);
    assert.match(combined(res), /install|stryker|mutator|npm/i);

    const after = readEvidenceLines(dir, id);
    assert.equal(after.length, before, 'tool-missing must not append mutation- evidence');
    assert.equal(mutationLines(after).length, 0);
  } finally {
    cleanup(dir);
  }
});

test('over-append-only: second mutation run appends; first line bytes unchanged', () => {
  const dir = spaceRoot();
  const id = 'M9G7MUT07';
  try {
    writeMutationGitFixture(dir, id, { changedFiles: 1 });

    const res1 = runMutation(dir, id, { SPACECRAFT_MUTATION_STUB: 'pass' });
    assert.equal(res1.code, 0, `first run\n${combined(res1)}`);

    const evidencePath = path.join(dir, '.space', 'missions', id, 'evidence.jsonl');
    const firstSnapshot = readFileSync(evidencePath, 'utf8');
    const firstLine = firstSnapshot.split('\n').find((l) => l.trim()) ?? '';
    assert.ok(firstLine.length > 0);

    const res2 = runMutation(dir, id, { SPACECRAFT_MUTATION_STUB: 'pass' });
    assert.equal(res2.code, 0, `second run\n${combined(res2)}`);

    const allText = readFileSync(evidencePath, 'utf8');
    const lines = allText.split('\n').filter((l) => l.trim());
    assert.equal(lines.length, 2, 'two append-only mutation evidence lines');
    assert.equal(lines[0], firstLine, 'first line bytes must be unchanged');
    assert.notEqual(lines[1], lines[0]);
    assert.equal(mutationLines(readEvidenceLines(dir, id)).length, 2);
  } finally {
    cleanup(dir);
  }
});
