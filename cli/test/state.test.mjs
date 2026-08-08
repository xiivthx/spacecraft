/**
 * Node CLI tests for set-state / state transitions and clarify-status.
 */
import assert from 'node:assert/strict';
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

function spaceRoot() {
  const dir = mkdtempSync(path.join(os.tmpdir(), 'spacecraft-state-'));
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

function writeCurrent(root, id) {
  writeFileSync(path.join(root, '.space', 'current'), `${id}\n`);
}

function readMissionState(root, id) {
  const data = readFileSync(
    path.join(root, '.space', 'missions', id, 'mission.json'),
    'utf8',
  );
  return JSON.parse(data).state;
}

function readClarifyStatus(root, id) {
  return readFileSync(
    path.join(root, '.space', 'missions', id, 'clarify-status'),
    'utf8',
  ).trim();
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
}

function cleanup(dir) {
  rmSync(dir, { recursive: true, force: true });
}

test('set-state happy path: active → planned → in_progress → ready → shipped', () => {
  const dir = spaceRoot();
  const id = 'M07STH01';
  try {
    writeMission(dir, id, 'active');
    for (const next of ['planned', 'in_progress', 'ready', 'shipped']) {
      const res = runCLI(dir, 'set-state', id, next);
      assertNotStub(res, `set-state ${id} ${next}`);
      assert.equal(res.code, 0, `set-state ${id} → ${next} exit=${res.code}\n${combined(res)}`);
      assert.equal(readMissionState(dir, id), next, `after set-state ${next}`);
    }
  } finally {
    cleanup(dir);
  }
});

test('state alias matches set-state', () => {
  const dir = spaceRoot();
  const id = 'M07STH02';
  try {
    writeMission(dir, id, 'active');
    const res = runCLI(dir, 'state', id, 'planned');
    assertNotStub(res, 'state');
    assert.doesNotMatch(combined(res), /unknown command/i);
    assert.equal(res.code, 0, `state exit=${res.code}\n${combined(res)}`);
    assert.equal(readMissionState(dir, id), 'planned');
  } finally {
    cleanup(dir);
  }
});

test('set-state single arg uses current mission', () => {
  const dir = spaceRoot();
  const id = 'M07STAR1';
  try {
    writeMission(dir, id, 'active');
    const use = runCLI(dir, 'use', id);
    assert.equal(use.code, 0, `use exit=${use.code}\n${combined(use)}`);

    const res = runCLI(dir, 'set-state', 'planned');
    assertNotStub(res, 'set-state planned');
    assert.equal(res.code, 0, `set-state planned exit=${res.code}\n${combined(res)}`);
    assert.equal(readMissionState(dir, id), 'planned');
  } finally {
    cleanup(dir);
  }
});

test('set-state invalid single-arg fails with clear error', () => {
  const dir = spaceRoot();
  const id = 'M07STAR4';
  try {
    writeMission(dir, id, 'active');
    writeCurrent(dir, id);

    const res = runCLI(dir, 'set-state', 'not-a-state');
    assertNotStub(res, 'set-state not-a-state');
    assert.notEqual(res.code, 0, 'invalid single-arg state must be rejected');
    assert.match(combined(res).toLowerCase(), /invalid state/);
    assert.equal(readMissionState(dir, id), 'active', 'state must not mutate on reject');
  } finally {
    cleanup(dir);
  }
});

test('invalid state name rejected', () => {
  const dir = spaceRoot();
  const id = 'M07STI01';
  try {
    writeMission(dir, id, 'active');
    const res = runCLI(dir, 'set-state', id, 'nonexistent');
    assertNotStub(res, 'set-state nonexistent');
    assert.notEqual(res.code, 0, 'invalid state must be rejected');
    assert.match(combined(res).toLowerCase(), /invalid state/);
    assert.equal(readMissionState(dir, id), 'active');
  } finally {
    cleanup(dir);
  }
});

test('invalid transitions rejected without mutating state', () => {
  const cases = [
    ['active', 'in_progress'],
    ['active', 'ready'],
    ['active', 'shipped'],
    ['planned', 'ready'],
    ['planned', 'shipped'],
    ['in_progress', 'shipped'],
    ['shipped', 'active'],
    ['shipped', 'ready'],
  ];

  for (const [from, to] of cases) {
    const dir = spaceRoot();
    const id = 'M07STI02';
    try {
      writeMission(dir, id, from);
      const res = runCLI(dir, 'set-state', id, to);
      assertNotStub(res, `set-state ${from} → ${to}`);
      assert.notEqual(res.code, 0, `expected reject ${from} → ${to}`);
      assert.match(
        combined(res).toLowerCase(),
        /invalid transition/,
        `clear transition error for ${from} → ${to}\n${combined(res)}`,
      );
      assert.equal(
        readMissionState(dir, id),
        from,
        `state mutated after rejected ${from} → ${to}`,
      );
    } finally {
      cleanup(dir);
    }
  }
});

test('blocked transitions round-trip', () => {
  const dir = spaceRoot();
  const id = 'M07STB01';
  try {
    writeMission(dir, id, 'in_progress');

    const toBlocked = runCLI(dir, 'set-state', id, 'blocked');
    assertNotStub(toBlocked, 'set-state blocked');
    assert.equal(
      toBlocked.code,
      0,
      `in_progress → blocked exit=${toBlocked.code}\n${combined(toBlocked)}`,
    );
    assert.equal(readMissionState(dir, id), 'blocked');

    const back = runCLI(dir, 'set-state', id, 'in_progress');
    assertNotStub(back, 'set-state in_progress');
    assert.equal(
      back.code,
      0,
      `blocked → in_progress exit=${back.code}\n${combined(back)}`,
    );
    assert.equal(readMissionState(dir, id), 'in_progress');
  } finally {
    cleanup(dir);
  }
});

test('clarify-status accepts open|clear|deferred', () => {
  const dir = spaceRoot();
  const id = 'M07CLR01';
  try {
    writeMission(dir, id, 'active');
    writeCurrent(dir, id);

    for (const status of ['open', 'clear', 'deferred']) {
      const res = runCLI(dir, 'clarify-status', status);
      assertNotStub(res, `clarify-status ${status}`);
      assert.equal(
        res.code,
        0,
        `clarify-status ${status} exit=${res.code}\n${combined(res)}`,
      );
      assert.equal(readClarifyStatus(dir, id), status);
      assert.match(combined(res), new RegExp(status));
    }
  } finally {
    cleanup(dir);
  }
});

test('clarify-status rejects invalid status with clear error', () => {
  const dir = spaceRoot();
  const id = 'M07CLR02';
  try {
    writeMission(dir, id, 'active');
    writeCurrent(dir, id);

    const res = runCLI(dir, 'clarify-status', 'maybe');
    assertNotStub(res, 'clarify-status maybe');
    assert.notEqual(res.code, 0, 'invalid clarify-status must be rejected');
    const out = combined(res).toLowerCase();
    assert.match(out, /invalid status/);
    assert.match(out, /open\|clear\|deferred|open.*clear.*deferred/);
  } finally {
    cleanup(dir);
  }
});
