/**
 * Setup CLI + interactive defaults (T4 / M99A9D1B) — RED.
 *
 * Expected CLI surface (coder GREEN — do not invent alternate flags):
 *
 *   spacecraft setup [--packs a,b] [--reconfigure]
 *     - Listed in help as `spacecraft setup` (see help-dispatch.test.mjs)
 *     - Dispatches to setupCmd → runProjectSetup (target = cwd; source = spacecraft repo root)
 *     - --packs a,b | --packs=a,b  → selectable pack ids (comma-separated); writes profile; reconcile/prune
 *     - --reconfigure              → deliberate pack-change when profile exists
 *                                    (with --packs: apply list; without: interactive)
 *     - SPACECRAFT_PACKS           → same as --packs when flag absent
 *     - spacecraft setup --help|-h → documents --packs and --reconfigure
 *
 * Expected new exports from cli/lib/setup.mjs (interactive seam; inject/stub prompts):
 *
 *   defaultInteractiveSelection(catalog?) → string[]
 *     - Frozen S5/E5: exactly ['quality'] (only quality pre-selected among selectable)
 *
 *   interactivePackChoices(catalog?) → Array<{
 *     id: string,
 *     status: 'selectable' | 'coming',
 *     enabled: boolean,   // false iff status === 'coming'
 *     selected: boolean,  // true iff id === 'quality' && selectable (default UI check)
 *   }>
 *     - Includes all catalog packs (selectable + coming)
 *     - Coming visible but enabled:false; never written into profile
 *
 * Fix-pass (review: TTY setup never prompts) — expected API for coder:
 *
 *   createTtyPromptSelect(ask?) → (ctx) => string[]
 *     - Factory for the CLI interactive promptSelect
 *     - ask(ctx) is injectable (tests stub); production default = real TTY multi-select
 *     - Must forward ctx.choices (interactivePackChoices shape) into ask — Must not
 *       ignore choices and return defaultInteractiveSelection(catalog) alone
 *     - Must return ask's selected selectable pack ids
 *     - S5 unchanged: defaults / choice.selected still quality-only
 *     - S6 unchanged: coming enabled:false; never returned as selected
 *
 *   setupCmd(args, cwd, sourceDir, deps?: { promptSelect?, tty? }) → number
 *     - Default promptSelect MUST be createTtyPromptSelect() — Must not hardcode
 *       promptSelect: ({ catalog }) => defaultInteractiveSelection(catalog)
 *     - deps.promptSelect / deps.tty override for tests (force TTY + stub ask)
 *
 * Oracles: approved-scenarios S4–S6 / Edge E4–E6; review.json TTY prompt finding.
 * Pbt skipped: no project pbt tool
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

import { loadCatalog } from '../lib/packs.mjs';
import { runProjectSetup } from '../lib/setup.mjs';

const repoRoot = path.resolve(path.dirname(fileURLToPath(import.meta.url)), '../..');
const entryPath = path.join(repoRoot, 'cli', 'spacecraft.mjs');
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

/** S4 frozen: reconfigure to quality only after frontend+quality. */
const S4_BEFORE_PACKS = ['frontend', 'quality'];
const S4_AFTER_PACKS = ['quality'];
const S4_REMOVED_SKILLS = PACK_SKILLS.frontend;
const S4_REMOVED_RULES = PACK_RULES.frontend;
const S4_KEPT_SKILLS = PACK_SKILLS.quality;
const S4_KEPT_RULES = [...PACK_RULES.quality, ALWAYS_ON_RULE];

/** S5 frozen: default checked selectable set. */
const S5_DEFAULT_SELECTION = ['quality'];

/** S6 frozen: coming pack ids (never in profile). */
const S6_COMING = ['iot', 'fpga', 'pcb', 'management'];
const S6_SELECTABLE = ['frontend', 'backend', 'database', 'embedded', 'quality'];

function sorted(arr) {
  return [...arr].sort();
}

