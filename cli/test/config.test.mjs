/**
 * T2 — .space/config.json criticFamily reader (M9G7II1F quality-tooling).
 */
import assert from 'node:assert/strict';
import { spawnSync } from 'node:child_process';
import {
  existsSync,
  mkdirSync,
  mkdtempSync,
  rmSync,
  writeFileSync,
} from 'node:fs';
import os from 'node:os';
import path from 'node:path';
import test from 'node:test';
import { fileURLToPath } from 'node:url';

const repoRoot = path.resolve(path.dirname(fileURLToPath(import.meta.url)), '../..');
const configModulePath = path.join(repoRoot, 'cli', 'lib', 'config.mjs');

async function loadConfig() {
  const mod = await import('../lib/config.mjs');
  assert.equal(typeof mod.readCriticFamily, 'function', 'readCriticFamily must be exported');
  assert.equal(typeof mod.readSpaceConfig, 'function', 'readSpaceConfig must be exported');
  return mod;
}

function spaceRoot() {
  const dir = mkdtempSync(path.join(os.tmpdir(), 'spacecraft-config-'));
  mkdirSync(path.join(dir, '.space'), { recursive: true });
  return dir;
}

function cleanup(dir) {
  rmSync(dir, { recursive: true, force: true });
}

function writeConfig(root, body) {
  writeFileSync(path.join(root, '.space', 'config.json'), body);
}

/** Invoke readSpaceConfig in a subprocess (malformed config must exit non-zero). */
function runReadSpaceConfig(spaceDir) {
  const spaceQ = JSON.stringify(spaceDir);
  const modQ = JSON.stringify(configModulePath);
  const script = `
import { readSpaceConfig } from ${modQ};
try {
  readSpaceConfig(${spaceQ});
  process.exit(0);
} catch (err) {
  console.error(err?.message ?? String(err));
  process.exit(1);
}
`;
  return spawnSync(process.execPath, ['--input-type=module', '-e', script], {
    encoding: 'utf8',
    cwd: repoRoot,
  });
}

test('pos-config-reader: valid criticFamily and unconfigured cases return null or string', async () => {
  const { readCriticFamily } = await loadConfig();
  const dir = spaceRoot();
  try {
    writeConfig(dir, '{"criticFamily":"gpt"}\n');
    assert.equal(readCriticFamily(dir), 'gpt');

    rmSync(path.join(dir, '.space', 'config.json'));
    assert.equal(readCriticFamily(dir), null, 'missing file must be unconfigured');

    writeConfig(dir, '{"otherKey":"ignored"}\n');
    assert.equal(readCriticFamily(dir), null, 'absent criticFamily key must be unconfigured');
  } finally {
    cleanup(dir);
  }
});

test('neg-malformed-config: invalid JSON or non-string criticFamily exits non-zero', async () => {
  await loadConfig();

  for (const [label, body] of [
    ['invalid JSON', '{not json'],
    ['non-string criticFamily', '{"criticFamily":1}\n'],
  ]) {
    const dir = spaceRoot();
    try {
      writeConfig(dir, body);
      const result = runReadSpaceConfig(dir);
      const out = `${result.stdout ?? ''}${result.stderr ?? ''}`;
      assert.notEqual(
        result.status,
        0,
        `${label}: expected non-zero exit\n${out}`,
      );
      assert.match(
        out,
        /spacecraft config:|malformed/i,
        `${label}: stderr must name config/malformed\n${out}`,
      );
    } finally {
      cleanup(dir);
    }
  }
});
