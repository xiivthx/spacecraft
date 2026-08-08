import {
  existsSync,
  mkdirSync,
  readdirSync,
  readFileSync,
  renameSync,
  unlinkSync,
  writeFileSync,
} from 'node:fs';
import path from 'node:path';
import { writeCurrent } from './mission.mjs';
import { currentFile, readCurrent } from './resolve.mjs';

function hasHelpFlag(args) {
  return args.some((a) => a === '--help' || a === '-h');
}

function roadmapPath(dir, id) {
  return path.join(dir, `${id}.json`);
}

function loadRoadmap(dir, id) {
  let data;
  try {
    data = readFileSync(roadmapPath(dir, id), 'utf8');
  } catch {
    throw new Error(`roadmap not found: ${id}`);
  }
  try {
    return JSON.parse(data);
  } catch (err) {
    throw new Error(`invalid roadmap: ${err.message}`);
  }
}

function saveRoadmap(dir, id, r) {
  try {
    writeFileSync(roadmapPath(dir, id), `${JSON.stringify(r, null, 2)}\n`);
  } catch (err) {
    console.error('spacecraft map:', err.message);
    return 1;
  }
  console.log(`Roadmap saved: ${id}`);
  return 0;
}

function slugTitle(title) {
  return title
    .toLowerCase()
    .replaceAll(' ', '-')
    .replace(/[^a-z0-9-]/g, '');
}

function flagValue(args, name) {
  for (let i = 0; i < args.length - 1; i += 1) {
    if (args[i] === name) return args[i + 1];
  }
  return '';
}

function mapNew(args, dir) {
  if (hasHelpFlag(args)) {
    console.log('Usage: spacecraft map new <title> [--desc <text>]');
    return 0;
  }
  if (args.length < 1) {
    console.error('Usage: spacecraft map new <title> [--desc <text>]');
    return 1;
  }
  const title = args[0];
  const desc = flagValue(args.slice(1), '--desc');
  const id = slugTitle(title);
  const now = new Date().toISOString().replace(/\.\d{3}Z$/, 'Z');
  const r = {
    id,
    title,
    description: desc,
    missions: [],
    issues: [],
    createdAt: now,
    updatedAt: now,
  };
  return saveRoadmap(dir, id, r);
}

function mapAdd(args, dir) {
  if (args.length < 2) {
    console.error('Usage: spacecraft map add <roadmap-id> <mission-id> [--desc <text>]');
    return 1;
  }
  const rid = args[0];
  const mid = args[1];
  const desc = flagValue(args.slice(2), '--desc') || mid;
  let r;
  try {
    r = loadRoadmap(dir, rid);
  } catch (err) {
    console.error('spacecraft map:', err.message);
    return 1;
  }
  const missions = Array.isArray(r.missions) ? [...r.missions] : [];
  missions.push({ id: mid, description: desc });
  r.missions = missions;
  r.updatedAt = new Date().toISOString().replace(/\.\d{3}Z$/, 'Z');
  return saveRoadmap(dir, rid, r);
}

function mapRemove(args, dir) {
  if (args.length < 2) {
    console.error('Usage: spacecraft map rm <roadmap-id> <mission-id>');
    return 1;
  }
  const rid = args[0];
  const mid = args[1];
  let r;
  try {
    r = loadRoadmap(dir, rid);
  } catch (err) {
    console.error('spacecraft map:', err.message);
    return 1;
  }
  const missions = Array.isArray(r.missions) ? r.missions : [];
  r.missions = missions.filter((m) => m && m.id !== mid);
  r.updatedAt = new Date().toISOString().replace(/\.\d{3}Z$/, 'Z');
  return saveRoadmap(dir, rid, r);
}

function mapList(dir) {
  let entries;
  try {
    entries = readdirSync(dir);
  } catch {
    return 0;
  }
  for (const name of entries) {
    if (!name.endsWith('.json')) continue;
    const id = name.slice(0, -'.json'.length);
    let r;
    try {
      r = loadRoadmap(dir, id);
    } catch {
      continue;
    }
    const missions = Array.isArray(r.missions) ? r.missions : [];
    const issues = Array.isArray(r.issues) ? r.issues : [];
    const title = typeof r.title === 'string' ? r.title : '';
    console.log(`${id.padEnd(30)} ${title} (${missions.length} missions, ${issues.length} issues)`);
  }
  return 0;
}

function mapShow(args, dir) {
  if (args.length < 1) {
    console.error('Usage: spacecraft map show <roadmap-id>');
    return 1;
  }
  let r;
  try {
    r = loadRoadmap(dir, args[0]);
  } catch (err) {
    console.error('spacecraft map:', err.message);
    return 1;
  }
  console.log(JSON.stringify(r, null, 2));
  return 0;
}

/** Incomplete display state, or "" if ready/shipped (complete). Missing → pending. */
function missionIncompleteState(missionPath) {
  let data;
  try {
    data = readFileSync(path.join(missionPath, 'mission.json'), 'utf8');
  } catch {
    return 'pending';
  }
  let mj;
  try {
    mj = JSON.parse(data);
  } catch {
    return 'pending';
  }
  const state = typeof mj.state === 'string' ? mj.state : '';
  if (state === '') return 'pending';
  if (state === 'ready' || state === 'shipped') return '';
  return state;
}

export function nextIncompleteOnRoadmap(spaceDir, r) {
  const missions = Array.isArray(r.missions) ? r.missions : [];
  for (const entry of missions) {
    if (!entry || typeof entry.id !== 'string' || entry.id === '') continue;
    const mid = entry.id;
    const desc = typeof entry.description === 'string' ? entry.description : '';
    if (existsSync(path.join(spaceDir, 'archive', mid))) continue;
    const st = missionIncompleteState(path.join(spaceDir, 'missions', mid));
    if (st === '') continue;
    return { id: mid, desc, state: st };
  }
  return null;
}

