/**
 * Validate mission artifacts and evidence (`spacecraft validate` / `val`).
 * Framing: not-doc-drift / not-10X-validate (mission evidence only).
 */

import { readFileSync, statSync } from 'node:fs';
import path from 'node:path';
import { outputSHA256Hex } from './evi.mjs';
import { missionDir, normalizeID } from './resolve.mjs';

function hasHelpFlag(args) {
  return args.some((a) => a === '--help' || a === '-h');
}

function printValidateHelp() {
  console.log('Usage: spacecraft validate [--strict] [mission-id]');
  console.log('');
  console.log('Validate mission artifacts and evidence');
  console.log('Framing: not-doc-drift / not-10X-validate');
  console.log('Checks mission dir files and evidence.jsonl; not docs ↔ mission drift.');
  console.log('');
  console.log('Options:');
  console.log('  --strict                 Require exitCode on evidence; done-task evidence');
  console.log('  --help, -h               Show this help');
}

function isJSONNumber(v) {
  return typeof v === 'number' && Number.isFinite(v);
}

function validateEvidence(evidencePath, strict, missionDirPath) {
  let data;
  try {
    data = readFileSync(evidencePath, 'utf8');
  } catch {
    console.log(`x ${'evidence'.padEnd(20)} missing: ${evidencePath}`);
    return false;
  }

  const required = ['label', 'command', 'output', 'ts'];
  let entries = 0;
  let ok = true;

  const lines = data.split('\n');
  for (let i = 0; i < lines.length; i++) {
    const line = lines[i];
    if (line.trim() === '') continue;

    let entry;
    try {
      entry = JSON.parse(line);
    } catch (err) {
      console.log(`x ${'evidence'.padEnd(20)} line ${i + 1} not valid JSON: ${err.message}`);
      return false;
    }
    if (entry === null || typeof entry !== 'object' || Array.isArray(entry)) {
      console.log(`x ${'evidence'.padEnd(20)} line ${i + 1} not valid JSON: expected object`);
      return false;
    }

    for (const field of required) {
      if (!(field in entry)) {
        console.log(
          `x ${'evidence'.padEnd(20)} line ${i + 1} missing required field ${JSON.stringify(field)}`,
        );
        ok = false;
      }
    }

    if (strict && !isJSONNumber(entry.exitCode)) {
      console.log(`x ${'evidence'.padEnd(20)} line ${i + 1} missing exitCode (number)`);
      ok = false;
    }

    if (typeof entry.outputHash === 'string') {
      let hashed;
      if (entry.outputTruncated === true) {
        const rel =
          typeof entry.outputRawPath === 'string' ? entry.outputRawPath : '';
        const sidecarPath = path.join(missionDirPath, rel);
        try {
          hashed = readFileSync(sidecarPath, 'utf8');
        } catch {
          console.log(
            `x ${'evidence'.padEnd(20)} line ${i + 1} outputHash sidecar missing: ${rel || '(no outputRawPath)'}`,
          );
          ok = false;
          entries++;
          continue;
        }
      } else {
        hashed = typeof entry.output === 'string' ? entry.output : '';
      }
      if (entry.outputHash !== outputSHA256Hex(hashed)) {
        console.log(`x ${'evidence'.padEnd(20)} line ${i + 1} outputHash mismatch`);
        ok = false;
      }
    }

    entries++;
  }

  if (strict && entries === 0) {
    console.log(`x ${'evidence'.padEnd(20)} strict mode requires ≥1 evidence entry`);
    return false;
  }

  if (!ok) return false;
  console.log(`ok ${'evidence'.padEnd(20)} ${entries} entries`);
  return true;
}

function validateStrictPlanEvidence(dir) {
  let plan;
  try {
    plan = JSON.parse(readFileSync(path.join(dir, 'plan.json'), 'utf8'));
  } catch (err) {
    console.log(`x ${'strict'.padEnd(20)} cannot load plan.json: ${err.message}`);
    return false;
  }

  const tasks = Array.isArray(plan?.tasks) ? plan.tasks : [];

  let evData;
  try {
    evData = readFileSync(path.join(dir, 'evidence.jsonl'), 'utf8');
  } catch (err) {
    console.log(`x ${'strict'.padEnd(20)} cannot load evidence.jsonl: ${err.message}`);
    return false;
  }

  const entries = [];
  for (const line of evData.split('\n')) {
    if (line.trim() === '') continue;
    let raw;
    try {
      raw = JSON.parse(line);
    } catch {
      continue;
    }
    const e = {
      label: typeof raw.label === 'string' ? raw.label : '',
      exitCode: 0,
      hasExit: false,
    };
    if (isJSONNumber(raw.exitCode)) {
      e.exitCode = raw.exitCode;
      e.hasExit = true;
    }
    entries.push(e);
  }

  let ok = true;
  for (const task of tasks) {
    if (task.status !== 'done') continue;
    const allowed = new Set(Array.isArray(task.evidence) ? task.evidence : []);
    let matched = false;
    for (const e of entries) {
      if (allowed.has(e.label) && e.hasExit && e.exitCode === 0) {
        matched = true;
        break;
      }
    }
    if (!matched) {
      const tid = task.id || '(unnamed)';
      const labs = Array.isArray(task.evidence) ? task.evidence : [];
      console.log(
        `x ${'strict'.padEnd(20)} done task ${tid} missing passing evidence (exitCode 0) for labels ${JSON.stringify(labs)}`,
      );
      ok = false;
    }
  }
  return ok;
}

export function valCmd(args, spaceDir, mid) {
  if (hasHelpFlag(args)) {
    printValidateHelp();
    return 0;
  }

  let strict = false;
  let resolvedMid = '';

  for (const a of args) {
    if (a === '--strict') {
      strict = true;
      continue;
    }
    if (a.startsWith('--')) continue;
    if (resolvedMid === '') {
      resolvedMid = normalizeID(a);
    }
  }

  if (resolvedMid === '') {
    resolvedMid = mid || '';
  }
  if (resolvedMid === '') {
    console.error(
      'spacecraft validate: no mission id - pass as argument or run from feat/<id>/ branch',
    );
    return 1;
  }

  const dir = missionDir(spaceDir, resolvedMid);
  let ok = true;

  const check = (file, desc) => {
    const p = path.join(dir, file);
    try {
      statSync(p);
      console.log(`ok ${desc.padEnd(20)} ${p}`);
    } catch {
      console.log(`x ${desc.padEnd(20)} missing: ${p}`);
      ok = false;
    }
  };

  const checkJSON = (file, desc) => {
    const p = path.join(dir, file);
    let data;
    try {
      data = readFileSync(p, 'utf8');
    } catch {
      console.log(`x ${desc.padEnd(20)} missing: ${p}`);
      ok = false;
      return;
    }
    try {
      JSON.parse(data);
      console.log(`ok ${desc.padEnd(20)} valid (${Buffer.byteLength(data, 'utf8')} bytes)`);
    } catch (err) {
      console.log(`x ${desc.padEnd(20)} invalid JSON: ${p} (${err.message})`);
      ok = false;
    }
  };

  check('spec.md', 'spec');
  checkJSON('mission.json', 'mission');
  checkJSON('plan.json', 'plan');

  if (!validateEvidence(path.join(dir, 'evidence.jsonl'), strict, dir)) {
    ok = false;
  }

  if (strict && !validateStrictPlanEvidence(dir)) {
    ok = false;
  }

  return ok ? 0 : 1;
}
