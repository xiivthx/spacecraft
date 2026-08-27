import { createHash } from 'node:crypto';
import {
  appendFileSync,
  existsSync,
  mkdirSync,
  readFileSync,
  realpathSync,
  writeFileSync,
} from 'node:fs';
import path from 'node:path';
import {
  EVI_OUTPUT_LIMIT,
  EVI_TRUNCATE_MARKER,
  evidenceRawRelPath,
  outputSHA256Hex,
} from './evi.mjs';
import { gatesAtOrAfter, readGatesVersion } from './gates.mjs';
import { missionDir, resolveActive } from './resolve.mjs';

const FREEZE_GATE = 'M9G7IHV3';

function readOptionalText(filePath) {
  try {
    return readFileSync(filePath, 'utf8');
  } catch {
    return null;
  }
}

function sha256File(absPath) {
  const data = readFileSync(absPath);
  return createHash('sha256').update(data).digest('hex');
}

function repoRootFromSpaceDir(spaceDir) {
  return path.dirname(spaceDir);
}

function normalizeRepoPath(repoRoot, filePath) {
  const root = safeRealpath(path.resolve(repoRoot));
  const abs = safeRealpath(
    path.isAbsolute(filePath)
      ? path.resolve(filePath)
      : path.resolve(repoRoot, filePath),
  );
  const rel = path.relative(root, abs);
  if (rel.startsWith('..') || path.isAbsolute(rel)) {
    throw new Error(`path outside repository: ${filePath}`);
  }
  return rel.split(path.sep).join('/');
}

function safeRealpath(p) {
  try {
    return realpathSync.native?.(p) ?? realpathSync(p);
  } catch {
    return path.resolve(p);
  }
}

function isTestRunLabel(label) {
  if (typeof label !== 'string') return false;
  return /^test-/.test(label) || label.includes('test-run');
}

function parseEvidenceLines(evidencePath) {
  const problems = [];
  const entries = [];
  let data;
  try {
    data = readFileSync(evidencePath, 'utf8');
  } catch {
    return { problems: ['missing evidence.jsonl'], entries: [] };
  }

  const lines = data.split('\n');
  for (let i = 0; i < lines.length; i++) {
    const line = lines[i];
    if (line.trim() === '') continue;
    try {
      const entry = JSON.parse(line);
      if (entry === null || typeof entry !== 'object' || Array.isArray(entry)) {
        problems.push(`evidence line ${i + 1} not valid JSON`);
        continue;
      }
      entries.push({ index: i, entry });
    } catch {
      problems.push(`evidence line ${i + 1} not valid JSON`);
    }
  }
  return { problems, entries };
}

function manifestOutputText(entry, missionDirPath) {
  if (entry.outputTruncated === true && typeof entry.outputRawPath === 'string') {
    const sidecar = path.join(missionDirPath, entry.outputRawPath);
    try {
      return readFileSync(sidecar, 'utf8');
    } catch {
      return null;
    }
  }
  return typeof entry.output === 'string' ? entry.output : null;
}

/**
 * @param {{ label?: string, output?: string, outputTruncated?: boolean, outputRawPath?: string }} entry
 * @param {string} missionDirPath
 * @returns {{ kind: string, files: Array<{ path: string, sha256: string }> }|null}
 */
export function parseFreezeManifest(entry, missionDirPath) {
  if (!entry || entry.label !== 'freeze') return null;
  const raw = manifestOutputText(entry, missionDirPath);
  if (raw === null) return null;
  try {
    const manifest = JSON.parse(raw);
    if (
      manifest === null ||
      typeof manifest !== 'object' ||
      Array.isArray(manifest) ||
      manifest.kind !== 'freeze-manifest' ||
      !Array.isArray(manifest.files)
    ) {
      return null;
    }
    return manifest;
  } catch {
    return null;
  }
}

function lastIndexByLabel(entries, predicate) {
  let last = -1;
  for (const { index, entry } of entries) {
    if (predicate(entry)) last = index;
  }
  return last;
}

function approvedScenariosSkipped(decisionsText, scenariosText) {
  const combined = `${decisionsText ?? ''}\n${scenariosText ?? ''}`;
  return /Approved-scenarios skipped:/.test(combined);
}

function approvedScenariosFrozen(scenariosText) {
  return scenariosText !== null && /Approved-scenarios:\s*frozen/.test(scenariosText);
}

function freezeRequired(decisionsText, scenariosText, missionGate) {
  if (!gatesAtOrAfter(missionGate, FREEZE_GATE)) return false;
  if (approvedScenariosSkipped(decisionsText, scenariosText)) return false;
  return approvedScenariosFrozen(scenariosText);
}

