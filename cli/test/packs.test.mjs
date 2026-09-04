/**
 * Pack catalog SoT + profile helpers (T1 / M99A9D1B).
 *
 * Expected public API from cli/lib/packs.mjs (coder exports):
 *   loadCatalog(catalogPath?) → { packs: Array<{ id, status, skills?, rules?, mcp? }> }
 *     - Default SoT: .cursor/spacecraft-packs.json (repo-relative or path arg)
 *     - status is "selectable" | "coming"
 *     - optional mcp: path relative to .cursor/ for pack MCP fragment
 *   expandPacks(packIds, catalog?, options?) → { skills: string[], rules: string[], mcp: string[] }
 *     - Union of mapped skills/rules for selectable pack ids
 *     - Always includes always-on rule 010-hard-contract.mdc in rules
 *     - mcp: resolved absolute fragment paths (union; skip packs without mcp)
 *     - Coming packs contribute no skill/rule/mcp paths (catalog empty arrays)
 *   validateProfile(profile, catalog?) → profile (or void)
 *     - Accepts { version: 1, packs: string[] } with unique selectable ids only
 *     - Throws Error for bad shape, unknown pack ids, or coming pack ids
 *
 * Oracles: design-contract.md frozen pack map + approved-scenarios S9 / E6.
 */
import assert from 'node:assert/strict';
import { spawnSync } from 'node:child_process';
import {
  existsSync,
  mkdtempSync,
  readFileSync,
  rmSync,
  writeFileSync,
} from 'node:fs';
import os from 'node:os';
import path from 'node:path';
import test from 'node:test';
import { fileURLToPath } from 'node:url';

import {
  expandPacks,
  loadCatalog,
  validateProfile,
} from '../lib/packs.mjs';

const repoRoot = path.resolve(path.dirname(fileURLToPath(import.meta.url)), '../..');
const CATALOG_PATH = path.join(repoRoot, '.cursor', 'spacecraft-packs.json');

/** Frozen design-contract selectable ids. */
const SELECTABLE = ['frontend', 'backend', 'database', 'embedded', 'quality', 'fpga'];
/** Frozen design-contract coming-only ids. */
const COMING = ['iot', 'pcb', 'management'];

/** Frozen pack → skills / rules (design-contract v1). Always-on not in pack rows. */
const PACK_MAP = {
  frontend: {
    skills: ['sc-web-frontend', 'sc-ux-design', 'sc-browser-probe'],
    rules: ['150-design.mdc'],
    mcp: 'mcp-packs/frontend.json',
  },
  backend: {
    skills: ['sc-web-backend'],
    rules: [],
  },
  database: {
    skills: ['sc-database'],
    rules: ['500-database.mdc'],
  },
  embedded: {
    skills: ['sc-firmware'],
    rules: [
      '600-firmware.mdc',
      '610-firmware-peripherals.mdc',
      '620-firmware-testing.mdc',
    ],
  },
  quality: {
    skills: [
      'sc-security',
      'sc-performance',
      'sc-solid',
      'sc-architect',
      'sc-diagram',
    ],
    rules: ['300-security.mdc', '400-performance.mdc'],
  },
  fpga: {
    skills: ['sc-rtl', 'sc-rtl-verify'],
    rules: [
      '700-rtl.mdc',
      '710-rtl-timing.mdc',
      '720-rtl-verify.mdc',
    ],
  },
};

const FRONTEND_MCP_FRAGMENT = path.join(repoRoot, '.cursor', 'mcp-packs', 'frontend.json');

/** Always-on rule (design-contract + S9). */
const ALWAYS_ON_RULE = '010-hard-contract.mdc';

/** S9: expand packs=[frontend,quality]. */
const S9_SKILLS = [
  'sc-web-frontend',
  'sc-ux-design',
  'sc-browser-probe',
  'sc-security',
  'sc-performance',
  'sc-solid',
  'sc-architect',
  'sc-diagram',
];
const S9_RULES = [
  '150-design.mdc',
  '300-security.mdc',
  '400-performance.mdc',
  ALWAYS_ON_RULE,
];

function sorted(arr) {
  return [...arr].sort();
}

function catalog() {
  return loadCatalog(CATALOG_PATH);
}

function byId(cat) {
  const map = new Map();
  for (const entry of cat.packs) {
    map.set(entry.id, entry);
  }
  return map;
}

function idsWithStatus(cat, status) {
  return cat.packs.filter((p) => p.status === status).map((p) => p.id);
}

// --- T1 acceptance 1 ---
test('catalog exposes selectable and coming packs; coming never resolve to skill or rule paths', () => {
  const cat = catalog();
  assert.ok(cat && Array.isArray(cat.packs), 'loadCatalog must return { packs: [] }');

  assert.deepEqual(
    sorted(idsWithStatus(cat, 'selectable')),
    sorted(SELECTABLE),
    'selectable pack ids must match design-contract',
  );
  assert.deepEqual(
    sorted(idsWithStatus(cat, 'coming')),
    sorted(COMING),
    'coming pack ids must match design-contract',
  );

  const entries = byId(cat);
  for (const id of COMING) {
    const entry = entries.get(id);
    assert.ok(entry, `coming pack ${id} must exist in catalog`);
    assert.equal(entry.status, 'coming');
    assert.deepEqual(
      entry.skills ?? [],
      [],
      `coming pack ${id} must not resolve to skill paths`,
    );
    assert.deepEqual(
      entry.rules ?? [],
      [],
      `coming pack ${id} must not resolve to rule paths`,
    );
  }
});

