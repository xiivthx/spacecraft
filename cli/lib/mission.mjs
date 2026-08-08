import {
  mkdirSync,
  readdirSync,
  readFileSync,
  writeFileSync,
} from 'node:fs';
import path from 'node:path';
import {
  currentFile,
  gitShowCurrentBranch,
  missionDir,
  missionExists,
  normalizeID,
  readCurrent,
  resolveActive,
} from './resolve.mjs';

const EPOCH_MS = Date.UTC(2026, 0, 1);

export function newMissionID(now = Date.now()) {
  let ms = now - EPOCH_MS;
  if (ms < 0) ms = 0;
  return `M${ms.toString(36).toUpperCase()}`;
}

export function writeCurrent(spaceDir, id) {
  mkdirSync(spaceDir, { recursive: true });
  writeFileSync(currentFile(spaceDir), `${id}\n`);
}

export function listMissionIDs(spaceDir) {
  const dir = path.join(spaceDir, 'missions');
  let entries;
  try {
    entries = readdirSync(dir, { withFileTypes: true });
  } catch {
    return [];
  }
  const ids = entries.filter((e) => e.isDirectory()).map((e) => e.name);
  ids.sort().reverse();
  return ids;
}

export function readMission(spaceDir, id) {
  const data = readFileSync(path.join(missionDir(spaceDir, id), 'mission.json'), 'utf8');
  return JSON.parse(data);
}

function writeJSON(filePath, value) {
  writeFileSync(filePath, `${JSON.stringify(value, null, 2)}\n`);
}

function countEvidence(spaceDir, id) {
  let data;
  try {
    data = readFileSync(path.join(missionDir(spaceDir, id), 'evidence.jsonl'), 'utf8');
  } catch {
    return 0;
  }
  let n = 0;
  for (const line of data.split('\n')) {
    if (line.trim() !== '') n += 1;
  }
  return n;
}

function nextStep(state) {
  switch (state) {
    case 'active':
    case 'planned':
      return '/sc-run';
    case 'in_progress':
      return '/sc-run (continue)';
    case 'ready':
      return '/sc-ship';
    case 'blocked':
      return 'resolve blockers';
    case 'shipped':
      return 'archive';
    default:
      return '/sc-run or spacecraft new';
  }
}

function hasHelpFlag(args) {
  return args.some((a) => a === '--help' || a === '-h');
}

export function initCmd(spaceDir) {
  try {
    mkdirSync(path.join(spaceDir, 'missions'), { recursive: true });
    mkdirSync(path.join(spaceDir, 'roadmaps'), { recursive: true });
  } catch (err) {
    console.error('spacecraft init:', err.message);
    return 1;
  }
  console.log('Spacecraft initialized at .space/');
  return 0;
}

export function newCmd(args, spaceDir) {
  if (hasHelpFlag(args)) {
    console.log('Usage: spacecraft new <title>');
    return 0;
  }
  const title = args.join(' ').trim();
  if (!title) {
    console.error('spacecraft new: missing mission title');
    return 1;
  }

  const id = newMissionID();
  const dir = missionDir(spaceDir, id);
  try {
    mkdirSync(path.join(dir, 'outputs'), { recursive: true });
  } catch (err) {
    console.error('spacecraft new:', err.message);
    return 1;
  }

  const now = new Date().toISOString().replace(/\.\d{3}Z$/, 'Z');
  const mission = {
    id,
    title,
    state: 'active',
    branches: [],
    createdAt: now,
  };
  try {
    writeJSON(path.join(dir, 'mission.json'), mission);
    writeFileSync(path.join(dir, 'spec.md'), `# ${title}\n\n## What\n\n## Why\n`);
    writeJSON(path.join(dir, 'plan.json'), {
      planName: '',
      missionId: id,
      tasks: [],
    });
    writeFileSync(path.join(dir, 'evidence.jsonl'), '');
    writeCurrent(spaceDir, id);
  } catch (err) {
    console.error('spacecraft new:', err.message);
    return 1;
  }

  console.log(`Created mission ${id}`);
  console.log('Next: /sc-run');
  return 0;
}

