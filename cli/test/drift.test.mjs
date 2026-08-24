/**
 * spacecraft drift — M9C0ZW5E T1–T4 (help, rules, skip, clean seed, no-write).
 *
 * Oracles: design-contract Public seams + Edge matrix; approved-scenarios frozen literals.
 *
 * Pbt skipped: not core logic
 */
import assert from 'node:assert/strict';
import { spawnSync } from 'node:child_process';
import {
  existsSync,
  mkdirSync,
  mkdtempSync,
  readdirSync,
  rmSync,
  writeFileSync,
} from 'node:fs';
import os from 'node:os';
import path from 'node:path';
import test from 'node:test';
import { fileURLToPath } from 'node:url';

const repoRoot = path.resolve(path.dirname(fileURLToPath(import.meta.url)), '../..');
const entryPath = path.join(repoRoot, 'cli', 'spacecraft.mjs');

/** decisions.md Test Ideas / approved-scenarios pos-broken-ref frozen path. */
const BROKEN_DOCS_SPEC_PATH = 'docs/specs/x.md';

/** decisions.md Test Ideas / approved-scenarios edge-orphan frozen path. */
const ORPHAN_EPIC_PATH = 'docs/epics/a.md';

function runDriftHelp() {
  return spawnSync(process.execPath, [entryPath, 'drift', '--help'], {
    encoding: 'utf8',
    cwd: repoRoot,
  });
}

function helpCombined(result) {
  return `${result.stdout ?? ''}${result.stderr ?? ''}`;
}

function spaceRoot() {
  const dir = mkdtempSync(path.join(os.tmpdir(), 'spacecraft-drift-'));
  mkdirSync(path.join(dir, '.space', 'missions'), { recursive: true });
  mkdirSync(path.join(dir, '.space', 'roadmaps'), { recursive: true });
  return dir;
}

function cleanup(dir) {
  rmSync(dir, { recursive: true, force: true });
}

/**
 * Minimal mission fixture (context/validate style).
 * @param {string} root
 * @param {string} id
 * @param {{ clarifyStatus?: string, specBody?: string }} [opts]
 */
function writeDriftMission(root, id, opts = {}) {
  const { clarifyStatus = 'clear', specBody = '# Spec\n' } = opts;
  const dir = path.join(root, '.space', 'missions', id);
  mkdirSync(dir, { recursive: true });
  writeFileSync(
    path.join(dir, 'mission.json'),
    `${JSON.stringify(
      {
        id,
        title: 'Drift Fixture Mission',
        state: 'active',
        createdAt: '2026-01-01T00:00:00Z',
        branches: [],
      },
      null,
      2,
    )}\n`,
  );
  writeFileSync(path.join(dir, 'spec.md'), specBody);
  writeFileSync(
    path.join(dir, 'plan.json'),
    `${JSON.stringify({ planName: 'test', missionId: id, tasks: [] }, null, 2)}\n`,
  );
  writeFileSync(path.join(dir, 'evidence.jsonl'), '');
  writeFileSync(path.join(dir, 'clarify-status'), `${clarifyStatus}\n`);
  writeFileSync(path.join(root, '.space', 'current'), `${id}\n`);
}

/**
 * @param {string} cwd
 * @param {string[]} [extraArgs]
 */
function runDrift(cwd, extraArgs = []) {
  const result = spawnSync(process.execPath, [entryPath, 'drift', ...extraArgs], {
    cwd,
    encoding: 'utf8',
  });
  return {
    stdout: result.stdout ?? '',
    stderr: result.stderr ?? '',
    code: result.status ?? 1,
  };
}

/** Rule 1 fixture: clear + Goal present, Verify heading absent. */
function missingVerifySpec() {
  return '# Spec\n\n## Goal\n\nCatch drift.\n';
}

