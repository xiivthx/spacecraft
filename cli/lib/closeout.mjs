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
    // Empty on this base: try next (e.g. post-merge on main, origin/main..HEAD still has release notes).
    if (log.stdout === '') {
      continue;
    }
    return [];
  }

  if (lastFailed) {
    return [
      'CHANGELOG check failed: neither main nor origin/main usable (or git unavailable)',
    ];
  }
  return ['no commits touch CHANGELOG.md since main or origin/main'];
}

const DISSENT_LABELS = new Set([
  'AGREE',
  'DISAGREE_EVIDENCE',
  'DISAGREE_CONCERN',
]);

function readOptionalText(filePath) {
  try {
    return readFileSync(filePath, 'utf8');
  } catch {
    return null;
  }
}

function readOptionalJSON(filePath) {
  const text = readOptionalText(filePath);
  if (text === null) return null;
  try {
    const data = JSON.parse(text);
    if (data === null || typeof data !== 'object' || Array.isArray(data)) {
      return null;
    }
    return data;
  } catch {
    return null;
  }
}

function itemHasDissentLabel(item) {
  if (item === null || typeof item !== 'object' || Array.isArray(item)) {
    return false;
  }
  for (const v of Object.values(item)) {
    if (typeof v === 'string' && DISSENT_LABELS.has(v.trim())) {
      return true;
    }
  }
  return false;
}

function claimsVerifiedOrReady(doc) {
  const verdict = typeof doc.verdict === 'string' ? doc.verdict : '';
  const status = typeof doc.status === 'string' ? doc.status : '';
  return verdict === 'VERIFIED' || status === 'ready';
}

function falseConsensusFromDoc(doc) {
  if (!claimsVerifiedOrReady(doc)) return false;
  const arrays = [];
  if (Array.isArray(doc.hunts) && doc.hunts.length > 0) arrays.push(doc.hunts);
  if (Array.isArray(doc.findings) && doc.findings.length > 0) {
    arrays.push(doc.findings);
  }
  for (const arr of arrays) {
    for (const item of arr) {
      if (!itemHasDissentLabel(item)) return true;
    }
  }
  return false;
}

/**
 * Disposition / judge-break leak predicates for a mission directory.
 * @param {string} missionDirPath
 * @returns {string[]}
 */
export function closeoutDispositionProblems(missionDirPath) {
  const problems = [];

  const judgeSummary = readOptionalJSON(
    path.join(missionDirPath, 'judge-summary.json'),
  );
  const review = readOptionalJSON(path.join(missionDirPath, 'review.json'));
  const judgeOrReview = judgeSummary ?? review;

  if (judgeOrReview) {
    if (falseConsensusFromDoc(judgeOrReview)) {
      problems.push(
        'false-consensus: VERIFIED/ready hunts or findings lack AGREE|DISAGREE_EVIDENCE|DISAGREE_CONCERN dissent labels',
      );
    }
    if (typeof judgeOrReview.builderRationale === 'string') {
      problems.push(
        'charitable-reviewer: free-text builderRationale present (structured-lines-only required)',
      );
    }
  }

  const decisions = readOptionalText(path.join(missionDirPath, 'decisions.md'));
  const evidenceText = readOptionalText(
    path.join(missionDirPath, 'evidence.jsonl'),
  );
  if (decisions !== null) {
    const mutationRequired =
      /Mutation:\s*required\b/.test(decisions) ||
      /Mutation:\s*high-risk\b/.test(decisions);
    const mutationSkipped = /Mutation skipped:/.test(decisions);
    let hasMutationEvidence = false;
    if (evidenceText !== null) {
      for (const line of evidenceText.split('\n')) {
        if (line.trim() === '') continue;
        try {
          const entry = JSON.parse(line);
          if (
            entry !== null &&
            typeof entry === 'object' &&
            !Array.isArray(entry) &&
            typeof entry.label === 'string' &&
            /mutation-/i.test(entry.label)
          ) {
            hasMutationEvidence = true;
            break;
          }
        } catch {
          // ignore malformed lines here; evidence problems own that gate
        }
      }
    }
    if (mutationRequired && !mutationSkipped && !hasMutationEvidence) {
      problems.push(
        'silent-mutation-skip: Mutation required/high-risk without Mutation skipped: or mutation- evidence',
      );
    }
  }

  const scenarios = readOptionalText(
    path.join(missionDirPath, 'approved-scenarios.md'),
  );
  if (scenarios !== null) {
    const frozen = /Approved-scenarios:\s*frozen/.test(scenarios);
    const decisionsText = decisions ?? '';
    const hasOracleChange = /Scenario oracle change:/.test(decisionsText);
    const combined = `${scenarios}\n${decisionsText}`;
    const thawOrEdit =
      /\bthawed\b/i.test(combined) || /Expected literal edited/i.test(combined);
    if (frozen && !hasOracleChange && thawOrEdit) {
      problems.push(
        'retroactive-oracle-change: frozen approved-scenarios thawed/edited without Scenario oracle change:',
      );
    }
  }

  return problems;
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

  if (existsSync(dir)) {
    problems.push(...closeoutDispositionProblems(dir));
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
