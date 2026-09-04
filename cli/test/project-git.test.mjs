/**
 * Node CLI / unit tests for project-git ensure helpers:
 * first .space create (git init + starter .gitignore when missing) vs merge into existing.
 */
import assert from 'node:assert/strict';
import { spawnSync } from 'node:child_process';
import {
  chmodSync,
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
const CODEGRAPH_DB = path.join('.codegraph', 'codegraph.db');
const STARTER_GITIGNORE = `.env
.env.*
!.env.example
node_modules/
dist/
build/
.DS_Store
.space/
.codegraph/
.impeccable/
`;

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

/**
 * PATH stub: fake `codegraph` that logs argv to a marker file and optionally
 * creates `.codegraph/codegraph.db` under project path arg or cwd.
 */
function installFakeCodegraph({ exitCode = 0, createDb = true } = {}) {
  const binDir = mkdtempSync(path.join(os.tmpdir(), 'spacecraft-codegraph-bin-'));
  const markerPath = path.join(binDir, 'codegraph-argv.log');
  const scriptPath = path.join(binDir, 'codegraph');
  const markerQ = JSON.stringify(markerPath);
  const createDbBlock = createDb
    ? `
root="."
for arg in "$@"; do
  case "$arg" in
    -*) ;;
    *)
      if [ -d "$arg" ]; then root="$arg"; break; fi
      ;;
  esac
done
mkdir -p "$root/.codegraph"
: > "$root/.codegraph/codegraph.db"
`
    : '';
  writeFileSync(
    scriptPath,
    `#!/bin/sh
printf '%s\\n' "$*" >> ${markerQ}
${createDbBlock}
exit ${exitCode}
`,
  );
  chmodSync(scriptPath, 0o755);
  return { binDir, markerPath, scriptPath };
}

function pathWithoutCodegraphBin() {
  const parts = (process.env.PATH || '').split(path.delimiter).filter(Boolean);
  return parts
    .filter((dir) => !existsSync(path.join(dir, 'codegraph')))
    .join(path.delimiter);
}

function withPath(binDirOrPath, fn) {
  const prev = process.env.PATH;
  process.env.PATH = binDirOrPath;
  try {
    return fn();
  } finally {
    process.env.PATH = prev;
  }
}

function captureStderr(fn) {
  const lines = [];
  const origError = console.error;
  const origWarn = console.warn;
  console.error = (...args) => {
    lines.push(args.map(String).join(' '));
  };
  console.warn = (...args) => {
    lines.push(args.map(String).join(' '));
  };
  try {
    const value = fn();
    return { value, stderr: lines.join('\n') };
  } finally {
    console.error = origError;
    console.warn = origWarn;
  }
}

function readMarker(markerPath) {
  if (!existsSync(markerPath)) return '';
  return readFileSync(markerPath, 'utf8');
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

test('ensureProjectReady: empty dir gets git, .space scaffold, and starter .gitignore', async () => {
  const dir = emptyProjectRoot();
  try {
    assert.equal(isGitRepo(dir), false, 'fixture must start without .git');
    assert.equal(existsSync(path.join(dir, '.space')), false, 'fixture must start without .space');

    const { ensureProjectReady } = await import('../lib/project-git.mjs');
    ensureProjectReady(dir);

    assert.ok(isGitRepo(dir), 'ensureProjectReady must git init when .git is missing');
    assertSpaceScaffold(dir);

    const gitignorePath = path.join(dir, '.gitignore');
    assert.ok(existsSync(gitignorePath), 'ensureProjectReady must write .gitignore');
    const written = readFileSync(gitignorePath, 'utf8');
    assert.equal(
      written,
      STARTER_GITIGNORE,
      'first .space create must write starter .gitignore',
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
    assert.notEqual(
      written,
      STARTER_GITIGNORE,
      'ensureSpaceIgnored must not replace whole .gitignore with starter',
    );
  } finally {
    cleanup(dir);
  }
});

test('isSpaceIgnoreLine / hasSpaceIgnored: recognize .space/* as space ignore', async () => {
  const { isSpaceIgnoreLine, hasSpaceIgnored } = await import('../lib/project-git.mjs');

  assert.equal(
    isSpaceIgnoreLine('.space/*'),
    true,
    'isSpaceIgnoreLine must treat .space/* as a space-ignore line',
  );
  assert.equal(
    hasSpaceIgnored('.space/*\n'),
    true,
    'hasSpaceIgnored must be true when .gitignore has .space/*',
  );
  assert.equal(
    hasSpaceIgnored('.space/*\n!.space/polish-backlog.md\n'),
    true,
    'hasSpaceIgnored must stay true with .space/* plus un-ignore exception',
  );
});

test('ensureSpaceIgnored: does not append .space/ when .space/* already present', async () => {
  const dir = emptyProjectRoot();
  try {
    for (const name of ['missions', 'archive', 'roadmaps']) {
      mkdirSync(path.join(dir, '.space', name), { recursive: true });
    }
    const existing = '.space/*\n!.space/polish-backlog.md\n';
    writeFileSync(path.join(dir, '.gitignore'), existing);

    const { ensureSpaceIgnored } = await import('../lib/project-git.mjs');
    ensureSpaceIgnored(dir);

    const written = readFileSync(path.join(dir, '.gitignore'), 'utf8');
    assert.equal(
      written,
      existing,
      'ensureSpaceIgnored must leave .space/* (+ exception) alone; must not append .space/',
    );
    assert.equal(
      (written.match(/^\.space\/$/gm) ?? []).length,
      0,
      'must not introduce a bare .space/ line when .space/* already covers the tree',
    );
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

test('ensureProjectReady: soft-runs codegraph init when index missing', async () => {
  const dir = emptyProjectRoot();
  const fake = installFakeCodegraph();
  try {
    assert.equal(existsSync(path.join(dir, CODEGRAPH_DB)), false);

    const { ensureProjectReady } = await import('../lib/project-git.mjs');
    withPath(`${fake.binDir}${path.delimiter}${process.env.PATH || ''}`, () => {
      ensureProjectReady(dir);
    });

    assert.ok(
      existsSync(path.join(dir, CODEGRAPH_DB)),
      'ensureProjectReady must leave .codegraph/codegraph.db after soft init',
    );
    const marker = readMarker(fake.markerPath);
    assert.match(marker, /\binit\b/, 'fake codegraph must be invoked with init');
    assert.ok(isGitRepo(dir), 'git init outcome must remain');
    assert.ok(existsSync(path.join(dir, '.gitignore')), 'starter .gitignore must remain');
    assert.equal(
      readFileSync(path.join(dir, '.gitignore'), 'utf8'),
      STARTER_GITIGNORE,
      'first ensure must still write starter .gitignore',
    );
  } finally {
    cleanup(dir);
    cleanup(fake.binDir);
  }
});

test('ensureProjectReady: skips codegraph init when index already exists', async () => {
  const dir = emptyProjectRoot();
  const fake = installFakeCodegraph();
  try {
    mkdirSync(path.join(dir, '.codegraph'), { recursive: true });
    writeFileSync(path.join(dir, CODEGRAPH_DB), '');

    const { ensureProjectReady } = await import('../lib/project-git.mjs');
    withPath(`${fake.binDir}${path.delimiter}${process.env.PATH || ''}`, () => {
      ensureProjectReady(dir);
    });

    const marker = readMarker(fake.markerPath);
    assert.equal(
      marker.trim(),
      '',
      'ensureProjectReady must not invoke codegraph when index exists',
    );
    assert.ok(isGitRepo(dir));
    assert.ok(existsSync(path.join(dir, '.gitignore')));
  } finally {
    cleanup(dir);
    cleanup(fake.binDir);
  }
});

test('ensureProjectReady: missing codegraph binary still succeeds', async () => {
  const dir = emptyProjectRoot();
  try {
    const { ensureProjectReady } = await import('../lib/project-git.mjs');
    const stripped = pathWithoutCodegraphBin();
    const { stderr } = captureStderr(() => {
      assert.doesNotThrow(() => {
        withPath(stripped, () => {
          ensureProjectReady(dir);
        });
      }, 'missing codegraph must not throw from ensureProjectReady');
    });

    assert.ok(isGitRepo(dir), 'must still create .git');
    assert.ok(existsSync(path.join(dir, '.gitignore')), 'must still write .gitignore');
    assert.equal(
      readFileSync(path.join(dir, '.gitignore'), 'utf8'),
      STARTER_GITIGNORE,
    );
    assert.match(
      stderr,
      /codegraph/i,
      'missing codegraph binary must warn on stderr',
    );
  } finally {
    cleanup(dir);
  }
});

test('ensureProjectReady: failed codegraph init warns and still succeeds', async () => {
  const dir = emptyProjectRoot();
  const fake = installFakeCodegraph({ exitCode: 1, createDb: false });
  try {
    const { ensureProjectReady } = await import('../lib/project-git.mjs');
    const { stderr } = captureStderr(() => {
      withPath(`${fake.binDir}${path.delimiter}${process.env.PATH || ''}`, () => {
        ensureProjectReady(dir);
      });
    });

    assert.ok(isGitRepo(dir), 'failed init must not block git');
    assert.ok(existsSync(path.join(dir, '.gitignore')), 'failed init must not block gitignore');
    assert.match(
      stderr,
      /codegraph/i,
      'failed codegraph init must warn on stderr',
    );
  } finally {
    cleanup(dir);
    cleanup(fake.binDir);
  }
});
