#!/usr/bin/env node
/**
 * Project-layer skills/rules reconcile: pack-selective or legacy all-selectable.
 *
 * Usage: node cli/lib/project-install.mjs <target-project-dir> <source-repo-dir>
 *
 * Pack resolution (T2):
 *   1. SPACECRAFT_PACKS=a,b → validate, write profile, install union
 *   2. else .cursor/spacecraft-profile.json → validate, install union
 *   3. else → all selectable catalog packs (legacy; no profile write; T3 fail-closes)
 *
 * Never installs lean-core skills, soft User-layer rules, or agents.
 */
import {
  copyFileSync,
  cpSync,
  existsSync,
  mkdirSync,
  readdirSync,
  readFileSync,
  rmSync,
  writeFileSync,
} from 'node:fs';
import path from 'node:path';
import { pathToFileURL } from 'node:url';

import {
  catalogCursorDir,
  expandPacks,
  loadCatalog,
  validateProfile,
} from './packs.mjs';

/** Keep in sync with scripts/global-sync.sh / install-cursor.sh LEAN_SKILLS. */
export const LEAN_SKILLS = Object.freeze([
  'sc-discuss',
  'sc-run',
  'sc-ship',
  'sc-quick',
  'sc-debug',
  'sc-mission',
  'sc-planning',
  'sc-tdd',
  'sc-verification',
  'sc-judge',
  'sc-clarify',
  'sc-git',
  'sc-search',
  'sc-storm',
  'sc-writer',
]);

/** Soft User-layer depth rules — never project. */
export const USER_LAYER_RULES = Object.freeze([
  '000-spacecraft.mdc',
  '026-intent-coach.mdc',
  '027-th-en-hil.mdc',
  '050-style.mdc',
  '100-conventions.mdc',
  '200-workflow.mdc',
]);

const PROFILE_NAME = 'spacecraft-profile.json';
const LEAN_SET = new Set(LEAN_SKILLS);
const USER_RULE_SET = new Set(USER_LAYER_RULES);

/**
 * Parse SPACECRAFT_PACKS / comma-separated pack list.
 * @param {string} raw
 * @returns {string[]}
 */
export function parsePacksEnv(raw) {
  return String(raw)
    .split(',')
    .map((s) => s.trim())
    .filter(Boolean);
}

/**
 * Read mcpServers map from a pack MCP fragment file.
 * @param {string} fragmentPath
 * @returns {Record<string, unknown>}
 */
export function loadMcpFragmentServers(fragmentPath) {
  if (!existsSync(fragmentPath)) {
    throw new Error(`missing MCP fragment ${fragmentPath}`);
  }
  const raw = readFileSync(fragmentPath, 'utf8').trim();
  if (!raw) return {};
  const data = JSON.parse(raw);
  if (!data || typeof data !== 'object' || Array.isArray(data)) {
    throw new Error(`${fragmentPath}: expected a JSON object`);
  }
  const servers = /** @type {{ mcpServers?: unknown }} */ (data).mcpServers;
  if (servers == null) return {};
  if (typeof servers !== 'object' || Array.isArray(servers)) {
    throw new Error(`${fragmentPath}: mcpServers must be an object`);
  }
  return /** @type {Record<string, unknown>} */ (servers);
}

/**
 * Union of skill dirs / rule basenames / MCP server names managed by selectable packs.
 * @param {{ packs: Array<{ id: string, status: string, skills?: string[], rules?: string[], mcp?: string }> }} catalog
 * @param {{ cursorDir?: string, catalogPath?: string }} [options]
 * @returns {{ skills: Set<string>, rules: Set<string>, mcpServers: Set<string> }}
 */
export function managedCatalogPaths(catalog, options = {}) {
  const skills = new Set();
  const rules = new Set();
  const mcpServers = new Set();
  const cursorDir =
    options.cursorDir ?? catalogCursorDir(options.catalogPath);

  for (const entry of catalog.packs) {
    if (entry.status !== 'selectable') continue;
    for (const s of entry.skills ?? []) skills.add(s);
    for (const r of entry.rules ?? []) rules.add(r);
    if (!entry.mcp) continue;
    const fragmentPath = path.resolve(cursorDir, entry.mcp);
    for (const name of Object.keys(loadMcpFragmentServers(fragmentPath))) {
      mcpServers.add(name);
    }
  }
  return { skills, rules, mcpServers };
}

/**
 * Load target mcp.json object (empty object when missing/blank). Mirrors mcp-merge.py load.
 * @param {string} targetPath
 * @returns {Record<string, unknown>}
 */
function loadMcpJson(targetPath) {
  if (!existsSync(targetPath)) return {};
  const text = readFileSync(targetPath, 'utf8').trim();
  if (!text) return {};
  const data = JSON.parse(text);
  if (!data || typeof data !== 'object' || Array.isArray(data)) {
    throw new Error(`${targetPath}: expected a JSON object`);
  }
  return /** @type {Record<string, unknown>} */ (data);
}

/**
 * Write mcp.json with trailing newline. Mirrors mcp-merge.py write.
 * @param {string} targetPath
 * @param {Record<string, unknown>} data
 */
