/**
 * Node CLI tests for evi/evidence JSONL capture: exitCode propagation,
 * outputHash integrity, and oversized-output truncate/sidecar behavior.
 */
import assert from 'node:assert/strict';
import { createHash } from 'node:crypto';
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

/** Independent oracle: SHA-256 of empty merged stdout+stderr. */
const EMPTY_OUTPUT_HASH =
  'e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855';

/** Independent oracle: SHA-256 of `echo hello` merged stdout+stderr (`hello\n`). */
const HELLO_OUTPUT = 'hello\n';
const HELLO_OUTPUT_HASH = createHash('sha256').update(HELLO_OUTPUT).digest('hex');

/** Discuss lock: truncate when output length exceeds this many bytes. */
const EVI_TRUNCATE_THRESHOLD = 65536;
/** Discuss lock: trailing marker on truncated JSONL `output`. */
const EVI_TRUNCATE_MARKER = '\n...[truncated]';

function sha256Hex(s) {
  return createHash('sha256').update(s, 'utf8').digest('hex');
}

/** Resolve sidecar path from evidence entry (`outputRawPath` relative to mission dir or absolute). */
function resolveOutputRawPath(root, missionId, entry) {
  assert.equal(
    typeof entry.outputRawPath,
    'string',
    'truncated evidence must record outputRawPath',
  );
  assert.ok(entry.outputRawPath.length > 0, 'outputRawPath must be non-empty');
  const missionRoot = path.join(root, '.space', 'missions', missionId);
  return path.isAbsolute(entry.outputRawPath)
    ? entry.outputRawPath
    : path.join(missionRoot, entry.outputRawPath);
}

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

// --- T2: truncate threshold / sidecar / path-unsafe label ---

test('evi keeps full output at and under truncate threshold without sidecar truncate flags', () => {
  const dir = spaceRoot();
  const id = 'M8S5EE9S01';
  try {
    writeMission(dir, id);

    const emptyRes = runCLI(dir, 'evi', '--mission', id, 'under-empty', '--', 'true');
    assertNotStub(emptyRes, 'evi under-empty');
    assert.equal(emptyRes.code, 0, `evi exit=${emptyRes.code}\n${combined(emptyRes)}`);
    const emptyEntry = readLastEvidence(dir, id);
    assert.equal(emptyEntry.output, '');
    assert.equal(emptyEntry.outputHash, EMPTY_OUTPUT_HASH);
    assert.notEqual(emptyEntry.outputTruncated, true);
    assert.equal(
      emptyRes.stdout,
      'Exit code: 0\n',
      'terminal must print full (empty) output then exit line',
    );

    const exact = 'y'.repeat(EVI_TRUNCATE_THRESHOLD);
    assert.equal(Buffer.byteLength(exact, 'utf8'), EVI_TRUNCATE_THRESHOLD);
    const exactHash = sha256Hex(exact);
    const exactRes = runCLI(
      dir,
      'evi',
      '--mission',
      id,
      'under-exact',
      '--',
      process.execPath,
      '-e',
      `process.stdout.write('y'.repeat(${EVI_TRUNCATE_THRESHOLD}))`,
    );
    assertNotStub(exactRes, 'evi under-exact');
    assert.equal(exactRes.code, 0, `evi exit=${exactRes.code}\n${combined(exactRes)}`);

    const exactEntry = readLastEvidence(dir, id);
    assert.equal(exactEntry.output, exact, 'exactly 65536 bytes must stay fully in JSONL');
    assert.equal(exactEntry.outputHash, exactHash);
    assert.notEqual(
      exactEntry.outputTruncated,
      true,
      'at threshold must not set outputTruncated:true',
    );
    assert.ok(
      exactRes.stdout.startsWith(exact),
      'terminal must still print the full output at threshold',
    );
  } finally {
    cleanup(dir);
  }
});

