/**
 * Discuss fast-path skill/doc contract (M9C0ZW2K T3).
 * Greppable string + dual-tree parity oracles from design-contract Edge matrix
 * E1–E7 and approved-scenarios S1–S7. No runtime eligibility helper.
 */
import assert from 'node:assert/strict';
import { readFileSync } from 'node:fs';
import path from 'node:path';
import test from 'node:test';
import { fileURLToPath } from 'node:url';

const repoRoot = path.resolve(path.dirname(fileURLToPath(import.meta.url)), '../..');

const CURSOR_SKILL = path.join(repoRoot, '.cursor/skills/sc-discuss/SKILL.md');
const PLUGIN_SKILL = path.join(
  repoRoot,
  'plugins/spacecraft/skills/sc-discuss/SKILL.md',
);
const CURSOR_SIZING = path.join(
  repoRoot,
  '.cursor/skills/sc-discuss/references/mission-sizing.md',
);
const PLUGIN_SIZING = path.join(
  repoRoot,
  'plugins/spacecraft/skills/sc-discuss/references/mission-sizing.md',
);

/** Frozen marker (design-contract Data structures / E1 / S1). */
const FAST_MARKER = 'Discuss path: fast';

/**
 * Eligibility greppable literals (T1 acceptance / S1 / E1).
 * One conjunction: single + non-visual + Verify present + empty blocking frontier.
 */
const ELIGIBILITY_NEEDLES = [
  'Sizing: single',
  'non-visual',
  'Verify present',
  'blocking frontier',
];

/** Soft-pass stand-in cue (E1 / S1 / T3-c): marker covers these four soft gates. */
const SOFT_PASS_GATES = ['lens', 'testability', 'strategy', 'RCRCRC'];

/** Legacy clear soft-pass markers (E5 / S5 / edge-legacy-skips). */
const LEGACY_SKIP_MARKERS = [
  'Lens pass skipped:',
  'Testability pass skipped:',
  'Strategy pass skipped:',
  'RCRCRC pass skipped:',
];

function readUtf8(filePath) {
  return readFileSync(filePath, 'utf8');
}

function skillPair() {
  return [
    { label: '.cursor/skills/sc-discuss/SKILL.md', text: readUtf8(CURSOR_SKILL) },
    {
      label: 'plugins/spacecraft/skills/sc-discuss/SKILL.md',
      text: readUtf8(PLUGIN_SKILL),
    },
  ];
}

function sizingPair() {
  return [
    {
      label: '.cursor/skills/sc-discuss/references/mission-sizing.md',
      text: readUtf8(CURSOR_SIZING),
    },
    {
      label: 'plugins/spacecraft/skills/sc-discuss/references/mission-sizing.md',
      text: readUtf8(PLUGIN_SIZING),
    },
  ];
}

function assertIncludes(text, needle, label) {
  assert.ok(
    text.includes(needle),
    `${label} must include ${JSON.stringify(needle)}`,
  );
}

/**
 * Must-not / forbid stamp: Must-not|Must not near FAST_MARKER and a reason needle
 * (E2–E4 / S2–S4 / T3-a). Window keeps assert on co-located contract text, not a
 * runtime eligibility helper.
 */
function assertMustNotFastForReason(text, reasonNeedle, label) {
  const mustNotRe = /Must[- ]not/gi;
  let match;
  let found = false;
  while ((match = mustNotRe.exec(text)) !== null) {
    const start = Math.max(0, match.index - 80);
    const end = Math.min(text.length, match.index + 500);
    const window = text.slice(start, end);
    if (
      window.includes(FAST_MARKER) &&
      window.toLowerCase().includes(reasonNeedle.toLowerCase())
    ) {
      found = true;
      break;
    }
  }
  assert.ok(
    found,
    `${label} must Must-not / Must not stamp ${JSON.stringify(FAST_MARKER)} when ${JSON.stringify(reasonNeedle)} (window ±Must-not)`,
  );
}