function tempTarget() {
  return mkdtempSync(path.join(os.tmpdir(), 'spacecraft-setup-cli-'));
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

/**
 * Seed target with frontend+quality via library (T3 path already GREEN).
 * CLI under test owns the reconfigure step.
 */
function seedFrontendQuality(targetDir) {
  runProjectSetup({
    targetDir,
    sourceDir: repoRoot,
    tty: false,
    packs: S4_BEFORE_PACKS,
    packsEnv: null,
    catalogPath: CATALOG_PATH,
    promptSelect: () => {
      throw new Error('seed must not prompt');
    },
  });
  const profile = readProfile(targetDir);
  assert.deepEqual(
    sorted(profile.packs),
    sorted(S4_BEFORE_PACKS),
    'seed: profile must be frontend,quality before S4 CLI',
  );
  for (const skill of S4_REMOVED_SKILLS) {
    assert.equal(skillInstalled(targetDir, skill), true, `seed missing ${skill}`);
  }
}

function runSetupCli(targetDir, args, env = {}) {
  return spawnSync(process.execPath, [entryPath, 'setup', ...args], {
    encoding: 'utf8',
    cwd: targetDir,
    env: { ...process.env, ...env },
  });
}

test('S4: spacecraft setup --reconfigure --packs quality prunes frontend and keeps quality', () => {
  const targetDir = tempTarget();
  try {
    seedFrontendQuality(targetDir);

    const result = runSetupCli(targetDir, ['--reconfigure', '--packs', 'quality']);
    const out = `${result.stdout ?? ''}${result.stderr ?? ''}`;
    assert.equal(
      result.status,
      0,
      `S4: setup --reconfigure --packs quality must exit 0\n${out}`,
    );

    assert.ok(existsSync(profilePath(targetDir)), 'S4: profile must remain');
    const profile = readProfile(targetDir);
    assert.equal(profile.version, 1);
    assert.deepEqual(
      sorted(profile.packs),
      sorted(S4_AFTER_PACKS),
      'S4: profile packs must be exactly ["quality"]',
    );

    for (const skill of S4_REMOVED_SKILLS) {
      assert.equal(
        skillInstalled(targetDir, skill),
        false,
        `S4: frontend skill ${skill} must be pruned`,
      );
    }
    for (const rule of S4_REMOVED_RULES) {
      assert.equal(
        ruleInstalled(targetDir, rule),
        false,
        `S4: frontend rule ${rule} must be pruned`,
      );
    }
    for (const skill of S4_KEPT_SKILLS) {
      assert.equal(
        skillInstalled(targetDir, skill),
        true,
        `S4: quality skill ${skill} must remain`,
      );
    }
    for (const rule of S4_KEPT_RULES) {
      assert.equal(
        ruleInstalled(targetDir, rule),
        true,
        `S4: quality/always-on rule ${rule} must remain`,
      );
    }
  } finally {
    rmSync(targetDir, { recursive: true, force: true });
  }
});

test('S5/S6: interactive default is quality only; coming packs disabled and not writable', async () => {
  const setupMod = await import('../lib/setup.mjs');
  assert.equal(
    typeof setupMod.defaultInteractiveSelection,
    'function',
    'setup.mjs must export defaultInteractiveSelection(catalog?)',
  );
  assert.equal(
    typeof setupMod.interactivePackChoices,
    'function',
    'setup.mjs must export interactivePackChoices(catalog?)',
  );

  const catalog = loadCatalog(CATALOG_PATH);

  const defaults = setupMod.defaultInteractiveSelection(catalog);
  assert.deepEqual(
    sorted(defaults),
    sorted(S5_DEFAULT_SELECTION),
    'S5/E5: defaultInteractiveSelection must be only quality',
  );

  const choices = setupMod.interactivePackChoices(catalog);
  assert.ok(Array.isArray(choices), 'interactivePackChoices must return an array');

  const byId = new Map(choices.map((c) => [c.id, c]));
  for (const id of S6_SELECTABLE) {
    const choice = byId.get(id);
    assert.ok(choice, `S6: selectable pack ${id} must appear in choices`);
    assert.equal(choice.status, 'selectable', `${id} status`);
    assert.equal(choice.enabled, true, `${id} must be enabled`);
    assert.equal(
      choice.selected,
      id === 'quality',
      `S5: only quality pre-selected (got selected=${choice.selected} for ${id})`,
    );
  }
  for (const id of S6_COMING) {
    const choice = byId.get(id);
    assert.ok(choice, `S6: coming pack ${id} must be visible in choices`);
    assert.equal(choice.status, 'coming', `${id} status`);
    assert.equal(choice.enabled, false, `S6: coming pack ${id} must be disabled`);
    assert.equal(choice.selected, false, `S6: coming pack ${id} must not be selected`);
  }

  // Coming cannot be written into profile (CLI --packs=iot / interactive select).
  const targetDir = tempTarget();
  try {
    mkdirSync(path.join(targetDir, '.cursor'), { recursive: true });
    const result = runSetupCli(targetDir, ['--packs', 'iot']);
    const out = `${result.stdout ?? ''}${result.stderr ?? ''}`;
    assert.doesNotMatch(
      out,
      /unknown command/i,
      `S6: setup must dispatch (not unknown) before rejecting coming packs\n${out}`,
    );
    assert.notEqual(
      result.status,
      0,
      `S6: --packs iot must exit non-zero\n${out}`,
    );
    if (existsSync(profilePath(targetDir))) {
      const profile = readProfile(targetDir);
      for (const id of S6_COMING) {
        assert.equal(
          profile.packs.includes(id),
          false,
          `S6: profile must never contain coming id ${id}`,
        );
      }
    }
    for (const id of S6_COMING) {
      assert.equal(
        existsSync(path.join(targetDir, '.cursor', 'skills', id)),
        false,
        `S6: must not invent skill dir for coming pack ${id}`,
      );
    }
  } finally {
    rmSync(targetDir, { recursive: true, force: true });
  }
});

/**
 * Review finding: setupCmd must not silently hardcode quality-only.
 * createTtyPromptSelect(ask) must invoke ask with interactivePackChoices and
 * return the user-selected packs (not defaultInteractiveSelection alone).
 */
test('fix: createTtyPromptSelect forwards interactivePackChoices to ask and returns user packs', async () => {
  const setupMod = await import('../lib/setup.mjs');
  assert.equal(
    typeof setupMod.createTtyPromptSelect,
    'function',
    'setup.mjs must export createTtyPromptSelect(ask?)',
  );

  const catalog = loadCatalog(CATALOG_PATH);
  const expectedChoices = setupMod.interactivePackChoices(catalog);
  const defaults = setupMod.defaultInteractiveSelection(catalog);
  assert.deepEqual(
    sorted(defaults),
    sorted(S5_DEFAULT_SELECTION),
    'S5 must stay quality-only defaults (do not weaken)',
  );

  /** User confirms frontend + quality (not silent quality-only). */
  const USER_SELECTED = ['frontend', 'quality'];
  let askCalls = 0;
  /** @type {object | null} */
  let seenCtx = null;
  const ask = (ctx) => {
    askCalls += 1;
    seenCtx = ctx;
    return USER_SELECTED;
  };

  const promptSelect = setupMod.createTtyPromptSelect(ask);
  assert.equal(typeof promptSelect, 'function', 'createTtyPromptSelect must return promptSelect');

  const selected = promptSelect({
    catalog,
    choices: expectedChoices,
    defaults,
  });

  assert.equal(askCalls, 1, 'createTtyPromptSelect must invoke ask exactly once');
  assert.ok(seenCtx, 'ask must receive a context object');
  assert.ok(Array.isArray(seenCtx.choices), 'ask ctx must include choices');
  assert.deepEqual(
    seenCtx.choices.map((c) => c.id).sort(),
    expectedChoices.map((c) => c.id).sort(),
    'ask must receive interactivePackChoices ids',
  );

  const byId = new Map(seenCtx.choices.map((c) => [c.id, c]));
  for (const id of S6_SELECTABLE) {
    const choice = byId.get(id);
    assert.ok(choice, `selectable ${id} must be in ask choices`);
    assert.equal(choice.enabled, true, `${id} must be enabled`);
    assert.equal(
      choice.selected,
      id === 'quality',
      `S5: only quality pre-selected in ask choices (got selected=${choice.selected} for ${id})`,
    );
  }
  for (const id of S6_COMING) {
    const choice = byId.get(id);
    assert.ok(choice, `coming ${id} must be visible in ask choices`);
    assert.equal(choice.enabled, false, `S6: coming ${id} must stay disabled`);
    assert.equal(choice.selected, false, `S6: coming ${id} must not be selected`);
  }

  assert.deepEqual(
    sorted(selected),
    sorted(USER_SELECTED),
    'createTtyPromptSelect must return ask selection (frontend+quality), not silent quality-only',
  );
});

/**
 * setupCmd TTY no-profile path must invoke injectable promptSelect and honor its
 * return value. Proves CLI wiring is not hardcoding defaultInteractiveSelection.
 */
test('fix: setupCmd TTY interactive invokes promptSelect and writes user-selected packs', async () => {
  const setupMod = await import('../lib/setup.mjs');
  assert.equal(
    typeof setupMod.setupCmd,
    'function',
    'setup.mjs must export setupCmd',
  );
  assert.equal(
    typeof setupMod.createTtyPromptSelect,
    'function',
    'setup.mjs must export createTtyPromptSelect for CLI wiring',
  );

  const USER_SELECTED = ['frontend', 'quality'];
  const targetDir = tempTarget();
  try {
    let promptCalls = 0;
    /** @type {object | null} */
    let seenCtx = null;

    const code = setupMod.setupCmd([], targetDir, repoRoot, {
      tty: true,
      promptSelect: setupMod.createTtyPromptSelect((ctx) => {
        promptCalls += 1;
        seenCtx = ctx;
        return USER_SELECTED;
      }),
    });

    const profileExists = existsSync(profilePath(targetDir));
    assert.equal(
      promptCalls,
      1,
      'setupCmd TTY no-profile must invoke promptSelect once (not hardcode quality-only)',
    );
    assert.ok(seenCtx?.choices, 'promptSelect ctx must include choices from interactivePackChoices');
    assert.equal(code, 0, 'setupCmd must exit 0 after interactive confirm');
    assert.equal(profileExists, true, 'setupCmd must write profile after prompt');
    const profile = readProfile(targetDir);
    assert.deepEqual(
      sorted(profile.packs),
      sorted(USER_SELECTED),
      'profile packs must be user selection from promptSelect, not silent ["quality"]',
    );
    for (const id of S6_COMING) {
      assert.equal(
        profile.packs.includes(id),
        false,
        `S6: profile must never contain coming id ${id}`,
      );
    }
  } finally {
    rmSync(targetDir, { recursive: true, force: true });
  }
});

test('T4-c: setup without --reconfigure documents deliberate pack-change via --reconfigure', () => {
  // Usage/help must document --reconfigure as the deliberate pack-change path.
  const help = spawnSync(process.execPath, [entryPath, 'setup', '--help'], {
    encoding: 'utf8',
    cwd: repoRoot,
  });
  const helpOut = `${help.stdout ?? ''}${help.stderr ?? ''}`;
  assert.equal(
    help.status,
    0,
    `setup --help must exit 0\n${helpOut}`,
  );
  assert.match(
    helpOut,
    /--reconfigure/,
    `setup --help must document --reconfigure\n${helpOut}`,
  );
  assert.match(
    helpOut,
    /--packs/,
    `setup --help must document --packs\n${helpOut}`,
  );

  // Profile present + setup without --reconfigure → silent reconcile (does not block);
  // deliberate change remains available via --reconfigure (proven by S4).
  const targetDir = tempTarget();
  try {
    mkdirSync(path.join(targetDir, '.cursor'), { recursive: true });
    writeFileSync(
      profilePath(targetDir),
      `${JSON.stringify({ version: 1, packs: ['quality'] }, null, 2)}\n`,
      'utf8',
    );

    const result = runSetupCli(targetDir, []);
    const out = `${result.stdout ?? ''}${result.stderr ?? ''}`;
    assert.equal(
      result.status,
      0,
      `profile present + setup (no --reconfigure) must silent-reconcile exit 0\n${out}`,
    );
    const profile = readProfile(targetDir);
    assert.deepEqual(
      sorted(profile.packs),
      sorted(['quality']),
      'silent setup must keep existing profile packs',
    );
  } finally {
    rmSync(targetDir, { recursive: true, force: true });
  }
});
