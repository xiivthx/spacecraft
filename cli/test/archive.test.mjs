/**
 * Node CLI tests for archive and roadmap next-hint after archive.
 */
import assert from 'node:assert/strict';
import { spawnSync } from 'node:child_process';
import {
  existsSync,
  mkdirSync,
  mkdtempSync,
  readdirSync,
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
  const dir = mkdtempSync(path.join(os.tmpdir(), 'spacecraft-archive-'));
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

function readCurrent(root) {
  const p = path.join(root, '.space', 'current');
  if (!existsSync(p)) return '';
  return readFileSync(p, 'utf8').trim();
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

test('archive suggests next roadmap mission and advances current', () => {
  const dir = spaceRoot();
  try {
    writeMission(dir, 'M07ARCH1', 'shipped');
    writeMission(dir, 'M07ARCH2', 'planned');
    writeCurrent(dir, 'M07ARCH1');

    const create = runCLI(dir, 'map', 'new', 'Archive Next');
    assertNotStub(create, 'map new');
    assert.equal(create.code, 0, `map new exit=${create.code}\n${combined(create)}`);

    let rid = 'archive-next';
    if (!existsSync(path.join(dir, '.space', 'roadmaps', `${rid}.json`))) {
      rid = findRoadmapID(dir);
    }

    assert.equal(
      runCLI(dir, 'map', 'add', rid, 'M07ARCH1', '--desc', 'shipped one').code,
      0,
    );
    assert.equal(
      runCLI(dir, 'map', 'add', rid, 'M07ARCH2', '--desc', 'planned next').code,
      0,
    );
    const use = runCLI(dir, 'map', 'use', rid);
    assert.equal(use.code, 0, `map use exit=${use.code}\n${combined(use)}`);

    const arch = runCLI(dir, 'archive', 'M07ARCH1');
    assertNotStub(arch, 'archive');
    assert.equal(arch.code, 0, `archive exit=${arch.code}\n${combined(arch)}`);
    assert.match(arch.stdout, /Archived mission M07ARCH1/);

    const wantHint = `Next mission on roadmap ${rid}: M07ARCH2: planned next (state=planned)`;
    assert.ok(
      arch.stdout.includes(wantHint),
      `want hint ${JSON.stringify(wantHint)}\nstdout=${arch.stdout}`,
    );
    assert.match(
      arch.stdout,
      /Suggested: new session → \/sc-discuss M07ARCH2 \(then \/sc-run\)/,
    );
    assert.equal(readCurrent(dir), 'M07ARCH2');
  } finally {
    cleanup(dir);
  }
});

test('archive last mission omits next hint and clears current', () => {
  const dir = spaceRoot();
  try {
    writeMission(dir, 'M07LAST1', 'shipped');
    writeCurrent(dir, 'M07LAST1');

    const create = runCLI(dir, 'map', 'new', 'Archive Last');
    assertNotStub(create, 'map new');
    assert.equal(create.code, 0, `map new exit=${create.code}\n${combined(create)}`);

    let rid = 'archive-last';
    if (!existsSync(path.join(dir, '.space', 'roadmaps', `${rid}.json`))) {
      rid = findRoadmapID(dir);
    }

    assert.equal(
      runCLI(dir, 'map', 'add', rid, 'M07LAST1', '--desc', 'only one').code,
      0,
    );

    const arch = runCLI(dir, 'archive', 'M07LAST1');
    assertNotStub(arch, 'archive');
    assert.equal(arch.code, 0, `archive exit=${arch.code}\n${combined(arch)}`);
    assert.doesNotMatch(arch.stdout, /Next mission on roadmap/);
    assert.equal(readCurrent(dir), '', `current=${JSON.stringify(readCurrent(dir))} want cleared`);
  } finally {
    cleanup(dir);
  }
});
