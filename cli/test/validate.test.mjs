/**
 * Node CLI tests for val/validate: evidence JSONL shape, outputHash integrity,
 * truncated sidecar hash checks, and --strict evidence / done-task gates.
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

/** Independent oracle: SHA-256 of "hi\n". */
const HI_OUTPUT_HASH =
  '98ea6e4f216f2fb4b69fff9b3a44842c38686ca685f3f55dc48c5d3fb1107be4';

/** Truncate marker aligned with evi discuss lock (trailing on JSONL `output`). */
const EVI_TRUNCATE_MARKER = '\n...[truncated]';

/** Fixture full raw for truncated-evidence validate cases (not from production). */
const TRUNC_FIXTURE_RAW =
  'FULL-RAW-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx\n';
/** Truncated JSONL `output` field (prefix + marker); hash differs from full raw. */
const TRUNC_FIXTURE_OUTPUT =
  `${TRUNC_FIXTURE_RAW.slice(0, 32)}${EVI_TRUNCATE_MARKER}`;
const TRUNC_FIXTURE_RAW_HASH = createHash('sha256')
  .update(TRUNC_FIXTURE_RAW, 'utf8')
  .digest('hex');
const TRUNC_FIXTURE_OUTPUT_HASH = createHash('sha256')
  .update(TRUNC_FIXTURE_OUTPUT, 'utf8')
  .digest('hex');
const TRUNC_FIXTURE_RAW_PATH = 'evidence-raw/2026-01-01T00-00-00Z-unit.log';

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

function writeSidecar(root, id, relPath, body) {
  const abs = path.join(root, '.space', 'missions', id, relPath);
  mkdirSync(path.dirname(abs), { recursive: true });
  writeFileSync(abs, body, 'utf8');
}

