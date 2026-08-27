/**
 * T2–T5 — spacecraft freeze / freeze-check CLI (M9G7IHV3 freeze-tooling).
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
  unlinkSync,
  writeFileSync,
} from 'node:fs';
import os from 'node:os';
import path from 'node:path';
import test from 'node:test';
import { fileURLToPath } from 'node:url';

const repoRoot = path.resolve(path.dirname(fileURLToPath(import.meta.url)), '../..');
const entryPath = path.join(repoRoot, 'cli', 'spacecraft.mjs');
const EVI_TRUNCATE_THRESHOLD = 65536;

function sha256Hex(content) {
  return createHash('sha256').update(content, 'utf8').digest('hex');
}

function spaceRoot() {
  const dir = mkdtempSync(path.join(os.tmpdir(), 'spacecraft-freeze-'));
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
        title: 'Freeze Test Mission',
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
  return dir;
}

function writeCurrent(root, id) {
  writeFileSync(path.join(root, '.space', 'current'), `${id}\n`);
}

function writeDecisions(root, id, body) {
  writeFileSync(path.join(root, '.space', 'missions', id, 'decisions.md'), body);
}

function writeApprovedScenarios(root, id, body) {
  writeFileSync(
    path.join(root, '.space', 'missions', id, 'approved-scenarios.md'),
    body,
  );
}

function writeEvidenceLine(root, id, line) {
  writeFileSync(path.join(root, '.space', 'missions', id, 'evidence.jsonl'), line);
}

function appendEvidenceLine(root, id, line) {
  const p = path.join(root, '.space', 'missions', id, 'evidence.jsonl');
  const prev = existsSync(p) ? readFileSync(p, 'utf8') : '';
  writeFileSync(p, `${prev}${line}`);
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

function runCLI(dir, args) {
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

function cleanup(dir) {
  rmSync(dir, { recursive: true, force: true });
}

/** Fixture: two test files + frozen approved-scenarios for freeze-cmd happy path. */
function writeFreezeFixture(root, id) {
  const missionDir = writeMission(root, id);
  writeCurrent(root, id);
  writeDecisions(root, id, '# Decisions\n\nGates version: M9G7IHV3\n');
  writeApprovedScenarios(
    root,
    id,
    '# Approved scenarios\n\n| id | status |\n| --- | --- |\n| P1 | frozen |\n\n## Freeze footer\n\n```\nApproved-scenarios: frozen-from-contract\n```\n',
  );

  mkdirSync(path.join(root, 'fixtures'), { recursive: true });
  const testA = path.join(root, 'fixtures', 'alpha.test.mjs');
  const testB = path.join(root, 'fixtures', 'beta.test.mjs');
  writeFileSync(testA, 'export const alpha = 1;\n');
  writeFileSync(testB, 'export const beta = 2;\n');

  return {
    missionDir,
    paths: [
      'fixtures/alpha.test.mjs',
      'fixtures/beta.test.mjs',
      `.space/missions/${id}/approved-scenarios.md`,
    ],
    files: { testA, testB },
  };
}

function parseFreezeManifest(entry) {
  assert.equal(entry.label, 'freeze');
  let manifest;
  if (entry.outputTruncated === true) {
    assert.equal(typeof entry.outputRawPath, 'string');
    const sidecar = path.join(
      path.dirname(path.join(repoRoot, 'noop')),
      entry.outputRawPath,
    );
    // sidecar path is relative to mission dir — resolved in tests that know root+id
    return { entry, manifest: null, sidecarRel: entry.outputRawPath };
  }
  manifest = JSON.parse(entry.output);
  assert.equal(manifest.kind, 'freeze-manifest');
  assert.ok(Array.isArray(manifest.files));
  return { entry, manifest };
}

// --- T2: freeze-cmd ---

