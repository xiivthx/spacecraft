/**
 * Pack-scoped MCP merge/unmerge via installProjectSurface.
 *
 * frontend pack → .cursor/mcp.json gains shadcn; prune without frontend removes
 * only pack-managed servers; unrelated user MCP stays.
 */
import assert from 'node:assert/strict';
import { spawnSync } from 'node:child_process';
import {
  mkdirSync,
  mkdtempSync,
  readFileSync,
  rmSync,
  writeFileSync,
  existsSync,
} from 'node:fs';
import os from 'node:os';
import path from 'node:path';
import test from 'node:test';
import { fileURLToPath } from 'node:url';

import {
  installProjectSurface,
  managedCatalogPaths,
} from '../lib/project-install.mjs';
import { loadCatalog } from '../lib/packs.mjs';

const repoRoot = path.resolve(path.dirname(fileURLToPath(import.meta.url)), '../..');
const CATALOG_PATH = path.join(repoRoot, '.cursor', 'spacecraft-packs.json');

function tempTarget() {
  return mkdtempSync(path.join(os.tmpdir(), 'spacecraft-pack-mcp-'));
}

function readMcp(targetDir) {
  const p = path.join(targetDir, '.cursor', 'mcp.json');
  if (!existsSync(p)) return null;
  return JSON.parse(readFileSync(p, 'utf8'));
}

test('managedCatalogPaths collects shadcn from frontend MCP fragment', () => {
  const catalog = loadCatalog(CATALOG_PATH);
  const managed = managedCatalogPaths(catalog, {
    catalogPath: CATALOG_PATH,
    cursorDir: path.join(repoRoot, '.cursor'),
  });
  assert.ok(managed.mcpServers.has('shadcn'), 'managed MCP must include shadcn');
});

test('install with frontend merges shadcn MCP; prune without frontend removes only shadcn', () => {
  const targetDir = tempTarget();
  try {
    mkdirSync(path.join(targetDir, '.cursor'), { recursive: true });
    writeFileSync(
      path.join(targetDir, '.cursor', 'mcp.json'),
      `${JSON.stringify(
        {
          mcpServers: {
            'user-keep-mcp': { command: 'true' },
          },
        },
        null,
        2,
      )}\n`,
      'utf8',
    );

    const withFrontend = installProjectSurface(targetDir, repoRoot, {
      catalogPath: CATALOG_PATH,
      packIds: ['frontend', 'quality'],
      writeProfile: true,
      allowLegacyAll: false,
    });
    assert.ok(withFrontend.mcp.includes('shadcn'), 'result.mcp must list shadcn');

    const afterFrontend = readMcp(targetDir);
    assert.ok(afterFrontend?.mcpServers?.shadcn, 'frontend pack must merge shadcn');
    assert.equal(afterFrontend.mcpServers.shadcn.command, 'npx');
    assert.deepEqual(afterFrontend.mcpServers.shadcn.args, ['shadcn@latest', 'mcp']);
    assert.ok(
      afterFrontend.mcpServers['user-keep-mcp'],
      'user MCP must remain after frontend merge',
    );

    const withoutFrontend = installProjectSurface(targetDir, repoRoot, {
      catalogPath: CATALOG_PATH,
      packIds: ['quality'],
      writeProfile: true,
      allowLegacyAll: false,
    });
    assert.equal(
      withoutFrontend.mcp.includes('shadcn'),
      false,
      'quality-only result.mcp must not list shadcn',
    );

    const afterPrune = readMcp(targetDir);
    assert.equal(
      Object.hasOwn(afterPrune?.mcpServers ?? {}, 'shadcn'),
      false,
      'reconfigure without frontend must remove shadcn',
    );
    assert.ok(
      afterPrune?.mcpServers?.['user-keep-mcp'],
      'reconfigure without frontend must keep unrelated user MCP',
    );
  } finally {
    rmSync(targetDir, { recursive: true, force: true });
  }
});

test('root mcp merge without frontend pack does not leave shadcn on target', () => {
  const targetDir = tempTarget();
  try {
    mkdirSync(path.join(targetDir, '.cursor'), { recursive: true });
    writeFileSync(
      path.join(targetDir, '.cursor', 'mcp.json'),
      `${JSON.stringify(
        { mcpServers: { 'user-keep-mcp': { command: 'true' } } },
        null,
        2,
      )}\n`,
      'utf8',
    );

    installProjectSurface(targetDir, repoRoot, {
      catalogPath: CATALOG_PATH,
      packIds: ['quality'],
      writeProfile: true,
      allowLegacyAll: false,
    });

    const rootLike = path.join(targetDir, 'root-like-mcp.json');
    writeFileSync(
      rootLike,
      readFileSync(path.join(repoRoot, '.cursor', 'mcp.json'), 'utf8'),
      'utf8',
    );
    const merge = spawnSync(
      'python3',
      [
        path.join(repoRoot, 'scripts', 'mcp-merge.py'),
        'merge',
        path.join(targetDir, '.cursor', 'mcp.json'),
        rootLike,
        '--strip-pack-mcp',
        path.join(repoRoot, '.cursor'),
      ],
      { encoding: 'utf8' },
    );
    assert.equal(merge.status, 0, `strip merge failed: ${merge.stderr}`);

    const after = readMcp(targetDir);
    assert.equal(
      Object.hasOwn(after?.mcpServers ?? {}, 'shadcn'),
      false,
      'install merge path must not copy shadcn from root without frontend pack',
    );
    assert.ok(after?.mcpServers?.['user-keep-mcp'], 'user MCP must remain');
  } finally {
    rmSync(targetDir, { recursive: true, force: true });
  }
});
