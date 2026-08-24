/**
 * Docs seed policy: ensureProjectReady / ensureSpaceIgnored seed shape,
 * stealth tokens, non-clobber, and .space/ gitignore (M9C0P519 S1–S4 / E1–E5).
 */
import assert from 'node:assert/strict';
import {
  existsSync,
  mkdirSync,
  mkdtempSync,
  readdirSync,
  readFileSync,
  rmSync,
  statSync,
  writeFileSync,
} from 'node:fs';
import os from 'node:os';
import path from 'node:path';
import test from 'node:test';
import { fileURLToPath } from 'node:url';

const repoRoot = path.resolve(path.dirname(fileURLToPath(import.meta.url)), '../..');
const templatesDocsRoot = path.join(repoRoot, 'templates', 'docs');

/** Frozen seed file set (design-contract + approved-scenarios S1). */
const SEED_REL_PATHS = [
  'docs/README.md',
  'docs/conventions/README.md',
  'docs/conventions/naming.md',
];

/** On-demand dirs that must remain absent after seed (S1/E1, E5). */
const ON_DEMAND_DIRS = [
  'docs/epics',
  'docs/specs',
  'docs/architecture',
  'docs/runbooks',
];

/**
 * Stealth forbidden tokens (design-contract E2 / approved-scenarios S2).
 * Case-insensitive; word boundaries for agent(s), mission(s), afk, prompt, llm.
 */
const FORBIDDEN_RE =
  /spacecraft|\bagents?\b|\bmissions?\b|\bafk\b|\/sc-|\bprompt\b|\bllm\b/gi;

const CUSTOM_README_MARKER = 'CUSTOM-DOCS-README';

function emptyProjectRoot() {
  return mkdtempSync(path.join(os.tmpdir(), 'spacecraft-docs-seed-'));
}

function cleanup(dir) {
  rmSync(dir, { recursive: true, force: true });
}

function listFilesRecursive(dir) {
  if (!existsSync(dir)) return [];
  const out = [];
  for (const name of readdirSync(dir)) {
    const full = path.join(dir, name);
    if (statSync(full).isDirectory()) {
      out.push(...listFilesRecursive(full));
    } else {
      out.push(full);
    }
  }
  return out;
}

function forbiddenMatches(text) {
  return [...text.matchAll(FORBIDDEN_RE)].map((m) => m[0]);
}

test('S1/E1: ensureProjectReady seeds docs/README.md and conventions stubs', async () => {
  const dir = emptyProjectRoot();
  try {
    const { ensureProjectReady } = await import('../lib/project-git.mjs');
    ensureProjectReady(dir);

    for (const rel of SEED_REL_PATHS) {
      const full = path.join(dir, rel);
      assert.ok(
        existsSync(full) && statSync(full).isFile(),
        `expected seed file ${rel} after ensureProjectReady`,
      );
      assert.ok(
        readFileSync(full, 'utf8').length > 0,
        `expected non-empty seed file ${rel}`,
      );
    }
  } finally {
    cleanup(dir);
  }
});

test('S1/E5: ensureProjectReady does not create on-demand docs directories', async () => {
  const dir = emptyProjectRoot();
  try {
    const { ensureProjectReady } = await import('../lib/project-git.mjs');
    ensureProjectReady(dir);

    // Seed must have run: otherwise "dirs absent" is vacuous (no write path).
    for (const rel of SEED_REL_PATHS) {
      assert.ok(
        existsSync(path.join(dir, rel)),
        `seed precondition failed: missing ${rel} (cannot judge on-demand absence)`,
      );
    }

    for (const rel of ON_DEMAND_DIRS) {
      assert.equal(
        existsSync(path.join(dir, rel)),
        false,
        `must not create on-demand dir ${rel}`,
      );
    }
  } finally {
    cleanup(dir);
  }
});

test('S2/E2: templates/docs has zero forbidden stealth tokens', () => {
  assert.ok(
    existsSync(templatesDocsRoot),
    'templates/docs must exist (T1 seed templates SoT)',
  );

  const files = listFilesRecursive(templatesDocsRoot);
  assert.ok(files.length > 0, 'templates/docs must contain at least one file');

  const hits = [];
  for (const file of files) {
    const text = readFileSync(file, 'utf8');
    for (const token of forbiddenMatches(text)) {
      hits.push(`${path.relative(repoRoot, file)}:${token}`);
    }
  }
  assert.deepEqual(
    hits,
    [],
    `forbidden stealth tokens in templates/docs: ${hits.join(', ')}`,
  );
});

