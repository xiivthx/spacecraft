import assert from 'node:assert/strict';
import { spawnSync } from 'node:child_process';
import { existsSync, readFileSync } from 'node:fs';
import path from 'node:path';
import test from 'node:test';
import { fileURLToPath } from 'node:url';

const repoRoot = path.resolve(path.dirname(fileURLToPath(import.meta.url)), '../..');
const entryPath = path.join(repoRoot, 'cli', 'spacecraft.mjs');

/** Canonical help entries that must appear as `spacecraft <cmd>` (kept surface). */
const KEPT_HELP_COMMANDS = [
  'init',
  'new',
  'missions',
  'use',
  'current',
  'resolve',
  'status',
  'context',
  'drift',
  'flow',
  'bind-branch',
  'set-state',
  'clarify-status',
  'evidence',
  'validate',
  'freeze',
  'freeze-check',
  'closeout-check',
  'ship-check',
  'archive',
  'map',
  'roadmap',
  'setup',
];

test('cli/spacecraft.mjs help lists kept commands', () => {
  assert.ok(existsSync(entryPath), 'cli/spacecraft.mjs must exist');

  const source = readFileSync(entryPath, 'utf8');
  assert.ok(
    source.startsWith('#!/usr/bin/env node'),
    'cli/spacecraft.mjs must start with #!/usr/bin/env node',
  );

  const result = spawnSync(process.execPath, [entryPath, 'help'], {
    encoding: 'utf8',
    cwd: repoRoot,
  });
  assert.equal(
    result.status,
    0,
    `help exit must be 0\nstderr=${result.stderr}\nstdout=${result.stdout}`,
  );

  const out = `${result.stdout ?? ''}${result.stderr ?? ''}`;

  const missing = KEPT_HELP_COMMANDS.filter((cmd) => !out.includes(`spacecraft ${cmd}`));
  assert.deepEqual(missing, [], `help missing kept commands: ${missing.join(', ')}\n${out}`);
});

/** Unknown command path: non-zero exit and an unknown message. */
function runCli(command) {
  return spawnSync(process.execPath, [entryPath, command], {
    encoding: 'utf8',
    cwd: repoRoot,
  });
}

function assertUnknownCommandReject(command) {
  const result = runCli(command);
  const err = result.stderr ?? '';
  const out = `${result.stdout ?? ''}${err}`;

  assert.notEqual(
    result.status,
    0,
    `${command} must exit non-zero\nstderr=${err}\nstdout=${result.stdout}`,
  );
  assert.match(out, /unknown/i, `${command} must report unknown\n${out}`);
}

test('unknown command exits non-zero with unknown message', () => {
  assertUnknownCommandReject('definitely-not-a-cmd');
});

/**
 * T4-a: `spacecraft setup` must dispatch (not unknown / not "not implemented" stub forever).
 * Full reconfigure/prune behavior lives in cli/test/setup-cli.test.mjs (S4).
 */
test('spacecraft setup is a known dispatched command', () => {
  const result = runCli('setup');
  const out = `${result.stdout ?? ''}${result.stderr ?? ''}`;
  assert.doesNotMatch(
    out,
    /unknown command/i,
    `setup must not be unknown\n${out}`,
  );
  assert.doesNotMatch(
    out,
    /not implemented/i,
    `setup must be implemented (dispatch to setupCmd / runProjectSetup)\n${out}`,
  );
});
