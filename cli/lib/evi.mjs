import { spawnSync } from 'node:child_process';
import { createHash } from 'node:crypto';
import { appendFileSync, mkdirSync } from 'node:fs';
import path from 'node:path';
import { missionDir, resolveActive } from './resolve.mjs';

export function outputSHA256Hex(s) {
  return createHash('sha256').update(s, 'utf8').digest('hex');
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

  const entry = {
    label,
    command: cmdArgs.join(' '),
    output,
    outputHash: outputSHA256Hex(output),
    exitCode,
    ts: new Date().toISOString().replace(/\.\d{3}Z$/, 'Z'),
  };

  const evidencePath = path.join(missionDir(spaceDir, resolvedMid), 'evidence.jsonl');
  try {
    mkdirSync(path.dirname(evidencePath), { recursive: true });
    appendFileSync(evidencePath, `${JSON.stringify(entry)}\n`);
  } catch (err) {
    console.error('spacecraft evidence:', err.message);
    return 1;
  }

  process.stdout.write(output);
  console.log(`Exit code: ${exitCode}`);

  return exitCode;
}
