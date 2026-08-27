import { existsSync, readFileSync } from 'node:fs';
import path from 'node:path';

function configFilePath(rootOrSpaceDir) {
  const inSpace = path.join(rootOrSpaceDir, 'config.json');
  if (existsSync(inSpace)) return inSpace;
  return path.join(rootOrSpaceDir, '.space', 'config.json');
}

/**
 * Read `.space/config.json` from project root or `.space` directory.
 * @param {string} rootOrSpaceDir
 * @returns {{ criticFamily: string | null }}
 */
export function readSpaceConfig(rootOrSpaceDir) {
  const configPath = configFilePath(rootOrSpaceDir);
  if (!existsSync(configPath)) {
    return { criticFamily: null };
  }

  let raw;
  try {
    raw = readFileSync(configPath, 'utf8');
  } catch (err) {
    throw new Error(`spacecraft config: cannot read ${configPath}: ${err.message}`);
  }

  let data;
  try {
    data = JSON.parse(raw);
  } catch {
    throw new Error('spacecraft config: malformed JSON in .space/config.json');
  }

  if (data === null || typeof data !== 'object' || Array.isArray(data)) {
    throw new Error('spacecraft config: malformed .space/config.json (expected object)');
  }

  if (!('criticFamily' in data)) {
    return { criticFamily: null };
  }

  const family = data.criticFamily;
  if (typeof family !== 'string' || family.trim() === '') {
    throw new Error(
      'spacecraft config: malformed criticFamily in .space/config.json (must be non-empty string)',
    );
  }

  return { criticFamily: family };
}

/**
 * @param {string} rootOrSpaceDir project root or `.space` directory
 * @returns {string | null}
 */
export function readCriticFamily(rootOrSpaceDir) {
  return readSpaceConfig(rootOrSpaceDir).criticFamily;
}
