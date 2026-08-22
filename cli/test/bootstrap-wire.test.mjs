/**
 * Bootstrap / install-project wire (T5 / M99A9D1B) — RED against shell entrypoints.
 *
 * Shared seam (coder GREEN — do not invent alternate entrypoints):
 *   scripts/install-cursor.sh  →  node cli/lib/setup.mjs <target> <source>
 *     (same pack-resolve + selective install as `spacecraft setup` reconcile /
 *      runProjectSetup; never bare project-install.mjs legacy-all)
 *   Makefile install-project    →  scripts/install-cursor.sh
 *   bootstrap.sh                →  scripts/install-cursor.sh
 *
 * Mode matrix via install-cursor (design-contract public seam + Edge E1–E3 /
 * approved-scenarios S1–S3):
 *   profile present → silent reconcile (no prompt; quality-only when packs=["quality"])
 *   no profile + non-TTY + no packs → exit ≠ 0; no all-domain skill tree
 *   no profile + SPACECRAFT_PACKS=frontend,quality → write profile; selective install
 *
 * Coder must change:
 *   - scripts/install-cursor.sh: replace `node …/project-install.mjs` domain install
 *     with `node …/cli/lib/setup.mjs "$TARGET_ABS" "$SRC_ABS"` (TTY from stdout;
 *     SPACECRAFT_PACKS from env). Keep hooks/MCP/.space scaffolding around it.
 *   - Makefile install-project / bootstrap.sh: keep calling install-cursor.sh
 *     (already do); do not call project-install.mjs directly.
 *   - Retire legacy-all as the install-cursor default when profile+packs absent.
 *
 * Pbt skipped: not core logic (thin bootstrap/Make wire)
 */
import assert from 'node:assert/strict';
import { spawnSync } from 'node:child_process';
import {
  existsSync,
  mkdirSync,
  mkdtempSync,
  readFileSync,
  rmSync,
  writeFileSync,
} from 'node:fs';
import os from 'node:os';
import path from 'node:path';
import test from 'node:test';
import { fileURLToPath } from 'node:url';

const repoRoot = path.resolve(path.dirname(fileURLToPath(import.meta.url)), '../..');
const INSTALL_CURSOR = path.join(repoRoot, 'scripts', 'install-cursor.sh');
const BOOTSTRAP = path.join(repoRoot, 'bootstrap.sh');
const MAKEFILE = path.join(repoRoot, 'Makefile');
const PROFILE_REL = path.join('.cursor', 'spacecraft-profile.json');
const ALWAYS_ON_RULE = '010-hard-contract.mdc';

/** Frozen design-contract pack → skills (v1). */
const PACK_SKILLS = Object.freeze({
  frontend: ['sc-web-frontend', 'sc-ux-design', 'sc-browser-probe'],
  backend: ['sc-web-backend'],
  database: ['sc-database'],
  embedded: ['sc-firmware'],
  quality: [
    'sc-security',
    'sc-performance',
    'sc-solid',
    'sc-architect',
    'sc-diagram',
  ],
});

/** Frozen design-contract pack → rules (v1; always-on separate). */
const PACK_RULES = Object.freeze({
  frontend: ['150-design.mdc'],
  backend: [],
  database: ['500-database.mdc'],
  embedded: [
    '600-firmware.mdc',
    '610-firmware-peripherals.mdc',
    '620-firmware-testing.mdc',
  ],
  quality: ['300-security.mdc', '400-performance.mdc'],
});

const S2_PACKS = ['frontend', 'quality'];
const S2_SKILLS = [...PACK_SKILLS.frontend, ...PACK_SKILLS.quality];
const S2_RULES = [...PACK_RULES.frontend, ...PACK_RULES.quality, ALWAYS_ON_RULE];
const S3_SKILLS = [...PACK_SKILLS.quality];
const S3_RULES = [...PACK_RULES.quality, ALWAYS_ON_RULE];
const S3_ABSENT_PACKS = ['frontend', 'backend', 'database', 'embedded'];

const ALL_DOMAIN_SKILLS = Object.freeze([
  ...new Set(Object.values(PACK_SKILLS).flat()),
]);

function sorted(arr) {
  return [...arr].sort();
}

function tempTarget(prefix = 'spacecraft-bootstrap-wire-') {
  return mkdtempSync(path.join(os.tmpdir(), prefix));
}