function postdatedFreezeProblem(entries) {
  const lastTestIdx = lastIndexByLabel(entries, (e) => isTestRunLabel(e.label));
  if (lastTestIdx === -1) return null;

  const freezeEntries = entries.filter((e) => e.entry.label === 'freeze');
  if (freezeEntries.length === 0) {
    return 'postdated-freeze: test-run evidence exists but no freeze event yet';
  }

  const latestFreeze = freezeEntries[freezeEntries.length - 1];
  if (latestFreeze.index > lastTestIdx) {
    return 'postdated-freeze: latest freeze event postdates test-run evidence in append-only log';
  }
  return null;
}

function wouldAppendPostdate(entries) {
  const lastTestIdx = lastIndexByLabel(entries, (e) => isTestRunLabel(e.label));
  return lastTestIdx !== -1;
}

function driftProblems(manifest, repoRoot, decisionsText) {
  const problems = [];
  const hasOracleChange = /Scenario oracle change:/.test(decisionsText ?? '');

  for (const item of manifest.files) {
    const rel = item.path;
    const abs = path.join(repoRoot, rel);
    if (!existsSync(abs)) {
      if (!hasOracleChange) {
        problems.push(
          `freeze-drift: ${rel} deleted (missing on disk) without Scenario oracle change:`,
        );
      }
      continue;
    }
    const current = sha256File(abs);
    if (current !== item.sha256) {
      if (!hasOracleChange) {
        problems.push(
          `freeze-drift: ${rel} hash changed without Scenario oracle change:`,
        );
      }
    }
  }
  return problems;
}

/**
 * Shared freeze-check predicates for freeze-check CLI and closeout.
 * @param {string} missionDirPath absolute mission directory
 * @param {string} spaceDir absolute .space directory
 * @param {string} [repoRoot] repository root (defaults to parent of spaceDir)
 * @returns {string[]}
 */
export function freezeCheckProblems(missionDirPath, spaceDir, repoRoot) {
  const root = repoRoot ?? repoRootFromSpaceDir(spaceDir);
  const decisionsText = readOptionalText(path.join(missionDirPath, 'decisions.md'));
  const scenariosText = readOptionalText(
    path.join(missionDirPath, 'approved-scenarios.md'),
  );
  const missionGate = decisionsText !== null ? readGatesVersion(decisionsText) : null;

  if (!gatesAtOrAfter(missionGate, FREEZE_GATE)) {
    return [];
  }
  if (approvedScenariosSkipped(decisionsText ?? '', scenariosText)) {
    return [];
  }
  if (!freezeRequired(decisionsText ?? '', scenariosText, missionGate)) {
    return [];
  }

  const evidencePath = path.join(missionDirPath, 'evidence.jsonl');
  const { problems: parseProblems, entries } = parseEvidenceLines(evidencePath);
  if (parseProblems.length > 0) {
    return parseProblems;
  }

  const postdated = postdatedFreezeProblem(entries);
  if (postdated) return [postdated];

  const freezeEntries = entries.filter((e) => e.entry.label === 'freeze');
  if (freezeEntries.length === 0) {
    return ['missing-freeze: Approved-scenarios frozen footer requires a freeze evidence event'];
  }

  const latestFreeze = freezeEntries[freezeEntries.length - 1].entry;
  const manifest = parseFreezeManifest(latestFreeze, missionDirPath);
  if (!manifest) {
    return ['freeze manifest invalid or unreadable in latest freeze evidence event'];
  }

  return driftProblems(manifest, root, decisionsText ?? '');
}

function utf8BytePrefix(s, maxBytes) {
  const buf = Buffer.from(s, 'utf8');
  if (buf.length <= maxBytes) return s;
  let end = maxBytes;
  while (end > 0 && (buf[end] & 0xc0) === 0x80) end--;
  return buf.subarray(0, end).toString('utf8');
}

function appendFreezeEvidence(mDir, manifestJson, pathsArg) {
  const rawOutput = manifestJson;
  const outputBytes = Buffer.byteLength(rawOutput, 'utf8');
  const outputHash = outputSHA256Hex(rawOutput);
  const ts = new Date().toISOString().replace(/\.\d{3}Z$/, 'Z');

  const entry = {
    label: 'freeze',
    command: `spacecraft freeze ${pathsArg}`,
    output: rawOutput,
    outputHash,
    exitCode: 0,
    ts,
  };

  if (outputBytes > EVI_OUTPUT_LIMIT) {
    const outputRawPath = evidenceRawRelPath('freeze', ts);
    const sidecarAbs = path.join(mDir, outputRawPath);
    mkdirSync(path.dirname(sidecarAbs), { recursive: true });
    writeFileSync(sidecarAbs, rawOutput, 'utf8');
    entry.output = utf8BytePrefix(rawOutput, EVI_OUTPUT_LIMIT) + EVI_TRUNCATE_MARKER;
    entry.outputTruncated = true;
    entry.outputBytes = outputBytes;
    entry.outputRawPath = outputRawPath;
  }

  const evidencePath = path.join(mDir, 'evidence.jsonl');
  appendFileSync(evidencePath, `${JSON.stringify(entry)}\n`);
}

