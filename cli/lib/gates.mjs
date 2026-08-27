/** Ordered milestone registry for grandfathering machine gates. */
export const GATES_REGISTRY = ['M9G7IHHW', 'M9G7IHON', 'M9G7IHV3', 'M9G7II1F'];

/**
 * Parse `Gates version: <id>` from decisions.md text.
 * @param {string} decisionsText
 * @returns {string|null} known milestone id or null when absent/unknown
 */
export function readGatesVersion(decisionsText) {
  if (typeof decisionsText !== 'string') return null;
  const match = decisionsText.match(/^Gates version:\s*(\S+)\s*$/m);
  if (!match) return null;
  const id = match[1].trim();
  return GATES_REGISTRY.includes(id) ? id : null;
}

/**
 * @param {string|null} missionGate from readGatesVersion
 * @param {string} requiredGate milestone id from GATES_REGISTRY
 * @returns {boolean}
 */
export function gatesAtOrAfter(missionGate, requiredGate) {
  if (!missionGate || !requiredGate) return false;
  const missionIdx = GATES_REGISTRY.indexOf(missionGate);
  const requiredIdx = GATES_REGISTRY.indexOf(requiredGate);
  if (missionIdx === -1 || requiredIdx === -1) return false;
  return missionIdx >= requiredIdx;
}
