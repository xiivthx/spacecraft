/**
 * Node CLI tests for evi/evidence JSONL capture: exitCode propagation and
 * outputHash integrity.
 */
import assert from 'node:assert/strict';
import { createHash } from 'node:crypto';
import { spawnSync } from 'node:child_process';
import {
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

/** Independent oracle: SHA-256 of empty merged stdout+stderr. */
const EMPTY_OUTPUT_HASH =
  'e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855';

/** Independent oracle: SHA-256 of `echo hello` merged stdout+stderr (`hello\n`). */
const HELLO_OUTPUT = 'hello\n';
const HELLO_OUTPUT_HASH = createHash('sha256').update(HELLO_OUTPUT).digest('hex');

function spaceRoot() {
  const dir = mkdtempSync(path.join(os.tmpdir(), 'spacecraft-evi-'));
  mkdirSync(path.join(dir, '.space', 'missions'), { recursive: true });
  mkdirSync(path.join(dir, '.space', 'roadmaps'), { recursive: true });
  return dir;
}

function writeMission(root, id, state = 'active') {
  const dir = path.join(root, '.space', 'missions', id);
  mkdirSync(dir, { recursive: true });
  writeFileSync(
    path.join(dir, 'mission.json'),
    `${JSON.stringify(
      {
        id,
        title: 'Test Mission',
        state,
        createdAt: '2026-01-01T00:00:00Z',
        branches: [],
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
}

function runCLI(dir, ...args) {
  const result = spawnSync(process.execPath, [entryPath, ...args], {
    cwd: dir,
    encoding: 'utf8',
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

function readLastEvidence(root, id) {
  const data = readFileSync(
    path.join(root, '.space', 'missions', id, 'evidence.jsonl'),
    'utf8',
  );
  const lines = data
    .split('\n')
    .map((l) => l.trim())
    .filter(Boolean);
  assert.ok(lines.length > 0, 'evidence.jsonl empty');
  return JSON.parse(lines[lines.length - 1]);
}

function cleanup(dir) {
  rmSync(dir, { recursive: true, force: true });
}

// --- Acceptance 1: JSONL fields + exitCode propagation ---

test('evi records exitCode in JSONL and propagates command exit', () => {
  const dir = spaceRoot();
  const id = 'M07EVI02';
  try {
    writeMission(dir, id);

    const res = runCLI(
      dir,
      'evi',
      '--mission',
      id,
      'exit-status',
      '--',
      'sh',
      '-c',
      'echo out; exit 7',
    );
    assertNotStub(res, 'evi exit-status');

    const entry = readLastEvidence(dir, id);
    assert.equal(entry.label, 'exit-status');
    assert.equal(typeof entry.command, 'string');
    assert.match(entry.command, /sh/);
    assert.equal(typeof entry.output, 'string');
    assert.match(entry.output, /out/);
    assert.equal(typeof entry.ts, 'string');
    assert.ok(entry.ts.length > 0, 'ts must be non-empty');
    assert.equal(entry.exitCode, 7);

    assert.equal(res.code, 7, `evi process exit=${res.code}, want 7\n${combined(res)}`);
  } finally {
    cleanup(dir);
  }
});

test('evi failing command returns nonzero', () => {
  const dir = spaceRoot();
  const id = 'M07EVI01';
  try {
    writeMission(dir, id);

    const res = runCLI(
      dir,
      'evi',
      '--mission',
      id,
      'fail-case',
      '--',
      'sh',
      '-c',
      'echo boom; exit 42',
    );
    assertNotStub(res, 'evi fail-case');
    assert.notEqual(res.code, 0, `evi with failing command returned 0\n${combined(res)}`);
    assert.equal(res.code, 42, `prefer exact captured exit 42; got ${res.code}`);
  } finally {
    cleanup(dir);
  }
});

test('evi success exits 0 and still records exitCode', () => {
  const dir = spaceRoot();
  const id = 'M07EVI03';
  try {
    writeMission(dir, id);

    const res = runCLI(dir, 'evi', '--mission', id, 'ok', '--', 'echo', 'hello');
    assertNotStub(res, 'evi ok');
    assert.equal(res.code, 0, `evi success exit=${res.code}\n${combined(res)}`);

    const entry = readLastEvidence(dir, id);
    assert.equal(entry.exitCode, 0);
    assert.equal(entry.label, 'ok');
    assert.equal(typeof entry.command, 'string');
    assert.equal(typeof entry.output, 'string');
    assert.equal(typeof entry.ts, 'string');
  } finally {
    cleanup(dir);
  }
});

// --- Acceptance 2: outputHash + aliases ---

test('evi and evidence aliases both write evidence', () => {
  const dir = spaceRoot();
  const id = 'M07EVI04';
  try {
    writeMission(dir, id);

    const evi = runCLI(dir, 'evi', '--mission', id, 'via-evi', '--', 'echo', 'a');
    assertNotStub(evi, 'evi alias');
    assert.equal(evi.code, 0, `evi exit=${evi.code}\n${combined(evi)}`);

    const ev = runCLI(dir, 'evidence', '--mission', id, 'via-evidence', '--', 'echo', 'b');
    assertNotStub(ev, 'evidence canonical');
    assert.equal(ev.code, 0, `evidence exit=${ev.code}\n${combined(ev)}`);

    const data = readFileSync(
      path.join(dir, '.space', 'missions', id, 'evidence.jsonl'),
      'utf8',
    );
    const lines = data
      .split('\n')
      .map((l) => l.trim())
      .filter(Boolean);
    assert.equal(lines.length, 2, `want 2 evidence lines, got ${lines.length}`);
    const labels = lines.map((l) => JSON.parse(l).label);
    assert.deepEqual(labels, ['via-evi', 'via-evidence']);
  } finally {
    cleanup(dir);
  }
});

test('evi records matching outputHash for stdout', () => {
  const dir = spaceRoot();
  const id = 'M07EVI05';
  try {
    writeMission(dir, id);

    const res = runCLI(dir, 'evi', '--mission', id, 'hash-case', '--', 'echo', 'hello');
    assertNotStub(res, 'evi hash-case');
    assert.equal(res.code, 0, `evi exit=${res.code}\n${combined(res)}`);

    const entry = readLastEvidence(dir, id);
    assert.equal(entry.output, HELLO_OUTPUT);
    assert.equal(
      entry.outputHash,
      HELLO_OUTPUT_HASH,
      `outputHash=${entry.outputHash}, want ${HELLO_OUTPUT_HASH}`,
    );
  } finally {
    cleanup(dir);
  }
});

test('evi records outputHash for empty stdout', () => {
  const dir = spaceRoot();
  const id = 'M07EVI06';
  try {
    writeMission(dir, id);

    const res = runCLI(dir, 'evi', '--mission', id, 'empty-hash', '--', 'true');
    assertNotStub(res, 'evi empty-hash');
    assert.equal(res.code, 0, `evi exit=${res.code}\n${combined(res)}`);

    const entry = readLastEvidence(dir, id);
    assert.equal(entry.output, '', 'empty merged output must record empty string');
    assert.equal(
      entry.outputHash,
      EMPTY_OUTPUT_HASH,
      `outputHash=${entry.outputHash}, want ${EMPTY_OUTPUT_HASH}`,
    );
  } finally {
    cleanup(dir);
  }
});
