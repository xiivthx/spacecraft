import { spawnSync } from 'node:child_process';
import { existsSync, readFileSync } from 'node:fs';
import path from 'node:path';
import { readMission } from './mission.mjs';
import { missionDir, resolveActive } from './resolve.mjs';

function isJSONNumber(v) {
  return typeof v === 'number' && Number.isFinite(v);
}

function runGit(cwd, args) {
  const result = spawnSync('git', args, {
    cwd,
    encoding: 'utf8',
  });
  return {
    ok: result.status === 0,
    stdout: (result.stdout ?? '').trim(),
    stderr: (result.stderr ?? '').trim(),
  };
}

export function closeoutEvidenceProblems(evidencePath) {
  let data;
  try {
    data = readFileSync(evidencePath, 'utf8');
  } catch {
    return ['missing evidence.jsonl'];
  }

  const required = ['label', 'command', 'output', 'ts'];
  let entries = 0;
  const problems = [];
  const lines = data.split('\n');

  for (let i = 0; i < lines.length; i++) {
    const line = lines[i];
    if (line.trim() === '') continue;

    let entry;
    try {
      entry = JSON.parse(line);
    } catch {
      problems.push(`evidence line ${i + 1} not valid JSON`);
      continue;
    }
    if (entry === null || typeof entry !== 'object' || Array.isArray(entry)) {
      problems.push(`evidence line ${i + 1} not valid JSON`);
      continue;
    }

    for (const field of required) {
      if (!(field in entry)) {
        problems.push(`evidence line ${i + 1} missing ${field}`);
      }
    }
    if (!isJSONNumber(entry.exitCode)) {
      problems.push(`evidence line ${i + 1} missing exitCode (number)`);
    }
    entries++;
  }

  if (entries === 0) {
    problems.push('no evidence captured');
  }
  return problems;
}

export function closeoutReviewProblems(reviewPath) {
  let data;
  try {
    data = readFileSync(reviewPath, 'utf8');
  } catch {
    return [];
  }

  let review;
  try {
    review = JSON.parse(data);
  } catch {
    return ['review.json invalid JSON'];
  }
  if (review === null || typeof review !== 'object' || Array.isArray(review)) {
    return ['review.json invalid JSON'];
  }

  const problems = [];
  const status = typeof review.status === 'string' ? review.status : '';
  if (status !== 'ready') {
    problems.push(`review.json status is ${JSON.stringify(status)}, expected "ready"`);
  }

  if (Array.isArray(review.findings)) {
    for (let i = 0; i < review.findings.length; i++) {
      const f =
        review.findings[i] !== null &&
        typeof review.findings[i] === 'object' &&
        !Array.isArray(review.findings[i])
          ? review.findings[i]
          : {};
      const sev = typeof f.severity === 'string' ? f.severity : '';
      const blocks = Boolean(f.blocksShip);
      problems.push(
        `review finding ${i + 1} present (severity=${JSON.stringify(sev)}, blocksShip=${blocks}); ship requires 0 findings (including warnings)`,
      );
    }
  }

  const rr = review.releaseReadiness;
  if (rr === null || typeof rr !== 'object' || Array.isArray(rr)) {
    problems.push('review.json releaseReadiness must be an object');
    return problems;
  }

  for (const key of ['changelog', 'specNote']) {
    const item = rr[key];
    if (item === null || typeof item !== 'object' || Array.isArray(item)) {
      problems.push(`releaseReadiness.${key} must be an object with status`);
      continue;
    }
    const st = typeof item.status === 'string' ? item.status : '';
    if (st !== 'ready') {
      problems.push(
        `releaseReadiness.${key} status is ${JSON.stringify(st)}, expected "ready"`,
      );
    }
  }
  return problems;
}

export function closeoutChangelogProblems(cwd) {
  const bases = ['main', 'origin/main'];
  let lastFailed = false;

  for (const base of bases) {
    const verify = runGit(cwd, ['rev-parse', '--verify', base]);
    if (!verify.ok) {
      lastFailed = true;
      continue;
    }
    const log = runGit(cwd, ['log', `${base}..HEAD`, '--', 'CHANGELOG.md']);
    if (!log.ok) {
      lastFailed = true;
      continue;
    }
    if (log.stdout === '') {
      return [`no commits touch CHANGELOG.md since ${base}`];
    }
    return [];
  }

  if (lastFailed) {
    return [
      'CHANGELOG check failed: neither main nor origin/main usable (or git unavailable)',
    ];
  }
  return ['no commits touch CHANGELOG.md'];
}

export function closeoutCmd(spaceDir, mid) {
  const id = resolveActive(spaceDir, mid);
  if (!id) {
    console.error('spacecraft closeout-check: no active mission');
    return 1;
  }

  let m;
  try {
    m = readMission(spaceDir, id);
  } catch (err) {
    console.error('spacecraft closeout-check:', err.message);
    return 1;
  }

  const dir = missionDir(spaceDir, id);
  const problems = [];

  for (const f of ['spec.md', 'plan.json', 'evidence.jsonl', 'review.json']) {
    if (!existsSync(path.join(dir, f))) {
      problems.push(`missing ${f}`);
    }
  }

  const state = typeof m.state === 'string' ? m.state : '';
  if (state !== 'ready' && state !== 'shipped') {
    problems.push(`state is ${state}, expected ready or shipped`);
  }

  try {
    const clarify = readFileSync(path.join(dir, 'clarify-status'), 'utf8').trim();
    if (clarify === 'open') {
      problems.push('clarify-status is open');
    }
  } catch {
    // absent clarify-status is fine
  }

  if (existsSync(path.join(dir, 'evidence.jsonl'))) {
    problems.push(...closeoutEvidenceProblems(path.join(dir, 'evidence.jsonl')));
  }
  if (existsSync(path.join(dir, 'review.json'))) {
    problems.push(...closeoutReviewProblems(path.join(dir, 'review.json')));
  }

  // SPACECRAFT_CLOSEOUT_SKIP_CHANGELOG=1 is for unit tests in temp dirs without
  // git history only. Production never sets this.
  if (process.env.SPACECRAFT_CLOSEOUT_SKIP_CHANGELOG !== '1') {
    problems.push(...closeoutChangelogProblems(path.dirname(spaceDir)));
  }

  if (problems.length > 0) {
    console.log(`Closeout blocked for ${id}:`);
    for (const p of problems) {
      console.log(`- ${p}`);
    }
    return 1;
  }
  console.log(`Closeout ready for ${id}.`);
  return 0;
}