function truncatedEvidenceLine({
  outputHash = TRUNC_FIXTURE_RAW_HASH,
  outputRawPath = TRUNC_FIXTURE_RAW_PATH,
  output = TRUNC_FIXTURE_OUTPUT,
} = {}) {
  return `${JSON.stringify({
    label: 'unit',
    command: 'echo truncated',
    output,
    ts: '2026-01-01T00:00:00Z',
    exitCode: 0,
    outputTruncated: true,
    outputBytes: Buffer.byteLength(TRUNC_FIXTURE_RAW, 'utf8'),
    outputRawPath,
    outputHash,
  })}\n`;
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

// --- T3: truncated sidecar hash + legacy non-truncated rules ---

assert.notEqual(
  TRUNC_FIXTURE_RAW_HASH,
  TRUNC_FIXTURE_OUTPUT_HASH,
  'fixture raw hash must differ from truncated-output hash (oracle for must-not-hash-truncated)',
);

test('validate accepts truncated evidence when sidecar bytes match outputHash', () => {
  const dir = spaceRoot();
  const id = 'M08VAL20';
  try {
    writeMission(dir, id);
    writeSidecar(dir, id, TRUNC_FIXTURE_RAW_PATH, TRUNC_FIXTURE_RAW);
    writeEvidence(dir, id, truncatedEvidenceLine());

    for (const cmd of ['val', 'validate']) {
      const res = runCLI(dir, cmd, id);
      assertNotStub(res, cmd);
      assert.equal(
        res.code,
        0,
        `${cmd} must accept outputTruncated:true when sidecar matches outputHash (must not hash truncated output); exit=${res.code}\n${combined(res)}`,
      );
    }
  } finally {
    cleanup(dir);
  }
});

test('validate rejects truncated evidence when sidecar is missing', () => {
  const dir = spaceRoot();
  const id = 'M08VAL21';
  try {
    writeMission(dir, id);
    writeEvidence(dir, id, truncatedEvidenceLine());

    for (const cmd of ['val', 'validate']) {
      const res = runCLI(dir, cmd, id);
      assertNotStub(res, cmd);
      assert.notEqual(
        res.code,
        0,
        `${cmd} must reject outputTruncated:true when sidecar is missing; exit=${res.code}\n${combined(res)}`,
      );
      const out = combined(res);
      assert.match(out, /line 1/, `${cmd} missing-sidecar message must identify evidence line\n${out}`);
      assert.match(
        out,
        /outputhash|hash|sidecar|outputrawpath|missing/i,
        `${cmd} missing-sidecar message must mention hash/sidecar\n${out}`,
      );
    }
  } finally {
    cleanup(dir);
  }
});

test('validate rejects truncated evidence when sidecar contents do not match outputHash', () => {
  const dir = spaceRoot();
  const id = 'M08VAL22';
  try {
    writeMission(dir, id);
    writeSidecar(dir, id, TRUNC_FIXTURE_RAW_PATH, 'WRONG-SIDECAR-BYTES\n');
    writeEvidence(dir, id, truncatedEvidenceLine());

    for (const cmd of ['val', 'validate']) {
      const res = runCLI(dir, cmd, id);
      assertNotStub(res, cmd);
      assert.notEqual(
        res.code,
        0,
        `${cmd} must reject outputTruncated:true when sidecar bytes mismatch outputHash; exit=${res.code}\n${combined(res)}`,
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

test('validate rejects truncated evidence when only truncated output field matches outputHash', () => {
  const dir = spaceRoot();
  const id = 'M08VAL23';
  try {
    writeMission(dir, id);
    writeSidecar(dir, id, TRUNC_FIXTURE_RAW_PATH, TRUNC_FIXTURE_RAW);
    // Wrong product state: outputHash is of truncated JSONL field, not full raw.
    // Validate must hash sidecar (or fail), never accept by hashing truncated `output`.
    writeEvidence(
      dir,
      id,
      truncatedEvidenceLine({ outputHash: TRUNC_FIXTURE_OUTPUT_HASH }),
    );

    for (const cmd of ['val', 'validate']) {
      const res = runCLI(dir, cmd, id);
      assertNotStub(res, cmd);
      assert.notEqual(
        res.code,
        0,
        `${cmd} must not accept by hashing truncated output when outputTruncated:true; exit=${res.code}\n${combined(res)}`,
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

test('validate rejects non-truncated mismatched outputHash (T3 legacy rule)', () => {
  const dir = spaceRoot();
  const id = 'M08VAL24';
  try {
    writeMission(dir, id);
    writeEvidence(
      dir,
      id,
      '{"label":"unit","command":"echo hi","output":"hi\\n","ts":"2026-01-01T00:00:00Z","exitCode":0,"outputHash":"0000000000000000000000000000000000000000000000000000000000000000"}\n',
    );

    const res = runCLI(dir, 'val', id);
    assertNotStub(res, 'val');
    assert.notEqual(
      res.code,
      0,
      `val must reject non-truncated outputHash mismatch; exit=${res.code}\n${combined(res)}`,
    );
  } finally {
    cleanup(dir);
  }
});

test('validate accepts legacy omit-hash and matching-hash non-truncated lines (T3)', () => {
  const dir = spaceRoot();
  const id = 'M08VAL25';
  try {
    writeMission(dir, id);
    writeEvidence(
      dir,
      id,
      [
        '{"label":"omit","command":"echo hi","output":"hi\\n","ts":"2026-01-01T00:00:00Z","exitCode":0}',
        `{"label":"match","command":"echo hi","output":"hi\\n","ts":"2026-01-01T00:00:00Z","exitCode":0,"outputHash":"${HI_OUTPUT_HASH}"}`,
        '',
      ].join('\n'),
    );

    const res = runCLI(dir, 'validate', id);
    assertNotStub(res, 'validate');
    assert.equal(
      res.code,
      0,
      `validate must accept legacy omit-hash + matching-hash lines; exit=${res.code}\n${combined(res)}`,
    );
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

// --- T6 / neg-validate-framing (approved-scenarios frozen literal) ---

/**
 * neg-validate-framing: spacecraft validate --help and README still frame
 * validate as mission evidence / not-doc-drift (not sold as doc-drift or
 * 10X validate). Oracles: design-contract Validate framing + approved-scenarios.
 */
test('neg-validate-framing: validate --help and README frame mission evidence / not-doc-drift', () => {
  const dir = spaceRoot();
  try {
    for (const cmd of ['validate', 'val']) {
      const res = runCLI(dir, cmd, '--help');
      assertNotStub(res, `${cmd} --help`);
      assert.equal(
        res.code,
        0,
        `${cmd} --help must exit 0\n${combined(res)}`,
      );
      const out = combined(res);
      assert.match(
        out,
        /Validate mission artifacts and evidence/,
        `${cmd} --help must describe mission artifacts and evidence\n${out}`,
      );
      assert.match(
        out,
        /not-doc-drift/,
        `${cmd} --help must keep not-doc-drift posture\n${out}`,
      );
      assert.match(
        out,
        /not-10X-validate/,
        `${cmd} --help must keep not-10X-validate posture\n${out}`,
      );
      // Negative sell: bare "10X validate" / "doc-drift" as product pitch
      assert.doesNotMatch(
        out,
        /\b10X validate\b/i,
        `${cmd} --help must not sell validate as 10X validate\n${out}`,
      );
      assert.doesNotMatch(
        out,
        /(?:^|[^-])doc-drift\b/m,
        `${cmd} --help must not sell validate as doc-drift (use not-doc-drift)\n${out}`,
      );
    }

    const readme = readFileSync(path.join(repoRoot, 'README.md'), 'utf8');
    assert.match(
      readme,
      /Validate mission artifacts and evidence/,
      'README must describe validate as mission artifacts and evidence',
    );
    assert.match(
      readme,
      /not-doc-drift/,
      'README must keep not-doc-drift posture for validate',
    );
    assert.doesNotMatch(
      readme,
      /\b10X validate\b/i,
      'README must not sell validate as 10X validate',
    );
  } finally {
    cleanup(dir);
  }
});
