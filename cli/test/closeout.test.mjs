/**
 * Node CLI tests for closeout-check / ship-check gates and judge-break wiring.
 */
import assert from 'node:assert/strict';
import { spawnSync } from 'node:child_process';
import {
  chmodSync,
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
const judgeBreakScript = path.join(repoRoot, 'scripts', 'check-judge-break.sh');

function spaceRoot() {
  const dir = mkdtempSync(path.join(os.tmpdir(), 'spacecraft-closeout-'));
  mkdirSync(path.join(dir, '.space', 'missions'), { recursive: true });
  mkdirSync(path.join(dir, '.space', 'roadmaps'), { recursive: true });
  return dir;
}

function writeMission(root, id, state = 'ready') {
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

function writeCurrent(root, id) {
  writeFileSync(path.join(root, '.space', 'current'), `${id}\n`);
}

function goodReview() {
  return {
    status: 'ready',
    findings: [],
    releaseReadiness: {
      changelog: { status: 'ready' },
      specNote: { status: 'ready' },
    },
  };
}

function writeReviewJSON(root, id, review) {
  writeFileSync(
    path.join(root, '.space', 'missions', id, 'review.json'),
    `${JSON.stringify(review, null, 2)}\n`,
  );
}

function writeEvidenceLine(root, id, line) {
  writeFileSync(
    path.join(root, '.space', 'missions', id, 'evidence.jsonl'),
    line,
  );
}

/** Ready mission that passes when SPACECRAFT_CLOSEOUT_SKIP_CHANGELOG=1. */
function writeReadyCloseoutMission(root, id) {
  writeMission(root, id, 'ready');
  writeCurrent(root, id);
  writeEvidenceLine(
    root,
    id,
    '{"label":"unit","command":"echo hi","output":"hi\\n","ts":"2026-01-01T00:00:00Z","exitCode":0}\n',
  );
  writeReviewJSON(root, id, goodReview());
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

function closeout(dir, cmd = 'closeout-check', { skipChangelog = true } = {}) {
  const env = skipChangelog
    ? { SPACECRAFT_CLOSEOUT_SKIP_CHANGELOG: '1' }
    : { SPACECRAFT_CLOSEOUT_SKIP_CHANGELOG: '' };
  return runCLI(dir, [cmd], env);
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

function writeNodeCliWrapper() {
  const wrapDir = mkdtempSync(path.join(os.tmpdir(), 'spacecraft-node-bin-'));
  const wrap = path.join(wrapDir, 'spacecraft');
  writeFileSync(
    wrap,
    `#!/bin/sh\nexec ${JSON.stringify(process.execPath)} ${JSON.stringify(entryPath)} "$@"\n`,
  );
  chmodSync(wrap, 0o755);
  return { wrapDir, wrap };
}

// --- Acceptance 1: failure gates ---

test('closeout-check fails missing review.json', () => {
  const dir = spaceRoot();
  const id = 'M07CLO01';
  try {
    writeMission(dir, id, 'ready');
    writeCurrent(dir, id);
    writeEvidenceLine(
      dir,
      id,
      '{"label":"unit","command":"echo","output":"x","ts":"2026-01-01T00:00:00Z","exitCode":0}\n',
    );

    const res = closeout(dir);
    assertNotStub(res, 'closeout-check');
    assert.notEqual(res.code, 0, `expected fail without review.json\n${combined(res)}`);
    assert.match(res.stdout, new RegExp(`Closeout blocked for ${id}`));
    assert.match(res.stdout, /missing review\.json/);
  } finally {
    cleanup(dir);
  }
});

test('closeout-check fails wrong state', () => {
  const dir = spaceRoot();
  const id = 'M07CLO02';
  try {
    writeReadyCloseoutMission(dir, id);
    const missionPath = path.join(dir, '.space', 'missions', id, 'mission.json');
    const m = JSON.parse(readFileSync(missionPath, 'utf8'));
    m.state = 'in_progress';
    writeFileSync(missionPath, `${JSON.stringify(m, null, 2)}\n`);

    const res = closeout(dir);
    assertNotStub(res, 'closeout-check');
    assert.notEqual(res.code, 0, `expected fail for wrong state\n${combined(res)}`);
    assert.match(res.stdout, /state is in_progress/);
  } finally {
    cleanup(dir);
  }
});

test('closeout-check fails clarify open', () => {
  const dir = spaceRoot();
  const id = 'M07CLO03';
  try {
    writeReadyCloseoutMission(dir, id);
    writeFileSync(
      path.join(dir, '.space', 'missions', id, 'clarify-status'),
      'open\n',
    );

    const res = closeout(dir);
    assertNotStub(res, 'closeout-check');
    assert.notEqual(res.code, 0, `expected fail for clarify open\n${combined(res)}`);
    assert.match(res.stdout, /clarify/i);
  } finally {
    cleanup(dir);
  }
});

test('closeout-check fails empty evidence', () => {
  const dir = spaceRoot();
  const id = 'M07CLO04';
  try {
    writeReadyCloseoutMission(dir, id);
    writeEvidenceLine(dir, id, '');

    const res = closeout(dir);
    assertNotStub(res, 'closeout-check');
    assert.notEqual(res.code, 0, `expected fail for empty evidence\n${combined(res)}`);
    assert.match(combined(res), /evidence/i);
  } finally {
    cleanup(dir);
  }
});

test('closeout-check fails missing exitCode', () => {
  const dir = spaceRoot();
  const id = 'M07CLO05';
  try {
    writeReadyCloseoutMission(dir, id);
    writeEvidenceLine(
      dir,
      id,
      '{"label":"unit","command":"echo","output":"x","ts":"2026-01-01T00:00:00Z"}\n',
    );

    const res = closeout(dir);
    assertNotStub(res, 'closeout-check');
    assert.notEqual(res.code, 0, `expected fail for missing exitCode\n${combined(res)}`);
    assert.match(res.stdout, /exitCode/);
  } finally {
    cleanup(dir);
  }
});

test('closeout-check fails critical finding', () => {
  const dir = spaceRoot();
  const id = 'M07CLO06';
  try {
    writeReadyCloseoutMission(dir, id);
    const r = goodReview();
    r.findings = [{ severity: 'critical', blocksShip: false, summary: 'bad' }];
    writeReviewJSON(dir, id, r);

    const res = closeout(dir);
    assertNotStub(res, 'closeout-check');
    assert.notEqual(res.code, 0, `expected fail for critical finding\n${combined(res)}`);
    assert.match(res.stdout, /critical/i);
  } finally {
    cleanup(dir);
  }
});

test('closeout-check fails blocksShip finding', () => {
  const dir = spaceRoot();
  const id = 'M07CLO07';
  try {
    writeReadyCloseoutMission(dir, id);
    const r = goodReview();
    r.findings = [{ severity: 'medium', blocksShip: true, summary: 'blocker' }];
    writeReviewJSON(dir, id, r);

    const res = closeout(dir);
    assertNotStub(res, 'closeout-check');
    assert.notEqual(res.code, 0, `expected fail for blocksShip\n${combined(res)}`);
    assert.ok(
      /blocksShip/.test(res.stdout) || /block/i.test(res.stdout),
      `want blocksShip problem\n${res.stdout}`,
    );
  } finally {
    cleanup(dir);
  }
});

test('closeout-check fails warning finding', () => {
  const dir = spaceRoot();
  const id = 'M07CLO12';
  try {
    writeReadyCloseoutMission(dir, id);
    const r = goodReview();
    r.findings = [{ severity: 'warning', blocksShip: false, summary: 'nit' }];
    writeReviewJSON(dir, id, r);

    const res = closeout(dir);
    assertNotStub(res, 'closeout-check');
    assert.notEqual(res.code, 0, `expected fail for warning finding\n${combined(res)}`);
    assert.match(res.stdout, /warning/i);
  } finally {
    cleanup(dir);
  }
});

test('closeout-check fails low severity finding', () => {
  const dir = spaceRoot();
  const id = 'M07CLO13';
  try {
    writeReadyCloseoutMission(dir, id);
    const r = goodReview();
    r.findings = [{ severity: 'low', blocksShip: false, summary: 'nit' }];
    writeReviewJSON(dir, id, r);

    const res = closeout(dir);
    assertNotStub(res, 'closeout-check');
    assert.notEqual(res.code, 0, `expected fail for low severity finding\n${combined(res)}`);
    assert.match(res.stdout, /low/i);
  } finally {
    cleanup(dir);
  }
});

test('closeout-check fails deferred CHANGELOG readiness', () => {
  const dir = spaceRoot();
  const id = 'M07CLO08';
  try {
    writeReadyCloseoutMission(dir, id);
    const r = goodReview();
    r.releaseReadiness = {
      changelog: { status: 'deferred' },
      specNote: { status: 'ready' },
    };
    writeReviewJSON(dir, id, r);

    const res = closeout(dir);
    assertNotStub(res, 'closeout-check');
    assert.notEqual(res.code, 0, `expected fail for deferred changelog\n${combined(res)}`);
    assert.match(res.stdout, /changelog/i);
  } finally {
    cleanup(dir);
  }
});

test('closeout-check fails without CHANGELOG skip in non-git dir', () => {
  const dir = spaceRoot();
  const id = 'M07CLO09';
  try {
    writeReadyCloseoutMission(dir, id);

    const res = closeout(dir, 'closeout-check', { skipChangelog: false });
    assertNotStub(res, 'closeout-check');
    assert.notEqual(res.code, 0, `expected fail without CHANGELOG skip\n${combined(res)}`);
    const out = combined(res).toUpperCase();
    assert.ok(
      out.includes('CHANGELOG') || out.includes('GIT'),
      `want CHANGELOG or git problem\n${combined(res)}`,
    );
  } finally {
    cleanup(dir);
  }
});

// --- Acceptance 2: happy path, ship-check alias, judge-break via Node ---

test('closeout-check passes full ready fixture', () => {
  const dir = spaceRoot();
  const id = 'M07CLO10';
  try {
    writeReadyCloseoutMission(dir, id);

    const res = closeout(dir);
    assertNotStub(res, 'closeout-check');
    assert.equal(res.code, 0, `expected pass\n${combined(res)}`);
    assert.match(res.stdout, new RegExp(`Closeout ready for ${id}`));
  } finally {
    cleanup(dir);
  }
});

test('ship-check alias dispatches', () => {
  const dir = spaceRoot();
  const id = 'M07CLO11';
  try {
    writeReadyCloseoutMission(dir, id);

    const res = closeout(dir, 'ship-check');
    assertNotStub(res, 'ship-check');
    assert.doesNotMatch(combined(res), /unknown command/i);
    assert.equal(res.code, 0, `ship-check exit=${res.code}\n${combined(res)}`);
  } finally {
    cleanup(dir);
  }
});

test('check-judge-break.sh rejects known-bad packs using Node CLI wrapper', () => {
  const { wrapDir, wrap } = writeNodeCliWrapper();
  try {
    const result = spawnSync('sh', [judgeBreakScript, repoRoot, wrap], {
      cwd: repoRoot,
      encoding: 'utf8',
    });
    const out = `${result.stdout ?? ''}${result.stderr ?? ''}`;
    assert.equal(
      result.status,
      0,
      `check-judge-break via Node CLI failed\nexit=${result.status}\n${out}`,
    );
    assert.match(out, /ok: judge-break/i);
  } finally {
    cleanup(wrapDir);
  }
});

// --- D6 disposition leak predicates (RED until closeoutDispositionProblems lands) ---

async function loadDispositionProblems() {
  const mod = await import('../lib/closeout.mjs');
  assert.equal(
    typeof mod.closeoutDispositionProblems,
    'function',
    'closeoutDispositionProblems must be exported from closeout.mjs',
  );
  return mod.closeoutDispositionProblems;
}

function writeDispositionMission(root, id, files) {
  writeReadyCloseoutMission(root, id);
  const dir = path.join(root, '.space', 'missions', id);
  for (const [name, body] of Object.entries(files)) {
    writeFileSync(path.join(dir, name), body);
  }
  return dir;
}

test('closeoutDispositionProblems rejects false-consensus without dissent labels', async () => {
  const closeoutDispositionProblems = await loadDispositionProblems();
  const dir = spaceRoot();
  const id = 'M07CLOFC';
  try {
    const missionDir = writeDispositionMission(dir, id, {
      'judge-summary.json': `${JSON.stringify(
        {
          verdict: 'VERIFIED',
          status: 'ready',
          hunts: [
            { id: 'hunt-scope', summary: 'claimed clean without dissent label' },
          ],
        },
        null,
        2,
      )}\n`,
    });

    const problems = closeoutDispositionProblems(missionDir);
    assert.ok(
      problems.some((p) => /false-consensus/i.test(p)),
      `want false-consensus problem, got ${JSON.stringify(problems)}`,
    );
  } finally {
    cleanup(dir);
  }
});

test('closeoutDispositionProblems rejects charitable-reviewer builderRationale', async () => {
  const closeoutDispositionProblems = await loadDispositionProblems();
  const dir = spaceRoot();
  const id = 'M07CLOCHR';
  try {
    const missionDir = writeDispositionMission(dir, id, {
      'judge-summary.json': `${JSON.stringify(
        {
          verdict: 'VERIFIED',
          status: 'ready',
          builderRationale: 'please trust the narrative',
          hunts: [],
        },
        null,
        2,
      )}\n`,
    });

    const problems = closeoutDispositionProblems(missionDir);
    assert.ok(
      problems.some((p) => /charitable-reviewer/i.test(p)),
      `want charitable-reviewer problem, got ${JSON.stringify(problems)}`,
    );
  } finally {
    cleanup(dir);
  }
});

test('closeoutDispositionProblems rejects silent-mutation-skip', async () => {
  const closeoutDispositionProblems = await loadDispositionProblems();
  const dir = spaceRoot();
  const id = 'M07CLOSMS';
  try {
    const missionDir = writeDispositionMission(dir, id, {
      'decisions.md': '# Decisions\n\n```\nMutation: required\n```\n',
    });

    const problems = closeoutDispositionProblems(missionDir);
    assert.ok(
      problems.some((p) => /silent-mutation-skip/i.test(p)),
      `want silent-mutation-skip problem, got ${JSON.stringify(problems)}`,
    );
  } finally {
    cleanup(dir);
  }
});

test('closeoutDispositionProblems rejects retroactive-oracle-change', async () => {
  const closeoutDispositionProblems = await loadDispositionProblems();
  const dir = spaceRoot();
  const id = 'M07CLOROC';
  try {
    const missionDir = writeDispositionMission(dir, id, {
      'decisions.md':
        '# Decisions\n\nExpected literal edited after freeze without Scenario oracle change.\n',
      'approved-scenarios.md':
        '# Approved scenarios\n\n| id | status | expected |\n| --- | --- | --- |\n| S1 | thawed | Expected literal edited |\n\n## Freeze footer\n\n```\nApproved-scenarios: frozen-from-contract\nFreeze stamp: 2026-01-01T00:00:00Z\n```\n',
    });

    const problems = closeoutDispositionProblems(missionDir);
    assert.ok(
      problems.some((p) => /retroactive-oracle-change/i.test(p)),
      `want retroactive-oracle-change problem, got ${JSON.stringify(problems)}`,
    );
  } finally {
    cleanup(dir);
  }
});

test('closeout-check rejects D6 false-consensus pack shape via CLI', () => {
  const dir = spaceRoot();
  const id = 'M07CLOFCcli';
  try {
    writeDispositionMission(dir, id, {
      'judge-summary.json': `${JSON.stringify(
        {
          verdict: 'VERIFIED',
          status: 'ready',
          hunts: [{ id: 'h1', summary: 'no dissent labels' }],
        },
        null,
        2,
      )}\n`,
    });

    const res = closeout(dir);
    assertNotStub(res, 'closeout-check');
    assert.notEqual(res.code, 0, `expected fail for false-consensus\n${combined(res)}`);
    assert.match(combined(res), /false-consensus/i);
  } finally {
    cleanup(dir);
  }
});
