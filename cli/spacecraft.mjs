#!/usr/bin/env node

import { existsSync } from 'node:fs';
import path from 'node:path';
import { fileURLToPath } from 'node:url';
import { archiveCmd } from './lib/archive.mjs';
import { closeoutCmd } from './lib/closeout.mjs';
import { contextCmd } from './lib/context.mjs';
import { driftCmd } from './lib/drift.mjs';
import { eviCmd } from './lib/evi.mjs';
import { freezeCheckCmd, freezeCmd } from './lib/freeze.mjs';
import { mapCmd } from './lib/map.mjs';
import {
  bindBranchCmd,
  currentCmd,
  flowCmd,
  initCmd,
  missionsCmd,
  newCmd,
  resolveCmd,
  statusCmd,
  useCmd,
} from './lib/mission.mjs';
import { mutationCmd } from './lib/mutation.mjs';
import { ensureProjectReady, ensureSpaceIgnored } from './lib/project-git.mjs';
import { resolveMission, spaceDirFromCwd } from './lib/resolve.mjs';
import { setupCmd } from './lib/setup.mjs';
import { clarifyStatusCmd, stateCmd } from './lib/state.mjs';
import { valCmd } from './lib/val.mjs';

const REPO_ROOT = path.resolve(path.dirname(fileURLToPath(import.meta.url)), '..');

const KEPT_COMMANDS = [
  'init',
  'new',
  'missions',
  'use',
  'current',
  'resolve',
  'status',
  'context',
  'drift',
  'flow',
  'bind-branch',
  'set-state',
  'clarify-status',
  'evidence',
  'validate',
  'freeze',
  'freeze-check',
  'closeout-check',
  'ship-check',
  'archive',
  'map',
  'roadmap',
  'setup',
  'mutation',
];

const ALIASES = {
  state: 'set-state',
  evi: 'evidence',
  val: 'validate',
  roadmap: 'map',
};

const IMPLEMENTED = new Set([
  'init',
  'new',
  'missions',
  'use',
  'current',
  'resolve',
  'status',
  'context',
  'drift',
  'flow',
  'bind-branch',
  'set-state',
  'clarify-status',
  'evidence',
  'validate',
  'freeze',
  'freeze-check',
  'closeout-check',
  'ship-check',
  'map',
  'archive',
  'setup',
  'mutation',
]);

function printHelp() {
  console.log('Usage: spacecraft <command> [args]');
  console.log('');
  console.log('Commands:');
  for (const cmd of KEPT_COMMANDS) {
    console.log(`  spacecraft ${cmd}`);
  }
}

function dispatch(command, args, spaceDir, cwd, mid) {
  switch (command) {
    case 'init':
      return initCmd(spaceDir);
    case 'new':
      return newCmd(args, spaceDir);
    case 'missions':
      return missionsCmd(spaceDir);
    case 'use':
      return useCmd(args, spaceDir);
    case 'current':
      return currentCmd(spaceDir);
    case 'resolve':
      return resolveCmd(args, spaceDir, mid);
    case 'status':
      return statusCmd(spaceDir, mid);
    case 'context':
      return contextCmd(args, spaceDir, cwd, mid);
    case 'drift':
      return driftCmd(args, spaceDir, cwd, mid);
    case 'flow':
      return flowCmd(spaceDir, mid);
    case 'bind-branch':
      return bindBranchCmd(args, spaceDir, cwd, mid);
    case 'set-state':
      return stateCmd(args, spaceDir, mid);
    case 'clarify-status':
      return clarifyStatusCmd(args, spaceDir, mid);
    case 'evidence':
      return eviCmd(args, spaceDir, mid);
    case 'validate':
      return valCmd(args, spaceDir, mid);
    case 'freeze':
      return freezeCmd(args, spaceDir, mid);
    case 'freeze-check':
      return freezeCheckCmd(args, spaceDir, mid);
    case 'closeout-check':
    case 'ship-check':
      return closeoutCmd(spaceDir, mid);
    case 'map':
      return mapCmd(args, spaceDir);
    case 'archive':
      return archiveCmd(args, spaceDir, mid);
    case 'setup':
      return setupCmd(args, cwd, REPO_ROOT);
    case 'mutation':
      return mutationCmd(spaceDir, mid, args);
    default:
      return null;
  }
}

const args = process.argv.slice(2);
const rawCommand = args[0];
const command = ALIASES[rawCommand] ?? rawCommand;

if (!rawCommand || rawCommand === 'help' || rawCommand === '--help' || rawCommand === '-h') {
  printHelp();
  process.exit(0);
}

if (!KEPT_COMMANDS.includes(command) && !KEPT_COMMANDS.includes(rawCommand)) {
  console.error(`spacecraft: unknown command '${rawCommand}'`);
  process.exit(1);
}

if (!IMPLEMENTED.has(command)) {
  console.error(`spacecraft ${command}: not implemented`);
  process.exit(1);
}

const cwd = process.cwd();
// context/drift are read-only: do not mutate project before run
// (otherwise absent-tree skip / omit-missing cases cannot observe absent docs).
if (command !== 'context' && command !== 'drift') {
  try {
    if (!existsSync(path.join(cwd, '.space'))) {
      ensureProjectReady(cwd);
    } else {
      ensureSpaceIgnored(cwd);
    }
  } catch (err) {
    console.error(`spacecraft: ${err.message}`);
    process.exit(1);
  }
}

const spaceDir = spaceDirFromCwd(cwd);
const mid = resolveMission(cwd);
const code = dispatch(command, args.slice(1), spaceDir, cwd, mid);
process.exit(code ?? 1);
