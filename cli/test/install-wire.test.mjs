/**
 * Install wire: make build / install-cli / install-binary put the Node CLI
 * entry on PATH with zero package.json CLI dependencies.
 */
import assert from 'node:assert/strict';
import { spawnSync } from 'node:child_process';
import { existsSync, readFileSync } from 'node:fs';
import path from 'node:path';
import test from 'node:test';
import { fileURLToPath } from 'node:url';

const repoRoot = path.resolve(path.dirname(fileURLToPath(import.meta.url)), '../..');
const entryRel = 'cli/spacecraft.mjs';
const entryPath = path.join(repoRoot, entryRel);

/** Tab-indented recipe lines under `target:` (skips blanks; stops at next rule). */
function makefileRecipe(makefile, target) {
  const lines = makefile.split('\n');
  const start = lines.findIndex((line) => new RegExp(`^${target}\\s*:`).test(line));
  assert.ok(start >= 0, `Makefile must define target ${target}`);
  const body = [];
  for (let i = start + 1; i < lines.length; i++) {
    const line = lines[i];
    if (line.startsWith('\t')) {
      body.push(line);
      continue;
    }
    if (line.trim() === '') continue;
    break;
  }
  return body.join('\n');
}

function assertNoCliPackageDeps() {
  for (const rel of ['package.json', 'cli/package.json']) {
    const pkgPath = path.join(repoRoot, rel);
    if (!existsSync(pkgPath)) continue;
    const pkg = JSON.parse(readFileSync(pkgPath, 'utf8'));
    const depKeys = [
      ...Object.keys(pkg.dependencies ?? {}),
      ...Object.keys(pkg.devDependencies ?? {}),
      ...Object.keys(pkg.optionalDependencies ?? {}),
      ...Object.keys(pkg.peerDependencies ?? {}),
    ];
    assert.deepEqual(
      depKeys,
      [],
      `${rel} must declare zero CLI package dependencies (got: ${depKeys.join(', ')})`,
    );
  }
  assert.equal(
    existsSync(path.join(repoRoot, 'node_modules')),
    false,
    'repo must not require node_modules for the CLI',
  );
}

test('make build and install-cli wire the Node CLI entry', () => {
  assert.ok(existsSync(entryPath), `${entryRel} must exist`);

  const makefile = readFileSync(path.join(repoRoot, 'Makefile'), 'utf8');
  const buildBody = makefileRecipe(makefile, 'build');
  const installCliBody = makefileRecipe(makefile, 'install-cli');

  assert.match(
    buildBody,
    /cli\/spacecraft\.mjs/,
    'make build must produce ./spacecraft from cli/spacecraft.mjs',
  );
  assert.doesNotMatch(
    buildBody,
    /\bgo\s+build\b/,
    'make build must wire the Node CLI entry for spacecraft',
  );

  assert.match(
    installCliBody,
    /\$\(LOCAL_BIN\)\/spacecraft/,
    'install-cli must link spacecraft into $(LOCAL_BIN)',
  );
  assert.doesNotMatch(
    installCliBody,
    /\bgo\s+build\b/,
    'install-cli must link the Node-wired binary without a compile step',
  );

  const installBinary = readFileSync(
    path.join(repoRoot, 'scripts', 'install-binary.sh'),
    'utf8',
  );
  assert.match(
    installBinary,
    /cli\/spacecraft\.mjs/,
    'install-binary.sh must install cli/spacecraft.mjs',
  );
  assert.doesNotMatch(
    installBinary,
    /\bgo\s+build\b/,
    'install-binary.sh must install the Node CLI entry',
  );

  assertNoCliPackageDeps();

  const help = spawnSync(process.execPath, [entryPath, 'help'], {
    encoding: 'utf8',
    cwd: repoRoot,
  });
  assert.equal(
    help.status,
    0,
    `node ${entryRel} help must exit 0\nstderr=${help.stderr}\nstdout=${help.stdout}`,
  );
  assert.match(
    `${help.stdout ?? ''}${help.stderr ?? ''}`,
    /spacecraft /i,
    'help smoke must mention spacecraft commands',
  );
});
