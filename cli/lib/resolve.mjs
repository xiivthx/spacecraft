import { spawnSync } from 'node:child_process';
import { existsSync, readFileSync, statSync } from 'node:fs';
import path from 'node:path';

export function missionDir(spaceDir, id) {
  return path.join(spaceDir, 'missions', id);
}

export function missionExists(spaceDir, id) {
  try {
    return statSync(missionDir(spaceDir, id)).isDirectory();
  } catch {
    return false;
  }
}

export function currentFile(spaceDir) {
  return path.join(spaceDir, 'current');
}

export function readCurrent(spaceDir) {
  try {
    return readFileSync(currentFile(spaceDir), 'utf8').trim();
  } catch {
    return '';
  }
}

export function normalizeID(sel) {
  return String(sel ?? '')
    .trim()
    .toUpperCase();
}

/** Mission ID from feat/<id>/… branch, or empty. */
export function resolveMission(cwd) {
  const result = spawnSync('git', ['branch', '--show-current'], {
    cwd,
    encoding: 'utf8',
  });
  const branch = (result.stdout ?? '').trim();
  if (!branch.startsWith('feat/')) return '';
  const parts = branch.split('/');
  if (parts.length >= 2 && parts[1].startsWith('M')) {
    return parts[1];
  }
  return '';
}

/** Prefer .space/current when the mission exists; else branch mid. */
export function resolveActive(spaceDir, mid) {
  const cur = readCurrent(spaceDir);
  if (cur && missionExists(spaceDir, cur)) return cur;
  if (mid && missionExists(spaceDir, mid)) return mid;
  return '';
}

export function gitShowCurrentBranch(cwd) {
  const result = spawnSync('git', ['branch', '--show-current'], {
    cwd,
    encoding: 'utf8',
  });
  if (result.status !== 0) return '';
  return (result.stdout ?? '').trim();
}

export function spaceDirFromCwd(cwd) {
  return path.join(cwd, '.space');
}

export function pathExists(p) {
  return existsSync(p);
}
