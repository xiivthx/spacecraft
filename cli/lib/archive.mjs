import { mkdirSync, renameSync } from 'node:fs';
import path from 'node:path';
import { suggestNextAfterArchive } from './map.mjs';
import { readMission } from './mission.mjs';
import {
  missionDir,
  missionExists,
  normalizeID,
  resolveActive,
} from './resolve.mjs';

function hasHelpFlag(args) {
  return args.some((a) => a === '--help' || a === '-h');
}

export function archiveCmd(args, spaceDir, mid) {
  if (hasHelpFlag(args)) {
    console.log('Usage: spacecraft archive [selector]');
    return 0;
  }
  let id = mid;
  for (const a of args) {
    if (!a.startsWith('--')) {
      id = normalizeID(a);
      break;
    }
  }
  if (!id) {
    id = resolveActive(spaceDir, mid);
  }
  if (!id || !missionExists(spaceDir, id)) {
    console.error('spacecraft archive: no mission to archive');
    return 1;
  }
  let m;
  try {
    m = readMission(spaceDir, id);
  } catch (err) {
    console.error('spacecraft archive:', err.message);
    return 1;
  }
  const state = typeof m.state === 'string' ? m.state : '';
  if (state !== 'shipped') {
    console.error(
      `spacecraft archive: mission ${id} state is ${state}; archive only shipped missions`,
    );
    return 1;
  }
  const archiveDir = path.join(spaceDir, 'archive');
  try {
    mkdirSync(archiveDir, { recursive: true });
    renameSync(missionDir(spaceDir, id), path.join(archiveDir, id));
  } catch (err) {
    console.error('spacecraft archive:', err.message);
    return 1;
  }
  console.log(`Archived mission ${id}`);
  suggestNextAfterArchive(spaceDir, id);
  return 0;
}