/** Rule 2 fixture: Goal+Verify present; cites missing docs/specs path. */
function brokenRefSpec() {
  return (
    '# Spec\n\n## Goal\n\nShip drift.\n\n## Verify\n\n' +
    `See ${BROKEN_DOCS_SPEC_PATH} for the product contract.\n`
  );
}

/** Clean Goal+Verify spec — isolates T3 rules from Rule 1 findings. */
function cleanGoalVerifySpec() {
  return '# Spec\n\n## Goal\n\nShip drift.\n\n## Verify\n\nConfirm drift report.\n';
}

/**
 * Rule 3 claim: docs map mentions conventions/ (design-contract Rule 3 claim signal).
 * Must not create docs/conventions/.
 */
function writeConventionsClaimingDocsMap(root) {
  mkdirSync(path.join(root, 'docs'), { recursive: true });
  writeFileSync(
    path.join(root, 'docs', 'README.md'),
    '# Docs map\n\n' +
      'Preferred read order includes conventions/ for shared engineering norms.\n' +
      '| `conventions/` | Naming, code style, review norms |\n',
  );
}

/**
 * Rule 4 fixture: non-empty docs/epics/a.md (decisions Test Ideas edge-orphan).
 */
function writeOrphanEpic(root) {
  mkdirSync(path.join(root, 'docs', 'epics'), { recursive: true });
  writeFileSync(path.join(root, ORPHAN_EPIC_PATH), '# Epic A\n\nUnused by mission.\n');
}

/**
 * T4 thin docs seed: README + conventions only (no epics/specs trees).
 * Isolates Rules 1–4: Goal+Verify on mission; conventions present; no orphans/broken refs.
 */
function writeThinDocsSeed(root) {
  mkdirSync(path.join(root, 'docs', 'conventions'), { recursive: true });
  writeFileSync(
    path.join(root, 'docs', 'README.md'),
    '# Docs map\n\n' +
      'Preferred read order includes conventions/ for shared engineering norms.\n' +
      '| `conventions/` | Naming, code style, review norms |\n',
  );
  writeFileSync(
    path.join(root, 'docs', 'conventions', 'README.md'),
    '# Conventions\n\nShared engineering norms.\n',
  );
}

/**
 * Sorted relative file paths under docs/ (posix-style).
 * @param {string} root
 * @returns {string[]}
 */
function listDocsRelPaths(root) {
  const docsRoot = path.join(root, 'docs');
  const out = [];
  function walk(abs, relBase) {
    if (!existsSync(abs)) return;
    let entries;
    try {
      entries = readdirSync(abs, { withFileTypes: true });
    } catch {
      return;
    }
    for (const ent of entries) {
      const rel = relBase ? `${relBase}/${ent.name}` : ent.name;
      const child = path.join(abs, ent.name);
      if (ent.isDirectory()) {
        walk(child, rel);
      } else if (ent.isFile()) {
        out.push(rel);
      }
    }
  }
  walk(docsRoot, '');
  out.sort((a, b) => (a < b ? -1 : a > b ? 1 : 0));
  return out;
}

/** design-contract Finding: stdout lines that count as findings (not skips). */
function findingLines(stdout) {
  return (stdout ?? '')
    .split('\n')
    .map((l) => l.trimEnd())
    .filter((l) => /^finding:/i.test(l));
}

/**
 * Init fixture git + commit docs/ so over-no-git-write can assert tracked set.
 * @param {string} root
 */
function commitDocsSeed(root) {
  const gitOpts = { cwd: root, encoding: 'utf8' };
  const init = spawnSync('git', ['init'], gitOpts);
  assert.equal(init.status, 0, `git init failed\n${init.stderr}`);
  const add = spawnSync('git', ['add', 'docs'], gitOpts);
  assert.equal(add.status, 0, `git add docs failed\n${add.stderr}`);
  const commit = spawnSync(
    'git',
    [
      '-c',
      'user.email=drift-fixture@example.com',
      '-c',
      'user.name=Drift Fixture',
      'commit',
      '-m',
      'thin docs seed',
    ],
    gitOpts,
  );
  assert.equal(commit.status, 0, `git commit failed\n${commit.stderr}`);
}