// --- S1 / E1 / pos-fast-marker + eligibility (T3-a/c) ---

test('S1/E1: both SKILL.md contain Discuss path: fast marker', () => {
  for (const { label, text } of skillPair()) {
    assertIncludes(text, FAST_MARKER, label);
  }
});

test('S1/E1: both SKILL.md document eligibility Sizing: single + non-visual + Verify present + blocking frontier', () => {
  for (const { label, text } of skillPair()) {
    for (const needle of ELIGIBILITY_NEEDLES) {
      assertIncludes(text, needle, label);
    }
  }
});

test('S1/E1: both SKILL.md document Discuss path: fast stands in for lens/testability/strategy/RCRCRC soft-pass', () => {
  for (const { label, text } of skillPair()) {
    assertIncludes(text, FAST_MARKER, label);
    const lower = text.toLowerCase();
    const markerIdx = text.indexOf(FAST_MARKER);
    assert.ok(markerIdx >= 0, `${label} missing ${FAST_MARKER}`);
    // Stand-in contract lives near the marker (E1 expected: soft-pass satisfied via marker).
    const window = text.slice(Math.max(0, markerIdx - 400), markerIdx + 600).toLowerCase();
    for (const gate of SOFT_PASS_GATES) {
      assert.ok(
        window.includes(gate.toLowerCase()) || lower.includes(gate.toLowerCase()),
        `${label} must tie ${FAST_MARKER} to ${gate} soft-pass stand-in`,
      );
    }
    assert.ok(
      /stand[\s-]*in|satisfies|soft-pass|clear/.test(window),
      `${label} must state ${FAST_MARKER} stand-in / soft-pass clear effect near marker`,
    );
  }
});

// --- S2–S4 / E2–E4 Must-not stamp (T3-a) ---

test('S2/E2 neg-roadmap-fast: both SKILL.md Must-not stamp Discuss path: fast for roadmap', () => {
  for (const { label, text } of skillPair()) {
    assertMustNotFastForReason(text, 'roadmap', label);
  }
});

test('S2/E2 neg-roadmap-fast: both SKILL.md Must-not stamp Discuss path: fast for phases', () => {
  for (const { label, text } of skillPair()) {
    assertMustNotFastForReason(text, 'phases', label);
  }
});

test('S3/E3 neg-soft-verify-fast: both SKILL.md Must-not stamp Discuss path: fast when Verify soft or missing', () => {
  for (const { label, text } of skillPair()) {
    // Reason needle from E3 / T1-b: soft or missing Verify
    const mustNotRe = /Must[- ]not/gi;
    let found = false;
    let match;
    while ((match = mustNotRe.exec(text)) !== null) {
      const start = Math.max(0, match.index - 80);
      const end = Math.min(text.length, match.index + 500);
      const window = text.slice(start, end);
      if (
        window.includes(FAST_MARKER) &&
        /Verify/i.test(window) &&
        /(soft|missing)/i.test(window)
      ) {
        found = true;
        break;
      }
    }
    assert.ok(
      found,
      `${label} must Must-not stamp ${FAST_MARKER} when Verify soft/missing`,
    );
  }
});

test('S4/E4 over-visual-fast: both SKILL.md Must-not / forbid Discuss path: fast when visual draft required', () => {
  for (const { label, text } of skillPair()) {
    assertMustNotFastForReason(text, 'visual', label);
  }
});

// --- S6 / E6 / pos-sizing-no-ask (T3-b) ---

test('S6/E6: both mission-sizing.md contain Must not ask / Must-not ask', () => {
  for (const { label, text } of sizingPair()) {
    assert.ok(
      /Must[- ]not ask/i.test(text),
      `${label} must include Must not ask / Must-not ask`,
    );
  }
});

test('S6/E6: both mission-sizing.md forbid one-vs-many / one mission vs many ask', () => {
  for (const { label, text } of sizingPair()) {
    assert.ok(
      /one[- ]?vs[- ]?many|one mission vs many|single.?vs.?multi/i.test(text),
      `${label} must forbid one-vs-many / one mission vs many ask`,
    );
  }
});

