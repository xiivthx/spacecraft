import { spawnSync } from 'node:child_process';
import { copyFileSync, existsSync, mkdirSync, readFileSync, writeFileSync } from 'node:fs';
import path from 'node:path';
import { fileURLToPath } from 'node:url';

const REPO_ROOT = path.resolve(path.dirname(fileURLToPath(import.meta.url)), '../..');
const GITIGNORE_TEMPLATE = path.join(REPO_ROOT, 'templates', 'gitignore');

const SPACE_DIRS = ['missions', 'archive', 'roadmaps'];

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
 * First `.space` create: scaffold dirs, git init if needed, overwrite `.gitignore`
 * from `templates/gitignore`.
 */
export function ensureProjectReady(projectRoot) {
  ensureSpaceDirs(projectRoot);
  ensureGitRepo(projectRoot);
  writeGitignoreFromTemplate(projectRoot);
}

/**
 * When `.space` already exists: ensure `.space/` is listed in `.gitignore`
 * without replacing the whole file.
 */
export function ensureSpaceIgnored(projectRoot) {
  const gitignorePath = path.join(projectRoot, '.gitignore');
  if (!existsSync(gitignorePath)) {
    writeFileSync(gitignorePath, '.space/\n');
    return;
  }
  const current = readFileSync(gitignorePath, 'utf8');
  if (hasSpaceIgnored(current)) return;
  const suffix = current.endsWith('\n') ? '.space/\n' : '\n.space/\n';
  writeFileSync(gitignorePath, `${current}${suffix}`);
}