function mapNext(args, dir) {
  if (args.length < 1) {
    console.error('Usage: spacecraft map next <roadmap-id>');
    return 1;
  }
  let r;
  try {
    r = loadRoadmap(dir, args[0]);
  } catch (err) {
    console.error('spacecraft map:', err.message);
    return 1;
  }
  const spaceDir = path.join(dir, '..');
  const next = nextIncompleteOnRoadmap(spaceDir, r);
  if (next) {
    console.log(`${next.id}: ${next.desc} (state=${next.state})`);
    return 0;
  }
  console.log('All missions complete.');
  return 0;
}

function roadmapContainsMission(r, missionID) {
  const missions = Array.isArray(r.missions) ? r.missions : [];
  return missions.some((m) => m && m.id === missionID);
}

function findRoadmapForArchivedMission(spaceDir, archivedID) {
  const roadmapsDir = path.join(spaceDir, 'roadmaps');
  try {
    const cur = readFileSync(path.join(spaceDir, 'current-roadmap'), 'utf8').trim();
    if (cur) {
      try {
        const rm = loadRoadmap(roadmapsDir, cur);
        if (roadmapContainsMission(rm, archivedID)) {
          return { rid: cur, r: rm };
        }
      } catch {
        // fall through to scan
      }
    }
  } catch {
    // no current-roadmap
  }
  let entries;
  try {
    entries = readdirSync(roadmapsDir);
  } catch {
    return null;
  }
  const ids = entries
    .filter((name) => name.endsWith('.json'))
    .map((name) => name.slice(0, -'.json'.length))
    .sort();
  for (const id of ids) {
    try {
      const rm = loadRoadmap(roadmapsDir, id);
      if (roadmapContainsMission(rm, archivedID)) {
        return { rid: id, r: rm };
      }
    } catch {
      // skip
    }
  }
  return null;
}

function clearCurrent(spaceDir) {
  try {
    unlinkSync(currentFile(spaceDir));
  } catch {
    // absent is fine
  }
}

/** Clear stale current, advance to next incomplete roadmap mission, print hint. */
export function suggestNextAfterArchive(spaceDir, archivedID) {
  if (readCurrent(spaceDir) === archivedID) {
    clearCurrent(spaceDir);
  }
  const found = findRoadmapForArchivedMission(spaceDir, archivedID);
  if (!found) return;
  const next = nextIncompleteOnRoadmap(spaceDir, found.r);
  if (!next) return;
  try {
    writeCurrent(spaceDir, next.id);
  } catch {
    // still print hint
  }
  console.log(
    `Next mission on roadmap ${found.rid}: ${next.id}: ${next.desc} (state=${next.state})`,
  );
  console.log(`Suggested: new session → /sc-discuss ${next.id} (then /sc-run)`);
}

function mapUse(args, spaceDir, roadmapsDir) {
  if (args.length < 1) {
    console.error('Usage: spacecraft map use <roadmap-id>');
    return 1;
  }
  const id = args[0];
  try {
    loadRoadmap(roadmapsDir, id);
  } catch (err) {
    console.error('spacecraft map:', err.message);
    return 1;
  }
  try {
    mkdirSync(spaceDir, { recursive: true });
    writeFileSync(path.join(spaceDir, 'current-roadmap'), `${id}\n`);
  } catch (err) {
    console.error('spacecraft map:', err.message);
    return 1;
  }
  console.log(`Selected roadmap ${id}`);
  return 0;
}

function mapCurrent(spaceDir) {
  let data;
  try {
    data = readFileSync(path.join(spaceDir, 'current-roadmap'), 'utf8');
  } catch {
    console.error('spacecraft map: no current roadmap');
    return 1;
  }
  const id = data.trim();
  if (!id) {
    console.error('spacecraft map: no current roadmap');
    return 1;
  }
  console.log(id);
  return 0;
}

function mapArchive(args, dir, spaceDir) {
  if (args.length < 1) {
    console.error('Usage: spacecraft map archive <roadmap-id>');
    return 1;
  }
  const id = args[0];
  const archiveDir = path.join(spaceDir, 'archive', 'roadmaps');
  try {
    mkdirSync(archiveDir, { recursive: true });
    renameSync(path.join(dir, `${id}.json`), path.join(archiveDir, `${id}.json`));
  } catch (err) {
    console.error('spacecraft map:', err.message);
    return 1;
  }
  console.log(`Archived roadmap: ${id}`);
  return 0;
}

export function mapCmd(args, spaceDir) {
  if (args.length < 1) {
    console.error(
      'Usage: spacecraft map <new|add|rm|ls|show|next|use|current|archive> [...]',
    );
    return 1;
  }
  const roadmapsDir = path.join(spaceDir, 'roadmaps');
  mkdirSync(roadmapsDir, { recursive: true });
  const sub = args[0];
  const rest = args.slice(1);
  switch (sub) {
    case 'new':
      return mapNew(rest, roadmapsDir);
    case 'add':
      return mapAdd(rest, roadmapsDir);
    case 'rm':
    case 'remove':
      return mapRemove(rest, roadmapsDir);
    case 'ls':
    case 'list':
      return mapList(roadmapsDir);
    case 'show':
      return mapShow(rest, roadmapsDir);
    case 'next':
      return mapNext(rest, roadmapsDir);
    case 'use':
      return mapUse(rest, spaceDir, roadmapsDir);
    case 'current':
      return mapCurrent(spaceDir);
    case 'archive':
      return mapArchive(rest, roadmapsDir, spaceDir);
    default:
      console.error(`spacecraft map: unknown subcommand ${JSON.stringify(sub)}`);
      return 1;
  }
}
