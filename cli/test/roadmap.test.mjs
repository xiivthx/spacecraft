/**
 * Node CLI tests for roadmap/map lifecycle (create, use, current, next).
 */
import assert from 'node:assert/strict';
import { spawnSync } from 'node:child_process';
import {
  existsSync,
  mkdirSync,
  mkdtempSync,
  readdirSync,
  rmSync,
  writeFileSync,
  readFileSync,
} from 'node:fs';
import os from 'node:os';
import path from 'node:path';
import test from 'node:test';
import { fileURLToPath } from 'node:url';

const repoRoot = path.resolve(path.dirname(fileURLToPath(import.meta.url)), '../..');
const entryPath = path.join(repoRoot, 'cli', 'spacecraft.mjs');

function spaceRoot() {
  const dir = mkdtempSync(path.join(os.tmpdir(), 'spacecraft-roadmap-'));
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
}

function findRoadmapID(dir) {
  const roadmapsDir = path.join(dir, '.space', 'roadmaps');
  const entries = readdirSync(roadmapsDir).filter((name) => name.endsWith('.json'));
  assert.ok(entries.length > 0, 'no roadmap json found');
  return path.basename(entries[0], '.json');
}

function cleanup(dir) {
  rmSync(dir, { recursive: true, force: true });
}

test('roadmap lifecycle: new|add|show|ls|rm|archive', () => {
  const dir = spaceRoot();
  try {
    writeMission(dir, 'M07MAP01', 'active');

    const create = runCLI(dir, 'roadmap', 'new', 'Cursor Native', '--desc', 'restore cli');
    assertNotStub(create, 'roadmap new');
    assert.equal(create.code, 0, `roadmap new exit=${create.code}\n${combined(create)}`);

    let rid = 'cursor-native';
    let roadmapPath = path.join(dir, '.space', 'roadmaps', `${rid}.json`);
    if (!existsSync(roadmapPath)) {
      rid = findRoadmapID(dir);
      roadmapPath = path.join(dir, '.space', 'roadmaps', `${rid}.json`);
    }
    assert.ok(existsSync(roadmapPath), `roadmap file missing after new\n${create.stdout}`);

    const add = runCLI(dir, 'roadmap', 'add', rid, 'M07MAP01', '--desc', 'first mission');
    assertNotStub(add, 'roadmap add');
    assert.equal(add.code, 0, `roadmap add exit=${add.code}\n${combined(add)}`);

    const show = runCLI(dir, 'roadmap', 'show', rid);
    assertNotStub(show, 'roadmap show');
    assert.equal(show.code, 0, `roadmap show exit=${show.code}\n${combined(show)}`);
    assert.match(show.stdout, /M07MAP01/);

    let ls = runCLI(dir, 'roadmap', 'ls');
    if (ls.code !== 0) {
      ls = runCLI(dir, 'roadmap', 'list');
    }
    assertNotStub(ls, 'roadmap ls');
    assert.equal(ls.code, 0, `roadmap ls/list exit=${ls.code}\n${combined(ls)}`);
    assert.ok(
      ls.stdout.includes(rid) || ls.stdout.includes('Cursor Native'),
      `ls missing roadmap\n${ls.stdout}`,
    );

    let rm = runCLI(dir, 'roadmap', 'rm', rid, 'M07MAP01');
    if (rm.code !== 0) {
      rm = runCLI(dir, 'roadmap', 'remove', rid, 'M07MAP01');
    }
    assertNotStub(rm, 'roadmap rm');
    assert.equal(rm.code, 0, `roadmap rm/remove exit=${rm.code}\n${combined(rm)}`);

    const show2 = runCLI(dir, 'roadmap', 'show', rid);
    assert.equal(show2.code, 0, `roadmap show after rm exit=${show2.code}\n${combined(show2)}`);
    try {
      const raw = JSON.parse(show2.stdout);
      const missions = raw.missions ?? [];
      assert.equal(missions.length, 0, `missions still present after rm: ${JSON.stringify(missions)}`);
    } catch {
      // Non-JSON show output is acceptable if missions are clearly gone.
      assert.doesNotMatch(show2.stdout, /M07MAP01/);
    }

    const arch = runCLI(dir, 'roadmap', 'archive', rid);
    assertNotStub(arch, 'roadmap archive');
    assert.equal(arch.code, 0, `roadmap archive exit=${arch.code}\n${combined(arch)}`);
    assert.ok(!existsSync(roadmapPath), 'roadmap still in roadmaps/ after archive');
  } finally {
    cleanup(dir);
  }
});

