/**
 * Lean prompt surface contract freeze (M9JQIN7V T1).
 * Greppable SoT literals that must survive lean refactors. Mirrors
 * scripts/check-workflow-*.sh, check-sc-planner-*.sh, check-sc-planning-*.sh
 * plus optional-lane / judge / probe / discuss / impeccable / AUTH tokens.
 * No Verify bars invented here.
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

/** ≤7 / <=7 / 7 tasks cap (check-*-hard-must.sh). */
function assertLeq7Cap(text, label) {
  assert.ok(
    /≤7|<=7|7 tasks/.test(text),
    `${label} must mention ≤7 / <=7 / 7 tasks`,
  );
}

/** Phrase "hard Must" (check-*-hard-must.sh; JS \s ≡ shell [[:space:]]). */
function assertHardMust(text, label) {
  assert.ok(
    /hard\s+Must/i.test(text),
    `${label} must state hard Must (not preference-only)`,
  );
}

/** Non-preference / non-soft stance (check-*-hard-must.sh). */
function assertNotPreferenceOnly(text, label) {
  assert.ok(
    /not\s+(a\s+)?preference|not\s+preference-only|not\s+soft/i.test(text),
    `${label} must reject preference-only / soft wording for ≤7 cap`,
  );
}

/** Explicit reject soft prefer ≤7 (check-*-reject-soft.sh). */
function assertRejectSoftPreferLeq7(text, label) {
  assert.ok(
    /reject.*prefer\s*(≤|<=)\s*7|prefer\s*(≤|<=)\s*7.*(reject|must\s+not|forbidden|disallowed)|soft\s+prefer\s*(≤|<=)\s*7|no\s+soft.*"prefer\s*(≤|<=)\s*7"|exclude.*prefer\s*(≤|<=)\s*7|Must\s+not:.*prefer\s*(≤|<=)\s*7/i.test(
      text,
    ),
    `${label} must explicitly reject soft prefer ≤7 (or prefer <=7)`,
  );
}

/** Explicit reject 8-9 exception band (check-*-reject-soft.sh). */
function assertRejectEightNineBand(text, label) {
  assert.ok(
    /reject.*(8[-–]9|8\s*[-–]\s*9).*(exception|band)|(8[-–]9|8\s*[-–]\s*9).*(exception\s+band).*(reject|must\s+not|forbidden|disallowed|no)|(no|reject|exclude|must\s+not).*(8[-–]9|8\s*[-–]\s*9).*(exception|band)|Must\s+not:.*(8[-–]9|8\s*[-–]\s*9)/i.test(
      text,
    ),
    `${label} must explicitly reject any 8-9 / 8–9 exception band`,
  );
}

function assertPlanPhaseN(text, label) {
  assert.ok(
    /plan-phaseN\.json|plan-phase<N>\.json|plan-phase[0-9]+\.json/.test(text),
    `${label} must name plan-phaseN.json (or plan-phase<N>.json / plan-phase1.json)`,
  );
}

function assertSameMission(text, label) {
  assert.ok(
    /same[- ]mission/i.test(text),
    `${label} must document same-mission phase split`,
  );
}

// --- 200-workflow.mdc (check-workflow-*.sh) ---

const WORKFLOW = '.cursor/rules/200-workflow.mdc';

test('200-workflow: ≤7 cap as hard Must (not preference-only)', () => {
  const text = readUtf8(WORKFLOW);
  assertLeq7Cap(text, WORKFLOW);
  assertHardMust(text, WORKFLOW);
  assertNotPreferenceOnly(text, WORKFLOW);
});

test('200-workflow: reject soft prefer ≤7 and 8-9 exception band', () => {
  const text = readUtf8(WORKFLOW);
  assertRejectSoftPreferLeq7(text, WORKFLOW);
  assertRejectEightNineBand(text, WORKFLOW);
});

test('200-workflow: plan-phaseN + same-mission + mission-sizing + map ban', () => {
  const text = readUtf8(WORKFLOW);
  assertPlanPhaseN(text, WORKFLOW);
  assertSameMission(text, WORKFLOW);
  assertIncludes(text, 'mission-sizing', WORKFLOW);
  assert.ok(/multi[- ]mission/i.test(text), `${WORKFLOW} must document multi-mission`);
  assert.ok(
    /spacecraft\s+map|map new/.test(text),
    `${WORKFLOW} must name spacecraft map / map new`,
  );
  assert.ok(
    /discuss owns|\/sc-discuss/i.test(text),
    `${WORKFLOW} must say discuss owns map or hand to /sc-discuss`,
  );
  assert.ok(
    /must not.*map new|never.*map new|planning must not/i.test(text),
    `${WORKFLOW} must forbid planning-owned map new`,
  );
});

