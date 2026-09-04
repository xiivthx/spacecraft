import { readFileSync } from 'node:fs';
import path from 'node:path';
import { fileURLToPath } from 'node:url';

const ALWAYS_ON_RULE = '010-hard-contract.mdc';

const DEFAULT_CATALOG_PATH = path.resolve(
  path.dirname(fileURLToPath(import.meta.url)),
  '../../.cursor/spacecraft-packs.json',
);

/**
 * Load pack catalog SoT from disk.
 * @param {string} [catalogPath]
 * @returns {{ packs: Array<{ id: string, status: string, skills?: string[], rules?: string[], mcp?: string }> }}
 */
export function loadCatalog(catalogPath = DEFAULT_CATALOG_PATH) {
  const raw = readFileSync(catalogPath, 'utf8');
  const data = JSON.parse(raw);
  if (!data || !Array.isArray(data.packs)) {
    throw new Error(`Invalid pack catalog at ${catalogPath}: expected { packs: [] }`);
  }
  return data;
}

function packById(catalog) {
  const map = new Map();
  for (const entry of catalog.packs) {
    map.set(entry.id, entry);
  }
  return map;
}

/**
 * Resolve .cursor dir that owns the catalog (fragments are relative to it).
 * @param {string} [catalogPath]
 * @returns {string}
 */
export function catalogCursorDir(catalogPath = DEFAULT_CATALOG_PATH) {
  return path.dirname(path.resolve(catalogPath));
}

/**
 * Union skills/rules/mcp fragments for selectable pack ids; always includes hard-contract rule.
 * @param {string[]} packIds
 * @param {{ packs: Array<{ id: string, status: string, skills?: string[], rules?: string[], mcp?: string }> }} [catalog]
 * @param {{ cursorDir?: string, catalogPath?: string }} [options]
 * @returns {{ skills: string[], rules: string[], mcp: string[] }}
 */
export function expandPacks(packIds, catalog, options = {}) {
  const cat = catalog ?? loadCatalog();
  const byId = packById(cat);
  const skills = new Set();
  const rules = new Set([ALWAYS_ON_RULE]);
  const mcp = new Set();
  const cursorDir =
    options.cursorDir ??
    catalogCursorDir(options.catalogPath ?? DEFAULT_CATALOG_PATH);

  for (const id of packIds) {
    const entry = byId.get(id);
    if (!entry) continue;
    for (const s of entry.skills ?? []) skills.add(s);
    for (const r of entry.rules ?? []) rules.add(r);
    if (entry.mcp) {
      mcp.add(path.resolve(cursorDir, entry.mcp));
    }
  }

  return { skills: [...skills], rules: [...rules], mcp: [...mcp] };
}

/**
 * Validate spacecraft-profile shape; only unique selectable pack ids allowed.
 * @param {unknown} profile
 * @param {{ packs: Array<{ id: string, status: string }> }} [catalog]
 * @returns {{ version: number, packs: string[] }}
 */
export function validateProfile(profile, catalog) {
  const cat = catalog ?? loadCatalog();
  const byId = packById(cat);

  if (!profile || typeof profile !== 'object' || Array.isArray(profile)) {
    throw new Error('Invalid profile: expected object with version and packs');
  }

  const { version, packs } = /** @type {{ version?: unknown, packs?: unknown }} */ (profile);

  if (version !== 1) {
    throw new Error('Invalid profile: version must be 1');
  }

  if (!Array.isArray(packs)) {
    throw new Error('Invalid profile: packs must be an array');
  }

  const seen = new Set();
  for (const id of packs) {
    if (typeof id !== 'string' || !id) {
      throw new Error('Invalid profile: pack ids must be non-empty strings');
    }
    if (seen.has(id)) {
      throw new Error(`Invalid profile: duplicate pack id "${id}"`);
    }
    seen.add(id);

    const entry = byId.get(id);
    if (!entry) {
      throw new Error(`Invalid profile: unknown pack id "${id}"`);
    }
    if (entry.status === 'coming') {
      throw new Error(`Invalid profile: coming pack id "${id}" cannot be selected`);
    }
    if (entry.status !== 'selectable') {
      throw new Error(`Invalid profile: pack id "${id}" is not selectable`);
    }
  }

  return { version: 1, packs: [...packs] };
}
