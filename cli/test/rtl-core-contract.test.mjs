/**
 * RTL core / intent contract freeze.
 * Greppable SoT from approved-scenarios.
 * Identity split stays in prompt-lean-contract.test.mjs.
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

const RTL_SKILL = '.cursor/skills/sc-rtl/SKILL.md';
const RTL_AGENT = '.cursor/agents/sc-rtl.md';
const CORE = '.cursor/skills/sc-rtl/references/core.md';
const INTENT_FPGA = '.cursor/skills/sc-rtl/references/intent-fpga.md';
const INTENT_TB = '.cursor/skills/sc-rtl/references/intent-tb.md';
const INTENT_CDC = '.cursor/skills/sc-rtl/references/intent-cdc.md';
const INTENT_FORMAL = '.cursor/skills/sc-rtl/references/intent-formal.md';
const ARCH = '.cursor/skills/sc-rtl/references/arch.md';
const INTENT_ASIC = '.cursor/skills/sc-rtl/references/intent-asic.md';
const INTENT_DFT = '.cursor/skills/sc-rtl/references/intent-dft.md';
const RULE_700 = '.cursor/rules/700-rtl.mdc';
const VERIFY_SKILL = '.cursor/skills/sc-rtl-verify/SKILL.md';
const SIGNOFF = '.cursor/skills/sc-rtl-verify/references/signoff.md';
const PACKS = '.cursor/spacecraft-packs.json';

const FSM_STATES = ['S_FETCH', 'S_EXECUTE', 'S_ATOMIC_RMW'];

const REQUIRED_REFS = [
  CORE,
  INTENT_FPGA,
  INTENT_TB,
  INTENT_CDC,
  INTENT_FORMAL,
  ARCH,
  INTENT_ASIC,
  INTENT_DFT,
];

const SKILL_REF_NEEDLES = [
  'core.md',
  'intent-fpga.md',
  'intent-tb.md',
  'intent-cdc.md',
  'intent-formal.md',
  'arch.md',
];

const INTENT_FILENAMES = [
  'intent-fpga.md',
  'intent-tb.md',
  'intent-cdc.md',
  'intent-formal.md',
];

const DUT_CORE_NEEDLES = ['sim≠synth', 'no latches', 'always_ff', 'always_comb'];

test('rtl references: core/intent/arch/stub files exist and are non-empty', () => {
  for (const file of REQUIRED_REFS) {
    assert.ok(existsSync(rel(file)), `${file} must exist`);
    assert.ok(readUtf8(file).trim().length > 0, `${file} must be non-empty`);
  }
});

test('SKILL.md: intent routing and named reference list', () => {
  const text = readUtf8(RTL_SKILL);
  for (const needle of SKILL_REF_NEEDLES) {
    assertIncludes(text, needle, RTL_SKILL);
  }
  assertIncludes(text, 'DUT | TB | constraints', RTL_SKILL);
});

test('SKILL.md: Cross-domain HIL dual evidence and physical observe', () => {
  const text = readUtf8(RTL_SKILL);
  assertIncludes(text, 'Cross-domain HIL', RTL_SKILL);
  assertIncludes(text, 'dual evidence', RTL_SKILL);
  assertIncludes(text, 'physical board observe equals `$display`', RTL_SKILL);
  assertIncludes(
    text,
    'After HIL RCA, append one greppable lesson to `.space/trust/lessons.md`',
    RTL_SKILL,
  );
});

test('sc-rtl agent: digital IC identity without MCU firmware', () => {
  const text = readUtf8(RTL_AGENT);
  assertIncludes(text, 'Digital IC designer', RTL_AGENT);
  assertIncludes(text, 'FPGA', RTL_AGENT);
  assertIncludes(text, 'Not MCU firmware', RTL_AGENT);
});

test('SKILL.md: no RV32 FSM table; states live in arch.md', () => {
  const skill = readUtf8(RTL_SKILL);
  const arch = readUtf8(ARCH);
  for (const state of FSM_STATES) {
    assert.ok(!skill.includes(state), `${RTL_SKILL} must not include ${JSON.stringify(state)}`);
    assertIncludes(arch, state, ARCH);
  }
});

test('core.md and 700-rtl.mdc: shared synthesizable DUT needles', () => {
  const core = readUtf8(CORE);
  const rule = readUtf8(RULE_700);
  for (const needle of DUT_CORE_NEEDLES) {
    assertIncludes(core, needle, CORE);
    assertIncludes(rule, needle, RULE_700);
  }
});

test('700-rtl.mdc: thin flow; no CPU ISA inventory', () => {
  const text = readUtf8(RULE_700);
  assertIncludes(
    text,
    'does not require Yosys/nextpnr as the only legal flow',
    RULE_700,
  );
  for (const state of FSM_STATES) {
    assert.ok(!text.includes(state), `${RULE_700} must not include ${JSON.stringify(state)}`);
  }
});

test('intent-fpga.md: synchronous reset, BRAM, DSP', () => {
  const text = readUtf8(INTENT_FPGA);
  assertIncludes(text, 'synchronous reset', INTENT_FPGA);
  assertIncludes(text, 'BRAM', INTENT_FPGA);
  assertIncludes(text, 'DSP', INTENT_FPGA);
});

test('intent-tb.md: DUT-versus-TB firewall', () => {
  const text = readUtf8(INTENT_TB);
  assertIncludes(text, '# delays', INTENT_TB);
  assertIncludes(text, 'initial', INTENT_TB);
  assertIncludes(text, '$display', INTENT_TB);
  assertIncludes(text, 'TB-only', INTENT_TB);
  assertIncludes(text, 'sc-tester', INTENT_TB);
  assertIncludes(text, 'sc-rtl-verify', INTENT_TB);
});

test('intent-cdc.md: 2FF, FIFO, gray', () => {
  const text = readUtf8(INTENT_CDC);
  assertIncludes(text, '2FF', INTENT_CDC);
  assertIncludes(text, 'FIFO', INTENT_CDC);
  assertIncludes(text, 'gray', INTENT_CDC);
});

test('intent-formal.md: formal and SVA', () => {
  const text = readUtf8(INTENT_FORMAL);
  assertIncludes(text, 'formal', INTENT_FORMAL);
  assertIncludes(text, 'SVA', INTENT_FORMAL);
});

test('sc-rtl-verify: L0–L5 layers; signoff links core and intent refs', () => {
  const skill = readUtf8(VERIFY_SKILL);
  for (const layer of ['L0', 'L1', 'L2', 'L3', 'L4', 'L5']) {
    assertIncludes(skill, layer, VERIFY_SKILL);
  }
  const signoff = readUtf8(SIGNOFF);
  assertIncludes(signoff, 'core.md', SIGNOFF);
  for (const needle of INTENT_FILENAMES) {
    assertIncludes(signoff, needle, SIGNOFF);
  }
});

test('sc-rtl-verify: system/integration sim before board HIL; physical HIL for silicon', () => {
  const skill = readUtf8(VERIFY_SKILL);
  assertIncludes(
    skill,
    'Prefer system/integration sim (DUT + software image) before board HIL',
    VERIFY_SKILL,
  );
  assertIncludes(
    skill,
    'Claim silicon ready without physical HIL evidence',
    VERIFY_SKILL,
  );
  const signoff = readUtf8(SIGNOFF);
  assertIncludes(signoff, 'System/integration sim', SIGNOFF);
  assertIncludes(signoff, 'Physical HIL', SIGNOFF);
});

test('intent-asic.md and intent-dft.md: coming stubs; ASIC reset vs FPGA 700', () => {
  const asic = readUtf8(INTENT_ASIC);
  const dft = readUtf8(INTENT_DFT);
  const rule = readUtf8(RULE_700);
  assertIncludes(asic, 'out of scope', INTENT_ASIC);
  assertIncludes(asic, 'coming', INTENT_ASIC);
  assertIncludes(asic, 'async-assert', INTENT_ASIC);
  assertIncludes(dft, 'out of scope', INTENT_DFT);
  assertIncludes(dft, 'coming', INTENT_DFT);
  assertIncludes(rule, 'synchronous reset', RULE_700);
});

test('fpga pack ships sc-rtl and sc-rtl-verify with 700/710/720', () => {
  const catalog = JSON.parse(readUtf8(PACKS));
  const fpga = catalog.packs.find((pack) => pack.id === 'fpga');
  assert.ok(fpga, `${PACKS} must have pack id "fpga"`);
  for (const skill of ['sc-rtl', 'sc-rtl-verify']) {
    assert.ok(
      Array.isArray(fpga.skills) && fpga.skills.includes(skill),
      `${PACKS} pack fpga must include skill ${skill}`,
    );
  }
  for (const rule of [
    '700-rtl.mdc',
    '710-rtl-timing.mdc',
    '720-rtl-verify.mdc',
  ]) {
    assert.ok(
      Array.isArray(fpga.rules) && fpga.rules.includes(rule),
      `${PACKS} pack fpga must include rule ${rule}`,
    );
  }
});
