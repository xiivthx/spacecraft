#!/usr/bin/env node
/**
 * Project setup mode resolve: TTY / non-TTY / profile / --packs / SPACECRAFT_PACKS.
 * Fail-closed when no profile and non-TTY without packs (no silent all-packs).
 */
import { existsSync, readFileSync, readSync } from 'node:fs';
import path from 'node:path';
import { fileURLToPath, pathToFileURL } from 'node:url';

import { loadCatalog, validateProfile } from './packs.mjs';
import { installProjectSurface, parsePacksEnv } from './project-install.mjs';

const PROFILE_REL = path.join('.cursor', 'spacecraft-profile.json');
const DEFAULT_QUALITY = Object.freeze(['quality']);

export const SetupErrorCode = Object.freeze({
  NON_TTY_NO_PACKS: 'NON_TTY_NO_PACKS',
});

export class SetupError extends Error {
  /**
   * @param {string} message
   * @param {string} code
   */
  constructor(message, code) {
    super(message);
    this.name = 'SetupError';
    this.code = code;
  }
}

/**
 * Frozen S5/E5 interactive default: only quality pre-selected.
 * @param {{ packs: Array<{ id: string, status: string }> }} [_catalog]
 * @returns {string[]}
 */
export function defaultInteractiveSelection(_catalog) {
  return [...DEFAULT_QUALITY];
}

/**
 * Interactive choice list: all catalog packs; coming disabled; quality selected.
 * @param {{ packs: Array<{ id: string, status: string }> }} [catalog]
 * @returns {Array<{ id: string, status: string, enabled: boolean, selected: boolean }>}
 */
export function interactivePackChoices(catalog) {
  const cat = catalog ?? loadCatalog();
  return cat.packs.map((entry) => {
    const status = entry.status === 'coming' ? 'coming' : 'selectable';
    const enabled = status !== 'coming';
    const selected = enabled && entry.id === 'quality';
    return { id: entry.id, status, enabled, selected };
  });
}

/**
 * Read one line from stdin (sync; no new deps).
 * @returns {string}
 */
function readStdinLine() {
  const chunks = [];
  const buf = Buffer.alloc(1);
  for (;;) {
    let n;
    try {
      n = readSync(0, buf, 0, 1, null);
    } catch {
      break;
    }
    if (n === 0) break;
    if (buf[0] === 0x0a) break; // LF
    if (buf[0] === 0x0d) continue; // CR
    chunks.push(Buffer.from(buf));
  }
  return Buffer.concat(chunks).toString('utf8');
}

/**
 * Real TTY multi-select: toggle by number, Enter alone to confirm.
 * Coming packs shown disabled and never selectable.
 * @param {{
 *   choices?: Array<{ id: string, status: string, enabled: boolean, selected: boolean }>,
 *   catalog?: { packs: Array<{ id: string, status: string }> },
 *   defaults?: string[],
 * }} ctx
 * @returns {string[]}
 */
function defaultTtyAsk(ctx) {
  const choices = (ctx.choices ?? interactivePackChoices(ctx.catalog)).map((c) => ({
    ...c,
  }));
  const selectable = choices.filter((c) => c.enabled);

  /** @type {Set<string>} */
  const selected = new Set(
    (ctx.defaults ?? defaultInteractiveSelection(ctx.catalog)).filter((id) =>
      selectable.some((c) => c.id === id),
    ),
  );
  // Sync choice.selected with working set (defaults / quality pre-check).
  for (const c of choices) {
    if (c.enabled) c.selected = selected.has(c.id);
  }

  const printMenu = () => {
    console.log('Select packs (toggle by number; Enter alone to confirm):');
    let num = 0;
    for (const c of choices) {
      if (!c.enabled) {
        console.log(`  —    ${c.id} (coming — not selectable)`);
        continue;
      }
      num += 1;
      const mark = selected.has(c.id) ? 'x' : ' ';
      console.log(`  [${mark}] ${num}. ${c.id}`);
    }
    process.stdout.write('> ');
  };

  for (;;) {
    printMenu();
    const line = readStdinLine().trim();
    if (!line) {
      const ids = selectable.filter((c) => selected.has(c.id)).map((c) => c.id);
      return ids.length ? ids : [...DEFAULT_QUALITY];
    }

    // Comma/space-separated: numbers toggle; known selectable ids set selection.
    const tokens = line.split(/[,\s]+/).filter(Boolean);
    const onlyIds = tokens.every((t) => selectable.some((c) => c.id === t));
    if (onlyIds) {
      selected.clear();
      for (const id of tokens) selected.add(id);
      continue;
    }

    for (const token of tokens) {
      const n = Number.parseInt(token, 10);
      if (!Number.isFinite(n) || String(n) !== token) continue;
      if (n < 1 || n > selectable.length) continue;
      const id = selectable[n - 1].id;
      if (selected.has(id)) selected.delete(id);
      else selected.add(id);
    }
  }
}

