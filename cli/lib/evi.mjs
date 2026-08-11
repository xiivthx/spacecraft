import { spawnSync } from 'node:child_process';
import { createHash } from 'node:crypto';
import { appendFileSync, mkdirSync, writeFileSync, writeSync } from 'node:fs';
import path from 'node:path';
import { missionDir, resolveActive } from './resolve.mjs';

/** Default max UTF-8 bytes kept in evidence.jsonl `output` before truncate+sidecar. */
export const EVI_OUTPUT_LIMIT = 65536;
/** Trailing marker on truncated JSONL `output` (align with ship-hook truncate style). */
export const EVI_TRUNCATE_MARKER = '\n...[truncated]';

export function outputSHA256Hex(s) {
  return createHash('sha256').update(s, 'utf8').digest('hex');
}

/** Path-safe label fragment for evidence-raw filenames. */
export function sanitizeEvidenceLabel(label) {
  const safe = String(label)
    .replace(/\.\./g, '')
    .replace(/[^A-Za-z0-9._-]+/g, '-')
    .replace(/^-+|-+$/g, '')
    .slice(0, 80);
  return safe || 'evidence';
}

/** Relative sidecar path under the mission dir: evidence-raw/<ts>-<safe-label>.log */
export function evidenceRawRelPath(label, ts) {
  const safeTs = String(ts).replace(/[:.]/g, '-');
  return `evidence-raw/${safeTs}-${sanitizeEvidenceLabel(label)}.log`;
}

/** UTF-8 prefix of at most maxBytes (does not split a multi-byte codepoint). */
function utf8BytePrefix(s, maxBytes) {
  const buf = Buffer.from(s, 'utf8');
  if (buf.length <= maxBytes) return s;
  let end = maxBytes;
  while (end > 0 && (buf[end] & 0xc0) === 0x80) end--;
  return buf.subarray(0, end).toString('utf8');
}

export function eviCmd(args, spaceDir, mid) {
  let label = '';
  let cmdArgs = [];
  let missionFlag = false;
  let resolvedMid = mid;

  for (let i = 0; i < args.length; i++) {
    if (args[i] === '--mission' && i + 1 < args.length) {
      resolvedMid = args[i + 1];
      missionFlag = true;
      i++;
    } else if (args[i] === '--') {
      cmdArgs = args.slice(i + 1);
      break;
    } else if (label === '' && !args[i].startsWith('--')) {
      label = args[i];
    }
  }

  if (!missionFlag) {
    resolvedMid = resolveActive(spaceDir, resolvedMid);
  }

  if (!resolvedMid) {
    console.error(
      'spacecraft evidence: no active mission - use --mission <id>, spacecraft use, or run from feat/<id>/ branch',
    );
    return 1;
  }
  if (!label || cmdArgs.length === 0) {
    console.error('Usage: spacecraft evidence [--mission <id>] <label> -- <command...>');
    return 1;
  }

  const cwd = path.dirname(spaceDir);
  const result = spawnSync(cmdArgs[0], cmdArgs.slice(1), {
    cwd,
    encoding: 'utf8',
    // Merge stdout and stderr into one buffer for the evidence record.
    stdio: ['ignore', 'pipe', 'pipe'],
  });

  let output = `${result.stdout ?? ''}${result.stderr ?? ''}`;
  let exitCode = 0;

  if (result.error) {
    exitCode = 127;
    output += `${result.error.message}\n`;
  } else if (result.status !== null && result.status !== 0) {
    exitCode = result.status;
  } else if (result.status === null) {
    // Signal termination - treat as nonzero failure.
    exitCode = result.signal ? 1 : 127;
  }

  const rawOutput = output;
  const outputBytes = Buffer.byteLength(rawOutput, 'utf8');
  // Hash full raw before any truncate of the JSONL field.
  const outputHash = outputSHA256Hex(rawOutput);
  const ts = new Date().toISOString().replace(/\.\d{3}Z$/, 'Z');
  const mDir = missionDir(spaceDir, resolvedMid);

  const entry = {
    label,
    command: cmdArgs.join(' '),
    output: rawOutput,
    outputHash,
    exitCode,
    ts,
  };

  if (outputBytes > EVI_OUTPUT_LIMIT) {
    const outputRawPath = evidenceRawRelPath(label, ts);
    const sidecarAbs = path.join(mDir, outputRawPath);
    try {
      mkdirSync(path.dirname(sidecarAbs), { recursive: true });
      writeFileSync(sidecarAbs, rawOutput, 'utf8');
    } catch (err) {
      console.error('spacecraft evidence:', err.message);
      return 1;
    }
    entry.output = utf8BytePrefix(rawOutput, EVI_OUTPUT_LIMIT) + EVI_TRUNCATE_MARKER;
    entry.outputTruncated = true;
    entry.outputBytes = outputBytes;
    entry.outputRawPath = outputRawPath;
  }

  const evidencePath = path.join(mDir, 'evidence.jsonl');
  try {
    mkdirSync(path.dirname(evidencePath), { recursive: true });
    appendFileSync(evidencePath, `${JSON.stringify(entry)}\n`);
  } catch (err) {
    console.error('spacecraft evidence:', err.message);
    return 1;
  }

  // writeSync: large outputs exceed the pipe buffer; async stdout.write can drop
  // the tail when the process exits before drain (exact 64KiB boundary).
  writeSync(1, rawOutput);
  console.log(`Exit code: ${exitCode}`);

  return exitCode;
}
