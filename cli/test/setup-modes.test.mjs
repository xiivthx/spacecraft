/**
 * Setup mode matrix (T3 / M99A9D1B) — RED against cli/lib/setup.mjs.
 *
 * Expected public API (coder GREEN — do not invent alternate seams):
 *
 *   SetupError extends Error
 *     .code: string  — typed fail codes (deep-assert ErrorCode)
 *   SetupErrorCode.NON_TTY_NO_PACKS === 'NON_TTY_NO_PACKS'
 *     — S1/E1: no profile + non-TTY + neither --packs nor SPACECRAFT_PACKS
 *
 *   resolveSetupMode({
 *     targetDir: string,
 *     tty: boolean,                 // inject; do not read process.stdout.isTTY alone in tests
 *     packs?: string | string[] | null,  // --packs a,b
 *     packsEnv?: string | null,     // SPACECRAFT_PACKS; null/'' = absent (do not fall back to process.env when explicitly null)
 *     reconfigure?: boolean,
 *   }) → {
 *     action: 'fail' | 'interactive' | 'reconcile',
 *     code?: typeof SetupErrorCode[keyof typeof SetupErrorCode],
 *     packIds?: string[],
 *     writeProfile: boolean,
 *     source: 'none' | 'packs' | 'env' | 'profile' | 'interactive',
 *     interactive: boolean,         // true only when action === 'interactive'
 *   }
 *
 *   runProjectSetup({
 *     targetDir: string,
 *     sourceDir: string,            // spacecraft repo (skills/rules source)
 *     tty: boolean,
 *     packs?: string | string[] | null,
 *     packsEnv?: string | null,
 *     reconfigure?: boolean,
 *     catalogPath?: string,
 *     promptSelect?: (ctx) => string[],  // interactive only; Must not run on silent profile reconcile
 *   }) → {
 *     packIds: string[],
 *     source: 'packs' | 'env' | 'profile' | 'interactive',
 *     prompted: boolean,            // false on S3 silent reconcile
 *     profilePath: string,          // …/.cursor/spacecraft-profile.json
 *   }
 *   On S1: throws SetupError with .code === SetupErrorCode.NON_TTY_NO_PACKS
 *   and Must not install all selectable-domain skill trees (no silent all-packs).
 *
 * Mode matrix (design-contract): kill legacy-all in project-install when setup owns resolve.
 * Oracles: approved-scenarios S1–S3 / Edge E1–E3.
 *
 * Pbt skipped: no project pbt tool
 */
import assert from 'node:assert/strict';
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

import {
  SetupError,
  SetupErrorCode,
  runProjectSetup,
} from '../lib/setup.mjs';

const repoRoot = path.resolve(path.dirname(fileURLToPath(import.meta.url)), '../..');
const CATALOG_PATH = path.join(repoRoot, '.cursor', 'spacecraft-packs.json');
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
  fpga: ['sc-rtl', 'sc-rtl-verify'],
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
  fpga: ['700-rtl.mdc', '710-rtl-timing.mdc', '720-rtl-verify.mdc'],
});

/** S2 frozen: SPACECRAFT_PACKS=frontend,quality */
const S2_PACKS = ['frontend', 'quality'];
const S2_SKILLS = [...PACK_SKILLS.frontend, ...PACK_SKILLS.quality];
const S2_RULES = [...PACK_RULES.frontend, ...PACK_RULES.quality, ALWAYS_ON_RULE];

/** S3 frozen: profile packs=["quality"] */
const S3_PACKS = ['quality'];
const S3_SKILLS = [...PACK_SKILLS.quality];
const S3_RULES = [...PACK_RULES.quality, ALWAYS_ON_RULE];
const S3_ABSENT_PACKS = ['frontend', 'backend', 'database', 'embedded'];

const ALL_DOMAIN_SKILLS = Object.freeze([
  ...new Set(Object.values(PACK_SKILLS).flat()),
]);

function sorted(arr) {
  return [...arr].sort();
}

