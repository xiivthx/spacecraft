/**
 * Ensure policy after templates/ removal: no product docs seed;
 * first create writes starter `.gitignore` with `.space/`.
 */
import assert from 'node:assert/strict';
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

/** Former seed targets — must stay absent after ensure. */
const FORMER_SEED_RELS = [
  'docs/README.md',
  'docs/conventions/README.md',
  'docs/conventions/naming.md',
];

function emptyProjectRoot() {
  return mkdtempSync(path.join(os.tmpdir(), 'spacecraft-docs-seed-'));
}

function cleanup(dir) {
  rmSync(dir, { recursive: true, force: true });
}

test('ensureProjectReady: does not seed product docs stubs', async () => {
  const dir = emptyProjectRoot();
  try {
    const { ensureProjectReady } = await import('../lib/project-git.mjs');
    ensureProjectReady(dir);

    for (const rel of FORMER_SEED_RELS) {
      assert.equal(
        existsSync(path.join(dir, rel)),
        false,
        `must not seed ${rel}`,
      );
    }
    assert.equal(
      existsSync(path.join(dir, 'docs')),
      false,
      'must not create docs/ tree',
    );
  } finally {
    cleanup(dir);
  }
});

test('ensureProjectReady: writes starter .gitignore with .space/', async () => {
  const dir = emptyProjectRoot();
  try {
    const { ensureProjectReady } = await import('../lib/project-git.mjs');
    ensureProjectReady(dir);

    const gitignorePath = path.join(dir, '.gitignore');
    assert.ok(existsSync(gitignorePath), 'ensure must write .gitignore');
    const text = readFileSync(gitignorePath, 'utf8');
    assert.match(text, /^\.space\/$/m, '.gitignore must include .space/');
    assert.equal(
      (text.match(/spacecraft/gi) ?? []).length,
      0,
      'starter .gitignore must not contain spacecraft',
    );
  } finally {
    cleanup(dir);
  }
});

test('ensureSpaceIgnored: does not seed product docs stubs', async () => {
  const dir = emptyProjectRoot();
  try {
    for (const name of ['missions', 'archive', 'roadmaps']) {
      mkdirSync(path.join(dir, '.space', name), { recursive: true });
    }
    writeFileSync(path.join(dir, '.gitignore'), 'node_modules/\n');

    const { ensureSpaceIgnored } = await import('../lib/project-git.mjs');
    ensureSpaceIgnored(dir);

    for (const rel of FORMER_SEED_RELS) {
      assert.equal(
        existsSync(path.join(dir, rel)),
        false,
        `ensureSpaceIgnored must not seed ${rel}`,
      );
    }
    assert.match(
      readFileSync(path.join(dir, '.gitignore'), 'utf8'),
      /^\.space\/$/m,
    );
  } finally {
    cleanup(dir);
  }
});

test('ensureProjectReady: does not overwrite pre-existing docs/README.md and does not add stubs', async () => {
  const dir = emptyProjectRoot();
  try {
    mkdirSync(path.join(dir, 'docs'), { recursive: true });
    const before = '# Local\nCUSTOM-DOCS-README\n';
    writeFileSync(path.join(dir, 'docs', 'README.md'), before);

    const { ensureProjectReady } = await import('../lib/project-git.mjs');
    ensureProjectReady(dir);

    assert.equal(
      readFileSync(path.join(dir, 'docs', 'README.md'), 'utf8'),
      before,
      'pre-existing docs/README.md must stay unchanged',
    );
    assert.equal(
      existsSync(path.join(dir, 'docs', 'conventions', 'README.md')),
      false,
      'must not seed conventions stubs',
    );
  } finally {
    cleanup(dir);
  }
});