test('freeze-cmd pos-freeze-happy: one freeze event with sha256 manifest (P1)', () => {
  const dir = spaceRoot();
  const id = 'M9G7FRZ01';
  try {
    const { paths } = writeFreezeFixture(dir, id);
    const before = readEvidenceLines(dir, id);
    assert.equal(before.length, 0);

    const res = runCLI(dir, ['freeze', '--mission', id, ...paths]);
    assertNotStub(res, 'freeze');
    assert.equal(res.code, 0, `freeze exit=${res.code}\n${combined(res)}`);

    const after = readEvidenceLines(dir, id);
    assert.equal(after.length, 1, 'exactly one evidence line appended');
    const { manifest } = parseFreezeManifest(after[0]);
    assert.equal(manifest.files.length, 3);

    const scenariosPath = `.space/missions/${id}/approved-scenarios.md`;
    const expected = new Map(
      paths.map((p) => [p, sha256Hex(readFileSync(path.join(dir, p), 'utf8'))]),
    );
    for (const item of manifest.files) {
      assert.match(item.path, /^[^/\\]/, 'paths must be repo-relative without leading slash');
      assert.doesNotMatch(item.path, /\.\./, 'paths must not contain ..');
      assert.match(item.path, /\//, 'paths must use forward slashes when nested');
      assert.equal(typeof item.sha256, 'string');
      assert.match(item.sha256, /^[0-9a-f]{64}$/);
      assert.equal(item.sha256, expected.get(item.path), `hash mismatch for ${item.path}`);
    }
    assert.ok(manifest.files.some((f) => f.path === scenariosPath));
  } finally {
    cleanup(dir);
  }
});

test('freeze-cmd neg-freeze-missing-file: nonexistent path exits non-zero, no event (N4)', () => {
  const dir = spaceRoot();
  const id = 'M9G7FRZ02';
  try {
    writeFreezeFixture(dir, id);

    const res = runCLI(dir, [
      'freeze',
      '--mission',
      id,
      'fixtures/does-not-exist.test.mjs',
    ]);
    assertNotStub(res, 'freeze');
    assert.notEqual(res.code, 0, `expected non-zero for missing file\n${combined(res)}`);

    const after = readEvidenceLines(dir, id);
    assert.equal(after.length, 0, 'must not append freeze event on missing file');
  } finally {
    cleanup(dir);
  }
});

test('freeze-cmd edge-freeze-large-manifest: >64KB uses sidecar and hash recheck passes (E1)', () => {
  const dir = spaceRoot();
  const id = 'M9G7FRZ03';
  try {
    const { paths, files } = writeFreezeFixture(dir, id);
    const bigContent = 'x'.repeat(EVI_TRUNCATE_THRESHOLD);
    writeFileSync(files.testA, bigContent);

    const res = runCLI(dir, ['freeze', '--mission', id, ...paths]);
    assertNotStub(res, 'freeze');
    assert.equal(res.code, 0, `freeze exit=${res.code}\n${combined(res)}`);

    const entry = readEvidenceLines(dir, id)[0];
    assert.equal(entry.label, 'freeze');
    assert.equal(entry.outputTruncated, true);
    assert.equal(typeof entry.outputRawPath, 'string');

    const sidecarPath = path.join(
      dir,
      '.space',
      'missions',
      id,
      entry.outputRawPath,
    );
    assert.ok(existsSync(sidecarPath), `sidecar missing at ${sidecarPath}`);
    const raw = readFileSync(sidecarPath, 'utf8');
    const manifest = JSON.parse(raw);
    assert.equal(manifest.kind, 'freeze-manifest');
    assert.ok(manifest.files.length >= 3);

    const check = runCLI(dir, ['freeze-check', '--mission', id]);
    assertNotStub(check, 'freeze-check');
    assert.equal(check.code, 0, `hash recheck must pass\n${combined(check)}`);
  } finally {
    cleanup(dir);
  }
});

// --- T3: freeze-check ---

test('freeze-check pos-unchanged-pass: freeze then unchanged files exit 0 (P2)', () => {
  const dir = spaceRoot();
  const id = 'M9G7FRZ04';
  try {
    const { paths } = writeFreezeFixture(dir, id);

    const freeze = runCLI(dir, ['freeze', '--mission', id, ...paths]);
    assertNotStub(freeze, 'freeze');
    assert.equal(freeze.code, 0, `freeze exit=${freeze.code}\n${combined(freeze)}`);

    const check = runCLI(dir, ['freeze-check', '--mission', id]);
    assertNotStub(check, 'freeze-check');
    assert.equal(check.code, 0, `unchanged freeze-check exit=${check.code}\n${combined(check)}`);
  } finally {
    cleanup(dir);
  }
});

test('freeze-check neg-drift-no-line: edited file without oracle line exits 1 with freeze-drift (N1)', () => {
  const dir = spaceRoot();
  const id = 'M9G7FRZ05';
  try {
    const { paths, files } = writeFreezeFixture(dir, id);

    const freeze = runCLI(dir, ['freeze', '--mission', id, ...paths]);
    assert.equal(freeze.code, 0, `freeze exit=${freeze.code}\n${combined(freeze)}`);

    writeFileSync(files.testA, 'export const alpha = 99;\n');

    const check = runCLI(dir, ['freeze-check', '--mission', id]);
    assertNotStub(check, 'freeze-check');
    assert.notEqual(check.code, 0, `expected drift fail\n${combined(check)}`);
    assert.match(combined(check), /freeze-drift/);
  } finally {
    cleanup(dir);
  }
});

test('freeze-check neg-retro-freeze: test evidence after freeze rejected as postdated-freeze (N2)', () => {
  const dir = spaceRoot();
  const id = 'M9G7FRZ06';
  try {
    const { paths } = writeFreezeFixture(dir, id);

    appendEvidenceLine(
      dir,
      id,
      '{"label":"test-run","command":"echo test","output":"ok\\n","ts":"2026-01-01T00:00:00Z","exitCode":0}\n',
    );

    const freeze = runCLI(dir, ['freeze', '--mission', id, ...paths]);
    assertNotStub(freeze, 'freeze');
    assert.notEqual(freeze.code, 0, `postdated freeze must fail\n${combined(freeze)}`);
    assert.match(combined(freeze), /postdated-freeze/);
  } finally {
    cleanup(dir);
  }
});

// --- T5: oracle cycle, docs skip, overlooked edges ---

test('pos-oracle-change-cycle: drift + Scenario oracle change + re-freeze passes (P3)', () => {
  const dir = spaceRoot();
  const id = 'M9G7FRZ07';
  try {
    const { paths, files } = writeFreezeFixture(dir, id);

    const freeze1 = runCLI(dir, ['freeze', '--mission', id, ...paths]);
    assert.equal(freeze1.code, 0, `initial freeze\n${combined(freeze1)}`);

    writeFileSync(files.testA, 'export const alpha = 42;\n');
    const drift = runCLI(dir, ['freeze-check', '--mission', id]);
    assert.notEqual(drift.code, 0);
    assert.match(combined(drift), /freeze-drift/);

    writeDecisions(
      dir,
      id,
      '# Decisions\n\nGates version: M9G7IHV3\n\nScenario oracle change: P3 - updated expected literal\n',
    );

    const freeze2 = runCLI(dir, ['freeze', '--mission', id, ...paths]);
    assert.equal(freeze2.code, 0, `re-freeze\n${combined(freeze2)}`);

    const lines = readEvidenceLines(dir, id);
    const freezeEvents = lines.filter((e) => e.label === 'freeze');
    assert.equal(freezeEvents.length, 2, 'prior freeze events retained');

    const check = runCLI(dir, ['freeze-check', '--mission', id]);
    assert.equal(check.code, 0, `oracle cycle pass\n${combined(check)}`);
  } finally {
    cleanup(dir);
  }
});

test('edge-docs-skip: Approved-scenarios skipped exempts freeze-check (E3)', () => {
  const dir = spaceRoot();
  const id = 'M9G7FRZ08';
  try {
    writeMission(dir, id);
    writeCurrent(dir, id);
    writeDecisions(dir, id, '# Decisions\n\nGates version: M9G7IHV3\n');
    writeApprovedScenarios(
      dir,
      id,
      '# Approved scenarios\n\nApproved-scenarios skipped: docs/prose-only\n',
    );

    const check = runCLI(dir, ['freeze-check', '--mission', id]);
    assertNotStub(check, 'freeze-check');
    assert.equal(check.code, 0, `docs skip must pass without freeze\n${combined(check)}`);
    assert.equal(readEvidenceLines(dir, id).length, 0, 'no freeze event required for docs skip');
  } finally {
    cleanup(dir);
  }
});

test('overlooked-freeze-edges O1 same-content-re-freeze: second freeze event allowed', () => {
  const dir = spaceRoot();
  const id = 'M9G7FRZ09';
  try {
    const { paths } = writeFreezeFixture(dir, id);

    const freeze1 = runCLI(dir, ['freeze', '--mission', id, ...paths]);
    assert.equal(freeze1.code, 0);

    const freeze2 = runCLI(dir, ['freeze', '--mission', id, ...paths]);
    assert.equal(freeze2.code, 0, `same-content re-freeze allowed\n${combined(freeze2)}`);

    const lines = readEvidenceLines(dir, id);
    assert.equal(lines.filter((e) => e.label === 'freeze').length, 2);

    const check = runCLI(dir, ['freeze-check', '--mission', id]);
    assert.equal(check.code, 0);
  } finally {
    cleanup(dir);
  }
});

test('overlooked-freeze-edges O2 path-normalization: manifest stores normalized repo-relative paths', () => {
  const dir = spaceRoot();
  const id = 'M9G7FRZ10';
  try {
    const { paths } = writeFreezeFixture(dir, id);
    mkdirSync(path.join(dir, 'fixtures', 'with spaces'), { recursive: true });
    const spaced = path.join(dir, 'fixtures', 'with spaces', 'spaced test.mjs');
    writeFileSync(spaced, 'export const spaced = true;\n');
    const relSpaced = 'fixtures/with spaces/spaced test.mjs';

    const absAlpha = path.join(dir, 'fixtures', 'alpha.test.mjs');
    const res = runCLI(dir, [
      'freeze',
      '--mission',
      id,
      absAlpha,
      relSpaced,
      paths[2],
    ]);
    assert.equal(res.code, 0, `freeze with mixed paths\n${combined(res)}`);

    const { manifest } = parseFreezeManifest(readEvidenceLines(dir, id)[0]);
    const stored = manifest.files.map((f) => f.path);
    assert.ok(stored.includes('fixtures/alpha.test.mjs'));
    assert.ok(stored.includes('fixtures/with spaces/spaced test.mjs'));
    for (const p of stored) {
      assert.doesNotMatch(p, /\\/);
      assert.doesNotMatch(p, /^\//);
      assert.doesNotMatch(p, /\.\./);
    }
  } finally {
    cleanup(dir);
  }
});

test('overlooked-freeze-edges E2 deleted-file: distinct deleted drift message', () => {
  const dir = spaceRoot();
  const id = 'M9G7FRZ11';
  try {
    const { paths, files } = writeFreezeFixture(dir, id);

    const freeze = runCLI(dir, ['freeze', '--mission', id, ...paths]);
    assert.equal(freeze.code, 0);

    unlinkSync(files.testA);

    const check = runCLI(dir, ['freeze-check', '--mission', id]);
    assert.notEqual(check.code, 0);
    assert.match(combined(check), /freeze-drift/);
    assert.match(combined(check), /deleted/i);
    assert.match(combined(check), /alpha\.test\.mjs|fixtures\/alpha/);
  } finally {
    cleanup(dir);
  }
});

test('overlooked-freeze-edges E5 empty-args: freeze with no paths is usage error', () => {
  const dir = spaceRoot();
  const id = 'M9G7FRZ12';
  try {
    writeFreezeFixture(dir, id);

    const res = runCLI(dir, ['freeze', '--mission', id]);
    assertNotStub(res, 'freeze');
    assert.notEqual(res.code, 0, `empty args must fail\n${combined(res)}`);
    assert.match(combined(res), /usage/i);
    assert.equal(readEvidenceLines(dir, id).length, 0);
  } finally {
    cleanup(dir);
  }
});

test('overlooked-freeze-edges O3 truncated-jsonl: torn last line reports parse problem', () => {
  const dir = spaceRoot();
  const id = 'M9G7FRZ13';
  try {
    const { paths } = writeFreezeFixture(dir, id);
    const freeze = runCLI(dir, ['freeze', '--mission', id, ...paths]);
    assert.equal(freeze.code, 0);

    const evidencePath = path.join(dir, '.space', 'missions', id, 'evidence.jsonl');
    const data = readFileSync(evidencePath, 'utf8');
    writeFileSync(evidencePath, `${data.slice(0, -5)}`);

    const check = runCLI(dir, ['freeze-check', '--mission', id]);
    assert.notEqual(check.code, 0);
    assert.match(combined(check), /parse|json|evidence/i);
    assert.doesNotMatch(combined(check), /Closeout ready/i);
  } finally {
    cleanup(dir);
  }
});

test('overlooked-freeze-edges O4 same-second-events: log position beats timestamp', () => {
  const dir = spaceRoot();
  const id = 'M9G7FRZ14';
  try {
    const { paths } = writeFreezeFixture(dir, id);
    const ts = '2026-01-01T12:00:00Z';

    appendEvidenceLine(
      dir,
      id,
      `{"label":"test-unit","command":"echo","output":"x\\n","ts":"${ts}","exitCode":0}\n`,
    );

    const freeze = runCLI(dir, ['freeze', '--mission', id, ...paths]);
    assertNotStub(freeze, 'freeze');
    assert.notEqual(freeze.code, 0, `freeze after test evidence must be postdated\n${combined(freeze)}`);
    assert.match(combined(freeze), /postdated-freeze/);

    // Same ts on freeze event should not override log order
    const lines = readEvidenceLines(dir, id);
    const testIdx = lines.findIndex((e) => e.label === 'test-unit');
    assert.ok(testIdx >= 0);
  } finally {
    cleanup(dir);
  }
});