/** @param {string} root @returns {string[]} */
function gitTrackedDocs(root) {
  const res = spawnSync('git', ['ls-files', 'docs'], {
    cwd: root,
    encoding: 'utf8',
  });
  assert.equal(res.status, 0, `git ls-files docs failed\n${res.stderr}`);
  return (res.stdout ?? '')
    .split('\n')
    .map((l) => l.trim())
    .filter(Boolean)
    .sort((a, b) => (a < b ? -1 : a > b ? 1 : 0));
}

/**
 * pos-drift-help: spacecraft drift --help exits 0
 * (approved-scenarios / design-contract Edge matrix).
 */
test('pos-drift-help: spacecraft drift --help exits 0', () => {
  const help = runDriftHelp();
  assert.equal(
    help.status,
    0,
    `drift --help must exit 0\nstderr=${help.stderr}\nstdout=${help.stdout}`,
  );
});

/**
 * pos-drift-help: help documents the four first-ship rules
 * (design-contract Data structures Rule 1–4 / decisions Frontier).
 */
test('pos-drift-help: spacecraft drift --help documents four first-ship rules', () => {
  const help = runDriftHelp();
  const helpOut = helpCombined(help);

  assert.equal(
    help.status,
    0,
    `drift --help must exit 0 before rule asserts\nstderr=${help.stderr}\nstdout=${help.stdout}`,
  );

  // Rule 1: clarify-status clear + Goal or Verify heading
  assert.match(
    helpOut,
    /clarify-status/i,
    `drift --help must document rule 1 (clarify-status)\n${helpOut}`,
  );
  assert.match(
    helpOut,
    /Goal/,
    `drift --help must document rule 1 (Goal heading)\n${helpOut}`,
  );
  assert.match(
    helpOut,
    /Verify/,
    `drift --help must document rule 1 (Verify heading)\n${helpOut}`,
  );

  // Rule 2: broken docs/specs/* or docs/epics/* refs
  assert.match(
    helpOut,
    /docs\/specs/,
    `drift --help must document rule 2 (docs/specs)\n${helpOut}`,
  );
  assert.match(
    helpOut,
    /docs\/epics/,
    `drift --help must document rule 2 (docs/epics)\n${helpOut}`,
  );

  // Rule 3: docs/conventions/ missing while claimed
  assert.match(
    helpOut,
    /docs\/conventions/,
    `drift --help must document rule 3 (docs/conventions)\n${helpOut}`,
  );

  // Rule 4: orphan under non-empty epics/specs
  assert.match(
    helpOut,
    /orphan/i,
    `drift --help must document rule 4 (orphan)\n${helpOut}`,
  );
});

/**
 * pos-drift-help: help documents skip-vs-find degrade
 * (approved-scenarios frozen: skip-vs-find; design-contract Help Must mention).
 */
test('pos-drift-help: spacecraft drift --help documents skip-vs-find degrade', () => {
  const help = runDriftHelp();
  const helpOut = helpCombined(help);

  assert.equal(
    help.status,
    0,
    `drift --help must exit 0 before skip-vs-find assert\nstderr=${help.stderr}\nstdout=${help.stdout}`,
  );

  assert.match(
    helpOut,
    /skip/i,
    `drift --help must document skip degrade\n${helpOut}`,
  );
  assert.match(
    helpOut,
    /finding/i,
    `drift --help must document finding vs skip distinction\n${helpOut}`,
  );
});

/**
 * pos-drift-help: help documents default exit 0
 * (approved-scenarios frozen: default exit `0`).
 */