// --- sc-planner.md (check-sc-planner-*.sh) ---

const PLANNER = '.cursor/agents/sc-planner.md';

test('sc-planner: ≤7 cap as hard Must (not preference-only)', () => {
  const text = readUtf8(PLANNER);
  assertLeq7Cap(text, PLANNER);
  assertHardMust(text, PLANNER);
  assertNotPreferenceOnly(text, PLANNER);
});

test('sc-planner: reject soft prefer ≤7 and 8-9 exception band', () => {
  const text = readUtf8(PLANNER);
  assertRejectSoftPreferLeq7(text, PLANNER);
  assertRejectEightNineBand(text, PLANNER);
});

test('sc-planner: plan-phaseN + discuss handoff + map ban', () => {
  const text = readUtf8(PLANNER);
  assertPlanPhaseN(text, PLANNER);
  assertSameMission(text, PLANNER);
  assertIncludes(text, '/sc-discuss', PLANNER);
  assertIncludes(text, 'mission-sizing', PLANNER);
  assert.ok(
    /multi[- ]mission|roadmap/i.test(text),
    `${PLANNER} must document multi-mission or roadmap handoff`,
  );
  assert.ok(
    /spacecraft\s+map/.test(text),
    `${PLANNER} must name spacecraft map`,
  );
  assert.ok(
    /never create|discuss owns|must not.*map/i.test(text),
    `${PLANNER} must forbid planner-owned map create`,
  );
});

// --- sc-planning/SKILL.md (check-sc-planning-*.sh) ---

const PLANNING = '.cursor/skills/sc-planning/SKILL.md';

test('sc-planning: ≤7 cap as hard Must (not preference-only)', () => {
  const text = readUtf8(PLANNING);
  assertLeq7Cap(text, PLANNING);
  assertHardMust(text, PLANNING);
  assertNotPreferenceOnly(text, PLANNING);
});

test('sc-planning: reject soft prefer ≤7 and 8-9 exception band', () => {
  const text = readUtf8(PLANNING);
  assertRejectSoftPreferLeq7(text, PLANNING);
  assertRejectEightNineBand(text, PLANNING);
});

test('sc-planning: same-mission plan-phaseN + multi-mission map ban', () => {
  const text = readUtf8(PLANNING);
  assertPlanPhaseN(text, PLANNING);
  assertSameMission(text, PLANNING);
  assertIncludes(text, '/sc-discuss', PLANNING);
  assertIncludes(text, 'mission-sizing', PLANNING);
  assert.ok(/multi[- ]mission/i.test(text), `${PLANNING} must document multi-mission`);
  assertIncludes(text, 'map new', PLANNING);
  assert.ok(
    /must not.*map new|never.*map new|do not.*map new/i.test(text),
    `${PLANNING} must forbid planning-owned map new`,
  );
});

// --- Optional-lane disposition prefixes ---

const LANE_PREFIXES = [
  ['.cursor/skills/sc-loop/SKILL.md', 'Loop watch:'],
  ['.cursor/skills/sc-automate-slack/SKILL.md', 'Automate-Slack:'],
  ['.cursor/skills/sc-canvas-sot/SKILL.md', 'Canvas-sot:'],
  ['.cursor/skills/sc-goal-roadmap/SKILL.md', 'Goal-roadmap:'],
  ['.cursor/skills/sc-post-ready-drain/SKILL.md', 'Post-ready drain:'],
  ['.cursor/skills/sc-split-to-prs/SKILL.md', 'Split-to-prs:'],
];

for (const [file, prefix] of LANE_PREFIXES) {
  test(`${file} keeps disposition prefix ${JSON.stringify(prefix)}`, () => {
    assertIncludes(readUtf8(file), prefix, file);
  });
}

// --- Judge / probe / discuss / impeccable / hard-contract / security ---

test('sc-judge: VERIFIED and REFUTED', () => {
  const file = '.cursor/skills/sc-judge/SKILL.md';
  const text = readUtf8(file);
  assertIncludes(text, 'VERIFIED', file);
  assertIncludes(text, 'REFUTED', file);
});

test('sc-browser-probe: PROBE: CLEAN', () => {
  const file = '.cursor/skills/sc-browser-probe/SKILL.md';
  assertIncludes(readUtf8(file), 'PROBE: CLEAN', file);
});

test('sc-discuss: Discuss path: fast and Roadmap contract skipped: pre-M1 map', () => {
  const file = '.cursor/skills/sc-discuss/SKILL.md';
  const text = readUtf8(file);
  assertIncludes(text, 'Discuss path: fast', file);
  assertIncludes(text, 'Roadmap contract skipped: pre-M1 map', file);
});

