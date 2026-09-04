import { spawnSync } from 'node:child_process';
import { existsSync, mkdirSync, readFileSync, writeFileSync } from 'node:fs';
import path from 'node:path';

const SPACE_DIRS = ['missions', 'archive', 'roadmaps'];

/**
 * Minimal first-create `.gitignore` when the file is missing.
 * Includes common secret/build ignores plus Spacecraft companion dirs.
 */
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

/** Companion dirs to merge into an existing `.gitignore` (never full overwrite). */
const COMPANION_IGNORE_LINES = ['.codegraph/', '.impeccable/'];

/** True when a non-comment gitignore line ignores `.space` / `.space/` / `.space/*`. */
export function isSpaceIgnoreLine(line) {
  const trimmed = line.trim();
  if (!trimmed || trimmed.startsWith('#')) return false;
  return trimmed === '.space' || trimmed === '.space/' || trimmed === '.space/*';
}

export function hasSpaceIgnored(gitignoreText) {
  return gitignoreText.split(/\r?\n/).some(isSpaceIgnoreLine);
}

function hasExactIgnoreLine(gitignoreText, entry) {
  const bare = entry.endsWith('/') ? entry.slice(0, -1) : entry;
  return gitignoreText.split(/\r?\n/).some((line) => {
    const trimmed = line.trim();
    if (!trimmed || trimmed.startsWith('#')) return false;
    return trimmed === entry || trimmed === bare;
  });
}

/** Append missing companion ignore lines without replacing the file. */
function ensureCompanionIgnores(projectRoot) {
  const gitignorePath = path.join(projectRoot, '.gitignore');
  if (!existsSync(gitignorePath)) return;
  let current = readFileSync(gitignorePath, 'utf8');
  let changed = false;
  for (const entry of COMPANION_IGNORE_LINES) {
    if (hasExactIgnoreLine(current, entry)) continue;
    current = current.endsWith('\n') ? `${current}${entry}\n` : `${current}\n${entry}\n`;
    changed = true;
  }
  if (changed) writeFileSync(gitignorePath, current);
}

function ensureSpaceDirs(projectRoot) {
  const spaceDir = path.join(projectRoot, '.space');
  for (const name of SPACE_DIRS) {
    mkdirSync(path.join(spaceDir, name), { recursive: true });
  }
}

function isGitRepo(projectRoot) {
  const result = spawnSync('git', ['rev-parse', '--git-dir'], {
    cwd: projectRoot,
    encoding: 'utf8',
  });
  if (result.error?.code === 'ENOENT') {
    throw new Error('git binary not found; install git and retry');
  }
  return result.status === 0;
}

function ensureGitRepo(projectRoot) {
  if (isGitRepo(projectRoot)) return;
  const result = spawnSync('git', ['init'], {
    cwd: projectRoot,
    encoding: 'utf8',
  });
  if (result.error?.code === 'ENOENT') {
    throw new Error('git binary not found; install git and retry');
  }
  if (result.status !== 0) {
    const detail = (result.stderr || result.stdout || '').trim();
    throw new Error(detail ? `git init failed: ${detail}` : 'git init failed');
  }
}

/**
 * Write starter `.gitignore` only when missing.
 * If present, merge Spacecraft lines via `ensureSpaceIgnored` + companions (no clobber).
 */
function writeStarterGitignore(projectRoot) {
  const gitignorePath = path.join(projectRoot, '.gitignore');
  if (!existsSync(gitignorePath)) {
    writeFileSync(gitignorePath, STARTER_GITIGNORE);
    return;
  }
  ensureSpaceIgnored(projectRoot);
  ensureCompanionIgnores(projectRoot);
}

/** Soft-run `codegraph init` when the index DB is missing; never throw. */
function ensureCodegraph(projectRoot) {
  if (existsSync(path.join(projectRoot, '.codegraph', 'codegraph.db'))) return;
  const result = spawnSync('codegraph', ['init'], {
    cwd: projectRoot,
    encoding: 'utf8',
  });
  if (result.error?.code === 'ENOENT') {
    console.error('warning: codegraph not found; skip index init');
    return;
  }
  if (result.status !== 0) {
    const detail = (result.stderr || result.stdout || result.error?.message || '').trim();
    console.error(
      detail
        ? `warning: codegraph init failed: ${detail}`
        : 'warning: codegraph init failed',
    );
  }
}

/**
 * First `.space` create: scaffold dirs, git init if needed, write starter `.gitignore`
 * only when missing (else merge Spacecraft ignores), soft-init codegraph when missing.
 * Does not seed product `docs/`.
 */
export function ensureProjectReady(projectRoot) {
  ensureSpaceDirs(projectRoot);
  ensureGitRepo(projectRoot);
  writeStarterGitignore(projectRoot);
  ensureCodegraph(projectRoot);
}

/**
 * When `.space` already exists: ensure `.space/` is listed in `.gitignore`
 * without replacing the whole file. Does not seed product `docs/`.
 */
export function ensureSpaceIgnored(projectRoot) {
  const gitignorePath = path.join(projectRoot, '.gitignore');
  if (!existsSync(gitignorePath)) {
    writeFileSync(gitignorePath, '.space/\n');
  } else {
    const current = readFileSync(gitignorePath, 'utf8');
    if (!hasSpaceIgnored(current)) {
      const suffix = current.endsWith('\n') ? '.space/\n' : '\n.space/\n';
      writeFileSync(gitignorePath, `${current}${suffix}`);
    }
  }
}