test('pos-drift-help: spacecraft drift --help documents default exit 0', () => {
  const help = runDriftHelp();
  const helpOut = helpCombined(help);

  assert.equal(
    help.status,
    0,
    `drift --help must exit 0 before default-exit assert\nstderr=${help.stderr}\nstdout=${help.stdout}`,
  );

  assert.match(
    helpOut,
    /default[^]*\bexit\b[^]*\b0\b|\bexit\b[^]*\b0\b[^]*default/i,
    `drift --help must document default exit 0\n${helpOut}`,
  );
});

/**
 * pos-drift-help: help documents --strict non-zero-on-findings
 * (approved-scenarios frozen: --strict non-zero-on-findings).
 */
test('pos-drift-help: spacecraft drift --help documents --strict non-zero-on-findings', () => {
  const help = runDriftHelp();
  const helpOut = helpCombined(help);

  assert.equal(
    help.status,
    0,
    `drift --help must exit 0 before --strict assert\nstderr=${help.stderr}\nstdout=${help.stdout}`,
  );

  assert.match(
    helpOut,
    /--strict/,
    `drift --help must document --strict\n${helpOut}`,
  );
  assert.match(
    helpOut,
    /non-?zero/i,
    `drift --help must document non-zero-on-findings\n${helpOut}`,
  );
  assert.match(
    helpOut,
    /finding/i,
    `drift --help must tie --strict non-zero to findings\n${helpOut}`,
  );
});

// --- T2 pos-missing-verify (Rule 1) ---

/**
 * pos-missing-verify: stdout finding when clarify-status clear and Verify missing
 * (design-contract Rule 1 / Edge matrix; approved-scenarios frozen: finding).
 */
test('pos-missing-verify: stdout finding when clear and Verify heading missing', () => {
  const dir = spaceRoot();
  const id = 'M9DRF001';
  try {
    writeDriftMission(dir, id, {
      clarifyStatus: 'clear',
      specBody: missingVerifySpec(),
    });

    const res = runDrift(dir);

    // design-contract Rule 1: Goal/Verify heading language on stdout finding
    assert.match(
      res.stdout,
      /Goal\/Verify|Verify.*heading|heading.*Verify/i,
      `stdout must find missing Verify/Goal heading\nstdout=${res.stdout}\nstderr=${res.stderr}`,
    );
  } finally {
    cleanup(dir);
  }
});

/**
 * pos-missing-verify: --strict non-zero when finding present
 * (approved-scenarios frozen: --strict non-zero).
 * Default exit 0 is stub-true today; assert after GREEN with finding present
 * (one condition: non-zero under --strict — posture vs default).
 */
test('pos-missing-verify: --strict exits non-zero when clear and Verify heading missing', () => {
  const dir = spaceRoot();
  const id = 'M9DRF002';
  try {
    writeDriftMission(dir, id, {
      clarifyStatus: 'clear',
      specBody: missingVerifySpec(),
    });

    const res = runDrift(dir, ['--strict']);
    assert.notEqual(
      res.code,
      0,
      `--strict must exit non-zero on Rule 1 finding\nstderr=${res.stderr}\nstdout=${res.stdout}`,
    );
  } finally {
    cleanup(dir);
  }
});

// --- T2 pos-broken-ref (Rule 2) ---

/**
 * pos-broken-ref: stdout finding when mission cites missing docs/specs path
 * (design-contract Rule 2 / Edge matrix; decisions Test Ideas docs/specs/x.md).
 */
test('pos-broken-ref: stdout finding when citing missing docs/specs path', () => {
  const dir = spaceRoot();
  const id = 'M9DRF004';
  try {
    writeDriftMission(dir, id, {
      clarifyStatus: 'clear',
      specBody: brokenRefSpec(),
    });

    const res = runDrift(dir);

    // decisions Test Ideas frozen path; design-contract Rule 2 broken docs ref
    assert.match(
      res.stdout,
      /docs\/specs\/x\.md/,
      `stdout must find broken docs path ${BROKEN_DOCS_SPEC_PATH}\nstdout=${res.stdout}\nstderr=${res.stderr}`,
    );
  } finally {
    cleanup(dir);
  }
});