/**
 * Factory for CLI interactive promptSelect.
 * @param {(ctx: object) => string[]} [ask] injectable; default = real TTY multi-select
 * @returns {(ctx: object) => string[]}
 */
export function createTtyPromptSelect(ask) {
  const askFn = typeof ask === 'function' ? ask : defaultTtyAsk;
  return (ctx) => {
    const choices = ctx.choices ?? interactivePackChoices(ctx.catalog);
    return askFn({ ...ctx, choices });
  };
}

/**
 * @param {string[]} args
 * @returns {{ packs: string | null, reconfigure: boolean, help: boolean }}
 */
export function parseSetupArgs(args) {
  let packs = null;
  let reconfigure = false;
  let help = false;
  for (let i = 0; i < args.length; i += 1) {
    const a = args[i];
    if (a === '--help' || a === '-h') {
      help = true;
      continue;
    }
    if (a === '--reconfigure') {
      reconfigure = true;
      continue;
    }
    if (a === '--packs') {
      const next = args[i + 1];
      if (next == null || next.startsWith('-')) {
        throw new SetupError(
          'spacecraft setup: --packs requires a comma-separated pack list',
          SetupErrorCode.NON_TTY_NO_PACKS,
        );
      }
      packs = next;
      i += 1;
      continue;
    }
    if (a.startsWith('--packs=')) {
      packs = a.slice('--packs='.length);
      continue;
    }
  }
  return { packs, reconfigure, help };
}

function printSetupHelp() {
  console.log('Usage: spacecraft setup [--packs a,b] [--reconfigure]');
  console.log('');
  console.log('Options:');
  console.log('  --packs a,b       Comma-separated selectable pack ids (writes profile)');
  console.log('  --reconfigure     Deliberate pack change when a profile already exists');
  console.log('  --help, -h        Show this help');
  console.log('');
  console.log('SPACECRAFT_PACKS=a,b  Same as --packs when the flag is absent');
}

/**
 * CLI entry: target = cwd; source = spacecraft repo root.
 * @param {string[]} args
 * @param {string} cwd
 * @param {string} [sourceDir]
 * @param {{ promptSelect?: (ctx: object) => string[], tty?: boolean }} [deps]
 * @returns {number}
 */
export function setupCmd(args, cwd, sourceDir, deps) {
  let parsed;
  try {
    parsed = parseSetupArgs(args);
  } catch (err) {
    const msg = err instanceof Error ? err.message : String(err);
    console.error(msg.startsWith('spacecraft setup:') ? msg : `spacecraft setup: ${msg}`);
    return 1;
  }

  if (parsed.help) {
    printSetupHelp();
    return 0;
  }

  const resolvedSource =
    sourceDir ??
    path.resolve(path.dirname(fileURLToPath(import.meta.url)), '../..');

  try {
    const result = runProjectSetup({
      targetDir: cwd,
      sourceDir: resolvedSource,
      tty: deps?.tty ?? Boolean(process.stdout.isTTY),
      packs: parsed.packs,
      packsEnv:
        process.env.SPACECRAFT_PACKS != null ? String(process.env.SPACECRAFT_PACKS) : null,
      reconfigure: parsed.reconfigure,
      promptSelect: deps?.promptSelect ?? createTtyPromptSelect(),
    });
    console.log(
      `  domain rules/skills (${result.source} packs=${result.packIds.join(',')}) -> ${path.join(cwd, '.cursor')}`,
    );
    return 0;
  } catch (err) {
    const msg = err instanceof Error ? err.message : String(err);
    console.error(`spacecraft setup: ${msg}`);
    return 1;
  }
}

/**
 * @param {string | string[] | null | undefined} packs
 * @returns {string[] | null}
 */
function normalizePacksFlag(packs) {
  if (packs == null) return null;
  if (Array.isArray(packs)) {
    const ids = packs.map((s) => String(s).trim()).filter(Boolean);
    return ids.length ? ids : null;
  }
  const raw = String(packs).trim();
  if (!raw) return null;
  return parsePacksEnv(raw);
}

/**
 * Explicit null/'' → absent. Do not fall back to process.env when null.
 * @param {string | null | undefined} packsEnv
 * @returns {string[] | null}
 */
function normalizePacksEnv(packsEnv) {
  if (packsEnv == null) return null;
  const trimmed = String(packsEnv).trim();
  if (!trimmed) return null;
  return parsePacksEnv(trimmed);
}

/**
 * @param {string} targetDir
 * @returns {string[]}
 */
function readProfilePacks(targetDir) {
  const profilePath = path.join(targetDir, PROFILE_REL);
  const raw = JSON.parse(readFileSync(profilePath, 'utf8'));
  if (!raw || !Array.isArray(raw.packs)) {
    throw new Error(`Invalid profile at ${profilePath}: expected { packs: [] }`);
  }
  return [...raw.packs];
}