test('S2/E2: seeded docs copies have zero forbidden stealth tokens', async () => {
  const dir = emptyProjectRoot();
  try {
    const { ensureProjectReady } = await import('../lib/project-git.mjs');
    ensureProjectReady(dir);

    const hits = [];
    for (const rel of SEED_REL_PATHS) {
      const full = path.join(dir, rel);
      assert.ok(existsSync(full), `seeded file missing: ${rel}`);
      for (const token of forbiddenMatches(readFileSync(full, 'utf8'))) {
        hits.push(`${rel}:${token}`);
      }
    }
    assert.deepEqual(
      hits,
      [],
      `forbidden stealth tokens in seeded docs: ${hits.join(', ')}`,
    );
  } finally {
    cleanup(dir);
  }
});

test('S3/E3: ensureProjectReady does not overwrite pre-existing docs/README.md', async () => {
  const dir = emptyProjectRoot();
  try {
    mkdirSync(path.join(dir, 'docs'), { recursive: true });
    const before = `# Local\n${CUSTOM_README_MARKER}\n`;
    writeFileSync(path.join(dir, 'docs', 'README.md'), before);

    const { ensureProjectReady } = await import('../lib/project-git.mjs');
    ensureProjectReady(dir);

    const after = readFileSync(path.join(dir, 'docs', 'README.md'), 'utf8');
    assert.equal(
      after,
      before,
      'pre-existing docs/README.md bytes must be unchanged (CUSTOM-DOCS-README)',
    );
  } finally {
    cleanup(dir);
  }
});

test('S3/E3 companion: with pre-existing README, ensure still seeds missing conventions stubs', async () => {
  const dir = emptyProjectRoot();
  try {
    mkdirSync(path.join(dir, 'docs'), { recursive: true });
    writeFileSync(
      path.join(dir, 'docs', 'README.md'),
      `# Local\n${CUSTOM_README_MARKER}\n`,
    );

    const { ensureProjectReady } = await import('../lib/project-git.mjs');
    ensureProjectReady(dir);

    assert.ok(
      existsSync(path.join(dir, 'docs', 'conventions', 'README.md')),
      'missing docs/conventions/README.md must still be seeded',
    );
    assert.ok(
      existsSync(path.join(dir, 'docs', 'conventions', 'naming.md')),
      'missing docs/conventions/naming.md must still be seeded',
    );
  } finally {
    cleanup(dir);
  }
});

test('S4/E4: ensureProjectReady leaves .space/ ignored in .gitignore', async () => {
  const dir = emptyProjectRoot();
  try {
    const { ensureProjectReady } = await import('../lib/project-git.mjs');
    ensureProjectReady(dir);

    const gitignorePath = path.join(dir, '.gitignore');
    assert.ok(existsSync(gitignorePath), 'ensure must write .gitignore');
    assert.match(
      readFileSync(gitignorePath, 'utf8'),
      /^\.space\/$/m,
      '.gitignore must include a .space/ ignore line after ensure',
    );
  } finally {
    cleanup(dir);
  }
});

test('ensureSpaceIgnored: seeds missing docs/README.md and conventions stubs', async () => {
  const dir = emptyProjectRoot();
  try {
    for (const name of ['missions', 'archive', 'roadmaps']) {
      mkdirSync(path.join(dir, '.space', name), { recursive: true });
    }
    writeFileSync(path.join(dir, '.gitignore'), 'node_modules/\n');

    const { ensureSpaceIgnored } = await import('../lib/project-git.mjs');
    ensureSpaceIgnored(dir);

    for (const rel of SEED_REL_PATHS) {
      const full = path.join(dir, rel);
      assert.ok(
        existsSync(full) && statSync(full).isFile(),
        `ensureSpaceIgnored must seed missing ${rel}`,
      );
    }
  } finally {
    cleanup(dir);
  }
});