function envWithoutPacks(extra = {}) {
  const env = { ...process.env, ...extra };
  // Only clear when caller did not pass SPACECRAFT_PACKS (S1/S3); S2 sets it in extra.
  if (!Object.hasOwn(extra, 'SPACECRAFT_PACKS')) {
    delete env.SPACECRAFT_PACKS;
  }
  return env;
}

function combined(res) {
  return `${res.stdout ?? ''}${res.stderr ?? ''}`;
}

function runInstallCursor(targetDir, env) {
  return spawnSync('sh', [INSTALL_CURSOR, targetDir, repoRoot], {
    encoding: 'utf8',
    cwd: repoRoot,
    env,
  });
}

function runMakeInstallProject(targetDir, env) {
  return spawnSync('make', ['install-project', `PROJECT=${targetDir}`], {
    encoding: 'utf8',
    cwd: repoRoot,
    env,
  });
}

function skillInstalled(targetDir, name) {
  return existsSync(path.join(targetDir, '.cursor', 'skills', name, 'SKILL.md'));
}

function ruleInstalled(targetDir, name) {
  return existsSync(path.join(targetDir, '.cursor', 'rules', name));
}

function profilePath(targetDir) {
  return path.join(targetDir, PROFILE_REL);
}

function assertNoAllPacksDomainInstall(targetDir) {
  const present = ALL_DOMAIN_SKILLS.filter((name) => skillInstalled(targetDir, name));
  assert.deepEqual(
    present,
    [],
    `S1: must not silent-install all-packs domain skills (found: ${present.join(', ')})`,
  );
}

/** Tab-indented recipe lines under `target:` (skips blanks; stops at next rule). */
function makefileRecipe(makefile, target) {
  const lines = makefile.split('\n');
  const start = lines.findIndex((line) => new RegExp(`^${target}\\s*:`).test(line));
  assert.ok(start >= 0, `Makefile must define target ${target}`);
  const body = [];
  for (let i = start + 1; i < lines.length; i++) {
    const line = lines[i];
    if (line.startsWith('\t')) {
      body.push(line);
      continue;
    }
    if (line.trim() === '') continue;
    break;
  }
  return body.join('\n');
}