function tempTarget() {
  return mkdtempSync(path.join(os.tmpdir(), 'spacecraft-setup-modes-'));
}

function profilePath(targetDir) {
  return path.join(targetDir, PROFILE_REL);
}

function readProfile(targetDir) {
  return JSON.parse(readFileSync(profilePath(targetDir), 'utf8'));
}

function skillInstalled(targetDir, name) {
  return existsSync(path.join(targetDir, '.cursor', 'skills', name, 'SKILL.md'));
}

function ruleInstalled(targetDir, name) {
  return existsSync(path.join(targetDir, '.cursor', 'rules', name));
}

function assertNoAllPacksDomainInstall(targetDir) {
  const present = ALL_DOMAIN_SKILLS.filter((name) => skillInstalled(targetDir, name));
  assert.deepEqual(
    present,
    [],
    `S1: must not silent-install all-packs domain skills (found: ${present.join(', ')})`,
  );
}

test('S1: no profile + non-TTY + no packs → SetupError NON_TTY_NO_PACKS and no all-packs install', () => {
  const targetDir = tempTarget();
  try {
    let thrown = null;
    try {
      runProjectSetup({
        targetDir,
        sourceDir: repoRoot,
        tty: false,
        packs: null,
        packsEnv: null,
        catalogPath: CATALOG_PATH,
        promptSelect: () => {
          throw new Error('S1 must not prompt');
        },
      });
    } catch (err) {
      thrown = err;
    }

    assert.ok(thrown, 'S1: runProjectSetup must fail (non-zero path)');
    assert.ok(
      thrown instanceof SetupError,
      `S1: expected SetupError instance, got ${thrown?.constructor?.name}`,
    );
    assert.equal(
      thrown.code,
      SetupErrorCode.NON_TTY_NO_PACKS,
      'S1/E1: ErrorCode must be NON_TTY_NO_PACKS',
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

test('S2: no profile + non-TTY + SPACECRAFT_PACKS=frontend,quality → write profile and install only those packs', () => {
  const targetDir = tempTarget();
  try {
    const result = runProjectSetup({
      targetDir,
      sourceDir: repoRoot,
      tty: false,
      packs: null,
      packsEnv: 'frontend,quality',
      catalogPath: CATALOG_PATH,
      promptSelect: () => {
        throw new Error('S2 non-TTY packs path must not prompt');
      },
    });

    assert.equal(result.prompted, false, 'S2: must not interactive-prompt');
    assert.ok(
      existsSync(profilePath(targetDir)),
      'S2: must write .cursor/spacecraft-profile.json',
    );
    const profile = readProfile(targetDir);
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
    for (const rule of [...PACK_RULES.database, ...PACK_RULES.embedded]) {
      assert.equal(
        ruleInstalled(targetDir, rule),
        false,
        `S2: must not install omitted-pack rule ${rule}`,
      );
    }
  } finally {
    rmSync(targetDir, { recursive: true, force: true });
  }
});

test('S3: profile present → silent reconcile from profile packs with no interactive prompt', () => {
  const targetDir = tempTarget();
  try {
    mkdirSync(path.join(targetDir, '.cursor'), { recursive: true });
    // Frozen S3 profile literal
    writeFileSync(
      profilePath(targetDir),
      `${JSON.stringify({ version: 1, packs: ['quality'] }, null, 2)}\n`,
      'utf8',
    );

    let promptCalls = 0;
    const result = runProjectSetup({
      targetDir,
      sourceDir: repoRoot,
      tty: true,
      packs: null,
      packsEnv: null,
      catalogPath: CATALOG_PATH,
      promptSelect: () => {
        promptCalls += 1;
        throw new Error('S3 silent reconcile must not call promptSelect');
      },
    });

    assert.equal(promptCalls, 0, 'S3: promptSelect must not be called');
    assert.equal(result.prompted, false, 'S3: prompted must be false');
    assert.equal(result.source, 'profile', 'S3: source must be profile');
    assert.deepEqual(sorted(result.packIds), sorted(S3_PACKS));

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
