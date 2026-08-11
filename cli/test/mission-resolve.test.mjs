/**
 * Node CLI tests for mission graph: init, new, missions, use, current, resolve,
 * status, flow, and bind-branch.
 */
import assert from 'node:assert/strict';
import { spawnSync } from 'node:child_process';
import {
  existsSync,
  mkdirSync,
  mkdtempSync,
  readFileSync,
  readdirSync,
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
  const dir = mkdtempSync(path.join(os.tmpdir(), 'spacecraft-mission-'));
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

function runGit(dir, args) {
  const globalCfg = path.join(dir, '.gitconfig');
  const xdgConfig = path.join(dir, '.config');
  mkdirSync(xdgConfig, { recursive: true });
  if (!existsSync(globalCfg)) writeFileSync(globalCfg, '');
  const result = spawnSync('git', args, {
    cwd: dir,
    encoding: 'utf8',
    env: {
      ...process.env,
      HOME: dir,
      GIT_CONFIG_GLOBAL: globalCfg,
      GIT_CONFIG_NOSYSTEM: '1',
      GIT_TERMINAL_PROMPT: '0',
      XDG_CONFIG_HOME: xdgConfig,
    },
  });
  if (result.status !== 0) {
    throw new Error(`git ${args.join(' ')} failed: ${result.stderr || result.stdout}`);
  }
}

function initGitRepo(dir, branch) {
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

function currentOverridesBranchDir() {
  const dir = spaceRoot();
  const currentID = 'M07CURA1';
  const branchID = 'M07CURB1';
  writeMission(dir, currentID);
  writeMission(dir, branchID);
  writeCurrent(dir, currentID);
  initGitRepo(dir, `feat/${branchID}/other-mission`);
  return { dir, currentID, branchID };
}

function cleanup(dir) {
  rmSync(dir, { recursive: true, force: true });
}

test('resolve from feat/<id>/… when .space/current unset', () => {
  const dir = spaceRoot();
  const id = 'M07RES01';
  try {
    writeMission(dir, id);
    initGitRepo(dir, `feat/${id}/restore-cli`);

    const res = runCLI(dir, 'resolve');
    assert.equal(res.code, 0, `resolve exit=${res.code}\n${combined(res)}`);
    assert.match(combined(res), new RegExp(`Mission: ${id}`));
  } finally {
    cleanup(dir);
  }
});

test('resolve: .space/current overrides feat branch', () => {
  const { dir, currentID, branchID } = currentOverridesBranchDir();
  try {
    const res = runCLI(dir, 'resolve');
    assert.equal(res.code, 0, `resolve exit=${res.code}\n${combined(res)}`);
    const out = combined(res);
    assert.match(out, new RegExp(`Mission: ${currentID}`));
    assert.doesNotMatch(out, new RegExp(`Mission: ${branchID}`));
  } finally {
    cleanup(dir);
  }
});

test('resolve explicit selector wins over current and branch', () => {
  const { dir, currentID, branchID } = currentOverridesBranchDir();
  const explicitID = 'M07EXPC1';
  try {
    writeMission(dir, explicitID);
    const res = runCLI(dir, 'resolve', explicitID);
    assert.equal(res.code, 0, `resolve exit=${res.code}\n${combined(res)}`);
    const out = combined(res);
    assert.match(out, new RegExp(`Mission: ${explicitID}`));
    assert.doesNotMatch(out, new RegExp(`Mission: ${currentID}`));
    assert.doesNotMatch(out, new RegExp(`Mission: ${branchID}`));
  } finally {
    cleanup(dir);
  }
});

test('resolve unknown selector fails', () => {
  const dir = spaceRoot();
  try {
    writeMission(dir, 'M07RES03');
    initGitRepo(dir, 'main');
    const res = runCLI(dir, 'resolve', 'M07DOESNOTEXIST');
    assert.notEqual(res.code, 0, 'resolve bad selector must fail');
    const out = combined(res);
    assert.doesNotMatch(out, /not implemented/i, 'must be real resolve, not stub');
    assert.match(out, /no mission matches/i);
  } finally {
    cleanup(dir);
  }
});

test('resolve: missing current falls through to branch', () => {
  const dir = spaceRoot();
  const branchID = 'M07CURB2';
  try {
    writeMission(dir, branchID);
    writeCurrent(dir, 'M07MISSING');
    initGitRepo(dir, `feat/${branchID}/other-mission`);

    const res = runCLI(dir, 'resolve');
    assert.equal(res.code, 0, `resolve exit=${res.code}\n${combined(res)}`);
    const out = combined(res);
    assert.match(out, new RegExp(`Mission: ${branchID}`));
    assert.doesNotMatch(out, /Mission: M07MISSING/);
  } finally {
    cleanup(dir);
  }
});

test('status and flow: current overrides branch', () => {
  const { dir, currentID, branchID } = currentOverridesBranchDir();
  try {
    for (const cmd of ['status', 'flow']) {
      const res = runCLI(dir, cmd);
      assert.equal(res.code, 0, `${cmd} exit=${res.code}\n${combined(res)}`);
      const out = combined(res);
      assert.match(out, new RegExp(`Mission: ${currentID}`), `${cmd} should prefer current`);
      assert.doesNotMatch(out, new RegExp(`Mission: ${branchID}`), `${cmd} must not prefer branch`);
    }
  } finally {
    cleanup(dir);
  }
});

test('status prints Pickup line when mission.json has pickup.next', () => {
  const dir = spaceRoot();
  const id = 'M08PICK1';
  const next = 'continue T2: evidence truncate';
  try {
    writeMission(dir, id);
    writeCurrent(dir, id);
    initGitRepo(dir, 'main');

    const missionPath = path.join(dir, '.space', 'missions', id, 'mission.json');
    const mission = JSON.parse(readFileSync(missionPath, 'utf8'));
    mission.pickup = {
      phase: 'run',
      next,
      updatedAt: '2026-08-10T00:00:00Z',
    };
    writeFileSync(missionPath, `${JSON.stringify(mission, null, 2)}\n`);

    const res = runCLI(dir, 'status');
    assert.equal(res.code, 0, `status exit=${res.code}\n${combined(res)}`);
    const out = combined(res);
    assert.match(out, /^Pickup: continue T2: evidence truncate$/m);
  } finally {
    cleanup(dir);
  }
});

test('status omits Pickup when mission.json has no pickup', () => {
  const dir = spaceRoot();
  const id = 'M08PICK0';
  try {
    writeMission(dir, id);
    writeCurrent(dir, id);
    initGitRepo(dir, 'main');

    const res = runCLI(dir, 'status');
    assert.equal(res.code, 0, `status exit=${res.code}\n${combined(res)}`);
    const out = combined(res);
    assert.match(out, new RegExp(`^Mission: ${id}$`, 'm'));
    assert.match(out, /^Title: Test Mission$/m);
    assert.match(out, /^State: active$/m);
    assert.match(out, /^Evidence: \d+$/m);
    assert.doesNotMatch(out, /^Pickup:/m);
  } finally {
    cleanup(dir);
  }
});

test('init creates .space/missions and .space/roadmaps', () => {
  const dir = mkdtempSync(path.join(os.tmpdir(), 'spacecraft-init-'));
  try {
    initGitRepo(dir, 'main');
    const res = runCLI(dir, 'init');
    assert.equal(res.code, 0, `init exit=${res.code}\n${combined(res)}`);
    assert.match(combined(res), /Spacecraft initialized at \.space\//);
    assert.ok(existsSync(path.join(dir, '.space', 'missions')));
    assert.ok(existsSync(path.join(dir, '.space', 'roadmaps')));
  } finally {
    cleanup(dir);
  }
});

test('new creates mission dir + artifacts and sets current', () => {
  const dir = spaceRoot();
  try {
    const res = runCLI(dir, 'new', 'Real Title');
    assert.equal(res.code, 0, `new exit=${res.code}\n${combined(res)}`);
    assert.match(combined(res), /Created mission M[0-9A-Z]+/);

    const missionsDir = path.join(dir, '.space', 'missions');
    const entries = readdirSync(missionsDir).filter((name) =>
      existsSync(path.join(missionsDir, name, 'mission.json')),
    );
    assert.equal(entries.length, 1, `new mission count=${entries.length} want 1`);

    const id = entries[0];
    const missionDir = path.join(missionsDir, id);
    for (const file of ['mission.json', 'spec.md', 'plan.json', 'evidence.jsonl']) {
      assert.ok(existsSync(path.join(missionDir, file)), `new must create ${file}`);
    }

    const mission = JSON.parse(readFileSync(path.join(missionDir, 'mission.json'), 'utf8'));
    assert.equal(mission.id, id);
    assert.equal(mission.title, 'Real Title');
    assert.equal(mission.state, 'active');

    const current = readFileSync(path.join(dir, '.space', 'current'), 'utf8').trim();
    assert.equal(current, id);
  } finally {
    cleanup(dir);
  }
});

test('missions lists missions; use/current round-trip .space/current', () => {
  const dir = spaceRoot();
  const id = 'M07USE01';
  try {
    writeMission(dir, id);
    initGitRepo(dir, 'main');

    const emptyish = runCLI(dir, 'missions');
    assert.equal(emptyish.code, 0, `missions exit=${emptyish.code}\n${combined(emptyish)}`);
    assert.match(combined(emptyish), new RegExp(`${id}`));
    assert.match(combined(emptyish), /Test Mission/);

    const useRes = runCLI(dir, 'use', id);
    assert.equal(useRes.code, 0, `use exit=${useRes.code}\n${combined(useRes)}`);
    assert.match(combined(useRes), new RegExp(`Selected mission ${id}`));
    assert.equal(readFileSync(path.join(dir, '.space', 'current'), 'utf8').trim(), id);

    const curRes = runCLI(dir, 'current');
    assert.equal(curRes.code, 0, `current exit=${curRes.code}\n${combined(curRes)}`);
    assert.equal(curRes.stdout.trim(), id);

    const listed = runCLI(dir, 'missions');
    assert.match(combined(listed), /\*/);
  } finally {
    cleanup(dir);
  }
});

test('bind-branch records current git branch on mission', () => {
  const dir = spaceRoot();
  const id = 'M07BIND1';
  const branch = `feat/${id}/bind-me`;
  try {
    writeMission(dir, id);
    initGitRepo(dir, branch);

    const res = runCLI(dir, 'bind-branch', id);
    assert.equal(res.code, 0, `bind-branch exit=${res.code}\n${combined(res)}`);
    assert.match(combined(res), new RegExp(`Bound branch ${branch} to mission ${id}`));

    const mission = JSON.parse(
      readFileSync(path.join(dir, '.space', 'missions', id, 'mission.json'), 'utf8'),
    );
    assert.ok(Array.isArray(mission.branches), 'mission.branches must be an array');
    assert.ok(mission.branches.includes(branch), `branches must include ${branch}`);
  } finally {
    cleanup(dir);
  }
});
