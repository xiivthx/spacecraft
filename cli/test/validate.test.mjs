/**
 * Node CLI tests for val/validate: evidence JSONL shape, outputHash integrity,
 * and --strict evidence / done-task gates.
 */
import assert from 'node:assert/strict';
import { spawnSync } from 'node:child_process';
import {
  mkdirSync,
  mkdtempSync,
  rmSync,
  writeFileSync,
} from 'node:fs';
import os from 'node:os';
import path from 'node:path';
import test from 'node:test';
import { fileURLToPath } from 'node:url';

const repoRoot = path.resolve(path.dirname(fileURLToPath(import.meta.url)), '../..');
const entryPath = path.join(repoRoot, 'cli', 'spacecraft.mjs');

/** Independent oracle: SHA-256 of "hi\n". */
const HI_OUTPUT_HASH =
  '98ea6e4f216f2fb4b69fff9b3a44842c38686ca685f3f55dc48c5d3fb1107be4';

function spaceRoot() {
  const dir = mkdtempSync(path.join(os.tmpdir(), 'spacecraft-val-'));
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

function writePlanWithTasks(root, id, tasks) {
  writeFileSync(
    path.join(root, '.space', 'missions', id, 'plan.json'),
    `${JSON.stringify({ planName: 'test', missionId: id, tasks }, null, 2)}\n`,
  );
}

function writeEvidence(root, id, body) {
  writeFileSync(
    path.join(root, '.space', 'missions', id, 'evidence.jsonl'),
    body,
  );
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

function cleanup(dir) {
  rmSync(dir, { recursive: true, force: true });
}

// --- Acceptance 1: JSONL shape + outputHash integrity ---

test('validate rejects malformed evidence JSONL', () => {
  const dir = spaceRoot();
  const id = 'M07VAL01';
  try {
    writeMission(dir, id);
    writeEvidence(
      dir,
      id,
      'not-json\n{"label":"ok","command":"echo","output":"x","ts":"2026-01-01T00:00:00Z"}\n',
    );

    for (const cmd of ['val', 'validate']) {
      const res = runCLI(dir, cmd, id);
      assertNotStub(res, cmd);
      assert.notEqual(
        res.code,
        0,
        `${cmd} accepted malformed evidence.jsonl; want nonzero\n${combined(res)}`,
      );
    }
  } finally {
    cleanup(dir);
  }
});

test('validate rejects evidence missing required fields', () => {
  const dir = spaceRoot();
  const id = 'M07VAL02';
  try {
    writeMission(dir, id);
    writeEvidence(dir, id, '{"foo":"bar"}\n');

    const res = runCLI(dir, 'val', id);
    assertNotStub(res, 'val');
    assert.notEqual(
      res.code,
      0,
      `val accepted evidence entry missing label/command/output/ts; want nonzero\n${combined(res)}`,
    );
  } finally {
    cleanup(dir);
  }
});

test('validate rejects mismatched outputHash', () => {
  const dir = spaceRoot();
  const id = 'M07VAL10';
  try {
    writeMission(dir, id);
    writeEvidence(
      dir,
      id,
      '{"label":"unit","command":"echo hi","output":"hi\\n","ts":"2026-01-01T00:00:00Z","exitCode":0,"outputHash":"0000000000000000000000000000000000000000000000000000000000000000"}\n',
    );

    for (const cmd of ['val', 'validate']) {
      const res = runCLI(dir, cmd, id);
      assertNotStub(res, cmd);
      assert.notEqual(
        res.code,
        0,
        `${cmd} accepted evidence with mismatched outputHash; want nonzero\n${combined(res)}`,
      );
      const out = combined(res);
      assert.match(out, /line 1/, `${cmd} mismatch message must identify evidence line\n${out}`);
      assert.match(
        out,
        /outputhash|hash/i,
        `${cmd} mismatch message must mention outputHash or hash\n${out}`,
      );
    }
  } finally {
    cleanup(dir);
  }
});

test('validate accepts matching outputHash', () => {
  const dir = spaceRoot();
  const id = 'M07VAL11';
  try {
    writeMission(dir, id);
    writeEvidence(
      dir,
      id,
      `{"label":"unit","command":"echo hi","output":"hi\\n","ts":"2026-01-01T00:00:00Z","exitCode":0,"outputHash":"${HI_OUTPUT_HASH}"}\n`,
    );

    const res = runCLI(dir, 'val', id);
    assertNotStub(res, 'val');
    assert.equal(
      res.code,
      0,
      `val matching outputHash exit=${res.code}\n${combined(res)}`,
    );
  } finally {
    cleanup(dir);
  }
});

test('validate accepts legacy evidence without outputHash', () => {
  const dir = spaceRoot();
  const id = 'M07VAL12';
  try {
    writeMission(dir, id);
    writeEvidence(
      dir,
      id,
      '{"label":"unit","command":"echo hi","output":"hi\\n","ts":"2026-01-01T00:00:00Z","exitCode":0}\n',
    );

    for (const cmd of ['val', 'validate']) {
      const res = runCLI(dir, cmd, id);
      assertNotStub(res, cmd);
      assert.equal(
        res.code,
        0,
        `${cmd} must accept well-formed evidence omitting outputHash; exit=${res.code}\n${combined(res)}`,
      );
    }
  } finally {
    cleanup(dir);
  }
});

// --- Acceptance 2: --strict evidence + done-task gates ---

test('validate --strict fails empty evidence', () => {
  const dir = spaceRoot();
  const id = 'M07VAL06';
  try {
    writeMission(dir, id);

    const res = runCLI(dir, 'validate', '--strict', id);
    assertNotStub(res, 'validate --strict');
    assert.notEqual(res.code, 0, `strict must fail empty evidence\n${combined(res)}`);
    assert.match(
      combined(res),
      /evidence/i,
      `want evidence mention\n${combined(res)}`,
    );
  } finally {
    cleanup(dir);
  }
});

test('validate --strict fails missing exitCode', () => {
  const dir = spaceRoot();
  const id = 'M07VAL07';
  try {
    writeMission(dir, id);
    writeEvidence(
      dir,
      id,
      '{"label":"unit","command":"echo","output":"x","ts":"2026-01-01T00:00:00Z"}\n',
    );

    const res = runCLI(dir, 'validate', '--strict', id);
    assertNotStub(res, 'validate --strict');
    assert.notEqual(res.code, 0, `strict must fail missing exitCode\n${combined(res)}`);
    assert.match(
      combined(res),
      /exitCode/,
      `want exitCode mention\n${combined(res)}`,
    );
  } finally {
    cleanup(dir);
  }
});

test('validate --strict fails done task without matching exitCode 0 evidence', () => {
  const dir = spaceRoot();
  const id = 'M07VAL08';
  try {
    writeMission(dir, id);
    writePlanWithTasks(dir, id, [
      {
        id: 'T1',
        title: 'Do thing',
        status: 'done',
        evidence: ['t1-pass'],
      },
    ]);
    writeEvidence(
      dir,
      id,
      '{"label":"other","command":"echo","output":"x","ts":"2026-01-01T00:00:00Z","exitCode":0}\n',
    );

    const res = runCLI(dir, 'validate', '--strict', id);
    assertNotStub(res, 'validate --strict');
    assert.notEqual(
      res.code,
      0,
      `strict must fail done task without matching evidence\n${combined(res)}`,
    );
  } finally {
    cleanup(dir);
  }
});

test('validate --strict passes done task with matching exitCode 0 evidence', () => {
  const dir = spaceRoot();
  const id = 'M07VAL09';
  try {
    writeMission(dir, id);
    writePlanWithTasks(dir, id, [
      {
        id: 'T1',
        title: 'Do thing',
        status: 'done',
        evidence: ['t1-pass'],
      },
    ]);
    writeEvidence(
      dir,
      id,
      '{"label":"t1-pass","command":"echo","output":"x","ts":"2026-01-01T00:00:00Z","exitCode":0}\n',
    );

    const res = runCLI(dir, 'validate', '--strict', id);
    assertNotStub(res, 'validate --strict');
    assert.equal(
      res.code,
      0,
      `strict pass exit=${res.code}\n${combined(res)}`,
    );
  } finally {
    cleanup(dir);
  }
});
