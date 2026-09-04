/**
 * Firmware core / OS / target contract freeze (M9M9K3HQ T7 / A13).
 * Greppable SoT literals from approved-scenarios S01–S14, S17–S18.
 * Identity split (firmware vs RTL) stays in prompt-lean-contract.test.mjs.
 */
import assert from 'node:assert/strict';
import { existsSync, readFileSync } from 'node:fs';
import path from 'node:path';
import test from 'node:test';
import { fileURLToPath } from 'node:url';

const repoRoot = path.resolve(path.dirname(fileURLToPath(import.meta.url)), '../..');

function rel(...parts) {
  return path.join(repoRoot, ...parts);
}

function readUtf8(...parts) {
  return readFileSync(rel(...parts), 'utf8');
}

function assertIncludes(text, needle, label) {
  assert.ok(
    text.includes(needle),
    `${label} must include ${JSON.stringify(needle)}`,
  );
}

const FIRMWARE_SKILL = '.cursor/skills/sc-firmware/SKILL.md';
const FIRMWARE_AGENT = '.cursor/agents/sc-firmware.md';
const CORE = '.cursor/skills/sc-firmware/references/core.md';
const OS_BARE = '.cursor/skills/sc-firmware/references/os-bare-metal.md';
const OS_RTOS = '.cursor/skills/sc-firmware/references/os-rtos.md';
const TARGET_CORTEX = '.cursor/skills/sc-firmware/references/target-cortex-m.md';
const TARGET_STM32 = '.cursor/skills/sc-firmware/references/target-stm32.md';
const TARGET_NRF = '.cursor/skills/sc-firmware/references/target-nrf52840.md';
const RULE_600 = '.cursor/rules/600-firmware.mdc';
const ARCHITECTURE = '.cursor/skills/sc-firmware/references/architecture.md';
const PACKS = '.cursor/spacecraft-packs.json';

const REQUIRED_REFS = [
  CORE,
  OS_BARE,
  OS_RTOS,
  TARGET_CORTEX,
  TARGET_STM32,
  TARGET_NRF,
];

const SKILL_REF_NEEDLES = [
  'references/core.md',
  'references/os-bare-metal.md',
  'references/os-rtos.md',
  'references/target-cortex-m.md',
  'references/target-stm32.md',
  'references/target-nrf52840.md',
];

test('firmware references: core/OS/ISA/target files exist; architecture.md retired', () => {
  for (const file of REQUIRED_REFS) {
    assert.ok(existsSync(rel(file)), `${file} must exist`);
  }
  assert.ok(!existsSync(rel(ARCHITECTURE)), `${ARCHITECTURE} must not exist`);
  assert.ok(readUtf8(OS_BARE).trim().length > 0, `${OS_BARE} must be non-empty`);
  assert.ok(readUtf8(OS_RTOS).trim().length > 0, `${OS_RTOS} must be non-empty`);
});

test('SKILL.md: OS class, needs-input, and reference routing', () => {
  const text = readUtf8(FIRMWARE_SKILL);
  for (const needle of SKILL_REF_NEEDLES) {
    assertIncludes(text, needle, FIRMWARE_SKILL);
  }
  assertIncludes(text, 'bare-metal | RTOS', FIRMWARE_SKILL);
  assertIncludes(text, 'needs-input', FIRMWARE_SKILL);
});

test('SKILL.md: Cross-domain HIL dual evidence and proof oracles', () => {
  const text = readUtf8(FIRMWARE_SKILL);
  assertIncludes(text, 'Cross-domain HIL', FIRMWARE_SKILL);
  assertIncludes(text, 'dual evidence', FIRMWARE_SKILL);
  assertIncludes(
    text,
    'host green ≠ HIL green without proof oracles on target',
    FIRMWARE_SKILL,
  );
  assertIncludes(text, 'blaming peer RTL or wiring', FIRMWARE_SKILL);
  assertIncludes(
    text,
    'After HIL RCA, append one greppable lesson to `.space/trust/lessons.md`',
    FIRMWARE_SKILL,
  );
});

