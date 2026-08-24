import { spawnSync } from 'node:child_process';
import { copyFileSync, existsSync, mkdirSync, readFileSync, writeFileSync } from 'node:fs';
import path from 'node:path';
import { fileURLToPath } from 'node:url';

const REPO_ROOT = path.resolve(path.dirname(fileURLToPath(import.meta.url)), '../..');
const GITIGNORE_TEMPLATE = path.join(REPO_ROOT, 'templates', 'gitignore');
const DOCS_TEMPLATE_ROOT = path.join(REPO_ROOT, 'templates', 'docs');

const SPACE_DIRS = ['missions', 'archive', 'roadmaps'];

/** Relpaths under `templates/docs/` → consumer `docs/` (seed only; no on-demand dirs). */
const DOCS_SEED_REL_PATHS = [
  'README.md',
  'conventions/README.md',
  'conventions/naming.md',
];

/** True when a non-comment gitignore line ignores `.space` / `.space/`. */
export function isSpaceIgnoreLine(line) {
  const trimmed = line.trim();
  if (!trimmed || trimmed.startsWith('#')) return false;
  return trimmed === '.space' || trimmed === '.space/';
}

export function hasSpaceIgnored(gitignoreText) {
  return gitignoreText.split(/\r?\n/).some(isSpaceIgnoreLine);
}

function gitignoreTemplatePath() {
  if (existsSync(GITIGNORE_TEMPLATE)) return GITIGNORE_TEMPLATE;
  throw new Error(`missing gitignore template at ${GITIGNORE_TEMPLATE}`);
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

function writeGitignoreFromTemplate(projectRoot) {
  const templatePath = gitignoreTemplatePath();
  copyFileSync(templatePath, path.join(projectRoot, '.gitignore'));
}

/**
 * Copy missing product-docs seed files from `templates/docs/`.
 * Never overwrites existing targets; does not create on-demand dirs.
 */
function ensureDocsSeeded(projectRoot) {
  for (const rel of DOCS_SEED_REL_PATHS) {
    const dest = path.join(projectRoot, 'docs', rel);
    if (existsSync(dest)) continue;
    const src = path.join(DOCS_TEMPLATE_ROOT, rel);
    if (!existsSync(src)) {
      throw new Error(`missing docs seed template at ${src}`);
    }
    mkdirSync(path.dirname(dest), { recursive: true });
    copyFileSync(src, dest);
  }
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
 * First `.space` create: scaffold dirs, git init if needed, overwrite `.gitignore`
 * from `templates/gitignore`, soft-init codegraph when missing, seed missing docs.
 */
export function ensureProjectReady(projectRoot) {
  ensureSpaceDirs(projectRoot);
  ensureGitRepo(projectRoot);
  writeGitignoreFromTemplate(projectRoot);
  ensureCodegraph(projectRoot);
  ensureDocsSeeded(projectRoot);
}

/**
 * When `.space` already exists: ensure `.space/` is listed in `.gitignore`
 * without replacing the whole file; also seed missing docs targets.
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
  ensureDocsSeeded(projectRoot);
}