function writeMcpJson(targetPath, data) {
  mkdirSync(path.dirname(targetPath), { recursive: true });
  writeFileSync(targetPath, `${JSON.stringify(data, null, 2)}\n`, 'utf8');
}

/**
 * Merge selected pack MCP servers into target; unmerge managed names not in selection.
 * Semantics match scripts/mcp-merge.py merge + unmerge (user servers untouched).
 * @param {string} targetPath
 * @param {Record<string, unknown>} selectedServers
 * @param {Set<string>} managedServerNames
 */
export function reconcilePackMcp(targetPath, selectedServers, managedServerNames) {
  const target = loadMcpJson(targetPath);
  const tgtServers = {
    ...(typeof target.mcpServers === 'object' &&
    target.mcpServers &&
    !Array.isArray(target.mcpServers)
      ? /** @type {Record<string, unknown>} */ (target.mcpServers)
      : {}),
  };
  const selectedNames = new Set(Object.keys(selectedServers));

  for (const [name, cfg] of Object.entries(selectedServers)) {
    tgtServers[name] = cfg;
  }

  for (const name of managedServerNames) {
    if (!selectedNames.has(name) && Object.hasOwn(tgtServers, name)) {
      delete tgtServers[name];
    }
  }

  const keys = Object.keys(target);
  if (Object.keys(tgtServers).length > 0) {
    target.mcpServers = tgtServers;
    writeMcpJson(targetPath, target);
    return;
  }

  if (!existsSync(targetPath)) {
    return;
  }

  if (keys.length === 0 || (keys.length === 1 && keys[0] === 'mcpServers')) {
    rmSync(targetPath, { force: true });
    return;
  }

  target.mcpServers = tgtServers;
  writeMcpJson(targetPath, target);
}

/**
 * Selectable pack ids from catalog (stable catalog order).
 * @param {{ packs: Array<{ id: string, status: string }> }} catalog
 * @returns {string[]}
 */
export function selectablePackIds(catalog) {
  return catalog.packs.filter((p) => p.status === 'selectable').map((p) => p.id);
}

/**
 * Resolve which packs to install and whether to write the profile.
 * Prefer cli/lib/setup.mjs (fail-closed) for bootstrap / setup path.
 * Direct callers may still hit legacy-all when no env/profile (T5/T6 retire).
 * @param {{
 *   targetDir: string,
 *   catalog: ReturnType<typeof loadCatalog>,
 *   packsEnv?: string | undefined,
 *   packIds?: string[],
 *   writeProfile?: boolean,
 *   source?: string,
 *   allowLegacyAll?: boolean,
 * }} opts
 * @returns {{ packIds: string[], writeProfile: boolean, source: string }}
 */
export function resolveInstallPacks({
  targetDir,
  catalog,
  packsEnv = process.env.SPACECRAFT_PACKS,
  packIds: explicitPackIds,
  writeProfile: explicitWrite,
  source: explicitSource,
  allowLegacyAll = true,
}) {
  if (Array.isArray(explicitPackIds)) {
    validateProfile({ version: 1, packs: explicitPackIds }, catalog);
    return {
      packIds: [...explicitPackIds],
      writeProfile: Boolean(explicitWrite),
      source: explicitSource ?? 'packs',
    };
  }

  const envRaw = packsEnv != null ? String(packsEnv).trim() : '';
  if (envRaw) {
    const packIds = parsePacksEnv(envRaw);
    validateProfile({ version: 1, packs: packIds }, catalog);
    return { packIds, writeProfile: true, source: 'env' };
  }

  const profilePath = path.join(targetDir, '.cursor', PROFILE_NAME);
  if (existsSync(profilePath)) {
    const raw = JSON.parse(readFileSync(profilePath, 'utf8'));
    const validated = validateProfile(raw, catalog);
    return { packIds: validated.packs, writeProfile: false, source: 'profile' };
  }

  if (!allowLegacyAll) {
    throw new Error(
      'No spacecraft-profile.json and no SPACECRAFT_PACKS; refuse silent all-packs install',
    );
  }

  // Legacy all-selectable (setup path sets allowLegacyAll: false / passes packIds).
  return {
    packIds: selectablePackIds(catalog),
    writeProfile: false,
    source: 'legacy-all',
  };
}

function writeProfile(targetDir, packIds) {
  const cursorDir = path.join(targetDir, '.cursor');
  mkdirSync(cursorDir, { recursive: true });
  const payload = `${JSON.stringify({ version: 1, packs: packIds }, null, 2)}\n`;
  writeFileSync(path.join(cursorDir, PROFILE_NAME), payload, 'utf8');
}

function pruneAgents(targetDir) {
  const agentsDir = path.join(targetDir, '.cursor', 'agents');
  if (!existsSync(agentsDir)) return;
  for (const name of readdirSync(agentsDir)) {
    if (!name.startsWith('sc-') || !name.endsWith('.md')) continue;
    rmSync(path.join(agentsDir, name), { force: true });
  }
}