export function freezeCmd(args, spaceDir, mid) {
  let resolvedMid = mid;
  let missionFlag = false;
  const filePaths = [];

  for (let i = 0; i < args.length; i++) {
    if (args[i] === '--mission' && i + 1 < args.length) {
      resolvedMid = args[i + 1];
      missionFlag = true;
      i++;
    } else if (!args[i].startsWith('--')) {
      filePaths.push(args[i]);
    }
  }

  if (!missionFlag) {
    resolvedMid = resolveActive(spaceDir, resolvedMid);
  }

  if (!resolvedMid) {
    console.error(
      'spacecraft freeze: no active mission - use --mission <id>, spacecraft use, or run from feat/<id>/ branch',
    );
    return 1;
  }

  if (filePaths.length === 0) {
    console.error('Usage: spacecraft freeze [--mission <id>] <path...>');
    return 1;
  }

  const repoRoot = repoRootFromSpaceDir(spaceDir);
  const mDir = missionDir(spaceDir, resolvedMid);
  const decisionsText = readOptionalText(path.join(mDir, 'decisions.md'));
  const scenariosText = readOptionalText(path.join(mDir, 'approved-scenarios.md'));
  const missionGate = decisionsText !== null ? readGatesVersion(decisionsText) : null;

  const evidencePath = path.join(mDir, 'evidence.jsonl');
  const { problems: parseProblems, entries } = parseEvidenceLines(evidencePath);
  if (parseProblems.length > 0) {
    for (const p of parseProblems) console.log(p);
    return 1;
  }

  if (
    gatesAtOrAfter(missionGate, FREEZE_GATE) &&
    !approvedScenariosSkipped(decisionsText ?? '', scenariosText) &&
    wouldAppendPostdate(entries)
  ) {
    console.log(
      'postdated-freeze: cannot append freeze after test-run evidence in append-only log',
    );
    return 1;
  }

  const manifestFiles = [];
  for (const fp of filePaths) {
    let rel;
    try {
      rel = normalizeRepoPath(repoRoot, fp);
    } catch (err) {
      console.error(`spacecraft freeze: ${err.message}`);
      return 1;
    }
    const abs = path.join(repoRoot, rel);
    if (!existsSync(abs)) {
      console.error(`spacecraft freeze: file not found: ${rel}`);
      return 1;
    }
    const content = readFileSync(abs);
    const entry = { path: rel, sha256: createHash('sha256').update(content).digest('hex') };
    if (content.length >= EVI_OUTPUT_LIMIT) {
      entry.content = content.toString('utf8');
    }
    manifestFiles.push(entry);
  }

  const manifest = JSON.stringify({
    kind: 'freeze-manifest',
    files: manifestFiles,
  });

  try {
    mkdirSync(mDir, { recursive: true });
    appendFreezeEvidence(mDir, manifest, filePaths.join(' '));
  } catch (err) {
    console.error('spacecraft freeze:', err.message);
    return 1;
  }

  console.log(manifest);
  return 0;
}

export function freezeCheckCmd(args, spaceDir, mid) {
  let resolvedMid = mid;
  let missionFlag = false;

  for (let i = 0; i < args.length; i++) {
    if (args[i] === '--mission' && i + 1 < args.length) {
      resolvedMid = args[i + 1];
      missionFlag = true;
      i++;
    }
  }

  if (!missionFlag) {
    resolvedMid = resolveActive(spaceDir, resolvedMid);
  }

  if (!resolvedMid) {
    console.error(
      'spacecraft freeze-check: no active mission - use --mission <id>, spacecraft use, or run from feat/<id>/ branch',
    );
    return 1;
  }

  const mDir = missionDir(spaceDir, resolvedMid);
  const repoRoot = repoRootFromSpaceDir(spaceDir);
  const problems = freezeCheckProblems(mDir, spaceDir, repoRoot);

  if (problems.length > 0) {
    for (const p of problems) console.log(p);
    return 1;
  }
  return 0;
}
