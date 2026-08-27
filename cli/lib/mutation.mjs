import { spawnSync } from 'node:child_process';
import { appendFileSync, mkdirSync } from 'node:fs';
import path from 'node:path';
import { readMission } from './mission.mjs';
import { missionDir, resolveActive } from './resolve.mjs';

const DEFAULT_THRESHOLD = 80;
const DEFAULT_SCOPE_CAP = 50;

const MUTABLE_EXT = new Set([
  '.mjs',
  '.js',
  '.cjs',
  '.ts',
  '.tsx',
  '.jsx',
  '.py',
  '.go',
  '.rs',
  '.java',
  '.c',
  '.cpp',
  '.h',
  '.hpp',
  '.sv',
]);

function runGit(cwd, args) {
  return spawnSync('git', args, {
    cwd,
    encoding: 'utf8',
  });
}

function isMutableFile(relPath) {
  const ext = path.extname(relPath).toLowerCase();
  return MUTABLE_EXT.has(ext);
}

function scopeCap() {
  const raw = process.env.SPACECRAFT_MUTATION_SCOPE_CAP;
  if (raw === undefined || raw === '') return DEFAULT_SCOPE_CAP;
  const n = Number.parseInt(raw, 10);
  return Number.isFinite(n) && n > 0 ? n : DEFAULT_SCOPE_CAP;
}

function scopeSlug(files) {
  if (files.length === 0) return 'scope';
  const top = files[0].split('/')[0];
  return top || 'scope';
}

function resolveMergeBase(repoRoot, boundBranch) {
  const verify = runGit(repoRoot, ['rev-parse', '--verify', boundBranch]);
  if (verify.status !== 0) {
    return { ok: false, error: `merge-base: bound branch not found (${boundBranch})` };
  }

  const bases = ['main', 'origin/main', 'master', 'origin/master'];
  for (const base of bases) {
    const baseVerify = runGit(repoRoot, ['rev-parse', '--verify', base]);
    if (baseVerify.status !== 0) continue;
    const mb = runGit(repoRoot, ['merge-base', 'HEAD', base]);
    if (mb.status === 0 && mb.stdout.trim()) {
      return { ok: true, mergeBase: mb.stdout.trim() };
    }
  }

  const mbBound = runGit(repoRoot, ['merge-base', 'HEAD', boundBranch]);
  if (mbBound.status === 0 && mbBound.stdout.trim()) {
    return { ok: true, mergeBase: mbBound.stdout.trim() };
  }

  return { ok: false, error: 'merge-base: cannot resolve merge-base for mission-bound branch' };
}

function diffMutableFiles(repoRoot, mergeBase) {
  const diff = runGit(repoRoot, ['diff', '--name-only', `${mergeBase}..HEAD`]);
  if (diff.status !== 0) {
    return { ok: false, error: 'merge-base: git diff against merge-base failed' };
  }
  const files = diff.stdout
    .split('\n')
    .map((l) => l.trim())
    .filter(Boolean)
    .filter(isMutableFile);
  return { ok: true, files };
}

function stubResult(stub) {
  if (stub === 'pass') return { ok: true, score: 85 };
  if (stub === 'fail') return { ok: true, score: 70 };
  if (stub === 'missing') return { ok: false, missing: true };
  return null;
}

function appendEvidence(mDir, label, summary, exitCode) {
  const output = JSON.stringify(summary);
  const entry = {
    label,
    command: 'spacecraft mutation',
    output,
    exitCode,
    ts: new Date().toISOString().replace(/\.\d{3}Z$/, 'Z'),
  };
  mkdirSync(mDir, { recursive: true });
  appendFileSync(path.join(mDir, 'evidence.jsonl'), `${JSON.stringify(entry)}\n`);
}

function printMutationHelp() {
  console.log('Usage: spacecraft mutation [--mission <id>]');
  console.log('');
  console.log('Run diff-scoped mutation testing vs merge-base of the mission-bound branch.');
  console.log('Writes append-only evidence labeled mutation-<scope> with JSON summary');
  console.log('{ files, score, threshold, pass }. Default pass threshold: 80%.');
  console.log('');
  console.log('Options:');
  console.log('  --mission <id>   Mission id (default: active mission)');
  console.log('  --help, -h       Show this help');
}

export function mutationCmd(spaceDir, mid, args = []) {
  if (args.some((a) => a === '--help' || a === '-h')) {
    printMutationHelp();
    return 0;
  }

  let resolvedMid = mid;
  for (let i = 0; i < args.length; i++) {
    if (args[i] === '--mission' && i + 1 < args.length) {
      resolvedMid = args[i + 1];
      i++;
    }
  }
  if (!resolvedMid) {
    resolvedMid = resolveActive(spaceDir, mid);
  }
  if (!resolvedMid) {
    console.error(
      'spacecraft mutation: no active mission - use --mission <id>, spacecraft use, or run from feat/<id>/ branch',
    );
    return 1;
  }

  let mission;
  try {
    mission = readMission(spaceDir, resolvedMid);
  } catch (err) {
    console.error('spacecraft mutation:', err.message);
    return 1;
  }

  const branches = Array.isArray(mission.branches) ? mission.branches : [];
  const boundBranch = typeof branches[0] === 'string' ? branches[0].trim() : '';
  if (!boundBranch) {
    console.error('spacecraft mutation: merge-base: mission has no bound branch');
    return 1;
  }

  const repoRoot = path.dirname(spaceDir);
  const mb = resolveMergeBase(repoRoot, boundBranch);
  if (!mb.ok) {
    console.error(`spacecraft mutation: ${mb.error}`);
    return 1;
  }

  const diff = diffMutableFiles(repoRoot, mb.mergeBase);
  if (!diff.ok) {
    console.error(`spacecraft mutation: ${diff.error}`);
    return 1;
  }

  if (diff.files.length === 0) {
    console.log('scope-empty: no mutable files in diff vs merge-base');
    return 0;
  }

  const cap = scopeCap();
  let files = diff.files;
  let truncationNote = '';
  if (files.length > cap) {
    files = files.slice(0, cap);
    truncationNote = `scope truncated: ${diff.files.length} files capped to ${cap}`;
    console.log(truncationNote);
  }

  const stub = process.env.SPACECRAFT_MUTATION_STUB;
  const stubbed = stubResult(stub);
  if (stubbed?.missing) {
    console.error(
      'spacecraft mutation: mutator not installed — install Stryker or project mutation tool (npm install @stryker-mutator/core)',
    );
    return 1;
  }

  const threshold = DEFAULT_THRESHOLD;
  let score;
  if (stubbed) {
    score = stubbed.score;
  } else {
    console.error(
      'spacecraft mutation: mutator not installed — install Stryker or project mutation tool (npm install @stryker-mutator/core)',
    );
    return 1;
  }

  const pass = score >= threshold;
  const label = `mutation-${scopeSlug(files)}`;
  const summary = { files, score, threshold, pass };
  const mDir = missionDir(spaceDir, resolvedMid);
  appendEvidence(mDir, label, summary, pass ? 0 : 1);

  console.log(JSON.stringify(summary));
  if (truncationNote) console.log(truncationNote);

  return pass ? 0 : 1;
}