/**
 * Resolve setup action without installing.
 * @param {{
 *   targetDir: string,
 *   tty: boolean,
 *   packs?: string | string[] | null,
 *   packsEnv?: string | null,
 *   reconfigure?: boolean,
 * }} opts
 */
export function resolveSetupMode({
  targetDir,
  tty,
  packs = null,
  packsEnv = null,
  reconfigure = false,
}) {
  const fromFlag = normalizePacksFlag(packs);
  if (fromFlag) {
    return {
      action: 'reconcile',
      packIds: fromFlag,
      writeProfile: true,
      source: 'packs',
      interactive: false,
    };
  }

  const fromEnv = normalizePacksEnv(packsEnv);
  if (fromEnv) {
    return {
      action: 'reconcile',
      packIds: fromEnv,
      writeProfile: true,
      source: 'env',
      interactive: false,
    };
  }

  const profilePath = path.join(targetDir, PROFILE_REL);
  const hasProfile = existsSync(profilePath);

  if (hasProfile && !reconfigure) {
    return {
      action: 'reconcile',
      packIds: readProfilePacks(targetDir),
      writeProfile: false,
      source: 'profile',
      interactive: false,
    };
  }

  if (tty) {
    return {
      action: 'interactive',
      writeProfile: true,
      source: 'interactive',
      interactive: true,
    };
  }

  return {
    action: 'fail',
    code: SetupErrorCode.NON_TTY_NO_PACKS,
    writeProfile: false,
    source: 'none',
    interactive: false,
  };
}

/**
 * Resolve mode then selective install (or throw SetupError).
 * Never uses project-install legacy-all: pack list always comes from resolve.
 * @param {{
 *   targetDir: string,
 *   sourceDir: string,
 *   tty: boolean,
 *   packs?: string | string[] | null,
 *   packsEnv?: string | null,
 *   reconfigure?: boolean,
 *   catalogPath?: string,
 *   promptSelect?: (ctx: object) => string[],
 * }} opts
 */
export function runProjectSetup({
  targetDir,
  sourceDir,
  tty,
  packs = null,
  packsEnv = null,
  reconfigure = false,
  catalogPath,
  promptSelect,
}) {
  const mode = resolveSetupMode({
    targetDir,
    tty,
    packs,
    packsEnv,
    reconfigure,
  });

  if (mode.action === 'fail') {
    throw new SetupError(
      'Non-interactive setup requires --packs or SPACECRAFT_PACKS when no spacecraft-profile.json exists',
      mode.code ?? SetupErrorCode.NON_TTY_NO_PACKS,
    );
  }

  const resolvedCatalogPath =
    catalogPath ?? path.join(sourceDir, '.cursor', 'spacecraft-packs.json');
  const catalog = loadCatalog(resolvedCatalogPath);

  let packIds = mode.packIds ? [...mode.packIds] : [];
  let source = mode.source;
  let prompted = false;
  let writeProfile = mode.writeProfile;

  if (mode.action === 'interactive') {
    const select =
      typeof promptSelect === 'function' ? promptSelect : createTtyPromptSelect();
    packIds = select({
      targetDir,
      sourceDir,
      catalog,
      catalogPath: resolvedCatalogPath,
      choices: interactivePackChoices(catalog),
      defaults: defaultInteractiveSelection(catalog),
    });
    prompted = true;
    source = 'interactive';
    writeProfile = true;
  }

  validateProfile({ version: 1, packs: packIds }, catalog);

  // Explicit packIds: setup owns resolve; never hit legacy-all in project-install.
  installProjectSurface(targetDir, sourceDir, {
    catalogPath: resolvedCatalogPath,
    packIds,
    writeProfile,
    source,
  });

  return {
    packIds,
    source,
    prompted,
    profilePath: path.join(targetDir, PROFILE_REL),
  };
}

function isMain() {
  const entry = process.argv[1];
  if (!entry) return false;
  return import.meta.url === pathToFileURL(path.resolve(entry)).href;
}

if (isMain()) {
  const target = process.argv[2];
  const source = process.argv[3];
  if (!target || !source || process.argv.length !== 4) {
    console.error('usage: setup.mjs <target-project-dir> <source-repo-dir>');
    process.exit(1);
  }
  try {
    const result = runProjectSetup({
      targetDir: path.resolve(target),
      sourceDir: path.resolve(source),
      tty: Boolean(process.stdout.isTTY),
      packs: null,
      packsEnv:
        process.env.SPACECRAFT_PACKS != null ? String(process.env.SPACECRAFT_PACKS) : null,
      promptSelect: createTtyPromptSelect(),
    });
    console.log(
      `  domain rules/skills (${result.source} packs=${result.packIds.join(',')}) -> ${path.resolve(target)}/.cursor`,
    );
  } catch (err) {
    const msg = err instanceof Error ? err.message : String(err);
    console.error(`setup: ${msg}`);
    process.exit(1);
  }
}