function installSkill(srcSkills, destSkills, name) {
  const from = path.join(srcSkills, name);
  const to = path.join(destSkills, name);
  if (!existsSync(from) || !existsSync(path.join(from, 'SKILL.md'))) {
    throw new Error(`missing source skill ${name} under ${srcSkills}`);
  }
  rmSync(to, { recursive: true, force: true });
  cpSync(from, to, { recursive: true });
}

/**
 * Reconcile project .cursor skills + rules + pack MCP for selected packs.
 * @param {string} targetDir
 * @param {string} sourceDir
 * @param {{
 *   catalogPath?: string,
 *   packsEnv?: string,
 *   packIds?: string[],
 *   writeProfile?: boolean,
 *   source?: string,
 *   allowLegacyAll?: boolean,
 * }} [options]
 * @returns {{ packIds: string[], skills: string[], rules: string[], mcp: string[], source: string }}
 */
export function installProjectSurface(targetDir, sourceDir, options = {}) {
  const catalogPath =
    options.catalogPath ?? path.join(sourceDir, '.cursor', 'spacecraft-packs.json');
  const catalog = loadCatalog(catalogPath);
  const cursorDir = catalogCursorDir(catalogPath);
  const resolved = resolveInstallPacks({
    targetDir,
    catalog,
    packsEnv: options.packsEnv,
    packIds: options.packIds,
    writeProfile: options.writeProfile,
    source: options.source,
    allowLegacyAll: options.allowLegacyAll,
  });
  const {
    skills: wantSkills,
    rules: wantRules,
    mcp: wantMcpPaths,
  } = expandPacks(resolved.packIds, catalog, { cursorDir, catalogPath });
  const managed = managedCatalogPaths(catalog, { cursorDir, catalogPath });
  const wantSkillSet = new Set(wantSkills);
  const wantRuleSet = new Set(wantRules);

  if (resolved.writeProfile) {
    writeProfile(targetDir, resolved.packIds);
  }

  const destRules = path.join(targetDir, '.cursor', 'rules');
  const destSkills = path.join(targetDir, '.cursor', 'skills');
  const srcRules = path.join(sourceDir, '.cursor', 'rules');
  const srcSkills = path.join(sourceDir, '.cursor', 'skills');

  mkdirSync(destRules, { recursive: true });
  mkdirSync(destSkills, { recursive: true });

  for (const rule of wantRules) {
    if (USER_RULE_SET.has(rule)) continue;
    const from = path.join(srcRules, rule);
    if (!existsSync(from)) {
      throw new Error(`missing source rule ${rule} under ${srcRules}`);
    }
    copyFileSync(from, path.join(destRules, rule));
  }

  // Prune spacecraft-managed rules not in the selected union (leave unrelated).
  if (existsSync(destRules)) {
    for (const name of readdirSync(destRules)) {
      if (!name.endsWith('.mdc')) continue;
      if (USER_RULE_SET.has(name)) {
        rmSync(path.join(destRules, name), { force: true });
        continue;
      }
      if (managed.rules.has(name) && !wantRuleSet.has(name)) {
        rmSync(path.join(destRules, name), { force: true });
      }
    }
  }

  for (const skill of wantSkills) {
    if (LEAN_SET.has(skill)) continue;
    installSkill(srcSkills, destSkills, skill);
  }

  // Prune managed + lean leftovers; leave unrelated user skill dirs alone.
  if (existsSync(destSkills)) {
    for (const name of readdirSync(destSkills)) {
      const dest = path.join(destSkills, name);
      if (LEAN_SET.has(name)) {
        rmSync(dest, { recursive: true, force: true });
        continue;
      }
      if (managed.skills.has(name) && !wantSkillSet.has(name)) {
        rmSync(dest, { recursive: true, force: true });
      }
    }
  }

  pruneAgents(targetDir);

  const selectedMcpServers = {};
  for (const fragmentPath of wantMcpPaths) {
    Object.assign(selectedMcpServers, loadMcpFragmentServers(fragmentPath));
  }
  reconcilePackMcp(
    path.join(targetDir, '.cursor', 'mcp.json'),
    selectedMcpServers,
    managed.mcpServers,
  );

  return {
    packIds: resolved.packIds,
    skills: wantSkills.filter((s) => !LEAN_SET.has(s)),
    rules: wantRules.filter((r) => !USER_RULE_SET.has(r)),
    mcp: Object.keys(selectedMcpServers),
    source: resolved.source,
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
    console.error('usage: project-install.mjs <target-project-dir> <source-repo-dir>');
    process.exit(1);
  }
  try {
    const result = installProjectSurface(path.resolve(target), path.resolve(source));
    const label =
      result.source === 'legacy-all'
        ? 'all selectable packs (legacy)'
        : result.source === 'env'
          ? `SPACECRAFT_PACKS=${result.packIds.join(',')}`
          : `profile packs=${result.packIds.join(',')}`;
    console.log(`  domain rules/skills (${label}) -> ${path.resolve(target)}/.cursor`);
  } catch (err) {
    console.error(`project-install: ${err.message}`);
    process.exit(1);
  }
}