test('A1: install-cursor / Makefile / bootstrap share setup.mjs pack-resolve path', () => {
  const installCursor = readFileSync(INSTALL_CURSOR, 'utf8');
  const bootstrap = readFileSync(BOOTSTRAP, 'utf8');
  const makefile = readFileSync(MAKEFILE, 'utf8');
  const installProjectBody = makefileRecipe(makefile, 'install-project');

  // Active node invocation must be setup.mjs (shared with spacecraft setup reconcile).
  assert.match(
    installCursor,
    /node\s+"\$SRC_ABS\/cli\/lib\/setup\.mjs"/,
    'install-cursor.sh must invoke cli/lib/setup.mjs for pack-resolve + selective install',
  );
  const codeLines = installCursor
    .split('\n')
    .filter((line) => !/^\s*#/.test(line))
    .join('\n');
  assert.doesNotMatch(
    codeLines,
    /node\s+[^\n]*cli\/lib\/project-install\.mjs/,
    'install-cursor.sh must not call project-install.mjs as the domain install entry (legacy-all)',
  );

  assert.match(
    installProjectBody,
    /scripts\/install-cursor\.sh/,
    'make install-project must call scripts/install-cursor.sh',
  );
  assert.match(
    bootstrap,
    /scripts\/install-cursor\.sh/,
    'bootstrap.sh must call scripts/install-cursor.sh',
  );
});

test('S1: install-cursor no profile + non-TTY + no packs → exit ≠ 0 and no all-packs', () => {
  const targetDir = tempTarget();
  try {
    const res = runInstallCursor(targetDir, envWithoutPacks());
    assert.notEqual(
      res.status,
      0,
      `S1/E1: install-cursor must fail closed\n${combined(res)}`,
    );
    assert.match(
      combined(res),
      /SPACECRAFT_PACKS|spacecraft-profile|Non-interactive setup/,
      'S1: fail path must surface setup fail-closed (packs-or-profile) message',
    );
    assert.equal(
      existsSync(profilePath(targetDir)),
      false,
      'S1: must not write spacecraft-profile.json on fail',
    );
    assertNoAllPacksDomainInstall(targetDir);
  } finally {
    rmSync(targetDir, { recursive: true, force: true });
  }
});

test('S1: make install-project no profile + non-TTY + no packs → exit ≠ 0 and no all-packs', () => {
  const targetDir = tempTarget();
  try {
    const res = runMakeInstallProject(targetDir, envWithoutPacks());
    assert.notEqual(
      res.status,
      0,
      `S1/E1: make install-project must fail closed\n${combined(res)}`,
    );
    assertNoAllPacksDomainInstall(targetDir);
  } finally {
    rmSync(targetDir, { recursive: true, force: true });
  }
});

test('S2: install-cursor SPACECRAFT_PACKS=frontend,quality → write profile and selective install', () => {
  const targetDir = tempTarget();
  try {
    const res = runInstallCursor(
      targetDir,
      envWithoutPacks({ SPACECRAFT_PACKS: 'frontend,quality' }),
    );
    assert.equal(
      res.status,
      0,
      `S2/E2: install-cursor with packs must exit 0\n${combined(res)}`,
    );
    // setup.mjs main log shape (not project-install legacy label).
    assert.match(
      combined(res),
      /\(env packs=/,
      'S2: must install via setup.mjs reconcile log (env packs=…), not project-install SPACECRAFT_PACKS=… label',
    );
    assert.doesNotMatch(
      combined(res),
      /all selectable packs \(legacy\)/,
      'S2: must not use legacy-all',
    );

    assert.ok(
      existsSync(profilePath(targetDir)),
      'S2: must write .cursor/spacecraft-profile.json',
    );
    const profile = JSON.parse(readFileSync(profilePath(targetDir), 'utf8'));
    assert.equal(profile.version, 1);
    assert.deepEqual(
      sorted(profile.packs),
      sorted(S2_PACKS),
      'S2: profile packs must be exactly frontend,quality',
    );

    for (const skill of S2_SKILLS) {
      assert.equal(
        skillInstalled(targetDir, skill),
        true,
        `S2: missing mapped skill ${skill}`,
      );
    }
    for (const rule of S2_RULES) {
      assert.equal(
        ruleInstalled(targetDir, rule),
        true,
        `S2: missing mapped/always-on rule ${rule}`,
      );
    }
    for (const skill of [
      ...PACK_SKILLS.backend,
      ...PACK_SKILLS.database,
      ...PACK_SKILLS.embedded,
    ]) {
      assert.equal(
        skillInstalled(targetDir, skill),
        false,
        `S2: must not install omitted-pack skill ${skill}`,
      );
    }
  } finally {
    rmSync(targetDir, { recursive: true, force: true });
  }
});

test('S3: install-cursor profile packs=[quality] → silent selective reconcile via setup path', () => {
  const installCursor = readFileSync(INSTALL_CURSOR, 'utf8');
  assert.match(
    installCursor,
    /node\s+"\$SRC_ABS\/cli\/lib\/setup\.mjs"/,
    'S3: silent reconcile must go through setup.mjs (shared with spacecraft setup)',
  );

  const targetDir = tempTarget();
  try {
    mkdirSync(path.join(targetDir, '.cursor'), { recursive: true });
    writeFileSync(
      profilePath(targetDir),
      `${JSON.stringify({ version: 1, packs: ['quality'] }, null, 2)}\n`,
      'utf8',
    );

    const res = runInstallCursor(targetDir, envWithoutPacks());
    assert.equal(
      res.status,
      0,
      `S3/E3: profile reconcile must exit 0\n${combined(res)}`,
    );
    assert.match(
      combined(res),
      /\(profile packs=quality\)/,
      'S3: must reconcile via setup.mjs profile source',
    );
    assert.doesNotMatch(
      combined(res),
      /all selectable packs \(legacy\)/,
      'S3: must not use legacy-all',
    );

    for (const skill of S3_SKILLS) {
      assert.equal(
        skillInstalled(targetDir, skill),
        true,
        `S3: missing quality skill ${skill}`,
      );
    }
    for (const rule of S3_RULES) {
      assert.equal(
        ruleInstalled(targetDir, rule),
        true,
        `S3: missing quality/always-on rule ${rule}`,
      );
    }
    for (const packId of S3_ABSENT_PACKS) {
      for (const skill of PACK_SKILLS[packId]) {
        assert.equal(
          skillInstalled(targetDir, skill),
          false,
          `S3: ${packId} skill ${skill} must be absent`,
        );
      }
      for (const rule of PACK_RULES[packId]) {
        assert.equal(
          ruleInstalled(targetDir, rule),
          false,
          `S3: ${packId} rule ${rule} must be absent`,
        );
      }
    }
  } finally {
    rmSync(targetDir, { recursive: true, force: true });
  }
});