/**
 * pos-broken-ref: --strict non-zero when finding present
 * (approved-scenarios frozen: --strict non-zero).
 * Default exit 0 is stub-true today; assert after GREEN with finding present.
 */
test('pos-broken-ref: --strict exits non-zero when citing missing docs/specs path', () => {
  const dir = spaceRoot();
  const id = 'M9DRF003';
  try {
    writeDriftMission(dir, id, {
      clarifyStatus: 'clear',
      specBody: brokenRefSpec(),
    });

    const res = runDrift(dir, ['--strict']);
    assert.notEqual(
      res.code,
      0,
      `--strict must exit non-zero on Rule 2 finding\nstderr=${res.stderr}\nstdout=${res.stdout}`,
    );
  } finally {
    cleanup(dir);
  }
});

// --- T3 edge-orphan (Rule 4) ---

/**
 * edge-orphan: stdout finding when non-empty docs/epics file is unused
 * (design-contract Rule 4 / Edge matrix; decisions Test Ideas docs/epics/a.md).
 */
test('edge-orphan: stdout finding when docs/epics/a.md never referenced', () => {
  const dir = spaceRoot();
  const id = 'M9DRF010';
  try {
    writeDriftMission(dir, id, {
      clarifyStatus: 'clear',
      specBody: cleanGoalVerifySpec(),
    });
    writeOrphanEpic(dir);

    const res = runDrift(dir);

    assert.match(
      res.stdout,
      /docs\/epics\/a\.md/,
      `stdout must find orphan ${ORPHAN_EPIC_PATH}\nstdout=${res.stdout}\nstderr=${res.stderr}`,
    );
  } finally {
    cleanup(dir);
  }
});

/**
 * edge-orphan: --strict non-zero when orphan finding present
 * (approved-scenarios frozen: --strict non-zero).
 * Default exit 0 is stub-true today; assert after GREEN with finding present.
 */
test('edge-orphan: --strict exits non-zero when docs/epics/a.md never referenced', () => {
  const dir = spaceRoot();
  const id = 'M9DRF011';
  try {
    writeDriftMission(dir, id, {
      clarifyStatus: 'clear',
      specBody: cleanGoalVerifySpec(),
    });
    writeOrphanEpic(dir);

    const res = runDrift(dir, ['--strict']);
    assert.notEqual(
      res.code,
      0,
      `--strict must exit non-zero on Rule 4 orphan finding\nstderr=${res.stderr}\nstdout=${res.stdout}`,
    );
  } finally {
    cleanup(dir);
  }
});

// --- T3 neg-skip-absent (skip degrade) ---

/**
 * neg-skip-absent: explicit skip line when optional docs/epics tree absent
 * (design-contract skip line / Edge matrix; approved-scenarios frozen: skip line).
 */
test('neg-skip-absent: stdout skip line when docs/epics tree absent', () => {
  const dir = spaceRoot();
  const id = 'M9DRF020';
  try {
    writeDriftMission(dir, id, {
      clarifyStatus: 'clear',
      specBody: cleanGoalVerifySpec(),
    });
    // No docs/epics or docs/specs trees — optional trees absent.

    const res = runDrift(dir);

    assert.match(
      res.stdout,
      /skip[^\n]*(docs\/)?epics|(docs\/)?epics[^\n]*skip/i,
      `stdout must emit explicit skip for absent docs/epics\nstdout=${res.stdout}\nstderr=${res.stderr}`,
    );
  } finally {
    cleanup(dir);
  }
});

/**
 * neg-skip-absent: --strict still emits skip (skips ≠ findings)
 * (approved-scenarios: skip not finding; --strict exit 0 when no findings).
 * One condition: skip language under --strict for absent docs/epics.
 * Exit 0 with only skips is stub-true today; assert after GREEN when skip lines exist
 * (skips must not flip --strict).
 */