test('S6/E6: both mission-sizing.md document auto-split when roadmap Must-when', () => {
  for (const { label, text } of sizingPair()) {
    assert.ok(
      /auto-?split|auto-?appl/i.test(text),
      `${label} must document auto-split / auto-apply when Must-when fires`,
    );
  }
});

test('S6/E6: both mission-sizing.md contain greppable can-split', () => {
  for (const { label, text } of sizingPair()) {
    assert.ok(/can-split/.test(text), `${label} must include greppable can-split`);
  }
});

test('S6/E6: both mission-sizing.md contain too big for one heuristic', () => {
  for (const { label, text } of sizingPair()) {
    assert.ok(
      /too big for one/i.test(text),
      `${label} must include too big for one`,
    );
  }
});

test('S6/E6: both mission-sizing.md require explicit human approval for map-id reuse', () => {
  for (const { label, text } of sizingPair()) {
    assert.ok(
      /human (explicitly )?approved|explicit human approval|explicitly approved reusing/i.test(
        text,
      ),
      `${label} must require explicit human approval for map-id reuse`,
    );
  }
});

// --- S7 / E7 / edge-plugin-mirror (T3-b) ---

test('S7/E7 edge-plugin-mirror: .cursor vs plugins SKILL.md are byte-identical', () => {
  const a = readFileSync(CURSOR_SKILL);
  const b = readFileSync(PLUGIN_SKILL);
  assert.equal(
    Buffer.compare(a, b),
    0,
    '.cursor and plugins sc-discuss/SKILL.md must be byte-identical',
  );
});

test('S7/E7 edge-plugin-mirror: .cursor vs plugins mission-sizing.md are byte-identical', () => {
  const a = readFileSync(CURSOR_SIZING);
  const b = readFileSync(PLUGIN_SIZING);
  assert.equal(
    Buffer.compare(a, b),
    0,
    '.cursor and plugins mission-sizing.md must be byte-identical',
  );
});

// --- S5 / E5 / edge-legacy-skips (T3-c) ---

test('S5/E5 edge-legacy-skips: both SKILL.md Verify/Exit keep legacy Lens/Testability/Strategy/RCRCRC pass or skip', () => {
  for (const { label, text } of skillPair()) {
    for (const needle of LEGACY_SKIP_MARKERS) {
      assertIncludes(text, needle, label);
    }
    assertIncludes(text, '## Lens pass', label);
    assertIncludes(text, '## Testability pass', label);
    assertIncludes(text, '## Strategy pass', label);
    assertIncludes(text, '## RCRCRC pass', label);
  }
});

test('S5/E5 edge-legacy-skips: legacy soft-pass path remains valid without requiring Discuss path: fast alone', () => {
  for (const { label, text } of skillPair()) {
    // Exit clear checklist (E5): legacy or-skip lines stay valid without fast marker.
    // If FAST_MARKER is added to Exit, it must be OR with legacy — not the sole gate.
    const exitIdx = text.indexOf('### Exit');
    assert.ok(exitIdx >= 0, `${label} must have ### Exit clear checklist`);
    const exitBlock = text.slice(exitIdx, exitIdx + 2000);
    assertIncludes(exitBlock, 'Lens pass skipped:', `${label} ### Exit`);
    assertIncludes(exitBlock, 'Testability pass skipped:', `${label} ### Exit`);
    assertIncludes(exitBlock, 'Strategy pass skipped:', `${label} ### Exit`);
    assertIncludes(exitBlock, 'RCRCRC pass skipped:', `${label} ### Exit`);
    if (exitBlock.includes(FAST_MARKER)) {
      assert.ok(
        /\bor\b/i.test(exitBlock),
        `${label} ### Exit: ${FAST_MARKER} must be OR with legacy skips, not sole requirement`,
      );
    }
  }
});
