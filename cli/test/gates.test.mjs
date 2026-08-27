/**
 * T1 — gates-version registry and grandfathering (M9G7IHV3 freeze-tooling).
 */
import assert from 'node:assert/strict';
import test from 'node:test';

async function loadGates() {
  const mod = await import('../lib/gates.mjs');
  assert.equal(typeof mod.gatesAtOrAfter, 'function', 'gatesAtOrAfter must be exported');
  assert.equal(typeof mod.readGatesVersion, 'function', 'readGatesVersion must be exported');
  return mod;
}

test('pos-gates-order: milestone registry ordering and gatesAtOrAfter boundaries', async () => {
  const { gatesAtOrAfter } = await loadGates();

  assert.equal(gatesAtOrAfter('M9G7IHHW', 'M9G7IHON'), false);
  assert.equal(gatesAtOrAfter('M9G7IHON', 'M9G7IHV3'), false);
  assert.equal(gatesAtOrAfter('M9G7IHV3', 'M9G7IHV3'), true);
  assert.equal(gatesAtOrAfter('M9G7IHV3', 'M9G7IHON'), true);
  assert.equal(gatesAtOrAfter('M9G7IHON', 'M9G7IHON'), true);
});

test('edge-grandfather-absent: missing Gates version line is pre-M1 (E4)', async () => {
  const { gatesAtOrAfter, readGatesVersion } = await loadGates();

  const decisionsText = '# Decisions\n\nNo gates line here.\n';
  const missionGate = readGatesVersion(decisionsText);
  assert.equal(
    gatesAtOrAfter(missionGate, 'M9G7IHV3'),
    false,
    'absent Gates version must not enable M9G7IHV3 checks',
  );
});

test('neg-gates-unknown: unknown or malformed Gates version does not silently enable checks', async () => {
  const { gatesAtOrAfter, readGatesVersion } = await loadGates();

  for (const text of [
    '# Decisions\n\nGates version: BOGUS\n',
    '# Decisions\n\nGates version:\n',
    '# Decisions\n\nGates version: M9G7IHV3-extra-garbage\n',
  ]) {
    const missionGate = readGatesVersion(text);
    assert.equal(
      gatesAtOrAfter(missionGate, 'M9G7IHV3'),
      false,
      `malformed "${text.trim()}" must not silently enable M9G7IHV3 checks`,
    );
  }
});

// --- T1: M9G7II1F quality-tooling gates milestone ---

test('pos-gates-m9g7ii1f: M9G7II1F registry entry and gatesAtOrAfter boundaries', async () => {
  const { GATES_REGISTRY, gatesAtOrAfter } = await loadGates();

  const m3Idx = GATES_REGISTRY.indexOf('M9G7IHV3');
  const m4Idx = GATES_REGISTRY.indexOf('M9G7II1F');
  assert.notEqual(m4Idx, -1, 'M9G7II1F must be in GATES_REGISTRY');
  assert.ok(m3Idx >= 0 && m4Idx > m3Idx, 'M9G7II1F must follow M9G7IHV3');
  assert.equal(gatesAtOrAfter('M9G7II1F', 'M9G7II1F'), true);
  assert.equal(gatesAtOrAfter('M9G7IHV3', 'M9G7II1F'), false);
});

test('over-gates-grandfather: pre-M9G7II1F gates do not enable M9G7II1F-only checks', async () => {
  const { gatesAtOrAfter, readGatesVersion } = await loadGates();

  const decisionsText = '# Decisions\n\nGates version: M9G7IHV3\n';
  assert.equal(readGatesVersion(decisionsText), 'M9G7IHV3');
  assert.equal(
    gatesAtOrAfter(readGatesVersion(decisionsText), 'M9G7II1F'),
    false,
    'M9G7IHV3 must not clear M9G7II1F gate',
  );

  for (const text of [
    '# Decisions\n\nNo gates line here.\n',
    '# Decisions\n\nGates version: BOGUS\n',
    '# Decisions\n\nGates version:\n',
  ]) {
    const missionGate = readGatesVersion(text);
    assert.equal(
      gatesAtOrAfter(missionGate, 'M9G7II1F'),
      false,
      `absent/unknown gates must not enable M9G7II1F checks (${text.trim()})`,
    );
  }
});