test('map alias lifecycle: new|add|ls', () => {
  const dir = spaceRoot();
  try {
    writeMission(dir, 'M07MAP02', 'active');

    const res = runCLI(dir, 'map', 'new', 'Alias Roadmap');
    assertNotStub(res, 'map new');
    assert.doesNotMatch(combined(res), /unknown command/i);
    assert.equal(res.code, 0, `map new exit=${res.code}\n${combined(res)}`);

    let rid = 'alias-roadmap';
    if (!existsSync(path.join(dir, '.space', 'roadmaps', `${rid}.json`))) {
      rid = findRoadmapID(dir);
    }

    const add = runCLI(dir, 'map', 'add', rid, 'M07MAP02');
    assertNotStub(add, 'map add');
    assert.equal(add.code, 0, `map add exit=${add.code}\n${combined(add)}`);

    const ls = runCLI(dir, 'map', 'ls');
    assertNotStub(ls, 'map ls');
    assert.equal(ls.code, 0, `map ls exit=${ls.code}\n${combined(ls)}`);
  } finally {
    cleanup(dir);
  }
});

test('map next skips ready and returns planned', () => {
  const dir = spaceRoot();
  try {
    writeMission(dir, 'M07NEXT1', 'planned');
    writeMission(dir, 'M07NEXT2', 'ready');

    const create = runCLI(dir, 'map', 'new', 'Next Order');
    assertNotStub(create, 'map new');
    assert.equal(create.code, 0, `map new exit=${create.code}\n${combined(create)}`);

    let rid = 'next-order';
    if (!existsSync(path.join(dir, '.space', 'roadmaps', `${rid}.json`))) {
      rid = findRoadmapID(dir);
    }

    const addReady = runCLI(dir, 'map', 'add', rid, 'M07NEXT2', '--desc', 'ready mission');
    assert.equal(addReady.code, 0, `map add ready exit=${addReady.code}\n${combined(addReady)}`);
    const addPlanned = runCLI(dir, 'map', 'add', rid, 'M07NEXT1', '--desc', 'planned mission');
    assert.equal(
      addPlanned.code,
      0,
      `map add planned exit=${addPlanned.code}\n${combined(addPlanned)}`,
    );

    const next = runCLI(dir, 'map', 'next', rid);
    assertNotStub(next, 'map next');
    assert.equal(next.code, 0, `map next exit=${next.code}\n${combined(next)}`);
    assert.match(next.stdout, /M07NEXT1: planned mission \(state=planned\)/);
    assert.doesNotMatch(next.stdout, /M07NEXT2/);
  } finally {
    cleanup(dir);
  }
});

test('map next reports all complete when none incomplete', () => {
  const dir = spaceRoot();
  try {
    writeMission(dir, 'M07DONE1', 'ready');
    writeMission(dir, 'M07DONE2', 'shipped');

    const create = runCLI(dir, 'map', 'new', 'All Done');
    assertNotStub(create, 'map new');
    assert.equal(create.code, 0, `map new exit=${create.code}\n${combined(create)}`);

    let rid = 'all-done';
    if (!existsSync(path.join(dir, '.space', 'roadmaps', `${rid}.json`))) {
      rid = findRoadmapID(dir);
    }

    assert.equal(
      runCLI(dir, 'map', 'add', rid, 'M07DONE1', '--desc', 'ready one').code,
      0,
    );
    assert.equal(
      runCLI(dir, 'map', 'add', rid, 'M07DONE2', '--desc', 'shipped one').code,
      0,
    );

    const next = runCLI(dir, 'map', 'next', rid);
    assertNotStub(next, 'map next');
    assert.equal(next.code, 0, `map next exit=${next.code}\n${combined(next)}`);
    assert.match(next.stdout, /All missions complete\./);
  } finally {
    cleanup(dir);
  }
});

test('map use/current round-trip current-roadmap', () => {
  const dir = spaceRoot();
  try {
    const create = runCLI(dir, 'map', 'new', 'Current Map');
    assertNotStub(create, 'map new');
    assert.equal(create.code, 0, `map new exit=${create.code}\n${combined(create)}`);

    let rid = 'current-map';
    if (!existsSync(path.join(dir, '.space', 'roadmaps', `${rid}.json`))) {
      rid = findRoadmapID(dir);
    }

    const use = runCLI(dir, 'map', 'use', rid);
    assertNotStub(use, 'map use');
    assert.equal(use.code, 0, `map use exit=${use.code}\n${combined(use)}`);

    const currentPath = path.join(dir, '.space', 'current-roadmap');
    assert.ok(existsSync(currentPath), 'current-roadmap missing after map use');
    assert.equal(readFileSync(currentPath, 'utf8').trim(), rid);

    const cur = runCLI(dir, 'map', 'current');
    assertNotStub(cur, 'map current');
    assert.equal(cur.code, 0, `map current exit=${cur.code}\n${combined(cur)}`);
    assert.equal(cur.stdout.trim(), rid);
  } finally {
    cleanup(dir);
  }
});