test('Impeccable path: active and Impeccable path: skipped in ux or orchestration', () => {
  const skill = '.cursor/skills/sc-ux-design/SKILL.md';
  const orch = '.cursor/skills/sc-ux-design/references/impeccable-orchestration.md';
  const combined = `${readUtf8(skill)}\n${readUtf8(orch)}`;
  assertIncludes(combined, 'Impeccable path: active', `${skill} or ${orch}`);
  assertIncludes(combined, 'Impeccable path: skipped', `${skill} or ${orch}`);
});

test('010-hard-contract: AUTH: and INTENT:', () => {
  const file = '.cursor/rules/010-hard-contract.mdc';
  const text = readUtf8(file);
  assertIncludes(text, 'AUTH:', file);
  assertIncludes(text, 'INTENT:', file);
});

test('sc-security: Sc-security fallback: pass', () => {
  const file = '.cursor/skills/sc-security/SKILL.md';
  assertIncludes(readUtf8(file), 'Sc-security fallback: pass', file);
});

// --- sc-debug lane (RCA opener; lean-core) ---

const DEBUG_SKILL = '.cursor/skills/sc-debug/SKILL.md';

test('sc-debug: slash lane + packs + no-ship', () => {
  const skill = readUtf8(DEBUG_SKILL);
  const workflow = readUtf8(WORKFLOW);
  const lean = readUtf8('cli/lib/project-install.mjs');
  assertIncludes(skill, 'Pack:', DEBUG_SKILL);
  assertIncludes(skill, 'hardware-mcu', DEBUG_SKILL);
  assertIncludes(skill, 'hardware-fpga', DEBUG_SKILL);
  assertIncludes(skill, '/sc-quick', DEBUG_SKILL);
  assert.ok(
    /Must not[\s\S]*[Mm]erge/m.test(skill),
    `${DEBUG_SKILL} must forbid merge from the debug lane`,
  );
  assertIncludes(workflow, '/sc-debug', WORKFLOW);
  assertIncludes(lean, "'sc-debug'", 'cli/lib/project-install.mjs');
  assert.ok(
    existsSync(rel('.cursor/skills/sc-debug/references/software.md')),
    'sc-debug software pack must exist',
  );
  assert.ok(
    existsSync(rel('.cursor/skills/sc-debug/references/hardware.md')),
    'sc-debug hardware pack must exist',
  );
  assert.ok(
    existsSync(rel('.cursor/skills/sc-debug/references/visual.md')),
    'sc-debug visual pack must exist',
  );
});

// --- Must-exist lean SoT references (required today) ---


const MUST_EXIST = [
  '.cursor/skills/sc-run/references/defect-finding.md',
  '.cursor/skills/sc-ux-design/references/impeccable-orchestration.md',
];

for (const file of MUST_EXIST) {
  test(`required reference exists: ${file}`, () => {
    assert.ok(existsSync(rel(file)), `${file} must exist`);
  });
}

// --- sc-firmware vs sc-rtl agent split (greppable SoT; not plugin twins) ---

const FIRMWARE_AGENT = '.cursor/agents/sc-firmware.md';
const RTL_AGENT = '.cursor/agents/sc-rtl.md';

test('sc-firmware vs sc-rtl: greppable embedded vs digital-IC split', () => {
  const firmware = readUtf8(FIRMWARE_AGENT);
  const rtl = readUtf8(RTL_AGENT);
  const workflow = readUtf8(WORKFLOW);

  assertIncludes(firmware, 'Embedded system engineer', FIRMWARE_AGENT);
  assertIncludes(firmware, 'Not FPGA/RTL', FIRMWARE_AGENT);
  assertIncludes(
    firmware,
    'FPGA / SystemVerilog / RTL / digital IC (Task `sc-rtl`)',
    FIRMWARE_AGENT,
  );

  assertIncludes(rtl, 'Digital IC designer', RTL_AGENT);
  assertIncludes(rtl, 'Not MCU firmware', RTL_AGENT);
  assertIncludes(
    rtl,
    'STM32 / HAL / CubeMX / MCU firmware (Task `sc-firmware`)',
    RTL_AGENT,
  );

  const taskLine = workflow.split('\n').find((line) => line.includes('Task:'));
  assert.ok(taskLine, `${WORKFLOW} Subagent Delegation must have a Task: line`);
  assertIncludes(taskLine, 'sc-firmware', `${WORKFLOW} Task:`);
  assertIncludes(taskLine, 'sc-rtl', `${WORKFLOW} Task:`);
});