// --- T1 acceptance 2 (+ S9) ---
test('expandPacks frozen pack→skill/rule map matches design-contract and S9 union', () => {
  const cat = catalog();
  const cursorDir = path.join(repoRoot, '.cursor');

  for (const id of SELECTABLE) {
    const got = expandPacks([id], cat, { cursorDir, catalogPath: CATALOG_PATH });
    const frozen = PACK_MAP[id];
    assert.deepEqual(
      sorted(got.skills),
      sorted(frozen.skills),
      `${id} skills must match design-contract`,
    );
    assert.deepEqual(
      sorted(got.rules),
      sorted([...frozen.rules, ALWAYS_ON_RULE]),
      `${id} rules must match design-contract plus always-on ${ALWAYS_ON_RULE}`,
    );
    const wantMcp = frozen.mcp ? [path.resolve(cursorDir, frozen.mcp)] : [];
    assert.deepEqual(
      sorted(got.mcp ?? []),
      sorted(wantMcp),
      `${id} mcp fragment paths must match catalog`,
    );
  }

  // S9: packs=[frontend,quality]
  const s9 = expandPacks(['frontend', 'quality'], cat, { cursorDir, catalogPath: CATALOG_PATH });
  assert.deepEqual(sorted(s9.skills), sorted(S9_SKILLS), 'S9 skills union');
  assert.deepEqual(sorted(s9.rules), sorted(S9_RULES), 'S9 rules union + always-on');
  assert.deepEqual(
    sorted(s9.mcp ?? []),
    sorted([FRONTEND_MCP_FRAGMENT]),
    'S9 mcp union includes frontend fragment only',
  );
});

test('frontend pack catalogs shadcn MCP fragment; root merge strips pack servers', () => {
  const cat = catalog();
  const entry = byId(cat).get('frontend');
  assert.equal(entry?.mcp, 'mcp-packs/frontend.json');
  assert.ok(existsSync(FRONTEND_MCP_FRAGMENT), 'frontend MCP fragment must exist');
  const fragment = JSON.parse(readFileSync(FRONTEND_MCP_FRAGMENT, 'utf8'));
  assert.ok(fragment.mcpServers?.shadcn, 'fragment must define shadcn server');
  assert.equal(fragment.mcpServers.shadcn.command, 'npx');
  assert.deepEqual(fragment.mcpServers.shadcn.args, ['shadcn@latest', 'mcp']);

  // Meta checkout may keep shadcn in root .cursor/mcp.json for local Cursor;
  // install-cursor must strip pack-managed names when merging into targets.
  const tmp = mkdtempSync(path.join(os.tmpdir(), 'spacecraft-mcp-strip-'));
  try {
    const targetMcp = path.join(tmp, 'target-mcp.json');
    const sourceMcp = path.join(tmp, 'source-mcp.json');
    writeFileSync(
      sourceMcp,
      `${JSON.stringify(
        {
          mcpServers: {
            shadcn: { command: 'npx', args: ['shadcn@latest', 'mcp'] },
            'always-on-demo': { command: 'true' },
          },
        },
        null,
        2,
      )}\n`,
      'utf8',
    );
    writeFileSync(
      targetMcp,
      `${JSON.stringify({ mcpServers: { 'user-keep-mcp': { command: 'true' } } }, null, 2)}\n`,
      'utf8',
    );

    const merge = spawnSync(
      'python3',
      [
        path.join(repoRoot, 'scripts', 'mcp-merge.py'),
        'merge',
        targetMcp,
        sourceMcp,
        '--strip-pack-mcp',
        path.join(repoRoot, '.cursor'),
      ],
      { encoding: 'utf8' },
    );
    assert.equal(
      merge.status,
      0,
      `mcp-merge --strip-pack-mcp must exit 0\nstderr=${merge.stderr}\nstdout=${merge.stdout}`,
    );

    const got = JSON.parse(readFileSync(targetMcp, 'utf8'));
    assert.equal(
      Object.hasOwn(got.mcpServers ?? {}, 'shadcn'),
      false,
      'strip-pack-mcp must not copy shadcn from root-like source into target',
    );
    assert.ok(
      got.mcpServers?.['always-on-demo'],
      'non-pack source servers must still merge',
    );
    assert.ok(
      got.mcpServers?.['user-keep-mcp'],
      'pre-existing user MCP must remain',
    );
  } finally {
    rmSync(tmp, { recursive: true, force: true });
  }
});

// --- T1 acceptance 3 ---
test('validateProfile accepts v1 selectable shape and rejects unknown or coming pack ids', () => {
  const cat = catalog();

  const ok = validateProfile({ version: 1, packs: ['frontend', 'quality'] }, cat);
  if (ok !== undefined) {
    assert.equal(ok.version, 1);
    assert.deepEqual(sorted(ok.packs), sorted(['frontend', 'quality']));
  }

  assert.throws(
    () => validateProfile({ version: 1, packs: ['iot'] }, cat),
    (err) => err instanceof Error,
    'coming pack id iot must be rejected with Error',
  );
  assert.throws(
    () => validateProfile({ version: 1, packs: ['not-a-pack'] }, cat),
    (err) => err instanceof Error,
    'unknown pack id must be rejected with Error',
  );
  assert.throws(
    () => validateProfile({ packs: ['quality'] }, cat),
    (err) => err instanceof Error,
    'missing version must be rejected with Error',
  );
  assert.throws(
    () => validateProfile({ version: 1, packs: 'quality' }, cat),
    (err) => err instanceof Error,
    'non-array packs must be rejected with Error',
  );
});