test('neg-skip-absent: --strict still emits skip line when docs/epics tree absent', () => {
  const dir = spaceRoot();
  const id = 'M9DRF021';
  try {
    writeDriftMission(dir, id, {
      clarifyStatus: 'clear',
      specBody: cleanGoalVerifySpec(),
    });

    const res = runDrift(dir, ['--strict']);

    assert.match(
      res.stdout,
      /skip[^\n]*(docs\/)?epics|(docs\/)?epics[^\n]*skip/i,
      `--strict must still emit skip for absent docs/epics (skips ≠ findings)\nstdout=${res.stdout}\nstderr=${res.stderr}`,
    );
  } finally {
    cleanup(dir);
  }
});

// --- T3 rule3-conventions (Rule 3) ---

/**
 * rule3-conventions: stdout finding when docs map claims conventions but tree missing
 * (design-contract Rule 3 claim signal / Edge matrix).
 */
test('rule3-conventions: stdout finding when docs map claims conventions and tree missing', () => {
  const dir = spaceRoot();
  const id = 'M9DRF030';
  try {
    writeDriftMission(dir, id, {
      clarifyStatus: 'clear',
      specBody: cleanGoalVerifySpec(),
    });
    writeConventionsClaimingDocsMap(dir);

    const res = runDrift(dir);

    assert.match(
      res.stdout,
      /docs\/conventions/,
      `stdout must find missing docs/conventions/ while claimed\nstdout=${res.stdout}\nstderr=${res.stderr}`,
    );
  } finally {
    cleanup(dir);
  }
});

/**
 * rule3-conventions: --strict non-zero when conventions finding present
 * (approved-scenarios frozen: --strict non-zero; same exit posture as other findings).
 * Default exit 0 is stub-true today; assert after GREEN with finding present.
 */
test('rule3-conventions: --strict exits non-zero when docs map claims conventions and tree missing', () => {
  const dir = spaceRoot();
  const id = 'M9DRF031';
  try {
    writeDriftMission(dir, id, {
      clarifyStatus: 'clear',
      specBody: cleanGoalVerifySpec(),
    });
    writeConventionsClaimingDocsMap(dir);

    const res = runDrift(dir, ['--strict']);
    assert.notEqual(
      res.code,
      0,
      `--strict must exit non-zero on Rule 3 conventions finding\nstderr=${res.stderr}\nstdout=${res.stdout}`,
    );
  } finally {
    cleanup(dir);
  }
});

// --- T4 over-clean-seed ---

/**
 * over-clean-seed: no finding lines for clean mission + thin docs seed
 * (approved-scenarios frozen: no findings; design-contract over-clean-seed).
 * One condition: stdout has zero finding: lines (skips for absent epics/specs OK).
 */
test('over-clean-seed: no finding lines for clean mission + thin docs seed', () => {
  const dir = spaceRoot();
  const id = 'M9DRF040';
  try {
    writeDriftMission(dir, id, {
      clarifyStatus: 'clear',
      specBody: cleanGoalVerifySpec(),
    });
    writeThinDocsSeed(dir);

    const res = runDrift(dir);
    const findings = findingLines(res.stdout);
    assert.deepEqual(
      findings,
      [],
      `clean seed must emit no finding lines\nfindings=${JSON.stringify(findings)}\nstdout=${res.stdout}\nstderr=${res.stderr}`,
    );
  } finally {
    cleanup(dir);
  }
});

/**
 * over-clean-seed: default exit 0
 * (approved-scenarios frozen: exit 0 default).
 */
test('over-clean-seed: default exit 0 for clean mission + thin docs seed', () => {
  const dir = spaceRoot();
  const id = 'M9DRF041';
  try {
    writeDriftMission(dir, id, {
      clarifyStatus: 'clear',
      specBody: cleanGoalVerifySpec(),
    });
    writeThinDocsSeed(dir);

    const res = runDrift(dir);
    assert.equal(
      res.code,
      0,
      `default drift must exit 0 on clean seed\nstderr=${res.stderr}\nstdout=${res.stdout}`,
    );
  } finally {
    cleanup(dir);
  }
});

