/**
 * Node CLI / unit tests for project-git ensure helpers:
 * first .space create (git init + gitignore overwrite) vs later append-only.
 */
import assert from 'node:assert/strict';
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
const gitignoreTemplatePath = path.join(repoRoot, 'templates', 'gitignore');

function emptyProjectRoot() {
  return mkdtempSync(path.join(os.tmpdir(), 'spacecraft-project-git-'));
}

function cleanup(dir) {
  rmSync(dir, { recursive: true, force: true });
}

function isolatedGitEnv(dir) {
  const globalCfg = path.join(dir, '.gitconfig');
  const xdgConfig = path.join(dir, '.config');
  mkdirSync(xdgConfig, { recursive: true });
  if (!existsSync(globalCfg)) writeFileSync(globalCfg, '');
  return {
    ...process.env,
    HOME: dir,
    GIT_CONFIG_GLOBAL: globalCfg,
    GIT_CONFIG_NOSYSTEM: '1',
    GIT_TERMINAL_PROMPT: '0',
    XDG_CONFIG_HOME: xdgConfig,
  };
}

function runGit(dir, args) {
  const result = spawnSync('git', args, {
    cwd: dir,
    encoding: 'utf8',
    env: isolatedGitEnv(dir),
  });
  if (result.status !== 0) {
    throw new Error(`git ${args.join(' ')} failed: ${result.stderr || result.stdout}`);
  }
  return result;
}

function isGitRepo(dir) {
  const result = spawnSync('git', ['rev-parse', '--git-dir'], {
    cwd: dir,
    encoding: 'utf8',
    env: isolatedGitEnv(dir),
  });
  return result.status === 0;
}

function runCLI(dir, ...args) {
  const result = spawnSync(process.execPath, [entryPath, ...args], {
    cwd: dir,
    encoding: 'utf8',
    env: isolatedGitEnv(dir),
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

function assertSpaceScaffold(dir) {
  for (const name of ['missions', 'archive', 'roadmaps']) {
    assert.ok(
      existsSync(path.join(dir, '.space', name)),
      `expected .space/${name} directory`,
    );
  }
}

test('ensureProjectReady: empty dir gets git, .space scaffold, and template .gitignore', async () => {
  const dir = emptyProjectRoot();
  try {
    assert.equal(isGitRepo(dir), false, 'fixture must start without .git');
    assert.equal(existsSync(path.join(dir, '.space')), false, 'fixture must start without .space');

    const { ensureProjectReady } = await import('../lib/project-git.mjs');
    ensureProjectReady(dir);

    assert.ok(isGitRepo(dir), 'ensureProjectReady must git init when .git is missing');
    assertSpaceScaffold(dir);

    assert.ok(
      existsSync(gitignoreTemplatePath),
      'templates/gitignore must exist (consumer SoT)',
    );
    const template = readFileSync(gitignoreTemplatePath, 'utf8');
    assert.match(template, /^\.space\/$/m, 'templates/gitignore must include a .space/ entry');

    const gitignorePath = path.join(dir, '.gitignore');
    assert.ok(existsSync(gitignorePath), 'ensureProjectReady must write .gitignore');
    const written = readFileSync(gitignorePath, 'utf8');
    assert.equal(
      written,
      template,
      'first .space create must overwrite .gitignore from templates/gitignore',
    );
    assert.match(written, /^\.space\/$/m);
  } finally {
    cleanup(dir);
  }
});

test('ensureSpaceIgnored: appends .space/ without overwriting custom .gitignore', async () => {
  const dir = emptyProjectRoot();
  try {
    for (const name of ['missions', 'archive', 'roadmaps']) {
      mkdirSync(path.join(dir, '.space', name), { recursive: true });
    }
    const marker = '# CUSTOM-MARKER-DO-NOT-DROP';
    writeFileSync(path.join(dir, '.gitignore'), `${marker}\nnode_modules/\n`);

    const { ensureSpaceIgnored } = await import('../lib/project-git.mjs');
    ensureSpaceIgnored(dir);

    const written = readFileSync(path.join(dir, '.gitignore'), 'utf8');
    assert.match(written, /CUSTOM-MARKER-DO-NOT-DROP/, 'must keep custom .gitignore content');
    assert.match(written, /^\.space\/$/m, 'must append .space/ when missing');

    if (existsSync(gitignoreTemplatePath)) {
      const template = readFileSync(gitignoreTemplatePath, 'utf8');
      assert.notEqual(
        written,
        template,
        'ensureSpaceIgnored must not replace whole .gitignore with the template',
      );
    }
  } finally {
    cleanup(dir);
  }
});

test('CLI init: no .space/.git scaffolds dirs and ensures git + .gitignore', () => {
  const dir = emptyProjectRoot();
  try {
    assert.equal(isGitRepo(dir), false);
    assert.equal(existsSync(path.join(dir, '.space')), false);

    const res = runCLI(dir, 'init');
    assert.equal(res.code, 0, `init exit=${res.code}\n${combined(res)}`);

    assert.ok(isGitRepo(dir), 'init path must create .git when missing');
    assertSpaceScaffold(dir);

    const gitignorePath = path.join(dir, '.gitignore');
    assert.ok(existsSync(gitignorePath), 'init ensure path must write .gitignore');
    assert.match(readFileSync(gitignorePath, 'utf8'), /^\.space\/$/m);
  } finally {
    cleanup(dir);
  }
});