test('evi truncates oversized JSONL output to sidecar while hashing and printing full raw', () => {
  const dir = spaceRoot();
  const id = 'M8S5EE9S02';
  try {
    writeMission(dir, id);

    const overLen = EVI_TRUNCATE_THRESHOLD + 1;
    const raw = 'z'.repeat(overLen);
    assert.equal(Buffer.byteLength(raw, 'utf8'), overLen);
    const fullHash = sha256Hex(raw);

    const res = runCLI(
      dir,
      'evi',
      '--mission',
      id,
      'over-limit',
      '--',
      process.execPath,
      '-e',
      `process.stdout.write('z'.repeat(${overLen}))`,
    );
    assertNotStub(res, 'evi over-limit');
    assert.equal(res.code, 0, `evi exit=${res.code}\n${combined(res)}`);

    const entry = readLastEvidence(dir, id);
    assert.equal(entry.outputTruncated, true);
    assert.equal(entry.outputBytes, overLen);
    assert.ok(
      entry.output.endsWith(EVI_TRUNCATE_MARKER),
      `JSONL output must end with truncate marker; got suffix ${JSON.stringify(entry.output.slice(-32))}`,
    );
    assert.notEqual(entry.output, raw, 'JSONL output must not keep the full raw');
    assert.equal(
      entry.outputHash,
      fullHash,
      `outputHash must be SHA-256 of full raw, want ${fullHash}, got ${entry.outputHash}`,
    );
    assert.notEqual(
      entry.outputHash,
      sha256Hex(entry.output),
      'outputHash must not be hash of truncated JSONL output alone',
    );

    const sidecarPath = resolveOutputRawPath(dir, id, entry);
    const evidenceRawDir = path.join(dir, '.space', 'missions', id, 'evidence-raw');
    assert.ok(
      sidecarPath.startsWith(evidenceRawDir + path.sep) ||
        path.dirname(sidecarPath) === evidenceRawDir,
      `outputRawPath must resolve under evidence-raw/; got ${sidecarPath}`,
    );
    assert.ok(existsSync(sidecarPath), `sidecar missing at ${sidecarPath}`);
    assert.equal(
      readFileSync(sidecarPath, 'utf8'),
      raw,
      'sidecar must contain the full raw output',
    );

    assert.ok(
      res.stdout.startsWith(raw),
      'terminal must still print the full raw output when truncated in JSONL',
    );
  } finally {
    cleanup(dir);
  }
});

test('evi path-unsafe label still writes sanitized sidecar under evidence-raw', () => {
  const dir = spaceRoot();
  const id = 'M8S5EE9S03';
  const unsafeLabel = '../evil:name/with spaces';
  try {
    writeMission(dir, id);

    const overLen = EVI_TRUNCATE_THRESHOLD + 1;
    const raw = 'w'.repeat(overLen);
    const res = runCLI(
      dir,
      'evi',
      '--mission',
      id,
      unsafeLabel,
      '--',
      process.execPath,
      '-e',
      `process.stdout.write('w'.repeat(${overLen}))`,
    );
    assertNotStub(res, 'evi path-unsafe label');
    assert.equal(res.code, 0, `evi exit=${res.code}\n${combined(res)}`);

    const entry = readLastEvidence(dir, id);
    assert.equal(entry.label, unsafeLabel);
    assert.equal(entry.outputTruncated, true);

    const sidecarPath = resolveOutputRawPath(dir, id, entry);
    const evidenceRawDir = path.join(dir, '.space', 'missions', id, 'evidence-raw');
    const rel = path.relative(evidenceRawDir, sidecarPath);
    assert.ok(
      rel && !rel.startsWith('..') && !path.isAbsolute(rel),
      `sidecar must stay under evidence-raw/; outputRawPath=${entry.outputRawPath} resolved=${sidecarPath}`,
    );
    assert.doesNotMatch(
      path.basename(sidecarPath),
      /[\\/]/,
      'sanitized sidecar basename must not contain path separators',
    );
    assert.ok(
      !String(entry.outputRawPath).includes('..'),
      'recorded outputRawPath must not retain path traversal segments',
    );
    assert.ok(existsSync(sidecarPath), `writable sidecar missing at ${sidecarPath}`);
    assert.equal(readFileSync(sidecarPath, 'utf8'), raw);
  } finally {
    cleanup(dir);
  }
});