/**
 * over-clean-seed: --strict exit 0 when no findings
 * (approved-scenarios frozen: exit 0 --strict).
 */
test('over-clean-seed: --strict exit 0 for clean mission + thin docs seed', () => {
  const dir = spaceRoot();
  const id = 'M9DRF042';
  try {
    writeDriftMission(dir, id, {
      clarifyStatus: 'clear',
      specBody: cleanGoalVerifySpec(),
    });
    writeThinDocsSeed(dir);

    const res = runDrift(dir, ['--strict']);
    assert.equal(
      res.code,
      0,
      `--strict must exit 0 on clean seed (no findings)\nstderr=${res.stderr}\nstdout=${res.stdout}`,
    );
  } finally {
    cleanup(dir);
  }
});

// --- T4 over-no-git-write ---

/**
 * over-no-git-write: docs/ on-disk file set unchanged after drift
 * (design-contract no product-git write; catches buggy writer of docs/drift-report.md etc.).
 * One condition: listDocsRelPaths before === after.
 */
test('over-no-git-write: docs/ file set unchanged after drift', () => {
  const dir = spaceRoot();
  const id = 'M9DRF050';
  try {
    writeDriftMission(dir, id, {
      clarifyStatus: 'clear',
      specBody: cleanGoalVerifySpec(),
    });
    writeThinDocsSeed(dir);

    const before = listDocsRelPaths(dir);
    const res = runDrift(dir);
    assert.equal(
      res.code,
      0,
      `drift must complete before docs/ inventory assert\nstderr=${res.stderr}\nstdout=${res.stdout}`,
    );
    const after = listDocsRelPaths(dir);
    assert.deepEqual(
      after,
      before,
      `drift must not create/delete docs/ files\nbefore=${JSON.stringify(before)}\nafter=${JSON.stringify(after)}\nstdout=${res.stdout}`,
    );
    assert.equal(
      existsSync(path.join(dir, 'docs', 'drift-report.md')),
      false,
      'drift must not write docs/drift-report.md',
    );
  } finally {
    cleanup(dir);
  }
});

/**
 * over-no-git-write: git-tracked docs/ set unchanged after drift
 * (approved-scenarios frozen: no new tracked files under docs/).
 * One condition: git ls-files docs before === after (and porcelain clean under docs/).
 */
test('over-no-git-write: git tracked docs/ unchanged after drift', () => {
  const dir = spaceRoot();
  const id = 'M9DRF051';
  try {
    writeDriftMission(dir, id, {
      clarifyStatus: 'clear',
      specBody: cleanGoalVerifySpec(),
    });
    writeThinDocsSeed(dir);
    commitDocsSeed(dir);

    const beforeTracked = gitTrackedDocs(dir);
    const res = runDrift(dir);
    assert.equal(
      res.code,
      0,
      `drift must complete before tracked-docs assert\nstderr=${res.stderr}\nstdout=${res.stdout}`,
    );
    const afterTracked = gitTrackedDocs(dir);
    assert.deepEqual(
      afterTracked,
      beforeTracked,
      `drift must not add tracked docs/ files\nbefore=${JSON.stringify(beforeTracked)}\nafter=${JSON.stringify(afterTracked)}`,
    );

    const porcelain = spawnSync('git', ['status', '--porcelain', '--', 'docs'], {
      cwd: dir,
      encoding: 'utf8',
    });
    assert.equal(porcelain.status, 0, `git status failed\n${porcelain.stderr}`);
    assert.equal(
      (porcelain.stdout ?? '').trim(),
      '',
      `docs/ working tree must stay clean after drift\n${porcelain.stdout}`,
    );
  } finally {
    cleanup(dir);
  }
});