test('sc-firmware agent: embedded identity without F746 default board', () => {
  const text = readUtf8(FIRMWARE_AGENT);
  assertIncludes(text, 'Embedded system engineer', FIRMWARE_AGENT);
  assertIncludes(text, 'MCU firmware', FIRMWARE_AGENT);
  assertIncludes(text, 'Not FPGA/RTL', FIRMWARE_AGENT);
  assert.ok(
    !text.includes('Default board STM32F746NG-Discovery'),
    `${FIRMWARE_AGENT} must not include ${JSON.stringify('Default board STM32F746NG-Discovery')}`,
  );
});

test('target-stm32.md: named boards; STM32H723ZG and STM32L412 in skill tree', () => {
  const stm32 = readUtf8(TARGET_STM32);
  assertIncludes(stm32, 'STM32F746NG-Discovery', TARGET_STM32);
  assertIncludes(stm32, 'NUCLEO-H723ZG', TARGET_STM32);
  assertIncludes(stm32, 'NUCLEO-L412KB', TARGET_STM32);

  const skillTree = `${stm32}\n${readUtf8(FIRMWARE_SKILL)}`;
  const partLabel = `${TARGET_STM32} or ${FIRMWARE_SKILL}`;
  assertIncludes(skillTree, 'STM32H723ZG', partLabel);
  assertIncludes(skillTree, 'STM32L412', partLabel);
});

test('target-nrf52840.md: nRF52840 DK without CubeMX/HAL', () => {
  const text = readUtf8(TARGET_NRF);
  assertIncludes(text, 'nRF52840', TARGET_NRF);
  assertIncludes(text, 'PCA10056', TARGET_NRF);
  assertIncludes(text, 'does not require CubeMX/HAL', TARGET_NRF);
});

test('target-cortex-m.md: M4 vs M7 cache, MPU, DMA', () => {
  const text = readUtf8(TARGET_CORTEX);
  assertIncludes(text, 'M4', TARGET_CORTEX);
  assertIncludes(text, 'M7', TARGET_CORTEX);
  assert.ok(
    text.includes('cache') || text.includes('D-cache'),
    `${TARGET_CORTEX} must include "cache" or "D-cache"`,
  );
  assertIncludes(text, 'MPU', TARGET_CORTEX);
  assertIncludes(text, 'DMA', TARGET_CORTEX);
});

test('600-firmware.mdc: Power-of-Ten-class Musts without CubeMX required', () => {
  const text = readUtf8(RULE_600);
  assertIncludes(text, 'no recursion', RULE_600);
  assertIncludes(text, 'no heap after init', RULE_600);
  assertIncludes(text, 'bounded loops', RULE_600);
  assertIncludes(text, 'static analysis', RULE_600);
  assertIncludes(text, 'does not require CubeMX as a Must', RULE_600);
});

test('core.md: FDIR safe-state handshake', () => {
  const text = readUtf8(CORE);
  assertIncludes(
    text,
    'failing assert / WDT / HardFault -> documented safe state recipe',
    CORE,
  );
  assertIncludes(text, 'FDIR', CORE);
});

test('SKILL.md: Linux/Yocto/userspace out of skill (iot pack)', () => {
  const text = readUtf8(FIRMWARE_SKILL);
  assertIncludes(text, 'Linux', FIRMWARE_SKILL);
  assertIncludes(text, 'Yocto', FIRMWARE_SKILL);
  assertIncludes(text, 'userspace process model', FIRMWARE_SKILL);
  assertIncludes(text, 'iot', FIRMWARE_SKILL);
});

test('embedded pack ships sc-firmware with 600/610/620', () => {
  const catalog = JSON.parse(readUtf8(PACKS));
  const embedded = catalog.packs.find((pack) => pack.id === 'embedded');
  assert.ok(embedded, `${PACKS} must have pack id "embedded"`);
  assert.ok(
    Array.isArray(embedded.skills) && embedded.skills.includes('sc-firmware'),
    `${PACKS} pack embedded must include skill sc-firmware`,
  );
  for (const rule of [
    '600-firmware.mdc',
    '610-firmware-peripherals.mdc',
    '620-firmware-testing.mdc',
  ]) {
    assert.ok(
      Array.isArray(embedded.rules) && embedded.rules.includes(rule),
      `${PACKS} pack embedded must include rule ${rule}`,
    );
  }
});