export function missionsCmd(spaceDir) {
  const ids = listMissionIDs(spaceDir);
  if (ids.length === 0) {
    console.log('No missions.');
    return 0;
  }
  const current = readCurrent(spaceDir);
  ids.forEach((id, i) => {
    let title = '(untitled)';
    let state = 'unknown';
    try {
      const m = readMission(spaceDir, id);
      if (typeof m.title === 'string' && m.title !== '') title = m.title;
      if (typeof m.state === 'string' && m.state !== '') state = m.state;
    } catch {
      // keep defaults
    }
    const marker = id === current ? ' *' : '';
    console.log(`${i + 1}. ${title} (${id}) state:${state}${marker}`);
  });
  return 0;
}

export function useCmd(args, spaceDir) {
  if (args.length === 0) {
    console.error('Usage: spacecraft use <number|id|title>');
    return 1;
  }
  const id = normalizeID(args[0]);
  if (!missionExists(spaceDir, id)) {
    console.error(`spacecraft use: no mission matches ${JSON.stringify(args[0])}`);
    return 1;
  }
  try {
    writeCurrent(spaceDir, id);
  } catch (err) {
    console.error('spacecraft use:', err.message);
    return 1;
  }
  console.log(`Selected mission ${id}`);
  return 0;
}

export function currentCmd(spaceDir) {
  const cur = readCurrent(spaceDir);
  if (!cur) {
    console.log('No current mission. Use spacecraft new then /sc-run.');
    return 0;
  }
  console.log(cur);
  return 0;
}

export function resolveCmd(args, spaceDir, mid) {
  let sel = '';
  for (const a of args) {
    if (!a.startsWith('--')) {
      sel = a;
      break;
    }
  }
  if (sel) {
    const id = normalizeID(sel);
    if (missionExists(spaceDir, id)) {
      console.log(`Mission: ${id}`);
      return 0;
    }
    console.error(`spacecraft resolve: no mission matches ${JSON.stringify(sel)}`);
    return 1;
  }
  const id = resolveActive(spaceDir, mid);
  if (id) {
    console.log(`Mission: ${id}`);
    return 0;
  }
  console.error('spacecraft resolve: no mission resolved');
  return 1;
}

export function statusCmd(spaceDir, mid) {
  const id = resolveActive(spaceDir, mid);
  if (!id) {
    console.log('No selected mission. Use spacecraft new then /sc-run.');
    return 0;
  }
  let m;
  try {
    m = readMission(spaceDir, id);
  } catch (err) {
    console.error('spacecraft status:', err.message);
    return 1;
  }
  console.log(`Mission: ${id}`);
  if (typeof m.title === 'string') console.log(`Title: ${m.title}`);
  if (typeof m.state === 'string') console.log(`State: ${m.state}`);
  console.log(`Evidence: ${countEvidence(spaceDir, id)}`);
  return 0;
}

export function flowCmd(spaceDir, mid) {
  const id = resolveActive(spaceDir, mid);
  if (!id) {
    console.log('No selected mission. Use spacecraft new then /sc-run.');
    return 0;
  }
  let m;
  try {
    m = readMission(spaceDir, id);
  } catch (err) {
    console.error('spacecraft flow:', err.message);
    return 1;
  }
  const state = typeof m.state === 'string' ? m.state : '';
  console.log(`Mission: ${id}`);
  console.log(`State: ${state}`);
  console.log(`Next: ${nextStep(state)}`);
  return 0;
}

export function bindBranchCmd(args, spaceDir, cwd, mid) {
  if (hasHelpFlag(args)) {
    console.log('Usage: spacecraft bind-branch [selector]');
    return 0;
  }
  let id = mid;
  if (args.length > 0) {
    id = normalizeID(args[0]);
  }
  if (!id) {
    id = resolveActive(spaceDir, mid);
  }
  if (!id || !missionExists(spaceDir, id)) {
    console.error('spacecraft bind-branch: no mission to bind');
    return 1;
  }
  const branch = gitShowCurrentBranch(cwd);
  if (!branch) {
    console.error('spacecraft bind-branch: not a git worktree or no current branch');
    return 1;
  }
  let m;
  try {
    m = readMission(spaceDir, id);
  } catch (err) {
    console.error('spacecraft bind-branch:', err.message);
    return 1;
  }
  const branches = Array.isArray(m.branches) ? [...m.branches] : [];
  if (!branches.includes(branch)) {
    branches.push(branch);
  }
  m.branches = branches;
  try {
    writeJSON(path.join(missionDir(spaceDir, id), 'mission.json'), m);
  } catch (err) {
    console.error('spacecraft bind-branch:', err.message);
    return 1;
  }
  console.log(`Bound branch ${branch} to mission ${id}`);
  return 0;
}
